package wgus

import (
	"github.com/j-muller/go-torrent-parser"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strings"
)

type Torrent struct {
	Files      []TorrentFile `xml:"file"`
	Blacklists []TorrentFile `xml:"blacklist"`
	Part       string        `xml:"part"`
	Hash       string        `xml:"hash"`
}

type IntegrityCheckProtocol struct {
	Protocol
	Torrents          []Torrent `xml:"torrents>torrent"`
	PrivatePTPEnabled bool      `xml:"torrents>private_ptp_enabled"`
	WebSeeds          []WebSeed `xml:"web_seeds>url"`
}

func GetIntegrityCheck(host string, gameID string, versions map[string]string, query url.Values) (*IntegrityCheckProtocol, error) {
	query.Set(`game_id`, gameID)
	query.Set(`protocol_version`, `1.6`)
	query.Set(`installation_id`, InstallationID)

	for part, version := range versions {
		query.Set(part+`_check_version`, version)
	}

	proto := &IntegrityCheckProtocol{}
	if err := request(proto, makeRequest(host, `/api/v1/integrity_check/`, query)); err == nil {
		if err := unpackError(proto.Protocol); err == nil {
			return proto, nil
		} else {
			return proto, err
		}
	} else {
		return nil, err
	}
}

func (s IntegrityCheckProtocol) GetReferenceRepository() (ReferenceRepository, error) {
	repository := ReferenceRepository{}
	var fetchErr error
	for _, file := range s.Torrents {
		if repo, err := file.GetReferenceRepository(); err == nil {
			for key, value := range repo {
				repository[key] = value
			}
		} else {
			fetchErr = err
		}
	}
	return repository, fetchErr
}

func (s Torrent) GetReferenceRepository() (ReferenceRepository, error) {
	var refErr error
	for _, file := range s.Files {
		if ref, err := file.GetReferenceRepository(); err == nil {
			return ref, nil
		} else {
			refErr = err
		}
	}
	return nil, refErr
}

func (s TorrentFile) GetReferenceRepository() (ReferenceRepository, error) {
	repo := ReferenceRepository{}
	return repo, repo.Fetch(strings.Replace(string(s), `.torrent`, `.filelist.txt`, -1))
}

func (s TorrentFile) Parse() (*gotorrentparser.Torrent, error) {
	if res, err := http.Get(string(s)); err == nil {
		defer res.Body.Close()
		if res.StatusCode == http.StatusOK {
			if torrent, err := gotorrentparser.Parse(res.Body); err == nil {
				return torrent, nil
			} else {
				return nil, err
			}
		} else {
			return nil, errors.New(res.Status)
		}
	} else {
		return nil, err
	}
}
