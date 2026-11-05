package dataworks

import (
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// DescribeDataQualityAlertRule queries a DataWorks data quality alert rule
// by its numeric id. The id is stored in Terraform state as a string and is
// forwarded verbatim to the RPC form parameter; the API accepts the long
// value whether encoded as a JSON number or a numeric string.
//
// The Get response wraps the rule object under $.DataQualityAlertRule. An
// empty or missing object is reported as NotFound so the resource Read can
// drop the instance from state.
func (s *Service) DescribeDataQualityAlertRule(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "GetDataQualityAlertRule"
	request := map[string]interface{}{
		"Id":       id,
		"RegionId": client.RegionId,
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("dataworks-public", "2024-05-18", action, nil, request, true)
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

	v, err := jsonpath.Get("$.DataQualityAlertRule", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.DataQualityAlertRule", response)
	}
	object, ok := v.(map[string]interface{})
	if !ok || len(object) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("DataQualityAlertRule", id), errs.NotFoundWithResponse, response)
	}
	return object, nil
}
