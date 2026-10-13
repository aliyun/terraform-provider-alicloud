package alicloud

import (
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlicloudCmsService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCmsServiceRead,

		DeprecationMessage: "This data source has been deprecated since v1.219.0 and will be removed in the future.",

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
func dataSourceAlicloudCmsServiceRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("CmsServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}
	action := "OpenCmsService"
	request := map[string]interface{}{}
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
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
		if IsExpectedErrors(err, []string{"ORDER.OPEND", "Has.effect.suit"}) {
			d.SetId("CmsServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cms_service", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId("CmsServiceHasBeenOpened")
	d.Set("status", "Opened")

	return nil
}
