package mongodb

import (
	"paymentprocessor/internal/domain/entity"
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

func (d *StripeCheckoutSession) ToDomainModel() (entity.StripeCheckoutSession, error) {
	return entity.StripeCheckoutSession{
		Id:        d.Id,
		SessionId: d.SessionId,
		Status:    d.Status,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}, nil
}
