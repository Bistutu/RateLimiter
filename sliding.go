package main

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

func NewSlidingLimiter(maxRequest int, window time.Duration) *SlidingLimiter {
	return &SlidingLimiter{
		maxRequests: maxRequest,
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

func main() {

	limiter := NewSlidingLimiter(1, 1*time.Second) // 创建一个新的滑动窗口限流器，每秒允许5个请求

	for i := 0; i < 100; i++ {
		go func() {
			if limiter.Allow() {
				println(i, "Request allowed") // 请求被允许
			} else {
				println("Request denied") // 请求被拒绝
			}
		}()
		time.Sleep(100 * time.Millisecond) // 每100毫秒发送一个请求
	}

	time.Sleep(2 * time.Second)
}
