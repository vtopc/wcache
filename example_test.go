package wcache_test

import (
	"context"
	"time"

	"github.com/vtopc/wcache"
)

func Example() {
	c := wcache.New(context.Background(), 100*time.Millisecond, wcache.PrintlnOnExpire)
	// with custom TTL:
	c.SetWithTTL(2, "to expire second", 200*time.Millisecond)
	// with default TTL:
	c.Set(1, "to expire first")

	time.Sleep(300 * time.Millisecond)

	// Output:
	// 1: to expire first
	// 2: to expire second
}
