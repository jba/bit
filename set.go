// Package bit implements operations on sets of bits.
package bit

// Set is a standard bitset, represented "densely"; in other words,
// using one bit per element. See the sparse sets in this package for
// more compact storage schemes.
type Set struct {
	set64s []Set64
}

// NewSet creates a set capable of representing values in the range
// [0, capacity-1), at least. It may allow values greater than capacity-1.
func NewSet(capacity int) *Set {
	return &Set{
		set64s: set64slice(capacity),
	}
}

func set64slice(capacity int) []Set64 {
	if capacity == 0 {
		return nil
	}
	if capacity < 0 {
		panic("negative capacity")
	}
	return make([]Set64, (capacity-1)/64+1)
}

func (s *Set) Capacity() int {
	return len(s.set64s)
}

func (s *Set) Add(i int) {
	u := uint(i)
	s.set64s[u/64].Add(uint8(u % 64))
}

func (s *Set) Remove(i int) {
	u := uint(i)
	s.set64s[u/64].Remove(uint8(u % 64))
}

func (s *Set) ChangeCapacity(newCapacity int) {
	newSets := set64slice(newCapacity)
	copy(newSets, s.set64s)
	s.set64s = newSets
}

func (s *Set) Clear() {
	for i := range s.set64s { // can't use _, t because it copies
		s.set64s[i].Clear()
	}
}

func (s1 *Set) IntersectWith(s2 *Set) {
	min := len(s1.set64s)
	if min > len(s2.set64s) {
		min = len(s2.set64s)
	}
	m := min / 64
	for i := 0; i < m; i++ {
		s1.set64s[i].IntersectWith(s2.set64s[i])
	}
	for i := m; i < len(s1.set64s); i++ {
		s1.set64s[i].Clear()
	}
}

