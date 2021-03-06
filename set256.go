package bit

import (
	"bytes"
	"fmt"
)

// A Set256 represents a set of integers in the range [0, 256).
// It does so more efficiently than a Set of capacity 256.
// For efficiency, the methods of Set256 perform no bounds checking on their
// arguments.
type Set256 struct {
	sets [4]Set64
}

func (s *Set256) Add(n uint8) {
	s.sets[n/64].Add(n % 64)
}

func (s *Set256) Remove(n uint8) {
	s.sets[n/64].Remove(n % 64)
}

func (s *Set256) Contains(n uint8) bool {
	return s.sets[n/64].Contains(n % 64)
}

func (s *Set256) Empty() bool {
	return s.sets[0].Empty() && s.sets[1].Empty() && s.sets[2].Empty() && s.sets[3].Empty()
}

func (s *Set256) Clear() {
	s.sets[0].Clear()
	s.sets[1].Clear()
	s.sets[2].Clear()
	s.sets[3].Clear()
}

func (s *Set256) Size() int {
	return s.sets[0].Size() + s.sets[1].Size() + s.sets[2].Size() + s.sets[3].Size()
}

func (Set256) Capacity() int {
	return 256
}

func (s1 *Set256) Equal(s2 *Set256) bool {
	return s1.sets[0] == s2.sets[0] &&
		s1.sets[1] == s2.sets[1] &&
		s1.sets[2] == s2.sets[2] &&
		s1.sets[3] == s2.sets[3]
}

// Position returns the 0-based position of n in the set. If
// the set is {3, 8, 15}, then the position of 8 is 1.
// If n is not in the set, returns 0, false.
// If not a member, return where it would go.
// The second return value reports whether n is a member of b.
func (b *Set256) Position(n uint8) (int, bool) {
	var pos int
	i := n / 64
	switch i {
	case 1:
		pos = b.sets[0].Size()
	case 2:
		pos = b.sets[0].Size() + b.sets[1].Size()
	case 3:
		pos = b.sets[0].Size() + b.sets[1].Size() + b.sets[2].Size()
	}
	p, ok := b.sets[i].Position(n % 64)
	return pos + p, ok
}

// c = a intersect b
// func (c *Set256) Intersect2(a, b *Set256) {
// 	c.sets[0] = a.sets[0] & b.sets[0]
// 	c.sets[1] = a.sets[1] & b.sets[1]
// 	c.sets[2] = a.sets[2] & b.sets[2]
// 	c.sets[3] = a.sets[3] & b.sets[3]
// }

// c cannot be one of sets
func (c *Set256) IntersectN(bs []*Set256) {
	if len(bs) == 0 {
		c.Clear()
		return
	}
	for i := 0; i < len(c.sets); i++ {
		c.sets[i] = bs[0].sets[i]
		for _, s := range bs[1:] {
			c.sets[i].IntersectWith(s.sets[i])
		}
	}
}

// Fill a with set elements, starting from start.
// Return the number added.
func (s *Set256) Elements(a []uint8, start uint8) int {
	if len(a) == 0 {
		return 0
	}
	si := start / 64
	n := elements(s.sets[si], a, start%64, si*64)
	for i := si + 1; i < 4; i++ {
		n += elements(s.sets[i], a[n:], 0, i*64)
	}
	return n
}

func elements(s Set64, a []uint8, start, high uint8) int {
	n := s.Elements(a, start)
	for i := 0; i < n; i++ {
		a[i] |= high
	}
	return n
}

func (s *Set256) Elements64(a []uint64, start uint8, high uint64) int {
	if len(a) == 0 {
		return 0
	}
	si := start / 64
	n := s.sets[si].Elements64(a, start%64, high|uint64(si*64))
	for i := si + 1; i < 4; i++ {
		n += s.sets[i].Elements64(a[n:], 0, high|uint64(i*64))
	}
	return n
}

func (s Set256) String() string {
	var a [256]uint64
	n := s.Elements64(a[:], 0, 0)
	if n == 0 {
		return "{}"
	}
	// TODO: use strings.Builder
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "{%d", a[0])
	for _, e := range a[1:n] {
		fmt.Fprintf(&buf, ", %d", e)
	}
	fmt.Fprint(&buf, "}")
	return buf.String()
}

// For subber, used in node:

func (s *Set256) add(e uint64) { s.Add(uint8(e)) }

func (s *Set256) remove(e uint64) bool {
	s.Remove(uint8(e))
	return s.Empty()
}

func (s *Set256) contains(e uint64) bool {
	return s.Contains(uint8(e))
}

func (s *Set256) size() int { return s.Size() }

func (s *Set256) memSize() uint64 { return memSize(*s) }

func (s *Set256) elements(a []uint64, start, high uint64) int {
	return s.Elements64(a, uint8(start), high)
}

func (s *Set256) equalSub(b subber) bool {
	return s.Equal(b.(*Set256))
}
