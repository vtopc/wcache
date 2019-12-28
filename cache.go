package wcache

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Cache struct {
	// thread safe:
	done       <-chan struct{}
	m          sync.Map
	expireFn   ExpireFn
	defaultTTL time.Duration

	// thread unsafe(TODO):
	CompareFn CompareFn

	// TODO: add reset TTL on hit/set
	// resetTTL bool
}

// CompareFn compares values on key collisions
type CompareFn func(old, new interface{}) (result interface{})

// ExpireFn a callback that will be called when record is expired
type ExpireFn func(key, value interface{})

type item struct {
	setter chan interface{}
	getter chan interface{}
	// TODO: add wait group
}

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
func (c *Cache) Set(key, value interface{}) error {
	return c.SetWithTTL(key, value, c.defaultTTL)
}

// SetWithTTL sets the value for a key.
func (c *Cache) SetWithTTL(key, value interface{}, ttl time.Duration) error {
	v, ok := c.m.Load(key)
	if ok {
		// update
		var ok bool
		i, ok := v.(item)
		if !ok {
			return errors.New("internal error: failed to assert item")
		}

		i.setter <- value

		return nil
	}

	// create new one
	i := newItem()
	go c.runVault(key, value, i, ttl)
	c.m.Store(key, i)

	return nil
}

// Get returns the value stored in the map for a key, or nil if no
// value is present.
// The ok result indicates whether value was found in the map.
func (c *Cache) Get(key interface{}) (value interface{}, ok bool) {
	v, ok := c.m.Load(key)
	if !ok {
		return nil, false
	}

	i := v.(item)
	value = <-i.getter

	return value, true
}

// TODO: add instead of context?
// func (c *Cache) Sync() {
// 	c.m.Range()
// }

// runVault - creates storage for value
func (c *Cache) runVault(key, value interface{}, i item, ttl time.Duration) {
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
			c.expireFn(key, value)
			return
		}
	}
}

func newItem() item {
	return item{
		setter: make(chan interface{}),
		getter: make(chan interface{}),
	}
}
