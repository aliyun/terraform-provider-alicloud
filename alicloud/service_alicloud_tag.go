package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type TagService struct {
	client *connectivity.AliyunClient
}

func (s *TagService) ListTagValues(key string) (object []interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListTagValues"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"Key":       key,
		"QueryType": "MetaTag",
	}
	values := make([]interface{}, 0)
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Tag", "2018-08-28", action, nil, request, true)
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
			if IsExpectedErrors(err, []string{"MetaTagKeyNotFound"}) {
				return object, WrapErrorf(Error(GetNotFoundMessage("TAG:MetaTag", key)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
			}
			return object, WrapErrorf(err, DefaultErrorMsg, key, action, AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.Values.Value", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, key, "$.Values.Value", response)
		}
		values = append(values, v.([]interface{})...)
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
		request["NextToken"] = response["NextToken"]
	}

	return values, nil
}

func (s *TagService) DescribeTagPolicy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetPolicy"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"PolicyId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Tag", "2018-08-28", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidParameter.PolicyId"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Tag:Policy", id)), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Policy", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Policy", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *TagService) DescribeTagPolicyAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListTargetsForPolicy"

	client := s.client
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return nil, WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"PolicyId": parts[0],
	}
	idExist := false
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Tag", "2018-08-28", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidParameter.PolicyId"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Tag:PolicyAttachment", id)), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.Targets", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Targets", response)
	}
	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("Tag:PolicyAttachment", id)), NotFoundWithResponse, response)
	}
	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["TargetId"]) == parts[1] && fmt.Sprint(v.(map[string]interface{})["TargetType"]) == parts[2] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("Tag:PolicyAttachment", id)), NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *TagService) DescribeTagValue(id string) (object []interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	parts, err := ParseResourceIdN(id, 2)
	if err != nil {
		return object, WrapError(err)
	}
	action := "ListTagValues"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"Key":       parts[1],
		"QueryType": "MetaTag",
	}
	values := make([]interface{}, 0)
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Tag", "2018-08-28", action, nil, request, true)
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
			if IsExpectedErrors(err, []string{"MetaTagKeyNotFound"}) {
				return object, WrapErrorf(Error(GetNotFoundMessage("TAG:MetaTag", parts[0])), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
			}
			return object, WrapErrorf(err, DefaultErrorMsg, parts[0], action, AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.Values.Value", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, parts[0], "$.Values.Value", response)
		}
		values = append(values, v.([]interface{})...)
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
		request["NextToken"] = response["NextToken"]
	}

	return values, nil
}
