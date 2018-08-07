package flagenv

import (
	"os"
	"testing"
	. "github.com/bborbe/assert"
)

func TestUint64Default(t *testing.T) {
	parameterPtr := Uint64("uint64parameter1", 1337, "usage")
	Parse()
	if err := AssertThat(*parameterPtr, Is(uint64(1337))); err != nil {
		t.Fatal(err)
	}
}

func TestUint64DefaultEnv(t *testing.T) {
	if err := os.Setenv("UINT64PARAMETER2", "1337"); err != nil {
		t.Fatal(err)
	}
	parameterPtr := Uint64("uint64parameter2", 42, "usage")
	Parse()
	if err := AssertThat(*parameterPtr, Is(uint64(1337))); err != nil {
		t.Fatal(err)
	}
}
