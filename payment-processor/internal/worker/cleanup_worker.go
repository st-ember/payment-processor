package worker

import (
	"context"
	"paymentprocessor/internal/domain"
	kafkaadapter "paymentprocessor/internal/infra/kafka"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CleanupWorker struct {
	sessionRepo domain.SessionRepo
	kafkaClient *kafkaadapter.KafkaClient
}

func NewCleanupWorker(sessionRepo domain.SessionRepo, kafkaClient *kafkaadapter.KafkaClient) *CleanupWorker {
	return &CleanupWorker{
		sessionRepo: sessionRepo,
		kafkaClient: kafkaClient,
	}
}

func (w *CleanupWorker) Cleanup(ctx context.Context) error {
	latestSessions, err := w.sessionRepo.ListLatest(ctx)
	if err != nil {
		return err
	}

	expiredIds := make([]primitive.ObjectID, 0, len(latestSessions)/10)
	for _, session := range latestSessions {
		if session.UpdatedAt.IsZero() {
			expiredIds = append(expiredIds, session.Id)
		}
	}

	if len(expiredIds) > 0 {
		err = w.sessionRepo.BulkSetExpire(ctx, expiredIds)
		if err != nil {
			return err
		}
	}

	// tell order ms about the expired ids
	batchMsg := map[string]interface{}{
		"timestamp":   time.Now().Format(time.RFC3339),
		"expired_ids": expiredIds,
	}

	err = w.kafkaClient.SendMessage(
		kafkaadapter.Topic.CheckoutStatusBatch,
		batchMsg,
	)
	if err != nil {
		return err
	}

	return nil

}
