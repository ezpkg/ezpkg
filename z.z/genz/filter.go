package ggen

import (
	"golang.org/x/tools/go/packages"
)

type FilteringPackage struct {
	PkgPath          string
	Imports          map[string]*packages.Package
	Directives       Directives
	InlineDirectives Directives

	ng *filterEngine
}

func (p *FilteringPackage) Include() {
	p.ng.IncludePackage(p.PkgPath)
}

type filteringPackage struct {
	PkgPath          string
	Imports          map[string]*packages.Package
	Directives       Directives
	InlineDirectives Directives
}

type includedPackage struct {
	PkgPath  string
	Included []bool
}

type FilterEngine interface {

	// Plugin should use the embedded logger to log messages.
	Logger

	// IncludePackage indicates that the given package will be included for generating. It will be returned later in
	// Engine.GeneratingPackages(). If it does not exist, an error with be returned later.
	IncludePackage(pkgPath string)

	// ParsePackage indicates that the given package should be parsed. If it does not exist, an error with be returned
	// later.
	//
	// Sometimes, the plugin depends on some specific package that are not transitive imported by IncludePackage. In
	// this case, the plugin can call ParsePackage to parse these packages. This should happen sparsely and not all
	// plugins need to call this function.
	ParsePackage(pkgPath string)

	// ParsePackages takes a pattern and add the matched packages to ParsingPackages. Pattern is the go package pattern:
	//
	//     example.com/...
	//     github.com/path/...
	//
	// ParsePackages indicates that the given package pattern should be parsed. If the pattern does not match any
	// package, an error with be returned later.
	//
	// Sometimes, the plugin depends on some specific package that are not transitive imported by IncludePackage. In
	// this case, the plugin can call ParsePackage to parse these packages. This should happen sparsely and not all
	// plugins need to call this function.
	//
	// A program may require many packages. So it's best to only include packages related to your work. For example if
	// all your packages are put under github.com/yourcompany, it's recommended to call
	// ParsePackages("github.com/yourcompany/...") to avoid parsing unnecessary packages. Or if your plugin only care
	// about protobuf files, you can call ParsePackages("github.com/yourcompany/models/protobuf/...") to include only
	// protobuf generated packages.
	ParsePackages(patterns ...string)

	// ParsingPackages returns a list of packages that are processing. The plugin can can use
	// FilteringPackage.Directives or FilteringPackage.Imports to determine which packages should be available to
	// Generate, then calls IncludePackage on those packages.
	ParsingPackages() []*FilteringPackage
}

var _ FilterEngine = &filterEngine{}

type filterEngine struct {
	zLogger

	ng       *engine
	plugin   *pluginStruct
	pkgs     []filteringPackage
	pkgMap   map[string][]bool
	patterns *[]string
}

func (ng *filterEngine) Logger() Logger {
	return ng.zLogger
}

// IncludePackage indicates that the given package will be included for generating. It will be returned later in
// Engine.GeneratingPackages(). If it does not exist, an error with be returned later.
func (ng *filterEngine) IncludePackage(pkgPath string) {
	if pkgPath == "" {
		panic("invalid package path")
	}
	flags := ng.pkgMap[pkgPath]
	if flags == nil {
		flags = make([]bool, len(ng.ng.plugins))
		ng.pkgMap[pkgPath] = flags
	}
	flags[ng.plugin.index] = true
}

// ParsePackage indicates that the given package should be parsed. If it does not exist, an error with be returned
// later. The parsed packages are returned by calling ParsingPackages.
func (ng *filterEngine) ParsePackage(pkgPath string) {
	if pkgPath == "" {
		panic("invalid package path")
	}
	_, ok := ng.pkgMap[pkgPath]
	if !ok {
		ng.pkgMap[pkgPath] = nil
	}
}

// ParsePackages takes a pattern and add the matched packages to ParsingPackages. Pattern is the go package pattern:
//
//	example.com/...
//	github.com/path/...
//
// A program may require many packages. So it's best to only include packages related to your work. For example if all
// your packages are put under github.com/yourcompany, it's recommended to call
// ParsePackages("github.com/yourcompany/...") to avoid parsing unnecessary packages. Or if your plugin are only care
// about protobuf files, you can call ParsePackages("github.com/yourcompany/models/protobuf/...") to include only
// protobuf generated packages.
func (ng *filterEngine) ParsePackages(patterns ...string) {
	for _, p := range patterns {
		if p == "" {
			continue
		}
		*ng.patterns = append(*ng.patterns, p)
	}
}

// ParsingPackages returns a list of packages provided to ParsePackage(s). The plugin can use
// FilteringPackage.Directives or FilteringPackage.Imports to determine which packages should be available to
// Generate, then calls IncludePackage on those packages.
func (ng *filterEngine) ParsingPackages() []*FilteringPackage {
	buf := make([]FilteringPackage, len(ng.pkgs))
	res := make([]*FilteringPackage, len(ng.pkgs))
	for i, pkg := range ng.pkgs {
		p := FilteringPackage{
			PkgPath:          pkg.PkgPath,
			Imports:          pkg.Imports,
			Directives:       cloneDirectives(pkg.Directives),
			InlineDirectives: cloneDirectives(pkg.InlineDirectives),

			ng: ng,
		}
		buf[i] = p
		res[i] = &buf[i]
	}
	return res
}
