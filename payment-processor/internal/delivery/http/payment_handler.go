package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	jwthelper "paymentprocessor/internal/lib/jwt"
	"paymentprocessor/internal/usecase"
	"strings"

	"github.com/stripe/stripe-go/v72"
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
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	event := stripe.Event{}

	if err := json.Unmarshal(payload, &event); err != nil {
		// kafka: request isn't stripe event
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var session stripe.CheckoutSession
	err = json.Unmarshal(event.Data.Raw, &session)
	if err != nil {
		// kafka: event isn't checkoutSession
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}

	err = h.usecase.ConfirmPayment(r.Context(), session)
	if err != nil {
		// kafka: failed to confirm payment
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	// 	Retry webhook delivery: Stripe retries if your webhook fails (use idempotency)
	// Dead-letter queue for failed webhooks
	// Add handler to manually re-sync payment status from Stripe (e.g., for audits)
	// Kafka event deduping (Stripe webhooks can fire twice)
}
