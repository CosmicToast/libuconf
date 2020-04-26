package libuconf

/*
Parse is a wrapper around all the built-in parse methods.

First, it will parse all of the standard files for your platform.
It will ignore "file missing" errors, but not any others (e.g setting errors).
If it encounters an error in the middle of a file, it will finish parsing the
file, but all other files will not be parsed.
You will only receive the *last* error in the *first* file that fails.

After that (assuming no errors), it will try to parse your environment.
If it encounters an error in your envirnoment, it will still parse the rest.
You will only receive the error for the *last* option that failed to be set.

Finally (assuming no errors), it will try to parse your flags (args).
See ParseFlags([]string) for more details as to how and when it might error out.

Once all parsing is done, it will check to see if you have a Help flag defined.
If you do, and it's set to true, or if there was an error, Parse() will call
Usage() for you.

If you want any other behavior, please write your own handling!
This is valid and encouraged.
I.e this function can also serve as an example as to how to set up your own.
*/
func (o *OptionSet) Parse(ss []string) error {
	var (
		help bool
		err  error
	)

	err = o.ParseStdToml()
	if err != nil {
		goto parse_finish
	}

	err = o.ParseEnv()
	if err != nil {
		goto parse_finish
	}

	err = o.ParseFlags(ss)
	if err != nil {
		goto parse_finish
	}

	// we don't need to try to find "help" if there is an active error
	o.VisitFlag(func(vv FlagOpt) {
		if v, ok := vv.(*HelpOpt); ok {
			help = help || bool(*v)
		}
	})

parse_finish:
	if err != nil || help {
		o.Usage()
	}
	return err
}
