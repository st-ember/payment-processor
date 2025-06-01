package request

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StartPaymentReq struct {
	OrderId       primitive.ObjectID
	StripePriceId string
	Quantity      int64
	UserId        primitive.ObjectID
}
