package encoding

import (
	"fmt"
	"github.com/ghodss/yaml"
)

type Version struct {
	Version int `json:"version,omitempty"`
}

func GetVersion(data []byte) (int, error) {
	var v Version
	err := yaml.Unmarshal(data, &v)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal OpenCompose version: %s", err)
	}
	return v.Version, nil
}
