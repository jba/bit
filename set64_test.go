package bit

import 	"testing"


func TestSet64(t *testing.T) {
	var s Set64
	s.Add(3)
	if uint64(s) != 8 {
		t.Errorf("Add(3): got %v, want 8", s)
	}
}

