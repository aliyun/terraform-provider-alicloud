// Package ecd holds the resources and data sources of the ECD (Elastic Cloud
// Desktop) product that have been split out of the flat root package.
package ecd

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// Service wraps the ECD API calls used by the resources and data sources of
// this package. The client field stays unexported; callers outside the package
// build values with NewService.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

// DescribeEcdAdConnectorDirectory queries a single AD connector directory by id.
// Moved from the flat layout's EcdService as part of the per-type split.
func (s *Service) DescribeEcdAdConnectorDirectory(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeDirectories"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"DirectoryId": []string{id},
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Directories", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Directories", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("ECD", id), errs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["DirectoryId"]) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("ECD", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

// EcdAdConnectorDirectoryStateRefreshFunc returns a state refresh func that
// follows an AD connector directory through its asynchronous lifecycle: creation
// settles on REGISTERED, while deletion makes the directory disappear (NotFound).
// It mirrors EcdAdConnectorOfficeSiteStateRefreshFunc of the sibling office site
// resource, which already waited for these async terminal states in the flat
// layout; the directory resource previously returned immediately after the RPC,
// which let dependent resources (VSwitch/VPC) be destroyed before the directory
// finished registering or being deleted, producing DependencyViolation.
func (s *Service) EcdAdConnectorDirectoryStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEcdAdConnectorDirectory(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}
