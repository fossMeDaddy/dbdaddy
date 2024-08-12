package types

import "sync"

type SafeVar[T any] struct {
	val T
	mx  sync.Mutex
}

func (s *SafeVar[T]) Get() T {
	return s.val
}

func (s *SafeVar[T]) Set(newVal T) {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.val = newVal
}

func NewSafeVar[T any](val T) SafeVar[T] {
	return SafeVar[T]{
		val: val,
		mx:  sync.Mutex{},
	}
}
