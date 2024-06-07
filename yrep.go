package yrep

import (
	"bytes"
	"errors"

	"github.com/goccy/go-yaml"
)

var ErrNotReplaced = errors.New("not replaced")

type yrep struct {
	replaceFuncs []ReplaceFunc
}

type ReplaceFunc func(in any) (any, error)

func Apply(in []byte, fns ...ReplaceFunc) ([]byte, error) {
	y := &yrep{}
	y.replaceFuncs = append(y.replaceFuncs, fns...)

	var m yaml.MapSlice

	if err := yaml.NewDecoder(bytes.NewBuffer(in), yaml.UseOrderedMap()).Decode(&m); err != nil {
		return nil, err
	}

	replaced, err := y.replace(m)
	if err != nil && !errors.Is(err, ErrNotReplaced) {
		return nil, err
	}

	b, err := yaml.Marshal(replaced)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (y *yrep) replace(in any) (any, error) {
	var err error
	for _, fn := range y.replaceFuncs {
		in, err = fn(in)
		if err != nil && !errors.Is(err, ErrNotReplaced) {
			return nil, err
		}
	}

	switch v := in.(type) {
	case yaml.MapSlice:
		for i, item := range v {
			replaced, err := y.replace(item.Value)
			if err != nil && !errors.Is(err, ErrNotReplaced) {
				return nil, err
			}
			v[i].Value = replaced
		}
		return v, nil
	case []any:
		for i, vv := range v {
			replaced, err := y.replace(vv)
			if err != nil && !errors.Is(err, ErrNotReplaced) {
				return nil, err
			}
			v[i] = replaced
		}
		return v, nil
	}

	return in, ErrNotReplaced
}
