package worker

import (
	"context"
	"paymentprocessor/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CleanupWorker struct {
	sessionRepo domain.SessionRepo
}

func NewCleanupWorker(sessionRepo domain.SessionRepo) *CleanupWorker {
	return &CleanupWorker{sessionRepo: sessionRepo}
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

	return nil

}
