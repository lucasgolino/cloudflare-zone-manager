package czm

import (
	cfgo "github.com/cloudflare/cloudflare-go"
	SLog "github.com/quan-to/slog"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type DNSMetadata struct {
	Key  string `yaml:"key"`
	Data string `yaml:"data"`
}


type Dns struct {
	ID	string
	Name string `yaml:"name"`
	Dtype string `yaml:"dtype"`
	Content string `yaml:"content"`
	Proxied bool   `yaml:"proxied"`
	Rules Rules `yaml:"rules"`
	TTL int `yaml:"ttl"`
	Module Module `yaml:"module"`
}

type Zone struct {
	Id       string `yaml:"id"`
	Hostname string `yaml:"hostname"`
	Dns []Dns `yaml:"dns"`
	DNSRecords []cfgo.DNSRecord

}

type ConfigMap struct {
	Cloudflare Cloudflare `yaml:"cloudflare"`
	Zones []Zone `yaml:"zones"`
}

func (cMap *ConfigMap) ReadConfigMap() {
	yamlFile, err := ioutil.ReadFile(CONFIG_MAP_PATH)

	if err != nil {
		SLog.Scope("ReadZones").Error(err)
	}

	err = yaml.Unmarshal(yamlFile, &cMap)

	if err != nil {
		SLog.Scope("ReadZones").Error(err)
	}
}
