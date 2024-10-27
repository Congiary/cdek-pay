# Golang CDEK Pay API client

## Installation
```shell
go get github.com/congiary/cdek-pay@latest
```

## Usage

#### Create client
Provide terminal key and password from terminal settings page.
```go
client := cdek_pay.NewClient(login, secretKey)
```

#### Create payment
```go
var receiptItems []cdek_pay.ReceiptItem
receiptItems = append(receiptItems, cdek_pay.ReceiptItem{
    ID:            "10",
    Name:          "Soul",
    Price:         100,
    Quantity:      1,
    Sum:           100,
    PaymentObject: 1,
})

order := cdek_pay.PaymentOrder{
    PayFor:         "Something",
    Currency:       "TST",
    PayAmount:      100,
    ReceiptDetails: receiptItems,
}
req := &cdek_pay.InitRequest{
    PaymentOrder: order,
}

response, err := client.InitPayment(context.TODO(), req)
```

## References
The code in this repo based on some code from [nikita-vanyasin/tinkoff](https://github.com/nikita-vanyasin/tinkoff)