package hbr

import (
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceService() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceServiceRead,
		Schema: map[string]*schema.Schema{
			"enable": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Off",
				ValidateFunc: validation.StringInSlice([]string{"On", "Off"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DataSourceServiceRead(d *schema.ResourceData, meta interface{}) error {
	action := "OpenHbrService"
	request := map[string]interface{}{}
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("HbrServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}

	client := meta.(*connectivity.AliyunClient)
	err := retry.Retry(3*time.Minute, func() *retry.RetryError {
		response, err := client.RpcPostWithEndpoint("hbr", "2017-09-08", action, nil, request, false, connectivity.OpenHbrService)
		if err != nil {
			if helper.NeedRetry(err) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, nil)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ORDER.OPEND"}) {
			d.SetId("HbrServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_hbr_service", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId("HbrServiceHasBeenOpened")
	d.Set("status", "Opened")

	return nil
}
