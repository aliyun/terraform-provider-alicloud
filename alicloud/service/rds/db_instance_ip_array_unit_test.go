package rds

import (
	"fmt"
	"testing"
	"time"

	sdkerrors "github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestUnitAliCloudDBInstanceFlattenIpArrays(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{
			"DBInstanceIPArrayName":      "default",
			"DBInstanceIPArrayAttribute": "",
			"SecurityIPList":             "10.0.0.1",
		},
		map[string]interface{}{
			"DBInstanceIPArrayName":      "office",
			"DBInstanceIPArrayAttribute": "",
			"SecurityIPList":             "10.10.10.0/24,10.10.11.0/24",
		},
		map[string]interface{}{
			"DBInstanceIPArrayName":      "partner",
			"DBInstanceIPArrayAttribute": "",
			"SecurityIPList":             "",
		},
		map[string]interface{}{
			"DBInstanceIPArrayName":      "DMS",
			"DBInstanceIPArrayAttribute": "hidden",
			"SecurityIPList":             "100.64.0.0/10",
		},
		map[string]interface{}{
			"DBInstanceIPArrayName":      "template_group_template_suffix",
			"DBInstanceIPArrayAttribute": "template",
			"SecurityIPList":             "10.20.0.0/24",
		},
	}

	got := flattenDBInstanceIpArrays(input, "default")
	assert.Equal(t, 2, len(got))
	assert.Equal(t, "office", got[0]["name"])
	assert.Equal(t, []string{"10.10.10.0/24", "10.10.11.0/24"}, got[0]["security_ips"])
	assert.Equal(t, "partner", got[1]["name"])
	assert.Equal(t, []string{}, got[1]["security_ips"])
	assert.Equal(t, 2, len(got[0]))

	assert.Equal(t, 0, len(flattenDBInstanceIpArrays(nil, "default")))
}

func TestUnitAliCloudDBInstanceFilterIpArraysByNames(t *testing.T) {
	input := []map[string]interface{}{
		{"name": "office", "security_ips": []string{"10.0.0.0/24"}},
		{"name": "external", "security_ips": []string{"10.1.0.0/24"}},
	}

	// No group is tracked yet, so nothing is adopted from the API.
	assert.Equal(t, 0, len(filterDBInstanceIpArraysByNames(input, nil)))
	assert.Equal(t, 0, len(filterDBInstanceIpArraysByNames(input, schema.NewSet(schema.HashString, nil))))

	existing := schema.NewSet(func(v interface{}) int {
		group, _ := v.(map[string]interface{})
		name, _ := group["name"].(string)
		return schema.HashString(name)
	}, []interface{}{
		map[string]interface{}{"name": "office"},
	})

	got := filterDBInstanceIpArraysByNames(input, existing)
	assert.Equal(t, 1, len(got))
	assert.Equal(t, "office", got[0]["name"])
}

func TestUnitAliCloudDBInstanceProtectedIpArrayNames(t *testing.T) {
	assert.True(t, isProtectedDBInstanceIpArrayName("default", "default"))
	assert.True(t, isProtectedDBInstanceIpArrayName("Default", "mygroup"))
	assert.True(t, isProtectedDBInstanceIpArrayName("mygroup", "mygroup"))
	assert.False(t, isProtectedDBInstanceIpArrayName("office", "mygroup"))
}

func TestUnitAliCloudDBInstanceFindIpArray(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{
			"DBInstanceIPArrayName":      "office",
			"DBInstanceIPArrayAttribute": "",
		},
		map[string]interface{}{
			"DBInstanceIPArrayName":      "DMS",
			"DBInstanceIPArrayAttribute": "hidden",
		},
	}

	group, ok := findDBInstanceIpArray(input, "DMS")
	assert.True(t, ok)
	assert.Equal(t, "hidden", group["DBInstanceIPArrayAttribute"])

	_, ok = findDBInstanceIpArray(input, "partner")
	assert.False(t, ok)
}

