package drds

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// The root alicloud package cannot be imported from a service package, so the
// encapsulated polardbx calls live here with the resource. DrdsServiceV2 was the
// legacy root-package name; the same methods are now exposed on Service via
// NewService as the single external entry point for the drds product package.

type Service struct {
	client *connectivity.AliyunClient
}

func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

func (s *Service) PolardbxInstanceAsynJobs(d *schema.ResourceData, resp map[string]interface{}) (object map[string]interface{}, err error) {

	client := s.client
	id := d.Id()
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeDBInstanceAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	query["DBInstanceName"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardbx", "2020-02-02", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"DBInstance.NotFound", "RAM.Permission.Denied"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("PolardbXInstance", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.DBInstance", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.DBInstance", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *Service) PolardbxInstanceJobStateRefreshFunc(d *schema.ResourceData, response map[string]interface{}, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.PolardbxInstanceAsynJobs(d, response)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		currentStatus := fmt.Sprint(object["Status"])
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

func (s *Service) PolardbxInstanceAsynDeleteJobs(d *schema.ResourceData, resp map[string]interface{}) (object map[string]interface{}, err error) {

	client := s.client
	id := d.Id()
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeTasks"
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	query["DBInstanceId"] = d.Id()
	query["StartTime"] = d.Get("create_time")
	query["EndTime"] = d.Get("create_time")
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardbx", "2020-02-02", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"DBInstance.NotFound", "RAM.Permission.Denied"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("PolardbXInstance", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Items", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Items", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *Service) PolardbxInstanceDeleteJobStateRefreshFunc(d *schema.ResourceData, response map[string]interface{}, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.PolardbxInstanceAsynDeleteJobs(d, response)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		currentStatus := fmt.Sprint(object["Status"])
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribePolardbxInstance <<< Encapsulated get interface for Drds PolardbxInstance.

func (s *Service) DescribePolardbxInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBInstanceName"] = id
	request["RegionId"] = client.RegionId
	action := "DescribeDBInstanceAttribute"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardbx", "2020-02-02", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"DBInstance.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("PolardbxInstance", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.DBInstance", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.DBInstance", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) PolardbxInstanceStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.PolardbxInstanceStateRefreshFuncWithApi(id, field, failStates, s.DescribePolardbxInstance)
}

func (s *Service) PolardbxInstanceStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		object["PayType"] = convertPolardbxInstanceDBInstancePayTypeResponse(object["PayType"])
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

func (s *Service) DescribeAsyncPolardbxInstanceStateRefreshFunc(d *schema.ResourceData, res map[string]interface{}, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAsyncDescribeTasks(d, res)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
		}
		if object == nil {
			return object, "", nil
		}
		object["PayType"] = convertPolardbxInstanceDBInstancePayTypeResponse(object["PayType"])
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
				if _err, ok := object["error"]; ok {
					return _err, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
				}
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribePolardbxInstance >>> Encapsulated.
// DescribeAsyncDescribeTasks <<< Encapsulated for Drds.
func (s *Service) DescribeAsyncDescribeTasks(d *schema.ResourceData, res map[string]interface{}) (object map[string]interface{}, err error) {
	client := s.client
	id := d.Id()
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBInstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	action := "DescribeTasks"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardbx", "2020-02-02", action, query, request, true)

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
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

// DescribeAsyncDescribeTasks >>> Encapsulated.

// QueryAvailableInstanceByInstanceID resolves the real BSS billing identity
// (ProductCode/ProductType/SubscriptionType) of a subscription instance solely
// from its InstanceID, using the domestic BSS endpoint with no
// ProductCode/ProductType filter and no international-endpoint fallback.
//
// The billing ProductCode of a PolarDB-X 2.0 subscription instance is "drds"
// (ProductType "drds_polarxpre_public_cn"), not the RPC service name
// "polardbx". Filtering by a wrong ProductCode yields an empty domestic list,
// and the previous fallback to the international endpoint produced a misleading
// "account not exists" (40001) for domestic accounts. This helper avoids both
// traps by querying only on InstanceIDs and surfacing NotFound when the instance
// is not a subscription resource under this account.
func (s *Service) QueryAvailableInstanceByInstanceID(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "QueryAvailableInstances"
	request := map[string]interface{}{
		"InstanceIDs": id,
	}
	// Intentionally no ProductCode/ProductType filter and no Region: the domestic
	// BSS resolves the instance by InstanceIDs. No international fallback: an
	// empty domestic list means the instance is not a subscription resource under
	// this account, not that the account does not exist internationally.
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, "")
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

	// Validate the Success flag before reading Data.InstanceList: BSS may return
	// a 200 response with Success=false (e.g. NotApplicable/permission denied for
	// the queried instance). Without this check an empty list would be misread as
	// NotFound instead of an API-level rejection, masking real errors.
	if fmt.Sprint(response["Success"]) == "null" {
		return object, errs.WrapErrorf(fmt.Errorf("%s failed, response: %v", action, response), errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) != "true" {
		return object, errs.WrapErrorf(fmt.Errorf("%s failed, response: %v", action, response), errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data.InstanceList", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data.InstanceList", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("PolarDBXInstance", id), errs.NotFoundWithResponse, response)
	}
	if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["InstanceID"]) != id {
		return object, errs.WrapErrorf(errs.NotFoundErr("PolarDBXInstance", id), errs.NotFoundWithResponse, response)
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func convertPolardbxInstanceDBInstancePayTypeResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "Postpaid":
		return "POSTPAY"
	case "Prepaid":
		return "PREPAY"
	}
	return source
}
