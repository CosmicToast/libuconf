package libuconf

import "fmt"

// HelpOpt is a special Option to provide -h and --help
type HelpOpt bool

// Help adds --help and -h flags to the Option Set
func (o *OptionSet) Help() *bool {
	out := new(bool)
	o.Var((*HelpOpt)(out))
	return out
}

// ---- FlagOpt

// Bool returns whether or not this option is a boolean
func (*HelpOpt) Bool() bool {
	return true
}

// Flag returns the long-form flag for this option
func (r *HelpOpt) Flag() string {
	return "help"
}

// Help returns the help string for this option
func (r *HelpOpt) Help() string {
	return "view this help message"
}

// ShortFlag returns the short-form flag for this option
func (r *HelpOpt) ShortFlag() rune {
	return 'h'
}

// ---- Setter

// Set sets this option's value
func (r *HelpOpt) Set(vv interface{}) error {
	if v, ok := vv.(bool); ok {
		*r = HelpOpt(v)
		return nil
	}
	return fmt.Errorf("%w: help to %+v", ErrSet, vv)
}
