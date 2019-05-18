package wgus

import (
	"testing"
)

func TestGetPatchesChain(t *testing.T) {
	_, err := GetPatchesChain(`wgus-wotru.wargaming.net`, `WOT.RU.PRODUCTION`, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
}
