package wcache

type item struct {
	setter chan interface{}
	getter chan interface{}
}

func newItem() item {
	return item{
		setter: make(chan interface{}),
		getter: make(chan interface{}),
	}
}

func (i item) get() interface{} {
	return <-i.getter
}

func (i item) set(value interface{}) {
	i.setter <- value
}
