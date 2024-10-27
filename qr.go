package cdek_pay

import (
	"context"
)

type GetQRResponse struct {
	QRLink  string `json:"qr_link"`  // Ссылка на платежную форму с QR кодом
	QRImage string `json:"qr_image"` // Картинка QR кода в base64 формате
	OrderID int    `json:"order_id"` // Сокращенный идентификатор заказа в системе CDEKFIN
}

func (c *Client) InitQRPayment(ctx context.Context, request *InitRequest) (*GetQRResponse, error) {
	response, err := c.PostRequestWithContext(ctx, "sbp_qrs", request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var res GetQRResponse
	err = c.decodeResponse(response, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