func TestUnitAliCloudDBInstanceIpArrayGroupLimit(t *testing.T) {
	resource := resourceAliCloudDBInstance()
	block, ok := resource.Schema["db_instance_ip_arrays"]
	assert.True(t, ok)
	assert.Equal(t, DBInstanceIPArrayGroupLimit-1, block.MaxItems)
	assert.Equal(t, 199, block.MaxItems)
}

func TestUnitAliCloudDBInstanceModifySecurityIpsRetryGate(t *testing.T) {
	cases := []struct {
		code string
		want bool
	}{
		{"IncorrectMasterDBInstanceState", true},
		{"IncorrectDBInstanceState", true},
		{"InvalideStatus.Format", true},
		{"InvalidSecurityIps.MalFormed", false},
		{"InvalidGroupName.DuplicatedWithTemplate", false},
		{"InvalidSecurityIPListGroup.QuotaExceeded", false},
	}
	for _, item := range cases {
		err := sdkerrors.NewServerError(400, fmt.Sprintf(`{"Code":%q,"Message":"test"}`, item.code), "")
		assert.Equal(t, item.want, isRetryableModifySecurityIpsError(err), item.code)
	}
	assert.False(t, isRetryableModifySecurityIpsError(nil))
}

type fakeSecurityIpsGroupAccess struct {
	present                bool
	coverErr               error
	deleteErr              error
	absentErr              error
	absentErrAfterCalls    int
	coverRemoves           bool
	deleteRemoves          bool
	coverDelay             time.Duration
	visibleReads           int
	visibleReadsAfterWrite int
	modes                  []string
	ips                    []string
	absentBudgets          []time.Duration
	groupAbsentCalls       int
}

func (f *fakeSecurityIpsGroupAccess) groupAbsent(instanceId, groupName string, timeout time.Duration) (bool, error) {
	f.absentBudgets = append(f.absentBudgets, timeout)
	f.groupAbsentCalls++
	if f.absentErr != nil && (f.absentErrAfterCalls == 0 || f.groupAbsentCalls > f.absentErrAfterCalls) {
		return false, f.absentErr
	}
	if f.visibleReads > 0 {
		f.visibleReads--
		return false, nil
	}
	return !f.present, nil
}

// fastAbsencePolling shrinks the removal poll window and interval so tests do
// not spend the production grace periods in real time.
func fastAbsencePolling(t *testing.T) {
	t.Helper()
	oldInterval := dbSecurityIpsAbsencePollInterval
	oldCoverWindow := dbSecurityIpsCoverAbsenceWindow
	oldDeleteWindow := dbSecurityIpsDeleteAbsenceWindow
	dbSecurityIpsAbsencePollInterval = time.Millisecond
	dbSecurityIpsCoverAbsenceWindow = 30 * time.Millisecond
	dbSecurityIpsDeleteAbsenceWindow = 30 * time.Millisecond
	t.Cleanup(func() {
		dbSecurityIpsAbsencePollInterval = oldInterval
		dbSecurityIpsCoverAbsenceWindow = oldCoverWindow
		dbSecurityIpsDeleteAbsenceWindow = oldDeleteWindow
	})
}

func (f *fakeSecurityIpsGroupAccess) modifyGroup(instanceId, groupName, ips, mode string, timeout time.Duration) error {
	f.modes = append(f.modes, mode)
	f.ips = append(f.ips, ips)
	if mode == "Cover" {
		if f.coverDelay > 0 {
			time.Sleep(f.coverDelay)
		}
		if f.coverErr != nil {
			return f.coverErr
		}
		if f.coverRemoves {
			f.present = false
			if f.visibleReadsAfterWrite > 0 {
				f.visibleReads = f.visibleReadsAfterWrite
			}
		}
		return nil
	}
	if f.deleteErr != nil {
		return f.deleteErr
	}
	if f.deleteRemoves {
		f.present = false
		if f.visibleReadsAfterWrite > 0 {
			f.visibleReads = f.visibleReadsAfterWrite
		}
	}
	return nil
}

