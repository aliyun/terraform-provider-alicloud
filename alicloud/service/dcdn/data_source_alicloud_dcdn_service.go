package dcdn

import (
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlicloudDcdnService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDcdnServiceRead,

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
func dataSourceAlicloudDcdnServiceRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("DcdnServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}
	action := "OpenDcdnService"
	request := map[string]interface{}{
		"BillType":          "PayByTraffic",
		"WebsocketBillType": "websockettraffic",
	}
	conn, err := meta.(*connectivity.AliyunClient).NewTeaCommonClient(connectivity.OpenDcdnService)
	if err != nil {
		return errs.WrapError(err)
	}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err := conn.DoRequest(helper.StringPointer(action), nil, helper.StringPointer("POST"), helper.StringPointer("2018-01-15"), helper.StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"QPS Limit Exceeded"}) || helper.NeedRetry(err) {
				return retry.RetryableError(err)
			}
			logs.AddDebug(action, response, nil)
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, nil)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"DcdnService.HasOpened", "ORDER.OPEND"}) {
			d.SetId("DcdnServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_dcdn_service", action, errs.AlibabaCloudSdkGoERROR)
	}
	d.SetId("DcdnServiceHasBeenOpened")
	d.Set("status", "Opened")

	return nil
}
