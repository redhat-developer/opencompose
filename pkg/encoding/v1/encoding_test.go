package v1

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"

	"github.com/redhat-developer/opencompose/pkg/goutil"
	"github.com/redhat-developer/opencompose/pkg/object"
	"gopkg.in/yaml.v2"

	"k8s.io/client-go/pkg/util/intstr"

	api_v1 "k8s.io/client-go/pkg/api/v1"
)

func UriAddrFromString(s string) *Fqdn {
	return (*Fqdn)(&s)
}

func UriPathAddrFromString(s string) *PathRegex {
	return (*PathRegex)(&s)
}

func TestPortMapping_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		Name    string
		Succeed bool
		RawPort string
		Port    *PortMapping
	}{
		{"Only ContainerPort", true, "5000", &PortMapping{ContainerPort: 5000, ServicePort: 5000}},
		{"ContainerPort:ServicePort", true, "5000:80", &PortMapping{ContainerPort: 5000, ServicePort: 80}},
		{"Empty Portmapping", true, "", &PortMapping{}}, // UnmarshalYAML won't be even called for empty strings
		{"Failed Portmapping x5000", false, "x5000", nil},
		{"Failed Portmapping 5000", false, "5000:", nil},
		{"Failed Portmapping x5000:", false, "x5000:", nil},
		{"Failed Portmapping :80", false, ":80", nil},
		{"Failed Portmapping :80x", false, ":80x", nil},
		{"Failed Portmapping x:80x", false, "x:80x", nil},
		{"Failed Portmapping x:80", false, "x:80", nil},
		{"Failed Portmapping :8080:80", false, ":8080:80", nil},
		{"Failed Portmapping x:8080:80", false, "x:8080:80", nil},
		{"Failed Portmapping x:x8080:x80", false, "x:x8080:x80", nil},
		{"Failed Portmapping :8080:x80", false, ":8080:x80", nil},
		{"Failed Portmapping :8080:80:", false, ":8080:80:", nil},
		{"Failed Portmapping 8080:80:", false, "8080:80:", nil},
		{"Failed Portmapping 8080:80::", false, "8080:80::", nil},
		{"Failed Portmapping :", false, ":", nil},
		{"Failed Portmapping ::", false, "::", nil},
		{"Failed Portmapping :::", false, ":::", nil},
		{"Failed Portmapping ::::", false, "::::", nil},
		{"Failed Portmapping :::::", false, ":::::", nil},
		{"Failed Portmapping ::1:80", false, "::1:80", nil},
		{"Failed Portmapping ::1:8080:80", false, "::1:8080:80", nil},
		{"Failed Portmapping ::1:8080:80:", false, "::1:8080:80:", nil},
	}

	for _, tt := range tests {
		t.Log("Test case: ", tt.Name)
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

func TestPortType_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		Name        string
		Succeed     bool
		RawPortType string
		PortType    object.PortType
	}{
		{"Success Empty Porttype", true, "", object.PortType_Internal}, // UnmarshalYAML won't be even called for empty strings -> default value
		{"Success Internal Porttype", true, "internal", object.PortType_Internal},
		{"Success External Porttype", true, "external", object.PortType_External},
		{"Failed Porttype - 'internal '", false, "'internal '", 0},
		{"Failed Porttype - ' internal'", false, "' internal'", 0},
		{"Failed Porttype - ' internal '", false, "' internal '", 0},
		{"Failed Porttype - 'external '", false, "'external '", 0},
		{"Failed Porttype - ' external'", false, "' external'", 0},
		{"Failed Porttype - ' external '", false, "' external '", 0},
		{"Failed Porttype - 'something '", false, "'something '", 0},
		{"Failed Porttype - ' something'", false, "' something'", 0},
		{"Failed Porttype - ' something '", false, "' something '", 0},
	}

	for _, tt := range tests {
		t.Log("Test case: ", tt.Name)
		var pt PortType
		err := yaml.Unmarshal([]byte(tt.RawPortType), &pt)
		if err != nil {
			if tt.Succeed {
				t.Errorf("Failed to unmarshal port type %q: %s", tt.RawPortType, err)
			}
			continue
		}

		if !tt.Succeed {
			t.Errorf("Expected port type %#v to fail!", tt.RawPortType)
			continue
		}

		if object.PortType(pt) != tt.PortType {
			t.Errorf("Expected port type %#v, got %#v", tt.PortType, pt)
			continue
		}
	}
}

