package eas_test

import (
	"errors"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/eas"
)

// TestUnitEasIsResourceNotFound exercises the EAS-08 not-found classification
// used by Read/Update/Delete to decide whether a DescribeResource miss means the
// resource is gone (remove from state / treat as deleted) versus a transient
// error that must surface as a failure. A misclassification in either direction
// is harmful: a transient error read as not-found drops a still-existing resource
// from state, while a genuine not-found read as a hard error leaves dangling
// state. The table covers the raw SDK error shapes the Roa call returns — a
// *tea.SDKError and a *errors.ServerError — plus the *errs.ComplexError wrapper
// that errs.WrapErrorf adds, recursing to the original cause.
func TestUnitEasIsResourceNotFound(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil error is not not-found", nil, false},
		{"sdk 404 status is not-found", &tea.SDKError{StatusCode: tea.Int(404)}, true},
		{"sdk message matching NotFound is not-found", &tea.SDKError{Message: tea.String("Resource NotFound: the group does not exist")}, true},
		{"sdk NotFound code is not-found via code substring", &tea.SDKError{Code: tea.String("InvalidResource.NotFound")}, true},
		{"sdk throttling is not not-found", &tea.SDKError{Code: tea.String("Throttling")}, false},
		{"sdk 403 is not not-found", &tea.SDKError{StatusCode: tea.Int(403)}, false},
		{"complex wrapping a not-found code is not-found", &errs.ComplexError{Cause: &tea.SDKError{Code: tea.String("InvalidResource.NotFound")}}, true},
		{"complex wrapping throttling is not not-found", &errs.ComplexError{Cause: &tea.SDKError{Code: tea.String("Throttling")}}, false},
		{"complex wrapping a 404 sdk error is not-found", &errs.ComplexError{Cause: &tea.SDKError{StatusCode: tea.Int(404)}}, true},
		{"plain go error is not not-found", errors.New("something went wrong"), false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := eas.IsResourceNotFound(tc.err); got != tc.want {
				t.Fatalf("IsResourceNotFound(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}
