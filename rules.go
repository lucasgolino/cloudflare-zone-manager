package czm

type Rules struct {
	NotExist string `yaml:"not-exist"` // create | skip
	Update string `yaml:"update"` // always | if-different | never
}

func (r *Rules) VerifyRule(tag string) (bool) {
	switch tag {
	case RULES_NEXISTS_TAG:
		if r.NotExist == RULES_NEXISTS_CREATE {
			return false
		} else if r.NotExist == RULES_NEXISTS_SKIP {
			return true
		}

		return false

	case RULES_UPDATE_TAG:
		if r.Update == RULES_UPDATE_ALWAYS  {
			return false
		} else if r.Update == RULES_UPDATE_NEVER {
			return true
		}
		return false

	default:
		return true
	}

	return true
}