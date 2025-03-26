package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type DrdsServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeDrdsPolardbxInstance <<< Encapsulated get interface for Drds PolardbXInstance.

func (s *DrdsServiceV2) DescribeDrdsPolardbxInstance(id string) (object map[string]interface{}, err error) {

	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeDBInstanceAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DBInstanceName"] = id
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("polardbx", "2020-02-02", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"DBInstance.NotFound", "RAM.Permission.Denied"}) {
			return object, WrapErrorf(NotFoundErr("PolardbXInstance", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.DBInstance", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DBInstance", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *DrdsServiceV2) DrdsPolardbxInstanceStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDrdsPolardbxInstance(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		currentStatus := fmt.Sprint(object[field])
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
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

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("polardbx", "2020-02-02", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"DBInstance.NotFound", "RAM.Permission.Denied"}) {
			return object, WrapErrorf(NotFoundErr("PolardbXInstance", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.DBInstance", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DBInstance", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *DrdsServiceV2) DrdsPolardbxInstanceJobStateRefreshFunc(d *schema.ResourceData, response map[string]interface{}, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DrdsPolardbxInstanceAsynJobs(d, response)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		currentStatus := fmt.Sprint(object["Status"])
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
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

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("polardbx", "2020-02-02", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"DBInstance.NotFound", "RAM.Permission.Denied"}) {
			return object, WrapErrorf(NotFoundErr("PolardbXInstance", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Items", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Items", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *DrdsServiceV2) DrdsPolardbxInstanceDeleteJobStateRefreshFunc(d *schema.ResourceData, response map[string]interface{}, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DrdsPolardbxInstanceAsynDeleteJobs(d, response)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		currentStatus := fmt.Sprint(object["Status"])
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeDrdsPolardbxInstance >>> Encapsulated.
