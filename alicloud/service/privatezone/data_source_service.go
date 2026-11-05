package privatezone

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServiceRead,
		Schema: map[string]*schema.Schema{
			"enable": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Off",
				ValidateFunc: validators.StringInSlice([]string{"On", "Off"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceServiceRead(d *schema.ResourceData, meta interface{}) error {
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

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(3*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			if !client.IsInternationalAccount() && errs.IsExpectedErrors(err, []string{"NotApplicable"}) {
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, nil)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"SYSTEM.SALE_VALIDATE_NO_SPECIFIC_CODE_FAILED", "ORDER.OPEND"}) {
			d.SetId("PvtzServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_pvtz_service", action, errs.AlibabaCloudSdkGoERROR)
	}

	if response["Data"] != nil {
		d.SetId(fmt.Sprintf("%v", response["Data"].(map[string]interface{})["OrderId"]))
	} else {
		log.Printf("[ERROR] When opening pvtz service, invoking %s got an nil data. Response: %s.", action, response)
		d.SetId("PvtzServiceHasBeenOpened")
	}

	pvtzService := Service{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"OPENED"}, 5*time.Minute, 5*time.Second, pvtzService.ServiceStatusStateRefreshFunc([]string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	d.Set("status", "Opened")

	return nil
}
