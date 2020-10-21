package wgus

import (
	"github.com/pkg/errors"
	"regexp"
)

const (
	ProtocolVersion         = `1.8`
	MetadataProtocolVersion = `6.5`
	InstallationID          = `go-wgus`
)

var (
	ModsFolderRegex      = regexp.MustCompile(`mods/(.+?)/readme\.txt`)
	ErrUnknownClientType = errors.New(`unknown client_type, try to check metadata`)
)
