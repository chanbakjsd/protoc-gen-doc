package doc

import (
	"fmt"
	"os"
	"strings"
)

// PruneTypes prunes unused types from the package by walking all endpoints.
// Types that are not used in requests or responses are removed.
func PruneTypes(pkgs []*Package) {
	usedTypes := make(map[string]bool)
	for _, pkg := range pkgs {
		for _, srv := range pkg.Services {
			for _, endpoint := range srv.Endpoints {
				markUsedTypes(pkgs, usedTypes, endpoint.Request, true)
				markUsedTypes(pkgs, usedTypes, endpoint.Response, true)
			}
		}
	}
	for _, pkg := range pkgs {
		newTyp := make(map[string]Type)
		for k, v := range pkg.Types {
			if usedTypes[k] {
				newTyp[k] = v
			}
		}
		pkg.Types = newTyp
	}
}

// markUsedTypes marks the specified types and all types it refers to as used in the usedTypes map.
func markUsedTypes(pkgs []*Package, usedTypes map[string]bool, t Type, topLevel bool) {
	if !topLevel {
		if t, ok := t.(NamedType); ok {
			usedTypes[t.TypeName()] = true
		}
	}
	switch t := t.(type) {
	default:
		panic(fmt.Sprintf("unknown type: %T", t))
	case *Enum, *Basic:
	case *Array:
		markUsedTypes(pkgs, usedTypes, t.Value, false)
	case *Map:
		markUsedTypes(pkgs, usedTypes, t.Key, false)
		markUsedTypes(pkgs, usedTypes, t.Value, false)
	case *Ref:
		typ := resolveRef(pkgs, t)
		markUsedTypes(pkgs, usedTypes, typ, false)
	case *Message:
		for _, f := range t.Fields {
			markUsedTypes(pkgs, usedTypes, f.Type, false)
		}
	}
}

// resolveRef resolves the reference type given the packages. It panics if the
// reference is not resolved.
func resolveRef(pkgs []*Package, ref *Ref) Type {
	for _, pkg := range pkgs {
		prefix := pkg.ID + "."
		if !strings.HasPrefix(ref.Name, prefix) {
			continue
		}
		name := strings.TrimPrefix(ref.Name, prefix)
		typ, ok := pkg.Types[name]
		if ok {
			return typ
		}
	}
	for _, pkg := range pkgs {
		prefix := pkg.ID + "."
		if !strings.HasPrefix(ref.Name, prefix) {
			continue
		}
		fmt.Fprintln(os.Stderr, pkg.ID+" GOOD")
		keys := make([]string, 0)
		for k := range pkg.Types {
			keys = append(keys, k)
		}
		fmt.Fprintln(os.Stderr, "Keys:", keys)
	}
	panic("ref type not found: " + ref.Name)
}
