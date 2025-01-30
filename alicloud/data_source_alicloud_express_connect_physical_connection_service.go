package alicloud

import (
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudExpressConnectPhysicalConnectionService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudPhysicalConnectionServiceRead,

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
func dataSourceAlicloudPhysicalConnectionServiceRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("PhysicalConnectionServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}
	client := meta.(*connectivity.AliyunClient)
	action := "OpenPhysicalConnectionService"
	var response map[string]interface{}
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				return resource.RetryableError(err)
			}
			addDebug(action, response, nil)
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, nil)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"PURCHASE_QUANTITY_LIMIT", "InvalidOperation.OrderOpened"}) {
			d.SetId("PhysicalConnectionServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_fc_service", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId("PhysicalConnectionServiceHasBeenOpened")
	d.Set("status", "Opened")

	return nil
}
