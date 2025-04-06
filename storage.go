package main

import (
	"sync"
)

var (
	payments = make(map[string]*Payment)
	mu       sync.RWMutex
)

func SavePayment(p *Payment) {
	mu.Lock()
	defer mu.Unlock()
	payments[p.ID] = p
}

func GetPayment(id string) (*Payment, bool) {
	mu.RLock()
	defer mu.RUnlock()
	p, ok := payments[id]
	return p, ok
}

func UpdatePaymentStatus(id string, status string) {
	mu.Lock()
	defer mu.Unlock()
	if p, ok := payments[id]; ok {
		p.Status = status
	}
}