func TestPort_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		Name    string
		Succeed bool
		RawPort string
		Port    Port
	}{
		{"Empty Port mapping", true, "", Port{}}, // UnmarshalYAML won't be even called for empty strings -> default value
		{
			"Port mapping with empty host",
			true,
			`
port: 5000:80
host: ""
`,
			Port{
				Port: PortMapping{ContainerPort: 5000, ServicePort: 80},
				Host: UriAddrFromString(""),
				Path: UriPathAddrFromString(""), // path defaults to "" when host is validated/set
			},
		},
		{
			"Port mapping with hostname",
			true,
			`
port: 5000:80
host: "subdomain.127.0.0.1.nip.io"
`,
			Port{
				Port: PortMapping{ContainerPort: 5000, ServicePort: 80},
				Host: UriAddrFromString("subdomain.127.0.0.1.nip.io"),
				Path: UriPathAddrFromString(""), // path defaults to "" when host is validated/set
			},
		},
		{
			"Port mapping with path and empty host",
			false, //you have to specify host
			`
port: 5000:80
path: "/admin"
`,
			Port{},
		},
		{
			"Port mapping with empty path",
			false, //you have to specify host
			`
port: 5000:80
path: ""
`,
			Port{},
		},
		{
			"Port mapping with empty host as well path",
			true,
			`
port: 5000:80
host: ""
path: ""
`,
			Port{
				Port: PortMapping{ContainerPort: 5000, ServicePort: 80},
				Host: UriAddrFromString(""),
				Path: UriPathAddrFromString(""),
			},
		},
		{
			"Port mapping with path and empty host",
			true,
			`
port: 5000:80
host: ""
path: "/admin"
`,
			Port{
				Port: PortMapping{ContainerPort: 5000, ServicePort: 80},
				Host: UriAddrFromString(""),
				Path: UriPathAddrFromString("/admin"),
			},
		},
		{
			"Port mapping with host and Path",
			true,
			`
port: 5000:80
host: "subdomain.127.0.0.1.nip.io"
path: "/admin"
`,
			Port{
				Port: PortMapping{ContainerPort: 5000, ServicePort: 80},
				Host: UriAddrFromString("subdomain.127.0.0.1.nip.io"),
				Path: UriPathAddrFromString("/admin"),
			},
		},
	}

	for _, tt := range tests {
		t.Log("Test case: ", tt.Name)
		t.Run(tt.Name, func(t *testing.T) {
			var p Port
			err := yaml.Unmarshal([]byte(tt.RawPort), &p)
			if err != nil {
				if tt.Succeed {
					t.Fatalf("Failed to unmarshal port %q: %s", tt.RawPort, err)
				}
				return
			}

			if !tt.Succeed {
				t.Fatal(spew.Errorf("Expected port %#+v to fail!", tt.RawPort))
			}

			if !reflect.DeepEqual(p, tt.Port) {
				t.Fatal(spew.Errorf("Expected:\n%#+v\n, got:\n%#+v", tt.Port, p))
			}
		})
	}
}

