package rds

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnitDBProxyMultiAZGuard(t *testing.T) {
	multiAZ := []map[string]interface{}{
		{"zone_id": "cn-hangzhou-h"},
		{"zone_id": "cn-hangzhou-i"},
	}
	singleAZ := []map[string]interface{}{{"zone_id": "cn-hangzhou-h"}}

	for _, tc := range []struct {
		name          string
		networkChange bool
		nodes         bool
		layout        []map[string]interface{}
		errContains   string
	}{
		{
			// The network-change refusal is unconditional: an explicit target
			// db_proxy_nodes must not bypass the multi-AZ guard.
			name:          "multi-AZ network change with configured nodes is refused",
			networkChange: true,
			nodes:         true,
			layout:        multiAZ,
			errContains:   "changing vswitch_id or vpc_id in place is not supported",
		},
		{
			name:          "multi-AZ network change without configured nodes is refused",
			networkChange: true,
			nodes:         false,
			layout:        multiAZ,
			errContains:   "changing vswitch_id or vpc_id in place is not supported",
		},
		{
			// A pure specification update cannot send the required layout, so it
			// has to fail locally instead of reaching ModifyDBProxyInstance.
			name:          "multi-AZ specification update without configured nodes is refused",
			networkChange: false,
			nodes:         false,
			layout:        multiAZ,
			errContains:   "db_proxy_nodes must be configured when a multi-AZ proxy is modified",
		},
		{
			name:          "multi-AZ specification update with configured nodes is allowed",
			networkChange: false,
			nodes:         true,
			layout:        multiAZ,
		},
		{
			name:          "single-AZ network change without configured nodes is allowed",
			networkChange: true,
			nodes:         false,
			layout:        singleAZ,
		},
		{
			name:   "unknown layout is not treated as multi-AZ",
			layout: nil,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := dbProxyMultiAZGuard(tc.networkChange, tc.nodes, tc.layout)
			if tc.errContains == "" {
				assert.NoError(t, err)
				return
			}
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.errContains)
		})
	}
}

func TestUnitDBProxyRemainingUpdateTimeout(t *testing.T) {
	remaining, err := remainingDBProxyUpdateTimeout(time.Now().Add(2*time.Second), 30*time.Minute)
	assert.NoError(t, err)
	assert.Positive(t, remaining)
	assert.LessOrEqual(t, remaining, 2*time.Second)

	remaining, err = remainingDBProxyUpdateTimeout(time.Now().Add(-time.Second), 30*time.Minute)
	assert.Error(t, err)
	assert.Zero(t, remaining)
	assert.Contains(t, err.Error(), "update timeout was exhausted")
}
