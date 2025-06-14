package kafkaadapter

var Topic = struct {
	PaymentCheckout     string
	CheckoutStatus      string
	CheckoutStatusBatch string
	PaymentError        string
}{
	PaymentCheckout:     "log.payment.checkout",
	CheckoutStatus:      "payment.checkout_session.status",
	CheckoutStatusBatch: "payment.checkout_session.status.batch",
	PaymentError:        "log.payment.error",
}