func TestUnitAliCloudDBInstanceDeleteSecurityIpsCoverPath(t *testing.T) {
	fastAbsencePolling(t)
	access := &fakeSecurityIpsGroupAccess{present: true, coverRemoves: true}
	err := deleteDBSecurityIpsGroup("rm-test", "office", "10.0.0.1", time.Minute, access)
	assert.NoError(t, err)
	assert.Equal(t, []string{"Cover"}, access.modes)
	assert.Equal(t, []string{""}, access.ips)
}

func TestUnitAliCloudDBInstanceDeleteSecurityIpsDeleteFallback(t *testing.T) {
	fastAbsencePolling(t)
	access := &fakeSecurityIpsGroupAccess{
		present:       true,
		coverErr:      sdkerrors.NewServerError(500, `{"Code":"InternalError","Message":"cover failed"}`, ""),
		deleteRemoves: true,
	}
	err := deleteDBSecurityIpsGroup("rm-test", "office", "10.0.0.1", time.Minute, access)
	assert.NoError(t, err)
	assert.Equal(t, []string{"Cover", "Delete"}, access.modes)
	assert.Equal(t, []string{"", "10.0.0.1"}, access.ips)
}

func TestUnitAliCloudDBInstanceDeleteSecurityIpsMalformedStillPresent(t *testing.T) {
	fastAbsencePolling(t)
	access := &fakeSecurityIpsGroupAccess{
		present:       true,
		coverErr:      sdkerrors.NewServerError(400, `{"Code":"InvalidSecurityIps.MalFormed","Message":"malformed"}`, ""),
		deleteRemoves: true,
	}
	err := deleteDBSecurityIpsGroup("rm-test", "office", "10.0.0.1", time.Minute, access)
	assert.NoError(t, err)
	assert.Equal(t, []string{"Cover", "Delete"}, access.modes)
}

func TestUnitAliCloudDBInstanceDeleteSecurityIpsStillExists(t *testing.T) {
	fastAbsencePolling(t)
	access := &fakeSecurityIpsGroupAccess{present: true}
	err := deleteDBSecurityIpsGroup("rm-test", "office", "10.0.0.1", time.Minute, access)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "still exists after it was removed")
}

func TestUnitAliCloudDBInstanceDeleteSecurityIpsDeleteError(t *testing.T) {
	access := &fakeSecurityIpsGroupAccess{
		present:   true,
		coverErr:  sdkerrors.NewServerError(403, `{"Code":"IncorrectDBInstanceState","Message":"state"}`, ""),
		deleteErr: sdkerrors.NewServerError(400, `{"Code":"InvalidParameter","Message":"bad request"}`, ""),
	}
	err := deleteDBSecurityIpsGroup("rm-test", "office", "10.0.0.1", time.Minute, access)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

func TestUnitAliCloudDBInstanceDeleteSecurityIpsBudgetExhausted(t *testing.T) {
	access := &fakeSecurityIpsGroupAccess{present: true}
	err := deleteDBSecurityIpsGroup("rm-test", "office", "10.0.0.1", 0, access)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "budget is exhausted")
	assert.Empty(t, access.modes)
	assert.Empty(t, access.absentBudgets)
}

func TestUnitAliCloudDBInstanceDeleteSecurityIpsBudgetPassedDown(t *testing.T) {
	fastAbsencePolling(t)
	access := &fakeSecurityIpsGroupAccess{present: true, coverRemoves: false, coverDelay: 1200 * time.Millisecond}
	err := deleteDBSecurityIpsGroup("rm-test", "office", "10.0.0.1", 3*time.Second, access)
	assert.Error(t, err)
	assert.Equal(t, []string{"Cover", "Delete"}, access.modes)
	assert.NotEmpty(t, access.absentBudgets)
	first := access.absentBudgets[0]
	last := access.absentBudgets[len(access.absentBudgets)-1]
	assert.True(t, first <= 3*time.Second && first > 0)
	assert.True(t, last < first, "every read-back must receive the remaining budget, got %s then %s", first, last)
	for i := 1; i < len(access.absentBudgets); i++ {
		assert.LessOrEqual(t, access.absentBudgets[i], access.absentBudgets[i-1], "the budget must never be rebuilt between reads")
	}
}

