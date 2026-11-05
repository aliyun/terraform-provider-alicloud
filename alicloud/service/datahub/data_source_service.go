package datahub

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

func DataSourceService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServiceRead,

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
func dataSourceServiceRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("DatahubServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}
	action := "OpenDataHubService"
	request := map[string]interface{}{}
	conn, err := meta.(*connectivity.AliyunClient).NewTeaCommonClient(connectivity.OpenDatahubService)

	if err != nil {
		return errs.WrapError(err)
	}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err := conn.DoRequest(helper.StringPointer(action), nil, helper.StringPointer("POST"), helper.StringPointer("2019-11-20"), helper.StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			d.SetId("DatahubServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_datahub_service", action, errs.AlibabaCloudSdkGoERROR)
	}
	d.SetId("DatahubServiceHasBeenOpened")
	d.Set("status", "Opened")

	return nil
}
