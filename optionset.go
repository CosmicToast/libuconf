package libuconf

import (
	"fmt"
	"os"
	"strings"
)

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

// VisitEnv visits EnvOpts
func (o *OptionSet) VisitEnv(f func(EnvOpt)) {
	o.Visit(func(vv Setter) {
		if v, ok := vv.(EnvOpt); ok {
			f(v)
		}
	})
}

// VisitFlag visits FlagOpts
func (o *OptionSet) VisitFlag(f func(FlagOpt)) {
	o.Visit(func(vv Setter) {
		if v, ok := vv.(FlagOpt); ok {
			f(v)
		}
	})
}

// VisitToml visits TomlOpts
func (o *OptionSet) VisitToml(f func(TomlOpt)) {
	o.Visit(func(vv Setter) {
		if v, ok := vv.(TomlOpt); ok {
			f(v)
		}
	})
}

// FindLongFlag returns a FlagOpt based on the "Flag" property
func (o *OptionSet) FindLongFlag(s string) (res FlagOpt) {
	o.VisitFlag(func(f FlagOpt) {
		if f.Flag() == s {
			res = f
			return
		}
	})
	return
}

// FindShortFlag returns a FlagOpt based on the "ShortFlag" property
func (o *OptionSet) FindShortFlag(c rune) (res FlagOpt) {
	if c == 0 {
		return nil
	}
	o.VisitFlag(func(f FlagOpt) {
		if f.ShortFlag() == c {
			res = f
			return
		}
	})
	return
}

// wrapper for fmt.Fprintf(os.Stderr, ...)
// this is here in case I want to make this configurable later
func (o *OptionSet) print(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(os.Stderr, format, a...)
}

// print a single flag
// "  "
// --flag | -f -> flag part
// buffer " " up to flen
// Help up to llen-1 + "-" if continuation
// flen -> continue help
func (o *OptionSet) printflag(flen, llen int, f FlagOpt) string {
	var s strings.Builder
	s.WriteString("  --")
	s.WriteString(f.Flag())
	if sf := f.ShortFlag(); sf != 0 {
		s.WriteString(" | -")
		s.WriteRune(sf)
	}

	for i := flen - s.Len(); i > 0; i-- {
		s.WriteRune(' ') // align
	}
	s.WriteRune(' ')

	hl := llen - flen - 1
	h := f.Help()

	// add value to help, if it's there
	if v, ok := f.(Getter); ok {
		h += fmt.Sprintf(" (%v)", v.Get())
	}

	for { // write help
		if len(h) <= hl { // it all fits
			s.WriteString(h)
			break
		} // it doesn't fit
		if h[hl] == ' ' || h[hl] == '\n' { // we're at a word-boundary
			s.WriteString(h[:hl])
			s.WriteRune('\n')
			h = h[hl+1:] // skip the boundary
			continue
		}
		s.WriteString(h[:hl-1]) // leave 1 for '-'
		s.WriteRune('-')
		s.WriteRune('\n')
		h = h[hl:] // do not skip anything
	}
	s.WriteRune('\n')

	return s.String()
}

// Usage will use AppName and the Options that are FlagOpts to print usage info
// The format looks like this:
/*
	AppName:
		--flag | -f Really long help m-
							  essage. (value)
		--flag      Still here (value).
		--flag | -f Also still here.
*/
// TODO: long help messages are not currently broken down.
func (o *OptionSet) Usage() {
	o.print("%s:\n", o.AppName)
	flen := 0
	llen := 80 // TODO: don't default to 80

	o.VisitFlag(func(v FlagOpt) {
		vlen := len(v.Flag()) + 4 // "  --" -> 4
		if v.ShortFlag() != 0 {
			vlen += 5 // " | -x"
		}

		if vlen > flen {
			flen = vlen
		}
	})

	o.VisitFlag(func(f FlagOpt) {
		o.print(o.printflag(flen, llen, f))
	})
}
