package v1

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"

	"github.com/redhat-developer/opencompose/pkg/goutil"
	"github.com/redhat-developer/opencompose/pkg/object"
	"gopkg.in/yaml.v2"
)

func UriAddrFromString(s string) *Fqdn {
	return (*Fqdn)(&s)
}

func UriPathAddrFromString(s string) *PathRegex {
	return (*PathRegex)(&s)
}

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

func TestPortType_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		Succeed     bool
		RawPortType string
		PortType    object.PortType
	}{
		{true, "", object.PortType_Internal}, // UnmarshalYAML won't be even called for empty strings -> default value
		{true, "internal", object.PortType_Internal},
		{true, "external", object.PortType_External},
		{false, "'internal '", 0},
		{false, "' internal'", 0},
		{false, "' internal '", 0},
		{false, "'external '", 0},
		{false, "' external'", 0},
		{false, "' external '", 0},
		{false, "'something '", 0},
		{false, "' something'", 0},
		{false, "' something '", 0},
	}

	for _, tt := range tests {
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
		Succeed bool
		RawPort string
		Port    Port
	}{
		{true, "", Port{}}, // UnmarshalYAML won't be even called for empty strings -> default value
		{
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
			false, //you have to specify host
			`
port: 5000:80
path: "/admin"
`,
			Port{},
		},
		{
			false, //you have to specify host
			`
port: 5000:80
path: ""
`,
			Port{},
		},
		{
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
		t.Run("", func(t *testing.T) {
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
volumeName: test-volume
mountPath: /foo/bar
volumeSubPath: some/path
readOnly: true
`,
			&Mount{
				VolumeName:    "test-volume",
				MountPath:     "/foo/bar",
				VolumeSubPath: goutil.StringAddr("some/path"),
				ReadOnly:      goutil.BoolAddr(true),
			},
		},

		{
			"Optional fields not given",
			true, `
volumeName: test-volume
mountPath: /foo/bar
`,
			&Mount{
				VolumeName: "test-volume",
				MountPath:  "/foo/bar",
			},
		},

		{
			"Giving bool value as 'foobar', should fail",
			false, `
volumeName: test-volume
mountPath: /foo/bar
readOnly: foobar
`,
			nil,
		},

		{
			"Giving an extra field which does not exist",
			false, `
volumeName: test-volume
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
volumeName: test-volume
readOnly: true
`,
			&Mount{
				VolumeName: "test-volume",
				ReadOnly:   goutil.BoolAddr(true),
			},
		},
	}

	for _, test := range tests {
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

func TestEmptyDirVolume_UnmarshalYAML(t *testing.T) {

	tests := []struct {
		Succeed     bool
		RawEmptyDir string
		EmptyDir    *EmptyDirVolume
	}{
		{true, "name: empty", &EmptyDirVolume{Name: "empty"}},
		{
			false, `
name: empty
excess: field
`,
			nil,
		},
		{true, "", &EmptyDirVolume{}}, // UnmarshalYAML won't be even called for empty strings -> default value
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
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

func TestService_UnmarshalYAML(t *testing.T) {
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
`,
			&Service{
				Name:     "frontend",
				Replicas: goutil.Int32Addr(3),
				Containers: []Container{
					{
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
`,
			nil,
		},

		{
			"Providing only mandatory fields, omitting the optional ones",
			true, `
name: frontend
containers:
- image: tomaskral/kompose-demo-frontend:test
`,
			&Service{
				Name: "frontend",
				Containers: []Container{
					{
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
  mounts:
  - volumeName: test-volume
    mountPath: /foo/bar
    volumeSubPath: some/path
    readOnly: true
`,
			&Service{
				Name: "frontend",
				Containers: []Container{
					{
						Image: "tomaskral/kompose-demo-frontend:test",
						Mounts: []Mount{
							{
								VolumeName:    "test-volume",
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
emptyDirVolumes:
- name: empty
`,
			&Service{
				Name: "frontend",
				Containers: []Container{
					{
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
    mounts:
    - volumeName: test-volume
      mountPath: /foo/bar
      volumeSubPath: some/path
      readOnly: true
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
										VolumeName:    "test-volume",
										MountPath:     "/foo/bar",
										VolumeSubPath: "some/path",
										ReadOnly:      true,
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
		{
			true, `
version: 0.1-dev
services:
- name: helloworld
  replicas: 2
  containers:
  - image: tomaskral/nonroot-nginx
`,
			&object.OpenCompose{
				Version: Version,
				Services: []object.Service{
					{
						Name:     "helloworld",
						Replicas: goutil.Int32Addr(2),
						Containers: []object.Container{
							{
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
    mounts:
    - volumeName: test-volume
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
volumes:
- name: db
  size: 5Gi
  storageClass: fast
`,
			nil,
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
