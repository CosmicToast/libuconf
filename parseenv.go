package libuconf

import (
	"os"
	"strings"
)

// ParseEnv will look for environment variables APPNAME_`option.Env()`
// If one is found, it will try to set it.
func (o *OptionSet) ParseEnv() (e error) {
	prefix := strings.ToUpper(o.AppName)
	o.VisitEnv(func(v EnvOpt) {
		env := prefix + "_" + v.Env()
		res, ok := os.LookupEnv(env)
		if !ok {
			return // continue
		}
		if err := v.Set(res); err != nil {
			e = err
		}
	})
	return
}
