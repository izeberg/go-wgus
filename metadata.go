package wgus

import "net/url"

type LaunchOption struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",innerxml"`
}

type Application struct {
	ID     string `xml:",innerxml"`
	Public bool   `xml:"public,attr"`
}

type ClientPart struct {
	ID        string `xml:"id,attr"`
	AppType   string `xml:"app_type,attr"`
	Integrity bool   `xml:"integrity,attr"`
	Lang      bool   `xml:"lang,attr"`
}

type ClientType struct {
	ID    string       `xml:"id,attr"`
	Parts []ClientPart `xml:"client_parts>client_part"`
}

type Metadata struct {
	App                 Application    `xml:"app_id"`
	ChainID             string         `xml:"chain_id"`
	SupportedLanguages  string         `xml:"supported_languages"`
	DefaultLanguage     string         `xml:"default_language"`
	Name                string         `xml:"name"`
	FsName              string         `xml:"fs_name"`
	ShortcutName        string         `xml:"shortcut_name"`
	ExecutableName      string         `xml:"executable_name"`
	LaunchOptions       []LaunchOption `xml:"launch_options>option"`
	MutexName           string         `xml:"mutex_name"`
	KeepPatchesInterval int            `xml:"keep_patches_interval"`
	LoginEnabled        bool           `xml:"login_enabled"`
	ClientTypes       struct{
		Types []ClientType `xml:"client_type"`
		Default string       `xml:"default,attr"`
	} `xml:"client_types"`
}

func (s Metadata) GetClientType(id string) *ClientType {
	for _, t := range s.ClientTypes.Types {
		if t.ID == id {
			return &t
		}
	}
	return nil
}

type MetadataProtocol struct {
	Protocol
	Version     string         `xml:"version"`
	Metadata Metadata `xml:"predefined_section"`
	GeneratedSection interface {} `xml:"generated_section"`

	// TODO: Implement keys
	// example https://wgus-wotru.wargaming.net/api/v1/metadata/?guid=WOT.RU.PRODUCTION&protocol_version=5.15&chain_id=unknown
}

func GetMetadata(host string, guid string, chainID string) (*MetadataProtocol, error) {
	if chainID == `` {
		chainID = `unknown`
	}

	query := url.Values{}
	query.Set(`guid`, guid)
	query.Set(`chain_id`, chainID)
	query.Set(`protocol_version`, MetadataProtocolVersion)

	proto := &MetadataProtocol{}
	if err := request(proto, makeRequest(host, `/api/v1/metadata/`, query)); err == nil {
		if err := unpackError(proto.Protocol); err == nil {
			return proto, nil
		} else {
			return proto, err
		}
	} else {
		return nil, err
	}
}
