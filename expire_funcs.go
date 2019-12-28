package wcache

// ExpireFn a callback that will be called when record is expired
type ExpireFn func(key, value interface{})

// NoopExpire does nothing
func NoopExpire(key, value interface{}) {}

// TODO: add ChanExpire
