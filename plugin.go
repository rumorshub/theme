package theme

import (
	"os"
	"sync"

	et "github.com/gowool/extends-template"
	"github.com/roadrunner-server/endure/v2/dep"
	"github.com/roadrunner-server/errors"
)

const PluginName = "theme"

type Plugin struct {
	mu   sync.RWMutex
	once sync.Once

	cfg *Config

	loaders []et.Loader
	env     *Environment
}

func (p *Plugin) Init(cfg Configurer, logger Logger) error {
	const op = errors.Op("theme_plugin_init")
	if !cfg.Has(PluginName) {
		return errors.E(op, errors.Disabled)
	}

	if err := cfg.UnmarshalKey(PluginName, &p.cfg); err != nil {
		return errors.E(op, err)
	}

	p.cfg.InitDefaults()

	log := logger.NamedLogger(PluginName)

NEXT:
	for _, loaderCfg := range p.cfg.Loaders {
		loader := et.NewFilesystemLoader(os.DirFS(loaderCfg.Dir))

		for namespace, paths := range loaderCfg.Paths {
			if err := loader.SetPaths(namespace, paths...); err != nil {
				log.Warn("loader set namespace paths", "err", err, "dir", loaderCfg.Dir, "namespace", namespace, "paths", paths)

				continue NEXT
			}
		}

		p.loaders = append(p.loaders, loader)
	}

	return nil
}

func (p *Plugin) Name() string {
	return PluginName
}

func (p *Plugin) Collects() []*dep.In {
	return []*dep.In{
		dep.Fits(func(pp any) {
			loader := pp.(et.Loader)

			p.mu.Lock()
			p.loaders = append(p.loaders, loader)
			p.mu.Unlock()
		}, (*et.Loader)(nil)),
	}
}

func (p *Plugin) Provides() []*dep.Out {
	return []*dep.Out{
		dep.Bind((*Theme)(nil), p.Theme),
	}
}

func (p *Plugin) Theme() *Environment {
	p.once.Do(func() {
		p.mu.Lock()
		p.env = NewEnvironment(p.loaders)
		p.env.Debug(p.cfg.Debug).Delims(p.cfg.Delims.Left, p.cfg.Delims.Right).Global(p.cfg.Global...)
		p.mu.Unlock()
	})

	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.env
}
