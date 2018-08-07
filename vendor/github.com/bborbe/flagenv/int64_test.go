package flagenv

import (
	"os"
	"testing"
	. "github.com/bborbe/assert"
)

func TestInt64Default(t *testing.T) {
	parameterPtr := Int64("int64parameter1", 1337, "usage")
	Parse()
	if err := AssertThat(*parameterPtr, Is(int64(1337))); err != nil {
		t.Fatal(err)
	}
}

func TestInt64DefaultEnv(t *testing.T) {
	if err := os.Setenv("INT64PARAMETER2", "1337"); err != nil {
		t.Fatal(err)
	}
	parameterPtr := Int64("int64parameter2", 42, "usage")
	Parse()
	if err := AssertThat(*parameterPtr, Is(int64(1337))); err != nil {
		t.Fatal(err)
	}
}
