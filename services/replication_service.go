package services

import (
	"context"
	"log"
	"time"
)

// ReplicationService handles background replication of counter updates
type ReplicationService struct {
	counter *DistributedCounter
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewReplicationService creates a new instance of ReplicationService
func NewReplicationService(counter *DistributedCounter) *ReplicationService {
	ctx, cancel := context.WithCancel(context.Background())
	return &ReplicationService{
		counter: counter,
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Start begins the background replication process
func (rs *ReplicationService) Start() {
	go rs.replicationLoop()
	log.Println("Counter replication service started")
}

// Stop stops the background replication process
func (rs *ReplicationService) Stop() {
	rs.cancel()
	log.Println("Counter replication service stopped")
}

// replicationLoop runs the replication process in a loop
func (rs *ReplicationService) replicationLoop() {
	ticker := time.NewTicker(5 * time.Second) // Replicate every 5 seconds
	defer ticker.Stop()

	for {
		select {
		case <-rs.ctx.Done():
			return
		case <-ticker.C:
			if err := rs.counter.ReplicateToMongoDB(); err != nil {
				log.Printf("Replication error: %v", err)
			}
		}
	}
}
