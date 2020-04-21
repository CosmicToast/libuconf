package libuconf

// OptionSet represents a set of options
type OptionSet struct {
	AppName string
	Options []Setter // least common denominator
	Args    []string
}

// Var adds an option to the OptionSet
// It is expected to be pre-set to its default value
func (o *OptionSet) Var(opt Setter) {
	o.Options = append(o.Options, opt)
}

// Visit visits each option, passing it to the argument function
func (o *OptionSet) Visit(f func(Setter)) {
	for _, opt := range o.Options {
		f(opt)
	}
}

// FindLongFlag returns a FlagOpt based on the "Flag" property
func (o *OptionSet) FindLongFlag(s string) (res FlagOpt) {
	o.Visit(func(opt Setter) {
		if oo, ok := opt.(FlagOpt); ok && oo.Flag() == s {
			res = oo
		}
	})
	return
}

// FindShortFlag returns a FlagOpt based on the "ShortFlag" property
func (o *OptionSet) FindShortFlag(c rune) (res FlagOpt) {
	if c == 0 {
		return nil
	}
	o.Visit(func(opt Setter) {
		if oo, ok := opt.(FlagOpt); ok && oo.ShortFlag() == c {
			res = oo
		}
	})
	return
}
