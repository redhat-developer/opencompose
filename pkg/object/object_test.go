package object

import (
	"fmt"
	"strings"
	"testing"

	"github.com/redhat-developer/opencompose/pkg/goutil"
	"k8s.io/client-go/pkg/util/intstr"

	api_v1 "k8s.io/client-go/pkg/api/v1"
)

const (
	Version = "0.1-dev" // TODO: replace with "1" once we reach that point
)

func TestService_EmptyDirVolumeExists(t *testing.T) {
	tests := []struct {
		ExpectedSuccess bool
		Search          string
		Service         *Service
	}{
		{
			true,
			"foo",
			&Service{
				EmptyDirVolumes: []EmptyDirVolume{
					{Name: "one"},
					{Name: "two"},
					{Name: "three"},
					{Name: "foo"},
				},
			},
		},
		{
			false,
			"foo",
			&Service{
				EmptyDirVolumes: []EmptyDirVolume{
					{Name: "one"},
					{Name: "two"},
					{Name: "three"},
				},
			},
		},
	}

	for _, test := range tests {
		output := test.Service.EmptyDirVolumeExists(test.Search)
		if output != test.ExpectedSuccess {
			t.Errorf("Expected output: %v but got: %v, for searching %q in emptyDirVolume: %+v", test.ExpectedSuccess, output, test.Search, test.Service.EmptyDirVolumes)
		}
	}
}

func TestOpenCompose_VolumeExists(t *testing.T) {
	tests := []struct {
		ExpectedSuccess bool
		Search          string
		OpenCompose     *OpenCompose
	}{
		{
			true,
			"foo",
			&OpenCompose{
				Volumes: []Volume{
					{Name: "one"},
					{Name: "two"},
					{Name: "three"},
					{Name: "foo"},
				},
			},
		},
		{
			false,
			"foo",
			&OpenCompose{
				Volumes: []Volume{
					{Name: "one"},
					{Name: "two"},
					{Name: "three"},
				},
			},
		},
	}

	for _, test := range tests {
		output := test.OpenCompose.VolumeExists(test.Search)
		if output != test.ExpectedSuccess {
			t.Errorf("Expected output: %v but got: %v, for searching %q in volumes: %+v", test.ExpectedSuccess, output, test.Search, test.OpenCompose.Volumes)
		}
	}
}

func TestValidatePortNumOrName(t *testing.T) {
	passTests := []intstr.IntOrString{
		intstr.FromInt(1),
		intstr.FromInt(1000),
		intstr.FromInt(65535),

		intstr.FromString("telnet"),
		intstr.FromString("re-mail-ck"),
		intstr.FromString("pop3"),
		intstr.FromString("a"),
		intstr.FromString("a-1"),
		intstr.FromString("1-a"),
		intstr.FromString("a-1-b-2-c"),
		intstr.FromString("1-a-2-b-3"),
	}

	for _, test := range passTests {
		err := validatePortNumOrName(test)
		if err != nil {
			t.Errorf("expected to pass, but failed for port - %v, got: %v", test.String(), err)
			continue
		}
		t.Logf("expected to pass and passed for port - %v", test.String())
	}

	failTests := []intstr.IntOrString{
		intstr.FromInt(-1),
		intstr.FromInt(0),
		intstr.FromInt(65536),
		intstr.FromInt(100000),

		intstr.FromString("longerthan15characters"),
		intstr.FromString(""),
		intstr.FromString(strings.Repeat("a", 16)),
		intstr.FromString("12345"),
		intstr.FromString("1-2-3-4"),
		intstr.FromString("-begin"),
		intstr.FromString("end-"),
		intstr.FromString("two--hyphens"),
		intstr.FromString("whois++"),
	}

	for _, test := range failTests {
		err := validatePortNumOrName(test)
		if err == nil {
			t.Errorf("expected to fail, but passed for port - %v", test.String())
			continue
		}
		t.Logf("expected to fail and failed for port - %v, got: %v", test.String(), err)
	}
}

