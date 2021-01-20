package wgus

import (
	"encoding/xml"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type ChangedGameInfo struct {
	RedirectApplicationID string `xml:"redirect_application_id"`
	RedirectURL           string `xml:"redirect_url"`
}

type Error struct {
	ErrorCode        int              `xml:"code"`
	ErrorDescription string           `xml:"description"`
	ChangedGameInfo  *ChangedGameInfo `xml:"changed_game_info"`
}

type Protocol struct {
	Error
	Name    string `xml:"name,attr"`
	Version string `xml:"version,attr"`
}

type TorrentFile string

func makeRequest(host, path string, query url.Values) url.URL {
	return url.URL{
		Scheme:   `http`,
		Host:     host,
		Path:     path,
		RawQuery: query.Encode(),
	}
}

func querySetDefault(query url.Values, key, value string) {
	if _, found := query[key]; !found {
		query.Set(key, value)
	}
}

func unpackError(proto Protocol) error {
	if proto.ErrorCode != 0 {
		return errors.New(proto.ErrorDescription + `:` + strconv.Itoa(proto.ErrorCode))
	}
	if proto.ChangedGameInfo != nil {
		return errors.New(`ChangedGameInfo:` + proto.ChangedGameInfo.RedirectApplicationID + `@` + proto.ChangedGameInfo.RedirectURL)
	}
	return nil
}

func request(proto interface{}, requestURL url.URL) error {
	if resp, err := http.Get(requestURL.String()); err == nil {
		if data, err := ioutil.ReadAll(resp.Body); err == nil {
			if err := xml.Unmarshal(data, proto); err == nil {
				if resp.StatusCode != http.StatusOK {
					return errors.New(resp.Status)
				}
				return nil
			} else {
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}
}
