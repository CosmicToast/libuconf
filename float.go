package libuconf

import (
	"fmt"
	"reflect"
	"strconv"
)

// ensure interface compliance
var (
	_ EnvOpt  = (*FloatOpt)(nil)
	_ FlagOpt = (*FloatOpt)(nil)
	_ Getter  = (*FloatOpt)(nil)
	_ Setter  = (*FloatOpt)(nil)
	_ TomlOpt = (*FloatOpt)(nil)
)

// FloatOpt represents a long float Option.
type FloatOpt struct {
	help  string
	name  string
	sname rune
	val   *float64
}

// ---- integration with OptionSet

// FloatVar adds a FloatOpt to the OptionSet
func (o *OptionSet) FloatVar(out *float64,
	name string,
	val float64,
	help string) *FloatOpt {
	return o.ShortFloatVar(out, name, 0, val, help)
}

// Float adds a IntOpt to the OptionSet
func (o *OptionSet) Float(name string,
	val float64,
	help string) *float64 {
	return o.ShortFloat(name, 0, val, help)
}

// ShortFloatVar adds a IntOpt to the OptionSet
func (o *OptionSet) ShortFloatVar(out *float64,
	name string, sname rune,
	val float64,
	help string) *FloatOpt {

	*out = val
	sopt := &FloatOpt{help, name, sname, out}
	o.Var(sopt)
	return sopt
}

// ShortFloat adds a FloatOpt to the Option Set
func (o *OptionSet) ShortFloat(name string, sname rune,
	val float64,
	help string) *float64 {

	var out float64
	o.ShortFloatVar(&out, name, sname, val, help)
	return &out
}

// ---- EnvOpt

// Env returns the option's environment search string.
// For example, if the app name is "APP" and Env() returns "FOO", we will look
// for the environment variable "APP_FOO".
func (r *FloatOpt) Env() string {
	return env(r)
}

// ---- FlagOpt

// Bool returns whether or not this option is a boolean.
// This is important because ParseFlags() will handle them differently.
func (*FloatOpt) Bool() bool {
	return true
}

// Flag returns the long-form flag for this option.
// All strings are valid, but non-printable ones aren't useful.
func (r *FloatOpt) Flag() string {
	return r.name
}

// Help returns the help string for this option.
// It is only used in the Usage() call.
func (r *FloatOpt) Help() string {
	return r.help
}

// ShortFlag returns the short-form flag for this option.
// All runes are valid, but non-printable ones aren't useful.
// 0 means "disabled".
func (r *FloatOpt) ShortFlag() rune {
	return r.sname
}

// ---- Getter

// Get returns the internal value.
// This is primarily used in Usage() to show the current value for options.
func (r *FloatOpt) Get() interface{} {
	return *r.val
}

// ---- Setter

// Set sets this option's value.
// It should be able to handle whatever type it might receive, which means it
// must at least handle strings for being usable in ParseFlags.
func (r *FloatOpt) Set(vv interface{}) error {
	val := reflect.ValueOf(vv)
	switch v := vv.(type) {
	case string:
		f, e := strconv.ParseFloat(v, 64)
		if e != nil {
			return fmt.Errorf("%w: to %+v", ErrSet, v)
		}
		*r.val = f
	case int8, int16, int32, int64, int:
		*r.val = float64(val.Int())
	case uint8, uint16, uint32, uint64, uint:
		*r.val = float64(val.Uint())
	case float32, float64:
		*r.val = val.Float() // WARNING: WILL TRUNCATE
	default:
		return fmt.Errorf("%w: to %+v", ErrSet, vv)
	}
	return nil
}

// ---- TomlOpt

// Toml returns the option's config file search string.
// It's passed as-is to toml.Tree.Get().
func (r *FloatOpt) Toml() string {
	return _toml(r)
}
