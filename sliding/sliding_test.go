package sliding

import (
	"fmt"
	"testing"
	"time"
)

func TestSlidingLimiter_Allow(t *testing.T) {
	limiter := NewSlidingLimiter(1, 1*time.Second)
	for {
		if limiter.Allow() {
			fmt.Println(time.Now().Format("15:04:05"))
		}
		time.Sleep(100 * time.Millisecond)
	}
}
