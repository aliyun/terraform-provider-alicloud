package alicloud

import (
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlicloudMaxcomputeService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudMaxcomputeServiceRead,
		DeprecationMessage: "This data source has been deprecated since v1.228.0. The OpenMaxComputeService API it relied on " +
			"has been decommissioned, so it can no longer activate the MaxCompute service and only checks whether the service " +
			"is activated. Please use the resource alicloud_max_compute_quota instead.",

		Schema: map[string]*schema.Schema{
			"enable": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"On", "Off"}, false),
				Optional:     true,
				Default:      "Off",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func dataSourceAlicloudMaxcomputeServiceRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("MaxcomputeServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}
	// The legacy OpenMaxComputeService action (version 2019-06-12) has been
	// decommissioned and now returns 404 InvalidAction.NotFound, so activation
	// status is detected instead: activating MaxCompute always creates at least
	// the default pay-as-you-go quota, so a non-empty quota list means the
	// service has been activated in this region.
	client := meta.(*connectivity.AliyunClient)
	action := "/api/v1/quotas"
	var response map[string]interface{}
	query := make(map[string]*string)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RoaGet("MaxCompute", "2022-01-04", action, query, nil, nil)
		if err != nil {
			if IsExpectedErrors(err, []string{"QPS Limit Exceeded"}) || NeedRetry(err) {
				return resource.RetryableError(err)
			}
			addDebug(action, resp, nil)
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, resp, nil)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_maxcompute_service", action, AlibabaCloudSdkGoERROR)
	}
	quotaList, err := jsonpath.Get("$.data.quotaInfoList", response)
	if err == nil && quotaList != nil {
		if quotas, ok := quotaList.([]interface{}); ok && len(quotas) > 0 {
			d.SetId("MaxcomputeServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
	}
	return WrapErrorf(NotFoundErr("MaxcomputeService", "MaxcomputeServiceHasNotBeenOpened"),
		"the MaxCompute service has not been activated in this region, and it can no longer be activated through this data source "+
			"because the OpenMaxComputeService API has been decommissioned; please activate MaxCompute in the Alibaba Cloud console")
}
