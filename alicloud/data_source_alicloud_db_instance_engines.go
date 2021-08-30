package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudDBInstanceEngines() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDBInstanceEnginesRead,

		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      PostPaid,
				ValidateFunc: validation.StringInSlice([]string{string(PostPaid), string(PrePaid)}, false),
			},
			"engine": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"MySQL", "SQLServer", "PostgreSQL", "PPAS", "MariaDB"}, false),
				Default:      "MySQL",
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"multi_zone": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Computed values.
			"instance_engines": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_ids": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"sub_zone_ids": {
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Computed: true,
									},
								},
							},
							Computed: true,
						},
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"category": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDBInstanceEnginesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeAvailableZones"
	request := map[string]interface{}{
		"RegionId": client.RegionId,
		"Engine":   d.Get("engine").(string),
		"SourceIp": client.SourceIp,
	}
	if v, ok := d.GetOk("engine_version"); ok && v.(string) != "" {
		request["EngineVersion"] = v.(string)
	}
	if v, ok := d.GetOk("zone_id"); ok && v.(string) != "" {
		request["ZoneId"] = v.(string)
	}
	instanceChargeType := d.Get("instance_charge_type").(string)
	if instanceChargeType == string(PostPaid) {
		request["InstanceChargeType"] = string(Postpaid)
	} else {
		request["InstanceChargeType"] = string(Prepaid)
	}
	multiZone := d.Get("multi_zone").(bool)
	var ids []string
	var s []map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_db_instance_engines", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.AvailableZones", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.AvailableZones", response)
	}
	engine, engineGot := d.GetOk("engine")
	engineVersion, engineVersionGot := d.GetOk("engine_version")
	for _, r := range resp.([]interface{}) {
		availableZoneItem := r.(map[string]interface{})

		zoneId := fmt.Sprint(availableZoneItem["ZoneId"])
		if multiZone && !strings.Contains(zoneId, MULTI_IZ_SYMBOL) {
			continue
		}
		zoneIds := make([]map[string]interface{}, 0)
		zoneIds = append(zoneIds, map[string]interface{}{
			"id":           zoneId,
			"sub_zone_ids": splitMultiZoneId(zoneId),
		})

		for _, r := range availableZoneItem["SupportedEngines"].([]interface{}) {
			supportedEngineItem := r.(map[string]interface{})
			if engineGot && engine.(string) != "" && engine.(string) != fmt.Sprint(supportedEngineItem["Engine"]) {
				continue
			}
			for _, r := range supportedEngineItem["SupportedEngineVersions"].([]interface{}) {
				supportedEngineVersionItem := r.(map[string]interface{})
				if engineVersionGot && engineVersion.(string) != "" && engineVersion.(string) != fmt.Sprint(supportedEngineVersionItem["Version"]) {
					continue
				}
				for _, r := range supportedEngineVersionItem["SupportedCategorys"].([]interface{}) {
					supportedCategoryItem := r.(map[string]interface{})
					if fmt.Sprint(supportedCategoryItem["Category"]) == "" {
						continue
					}
					mapping := map[string]interface{}{
						"zone_ids":       zoneIds,
						"engine":         fmt.Sprint(supportedEngineItem["Engine"]),
						"engine_version": fmt.Sprint(supportedEngineVersionItem["Version"]),
						"category":       fmt.Sprint(supportedCategoryItem["Category"]),
					}
					s = append(s, mapping)
					ids = append(ids, fmt.Sprint(supportedEngineItem["Engine"]))
				}
			}
		}
	}

	d.SetId(dataResourceIdHash(ids))
	err = d.Set("instance_engines", s)
	if err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		err = writeToFile(output.(string), s)
		if err != nil {
			return WrapError(err)
		}
	}
	return nil
}
