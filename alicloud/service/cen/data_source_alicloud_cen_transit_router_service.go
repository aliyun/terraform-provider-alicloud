package cen

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

func dataSourceAlicloudCenTransitRouterService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCenTransitRouterServiceRead,
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
func dataSourceAlicloudCenTransitRouterServiceRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("CenTransitRouterServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}

	client := meta.(*connectivity.AliyunClient)
	var err error

	var response map[string]interface{}
	request := map[string]interface{}{}
	action := "CheckTransitRouterService"
	request["ClientToken"] = helper.BuildClientToken(action)
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(3*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, false)
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
		if errs.IsExpectedErrors(err, []string{"USER_NOT_OPEN_TR_SERVICE", "Forbbiden.TransitRouterServiceNotOpen"}) {
			action = "OpenTransitRouterService"
			request["ClientToken"] = helper.BuildClientToken(action)
			err = retry.Retry(3*time.Minute, func() *retry.RetryError {
				response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, false)
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
				if errs.IsExpectedErrors(err, []string{"ORDER.OPEND"}) {
					d.SetId("CenTransitRouterServiceHasBeenOpened")
					d.Set("status", "Opened")
					return nil
				}
				return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cen_transit_router_service", action, errs.AlibabaCloudSdkGoERROR)
			}

			d.SetId("CenTransitRouterServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cen_transit_router_service", action, errs.AlibabaCloudSdkGoERROR)
	}

	if v, ok := response["Enabled"]; ok && v.(bool) {
		d.SetId("CenTransitRouterServiceHasBeenOpened")
		d.Set("status", "Opened")
		return nil
	}

	d.SetId("CenTransitRouterServiceHasNotBeenOpened")
	d.Set("status", "")
	return nil
}
