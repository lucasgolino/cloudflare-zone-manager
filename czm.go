package czm

import "github.com/quan-to/slog"

type CloudflareZoneManager struct {
	Config   Config
	Services struct {
		CF *Cloudflare
		Reporter string
	}
	SLog *slog.Instance
}

func (e *CloudflareZoneManager) Init() {
	e.SLog = slog.Scope(`CZM`)

	e.Config = ReadConfig()
	e.InitServices()

	e.VerifyAndUpdateZones()
}

func (e *CloudflareZoneManager) InitServices() {
	cf := Cloudflare{}
	e.Services.CF = cf.Initialize(e.Config)
}

func (e *CloudflareZoneManager) VerifyAndUpdateZones() {
	var zones = e.Config.Zones

	for _, zone := range zones {
		if !e.Services.CF.ExistsZone(zone.Id, zone.Hostname) {
			continue
		}

		dnsRecords := e.Services.CF.LoadDnsRecords(zone.Id)

		for _, dns := range zone.Dns {
			dnsData, exists := e.Services.CF.ExistsDnsRule(dnsRecords, dns.Name)

			if !exists {
				e.SLog.Warn(`DNS Zone: %s dosen't exists... skip`, dns.Name)
				continue
			}

			e.SLog.Info(dnsData)
		}
	}
}


