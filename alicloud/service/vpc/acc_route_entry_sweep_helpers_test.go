package vpc_test

// Sweep helpers for the migrated alicloud_route_entry resource. The baseline
// root package alicloud hosted these as package-internal symbols: the
// unexported var routeEntryTaskLock in resource_alicloud_route_entry.go and
// the unexported func rdsErrorHasCode (plus its rdsRetryableQueryError
// unwrap type) in errors.go / service_alicloud_rds.go. They were reachable
// from the baseline route_entry test because that test lived in package
// alicloud. The migration moves the resource into package vpc
// (alicloud/service/vpc) and the test into the external package vpc_test,
// which can see neither vpc's unexported routeEntryTaskLock nor the root's
// unexported rdsErrorHasCode. rdsErrorHasCode (plus its rdsRetryableQueryError
// unwrap type) is re-hosted here as a vpc_test-local helper so the migrated
// sweep logic (sweepVpcRouteTableEntries) stays byte-identical.
//
// routeEntryTaskLock is NOT re-hosted as a fresh mutex: it shares the
// provider's real CreateRouteEntry/DeleteRouteEntry serialization mutex via
// the export_test.go bridge RouteEntryTaskLockForTest (package vpc), so the
// test's .Lock()/.Unlock() at resource_alicloud_route_entry_test.go queue on
// the same mutex the provider's Create/Delete call sites hold — restoring the
// baseline fidelity (where the root test and production shared one package var)
// instead of a disconnected lock that only serializes the test process. The
// pointer assignment (var = vpc.RouteEntryTaskLockForTest) binds the shared
// *sync.Mutex; the production var and its call sites are unchanged.
// ComplexError and ProviderError resolve via aliases_test.go. The rds-retry
// unwrap branch is unreachable for VPC RPC errors (client.RpcPost wraps
// *errors.ServerError, never rdsRetryableQueryError) but is retained for
// fidelity.

import (
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	vpc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/vpc"
	"github.com/denverdino/aliyungo/common"
)

var routeEntryTaskLock = vpc.RouteEntryTaskLockForTest

type rdsRetryableQueryError struct {
	cause error
}

func (e *rdsRetryableQueryError) Error() string { return e.cause.Error() }
func (e *rdsRetryableQueryError) Unwrap() error { return e.cause }

// rdsErrorHasCode matches structured error codes and unwraps query retry
// metadata. Transcribed verbatim from the root alicloud/errors.go original.
func rdsErrorHasCode(err error, codes ...string) bool {
	for {
		if e, ok := err.(*ComplexError); ok {
			err = e.Cause
			continue
		}
		if e, ok := err.(*rdsRetryableQueryError); ok {
			err = e.cause
			continue
		}
		break
	}
	var code string
	switch e := err.(type) {
	case *errors.ServerError:
		code = e.ErrorCode()
	case *tea.SDKError:
		code = tea.StringValue(e.Code)
	case *ProviderError:
		code = e.ErrorCode()
	case *common.Error:
		code = e.Code
	default:
		return false
	}
	for _, expected := range codes {
		if code == expected {
			return true
		}
	}
	return false
}
