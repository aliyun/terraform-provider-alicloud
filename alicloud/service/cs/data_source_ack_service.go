package cs

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

func DataSourceAckService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAckServiceRead,

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
func dataSourceAckServiceRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("AckServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}
	action := "OpenAckService"
	query := map[string]*string{
		"type": helper.StringPointer(d.Get("type").(string)),
	}
	conn, err := meta.(*connectivity.AliyunClient).NewTeaRoaCommonClient(connectivity.OpenAckService)
	if err != nil {
		return errs.WrapError(err)
	}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err := conn.DoRequestWithAction(helper.StringPointer(action), helper.StringPointer("2015-12-15"), nil, helper.StringPointer("POST"), helper.StringPointer("AK"), helper.String("/service/open"), query, nil, nil, &util.RuntimeOptions{})
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
		if errs.IsExpectedErrors(err, []string{"ORDER.OPEND"}) {
			d.SetId("AckServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_ack_service", action, errs.AlibabaCloudSdkGoERROR)
	}
	d.SetId("AckServiceHasBeenOpened")
	d.Set("status", "Opened")

	return nil
}
