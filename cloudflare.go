package czm

import (
	cfgo "github.com/cloudflare/cloudflare-go"
	"github.com/quan-to/slog"
)

type Cloudflare struct {
	APIKey string
	Email string
	SLog *slog.Instance
	Api *cfgo.API
	UserInfo cfgo.User
	cfg Config
}

func (cf *Cloudflare) Initialize(cfg Config) (*Cloudflare) {
	cf.SLog = slog.Scope("Cloudflare-API")

	cf.APIKey = cfg.Cloudflare.Apikey
	cf.Email = cfg.Cloudflare.Email

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

func (cf *Cloudflare) ExistsZone(zoneId string, zoneName string) (bool) {
	zone, err := cf.Api.ZoneDetails(zoneId)

	if err != nil {
		cf.SLog.Error(err)
		return false
	}

	if zone.Name != zoneName {
		cf.SLog.Error(`Zone has a different hostname, mark skip`)
		return false
	}

	return true
}

func (cf *Cloudflare) LoadDnsRecords(zoneId string) ([]cfgo.DNSRecord) {
	records, err := cf.Api.DNSRecords(zoneId, cfgo.DNSRecord{})

	if err != nil {
		cf.SLog.Fatal(err)
	}

	return records
}

func (cf *Cloudflare) ExistsDnsRule(dnsRecords []cfgo.DNSRecord, dnsName string) (cfgo.DNSRecord, bool) {
	for _, value := range dnsRecords {
		if value.Name == dnsName {
			return value, true
		}
	}

	return cfgo.DNSRecord{}, false
}