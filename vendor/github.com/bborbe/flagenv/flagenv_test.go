package flagenv

import (
	"testing"

	"os"

	"time"

	. "github.com/bborbe/assert"
)

func TestParameterToEnvName(t *testing.T) {
	err := AssertThat(parameterNameToEnvName("name"), Is("NAME"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestParameterToEnvNameMinusToUnderscore(t *testing.T) {
	err := AssertThat(parameterNameToEnvName("log-level"), Is("LOG_LEVEL"))
	if err != nil {
		t.Fatal(err)
	}
}

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
