package flagenv

import (
	"flag"
	"log"
	"os"
	"strconv"
)

func Float64(name string, value float64, usage string) *float64 {
	return flag.Float64(name, envFloat64(parameterNameToEnvName(name), value), usage)
}

func Float64Var(p *float64, name string, value float64, usage string) {
	flag.Float64Var(p, name, envFloat64(parameterNameToEnvName(name), value), usage)
}

func envFloat64(key string, def float64) float64 {
	if env := os.Getenv(key); env != "" {
		val, err := strconv.ParseFloat(env, 64)
		if err != nil {
			log.Printf("invalid value for %v: using default: %v", key, def)
			return def
		}
		return val
	}
	return def
}
