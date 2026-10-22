// Package esa. This file is generated automatically. Please do not modify it manually, thank you!
package esa

import (
	tfhelper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	tflogs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	tferrs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"

	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfretry "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type EsaServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeEsaImageTransform <<< Encapsulated get interface for Esa ImageTransform.

func (s *EsaServiceV2) DescribeEsaImageTransform(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = tferrs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
		return nil, err
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ConfigId"] = parts[1]
	query["SiteId"] = parts[0]

	action := "GetImageTransform"

	wait := tfhelper.IncrementalWait(3*time.Second, 5*time.Second)
	err = tfretry.Retry(1*time.Minute, func() *tfretry.RetryError {
		response, err = client.RpcGet("ESA", "2024-09-10", action, query, request)

		if err != nil {
			if tferrs.IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests"}) || tfhelper.NeedRetry(err) {
				wait()
				return tfretry.RetryableError(err)
			}
			return tfretry.NonRetryableError(err)
		}
		return nil
	})
	tflogs.AddDebug(action, response, request)
	if err != nil {
		return object, tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, id, action, tferrs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

// DescribeEsaImageTransform >>> Encapsulated.
