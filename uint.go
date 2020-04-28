package libuconf

import (
	"fmt"
	"reflect"
	"strconv"
)

// ensure interface compliance
var (
	_ EnvOpt  = (*UintOpt)(nil)
	_ FlagOpt = (*UintOpt)(nil)
	_ Getter  = (*UintOpt)(nil)
	_ Setter  = (*UintOpt)(nil)
	_ TomlOpt = (*UintOpt)(nil)
)

// UintOpt represents an unsigned 64-bit integer Option.
type UintOpt struct {
	help  string
	name  string
	sname rune
	val   *uint64
}

// NewUintOpt instantiates a UintOpt and returns needed implementation details.
func NewUintOpt(name string, sname rune, val uint64, help string) (*UintOpt, *uint64) {
	out := UintOpt{help, name, sname, &val}
	return &out, out.val
}

// ---- EnvOpt

// Env returns the option's environment search string.
// For example, if the app name is "APP" and Env() returns "FOO", we will look
// for the environment variable "APP_FOO".
func (r *UintOpt) Env() string {
	return env(r)
}

// ---- FlagOpt

// Bool returns whether or not this option is a boolean.
// This is important because ParseFlags() will handle them differently.
func (*UintOpt) Bool() bool {
	return true
}

// Flag returns the long-form flag for this option.
// All strings are valid, but non-printable ones aren't useful.
func (r *UintOpt) Flag() string {
	return r.name
}

// Help returns the help string for this option.
// It is only used in the Usage() call.
func (r *UintOpt) Help() string {
	return r.help
}

// ShortFlag returns the short-form flag for this option.
// All runes are valid, but non-printable ones aren't useful.
// 0 means "disabled".
func (r *UintOpt) ShortFlag() rune {
	return r.sname
}

// ---- Getter

// Get returns the internal value.
// This is primarily used in Usage() to show the current value for options.
func (r *UintOpt) Get() interface{} {
	return *r.val
}

// ---- Setter

// Set sets this option's value.
// It should be able to handle whatever type it might receive, which means it
// must at least handle strings for being usable in ParseFlags.
func (r *UintOpt) Set(vv interface{}) error {
	val := reflect.ValueOf(vv)
	switch v := vv.(type) {
	case string:
		i, e := strconv.ParseUint(v, 0, 0)
		if e != nil {
			return fmt.Errorf("%w: to %+v", ErrSet, v)
		}
		*r.val = i
	case int8, int16, int32, int64, int:
		*r.val = uint64(val.Int())
	case uint8, uint16, uint32, uint64, uint:
		*r.val = val.Uint()
	case float32, float64:
		*r.val = uint64(val.Float()) // WARNING: WILL TRUNCATE
	default:
		return fmt.Errorf("%w: to %+v", ErrSet, vv)
	}
	return nil
}

// ---- TomlOpt

// Toml returns the option's config file search string.
// It's passed as-is to toml.Tree.Get().
func (r *UintOpt) Toml() string {
	return _toml(r)
}
