package analyzer

import (
	"container/list"
	"sync"
)

type Store struct {
	mux  *sync.RWMutex
	data *list.List
	max  int
}

// NewStore return a new store
func NewStore(max int) *Store {
	data := list.New()
	return &Store{
		max:  max,
		data: data,
	}
}

//GetSome return some data
func (s *Store) GetSome(length int) []SystemStatus {
	s.mux.RLock()
	defer s.mux.RUnlock()
	var status []SystemStatus
	var counter int
	for t := s.data.Back(); t != nil && counter < length; t = t.Prev() {
		if s, ok := t.Value.(SystemStatus); ok == true {
			status = append(status, s)
		}
		counter++
	}
	return status
}

//GetLastOne return only one last data
func (s *Store) GetLastOne() (status SystemStatus) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	value := s.data.Back()
	if status, ok := value.Value.(SystemStatus); ok == true {
		return status
	}
	return
}

// Clear remove all the data
func (s *Store) Clear() {
	s.mux.Lock()
	defer s.mux.Unlock()
	for t := s.data.Front(); t != nil; t = t.Next() {
		s.data.Remove(t)
	}
}

// Put append a new data to list
func (s *Store) Put(status *SystemStatus) bool {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.data.Len() == s.max {
		s.data.Remove(s.data.Front())
		s.data.PushBack(status)
		return true
	} else if s.data.Len() < s.max {
		s.data.PushBack(status)
		return true
	}
	return false
}
