package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type ResourcesharingService struct {
	client *connectivity.AliyunClient
}

func (s *ResourcesharingService) DescribeResourceManagerResourceShare(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListResourceShares"
	client := s.client

	request := map[string]interface{}{
		"ResourceOwner":    "Self",
		"ResourceShareIds": []string{id},
		"MaxResults":       1,
	}

	idExist := false
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("ResourceSharing", "2020-01-10", action, nil, request, true)
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

		resp, err := jsonpath.Get("$.ResourceShares", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ResourceShares", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, WrapErrorf(NotFoundErr("ResourceManager:ResourceShare", id), NotFoundWithResponse, response)
		}

		for _, v := range resp.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["ResourceShareId"]) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	if !idExist {
		return object, WrapErrorf(NotFoundErr("ResourceManager:ResourceShare", id), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *ResourcesharingService) ResourceManagerResourceShareStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeResourceManagerResourceShare(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["ResourceShareStatus"].(string) == failState {
				return object, object["ResourceShareStatus"].(string), WrapError(Error(FailedToReachTargetStatus, object["ResourceShareStatus"].(string)))
			}
		}

		return object, object["ResourceShareStatus"].(string), nil
	}
}

func (s *ResourcesharingService) DescribeResourceManagerSharedResource(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListResourceShareAssociations"
	client := s.client

	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":         s.client.RegionId,
		"AssociationType":  "Resource",
		"ResourceShareIds": []string{parts[0]},
		"ResourceId":       parts[1],
	}

	idExist := false
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("ResourceSharing", "2020-01-10", action, nil, request, true)
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

		resp, err := jsonpath.Get("$.ResourceShareAssociations", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ResourceShareAssociations", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, WrapErrorf(NotFoundErr("ResourceManager:SharedResource", id), NotFoundWithResponse, response)
		}

		for _, v := range resp.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["ResourceShareId"]) == parts[0] && fmt.Sprint(v.(map[string]interface{})["EntityId"]) == parts[1] && fmt.Sprint(v.(map[string]interface{})["EntityType"]) == parts[2] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	if !idExist {
		return object, WrapErrorf(NotFoundErr("ResourceManager:SharedResource", id), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *ResourcesharingService) ResourceManagerSharedResourceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeResourceManagerSharedResource(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["AssociationStatus"].(string) == failState {
				return object, object["AssociationStatus"].(string), WrapError(Error(FailedToReachTargetStatus, object["AssociationStatus"].(string)))
			}
		}
		return object, object["AssociationStatus"].(string), nil
	}
}

func (s *ResourcesharingService) DescribeResourceManagerSharedTarget(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListResourceShareAssociations"
	client := s.client

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":         s.client.RegionId,
		"AssociationType":  "Target",
		"ResourceShareIds": []string{parts[0]},
		"Target":           parts[1],
	}

	idExist := false
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("ResourceSharing", "2020-01-10", action, nil, request, true)
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

		resp, err := jsonpath.Get("$.ResourceShareAssociations", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ResourceShareAssociations", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, WrapErrorf(NotFoundErr("ResourceManager:SharedTarget", id), NotFoundWithResponse, response)
		}

		for _, v := range resp.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["ResourceShareId"]) == parts[0] && fmt.Sprint(v.(map[string]interface{})["EntityId"]) == parts[1] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	if !idExist {
		return object, WrapErrorf(NotFoundErr("ResourceManager:SharedTarget", id), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *ResourcesharingService) ResourceManagerSharedTargetStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeResourceManagerSharedTarget(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["AssociationStatus"].(string) == failState {
				return object, object["AssociationStatus"].(string), WrapError(Error(FailedToReachTargetStatus, object["AssociationStatus"].(string)))
			}
		}
		return object, object["AssociationStatus"].(string), nil
	}
}
