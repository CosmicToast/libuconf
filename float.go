package libuconf

import (
	"fmt"
	"reflect"
	"strconv"
)

// ensure interface compliance
var (
	_ EnvOpt  = &FloatOpt{}
	_ FlagOpt = &FloatOpt{}
	_ Getter  = &FloatOpt{}
	_ Setter  = &FloatOpt{}
	_ TomlOpt = &FloatOpt{}
)

// FloatOpt represents an unsigned integer Option
type FloatOpt struct {
	help  string
	name  string
	sname rune
	val   *float64
}

// ---- integration with OptionSet

// FloatVar adds a FloatOpt to the OptionSet
func (o *OptionSet) FloatVar(out *float64, name string, val float64, help string) {
	o.ShortFloatVar(out, name, 0, help)
}

// Float adds a IntOpt to the OptionSet
func (o *OptionSet) Float(name string, val float64, help string) *float64 {
	return o.ShortFloat(name, 0, val, help)
}

// ShortFloatVar adds a IntOpt to the OptionSet
func (o *OptionSet) ShortFloatVar(out *float64, name string, sname rune, help string) {
	sopt := &FloatOpt{help, name, sname, out}
	o.Var(sopt)
}

// ShortFloat adds a FloatOpt to the Option Set
func (o *OptionSet) ShortFloat(name string, sname rune, val float64, help string) *float64 {
	out := &val
	o.ShortFloatVar(out, name, sname, help)
	return out
}

// ---- EnvOpt

// Env returns the option's environment search string
// For example, if the app name is APP and Env() returns "FOO"
// We will look for an env var APP_FOO
func (b *FloatOpt) Env() string {
	return env(b)
}

// ---- FlagOpt

// Bool returns whether or not this option is a boolean
func (*FloatOpt) Bool() bool {
	return true
}

// Flag returns the long-form flag for this option
func (b *FloatOpt) Flag() string {
	return b.name
}

// Help returns the help string for this option
func (b *FloatOpt) Help() string {
	return b.help
}

// ShortFlag returns the short-form flag for this option
func (b *FloatOpt) ShortFlag() rune {
	return b.sname
}

// ---- Getter

// Get returns the internal value
func (b *FloatOpt) Get() interface{} {
	return *b.val
}

// ---- Setter

// Set sets this option's value
func (b *FloatOpt) Set(vv interface{}) error {
	val := reflect.ValueOf(vv)
	switch v := vv.(type) {
	case string:
		f, e := strconv.ParseFloat(v, 64)
		if e != nil {
			return fmt.Errorf("%w: to %+v", ErrSet, v)
		}
		*b.val = f
	case int8, int16, int32, int64, int:
		*b.val = float64(val.Int())
	case uint8, uint16, uint32, uint64, uint:
		*b.val = float64(val.Uint())
	case float32, float64:
		*b.val = val.Float() // WARNING: WILL TRUNCATE
	default:
		return fmt.Errorf("%w: to %+v", ErrSet, vv)
	}
	return nil
}

// ---- TomlOpt

// Toml returns the option's config file search string
// It's passed as-is to toml.Tree.Get()
func (b *FloatOpt) Toml() string {
	return _toml(b)
}
