package mongodb

import (
	"paymentprocessor/internal/domain/enum"
	"paymentprocessor/internal/domain/payment"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StripeCheckoutSession struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"` // order id
	SessionId string             `bson:"session_id"`
	Status    string             `bson:"status"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func (d *StripeCheckoutSession) ToDomainModel() (payment.StripeCheckoutSession, error) {
	status, err := enum.ParseStripeStatus(d.Status)

	if err != nil {
		return payment.StripeCheckoutSession{}, err
	}

	return payment.StripeCheckoutSession{
		Id:        d.Id,
		SessionId: d.SessionId,
		Status:    status,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}, nil
}