func TestUnitAliCloudDBInstanceDeleteSecurityIpsNoFallbackAfterBudgetExhausted(t *testing.T) {
	access := &fakeSecurityIpsGroupAccess{
		present:    true,
		coverDelay: 400 * time.Millisecond,
		coverErr:   sdkerrors.NewServerError(500, `{"Code":"InternalError","Message":"cover failed"}`, ""),
	}
	err := deleteDBSecurityIpsGroup("rm-test", "office", "10.0.0.1", 200*time.Millisecond, access)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "budget is exhausted")
	assert.Equal(t, []string{"Cover"}, access.modes)
}

func TestUnitAliCloudDBInstanceDeleteSecurityIpsDelayedVisibility(t *testing.T) {
	fastAbsencePolling(t)
	access := &fakeSecurityIpsGroupAccess{present: true, coverRemoves: true, visibleReads: 3}
	err := deleteDBSecurityIpsGroup("rm-test", "office", "10.0.0.1", time.Minute, access)
	assert.NoError(t, err)
	assert.Equal(t, []string{"Cover"}, access.modes)
	assert.GreaterOrEqual(t, access.groupAbsentCalls, 4, "a group that stays visible briefly must be polled again")
}

func TestUnitAliCloudDBInstanceDeleteSecurityIpsAbsenceErrorPropagates(t *testing.T) {
	fastAbsencePolling(t)
	access := &fakeSecurityIpsGroupAccess{
		present:   true,
		absentErr: sdkerrors.NewServerError(500, `{"Code":"InternalError","Message":"read failed"}`, ""),
	}
	err := deleteDBSecurityIpsGroup("rm-test", "office", "10.0.0.1", time.Minute, access)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "read failed")
	assert.Empty(t, access.modes, "a failed read must stop before any write")
}

// A group removed by the Delete fallback can stay visible for a short while.
// The fallback must poll after the write instead of deciding on one read.
func TestUnitAliCloudDBInstanceDeleteSecurityIpsDeleteFallbackDelayedVisibility(t *testing.T) {
	fastAbsencePolling(t)
	access := &fakeSecurityIpsGroupAccess{
		present:                true,
		coverErr:               sdkerrors.NewServerError(500, `{"Code":"InternalError","Message":"cover failed"}`, ""),
		deleteRemoves:          true,
		visibleReadsAfterWrite: 3,
	}
	err := deleteDBSecurityIpsGroup("rm-test", "office", "10.0.0.1", time.Minute, access)
	assert.NoError(t, err)
	assert.Equal(t, []string{"Cover", "Delete"}, access.modes)
	assert.GreaterOrEqual(t, access.groupAbsentCalls, 5, "the Delete fallback must poll after the write until the group disappears")
	for i := 1; i < len(access.absentBudgets); i++ {
		assert.LessOrEqual(t, access.absentBudgets[i], access.absentBudgets[i-1], "the budget must never be rebuilt between reads")
	}
}

// A read error that only appears while polling after the Cover write must
// reach the caller. The old single-read implementation could not observe it.
func TestUnitAliCloudDBInstanceDeleteSecurityIpsPostCoverPollReadErrorPropagates(t *testing.T) {
	fastAbsencePolling(t)
	access := &fakeSecurityIpsGroupAccess{
		present:                true,
		coverRemoves:           true,
		visibleReadsAfterWrite: 2,
		absentErr:              sdkerrors.NewServerError(500, `{"Code":"InternalError","Message":"read failed"}`, ""),
		absentErrAfterCalls:    3,
	}
	err := deleteDBSecurityIpsGroup("rm-test", "office", "10.0.0.1", time.Minute, access)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "read failed")
	assert.Equal(t, []string{"Cover"}, access.modes, "the write happens before the failing poll read")
	assert.GreaterOrEqual(t, access.groupAbsentCalls, 4, "the poll must read again before the error is surfaced")
}

