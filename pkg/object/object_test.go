package object

import (
	"fmt"
	"testing"

	"github.com/redhat-developer/opencompose/pkg/goutil"
)

const (
	Version = "0.1-dev" // TODO: replace with "1" once we reach that point
)

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
			"Valid replica value: -1",
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
