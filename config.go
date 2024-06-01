package czm

import (
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	SLog "github.com/quan-to/slog"
	"gopkg.in/yaml.v2"
	"os"
)

type DNSMetadata struct {
	Key  string `yaml:"key"`
	Data string `yaml:"data"`
}

type Dns struct {
	ID      string
	Name    string `yaml:"name"`
	Dtype   string `yaml:"dtype"`
	Content string `yaml:"content"`
	Proxied *bool  `yaml:"proxied"`
	Rules   Rules  `yaml:"rules"`
	TTL     int    `yaml:"ttl"`
	Module  Module `yaml:"module"`
}

type Zone struct {
	Id         string `yaml:"id"`
	Hostname   string `yaml:"hostname"`
	Dns        []Dns  `yaml:"dns"`
	DNSRecords []cloudflare.DNSRecord
}

type ConfigMap struct {
	Cloudflare Cloudflare `yaml:"cloudflare"`
	Zones      []Zone     `yaml:"zones"`
}

func (cMap *ConfigMap) ReadConfigMap() {
	log := SLog.Scope(fmt.Sprintf("LoadConfig: %s", CONFIG_MAP_PATH))
	yamlFile, err := os.ReadFile(CONFIG_MAP_PATH)

	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(yamlFile, &cMap)

	if err != nil {
		log.Fatal(err)
	}
}
