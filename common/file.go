// Copyright 2022 The block-node Authors
// This file is part of the block-node library.
//
// The block-node library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The block-node library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the block-node library. If not, see <http://www.gnu.org/licenses/>.
package common

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/block-api/block-node/log"
	"gopkg.in/yaml.v3"
)

var (
	ErrInvalidFileType = errors.New("invalid file type")
)

type Type string

const (
	JSON Type = "JSON"
	YML  Type = "YML"
)

type File struct {
	file     *os.File
	content  []byte
	fileType Type
}

func OpenFile(filePath string, fileType Type) (*File, error) {
	log.New("open file: "+filePath, log.TypeDebug)

	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, _ := ioutil.ReadAll(file)

	return &File{
		file:     file,
		content:  content,
		fileType: fileType,
	}, nil
}

func (f *File) Parse(out interface{}) error {
	if f.fileType == JSON {
		err := f.parseJSON(out)

		if err != nil {
			return err
		}

		return nil
	}

	if f.fileType == YML {
		err := f.parseYML(out)
		if err != nil {
			return err
		}

		return nil
	}

	return ErrInvalidFileType
}

func (f *File) parseJSON(out interface{}) error {
	errUnmarshal := json.Unmarshal(f.content, out)
	if errUnmarshal != nil {
		return errUnmarshal
	}

	return nil
}

func (f *File) parseYML(out interface{}) error {
	err := yaml.Unmarshal(f.content, out)
	if err != nil {
		return err
	}

	return nil
}
