package wgus

import "github.com/pkg/errors"

const (
	ProtocolVersion = `1.8`
	MetadataProtocolVersion = `5.15`
	InstallationID = `go-wgus`
)


var (
	ErrUnknownClientType = errors.New(`unknown client_type, try to check metadata`)
)