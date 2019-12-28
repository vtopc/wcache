package wcache

// CompareFn compares values on key collisions
type CompareFn func(old, new interface{}) (result interface{})

// MaxInt64 a CompareFn typed function that will return biggest value as result
func MaxInt64(old, new interface{}) (result interface{}) {
	o := old.(int64)
	n := new.(int64)

	if n > o {
		return n
	}

	return o
}

// AddInt64 a CompareFn typed function that will add new value to the old one and return as result
func AddInt64(old, new interface{}) (result interface{}) {
	o := old.(int64)
	n := new.(int64)

	return o + n
}
