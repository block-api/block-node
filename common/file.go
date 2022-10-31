package common

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/block-api/block-node/errors"
	"github.com/block-api/block-node/log"
	"gopkg.in/yaml.v3"
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

	content, _ := ioutil.ReadAll(file)

	return &File{
		file:     file,
		content:  content,
		fileType: fileType,
	}, nil
}

func (f *File) Close() {
	if f.file != nil {
		_ = f.file.Close()
	}
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

	return errors.ErrInvalidFileType
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