func TestEnvVariable_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		Name      string
		Succeed   bool
		RawEnvVar string
		EnvVar    *EnvVariable
	}{
		{
			"Success Value is given",
			true,
			`
name: "KEY"
value: "value string "
`,
			&EnvVariable{Key: "KEY", Value: "value string "},
		},
		{
			"Success value is given with space",
			true,
			`
name: "KEY"
value: " value"
`,
			&EnvVariable{Key: "KEY", Value: " value"},
		},
		{
			"Success key is given with space",
			true,
			`
name: " KEY"
value: "value"
`,
			&EnvVariable{Key: " KEY", Value: "value"},
		},

		{
			"Failed value is not given",
			true,
			`
name: KEY
value: ""
`,
			&EnvVariable{Key: "KEY", Value: ""},
		},
		{
			"Failed value is not string",
			false,
			`
name: KEY
key: extra_field
`,
			nil,
		},
	}

	for _, tt := range tests {
		t.Log("Test case: ", tt.Name)
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

func TestMount_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		Name     string
		Succeed  bool
		RawMount string
		Mount    *Mount
	}{
		{
			"All fields given",
			true, `
volumeRef: test-volume
mountPath: /foo/bar
volumeSubPath: some/path
readOnly: true
`,
			&Mount{
				VolumeRef:     "test-volume",
				MountPath:     "/foo/bar",
				VolumeSubPath: goutil.StringAddr("some/path"),
				ReadOnly:      goutil.BoolAddr(true),
			},
		},

		{
			"Optional fields not given",
			true, `
volumeRef: test-volume
mountPath: /foo/bar
`,
			&Mount{
				VolumeRef: "test-volume",
				MountPath: "/foo/bar",
			},
		},

		{
			"Giving bool value as 'foobar', should fail",
			false, `
volumeRef: test-volume
mountPath: /foo/bar
readOnly: foobar
`,
			nil,
		},

		{
			"Giving an extra field which does not exist",
			false, `
volumeRef: test-volume
mountPath: /foo/bar
foo: bar
`,
			nil,
		},

		{
			"No fields given", // UnmarshalYAML won't be even called for empty strings -> default value
			true,
			"",
			&Mount{},
		},

		{
			"Not giving a required field",
			true, `
volumeRef: test-volume
readOnly: true
`,
			&Mount{
				VolumeRef: "test-volume",
				ReadOnly:  goutil.BoolAddr(true),
			}}}
	for _, test := range tests {
		t.Log("Test case: ", test.Name)
		t.Run(test.Name, func(t *testing.T) {
			var mount Mount
			err := yaml.Unmarshal([]byte(test.RawMount), &mount)
			if err != nil {
				if test.Succeed {
					t.Errorf("failed to unmarshal 'Mount': %#v\nerror: %v", test.RawMount, err)
				}
				return
			}

			if !test.Succeed {
				t.Fatalf("Expected %#v to fail, but succeeded! Mount object looks like: %#v", test.RawMount, mount)
			}

			if !reflect.DeepEqual(mount, *test.Mount) {
				t.Fatalf("Expected %#v\ngot %#v", *test.Mount, mount)
			}
		})
	}
}

