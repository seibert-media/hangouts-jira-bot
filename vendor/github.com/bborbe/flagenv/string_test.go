package flagenv

import (
	"os"
	"testing"
	. "github.com/bborbe/assert"
)

func TestStringDefault(t *testing.T) {
	parameterPtr := String("parameter1", "bar", "usage")
	Parse()
	if err := AssertThat(*parameterPtr, Is("bar")); err != nil {
		t.Fatal(err)
	}
}

func TestStringDefaultEnv(t *testing.T) {
	if err := os.Setenv("PARAMETER2", "foo"); err != nil {
		t.Fatal(err)
	}
	parameterPtr := String("parameter2", "bar", "usage")
	Parse()
	if err := AssertThat(*parameterPtr, Is("foo")); err != nil {
		t.Fatal(err)
	}
}
