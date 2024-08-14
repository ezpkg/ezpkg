package ggen

import (
	"context"
	"os"

	"ezpkg.io/logz"
)

type GenerateFileNameInput struct {
	PluginName string
}

type Config struct {
	Plugins []Plugin

	// Map of enabled plugins. Leave this nil to enable all plugins.
	EnabledPlugins map[string]bool

	// default to "zz_generated.{{.Name}}.go"
	GenerateFileName func(GenerateFileNameInput) string

	CleanOnly bool

	Namespace string

	GoimportsArgs []string

	BuildTags []string

	Logger Logger
}

func (c *Config) RegisterPlugin(plugins ...Plugin) {
	c.Plugins = append(c.Plugins, plugins...)
}

func (c *Config) EnablePlugin(names ...string) {
	if c.EnabledPlugins == nil {
		c.EnabledPlugins = map[string]bool{}
	}
	for _, name := range names {
		c.EnabledPlugins[name] = true
	}
}

func Start(ctx context.Context, cfg Config, patterns ...string) error {
	logger := cfg.Logger
	if logger == nil {
		logger = logz.DefaultLogger(os.Stderr)
	}

	ng := newEngine(cfg.Logger)
	return ng.start(ctx, cfg, patterns...)
}
