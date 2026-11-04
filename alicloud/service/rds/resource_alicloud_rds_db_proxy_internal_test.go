package rds

import (
	"testing"
	"time"

	"github.com/alibabacloud-go/tea/tea"
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

func TestUnitDBProxyModifyRetryError(t *testing.T) {
	newSDKError := func(code string, statusCode int) error {
		return tea.NewSDKError(map[string]interface{}{
			"code":       code,
			"message":    code,
			"statusCode": statusCode,
			"data":       map[string]interface{}{"Code": code},
		})
	}

	for _, tc := range []struct {
		name      string
		err       error
		shutdown  bool
		retryable bool
		success   bool
	}{
		{
			// The proxy reports MaxscaleInternalError (HTTP 400) while it is
			// still changing, so both actions have to retry it.
			name:      "startup retries MaxscaleInternalError",
			err:       newSDKError("MaxscaleInternalError", 400),
			retryable: true,
		},
		{
			name:      "shutdown retries MaxscaleInternalError",
			err:       newSDKError("MaxscaleInternalError", 400),
			shutdown:  true,
			retryable: true,
		},
		{
			name:      "startup retries MaxscaleCreating",
			err:       newSDKError("MaxscaleCreating", 404),
			retryable: true,
		},
		{
			// MaxscaleCreating is HTTP 404. Shutdown must retry it instead of
			// reading the status as a proxy that is already gone.
			name:      "shutdown retries MaxscaleCreating instead of treating it as deleted",
			err:       newSDKError("MaxscaleCreating", 404),
			shutdown:  true,
			retryable: true,
		},
		{
			name:     "shutdown keeps treating a missing proxy as deleted",
			err:      newSDKError("InvalidDBInstanceId.NotFound", 404),
			shutdown: true,
			success:  true,
		},
		{
			name: "startup does not treat a missing proxy as success",
			err:  newSDKError("InvalidDBInstanceId.NotFound", 404),
		},
		{
			name:      "shutdown keeps retrying OperationDenied.DBStatus",
			err:       newSDKError("OperationDenied.DBStatus", 400),
			shutdown:  true,
			retryable: true,
		},
		{
			name: "startup does not retry OperationDenied.DBStatus",
			err:  newSDKError("OperationDenied.DBStatus", 400),
		},
		{
			name: "terminal error fails fast",
			err:  newSDKError("InvalidDBProxyStatus", 400),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			retryErr := dbProxyModifyRetryError(tc.err, tc.shutdown)
			if tc.success {
				assert.Nil(t, retryErr)
				return
			}
			assert.NotNil(t, retryErr)
			assert.Equal(t, tc.retryable, retryErr.Retryable)
			assert.ErrorIs(t, retryErr.Err, tc.err)
		})
	}
}

func TestUnitDBProxyUpdateBudgetSharedDeadline(t *testing.T) {
	total := 30 * time.Minute
	now := time.Date(2026, 11, 4, 12, 0, 0, 0, time.UTC)
	budget := newDBProxyUpdateBudget(func() time.Time { return now }, total)

	var timeouts []time.Duration
	consume := func(elapsed time.Duration) func(timeout time.Duration) error {
		return func(timeout time.Duration) error {
			timeouts = append(timeouts, timeout)
			now = now.Add(elapsed)
			return nil
		}
	}

	// Two sequential steps of one update: the first spends five minutes, so the
	// second must receive the remaining 25 minutes of the same deadline instead
	// of a freshly restarted full update timeout.
	assert.NoError(t, runDBProxyUpdateSteps(budget, consume(5*time.Minute), consume(0)))
	assert.Equal(t, []time.Duration{total, total - 5*time.Minute}, timeouts)
	assert.NotEqual(t, total, timeouts[1])

	// Once the shared budget is exhausted, the following write or wait must not
	// start at all.
	now = now.Add(total)
	started := false
	err := runDBProxyUpdateSteps(budget, func(time.Duration) error {
		started = true
		return nil
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update timeout was exhausted")
	assert.False(t, started)
}
