package wgus

import (
	"github.com/pkg/errors"
)

const (
	ProtocolVersion         = `1.9`
	MetadataProtocolVersion = `7.2`
	InstallationID          = `go-wgus`
)

var ErrUnknownClientType = errors.New(`unknown client_type, try to check metadata`)
