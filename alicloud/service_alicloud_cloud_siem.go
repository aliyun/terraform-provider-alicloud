package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// CloudSiemService wraps the connectivity client for ThreatDetection Vendor
// resources that are backed by the cloud-siem/2024-12-12 RPC API.
type CloudSiemService struct {
	client *connectivity.AliyunClient
}

// DescribeThreatDetectionVendor fetches a single vendor by id via ListVendors.
// ListVendors does not honor the VendorIds filter (verified against the live
// API: the full vendor list is always returned), so the match is done
// client-side by exact VendorId.
// ListVendors requires Lang; pass the configured lang, defaulting to "en".
func (s *CloudSiemService) DescribeThreatDetectionVendor(id string, lang ...string) (object map[string]interface{}, err error) {
	client := s.client

	langValue := "en"
	if len(lang) > 0 && lang[0] != "" {
		langValue = lang[0]
	}
	request := map[string]interface{}{
		"Lang": langValue,
	}

	var response map[string]interface{}
	action := "ListVendors"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPost("cloud-siem", "2024-12-12", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	vendorsRaw, err := jsonpath.Get("$.Vendors", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Vendors", response)
	}
	vendors, ok := vendorsRaw.([]interface{})
	if !ok {
		return object, WrapErrorf(fmt.Errorf("ThreatDetection vendor %s is not found", id), NotFoundWithResponse, response)
	}
	for _, v := range vendors {
		vendor, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		if fmt.Sprint(vendor["VendorId"]) == id {
			return vendor, nil
		}
	}
	return object, WrapErrorf(fmt.Errorf("ThreatDetection vendor %s is not found", id), NotFoundWithResponse, response)
}
