//TODO: use sync.Pool?
package bit

/* A bitset for uint64s that are expected to be assigned consecutively, as unique IDs from 
  a counter.  The items they represent get deleted over time, so we expect a relatively
  small range to be occupied at any given time.
*/
// The high-order 32 bits are unlikely to have more than a handful of values at any time,
// so we use a sorted sequence of bitsets for 32-bit values.
// type BitSet64 struct {
// 	items []item64
// }

// type item64 struct {
// 	value uint32 
// 	set *BitSet32
// }

// func (b *BitSet64) Add(n uint64) bool {
// 	h := n >> 32
// 	var set32 *BitSet32
// 	for i, it := range b.items {
// 		if it.value == h {
// 			set32 = it.set
// 			break
// 		} else if it.value > h {
// 			newItems := make([]item64, len(b.items) + 1)
// 			copy(newItems, b.items[:i])
// 			set32 = NewBitSet32()
// 			newItems[i] = item64{h, set32}
// 			newItems[i].set.Add(uint32(n))
// 			copy(newItems[i+1:], b.items[i:])
// 			b.items = newItems
// 			break
// 		}
// 	}
// 	if set32 == nil {
// 		set32 = NewBitSet32()
// 		b.items = append(b.items, item64{h, set32})
// 	}
// 	set32.Add(uint32(n))
// }

// func (b *BitSet64) Remove(n uint64) {
// }

// func (b *BitSet64) Contains(n uint64) bool {
// 	h := n >> 32
// 	for i, it := range b.items {
// 		if it.value == h {
// 			return it.set.Contains(uint32(n & low32mask))
// 		} else if it.value > h {
// 			return false
// 		}
// 	}
// 	return false
// }

// func (b *BitSet64) Empty() bool {
// 	return len(b.items) == 0
// }

// func (b *BitSet64) Size() int64 {
// 	var s int64
// 	for _, it := range b.items {
// 		s += it.set.Size()
// 	}
// 	return s
// }

// // Size of data structure in bytes.
// func (b *BitSet64) MemSize() int64 {
// 	var s int64 = sizeof(*b)
// 	// TODO: add size of each item?
// 	for _, it := range b.items {
// 		s += it.set.MemSize()
// 	}
// 	return s
// }


// func (b *BitSet64) Elements(a []uint64, start uint64) int {
		
	
//}
