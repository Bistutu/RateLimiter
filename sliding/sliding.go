package sliding

import (
	"sync"
	"time"
)

// SlidingLimiter 滑动窗口限流器
type SlidingLimiter struct {
	mu          sync.Mutex
	maxRequests int
	window      time.Duration
	timestamps  []time.Time
}

func NewSlidingLimiter(maxRequests int, window time.Duration) *SlidingLimiter {
	return &SlidingLimiter{
		maxRequests: maxRequests,
		window:      window,
		timestamps:  make([]time.Time, 0),
	}
}

// Allow 请求是否被运行
func (l *SlidingLimiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	start := now.Add(-l.window)
	// 获取前 window 窗口中的请求数
	i := 0
	for i < len(l.timestamps) && l.timestamps[i].Before(start) {
		i++
	}
	l.timestamps = l.timestamps[i:]

	// 检查是否达到最大请求数
	if len(l.timestamps) < l.maxRequests {
		l.timestamps = append(l.timestamps, time.Now())
		return true
	}
	return false
}
