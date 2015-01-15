package bit

import (
	"runtime"
	"strconv"
	"testing"
	"math/rand"
	"fmt"
)

func TestSparseBasics(t *testing.T) {
	check := func(b bool) {
		if !b {
			_, _, line, ok := runtime.Caller(1)
			var ln string
			if !ok {
				ln = "???"
			} else {
				ln = strconv.Itoa(line)
			}
			t.Fatalf("line %s failed", ln)
		}
	}
 	var s SparseSet

 	check(s.Empty())
	s.Add(0)
	check(!s.Empty())
	check(s.Contains(0))
	check(!s.Contains(1))

	s.Add(492409)
	check(!s.Empty())
	check(s.Contains(0))
	check(!s.Contains(1))
	check(s.Contains(492409))

	s.Remove(0)
	check(!s.Empty())
	check(!s.Contains(0))
	check(!s.Contains(1))
	check(s.Contains(492409))

	s.Remove(492409)
	check(s.Empty())
	check(!s.Contains(0))
	check(!s.Contains(1))
	check(!s.Contains(492409))
}

func randUint64() uint64 {
	lo := uint64(rand.Uint32())
	hi := uint64(rand.Uint32())
	return (hi << 32) | lo
}

func TestLots(t *testing.T) {
	var s SparseSet
	nums := make([]uint64, 1e3)
	for i := 0; i < len(nums); i++ {
		nums[i] = randUint64()
	}

	for i, n := range nums {
		if s.Size() != i {
			t.Fatalf("s.Size() = %d, want %d", s.Size(), i)
		}
		s.Add(n)

	}
	fmt.Printf("size=%4d, bytes=%6d\n", s.Size(), s.MemSize())
	for _, n := range nums {
		if !s.Contains(n) {
			t.Errorf("does not contain %d", n)
		}
	}
	for i, n := range nums {
		got := s.Size()
		want := len(nums) - i
		if got != want {
			t.Fatalf("s.Size() = %d, want %d", got, want)
		}
		s.Remove(n)
	}
	for _, n := range nums {
		if s.Contains(n) {
			t.Errorf("does contain %d", n)
		}
	}
}

	
	
func TestConsecutive(t *testing.T) {
	for _, start := range []uint64{0, 100, 1e8} {
		for _, sz := range []int{0, 1, 2, 3, 4, 5, 64, 256, 512, 1000, 10000, 100000} {
			var s SparseSet
			for i := 0; i < sz; i++ {
				s.Add(uint64(i) + start)
			}
			fmt.Printf("consec: size=%d, start=%d, bytes=%d\n", s.Size(), start, s.MemSize())
		}
	}

	fmt.Printf("memsize of set256: %d\n", memSize(Set256{}))
	fmt.Printf("memsize of node: %d\n", memSize(node{}))
	fmt.Printf("memsize of *node: %d\n", memSize(&node{}))
}
func TestMemSize(t *testing.T) {
	var s SparseSet
	for i := 0; i < 3; i++ {
		fmt.Printf("Add %d\n", i)
		s.Add(uint64(i))
	}
	fmt.Printf("consec: size=%d, bytes=%d\n", s.Size(), s.MemSize())
}