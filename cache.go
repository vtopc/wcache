package wcache

import (
	"context"
	"sync"
	"time"
)

type Cache struct {
	expireFn   ExpireFn
	defaultTTL time.Duration
	m          mapType

	// thread safe:
	mu         sync.RWMutex
	globalDone <-chan struct{} // context cancel
	wg         sync.WaitGroup

	// thread unsafe(TODO):
	CompareFn CompareFn

	// TODO: add reset TTL on hit/set
	// resetTTL bool
}

type mapType map[string]item

// New creates fully functional cache.
//
// ctx is used for shutdown.
func New(ctx context.Context, defaultTTL time.Duration, expireFn ExpireFn) *Cache {
	if expireFn == nil {
		panic("expireFn can't be nil")
	}

	return &Cache{
		m:          make(mapType),
		globalDone: ctx.Done(),
		defaultTTL: defaultTTL,
		expireFn:   expireFn,
	}
}

// SetWithTTL sets the value for a key with default TTL
func (c *Cache) Set(key string, value interface{}) {
	c.SetWithTTL(key, value, c.defaultTTL)
}

// SetWithTTL sets the value for a key.
func (c *Cache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	i, found := c.m[key]
	if found {
		// update
		i.set(value)

		return
	}

	// create new one
	c.wg.Add(1)
	i = newItem()
	c.m[key] = i
	go c.runVault(key, value, i, ttl)
}

// Get returns the value stored in the map for a key, or nil if no
// value is present.
// The found result indicates whether value was found in the cache.
func (c *Cache) Get(key string) (value interface{}, found bool) {
	c.mu.RLock()
	i, found := c.m[key]
	c.mu.RUnlock()

	if !found {
		return nil, false
	}

	return i.get(), true
}

// Delete deletes the value for a key.
func (c *Cache) Delete(key string) {
	i, found := c.delete(key)
	if !found {
		// nothing to do
		return
	}

	i.delete()
}

// Done returns a channel that will be closed when work done(expireFn called for all records).
// This channel would not be closed if context is not canceled(expired, etc.).
func (c *Cache) Done() <-chan struct{} {
	done := make(chan struct{})

	go func() {
		<-c.globalDone
		// after context is canceled, wait for all vaults to be closed:
		c.wg.Wait()
		close(done)
	}()

	return done
}

func (c *Cache) delete(key string) (item, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	i, found := c.m[key]
	if !found {
		return i, false
	}

	delete(c.m, key)

	return i, true
}

// runVault - creates storage for value
func (c *Cache) runVault(key string, value interface{}, i item, ttl time.Duration) {
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
			c.delete(key)
			c.expireFn(key, value)
			return

		case <-c.globalDone:
			// global shutdown
			c.expireFn(key, value)
			return

		case <-i.done:
			// item deleted
			c.expireFn(key, value)
			return
		}
	}
}
