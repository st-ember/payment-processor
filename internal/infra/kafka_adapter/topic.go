package kafkaadapter

var Topic = struct {
	LogPaymentCheckout           string
	PaymentCheckoutSessionStatus string
}{
	LogPaymentCheckout:           "log.payment.checkout",
	PaymentCheckoutSessionStatus: "payment.checkout_session.status",
}
