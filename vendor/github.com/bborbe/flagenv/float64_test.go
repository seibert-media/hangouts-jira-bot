package flagenv

import (
	"os"
	"testing"
	. "github.com/bborbe/assert"
)

func TestFloat64Default(t *testing.T) {
	parameterPtr := Float64("float64parameter1", 1337, "usage")
	Parse()
	if err := AssertThat(*parameterPtr, Is(float64(1337))); err != nil {
		t.Fatal(err)
	}
}

func TestFloat64DefaultEnv(t *testing.T) {
	if err := os.Setenv("FLOAT64PARAMETER2", "1337"); err != nil {
		t.Fatal(err)
	}
	parameterPtr := Float64("float64parameter2", 42, "usage")
	Parse()
	if err := AssertThat(*parameterPtr, Is(float64(1337))); err != nil {
		t.Fatal(err)
	}
}
