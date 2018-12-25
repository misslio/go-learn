package simplepath

import "testing"

func TestSqrt1(t *testing.T) {
    r := Sqrt(4)
    if r != 2 {
        t.Errorf("Sqrt(2) failed. Got %d. expected 2", r)
    }
}
