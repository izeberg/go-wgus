package wgus

import "net/url"

type File struct {
	Name         string `xml:"name"`
	Size         int    `xml:"size"`
	UnpackedSize int    `xml:"unpacked_size"`
}

type Patch struct {
	Files       []File        `xml:"files>file"`
	TorrentHash string        `xml:"torrent>hash"`
	TorrentURLs []TorrentFile `xml:"torrent>urls>url"`
	VersionTo   string        `xml:"version_to"`
	Part        string        `xml:"part"`
}

type WebSeed struct {
	Threads int    `xml:"threads,attr"`
	Name    string `xml:"name,attr"`
	URL     string `xml:",innerxml"`
}

type PatchesChain struct {
	Type        string    `xml:"type,attr"`
	VersionName string    `xml:"version_name"`
	Patches     []Patch   `xml:"patch"`
	WebSeeds    []WebSeed `xml:"web_seeds>url"`
}

type PatchesChainProtocol struct {
	Protocol
	PatchesChain []PatchesChain `xml:"patches_chain"`
}

func (s PatchesChainProtocol) GetPatchesChain(chainType string) *PatchesChain {
	for _, chain := range s.PatchesChain {
		if chain.Type == chainType {
			return &chain
		}
	}
	return nil
}

func GetPatchesChain(host string, gameID string, versions map[string]string, query url.Values, meta *MetadataProtocol) (*PatchesChainProtocol, error) {
	if query == nil {
		query = url.Values{}
	}

	if meta == nil {
		if _meta, err := GetMetadata(host, gameID, ``); err != nil {
			return nil, err
		} else {
			meta = _meta
		}
	}

	query.Set(`game_id`, gameID)
	query.Set(`protocol_version`, ProtocolVersion)
	query.Set(`metadata_protocol_version`, MetadataProtocolVersion)
	query.Set(`installation_id`, InstallationID)

	querySetDefault(query, `client_type`, meta.Metadata.DefaultClientType)
	querySetDefault(query, `lang`, meta.DefaultLanguage)
	querySetDefault(query, `metadata_version`, meta.MetadataVersion)

	if clientType := meta.Metadata.GetClientType(query.Get(`client_type`)); clientType != nil {
		for _, part := range clientType.Parts {
			version := `0`
			if v, found := versions[part.ID]; found {
				version = v
			}
			querySetDefault(query, part.ID+`_current_version`, version)
		}
	} else {
		return nil, ErrUnknownClientType
	}

	proto := &PatchesChainProtocol{}
	if err := request(proto, makeRequest(host, `/api/v1/patches_chain/`, query)); err == nil {
		if err := unpackError(proto.Protocol); err == nil {
			return proto, nil
		} else {
			return proto, err
		}
	} else {
		return nil, err
	}
}
