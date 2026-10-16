package eas

import (
	"fmt"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// easProductCode and easAPIVersion identify the EAS (PAI) Roa API used by the
// resource group resource. The EAS API path embeds the cluster id, which is
// the provider region; the resource id alone addresses a resource within the
// configured region.
const (
	easProductCode = "EAS"
	easAPIVersion  = "2021-07-01"
)

// EasService holds the shared EAS API calls used by the eas service package.
type EasService struct {
	client *connectivity.AliyunClient
}

// NewEasService constructs an EasService bound to the provider client.
func NewEasService(client *connectivity.AliyunClient) *EasService {
	return &EasService{client: client}
}

// DescribeResource fetches a single EAS resource group by id. Callers detect
// the not-found case with isResourceNotFound.
func (s *EasService) DescribeResource(id string) (object map[string]interface{}, err error) {
	client := s.client
	action := fmt.Sprintf("/api/v2/resources/%s/%s", client.RegionId, id)
	query := make(map[string]*string)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, e := client.RoaGet(easProductCode, easAPIVersion, action, query, nil, nil)
		if e != nil {
			if isResourceNotFound(e) {
				return retry.NonRetryableError(e)
			}
			if helper.NeedRetry(e) {
				wait()
				return retry.RetryableError(e)
			}
			return retry.NonRetryableError(e)
		}
		object = response
		return nil
	})
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	return object, nil
}

// WaitForResource polls DescribeResource until Status reaches one of target
// states or one of failStates. A not-found result during a readiness wait is a
// failure, because the resource vanished before becoming ready.
func (s *EasService) WaitForResource(id string, target, failStates []string, timeout time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Target:     target,
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
		Refresh: func() (interface{}, string, error) {
			object, e := s.DescribeResource(id)
			if e != nil {
				if isResourceNotFound(e) {
					return nil, "", fmt.Errorf("the resource %s does not exist while waiting", id)
				}
				return nil, "", e
			}
			status, _ := object["Status"].(string)
			for _, fs := range failStates {
				if status == fs {
					return nil, "", fmt.Errorf("the resource %s reached fail state %s", id, fs)
				}
			}
			return object, status, nil
		},
	}
	_, err := stateConf.WaitForState()
	return err
}

// WaitForResourceGone polls DescribeResource until the resource group no longer
// exists, i.e. DescribeResource reports not-found. Used after Delete so a
// follow-up create in the same name does not race a still-deleting one.
func (s *EasService) WaitForResourceGone(id string, timeout time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Target:     []string{"Gone"},
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
		Refresh: func() (interface{}, string, error) {
			object, e := s.DescribeResource(id)
			if e != nil {
				if isResourceNotFound(e) {
					return nil, "Gone", nil
				}
				return nil, "", e
			}
			status, _ := object["Status"].(string)
			return object, status, nil
		},
	}
	_, err := stateConf.WaitForState()
	return err
}

// isResourceNotFound reports whether err carries the EAS not-found condition.
// The roaRequest returns the raw SDK error (a *tea.SDKError or
// *errors.ServerError); errs.WrapErrorf later wraps it in a *errs.ComplexError.
// errs.NotFoundError covers the 404 status and the standard NotFound message
// patterns; the code substring match covers EAS-style codes such as
// "InvalidResource.NotFound" that are not HTTP 404 and whose message does not
// match the generic instance-not-found wording. This mirrors the root-package
// IsExpectedErrors check used by the SDKv2 eas resources.
func isResourceNotFound(err error) bool {
	if err == nil {
		return false
	}
	if errs.NotFoundError(err) {
		return true
	}
	return errorHasCodeSubstring(err, "NotFound")
}

// IsResourceNotFound is the exported form of isResourceNotFound, for use by
// external test packages (eas_test) that verify a resource is gone after destroy.
func IsResourceNotFound(err error) bool {
	return isResourceNotFound(err)
}

// errorHasCodeSubstring recurses through *errs.ComplexError causes and reports
// whether any tea SDK error code/data or SDK ServerError code/message contains
// the substring.
func errorHasCodeSubstring(err error, sub string) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*errs.ComplexError); ok {
		return errorHasCodeSubstring(e.Cause, sub)
	}
	if e, ok := err.(*tea.SDKError); ok {
		if e.Code != nil && strings.Contains(*e.Code, sub) {
			return true
		}
		if e.Data != nil && strings.Contains(*e.Data, sub) {
			return true
		}
		return false
	}
	if e, ok := err.(*errors.ServerError); ok {
		return strings.Contains(e.ErrorCode(), sub) || strings.Contains(e.Message(), sub)
	}
	return false
}

// incrementalWait mirrors the root-package helper used by the SDKv2 eas
// resources: a short, gently increasing backoff for throttling retries.
func incrementalWait(initial, max time.Duration) func() {
	period := initial
	return func() {
		time.Sleep(period)
		if period < max {
			period *= 2
		}
	}
}
