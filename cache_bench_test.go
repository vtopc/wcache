package wcache

import (
	"context"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func BenchmarkSameKey(b *testing.B) {
	const (
		key   = "exists"
		value = "some test value"
	)

	c := New(context.Background(), time.Hour, NoopExpire)

	b.Run("Set", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			c.Set(key, value)
		}
	})

	b.Run("Get", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			c.Get(key)
		}
	})

	b.Run("Set_Delete", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			c.Set(key, value)
			c.Delete(key)
		}
	})
}

func BenchmarkRandomKeys(b *testing.B) {
	const (
		value = "some test value"
		keys  = 1000
	)

	c := New(context.Background(), time.Hour, NoopExpire)

	b.Run("Set", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			c.Set(strconv.Itoa(rand.Intn(keys)), value)
		}
	})

	b.Run("Get", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			c.Get(strconv.Itoa(rand.Intn(keys)))
		}
	})

	b.Run("Set_Delete", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			key := strconv.Itoa(rand.Intn(keys))
			c.Set(key, value)
			c.Delete(key)
		}
	})
}

func BenchmarkSetSameKeyWithTTL(b *testing.B) {
	const (
		key   = "exists"
		value = "some test value"
	)
	ttl := time.Microsecond

	c := New(context.Background(), ttl, NoopExpire)

	for n := 0; n < b.N; n++ {
		c.Set(key, value)
		time.Sleep(ttl)
	}
}

func BenchmarkCompareFn(b *testing.B) {
	const (
		key   = "exists"
		value = "some test value"
	)

	compareFn := func(old, new interface{}) (result interface{}) {
		time.Sleep(time.Microsecond)
		return new
	}

	c := New(context.Background(), time.Hour, NoopExpire)
	c.CompareFn = compareFn

	for n := 0; n < b.N; n++ {
		c.Set(key, value)
	}
}
