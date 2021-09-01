package cache

import (
	"context"
	"log"

	lru "github.com/hashicorp/golang-lru"
)

type LRU struct {
	lru *lru.Cache
}

func New(size int) *LRU {
	cache, err := lru.New(size)
	if err != nil {
		// An error is only returned for non-positive cache size
		// and we already checked for that.
		log.Panicf("unexpected error creating cache: %v", err)
	}
	return &LRU{cache}
}

func (l LRU) Get(ctx context.Context, key string) (value interface{}, ok bool) {
	return l.lru.Get(key)
}

func (l LRU) Add(ctx context.Context, key string, value interface{}) {
	l.lru.Add(key, value)
}
