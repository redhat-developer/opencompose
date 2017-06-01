package kubernetes

import (
	"fmt"

	"encoding/base64"
	"io/ioutil"

	"github.com/redhat-developer/opencompose/pkg/object"
	transformUtil "github.com/redhat-developer/opencompose/pkg/transform/util"
	"github.com/redhat-developer/opencompose/pkg/util"
	_ "k8s.io/client-go/pkg/api/install"
	"k8s.io/client-go/pkg/api/resource"
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
		serviceLabels := map[string]string(o.Labels)
		return &api_v1.Service{
			ObjectMeta: api_v1.ObjectMeta{
				Name: o.Name,
				Labels: *util.MergeMaps(
					// The map containing `"service": o.Name` should always be
					// passed later to avoid being overridden by util.MergeMaps()
					&serviceLabels,
					&map[string]string{
						"service": o.Name,
					},
				),
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

// Create k8s ingresses for OpenCompose service
func (t *Transformer) CreateIngresses(o *object.Service) ([]runtime.Object, error) {
	result := []runtime.Object{}
	serviceLabels := map[string]string(o.Labels)

	i := &ext_v1beta1.Ingress{
		ObjectMeta: api_v1.ObjectMeta{
			Name: o.Name,
			Labels: *util.MergeMaps(
				// The map containing `"service": o.Name` should always be
				// passed later to avoid being overridden by util.MergeMaps()
				&serviceLabels,
				&map[string]string{
					"service": o.Name,
				},
			),
		},
	}

	for _, c := range o.Containers {
		// We don't want to generate ingress if there are no ports to be mapped
		if len(c.Ports) == 0 {
			continue
		}

		for _, p := range c.Ports {
			if p.Host == nil {
				// Not Ingress
				continue
			}

			host := *p.Host
			var rule *ext_v1beta1.IngressRule
			for idx := range i.Spec.Rules {
				r := &i.Spec.Rules[idx]
				if r.Host == host {
					rule = r
					break
				}
			}
			if rule == nil {
				rule = &ext_v1beta1.IngressRule{
					Host: host,
					IngressRuleValue: ext_v1beta1.IngressRuleValue{
						HTTP: &ext_v1beta1.HTTPIngressRuleValue{},
					},
				}
				i.Spec.Rules = append(i.Spec.Rules, *rule)
			}

			rule.HTTP.Paths = append(rule.HTTP.Paths, ext_v1beta1.HTTPIngressPath{
				Path: p.Path,
				Backend: ext_v1beta1.IngressBackend{
					ServiceName: o.Name,
					ServicePort: intstr.FromInt(p.Port.ServicePort),
				},
			})
		}
	}

	if len(i.Spec.Rules) > 0 {
		result = append(result, i)
	}

	return result, nil
}

// Create k8s deployments for OpenCompose service
func (t *Transformer) CreateDeployments(s *object.Service) ([]runtime.Object, error) {
	result := []runtime.Object{}
	serviceLabels := map[string]string(s.Labels)

	d := &ext_v1beta1.Deployment{
		ObjectMeta: api_v1.ObjectMeta{
			Name: s.Name,
			Labels: *util.MergeMaps(
				// The map containing `"service": s.Name` should always be
				// passed later to avoid being overridden by util.MergeMaps()
				&serviceLabels,
				&map[string]string{
					"service": s.Name,
				},
			),
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
					Labels: *util.MergeMaps(
						// The map containing `"service": s.Name` should always be
						// passed later to avoid being overridden by util.MergeMaps()
						&serviceLabels,
						&map[string]string{
							"service": s.Name,
						},
					),
				},
				Spec: api_v1.PodSpec{},
			},
		},
	}

	d.Spec.Replicas = s.Replicas

	for i, c := range s.Containers {
		kc := api_v1.Container{
			Name:  fmt.Sprintf("%s-%d", s.Name, i),
			Image: c.Image,
		}

		for _, e := range c.Environment {
			env := api_v1.EnvVar{
				Name: e.Key,
			}

			if e.Value != nil {
				env.Value = *e.Value
			}

			if e.SecretRef != nil {

				secretDef, err := transformUtil.SecretRefToSecretDef(e.SecretRef)
				if err != nil {
					return nil, fmt.Errorf("Invalid secret %v specified in environment: %v", *e.SecretRef, err)
				}

				env.ValueFrom = &api_v1.EnvVarSource{
					SecretKeyRef: &api_v1.SecretKeySelector{
						LocalObjectReference: api_v1.LocalObjectReference{
							Name: secretDef.SecretName,
						},
						Key: secretDef.DataKey,
					},
				}
			}

			kc.Env = append(kc.Env, env)
		}

		for _, p := range c.Ports {
			kc.Ports = append(kc.Ports, api_v1.ContainerPort{
				Name:          fmt.Sprintf("port-%d", p.Port.ContainerPort),
				ContainerPort: int32(p.Port.ContainerPort),
			})
		}

		// TODO: It is assumed that the check is done about the existence of volume in root level volume section
		for _, mount := range c.Mounts {
			volumeMount := api_v1.VolumeMount{
				ReadOnly:  mount.ReadOnly,
				MountPath: mount.MountPath,
				SubPath:   mount.VolumeSubPath,
			}

			if mount.VolumeRef != nil {
				volumeMount.Name = *mount.VolumeRef
			}

			if mount.SecretRef != nil {
				volumeMount.ReadOnly = true

				secretDef, err := transformUtil.SecretRefToSecretDef(mount.SecretRef)
				if err != nil {
					return nil, fmt.Errorf("Invalid secret %v specified in volume: %v", *mount.SecretRef, err)
				}

				secretVolume := api_v1.Volume{
					Name: secretDef.SecretName,
					VolumeSource: api_v1.VolumeSource{
						Secret: &api_v1.SecretVolumeSource{
							SecretName: secretDef.SecretName,
							Items: []api_v1.KeyToPath{
								{
									Key:  secretDef.DataKey,
									Path: secretDef.DataKey,
								},
							},
						},
					},
				}

				d.Spec.Template.Spec.Volumes = append(d.Spec.Template.Spec.Volumes, secretVolume)

			}

			kc.VolumeMounts = append(kc.VolumeMounts, volumeMount)

			// if this mount does not exist in emptydir then this is coming from root level volumes directive
			// if tomorrow we add support for ConfigMaps or Secrets mounted as volumes the check should be done
			// here to see if it is not coming from configMaps or Secrets
			if (mount.VolumeRef != nil && !s.EmptyDirVolumeExists(*mount.VolumeRef)) || (mount.SecretRef == nil) {
				volume := api_v1.Volume{
					Name: *mount.VolumeRef,
					VolumeSource: api_v1.VolumeSource{
						PersistentVolumeClaim: &api_v1.PersistentVolumeClaimVolumeSource{
							ClaimName: *mount.VolumeRef,
						},
					},
				}
				d.Spec.Template.Spec.Volumes = append(d.Spec.Template.Spec.Volumes, volume)
			}
		}

		d.Spec.Template.Spec.Containers = append(d.Spec.Template.Spec.Containers, kc)
	}

	// make entry of emptydir in deployment volume directive
	for _, emptyDir := range s.EmptyDirVolumes {
		volume := api_v1.Volume{
			Name: emptyDir.Name,
			VolumeSource: api_v1.VolumeSource{
				EmptyDir: &api_v1.EmptyDirVolumeSource{},
			},
		}
		d.Spec.Template.Spec.Volumes = append(d.Spec.Template.Spec.Volumes, volume)
	}

	result = append(result, d)

	return result, nil
}

// Create Kubernetes Persistent Volume Claim
func (t *Transformer) CreatePVC(volume object.Volume) (runtime.Object, error) {

	size, err := resource.ParseQuantity(volume.Size)
	if err != nil {
		return nil, err
	}

	pvc := &api_v1.PersistentVolumeClaim{
		ObjectMeta: api_v1.ObjectMeta{
			Name: volume.Name,
		},
		Spec: api_v1.PersistentVolumeClaimSpec{
			Resources: api_v1.ResourceRequirements{
				Requests: api_v1.ResourceList{
					api_v1.ResourceStorage: size,
				},
			},
		},
	}

	switch volume.AccessMode {
	case "ReadWriteOnce":
		pvc.Spec.AccessModes = []api_v1.PersistentVolumeAccessMode{api_v1.ReadWriteOnce}
	case "ReadOnlyMany":
		pvc.Spec.AccessModes = []api_v1.PersistentVolumeAccessMode{api_v1.ReadOnlyMany}
	case "ReadWriteMany":
		pvc.Spec.AccessModes = []api_v1.PersistentVolumeAccessMode{api_v1.ReadWriteMany}
	default:
		return nil, fmt.Errorf("invalid accessMode: %q, must be either %q, %q or %q", volume.AccessMode, "ReadWriteOnce", "ReadOnlyMany", "ReadWriteMany")
	}

	if volume.StorageClass != nil {
		pvc.ObjectMeta.Annotations = make(map[string]string)
		pvc.ObjectMeta.Annotations["volume.beta.kubernetes.io/storage-class"] = *volume.StorageClass
	}

	return pvc, nil
}

func (t *Transformer) CreateSecret(secret *object.Secret) (runtime.Object, error) {

	kubeSecretData := make(map[string][]byte)

	for _, data := range secret.Data {
		switch {

		case data.Plaintext != nil:
			kubeSecretData[data.Key] = []byte(*data.Plaintext)

		case data.File != nil:
			fileData, err := ioutil.ReadFile(*data.File)
			if err != nil {
				return nil, fmt.Errorf("Unable to read %v: %v", *data.File, err)
			}
			kubeSecretData[data.Key] = fileData

		case data.Base64 != nil:
			decodedData, err := base64.StdEncoding.DecodeString(*data.Base64)
			if err != nil {
				return nil, fmt.Errorf("Unable to base64 decode data %v: %v", *data.Base64, err)
			}
			kubeSecretData[data.Key] = decodedData
		}
	}

	sec := &api_v1.Secret{
		ObjectMeta: api_v1.ObjectMeta{
			Name: secret.Name,
		},
		Data: kubeSecretData,
	}

	// TODO: add secret validation here
	// TODO: refer https://github.com/redhat-developer/opencompose/issues/139

	return sec, nil
}

func (t *Transformer) TransformServices(services []object.Service) ([]runtime.Object, error) {
	result := []runtime.Object{}

	for _, service := range services {
		// create k8s services
		objects, err := t.CreateServices(&service)
		if err != nil {
			return nil, fmt.Errorf("failed to generate services: %s", err)
		}
		result = append(result, objects...)

		// create k8s ingresses
		objects, err = t.CreateIngresses(&service)
		if err != nil {
			return nil, fmt.Errorf("failed to generate ingresses: %s", err)
		}
		result = append(result, objects...)

		// create k8s deployments
		objects, err = t.CreateDeployments(&service)
		if err != nil {
			return nil, fmt.Errorf("failed to generate deployments: %s", err)
		}
		result = append(result, objects...)
	}

	return result, nil
}

func (t *Transformer) TransformVolumes(volumes []object.Volume) ([]runtime.Object, error) {
	result := []runtime.Object{}

	for _, volume := range volumes {
		// create pvc
		object, err := t.CreatePVC(volume)
		if err != nil {
			return nil, fmt.Errorf("failed to create PVC for volume %q: %s", volume.Name, err)
		}

		result = append(result, object)
	}

	return result, nil
}

func (t *Transformer) TransformSecrets(secrets []object.Secret) ([]runtime.Object, error) {
	result := []runtime.Object{}

	// create secrets
	for _, secret := range secrets {
		object, err := t.CreateSecret(&secret)
		if err != nil {
			return nil, fmt.Errorf("failed to create secret %v: %v", secret.Name, err)
		}

		result = append(result, object)
	}

	return result, nil
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

	// secrets
	secretObjects, err := t.TransformSecrets(o.Secrets)
	if err != nil {
		return nil, fmt.Errorf("failed to transform secrets: %s", err)
	}
	result = append(result, secretObjects...)

	return result, nil
}
