package v1

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/tnozicka/opencompose/pkg/encoding/util"
	"github.com/tnozicka/opencompose/pkg/object"
)

type Port string

func (raw *Port) Unmarshal() (*object.Port, error) {
	// TODO: support binding by address
	p := &object.Port{}
	var err error

	sliceBySlash := strings.Split(string(*raw), "/")
	switch l := len(sliceBySlash); l {
	case 1:
		// no protocol specified; we want to retain the information about it being empty
	case 2:
		// [1] - protocol
		p.Protocol = sliceBySlash[1]
		switch p.Protocol {
		case "tcp":
			// ok
		case "udp":
			// ok
		case "":
			return nil, fmt.Errorf("failed to unmarshal port '%s': invalid format (no protocol, but protocol separator specified)", raw)
		default:
			return nil, fmt.Errorf("failed to unmarshal port '%s': invalid protocol '%s'", raw, p.Protocol)
		}

	default:
		return nil, fmt.Errorf("failed to unmarshal port '%s': unable to parse protocol", raw)
	}

	sliceByColumn := strings.Split(sliceBySlash[0], ":")
	switch l := len(sliceByColumn); l {
	case 3:
		// [0]=ContainerPort [1]=HostPort [2]=ServicePort
		p.ServicePort, err = strconv.Atoi(sliceByColumn[2])
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal port (service) '%s': %s", raw, err)
		}
		fallthrough
	case 2:
		// [0]=ContainerPort [1]=HostPort==ServicePort
		p.HostPort, err = strconv.Atoi(sliceByColumn[1])
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal port (host) '%s': %s", raw, err)
		}
		fallthrough
	case 1:
		// [0] ContainerPort==HostPort==ServicePort
		p.ContainerPort, err = strconv.Atoi(sliceByColumn[0])
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal port (container) '%s': %s", raw, err)
		}
	case 0:
		return nil, fmt.Errorf("failed to unmarshal port '%s': no items found", raw)
	default:
		return nil, fmt.Errorf("failed to unmarshal port '%s': too many items (%d)", raw, l)
	}

	return p, nil
}

type Mapping struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
	Port Port   `json:"port,omitempty"`
}

func (raw *Mapping) Unmarshal() (*object.Mapping, error) {
	// name, type
	m := &object.Mapping{
		Name: raw.Name,
		Type: raw.Type,
	}

	// validate name
	if err := util.ValidateResourceName(m.Name); err != nil {
		return nil, fmt.Errorf("service name: %s", err)
	}

	// TODO: validate type

	// port
	var err error
	port, err := raw.Port.Unmarshal()
	if err != nil {
		return nil, err
	}
	m.Port = *port

	return m, nil
}

type EnvVariable string

func (raw *EnvVariable) Unmarshal() (*object.EnvVariable, error) {
	splitSlice := strings.SplitN(string(*raw), "=", 2)

	if len(splitSlice) != 2 {
		return nil, fmt.Errorf("failed to unmarshal environment variable '%s'", string(*raw))
	}

	if splitSlice[0] == "" {
		return nil, fmt.Errorf("failed to unmarshal environment variable '%s': no key", string(*raw))
	}

	e := &object.EnvVariable{
		Key:   strings.TrimSpace(splitSlice[0]),
		Value: splitSlice[1],
	}

	return e, nil
}

type Container struct {
	Name     string        `json:"name"`
	Image    string        `json:"image"`
	Env      []EnvVariable `json:"env,omitempty"`
	Mappings []Mapping     `json:"mappings,omitempty"`
}

func (raw *Container) Unmarshal() (*object.Container, error) {
	// name, image
	c := &object.Container{
		Name:  raw.Name,
		Image: raw.Image,
	}

	// validate name
	if err := util.ValidateResourceName(c.Name); err != nil {
		return nil, fmt.Errorf("service name: %s", err)
	}

	// TODO: validate image ref

	// environment
	for _, rawEnv := range raw.Env {
		envVar, err := rawEnv.Unmarshal()
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal environment variable: %s", err)
		}
		c.Environment = append(c.Environment, *envVar)
	}

	// mappings
	for _, rawMapping := range raw.Mappings {
		envMapping, err := rawMapping.Unmarshal()
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal mapping: %s", err)
		}
		c.Mappings = append(c.Mappings, *envMapping)
	}

	return c, nil
}

type Service struct {
	Name       string      `json:"name"`
	Containers []Container `json:"containers"`
}

func (raw *Service) Unmarshal() (*object.Service, error) {
	// name
	s := &object.Service{
		Name: raw.Name,
	}

	// validate name
	if err := util.ValidateResourceName(s.Name); err != nil {
		return nil, fmt.Errorf("service name: %s", err)
	}

	// containers
	for _, rawContainer := range raw.Containers {
		container, err := rawContainer.Unmarshal()
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal container: %s", err)
		}
		s.Containers = append(s.Containers, *container)
	}

	return s, nil
}

type Volume struct {
	Name string `json:"name"`
	Size string `json:"size,omitempty"`
	Mode string `json:"mode,omitempty"`
}

func (raw *Volume) Unmarshal() (*object.Volume, error) {
	v := object.Volume(*raw)

	// validate name
	if err := util.ValidateResourceName(v.Name); err != nil {
		return nil, fmt.Errorf("volume name: %s", err)
	}

	// TODO: validate size

	// TODO: validate mode

	return &v, nil
}

type OpenCompose struct {
	Version  int       `json:"version,omitempty"`
	Services []Service `json:"services"`
	Volumes  []Volume  `json:"volumes,omitempty"`
}

func (raw *OpenCompose) Unmarshal() (*object.OpenCompose, error) {
	// version
	if raw.Version != 1 {
		return nil, fmt.Errorf("unmarshal OpenCompose (version 1): unsupported version %d", raw.Version)
	}
	o := &object.OpenCompose{
		Version: raw.Version,
	}

	// services
	for _, rawService := range raw.Services {
		service, err := rawService.Unmarshal()
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal service: %s", err)
		}
		o.Services = append(o.Services, *service)
	}

	// volumes
	for _, rawVolume := range raw.Volumes {
		volume, err := rawVolume.Unmarshal()
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal volume: %s", err)
		}
		o.Volumes = append(o.Volumes, *volume)
	}

	return o, nil
}

type Decoder struct{}

// Unmarshals OpenCompose file into object.OpenCompose struct
// It does not add any additional (or default) values so it can be marshaled back
// to give the same result
// Currently it does not check for excess fields - this is an issue of yaml library
// and there is already accepted proposal for Go 1.9 about json alternative
// https://github.com/golang/go/issues/15314 so hopefully yaml gets something similar
// otherwise we have to ditch the decoder and write our own using reflect
func (u *Decoder) Unmarshal(data []byte) (*object.OpenCompose, error) {
	rawOpenCompose := &OpenCompose{}
	// TODO: check for excess fields (see above)
	err := yaml.Unmarshal(data, rawOpenCompose)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal OpenCompose: %s", err)
	}

	// OpenCompose spec isn't just pure yaml. It has other formats embedded inside
	// and we need to unmarshal them as well
	openCompose, err := rawOpenCompose.Unmarshal()
	if err != nil {
		return nil, err
	}

	return openCompose, nil
}
