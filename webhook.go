package cdek_pay

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

type Webhook struct {
	Payment   PaymentEvent `json:"payment"`
	Signature string       `json:"signature"`
}

type PaymentEvent struct {
	Amount    int    `json:"pay_amount"`
	AccessKey string `json:"access_key"`
	Currency  string `json:"currency"`
	Id        int    `json:"id"`
	OrderId   int    `json:"order_id"`
}

func (n *Webhook) GetValuesForSignature() map[string]interface{} {
	return FlattenStructToMap(n.Payment, "")
}

func (c *Client) ParseWebhook(requestBody io.Reader) (*Webhook, error) {
	bytes, err := ioutil.ReadAll(requestBody)
	if err != nil {
		return nil, err
	}
	var webhook Webhook
	err = json.Unmarshal(bytes, &webhook)
	if err != nil {
		return nil, err
	}

	valuesForSignature := webhook.GetValuesForSignature()
	signature := generateSignature(valuesForSignature, c.secretKey)
	if signature != webhook.Signature {
		valsForTokenJSON, _ := json.Marshal(valuesForSignature)
		return nil, fmt.Errorf("invalid signature: expected %s got %s.\nValues for signature: %s.\nWebhook: %s", signature, webhook.Signature, valsForTokenJSON, string(bytes))
	}

	return &webhook, nil
}

func (c *Client) GetWebhookSuccessResponse() string {
	return "OK"
}
