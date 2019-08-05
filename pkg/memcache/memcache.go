package memcache

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Cache is a thread-safe store for objects which are valid until a
// client-specified time.
type Cache interface {
	// Put stores the provided value under the key. The value is valid until
	// the specified time, becoming invalid at the specific time.
	Put(key string, value interface{}, validUntil time.Time)

	// Get returns the value stored under the provided key. If the value under
	// the provided key is no longer valid, a nil value is returned and the
	// valid boolean is set to false. If the provided key doesn't exist, an
	// error is returned.
	Get(key string) (value interface{}, valid bool, err error)
}

// cache is an in-memory implementation of Cache.
type cache struct {
	m            map[string]*expiringVal
	sync.RWMutex // guards m
}

// expiringVal is an arbitrary type with an associated validity lifetime.
type expiringVal struct {
	v          interface{}
	validUntil time.Time
}

// New initializes and returns a new cache.
func New() Cache {
	return &cache{
		m: make(map[string]*expiringVal),
	}
}

func (c *cache) Put(k string, v interface{}, validUntil time.Time) {
	c.Lock()
	defer c.Unlock()

	c.m[k] = &expiringVal{
		v:          v,
		validUntil: validUntil,
	}
}

func (c *cache) Get(k string) (interface{}, bool, error) {
	c.RLock()
	defer c.RUnlock()

	v, ok := c.m[k]
	if !ok {
		return nil, false, errors.New(fmt.Sprintf("key does not exist in cache: %s", k))
	}

	// the validity lifetime of a cached object is exclusive
	if !time.Now().Before(v.validUntil) {
		return nil, false, nil
	}

	return v.v, true, nil
}
