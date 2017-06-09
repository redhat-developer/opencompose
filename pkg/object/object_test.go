package object

import (
	"fmt"
	"testing"

	"bytes"
	"log"
	"os"
	"strings"

	"github.com/redhat-developer/opencompose/pkg/goutil"
)

const (
	Version = "0.1-dev" // TODO: replace with "1" once we reach that point
)

func captureStderr(f func()) string {
	var buffer bytes.Buffer
	log.SetOutput(&buffer)
	f()
	log.SetOutput(os.Stderr)
	return buffer.String()
}

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
						VolumeRef: goutil.StringAddr("invalid_mount_name"),
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
						VolumeRef: goutil.StringAddr("test"),
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
						VolumeRef: goutil.StringAddr("test1"),
						MountPath: "/foo/bar",
					},
					{
						VolumeRef: goutil.StringAddr("test2"),
						MountPath: "/foo/bar",
					},
				},
			},
		},
		{
			"passing both mountRef and secretRef",
			false,
			&Container{
				Mounts: []Mount{
					{
						VolumeRef: goutil.StringAddr("mountDef"),
						SecretRef: &SecretDef{
							SecretName: "foo",
							DataKey:    "bar",
						},
						MountPath: "/foo/bar",
					},
				},
			},
		},
		{
			"passing same mountPath in volumeRef and secretRef",
			false,
			&Container{
				Mounts: []Mount{
					{
						VolumeRef: goutil.StringAddr("mount"),
						MountPath: "/foo/bar",
					},
					{
						SecretRef: &SecretDef{
							SecretName: "foo",
							DataKey:    "bar",
						},
						MountPath: "/foo/bar",
					},
				},
			},
		},
		{
			"passing valid mount with secretRef",
			true,
			&Container{
				Mounts: []Mount{
					{
						SecretRef: &SecretDef{
							SecretName: "foo",
							DataKey:    "bar",
						},
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
						Value: goutil.StringAddr("value"),
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
						Value: goutil.StringAddr("va=lue"),
					},
				},
			},
		},
		{
			"passing a valid environment variable with key and value",
			true,
			&Container{
				Environment: []EnvVariable{
					{
						Key:   "key",
						Value: goutil.StringAddr("value"),
					},
				},
			},
		},
		{
			"passing both, value and secretRef, in environment variable",
			false,
			&Container{
				Environment: []EnvVariable{
					{
						Key:   "key",
						Value: goutil.StringAddr("value"),
						SecretRef: &SecretDef{
							SecretName: "foo",
							DataKey:    "bar",
						},
					},
				},
			},
		},
		{
			"passing neither value and secretRef in environment vairable",
			false,
			&Container{
				Environment: []EnvVariable{
					{
						Key: "key",
					},
				},
			},
		},
		{
			"passing a valid environment variable with key and secretRef",
			true,
			&Container{
				Environment: []EnvVariable{
					{
						Key: "key",
						SecretRef: &SecretDef{
							SecretName: "foo",
							DataKey:    "bar",
						},
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
				VolumeRef: goutil.StringAddr("invalid_mount_name"),
				MountPath: "/foo/bar",
			},
		},
		{
			"valid name",
			true,
			&Mount{
				VolumeRef: goutil.StringAddr("validmountname"),
				MountPath: "/foo/bar",
			},
		},
		{
			"not an absolute path",
			false,
			&Mount{
				VolumeRef: goutil.StringAddr("test"),
				MountPath: "foo/bar",
			},
		},
		{
			"absolute path",
			true,
			&Mount{
				VolumeRef: goutil.StringAddr("test"),
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

func TestSecretData_Validate(t *testing.T) {
	tests := []struct {
		Name            string
		ExpectedSuccess bool
		SecretData      *SecretData
	}{
		{
			"Only Key is specified, should fail",
			false,
			&SecretData{
				Key: "secretKey",
			},
		},
		{
			"Plaintext and Base64 are specified together, should fail",
			false,
			&SecretData{
				Key:       "secretKey",
				Plaintext: goutil.StringAddr("plaintextData"),
				Base64:    goutil.StringAddr("base64EncodedData"),
			},
		},
		{
			"Plaintext and File are specified together, should fail",
			false,
			&SecretData{
				Key:       "secretKey",
				Plaintext: goutil.StringAddr("plaintextData"),
				File:      goutil.StringAddr("filePath"),
			},
		},
		{
			"Base64 and File are specified together, should fail",
			false,
			&SecretData{
				Key:    "secretKey",
				File:   goutil.StringAddr("filePath"),
				Base64: goutil.StringAddr("base64EncodedData"),
			},
		},
		{
			"Passing Key and Plaintext, should pass",
			true,
			&SecretData{
				Key:       "secretKey",
				Plaintext: goutil.StringAddr("plaintextData"),
			},
		},
		{
			"Passing Key and Base64, should pass",
			true,
			&SecretData{
				Key:    "secretKey",
				Base64: goutil.StringAddr("base64EncodedData"),
			},
		},
		{
			"Passing Key and File, should pass",
			true,
			&SecretData{
				Key:    "secretKey",
				Base64: goutil.StringAddr("filePath"),
			},
		},
	}
	for _, test := range tests {
		{
			err := test.SecretData.validate()
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
		}
	}
}

func TestLogOpenCompose_Validate(t *testing.T) {
	name := "test-service"
	image := "test-image"
	secretRef := &SecretDef{
		SecretName: "secretName",
		DataKey:    "dataKey",
	}
	secretName := "secretName"
	dataKey := "dataKey"

	tests := []struct {
		Name        string
		openCompose *OpenCompose
		logString   string
	}{
		{
			"No root level secrets - env",
			&OpenCompose{
				Version: Version,
				Services: []Service{
					{
						Name: name,
						Containers: []Container{
							{
								Image: image,
								Environment: []EnvVariable{
									{
										Key:       "envKey",
										SecretRef: secretRef,
									},
								},
							},
						},
					},
				},
			},
			"no root level secrets",
		},
		{
			"secretName and dataKey not defined at root level - env",
			&OpenCompose{
				Version: Version,
				Services: []Service{
					{
						Name: name,
						Containers: []Container{
							{
								Image: image,
								Environment: []EnvVariable{
									{
										Key:       "envKey",
										SecretRef: secretRef,
									},
								},
							},
						},
					},
				},
				Secrets: []Secret{
					{
						Name: "notSecretName",
						Data: []SecretData{
							{
								Key:       "notDataKey",
								Plaintext: goutil.StringAddr("randomData"),
							},
						},
					},
				},
			},
			"does not correspond to any root level secret",
		},
		{
			"Secret name present, secret data key absent - env",
			&OpenCompose{
				Version: Version,
				Services: []Service{
					{
						Name: name,
						Containers: []Container{
							{
								Image: image,
								Environment: []EnvVariable{
									{
										Key:       "envKey",
										SecretRef: secretRef,
									},
								},
							},
						},
					},
				},
				Secrets: []Secret{
					{
						Name: secretName,
						Data: []SecretData{
							{
								Key:       "invalidKey",
								Plaintext: goutil.StringAddr("randomData"),
							},
						},
					},
				},
			},
			fmt.Sprintf("secret name: %v found, but the corresponding data key: %v is missing", secretName, dataKey),
		},
		{
			"Secret name and secret data key, both present in root level secrets - env",
			&OpenCompose{
				Version: Version,
				Services: []Service{
					{
						Name: name,
						Containers: []Container{
							{
								Image: image,
								Environment: []EnvVariable{
									{
										Key:       "envKey",
										SecretRef: secretRef,
									},
								},
							},
						},
					},
				},
				Secrets: []Secret{
					{
						Name: secretName,
						Data: []SecretData{
							{
								Key:       dataKey,
								Plaintext: goutil.StringAddr("randomData"),
							},
						},
					},
				},
			},
			"",
		},
		{
			"No root level secrets - mount",
			&OpenCompose{
				Version: Version,
				Services: []Service{
					{
						Name: name,
						Containers: []Container{
							{
								Image: image,
								Mounts: []Mount{
									{
										MountPath: "/foo/bar",
										SecretRef: secretRef,
									},
								},
							},
						},
					},
				},
			},
			"no root level secrets",
		},
		{
			"secretName and dataKey not defined at root level - env",
			&OpenCompose{
				Version: Version,
				Services: []Service{
					{
						Name: name,
						Containers: []Container{
							{
								Image: image,
								Mounts: []Mount{
									{
										MountPath: "/foo/bar",
										SecretRef: secretRef,
									},
								},
							},
						},
					},
				},
				Secrets: []Secret{
					{
						Name: "notSecretName",
						Data: []SecretData{
							{
								Key:       "notDataKey",
								Plaintext: goutil.StringAddr("randomData"),
							},
						},
					},
				},
			},
			"does not correspond to any root level secret",
		},
		{
			"Secret name present, secret data key absent - env",
			&OpenCompose{
				Version: Version,
				Services: []Service{
					{
						Name: name,
						Containers: []Container{
							{
								Image: image,
								Mounts: []Mount{
									{
										MountPath: "/foo/bar",
										SecretRef: secretRef,
									},
								},
							},
						},
					},
				},
				Secrets: []Secret{
					{
						Name: secretName,
						Data: []SecretData{
							{
								Key:       "invalidKey",
								Plaintext: goutil.StringAddr("randomData"),
							},
						},
					},
				},
			},
			fmt.Sprintf("secret name: %v found, but the corresponding data key: %v is missing", secretName, dataKey),
		},
		{
			"Secret name and secret data key, both present in root level secrets - env",
			&OpenCompose{
				Version: Version,
				Services: []Service{
					{
						Name: name,
						Containers: []Container{
							{
								Image: image,
								Mounts: []Mount{
									{
										MountPath: "/foo/bar",
										SecretRef: secretRef,
									},
								},
							},
						},
					},
				},
				Secrets: []Secret{
					{
						Name: secretName,
						Data: []SecretData{
							{
								Key:       dataKey,
								Plaintext: goutil.StringAddr("randomData"),
							},
						},
					},
				},
			},
			"",
		},
	}

	for _, test := range tests {
		{
			t.Run(fmt.Sprintf("Running test: %q", test.Name), func(t *testing.T) {
				stderr := captureStderr(func() { test.openCompose.Validate() })
				if !strings.Contains(stderr, test.logString) {
					t.Errorf("The STDERR output \n%v\n does not contain \n%v", stderr, test.logString)
				}
				if test.logString == "" && stderr != "" {
					t.Errorf("Expected nothing in STDERR, but got %v", stderr)
				}
			})
		}
	}
}

func TestOpenCompose_Validate(t *testing.T) {
	name := "test-service"
	image := "test-image"
	mountName := "mountname"
	mountPath := "/foo/bar"
	secretRef := &SecretDef{
		SecretName: "fooSec",
		DataKey:    "fooKey",
	}

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
			"mount - volumeRef given, but not referenced anywhere",
			false,
			&OpenCompose{
				Version: Version,
				Services: []Service{
					{
						Name: name,
						Containers: []Container{
							{
								Image: image,
								Mounts: []Mount{
									{
										VolumeRef: &mountName,
										MountPath: mountPath,
									},
								},
							},
						},
					},
				},
			},
		},

		{
			"mount - volumeRef given, referenced in emptydir volume",
			true,
			&OpenCompose{
				Version: Version,
				Services: []Service{
					{
						Name: name,
						Containers: []Container{
							{
								Image: image,
								Mounts: []Mount{
									{
										VolumeRef: &mountName,
										MountPath: mountPath,
									},
								},
							},
						},
						EmptyDirVolumes: []EmptyDirVolume{
							{
								Name: mountName,
							},
						},
					},
				},
			},
		},

		{
			"mount - volumeRef given, referenced in root level volumes",
			true,
			&OpenCompose{
				Version: Version,
				Services: []Service{
					{
						Name: name,
						Containers: []Container{
							{
								Image: image,
								Mounts: []Mount{
									{
										VolumeRef: &mountName,
										MountPath: mountPath,
									},
								},
							},
						},
					},
				},
				Volumes: []Volume{
					{
						Name:       mountName,
						AccessMode: "ReadWriteOnce",
						Size:       "100Mi",
					},
				},
			},
		},

		{
			"mount - valid secretRef",
			true,
			&OpenCompose{
				Version: Version,
				Services: []Service{
					{
						Name: name,
						Containers: []Container{
							{
								Image: image,
								Mounts: []Mount{
									{
										SecretRef: secretRef,
										MountPath: mountPath,
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
			}
		})
	}
}
