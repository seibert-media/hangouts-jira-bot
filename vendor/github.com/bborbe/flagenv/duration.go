package flagenv

import (
	"flag"
	"os"
	"time"

	"github.com/golang/glog"
)

func Duration(name string, value time.Duration, usage string) *time.Duration {
	return flag.Duration(name, envDuration(parameterNameToEnvName(name), value), usage)
}

func DurationVar(p *time.Duration, name string, value time.Duration, usage string) {
	flag.DurationVar(p, name, envDuration(parameterNameToEnvName(name), value), usage)
}

func envDuration(key string, def time.Duration) time.Duration {
	if env := os.Getenv(key); env != "" {
		val, err := time.ParseDuration(env)
		if err != nil {
			glog.V(2).Infof("invalid value for %q: using default: %q", key, def)
			return def
		}
		return val
	}
	return def
}
