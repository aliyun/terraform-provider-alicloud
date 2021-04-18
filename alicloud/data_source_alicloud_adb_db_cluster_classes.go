package alicloud

import (
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudAdbDbClusterClasses() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAdbDbClusterClassesRead,

		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "PayAsYouGo",
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"available_zone_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"classes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudAdbDbClusterClassesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeAvailableResource"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v.(string)
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["ChargeType"] = convertAdbDbClusterDBClusterPaymentTypeRequest(v.(string))
	}

	response := make(map[string]interface{})
	conn, err := client.NewAdsClient()
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_adb_db_cluster_classes", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	resp, err := jsonpath.Get("$.AvailableZoneList", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.AvailableZoneList", response)
	}

	result := resp.([]interface{})
	objects := make([]map[string]interface{}, 0)
	zoneIds := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range result {
		item := v.(map[string]interface{})
		objects = append(objects, item)
	}

	for _, object := range objects {
		classes := make([]string, 0)
		mapping := map[string]interface{}{
			"zone_id": object["ZoneId"],
		}
		zoneIds = append(zoneIds, object["ZoneId"].(string))
		if supportedMode, ok := object["SupportedMode"]; ok {
			if supportedModes, ok := supportedMode.([]interface{}); ok {
				for _, supportedMode := range supportedModes {
					if supportedSerialList, ok := supportedMode.(map[string]interface{})["SupportedSerialList"]; ok {
						if supportedSerials, ok := supportedSerialList.([]interface{}); ok {
							for _, supportedSerial := range supportedSerials {
								if supportedInstanceClassList, ok := supportedSerial.(map[string]interface{})["SupportedInstanceClassList"]; ok {
									if supportedInstanceClasses, ok := supportedInstanceClassList.([]interface{}); ok {
										for _, supportedInstanceClass := range supportedInstanceClasses {
											classes = append(classes, supportedInstanceClass.(map[string]interface{})["InstanceClass"].(string))
										}
									}
								}
							}
						}
					}
				}
			}
		}
		mapping["classes"] = classes
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(zoneIds))
	if err := d.Set("available_zone_list", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
