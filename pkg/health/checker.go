package health

import (
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Checker struct {
	db    *gorm.DB
	redis *redis.Client
	mu    sync.RWMutex
	state map[string]Status
}

type Status struct {
	Healthy bool
	Message string
	Time    time.Time
}

func NewChecker(db *gorm.DB, redis *redis.Client) *Checker {
	c := &Checker{
		db:    db,
		redis: redis,
		state: make(map[string]Status),
	}
	go c.startPeriodicChecks()
	return c
}

func (c *Checker) startPeriodicChecks() {
	ticker := time.NewTicker(30 * time.Second)
	for range ticker.C {
		c.checkDatabase()
		c.checkRedis()
		c.checkDiskSpace()
	}
}

func (c *Checker) checkDatabase() {
	sqlDB, err := c.db.DB()
	status := Status{Time: time.Now()}

	if err != nil {
		status.Healthy = false
		status.Message = "Failed to get database instance"
	} else if err := sqlDB.Ping(); err != nil {
		status.Healthy = false
		status.Message = "Database ping failed"
	} else {
		status.Healthy = true
		status.Message = "Database is healthy"
	}

	c.mu.Lock()
	c.state["database"] = status
	c.mu.Unlock()
}
