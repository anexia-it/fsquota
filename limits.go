package fsquota

import "sync"

// Limits contains quota limits
type Limits struct {
	// Byte usage limits
	Bytes Limit

	// File count limits
	Files Limit
}

// Limit represents a combined hard and soft limit
type Limit struct {
	mu   sync.Mutex
	soft *uint64
	hard *uint64
}

// SetHard sets the hard limit
func (l *Limit) SetHard(limit uint64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.hard = &limit
}

// GetHard retrieves the hard limit
func (l *Limit) GetHard() (limit uint64) {
	limit, _, _ = l.getValues()
	return
}

// SetSoft sets the soft limit
func (l *Limit) SetSoft(limit uint64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.soft = &limit
}

// GetSoft retrieves the soft limit
func (l *Limit) GetSoft() (limit uint64) {
	_, limit, _ = l.getValues()
	return
}

func (l *Limit) getValues() (hard, soft uint64, ok bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.hard != nil {
		hard = *l.hard
		ok = true
	}

	if l.soft != nil {
		soft = *l.soft
		ok = true
	}

	return
}
