package kafkaadapter

var Topic = struct {
	LogPaymentCheckout  string
	CheckoutStatus      string
	CheckoutStatusBatch string
}{
	LogPaymentCheckout:  "log.payment.checkout",
	CheckoutStatus:      "payment.checkout_session.status",
	CheckoutStatusBatch: "payment.checkout_session.status.batch",
}
