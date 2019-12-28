package wcache

type item struct {
	setter chan interface{}
	getter chan interface{}
	done   chan struct{}
}

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
