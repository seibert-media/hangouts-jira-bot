package flagenv

import (
	"os"
	"testing"
	. "github.com/bborbe/assert"
)

func TestUintDefault(t *testing.T) {
	parameterPtr := Uint("uintparameter1", 1337, "usage")
	Parse()
	if err := AssertThat(*parameterPtr, Is(uint(1337))); err != nil {
		t.Fatal(err)
	}
}

func TestUintDefaultEnv(t *testing.T) {
	if err := os.Setenv("UINTPARAMETER2", "1337"); err != nil {
		t.Fatal(err)
	}
	parameterPtr := Uint("uintparameter2", 42, "usage")
	Parse()
	if err := AssertThat(*parameterPtr, Is(uint(1337))); err != nil {
		t.Fatal(err)
	}
}
