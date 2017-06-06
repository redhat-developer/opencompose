package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/redhat-developer/opencompose/pkg/encoding/util"
	"github.com/redhat-developer/opencompose/pkg/goutil"
	"github.com/redhat-developer/opencompose/pkg/object"
	"gopkg.in/yaml.v2"

	api_v1 "k8s.io/client-go/pkg/api/v1"
)

const (
	Version = "0.1-dev" // TODO: replace with "1" once we reach that point
)

type ResourceName string

func (rn *ResourceName) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var name string
	err := unmarshal(&name)
	if err != nil {
		return err
	}

	// validate name
	if err := util.ValidateResourceName(name); err != nil {
		return fmt.Errorf("failed to unmarshal ResourceName - invalid name: %s", err)
	}

	*rn = ResourceName(name)

	return nil
}

type PortMapping struct {
	ContainerPort int
	ServicePort   int
}

func (pm *PortMapping) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal(&s)
	if err != nil {
		return err
	}

	sliceByColumn := strings.Split(s, ":")
	l := len(sliceByColumn)
	switch l {
	case 2:
		// [0]=ContainerPort [1]=ServicePort
		pm.ServicePort, err = strconv.Atoi(sliceByColumn[l-1])
		if err != nil {
			return fmt.Errorf("failed to unmarshal port (service) %q: %s", s, err)
		}
		fallthrough
	case 1:
		// [0] ContainerPort
		pm.ContainerPort, err = strconv.Atoi(sliceByColumn[0])
		if err != nil {
			return fmt.Errorf("failed to unmarshal port (container) %q: %s", s, err)
		}
	case 0:
		return fmt.Errorf("failed to unmarshal port %q: no items found", s)
	default:
		return fmt.Errorf("failed to unmarshal port %q: too many items (%d)", s, l)
	}

	// Fill in default ports by deduction
	switch l {
	case 1:
		// [0] ContainerPort==ServicePort
		pm.ServicePort = pm.ContainerPort
	}

	return nil
}

type PortType object.PortType

func (pt *PortType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal(&s)
	if err != nil {
		return err
	}

	switch s {
	case "internal":
		*pt = PortType(object.PortType_Internal)
	case "external":
		*pt = PortType(object.PortType_External)
	default:
		return fmt.Errorf("failed to unmarshal port type: invalid port type %q", s)
	}

	return nil
}

// Fully qualified domain name as defined by RFC 3986
type Fqdn string

// TODO: Add Fqdn unmarshalling to validate it

// An extended POSIX regex as defined by IEEE Std 1003.1, (i.e this follows the egrep/unix syntax, not the perl syntax)
type PathRegex string

// TODO: Add PathRegex unmarshalling to validate it

type Port struct {
	Port PortMapping `yaml:"port"`
	Type PortType    `yaml:"type,omitempty"`
	Host *Fqdn       `yaml:"host,omitempty"`
	Path *PathRegex  `yaml:"path,omitempty"`
}

func (v *Port) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type PortAlias Port
	var st struct {
		PortAlias `yaml:",inline"`
		Leftovers map[string]interface{} `yaml:",inline"` // Catches all undefined fields and must be empty after parsing.
	}
	err := unmarshal(&st)
	if err != nil {
		return err
	}

	if len(st.Leftovers) > 0 {
		return util.NewExcessKeysErrorFromMap("Port", st.Leftovers)
	}

	*v = Port(st.PortAlias)

	// Setting "path" requires specifying "host"
	if v.Host == nil && v.Path != nil {
		return errors.New("failed to unmarshal port: 'host' not specified: setting 'path' requires specifying 'host'")
	}

	// If there is no path specified it implies ""
	if v.Host != nil && v.Path == nil {
		v.Path = new(PathRegex)
	}

	return nil
}

type EnvVariable struct {
	Key   string `yaml:"name"`
	Value string `yaml:"value"`
}

