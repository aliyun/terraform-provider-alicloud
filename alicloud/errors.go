package alicloud

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/fc-go-sdk"
	"github.com/denverdino/aliyungo/common"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
)

// The helpers below moved to the leaf packages errs and logs; these
// aliases keep the flat layout compiling unchanged.
//
// WrapError/WrapErrorf must stay var aliases rather than forwarding
// functions: they read runtime.Caller(1) to attribute errors to the
// caller, and a wrapper would insert a frame and misattribute every
// error in the provider. The same holds for addDebug in common.go.
//
// ComplexError and ProviderError must be type aliases — the "=" is
// load-bearing: the type assertions in IsExpectedErrors and
// rdsErrorHasCode below must keep matching values wrapped by errs,
// or 404 handling and error-code exemptions silently stop working.
type ComplexError = errs.ComplexError
type ProviderError = errs.ProviderError
type ErrorSource = errs.ErrorSource

var (
	Error                      = errs.Error
	WrapError                  = errs.WrapError
	WrapErrorf                 = errs.WrapErrorf
	WrapComplexError           = errs.WrapComplexError
	NotFoundError              = errs.NotFoundError
	NotFoundErr                = errs.NotFoundErr
	GetNotFoundErrorFromString = errs.GetNotFoundErrorFromString
	GetTimeErrorFromString     = errs.GetTimeErrorFromString
	NewProviderError           = errs.NewProviderError
	NeedRetry                  = helper.NeedRetry
)

const (
	AlibabaCloudSdkGoERROR = errs.AlibabaCloudSdkGoERROR
	AliyunLogGoSdkERROR    = errs.AliyunLogGoSdkERROR
	AliyunDatahubSdkGo     = errs.AliyunDatahubSdkGo
	AliyunOssGoSdk         = errs.AliyunOssGoSdk
	FcGoSdk                = errs.FcGoSdk
	DenverdinoAliyungo     = errs.DenverdinoAliyungo
	AliyunTablestoreGoSdk  = errs.AliyunTablestoreGoSdk
	AliMnsERROR            = errs.AliMnsERROR
	ProviderERROR          = errs.ProviderERROR
	DefaultErrorMsg        = errs.DefaultErrorMsg
	DataDefaultErrorMsg    = errs.DataDefaultErrorMsg
	NotFoundWithResponse   = errs.NotFoundWithResponse
	FailedGetAttributeMsg  = errs.FailedGetAttributeMsg
	DefaultDebugMsg        = logs.DefaultDebugMsg
)

const (
	// common
	NotFound                = "NotFound"
	ResourceNotfound        = "ResourceNotfound"
	ServiceUnavailable      = "ServiceUnavailable"
	InstanceNotFound        = "Instance.Notfound"
	ForbiddenInstance       = "Forbidden.InstanceNotFound"
	MessageInstanceNotFound = "instance is not found"
	Throttling              = "Throttling"
	ThrottlingUser          = "Throttling.User"

	// RAM Instance Not Found
	RamInstanceNotFound   = "Forbidden.InstanceNotFound"
	AliyunGoClientFailure = "AliyunGoClientFailure"

	LogClientTimeout = "Client.Timeout exceeded while awaiting headers"

	InvalidFileSystemStatus_Ordering = "InvalidFileSystemStatus.Ordering"
	NotFoundArticle                  = "not found article by given param"
)

