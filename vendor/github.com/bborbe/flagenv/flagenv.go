package flagenv

import (
	"flag"
	"strings"
)

func parameterNameToEnvName(name string) string {
	return strings.Replace(strings.ToUpper(name), "-", "_", -1)
}

func Parse() {
	flag.Parse()
}

// PrintDefaults wraps flag.PrintDefaults
func PrintDefaults() {
	flag.CommandLine.PrintDefaults()
}