func TestInterfaceToProbe(t *testing.T) {
	tests := []struct {
		Name     string
		Succeed  bool
		RawProbe string
		Probe    *api_v1.Probe
	}{
		{
			"Probe with exec",
			true, `
exec:
  command:
  - cat
  - /tmp/healthy
initialDelaySeconds: 5
periodSeconds: 5`,
			&api_v1.Probe{
				Handler: api_v1.Handler{
					Exec: &api_v1.ExecAction{
						Command: []string{"cat", "/tmp/healthy"},
					},
				},
				InitialDelaySeconds: 5,
				PeriodSeconds:       5,
			},
		},
		{
			"Probe with httpGet",
			true, `
httpGet:
  scheme: HTTP`,
			&api_v1.Probe{
				Handler: api_v1.Handler{
					HTTPGet: &api_v1.HTTPGetAction{
						Scheme: api_v1.URISchemeHTTP,
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var input interface{}
			err := yaml.Unmarshal([]byte(test.RawProbe), &input)
			if err != nil {
				t.Fatalf("error unmarshalling input to interface: %v", err)
			}

			gotProbe, err := interfaceToProbe(input)
			if err != nil {
				if test.Succeed {
					t.Errorf("failed to convert raw interface: %q to Probe, error: %v\n", test.RawProbe, err)
				}
				return
			}

			if !test.Succeed {
				t.Fatalf("expected %q to fail, but succeeded! Probe object looks like: %s", test.RawProbe, spew.Sprint(gotProbe))
			}

			if !reflect.DeepEqual(test.Probe, gotProbe) {
				t.Fatalf("expected %#v\ngot %#v", spew.Sprint(test.Probe), spew.Sprint(gotProbe))
			}

		})

	}
}

func TestHealth_UnmarshalYAML(t *testing.T) {

	tests := []struct {
		Name      string
		Succeed   bool
		RawHealth string
		Health    *Health
	}{
		{
			"Given all fields1",
			true, `
livenessProbe:
  exec:
    command:
    - cat
    - /tmp/healthy
  initialDelaySeconds: 5
  periodSeconds: 5
  timeoutSeconds: 1
readinessProbe:
  httpGet:
    path: /healthz
    port: 8080
    httpHeaders:
    - name: X-Custom-Header
      value: Awesome
  initialDelaySeconds: 3
  periodSeconds: 3`,
			&Health{
				LivenessProbe: &api_v1.Probe{
					Handler: api_v1.Handler{
						Exec: &api_v1.ExecAction{
							Command: []string{"cat", "/tmp/healthy"},
						},
					},
					InitialDelaySeconds: 5,
					PeriodSeconds:       5,
					TimeoutSeconds:      1,
				},
				ReadinessProbe: &api_v1.Probe{
					Handler: api_v1.Handler{
						HTTPGet: &api_v1.HTTPGetAction{
							Path: "/healthz",
							Port: intstr.FromInt(8080),
							HTTPHeaders: []api_v1.HTTPHeader{
								{
									Name:  "X-Custom-Header",
									Value: "Awesome",
								},
							},
						},
					},
					InitialDelaySeconds: 3,
					PeriodSeconds:       3,
				},
			},
		},
		{
			"Given all fields2",
			true, `
readinessProbe:
  tcpSocket:
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 10
livenessProbe:
  tcpSocket:
    port: 8080
  initialDelaySeconds: 15
  periodSeconds: 20`,
			&Health{
				ReadinessProbe: &api_v1.Probe{
					Handler: api_v1.Handler{
						TCPSocket: &api_v1.TCPSocketAction{
							Port: intstr.FromInt(8080),
						},
					},
					InitialDelaySeconds: 5,
					PeriodSeconds:       10,
				},
				LivenessProbe: &api_v1.Probe{
					Handler: api_v1.Handler{
						TCPSocket: &api_v1.TCPSocketAction{
							Port: intstr.FromInt(8080),
						},
					},
					InitialDelaySeconds: 15,
					PeriodSeconds:       20,
				},
			},
		},
		{
			"Only required fields given",
			true, `
readinessProbe:
  httpGet:
    port: 8080`,
			&Health{
				ReadinessProbe: &api_v1.Probe{
					Handler: api_v1.Handler{
						HTTPGet: &api_v1.HTTPGetAction{

							Port: intstr.FromInt(8080),
						},
					},
				},
			},
		},
		{
			"Extra field given", // FIXME: this should fail ideally, see if an extra field is given
			true, `
readinessProbe:
  httpGet:
    port: 8080
  foo: bar`,
			&Health{
				ReadinessProbe: &api_v1.Probe{
					Handler: api_v1.Handler{
						HTTPGet: &api_v1.HTTPGetAction{

							Port: intstr.FromInt(8080),
						},
					},
				},
			},
		},
		{
			"Wrong field type given",
			false, `
readinessProbe:
  httpGet: foobar`,
			nil,
		},
		{
			"No fields given", // UnmarshalYAML won't be even called for empty strings -> default value
			true,
			``,
			&Health{},
		},
		{
			"Required field not given", // FIXME: this should fail ideally, see if the required fields not given
			true, `
readinessProbe:
  initialDelaySeconds: 5`,
			&Health{
				ReadinessProbe: &api_v1.Probe{
					InitialDelaySeconds: 5,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			var health Health
			err := yaml.Unmarshal([]byte(test.RawHealth), &health)
			if err != nil {
				if test.Succeed {
					t.Errorf("failed to unmarshal 'Health': %q\nError: %v", test.RawHealth, err)
				}
				return
			}

			if !test.Succeed {
				t.Fatalf("Expected %q to fail, but succeeded! Health object looks like: %s", test.RawHealth, spew.Sprint(health))
			}

			if !reflect.DeepEqual(health, *test.Health) {
				t.Fatalf("Expected: %s\nGot: %s", spew.Sprint(test.Health), spew.Sprint(health))
			}
		})
	}
}

func TestEmptyDirVolume_UnmarshalYAML(t *testing.T) {

	tests := []struct {
		Name        string
		Succeed     bool
		RawEmptyDir string
		EmptyDir    *EmptyDirVolume
	}{
		{"name provided", true, "name: empty", &EmptyDirVolume{Name: "empty"}},
		{
			"nothing provided",
			false, `
name: empty
excess: field
`,
			nil,
		},
		{"Blank string provided", true, "", &EmptyDirVolume{}}, // UnmarshalYAML won't be even called for empty strings -> default value
	}

	for _, test := range tests {
		t.Log("Test case: ", test.Name)
		t.Run(test.Name, func(t *testing.T) {
			var emptyDir EmptyDirVolume
			err := yaml.Unmarshal([]byte(test.RawEmptyDir), &emptyDir)
			if err != nil {
				if test.Succeed {
					t.Errorf("failed to unmarshal 'EmptyDirVolume': %#v\nError: %v", test.RawEmptyDir, err)
				}
				return
			}

			if !test.Succeed {
				t.Fatalf("Expected %#v to fail, but succeeded! EmptyDirVolume object looks like: %#v", test.RawEmptyDir, emptyDir)
			}

			if !reflect.DeepEqual(emptyDir, *test.EmptyDir) {
				t.Fatalf("Expected %#v\ngot %#v", *test.EmptyDir, emptyDir)
			}
		})
	}
}

func TestLabels_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name      string
		Succeed   bool
		RawLabels string
		Labels    *Labels
	}{
		{
			"Providing valid label strings",
			true,
			`
key1: value1
key2: value2
key3:
key4: value4
`,
			&Labels{
				"key1": "value1",
				"key2": "value2",
				"key3": "",
				"key4": "value4",
			},
		},
	}
	for _, tt := range tests {
		var labels Labels
		err := yaml.Unmarshal([]byte(tt.RawLabels), &labels)
		if err != nil {
			if tt.Succeed {
				t.Errorf("Failed to unmarshal %#v; error %#v", tt.RawLabels, err)
			}
			continue
		}

		if !tt.Succeed {
			t.Errorf("Expected %#v to fail!", tt.RawLabels)
			continue
		}

		if !reflect.DeepEqual(labels, *tt.Labels) {
			t.Errorf("Expected %#v, got %#v", *tt.Labels, labels)
			continue
		}
	}
}

func TestService_UnmarshalYAML(t *testing.T) {
	containerName := ResourceName("test")
	tests := []struct {
		Name       string
		Succeed    bool
		RawService string
		Service    *Service
	}{
		{
			"Replica as positive int",
			true, `
name: frontend
replicas: 3
containers:
- image: tomaskral/kompose-demo-frontend:test
  name: test
`,
			&Service{
				Name:     "frontend",
				Replicas: goutil.Int32Addr(3),
				Containers: []Container{
					{
						Name:  containerName,
						Image: "tomaskral/kompose-demo-frontend:test",
					},
				},
			},
		},

		{
			"Replica as 'string'",
			false, `
name: frontend
replicas: notint
containers:
- image: tomaskral/kompose-demo-frontend:test
  name: test
`,
			nil,
		},

		{
			"Providing only mandatory fields, omitting the optional ones",
			true, `
name: frontend
containers:
- image: tomaskral/kompose-demo-frontend:test
  name: test
`,
			&Service{
				Name: "frontend",
				Containers: []Container{
					{
						Name:  containerName,
						Image: "tomaskral/kompose-demo-frontend:test",
					},
				},
			},
		},

		{
			"Checking mounts works when integrated with services",
			true, `
name: frontend
containers:
- image: tomaskral/kompose-demo-frontend:test
  name: test
  mounts:
  - volumeRef: test-volume
    mountPath: /foo/bar
    volumeSubPath: some/path
    readOnly: true
`,
			&Service{
				Name: "frontend",
				Containers: []Container{
					{
						Name:  containerName,
						Image: "tomaskral/kompose-demo-frontend:test",
						Mounts: []Mount{
							{
								VolumeRef:     "test-volume",
								MountPath:     "/foo/bar",
								VolumeSubPath: goutil.StringAddr("some/path"),
								ReadOnly:      goutil.BoolAddr(true),
							},
						},
					},
				},
			},
		},

		{
			"Integrate emptyDirVolume with service",
			true, `
name: frontend
containers:
- image: tomaskral/kompose-demo-frontend:test
  name: test
emptyDirVolumes:
- name: empty
`,
			&Service{
				Name: "frontend",
				Containers: []Container{
					{
						Name:  containerName,
						Image: "tomaskral/kompose-demo-frontend:test",
					},
				},
				EmptyDirVolumes: []EmptyDirVolume{
					{
						Name: "empty",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Log("Test case: ", test.Name)
		t.Run(test.Name, func(t *testing.T) {
			var service Service
			err := yaml.Unmarshal([]byte(test.RawService), &service)
			if err != nil {
				if test.Succeed {
					t.Errorf("Failed to unmarshal: %#v\nError: %#v", test.RawService, err)
				}
				return
			}

			if !test.Succeed {
				t.Fatalf("Expected %#v to fail, but succeeded!", test.RawService)
			}

			if !reflect.DeepEqual(service, *test.Service) {
				t.Fatalf("Expected %#v\ngot %#v", *test.Service, service)
			}
		})
	}
}

func TestVolume_UnmarshalYAML(t *testing.T) {

	storageClass := ResourceName("fast")

	tests := []struct {
		Name      string
		Succeed   bool
		RawVolume string
		Volume    *Volume
	}{
		{
			"All fields given",
			true, `
name: db
size: 3Gi
accessMode: ReadWriteMany
storageClass: fast
`,
			&Volume{
				Name:         ResourceName("db"),
				Size:         "3Gi",
				AccessMode:   "ReadWriteMany",
				StorageClass: &storageClass,
			},
		},

		{
			"Optional fields not given",
			true, `
name: db
size: 3Gi
accessMode: ReadWriteMany
`,
			&Volume{
				Name:       ResourceName("db"),
				Size:       "3Gi",
				AccessMode: "ReadWriteMany",
			},
		},

		{
			"Extra field given",
			false, `
name: db
size: 3Gi
accessMode: ReadWriteMany
excess: key
`,
			nil,
		},

		{
			"No fields given", // UnmarshalYAML won't be even called for empty strings -> default value
			true,
			"",
			&Volume{},
		},
	}

	for _, test := range tests {
		t.Log("Test case: ", test.Name)
		t.Run(test.Name, func(t *testing.T) {
			var volume Volume
			err := yaml.Unmarshal([]byte(test.RawVolume), &volume)
			if err != nil {
				if test.Succeed {
					t.Errorf("failed to unmarshal: %#v\nerror: %#v", test.RawVolume, err)
				}
				return
			}

			if !test.Succeed {
				t.Fatalf("Expected %#v to fail, but succeeded!", test.RawVolume)
			}

			if !reflect.DeepEqual(volume, *test.Volume) {
				t.Fatalf("Expected %#v\ngot %#v", *test.Volume, volume)
			}
		})
	}
}

func TestDecoder_Decode(t *testing.T) {
	// TODO: make better tests w.r.t excess keys in all possible places
	// TODO: add checking for proper error because tests can fail for other than expected reasons

	containerName := "test"
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
    name: test
    env:
    - name: KEY
      value: value
    - name: KEY2
      value: value2
    ports:
    - port: 5000:80
    - port: 5001:81
    mounts:
    - volumeRef: test-volume
      mountPath: /foo/bar
      volumeSubPath: some/path
      readOnly: true
    health:
      readinessProbe:
        httpGet:
          path: /healthz
          port: 8080
          httpHeaders:
          - name: X-Custom-Header
            value: Awesome
        initialDelaySeconds: 3
        periodSeconds: 3
  emptyDirVolumes:
  - name: empty
volumes:
- name: data
  size: 1Gi
  accessMode: ReadWriteOnce
  storageClass: fast
`,
			&object.OpenCompose{
				Version: Version,
				Services: []object.Service{
					{
						Name: "frontend",
						Containers: []object.Container{
							{
								Name:  containerName,
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
										Type: object.PortType_Internal,
									},
									{
										Port: object.PortMapping{
											ContainerPort: 5001,
											ServicePort:   81,
										},
									},
								},
								Mounts: []object.Mount{
									{
										VolumeRef:     "test-volume",
										MountPath:     "/foo/bar",
										VolumeSubPath: "some/path",
										ReadOnly:      true,
									},
								},
								Health: object.Health{
									ReadinessProbe: &api_v1.Probe{
										Handler: api_v1.Handler{
											HTTPGet: &api_v1.HTTPGetAction{
												Path: "/healthz",
												Port: intstr.FromInt(8080),
												HTTPHeaders: []api_v1.HTTPHeader{
													{
														Name:  "X-Custom-Header",
														Value: "Awesome",
													},
												},
											},
										},
										InitialDelaySeconds: 3,
										PeriodSeconds:       3,
									},
								},
							},
						},
						EmptyDirVolumes: []object.EmptyDirVolume{
							{
								Name: "empty",
							},
						},
					},
				},
				Volumes: []object.Volume{
					{
						Name:         "data",
						Size:         "1Gi",
						AccessMode:   "ReadWriteOnce",
						StorageClass: goutil.StringAddr("fast"),
					},
				},
			},
		},
		{
			true, `
version: 0.1-dev
services:
- name: helloworld
  containers:
  - image: tomaskral/nonroot-nginx
    name: test
    ports:
    - port: 8080
      type: external
`,
			&object.OpenCompose{
				Version: Version,
				Services: []object.Service{
					{
						Name: "helloworld",
						Containers: []object.Container{
							{
								Name:  containerName,
								Image: "tomaskral/nonroot-nginx",
								Ports: []object.Port{
									{
										Port: object.PortMapping{
											ContainerPort: 8080,
											ServicePort:   8080,
										},
										Type: object.PortType_External,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			true, `
version: 0.1-dev
services:
- name: helloworld
  containers:
  - image: tomaskral/nonroot-nginx
    name: test
    ports:
    - port: 8080
      type: internal
`,
			&object.OpenCompose{
				Version: Version,
				Services: []object.Service{
					{
						Name: "helloworld",
						Containers: []object.Container{
							{
								Name:  containerName,
								Image: "tomaskral/nonroot-nginx",
								Ports: []object.Port{
									{
										Port: object.PortMapping{
											ContainerPort: 8080,
											ServicePort:   8080,
										},
										Type: object.PortType_Internal,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			true, `
version: 0.1-dev
services:
- name: helloworld
  containers:
  - image: tomaskral/nonroot-nginx
    name: test
    ports:
    - port: 8080
      host: hw-nginx.127.0.0.1.nip.io
      path: /admin
`,
			&object.OpenCompose{
				Version: Version,
				Services: []object.Service{
					{
						Name: "helloworld",
						Containers: []object.Container{
							{
								Name:  containerName,
								Image: "tomaskral/nonroot-nginx",
								Ports: []object.Port{
									{
										Port: object.PortMapping{
											ContainerPort: 8080,
											ServicePort:   8080,
										},
										Type: object.PortType_Internal,
										Host: goutil.StringAddr("hw-nginx.127.0.0.1.nip.io"),
										Path: "/admin",
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
services:
- name: frontend
  containers:
  - image: tomaskral/kompose-demo-frontend:test
    name: test
    env:
	- name: KEY
	  value: value
	- name: KEY2
	  value: value2
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
    name: test
  - image: tomaskral/kompose-demo-frontend:test
	env:
	- name: KEY
	  value: value
	- name: KEY2
	  value: value2
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
services:
- name: frontend
  containers:
  - image: tomaskral/kompose-demo-frontend:test
    name: test
    env:
    - name: KEY
      value: value
    - name: KEY2
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
    name: test
    env:
    - name: KEY
      value: value
    - value: value
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
    name: test
    env:
    - name: KEY
      value: value
    - name:
      value: value
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
    name: test
    env:
    - name: KEY
      value: value
    - name: KEY2
      value:
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
    name: test
    env:
    - name: KEY
      value: value
    - name: KEY2
      value: value2
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
								Name:  containerName,
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
		{
			true, `
version: 0.1-dev
services:
- name: helloworld
  replicas: 2
  containers:
  - image: tomaskral/nonroot-nginx
    name: test
`,
			&object.OpenCompose{
				Version: Version,
				Services: []object.Service{
					{
						Name:     "helloworld",
						Replicas: goutil.Int32Addr(2),
						Containers: []object.Container{
							{
								Name:  containerName,
								Image: "tomaskral/nonroot-nginx",
							},
						},
					},
				},
			},
		},
		{ // testing mounts, one required value "mountPath" -  is not given
			false, `
version: 0.1-dev
services:
- name: helloworld
  replicas: 2
  containers:
  - image: tomaskral/nonroot-nginx
    name: test
    mounts:
    - volumeRef: test-volume
      readOnly: true
`,
			nil,
		},
		{ // testing Volumes, one required value "accessMode" - is not given
			false, `
version: 0.1-dev
services:
- name: helloworld
  containers:
  - image: tomaskral/nonroot-nginx
    name: test
volumes:
- name: db
  size: 5Gi
  storageClass: fast
`,
			nil,
		},
		{
			true,
			`
version: 0.1-dev
services:
- name: helloworld
  replicas: 2
  containers:
  - image: tomaskral/nonroot-nginx
    name: test
  labels:
      key1: value1
      key2: value2
      key3:
      key4: value4
`,
			&object.OpenCompose{
				Version: Version,
				Services: []object.Service{
					{
						Name:     "helloworld",
						Replicas: goutil.Int32Addr(2),
						Containers: []object.Container{
							{
								Name:  containerName,
								Image: "tomaskral/nonroot-nginx",
							},
						},
						Labels: object.Labels{
							"key1": "value1",
							"key2": "value2",
							"key3": "",
							"key4": "value4",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			data := []byte(tt.File)
			openCompose, err := (&Decoder{}).Decode(data)
			if err != nil {
				if tt.Succeed {
					t.Fatalf("Failed to unmarshal %#v; error %v", tt.File, err)
				}
				t.Logf("Expected to fail and failed as: %v", err)
				return
			}

			if !tt.Succeed {
				t.Fatal(spew.Errorf("Expected %#+v to fail!", tt.File))
			}

			if !reflect.DeepEqual(openCompose, tt.OpenCompose) {
				t.Fatal(spew.Errorf("Expected:\n%#+v\n, got:\n%#+v", tt.OpenCompose, openCompose))
			}
		})
	}
}
