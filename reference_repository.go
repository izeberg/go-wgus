package wgus

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type ReferenceRepositoryFile struct {
	Name string
	CRC  int64
	Size int
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

func (s ReferenceRepository) GetModsFolder() string {
	for key := range s {
		if results := ModsFolderRegex.FindStringSubmatch(key); len(results) > 0 {
			return results[1]
		}
	}
	return ``
}
