package flagenv

import (
	"flag"
	"os"
	"strconv"

	"github.com/golang/glog"
)

func Int64(name string, value int64, usage string) *int64 {
	return flag.Int64(name, envInt64(parameterNameToEnvName(name), value), usage)
}

func Int64Var(p *int64, name string, value int64, usage string) {
	flag.Int64Var(p, name, envInt64(parameterNameToEnvName(name), value), usage)
}

func envInt64(key string, def int64) int64 {
	if env := os.Getenv(key); env != "" {
		val, err := strconv.Atoi(env)
		if err != nil {
			glog.V(2).Infof("invalid value for %q: using default: %q", key, def)
			return def
		}
		return int64(val)
	}
	return def
}
