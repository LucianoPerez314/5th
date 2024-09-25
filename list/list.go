package list

import (
	"cmp"
	"sync"
	"sync/atomic"
)

type node[K cmp.Ordered, V any] struct {
	sync.Mutex
	key    K
	item   V
	marked atomic.Bool
	next   atomic.Pointer[node[K, V]]
}

type List[K cmp.Ordered, V any] struct {
	head *node[K, V]
}

func New[K cmp.Ordered, V any](minKey K, maxKey K) List[K, V] {
	tail := new(node[K, V])
	tail.key = maxKey
	head := new(node[K, V])
	head.key = minKey
	head.next.Store(tail)
	list := List[K, V]{head}
	return list
}

func (l List[K, V]) Find(key K) (V, bool) {
	curr := l.head
	for curr.key < key {
		curr = curr.next.Load()
	}
	return curr.item, curr.key == key
}

func (l List[K, V]) Insert(key K, item V) bool {
	for {
		pred := l.head
		curr := pred.next.Load()
		for curr.key < key {
			pred = curr
			curr = curr.next.Load()
		}
		pred.Lock()
		curr.Lock()
		if !l.validate(pred, curr) {
			curr.Unlock()
			pred.Unlock()
			continue
		}
		result := false
		if key != curr.key {
			var newNode node[K, V]
			newNode.key = key
			newNode.item = item
			newNode.next.Store(curr)
			pred.next.Store(&newNode)
			result = true
		}

		curr.Unlock()
		pred.Unlock()
		return result
	}
}

func (l List[K, V]) Remove(key K) (V, bool) {
	for {
		pred := l.head
		curr := pred.next.Load()
		for curr.key < key {
			pred = curr
			curr = curr.next.Load()
		}

		pred.Lock()
		curr.Lock()

		if !l.validate(pred, curr) {
			curr.Unlock()
			pred.Unlock()
			continue
		}

		result := false
		if key == curr.key {
			curr.marked.Store(true)
			next := curr.next.Load()
			pred.next.Store(next)
			result = true
		}

		curr.Unlock()
		pred.Unlock()
		return curr.item, result
	}
}

func (l List[K, V]) validate(pred, curr *node[K, V]) bool {
	return !pred.marked.Load() && !curr.marked.Load() && pred.next.Load() == curr
}
