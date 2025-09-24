package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudPvtzService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudPvtzServiceRead,
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

func dataSourceAliCloudPvtzServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("PvtzServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}

	var response map[string]interface{}
	var endpoint string
	var err error

	action := "CreateInstance"

	request := map[string]interface{}{
		"ProductCode":       "pvtz",
		"ProductType":       "pvtzpost",
		"SubscriptionType":  "PayAsYouGo",
		"Parameter.1.Code":  "CommodityType",
		"Parameter.1.Value": "pvtz",
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, nil)

	if err != nil {
		if IsExpectedErrors(err, []string{"SYSTEM.SALE_VALIDATE_NO_SPECIFIC_CODE_FAILED", "ORDER.OPEND"}) {
			d.SetId("PvtzServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_pvtz_service", action, AlibabaCloudSdkGoERROR)
	}

	if response["Data"] != nil {
		d.SetId(fmt.Sprintf("%v", response["Data"].(map[string]interface{})["OrderId"]))
	} else {
		log.Printf("[ERROR] When opening pvtz service, invoking %s got an nil data. Response: %s.", action, response)
		d.SetId("PvtzServiceHasBeenOpened")
	}

	pvtzService := PvtzService{client}
	stateConf := BuildStateConf([]string{}, []string{"OPENED"}, 5*time.Minute, 5*time.Second, pvtzService.PvtzServiceStatusStateRefreshFunc([]string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	d.Set("status", "Opened")

	return nil
}
