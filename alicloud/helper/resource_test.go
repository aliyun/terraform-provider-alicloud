package helper

import (
	"fmt"
	"testing"

	sdkerrors "github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
)

func TestUnitIsExpectedErrors(t *testing.T) {
	codeErr := error(sdkerrors.NewServerError(400, `{"Code":"InvalidSnapshotLock.NotFound","Message":"The specified snapshot lock is not found."}`, ""))
	if !IsExpectedErrors(codeErr, []string{"InvalidSnapshotLock.NotFound"}) {
		t.Fatalf("expected ServerError with matching ErrorCode to be expected, got: %#v", codeErr)
	}
	if IsExpectedErrors(codeErr, []string{"InvalidSnapshotId.NotFound"}) {
		t.Fatalf("expected unrelated ErrorCode not to match, got: %#v", codeErr)
	}

	messageErr := error(sdkerrors.NewServerError(400, `{"Code":"Mock.Code","Message":"operation failed: InvalidSnapshotId.NotFound for mock"}`, ""))
	if !IsExpectedErrors(messageErr, []string{"InvalidSnapshotId.NotFound"}) {
		t.Fatalf("expected ServerError whose Message contains the code to be expected, got: %#v", messageErr)
	}

	wrapped := &errs.ComplexError{Cause: codeErr, Err: fmt.Errorf("Resource alicloud_ecs_snapshot_lock LockSnapshot Failed!!! [SDK alibaba-cloud-sdk-go ERROR]"), Path: "mock.go", Line: 1}
	if !IsExpectedErrors(wrapped, []string{"InvalidSnapshotLock.NotFound"}) {
		t.Fatalf("expected ComplexError to unwrap to its Cause for matching, got: %#v", wrapped)
	}

	if !IsExpectedErrors(fmt.Errorf("request failed with InvalidSnapshotLock.NotFound"), []string{"InvalidSnapshotLock.NotFound"}) {
		t.Fatal("expected plain error whose message contains the code to be expected")
	}
	if IsExpectedErrors(fmt.Errorf("boom"), []string{"InvalidSnapshotLock.NotFound"}) {
		t.Fatal("expected plain error without the code not to match")
	}
	if IsExpectedErrors(nil, []string{"InvalidSnapshotLock.NotFound"}) {
		t.Fatal("expected nil error not to match")
	}
}

func TestUnitExpandSingletonToList(t *testing.T) {
	got := ExpandSingletonToList("snap-abc")
	if len(got) != 1 || got[0] != "snap-abc" {
		t.Fatalf("expected []interface{}{\"snap-abc\"}, got %#v", got)
	}
}

func TestUnitIntBetween(t *testing.T) {
	if _, es := IntBetween(0, 72)(48, "cool_off_period"); len(es) > 0 {
		t.Fatalf("expected in-range value to validate, got %v", es)
	}
	for _, v := range []int{-1, 73} {
		if _, es := IntBetween(0, 72)(v, "cool_off_period"); len(es) == 0 {
			t.Fatalf("expected out-of-range value %d to fail validation", v)
		}
	}
	if _, es := IntBetween(0, 72)("48", "cool_off_period"); len(es) == 0 {
		t.Fatal("expected non-int value to fail validation")
	}
}

func TestUnitStringInSlice(t *testing.T) {
	if _, es := StringInSlice([]string{"compliance"}, false)("compliance", "lock_mode"); len(es) > 0 {
		t.Fatalf("expected exact match to validate, got %v", es)
	}
	if _, es := StringInSlice([]string{"compliance"}, false)("Compliance", "lock_mode"); len(es) == 0 {
		t.Fatal("expected case mismatch to fail validation when ignoreCase is false")
	}
	if _, es := StringInSlice([]string{"compliance"}, true)("Compliance", "lock_mode"); len(es) > 0 {
		t.Fatalf("expected case-insensitive match to validate, got %v", es)
	}
	if _, es := StringInSlice([]string{"compliance"}, false)("ignored", "lock_mode"); len(es) == 0 {
		t.Fatal("expected unmatched value to fail validation")
	}
	if _, es := StringInSlice([]string{"compliance"}, false)(42, "lock_mode"); len(es) == 0 {
		t.Fatal("expected non-string value to fail validation")
	}
}
