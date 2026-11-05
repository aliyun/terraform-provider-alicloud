package esa

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	helper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// DescribeApiSecToken <<< Encapsulated get interface for Esa ApiSecToken.
//
// The composite id "SiteId:TokenId" is split into its identity parts.
// GetApiSecToken returns the TokenInfos array for the site; the matching
// token is selected by TokenId. ApiSecToken APIs are Private and their
// NotFound codes / response envelope are not fully enumerated in the frozen
// CloudSpec evidence (cloudspecCommit 5ff43bdc), so the array-find-no-match
// path is the authoritative NotFound signal, supplemented by the common ESA
// error codes ("SiteNotFound", "101") already used by sibling ESA readers.
func (s *EsaServiceV2) DescribeApiSecToken(id string) (object map[string]interface{}, err error) {
	client := s.client
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return object, fmt.Errorf("invalid ApiSecToken resource id %q: expected SiteId:TokenId", id)
	}
	siteId := parts[0]
	tokenId := parts[1]

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["SiteId"] = siteId
	query["RegionId"] = client.RegionId
	action := "GetApiSecToken"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("ESA", "2024-09-10", action, query, request)
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
		if errs.IsExpectedErrors(err, []string{"SiteNotFound", "101"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("ApiSecToken", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.TokenInfos[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TokenInfos[*]", response)
	}
	rows, ok := v.([]interface{})
	if !ok {
		return object, errs.WrapErrorf(errs.NotFoundErr("ApiSecToken", id), errs.NotFoundMsg, response)
	}
	for _, row := range rows {
		if item, ok := row.(map[string]interface{}); ok && fmt.Sprint(item["TokenId"]) == tokenId {
			return item, nil
		}
	}
	return object, errs.WrapErrorf(errs.NotFoundErr("ApiSecToken", id), errs.NotFoundMsg, response)
}

// DescribeApiSecToken >>> Encapsulated.
