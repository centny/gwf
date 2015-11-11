package pool

import (
	"sync"
)

type Queue struct {
	Lock bool
	lock sync.RWMutex
	objs []interface{}
}

func NewQueue(lock bool) *Queue {
	return &Queue{
		Lock: lock,
		lock: sync.RWMutex{},
	}
}

func (q *Queue) Push(objs ...interface{}) {
	if q.Lock {
		q.lock.Lock()
		defer q.lock.Unlock()
	}
	q.objs = append(q.objs, objs...)
}
func (q *Queue) Poll() interface{} {
	if q.Lock {
		q.lock.Lock()
		defer q.lock.Unlock()
	}
	if len(q.objs) < 1 {
		return nil
	}
	v := q.objs[0]
	q.objs = q.objs[1:]
	return v
}
func (q *Queue) Pollv() interface{} {
	v := q.Poll()
	if v == nil {
		panic("not more object on quque")
	}
	return v
}
func (q *Queue) Len() int {
	return len(q.objs)
}