var SlbIsBusy = []string{"SystemBusy", "OperationBusy", "ServiceIsStopping", "BackendServer.configuring", "ServiceIsConfiguring"}
var EcsNotFound = []string{"InvalidInstanceId.NotFound", "Forbidden.InstanceNotFound"}
var DiskInvalidOperation = []string{"IncorrectDiskStatus", "IncorrectInstanceStatus", "OperationConflict", "InternalError", "InvalidOperation.Conflict", "IncorrectDiskStatus.Initializing"}
var NetworkInterfaceInvalidOperations = []string{"InvalidOperation.InvalidEniState", "InvalidOperation.InvalidEcsState", "InvalidOperation.InvalidEniType", "InvalidOperation.HasMemberEniAttached", "OperationConflict", "ServiceUnavailable", "InternalError"}
var OperationDeniedDBStatus = []string{"LockTimeout", "InstanceConnectTimeoutFault", "ConcurrentTaskExceeded", "OperationDenied.DBStatus", "Database.ConnectError", "OperationDenied.DBInstanceStatus", "OperationDenied.DBClusterStatus", "InternalError", "OperationDenied.OutofUsage", "IncorrectDBInstanceState"}
var DBReadInstanceNotReadyStatus = []string{"OperationDenied.ReadDBInstanceStatus", "OperationDenied.MasterDBInstanceState", "ReadDBInstance.Mismatch"}
var NasNotFound = []string{"InvalidMountTarget.NotFound", "InvalidFileSystem.NotFound", "Forbidden.NasNotFound", "InvalidLBid.NotFound", "VolumeUnavailable"}
var SnapshotInvalidOperations = []string{"OperationConflict", "ServiceUnavailable", "InternalError", "SnapshotCreatedDisk", "SnapshotCreatedImage"}
var DiskNotSupportOnlineChangeErrors = []string{"InvalidDiskCategory.NotSupported", "InvalidRegion.NotSupport", "IncorrectInstanceStatus", "IncorrectDiskStatus", "InvalidOperation.InstanceTypeNotSupport"}
var DBInstanceTDEErrors = []string{"InvaildEngineInRegion.ValueNotSupported", "InstanceEngineType.NotSupport", "OperationDenied.DBInstanceType", "IncorrectDBInstanceType", "IncorrectEngineVersion", "DBSizeExceeded", "InvalidDBName.NotFound", "DbossGeneralError"}

// details at: https://help.aliyun.com/document_detail/27300.html
var OtsTableIsTemporarilyUnavailable = []string{"no such host", "OTSServerBusy", "OTSPartitionUnavailable", "OTSInternalServerError",
	"OTSTimeout", "OTSServerUnavailable", "OTSRowOperationConflict", "OTSTableNotReady", "OTSNotEnoughCapacityUnit", "Too frequent table operations."}

var OtsTunnelIsTemporarilyUnavailable = []string{"no such host", "OTSTunnelServerUnavailable", "OTSPermissionDenied"}
var OtsSecondaryIndexIsTemporarilyUnavailable = []string{"no such host", "OTSServerUnavailable"}
var OtsSearchIndexIsTemporarilyUnavailable = []string{"no such host", "OTSServerUnavailable"}

// rdsErrorHasCode matches structured RDS error codes and unwraps query retry metadata.
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

