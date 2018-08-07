package flagenv

import (
	"flag"
	"os"
	"strconv"

	"github.com/golang/glog"
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
			glog.V(2).Infof("invalid value for %q: using default: %q", key, def)
			return def
		}
		return val
	}
	return def
}
