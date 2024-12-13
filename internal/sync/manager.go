package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type SyncManager struct {
	db         *gorm.DB
	redis      *redis.Client
	workerPool *WorkerPool
	syncQueue  chan SyncJob
	wg         sync.WaitGroup
}

type SyncJob struct {
	UserID    uuid.UUID
	Operation string
	Data      interface{}
	Timestamp time.Time
}

func NewSyncManager(db *gorm.DB, redis *redis.Client) *SyncManager {
	sm := &SyncManager{
		db:         db,
		redis:      redis,
		syncQueue:  make(chan SyncJob, 1000),
		workerPool: NewWorkerPool(20), // 20 workers
	}
	sm.startWorkers()
	return sm
}

func (sm *SyncManager) startWorkers() {
	for i := 0; i < sm.workerPool.Size(); i++ {
		sm.wg.Add(1)
		go sm.worker()
	}
}

func (sm *SyncManager) worker() {
	defer sm.wg.Done()
	for job := range sm.syncQueue {
		sm.processSync(job)
	}
}

func (sm *SyncManager) processSync(job SyncJob) {
	ctx := context.Background()
	key := fmt.Sprintf("sync:%s:%s", job.UserID, job.Operation)

	// Use Redis as a distributed lock
	lock := sm.redis.SetNX(ctx, key+":lock", "1", 30*time.Second)
	if !lock.Val() {
		return // Another worker is processing this sync
	}
	defer sm.redis.Del(ctx, key+":lock")

	// Start transaction
	tx := sm.db.Begin()
	if tx.Error != nil {
		return
	}

	// Process based on operation type
	var err error
	switch job.Operation {
	case "exam_submission":
		err = sm.processExamSubmission(tx, job)
	case "user_progress":
		err = sm.processUserProgress(tx, job)
	}

	if err != nil {
		tx.Rollback()
		// Store failed sync for retry
		sm.storeFailedSync(job)
		return
	}

	tx.Commit()
}
