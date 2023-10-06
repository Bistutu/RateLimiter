package main

import (
	"sync"
	"time"

	"github.com/golang/groupcache/lru"

	"RateLimiter/sliding"
)

type Robot struct {
	cache     *lru.Cache // 最近最少使用淘汰
	cacheSize int
	limiterMu sync.Mutex
}

func NewRobot(size int) *Robot {
	return &Robot{
		cache:     lru.New(size),
		cacheSize: size,
	}
}

func (r *Robot) getOrCreateLimiter(groupID string) *sliding.SlidingLimiter {
	r.limiterMu.Lock()
	defer r.limiterMu.Unlock()
	val, exists := r.cache.Get(groupID)
	if exists {
		return val.(*sliding.SlidingLimiter)
	}
	limiter := sliding.NewSlidingLimiter(3, 1*time.Second)
	r.cache.Add(groupID, limiter)
	return limiter
}

func (r *Robot) SendMessageToGroup(groupID string, message string) bool {
	limiter := r.getOrCreateLimiter(groupID)
	if !limiter.Allow() {
		return false
	}
	// TODO: Add actual logic to send the message to the group
	return true
}

func main() {
	robot := NewRobot(1000) // Assuming a maximum of 1000 active groups
	groupID := "exampleGroup"
	message := "Hello, World!"
	success := robot.SendMessageToGroup(groupID, message)
	if success {
		// Message was successfully sent
	} else {
		// Rate limit exceeded
	}
}
