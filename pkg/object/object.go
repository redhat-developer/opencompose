package object

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
