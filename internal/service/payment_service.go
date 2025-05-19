package service

import (
	"context"
	"encoding/json"
	"paymentprocessor/internal/enums"
	"paymentprocessor/internal/mongo"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/redis/go-redis"
	"github.com/stripe/stripe-go/v72"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentService struct {
	orderStatusUpdater mongo.OrderStatusUpdater
	orderGetter        mongo.OrderGetter
	productGetter      mongo.ProductGetter
	producer           *kafka.Producer
}

func NewPaymentSerivce(
	orderStatusUpdater storage.OrderStatusUpdater,
	sessionStarter stripeadapter.SessionStarter,
	messageSender kafkaadapter.MessageSender,
	redisAdapter interface {
		redisadapter.StringRecordStorer
		redisadapter.StringRecordGetter
	},
	p *kafka.Producer, rdb *redis.Client) *PaymentService {
	return &PaymentService{
		orderStatusUpdater: orderStatusUpdater,
		sessionStarter:     sessionStarter,
		messageSender:      messageSender,
		stringRecordStorer: redisAdapter,
		stringRecordGetter: redisAdapter,
		producer:           p,
	}
}

func (s *PaymentService) ProcessPayment(
	ctx context.Context,
	orderId primitive.ObjectID,
	userId string) (string, error) {

	order, err := s.orderGetter.GetOrder(ctx, orderId)
	if err != nil {
		return "", err
	}

	productObjectId, err := primitive.ObjectIDFromHex(order.Product.String())
	if err != nil {
		return "", err
	}

	product, err := s.productGetter.GetProduct(ctx, productObjectId)
	if err != nil {
		return "", err
	}

	// construce stripe param for checkout request
	stripeParams := []*stripe.CheckoutSessionLineItemParams{
		{
			Price:    stripe.String(product.StripePriceId),
			Quantity: stripe.Int64(1),
		},
	}

	// tell stripe about the expected payment
	checkoutId, checkoutUrl, err := startSession(stripeParams)
	if err != nil {
		return "", err
	}

	// set order to payment pending
	err = s.orderStatusUpdater.UpdateOrderStatus(
		ctx, orderId,
		enums.PaymentPending)
	if err != nil {
		return "", err
	}

	// tell kafka about the checkout id we got from stripe
	// checkoutId as key, orderId as value
	// for cleaner abstraction to check after we get webhook reply
	err = sendMessage(
		s.producer,
		"stripe.checkout_session",
		[]byte(checkoutId),
		[]byte(orderId.Hex()))
	if err != nil {
		return "", err
	}

	logVal := map[string]string{
		"level":     "info",
		"timestamp": time.Now().Format(time.RFC3339),
		"message":   "Stripe session started",
		"productId": product.Id.String(),
		"userId":    userId,
	}

	logValBytes, err := json.Marshal(logVal)
	if err != nil {
		return "", err
	}

	// log stripe session event to kafka to store in db
	err = sendMessage(
		s.producer,
		"logs.payment.checkout",
		[]byte(orderId.Hex()),
		logValBytes,
	)
	return checkoutUrl, nil
}

type PaymentProcessor interface {
	ProcessPayment(
		ctx context.Context,
		orderId primitive.ObjectID,
		userId string) (string, error)
}
