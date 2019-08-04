package czm

import (
	cfgo "github.com/cloudflare/cloudflare-go"
	"github.com/quan-to/slog"
)

type Cloudflare struct {
	Email  string `yaml:"email"`
	APIKey string `yaml:"api_key"`
	SLog     *slog.Instance
	Api      *cfgo.API
	UserInfo cfgo.User
}

func (cf *Cloudflare) Initialize() (*Cloudflare) {
	cf.SLog = slog.Scope("Cloudflare-API")

	api, err := cfgo.New(cf.APIKey, cf.Email)
	cf.Api = api

	if err != nil {
		cf.SLog.Fatal(err)
	} else {
		cf.SLog.Info(`Cloudflare connection is setting!`)
	}

	cf.SLog.Info(`Fetching User informations`)
	cf.UserInfo, err = cf.Api.UserDetails()

	if err != nil {
		cf.SLog.Error(`Fail to fetch user informations`)
		cf.SLog.Fatal(err)
	}

	return cf
}

func (cf *Cloudflare) ExistsZone(zone *Zone) (bool) {
	zoneResponse, err := cf.Api.ZoneDetails(zone.Id)

	if err != nil {
		cf.SLog.Error(err)
		return false
	}

	if zoneResponse.Name != zone.Hostname {
		cf.SLog.Error(`Zone has a different hostname, mark skip`)
		return false
	}

	return true
}

func (cf *Cloudflare) LoadDnsRecords(zone *Zone) {
	cf.SLog.Info(`Loading DNS Records for %s zone`, zone.Hostname)
	records, err := cf.Api.DNSRecords(zone.Id, cfgo.DNSRecord{})

	if err != nil {
		cf.SLog.Fatal(err)
	}

	zone.DNSRecords = records
}

func (cf *Cloudflare) ExistsDnsRule(zone *Zone, dns *Dns) (bool) {
	for _, value := range zone.DNSRecords  {
		if value.Name == dns.Name {
			return true
		}
	}

	return false
}

func (cf *Cloudflare) CreateDnsRule(zone *Zone, dns *Dns) (error) {
	if dns.Content == "" {
		dns.Content = dns.Module.Resolve()
	}

	dres, err := cf.Api.CreateDNSRecord(zone.Id, cfgo.DNSRecord{
		ID: dns.ID,
		Type: dns.Dtype,
		Name: dns.Name,
		Content: dns.Content,
		Proxied: dns.Proxied,
	})

	if err == nil {
		dns.ID = dres.Result.ID
	}

	return err
}