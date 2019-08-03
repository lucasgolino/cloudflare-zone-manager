package czm

import (
	SLog "github.com/quan-to/slog"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Cloudflare struct {
		Email  string `yaml:"email"`
		Apikey string `yaml:"api_key"`
	} `yaml:"cloudflare"`
	Zones []struct {
		Id       string `yaml:"id"`
		Hostname string `yaml:"hostname"`
		Dns      []struct {
			Name    string `yaml:"name"`
			Dtype   string `yaml:"dtype"`
			Content string `yaml:"content"`
			Proxied bool   `yaml:"proxied"`
			Module  struct {
				Name     string `yaml:"name"`
				Metadata []struct {
					Key  string `yaml:"key"`
					Data string `yaml:"Data"`
				} `yaml:"metadata"`
			} `yaml:"module"`
		} `yaml:"dns"`
	} `yaml:"zones"`
}

func ReadConfig() (Config) {
	var config Config

	yamlFile, err := ioutil.ReadFile("../dns.yaml")

	if err != nil {
		SLog.Scope("ReadZones").Error(err)
	}

	err = yaml.Unmarshal(yamlFile, &config)

	if err != nil {
		SLog.Scope("ReadZones").Error(err)
	}

	return config
}
