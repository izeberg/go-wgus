package wgus

import (
	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type ReferenceRepositoryFile struct {
	Name string
	CRC  int64
	SHA1 string
	Size int
}

func (s ReferenceRepositoryFile) Hash() interface{} {
	if s.SHA1 != `` {
		return s.SHA1
	}
	if s.CRC != 0 {
		return s.CRC
	}
	return s.Size
}

type ReferenceRepository map[string]ReferenceRepositoryFile

func (s ReferenceRepository) Parse(data []byte) error {
	for _, fileInfoStr := range strings.Split(string(data), "\n") {
		if fileInfoStr != `` {
			fileInfo := strings.Split(fileInfoStr, "\t") // filename size crc32
			if len(fileInfo) == 3 {
				name := strings.TrimSpace(fileInfo[0])
				if size, err := strconv.Atoi(strings.TrimSpace(fileInfo[1])); err == nil {
					if crc, err := strconv.ParseInt(strings.TrimSpace(fileInfo[2]), 16, 64); err == nil {
						file := ReferenceRepositoryFile{
							Name: name,
							CRC:  crc,
							Size: size,
						}
						s[strings.TrimSpace(fileInfo[0])] = file
					}
				}
			}
		}
	}
	return nil
}

func (s ReferenceRepository) Fetch(url string) error {
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			if data, err := ioutil.ReadAll(resp.Body); err == nil {
				return s.Parse(data)
			} else {
				return err
			}
		} else {
			return errors.New(resp.Status)
		}
	} else {
		return err
	}
}

func (s ReferenceRepository) FetchTorrent(url TorrentFile) error {
	if resp, err := http.Get(string(url)); err == nil {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			if meta, err := metainfo.Load(resp.Body); err == nil {
				info := &metainfo.Info{}
				if err := bencode.Unmarshal(meta.InfoBytes, info); err != nil {
					return err
				}
				for i, fileInfo := range info.Files {
					displayPath := fileInfo.DisplayPath(info)
					s[displayPath] = ReferenceRepositoryFile{
						Name: displayPath,
						SHA1:  info.Piece(i).Hash().HexString(),
						Size: int(fileInfo.Length),
					}
				}
				return nil
			} else {
				return err
			}
		} else {
			return errors.New(resp.Status)
		}
	} else {
		return err
	}
}

func (s ReferenceRepository) GetModsFolder() string {
	for key := range s {
		if results := ModsFolderRegex.FindStringSubmatch(key); len(results) > 0 {
			return results[1]
		}
	}
	return ``
}
