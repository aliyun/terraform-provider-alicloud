package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudDBInstanceClassInfos() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDBInstanceClassInfosRead,

		Schema: map[string]*schema.Schema{
			"commodity_code": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"bards", "rds", "rords", "rds_rordspre_public_cn", "bards_intl", "rds_intl", "rords_intl", "rds_rordspre_public_intl"}, false),
			},
			"db_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"BUY", "UPGRADE", "RENEW", "CONVERT"}, false),
			},
			"ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// response value
			"infos": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_iombps": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"class_code": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"max_connections": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"class_group": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cpu": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"instruction_set_arch": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"memory_class": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"reference_price": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"max_iops": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDBInstanceClassInfosRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "ListClasses"
	request := map[string]interface{}{
		"RegionId":      client.RegionId,
		"SourceIp":      client.SourceIp,
		"CommodityCode": d.Get("commodity_code"),
		"OrderType":     d.Get("order_type"),
	}
	instanceIdFlag := false
	roCommodityArray := []string{"rords", "rds_rordspre_public_cn", "rords_intl", "rds_rordspre_public_intl"}
	for _, ro := range roCommodityArray {
		if ro == d.Get("commodity_code") {
			instanceIdFlag = true
		}
	}
	if instanceIdFlag {
		if v, ok := d.GetOk("db_instance_id"); ok {
			request["DBInstanceId"] = v.(string)
		}
	}
	var ids []string
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_db_instance_class_infos", action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Items", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Items", response)
	}
	result := make([]map[string]interface{}, 0)
	for _, r := range resp.([]interface{}) {
		instanceClassInfoItem := r.(map[string]interface{})
		mapping := map[string]interface{}{}
		if v, ok := instanceClassInfoItem["MaxIOMBPS"]; ok && v != "" {
			mapping["max_iombps"] = fmt.Sprint(v)
		}
		if v, ok := instanceClassInfoItem["ClassCode"]; ok && v != "" {
			mapping["class_code"] = fmt.Sprint(v)
		}
		if v, ok := instanceClassInfoItem["MaxConnections"]; ok && v != "" {
			mapping["max_connections"] = fmt.Sprint(v)
		}
		if v, ok := instanceClassInfoItem["ClassGroup"]; ok && v != "" {
			mapping["class_group"] = fmt.Sprint(v)
		}
		if v, ok := instanceClassInfoItem["Cpu"]; ok && v != "" {
			mapping["cpu"] = fmt.Sprint(v)
		}
		if v, ok := instanceClassInfoItem["MemoryClass"]; ok && v != "" {
			mapping["memory_class"] = fmt.Sprint(v)
		}
		if v, ok := instanceClassInfoItem["MaxIOPS"]; ok && v != "" {
			mapping["max_iops"] = fmt.Sprint(v)
		}
		if v, ok := instanceClassInfoItem["InstructionSetArch"]; ok && v != "" {
			mapping["instruction_set_arch"] = fmt.Sprint(v)
		}
		if v, ok := instanceClassInfoItem["ReferencePrice"]; ok && v != "" {
			mapping["reference_price"] = fmt.Sprint(v)
		}

		ids = append(ids, fmt.Sprint(instanceClassInfoItem["ClassCode"]))
		result = append(result, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("infos", result); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), result)
	}

	return nil
}
