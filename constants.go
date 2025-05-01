package wgus

import (
	"errors"
)

const (
	ProtocolVersion         = `100500.6969696`
	MetadataProtocolVersion = `100500.6969696`
	InstallationID          = `go-wgus`
)

var ErrUnknownClientType = errors.New(`unknown client_type, try to check metadata`)
