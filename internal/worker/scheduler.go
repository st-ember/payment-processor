package worker

import (
	"context"
	"paymentprocessor/internal/domain"
	"time"
)

type Scheduler struct {
	sessionRepo domain.SessionRepo
}

func NewScheduler(sessionRepo domain.SessionRepo) *Scheduler {
	return &Scheduler{sessionRepo: sessionRepo}
}

func (s *Scheduler) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cleanupTicker := time.NewTicker(12 * time.Hour)
	defer cleanupTicker.Stop()

	cleanupWorker := NewCleanupWorker(s.sessionRepo)

	for range cleanupTicker.C {
		cleanupWorker.Cleanup(ctx)
	}
}
