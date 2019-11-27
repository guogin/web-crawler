package main

import "sync"

type SafeCounter struct {
	m   map[string]int
	mux sync.Mutex
}

func (counter *SafeCounter) GetAndIncrement(key string) int {
	counter.mux.Lock()
	defer counter.mux.Unlock()

	v, _ := counter.m[key]
	counter.m[key]++
	return v
}
