package cdek_pay

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
)

type Config struct {
	httpClient *http.Client

	login     string
	secretKey string
	baseURL   string
}

func WithLogin(login string) func(*Config) {
	return func(config *Config) {
		config.login = login
	}
}

func WithSecretKey(secretKey string) func(*Config) {
	return func(config *Config) {
		config.secretKey = secretKey
	}
}

func WithBaseURL(baseURL string) func(*Config) {
	return func(config *Config) {
		config.baseURL = baseURL
	}
}

func WithHTTPClient(c *http.Client) func(*Config) {
	return func(config *Config) {
		config.httpClient = c
	}
}

// Client is the main entity which executes requests against the Tinkoff Acquiring API endpoint
type Client struct {
	Config
}

// NewClient returns new Client instance
func NewClient(login string, secretKey string) *Client {
	return NewClientWithOptions(
		WithLogin(login),
		WithSecretKey(secretKey),
	)
}

func NewClientWithOptions(cfgOption ...func(*Config)) *Client {
	defaultConfig := Config{
		httpClient: http.DefaultClient,
		baseURL:    "https://secure.cdekfin.ru/merchant_api/",
	}
	cfg := defaultConfig

	for _, opt := range cfgOption {
		opt(&cfg)
	}

	return &Client{
		Config: cfg,
	}
}

// SetBaseURL allows to change default API endpoint
func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

func (c *Client) decodeResponse(response *http.Response, result interface{}) error {
	if response.StatusCode != 200 {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		bodyString := string(bodyBytes)
		return errors.New(fmt.Sprintf("%v - '%s'", response.StatusCode, bodyString))
	}
	return json.NewDecoder(response.Body).Decode(result)
}

// Deprecated: use PostRequestWithContext instead
func (c *Client) PostRequest(url string, request RequestInterface) (*http.Response, error) {
	return c.PostRequestWithContext(context.Background(), url, request)
}

func (c *Client) secureRequest(request RequestInterface) {
	request.SetLogin(c.login)

	v := request.GetValuesForSignature()
	request.SetSignature(strings.ToUpper(generateSignature(v, c.secretKey)))
}

// PostRequestWithContext will automatically sign the request with token
// Use BaseRequest type to implement any API request
func (c *Client) PostRequestWithContext(ctx context.Context, url string, request RequestInterface) (*http.Response, error) {
	c.secureRequest(request)
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.httpClient.Do(req)
}

// Функция для генерации подписи
func generateSignature(data map[string]interface{}, secretKey string) string {
	// Сортировка ключей
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Конкатенация значений в строку
	signatureStr := ""
	for _, k := range keys {
		signatureStr += fmt.Sprintf("%v|", data[k])
	}
	signatureStr += secretKey

	// Генерация SHA-256
	hash := sha256.Sum256([]byte(signatureStr))
	return hex.EncodeToString(hash[:])
}
