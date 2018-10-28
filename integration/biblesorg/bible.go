package biblepassageapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var _ Bible = BibleBiblesOrg{}

type BibleBiblesOrg struct {
	apiKey string

	// translation details
	id         string
	nameShort  string
	nameCommon string
	nameLong   string
	copyright  string
}

func (b BibleBiblesOrg) Source() string {
	return "biblesorg"
}

func (b BibleBiblesOrg) NameShort() string {
	return b.nameShort
}

func (b BibleBiblesOrg) NameCommon() string {
	return b.nameCommon
}

func (b BibleBiblesOrg) NameLong() string {
	return b.nameLong
}

func (b BibleBiblesOrg) GetPassage(reference string) (*Passage, error) {
	q := url.Values{}
	q.Add("q[]", reference)

	u := url.URL{
		Scheme:   "https",
		Host:     "bibles.org",
		Path:     fmt.Sprintf("/v2/%s/passages.js", b.id),
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
		return nil, fmt.Errorf("reference %s returned %d passages, expected %d", reference, passageCount, expectedPassageCount)
	}

	html := ""
	for _, passage := range searchRes.Response.Search.Result.Passages {
		html += passage.Text
	}

	copyright := b.copyright
	if copyright == "" {
		copyright = transposePassageCopyright(searchRes.Response.Search.Result.Passages[0].Copyright)
	}

	passage := Passage{
		HTML:         transposePassageHTML(html),
		TrackingCode: searchRes.Response.Meta.FumsNoscript,
		Copyright:    copyright,
	}

	return &passage, nil
}

func transposePassageHTML(s string) string {
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
