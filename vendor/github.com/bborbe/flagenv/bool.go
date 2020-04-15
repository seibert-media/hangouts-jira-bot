package flagenv

import (
	"flag"
	"os"
	"strconv"
)

func Bool(name string, value bool, usage string) *bool {
	return flag.Bool(name, envBool(parameterNameToEnvName(name), value), usage)
}

func BoolVar(p *bool, name string, value bool, usage string) {
	flag.BoolVar(p, name, envBool(parameterNameToEnvName(name), value), usage)
}

func envBool(key string, def bool) bool {
	if env := os.Getenv(key); env != "" {
		res, err := strconv.ParseBool(env)
		if err != nil {
			return def
		}
		return res
	}
	return def
}
