// Package ehpc. This file is generated automatically. Please do not modify it manually, thank you!
package ehpc

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

type EhpcServiceV2 struct {
	client *connectivity.AliyunClient
}

// NewEhpcServiceV2 is the exported constructor used by the package's external
// tests (package ehpc_test); production wiring reaches the service through
// inline EhpcServiceV2{client} construction.
func NewEhpcServiceV2(client *connectivity.AliyunClient) *EhpcServiceV2 {
	return &EhpcServiceV2{client: client}
}

// DescribeEhpcClusterV2 <<< Encapsulated get interface for Ehpc ClusterV2.

func (s *EhpcServiceV2) DescribeEhpcClusterV2(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ClusterId"] = id

	action := "GetCluster"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("EHPC", "2024-07-30", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"ClusterNotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("ClusterV2", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}
func (s *EhpcServiceV2) DescribeClusterV2ListSharedStorages(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ClusterId"] = id

	action := "ListSharedStorages"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("EHPC", "2024-07-30", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"ClusterNotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("ClusterV2", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *EhpcServiceV2) EhpcClusterV2StateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.EhpcClusterV2StateRefreshFuncWithApi(id, field, failStates, s.DescribeEhpcClusterV2)
}

func (s *EhpcServiceV2) EhpcClusterV2StateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		object["Manager.DirectoryService.Version"] = convertEhpcClusterV2ManagerDirectoryServiceVersionResponse(object["Manager.DirectoryService.Version"])
		object["Manager.DNS.Version"] = convertEhpcClusterV2ManagerDNSVersionResponse(object["Manager.DNS.Version"])
		object["Manager.Scheduler.Type"] = convertEhpcClusterV2ManagerSchedulerTypeResponse(object["Manager.Scheduler.Type"])
		object["Manager.Scheduler.Version"] = convertEhpcClusterV2ManagerSchedulerVersionResponse(object["Manager.Scheduler.Version"])
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

// DescribeEhpcClusterV2 >>> Encapsulated.

// DescribeEhpcQueue <<< Encapsulated get interface for Ehpc Queue.

func (s *EhpcServiceV2) DescribeEhpcQueue(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
		return nil, err
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ClusterId"] = parts[0]
	request["QueueName"] = parts[1]

	action := "GetQueue"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("EHPC", "2024-07-30", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"ClusterNotFound", "InvalidParams"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Queue", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Queue", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Queue", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *EhpcServiceV2) EhpcQueueStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.EhpcQueueStateRefreshFuncWithApi(id, field, failStates, s.DescribeEhpcQueue)
}

func (s *EhpcServiceV2) EhpcQueueStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
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

// DescribeEhpcQueue >>> Encapsulated.

// DescribeEhpcUser <<< Encapsulated get interface for Ehpc User.

func (s *EhpcServiceV2) DescribeEhpcUser(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
		return nil, err
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ClusterId"] = parts[0]
	query["UserName"] = parts[1]

	action := "GetUser"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("EHPC", "2024-07-30", action, query, request)

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
		if errs.IsExpectedErrors(err, []string{"UserNotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("User", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.User", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.User", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *EhpcServiceV2) EhpcUserStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.EhpcUserStateRefreshFuncWithApi(id, field, failStates, s.DescribeEhpcUser)
}

func (s *EhpcServiceV2) EhpcUserStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
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

// DescribeEhpcUser >>> Encapsulated.