func (raw *EnvVariable) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type EnvVariableAlias EnvVariable
	var st struct {
		EnvVariableAlias `yaml:",inline"`
		Leftovers        map[string]interface{} `yaml:",inline"` // Catches all undefined fields and must be empty after parsing.
	}
	err := unmarshal(&st)
	if err != nil {
		return err
	}

	if len(st.Leftovers) > 0 {
		return util.NewExcessKeysErrorFromMap("Env", st.Leftovers)
	}

	*raw = EnvVariable(st.EnvVariableAlias)

	return nil
}

type Labels object.Labels

func (lb *Labels) UnmarshalYAML(unmarshal func(interface{}) error) error {
	labelMap := make(map[string]string)
	if err := unmarshal(&labelMap); err != nil {
		return err
	}

	*lb = Labels(labelMap)

	return nil
}

type ImageRef string

// FIXME: implement ImageRef unmarshalling

type Mount struct {
	VolumeRef ResourceName `yaml:"volumeRef"`
	MountPath string       `yaml:"mountPath"`
	// these are optional fields so making them as pointer because it helps
	// to identify whether these fields were given by user or not
	// if these are not pointer then it is hard to identify what was given
	// by user and what is the default value
	VolumeSubPath *string `yaml:"volumeSubPath,omitempty"`
	ReadOnly      *bool   `yaml:"readOnly,omitempty"`
}

func (m *Mount) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type MountAlias Mount
	var st struct {
		MountAlias `yaml:",inline"`
		Leftovers  map[string]interface{} `yaml:",inline"` // Catches all undefined fields and must be empty after parsing.
	}
	if err := unmarshal(&st); err != nil {
		return err
	}

	if len(st.Leftovers) > 0 {
		return util.NewExcessKeysErrorFromMap("Mount", st.Leftovers)
	}

	*m = Mount(st.MountAlias)

	return nil
}

type Health struct {
	// Data holder for ReadinessProbe while parsing
	// Data from the yaml file will be read into this field
	ReadinessProbeData interface{} `yaml:"readinessProbe,omitempty"`
	// After certain processing the data in ReadinessProbeData
	// will be populated into ReadinessProbe for further use
	ReadinessProbe *api_v1.Probe

	LivenessProbeData interface{} `yaml:"livenessProbe,omitempty"`
	LivenessProbe     *api_v1.Probe
}

// If given an interface which has JSONified data of type Probe
// this function will read the interface and give concrete
// data strcuture pointer.
func interfaceToProbe(i interface{}) (*api_v1.Probe, error) {
	i = util.InterfaceToJSON(i)

	var b []byte
	var err error
	if b, err = json.Marshal(i); err != nil {
		return nil, fmt.Errorf("error: marshalling interface to bytes: %v", err)
	}
	var p api_v1.Probe
	if err = json.Unmarshal(b, &p); err != nil {
		return nil, fmt.Errorf("error: unmarshalling bytes to Probe: %v", err)
	}
	return &p, nil
}

func (h *Health) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type HealthAlias Health
	var st struct {
		HealthAlias `yaml:",inline"`
		Leftovers   map[string]interface{} `yaml:",inline"` // Catches all undefined fields and must be empty after parsing.
	}

	err := unmarshal(&st)
	if err != nil {
		return err
	}

	if len(st.Leftovers) > 0 {
		return util.NewExcessKeysErrorFromMap("Health", st.Leftovers)
	}

	*h = Health(st.HealthAlias)

	// extract the data from interface into concrete data type 'Probe'
	if h.ReadinessProbeData != nil {
		h.ReadinessProbe, err = interfaceToProbe(h.ReadinessProbeData)
		if err != nil {
			return fmt.Errorf("readinessProbe: %v", err)
		}
		h.ReadinessProbeData = interface{}(nil)
	}

	// extract the data from interface into concrete data type 'Probe'
	if h.LivenessProbeData != nil {
		h.LivenessProbe, err = interfaceToProbe(h.LivenessProbeData)
		if err != nil {
			return fmt.Errorf("livenessProbe: %v", err)
		}
		h.LivenessProbeData = interface{}(nil)
	}

	// TODO: Right now we have no way of finding if the excess keys are given
	// by the user, since we are doing the whole conversion from YAML to JSON
	// and then parsing it into the internal k8s structs

	return nil
}

