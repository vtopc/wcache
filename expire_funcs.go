package wcache

import (
	"fmt"
)

// ExpireFn a callback that will be called when record is expired
type ExpireFn func(key, value interface{})

// NoopExpire does nothing
func NoopExpire(key, value interface{}) {}

// PrintOnExpire a dummy ExpireFn that will print key value when record is expired
func PrintlnOnExpire(key, value interface{}) {
	fmt.Printf("%d: %s\n", key, value)
}

// TODO: add ChanExpire
