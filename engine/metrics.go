package engine

import "sync"

var WinsMap = NewSafeMap()

type SafeMap struct {
	mu   sync.Mutex // Guards access to the internal map
	data map[int]uint64
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		data: make(map[int]uint64, 3),
	}
}

func (sm *SafeMap) Inc(key int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data[key]++
}

func (sm *SafeMap) Map() map[int]uint64 {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return sm.data
}
