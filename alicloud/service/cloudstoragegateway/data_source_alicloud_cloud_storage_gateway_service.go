package cloudstoragegateway

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

func DataSourceAlicloudCloudStorageGatewayService() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceAlicloudCloudStorageGatewayServiceRead,

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
func DataSourceAlicloudCloudStorageGatewayServiceRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("CloudStorageGatewayServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}
	action := "OpenSgwService"
	conn, err := meta.(*connectivity.AliyunClient).NewTeaCommonClient(connectivity.OpenCloudStorageGatewayService)
	if err != nil {
		return errs.WrapError(err)
	}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err := conn.DoRequest(helper.StringPointer(action), nil, helper.StringPointer("POST"), helper.StringPointer("2018-05-11"), helper.StringPointer("AK"), nil, nil, &util.RuntimeOptions{})
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
		if errs.IsExpectedErrors(err, []string{"OrderOpend"}) {
			d.SetId("CloudStorageGatewayServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cloud_storage_gateway_service", action, errs.AlibabaCloudSdkGoERROR)
	}
	d.SetId("CloudStorageGatewayServiceHasBeenOpened")
	d.Set("status", "Opened")

	return nil
}
