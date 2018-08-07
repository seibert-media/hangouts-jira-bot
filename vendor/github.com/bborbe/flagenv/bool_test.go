package flagenv

import (
	"os"
	"testing"
	. "github.com/bborbe/assert"
)

func TestBoolDefaultFalse(t *testing.T) {
	parameterPtr := Bool("parameter7", false, "usage")
	Parse()
	if err := AssertThat(*parameterPtr, Is(false)); err != nil {
		t.Fatal(err)
	}
}

func TestBoolDefaultFalseEnv(t *testing.T) {
	if err := os.Setenv("PARAMETER8", "false"); err != nil {
		t.Fatal(err)
	}
	parameterPtr := Bool("parameter8", true, "usage")
	Parse()
	if err := AssertThat(*parameterPtr, Is(false)); err != nil {
		t.Fatal(err)
	}
}

func TestBoolDefaultTrue(t *testing.T) {
	parameterPtr := Bool("parameter5", true, "usage")
	Parse()
	if err := AssertThat(*parameterPtr, Is(true)); err != nil {
		t.Fatal(err)
	}
}

func TestBoolDefaultTrueEnv(t *testing.T) {
	if err := os.Setenv("PARAMETER6", "true"); err != nil {
		t.Fatal(err)
	}
	parameterPtr := Bool("parameter6", false, "usage")
	Parse()
	if err := AssertThat(*parameterPtr, Is(true)); err != nil {
		t.Fatal(err)
	}
}
