package fcv2

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// Service wraps the FCV2 API calls used by the resources of this package.
// The root alicloud package cannot be imported from a service package, so
// the encapsulated calls live here with the resources.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor: the client field stays unexported,
// so callers outside the package build values with it.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

// DescribeFcv2Function <<< Encapsulated get interface for Fcv2 Function.

func (s *Service) DescribeFcv2Function(id string) (object map[string]interface{}, err error) {

	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	functionName := parts[1]
	serviceName := parts[0]
	action := fmt.Sprintf("/2021-04-06/services/%s/functions/%s", serviceName, functionName)
	request = make(map[string]interface{})
	query = make(map[string]*string)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("FC-Open", "2021-04-06", action, query, nil, nil)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"FunctionNotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Function", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *Service) Fcv2FunctionStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeFcv2Function(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		currentStatus := fmt.Sprint(object[field])
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeFcv2Function >>> Encapsulated.

// DescribeFunctionOnDemandConfig <<< Encapsulated get interface for the
// FunctionOnDemandConfig. The id is serviceName:functionName:qualifier, where
// qualifier is a service alias or "LATEST" (the API rejects bare versions).
func (s *Service) DescribeFunctionOnDemandConfig(id string) (object map[string]interface{}, err error) {

	client := s.client
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		return object, errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	serviceName := parts[0]
	functionName := parts[1]
	qualifier := parts[2]
	action := fmt.Sprintf("/2021-04-06/services/%s/functions/%s/on-demand-config", serviceName, functionName)
	query = make(map[string]*string)
	query["qualifier"] = helper.StringPointer(qualifier)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		object, err = client.RoaGet("FC-Open", "2021-04-06", action, query, nil, nil)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, object, query)
		return nil
	})

	if err != nil {
		// FC-Open returns 404 once the on-demand config is gone; the tea-rpc
		// client surfaces that as a *tea.SDKError whose StatusCode is 404, which
		// errs.NotFoundError reaches through the wrapped ComplexError's Cause.
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return object, nil
}

// DescribeFunctionOnDemandConfig >>> Encapsulated.
