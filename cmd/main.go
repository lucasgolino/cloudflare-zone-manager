package main

import czm "github.com/lucasgolino/cloudflare-zone-manager"

func main() {
	var CloudflareZM = czm.CloudflareZoneManager{}

	CloudflareZM.Init()
}
