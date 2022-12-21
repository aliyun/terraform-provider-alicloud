package alicloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudThreatDetectionLogShipper() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudThreatDetectionLogShipperRead,
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
			"open_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auth_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"buy_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sls_project_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sls_service_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudThreatDetectionLogShipperRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionService := ThreatDetectionService{client}
	var response map[string]interface{}
	action := "ModifyOpenLogShipper"
	request := make(map[string]interface{})
	conn, err := client.NewThreatdetectionClient()
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("ThreatDetectionLogShipperHasNotBeenOpened")
		d.Set("status", "")
		d.Set("open_status", "")
		d.Set("auth_status", "")
		d.Set("buy_status", "")
		d.Set("sls_project_status", "")
		d.Set("sls_service_status", "")
		return nil
	}

	request["From"] = "sas"

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-12-03"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_threat_detection_log_shipper", action, AlibabaCloudSdkGoERROR)
	}

	object, err := threatDetectionService.DescribeThreatDetectionLogShipper(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("open_status", object["OpenStatus"])
	d.Set("auth_status", object["AuthStatus"])
	d.Set("buy_status", object["BuyStatus"])
	d.Set("sls_project_status", object["SlsProjectStatus"])
	d.Set("sls_service_status", object["SlsServiceStatus"])

	d.SetId("ThreatDetectionLogShipperHasBeenOpened")

	d.Set("status", "Opened")

	return nil
}
