package main

import (
	"net/http"
	"paymentprocessor/internal/handler"
)

func main() {
	http.HandleFunc("/payment/start", handler.PaymentStart)
	http.HandleFunc("/payment/confirm", handler.PaymentConfirmation)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}
}
