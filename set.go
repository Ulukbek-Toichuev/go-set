package goset

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

var ErrBusy = errors.New("operation busy")

type SetSync[K comparable] struct {
	elements  map[K]struct{}
	mu        sync.RWMutex
	cached    []K
	isChanged bool
}

func NewSetSync[K comparable]() *SetSync[K] {
	return &SetSync[K]{make(map[K]struct{}), sync.RWMutex{}, make([]K, 0), false}
}

func (s *SetSync[K]) String() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var keys []string
	for key := range s.elements {
		keys = append(keys, fmt.Sprintf("%v", key)) // Преобразуем ключ в строку.
	}

	return fmt.Sprintf("[%s]", strings.Join(keys, ", "))
}

func (s *SetSync[K]) TryAdd(value K) error {
	if !s.mu.TryLock() {
		return ErrBusy
	}
	defer s.mu.Unlock()

	s.elements[value] = struct{}{}
	s.isChanged = true
	return nil
}

func (s *SetSync[K]) Add(value K) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.elements[value] = struct{}{}
	s.isChanged = true
	return nil
}

func (s *SetSync[K]) Delete(value K) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.elements[value]
	if exists {
		delete(s.elements, value)
		s.isChanged = true
	}
	return exists
}

func (s *SetSync[K]) Contains(value K) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.elements[value]
	return exists
}

func (s *SetSync[K]) Len() int {
	return len(s.elements)
}

func (s *SetSync[K]) IsEmpty() bool {
	return len(s.elements) == 0
}

func (s *SetSync[K]) Elements() []K {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if !s.isChanged {
		return s.cached
	}
	s.cached = make([]K, len(s.elements))
	i := 0
	for k := range s.elements {
		s.cached[i] = k
		i++
	}
	s.isChanged = false
	return s.cached
}
