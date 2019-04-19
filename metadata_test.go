package wgus

import (
	"log"
)

func ExampleGetMetadataGetMetadata() {
	meta, err := GetMetadata(`wgus-wotru.wargaming.net`, `WOT.RU.PRODUCTION`, ``)
	if err != nil {
		panic(err)
	}
	log.Println(meta)
}