type Container struct {
	Image  ImageRef      `yaml:"image"`
	Env    []EnvVariable `yaml:"env,omitempty"`
	Ports  []Port        `yaml:"ports,omitempty"`
	Mounts []Mount       `yaml:"mounts,omitempty"`
	Health *Health       `yaml:"health,omitempty"`
}

func (c *Container) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type ContainerAlias Container
	var st struct {
		ContainerAlias `yaml:",inline"`
		Leftovers      map[string]interface{} `yaml:",inline"` // Catches all undefined fields and must be empty after parsing.
	}
	err := unmarshal(&st)
	if err != nil {
		return err
	}

	if len(st.Leftovers) > 0 {
		return util.NewExcessKeysErrorFromMap("Container", st.Leftovers)
	}

	*c = Container(st.ContainerAlias)

	return nil
}

type EmptyDirVolume struct {
	Name ResourceName `yaml:"name"`
}

func (e *EmptyDirVolume) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type EmptyDirVolumeAlias EmptyDirVolume
	var st struct {
		EmptyDirVolumeAlias `yaml:",inline"`
		Leftovers           map[string]interface{} `yaml:",inline"` // Catches all undefined fields and must be empty after parsing.
	}
	if err := unmarshal(&st); err != nil {
		return err
	}

	if len(st.Leftovers) > 0 {
		return util.NewExcessKeysErrorFromMap("EmptyDirVolume", st.Leftovers)
	}

	*e = EmptyDirVolume(st.EmptyDirVolumeAlias)

	return nil
}

type Service struct {
	Name            ResourceName     `yaml:"name"`
	Containers      []Container      `yaml:"containers"`
	Replicas        *int32           `yaml:"replicas,omitempty"`
	EmptyDirVolumes []EmptyDirVolume `yaml:"emptyDirVolumes,omitempty"`
	Labels          Labels           `yaml:"labels,omitempty"`
}

func (s *Service) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type ServiceAlias Service
	var st struct {
		ServiceAlias `yaml:",inline"`
		Leftovers    map[string]interface{} `yaml:",inline"` // Catches all undefined fields and must be empty after parsing.
	}
	err := unmarshal(&st)
	if err != nil {
		return err
	}

	if len(st.Leftovers) > 0 {
		return util.NewExcessKeysErrorFromMap("Service", st.Leftovers)
	}

	*s = Service(st.ServiceAlias)

	return nil
}

type Volume struct {
	Name         ResourceName  `yaml:"name"`
	Size         string        `yaml:"size"`
	AccessMode   string        `yaml:"accessMode"`
	StorageClass *ResourceName `yaml:"storageClass,omitempty"`
}

func (v *Volume) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type VolumeAlias Volume
	var st struct {
		VolumeAlias `yaml:",inline"`
		Leftovers   map[string]interface{} `yaml:",inline"` // Catches all undefined fields and must be empty after parsing.
	}
	err := unmarshal(&st)
	if err != nil {
		return err
	}

	if len(st.Leftovers) > 0 {
		return util.NewExcessKeysErrorFromMap("Volume", st.Leftovers)
	}

	*v = Volume(st.VolumeAlias)

	return nil
}

type VersionString string

func (vs *VersionString) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v string
	err := unmarshal(&v)
	if err != nil {
		return err
	}

	if v != Version {
		return fmt.Errorf("can't unmarshal OpenCompose version - expected %q, got %q", Version, v)
	}

	*vs = VersionString(v)

	return nil
}

type OpenCompose struct {
	Version  VersionString `yaml:"version"`
	Services []Service     `yaml:"services"`
	Volumes  []Volume      `yaml:"volumes,omitempty"`
}

