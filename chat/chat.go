package chat

import (
	"sync"

	"github.com/fumiama/deepinfra"
	"github.com/fumiama/deepinfra/model"
)

type Log struct {
	mu            sync.RWMutex
	cap           int
	sep           string
	defaultprompt string
	namel, namer  string
	atprefix      string
	m             map[int64][]*batch
}

func NewLog(cap int, sep, defaultprompt, namel, namer, atprefix string) Log {
	if cap < 2 {
		panic("cap cannot < 2")
	}
	if cap%2 != 0 {
		panic("cap % 2 must be 0")
	}
	return Log{
		cap:           cap,
		sep:           sep,
		defaultprompt: defaultprompt,
		namel:         namel,
		namer:         namer,
		atprefix:      atprefix,
		m:             make(map[int64][]*batch, 64),
	}
}

func (l *Log) Add(grp int64, usr, txt string, isbot, isatme bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	msgs, ok := l.m[grp]
	if !ok {
		msgs = make([]*batch, 1, l.cap)
		msgs[0] = l.newbatch(l.cap).add(item{
			isatme: isatme,
			usr:    usr, txt: txt,
		})
		l.m[grp] = msgs
		return
	}
	isprevusr := len(msgs)%2 != 0
	if (isprevusr && !isbot) || (!isprevusr && isbot) { // is same
		_ = msgs[len(msgs)-1].add(item{
			isatme: isatme,
			usr:    usr, txt: txt,
		})
		return
	}
	if len(msgs) < cap(msgs) {
		msgs = append(msgs, l.newbatch(l.cap).add(item{
			isatme: isatme,
			usr:    usr, txt: txt,
		}))
		l.m[grp] = msgs
		return
	}
	copy(msgs, msgs[2:])
	msgs[len(msgs)-2] = l.newbatch(l.cap).add(item{
		isatme: isatme,
		usr:    usr, txt: txt,
	})
	l.m[grp] = msgs[:len(msgs)-1]
}

func (l *Log) Modelize(p model.Protocol, grp int64, sysp string, isusersystem bool) deepinfra.Model {
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
func Modelize[T any](l *Log, grp int64, f func(int, string) T) []T {
	l.mu.RLock()
	defer l.mu.RUnlock()
	sz := len(l.m[grp])
	if sz == 0 {
		return []T{f(0, l.defaultprompt)}
	}
	t := make([]T, sz)
	for i, msg := range l.m[grp] {
		t[i] = f(i, msg.String())
	}
	return t
}
