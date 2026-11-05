package ram

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeImsUserTags and SetImsUserTags migrate the IMS (2019-08-15) user
// tagging helpers that the root alicloud ImsServiceV2 exposed. Tags on a RAM
// user are managed by the IMS tagging API with ResourceType=user, because the
// legacy RAM (2015-05-01) TagResources API does not enumerate "user" as a
// taggable resource type. The methods now live on RamServiceV2 so the ram
// subpackage does not reach back into the root alicloud package (which would be
// a circular import).

// DescribeImsUserTags returns the raw ListTagResources response for a RAM user.
func (s *RamServiceV2) DescribeImsUserTags(userId string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceType"] = "user"
	request["ResourceId.1"] = userId

	action := "ListTagResources"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ims", "2019-08-15", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"EntityNotExist.User"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("RamUser", userId), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, userId, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

// SetImsUserTags adds and removes tags on a RAM user through the IMS tagging API.
func (s *RamServiceV2) SetImsUserTags(d *schema.ResourceData) error {
	if d.HasChange("tags") {
		client := s.client
		var response map[string]interface{}
		var err error

		added, removed := tags.ParsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !tags.IgnoredTags(v, "") {
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
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("Ims", "2019-08-15", action, query, request, true)
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
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
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
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("Ims", "2019-08-15", action, query, request, true)
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
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
		}
	}
	return nil
}
