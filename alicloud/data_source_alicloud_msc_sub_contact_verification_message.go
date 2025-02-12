package alicloud

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudMscSubContactVerificationMessage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudMscSubContactVerificationMessageRead,
		Schema: map[string]*schema.Schema{
			"contact_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2}),
				ForceNew:     true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudMscSubContactVerificationMessageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "SendVerificationMessage"
	var response map[string]interface{}
	request := make(map[string]interface{})
	var err error
	request["ContactId"] = d.Get("contact_id")
	request["Type"] = d.Get("type")
	request["Locale"] = "en"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcGet("MscOpenSubscription", "2021-07-13", action, request, nil)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		d.SetId("SendMscSubContactVerificationMessage")
		d.Set("status", "Failed")
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_msc_sub_contact_verification_message", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId("SendMscSubContactVerificationMessage")
	d.Set("status", "Success")
	return nil
}
