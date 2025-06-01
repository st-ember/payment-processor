package payment

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConcludedPayment struct {
	Id          primitive.ObjectID
	OrderId     primitive.ObjectID
	ConcludedAt time.Time
}
