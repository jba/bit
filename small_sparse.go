package bit

type SmallSparseSet struct {
	root *node // compact radix tree 3 levels deep.
}

func (s *SmallSparseSet) Add(n uint32) {
	if s.root == nil {
		s.root = &node{shift: 24}
	}
	s.root.add(uint64(n))
}

func (s *SmallSparseSet) Remove(n uint32) {
	if s.root.remove(uint64(n)) {
		s.root = nil
	}
}

func (s *SmallSparseSet) Contains(n uint32) bool {
	return s.root.contains(uint64(n))
}

func (s *SmallSparseSet) Empty() bool {
	return s.root == nil
}


