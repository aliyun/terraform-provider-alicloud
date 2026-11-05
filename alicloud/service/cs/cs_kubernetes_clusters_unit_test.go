package cs

import (
	"testing"
)

func TestUnitCsNodePrivateIp(t *testing.T) {
	tests := []struct {
		name      string
		ipAddress []string
		want      string
	}{
		{
			name:      "nil address list",
			ipAddress: nil,
			want:      "",
		},
		{
			name:      "empty address list",
			ipAddress: []string{},
			want:      "",
		},
		{
			name:      "single address",
			ipAddress: []string{"192.168.0.1"},
			want:      "192.168.0.1",
		},
		{
			name:      "multiple addresses keep the first one",
			ipAddress: []string{"192.168.0.1", "192.168.0.2"},
			want:      "192.168.0.1",
		},
		{
			name:      "empty first address is preserved",
			ipAddress: []string{""},
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CsNodePrivateIp(tt.ipAddress); got != tt.want {
				t.Errorf("csNodePrivateIp(%v) want %q got %q", tt.ipAddress, tt.want, got)
			}
		})
	}
}
