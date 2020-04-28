package libuconf

import "fmt"

// ensure interface compliance
var (
	_ FlagOpt = (*HelpOpt)(nil)
	_ Setter  = (*HelpOpt)(nil)
)

// HelpOpt is a special Option to provide -h and --help.
type HelpOpt bool

// ---- FlagOpt

// Bool returns whether or not this option is a boolean.
// This is important because ParseFlags() will handle them differently.
//
// This is always true for HelpOpt.
func (*HelpOpt) Bool() bool {
	return true
}

// Flag returns the long-form flag for this option.
// All strings are valid, but non-printable ones aren't useful.
//
// This is always "help" for HelpOpt.
func (r *HelpOpt) Flag() string {
	return "help"
}

// Help returns the help string for this option.
// It is only used in the Usage() call.
//
// It is always "view this help message" for HelpOpt.
func (r *HelpOpt) Help() string {
	return "view this help message"
}

// ShortFlag returns the short-form flag for this option.
// All runes are valid, but non-printable ones aren't useful.
// 0 means "disabled".
//
// This is always 'h' for HelpOpt.
func (r *HelpOpt) ShortFlag() rune {
	return 'h'
}

// ---- Setter

// Set sets this option's value.
// It should be able to handle whatever type it might receive, which means it
// must at least handle strings for being usable in ParseFlags.
//
// It only consumes boolean options, but since it only implements FlagOpt, the only
// time that will happen is when it's being set implicitly.
// The side-effect of this is that it will never consume a cli flag value.
func (r *HelpOpt) Set(vv interface{}) error {
	if v, ok := vv.(bool); ok {
		*r = HelpOpt(v)
		return nil
	}
	return fmt.Errorf("%w: help to %+v", ErrSet, vv)
}
