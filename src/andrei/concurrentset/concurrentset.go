package concurrentset

import "sync"

type ConcurrentStringSet struct {
	mutex *sync.RWMutex
	m map[string]bool
}

func NewConcurrentStringSet() ConcurrentStringSet {
	return ConcurrentStringSet{mutex: &sync.RWMutex{}, m: map[string]bool{}}
}

func (cs ConcurrentStringSet) Put(s string){
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	cs.m[s] = true
}

func (cs ConcurrentStringSet) Contains(s string) bool{
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	_, alreadySeen := cs.m[s]
	return alreadySeen
}
