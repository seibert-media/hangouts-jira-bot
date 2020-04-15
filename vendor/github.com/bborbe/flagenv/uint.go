package flagenv

import (
	"flag"
	"log"
	"os"
	"strconv"
)

func Uint(name string, value uint, usage string) *uint {
	return flag.Uint(name, envUint(parameterNameToEnvName(name), value), usage)
}

func UintVar(p *uint, name string, value uint, usage string) {
	flag.UintVar(p, name, envUint(parameterNameToEnvName(name), value), usage)
}

func envUint(key string, def uint) uint {
	if env := os.Getenv(key); env != "" {
		val, err := strconv.Atoi(env)
		if err != nil {
			log.Printf("invalid value for %q: using default: %q", key, def)
			return def
		}
		return uint(val)
	}
	return def
}
