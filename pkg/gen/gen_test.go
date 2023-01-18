package gen

import (
	"strings"
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
	if !strings.HasPrefix(r[nodeID], "NodeID") {
		t.Error(nodeID + " not found in GenerateKeys return")
	}
}
