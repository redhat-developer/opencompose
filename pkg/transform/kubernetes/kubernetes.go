package kubernetes

import (
	"fmt"

	"github.com/tnozicka/opencompose/pkg/object"
	api_v1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/runtime"
)

type Transformer struct{}

// Create k8s services for OpenCompose service
func (t *Transformer) CreateServices(o *object.Service) ([]runtime.Object, error) {
	// TODO: go through all container mapping

	result := []runtime.Object{
		&api_v1.Service{
			ObjectMeta: api_v1.ObjectMeta{
				Name: "test",
			},
		},
	}

	return result, nil
}

// Create k8s deployments for OpenCompose service
func (t *Transformer) CreateDeployments(o *object.Service) ([]runtime.Object, error) {
	// TODO: implement
	return nil, nil
}

func (t *Transformer) TransformServices(services []object.Service) ([]runtime.Object, error) {
	result := []runtime.Object{}

	for _, service := range services {
		// create k8s services
		objects, err := t.CreateServices(&service)
		if err != nil {
			return nil, fmt.Errorf("failed to transform service: %s", err)
		}
		result = append(result, objects...)

		// create k8s deployments
		objects, err = t.CreateDeployments(&service)
		if err != nil {
			return nil, fmt.Errorf("failed to create deployments: %s", err)
		}
		result = append(result, objects...)
	}

	return result, nil
}

func (t *Transformer) TransformVolumes(volumes []object.Volume) ([]runtime.Object, error) {
	return nil, nil
}

func (t *Transformer) Transform(o *object.OpenCompose) ([]runtime.Object, error) {
	result := []runtime.Object{}

	// services
	serviceObjects, err := t.TransformServices(o.Services)
	if err != nil {
		return nil, fmt.Errorf("failed to transform services: %s", err)
	}
	result = append(result, serviceObjects...)

	// volumes
	volumeObjects, err := t.TransformVolumes(o.Volumes)
	if err != nil {
		return nil, fmt.Errorf("failed to transform volumes: %s", err)
	}
	result = append(result, volumeObjects...)

	return result, nil
}
