package limiter

import (
	"time"
)

type Limiter struct {
	rate  int           // tps in a second
	begin time.Time     // time start
	cycle time.Duration // time recycle period
	count int
}

func (l *Limiter) Allow() bool {
	if l.rate == 0 {
		return false
	}
	if l.count == l.rate-1 {
		now := time.Now()
		if now.Sub(l.begin) >= l.cycle {
			l.Reset(now)
			return true
		} else {
			return false
		}
	} else {
		l.count++
		return true
	}
}

func (l *Limiter) Set(r int, cycle time.Duration) {
	l.rate = r
	l.begin = time.Now()
	l.cycle = cycle
	l.count = 0
}

func (l *Limiter) Reset(t time.Time) {
	l.begin = t
	l.count = 0
}