func IsExpectedErrors(err error, expectCodes []string) bool {
	if err == nil {
		return false
	}

	if e, ok := err.(*ComplexError); ok {
		return IsExpectedErrors(e.Cause, expectCodes)
	}

	if e, ok := err.(*tea.SDKError); ok {
		// Use tea.StringValue to avoid nil-pointer panics when Code or Data is
		// nil (e.g. a backend validation error with Code set but Data unset).
		// Guard HasPrefix against an empty codeVal: strings.HasPrefix(code, "")
		// is always true and would falsely match every expected code.
		codeVal := tea.StringValue(e.Code)
		dataVal := tea.StringValue(e.Data)
		for _, code := range expectCodes {
			if codeVal == code || (codeVal != "" && strings.HasPrefix(code, codeVal)) {
				return true
			}
			if dataVal != "" && strings.Contains(dataVal, code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(*errors.ServerError); ok {
		for _, code := range expectCodes {
			if e.ErrorCode() == code || strings.Contains(e.Message(), code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(*ProviderError); ok {
		for _, code := range expectCodes {
			if e.ErrorCode() == code || strings.Contains(e.Message(), code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(*common.Error); ok {
		for _, code := range expectCodes {
			if e.Code == code || fmt.Sprint(e.StatusCode) == code || strings.Contains(e.Message, code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(*sls.Error); ok {
		for _, code := range expectCodes {
			if e.Code == code || strings.Contains(e.Message, code) || strings.Contains(e.String(), code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(oss.ServiceError); ok {
		for _, code := range expectCodes {
			if e.Code == code || strings.Contains(e.Message, code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(*fc.ServiceError); ok {
		for _, code := range expectCodes {
			if e.ErrorCode == code || strings.Contains(e.ErrorMessage, code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(*datahub.DatahubClientError); ok {
		for _, code := range expectCodes {
			if e.Code == code || strings.Contains(e.Message, code) {
				return true
			}
		}
		return false
	}

	for _, code := range expectCodes {
		if strings.Contains(err.Error(), code) {
			return true
		}
	}
	return false
}

func NoCodeRegexRetry(err error) bool {
	if err == nil {
		return false
	}

	postRegex := regexp.MustCompile("^Post [\"]*https://.*")
	if postRegex.MatchString(err.Error()) {
		return true
	}

	if e, ok := err.(*errors.ClientError); ok && e.ErrorCode() == errors.TimeoutErrorCode {
		return true
	}

	throttlingRegex := regexp.MustCompile("Throttling")

	if e, ok := err.(*tea.SDKError); ok {
		if strings.Contains(*e.Message, "code: 500, 您已开通过") {
			return false
		}

		if strings.Contains(*e.Message, "Client.Timeout") {
			return true
		}

		if *e.Code == ServiceUnavailable || *e.Code == "Rejected.Throttling" || throttlingRegex.MatchString(*e.Code) {
			return true
		}
	}

	if e, ok := err.(*errors.ServerError); ok {
		return e.ErrorCode() == ServiceUnavailable || e.ErrorCode() == "Rejected.Throttling" || throttlingRegex.MatchString(e.ErrorCode())
	}

	if e, ok := err.(*common.Error); ok {
		return e.Code == ServiceUnavailable || e.Code == "Rejected.Throttling" || throttlingRegex.MatchString(e.Code)
	}

	return false
}

func IsExpectedErrorCodes(code string, errorCodes []string) bool {
	if code == "" {
		return false
	}
	for _, v := range errorCodes {
		if v == code {
			return true
		}
	}
	return false
}

func GetNotFoundMessage(product, id string) string {
	return fmt.Sprintf("The specified %s %s is not found.", product, id)
}

func GetTimeoutMessage(product, status string) string {
	return fmt.Sprintf("Waiting for %s %s is timeout.", product, status)
}

func GetCreateFailedMessage(product string) string {
	return fmt.Sprintf("The specified %s is create failed.", product)
}

// A default message of ComplexError's Err. It is format to Resource <resource-id> <operation> Failed!!! <error source>
const ResponseCodeMsg = "Resource %s %s Failed!!! %v"
const RequestIdMsg = "RequestId: %s"
const NotFoundMsg = ResourceNotfound + "!!! %s"
const NotFoundWithError = ResourceNotfound + "!!! Error: %v"
const DefaultTimeoutMsg = "Resource %s %s Timeout!!! %s"
const DeleteTimeoutMsg = "Resource %s Still Exists. %s Timeout!!! %s"
const WaitTimeoutMsg = "Resource %s %s Timeout In %d Seconds. Got: %s Expected: %s !!! %s"
const SweepDefaultErrorMsg = "Sweep %s %s Failed!!!"
const IdMsg = "Resource id：%s "
const FailedToReachTargetStatus = "Failed to reach target status. Last status: %s."
const FailedToReachTargetStatusWithResponse = "Resource %s failed to reach target status. Last response: %s"
const FailedToReachTargetStatusWithError = "Resource %s failed to reach target status. Last error: %s"

const FailedToReachTargetStatusWithRequestId = FailedToReachTargetStatus + " Last RequestId: %s."
const FailedToReachTargetAttribute = "Failed to reach value for target attribute. Current value is %s."
const RequiredWhenMsg = "attribute '%s' is required when '%s' is %v"
const InvalidAttributeValue = "attribute '%s' must be set to the value '%s'."
