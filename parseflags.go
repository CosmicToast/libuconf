package libuconf

import (
	"errors"
	"fmt"
)

// Errors
var (
	ErrNoVal = errors.New("missing value")
	ErrSet   = errors.New("failed to set")
)

func errNoVal(name string) error {
	return fmt.Errorf("%w: %s", ErrNoVal, name)
}
func errSet(name, val string) error {
	return fmt.Errorf("%w: %s to %s", ErrSet, name, val)
}

/* ParseFlags parses a set of args (ss) and modifies the attached options
These formats are allowed:
 -abc value (where a and b are bool types)
 -dvalue
 --a.b.c=value
 --d.e.f value
 --g.h.i (where g.h.i is a bool type)
notably missing is -a=value, but you shouldn't need it

Note that non-registered flags are valid.
For example, if you did not register a "foobar" flag, "--foobar=y" will be added
to Args.
*/
func (o *OptionSet) ParseFlags(ss []string) error {
	var (
		done, handled bool
		handling      FlagOpt
		name          string // name on handling
	)

	for _, s := range ss {
		if done { // we finished parsing args, so just pass the rest through
			o.Args = append(o.Args, s)
			continue
		}

		var (
			opt, val1  = o.getLongFlag(s)
			opts, val2 = o.getShortFlags(s)
		)

		// we never did the last flag!
		if (opt != nil || opts != nil) && handling != nil && !handled {
			if !handling.Bool() { // needed a value, but we have a flag
				return errNoVal(name)
			}
			if err := handling.Set(true); err != nil { // can't set to true
				return errSet(name, "true")
			}
		}

		// long flag
		if opt != nil {
			handling = opt
			handled = false
			name = handling.Flag()

			if len(val1) > 0 { // it came with a value!
				if err := opt.Set(val1); err != nil {
					return errSet(name, val1)
				}
				handled = true
			}
			continue
		}

		// short flag
		if l := len(opts); l > 0 {
			for i, v := range opts { // range of nil = skip
				handling = v
				handled = false
				name = string(handling.ShortFlag())

				if i != l-1 { // NOT the last opt - set bool to true, guaranteed by func
					if err := v.Set(true); err != nil {
						return errSet(name, "true")
					}
					continue
				} // it IS the last opt, check for value
				if val2 != "" { // there is one!
					handled = true
					if err := v.Set(val2); err != nil {
						return errSet(name, val2)
					}
				}
			}
			continue
		}

		// are we done handling flags?
		if len(s) == 2 && s[0] == '-' && s[1] == '-' {
			done = true
			continue
		}

		// ok, so it *must* be a value or an arg
		if handled || handling == nil {
			o.Args = append(o.Args, s) // this is easy
			continue
		}

		// so the previous flag isn't handled, but we have a value
		// so if the value fails to apply, bools should just be set
		err := handling.Set(s)             // try to set
		if err != nil && handling.Bool() { // it failed, but the flag is bool
			o.Args = append(o.Args, s) // save s before we override it
			err = handling.Set(true)   // try to set again, resetting err
			s = "true"                 // set s to true for errSet
		}
		if err != nil { // if there's still an error, error out
			return errSet(name, s)
		}
		handled = true // there was no error
	}

	// the last string was a flag, and we didn't handle it
	if !handled && handling != nil {
		if handling.Bool() { // if it's a boolean, set it to true
			if err := handling.Set(true); err != nil { // can't set to true
				return errSet(name, "true")
			}
		} else {
			return errNoVal(name)
		}
	}

	return nil
}

func (o *OptionSet) getLongFlag(s string) (FlagOpt, string) {
	if len(s) < 3 || s[0] != '-' || s[1] != '-' {
		// not a long flag
		return nil, ""
	}
	var (
		flag    []rune
		val     []rune
		foundeq bool
	)
	for i, v := range s {
		if i < 2 { // skip "--"
			continue
		}
		if !foundeq && v == '=' {
			foundeq = true
			continue
		}
		if foundeq {
			val = append(val, v)
		} else {
			flag = append(flag, v)
		}
	}
	opt := o.FindLongFlag(string(flag))
	return opt, string(val)
}

func (o *OptionSet) getShortFlags(s string) ([]FlagOpt, string) {
	if len(s) < 2 || s[0] != '-' || s[1] == '-' {
		// not a short flag
		return nil, ""
	}
	var (
		res  []FlagOpt
		val  []rune
		done bool // the rest of the string is the val
	)
	for i, v := range s {
		if i == 0 { // skip "-"
			continue
		}
		if done { // could move this into the if, but meh
			val = append(val, v)
			continue
		}
		opt := o.FindShortFlag(v)
		// no flag, or last flag reqs value
		if opt == nil || (res != nil && !res[len(res)-1].Bool()) {
			done = true
			val = append(val, v)
			continue
		} else { // we have a flag, and we can add it
			res = append(res, opt)
		}
	}
	out := ""
	if val != nil {
		out = string(val)
	}
	return res, out
}
