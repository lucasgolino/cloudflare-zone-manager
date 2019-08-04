package czm

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/quan-to/slog"
	"os"
	"path/filepath"
	"plugin"
	"strings"
)

type Module struct {
	Name     string `yaml:"name"`
	Metadata []DNSMetadata `yaml:"metadata"`
	Mods *Modules
}

type Plugin interface {
	Resolve(interface{}) string
}

func (m *Module) Resolve() (string) {

	for _, mod := range m.Mods.Loaded {
		if mod.Name == m.Name {
			SLog := slog.Scope("Module").SubScope(fmt.Sprint(mod.Name, `.Resolve()`))
			symPlugin, err := mod.Plugin.Lookup("Plugin")

			if err != nil {
				SLog.Fatal(err)
			}

			var plugin Plugin
			plugin, ok := symPlugin.(Plugin)

			if !ok {
				SLog.Fatal("unexpected type from module symbol")
			}

			return plugin.Resolve("")
		}
	}

	return ""
}

type ModuleLoad struct {
	Name string
	Plugin *plugin.Plugin
}

type Modules struct {
	Loaded []*ModuleLoad
	SLog *slog.Instance
}

func (m *Modules) LoadAllModules() {
	m.SLog = slog.Scope("Modules")

	err := filepath.Walk(CONFIG_MOD_PATH, func(path string, info os.FileInfo, err error) error {

		if filepath.Ext(path) == ".so" {
			plugin, err := plugin.Open(path)

			if err != nil {
				return err
			}

			Mod := &ModuleLoad{
				Name: strings.TrimSuffix(info.Name(), filepath.Ext(path)),
				Plugin: plugin,
			}

			m.SLog.Log(`Module %s %s`, aurora.Yellow(Mod.Name), aurora.Cyan("is loaded"))

			m.Loaded = append(m.Loaded, Mod)
		}

		return nil
	})

	if err != nil {
		m.SLog.Fatal(err)
	}

	m.SLog.Info(`All Modules are loaded`)
}