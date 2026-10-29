package polardbx

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

type DrdsServiceV2 struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor used by external test packages
// (package polardbx_test) to build a DrdsServiceV2, since the client field
// is unexported. It mirrors the established migrated-package convention.
func NewService(client *connectivity.AliyunClient) *DrdsServiceV2 {
	return &DrdsServiceV2{client: client}
}

func (s *DrdsServiceV2) DrdsPolardbxInstanceAsynJobs(d *schema.ResourceData, resp map[string]interface{}) (object map[string]interface{}, err error) {

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

func (s *DrdsServiceV2) DrdsPolardbxInstanceJobStateRefreshFunc(d *schema.ResourceData, response map[string]interface{}, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DrdsPolardbxInstanceAsynJobs(d, response)
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

func (s *DrdsServiceV2) DrdsPolardbxInstanceAsynDeleteJobs(d *schema.ResourceData, resp map[string]interface{}) (object map[string]interface{}, err error) {

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

func (s *DrdsServiceV2) DrdsPolardbxInstanceDeleteJobStateRefreshFunc(d *schema.ResourceData, response map[string]interface{}, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DrdsPolardbxInstanceAsynDeleteJobs(d, response)
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

// DescribeDrdsPolardbxInstance <<< Encapsulated get interface for Drds PolardbxInstance.

func (s *DrdsServiceV2) DescribeDrdsPolardbxInstance(id string) (object map[string]interface{}, err error) {
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

func (s *DrdsServiceV2) DrdsPolardbxInstanceStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.DrdsPolardbxInstanceStateRefreshFuncWithApi(id, field, failStates, s.DescribeDrdsPolardbxInstance)
}

func (s *DrdsServiceV2) DrdsPolardbxInstanceStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		object["PayType"] = convertDrdsPolardbxInstanceDBInstancePayTypeResponse(object["PayType"])
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

func (s *DrdsServiceV2) DescribeAsyncDrdsPolardbxInstanceStateRefreshFunc(d *schema.ResourceData, res map[string]interface{}, field string, failStates []string) retry.StateRefreshFunc {
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
		object["PayType"] = convertDrdsPolardbxInstanceDBInstancePayTypeResponse(object["PayType"])
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

// DescribeDrdsPolardbxInstance >>> Encapsulated.
// DescribeAsyncDescribeTasks <<< Encapsulated for Drds.
func (s *DrdsServiceV2) DescribeAsyncDescribeTasks(d *schema.ResourceData, res map[string]interface{}) (object map[string]interface{}, err error) {
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
