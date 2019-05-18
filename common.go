package wgus

import (
	"encoding/xml"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Error struct {
	ErrorCode        int    `xml:"code"`
	ErrorDescription string `xml:"description"`
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
		return errors.New(proto.ErrorDescription)
	}
	return nil
}

func request(proto interface{}, requestURL url.URL) error {
	if resp, err := http.Get(requestURL.String()); err == nil {
		if data, err := ioutil.ReadAll(resp.Body); err == nil {
			if err := xml.Unmarshal(data, proto); err == nil {
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
