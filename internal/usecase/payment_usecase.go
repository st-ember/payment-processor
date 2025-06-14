package usecase

import (
	"context"
	"paymentprocessor/internal/domain"
	"paymentprocessor/internal/domain/enum"
	"paymentprocessor/internal/domain/request"
	kafkaadapter "paymentprocessor/internal/infra/kafka"
	redisadapter "paymentprocessor/internal/infra/redis"
	stripeadapter "paymentprocessor/internal/lib/stripe"
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
	checkoutMsg := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"status":    enum.Open.String(),
	}

	err = u.kafkaClient.SendMessage(
		kafkaadapter.Topic.CheckoutStatus,
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
	logMsg := map[string]interface{}{
		"level":     "info",
		"timestamp": time.Now().Format(time.RFC3339),
		"order_id":  req.OrderId.Hex(),
		"price_id":  req.StripePriceId,
		"user_id":   req.UserId.Hex(),
		"message":   "Stripe session started",
	}

	err = u.kafkaClient.SendMessage(
		kafkaadapter.Topic.LogPaymentCheckout,
		logMsg,
	)
	if err != nil {
		return "", err
	}

	return checkoutUrl, nil
}

func (u *PaymentUsecase) ConfirmPayment(ctx context.Context, session stripe.CheckoutSession) error {
	sessionId := session.ID
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
	msg := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"status":    string(session.Status),
		"order_id":  orderId,
	}

	err = u.kafkaClient.SendMessage(
		kafkaadapter.Topic.CheckoutStatus,
		msg,
	)
	if err != nil {
		return err
	}

	// update status in db
	err = u.sessionRepo.UpdateStatus(ctx, sessionId, session.Status)
	if err != nil {
		return err
	}

	return nil
}