// A read error that only appears while polling after the Delete fallback write
// must reach the caller instead of being reported as the group still existing.
func TestUnitAliCloudDBInstanceDeleteSecurityIpsPostDeletePollReadErrorPropagates(t *testing.T) {
	fastAbsencePolling(t)
	access := &fakeSecurityIpsGroupAccess{
		present:                true,
		coverErr:               sdkerrors.NewServerError(500, `{"Code":"InternalError","Message":"cover failed"}`, ""),
		deleteRemoves:          true,
		visibleReadsAfterWrite: 1,
		absentErr:              sdkerrors.NewServerError(500, `{"Code":"InternalError","Message":"read failed"}`, ""),
		absentErrAfterCalls:    2,
	}
	err := deleteDBSecurityIpsGroup("rm-test", "office", "10.0.0.1", time.Minute, access)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "read failed")
	assert.Equal(t, []string{"Cover", "Delete"}, access.modes, "the fallback write happens before the failing poll read")
	assert.GreaterOrEqual(t, access.groupAbsentCalls, 3, "the poll must read again before the error is surfaced")
}

func TestUnitAliCloudDBInstanceLifecycleBudget(t *testing.T) {
	left, ok := remainingLifecycleBudget(time.Now().Add(time.Minute))
	assert.True(t, ok)
	assert.True(t, left > 0 && left <= time.Minute)

	left, ok = remainingLifecycleBudget(time.Now().Add(-time.Second))
	assert.False(t, ok)
	assert.True(t, left <= 0)
}

func TestUnitAliCloudDBInstanceStateRefreshWithinDeadlineExhausted(t *testing.T) {
	service := &RdsService{}
	refresh := service.RdsDBInstanceStateRefreshFuncWithinDeadline("rm-test", []string{"Deleting"}, time.Now().Add(-time.Second))
	_, _, err := refresh()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "budget is exhausted")
}

type fakeSecurityIpsReader struct {
	called bool
	got    time.Duration
}

func (f *fakeSecurityIpsReader) DescribeDBSecurityIpsWithTimeout(instanceId string, timeout time.Duration) ([]interface{}, error) {
	f.called = true
	f.got = timeout
	return []interface{}{}, nil
}

func TestUnitAliCloudDBInstanceEntryReadUsesRemainingBudget(t *testing.T) {
	reader := &fakeSecurityIpsReader{}
	_, err := loadDBInstanceIpArraysWithinBudget(reader, "rm-test", 2*time.Second)
	assert.NoError(t, err)
	assert.True(t, reader.called)
	assert.Equal(t, 2*time.Second, reader.got)
}

func TestUnitAliCloudDBInstanceEntryReadBudgetExhausted(t *testing.T) {
	reader := &fakeSecurityIpsReader{}
	_, err := loadDBInstanceIpArraysWithinBudget(reader, "rm-test", 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "budget is exhausted")
	assert.False(t, reader.called, "an exhausted budget must abort before any read or write")
}

func TestUnitAliCloudDBInstanceTemplateGroupNames(t *testing.T) {
	ipArrays := []interface{}{
		map[string]interface{}{"DBInstanceIPArrayName": "office", "DBInstanceIPArrayAttribute": ""},
		map[string]interface{}{"DBInstanceIPArrayName": "tmpl_template_suffix", "DBInstanceIPArrayAttribute": "template"},
		map[string]interface{}{"DBInstanceIPArrayName": "", "DBInstanceIPArrayAttribute": "template"},
		"not-a-map",
	}
	assert.Equal(t, map[string]bool{"tmpl_template_suffix": true}, templateDBInstanceIpArrayGroupNames(ipArrays))
	assert.Empty(t, templateDBInstanceIpArrayGroupNames(nil))
}

func TestUnitAliCloudDBInstanceTemplateGroupConflict(t *testing.T) {
	protected := map[string]bool{"tmpl_template_suffix": true}

	name, conflict := findDBInstanceIpArrayTemplateConflict([]string{"office", "tmpl_template_suffix"}, protected)
	assert.True(t, conflict)
	assert.Equal(t, "tmpl_template_suffix", name)

	name, conflict = findDBInstanceIpArrayTemplateConflict([]string{"office", "partner"}, protected)
	assert.False(t, conflict)
	assert.Empty(t, name)
}
