package mscopensubscription

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

func DataSourceContactVerificationMessage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceContactVerificationMessageRead,
		Schema: map[string]*schema.Schema{
			"contact_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2}),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceContactVerificationMessageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "SendVerificationMessage"
	var response map[string]interface{}
	request := make(map[string]interface{})
	var err error
	request["ContactId"] = d.Get("contact_id")
	request["Type"] = d.Get("type")
	request["Locale"] = "en"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcGet("MscOpenSubscription", "2021-07-13", action, request, nil)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		d.SetId("SendMscSubContactVerificationMessage")
		d.Set("status", "Failed")
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_msc_sub_contact_verification_message", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId("SendMscSubContactVerificationMessage")
	d.Set("status", "Success")
	return nil
}
