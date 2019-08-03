package main

import czm "golinux.network/tools/cloudflare-zone-manager"

func main() {
	var CloudflareZM = czm.CloudflareZoneManager{}

	CloudflareZM.Init()
}
