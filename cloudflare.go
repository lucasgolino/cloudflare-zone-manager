package czm

import (
	cfgo "github.com/cloudflare/cloudflare-go"
	"github.com/logrusorgru/aurora"
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

func (cf *Cloudflare) LoadDNSRecords(zone *Zone) {
	cf.SLog.Info(`Loading DNS Records for %s zone`, zone.Hostname)
	records, err := cf.Api.DNSRecords(zone.Id, cfgo.DNSRecord{})

	if err != nil {
		cf.SLog.Fatal(err)
	}

	zone.DNSRecords = records
}

func (cf *Cloudflare) ExistsDNSRecord(zone *Zone, dns *Dns) (bool) {
	for _, value := range zone.DNSRecords  {
		if value.Name == dns.Name {
			dns.ID = value.ID
			return true
		}
	}

	return false
}

func (cf *Cloudflare) CreateDNSRecord(zone *Zone, dns *Dns) (error) {
	if dns.Content == "" {
		cf.SLog.SubScope(dns.Name).Info("Resolving module")
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

	cf.SLog.SubScope(dns.Name).Info("Zone created")

	return err
}

func (cf *Cloudflare) DNSRecordHasDiff(zone *Zone, dns *Dns) (bool) {
	dnsData, err := cf.Api.DNSRecord(zone.Id, dns.ID)

	if err != nil {
		cf.SLog.Error(err)
		return false
	}

	if dnsData.Name != dns.Name {
		cf.SLog.SubScope(dns.Name).Info("Has different name")
		return true
	}

	if dnsData.Proxiable {
		if dns.Proxied != dnsData.Proxied {
			cf.SLog.SubScope(dns.Name).Info("Has mark with different Proxied")
			return true
		}
	}

	if dnsData.Content != dns.Content {
		var log = cf.SLog.SubScope(dns.Name)
		log.Info("Has mark with different Content")
		log.Log(` diff(%s, %s%s)`, aurora.Red(dnsData.Content), aurora.Green(dns.Content), aurora.Cyan(""))
		return true
	}

	if dnsData.Type != dns.Dtype {
		cf.SLog.SubScope(dns.Name).Info("Has mark with different Type")
		return true
	}

	if dnsData.TTL != dns.TTL {
		cf.SLog.SubScope(dns.Name).Info("Has mark with different TTL")
		return true
	}

	return false
}

func (cf *Cloudflare) UpdateDNSRecord(zone *Zone, dns *Dns) (error) {

	cf.SLog.SubScope(dns.Name).Info(`Updating DNSRecord`)

	err := cf.Api.UpdateDNSRecord(zone.Id, dns.ID, cfgo.DNSRecord{
		Name: dns.Name,
		Type: dns.Dtype,
		Content: dns.Content,
		TTL: dns.TTL,
		Proxied: dns.Proxied,
	})

	return err
}