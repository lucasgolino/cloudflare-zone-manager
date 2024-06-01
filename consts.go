package czm

import "os"

const (
	RULES_NEXISTS_TAG    = "NOTEXIST"
	RULES_NEXISTS_CREATE = "create"
	RULES_NEXISTS_SKIP   = "skip"

	RULES_UPDATE_TAG    = "UPDATE"
	RULES_UPDATE_ALWAYS = "always"
	RULES_UPDATE_NEVER  = "never"
)

var (
	CONFIG_MAP_PATH = os.Getenv("CONFIG_MAP_PATH")
	CONFIG_MOD_PATH = os.Getenv("CONFIG_MOD_PATH")
)
