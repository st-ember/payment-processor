package handler

import "net/http"

func PaymentStart(w http.ResponseWriter, r *http.Request) {
	// get order id
	// set order to payment pending
	// tell stripe about the expected payment
	// tell kafka about the expected payment and the checkout id we got from stripe
	// return the Stripe checkoutSession.url to frontend
}

func PaymentConfirmation(w http.ResponseWriter, r *http.Request) {
	// listen to stripe's confirmation
	// set order's status to the right one
	// tell kafka about the confirmation we got

	// 	Retry webhook delivery: Stripe retries if your webhook fails (use idempotency)
	// Dead-letter queue for failed webhooks
	// Add handler to manually re-sync payment status from Stripe (e.g., for audits)
	// Kafka event deduping (Stripe webhooks can fire twice)
}
