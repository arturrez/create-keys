package gen

import (
	"testing"
)

func TestGenerateKeys(t *testing.T) {
	r, err := GenerateKeys()
	if err != nil {
		t.Error("GenerateKeys returns error")
	}
	if _, ok := r[stakerCert]; !ok {
		t.Error(stakerCert + " not found in GenerateKeys return")
	}
}
