package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type DbfsService struct {
	client *connectivity.AliyunClient
}

func (s *DbfsService) DescribeDbfsInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListDbfs"

	client := s.client
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}

	idExist := false
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("DBFS", "2020-04-18", action, nil, request, true)
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

		resp, err := jsonpath.Get("$.DBFSInfo", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DBFSInfo", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("Dbfs:Instance", id)), NotFoundWithResponse, response)
		}

		for _, v := range resp.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["FsId"]) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if len(resp.([]interface{})) < request["PageSize"].(int) {
			break
		}

		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("Dbfs:Instance", id)), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *DbfsService) DescribeDbfsInstanceAttachment(id string) (object map[string]interface{}, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return object, WrapError(err)
	}

	object, err = s.DescribeDbfsInstance(parts[0])
	if err != nil {
		return object, WrapError(err)
	}

	idExist := false
	if ecsList, ok := object["EcsList"]; ok {
		for _, ecs := range ecsList.([]interface{}) {
			ecsArg := ecs.(map[string]interface{})

			if ecsId, ok := ecsArg["EcsId"]; ok && fmt.Sprint(ecsId) == parts[1] {
				idExist = true
				return object, nil
			}
		}
	}

	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("Dbfs:InstanceAttachment", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(object["RequestId"]))
	}

	return object, nil
}

func (s *DbfsService) DbfsInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDbfsInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}

		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *DbfsService) DescribeDbfsSnapshot(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListSnapshot"

	client := s.client
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}

	idExist := false
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("DBFS", "2020-04-18", action, nil, request, true)
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

		resp, err := jsonpath.Get("$.Snapshots", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Snapshots", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("Dbfs:Snapshot", id)), NotFoundWithResponse, response)
		}

		for _, v := range resp.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["SnapshotId"]) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if len(resp.([]interface{})) < request["PageSize"].(int) {
			break
		}

		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("Dbfs:Snapshot", id)), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *DbfsService) DbfsSnapshotStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDbfsSnapshot(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}

		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *DbfsService) DescribeDbfsServiceLinkedRole(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetServiceLinkedRole"
	request := map[string]interface{}{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("DBFS", "2020-04-18", action, request, nil, true)
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
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *DbfsService) DbfsServiceLinkedRoleStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDbfsServiceLinkedRole(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["DbfsLinkedRole"]) == failState {
				return object, fmt.Sprint(object["DbfsLinkedRole"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["DbfsLinkedRole"])))
			}
		}
		return object, fmt.Sprint(object["DbfsLinkedRole"]), nil
	}
}

func (s *DbfsService) DescribeDbfsAutoSnapShotPolicy(id string) (object map[string]interface{}, err error) {
	client := s.client
	request := map[string]interface{}{
		"PolicyId": id,
		"RegionId": s.client.RegionId,
	}

	var response map[string]interface{}
	action := "GetAutoSnapshotPolicy"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPost("DBFS", "2020-04-18", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"AutoSnapshotPolicyNotFound"}) {
			return object, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}
	return v.(map[string]interface{}), nil
}
