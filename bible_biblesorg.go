package biblepassageapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type translationDetails struct {
	ID         string
	Name       string
	NameCommon string
	NameShort  string
	Copyright  string
}

var (
	CEV  = "CEV"
	GNT  = "GNT"
	NASB = "NASB"
	AMP  = "AMP"
	MSG  = "MSG"
)

var translations = map[string]translationDetails{
	CEV: {
		ID:         "eng-CEVD",
		NameShort:  "CEV",
		NameCommon: "Contemporary English Version",
		Name:       "2006 Contemporary English Version, Second Edition (US Version)",
	},
	GNT: {
		ID:         "eng-GNTD",
		NameShort:  "GNT",
		NameCommon: "Good News Translation",
		Name:       "1992 Good News Translation, Second Edition (US Version)",
	},
	NASB: {
		ID:         "eng-NASB",
		NameShort:  "NASB",
		NameCommon: "New American Standard Bible",
		Name:       "1995, New American Standard Bible",
	},
	AMP: {
		ID:         "eng-AMP",
		NameShort:  "AMP",
		NameCommon: "Amplified Bible",
		Name:       "Amplified Bible",
	},
	MSG: {
		ID:         "eng-MSG",
		NameShort:  "MSG",
		NameCommon: "The Message",
		Name:       "The Message",
		Copyright:  "Scripture taken from The Message. Copyright © 1993, 1994, 1995, 1996, 2000, 2001, 2002. Used by permission of NavPress Publishing Group.",
	},
}

type BibleBiblesOrg struct {
	apiKey      string
	translation translationDetails
}

func NewBiblesOrg(apiKey, translation string) BibleBiblesOrg {
	return BibleBiblesOrg{apiKey: apiKey, translation: translations[translation]}
}

func (b BibleBiblesOrg) Source() string {
	return "biblesorg"
}

func (b BibleBiblesOrg) NameShort() string {
	return b.translation.NameShort
}

func (b BibleBiblesOrg) NameCommon() string {
	return b.translation.NameCommon
}

func (b BibleBiblesOrg) Name() string {
	return b.translation.Name
}

func (b BibleBiblesOrg) GetPassage(reference string) (*Passage, error) {
	q := url.Values{}
	q.Add("q[]", reference)

	u := url.URL{
		Scheme:   "https",
		Host:     "bibles.org",
		Path:     fmt.Sprintf("/v2/%s/passages.js", b.translation.ID),
		RawQuery: q.Encode(),
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", u.String(), nil)
	req.SetBasicAuth(b.apiKey, "X")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var searchRes struct {
		Response struct {
			Search struct {
				Result struct {
					Passages []struct {
						Copyright string
						Text      string
					}
				}
			}
			Meta struct {
				FumsNoscript string `json:"fums_noscript"`
			}
		}
	}
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&searchRes)
	if err != nil {
		return nil, fmt.Errorf("error json decoding %d response from bibles.org api: %v", res.StatusCode, err)
	}

	expectedPassageCount := strings.Count(reference, ",") + 1
	passageCount := len(searchRes.Response.Search.Result.Passages)
	if passageCount != expectedPassageCount {
		return nil, fmt.Errorf("Reference %s returned %d passages, expected %d.", reference, passageCount, expectedPassageCount)
	}

	html := ""
	for _, passage := range searchRes.Response.Search.Result.Passages {
		html += passage.Text
	}

	passage := Passage{
		Html:         transposePassageHtml(html),
		TrackingCode: searchRes.Response.Meta.FumsNoscript,
		Copyright: func() string {
			if b.translation.Copyright != "" {
				return b.translation.Copyright
			} else {
				return transposePassageCopyright(searchRes.Response.Search.Result.Passages[0].Copyright)
			}
		}(),
	}

	return &passage, nil
}

func transposePassageHtml(s string) string {
	s = strings.Replace(s, "<sup", " <sup", -1)
	s = strings.Replace(s, "</sup>", "</sup> ", -1)
	s = strings.Replace(s, "[Lord]", "Lord", -1)
	return s
}

func transposePassageCopyright(s string) string {
	s = regexp.MustCompile("(&#169;)").ReplaceAllString(s, "</br>$1")
	s = regexp.MustCompile("(</p>)").ReplaceAllString(s, "</br>Used with permission.$1")
	s = regexp.MustCompile(",([^ ])").ReplaceAllString(s, ", $1")
	return s
}
