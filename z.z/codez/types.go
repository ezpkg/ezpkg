package codez

import (
	"golang.org/x/tools/go/packages"
)

type Config struct {
	packages.Config
}

func Cfg() *Config {
	return &Config{
		Config: packages.Config{
			Mode: packages.NeedName |
				packages.NeedFiles |
				packages.NeedCompiledGoFiles |
				packages.NeedImports |
				packages.NeedDeps |
				packages.NeedExportFile |
				packages.NeedTypes |
				packages.NeedSyntax |
				packages.NeedTypesInfo |
				packages.NeedTypesSizes |
				packages.NeedModule |
				packages.NeedEmbedFiles |
				packages.NeedEmbedPatterns,
		},
	}
}

func (c *Config) Unwrap() *packages.Config {
	return &c.Config
}
