package wcache_test

import (
	"context"
	"fmt"
	"time"

	"github.com/vtopc/wcache"
)

func ExampleCache() {
	// will be called when record is expired:
	expireFn := func(key, value interface{}) {
		fmt.Printf("%d: %s\n", key, value)
	}

	c := wcache.New(context.Background(), 100*time.Millisecond, expireFn)
	// with custom TTL:
	_ = c.SetWithTTL(2, "to expire second", 200*time.Millisecond)
	// with default TTL:
	_ = c.Set(1, "to expire first")

	time.Sleep(300 * time.Millisecond)

	// Output:
	// 1: to expire first
	// 2: to expire second
}

func ExampleCache_Sync() {

}
