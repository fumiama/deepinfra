package chat

import (
	"fmt"
	"sync"

	"github.com/fumiama/deepinfra"
	"github.com/fumiama/deepinfra/model"
)

type Log[T fmt.Stringer] struct {
	mu                 sync.RWMutex
	batchcap, itemscap int
	itemsep            string
	defaultprompt      string
	m                  map[int64][]*batch[T]
}

func NewLog[T fmt.Stringer](batchcap, itemscap int, itemsep, defaultprompt string) Log[T] {
	if batchcap < 2 {
		panic("batchcap cannot < 2")
	}
	if batchcap%2 != 0 {
		panic("batchcap % 2 must be 0")
	}
	if itemscap < 1 {
		panic("itemscap cannot < 1")
	}
	return Log[T]{
		batchcap:      batchcap,
		itemscap:      itemscap,
		itemsep:       itemsep,
		defaultprompt: defaultprompt,
		m:             make(map[int64][]*batch[T], 64),
	}
}

func (l *Log[T]) Add(grp int64, item T, isbot bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	msgs, ok := l.m[grp]
	if !ok {
		msgs = make([]*batch[T], 1, l.batchcap)
		msgs[0] = l.newbatch(l.itemscap).add(item)
		l.m[grp] = msgs
		return
	}
	isprevusr := len(msgs)%2 != 0
	if (isprevusr && !isbot) || (!isprevusr && isbot) { // is same
		_ = msgs[len(msgs)-1].add(item)
		return
	}
	if len(msgs) < cap(msgs) {
		msgs = append(msgs, l.newbatch(l.itemscap).add(item))
		l.m[grp] = msgs
		return
	}
	copy(msgs, msgs[2:])
	msgs[len(msgs)-2] = l.newbatch(l.itemscap).add(item)
	l.m[grp] = msgs[:len(msgs)-1]
}

func (l *Log[T]) Modelize(p model.Protocol, grp int64, sysp string, isusersystem bool) deepinfra.Model {
	m := p
	if sysp != "" && !isusersystem {
		m.System(sysp)
	}
	l.mu.RLock()
	defer l.mu.RUnlock()
	sz := len(l.m[grp])
	if sz == 0 {
		return m.User(l.defaultprompt)
	}
	for i, msg := range l.m[grp] {
		if i%2 == 0 { // is user
			if i == 0 && isusersystem {
				_ = m.User(sysp + "\n\n" + msg.String())
				continue
			}
			_ = m.User(msg.String())
			continue
		}
		_ = m.Assistant(msg.String())
	}
	return m
}

// Modelize into any type from index and message
func Modelize[X any, T fmt.Stringer](l *Log[T], grp int64, f func(int, string) X) []X {
	l.mu.RLock()
	defer l.mu.RUnlock()
	sz := len(l.m[grp])
	if sz == 0 {
		return []X{f(0, l.defaultprompt)}
	}
	t := make([]X, sz)
	for i, msg := range l.m[grp] {
		t[i] = f(i, msg.String())
	}
	return t
}

// Reset clears all conversation logs while preserving configuration
func (l *Log[T]) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.m = make(map[int64][]*batch[T], 64)
}

// ResetIn removes specified groups from the conversation logs
func (l *Log[T]) ResetIn(grps ...int64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, grp := range grps {
		delete(l.m, grp)
	}
}
