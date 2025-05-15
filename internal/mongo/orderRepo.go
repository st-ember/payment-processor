package mongo

import (
	"context"
	"paymentprocessor/internal/enums"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type orderRepo struct {
	collection *mongo.Collection
}

func NewOrderRepository(client *mongo.Client, dbName, collectionName string) *orderRepo {
	coll := client.Database(dbName).Collection(collectionName)
	return &orderRepo{collection: coll}
}

func (r *orderRepo) UpdateOrderStatus(ctx context.Context, orderId primitive.ObjectID, newStatus enums.OrderStatus) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"status": newStatus,
		},
	}

	_, err := r.collection.UpdateByID(ctx, orderId, update)
	if err != nil {
		return err
	}

	return nil
}

type OrderStatusUpdater interface {
	UpdateOrderStatus(ctx context.Context, orderId primitive.ObjectID, newStatus enums.OrderStatus) error
}
