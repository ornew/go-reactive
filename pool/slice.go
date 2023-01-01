package pool

import (
	"sync"
)

type Slice[T any] struct {
	pool sync.Pool
}

func (p *Slice[T]) Get(l int) ([]T, func([]T)) {
	ptr := p.pool.Get().(*[]T)
	arr := *ptr
	arr = arr[0:0]
	if l > cap(arr) {
		arr = make([]T, 0, l*2)
	}
	put := func(arr []T) {
		*ptr = arr
		p.pool.Put(ptr)
	}
	return arr, put
}

func NewSlice[T any](icap int) Slice[T] {
	return Slice[T]{
		pool: sync.Pool{
			New: func() interface{} {
				stack := make([]T, 0, icap)
				return &stack
			},
		},
	}
}
