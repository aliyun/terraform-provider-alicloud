// Package errs holds the provider's error helpers and message formats.
// WrapError/WrapErrorf/ComplexError live here as the implementation of
// record; package alicloud re-exports them as var aliases so the flat
// layout keeps compiling, and service packages import this package
// directly.
package errs

import (
	"fmt"
	"log"
	"regexp"
	"runtime"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	otsTunnel "github.com/aliyun/aliyun-tablestore-go-sdk/tunnel"
	"github.com/aliyun/fc-go-sdk"
	"github.com/denverdino/aliyungo/common"
)

const ResourceNotfound = "ResourceNotfound"
const RequestIdMsg = "RequestId: %s"
const NotFoundMsg = ResourceNotfound + "!!! %s"

const NotFound = "NotFound"
const InstanceNotFound = "Instance.Notfound"
const RamInstanceNotFound = "Forbidden.InstanceNotFound"
const MessageInstanceNotFound = "instance is not found"

type ComplexError struct {
	Cause error
	Err   error
	Path  string
	Line  int
}

func (e ComplexError) Error() string {
	if e.Cause == nil {
		e.Cause = Error("<nil cause>")
	}
	if e.Err == nil {
		return fmt.Sprintf("\u001B[31m[ERROR]\u001B[0m %s:%d:\n%s", e.Path, e.Line, e.Cause.Error())
	}
	return fmt.Sprintf("\u001B[31m[ERROR]\u001B[0m %s:%d: %s:\n%s", e.Path, e.Line, e.Err.Error(), e.Cause.Error())
}

