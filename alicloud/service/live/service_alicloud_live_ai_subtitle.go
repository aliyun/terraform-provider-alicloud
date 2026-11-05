package live

import (
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// DescribeLiveAISubtitle fetches a single AI subtitle template by SubtitleId.
//
// DescribeLiveAISubtitle is a List-style RPC: even with SubtitleId set, it answers with
// SubtitleConfigs.SubtitleConfig[], so the matching item is selected client-side. An empty
// list means the template no longer exists, which is reported as NotFound so the Framework
// can remove it from state.
func (s *Service) DescribeLiveAISubtitle(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeLiveAISubtitle"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SubtitleId"] = id
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("live", "2016-11-01", action, query, request, true)
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
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.SubtitleConfigs.SubtitleConfig[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.SubtitleConfigs.SubtitleConfig[*]", response)
	}
	items, ok := v.([]interface{})
	if !ok || len(items) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("LiveAISubtitle", id), errs.NotFoundMsg, response)
	}
	for _, item := range items {
		if m, ok := item.(map[string]interface{}); ok {
			if sid, _ := m["SubtitleId"].(string); sid == id {
				return m, nil
			}
		}
	}
	return object, errs.WrapErrorf(errs.NotFoundErr("LiveAISubtitle", id), errs.NotFoundMsg, response)
}
