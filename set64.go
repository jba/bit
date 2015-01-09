package bit

// A Set64 is an efficient representation of a bitset that can represent [0, 64).
// For efficiency, the methods of Set64 perform no bounds checking on their arguments.
type Set64 uint64

func (s *Set64) Add(u uint8) {
	*s |= 1 << u
}

func (s *Set64) Remove(u uint8) {
	*s &= ^(1 << u)
}

func (s *Set64) Clear() {
	*s = 0
}

func (s1 *Set64) IntersectWith(s2 Set64) {
	*s1 &= s2
}

func (s1 *Set64) UnionWith(s2 Set64) {
	*s1 |= s2
}

//func (s Set64) Size() int {
	
	
