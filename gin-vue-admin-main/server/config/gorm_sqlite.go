package config

import (
	"path/filepath"
)

type Sqlite struct {
	GeneralDB `yaml:",inline" mapstructure:",squash"`
}

func (s *Sqlite) Dsn() string {
	if len(s.Path) > 3 && s.Path[len(s.Path)-3:] == ".db" {
		return s.Path
	}
	return filepath.Join(s.Path, s.Dbname+".db")
}
