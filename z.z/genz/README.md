# ggen - Code generation for Go

## Install

```
go get ezpkg.io/genz
```

## Usage

This library helps write your own code generator script.

### Plugin Structure

Plugin must implement the `ggen.Plugin` interface:

```go
type Plugin interface {

	// Name returns name of the plugin. Each plugin must have a different name.
	Name() string

	// Filter is called to determine which packages will be parsed and which will be skipped. It will be called before
	// Generate. It received a FilterEngine and need to call the following methods:
	//
	//     ParsingPackages: Get a list of all packages available for IncludePackage.
	//                      The plugin can only include packages in this list.
	//     IncludePackage:  Make the package available for Generate.
	Filter(FilterEngine) error

	// Generate is called to actually generate code for the given packages. Only packages passed to
	// FilterEngine.IncludePackage are available for Generate.
	Generate(Engine) error
}
```

## LICENSE

[MIT](https://ezpkg.io/genz/blob/master/LICENSE)
