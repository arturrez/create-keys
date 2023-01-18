package gen

import (
	"fmt"
	"strings"
	"testing"
)

func TestGenerateKeys(t *testing.T) {
	r, err := GenerateKeys()
	fmt.Println(r)
	if err != nil {
		t.Error("GenerateKeys returns error")
	}
	if _, ok := r[stakerCert]; !ok {
		t.Error(stakerCert + " not found in GenerateKeys return")
	}
	if !strings.HasPrefix("NodeID", r[nodeID]) {
		t.Error(nodeID + " not found in GenerateKeys return")
	}
}
