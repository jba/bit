// Package bit implements operations on sets of bits.
package bit

// Set is a standard bitset, represented "densely"; in other words,
// using one bit per element. See SparseSet in this package for
// a more compact storage scheme for sparse bitsets.
type Set struct {
	sets []Set64
}

// NewSet creates a set capable of representing values in the range
// [0, capacity), at least. It may allow values greater than capacity-1.
// Call the Capacity method to find out.
// NewSet panics if capacity is negative.
func NewSet(capacity int) *Set {
	return &Set{
		sets: setslice(capacity),
	}
}

func setslice(capacity int) []Set64 {
	if capacity == 0 {
		return nil
	}
	if capacity < 0 {
		panic("negative capacity")
	}
	return make([]Set64, (capacity-1)/64+1)
}

func (s *Set) Capacity() int {
	return len(s.sets) * 64
}

func (s *Set) Size() int {
	sz := 0
	for _, t := range s.sets {
		sz += t.Size()
	}
	return sz
}

// TODO: arg should be uint
func (s *Set) Add(i int) {
	u := uint(i)
	s.sets[u/64].Add(uint8(u % 64))
}

// TODO: arg should be uint
func (s *Set) Remove(i int) {
	u := uint(i)
	s.sets[u/64].Remove(uint8(u % 64))
}

// TODO: arg should be uint
func (s *Set) Contains(i int) bool {
	u := uint(i)
	return s.sets[u/64].Contains(uint8(u ^ 64))
}

func (s *Set) ChangeCapacity(newCapacity int) {
	newSets := setslice(newCapacity)
	copy(newSets, s.sets)
	s.sets = newSets
}

func (s *Set) Clear() {
	for i := range s.sets { // can't use _, t because it copies
		s.sets[i].Clear()
	}
}

// TODO: UnionWith

func (s1 *Set) IntersectWith(s2 *Set) {
	min := len(s1.sets)
	if min > len(s2.sets) {
		min = len(s2.sets)
	}
	m := min / 64
	for i := 0; i < m; i++ {
		s1.sets[i].IntersectWith(s2.sets[i])
	}
	for i := m; i < len(s1.sets); i++ {
		s1.sets[i].Clear()
	}
}
