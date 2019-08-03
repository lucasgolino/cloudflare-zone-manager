package czm

import "github.com/quan-to/slog"

type CloudflareZoneManager struct {
	ConfigMap ConfigMap
	Env string
	SLog *slog.Instance
	Srv struct {
		CF *Cloudflare
		Reporter string
	}
}

func (e *CloudflareZoneManager) Init() {
	e.SLog = slog.Scope(`CZM`)

	e.ConfigMap = ReadConfigMap()
	e.InitServices()

	e.VerifyAndUpdateZones()
}

func (e *CloudflareZoneManager) InitServices() {
	e.Srv.CF = e.ConfigMap.Cloudflare.Initialize()
}

func (e *CloudflareZoneManager) VerifyAndUpdateZones() {
	var zones = &e.ConfigMap.Zones

	for _, zone := range *zones {
		if !e.Srv.CF.ExistsZone(&zone) {
			continue
		}

		e.Srv.CF.LoadDnsRecords(&zone)

		for _, dns := range zone.Dns {
			exists := e.Srv.CF.ExistsDnsRule(&zone, &dns)
			rules := Rules{
				dns.Rules.NotExist,
				dns.Rules.Update,
			}

			if !exists {
				if rules.VerifyRule(RULES_NEXISTS_TAG) {
					e.SLog.Warn(`DNS Zone: %s dosen't exists... skip`, dns.Name)
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


