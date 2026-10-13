package alicloud

import (
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
	conn, err := meta.(*connectivity.AliyunClient).NewTeaRoaCommonClient(connectivity.OpenAckService)
	if err != nil {
		return WrapError(err)
	}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err := conn.DoRequestWithAction(StringPointer(action), StringPointer("2015-12-15"), nil, StringPointer("POST"), StringPointer("AK"), String("/service/open"), query, nil, nil, &util.RuntimeOptions{})
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
