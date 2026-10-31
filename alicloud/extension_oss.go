package alicloud

import (
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"syscall"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type LifecycleRuleStatus string

const (
	ExpirationStatusEnabled  = LifecycleRuleStatus("Enabled")
	ExpirationStatusDisabled = LifecycleRuleStatus("Disabled")
)

func ossNotFoundError(err error) bool {
	if e, ok := err.(oss.ServiceError); ok &&
		(e.StatusCode == 404 || strings.HasPrefix(e.Code, "NoSuch") || strings.HasPrefix(e.Message, "No Row found")) {
		return true
	}
	return false
}

// ossNoSuchBucketError reports whether err is precisely the OSS NoSuchBucket
// service error. Unlike ossNotFoundError it never matches a bare 404.
func ossNoSuchBucketError(err error) bool {
	var serviceErr oss.ServiceError
	if errors.As(err, &serviceErr) {
		return serviceErr.Code == "NoSuchBucket"
	}
	return false
}

func ossRetryableServiceError(err error) bool {
	var serviceErr oss.ServiceError
	if errors.As(err, &serviceErr) {
		if serviceErr.StatusCode == http.StatusBadRequest && serviceErr.Code == "RequestTimeout" {
			return strings.Contains(serviceErr.Message, "Request timeout")
		}
		return serviceErr.Code != "NotImplemented" &&
			(serviceErr.StatusCode >= 500 || serviceErr.Code == "Throttling")
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	return errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) ||
		errors.Is(err, syscall.ECONNRESET) || errors.Is(err, syscall.ECONNREFUSED) ||
		errors.Is(err, syscall.EPIPE)
}
