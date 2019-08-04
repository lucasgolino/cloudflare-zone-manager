package czm

import (
	"github.com/quan-to/slog"
)

type CloudflareZoneManager struct {
	ConfigMap ConfigMap
	SLog *slog.Instance
	Srv struct {
		CF *Cloudflare
		Reporter string
		Mods Modules
	}
}

func (e *CloudflareZoneManager) Init() {
	e.SLog = slog.Scope(`CZM`)

	e.ConfigMap.ReadConfigMap()
	e.InitServices()

	e.VerifyAndUpdateZones()
}

func (e *CloudflareZoneManager) InitServices() {
	e.Srv.CF = e.ConfigMap.Cloudflare.Initialize()
	e.Srv.Mods.LoadAllModules()
}

func (e *CloudflareZoneManager) VerifyAndUpdateZones() {
	var zones = &e.ConfigMap.Zones

	for _, zone := range *zones {
		zoneLog := e.SLog.SubScope(zone.Hostname)

		if !e.Srv.CF.ExistsZone(&zone) {
			continue
		}

		e.Srv.CF.LoadDnsRecords(&zone)

		for _, dns := range zone.Dns {
			dns.Module.Mods = &e.Srv.Mods // append address to dns mods
			exists := e.Srv.CF.ExistsDnsRule(&zone, &dns)
			rules := Rules{
				dns.Rules.NotExist,
				dns.Rules.Update,
			}

			if !exists {
				if rules.VerifyRule(RULES_NEXISTS_TAG) {
					zoneLog.Warn(`DNS Zone dosen't exists... skip`)
					continue
				}

				err := e.Srv.CF.CreateDnsRule(&zone, &dns)

				if err != nil {
					e.SLog.Error(err)
					continue
				}
			}
		}
	}
}


