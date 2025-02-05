package queue

import "sync"

type Queue struct {
	lock  sync.RWMutex
	items []interface{}
}

func (q *Queue) Enqueue(item interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.items = append(q.items, item)
}

func (q *Queue) Dequeue() (interface{}, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(q.items) == 0 {
		return nil, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

func (q *Queue) Front() (interface{}, bool) {
	q.lock.RLock()
	defer q.lock.RUnlock()
	if len(q.items) == 0 {
		return nil, false
	}
	return q.items[0], true
}

func (q *Queue) IsEmpty() bool {
	q.lock.RLock()
	defer q.lock.RUnlock()
	return len(q.items) == 0
}

func (q *Queue) Size() int {
	q.lock.RLock()
	defer q.lock.RUnlock()
	return len(q.items)
}