func (oc *OpenCompose) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type OpenComposeAlias OpenCompose
	var st struct {
		OpenComposeAlias `yaml:",inline"`
		Leftovers        map[string]interface{} `yaml:",inline"` // Catches all undefined fields and must be empty after parsing.
	}
	err := unmarshal(&st)
	if err != nil {
		return err
	}

	if len(st.Leftovers) > 0 {
		return util.NewExcessKeysErrorFromMap("OpenCompose", st.Leftovers)
	}

	*oc = OpenCompose(st.OpenComposeAlias)

	return nil
}

type Decoder struct{}

// Unmarshals OpenCompose file into object.OpenCompose struct
// It does not add any additional (or default) values so it can be marshaled back
// to give the same result
// Currently it does not check for excess fields - this is an issue of yaml library
// and there is already accepted proposal for Go 1.9 about json alternative
// https://github.com/golang/go/issues/15314 so hopefully yaml gets something similar
// otherwise we have to ditch the decoder and write our own using reflect
func (d *Decoder) Decode(data []byte) (*object.OpenCompose, error) {
	var v1 OpenCompose
	// TODO: check for excess fields (see above)
	err := yaml.Unmarshal(data, &v1)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal OpenCompose: %s", err)
	}

	// UnmarshalYAML can't check for empty values because in that case it won't get even called
	// We have to do it here manually
	err = util.ValidateRequiredFields(v1)
	if err != nil {
		return nil, err
	}

	// convert it from our version to internal definitions
	openCompose := &object.OpenCompose{
		Version: string(v1.Version),
	}

	// convert services
	for _, s := range v1.Services {
		os := object.Service{
			Name:   string(s.Name),
			Labels: object.Labels(s.Labels),
		}

		os.Replicas = s.Replicas

		// convert containers
		for _, c := range s.Containers {
			oc := object.Container{
				Image: string(c.Image),
			}

			// convert ports
			for _, p := range c.Ports {
				oc.Ports = append(oc.Ports, object.Port{
					Port: object.PortMapping{
						ContainerPort: p.Port.ContainerPort,
						ServicePort:   p.Port.ServicePort,
					},
					Type: object.PortType(p.Type),
					Host: (*string)(p.Host),
					Path: goutil.StringOrEmpty((*string)(p.Path)),
				})
			}

			// convert mounts
			for _, m := range c.Mounts {
				mount := object.Mount{
					VolumeRef: string(m.VolumeRef),
					MountPath: string(m.MountPath),
				}

				if m.VolumeSubPath != nil {
					mount.VolumeSubPath = string(*m.VolumeSubPath)
				}

				if m.ReadOnly != nil {
					mount.ReadOnly = *m.ReadOnly
				}

				oc.Mounts = append(oc.Mounts, mount)
			}

			if c.Health != nil {
				oc.Health.LivenessProbe = c.Health.LivenessProbe
				oc.Health.ReadinessProbe = c.Health.ReadinessProbe
			}

			// convert env
			for _, e := range c.Env {
				oc.Environment = append(oc.Environment, object.EnvVariable{
					Key:   e.Key,
					Value: e.Value,
				})
			}

			os.Containers = append(os.Containers, oc)
		}

		// Add emptyDirVolumes
		for _, emptydir := range s.EmptyDirVolumes {
			os.EmptyDirVolumes = append(os.EmptyDirVolumes, object.EmptyDirVolume{
				Name: string(emptydir.Name),
			})
		}

		openCompose.Services = append(openCompose.Services, os)
	}

	// convert volumes
	// TODO: remove the redundant sting conversion
	for _, v := range v1.Volumes {
		ov := object.Volume{
			Name:       string(v.Name),
			Size:       v.Size,
			AccessMode: v.AccessMode,
		}

		if v.StorageClass != nil {
			storageClass := string(*v.StorageClass)
			ov.StorageClass = &storageClass
		}

		openCompose.Volumes = append(openCompose.Volumes, ov)
	}

	return openCompose, nil
}
