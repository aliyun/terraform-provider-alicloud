package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type DtsService struct {
	client *connectivity.AliyunClient
}

func (s *DtsService) DescribeDtsJobMonitorRule(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewDtsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeJobMonitorRule"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"DtsJobId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &runtime)
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
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *DtsService) ListTagResources(id string, resourceType string) (object interface{}, err error) {
	conn, err := s.client.NewDtsClient()
	if err != nil {
		return nil, WrapError(err)
	}
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
			response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{Throttling}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			v, err := jsonpath.Get("$.TagResources.TagResource", response)
			if err != nil {
				return resource.NonRetryableError(WrapErrorf(err, FailedGetAttributeMsg, id, "$.TagResources.TagResource", response))
			}
			if v != nil {
				tags = append(tags, v.([]interface{})...)
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

func (s *DtsService) SetResourceTags(d *schema.ResourceData, resourceType string) error {

	if d.HasChange("tags") {
		added, removed := parsingTags(d)
		conn, err := s.client.NewDtsClient()
		if err != nil {
			return WrapError(err)
		}

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
				"ResourceId.1": d.Get("dts_instance_id"),
			}
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					if IsThrottling(err) {
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
				"ResourceId.1": d.Get("dts_instance_id"),
			}
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					if IsThrottling(err) {
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

func (s *DtsService) DescribeDtsSubscriptionJob(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewDtsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeDtsJobDetail"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"DtsJobId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"Forbidden.InstanceNotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("DTS:SubscriptionJob", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	if object["Status"] == "Starting" {
		object["Status"] = "Normal"
	}
	return object, nil
}

func (s *DtsService) DescribeDtsJobDetail(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewDtsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeDtsJobDetail"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"DtsJobId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &runtime)
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
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *DtsService) DtsSubscriptionJobStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDtsSubscriptionJob(id)
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
