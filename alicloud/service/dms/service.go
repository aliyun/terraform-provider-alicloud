// Package dms. This file is generated automatically. Please do not modify it manually, thank you!
package dms

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// Service wraps the Dms API calls used by the resources of this package.
// The root alicloud package cannot be imported from a service package, so
// the encapsulated calls live here with the resources.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor: the client field stays
// unexported, so callers outside the package (including the acceptance
// tests) build values with it.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

// DescribeAirflow <<< Encapsulated get interface for Dms Airflow.

func (s *Service) DescribeAirflow(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AirflowId"] = parts[1]
	request["WorkspaceId"] = parts[0]
	request["RegionId"] = client.RegionId
	action := "GetAirflow"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dms", "2025-04-14", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"InvalidAirflow.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Airflow", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Root", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Root", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) AirflowStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAirflow(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeAirflow >>> Encapsulated.
