package stack

import (
	"context"
	"log"
	"time"
)

type Scheduler struct {
	interval time.Duration
	service  *Service
}

func NewScheduler(interval time.Duration, service *Service) *Scheduler {
	return &Scheduler{interval: interval, service: service}
}

func (s *Scheduler) Run(ctx context.Context) {
	if s == nil || s.service == nil {
		log.Printf("level=ERROR msg=\"scheduler run skipped due to nil dependency\"")
		return
	}

	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("level=ERROR msg=\"scheduler panic recovered\" panic=%v", rec)
		}
	}()

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.service.CleanupExpiredAndOrphaned(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.service.CleanupExpiredAndOrphaned(ctx)
		}
	}
}
