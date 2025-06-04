package http

import (
	"encoding/json"
	"net/http"
	jwthelper "paymentprocessor/internal/lib/jwt"
	"paymentprocessor/internal/usecase"
	"strings"
)

type PaymentHandler struct {
	usecase   *usecase.PaymentUsecase
	jwtHelper *jwthelper.JWTHelper
}

func NewPaymentHandler(usecase *usecase.PaymentUsecase, jwtHelper *jwthelper.JWTHelper) *PaymentHandler {
	return &PaymentHandler{usecase: usecase, jwtHelper: jwtHelper}
}

func (h *PaymentHandler) PaymentStart(w http.ResponseWriter, r *http.Request) {
	// check jwt
	var token string
	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		token = parts[1]
	} else {
		// Handle missing or malformed token
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	_, err := h.jwtHelper.ParseToken(token)
	if err != nil {
		http.Error(w, "Invalid claims", http.StatusUnauthorized)
	}

	// get order id
	var req StartPaymentReq
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// return the Stripe checkoutSession.url to frontend
	useCaseReq := req.ToUsecaseRequest()
	checkOutUrl, err := h.usecase.ProcessPayment(r.Context(), useCaseReq)
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

	// 	Retry webhook delivery: Stripe retries if your webhook fails (use idempotency)
	// Dead-letter queue for failed webhooks
	// Add handler to manually re-sync payment status from Stripe (e.g., for audits)
	// Kafka event deduping (Stripe webhooks can fire twice)
}
