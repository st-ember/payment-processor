package usecase

import (
	"context"
	"encoding/json"
	"paymentprocessor/internal/domain"
	"paymentprocessor/internal/domain/enum"
	"paymentprocessor/internal/domain/request"
	kafkaadapter "paymentprocessor/internal/infra/kafka_adapter"
	redisadapter "paymentprocessor/internal/infra/redis_adapter"
	stripeadapter "paymentprocessor/internal/infra/stripe_adapter"
	"time"

	"github.com/stripe/stripe-go/v72"
)

type PaymentUsecase struct {
	sessionRepo domain.SessionRepo
	redisUtil   *redisadapter.RedisUtil
	kafkaClient *kafkaadapter.KafkaClient
	stripeUtil  *stripeadapter.CheckoutSessionUtil
}

func NewPaymentUsecase(
	sessionRepo domain.SessionRepo,
	redisUtil *redisadapter.RedisUtil,
	kafkaClient *kafkaadapter.KafkaClient,
	stripeUtil *stripeadapter.CheckoutSessionUtil) *PaymentUsecase {
	return &PaymentUsecase{
		sessionRepo: sessionRepo,
		redisUtil:   redisUtil,
		kafkaClient: kafkaClient,
		stripeUtil:  stripeUtil,
	}
}

func (u *PaymentUsecase) ProcessPayment(ctx context.Context, req request.StartPaymentReq) (string, error) {
	// construce stripe param for checkout request
	stripeParams := []*stripe.CheckoutSessionLineItemParams{
		{
			Price:    stripe.String(req.StripePriceId),
			Quantity: stripe.Int64(req.Quantity),
		},
	}

	// tell stripe about the expected payment
	checkoutId, checkoutUrl, err := u.stripeUtil.StartSession(stripeParams)
	if err != nil {
		return "", err
	}

	// tell order ms about the payment status through Kafka
	checkoutVal := map[string]string{
		"timestamp": time.Now().Format(time.RFC3339),
		"status":    string(rune(enum.Open)),
	}

	checkoutValBytes, err := json.Marshal(checkoutVal)
	if err != nil {
		return "", err
	}

	err = u.kafkaClient.SendMessage(
		kafkaadapter.Topic.PaymentCheckoutSessionStatus,
		[]byte(req.OrderId.Hex()),
		checkoutValBytes,
	)
	if err != nil {
		return "", err
	}

	// tell redis about the checkout id we got from stripe
	// checkoutId as key, orderId as value
	// Todo: change to publish for topic-like behavior
	err = u.redisUtil.SetStripeSession(ctx, checkoutId, req.OrderId.Hex())
	if err != nil {
		return "", err
	}

	// Todo: also log to db in case Redis crashes

	// log stripe session event to Kafka
	logVal := map[string]string{
		"level":     "info",
		"timestamp": time.Now().Format(time.RFC3339),
		"order_id":  req.OrderId.Hex(),
		"price_id":  req.StripePriceId,
		"user_id":   req.UserId.Hex(),
		"message":   "Stripe session started",
	}

	logValBytes, err := json.Marshal(logVal)
	if err != nil {
		return "", err
	}

	err = u.kafkaClient.SendMessage(
		kafkaadapter.Topic.LogPaymentCheckout,
		[]byte(req.OrderId.Hex()),
		logValBytes,
	)
	if err != nil {
		return "", err
	}

	return checkoutUrl, nil
}

func (u *PaymentUsecase) ConfirmPayment(ctx context.Context, sessionId string, status enum.StripeStatus) error {
	// get orderId from redis
	orderId, err := u.redisUtil.GetOrderIdFromSession(ctx, sessionId)
	if err != nil {
		return err
	}

	// confirm orderId exists in db

	// set order's status to the right one

	return nil
}

type Processor interface {
	ProcessPayment(
		ctx context.Context,
		req request.StartPaymentReq) (string, error)
}

type Confirmer interface {
	ConfirmPayment(ctx context.Context, sessionId string, status enum.StripeStatus) error
}
