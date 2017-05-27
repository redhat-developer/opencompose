package util

import (
	"reflect"
	"testing"
)

func TestSecretRefToSecretDef(t *testing.T) {
	tests := []struct {
		Name      string
		secretRef string
		secretDef *SecretDef
	}{
		{
			"Test valid secretRef",
			"secretname/datakey",
			&SecretDef{
				SecretName: "secretname",
				DataKey:    "datakey",
			},
		},
		{
			"Test invalid secretRef, no '/'",
			"secretname,datakey",
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			sDef, _ := SecretRefToSecretDef(&test.secretRef)

			if !reflect.DeepEqual(sDef, test.secretDef) {
				t.Errorf("Expected -\n%v\nGot -\n%v", *test.secretDef, *sDef)
			}
		})
	}
}
