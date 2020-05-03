/*
Package libuconf is a small library to handle all sorts of configuration tasks.
More specifically, it is an extendable scaffolding system for handling
configuration in a unified way, with built-ins to be useful without extending.

Scaffolding

At the core of libuconf is the OptionSet and the interfaces it consumes.
Setter is the least common denominator: allowing setting values to options.
Beyond that, there are interfaces for each Parse function.
For example, the EnvOpt interface provides everything needed by ParseEnv.

Regardless of how many extensions are plugged into libuconf, all you need to do
to utilize it is plug a value that implements the needed interfaces into
OptionSet.Var(), which will make it available to the Parse family of functions.

Built-In Systems

Libuconf comes with three built-in systems for configuration handling.
EnvOpt, used with ParseEnv,
FlagOpt, used with ParseFlags, and
TomlOpt, used with ParseTomlFile.

Because of how libuconf works, utilizing them for custom values is simple:
all you need to do is implement the appropriate interface.
For an example, you can look at any of the built-in types.

Built-In Types

Libuconf comes with several types that implement the needed interfaces for the
built-in systems.
These types are BoolOpt (a boolean option), FloatOpt (a float64 option),
HelpOpt (a special help option, only implementing FlagOpt),
IntOpt (an int64 option), StringOpt (a string option) and
UintOpt (a uint64 option).

All of these built-in types have special integration functions within
OptionSet.
For example, BoolOpt has OptionSet.Bool and OptionSet.BoolVar.
These integration functions, however, are entirely optional.
It is entirely valid to create a BoolOpt (or any Opt) by hand, and call
OptionSet.Var against it.
In fact, this is what half of these integration functions do.

Usage

First, instantiate an OptionSet using NewOptionSet.
The argument should be the name of your application.
Note that this name is significant!
If you use os.Args[0], users will lose their configuration file access if they
rename the binary, because AppName is used in ParseStdFiles and Usage.
As such, it is recommended to use a constant value for the AppName.

Then, add any options you want.
If you are using your own Option types, you should add them using
OptionSet.Var.
Otherwise, you may take advantage of the integration functions.
Note that with the built-in types, the environment variables and configuration
file keys searched are derived from the long flag.
You can always pass 0 as the shortname to disable the short flag for a flag.

If you want to limit the applicability of some options, you will need to extend
libuconf.

Extending - Option Types

To extend libuconf's available options (or de-extend them, restricting parsing)
all you need to do is implement the interfaces you want to be used.
Feel free to take inspiration from the built-in types.
If you want to restrict some string option to only be available in a
configuration file, for example, you can copy string.go and remove the
functions related to FlagOpt and EnvOpt.
Note that you are fully free to do this with no licensing issues, because
libuconf is released under the unlicense.

Extending - Systems

To extend libuconf's systems (for example, add support for setting options
through consul), there is a bit more work to do.

You will need to subclass OptionSet (or fork libuconf and just add to it).
Then, you will need to create a new interface (in this example, ConsulOpt) and
a new OptionSet.ParseConsul function.
The ParseConsul function should only operate on ConsulOpt members of OptionSet
(see ParseEnv for a succint example, implementing a VisitConsul is encouraged)
and should only use functions available in ConsulOpt to handle the parsing.
This means that ConsulOpt interface should contain everything necessary to
perform the work.
*/
package libuconf