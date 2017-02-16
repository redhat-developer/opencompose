package encoding

import (
	"fmt"

	"github.com/redhat-developer/opencompose/pkg/encoding/v1"
	"github.com/redhat-developer/opencompose/pkg/object"
)

type Decoder interface {
	Decode([]byte) (*object.OpenCompose, error)
}

func GetDecoderFor(data []byte) (Decoder, error) {
	version, err := GetVersion(data)
	if err != nil {
		return nil, err
	}

	switch version {
	case v1.Version:
		return &v1.Decoder{}, nil
	default:
		return nil, fmt.Errorf("unsupported version %q", version)
	}
}
