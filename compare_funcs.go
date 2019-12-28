package wcache

// CompareFn compares values on key collisions
type CompareFn func(old, new interface{}) (result interface{})
