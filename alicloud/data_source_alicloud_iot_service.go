package alicloud

import (
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlicloudIotService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudIotServiceRead,

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
func dataSourceAlicloudIotServiceRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("IotServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}
	action := "OpenIotService"
	request := map[string]interface{}{}
	client := meta.(*connectivity.AliyunClient)
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err := client.RpcPostWithEndpoint("Iot", "2018-01-20", action, nil, request, false, connectivity.OpenIotService)
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
			d.SetId("IotServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_iot_service", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId("IotServiceHasBeenOpened")
	d.Set("status", "Opened")

	return nil
}
