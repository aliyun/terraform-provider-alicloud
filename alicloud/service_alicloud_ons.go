package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type OnsService struct {
	client *connectivity.AliyunClient
}

func (s *OnsService) DescribeOnsInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "OnsInstanceBaseInfo"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": id,
	}
	response, err = client.RpcPost("Ons", "2019-02-14", action, nil, request, true)
	if err != nil {
		if IsExpectedErrors(err, []string{"INSTANCE_NOT_FOUNDError", "InvalidDomainName.NoExist"}) {
			err = WrapErrorf(NotFoundErr("OnsInstance", id), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.InstanceBaseInfo", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.InstanceBaseInfo", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *OnsService) DescribeOnsTopic(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "OnsTopicList"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": parts[0],
		"Topic":      parts[1],
	}
	response, err = client.RpcPost("Ons", "2019-02-14", action, nil, request, true)
	if err != nil {
		if IsExpectedErrors(err, []string{"AUTH_RESOURCE_OWNER_ERROR", "INSTANCE_NOT_FOUND"}) {
			err = WrapErrorf(NotFoundErr("OnsTopic", id), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Data.PublishInfoDo", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data.PublishInfoDo", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("Ons", id), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["Topic"].(string) != parts[1] {
			return object, WrapErrorf(NotFoundErr("Ons", id), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *OnsService) OceanOnsTopicStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeOnsTopic(id)
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)
		if field == "$.CreateTime" {
			if currentStatus != "" {
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

func (s *OnsService) DescribeOnsGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "OnsGroupList"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"GroupId":    parts[1],
		"InstanceId": parts[0],
		"GroupType":  "all",
	}
	response, err = client.RpcPost("Ons", "2019-02-14", action, nil, request, true)
	if err != nil {
		if IsExpectedErrors(err, []string{"AUTH_RESOURCE_OWNER_ERROR", "INSTANCE_NOT_FOUND"}) {
			err = WrapErrorf(NotFoundErr("OnsGroup", id), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Data.SubscribeInfoDo", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data.SubscribeInfoDo", response)
	}
	var exist bool
	var index int = 0
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("Ons", id), NotFoundWithResponse, response)
	} else {
		// special handling for fuzzy matching
		onsGroupList := v.([]interface{})
		for i, onsGroup := range onsGroupList {
			if onsGroup.(map[string]interface{})["GroupId"].(string) == parts[1] {
				exist = true
				index = i
				break
			}
		}
		if !exist {
			return object, WrapErrorf(NotFoundErr("Ons", id), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[index].(map[string]interface{})
	return object, nil
}

func (s *OnsService) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	var parts []string
	var err error
	if resourceType == "GROUP" || resourceType == "TOPIC" {
		parts, err = ParseResourceId(d.Id(), 2)
		if err != nil {
			return WrapError(err)
		}
	}
	client := s.client
	if d.HasChange("tags") {
		added, removed := parsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action := "UntagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
			}
			if resourceType == "INSTANCE" {
				request["ResourceId.1"] = d.Id()
			}
			if resourceType == "GROUP" || resourceType == "TOPIC" {
				request["InstanceId"] = parts[0]
				request["ResourceId.1"] = parts[1]
			}
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := client.RpcPost("Ons", "2019-02-14", action, nil, request, false)
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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}
		if len(added) > 0 {
			action := "TagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
			}
			if resourceType == "INSTANCE" {
				request["ResourceId.1"] = d.Id()
			}
			if resourceType == "GROUP" || resourceType == "TOPIC" {
				request["InstanceId"] = parts[0]
				request["ResourceId.1"] = parts[1]
			}
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := client.RpcPost("Ons", "2019-02-14", action, nil, request, false)
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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}
		d.SetPartial("tags")
	}
	return nil
}

func (s *OnsService) OnsTopicStatus(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "OnsTopicStatus"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": parts[0],
		"Topic":      parts[1],
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(11*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ons", "2019-02-14", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrors(err, []string{"AUTH_RESOURCE_OWNER_ERROR", "INSTANCE_NOT_FOUND"}) {
				err = WrapErrorf(NotFoundErr("OnsTopic", id), NotFoundMsg, ProviderERROR)
				return resource.NonRetryableError(err)
			}
			err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		v, err := jsonpath.Get("$.Data", response)
		if err != nil {
			return resource.NonRetryableError(WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response))
		}
		object = v.(map[string]interface{})
		return nil
	})
	return
}

func (s *OnsService) ListTagResources(id string, resourceType string) (object interface{}, err error) {
	client := s.client
	action := "ListTagResources"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"ResourceType": resourceType,
		"ResourceId.1": id,
	}
	tags := make([]interface{}, 0)
	var response map[string]interface{}

	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err := client.RpcPost("Ons", "2019-02-14", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{Throttling}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			v, err := jsonpath.Get("$.TagResources", response)
			if err != nil {
				return resource.NonRetryableError(WrapErrorf(err, FailedGetAttributeMsg, id, "$.TagResources.TagResource", response))
			}
			if v != nil {
				if v != nil {
					tags = append(tags, v.([]interface{})...)
				}
			}
			return nil
		})
		if err != nil {
			err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
			return
		}
		if response["NextToken"] == nil {
			break
		}
		request["NextToken"] = response["NextToken"]
	}

	return tags, nil
}
