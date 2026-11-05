package sslcertificateservice

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// NewCasService constructs the CAS (SSL Certificates Service) helper bound to
// the given AliyunClient. CasService itself is declared in service_cas.go;
// the constructor lives here so the CloudAccess resource can obtain one
// without touching the certificate-only file.
func NewCasService(client *connectivity.AliyunClient) *CasService {
	return &CasService{client: client}
}

// ListCloudAccess calls ListCloudAccess (cas/2020-04-07) and returns every
// CloudAccess entry. cloudName and secretId are optional server-side filters
// that narrow the result set; pass empty strings to list all entries. The
// method pages through the full result set automatically.
func (s *CasService) ListCloudAccess(cloudName, secretId string) ([]map[string]interface{}, error) {
	client := s.client
	var entries []map[string]interface{}
	currentPage := 1
	showSize := 50

	for {
		request := map[string]interface{}{
			"RegionId":    client.RegionId,
			"CurrentPage": currentPage,
			"ShowSize":    showSize,
		}
		if cloudName != "" {
			request["CloudName"] = cloudName
		}
		if secretId != "" {
			request["SecretId"] = secretId
		}

		action := "ListCloudAccess"
		var response map[string]interface{}
		var err error
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("cas", "2020-04-07", action, nil, request, true)
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
			return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, "ListCloudAccess", action, errs.AlibabaCloudSdkGoERROR)
		}

		v, err := jsonpath.Get("$.CloudAccessList", response)
		if err != nil {
			return nil, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, "ListCloudAccess", "$.CloudAccessList", response)
		}
		list, ok := v.([]interface{})
		if !ok || len(list) == 0 {
			break
		}
		for _, item := range list {
			if entry, ok := item.(map[string]interface{}); ok {
				entries = append(entries, entry)
			}
		}

		totalCount := formatInt(response["TotalCount"])
		if currentPage*showSize >= totalCount || len(list) < showSize {
			break
		}
		currentPage++
	}

	return entries, nil
}

// DescribeCloudAccess finds a single CloudAccess entry by accessId. cloudName
// and secretId are optional filters forwarded to ListCloudAccess to narrow
// the search; pass empty strings when the caller has no natural key besides
// the AccessId.
func (s *CasService) DescribeCloudAccess(accessId, cloudName, secretId string) (map[string]interface{}, error) {
	entries, err := s.ListCloudAccess(cloudName, secretId)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if formatAccessId(entry["AccessId"]) == accessId {
			return entry, nil
		}
	}
	return nil, errs.WrapErrorf(errs.NotFoundErr("SslCertificatesService:CloudAccess", accessId), errs.NotFoundWithResponse, nil)
}

// formatAccessId normalises the AccessId value returned by ListCloudAccess
// (which the tea-rpc client may decode as json.Number, float64, int or
// int64) into the string form stored in Terraform state.
func formatAccessId(v interface{}) string {
	switch n := v.(type) {
	case float64:
		return strconv.FormatInt(int64(n), 10)
	case json.Number:
		return n.String()
	case int:
		return strconv.Itoa(n)
	case int64:
		return strconv.FormatInt(n, 10)
	case string:
		return n
	}
	return ""
}

// formatInt normalises a numeric response field (TotalCount etc.) into an int.
func formatInt(v interface{}) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case json.Number:
		if i, err := n.Int64(); err == nil {
			return int(i)
		}
	case int:
		return n
	case int64:
		return int(n)
	case string:
		if i, err := strconv.Atoi(n); err == nil {
			return i
		}
	}
	return 0
}
