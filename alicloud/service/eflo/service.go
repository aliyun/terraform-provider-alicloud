package eflo

import (
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// Service wraps the Eflo API calls used by the data sources of this package.
// Eflo (灵骏) is the managed HPC network service; the RPC product/version pair
// is eflo/2022-05-30 for every call here.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor: the client field stays unexported, so
// the data source in the same package builds values with it.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

// defaultErRouteEntriesPageSize mirrors PageSizeLarge (50) used by the SDKv2
// eflo data sources when paging ListErRouteEntries.
const defaultErRouteEntriesPageSize = 50

// ListErRouteEntries pages through the ListErRouteEntries API for the Er named
// by request["ErId"]. When singlePage is true only the page named by
// request["PageNumber"]/request["PageSize"] is fetched, matching the legacy
// page_number single-page mode of the SDKv2 eflo data sources. The response
// wraps the entry list under $.Content.Data and the total under $.Content.Total.
func (s *Service) ListErRouteEntries(request map[string]interface{}, singlePage bool) (objects []map[string]interface{}, total int, err error) {
	client := s.client
	const action = "ListErRouteEntries"
	objects = make([]map[string]interface{}, 0)

	pageSize := defaultErRouteEntriesPageSize
	if v, ok := request["PageSize"]; ok {
		if ps, ok2 := v.(int); ok2 && ps > 0 {
			pageSize = ps
		}
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	for {
		var response map[string]interface{}
		rerr := retry.Retry(5*time.Minute, func() *retry.RetryError {
			resp, e := client.RpcPost("eflo", "2022-05-30", action, nil, request, true)
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
		logs.AddDebug(action, response, request)
		if rerr != nil {
			return objects, total, errs.WrapErrorf(rerr, errs.DataDefaultErrorMsg, "alicloud_eflo_er_route_entries", action, errs.AlibabaCloudSdkGoERROR)
		}

		if totalRaw, jerr := jsonpath.Get("$.Content.Total", response); jerr == nil {
			if tv, ok := totalRaw.(float64); ok {
				total = int(tv)
			}
		}
		dataRaw, derr := jsonpath.Get("$.Content.Data", response)
		if derr != nil {
			return objects, total, errs.WrapErrorf(derr, errs.FailedGetAttributeMsg, "alicloud_eflo_er_route_entries", "$.Content.Data", response)
		}
		list, _ := dataRaw.([]interface{})
		for _, v := range list {
			if item, ok := v.(map[string]interface{}); ok {
				objects = append(objects, item)
			}
		}
		if singlePage || len(list) < pageSize {
			break
		}
		pageNumber := 1
		if v, ok := request["PageNumber"]; ok {
			if n, ok2 := v.(int); ok2 {
				pageNumber = n
			}
		}
		request["PageNumber"] = pageNumber + 1
	}
	return objects, total, nil
}
