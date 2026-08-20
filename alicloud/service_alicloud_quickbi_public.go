package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type QuickbiPublicService struct {
	client *connectivity.AliyunClient
}

func (s *QuickbiPublicService) DescribeQuickBiUser(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "QueryUserInfoByUserId"
	request := map[string]interface{}{
		"UserId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("quickbi-public", "2020-08-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"User.Not.In.Organization"}) {
			return object, WrapErrorf(NotFoundErr("QuickBI:User", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Result", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Result", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *QuickbiPublicService) QueryUserInfoByUserId(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "QueryUserInfoByUserId"
	request := map[string]interface{}{
		"UserId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("quickbi-public", "2020-08-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"User.Not.In.Organization"}) {
			return object, WrapErrorf(NotFoundErr("QuickBI:User", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Result", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Result", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
