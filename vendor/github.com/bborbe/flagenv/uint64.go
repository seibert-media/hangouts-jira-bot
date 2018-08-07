package flagenv

import (
	"flag"
	"os"
	"strconv"

	"github.com/golang/glog"
)

func Uint64(name string, value uint64, usage string) *uint64 {
	return flag.Uint64(name, envUint64(parameterNameToEnvName(name), value), usage)
}

func Uint64Var(p *uint64, name string, value uint64, usage string) {
	flag.Uint64Var(p, name, envUint64(parameterNameToEnvName(name), value), usage)
}

func envUint64(key string, def uint64) uint64 {
	if env := os.Getenv(key); env != "" {
		val, err := strconv.Atoi(env)
		if err != nil {
			glog.V(2).Infof("invalid value for %q: using default: %q", key, def)
			return def
		}
		return uint64(val)
	}
	return def
}
