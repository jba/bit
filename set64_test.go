package bit

import 	"testing"


func TestSet64(t *testing.T) {
	var s Set64
	s.Add(3)
	if uint64(s) != 8 {
		t.Errorf("Add(3): got %v, want 8", s)
	}
}

func TestPosition(t *testing.T) {
	var s Set64
	s.Add(3)
	s.Add(17)
	for _, test := range []struct{
		n uint8
		pos int
		in bool
	} {
		{0, 0, false},
		{1, 0, false},
		{2, 0, false},
		{3, 0, true},
		{4, 1, false},
		{10, 1, false},
		{16, 1, false},
		{17, 1, true},
		{20, 2, false},
	} {
		gotPos, gotIn := s.Position(test.n)
		if gotPos != test.pos || gotIn != test.in {
			t.Errorf("Position(%d) = (%d, %t), want (%d, %t)", test.n, gotPos, gotIn, test.pos, test.in)
		}
	}
}

