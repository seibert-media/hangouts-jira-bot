package flagenv

import (
	"os"
	"testing"
	"time"
	. "github.com/bborbe/assert"
)

func TestDurationDefault(t *testing.T) {
	parameterPtr := Duration("parameter9", time.Minute*8, "usage")
	Parse()
	if err := AssertThat(*parameterPtr, Is(time.Minute*8)); err != nil {
		t.Fatal(err)
	}
}

func TestDurationDefaultEnv(t *testing.T) {
	if err := os.Setenv("PARAMETER10", "8m"); err != nil {
		t.Fatal(err)
	}
	parameterPtr := Duration("parameter10", time.Minute*7, "usage")
	Parse()
	if err := AssertThat(*parameterPtr, Is(time.Minute*8)); err != nil {
		t.Fatal(err)
	}
}
