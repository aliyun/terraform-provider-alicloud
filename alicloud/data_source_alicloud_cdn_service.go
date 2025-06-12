package alicloud

import (
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"fmt"
)

func dataSourceAliCloudCdnService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCdnServiceRead,

		Schema: map[string]*schema.Schema{
			"enable": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Off",
				ValidateFunc: StringInSlice([]string{"On", "Off"}, false),
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "PayByTraffic",
				ValidateFunc: StringInSlice([]string{"PayByTraffic", "PayByBandwidth"}, false),
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
		return WrapError(err)
	}
	response, err := conn.DoRequest(StringPointer("DescribeCdnService"), nil, StringPointer("POST"), StringPointer("2018-05-10"), StringPointer("AK"), nil, nil, &util.RuntimeOptions{})
	addDebug("DescribeCdnService", response, nil)
	if err != nil {
		if IsExpectedErrors(err, []string{"OperationDenied", "CdnServiceNotFound"}) {
			log.Printf("[DEBUG] Datasource alicloud_cdn_service DescribeCdnService Failed!!! %s", err)
		} else {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cdn_service", "DescribeCdnService", AlibabaCloudSdkGoERROR)
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
			return WrapError(fmt.Errorf("Field 'internet_charge_type' is required when 'enable' is 'On'."))
		}
		requestBody := map[string]interface{}{"InternetChargeType": chargeType}

		isUpdateChargeType := false
		checkAndUpdate := func(key string) bool {
			if v, ok := response[key]; ok && fmt.Sprint(v) != "" {
				return chargeType != fmt.Sprint(v)
			}
			return false
		}

		if checkAndUpdate("ChangingChargeType") || checkAndUpdate("InternetChargeType") {
			isUpdateChargeType = true
		}

		if opened && isUpdateChargeType {
			resp, err := conn.DoRequest(StringPointer("ModifyCdnService"), nil, StringPointer("POST"), StringPointer("2018-05-10"), StringPointer("AK"), nil, requestBody, &util.RuntimeOptions{})

			addDebug("ModifyCdnService", resp, nil)
			if err != nil {
				return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cdn_service", "ModifyCdnService", AlibabaCloudSdkGoERROR)
			}
		}

		if !opened {
			resp, err := conn.DoRequest(StringPointer("OpenCdnService"), nil, StringPointer("POST"), StringPointer("2018-05-10"), StringPointer("AK"), nil, requestBody, &util.RuntimeOptions{})
			addDebug("OpenCdnService", resp, nil)
			if err != nil {
				if IsExpectedErrors(err, []string{"CdnService.HasOpened"}) {
					log.Printf("[DEBUG] Datasource alicloud_cdn_service has opened!!!")
					opened = true
				} else {
					return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cdn_service", "OpenCdnService", AlibabaCloudSdkGoERROR)
				}
			} else {
				opened = true
			}
		}

		action := "DescribeCdnService"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(4*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-10"), StringPointer("AK"), nil, nil, &runtime)
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"CdnServiceNotFound"}) {
					wait()
					return resource.RetryableError(err)
				}
				if IsExpectedErrors(err, []string{"OperationDenied"}) {
					log.Printf("[DEBUG] Datasource alicloud_cdn_service DescribeCdnService Failed!!! %s", err)
					return nil
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, nil)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cdn_service", action, AlibabaCloudSdkGoERROR)
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
