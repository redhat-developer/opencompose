package object

import (
	"fmt"

	"strings"

	"path"

	"k8s.io/client-go/pkg/api/resource"
	"k8s.io/client-go/pkg/util/validation"
)

type PortMapping struct {
	ContainerPort int
	ServicePort   int
}

type PortType int

const (
	PortType_Internal PortType = iota
	PortType_External
)

type Port struct {
	Port PortMapping
	Type PortType
	Host *string
	Path string
}

type EnvVariable struct {
	Key   string
	Value string
}

type Mount struct {
	VolumeRef     string
	MountPath     string
	VolumeSubPath string
	ReadOnly      bool
}

type Labels map[string]string

type Container struct {
	Image       string
	Environment []EnvVariable
	Ports       []Port
	Mounts      []Mount
}

type EmptyDirVolume struct {
	Name string
}

type Service struct {
	Name            string
	Containers      []Container
	Replicas        *int32
	EmptyDirVolumes []EmptyDirVolume
	Labels          Labels
}

type Volume struct {
	Name         string
	Size         string
	AccessMode   string
	StorageClass *string
}

type OpenCompose struct {
	Version  string
	Services []Service
	Volumes  []Volume
}

// Given the name of 'emptyDirVolume' this function searches
// if service receiver has that 'emptyDirVolume'
func (s *Service) EmptyDirVolumeExists(name string) bool {
	for _, emptyDirVolume := range s.EmptyDirVolumes {
		if name == emptyDirVolume.Name {
			return true
		}
	}
	return false
}

// Given name of root level 'volume' this function searches
// if opencompose receiver has that 'volume'.
func (o *OpenCompose) VolumeExists(name string) bool {
	for _, volume := range o.Volumes {
		if name == volume.Name {
			return true
		}
	}
	return false
}

// Documentation about the valid identifiers can be found at
// https://github.com/kubernetes/community/blob/master/contributors/design-proposals/identifiers.md
func validateName(name string) error {
	if errs := validation.IsDNS1123Subdomain(name); len(errs) != 0 {
		return fmt.Errorf("%s", strings.Join(errs, "\n"))
	}
	return nil
}

func (e *EnvVariable) validate() error {
	// TODO: add more validation tests besides checking for '='
	if strings.Contains(e.Key, "=") {
		return fmt.Errorf("Illegal character '=' in environment variable key: %v", e.Key)
	}

	if strings.Contains(e.Value, "=") {
		return fmt.Errorf("Illegal character '=' in environment variable value: %v", e.Value)
	}
	return nil
}

func (m *Mount) validate() error {

	// validate volumeRef
	if err := validateName(m.VolumeRef); err != nil {
		return fmt.Errorf("mount %q: invalid name, %v", m.VolumeRef, err)
	}

	// mountPath should be absolute
	if !path.IsAbs(m.MountPath) {
		return fmt.Errorf("mount %q: mountPath %q: is not absolute path", m.VolumeRef, m.MountPath)
	}

	// validate volumeSubPath
	// TODO: if there is someway to do it

	return nil
}

func (c *Container) validate() error {

	// validate image name
	// TODO: implement me
	// validate Ports
	// TODO: implement me

	for _, env := range c.Environment {
		if err := env.validate(); err != nil {
			return fmt.Errorf("failed to validate environment variable: %v", err)
		}
	}

	// validate Mounts
	allMounts := make(map[string]string)
	for _, mount := range c.Mounts {
		if err := mount.validate(); err != nil {
			return fmt.Errorf("failed to validate mount: %v", err)
		}

		// mountPath should not collide, which means you should not do multiple mounts in same place
		if v, ok := allMounts[mount.MountPath]; ok {
			return fmt.Errorf("mount %q: mountPath %q: cannot have same mountPath as %q", mount.VolumeRef, mount.MountPath, v)
		}
		allMounts[mount.MountPath] = mount.VolumeRef
	}
	return nil
}

func (s *Service) validate() error {
	// validate service name, like it cannot have underscores, etc.
	if err := validateName(s.Name); err != nil {
		return fmt.Errorf("invalid name, %v", err)
	}

	// validate containers
	for cno, cnt := range s.Containers {
		if err := cnt.validate(); err != nil {
			return fmt.Errorf("container#%d: %v", cno+1, err)
		}
	}

	// validate replicas
	if s.Replicas != nil && *s.Replicas < 0 {
		return fmt.Errorf("%s", "'replicas' can't be negative")
	}

	// validate emptyDirVolume
	for _, e := range s.EmptyDirVolumes {
		if err := validateName(e.Name); err != nil {
			return fmt.Errorf("emptyDirVolume %q: invalid name, %v", e.Name, err)
		}
	}

	// validate label values
	for _, v := range s.Labels {
		errString := validation.IsValidLabelValue(v)
		if errString != nil {
			return fmt.Errorf("Invalid label value: %v", errString)
		}
	}

	return nil
}

func validateVolumeMode(volumeMode string) error {
	switch volumeMode {
	case "ReadWriteOnce", "ReadOnlyMany", "ReadWriteMany":
	default:
		return fmt.Errorf("invalid accessMode: %q, must be either %q, %q or %q", volumeMode, "ReadWriteOnce", "ReadOnlyMany", "ReadWriteMany")
	}
	return nil
}

func (v *Volume) validate() error {
	// validate volume name
	if err := validateName(v.Name); err != nil {
		return fmt.Errorf("invalid name, %v", err)
	}

	// validate volume size
	if _, err := resource.ParseQuantity(v.Size); err != nil {
		return fmt.Errorf("size %q: %v", v.Size, err)
	}

	// validate volume access mode
	if err := validateVolumeMode(v.AccessMode); err != nil {
		return err
	}

	if v.StorageClass != nil {
		if err := validateName(*v.StorageClass); err != nil {
			return fmt.Errorf("storageClass %q: invalid name, %v", *v.StorageClass, err)
		}
	}

	return nil
}

// Does high level (mostly semantic) validation of OpenCompose
// (e.g. it checks internal object references)
func (o *OpenCompose) Validate() error {
	// validating services
	for _, service := range o.Services {
		if err := service.validate(); err != nil {
			return fmt.Errorf("service %q: %v", service.Name, err)
		}

		// validate if the mounts are specified in root level volumes
		// or emptydirvolumes, error out if not found anywhere
		for cno, container := range service.Containers {
			for _, mount := range container.Mounts {
				if !o.VolumeExists(mount.VolumeRef) && !service.EmptyDirVolumeExists(mount.VolumeRef) {
					return fmt.Errorf("volume mount %q in service %q in container#%d does not correspond to either 'root level volume' or 'emptydir volume'",
						mount.VolumeRef, service.Name, cno+1)
				}
			}
		}
	}

	// validate volumes
	for _, volume := range o.Volumes {
		if err := volume.validate(); err != nil {
			return fmt.Errorf("volume %q: %v", volume.Name, err)
		}
	}

	return nil
}