func TestValidateExec(t *testing.T) {
	tests := []struct {
		Name            string
		ExpectedSuccess bool
		Exec            *api_v1.ExecAction
	}{
		{
			"Command given",
			true,
			&api_v1.ExecAction{
				Command: []string{"cat", "/tmp/healthz"},
			},
		},
		{
			"Command empty",
			false,
			&api_v1.ExecAction{
				Command: []string{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := validateExec(test.Exec)
			if err != nil && test.ExpectedSuccess {
				// failing condition
				t.Fatalf("Expected success but failed with error: %v", err)
			} else if err == nil && !test.ExpectedSuccess {
				// failing condition
				t.Fatal("Expected to fail but passed.")
			} else if !test.ExpectedSuccess && err != nil {
				// passing condition
				t.Logf("Failed with error: %v", err)
			}
		})
	}
}

func TestValidateHTTPGet(t *testing.T) {
	tests := []struct {
		Name            string
		ExpectedSuccess bool
		HTTPGet         *api_v1.HTTPGetAction
	}{
		{
			"Invalid scheme",
			false,
			&api_v1.HTTPGetAction{
				Scheme: "http",
				Port:   intstr.FromInt(80),
			},
		},
		{
			"Valid scheme",
			true,
			&api_v1.HTTPGetAction{
				Scheme: "HTTP",
				Port:   intstr.FromInt(80),
			},
		},
		{
			"Valid HTTPGet object",
			true,
			&api_v1.HTTPGetAction{
				Path:   "/healthz",
				Port:   intstr.FromInt(8080),
				Scheme: "HTTPS",
			},
		},
		{
			"Invalid HTTP Header",
			false,
			&api_v1.HTTPGetAction{
				HTTPHeaders: []api_v1.HTTPHeader{
					{
						Name:  "X-Forwarded-For:",
						Value: "X-Forwarded-For:",
					},
				},
				Port: intstr.FromInt(8080),
			},
		},
		{
			"Valid HTTP Header",
			true,
			&api_v1.HTTPGetAction{
				HTTPHeaders: []api_v1.HTTPHeader{
					{
						Name:  "X-Forwarded-For",
						Value: "X-Forwarded-For",
					},
				},
				Port: intstr.FromInt(8080),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := validateHTTPGet(test.HTTPGet)
			if err != nil && test.ExpectedSuccess {
				// failing condition
				t.Fatalf("Expected success but failed with error: %v", err)
			} else if err == nil && !test.ExpectedSuccess {
				// failing condition
				t.Fatal("Expected to fail but passed.")
			} else if !test.ExpectedSuccess && err != nil {
				// passing condition
				t.Logf("Failed with error: %v", err)
			}
		})
	}
}

func TestValidateTCPSocket(t *testing.T) {
	tests := []struct {
		ExpectedSuccess bool
		TCPSocket       *api_v1.TCPSocketAction
	}{
		{
			false,
			&api_v1.TCPSocketAction{Port: intstr.FromInt(0)},
		},
		{
			true,
			&api_v1.TCPSocketAction{Port: intstr.FromInt(1)},
		},
		{
			true,
			&api_v1.TCPSocketAction{Port: intstr.FromInt(65535)},
		},
		{
			false,
			&api_v1.TCPSocketAction{Port: intstr.FromInt(65536)},
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			err := validateTCPSocket(test.TCPSocket)
			if err != nil && test.ExpectedSuccess {
				// failing condition
				t.Fatalf("Expected success but failed with error: %v", err)
			} else if err == nil && !test.ExpectedSuccess {
				// failing condition
				t.Fatal("Expected to fail but passed.")
			} else if !test.ExpectedSuccess && err != nil {
				// passing condition
				t.Logf("Failed with error: %v", err)
			}
		})
	}
}

