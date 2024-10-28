package cdek_pay

import (
	"context"
)

const (
	PaymentColumnId   = "id"
	PaymentColumnTime = "payment_time"
)

const (
	PaymentDirectionASC  = "ASC"
	PaymentDirectionDESC = "DESC"
)

const (
	PaymentStatusSuccess               = "success"                // успешный платеж
	PaymentStatusCancelled             = "cancelled"              // возврат, данный платеж создается после успешного возврата
	PaymentStatusSuccessCancellation   = "success_cancellation"   // успешная отмена
	PaymentStatusCancellationRequested = "cancellation_requested" // запрошена отмена
)

type getPaymentsRequest struct {
	BaseRequest
	GetPaymentsRequest
}

// todo: сделать как-то более красиво структуры
type GetPaymentsRequest struct {
	Page      int    `json:"p[page]" url:"p[page]"`
	PerPage   int    `json:"p[per_page]" url:"p[per_page]"`
	Column    string `json:"o[column]" url:"o[column]"`
	Direction string `json:"o[direction]" url:"o[direction]"`

	OrderId   *int `json:"q[order_id],omitempty" url:"q[order_id],omitempty"`
	AccessKey *int `json:"q[access_key],omitempty" url:"q[access_key],omitempty"`
}

func (i *getPaymentsRequest) GetValuesForSignature() map[string]interface{} {
	return FlattenStructToMap(GetPaymentsRequest{
		Page:      i.Page,
		PerPage:   i.PerPage,
		Column:    i.Column,
		Direction: i.Direction,
		OrderId:   i.OrderId,
		AccessKey: i.AccessKey,
	}, "")
}

type GetPaymentsResponse struct {
	TotalPayments int       `json:"total_payments"` // Общее количество платежей
	CurrentPage   int       `json:"current_page"`   // Текущая страница списка платежей
	TotalPages    int       `json:"total_pages"`    // Общее количество страниц списка платежей
	Payments      []Payment `json:"payments"`       // Массив объектов платежей
}

type Payment struct {
	ID          int    `json:"id"`           // Идентификатор платежа
	OrderID     int    `json:"order_id"`     // Сокращенный идентификатор заказа в системе CDEKFIN
	AccessKey   string `json:"access_key"`   // Уникальный идентификатор заказа в системе CDEKFIN
	Currency    string `json:"currency"`     // Валюта платежа
	PayAmount   int    `json:"pay_amount"`   // Сумма платежа в копейках (отрицательная для возвратов)
	Status      string `json:"status"`       // Статус платежа
	PaymentTime int    `json:"payment_time"` // Дата и время платежа в формате Timestamp
}

// GetPayments возвращает список платежей и возвратов
func (c *Client) GetPayments(ctx context.Context, request *GetPaymentsRequest) (*GetPaymentsResponse, error) {
	req := getPaymentsRequest{
		GetPaymentsRequest: *request,
	}
	response, err := c.GetRequestWithContext(ctx, "payments", &req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var res GetPaymentsResponse
	err = c.decodeResponse(response, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
