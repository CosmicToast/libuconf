package libuconf

import "strings"

// EnvOpt describes an option that can be set from the environment
type EnvOpt interface {
	Setter
	Env() string // example: A_B_C
}

// FlagOpt describes an option that can be set using a flag
type FlagOpt interface {
	Setter
	Bool() bool
	Flag() string    // example: a.b.c
	Help() string    // UI-only
	ShortFlag() rune // can be 0: means disabled; example: x
}

// Getter describes an option that can return its value
// All of the built-in option types will implement this
// Primarily used in tests, but you may find it useful!
type Getter interface {
	Get() interface{}
}

// Setter describes a type that can set() itself
type Setter interface {
	Set(interface{}) error
}

// TomlOpt describes an option that can be set in a TOML file
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
