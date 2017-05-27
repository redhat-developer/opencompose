package util

import (
	"fmt"
	"strings"
)

type SecretDef struct {
	SecretName string
	DataKey    string
}

func SecretRefToSecretDef(secretRef *string) (*SecretDef, error) {
	splitSecretRef := strings.Split(*secretRef, "/")
	if len(splitSecretRef) != 2 {
		return nil, fmt.Errorf("invalid secret syntax, use 'secret: <secret_name>/<data_key>'")
	}

	return &SecretDef{
		SecretName: splitSecretRef[0],
		DataKey:    splitSecretRef[1],
	}, nil

}
