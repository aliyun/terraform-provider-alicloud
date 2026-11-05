package live

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"time"
)

// DescribeLiveStreamMerge <<< Encapsulated get interface for Live StreamMerge.
//
// DescribeLiveStreamMerge is a List-shaped Read: it returns
// LiveStreamMergeList.LiveStreamMerge[] filtered by the optional AppName,
// StreamName and Protocol query parameters. The caller passes the resource
// identity (DomainName + AppName + StreamName) plus the configured Protocol so
// the server narrows the result; the method still matches StreamName on the
// client to be robust against a non-exhaustive server filter, and reports
// NotFound when no entry remains.
func (s *Service) DescribeLiveStreamMerge(domainName, appName, streamName, protocol string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DomainName"] = domainName
	request["AppName"] = appName
	request["StreamName"] = streamName
	request["RegionId"] = client.RegionId
	if protocol != "" {
		request["Protocol"] = protocol
	}
	action := "DescribeLiveStreamMerge"

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
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, fmt.Sprintf("%s:%s:%s", domainName, appName, streamName), action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.LiveStreamMergeList.LiveStreamMerge[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, fmt.Sprintf("%s:%s:%s", domainName, appName, streamName), "$.LiveStreamMergeList.LiveStreamMerge[*]", response)
	}

	items, ok := v.([]interface{})
	if !ok || len(items) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("LiveStreamMerge", fmt.Sprintf("%s:%s:%s", domainName, appName, streamName)), errs.NotFoundMsg, response)
	}

	for _, row := range items {
		item, ok := row.(map[string]interface{})
		if !ok {
			continue
		}
		if fmt.Sprint(item["StreamName"]) == streamName && fmt.Sprint(item["AppName"]) == appName {
			return item, nil
		}
	}

	return object, errs.WrapErrorf(errs.NotFoundErr("LiveStreamMerge", fmt.Sprintf("%s:%s:%s", domainName, appName, streamName)), errs.NotFoundMsg, response)
}

// DescribeLiveStreamMerge >>> Encapsulated.
