package object

import (
	"fmt"
	"testing"

	"github.com/redhat-developer/opencompose/pkg/goutil"
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
