package main

import (
	"testing"
)

func TestDeploymentMetricFamilies(t *testing.T) {
	f := deploymentMetricFamilies(nil, nil)
	if f[0].GenerateFunc == nil {
		t.Errorf("GenerateFunc is nil")
	} else {
		t.Logf("GenerateFunc is not nil")
	}
}
