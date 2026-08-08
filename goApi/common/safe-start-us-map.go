package common

import "sync"

type SafeStartUsMap struct {
	M  map[string]int64
	Mu sync.RWMutex
}

// func (m *SafeStartUsMap) Get(key string) (int64, bool) {
// 	m.Mu.RLock()
// 	defer m.Mu.RUnlock()
// 	val, ok := m.M[key]
// 	return val, ok
// }

// func (m *SafeStartUsMap) Set(key string, value int64) {
// 	m.Mu.Lock()
// 	defer m.Mu.Unlock()
// 	m.M[key] = value
// }
func (m *SafeStartUsMap) GetAndSet(key string, incrVal int64) (oldVal, newVal int64, ok bool) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	oldVal, ok = m.M[key]
	if !ok {
		m.M[key] = 0
		ok = true
	}
	m.M[key] += incrVal
	newVal = m.M[key]
	return oldVal, newVal, ok
}
