package libuconf

import "strings"

/*
EnvOpt describes an option that can be set from the environment.

It is used by ParseEnv() to determine what environment variables to read.
For example, if the app name is "APP" and Env() returns "FOO", we will look for
the environment variable "APP_FOO".
*/
type EnvOpt interface {
	Setter
	Env() string // example: A_B_C
}

/*
FlagOpt describes an option that can be set using a flag.
It is used by ParseFlags.

If Bool returns true, ParseFlags will know that a value is technically optional
for the given flag.
If it cannot find a value for it, it will thus pass "true" (bool) to Set().

Flag and ShortFlag determine what ParseFlags looks for on the command line, as
well as what Usage will print out.

Finally, Help changes the help string in Usage.
*/
type FlagOpt interface {
	Setter
	Bool() bool
	Flag() string // example: a.b.c
	Help() string
	ShortFlag() rune // can be 0: means disabled; example: x
}

/*
Getter describes an option that can return its *current* value.

This is used in tests, as well as in Usage.
Usage in particular allows you to observe the current value of your flags, which
is useful to see whether or not your config files (for example) are being
applied.

It may also be of interest to consumers of the library, via Visit.
*/
type Getter interface {
	Get() interface{}
}

/*
Setter describes a type that can be Set.
All options must implement Setter.

Set takes an interface, and should use type switch and/or reflection to take the
appropriate course of action based on what ended up being passed to it.
It is highly recommended to at least handle strings.

If your option implements FlagOpt and Bool returns true, Set is also expected to
handle booleans.
*/
type Setter interface {
	Set(interface{}) error
}

/*
TomlOpt describes an option that can be set by a value in a TOML file.

The return value of Toml should be the exact search query to plug in.
Set will be called with whatever the get call returns (unless it's nil).
*/
type TomlOpt interface {
	Setter
	Toml() string // a TOML query string; example: a.b.c
}

func env(f FlagOpt) string {
	return strings.ToUpper(strings.ReplaceAll(f.Flag(), ".", "_"))
}

func _toml(f FlagOpt) string {
	return f.Flag()
}
