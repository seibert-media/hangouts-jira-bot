package flagenv

import (
	"flag"
	"os"
	"strconv"

	"github.com/golang/glog"
)

func Int(name string, value int, usage string) *int {
	return flag.Int(name, envInt(parameterNameToEnvName(name), value), usage)
}

func IntVar(p *int, name string, value int, usage string) {
	flag.IntVar(p, name, envInt(parameterNameToEnvName(name), value), usage)
}

func envInt(key string, def int) int {
	if env := os.Getenv(key); env != "" {
		val, err := strconv.Atoi(env)
		if err != nil {
			glog.V(2).Infof("invalid value for %q: using default: %q", key, def)
			return def
		}
		return val
	}
	return def
}
