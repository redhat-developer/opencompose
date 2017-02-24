package kubernetes

import (
	"reflect"
	"testing"

	"github.com/redhat-developer/opencompose/pkg/object"
	api_v1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/runtime"
	"k8s.io/client-go/pkg/util/intstr"
)

func TestTransformer_CreateServices(t *testing.T) {
	name := "test"
	sMeta := api_v1.ObjectMeta{
		Name: name,
		Labels: map[string]string{
			"service": name,
		},
	}
	sSelector := map[string]string{
		"service": name,
	}

	tests := []struct {
		Succeed            bool
		Service            *object.Service
		KubernetesServices []runtime.Object
	}{
		{
			true,
			&object.Service{
				Name: name,
				Containers: []object.Container{
					{
						Ports: []object.Port{},
					},
				},
			},
			[]runtime.Object{},
		},
		{
			true,
			&object.Service{
				Name: name,
				Containers: []object.Container{
					{
						Ports: []object.Port{
							{
								Port: object.PortMapping{
									ContainerPort: 8080,
									ServicePort:   80,
								},
								Type: object.PortType(object.PortType_Internal),
							},
						},
					},
				},
			},
			[]runtime.Object{
				&api_v1.Service{
					ObjectMeta: sMeta,
					Spec: api_v1.ServiceSpec{
						Selector: sSelector,
						Ports: []api_v1.ServicePort{
							{
								Name:       "port-80",
								Port:       int32(80),
								TargetPort: intstr.FromInt(8080),
							},
						},
						Type: api_v1.ServiceTypeClusterIP,
					},
				},
			},
		},
		{
			true,
			&object.Service{
				Name: name,
				Containers: []object.Container{
					{
						Ports: []object.Port{
							{
								Port: object.PortMapping{
									ContainerPort: 8080,
									ServicePort:   80,
								},
								Type: object.PortType(object.PortType_External),
							},
						},
					},
				},
			},
			[]runtime.Object{
				&api_v1.Service{
					ObjectMeta: sMeta,
					Spec: api_v1.ServiceSpec{
						Selector: sSelector,
						Ports: []api_v1.ServicePort{
							{
								Name:       "port-80",
								Port:       int32(80),
								TargetPort: intstr.FromInt(8080),
							},
						},
						Type: api_v1.ServiceTypeLoadBalancer,
					},
				},
			},
		},
		{
			true,
			&object.Service{
				Name: name,
				Containers: []object.Container{
					{
						Ports: []object.Port{
							{
								Port: object.PortMapping{
									ContainerPort: 8080,
									ServicePort:   80,
								},
								Type: object.PortType(object.PortType_Internal),
							},
							{
								Port: object.PortMapping{
									ContainerPort: 5000,
									ServicePort:   443,
								},
								Type: object.PortType(object.PortType_External),
							},
							{
								Port: object.PortMapping{
									ContainerPort: 8081,
									ServicePort:   81,
								},
								Type: object.PortType(object.PortType_Internal),
							},
							{
								Port: object.PortMapping{
									ContainerPort: 5001,
									ServicePort:   444,
								},
								Type: object.PortType(object.PortType_External),
							},
						},
					},
				},
			},
			[]runtime.Object{
				&api_v1.Service{
					ObjectMeta: sMeta,
					Spec: api_v1.ServiceSpec{
						Selector: sSelector,
						Ports: []api_v1.ServicePort{
							{
								Name:       "port-80",
								Port:       int32(80),
								TargetPort: intstr.FromInt(8080),
							},
							{
								Name:       "port-81",
								Port:       int32(81),
								TargetPort: intstr.FromInt(8081),
							},
						},
						Type: api_v1.ServiceTypeClusterIP,
					},
				},
				&api_v1.Service{
					ObjectMeta: sMeta,
					Spec: api_v1.ServiceSpec{
						Selector: sSelector,
						Ports: []api_v1.ServicePort{
							{
								Name:       "port-443",
								Port:       int32(443),
								TargetPort: intstr.FromInt(5000),
							},
							{
								Name:       "port-444",
								Port:       int32(444),
								TargetPort: intstr.FromInt(5001),
							},
						},
						Type: api_v1.ServiceTypeLoadBalancer,
					},
				},
			},
		},
		{
			true,
			&object.Service{
				Name: name,
				Containers: []object.Container{
					{
						Ports: []object.Port{
							{
								Port: object.PortMapping{
									ContainerPort: 8080,
									ServicePort:   80,
								},
								Type: object.PortType(object.PortType_Internal),
							},
							{
								Port: object.PortMapping{
									ContainerPort: 5000,
									ServicePort:   443,
								},
								Type: object.PortType(object.PortType_External),
							},
							{
								Port: object.PortMapping{
									ContainerPort: 8081,
									ServicePort:   81,
								},
								Type: object.PortType(object.PortType_Internal),
							},
							{
								Port: object.PortMapping{
									ContainerPort: 5001,
									ServicePort:   444,
								},
								Type: object.PortType(object.PortType_External),
							},
						},
					},
					{
						Ports: []object.Port{
							{
								Port: object.PortMapping{
									ContainerPort: 9080,
									ServicePort:   90,
								},
								Type: object.PortType(object.PortType_Internal),
							},
							{
								Port: object.PortMapping{
									ContainerPort: 6000,
									ServicePort:   543,
								},
								Type: object.PortType(object.PortType_External),
							},
							{
								Port: object.PortMapping{
									ContainerPort: 9081,
									ServicePort:   91,
								},
								Type: object.PortType(object.PortType_Internal),
							},
							{
								Port: object.PortMapping{
									ContainerPort: 6001,
									ServicePort:   544,
								},
								Type: object.PortType(object.PortType_External),
							},
						},
					},
				},
			},
			[]runtime.Object{
				&api_v1.Service{
					ObjectMeta: sMeta,
					Spec: api_v1.ServiceSpec{
						Selector: sSelector,
						Ports: []api_v1.ServicePort{
							{
								Name:       "port-80",
								Port:       int32(80),
								TargetPort: intstr.FromInt(8080),
							},
							{
								Name:       "port-81",
								Port:       int32(81),
								TargetPort: intstr.FromInt(8081),
							},
							{
								Name:       "port-90",
								Port:       int32(90),
								TargetPort: intstr.FromInt(9080),
							},
							{
								Name:       "port-91",
								Port:       int32(91),
								TargetPort: intstr.FromInt(9081),
							},
						},
						Type: api_v1.ServiceTypeClusterIP,
					},
				},
				&api_v1.Service{
					ObjectMeta: sMeta,
					Spec: api_v1.ServiceSpec{
						Selector: sSelector,
						Ports: []api_v1.ServicePort{
							{
								Name:       "port-443",
								Port:       int32(443),
								TargetPort: intstr.FromInt(5000),
							},
							{
								Name:       "port-444",
								Port:       int32(444),
								TargetPort: intstr.FromInt(5001),
							},
							{
								Name:       "port-543",
								Port:       int32(543),
								TargetPort: intstr.FromInt(6000),
							},
							{
								Name:       "port-544",
								Port:       int32(544),
								TargetPort: intstr.FromInt(6001),
							},
						},
						Type: api_v1.ServiceTypeLoadBalancer,
					},
				},
			},
		},
	}

	transformer := Transformer{}
	for _, tt := range tests {
		ks, err := transformer.CreateServices(tt.Service)
		if err != nil {
			if tt.Succeed {
				t.Errorf("Failed to create services from %+v: %s", tt.Service, err)
			}
			continue
		}

		if !tt.Succeed {
			t.Errorf("Expected service %+v to fail!", tt.Service)
			continue
		}

		if !reflect.DeepEqual(ks, tt.KubernetesServices) {
			t.Errorf("Expected\n%+v\n, got\n%+v", tt.KubernetesServices, ks)
			continue
		}
	}
}
