package domain

import (
	"context"
	"paymentprocessor/internal/domain/enum"
	"paymentprocessor/internal/domain/payment"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SessionRepo interface {
	Insert(ctx context.Context, orderId primitive.ObjectID, sessionId string) error

	GetBySessionId(ctx context.Context, sessionId string) (payment.StripeCheckoutSession, error)

	UpdateStatus(ctx context.Context, sessionId string, newStatus enum.StripeStatus) error

	Delete(ctx context.Context, sessionId primitive.ObjectID) error
}
