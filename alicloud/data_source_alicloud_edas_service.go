package alicloud

import (
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"fmt"

	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudEdasService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEdasServiceRead,

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
func dataSourceAlicloudEdasServiceRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("EdasServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}

	request := map[string]interface{}{
		"ProductCode":       "edas",
		"SubscriptionType":  "PayAsYouGo",
		"ProductType":       "edaspostpay",
		"Parameter.1.Code":  "env",
		"Parameter.1.Value": "env_public",
	}
	conn, err := meta.(*connectivity.AliyunClient).NewTeaCommonClient(connectivity.OpenBssService)
	if err != nil {
		return WrapError(err)
	}
	response, err := conn.DoRequest(StringPointer("CreateInstance"), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

	addDebug("CreateInstance", response, nil)
	if err != nil {
		if IsExpectedErrors(err, []string{"SYSTEM.SALE_VALIDATE_NO_SPECIFIC_CODE_FAILED"}) {
			d.SetId("EdasServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_edas_service", "CreateInstance", AlibabaCloudSdkGoERROR)
	}

	if response["Success"] != nil && !response["Success"].(bool) {
		if response["Code"] != nil && response["Code"].(string) == "SYSTEM.SALE_VALIDATE_NO_SPECIFIC_CODE_FAILED" {
			d.SetId("EdasServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return WrapErrorf(fmt.Errorf("%s", response), DataDefaultErrorMsg, "alicloud_edas_service", "CreateInstance", AlibabaCloudSdkGoERROR)
	}

	if response["Data"] != nil {
		d.SetId(fmt.Sprintf("%v", response["Data"].(map[string]interface{})["OrderId"]))
	} else {
		log.Printf("[ERROR] When opening EDAS service, invoking CreateInstance got an nil data. Response: %s.", response)
		d.SetId("EdasServiceHasBeenOpened")
	}
	d.Set("status", "Opened")

	return nil
}
