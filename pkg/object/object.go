package object

import (
	"fmt"

	"strings"

	"path"

	"log"

	"github.com/redhat-developer/opencompose/pkg/transform/util"
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
	Key       string
	Value     *string
	SecretRef *string
}

type Mount struct {
	VolumeRef     *string
	MountPath     string
	VolumeSubPath string
	ReadOnly      bool
	SecretRef     *string
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

type Secret struct {
	Name string
	Data []SecretData
}

type SecretData struct {
	Key       string
	Plaintext *string
	Base64    *string
	File      *string
}

type OpenCompose struct {
	Version  string
	Services []Service
	Volumes  []Volume
	Secrets  []Secret
}

// This struct is not a part of the OpenCompose spec, but this is used to
// store the data input metadata which is used in other parts of the code
// to make decisions.
// For instance, this is being set in cmd.GetValidatedObject(), and getting
// in v1.Decode() to convert the relative path of the secret file from the
// OpenCompose file, and convert it to an absolute path
type Input struct {
	Data     []byte
	STDIN    bool
	URL      *string
	FilePath *string
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

	if e.Value != nil && strings.Contains(*e.Value, "=") {
		return fmt.Errorf("Illegal character '=' in environment variable value: %v", e.Value)
	}

	// make sure Value and SecretRef are not supplied at the same time
	// also, that at least one of the both should be specified
	// the following makes sure that exactly one of them exists at a given time
	if (e.SecretRef != nil) == (e.Value != nil) {
		return fmt.Errorf("Exactly one from 'value' or 'secretRef' must be specified for the environment variable key: %v", e.Key)
	}

	return nil
}

func (m *Mount) validate() error {

	// make sure mountRef and SecretRef are not supplied at the same time
	// also, that at least one of the both should be specified
	// the following makes sure that exactly one of them exists at a given time
	if (m.SecretRef != nil) == (m.VolumeRef != nil) {
		return fmt.Errorf("Exactly one from 'mountRef' or 'secretRef' must be specified for the mountPath: %v", m.MountPath)
	}

	if m.VolumeRef != nil {
		// validate volumeRef
		if err := validateName(*m.VolumeRef); err != nil {
			return fmt.Errorf("mount %q: invalid name, %v", *m.VolumeRef, err)
		}
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

		var mountType *string
		if mount.VolumeRef != nil {
			mountType = mount.VolumeRef
		}
		if mount.SecretRef != nil {
			mountType = mount.SecretRef
		}

		// mountPath should not collide, which means you should not do multiple mounts in same place
		if v, ok := allMounts[mount.MountPath]; ok {
			return fmt.Errorf("mount %q: mountPath %q: cannot have same mountPath as %q", *mountType, mount.MountPath, v)
		}
		allMounts[mount.MountPath] = *mountType

		// validate volumeSubPath
		// TODO: if there is someway to do it
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

func (sd *SecretData) validate() error {
	var count int
	if sd.Plaintext != nil {
		count++
	}
	if sd.Base64 != nil {
		count++
	}
	if sd.File != nil {
		count++
	}

	switch count {
	case 0:
		return fmt.Errorf("Please set one of plaintext, base64 or file field for the secret key: %v", sd.Key)
	case 2, 3:
		return fmt.Errorf("Only one of plaintext, base64 or file fields can be set at a time for the secret key: %v", sd.Key)
	default:
		return fmt.Errorf("Something went wrong with counting the secret fields for the secret key: %v", sd.Key)
	}
}

func validateSecretRef(sRef *string, oSec *[]Secret) error {
	// validate secretRef syntax
	secretDef, err := util.SecretRefToSecretDef(sRef)
	if err != nil {
		return fmt.Errorf("unable to verify secretRef syntax: %v", err)
	}

	if len(*oSec) == 0 {
		log.Printf("There are no root level secrets defined, assuming the corresponding secret exists in the cluster: %v", *sRef)
	} else {
		secretFound := false
		for _, secret := range *oSec {
			if secretDef.SecretName == secret.Name {
				secretFound = true

				dataKeyFound := false
				for _, secData := range secret.Data {
					if secretDef.DataKey == secData.Key {
						dataKeyFound = true
						break
					}
				}
				if !dataKeyFound {
					log.Printf("Root level secret name : %v found, but the provided data key : %v is missing, assuming the corresponding secret exists in the cluster", secret.Name, secretDef.DataKey)
				}
				break
			}
		}
		if !secretFound {
			log.Printf("secretRef %v does not correspond to any root level secret, assuming the corresponding secret exists in the cluster", *sRef)
		}
	}

	return nil
}

func (s *Secret) validate() error {

	if err := validateName(s.Name); err != nil {
		return fmt.Errorf("invalid secret name, %v", err)
	}

	for _, secretData := range s.Data {
		if err := secretData.validate(); err != nil {
			return fmt.Errorf("failed to validate secret key %v: %v", secretData.Key, err)
		}
	}

	return nil
}

// Does high level (mostly semantic) validation of OpenCompose
// (e.g. it checks internal object references)
func (o *OpenCompose) Validate() error {
	log.SetFlags(0)
	// validating services
	for _, service := range o.Services {
		if err := service.validate(); err != nil {
			return fmt.Errorf("service %q: %v", service.Name, err)
		}

		// validate if the mounts are specified in root level volumes
		// or emptydirvolumes, error out if not found anywhere
		for cno, container := range service.Containers {
			for _, mount := range container.Mounts {
				if mount.VolumeRef != nil && !o.VolumeExists(*mount.VolumeRef) && !service.EmptyDirVolumeExists(*mount.VolumeRef) {
					return fmt.Errorf("volume mount %q in service %q in container#%d does not correspond to any of 'root level volume' or 'emptydir volume'", *mount.VolumeRef, service.Name, cno+1)
				}

				if mount.SecretRef != nil {
					if err := validateSecretRef(mount.SecretRef, &o.Secrets); err != nil {
						return fmt.Errorf("Failed to validate secretRef: %v", *mount.SecretRef)
					}
				}
			}
			for _, env := range container.Environment {
				if env.SecretRef != nil {
					if err := validateSecretRef(env.SecretRef, &o.Secrets); err != nil {
						return fmt.Errorf("Failed to validate secretRef: %v", *env.SecretRef)
					}
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

	// validate root level secrets
	for _, secret := range o.Secrets {
		if err := secret.validate(); err != nil {
			return fmt.Errorf("failed to validate secret %v: %v", secret.Name, err)
		}
	}

	return nil
}
