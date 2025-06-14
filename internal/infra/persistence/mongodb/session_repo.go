package mongodb

import (
	"context"
	"paymentprocessor/internal/domain/entity"
	"time"

	"github.com/stripe/stripe-go/v72"
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
		Status:    string(stripe.CheckoutSessionStatusOpen),
		CreatedAt: time.Now(),
	}

	_, err := r.collection.InsertOne(ctx, session)
	if err != nil {

		return err
	}

	return nil
}

func (r *SessionRepo) GetBySessionId(ctx context.Context, sessionId string) (entity.StripeCheckoutSession, error) {
	var session StripeCheckoutSession
	filter := bson.M{"session_id": sessionId}
	err := r.collection.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		return entity.StripeCheckoutSession{}, err
	}

	paymentModel, err := session.ToDomainModel()
	if err != nil {
		return entity.StripeCheckoutSession{}, err
	}

	return paymentModel, nil
}

func (r *SessionRepo) UpdateStatus(ctx context.Context, sessionId string, newStatus stripe.CheckoutSessionStatus) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{
		"session_id": sessionId,
	}
	update := bson.M{
		"$set": bson.M{
			"status":     newStatus,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *SessionRepo) BulkSetExpire(ctx context.Context, ids []primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	now := time.Now()
	filter := bson.M{
		"_id": bson.M{"$in": ids},
	}
	update := bson.M{
		"$set": bson.M{
			"status":     stripe.CheckoutSessionStatusExpired,
			"updated_at": now,
		},
	}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *SessionRepo) ListLatest(ctx context.Context) ([]entity.StripeCheckoutSession, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var sessions []entity.StripeCheckoutSession

	timeThreshold := time.Now().Add(-25 * time.Hour)
	filter := bson.M{
		"created_at": bson.M{"$gte": timeThreshold},
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return sessions, err
	}

	err = cursor.All(ctx, &sessions)
	if err != nil {
		return sessions, err
	}

	return sessions, nil
}
