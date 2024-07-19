package ggen

import (
	"os"

	"github.com/iolivernguyen/ggen/ggen/logging"
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

	LogLevel   LogLevel
	LogHandler LogHandler
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

func (c *Config) defaultLogHandler() LogHandler {
	handler := defaultLogHandler{
		w:     os.Stderr,
		level: c.LogLevel,
	}
	return handler
}

func Start(cfg Config, patterns ...string) error {
	if cfg.LogHandler == nil {
		cfg.LogHandler = cfg.defaultLogHandler()
	}
	logger = logging.NewLogger(cfg.LogHandler)

	ng := newEngine(logger)
	return ng.start(cfg, patterns...)
}
