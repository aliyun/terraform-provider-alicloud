package alicloud

import (
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudMaxcomputeService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudMaxcomputeServiceRead,

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
	action := "OpenMaxComputeService"
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"Region": client.RegionId,
	}
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := client.RpcPostWithEndpoint("MaxCompute", "2019-06-12", action, nil, request, false, connectivity.OpenMaxcomputeService)
		if err != nil {
			if IsExpectedErrors(err, []string{"QPS Limit Exceeded"}) || NeedRetry(err) {
				return resource.RetryableError(err)
			}
			addDebug(action, response, nil)
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, nil)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ORDER.OPEND"}) {
			d.SetId("MaxcomputeServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_maxcompute_service", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId("MaxcomputeServiceHasBeenOpened")
	d.Set("status", "Opened")

	return nil
}
