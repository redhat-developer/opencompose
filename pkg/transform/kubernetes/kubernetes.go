package kubernetes

import (
	"fmt"

	"github.com/redhat-developer/opencompose/pkg/object"
	_ "k8s.io/client-go/pkg/api/install"
	api_v1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/runtime"
	"k8s.io/client-go/pkg/util/intstr"
)

type Transformer struct{}

// Create k8s services for OpenCompose service
func (t *Transformer) CreateServices(o *object.Service) ([]runtime.Object, error) {
	result := []runtime.Object{}

	for _, c := range o.Containers {
		// We don't want to generate service if there are no ports to be mapped
		if len(c.Ports) == 0 {
			continue
		}

		s := &api_v1.Service{
			ObjectMeta: api_v1.ObjectMeta{
				Name: o.Name,
				Labels: map[string]string{
					"service": o.Name,
				},
			},
			Spec: api_v1.ServiceSpec{
				Selector: map[string]string{
					"service": o.Name,
				},
			},
		}
		for _, p := range c.Ports {
			s.Spec.Ports = append(s.Spec.Ports, api_v1.ServicePort{
				Name:       fmt.Sprintf("port-%d", p.Port.ServicePort),
				Port:       int32(p.Port.ServicePort),
				TargetPort: intstr.FromInt(p.Port.HostPort),
			})
		}
		result = append(result, s)
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
