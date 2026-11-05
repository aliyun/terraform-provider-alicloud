package realtimecompute

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// foasPackageAPIVersion is the ROA API version of the FOAS (Realtime Compute)
// package management APIs. Per the ACube operation mapping the POP product code
// is "foas" (distinct from the SDKv2 realtimecompute "foasconsole" 2019-06-01
// RPC and "ververica" 2022-07-18 ROA gateways); the package lifecycle APIs are
// exposed on the legacy 2018-11-11 ROA version.
const foasPackageAPIVersion = "2018-11-11"

// DescribeRealtimeComputePackage calls foas/2018-11-11 GetPackage to refresh a
// single Package. The Terraform id is "<projectName>:<packageName>", matching
// the ACube composite unique key ProjectName:PackageName.
func (s *Service) DescribeRealtimeComputePackage(id string) (object map[string]interface{}, err error) {
	client := s.client
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return nil, errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	projectName := parts[0]
	packageName := parts[1]
	action := fmt.Sprintf("/api/v2/projects/%s/packages/%s", projectName, packageName)
	// Per the ACube operation mapping, RegionId is passed in the request header
	// (requestIn: header) on all package lifecycle APIs. The "foas" endpoint is
	// resolved dynamically by the location service, so the query map is empty.
	query := map[string]*string{}
	header := map[string]*string{
		"RegionId": helper.StringPointer(client.RegionId),
	}

	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("foas", foasPackageAPIVersion, action, query, header, nil)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, nil)
	if err != nil {
		// GetPackage reports a missing package as a request error; treat any
		// NotFound-shaped error or an empty response body as gone, so refresh
		// and the import lifecycle stay consistent.
		if errs.NotFoundError(err) {
			return nil, errs.WrapErrorf(errs.NotFoundErr("Package", id), errs.NotFoundMsg, response)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Package", response)
	if err != nil || v == nil {
		return nil, errs.WrapErrorf(errs.NotFoundErr("Package", id), errs.NotFoundMsg, response)
	}
	return v.(map[string]interface{}), nil
}
