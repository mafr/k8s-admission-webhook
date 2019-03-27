package validator

import (
	core "k8s.io/api/core/v1"
	res "k8s.io/apimachinery/pkg/api/resource"
	"testing"
)

func TestValidateContainer(t *testing.T) {
	var tests = []struct {
		input     []core.Container
		validator CpuValidator
		expected  bool
	}{
		{
			[]core.Container{newContainer("100m", "300m")},
			CpuValidator{Max: "50m"},
			false,
		},
		{
			[]core.Container{newContainer("100m", "300m"), newContainer("100m", "300m")},
			CpuValidator{Max: "200m"},
			true,
		},
	}

	for i, test := range tests {
		if got := test.validator.validateCpu(test.input); got.Allowed != test.expected {
			t.Errorf("validateCpu(#%d): %v expected, got: %v", i, test.expected, got.Allowed)
		}
	}
}

func TestParseCpu(t *testing.T) {
	var tests = []struct {
		input    string
		expected int64
	}{
		{"2", 2000},
		{"150m", 150},
		{"0.5", 500},
	}

	for _, test := range tests {
		if got := parseCpu(test.input); got != test.expected {
			t.Errorf("parseCpu(%q): %d expected, got: %d", test.input, test.expected, got)
		}
	}
}

func newContainer(cpuRequest, cpuLimit string) core.Container {
	return core.Container{
		Resources: core.ResourceRequirements{
			Requests: core.ResourceList{core.ResourceCPU: res.MustParse(cpuRequest)},
			Limits:   core.ResourceList{core.ResourceCPU: res.MustParse(cpuLimit)},
		},
	}
}
