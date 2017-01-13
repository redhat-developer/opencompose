package v1

import (
	"github.com/tnozicka/opencompose/pkg/object"
	"reflect"
	"testing"
)

func TestPortUnmarshal(t *testing.T) {
	tests := []struct {
		Succeed bool
		RawPort Port
		Port    *object.Port
	}{
		{true, "5000", &object.Port{ContainerPort: 5000}},
		{true, "5000:80", &object.Port{ContainerPort: 5000, HostPort: 80}},
		{true, "5000:8080:80", &object.Port{ContainerPort: 5000, HostPort: 8080, ServicePort: 80}},
		{true, "5000:8080:80/tcp", &object.Port{ContainerPort: 5000, HostPort: 8080, ServicePort: 80, Protocol: "tcp"}},
		{false, "", nil},
		{false, "x5000", nil},
		{false, "5000:8080:80/", nil},
		//{false, "5000:8080:80/tcp:90", nil},
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
		port, err := tt.RawPort.Unmarshal()
		if err != nil {
			if tt.Succeed {
				t.Errorf("Failed to unmarshal %q; error %q", tt.RawPort, err)
			}
			continue
		}

		if !tt.Succeed {
			t.Errorf("Expected %#v to fail!", tt.RawPort)
			continue
		}

		if !reflect.DeepEqual(port, tt.Port) {
			t.Errorf("Expected %#v, got %#v", tt.Port, port)
			continue
		}
	}
}

func TestEnvironmentVariableUnmarshal(t *testing.T) {
	tests := []struct {
		Succeed   bool
		RawEnvVar EnvVariable
		EnvVar    *object.EnvVariable
	}{
		{true, "KEY=value string ", &object.EnvVariable{Key: "KEY", Value: "value string "}},
		{true, "KEY= value", &object.EnvVariable{Key: "KEY", Value: " value"}},
		{true, "KEY =value", &object.EnvVariable{Key: "KEY", Value: "value"}},
		{true, "KEY==value", &object.EnvVariable{Key: "KEY", Value: "=value"}},
		{true, "KEY=", &object.EnvVariable{Key: "KEY", Value: ""}},
		{false, "KEY", nil},
		{false, "=KEYvalue", nil},
		{false, "=KEY=value", nil},
		{false, "=KEY=value=", nil},
	}

	for _, tt := range tests {
		envVar, err := tt.RawEnvVar.Unmarshal()
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

		if !reflect.DeepEqual(envVar, tt.EnvVar) {
			t.Errorf("Expected %#v, got %#v", tt.EnvVar, envVar)
			continue
		}
	}
}

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		Succeed     bool
		File        string
		OpenCompose *object.OpenCompose
	}{
		{
			true, `
version: 1
services:
- name: frontend
  containers:
  - name: frontend
    image: tomaskral/kompose-demo-frontend:test
    env:
    - KEY=value
    - KEY2=value2
    mappings:
    - port: 5000:8080:80/tcp
      type: LoadBalancer
      name: some-name
    - port: 5001:8081:81
      type: ClusterIp
      name: some-name2
volumes:
- name: data
  size: 1Gi
  mode: ReadWriteOnce
`,
			&object.OpenCompose{
				Version: 1,
				Services: []object.Service{
					{
						Name: "frontend",
						Containers: []object.Container{
							{
								Name:  "frontend",
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
								Mappings: []object.Mapping{
									{
										Port: object.Port{
											ContainerPort: 5000,
											HostPort:      8080,
											ServicePort:   80,
											Protocol:      "tcp",
										},
										Type: "LoadBalancer",
										Name: "some-name",
									},
									{
										Port: object.Port{
											ContainerPort: 5001,
											HostPort:      8081,
											ServicePort:   81,
										},
										Type: "ClusterIp",
										Name: "some-name2",
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
	}

	for _, tt := range tests {
		data := []byte(tt.File)
		openCompose, err := (&Decoder{}).Unmarshal(data)
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
			t.Errorf("Expected %#v, got %#v", tt.OpenCompose, openCompose)
			continue
		}
	}

}
