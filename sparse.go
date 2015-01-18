//TODO: use sync.Pool?
package bit

import (
	"bytes"
	"fmt"
	"reflect"
)

type SparseSet struct {
	root *node // compact radix tree 7 levels deep.
}

func NewSparseSet(els ...uint64) *SparseSet {
	s := &SparseSet{}
	for _, e := range els {
		s.Add(e)
	}
	return s
}

func (s *SparseSet) Add(n uint64) {
	if s.root == nil {
		s.root = &node{shift: 64 - 8}
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

func (s *SparseSet) Clear() {
	s.root = nil
}

func (s1 *SparseSet) Equal(s2 *SparseSet) bool {
	if s1.root == nil || s2.root == nil {
		return s1.root == s2.root
	}
	return s1.root.equal(s2.root)
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

func (s *SparseSet) Elements(a []uint64, start uint64) int {
	if s.root == nil {
		return 0
	}
	return s.root.elements(a, start, 0)
}

// s becomes the intersection of the ss. It must not be
// one of the ss, and it is not part of the intersection.
func (s *SparseSet) Intersect(ss ...*SparseSet) {
	s.Clear()
	var nodes []*node
	for _, t := range ss {
		if t.Empty() {
			return
		}
		nodes = append(nodes, t.root)
	}
	s.root = intersectNodes(nodes)
}

func (s SparseSet) String() string {
	if s.Empty() {
		return "{}"
	}
	els := make([]uint64, s.Size())
	s.Elements(els, 0)
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "{%d", els[0])
	for _, e := range els[1:] {
		fmt.Fprintf(&buf, ", %d", e)
	}
	fmt.Fprint(&buf, "}")
	return buf.String()
}


	
		
