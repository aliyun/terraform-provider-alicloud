package alicloud

import (
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudGovernanceService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudGovernanceServiceRead,
		Schema: map[string]*schema.Schema{
			"enable": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Off",
				ValidateFunc: validation.StringInSlice([]string{"On", "Off"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAliCloudGovernanceServiceRead(d *schema.ResourceData, meta interface{}) error {
	action := "OpenGovernanceService"
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("GovernanceServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}

	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{}
	query := map[string]interface{}{}
	query["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		response, err := client.RpcPost("governance", "2021-01-20", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"Order.Opend"}) {
			d.SetId("GovernanceServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_governance_service", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId("GovernanceServiceHasBeenOpened")
	d.Set("status", "Opened")

	return nil
}
