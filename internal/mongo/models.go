package mongo

import (
	"paymentprocessor/internal/enums"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	Id          uuid.UUID         `bson:"_id,omitempty"`
	Product     uuid.UUID         `bson:"product"`
	Customer    uuid.UUID         `bson:"customer"`
	PurchasedAt time.Time         `bson:"purchased_at"`
	Status      enums.OrderStatus `bson:"status"`
}

type Product struct {
	Id            uuid.UUID `bson:"_id,omitempt"`
	Category      uuid.UUID `bson:"category"`
	Name          string    `bson:"name"`
	Description   string    `bson:"description"`
	Merchant      uuid.UUID `bson:"merchant"`
	Price         float64   `bson:"price"`
	CreatedAt     time.Time `bson:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at"`
	StripePriceId string    `bson:"stripe_price_id"`
}
