package oceanbasepro

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

// Service converges the former root-package OceanBaseServiceV2 and
// OceanBaseProService into a single service for the OceanBasePro product.
// The method bodies are migrated verbatim from service_alicloud_ocean_base_v2.go
// and service_alicloud_ocean_base_pro.go; only the package, receiver and leaf
// helper qualifications changed.
type Service struct {
	client *connectivity.AliyunClient
}

func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

// DescribeOceanBaseInstance <<< Encapsulated get interface for OceanBase Instance.
// Body migrated verbatim from OceanBaseServiceV2.DescribeOceanBaseInstance
// (resource Read/Create/Update/Delete and the ACC reflection path use this name).
func (s *Service) DescribeOceanBaseInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeInstance"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = id

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("OceanBasePro", "2019-09-01", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"IllegalOperation.Resource", "UnknownError"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Instance", id), errs.NotFoundMsg, response)
		}
		logs.AddDebug(action, response, request)
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Instance", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Instance", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) DescribeDescribeInstances(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeInstances"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = id

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("OceanBasePro", "2019-09-01", action, query, request, true)

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
		logs.AddDebug(action, response, request)
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Instances[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Instances[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Instance", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *Service) OceanBaseInstanceStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeOceanBaseInstance(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOceanBaseInstanceDetail fetches a single instance for the data source's
// enable_details path. The body is migrated verbatim from the former
// OceanBaseProService.DescribeOceanBaseInstance (RegionId in the request, 5m
// retry, Instance-presence NotFound). It is renamed because the resource/ACC
// reflection path owns the DescribeOceanBaseInstance name with a different retry
// budget and NotFound contract; both behaviors are preserved unchanged.
func (s *Service) DescribeOceanBaseInstanceDetail(id string) (object map[string]interface{}, err error) {
	client := s.client
	request := map[string]interface{}{
		"InstanceId": id,
		"RegionId":   s.client.RegionId,
	}

	var response map[string]interface{}
	action := "DescribeInstance"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("OceanBasePro", "2019-09-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if _, ok := response["Instance"]; !ok {
		return object, errs.WrapErrorf(errs.NotFoundErr("Instance", id), errs.NotFoundWithResponse, response)
	}
	v, err := jsonpath.Get("$.Instance", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Instance", response)
	}
	return v.(map[string]interface{}), nil
}

// DescribeOceanBaseInstance >>> Encapsulated.
