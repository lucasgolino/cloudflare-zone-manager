package czm

import (
	"context"
	"github.com/cloudflare/cloudflare-go"
	"github.com/logrusorgru/aurora"
	"github.com/quan-to/slog"
)

type Cloudflare struct {
	Email    string `yaml:"email"`
	APIKey   string `yaml:"api_key"`
	SLog     *slog.Instance
	Api      *cloudflare.API
	UserInfo cloudflare.User
	ctx      context.Context
}

func (cf *Cloudflare) Initialize() *Cloudflare {
	cf.ctx = context.Background()
	cf.SLog = slog.Scope("Cloudflare-API")

	api, err := cloudflare.New(cf.APIKey, cf.Email)
	cf.Api = api

	if err != nil {
		cf.SLog.Fatal(err)
	}

	cf.SLog.Info(`Cloudflare connection successful created`)
	cf.SLog.Info(`Fetching user information`)

	cf.UserInfo, err = cf.Api.UserDetails(cf.ctx)

	if err != nil {
		cf.SLog.Error(`Fail to fetch user information`)
		cf.SLog.Fatal(err)
	}

	return cf
}

func (cf *Cloudflare) ExistsZone(zone *Zone) bool {
	zoneResponse, err := cf.Api.ZoneDetails(cf.ctx, zone.Id)

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

func (cf *Cloudflare) LoadDNSRecords(zone *Zone) {
	cf.SLog.Info(`Loading DNS Records for %s zone`, zone.Hostname)
	zoneId := cloudflare.ZoneIdentifier(zone.Id)
	records, _, err := cf.Api.ListDNSRecords(cf.ctx, zoneId, cloudflare.ListDNSRecordsParams{})

	if err != nil {
		cf.SLog.Fatal(err)
	}

	zone.DNSRecords = records
}

func (cf *Cloudflare) ExistsDNSRecord(zone *Zone, dns *Dns) bool {
	for _, value := range zone.DNSRecords {
		if value.Name == dns.Name {
			dns.ID = value.ID
			return true
		}
	}

	return false
}

func (cf *Cloudflare) CreateDNSRecord(zone *Zone, dns *Dns) error {
	if dns.Content == "" {
		cf.SLog.SubScope(dns.Name).Info("Resolving module")
		dns.Content = dns.Module.Resolve()
	}

	zoneId := cloudflare.ZoneIdentifier(zone.Id)
	response, err := cf.Api.CreateDNSRecord(cf.ctx, zoneId, cloudflare.CreateDNSRecordParams{
		ID:      dns.ID,
		Type:    dns.Dtype,
		Name:    dns.Name,
		Content: dns.Content,
		Proxied: dns.Proxied,
	})

	if err == nil {
		dns.ID = response.ID
	}

	cf.SLog.SubScope(dns.Name).Info("Zone created")

	return err
}

func (cf *Cloudflare) DNSRecordHasDiff(zone *Zone, dns *Dns) bool {
	log := cf.SLog.SubScope(dns.Name)
	zoneId := cloudflare.ZoneIdentifier(zone.Id)
	response, err := cf.Api.GetDNSRecord(cf.ctx, zoneId, dns.ID)

	hasDiff := false

	if err != nil {
		log.Error(err)
		return false
	}

	if response.Name != dns.Name {
		log.Info("DNS Zone has different name")
		hasDiff = true
	}

	if response.Proxiable {
		if response.Proxied != dns.Proxied {
			log.Info("")
			hasDiff = true
		}
	}

	if response.Content != dns.Content {
		log.Info("DNS Zone has different payload")
		log.Log(` diff(%s, %s)%s`, aurora.Red(response.Content), aurora.Green(dns.Content), aurora.Cyan(""))
		hasDiff = true
	}

	if response.Type != dns.Dtype {
		log.Info("DNS Zone has different type")
		hasDiff = true
	}

	if response.TTL != dns.TTL {
		log.Info("DNS Zone has different TTL")
		hasDiff = true
	}

	return hasDiff
}

func (cf *Cloudflare) UpdateDNSRecord(zone *Zone, dns *Dns) error {

	cf.SLog.SubScope(dns.Name).Info(`Updating Record`)
	zoneId := cloudflare.ZoneIdentifier(zone.Id)
	_, err := cf.Api.UpdateDNSRecord(cf.ctx, zoneId, cloudflare.UpdateDNSRecordParams{
		ID:      dns.ID,
		Name:    dns.Name,
		Type:    dns.Dtype,
		Content: dns.Content,
		TTL:     dns.TTL,
		Proxied: dns.Proxied,
	})

	return err
}