func TestPositiveNumber(t *testing.T) {
	tests := []struct {
		ExpectedSuccess bool
		Number          int32
	}{
		{true, 0},
		{true, 1},
		{false, -1},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			err := positiveNumber(test.Number)
			if err != nil && test.ExpectedSuccess {
				// failing condition
				t.Fatalf("Expected success but failed with error: %v", err)
			} else if err == nil && !test.ExpectedSuccess {
				// failing condition
				t.Fatal("Expected to fail but passed.")
			} else if !test.ExpectedSuccess && err != nil {
				// passing condition
				t.Logf("Failed with error: %v", err)
			}
		})
	}
}

func TestValidateProbes(t *testing.T) {
	tests := []struct {
		Name            string
		ExpectedSuccess bool
		Probe           *api_v1.Probe
	}{
		{
			"No probe given",
			true,
			nil,
		},
		{
			"Multiple handlers given, 'exec' and 'tcpSocket'",
			false,
			&api_v1.Probe{
				Handler: api_v1.Handler{
					Exec: &api_v1.ExecAction{
						Command: []string{"cat", "/tmp/healthz"},
					},
					TCPSocket: &api_v1.TCPSocketAction{
						Port: intstr.FromInt(80),
					},
				},
			},
		},
		{
			"Multiple handlers given, 'exec' and 'httpGet'",
			false,
			&api_v1.Probe{
				Handler: api_v1.Handler{
					Exec: &api_v1.ExecAction{
						Command: []string{"cat", "/tmp/healthz"},
					},
					HTTPGet: &api_v1.HTTPGetAction{
						Port: intstr.FromInt(80),
					},
				},
			},
		},
		{
			"Multiple handlers given, 'tcpSocket' and 'httpGet'",
			false,
			&api_v1.Probe{
				Handler: api_v1.Handler{
					HTTPGet: &api_v1.HTTPGetAction{
						Port: intstr.FromInt(80),
					},
					TCPSocket: &api_v1.TCPSocketAction{
						Port: intstr.FromInt(80),
					},
				},
			},
		},
		{
			"All handlers given",
			false,
			&api_v1.Probe{
				Handler: api_v1.Handler{
					HTTPGet: &api_v1.HTTPGetAction{
						Port: intstr.FromInt(80),
					},
					TCPSocket: &api_v1.TCPSocketAction{
						Port: intstr.FromInt(80),
					},
					Exec: &api_v1.ExecAction{
						Command: []string{"cat", "/tmp/healthz"},
					},
				},
			},
		},
		{
			"Normal Probe given",
			true,
			&api_v1.Probe{
				Handler: api_v1.Handler{
					HTTPGet: &api_v1.HTTPGetAction{
						Port: intstr.FromInt(8080),
						Path: "/healthz",
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
		{
			"Negative value for the positive fields",
			false,
			&api_v1.Probe{
				Handler: api_v1.Handler{
					HTTPGet: &api_v1.HTTPGetAction{
						Port: intstr.FromInt(8080),
					},
				},
				InitialDelaySeconds: -1,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := validateProbes(test.Probe)
			if err != nil && test.ExpectedSuccess {
				// failing condition
				t.Fatalf("Expected success but failed with error: %v", err)
			} else if err == nil && !test.ExpectedSuccess {
				// failing condition
				t.Fatal("Expected to fail but passed.")
			} else if !test.ExpectedSuccess && err != nil {
				// passing condition
				t.Logf("Failed with error: %v", err)
			}
		})
	}
}

func TestHealth_Validate(t *testing.T) {
	tests := []struct {
		Name            string
		ExpectedSuccess bool
		Health          *Health
	}{
		{
			"Valid readinessProbe",
			true,
			&Health{
				ReadinessProbe: &api_v1.Probe{
					Handler: api_v1.Handler{
						HTTPGet: &api_v1.HTTPGetAction{
							Port: intstr.FromInt(8080),
							Path: "/healthz",
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
			"Valid readinessProbe and valid livenessProbe",
			true,
			&Health{
				ReadinessProbe: &api_v1.Probe{
					Handler: api_v1.Handler{
						HTTPGet: &api_v1.HTTPGetAction{
							Port: intstr.FromInt(8080),
							Path: "/healthz",
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
				LivenessProbe: &api_v1.Probe{
					Handler: api_v1.Handler{
						Exec: &api_v1.ExecAction{
							Command: []string{"cat", "/tmp/healthz"},
						},
					},
				},
			},
		},
		{
			"Valid readinessProbe and invalid livenssProbe",
			false,
			&Health{
				ReadinessProbe: &api_v1.Probe{
					Handler: api_v1.Handler{
						HTTPGet: &api_v1.HTTPGetAction{
							Port: intstr.FromInt(8080),
							Path: "/healthz",
						},
					},
					InitialDelaySeconds: 3,
					PeriodSeconds:       3,
				},
				LivenessProbe: &api_v1.Probe{
					Handler: api_v1.Handler{
						Exec: &api_v1.ExecAction{
							Command: []string{},
						},
					},
				},
			},
		},
		{
			"Invalid readinessProbe and valid livenssProbe",
			false,
			&Health{
				ReadinessProbe: &api_v1.Probe{
					Handler: api_v1.Handler{
						HTTPGet: &api_v1.HTTPGetAction{
							Port: intstr.FromInt(80809),
							Path: "/healthz",
						},
					},
					InitialDelaySeconds: 3,
					PeriodSeconds:       3,
				},
				LivenessProbe: &api_v1.Probe{
					Handler: api_v1.Handler{
						Exec: &api_v1.ExecAction{
							Command: []string{"cat", "/tmp/healthz"},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := test.Health.validate()
			if err != nil && test.ExpectedSuccess {
				// failing condition
				t.Fatalf("Expected success but failed with error: %v", err)
			} else if err == nil && !test.ExpectedSuccess {
				// failing condition
				t.Fatal("Expected to fail but passed.")
			} else if !test.ExpectedSuccess && err != nil {
				// passing condition
				t.Logf("Failed with error: %v", err)
			}
		})
	}
}

func TestContainer_Validate(t *testing.T) {
	tests := []struct {
		Name            string
		ExpectedSuccess bool
		Container       *Container
	}{
		{
			"invalid mount name",
			false,
			&Container{
				Mounts: []Mount{
					{
						VolumeRef: "invalid_mount_name",
					},
				},
			},
		},
		{
			"mountPath not absolute",
			false,
			&Container{
				Mounts: []Mount{
					{
						VolumeRef: "test",
						MountPath: "foo/bar",
					},
				},
			},
		},
		{
			"same mountPath in multiple mounts",
			false,
			&Container{
				Mounts: []Mount{
					{
						VolumeRef: "test1",
						MountPath: "/foo/bar",
					},
					{
						VolumeRef: "test2",
						MountPath: "/foo/bar",
					},
				},
			},
		},
		{
			"passing '=' in environment variable key",
			false,
			&Container{
				Environment: []EnvVariable{
					{
						Key:   "ke=y",
						Value: "value",
					},
				},
			},
		},
		{
			"passing '=' in environment variable value",
			false,
			&Container{
				Environment: []EnvVariable{
					{
						Key:   "key",
						Value: "va=lue",
					},
				},
			},
		},
		{
			"passing a valid environment variable",
			true,
			&Container{
				Environment: []EnvVariable{
					{
						Key:   "key",
						Value: "value",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := test.Container.validate()
			if err != nil && test.ExpectedSuccess {
				// failing condition
				t.Fatalf("Expected success but failed with error: %v", err)
			} else if err == nil && !test.ExpectedSuccess {
				// failing condition
				t.Fatal("Expected to fail but passed.")
			} else if !test.ExpectedSuccess && err != nil {
				// passing condition
				t.Logf("Failed with error: %v", err)
			}
		})
	}
}

func TestService_Validate(t *testing.T) {
	tests := []struct {
		Name            string
		ExpectedSuccess bool
		Service         *Service
	}{
		{
			"invalid name of service",
			false,
			&Service{
				Name: "foo_bar",
			},
		},
		{
			"negative replica count",
			false,
			&Service{
				Name:     "test",
				Replicas: goutil.Int32Addr(-2),
			},
		},
		{
			"invalid emptyDirVolume name",
			false,
			&Service{
				Name: "test",
				EmptyDirVolumes: []EmptyDirVolume{
					{Name: "foo_bar"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := test.Service.validate()
			if err != nil && test.ExpectedSuccess {
				// failing condition
				t.Fatalf("Expected success but failed with error: %v", err)
			} else if err == nil && !test.ExpectedSuccess {
				// failing condition
				t.Fatal("Expected to fail but passed.")
			} else if !test.ExpectedSuccess && err != nil {
				// passing condition
				t.Logf("Failed with error: %v", err)
			}
		})
	}
}

func TestValidateVolumeMode(t *testing.T) {
	tests := []struct {
		ExpectedSuccess bool
		Input           string
	}{
		{true, "ReadWriteOnce"},
		{true, "ReadOnlyMany"},
		{true, "ReadWriteMany"},
		{false, "foo"},
	}

	for _, test := range tests {
		err := validateVolumeMode(test.Input)
		if err != nil && test.ExpectedSuccess {
			t.Errorf("Expected success but failed with error: %v", err)
		} else if err == nil && !test.ExpectedSuccess {
			t.Error("Expected to fail but passed.")
		}
	}
}

func TestVolume_Validate(t *testing.T) {

	storageClass := "fast"
	invalidStorageClass := "foo_foo"

	tests := []struct {
		Name            string
		ExpectedSuccess bool
		Volume          *Volume
	}{
		{
			"All fields given as valid input",
			true,
			&Volume{
				Name:         "testvol",
				Size:         "5Gi",
				AccessMode:   "ReadWriteMany",
				StorageClass: &storageClass,
			},
		},

		{
			"Invalid volume name given, should fail",
			false,
			&Volume{
				Name: "test_vol",
			},
		},

		{
			"Invalid volume size, should fail",
			false,
			&Volume{
				Name: "testvol",
				Size: "5foo",
			},
		},

		{
			"Invalid volume access mode, should fail",
			false,
			&Volume{
				Name:       "testsvol",
				Size:       "5Gi",
				AccessMode: "foo",
			},
		},

		{
			"Invalid volume storage class, should fail",
			false,
			&Volume{
				Name:         "testsvol",
				Size:         "5Gi",
				AccessMode:   "ReadWriteMany",
				StorageClass: &invalidStorageClass,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := test.Volume.validate()
			if test.ExpectedSuccess && err != nil {
				// failing condition
				t.Fatalf("Expected success but failed as: %v", err)
			} else if !test.ExpectedSuccess && err == nil {
				// failing condition
				t.Fatal("Expected failure but passed.")
			} else if !test.ExpectedSuccess && err != nil {
				// passing condition
				t.Logf("Failed with error: %v", err)
			}
		})
	}
}

func TestMount_Validate(t *testing.T) {
	tests := []struct {
		Name            string
		ExpectedSuccess bool
		Mount           *Mount
	}{
		{
			"invalid name",
			false,
			&Mount{
				VolumeRef: "invalid_mount_name",
				MountPath: "/foo/bar",
			},
		},
		{
			"valid name",
			true,
			&Mount{
				VolumeRef: "validmountname",
				MountPath: "/foo/bar",
			},
		},
		{
			"not an absolute path",
			false,
			&Mount{
				VolumeRef: "test",
				MountPath: "foo/bar",
			},
		},
		{
			"absolute path",
			true,
			&Mount{
				VolumeRef: "test",
				MountPath: "/foo/bar",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := test.Mount.validate()
			if test.ExpectedSuccess && err != nil {
				// failing condition
				t.Fatalf("Expected success but failed as: %v", err)
			} else if !test.ExpectedSuccess && err == nil {
				// failing condition
				t.Fatal("Expected failure but passed.")
			} else if !test.ExpectedSuccess && err != nil {
				// passing condition
				t.Logf("Failed with error: %v", err)
			}
		})
	}

}

func TestOpenCompose_Validate(t *testing.T) {
	name := "test-service"
	image := "test-image"

	tests := []struct {
		Name            string
		ExpectedSuccess bool
		openCompose     *OpenCompose
	}{
		{
			"Empty replica value",
			true,
			&OpenCompose{
				Version: Version,
				Services: []Service{
					{
						Name: name,
						Containers: []Container{
							{
								Image: image,
							},
						},
						Replicas: nil,
					},
				},
			},
		},

		{
			"Valid replica value: 0",
			true,
			&OpenCompose{
				Version: Version,
				Services: []Service{
					{
						Name: name,
						Containers: []Container{
							{
								Image: image,
							},
						},
						Replicas: goutil.Int32Addr(0),
					},
				},
			},
		},

		{
			"Valid replica value: 2",
			true,
			&OpenCompose{
				Version: Version,
				Services: []Service{
					{
						Name: name,
						Containers: []Container{
							{
								Image: image,
							},
						},
						Replicas: goutil.Int32Addr(2),
					},
				},
			},
		},

		{
			"Invalid replica value: -1",
			false,
			&OpenCompose{
				Version: Version,
				Services: []Service{
					{
						Name: name,
						Containers: []Container{
							{
								Image: image,
							},
						},
						Replicas: goutil.Int32Addr(-1),
					},
				},
			},
		},

		{
			"Valid labels",
			true,
			&OpenCompose{
				Version: Version,
				Services: []Service{
					{
						Name: name,
						Containers: []Container{
							{
								Image: image,
							},
						},
						Labels: Labels{
							"key1": "value1",
							"key2": "value2",
						},
					},
				},
			},
		},

		{
			"Invalid label values",
			false,
			&OpenCompose{
				Version: Version,
				Services: []Service{
					{
						Name: name,
						Containers: []Container{
							{
								Image: image,
							},
						},
						Labels: Labels{
							"key1": "garbage^value",
							"key2": "value2",
						},
					},
				},
			},
		},

		{
			"Valid Health values",
			true,
			&OpenCompose{
				Version: Version,
				Services: []Service{
					{
						Name: name,
						Containers: []Container{
							{
								Image: image,
								Health: Health{
									ReadinessProbe: &api_v1.Probe{
										Handler: api_v1.Handler{
											HTTPGet: &api_v1.HTTPGetAction{
												Port: intstr.FromInt(8080),
												Path: "/healthz",
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
									LivenessProbe: &api_v1.Probe{
										Handler: api_v1.Handler{
											Exec: &api_v1.ExecAction{
												Command: []string{"cat", "/tmp/healthz"},
											},
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
			"Invalid Health value",
			false,
			&OpenCompose{
				Version: Version,
				Services: []Service{
					{
						Name: name,
						Containers: []Container{
							{
								Image: image,
								Health: Health{
									ReadinessProbe: &api_v1.Probe{
										Handler: api_v1.Handler{
											HTTPGet: &api_v1.HTTPGetAction{
												Port: intstr.FromInt(808000),
												Path: "/healthz",
											},
										},
										InitialDelaySeconds: 3,
										PeriodSeconds:       3,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Running test: %q", test.Name), func(t *testing.T) {
			err := test.openCompose.Validate()

			if test.ExpectedSuccess && err != nil {
				t.Errorf("Expected success but failed as: %v", err)
			} else if !test.ExpectedSuccess && err == nil {
				t.Error("Expected failure but passed.")
			} else if !test.ExpectedSuccess && err != nil {
				t.Logf("Failed with error: %v", err)
			}
		})
	}
}
