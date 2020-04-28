package libuconf

import "fmt"

// ensure interface compliance
var (
	_ EnvOpt  = (*StringOpt)(nil)
	_ FlagOpt = (*StringOpt)(nil)
	_ Getter  = (*StringOpt)(nil)
	_ Setter  = (*StringOpt)(nil)
	_ TomlOpt = (*StringOpt)(nil)
)

// StringOpt represents a string Option.
type StringOpt struct {
	help  string
	name  string
	sname rune
	val   *string
}

// NewStringOpt instantiates a StringOpt and returns needed implementation details.
func NewStringOpt(name string, sname rune, val string, help string) (*StringOpt, *string) {
	out := StringOpt{help, name, sname, &val}
	return &out, out.val
}

// ---- EnvOpt

// Env returns the option's environment search string.
// For example, if the app name is "APP" and Env() returns "FOO", we will look
// for the environment variable "APP_FOO".
func (r *StringOpt) Env() string {
	return env(r)
}

// ---- FlagOpt

// Bool returns whether or not this option is a boolean.
// This is important because ParseFlags() will handle them differently.
func (*StringOpt) Bool() bool {
	return false
}

// Flag returns the long-form flag for this option.
// All strings are valid, but non-printable ones aren't useful.
func (r *StringOpt) Flag() string {
	return r.name
}

// Help returns the help string for this option.
// It is only used in the Usage() call.
func (r *StringOpt) Help() string {
	return r.help
}

// ShortFlag returns the short-form flag for this option.
// All runes are valid, but non-printable ones aren't useful.
// 0 means "disabled".
func (r *StringOpt) ShortFlag() rune {
	return r.sname
}

// ---- Getter

// Get returns the internal value.
// This is primarily used in Usage() to show the current value for options.
func (r *StringOpt) Get() interface{} {
	return *r.val
}

// ---- Setter

// Set sets this option's value.
// It should be able to handle whatever type it might receive, which means it
// must at least handle strings for being usable in ParseFlags.
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

// Toml returns the option's config file search string.
// It's passed as-is to toml.Tree.Get().
func (r *StringOpt) Toml() string {
	return _toml(r)
}
