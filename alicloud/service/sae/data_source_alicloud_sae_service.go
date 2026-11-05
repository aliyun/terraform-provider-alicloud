package sae

import (
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlicloudSaeService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSaeServiceRead,

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
func dataSourceAlicloudSaeServiceRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("SaeServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}
	action := "OpenSaeService"
	request := map[string]*string{}
	conn, err := meta.(*connectivity.AliyunClient).NewTeaRoaCommonClient(connectivity.OpenSaeService)
	if err != nil {
		return errs.WrapError(err)
	}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err := conn.DoRequestWithAction(helper.StringPointer(action), helper.StringPointer("2019-05-06"), nil, helper.StringPointer("POST"), helper.StringPointer("AK"), helper.String("/service/open"), request, nil, nil, &util.RuntimeOptions{})
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
		if errs.IsExpectedErrors(err, []string{"ORDER.OPEND", "SYSTEM.SALE_VALIDATE_NO_SPECIFIC_CODE_FAILED"}) {
			d.SetId("SaeServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_sae_service", action, errs.AlibabaCloudSdkGoERROR)
	}
	d.SetId("SaeServiceHasBeenOpened")
	d.Set("status", "Opened")

	return nil
}
