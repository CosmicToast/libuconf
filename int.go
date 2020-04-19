package libuconf

import (
	"fmt"
	"strconv"
  "reflect"
)

// ensure interface compliance
var (
	_ EnvOpt  = &IntOpt{}
	_ FlagOpt = &IntOpt{}
	_ Getter  = &IntOpt{}
	_ Setter  = &IntOpt{}
	_ TomlOpt = &IntOpt{}
)

// IntOpt represents an integer Option
type IntOpt struct {
	help  string
	name  string
	sname rune
	val   *int64
}

// ---- integration with OptionSet

// IntVar adds a IntOpt to the OptionSet
func (o *OptionSet) IntVar(out *int64, name string, val int64, help string) {
	o.ShortIntVar(out, name, 0, help)
}

// Int adds a IntOpt to the OptionSet
func (o *OptionSet) Int(name string, val int64, help string) *int64 {
	return o.ShortInt(name, 0, val, help)
}

// ShortIntVar adds a IntOpt to the OptionSet
func (o *OptionSet) ShortIntVar(out *int64, name string, sname rune, help string) {
	sopt := &IntOpt{help, name, sname, out}
	o.Var(sopt)
}

// ShortInt adds a IntOpt to the Option Set
func (o *OptionSet) ShortInt(name string, sname rune, val int64, help string) *int64 {
	out := &val
	o.ShortIntVar(out, name, sname, help)
	return out
}

// ---- EnvOpt

// Env returns the option's environment search string
// For example, if the app name is APP and Env() returns "FOO"
// We will look for an env var APP_FOO
func (b *IntOpt) Env() string {
	return env(b)
}

// ---- FlagOpt

// Bool returns whether or not this option is a boolean
func (*IntOpt) Bool() bool {
	return true
}

// Flag returns the long-form flag for this option
func (b *IntOpt) Flag() string {
	return b.name
}

// Help returns the help string for this option
func (b *IntOpt) Help() string {
	return b.help
}

// ShortFlag returns the short-form flag for this option
func (b *IntOpt) ShortFlag() rune {
	return b.sname
}

// ---- Getter

// Get returns the internal value
func (b *IntOpt) Get() interface{} {
	return *b.val
}

// ---- Setter

// Set sets this option's value
func (b *IntOpt) Set(vv interface{}) error {
  val := reflect.ValueOf(vv)
	switch v := vv.(type) {
	case string:
		i, e := strconv.ParseInt(v, 0, 0)
    if e != nil {
      return fmt.Errorf("%w: to %+v", ErrSet, v)
    }
    *b.val = i
  case int8, int16, int32, int64, int:
    *b.val = val.Int()
  case uint8, uint16, uint32, uint64, uint:
    *b.val = int64(val.Uint()) // WARNING: MAY TRUNCATE
  case float32, float64:
    *b.val = int64(val.Float()) // WARNING: WILL TRUNCATE
	default:
		return fmt.Errorf("%w: to %+v", ErrSet, vv)
	}
	return nil
}

// ---- TomlOpt

// Toml returns the option's config file search string
// It's passed as-is to toml.Tree.Get()
func (b *IntOpt) Toml() string {
	return _toml(b)
}
