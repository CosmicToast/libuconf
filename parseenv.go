package libuconf

import (
	"os"
	"strings"
)

// ParseEnv will look for environment variables APPNAME_`option.Env()`
// If one is found, it will try to set it.
func (o *OptionSet) ParseEnv() error {
	prefix := strings.ToUpper(o.AppName)
	for _, vv := range o.Options {
		v, ok := vv.(EnvOpt)
		if !ok {
			continue
		}

		env := prefix + "_" + v.Env()
		res, ok := os.LookupEnv(env)
		if !ok {
			continue
		}

		if err := v.Set(res); err != nil {
			return err
		}
	}
	return nil
}
