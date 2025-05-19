package handler

import (
	"encoding/json"
	"net/http"
	"paymentprocessor/internal/service"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RequestBody struct {
	OrderId string `json:"order_id"`
	UserId  string `json:"user_id"`
}

type PaymentHandler struct {
	paymentProcessor service.PaymentProcessor
}

func NewPaymentHandler(paymentProcessor service.PaymentProcessor) *PaymentHandler {
	return &PaymentHandler{paymentProcessor: paymentProcessor}
}

func (h *PaymentHandler) PaymentStart(w http.ResponseWriter, r *http.Request) {
	// get order id
	var body RequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	orderId := body.OrderId
	orderObjectID, err := primitive.ObjectIDFromHex(orderId)
	if err != nil {
		http.Error(w, "Cannot parse order_id", http.StatusInternalServerError)
		return
	}

	// return the Stripe checkoutSession.url to frontend
	checkOutUrl, err := h.paymentProcessor.ProcessPayment(r.Context(), orderObjectID, body.UserId)
	if err != nil {
		http.Error(w, "Failed to process payment", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"checkout_url": checkOutUrl,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal session url", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write(jsonData)
}

func (h *PaymentHandler) PaymentConfirmation(w http.ResponseWriter, r *http.Request) {
	// listen to stripe's confirmation

	// tell kafka about the confirmation we got
	// set order's status to the right one

	// 	Retry webhook delivery: Stripe retries if your webhook fails (use idempotency)
	// Dead-letter queue for failed webhooks
	// Add handler to manually re-sync payment status from Stripe (e.g., for audits)
	// Kafka event deduping (Stripe webhooks can fire twice)
}
