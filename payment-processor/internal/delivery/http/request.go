package http

import (
	"paymentprocessor/internal/domain/request"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StartPaymentReq struct {
	OrderId       primitive.ObjectID `json:"order_id"`
	StripePriceId string             `json:"stripe_price_id"`
	Quantity      int64              `json:"quantity"`
	UserId        primitive.ObjectID `json:"user_id"`
}

// map http model to domain model
func (h *StartPaymentReq) ToUsecaseRequest() request.StartPaymentReq {
	return request.StartPaymentReq{
		OrderId:       h.OrderId,
		StripePriceId: h.StripePriceId,
		Quantity:      h.Quantity,
		UserId:        h.UserId,
	}
}
