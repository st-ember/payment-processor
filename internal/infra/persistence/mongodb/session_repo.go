package mongodb

import (
	"context"
	"paymentprocessor/internal/domain/enum"
	"paymentprocessor/internal/domain/payment"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SessionRepo struct {
	collection *mongo.Collection
}

func NewSessionRepository(db *mongo.Database) *SessionRepo {
	coll := db.Collection("stripe_checkout_sessions")
	return &SessionRepo{collection: coll}
}

func (r *SessionRepo) Insert(ctx context.Context, orderId primitive.ObjectID, sessionId string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	session := StripeCheckoutSession{
		Id:        orderId,
		SessionId: sessionId,
		Status:    enum.Open.String(),
		CreatedAt: time.Now(),
	}

	_, err := r.collection.InsertOne(ctx, session)
	if err != nil {
		return err
	}

	return nil
}

func (r *SessionRepo) GetById(ctx context.Context, sessionId primitive.ObjectID) (payment.StripeCheckoutSession, error) {
	var session StripeCheckoutSession
	filter := bson.M{"_id": sessionId}
	err := r.collection.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		return payment.StripeCheckoutSession{}, err
	}

	paymentModel, err := session.ToDomainModel()
	if err != nil {
		return payment.StripeCheckoutSession{}, err
	}

	return paymentModel, nil
}

func (r *SessionRepo) UpdateStatus(ctx context.Context, sessionId primitive.ObjectID, newStatus enum.StripeStatus) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"status": newStatus,
		},
	}

	_, err := r.collection.UpdateByID(ctx, sessionId, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *SessionRepo) Delete(ctx context.Context, sessionId primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{
		"_id": sessionId,
	}

	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
