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
func (o *OptionSet) BoolVar(out *bool, name string, val bool, help string) {
	o.ShortBoolVar(out, name, 0, help)
}

// Bool adds a BoolOpt to the OptionSet
func (o *OptionSet) Bool(name string, val bool, help string) *bool {
	return o.ShortBool(name, 0, val, help)
}

// ShortBoolVar adds a BoolOpt to the OptionSet
func (o *OptionSet) ShortBoolVar(out *bool, name string, sname rune, help string) {
	sopt := &BoolOpt{help, name, sname, out}
	o.Var(sopt)
}

// ShortBool adds a BoolOpt to the Option Set
func (o *OptionSet) ShortBool(name string, sname rune, val bool, help string) *bool {
	out := &val
	o.ShortBoolVar(out, name, sname, help)
	return out
}

// ---- EnvOpt

// Env returns the option's environment search string
// For example, if the app name is APP and Env() returns "FOO"
// We will look for an env var APP_FOO
func (b *BoolOpt) Env() string {
	return env(b)
}

// ---- FlagOpt

// Bool returns whether or not this option is a boolean
func (*BoolOpt) Bool() bool {
	return true
}

// Flag returns the long-form flag for this option
func (b *BoolOpt) Flag() string {
	return b.name
}

// Help returns the help string for this option
func (b *BoolOpt) Help() string {
	return b.help
}

// ShortFlag returns the short-form flag for this option
func (b *BoolOpt) ShortFlag() rune {
	return b.sname
}

// ---- Getter

// Get returns the internal value
func (b *BoolOpt) Get() interface{} {
	return *b.val
}

// ---- Setter

// Set sets this option's value
func (b *BoolOpt) Set(vv interface{}) error {
	switch v := vv.(type) {
	case string:
		switch strings.ToLower(v) {
		case "true":
			*b.val = true
		case "false":
			*b.val = false
		default:
			return fmt.Errorf("%w: to %s", ErrSet, v)
		}
	case bool:
		*b.val = v
	default:
		return fmt.Errorf("%w: to %+v", ErrSet, vv)
	}
	return nil
}

// ---- TomlOpt

// Toml returns the option's config file search string
// It's passed as-is to toml.Tree.Get()
func (b *BoolOpt) Toml() string {
	return _toml(b)
}
