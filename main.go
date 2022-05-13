package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/chanbakjsd/protoc-gen-doc/config"
	"github.com/chanbakjsd/protoc-gen-doc/doc"
	"github.com/chanbakjsd/protoc-gen-doc/generate"
	"github.com/chanbakjsd/protoc-gen-doc/proto"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	protogen.Options{}.Run(run)
}

func run(p *protogen.Plugin) error {
	cfg, err := readConfig(p.Request.GetParameter())
	if err != nil {
		return err
	}
	sort.Slice(p.Files, func(i int, j int) bool {
		return p.Files[i].Proto.GetName() < p.Files[j].Proto.GetName()
	})
	pkgs := make([]*doc.Package, 0, len(p.Files))
	for _, f := range p.Files {
		pkg := proto.ConvertFile(f)
		pkgs = append(pkgs, pkg)
	}
	doc.PruneTypes(pkgs)
	genPkgs := make([]*doc.Package, 0, len(pkgs))
	for i, f := range p.Files {
		if f.Generate {
			genPkgs = append(genPkgs, pkgs[i])
		}
	}
	tags, err := generate.Tags(cfg, genPkgs)
	if err != nil {
		return err
	}
	for name, tag := range tags {
		f := p.NewGeneratedFile(name+".json", "")
		if err := json.NewEncoder(f).Encode(tag); err != nil {
			return err
		}
	}
	return nil
}

func readConfig(param string) (*config.Config, error) {
	if !strings.HasPrefix(param, "config=") {
		return nil, fmt.Errorf("config location not provided, pass param in the form of 'config=...'")
	}
	file := strings.TrimPrefix(param, "config=")
	return config.LoadFile(file)
}
