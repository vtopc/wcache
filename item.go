package wcache

type item struct {
	setter chan interface{}
	getter chan interface{}
	done   chan struct{}
}

// TODO:
// BenchmarkCompareFn               	     500	   2473886 ns/op	       2 B/op	       0 allocs/op
// BenchmarkCompareFn-2             	     500	   2606114 ns/op	       2 B/op	       0 allocs/op
// BenchmarkCompareFn-4             	     500	   2592972 ns/op	       2 B/op	       0 allocs/op
// BenchmarkCompareFn-8             	     500	   2548593 ns/op	       3 B/op	       0 allocs/op
// BenchmarkCompareFn-16            	     500	   2597770 ns/op	       2 B/op	       0 allocs/op
func newItem() item {
	return item{
		setter: make(chan interface{}),
		getter: make(chan interface{}),
		done:   make(chan struct{}),
	}
}

func (i item) get() interface{} {
	return <-i.getter
}

func (i item) set(value interface{}) {
	i.setter <- value
}

func (i item) delete() {
	close(i.done)
}
