package kubernetes

import (
	"fmt"
	"reflect"
	"testing"

	"encoding/base64"

	"github.com/davecgh/go-spew/spew"
	"github.com/redhat-developer/opencompose/pkg/goutil"
	"github.com/redhat-developer/opencompose/pkg/object"
	api_v1 "k8s.io/client-go/pkg/api/v1"
	ext_v1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/pkg/runtime"
	"k8s.io/client-go/pkg/util/intstr"
)

var (
	name = "test"
	meta = api_v1.ObjectMeta{
		Name: name,
		Labels: map[string]string{
			"service": name,
		},
	}
)

func TestTransformer_CreateSecret(t *testing.T) {
	secretName := "secretname"
	plaintextValue := "secretValue"
	base64Value := "d29yZHByZXNz"
	decodedBase64Value, err := base64.StdEncoding.DecodeString(base64Value)
	if err != nil {
		t.Errorf("Unable to decode the provided base64 encoded value: %v", base64Value)
	}
	tests := []struct {
		Succeed          bool
		Secret           *object.Secret
		KubernetesSecret runtime.Object
	}{
		{
			true,
			&object.Secret{
				Name: secretName,
				Data: []object.SecretData{
					{
						Key:       "secretKey1",
						Plaintext: &plaintextValue,
					},
					{
						Key:    "secretKey2",
						Base64: &base64Value,
					},
				},
			},
			&api_v1.Secret{
				ObjectMeta: api_v1.ObjectMeta{
					Name: secretName,
				},
				Data: map[string][]byte{
					"secretKey1": []byte(plaintextValue),
					"secretKey2": decodedBase64Value,
				},
			},
		},
	}

	transformer := Transformer{}
	for _, tt := range tests {
		ks, err := transformer.CreateSecret(tt.Secret)
		if err != nil {
			if tt.Succeed {
				t.Errorf("Failed to create secret from %+v: %s", tt.Secret, err)
			}
			continue
		}

		if !tt.Succeed {
			t.Errorf("Expected service %+v to fail!", tt.Secret)
			continue
		}

		if !reflect.DeepEqual(ks, tt.KubernetesSecret) {
			t.Errorf("Expected\n%+v\n, got\n%+v", tt.KubernetesSecret, ks)
			continue
		}
	}

}

