package alicloud

import (
	rpc "github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudLogService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudLogServiceRead,

		Schema: map[string]*schema.Schema{
			"enable": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"On", "Off"}, false),
				Optional:     true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func dataSourceAlicloudLogServiceRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("LogServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}

	client := meta.(*connectivity.AliyunClient)

	response, err := client.NewTeaCommonClientWithEndpoint(connectivity.OpenSlsService, func(teaClient *rpc.Client) (map[string]interface{}, error) {
		return teaClient.DoRequest(StringPointer("GetSlsService"), nil, StringPointer("POST"), StringPointer("2019-10-23"), StringPointer("AK"), nil, nil, &util.RuntimeOptions{})
	})
	addDebug("GetSlsService", response, nil)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_log_service", "GetLogService", AlibabaCloudSdkGoERROR)
	}

	if response["Success"] != nil && !response["Success"].(bool) {
		return WrapErrorf(fmt.Errorf("%s", response), DataDefaultErrorMsg, "alicloud_log_service", "GetLogService", AlibabaCloudSdkGoERROR)
	}
	if response["Enabled"] != nil && response["Enabled"].(bool) {
		d.SetId("LogServiceHasBeenOpened")
		d.Set("status", "Opened")
		return nil
	}

	response, err = client.NewTeaCommonClientWithEndpoint(connectivity.OpenSlsService, func(teaClient *rpc.Client) (map[string]interface{}, error) {
		return teaClient.DoRequest(StringPointer("OpenSlsService"), nil, StringPointer("POST"), StringPointer("2019-10-23"), StringPointer("AK"), nil, nil, &util.RuntimeOptions{})
	})
	addDebug("OpenSlsService", response, nil)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_log_service", "OpenLogService", AlibabaCloudSdkGoERROR)
	}

	if response["Success"] != nil && !response["Success"].(bool) {
		return WrapErrorf(fmt.Errorf("%s", response), DataDefaultErrorMsg, "alicloud_log_service", "OpenLogService", AlibabaCloudSdkGoERROR)
	}
	d.SetId("LogServiceHasBeenOpened")
	d.Set("status", "Opened")

	return nil
}
