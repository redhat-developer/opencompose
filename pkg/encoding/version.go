package encoding

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

type Version struct {
	Version string `yaml:"version,omitempty"`
}

func GetVersion(data []byte) (string, error) {
	var v Version
	err := yaml.Unmarshal(data, &v)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal OpenCompose version: %s", err)
	}
	return v.Version, nil
}