func Error(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

// Return a ComplexError which including error occurred file and path
func WrapError(cause error) error {
	if cause == nil {
		return nil
	}
	_, filepath, line, ok := runtime.Caller(1)
	if !ok {
		log.Printf("\u001B[31m[ERROR]\u001B[0m runtime.Caller error in WrapError.")
		return WrapComplexError(cause, nil, "", -1)
	}
	parts := strings.Split(filepath, "/")
	if len(parts) > 3 {
		filepath = strings.Join(parts[len(parts)-3:], "/")
	}
	return WrapComplexError(cause, nil, filepath, line)
}

// Return a ComplexError which including extra error message, error occurred file and path
func WrapErrorf(cause error, msg string, args ...interface{}) error {
	if cause == nil && strings.TrimSpace(msg) == "" {
		return nil
	}
	_, filepath, line, ok := runtime.Caller(1)
	if !ok {
		log.Printf("\u001B[31m[ERROR]\u001B[0m runtime.Caller error in WrapErrorf.")
		return WrapComplexError(cause, Error("%s", msg), "", -1)
	}
	parts := strings.Split(filepath, "/")
	if len(parts) > 3 {
		filepath = strings.Join(parts[len(parts)-3:], "/")
	}
	// The second parameter of args is requestId, if the error message is NotFoundMsg the requestId need to be returned.
	if msg == NotFoundMsg && len(args) == 2 {
		msg += RequestIdMsg
	}
	return WrapComplexError(cause, fmt.Errorf(msg, args...), filepath, line)
}

func WrapComplexError(cause, err error, filepath string, fileline int) error {
	return &ComplexError{
		Cause: cause,
		Err:   err,
		Path:  filepath,
		Line:  fileline,
	}
}

func NotFoundErr(args ...interface{}) error {
	return fmt.Errorf(notFoundFmt, args...)
}

const notFoundFmt = "The specified %s %s is not found."

// A default message of ComplexError's Err. It is format to Resource <resource-id> <operation> Failed!!! <error source>
const DefaultErrorMsg = "Resource %s %s Failed!!! %s"
const DataDefaultErrorMsg = "Datasource %s %s Failed!!! %s"
const NotFoundWithResponse = ResourceNotfound + "!!! Response: %v"
const FailedGetAttributeMsg = "Getting resource %s attribute by path %s failed!!! Response: %v."

type ErrorSource string

const (
	AlibabaCloudSdkGoERROR = ErrorSource("[SDK alibaba-cloud-sdk-go ERROR]")
	AliyunLogGoSdkERROR    = ErrorSource("[SDK aliyun-log-go-sdk ERROR]")
	AliyunDatahubSdkGo     = ErrorSource("[SDK aliyun-datahub-sdk-go ERROR]")
	AliyunOssGoSdk         = ErrorSource("[SDK aliyun-oss-go-sdk ERROR]")
	FcGoSdk                = ErrorSource("[SDK fc-go-sdk ERROR]")
	DenverdinoAliyungo     = ErrorSource("[SDK denverdino/aliyungo ERROR]")
	AliyunTablestoreGoSdk  = ErrorSource("[SDK aliyun-tablestore-go-sdk ERROR]")
	AliMnsERROR            = ErrorSource("[SDK ali_mns ERROR]")
	ProviderERROR          = ErrorSource("[Provider ERROR]")
)

// An Error represents a custom error for Terraform failure response
type ProviderError struct {
	errorCode string
	message   string
}

func (e *ProviderError) Error() string {
	return fmt.Sprintf("[ERROR] Terraform Alicloud Provider Error: Code: %s Message: %s", e.errorCode, e.message)
}

func (err *ProviderError) ErrorCode() string {
	return err.errorCode
}

func (err *ProviderError) Message() string {
	return err.message
}

// NewProviderError is the exported constructor: the fields stay
// unexported, so callers outside this package (including package
// alicloud, through its aliases) build values with it.
func NewProviderError(errorCode, message string) *ProviderError {
	return &ProviderError{errorCode: errorCode, message: message}
}

func GetNotFoundErrorFromString(str string) error {
	return &ProviderError{
		errorCode: InstanceNotFound,
		message:   str,
	}
}

func GetTimeErrorFromString(str string) error {
	return &ProviderError{
		errorCode: "WaitForTimeout",
		message:   str,
	}
}

func NotFoundError(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*ComplexError); ok {
		if e.Err != nil && strings.HasPrefix(e.Err.Error(), ResourceNotfound) {
			return true
		}
		return NotFoundError(e.Cause)
	}

	if e, ok := err.(*tea.SDKError); ok {
		return tea.IntValue(e.StatusCode) == 404 || regexp.MustCompile(NotFound).MatchString(tea.StringValue(e.Message))
	}

	if e, ok := err.(*errors.ServerError); ok {
		return e.ErrorCode() == InstanceNotFound || e.ErrorCode() == RamInstanceNotFound || e.ErrorCode() == NotFound || e.HttpStatus() == 404 || strings.Contains(strings.ToLower(e.Message()), MessageInstanceNotFound)
	}

	if e, ok := err.(*ProviderError); ok {
		return e.ErrorCode() == InstanceNotFound || e.ErrorCode() == RamInstanceNotFound || e.ErrorCode() == NotFound || strings.Contains(strings.ToLower(e.Message()), MessageInstanceNotFound)
	}

	if e, ok := err.(*common.Error); ok {
		return e.Code == InstanceNotFound || e.Code == RamInstanceNotFound || e.Code == NotFound || strings.Contains(strings.ToLower(e.Message), MessageInstanceNotFound)
	}

	if e, ok := err.(oss.ServiceError); ok {
		return e.StatusCode == 404 || strings.HasPrefix(e.Code, "NoSuch") || strings.HasPrefix(e.Message, "No Row found")
	}

	if e, ok := err.(*otsTunnel.TunnelError); ok {
		if e.Code == otsTunnel.ErrCodeParamInvalid && strings.Contains(strings.ToLower(e.Message), "tunnel not exist") {
			return true
		}
		if e.Code == otsTunnel.ErrCodePermissionDenied {
			lowerMsg := strings.ToLower(e.Message)
			if strings.Contains(lowerMsg, "is not found") || strings.Contains(lowerMsg, "instance is not running") {
				return true
			}
		}
	}

	return false
}

// Error codes shared by retry policies and error matching. Package
// alicloud re-exports them as const aliases.
const (
	ServiceUnavailable = "ServiceUnavailable"
	ForbiddenInstance  = "Forbidden.InstanceNotFound"
	Throttling         = "Throttling"
	ThrottlingUser     = "Throttling.User"

	AliyunGoClientFailure = "AliyunGoClientFailure"

	LogClientTimeout = "Client.Timeout exceeded while awaiting headers"

	InvalidFileSystemStatus_Ordering = "InvalidFileSystemStatus.Ordering"
	NotFoundArticle                  = "not found article by given param"
)

// A default message of ComplexError's Err. It is format to Resource <resource-id> <operation> Failed!!! <error source>
const ResponseCodeMsg = "Resource %s %s Failed!!! %v"
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

// Retryable error-code lists fed into IsExpectedErrors by the resources.
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

func IsExpectedErrors(err error, expectCodes []string) bool {
	if err == nil {
		return false
	}

	if e, ok := err.(*ComplexError); ok {
		return IsExpectedErrors(e.Cause, expectCodes)
	}

	if e, ok := err.(*tea.SDKError); ok {
		for _, code := range expectCodes {
			// The second statement aims to match the tea sdk history bug
			if *e.Code == code || strings.HasPrefix(code, *e.Code) || strings.Contains(*e.Data, code) {
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
