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

// CompareFn compares values on collisions
type CompareFn func(old, new interface{}) (result interface{})

// ExpireFn will be called when record is expired
type ExpireFn func(key, value interface{})

type item struct {
	update chan<- interface{}
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
	v, exists := c.m.Load(key)
	var i item
	if exists {
		// update
		var ok bool
		i, ok = v.(item)
		if !ok {
			return errors.New("internal error: failed to assert item")
		}
	} else {
		// create new one
		update := make(chan interface{}, 1)
		go c.newVault(key, update, ttl)
		i = item{update: update}
		c.m.Store(key, i)
	}

	i.update <- value

	return nil
}

// TODO:
// // Get returns the value stored in the map for a key, or nil if no
// // value is present.
// // The ok result indicates whether value was found in the map.
// func (c *Cache) Get(key interface{}) (value interface{}, ok bool) {
// 	return c.m.Load(key)
// }

// TODO: add instead of context?
// func (c *Cache) Sync() {
// 	c.m.Range()
// }

// newVault - creates storage for value
func (c *Cache) newVault(key interface{}, update <-chan interface{}, ttl time.Duration) {
	var value interface{}

	timer := time.NewTimer(ttl)
	defer timer.Stop()

	for {
		select {
		case newValue := <-update:
			if c.CompareFn != nil {
				value = c.CompareFn(value, newValue)
			} else {
				// overwrite
				value = newValue
			}

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
