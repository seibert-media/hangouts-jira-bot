package flagenv

import "flag"

func Var(value flag.Value, name string, usage string) {
	flag.Var(value, name, usage)
}
