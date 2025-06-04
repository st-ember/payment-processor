package entity

import (
	"paymentprocessor/internal/domain/enum"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StripeCheckoutSession struct {
	Id        primitive.ObjectID // order_id
	SessionId string
	CreatedAt time.Time
	UpdatedAt time.Time
	Status    enum.StripeStatus
}
