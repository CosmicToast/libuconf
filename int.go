package libuconf

import (
	"fmt"
	"reflect"
	"strconv"
)

// ensure interface compliance
var (
	_ EnvOpt  = (*IntOpt)(nil)
	_ FlagOpt = (*IntOpt)(nil)
	_ Getter  = (*IntOpt)(nil)
	_ Setter  = (*IntOpt)(nil)
	_ TomlOpt = (*IntOpt)(nil)
)

// IntOpt represents a 64-bit integer Option.
type IntOpt struct {
	help  string
	name  string
	sname rune
	val   *int64
}

// NewIntOpt instantiates a IntOpt and returns needed implementation details.
func NewIntOpt(name string, sname rune, val int64, help string) (*IntOpt, *int64) {
	out := IntOpt{help, name, sname, &val}
	return &out, out.val
}

// ---- EnvOpt

// Env returns the option's environment search string.
// For example, if the app name is "APP" and Env() returns "FOO", we will look
// for the environment variable "APP_FOO".
func (r *IntOpt) Env() string {
	return env(r)
}

// ---- FlagOpt

// Bool returns whether or not this option is a boolean.
// This is important because ParseFlags() will handle them differently.
func (*IntOpt) Bool() bool {
	return true
}

// Flag returns the long-form flag for this option.
// All strings are valid, but non-printable ones aren't useful.
func (r *IntOpt) Flag() string {
	return r.name
}

// Help returns the help string for this option.
// It is only used in the Usage() call.
func (r *IntOpt) Help() string {
	return r.help
}

// ShortFlag returns the short-form flag for this option.
// All runes are valid, but non-printable ones aren't useful.
// 0 means "disabled".
func (r *IntOpt) ShortFlag() rune {
	return r.sname
}

// ---- Getter

// Get returns the internal value.
// This is primarily used in Usage() to show the current value for options.
func (r *IntOpt) Get() interface{} {
	return *r.val
}

// ---- Setter

// Set sets this option's value.
// It should be able to handle whatever type it might receive, which means it
// must at least handle strings for being usable in ParseFlags.
func (r *IntOpt) Set(vv interface{}) error {
	val := reflect.ValueOf(vv)
	switch v := vv.(type) {
	case string:
		i, e := strconv.ParseInt(v, 0, 0)
		if e != nil {
			return fmt.Errorf("%w: to %+v", ErrSet, v)
		}
		*r.val = i
	case int8, int16, int32, int64, int:
		*r.val = val.Int()
	case uint8, uint16, uint32, uint64, uint:
		*r.val = int64(val.Uint()) // WARNING: MAY TRUNCATE
	case float32, float64:
		*r.val = int64(val.Float()) // WARNING: WILL TRUNCATE
	default:
		return fmt.Errorf("%w: to %+v", ErrSet, vv)
	}
	return nil
}

// ---- TomlOpt

// Toml returns the option's config file search string.
// It's passed as-is to toml.Tree.Get().
func (r *IntOpt) Toml() string {
	return _toml(r)
}
