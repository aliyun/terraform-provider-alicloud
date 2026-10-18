package cms

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// Service wraps the CMS API calls used by the data sources of this package.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor: the client field stays unexported so
// callers outside the package (including acceptance tests) build values with it.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

// ListAlertMetricGroups queries the CMS /alertMetricGroups endpoint with cursor
// pagination (maxResults + nextToken), retrying transient errors per page, and
// returns the raw group objects from `data`. Client-side filtering (e.g. ids)
// is the caller's responsibility.
func (s *Service) ListAlertMetricGroups(datasourceType string, includeDetails bool) ([]map[string]interface{}, error) {
	client := s.client
	const action = "/alertMetricGroups"

	query := map[string]*string{
		"maxResults": helper.StringPointer(strconv.Itoa(helper.PageSizeLarge)),
	}
	if datasourceType != "" {
		query["datasourceType"] = helper.StringPointer(datasourceType)
	}
	if includeDetails {
		query["includeDetails"] = helper.StringPointer(fmt.Sprint(includeDetails))
	}

	var objects []map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	for {
		var response map[string]interface{}
		err := retry.Retry(5*time.Minute, func() *retry.RetryError {
			resp, e := client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)
			if e != nil {
				if helper.NeedRetry(e) {
					wait()
					return retry.RetryableError(e)
				}
				return retry.NonRetryableError(e)
			}
			response = resp
			return nil
		})
		logs.AddDebug(action, response, query)
		if err != nil {
			return objects, errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cms_alert_metric_groups", action, errs.AlibabaCloudSdkGoERROR)
		}

		data, _ := jsonpath.Get("$.data[*]", response)
		if items, ok := data.([]interface{}); ok {
			for _, v := range items {
				if item, ok := v.(map[string]interface{}); ok {
					objects = append(objects, item)
				}
			}
		}

		nextToken, ok := response["nextToken"].(string)
		if !ok || nextToken == "" {
			break
		}
		query["nextToken"] = helper.StringPointer(nextToken)
	}
	return objects, nil
}
