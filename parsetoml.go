package libuconf

import "github.com/pelletier/go-toml"

// ParseTomlFile parses a toml file, looking for options that implement TomlOpt
func (o *OptionSet) ParseTomlFile(path string) error {
	tree, err := toml.LoadFile(path)
	if err != nil {
		return err
	}
	return o.parseTomlTree(tree)
}

// ParseTomlFiles is a convenience function to run multiple ParseTomlFile s
func (o *OptionSet) ParseTomlFiles(paths ...string) error {
	for _, v := range paths {
		if err := o.ParseTomlFile(v); err != nil {
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
		if err := v.Set(out); err != nil {
			e = err
			return
		}
	})
	return
}
