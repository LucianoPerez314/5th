package test

import (
	"cmp"
)

type ConcurrentList[K cmp.Ordered, V any] interface {
	Find(K) (V, bool)
	Insert(K, V) bool
	Remove(K) (V, bool)
}

func RunSequentialTests[K cmp.Ordered, V any](l ConcurrentList[K, V]) {

}

func RunConcurrentTests[K cmp.Ordered, V any](l ConcurrentList[K, V]) {

}
