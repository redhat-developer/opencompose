package v1

import (
	"reflect"
	"testing"

	"github.com/redhat-developer/opencompose/pkg/object"
	"gopkg.in/yaml.v2"
)

func TestPortMapping_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		Succeed bool
		RawPort string
		Port    *PortMapping
	}{
		{true, "5000", &PortMapping{ContainerPort: 5000, ServicePort: 5000}},
		{true, "5000:80", &PortMapping{ContainerPort: 5000, ServicePort: 80}},
		{true, "", &PortMapping{}}, // UnmarshalYAML won't be even called for empty strings
		{false, "x5000", nil},
		{false, "5000:", nil},
		{false, "x5000:", nil},
		{false, ":80", nil},
		{false, ":80x", nil},
		{false, "x:80x", nil},
		{false, "x:80", nil},
		{false, ":8080:80", nil},
		{false, "x:8080:80", nil},
		{false, "x:x8080:x80", nil},
		{false, ":8080:x80", nil},
		{false, ":8080:80:", nil},
		{false, "8080:80:", nil},
		{false, "8080:80::", nil},
		{false, ":", nil},
		{false, "::", nil},
		{false, ":::", nil},
		{false, "::::", nil},
		{false, ":::::", nil},
		{false, "::1:80", nil},
		{false, "::1:8080:80", nil},
		{false, "::1:8080:80:", nil},
	}

	for _, tt := range tests {
		var pm PortMapping
		err := yaml.Unmarshal([]byte(tt.RawPort), &pm)
		if err != nil {
			if tt.Succeed {
				t.Errorf("Failed to unmarshal %q: %s", tt.RawPort, err)
			}
			continue
		}

		if !tt.Succeed {
			t.Errorf("Expected %#v to fail!", tt.RawPort)
			continue
		}

		if !reflect.DeepEqual(pm, *tt.Port) {
			t.Errorf("Expected %#v, got %#v", *tt.Port, pm)
			continue
		}
	}
}

func TestEnvVariable_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		Succeed   bool
		RawEnvVar string
		EnvVar    *EnvVariable
	}{
		{true, "'KEY=value string '", &EnvVariable{Key: "KEY", Value: "value string "}},
		{true, "KEY= value", &EnvVariable{Key: "KEY", Value: " value"}},
		{true, "KEY =value", &EnvVariable{Key: "KEY", Value: "value"}},
		{true, "KEY==value", &EnvVariable{Key: "KEY", Value: "=value"}},
		{true, "KEY=", &EnvVariable{Key: "KEY", Value: ""}},
		{false, "KEY", nil},
		{false, "=KEYvalue", nil},
		{false, "=KEY=value", nil},
		{false, "=KEY=value=", nil},
	}

	for _, tt := range tests {
		var envVar EnvVariable
		err := yaml.Unmarshal([]byte(tt.RawEnvVar), &envVar)
		if err != nil {
			if tt.Succeed {
				t.Errorf("Failed to unmarshal %#v; error %#v", tt.RawEnvVar, err)
			}
			continue
		}

		if !tt.Succeed {
			t.Errorf("Expected %#v to fail!", tt.RawEnvVar)
			continue
		}

		if !reflect.DeepEqual(envVar, *tt.EnvVar) {
			t.Errorf("Expected %#v, got %#v", *tt.EnvVar, envVar)
			continue
		}
	}
}

func TestDecoder_Decode(t *testing.T) {
	// TODO: make better tests w.r.t excess keys in all possible places
	// TODO: add checking for proper error because tests can fail for other than expected reasons
	tests := []struct {
		Succeed     bool
		File        string
		OpenCompose *object.OpenCompose
	}{
		{
			true, `
version: 0.1-dev
services:
- name: frontend
  containers:
  - image: tomaskral/kompose-demo-frontend:test
    env:
    - KEY=value
    - KEY2=value2
    ports:
    - port: 5000:80
    - port: 5001:81
volumes:
- name: data
  size: 1Gi
  mode: ReadWriteOnce
`,
			&object.OpenCompose{
				Version: Version,
				Services: []object.Service{
					{
						Name: "frontend",
						Containers: []object.Container{
							{
								Image: "tomaskral/kompose-demo-frontend:test",
								Environment: []object.EnvVariable{
									{
										Key:   "KEY",
										Value: "value",
									},
									{
										Key:   "KEY2",
										Value: "value2",
									},
								},
								Ports: []object.Port{
									{
										Port: object.PortMapping{
											ContainerPort: 5000,
											ServicePort:   80,
										},
									},
									{
										Port: object.PortMapping{
											ContainerPort: 5001,
											ServicePort:   81,
										},
									},
								},
							},
						},
					},
				},
				Volumes: []object.Volume{
					{
						Name: "data",
						Size: "1Gi",
						Mode: "ReadWriteOnce",
					},
				},
			},
		},
		{
			false, `
version: 0.1-dev
services:
- name: frontend
  containers:
  - image: tomaskral/kompose-demo-frontend:test
    env:
    - KEY=value
    - KEY2=value2
    ports:
    - port: 5000:80
    - port: 5001:81
  - EXCESSKEY: some value
`,
			nil,
		},
		{
			false, `
version: 0.1-dev
services:
- name: frontend
  containers:
  - image: tomaskral/kompose-demo-frontend:test
	env:
	- KEY=value
	- KEY2=value2
	ports:
	- port: 5000:80
	- port: 5001:81
volumes:
- name: data
  size: 1Gi
  mode: ReadWriteOnce
  EXCESSKEY: some value
`,
			nil,
		},
		{
			false, `
version: 0.1-dev
services: []
volumes: []
EXCESSKEY: some value
`,
			nil,
		},
		{
			true, `
version: 0.1-dev
services:
- name: frontend
  containers:
  - image: tomaskral/kompose-demo-frontend:test
    env:
    - KEY=value
    - KEY2=value2
    ports:
    - port: 5000:80
    - port: 5001:81
volumes: []
`,
			&object.OpenCompose{
				Version: Version,
				Services: []object.Service{
					{
						Name: "frontend",
						Containers: []object.Container{
							{
								Image: "tomaskral/kompose-demo-frontend:test",
								Environment: []object.EnvVariable{
									{
										Key:   "KEY",
										Value: "value",
									},
									{
										Key:   "KEY2",
										Value: "value2",
									},
								},
								Ports: []object.Port{
									{
										Port: object.PortMapping{
											ContainerPort: 5000,
											ServicePort:   80,
										},
									},
									{
										Port: object.PortMapping{
											ContainerPort: 5001,
											ServicePort:   81,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			false, `
version: 0.1-dev
services: []
volumes: []
`,
			nil,
		},
		{
			false,
			"",
			nil,
		},
	}

	for _, tt := range tests {
		data := []byte(tt.File)
		openCompose, err := (&Decoder{}).Decode(data)
		if err != nil {
			if tt.Succeed {
				t.Errorf("Failed to unmarshal %#v; error %#v", tt.File, err)
			}
			continue
		}

		if !tt.Succeed {
			t.Errorf("Expected %#v to fail!", tt.File)
			continue
		}

		if !reflect.DeepEqual(openCompose, tt.OpenCompose) {
			t.Errorf("Expected: %#v; got: %#v", tt.OpenCompose, openCompose)
			continue
		}
	}

}
