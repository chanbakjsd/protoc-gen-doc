package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kenshaw/ini"
	"github.com/kenshaw/ini/parser"
)

// LoadFile loads the configuration at the specified path.
func LoadFile(file string) (*Config, error) {
	dir, _ := filepath.Split(file)
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Load(dir, f)
}

// Load loads the configuration, assuming that paths are relative to the
// provided folderPath.
func Load(folderPath string, r io.Reader) (*Config, error) {
	f, err := ini.Load(r)
	if err != nil {
		return nil, err
	}
	cfg := &Config{
		Sections: make(map[string]Section),
	}
	for _, s := range f.AllSections() {
		if s.Name() == "" {
			continue
		}
		sect, err := loadSection(folderPath, s)
		if err != nil {
			return nil, err
		}
		cfg.Sections[s.Name()] = sect
	}
	return cfg, nil
}

// loadSection loads the configuration section.
func loadSection(folderPath string, s *parser.Section) (Section, error) {
	sect := Section{}
	for _, k := range s.Keys() {
		v := s.Get(k)
		switch k {
		default:
			return Section{}, fmt.Errorf("unknown key %q in section %q", k, s.Name())
		case "name":
			sect.DisplayName = v
		case "packages":
			if s.Name() == "default" {
				return Section{}, fmt.Errorf("default section may not specify packages")
			}
			sect.Packages = strings.Split(v, ",")
		case "preamble":
			f, err := os.Open(filepath.Join(folderPath, v))
			if err != nil {
				return Section{}, fmt.Errorf("cannot open preamble file: %w", err)
			}
			content, err := io.ReadAll(f)
			if err != nil {
				return Section{}, fmt.Errorf("cannot read preamble file: %w", err)
			}
			sect.PreambleContent = string(content)
		case "weight":
			var err error
			sect.Weight, err = strconv.Atoi(v)
			if err != nil {
				return Section{}, fmt.Errorf("weight not a number in section %q", s.Name())
			}
		}
	}
	return sect, nil
}
