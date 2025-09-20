package chat

import (
	"fmt"
	"strings"
)

type batch[T fmt.Stringer] struct {
	lst   *Log[T]
	items []T
}

func (l *Log[T]) newbatch(sz int) *batch[T] {
	if sz == 0 {
		panic("sz cannot be 0")
	}
	return &batch[T]{
		lst:   l,
		items: make([]T, 0, sz),
	}
}

func (cl *batch[T]) String() string {
	sb := strings.Builder{}
	for _, item := range cl.items {
		sb.WriteString(cl.lst.itemsep)
		sb.WriteString(item.String())
	}
	return sb.String()[len(cl.lst.itemsep):]
}

// add without mutex
func (cl *batch[T]) add(item T) *batch[T] {
	v := cl.items
	defer func() {
		cl.items = v
	}()
	if cap(v) == 1 {
		v[0] = item
		return cl
	}
	if len(v) < cap(v) {
		v = append(v, item)
		return cl
	}
	copy(v, v[1:])
	v[len(v)-1] = item
	return cl
}
