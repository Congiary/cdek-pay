package cdek_pay

type BaseRequest struct {
	Login     string `json:"login" url:"login"`
	Signature string `json:"signature" url:"signature"`
}

type RequestInterface interface {
	GetValuesForSignature() map[string]interface{}
	SetLogin(login string)
	SetSignature(signature string)
}

func (r *BaseRequest) SetLogin(key string) {
	r.Login = key
}

func (r *BaseRequest) SetSignature(token string) {
	r.Signature = token
}
