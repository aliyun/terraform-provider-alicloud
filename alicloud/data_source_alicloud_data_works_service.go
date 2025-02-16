package alicloud

import (
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudDataWorksService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDataWorksServiceRead,

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
func dataSourceAlicloudDataWorksServiceRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("DataWorksServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}
	action := "OpenDataWorksStandardService"
	request := map[string]interface{}{
		"Region": meta.(*connectivity.AliyunClient).RegionId,
	}
	client := meta.(*connectivity.AliyunClient)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := client.RpcPostWithEndpoint("dataworks-public", "2018-06-01", action, nil, request, false, connectivity.DataWorksService)
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
			d.SetId("DataWorksServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dataworks_service", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId("DataWorksServiceHasBeenOpened")
	d.Set("status", "Opened")

	return nil
}
