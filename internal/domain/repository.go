package domain

import (
	"context"
	"paymentprocessor/internal/domain/entity"

	"github.com/stripe/stripe-go/v72"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SessionRepo interface {
	Insert(ctx context.Context, orderId primitive.ObjectID, sessionId string) error

	GetBySessionId(ctx context.Context, sessionId string) (entity.StripeCheckoutSession, error)

	UpdateStatus(ctx context.Context, sessionId string, newStatus stripe.CheckoutSessionStatus) error

	BulkSetExpire(ctx context.Context, ids []primitive.ObjectID) error

	ListLatest(ctx context.Context) ([]entity.StripeCheckoutSession, error)
}
