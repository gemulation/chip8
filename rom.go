package chip8

import (
	"io/ioutil"
	"path"

	"github.com/pkg/errors"
)

type ROM struct {
	Name string
	Data []byte
}

func NewROM(filename string) (*ROM, error) {
	name := path.Base(filename)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load rom")
	}
	return &ROM{Name: name, Data: data}, nil
}
