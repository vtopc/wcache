package wcache_test

import (
	"context"
	"fmt"
	"time"

	"github.com/vtopc/wcache"
)

func ExampleCompareFn() {
	const key = "test"

	c := wcache.New(context.Background(), time.Hour, wcache.NoopExpire)
	c.CompareFn = wcache.MaxInt64

	c.Set(key, int64(1))
	c.Set(key, int64(2))

	v, _ := c.Get(key)
	fmt.Println(v)

	// Output:
	// 2
}
