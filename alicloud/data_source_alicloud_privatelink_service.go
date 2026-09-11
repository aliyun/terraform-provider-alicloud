package alicloud

import (
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudPrivateLinkService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudPrivateLinkServiceRead,
		Schema: map[string]*schema.Schema{
			"enable": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Off",
				ValidateFunc: StringInSlice([]string{"On", "Off"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAliCloudPrivateLinkServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("PrivateLinkServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}

	var response map[string]interface{}
	var err error

	action := "OpenPrivateLinkService"
	request := map[string]interface{}{}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Privatelink", "2020-04-15", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"QPS Limit Exceeded"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, nil)

	if err != nil {
		if IsExpectedErrors(err, []string{"OrderOpend"}) {
			d.SetId("PrivateLinkServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_privatelink_service", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId("PrivateLinkServiceHasBeenOpened")

	d.Set("status", "Opened")

	return nil
}
