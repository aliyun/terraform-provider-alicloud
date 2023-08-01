package alicloud

import (
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudVpcFlowLogService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudVpcFlowLogServiceRead,
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
func dataSourceAliCloudVpcFlowLogServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}

	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("VpcFlowLogServiceHasNotBeenOpened")
		d.Set("status", "")

		return nil
	}

	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}

	action := "GetFlowLogServiceStatus"
	getFlowLogServiceStatusReq := map[string]interface{}{
		"RegionId":    client.RegionId,
		"ClientToken": buildClientToken("GetFlowLogServiceStatus"),
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, getFlowLogServiceStatusReq, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, getFlowLogServiceStatusReq)

	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_vpc_flow_log_service", action, AlibabaCloudSdkGoERROR)
	}

	if v, ok := response["Enabled"]; ok && v.(bool) {
		d.SetId("VpcFlowLogServiceHasBeenOpened")
		d.Set("status", "Opened")

		return nil
	}

	action = "OpenFlowLogService"
	openFlowLogServiceReq := map[string]interface{}{
		"RegionId":    client.RegionId,
		"ClientToken": buildClientToken("OpenFlowLogService"),
	}

	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, openFlowLogServiceReq, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, openFlowLogServiceReq)

	if err != nil {
		if IsExpectedErrors(err, []string{"OperationFailed.ExceedPurchaseLimit"}) {
			d.SetId("VpcFlowLogServiceHasBeenOpened")
			d.Set("status", "Opened")

			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_vpc_flow_log_service", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId("VpcFlowLogServiceHasBeenOpened")
	d.Set("status", "Opened")

	return nil
}
