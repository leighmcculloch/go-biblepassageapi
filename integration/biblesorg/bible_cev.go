package biblepassageapi

func NewBiblesOrgCEV(apiKey string) BibleBiblesOrg {
	return BibleBiblesOrg{
		apiKey:     apiKey,
		id:         "eng-CEVD",
		nameShort:  "CEV",
		nameCommon: "Contemporary English Version",
		nameLong:   "2006 Contemporary English Version, Second Edition (US Version)",
	}
}
