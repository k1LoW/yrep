package yrep

import (
	"fmt"
	"testing"

	"github.com/goccy/go-yaml"
)

func TestApply(t *testing.T) {
	tests := []struct {
		in           string
		replaceFuncs []ReplaceFunc
		want         string
	}{
		{
			`a: 1
b: 2
c:
  c1: 3
  c2: 4
d:
- 5
- 6
`,
			nil,
			`a: 1
b: 2
c:
  c1: 3
  c2: 4
d:
- 5
- 6
`,
		},
		{
			`a: 1
b: 2
c:
  c1: 3
  c2: 4
d:
- 5
- 6
e:
- e1: 7
- e1: 8
  e2: 9
`,
			[]ReplaceFunc{numberToString},
			`a: "1"
b: "2"
c:
  c1: "3"
  c2: "4"
d:
- "5"
- "6"
e:
- e1: "7"
- e1: "8"
  e2: "9"
`,
		},
		{
			`a: 1
b: 2
c:
  c1: 3
  c2: 4
d:
- 5
- 6
e:
- e1: 7
- e1: 8
  e2: 9
`,
			[]ReplaceFunc{numberToString, strOneToNumber},
			`a: 1
b: "2"
c:
  c1: "3"
  c2: "4"
d:
- "5"
- "6"
e:
- e1: "7"
- e1: "8"
  e2: "9"
`,
		},
		{
			`a: 1
b: 2
c:
  c1: 1
  c2: 2
d:
- 1
- 2
`,
			[]ReplaceFunc{oneToMap},
			`a:
  one: 1
b: 2
c:
  c1:
    one: 1
  c2: 2
d:
- one: 1
- 2
`,
		},
		{
			`a: 1
b: 2
c:
  b: 3
  d: 4
`,
			[]ReplaceFunc{deleteKeyB},
			`a: 1
c:
  d: 4
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := Apply([]byte(tt.in), tt.replaceFuncs...)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if string(got) != tt.want {
				t.Errorf("got\n-----\n%v\n\nwant\n-----\n%v\n", string(got), tt.want)
			}
		})
	}
}

func numberToString(in any) (any, error) {
	switch v := in.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v), nil
	default:
		return in, ErrNotReplaced
	}
}

func strOneToNumber(in any) (any, error) {
	switch v := in.(type) {
	case string:
		if v == "1" {
			return 1, nil
		}
		return in, ErrNotReplaced
	default:
		return in, ErrNotReplaced
	}
}

func oneToMap(in any) (any, error) {
	switch v := in.(type) {
	case uint64:
		if v == 1 {
			return map[string]any{"one": v}, nil
		}
		return in, ErrNotReplaced
	default:
		return in, ErrNotReplaced
	}
}

func deleteKeyB(in any) (any, error) {
	v, ok := in.(yaml.MapSlice)
	if !ok {
		return in, ErrNotReplaced
	}
	var deleted yaml.MapSlice
	found := false
	for i, item := range v {
		if item.Key == "b" {
			found = true
			continue
		}
		deleted = append(deleted, v[i])
	}
	if !found {
		return in, ErrNotReplaced
	}
	return deleted, nil
}
