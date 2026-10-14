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
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	otsTunnel "github.com/aliyun/aliyun-tablestore-go-sdk/tunnel"
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
