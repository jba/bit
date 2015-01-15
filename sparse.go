//TODO: use sync.Pool?
package bit

import "reflect"

type SparseSet struct {
	root *node // compact radix tree 7 levels deep.
}

func (s *SparseSet) Add(n uint64) {
	if s.root == nil {
		s.root = &node{shift: 64-8}
	}
	s.root.add(n)
}

func (s *SparseSet) Remove(n uint64) {
	if s.root == nil {
		return
	}
	if s.root.remove(uint64(n)) {
		s.root = nil
	}
}

func (s *SparseSet) Contains(n uint64) bool {
	if s.root == nil {
		return false
	}
	return s.root.contains(n)
}

func (s *SparseSet) Empty() bool {
	return s.root == nil
}

func (s *SparseSet) Size() int {
	if s.root == nil {
		return 0
	}
	return s.root.size()
}

func (s *SparseSet) MemSize() uint64 {
	sz := memSize(*s)
	if s.root != nil {
		sz += s.root.memSize()
	}
	return sz
}

func memSize(x interface{}) uint64 {
	return uint64(reflect.TypeOf(x).Size())
}
