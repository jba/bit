package bit

// A node is a compact radix tree element.
// It behaves like a 256-element array of subnodes, indexed by
// one byte of the element. In fact, only the non-empty
// subnodes are represented; the bitset stores this set
// and subnodes slice contains the non-empty subnodes in order.
type node struct {
	shift    uint // how many bits to shift elements
	bitset   Set256
	subnodes []subnode // if shift > 0
}

type subnode struct {
	index uint8 // the index in the full 256-element array
	sub   subber
}

type subber interface {
	add(uint64)
	remove(uint64) bool // return true if empty
	contains(uint64) bool
	elements(a []uint64, start, high uint64) int
	size() int
	memSize() uint64
	equalSub(subber) bool
}

func (n *node) newSubber() subber {
	if n.shift == 8 {
		return &Set256{}
	} else {
		return &node{shift: n.shift - 8}
	}
}

func (n *node) add(e uint64) {
	index := uint8(e >> n.shift)
	pos, found := n.bitset.Position(index)
	if !found {
		n.bitset.Add(index)
	}
	var sub subber
	if found {
		sub = n.subnodes[pos].sub
	} else {
		sub = n.newSubber()
		newsub := make([]subnode, len(n.subnodes)+1)
		copy(newsub, n.subnodes[:pos])
		newsub[pos] = subnode{index: index, sub: sub}
		copy(newsub[pos+1:], n.subnodes[pos:])
		n.subnodes = newsub
	}
	sub.add(e)
}

func (n *node) remove(e uint64) (empty bool) {
	// assert node is not empty
	index := uint8(e >> n.shift)
	pos, found := n.bitset.Position(index)
	if !found {
		return false // we weren't empty coming in
	}
	sub := n.subnodes[pos].sub
	if sub.remove(e) {
		if len(n.subnodes) == 1 {
			// No need to clean up, we're finished.
			return true
		}
		copy(n.subnodes[pos:], n.subnodes[pos+1:])
		// TODO: really shrink memory
		n.subnodes = n.subnodes[:len(n.subnodes)-1]
		n.bitset.Remove(index)
	}
	return false
}

func (n *node) contains(e uint64) bool {
	index := uint8(e >> n.shift)
	p, found := n.bitset.Position(index)
	if !found {
		return false
	}
	return n.subnodes[p].sub.contains(e)
}

func (n1 *node) equal(n2 *node) bool {
	if !n1.bitset.Equal(&n2.bitset) {
		return false
	}
	for i, sn1 := range n1.subnodes {
		if !sn1.sub.equalSub(n2.subnodes[i].sub) {
			return false
		}
	}
	return true
}

func (n1 *node) equalSub(s subber) bool {
	return n1.equal(s.(*node))
}


func (n *node) size() int {
	t := 0
	for _, s := range n.subnodes {
		t += s.sub.size()
	}
	return t
}

func (n *node) memSize() uint64 {
	sz := memSize(*n)
	for _, s := range n.subnodes {
		sz += memSize(s)
		sz += s.sub.memSize()
	}
	return sz
}

func (n *node) elements(a []uint64, start, high uint64) int {
	hi := func(i int) uint64 {
		return high | (uint64(n.subnodes[i].index) << n.shift)
	}

	var total int
	si := uint8(start >> n.shift)
	p, found := n.bitset.Position(si)
	if found {
		total = n.subnodes[p].sub.elements(a, start, hi(p))
		p++
	}
	for i := p; i < len(n.subnodes); i++ {
		total += n.subnodes[i].sub.elements(a[total:], 0, hi(i))
	}
	return total
}

// func (c *node) intersect(a, b, *node) {
// 	// We have to be careful because c might be a or b.
// 	// TODO: try to reuse c's items slice.
// 	if a == nil || b == nil {
// 		c.items = nil
// 		return
// 	}
// 	i, j := 0, 0
// 	ai := a.items
// 	bi := b.items
// 	c.items = nil  // if c != a or b, we need to release back to pool?
// 	for i < len(ai) && j < len(bi) {
// 		d := ai[i].pos - bi[j].pos
// 		switch {
// 		case d < 0:
// 			i++
// 		case d > 0:
// 			j++
// 		default: // equal
// 			it := item{pos: pos}
// 			if ai[i].node != nil {
// 				node := node{shift: ai[i].node.shift}
// 				node.intersect(ai[i].node, bi[j].node)
// 				if !node.Empty() {
// 					it.node = &node
// 					c.items = append(c.items, it)
// 				}
// 			} else { // ai[i].set != nil
// 				var bs Set256
// 				bs.Intersect(ai[i].set, bi[j].set)
// 				if !bs.Empty() {
// 					it.set = &bs
// 					c.items = append(c.items, it)
// 				}
// 			}
// 		}
// 	}
// 	// Reconstruct the set from the items.
// 	c.set.Clear()
// 	for _, it := range c.items {
// 		c.set.Add(it.pos)
// 	}
// }

func intersectNodes(nodes []*node) *node {
	var bsets [256]*Set256
	for i, n := range nodes {
		bsets[i] = &n.bitset
	}
	var bset Set256
	bset.IntersectN(bsets[:len(nodes)])
	if bset.Empty() {
		return nil
	}
	// posSet contains the indices of the intersection.
	// At this point we know that there is at least one node,
	// and none of the nodes are empty.
	result := &node{
		shift: nodes[0].shift,
		bitset: bset,
	}
	var indices [256]uint8
	size := bset.Elements(indices[:], 0)
	var subnodes [256]*node
	var subsets [256]*Set256
	isSets := (nodes[0].shift == 8)
	for _, index := range indices[:size] {
		for i, n := range nodes {
			p, found := n.bitset.Position(index)
			if !found {
				panic("intersectNodes: index not found")
			}
			sub := n.subnodes[p].sub
			if isSets {
				subsets[i] = sub.(*Set256)
			} else {
				subnodes[i] = sub.(*node)
			}
		}
		var newsub subber
		if isSets {
			var bs Set256
			bs.IntersectN(subsets[:len(nodes)])
			if !bs.Empty() {
				newsub = &bs
			}
		} else {
			in := intersectNodes(subnodes[:len(nodes)])
			if in != nil {
				newsub = in
			}
		}
		if newsub != nil {
			result.subnodes = append(result.subnodes,
				subnode{index: index, sub: newsub})
		} else {
			// Although all the nodes have an item at this position,
			// the intersection of those items is empty.
			result.bitset.Remove(index)
		}
	}
	if result.bitset.Empty() {
		return nil
	}
	return result
}

