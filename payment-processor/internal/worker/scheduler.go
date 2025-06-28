package worker

import (
	"context"
	"paymentprocessor/internal/domain"
	kafkaadapter "paymentprocessor/internal/infra/kafka"
	"time"
)

type Scheduler struct {
	sessionRepo domain.SessionRepo
	kafkaClient *kafkaadapter.KafkaClient
}

func NewScheduler(sessionRepo domain.SessionRepo, kafkaClient *kafkaadapter.KafkaClient) *Scheduler {
	return &Scheduler{
		sessionRepo: sessionRepo,
		kafkaClient: kafkaClient,
	}
}

func (s *Scheduler) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cleanupTicker := time.NewTicker(12 * time.Hour)
	defer cleanupTicker.Stop()

	cleanupWorker := NewCleanupWorker(s.sessionRepo, s.kafkaClient)

	for range cleanupTicker.C {
		cleanupWorker.Cleanup(ctx)
	}
}