func TestTransformer_CreateServices(t *testing.T) {
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
					ObjectMeta: meta,
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
					ObjectMeta: meta,
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
					ObjectMeta: meta,
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
					ObjectMeta: meta,
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
					ObjectMeta: meta,
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
					ObjectMeta: meta,
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

func TestTransformer_CreateIngresses(t *testing.T) {
	name := "test"

	tests := []struct {
		Succeed             bool
		Service             *object.Service
		KubernetesIngresses []runtime.Object
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
								Port: object.PortMapping{ServicePort: 80},
								Host: goutil.StringAddr("alpha.127.0.0.1.nip.io"),
								Path: "/admin",
							},
						},
					},
				},
			},
			[]runtime.Object{
				&ext_v1beta1.Ingress{
					ObjectMeta: meta,
					Spec: ext_v1beta1.IngressSpec{
						Rules: []ext_v1beta1.IngressRule{
							{
								Host: "alpha.127.0.0.1.nip.io",
								IngressRuleValue: ext_v1beta1.IngressRuleValue{
									HTTP: &ext_v1beta1.HTTPIngressRuleValue{
										Paths: []ext_v1beta1.HTTPIngressPath{
											{
												Path: "/admin",
												Backend: ext_v1beta1.IngressBackend{
													ServiceName: name,
													ServicePort: intstr.FromInt(80),
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
		},
		{
			true,
			&object.Service{
				Name: name,
				Containers: []object.Container{
					{
						Ports: []object.Port{
							{
								Port: object.PortMapping{ServicePort: 80},
								Host: goutil.StringAddr("alpha.127.0.0.1.nip.io"),
								Path: "/web1",
							},
							{
								Port: object.PortMapping{ServicePort: 80},
								Host: goutil.StringAddr("beta.127.0.0.1.nip.io"),
								Path: "/w1",
							},
							{
								Port: object.PortMapping{ServicePort: 80},
								Host: goutil.StringAddr("alpha.127.0.0.1.nip.io"),
								Path: "/web2",
							},
							{
								Port: object.PortMapping{ServicePort: 80},
								Host: goutil.StringAddr("beta.127.0.0.1.nip.io"),
								Path: "/w2",
							},
							{
								Port: object.PortMapping{ServicePort: 443},
								Host: goutil.StringAddr("alpha.127.0.0.1.nip.io"),
								Path: "/admin1",
							},
							{
								Port: object.PortMapping{ServicePort: 443},
								Host: goutil.StringAddr("beta.127.0.0.1.nip.io"),
								Path: "/adm1",
							},
						},
					},
					{
						Ports: []object.Port{
							{
								Port: object.PortMapping{ServicePort: 443},
								Host: goutil.StringAddr("alpha.127.0.0.1.nip.io"),
								Path: "/admin2",
							},
							{
								Port: object.PortMapping{ServicePort: 443},
								Host: goutil.StringAddr("beta.127.0.0.1.nip.io"),
								Path: "/adm2",
							},
						},
					},
				},
			},
			[]runtime.Object{
				&ext_v1beta1.Ingress{
					ObjectMeta: meta,
					Spec: ext_v1beta1.IngressSpec{
						Rules: []ext_v1beta1.IngressRule{
							{
								Host: "alpha.127.0.0.1.nip.io",
								IngressRuleValue: ext_v1beta1.IngressRuleValue{
									HTTP: &ext_v1beta1.HTTPIngressRuleValue{
										Paths: []ext_v1beta1.HTTPIngressPath{
											{
												Path: "/web1",
												Backend: ext_v1beta1.IngressBackend{
													ServiceName: name,
													ServicePort: intstr.FromInt(80),
												},
											},
											{
												Path: "/web2",
												Backend: ext_v1beta1.IngressBackend{
													ServiceName: name,
													ServicePort: intstr.FromInt(80),
												},
											},
											{
												Path: "/admin1",
												Backend: ext_v1beta1.IngressBackend{
													ServiceName: name,
													ServicePort: intstr.FromInt(443),
												},
											},
											{
												Path: "/admin2",
												Backend: ext_v1beta1.IngressBackend{
													ServiceName: name,
													ServicePort: intstr.FromInt(443),
												},
											},
										},
									},
								},
							},
							{
								Host: "beta.127.0.0.1.nip.io",
								IngressRuleValue: ext_v1beta1.IngressRuleValue{
									HTTP: &ext_v1beta1.HTTPIngressRuleValue{
										Paths: []ext_v1beta1.HTTPIngressPath{
											{
												Path: "/w1",
												Backend: ext_v1beta1.IngressBackend{
													ServiceName: name,
													ServicePort: intstr.FromInt(80),
												},
											},
											{
												Path: "/w2",
												Backend: ext_v1beta1.IngressBackend{
													ServiceName: name,
													ServicePort: intstr.FromInt(80),
												},
											},
											{
												Path: "/adm1",
												Backend: ext_v1beta1.IngressBackend{
													ServiceName: name,
													ServicePort: intstr.FromInt(443),
												},
											},
											{
												Path: "/adm2",
												Backend: ext_v1beta1.IngressBackend{
													ServiceName: name,
													ServicePort: intstr.FromInt(443),
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
		},
	}

	transformer := Transformer{}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			ks, err := transformer.CreateIngresses(tt.Service)
			if err != nil {
				if tt.Succeed {
					t.Fatalf("Failed to create ingresses from %+v: %s", tt.Service, err)
				}
				return
			}

			if !tt.Succeed {
				t.Fatal(spew.Errorf("Expected ingresses %#+v to fail!", tt.Service))
			}

			if !reflect.DeepEqual(ks, tt.KubernetesIngresses) {
				t.Fatal(spew.Errorf("Expected:\n%#+v\n, got:\n%#+v", tt.KubernetesIngresses, ks))
			}
		})
	}
}

func TestTransformer_CreateDeployments(t *testing.T) {
	name := "test"
	podname := fmt.Sprintf("%s-%d", name, 0)
	image := "docker.io/test"
	sMeta := api_v1.ObjectMeta{
		Name: name,
		Labels: map[string]string{
			"service": name,
		},
	}
	strategy := ext_v1beta1.DeploymentStrategy{
		Type: ext_v1beta1.RollingUpdateDeploymentStrategyType,
	}

	tests := []struct {
		Name           string
		Succeed        bool
		Service        *object.Service
		K8sDeployments []runtime.Object
	}{
		{
			"When no replica field given",
			true,
			&object.Service{
				Name: name,
				Containers: []object.Container{
					{
						Image: image,
					},
				},
			},
			[]runtime.Object{
				&ext_v1beta1.Deployment{
					ObjectMeta: sMeta,
					Spec: ext_v1beta1.DeploymentSpec{
						Strategy: strategy,
						Template: api_v1.PodTemplateSpec{
							ObjectMeta: api_v1.ObjectMeta{
								Labels: map[string]string{
									"service": name,
								},
							},
							Spec: api_v1.PodSpec{
								Containers: []api_v1.Container{
									{
										Name:  podname,
										Image: image,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			"When valid replica value given",
			true,
			&object.Service{
				Name:     name,
				Replicas: goutil.Int32Addr(1),
				Containers: []object.Container{
					{
						Image: image,
					},
				},
			},
			[]runtime.Object{
				&ext_v1beta1.Deployment{
					ObjectMeta: sMeta,
					Spec: ext_v1beta1.DeploymentSpec{
						Strategy: strategy,
						Replicas: goutil.Int32Addr(1),
						Template: api_v1.PodTemplateSpec{
							ObjectMeta: api_v1.ObjectMeta{
								Labels: map[string]string{
									"service": name,
								},
							},
							Spec: api_v1.PodSpec{
								Containers: []api_v1.Container{
									{
										Name:  podname,
										Image: image,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	transformer := Transformer{}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			kd, err := transformer.CreateDeployments(test.Service)
			if err != nil {
				if test.Succeed {
					t.Errorf("Failed to create deployment from %#v\nErr: %s", test.Service, err)
				}
				return
			}

			if !test.Succeed {
				t.Errorf("Expected failure, but succeeded, service: %#v\nConverted k8s Deployment: %#v", test.Service, spew.Sprint(kd))
				return
			}

			if !reflect.DeepEqual(kd, test.K8sDeployments) {
				t.Errorf("Expected: %#v\nGot: %#v\n", spew.Sprint(test.K8sDeployments), spew.Sprint(kd))
				return
			}
		})
	}
}
