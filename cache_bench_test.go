package wcache

import (
	"context"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
)

func BenchmarkSetSameKey(b *testing.B) {
	const (
		key   = "exists"
		value = "some test value"
	)

	c := New(context.Background(), time.Hour, NoopExpire)

	for n := 0; n < b.N; n++ {
		c.Set(key, value)
	}
}

func BenchmarkSetABitRandomKeys(b *testing.B) {
	const value = "some test value"

	c := New(context.Background(), time.Hour, NoopExpire)

	for n := 0; n < b.N; n++ {
		c.Set(strconv.Itoa(rand.Intn(10)), value)
	}
}

func BenchmarkGetSameKey(b *testing.B) {
	const (
		key   = "exists"
		value = "some test value"
	)

	c := New(context.Background(), time.Hour, NoopExpire)
	c.Set(key, value)

	for n := 0; n < b.N; n++ {
		c.Get(key)
	}
}

func BenchmarkSetDeleteRandom(b *testing.B) {
	const value = "some test value"

	c := New(context.Background(), time.Hour, NoopExpire)

	for n := 0; n < b.N; n++ {
		c.Set(uuid.New().String(), value)
	}
}

func BenchmarkCompareFn(b *testing.B) {
	const (
		key   = "exists"
		value = "some test value"
	)

	compareFn := func(old, new interface{}) (result interface{}) {
		time.Sleep(2 * time.Millisecond)
		return new
	}

	c := New(context.Background(), time.Hour, NoopExpire)
	c.CompareFn = compareFn

	for n := 0; n < b.N; n++ {
		c.Set(key, value)
	}
}
