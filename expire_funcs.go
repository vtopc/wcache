package wcache

import (
	"fmt"
)

// ExpireFn a callback that will be called when record is expired
type ExpireFn func(key, value interface{})

type KeyValue struct {
	Key   interface{}
	Value interface{}
}

// NoopExpire does nothing
func NoopExpire(key, value interface{}) {}

// PrintOnExpire a dummy ExpireFn that will print key, value when record is expired
func PrintlnOnExpire(key, value interface{}) {
	fmt.Printf("%s: %v\n", key, value)
}

// ChanExpire returns a ExpireFn that will send key, value to the channel `ch` when record is expired
func ChanExpire(ch chan<- KeyValue) ExpireFn {
	return func(key, value interface{}) {
		ch <- KeyValue{
			Key:   key,
			Value: value,
		}
	}
}
