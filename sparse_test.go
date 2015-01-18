package bit

import (
	"math/rand"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"testing"
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

type uslice []uint64

func (u uslice) Len() int           { return len(u) }
func (u uslice) Swap(i, j int)      { u[i], u[j] = u[j], u[i] }
func (u uslice) Less(i, j int) bool { return u[i] < u[j] }

func TestSparseElements1(t *testing.T) {
	var s SparseSet
	els := []uint64{3, 17, 300, 12345, 1e8}
	for _, e := range els {
		s.Add(e)
	}
	if !s.Contains(1e8) {
		t.Fatal("no 1e8")
	}
	a := make([]uint64, len(els), len(els))
	n := s.Elements(a, 0)
	got := a[:n]
	if !reflect.DeepEqual(got, els) {
		t.Fatalf("got %v, want %v", got, els)
	}
}

func TestSparseElements2(t *testing.T) {
	var s SparseSet
	nums := make([]uint64, 1e3)
	for i := 0; i < len(nums); i++ {
		nums[i] = randUint64()
	}
	for _, n := range nums {
		s.Add(n)
	}
	sort.Sort(uslice(nums))
	a := make([]uint64, len(nums), len(nums))
	if s.Size() != len(nums) {
		t.Fatalf("size: got %d", s.Size())
	}

	n := s.Elements(a, 0)
	if n != len(nums) {
		t.Fatalf("len: got %d, want %d", n, len(nums))
	}
	if !reflect.DeepEqual(a[:n], nums) {
		t.Fatal("not equal")
	}
}

func set(x ...uint64) []uint64 { return x }

func TestString(t *testing.T) {
	for _, test := range []struct{
		els []uint64
		want string
	} {
		{nil, "{}"},
		{set(9), "{9}"},
		{set(9, 1e4, 99), "{9, 99, 10000}"},
	} {
		got := NewSparseSet(test.els...).String()
		if got != test.want {
			t.Errorf("%v: got %q, want %q", test.els, got, test.want)
		}
	}
}

func TestIntersect(t *testing.T) {
	for _, test := range []struct{
		els1, els2, want []uint64
	} {
		{nil, nil, nil},
		{set(9), nil, nil},
		{nil, set(9), nil},
		{set(9), set(9), set(9)},
		{set(9), set(9, 10), set(9)},
		{set(9, 99, 1e8), set(99, 1e8+1), set(99)},
	} {
		s1 := NewSparseSet(test.els1...)
		s2 := NewSparseSet(test.els2...)
		want := NewSparseSet(test.want...)
		var got SparseSet
		got.Intersect(s1, s2)
		if !got.Equal(want) {
			t.Errorf("%s & %s = %v, want %v", s1, s2, got, want)
		}
	}
}
	
		
		

// func TestConsecutive(t *testing.T) {
// 	for _, start := range []uint64{0, 100, 1e8} {
// 		for _, sz := range []int{0, 1, 2, 3, 4, 5, 64, 256, 512, 1000, 10000, 100000} {
// 			var s SparseSet
// 			for i := 0; i < sz; i++ {
// 				s.Add(uint64(i) + start)
// 			}
// 			fmt.Printf("consec: size=%d, start=%d, bytes=%d\n", s.Size(), start, s.MemSize())
// 		}
// 	}

// 	fmt.Printf("memsize of set256: %d\n", memSize(Set256{}))
// 	fmt.Printf("memsize of node: %d\n", memSize(node{}))
// 	fmt.Printf("memsize of *node: %d\n", memSize(&node{}))
// }
// func TestMemSize(t *testing.T) {
// 	var s SparseSet
// 	for i := 0; i < 3; i++ {
// 		fmt.Printf("Add %d\n", i)
// 		s.Add(uint64(i))
// 	}
// 	fmt.Printf("consec: size=%d, bytes=%d\n", s.Size(), s.MemSize())
// }
