package kubernetes

import (
	"fmt"

	"github.com/redhat-developer/opencompose/pkg/object"
	_ "k8s.io/client-go/pkg/api/install"
	api_v1 "k8s.io/client-go/pkg/api/v1"
	_ "k8s.io/client-go/pkg/apis/extensions/install"
	ext_v1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/pkg/runtime"
	"k8s.io/client-go/pkg/util/intstr"
)

type Transformer struct{}

// Create k8s services for OpenCompose service
func (t *Transformer) CreateServices(o *object.Service) ([]runtime.Object, error) {
	result := []runtime.Object{}

	Service := func() *api_v1.Service {
		return &api_v1.Service{
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
	}

	is := Service()
	is.Spec.Type = api_v1.ServiceTypeClusterIP

	es := Service()
	es.Spec.Type = api_v1.ServiceTypeLoadBalancer

	for _, c := range o.Containers {
		// We don't want to generate service if there are no ports to be mapped
		if len(c.Ports) == 0 {
			continue
		}

		for _, p := range c.Ports {
			var s *api_v1.Service
			switch p.Type {
			case object.PortType_Internal:
				s = is
			case object.PortType_External:
				s = es
			default:
				// There is a mistake in our code; and in Golang because it doesn't have strongly typed enumerations :)
				return result, fmt.Errorf("Internal error: unknown PortType %#v", p.Type)
			}

			s.Spec.Ports = append(s.Spec.Ports, api_v1.ServicePort{
				Name:       fmt.Sprintf("port-%d", p.Port.ServicePort),
				Port:       int32(p.Port.ServicePort),
				TargetPort: intstr.FromInt(p.Port.ContainerPort),
			})
		}
	}

	if len(is.Spec.Ports) > 0 {
		result = append(result, is)
	}

	if len(es.Spec.Ports) > 0 {
		result = append(result, es)
	}

	return result, nil
}

// Create k8s deployments for OpenCompose service
func (t *Transformer) CreateDeployments(o *object.Service) ([]runtime.Object, error) {
	result := []runtime.Object{}

	d := &ext_v1beta1.Deployment{
		ObjectMeta: api_v1.ObjectMeta{
			Name: o.Name,
			Labels: map[string]string{
				"service": o.Name,
			},
		},
		Spec: ext_v1beta1.DeploymentSpec{
			Strategy: ext_v1beta1.DeploymentStrategy{
				// TODO: make it configurable
				Type: ext_v1beta1.RollingUpdateDeploymentStrategyType,
				// TODO: make it configurable
				RollingUpdate: nil,
			},
			Template: api_v1.PodTemplateSpec{
				ObjectMeta: api_v1.ObjectMeta{
					Labels: map[string]string{
						"service": o.Name,
					},
				},
				Spec: api_v1.PodSpec{},
			},
		},
	}

	for i, c := range o.Containers {
		kc := api_v1.Container{
			Name:  fmt.Sprintf("%s-%d", o.Name, i),
			Image: c.Image,
		}

		for _, e := range c.Environment {
			kc.Env = append(kc.Env, api_v1.EnvVar{
				Name:  e.Key,
				Value: e.Value,
			})
		}

		for _, p := range c.Ports {
			kc.Ports = append(kc.Ports, api_v1.ContainerPort{
				Name:          fmt.Sprintf("port-%d", p.Port.ContainerPort),
				ContainerPort: int32(p.Port.ContainerPort),
			})
		}

		d.Spec.Template.Spec.Containers = append(d.Spec.Template.Spec.Containers, kc)

		result = append(result, d)
	}

	return result, nil
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
