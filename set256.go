package bit

// A Set256 represents a set of values in the range [0, 256)
// more efficiently than Set of capacity 256.
type Set256 struct {
	sets [4]Set64
}

func (b *Set256) Add(n uint8) {
	b.sets[n/64].Add(n % 64)
}

func (b *Set256) Remove(n uint8) {
	b.sets[n/64].Remove(n % 64)
}

func (b *Set256) Contains(n uint8) bool {
	return b.sets[n/64].Contains(n % 64)
}

func (b *Set256) Empty() bool {
	return b.sets[0].Empty() && b.sets[1].Empty() && b.sets[2].Empty() && b.sets[3].Empty()
}

func (b *Set256) Clear() {
	b.sets[0].Clear()
	b.sets[1].Clear()
	b.sets[2].Clear()
	b.sets[3].Clear()
}

func (b *Set256) Size() int {
	return b.sets[0].Size() + b.sets[1].Size() + b.sets[2].Size() + b.sets[3].Size()
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
func (c *Set256) Intersect2(a, b *Set256) {
	c.sets[0] = a.sets[0] & b.sets[0]
	c.sets[1] = a.sets[1] & b.sets[1]
	c.sets[2] = a.sets[2] & b.sets[2]
	c.sets[3] = a.sets[3] & b.sets[3]
}

// c cannot be one of sets
func (c *Set256) IntersectN(bs []*Set256) {
	if len(bs) == 0 {
		c.Clear()
		return
	}
	for i := 0; i < 4; i++ {
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
	n := elements(s.sets[si], a, start %64, si * 64)
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



// func (c *Set256) Elements64(a []uint64, start uint8, high uint64) int {
// 	if len(a) == 0 {
// 		return 0
// 	}
// 	i := start >> 6
// 	total := elements64(b.sets[i], a, start&0x3f, high|(i<<6))
// 	// TODO: compare perf with unrolling the loop.
// 	for ; i < 4; i++ {
// 		total += elements64(b.sets[i], a[total:], 0, high|(i<<6))
// 	}
// 	return total
// }

// func elements64(u uint64, a []uint64, startBit int, high uint64) int {
// 	i := 0
// 	// TODO: compare performance with making mask a loop variable.
// 	for b := startBit; b < 64 && i < len(a); b++ {
// 		if u & (1 << b) {
// 			a[i] = high | b
// 			i++
// 		}
// 	}
// 	return i
// }

// a 1 in the position indicated by the low 6 bits of n.
// TODO: confirm this is inlined.
func onebit(n uint8) uint64 {
	return 1 << (n & 0x3f)
}
