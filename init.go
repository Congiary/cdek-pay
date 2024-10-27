package cdek_pay

import (
	"context"
)

type InitRequest struct {
	BaseRequest

	PaymentOrder PaymentOrder `json:"payment_order"`
}

type PaymentOrder struct {
	PayFor           string             `json:"pay_for"`
	Currency         string             `json:"currency"`
	PayAmount        int                `json:"pay_amount"`
	ReceiptDetails   []ReceiptItem      `json:"receipt_details"`
	LinkLifeTime     *int               `json:"link_life_time"`     // время жизни платежной ссылки в минутах
	UserPhone        *string            `json:"user_phone"`         // номер телефона плательщика
	UserEmail        *string            `json:"user_email"`         // email адрес плательщика
	ReturnURLSuccess *string            `json:"return_url_success"` // URL для возврата после успешного платежа
	ReturnURLFail    *string            `json:"return_url_fail"`    // URL для возврата в случае ошибки
	PayForDetails    *map[string]string `json:"pay_for_details"`    // детальная информация о платеже
	QRLifeTime       *int               `json:"qr_life_time"`       // время жизни QR в минутах
}

type ReceiptItem struct {
	ID            string  `json:"id"`             // Идентификатор товара (уникальный в рамках одного чека)
	Name          string  `json:"name"`           // Наименование товара
	Price         int     `json:"price"`          // Цена за единицу товара в копейках
	Quantity      float64 `json:"quantity"`       // Количество единиц товара, до 2 знаков после запятой
	Sum           int     `json:"sum"`            // Сумма в копейках
	PaymentObject int     `json:"payment_object"` // Признак предмета расчета
}

func (i *InitRequest) GetValuesForSignature() map[string]interface{} {
	return FlattenStructToMap(i.PaymentOrder, "")
}

type InitResponse struct {
	OrderId   int    `json:"order_id"`
	AccessKey string `json:"access_key"`
	Link      string `json:"link"`
}

func (c *Client) InitPayment(ctx context.Context, request *InitRequest) (*InitResponse, error) {
	response, err := c.PostRequestWithContext(ctx, "payment_orders", request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var res InitResponse
	err = c.decodeResponse(response, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
