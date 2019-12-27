package wcache_test

import (
	"context"
	"fmt"
	"time"

	"github.com/vtopc/wcache"
)

func ExampleCache() {
	// will be called when record is expired:
	expireFn := func(value interface{}) {
		fmt.Println(value.(string))
	}

	c := wcache.New(context.Background(), time.Second, expireFn)
	// with custom TTL:
	_ = c.SetWithTTL(2, "to expire second", 2*time.Second)
	// with default TTL:
	_ = c.Set(1, "to expire first")

	time.Sleep(3 * time.Second)

	// Output:
	// to expire first
	// to expire second
}

func ExampleCache_Sync() {

}
