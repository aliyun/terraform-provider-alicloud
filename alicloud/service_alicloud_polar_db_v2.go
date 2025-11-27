// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type PolarDbServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribePolarDbParameterGroup <<< Encapsulated get interface for PolarDb ParameterGroup.

func (s *PolarDbServiceV2) DescribePolarDbParameterGroup(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ParameterGroupId"] = id
	request["RegionId"] = client.RegionId
	action := "DescribeParameterGroup"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"ParamGroupsNotExist"}) {
			return object, WrapErrorf(NotFoundErr("ParameterGroup", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.ParameterGroup[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ParameterGroup[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("ParameterGroup", id), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *PolarDbServiceV2) PolarDbParameterGroupStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.PolarDbParameterGroupStateRefreshFuncWithApi(id, field, failStates, s.DescribePolarDbParameterGroup)
}

func (s *PolarDbServiceV2) PolarDbParameterGroupStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
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
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribePolarDbParameterGroup >>> Encapsulated.

// DescribePolarDbExtension <<< Encapsulated get interface for PolarDb Extension.

func (s *PolarDbServiceV2) DescribePolarDbExtension(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 4 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 4, len(parts)))
		return nil, err
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBClusterId"] = parts[0]
	request["DBName"] = parts[2]
	request["ExtensionName"] = parts[3]

	action := "DescribeExtensions"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.InstalledExtensions[*]", response)
	if err != nil {
		return object, WrapErrorf(NotFoundErr("Extension", id), NotFoundMsg, response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("Extension", id), NotFoundMsg, response)
	}

	result, _ := v.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if fmt.Sprint(item["Owner"]) != parts[1] {
			continue
		}
		return item, nil
	}
	return object, WrapErrorf(NotFoundErr("Extension", id), NotFoundMsg, response)
}

func (s *PolarDbServiceV2) PolarDbExtensionStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.PolarDbExtensionStateRefreshFuncWithApi(id, field, failStates, s.DescribePolarDbExtension)
}

func (s *PolarDbServiceV2) PolarDbExtensionStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
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
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribePolarDbExtension >>> Encapsulated.

// DescribeAsyncDescribeExtensions <<< Encapsulated for PolarDb.
func (s *PolarDbServiceV2) DescribeAsyncDescribeExtensions(d *schema.ResourceData, res map[string]interface{}) (object map[string]interface{}, err error) {
	client := s.client
	id := d.Id()
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 4 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 4, len(parts)))
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ExtensionName"] = parts[3]
	request["DBClusterId"] = parts[0]
	request["DBName"] = parts[1]

	action := "DescribeExtensions"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return response, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

// DescribeAsyncDescribeExtensions >>> Encapsulated.
