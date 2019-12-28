package wcache

import (
	"context"
	"sync"
	"time"
)

type Cache struct {
	// thread safe:
	done       <-chan struct{}
	m          sync.Map
	expireFn   ExpireFn
	defaultTTL time.Duration
	wg         sync.WaitGroup

	// thread unsafe(TODO):
	CompareFn CompareFn

	// TODO: add reset TTL on hit/set
	// resetTTL bool
}

// CompareFn compares values on key collisions
type CompareFn func(old, new interface{}) (result interface{})

// ExpireFn a callback that will be called when record is expired
type ExpireFn func(key, value interface{})

// New creates fully functional cache.
//
// ctx is used for shutdown and triggering expireFn for all records ignoring the TTL.
func New(ctx context.Context, defaultTTL time.Duration, expireFn ExpireFn) *Cache {
	if expireFn == nil {
		panic("expireFn can't be nil")
	}

	return &Cache{
		done:       ctx.Done(),
		defaultTTL: defaultTTL,
		expireFn:   expireFn,
	}
}

// SetWithTTL sets the value for a key with default TTL
func (c *Cache) Set(key, value interface{}) {
	c.SetWithTTL(key, value, c.defaultTTL)
}

// SetWithTTL sets the value for a key.
func (c *Cache) SetWithTTL(key, value interface{}, ttl time.Duration) {
	i, found := c.get(key)
	if found {
		// update
		i.set(value)

		return
	}

	// create new one
	i = newItem()
	go c.runVault(key, value, i, ttl)
	c.m.Store(key, i)

	return
}

// Get returns the value stored in the map for a key, or nil if no
// value is present.
// The found result indicates whether value was found in the map.
func (c *Cache) Get(key interface{}) (value interface{}, found bool) {
	i, found := c.get(key)
	if !found {
		return nil, false
	}

	return i.get(), true
}

// TODO: add Delete with triggering expire func

// TODO: add instead of context?
// func (c *Cache) Sync() {
// 	c.m.Range()
// }

func (c *Cache) get(key interface{}) (item, bool) {
	v, found := c.m.Load(key)
	if !found {
		return item{}, false
	}

	i, ok := v.(item)
	if !ok {
		panic("internal error: failed to assert item")
	}

	return i, true
}

// runVault - creates storage for value
func (c *Cache) runVault(key, value interface{}, i item, ttl time.Duration) {
	c.wg.Add(1)
	defer c.wg.Done()

	timer := time.NewTimer(ttl)
	defer timer.Stop()

	for {
		select {
		case newValue := <-i.setter:
			if c.CompareFn != nil {
				value = c.CompareFn(value, newValue)
			} else {
				// overwrite
				value = newValue
			}

		case i.getter <- value:

		case <-timer.C:
			c.m.Delete(key)
			c.expireFn(key, value)
			return

		case <-c.done:
			// global shutdown
			c.expireFn(key, value)
			return
		}
	}
}
