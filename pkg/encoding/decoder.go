package encoding

import (
	"errors"
	"github.com/tnozicka/opencompose/pkg/encoding/v1"
	"github.com/tnozicka/opencompose/pkg/object"
)

type Decoder interface {
	Unmarshal([]byte) (*object.OpenCompose, error)
}

func GetDecoderFor(data []byte) (Decoder, error) {
	version, err := GetVersion(data)
	if err != nil {
		return nil, err
	}

	switch version {
	case 1:
		return &v1.Decoder{}, nil
	default:
		return nil, errors.New("unsupported version")
	}
}
