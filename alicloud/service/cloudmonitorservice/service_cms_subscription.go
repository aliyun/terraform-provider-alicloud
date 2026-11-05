package cloudmonitorservice

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// DescribeCmsSubscription fetches one Cms 2024-03-30 event subscription. The
// workspace query parameter is optional on the API; it is forwarded when the
// caller has it so that workspace-scoped subscriptions resolve consistently.
func (s *CmsService) DescribeCmsSubscription(id, workspace string) (object map[string]interface{}, err error) {
	client := s.client
	action := fmt.Sprintf("/subscriptions/%s", id)
	query := make(map[string]*string)
	if workspace != "" {
		query["workspace"] = helper.StringPointer(workspace)
	}

	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(cmsSubscriptionTimeout, func() *retry.RetryError {
		response, err = client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			if errs.NotFoundError(err) || errs.IsExpectedErrors(err, []string{"NotFound", "ResourceNotFound"}) {
				return retry.NonRetryableError(errs.WrapErrorf(errs.NotFoundErr("Cms:Subscription", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response)))
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, query)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	obj, err := jsonpath.Get("$.subscription", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.subscription", response)
	}
	if obj == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Cms:Subscription", id), errs.NotFoundWithResponse, response)
	}
	return obj.(map[string]interface{}), nil
}
