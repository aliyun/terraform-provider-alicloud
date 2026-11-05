package cloudmonitorservice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitCommonCompareCmsHybridMonitorFcTaskYamlConfigAreEquivalent(t *testing.T) {
	tests := []struct {
		name     string
		yaml1    string
		yaml2    string
		expected bool
		wantErr  bool
	}{
		{
			name:     "equal yamls with different order",
			expected: true,
			yaml1: `
products:
  - namespace: ns1
    metric_info:
      - metric_list: ["metric1", "metric2"]
  - namespace: ns2
    metric_info:
      - metric_list: ["metric3"]
`,
			yaml2: `
products:
  - namespace: ns2
    metric_info:
      - metric_list: ["metric3"]
  - namespace: ns1
    metric_info:
      - metric_list: ["metric1", "metric2"]
`,
		},
		{
			name:     "empty yamls",
			expected: true,
			yaml1:    "",
			yaml2:    "",
		},
		{
			name:    "one valid one invalid",
			wantErr: true,
			yaml1: `
products:
  - namespace: ns1
    metric_info:
      - metric_list: ["metric1"]
`,
			yaml2: "invalid: [yaml",
		},
		{
			name: "same namespace different metrics",
			yaml1: `
products:
  - namespace: ns1
    metric_info:
      - metric_list: ["metric1"]
`,
			yaml2: `
products:
  - namespace: ns1
    metric_info:
      - metric_list: ["metric1", "metric2"]
`,
		},
		{
			name: "completely different",
			yaml1: `
products:
  - namespace: ns1
    metric_info:
      - metric_list: ["metric1"]
`,
			yaml2: `
products:
  - namespace: ns2
    metric_info:
      - metric_list: ["metric2"]
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := compareCmsHybridMonitorFcTaskYamlConfigAreEquivalent(tt.yaml1, tt.yaml2)
			assert.Equal(t, tt.wantErr, err != nil, "unexpected YAML parsing outcome: %v", err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
