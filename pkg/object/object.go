package object

type PortMapping struct {
	ContainerPort int
	HostPort      int
	ServicePort   int
}

type Port struct {
	Port PortMapping
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

// Does high level (mostly semantic) validation of OpenCompose
// (e.g. it checks internal object references)
func (o *OpenCompose) Validate() error {
	// TODO: implement
	return nil
}
