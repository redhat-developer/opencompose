package object

import "fmt"

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

type Container struct {
	Image       string
	Environment []EnvVariable
	Ports       []Port
}

type Service struct {
	Name       string
	Containers []Container
	Replicas   *int32
}

type Volume struct {
	Name string
	Size string
	Mode string
}

type OpenCompose struct {
	Version  string
	Services []Service
	Volumes  []Volume
}

func (s *Service) Validate() error {
	// validate service name, like it cannot have underscores, etc.

	// validate containers

	// validate replicas
	if s.Replicas != nil && *s.Replicas < 0 {
		return fmt.Errorf("Replica count is negative in service: %q", s.Name)
	}

	return nil
}

// Does high level (mostly semantic) validation of OpenCompose
// (e.g. it checks internal object references)
func (o *OpenCompose) Validate() error {
	// validating services
	for _, service := range o.Services {
		if err := service.Validate(); err != nil {
			return err
		}
	}
	return nil
}
