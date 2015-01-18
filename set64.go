package bit

import (
	"bytes"
	"fmt"
)

// A Set64 is an efficient representation of a bitset that can represent [0, 64).
// For efficiency, the methods of Set64 perform no bounds checking on their arguments.
type Set64 uint64

func (s *Set64) Add(u uint8) {
	*s |= 1 << u
}

func (s *Set64) Remove(u uint8) {
	*s &^= (1 << u)
}

func (s *Set64) Contains(u uint8) bool {
	return (*s&(1<<u) != 0)
}

func (s Set64) Empty() bool {
	return s == 0
}

func (s Set64) Size() int {
	return CountOnes64(uint64(s))
}

func (Set64) Capacity() int {
	return 64
}

func (s *Set64) Clear() {
	*s = 0
}

// Position returns the 0-based position of n in the set. If the set
// is {3, 8, 15}, then the position of 8 is 1.  If n is not in the
// set, Position returns the position n would be at if it were a
// member. The second return value reports whether n is a member of
// s.
func (s Set64) Position(n uint8) (int, bool) {
	mask := uint64(1 << n)
	in := (uint64(s)&mask != 0)
	pos := CountOnes64(uint64(s) & (mask - 1))
	return pos, in
}

func (s1 *Set64) IntersectWith(s2 Set64) {
	*s1 &= s2
}

func (s1 *Set64) UnionWith(s2 Set64) {
	*s1 |= s2
}

func (s Set64) Elements(a []uint8, start uint8) int {
	if len(a) == 0 {
		return 0
	}
	i := 0
	for b := start; b < 64 && i < len(a); b++ {
		if s.Contains(b) {
			a[i] = b
			i++
		}
	}
	return i
}

func (s Set64) Elements64(a []uint64, start uint8, high uint64) int {
	if len(a) == 0 {
		return 0
	}
	i := 0
	for b := start; b < 64 && i < len(a); b++ {
		if s.Contains(b) {
			a[i] = high | uint64(b)
			i++
		}
	}
	return i
}

func (s Set64) String() string {
	var a [64]uint8
	n := s.Elements(a[:], 0)
	if n == 0 {
		return "{}"
	}
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "{%d", a[0])
	for _, e := range a[1:n] {
		fmt.Fprintf(&buf, ", %d", e)
	}
	fmt.Fprint(&buf, "}")
	return buf.String()
}
