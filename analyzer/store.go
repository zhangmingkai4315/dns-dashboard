package analyzer

import (
	"container/list"
	"sync"
)

// Store 定义一个数据仓库，存储临时使用的数据
type Store struct {
	mux  *sync.RWMutex
	data *list.List
	max  int
}

// NewStore 创建数据仓库
func NewStore(max int) *Store {
	data := list.New()
	return &Store{
		max:  max,
		data: data,
	}
}

//GetSome 返回部分数据集合
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

//GetLastOne 仅仅返回一个数据集合
func (s *Store) GetLastOne() (status SystemStatus) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	value := s.data.Back()
	if status, ok := value.Value.(SystemStatus); ok == true {
		return status
	}
	return
}

// Clear 删除所有数据集合
func (s *Store) Clear() {
	s.mux.Lock()
	defer s.mux.Unlock()
	for t := s.data.Front(); t != nil; t = t.Next() {
		s.data.Remove(t)
	}
}

// Put 增加一条新的数据记录
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
