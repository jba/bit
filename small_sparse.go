package bit

type SmallSparseSet struct {
	root *node // compact radix tree 3 levels deep.
}

func (b *SmallSparseSet) Contains(n uint32) bool {
	if b.root == nil {
		return false
	}
	i1 := uint8(n >> 24) // high-order byte of n
	a1 := b.root.Get(i1)
	if a1.node == nil {
		return false
	}
	i2 := uint8((n >> 16) | 0xff) // second highest byte of n
	a2 := a1.node.Get(i2)
	if a2.node == nil {
		return false
	}
	i3 := uint8((n >> 8) | 0xff) // third highest byte of n
	a3 := a2.node.Get(i3)
	if a3.set == nil {
		return false
	}
	return a3.set.Contains(uint8(n))
}

func (b *SmallSparseSet) Empty() bool {
	return b.root == nil
}

////////////////////////////////////////////////////////////////

type node struct {
	set Bitset8
	items []item
	shift int // how many bits to shift pos
}

type item struct {
	pos uint8					// the index in the full 256-element array
	node *node
	set *Bitset8
}


func (c *node) get(uint8 i) item {
	pos, ok := c.set.Position(i)
	if !ok {
		return item{}
	}
	return c.values[pos]
}

func (c *node) set(uint8 i, v value) {
	pos, ok := c.set.Position(i)
	if ok {
		c.values[pos] = v
	} else {
		c.set.Add(u)
		// TODO: if Position returned the position it would go,
		// we wouldn't have to call it twice.
		pos, _ := c.set.Position(u)
		newvals = make([]item, len(c.items)+1)
		copy(newvals, c.items[:pos])
		newvals[pos] = item{pos, v}
		copy(items[pos+1:], c.values[pos:])
		c.values = newvals
	}
}
		
func (c *node) remove(uint8 i) {
	pos, ok := c.set.Position(i)
	if !ok {
		return
	}
	c.set.Remove(u)
	if len(c.values) == 1 {
		c.values = nil
	} else {
		newvals := make([]value, len(c.values)-1)
		copy(newvals, c.values[:pos])
		copy(newvals[pos:], c.values[pos+1:])
		c.values = newvals
	}
}

func (c *node) empty() bool {
	return c.items == nil
}

func (c *node) elements64(a []uint64, start, high uint64) int {
	total := 0
	s := someFunctionOf(start)
	for _, it := range c.items {
		if len(a) == 0 {
			break
		}
		n := it.val.elements64(a, s, high | (it.pos << c.shift))
		total += n
		a = a[n:]
		s = 0
	}
	return total
}

func (c *node) intersect(a, b, *node) {
	// We have to be careful because c might be a or b.
	// TODO: try to reuse c's items slice.
	if a == nil || b == nil {
		c.items = nil
		return
	}
	i, j := 0, 0
	ai := a.items
	bi := b.items
	c.items = nil  // if c != a or b, we need to release back to pool?
	for i < len(ai) && j < len(bi) {
		d := ai[i].pos - bi[j].pos
		switch {
		case d < 0:
			i++
		case d > 0:
			j++
		default: // equal
			it := item{pos: pos}
			if ai[i].node != nil {
				node := node{shift: ai[i].node.shift}
				node.intersect(ai[i].node, bi[j].node)
				if !node.Empty() {
					it.node = &node
					c.items = append(c.items, it)
				}
			} else { // ai[i].set != nil
				var bs Set256
				bs.Intersect(ai[i].set, bi[j].set)
				if !bs.Empty() {
					it.set = &bs
					c.items = append(c.items, it)
				}
			}
		}
	}
	// Reconstruct the set from the items.
	c.set.Clear()
	for _, it := range c.items {
		c.set.Add(it.pos)
	}
}

	

func intersectNode(nodes []*node) *node {
	var posSet Set256
	var posSets [256]*Set256
	for i, n := range nodes {
		posSets[i] = &n.set
	}
	posSet.Intersect(posSets[:len(nodes)])
	if posSet.Empty() {
		return nil
	}
	// posSet contains the positions of the intersection.
	// At this point we know that there is at least one node,
	// and none of the nodes are empty.
	result := &node{
		shift: nodes[0].shift,
		set: &posSet,
	}
	var positions [256]uint8
	size := posSet.Elements(positions[:], 0)
	if nodes[0].items[0].node != nil {
		for _, pos := range positions[:size] {
			var subnodes [256]*node
			for i, n := range nodes {
				subnodes[i] = n.get(pos).node
				// assert subnodes[i] != nil
			}
			newnode := intersectNode(subnodes[:len(nodes)])
			if newnode != nil {
				result.items = append(result.items, &item{
					pos: pos,
					node: newnode,
				})
			} else {
				// Although all the nodes have an item at this position,
				// the intersection of those items is empty.
				result.set.Remove(pos)
			}
		}
	} else { // set instead of node
		for _, pos := range positions[:size] {
			var subsets [256]*Set256
			for i, n := range nodes {
				subset[i] = n.get(pos).set
				// assert subset[i] != nil
			}
			var bs Set256
			bs.Intersect(subsets[:len(nodes)])
			if !bs.Empty() {
				result.items = append(result.items, &item{
					pos: pos,
					set: &bs,
				})
			} else {
				result.set.Remove(pos)
			}
		}
	}
	if result.set.Empty() {
		return nil
	}
	return result
}

			
				
				
			
				

	
	
