package alicloud

import (
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudAckService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAckServiceRead,

		Schema: map[string]*schema.Schema{
			"enable": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"On", "Off"}, false),
				Optional:     true,
				Default:      "Off",
			},
			"type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"propayasgo", "edgepayasgo", "gspayasgo"}, false),
				Required:     true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func dataSourceAlicloudAckServiceRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("AckServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}
	action := "OpenAckService"
	query := map[string]*string{
		"type": StringPointer(d.Get("type").(string)),
	}
	client := meta.(*connectivity.AliyunClient)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := client.RoaPostWithApiNameEndpoint("CS", "2015-12-15", action, "/service/open", query, nil, nil, false, connectivity.OpenAckService)
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
			d.SetId("AckServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ack_service", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId("AckServiceHasBeenOpened")
	d.Set("status", "Opened")

	return nil
}
