package bit

//import "fmt"

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
	sub subber
}

type subber interface {
	add(uint64)
	remove(uint64) bool // return true if empty
	contains(uint64) bool
	elements(a []uint64, start, high uint64) int
	size() int
	memSize() uint64
}	

func (n *node) newSubber() subber {
	if n.shift == 8 {
		return &Set256{}
	} else {
		return &node{shift: n.shift-8}
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
		//fmt.Printf("node shift %d: grew to %d\n", n.shift, len(n.subnodes))
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
	var total int
	si := uint8(start >> n.shift)
	p, found := n.bitset.Position(si)
	if found {
		total = n.subnodes[p].sub.elements(a, start, uint64(n.subnodes[p].index) << n.shift)
	}
	for i := p+1; i < len(n.subnodes); i++ {
		total += n.subnodes[i].sub.elements(a[total:], 0, uint64(n.subnodes[i].index) << n.shift)
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

// func intersectNode(nodes []*node) *node {
// 	var posSet Set256
// 	var posSets [256]*Set256
// 	for i, n := range nodes {
// 		posSets[i] = &n.set
// 	}
// 	posSet.Intersect(posSets[:len(nodes)])
// 	if posSet.Empty() {
// 		return nil
// 	}
// 	// posSet contains the positions of the intersection.
// 	// At this point we know that there is at least one node,
// 	// and none of the nodes are empty.
// 	result := &node{
// 		shift: nodes[0].shift,
// 		set: &posSet,
// 	}
// 	var positions [256]uint8
// 	size := posSet.Elements(positions[:], 0)
// 	if nodes[0].items[0].node != nil {
// 		for _, pos := range positions[:size] {
// 			var subnodes [256]*node
// 			for i, n := range nodes {
// 				subnodes[i] = n.get(pos).node
// 				// assert subnodes[i] != nil
// 			}
// 			newnode := intersectNode(subnodes[:len(nodes)])
// 			if newnode != nil {
// 				result.items = append(result.items, &item{
// 					pos: pos,
// 					node: newnode,
// 				})
// 			} else {
// 				// Although all the nodes have an item at this position,
// 				// the intersection of those items is empty.
// 				result.set.Remove(pos)
// 			}
// 		}
// 	} else { // set instead of node
// 		for _, pos := range positions[:size] {
// 			var subsets [256]*Set256
// 			for i, n := range nodes {
// 				subset[i] = n.get(pos).set
// 				// assert subset[i] != nil
// 			}
// 			var bs Set256
// 			bs.Intersect(subsets[:len(nodes)])
// 			if !bs.Empty() {
// 				result.items = append(result.items, &item{
// 					pos: pos,
// 					set: &bs,
// 				})
// 			} else {
// 				result.set.Remove(pos)
// 			}
// 		}
// 	}
// 	if result.set.Empty() {
// 		return nil
// 	}
// 	return result
// }
