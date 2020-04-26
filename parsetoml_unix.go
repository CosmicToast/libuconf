// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package libuconf

import (
	"os"
	"path"
)

// ParseStdToml parses the standard configuration files for a platform, in order.
func (o *OptionSet) ParseStdToml() error {
	// dirs
	uhome, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	uconf, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	var (
		// filenames
		tfile = o.AppName + ".toml"
		trc   = o.AppName + "rc"
		thid  = "." + trc
	)

	/* /etc/app.toml        ->
	   /etc/app/config.toml ->
	   ~/.tomlrc            ->
	   ~/.config/app.toml   ->
	   ~/.config/app/config.toml
	*/
	return o.ParseTomlFiles(
		path.Join("/etc", tfile),
		path.Join("/etc", o.AppName, "config.toml"),
		path.Join(uhome, thid),
		path.Join(uconf, tfile),
		path.Join(uconf, o.AppName, "config.toml"),
	)
}
