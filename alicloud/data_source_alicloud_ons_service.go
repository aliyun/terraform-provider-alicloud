package alicloud

import (
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlicloudOnsService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudOnsServiceRead,

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
func dataSourceAlicloudOnsServiceRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("OnsServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}
	action := "OpenOnsService"
	request := map[string]interface{}{}
	client := meta.(*connectivity.AliyunClient)
	var err error
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err := client.RpcPostWithEndpoint("Ons", "2019-02-14", action, nil, request, false, connectivity.OpenOnsService)
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
		if IsExpectedErrors(err, []string{"OrderOpend"}) {
			d.SetId("OnsServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ons_service", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId("OnsServiceHasBeenOpened")
	d.Set("status", "Opened")

	return nil
}
