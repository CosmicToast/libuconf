package libuconf

/*
In this file we define various OptionSet integrations for the built-in Options.
The Type() functions are glorified wrappers around NewTypeOpt().
The TypeVar() functions are special handling, meant to emulate go's flag lib.

This will likely all be generated later.
*/

// ---- bool

// Bool adds a BoolOpt to the OptionSet.
// It returns a pointer to the output value.
func (o *OptionSet) Bool(name string, sname rune, val bool, help string) *bool {
	op, v := NewBoolOpt(name, sname, val, help)
	o.Var(op)
	return v
}

// BoolVar adds a BoolOpt to the OptionSet.
func (o *OptionSet) BoolVar(out *bool, name string, sname rune, val bool, help string) {
	*out = val
	op := &BoolOpt{help, name, sname, out}
	o.Var(op)
}

// ---- float

// Float adds a FloatOpt to the OptionSet.
// It returns a pointer to the output value.
func (o *OptionSet) Float(name string, sname rune, val float64, help string) *float64 {
	op, v := NewFloatOpt(name, sname, val, help)
	o.Var(op)
	return v
}

// FloatVar adds a FloatOpt to the OptionSet.
func (o *OptionSet) FloatVar(out *float64, name string, sname rune, val float64, help string) {
	*out = val
	op := &FloatOpt{help, name, sname, out}
	o.Var(op)
}

// ---- help

// Help adds --help and -h flags to the OptionSet.
func (o *OptionSet) Help() *bool {
	out := new(bool)
	o.Var((*HelpOpt)(out))
	return out
}

// ---- int

// Int adds a IntOpt to the OptionSet.
// It returns a pointer to the output value.
func (o *OptionSet) Int(name string, sname rune, val int64, help string) *int64 {
	op, v := NewIntOpt(name, sname, val, help)
	o.Var(op)
	return v
}

// IntVar adds a IntOpt to the OptionSet.
func (o *OptionSet) IntVar(out *int64, name string, sname rune, val int64, help string) {
	*out = val
	op := &IntOpt{help, name, sname, out}
	o.Var(op)
}

// ---- string

// String adds a StringOpt to the OptionSet.
// It returns a pointer to the output value.
func (o *OptionSet) String(name string, sname rune, val string, help string) *string {
	op, v := NewStringOpt(name, sname, val, help)
	o.Var(op)
	return v
}

// StringVar adds a StringOpt to the OptionSet.
func (o *OptionSet) StringVar(out *string, name string, sname rune, val string, help string) {
	*out = val
	op := &StringOpt{help, name, sname, out}
	o.Var(op)
}

// ---- uint

// Uint adds a UintOpt to the OptionSet.
// It returns a pointer to the output value.
func (o *OptionSet) Uint(name string, sname rune, val uint64, help string) *uint64 {
	op, v := NewUintOpt(name, sname, val, help)
	o.Var(op)
	return v
}

// UintVar adds a UintOpt to the OptionSet.
func (o *OptionSet) UintVar(out *uint64, name string, sname rune, val uint64, help string) {
	*out = val
	op := &UintOpt{help, name, sname, out}
	o.Var(op)
}
