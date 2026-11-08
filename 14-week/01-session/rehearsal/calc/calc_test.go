package calc

import (
	"errors"
	"testing"
)

func TestDivideReturnsTheQuotient(t *testing.T) {
	got, err := Divide(10, 2)
	if err != nil || got != 5 {
		t.Fatalf("Divide(10, 2) = %d, %v; want 5, nil", got, err)
	}
}

func TestDivideRejectsZeroDivisor(t *testing.T) {
	_, err := Divide(1, 0)
	if !errors.Is(err, ErrDivideByZero) {
		t.Fatalf("Divide(1, 0) error = %v; want ErrDivideByZero", err)
	}
}
