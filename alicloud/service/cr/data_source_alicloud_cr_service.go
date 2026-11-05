package cr

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

func dataSourceAlicloudCrService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCrServiceRead,

		Schema: map[string]*schema.Schema{
			"enable": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"On", "Off"}, false),
				Optional:     true,
				Default:      "Off",
			},
			"password": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringLenBetween(8, 32),
				Required:     true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func dataSourceAlicloudCrServiceRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("CrServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}
	action := "CreateUserInfo"
	request := map[string]interface{}{
		"User": map[string]interface{}{
			"Password": d.Get("password"),
		},
	}
	conn, err := meta.(*connectivity.AliyunClient).NewTeaRoaCommonClient(connectivity.OpenCrService)
	if err != nil {
		return errs.WrapError(err)
	}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err := conn.DoRequestWithAction(helper.StringPointer(action), helper.StringPointer("2016-06-07"), nil, helper.StringPointer("PUT"), helper.StringPointer("AK"), helper.String("/users"), nil, nil, request, &util.RuntimeOptions{})
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
			d.SetId("CrServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cr_service", action, errs.AlibabaCloudSdkGoERROR)
	}
	d.SetId("CrServiceHasBeenOpened")
	d.Set("status", "Opened")

	return nil
}
