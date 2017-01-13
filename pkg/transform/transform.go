package transform

import (
	"github.com/tnozicka/opencompose/pkg/object"
	"k8s.io/client-go/pkg/runtime"
)

// Transformer interface
type Transformer interface {
	// Transform OpenCompose into Kubernetes/OpenShift objects
	Transform(o *object.OpenCompose) ([]runtime.Object, error)
}
