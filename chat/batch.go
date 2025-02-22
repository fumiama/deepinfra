package chat

import "strings"

type batch struct {
	lst   *Log
	items []item
}

func (l *Log) newbatch(sz int) *batch {
	if sz == 0 {
		panic("sz cannot be 0")
	}
	return &batch{
		lst:   l,
		items: make([]item, 0, sz),
	}
}

func (cl *batch) String() string {
	sb := strings.Builder{}
	for _, item := range cl.items {
		sb.WriteString(cl.lst.sep)
		item.writeToBuilder(&sb, cl.lst.atprefix, cl.lst.namel, cl.lst.namer)
	}
	return sb.String()[len(cl.lst.sep):]
}

// add without mutex
func (cl *batch) add(item item) *batch {
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
