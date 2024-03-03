package cache

import (
	"sync"
	"time"
)

type inMemoryCache struct {
	dataMap map[int]inMemoryValue
	lock    *sync.Mutex
	clock   time.Time
}

type inMemoryValue struct {
	SetTime    int64
	Expiration int64
}

func New(clock time.Time) *inMemoryCache {
	return &inMemoryCache{
		dataMap: make(map[int]inMemoryValue, 0),
		lock:    &sync.Mutex{},
		clock:   clock,
	}
}

func (c *inMemoryCache) Add(key int, expiration int64) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.dataMap[key] = inMemoryValue{
		SetTime:    time.Now().Unix(),
		Expiration: expiration,
	}
	return nil
}

func (c *inMemoryCache) Get(key int) (bool, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	value, ok := c.dataMap[key]
	if ok && time.Now().Unix()-value.SetTime > value.Expiration {
		return false, nil
	}
	return ok, nil
}

func (c *inMemoryCache) Delete(key int) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.dataMap, key)
}
