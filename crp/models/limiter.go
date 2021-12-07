package models

import (
	"sync"
	"time"
)

type Limiter struct {
	IPs map[string]*Counter
	sync.Mutex
}

type Counter struct {
	Count    int
	LastSeen time.Time
}
