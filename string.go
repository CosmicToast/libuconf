package libuconf

import "fmt"

// ensure interface compliance
var (
	_ EnvOpt  = &StringOpt{}
	_ FlagOpt = &StringOpt{}
	_ Getter  = &StringOpt{}
	_ Setter  = &StringOpt{}
	_ TomlOpt = &StringOpt{}
)

// StringOpt represents a string Option
type StringOpt struct {
	help  string
	name  string
	sname rune
	val   *string
}

// ---- integration with OptionSet

// StringVar adds a StringOpt to the OptionSet
func (o *OptionSet) StringVar(out *string,
	name string,
	val string,
	help string) *StringOpt {
	return o.ShortStringVar(out, name, 0, val, help)
}

// String adds a StringOpt to the OptionSet
func (o *OptionSet) String(name string,
	val string,
	help string) *string {
	return o.ShortString(name, 0, val, help)
}

// ShortStringVar adds a StringOpt to the OptionSet
func (o *OptionSet) ShortStringVar(out *string,
	name string, sname rune,
	val string,
	help string) *StringOpt {

	*out = val
	sopt := &StringOpt{help, name, sname, out}
	o.Var(sopt)
	return sopt
}

// ShortString adds a StringOpt to the Option Set
func (o *OptionSet) ShortString(name string, sname rune,
	val string,
	help string) *string {

	var out string
	o.ShortStringVar(&out, name, sname, val, help)
	return &out
}

// ---- EnvOpt

// Env returns the option's environment search string
// For example, if the app name is APP and Env() returns "FOO"
// We will look for an env var APP_FOO
func (r *StringOpt) Env() string {
	return env(r)
}

// ---- FlagOpt

// Bool returns whether or not this option is a boolean
func (*StringOpt) Bool() bool {
	return false
}

// Flag returns the long-form flag for this option
func (r *StringOpt) Flag() string {
	return r.name
}

// Help returns the help string for this option
func (r *StringOpt) Help() string {
	return r.help
}

// ShortFlag returns the short-form flag for this option
func (r *StringOpt) ShortFlag() rune {
	return r.sname
}

// ---- Getter

// Get returns the internal value
func (r *StringOpt) Get() interface{} {
	return *r.val
}

// ---- Setter

// Set sets this option's value
func (r *StringOpt) Set(vv interface{}) error {
	switch v := vv.(type) {
	case string:
		*r.val = v
	default:
		*r.val = fmt.Sprint(vv)
	}
	return nil
}

// ---- TomlOpt

// Toml returns the option's config file search string
// It's passed as-is to toml.Tree.Get()
func (r *StringOpt) Toml() string {
	return _toml(r)
}
