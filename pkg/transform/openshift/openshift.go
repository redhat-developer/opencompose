package openshift

import (
	//"github.com/redhat-developer/opencompose/pkg/object"
	"github.com/redhat-developer/opencompose/pkg/transform/kubernetes"
	//"k8s.io/client-go/pkg/runtime"
)

type Transformer struct {
	kubernetes.Transformer
}

//func (t *Transformer) Transform(o *object.OpenCompose) ([]runtime.Object, error) {
//	return nil, nil
//}
