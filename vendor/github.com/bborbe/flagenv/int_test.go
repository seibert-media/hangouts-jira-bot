package flagenv

import (
	"os"
	"testing"
	. "github.com/bborbe/assert"
)

func TestIntDefault(t *testing.T) {
	parameterPtr := Int("parameter3", 1337, "usage")
	Parse()
	if err := AssertThat(*parameterPtr, Is(1337)); err != nil {
		t.Fatal(err)
	}
}

func TestIntDefaultEnv(t *testing.T) {
	if err := os.Setenv("PARAMETER4", "1337"); err != nil {
		t.Fatal(err)
	}
	parameterPtr := Int("parameter4", 42, "usage")
	Parse()
	if err := AssertThat(*parameterPtr, Is(1337)); err != nil {
		t.Fatal(err)
	}
}
