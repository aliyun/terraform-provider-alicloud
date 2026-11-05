package rds

import (
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/denverdino/aliyungo/common"
)

// rdsErrorHasCode matches structured RDS error codes and unwraps query retry
// metadata. It is the package-rds copy of the root helper; rdsRetryableQueryError
// (defined in service.go) moves with the RDS service, and this function keeps
// resolving it inside the split package. The root package keeps its own copy
// for non-RDS tests (vpc/route) that still reference the root helper.
//
// ComplexError and ProviderError are type aliases to errs.* (see aliases.go);
// the "=" is load-bearing: the type assertions here must keep matching values
// wrapped by errs, or error-code exemptions silently stop working.
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
