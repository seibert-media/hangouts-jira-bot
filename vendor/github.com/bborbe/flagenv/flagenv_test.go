package flagenv

import (
	"testing"
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
