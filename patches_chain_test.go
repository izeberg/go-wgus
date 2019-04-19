package wgus

import (
	"log"
)


func ExampleGetPatchesChain() {
	chain, err := GetPatchesChain(`wgus-wotru.wargaming.net`, `WOT.RU.PRODUCTION`, nil, nil)
	if err != nil {
		panic(err)
	}
	log.Println(chain)
}
