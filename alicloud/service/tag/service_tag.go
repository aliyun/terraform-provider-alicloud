package tag

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type TagService struct {
	client *connectivity.AliyunClient
}

func NewTagService(client *connectivity.AliyunClient) *TagService {
	return &TagService{client: client}
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
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Tag", "2018-08-28", action, nil, request, true)
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
			if errs.IsExpectedErrors(err, []string{"MetaTagKeyNotFound"}) {
				return object, errs.WrapErrorf(errs.NotFoundErr("TAG:MetaTag", key), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
			}
			return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, key, action, errs.AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.Values.Value", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, key, "$.Values.Value", response)
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
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Tag", "2018-08-28", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"InvalidParameter.PolicyId"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Tag:Policy", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Policy", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Policy", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *TagService) DescribeTagPolicyAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListTargetsForPolicy"

	client := s.client
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"PolicyId": parts[0],
	}
	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Tag", "2018-08-28", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"InvalidParameter.PolicyId"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Tag:PolicyAttachment", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.Targets", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Targets", response)
	}
	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Tag:PolicyAttachment", id), errs.NotFoundWithResponse, response)
	}
	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["TargetId"]) == parts[1] && fmt.Sprint(v.(map[string]interface{})["TargetType"]) == parts[2] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}
	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("Tag:PolicyAttachment", id), errs.NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *TagService) DescribeTagValue(id string) (object []interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	parts, err := helper.ParseResourceIdN(id, 2)
	if err != nil {
		return object, errs.WrapError(err)
	}
	action := "ListTagValues"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"Key":       parts[1],
		"QueryType": "MetaTag",
	}
	values := make([]interface{}, 0)
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Tag", "2018-08-28", action, nil, request, true)
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
			if errs.IsExpectedErrors(err, []string{"MetaTagKeyNotFound"}) {
				return object, errs.WrapErrorf(errs.NotFoundErr("TAG:MetaTag", parts[0]), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
			}
			return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, parts[0], action, errs.AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.Values.Value", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, parts[0], "$.Values.Value", response)
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
