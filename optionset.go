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

/* print a single flag
"  " ->
if there is a short flag at all, "-f, " for short flags, or padding w/o ->
--flag ->
buffer " " up to flen
Help() up to llen-1 + '-' if it's a non-word-boundary
flen -> continue Help()
*/
func (o *OptionSet) printflag(flen, llen int, short bool, f FlagOpt) string {
	var s strings.Builder
	s.WriteString("  ")
	if short {
		if sf := f.ShortFlag(); sf != 0 {
			s.WriteRune('-')
			s.WriteRune(sf)
			s.WriteString(", ")
		} else {
			s.WriteString("    ") // - f , ' '
		}
	}
	s.WriteString("--")
	s.WriteString(f.Flag())

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

/* Usage will use AppName and the Options that are FlagOpts to print usage info
The format looks like this if there are no short flags:
  AppName:
    --flag Really long help
           message.
    --ofl  Other really lon-
           g help message.
    --fl   Help message.

Or like this if there is at least one:
  AppName:
    -f, --flag Really long help
               message.
    -o, --of   Other really lon-
               g help message.
        --st   Standalone flag.
*/
func (o *OptionSet) Usage() {
	o.print("%s:\n", o.AppName)
	flen := 0
	llen := 80 // TODO: don't default to 80
	short := false

	o.VisitFlag(func(v FlagOpt) {
		vlen := len(v.Flag()) + 4 // "  --" -> 4
		if v.ShortFlag() != 0 {
			if !short {
				flen += 4 // - f , ' '; we didn't know we had short flags until now
			}
			short = true
		}
		if short {
			vlen += 4
		}

		if vlen > flen {
			flen = vlen
		}
	})

	o.VisitFlag(func(f FlagOpt) {
		o.print(o.printflag(flen, llen, short, f))
	})

	if len(o.Args) > 0 {
		o.print("Args: %q\n", o.Args)
	}
}
