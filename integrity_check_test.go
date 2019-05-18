package wgus

import (
	"net/url"
	"testing"
)

func TestGetIntegrityCheck(t *testing.T) {
	version := map[string]string{
		`client`:    `1.5.0.19635`,
		`hdcontent`: `1.5.0.19609`,
		`locale`:    `1.5.0.1145986`,
		`sdcontent`: `1.5.0.19609`,
	}
	params := url.Values{`locale_lang`: []string{`eu`}}
	check, err := GetIntegrityCheck(`wgus-wotru.wargaming.net`, `WOT.RU.PRODUCTION`, version, params)
	if err != nil {
		t.Fatal(err)
	}
	if ref, err := check.GetReferenceRepository(); err != nil {
		t.Fatal(err)
	} else {
		t.Log(ref)
	}

}
