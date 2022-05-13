package generate

import (
	"fmt"
	"strings"

	"github.com/chanbakjsd/protoc-gen-doc/config"
	"github.com/chanbakjsd/protoc-gen-doc/doc"
)

// Tags combines the packages provided into tags based on the specified config.
// Packages that are not specified will be placed in the `default` tag.
func Tags(cfg *config.Config, pkgs []*doc.Package) (map[string]*doc.Tag, error) {
	usedPkgs := make(map[string]bool)
	tags := make(map[string]*doc.Tag, len(cfg.Sections)+1)
	var defaultSection config.Section
	for tagName, sect := range cfg.Sections {
		if tagName == "default" {
			// default is handled after everything.
			defaultSection = sect
			continue
		}
		tag := &doc.Tag{
			Name:     sect.DisplayName,
			Preamble: sect.PreambleContent,
			Weight:   sect.Weight,
			Packages: make([]*doc.Package, 0, len(sect.Packages)),
		}
		for _, pkgName := range sect.Packages {
			pkg, err := findPkg(pkgs, pkgName, tagName)
			if err != nil {
				return nil, err
			}
			usedPkgs[pkg.ID] = true
			tag.Packages = append(tag.Packages, pkg)
		}
		tags[tagName] = tag
	}

	defaultPkgs := make([]*doc.Package, 0)
	for _, pkg := range pkgs {
		if !usedPkgs[pkg.ID] {
			defaultPkgs = append(defaultPkgs, pkg)
		}
	}
	tags["default"] = &doc.Tag{
		Name:     defaultSection.DisplayName,
		Preamble: defaultSection.PreambleContent,
		Weight:   defaultSection.Weight,
		Packages: defaultPkgs,
	}
	return tags, nil
}

func findPkg(pkgs []*doc.Package, pkgName, tagName string) (*doc.Package, error) {
	// Priority:
	// - Match on exact ID
	// - Match on exact name
	// - Match on trailing ID
	for _, pkg := range pkgs {
		if pkg.ID == pkgName {
			return pkg, nil
		}
	}
	for _, pkg := range pkgs {
		if pkg.Name == pkgName {
			return pkg, nil
		}
	}
	for _, pkg := range pkgs {
		if strings.HasSuffix(pkg.ID, "."+pkgName) {
			return pkg, nil
		}
	}
	return nil, fmt.Errorf("package %q in section %q not found", pkgName, tagName)

}
