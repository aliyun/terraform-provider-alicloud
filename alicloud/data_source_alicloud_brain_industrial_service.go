package alicloud

import (
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlicloudBrainIndustrialService() *schema.Resource {
	return &schema.Resource{
		Read:               dataSourceAlicloudBrainIndustrialServiceRead,
		DeprecationMessage: "This data source has been deprecated since v1.229.1 and will be removed in the future.",
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
func dataSourceAlicloudBrainIndustrialServiceRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("BrainIndustrialServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}
	action := "OpenCloudPidService"
	request := map[string]interface{}{}
	client := meta.(*connectivity.AliyunClient)
	var err error
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err := client.RpcPostWithEndpoint("brain-industrial", "2020-09-20", action, nil, request, false, connectivity.OpenBrainIndustrialService)
		if err != nil {
			if IsExpectedErrors(err, []string{"QPS Limit Exceeded"}) || NeedRetry(err) {
				return retry.RetryableError(err)
			}
			addDebug(action, response, nil)
			return retry.NonRetryableError(err)
		}
		addDebug(action, response, nil)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ORDER.OPEND"}) {
			d.SetId("BrainIndustrialServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_brain_industrial_service", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId("BrainIndustrialServiceHasBeenOpened")
	d.Set("status", "Opened")

	return nil
}
