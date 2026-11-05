package cdn

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAliCloudCdnService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCdnServiceRead,

		Schema: map[string]*schema.Schema{
			"enable": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Off",
				ValidateFunc: validators.StringInSlice([]string{"On", "Off"}, false),
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "PayByTraffic",
				ValidateFunc: validators.StringInSlice([]string{"PayByTraffic", "PayByBandwidth"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"opening_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"changing_charge_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"changing_affect_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func dataSourceAliCloudCdnServiceRead(d *schema.ResourceData, meta interface{}) error {
	opened := false
	enable := ""
	if v, ok := d.GetOk("enable"); ok {
		enable = v.(string)
	}

	conn, err := meta.(*connectivity.AliyunClient).NewTeaCommonClient(connectivity.OpenCdnService)
	if err != nil {
		return errs.WrapError(err)
	}
	response, err := conn.DoRequest(helper.StringPointer("DescribeCdnService"), nil, helper.StringPointer("POST"), helper.StringPointer("2018-05-10"), helper.StringPointer("AK"), nil, nil, &util.RuntimeOptions{})
	logs.AddDebug("DescribeCdnService", response, nil)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"OperationDenied", "CdnServiceNotFound"}) {
			log.Printf("[DEBUG] Datasource alicloud_cdn_service DescribeCdnService Failed!!! %s", err)
		} else {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cdn_service", "DescribeCdnService", errs.AlibabaCloudSdkGoERROR)
		}
	}

	if response["OpeningTime"] != nil && response["OpeningTime"].(string) != "" {
		opened = true
	}

	if enable == "On" {
		chargeType := ""
		if v, ok := d.GetOk("internet_charge_type"); ok {
			chargeType = v.(string)
		}
		if chargeType == "" {
			return errs.WrapError(fmt.Errorf("Field 'internet_charge_type' is required when 'enable' is 'On'."))
		}
		requestBody := map[string]interface{}{"InternetChargeType": chargeType}

		if !opened {
			resp, err := conn.DoRequest(helper.StringPointer("OpenCdnService"), nil, helper.StringPointer("POST"), helper.StringPointer("2018-05-10"), helper.StringPointer("AK"), nil, requestBody, &util.RuntimeOptions{})
			logs.AddDebug("OpenCdnService", resp, nil)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"CdnService.HasOpened"}) {
					log.Printf("[DEBUG] Datasource alicloud_cdn_service has opened!!!")
					opened = true
				} else {
					return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cdn_service", "OpenCdnService", errs.AlibabaCloudSdkGoERROR)
				}
			} else {
				opened = true
			}
		}

		action := "DescribeCdnService"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(4*time.Minute, func() *retry.RetryError {
			response, err = conn.DoRequest(helper.StringPointer(action), nil, helper.StringPointer("POST"), helper.StringPointer("2018-05-10"), helper.StringPointer("AK"), nil, nil, &runtime)
			if err != nil {
				if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"CdnServiceNotFound"}) {
					wait()
					return retry.RetryableError(err)
				}
				if errs.IsExpectedErrors(err, []string{"OperationDenied"}) {
					log.Printf("[DEBUG] Datasource alicloud_cdn_service DescribeCdnService Failed!!! %s", err)
					return nil
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, nil)
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cdn_service", action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	if opened {
		d.SetId("CdnServiceHasBeenOpened")
		d.Set("status", "Opened")
	} else {
		d.SetId("CdnServiceHasNotBeenOpened")
		d.Set("status", "")
	}

	d.Set("internet_charge_type", response["InternetChargeType"])
	d.Set("opening_time", response["OpeningTime"])
	d.Set("changing_charge_type", response["ChangingChargeType"])
	d.Set("changing_affect_time", response["ChangingAffectTime"])

	return nil
}
