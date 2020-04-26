package libuconf

import (
	"os"

	"github.com/pelletier/go-toml"
)

// ParseTomlFile parses a toml file, looking for options that implement TomlOpt
// If there is an error, it is returned, even if it's just a missing file
func (o *OptionSet) ParseTomlFile(path string) error {
	tree, err := toml.LoadFile(path)
	if err != nil {
		return err
	}
	return o.parseTomlTree(tree)
}

// ParseTomlFiles is a convenience function to run multiple ParseTomlFile s
// It also ignores any missing files, unlike ParseTomlFile
func (o *OptionSet) ParseTomlFiles(paths ...string) error {
	for _, v := range paths {
		if err := o.ParseTomlFile(v); err != nil {
			if os.IsNotExist(err) { // ignore missing files
				continue
			}
			return err
		}
	}
	return nil
}

// ParseTomlString parses a string containing TOML data
func (o *OptionSet) ParseTomlString(c string) error {
	tree, err := toml.Load(c)
	if err != nil {
		return err
	}
	return o.parseTomlTree(tree)
}

func (o *OptionSet) parseTomlTree(t *toml.Tree) (e error) {
	o.VisitToml(func(v TomlOpt) {
		out := t.Get(v.Toml())
		if out == nil {
			return
		}
		if err := v.Set(out); err != nil {
			e = err
		}
	})
	return
}
