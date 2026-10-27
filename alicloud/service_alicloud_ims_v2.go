package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ImsServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeImsUserTags <<< Encapsulated get tags for Ims RAM user.
// Tags on a RAM user are managed by the IMS (2019-08-15) tagging API with
// ResourceType=user and the resource id being the RAM user id, because the
// legacy RAM (2015-05-01) TagResources API does not enumerate "user" as a
// taggable resource type.
func (s *ImsServiceV2) DescribeImsUserTags(userId string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceType"] = "user"
	request["ResourceId.1"] = userId

	action := "ListTagResources"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ims", "2019-08-15", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"EntityNotExist.User"}) {
			return object, WrapErrorf(NotFoundErr("RamUser", userId), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, userId, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

// SetImsUserTags <<< Encapsulated tag function for Ims RAM user.
func (s *ImsServiceV2) SetImsUserTags(d *schema.ResourceData) error {
	if d.HasChange("tags") {
		client := s.client
		var response map[string]interface{}
		var err error

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
				"ResourceType": "user",
				"ResourceId.1": d.Id(),
			}
			query := make(map[string]interface{})
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("Ims", "2019-08-15", action, query, request, true)
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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}

		if len(added) > 0 {
			action := "TagResources"
			request := map[string]interface{}{
				"ResourceType": "user",
				"ResourceId.1": d.Id(),
			}
			query := make(map[string]interface{})
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("Ims", "2019-08-15", action, query, request, true)
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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}
	}
	return nil
}
