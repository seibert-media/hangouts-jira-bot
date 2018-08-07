package flagenv

import (
	"flag"
	"os"
)

func String(name string, value string, usage string) *string {
	return flag.String(name, envString(parameterNameToEnvName(name), value), usage)
}

func StringVar(p *string, name string, value string, usage string) {
	flag.StringVar(p, name, envString(parameterNameToEnvName(name), value), usage)
}

func envString(key, def string) string {
	if env := os.Getenv(key); env != "" {
		return env
	}
	return def
}
