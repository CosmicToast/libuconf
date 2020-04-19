package libuconf

import (
	"fmt"
	"reflect"
	"strconv"
)

// ensure interface compliance
var (
	_ EnvOpt  = &UintOpt{}
	_ FlagOpt = &UintOpt{}
	_ Getter  = &UintOpt{}
	_ Setter  = &UintOpt{}
	_ TomlOpt = &UintOpt{}
)

// UintOpt represents an unsigned integer Option
type UintOpt struct {
	help  string
	name  string
	sname rune
	val   *uint64
}

// ---- integration with OptionSet

// UintVar adds a UintOpt to the OptionSet
func (o *OptionSet) UintVar(out *uint64, name string, val uint64, help string) {
	o.ShortUintVar(out, name, 0, help)
}

// Uint adds a IntOpt to the OptionSet
func (o *OptionSet) Uint(name string, val uint64, help string) *uint64 {
	return o.ShortUint(name, 0, val, help)
}

// ShortUintVar adds a IntOpt to the OptionSet
func (o *OptionSet) ShortUintVar(out *uint64, name string, sname rune, help string) {
	sopt := &UintOpt{help, name, sname, out}
	o.Var(sopt)
}

// ShortUint adds a UintOpt to the Option Set
func (o *OptionSet) ShortUint(name string, sname rune, val uint64, help string) *uint64 {
	out := &val
	o.ShortUintVar(out, name, sname, help)
	return out
}

// ---- EnvOpt

// Env returns the option's environment search string
// For example, if the app name is APP and Env() returns "FOO"
// We will look for an env var APP_FOO
func (b *UintOpt) Env() string {
	return env(b)
}

// ---- FlagOpt

// Bool returns whether or not this option is a boolean
func (*UintOpt) Bool() bool {
	return true
}

// Flag returns the long-form flag for this option
func (b *UintOpt) Flag() string {
	return b.name
}

// Help returns the help string for this option
func (b *UintOpt) Help() string {
	return b.help
}

// ShortFlag returns the short-form flag for this option
func (b *UintOpt) ShortFlag() rune {
	return b.sname
}

// ---- Getter

// Get returns the internal value
func (b *UintOpt) Get() interface{} {
	return *b.val
}

// ---- Setter

// Set sets this option's value
func (b *UintOpt) Set(vv interface{}) error {
	val := reflect.ValueOf(vv)
	switch v := vv.(type) {
	case string:
		i, e := strconv.ParseUint(v, 0, 0)
		if e != nil {
			return fmt.Errorf("%w: to %+v", ErrSet, v)
		}
		*b.val = i
	case int8, int16, int32, int64, int:
		*b.val = uint64(val.Int())
	case uint8, uint16, uint32, uint64, uint:
		*b.val = val.Uint()
	case float32, float64:
		*b.val = uint64(val.Float()) // WARNING: WILL TRUNCATE
	default:
		return fmt.Errorf("%w: to %+v", ErrSet, vv)
	}
	return nil
}

// ---- TomlOpt

// Toml returns the option's config file search string
// It's passed as-is to toml.Tree.Get()
func (b *UintOpt) Toml() string {
	return _toml(b)
}
