package wgus

import (
	"log"
	"testing"
)

func TestGetMetadata(t *testing.T) {
	meta, err := GetMetadata(`wgus-wotru.wargaming.net`, `WOT.RU.PRODUCTION`, ``)
	if err != nil {
		panic(err)
	}
	log.Println(meta)
}
