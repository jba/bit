package bit

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func sampleSet64() Set64 {
	var s Set64
	s.Add(3)
	s.Add(63)
	s.Add(17)
	return s
}

func TestBasics(t *testing.T) {
	s := sampleSet64()
	want := "{3, 17, 63}"
	got := s.String()
	if got != want {
		t.Errorf("s.String() = %q, want %q", got, want)
	}
	if !cmp.Equal(naiveElementsUint8(&s), []uint8{3, 17, 63}) {
		t.Errorf("%s: wrong elements", s)
	}
	if s.Size() != 3 {
		t.Error("wrong size")
	}
	if s.Empty() {
		t.Error("shouldn't be empty")
	}
	var z Set64
	if !z.Empty() {
		t.Error("should be empty")
	}
}

func TestElements(t *testing.T) {
	var a [10]uint8
	s := sampleSet64()
	for _, test := range []struct {
		n     int
		start uint8
		want  []uint8
	}{
		{0, 0, []uint8{}},
		{0, 10, []uint8{}},
		{1, 0, []uint8{3}},
		{1, 5, []uint8{17}},
		{1, 27, []uint8{63}},
		{2, 0, []uint8{3, 17}},
		{2, 5, []uint8{17, 63}},
		{2, 39, []uint8{63}},
		{2, 63, []uint8{63}},
		{2, 83, []uint8{}},
		{3, 0, []uint8{3, 17, 63}},
		{3, 10, []uint8{17, 63}},
		{3, 99, []uint8{}},
	} {
		n := s.Elements(a[:test.n], test.start)
		got := a[:n]
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%+v: got %v, want %v", test, got, test.want)
		}
	}
}

func TestElements64(t *testing.T) {
	var a [10]uint64
	s := sampleSet64()
	for _, test := range []struct {
		n     int
		start uint64
		high  uint64
		want  []uint64
	}{
		{0, 0, 0, []uint64{}},
		{0, 10, 0, []uint64{}},
		{1, 0, 0, []uint64{3}},
		{1, 5, 0, []uint64{17}},
		{1, 27, 0, []uint64{63}},
		{2, 0, 0, []uint64{3, 17}},
		{2, 5, 0, []uint64{17, 63}},
		{2, 39, 0, []uint64{63}},
		{2, 63, 0, []uint64{63}},
		{2, 83, 0, []uint64{}},
		{3, 0, 0, []uint64{3, 17, 63}},
		{3, 10, 0, []uint64{17, 63}},
		{3, 99, 0, []uint64{}},
		{3, 10, 64, []uint64{64 + 17, 64 + 63}},
		{3, 0, 256, []uint64{256 + 3, 256 + 17, 256 + 63}},
	} {
		n := s.Elements64(a[:test.n], uint8(test.start), test.high)
		got := a[:n]
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%+v: got %v", test, got)
		}
	}
}

func TestPosition(t *testing.T) {
	s := sampleSet64()
	for _, test := range []struct {
		n   uint8
		pos int
		in  bool
	}{
		{0, 0, false},
		{1, 0, false},
		{2, 0, false},
		{3, 0, true},
		{4, 1, false},
		{10, 1, false},
		{16, 1, false},
		{17, 1, true},
		{20, 2, false},
		{62, 2, false},
		{63, 2, true},
	} {
		gotPos, gotIn := s.Position(test.n)
		if gotPos != test.pos || gotIn != test.in {
			t.Errorf("Position(%d) = (%d, %t), want (%d, %t)", test.n, gotPos, gotIn, test.pos, test.in)
		}
	}
}

func naiveElementsUint8(s interface {
	Capacity() int
	Contains(uint8) bool
}) []uint8 {
	var els []uint8
	for i := 0; i < s.Capacity(); i++ {
		u := uint8(i)
		if s.Contains(u) {
			els = append(els, u)
		}
	}
	return els
}
