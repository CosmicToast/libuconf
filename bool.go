package libuconf

import (
	"fmt"
	"strings"
)

// ensure interface compliance
var (
	_ EnvOpt  = &BoolOpt{}
	_ FlagOpt = &BoolOpt{}
	_ Getter  = &BoolOpt{}
	_ Setter  = &BoolOpt{}
	_ TomlOpt = &BoolOpt{}
)

// BoolOpt represents a boolean Option
type BoolOpt struct {
	help  string
	name  string
	sname rune
	val   *bool
}

// ---- integration with OptionSet

// BoolVar adds a BoolOpt to the OptionSet
func (o *OptionSet) BoolVar(out *bool,
	name string,
	val bool,
	help string) *BoolOpt {
	return o.ShortBoolVar(out, name, 0, val, help)
}

// Bool adds a BoolOpt to the OptionSet
func (o *OptionSet) Bool(name string,
	val bool,
	help string) *bool {
	return o.ShortBool(name, 0, val, help)
}

// ShortBoolVar adds a BoolOpt to the OptionSet
func (o *OptionSet) ShortBoolVar(out *bool,
	name string, sname rune,
	val bool,
	help string) *BoolOpt {

	*out = val
	sopt := &BoolOpt{help, name, sname, out}
	o.Var(sopt)
	return sopt
}

// ShortBool adds a BoolOpt to the Option Set
func (o *OptionSet) ShortBool(name string, sname rune,
	val bool,
	help string) *bool {

	var out bool
	o.ShortBoolVar(&out, name, sname, val, help)
	return &out
}

// ---- EnvOpt

// Env returns the option's environment search string
// For example, if the app name is APP and Env() returns "FOO"
// We will look for an env var APP_FOO
func (r *BoolOpt) Env() string {
	return env(r)
}

// ---- FlagOpt

// Bool returns whether or not this option is a boolean
func (*BoolOpt) Bool() bool {
	return true
}

// Flag returns the long-form flag for this option
func (r *BoolOpt) Flag() string {
	return r.name
}

// Help returns the help string for this option
func (r *BoolOpt) Help() string {
	return r.help
}

// ShortFlag returns the short-form flag for this option
func (r *BoolOpt) ShortFlag() rune {
	return r.sname
}

// ---- Getter

// Get returns the internal value
func (r *BoolOpt) Get() interface{} {
	return *r.val
}

// ---- Setter

// Set sets this option's value
func (r *BoolOpt) Set(vv interface{}) error {
	switch v := vv.(type) {
	case string:
		switch strings.ToLower(v) {
		case "true":
			*r.val = true
		case "false":
			*r.val = false
		default:
			return fmt.Errorf("%w: to %s", ErrSet, v)
		}
	case bool:
		*r.val = v
	default:
		return fmt.Errorf("%w: to %+v", ErrSet, vv)
	}
	return nil
}

// ---- TomlOpt

// Toml returns the option's config file search string
// It's passed as-is to toml.Tree.Get()
func (r *BoolOpt) Toml() string {
	return _toml(r)
}
