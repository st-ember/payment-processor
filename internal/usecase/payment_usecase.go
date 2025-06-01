package usecase

import (
	"context"
	"paymentprocessor/internal/domain"
	"paymentprocessor/internal/domain/enum"
	"paymentprocessor/internal/domain/request"
	kafkaadapter "paymentprocessor/internal/infra/kafka_adapter"
	redisadapter "paymentprocessor/internal/infra/redis_adapter"
	stripeadapter "paymentprocessor/internal/infra/stripe_adapter"
	"time"

	"github.com/redis/go-redis/v9"
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
	// construct stripe param for checkout request
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
	checkoutMsg := map[string]string{
		"timestamp": time.Now().Format(time.RFC3339),
		"status":    enum.Open.String(),
	}

	err = u.kafkaClient.SendMessage(
		kafkaadapter.Topic.PaymentCheckoutSessionStatus,
		[]byte(req.OrderId.Hex()),
		checkoutMsg,
	)
	if err != nil {
		return "", err
	}

	// tell redis about the checkout id we got from stripe
	err = u.redisUtil.SetStripeSession(ctx, checkoutId, req.OrderId.Hex())
	if err != nil {
		return "", err
	}

	// log to db in case Redis crashes
	err = u.sessionRepo.Insert(ctx, req.OrderId, checkoutId)
	if err != nil {
		return "", err
	}

	// log stripe session event to Kafka
	logMsg := map[string]string{
		"level":     "info",
		"timestamp": time.Now().Format(time.RFC3339),
		"order_id":  req.OrderId.Hex(),
		"price_id":  req.StripePriceId,
		"user_id":   req.UserId.Hex(),
		"message":   "Stripe session started",
	}

	err = u.kafkaClient.SendMessage(
		kafkaadapter.Topic.LogPaymentCheckout,
		[]byte(req.OrderId.Hex()),
		logMsg,
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
		if err == redis.Nil {
			// if sessionId does not exist in redis, try to find in db
			u.sessionRepo.GetBySessionId(ctx, sessionId)
		} else {
			return err
		}
	}

	// tell order ms about payment conclusion status
	msg := map[string]string{
		"timestamp": time.Now().Format(time.RFC3339),
		"status":    status.String(),
	}

	err = u.kafkaClient.SendMessage(
		kafkaadapter.Topic.PaymentCheckoutSessionStatus,
		[]byte(orderId),
		msg,
	)
	if err != nil {
		return err
	}

	// add concluded session to db

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
