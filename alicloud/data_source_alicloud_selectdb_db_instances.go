package alicloud

import (
	"encoding/json"
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudSelectDBDbInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSelectDBDbInstancesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_id": {
							Type:     schema.TypeString,
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
						"engine_minor_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_prepaid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory_prepaid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cache_size_prepaid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cluster_count_prepaid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cpu_postpaid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory_postpaid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cache_size_postpaid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cluster_count_postpaid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sub_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gmt_created": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gmt_modified": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gmt_expired": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lock_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lock_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudSelectDBDbInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	selectDBService := SelectDBService{client}

	tags := make([]map[string]interface{}, 0)
	if v, ok := d.GetOk("tags"); ok {
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, map[string]interface{}{
				"Key":   key,
				"Value": value.(string),
			})
		}
	}

	idsStr := ""
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if idsStr == "" {
				idsStr = vv.(string)
			} else {
				idsStr = idsStr + ":" + vv.(string)
			}
		}
	}

	instanceResp, err := selectDBService.DescribeSelectDBDbInstances(idsStr, tags)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_selectdb_db_instances", AlibabaCloudSdkGoERROR)
	}

	var objects []map[string]interface{}

	objects = append(objects, instanceResp...)

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		// summary
		mapping := map[string]interface{}{
			"db_instance_id": object["DBInstanceId"],
			"engine":         object["Engine"],
			"engine_version": object["EngineVersion"],

			"db_instance_description": object["Description"],
			"status":                  object["Status"],
			"payment_type":            convertChargeTypeToPaymentType(object["ChargeType"]),

			"region_id":    object["RegionId"],
			"zone_id":      object["ZoneId"],
			"vpc_id":       object["VpcId"],
			"vswitch_id":   object["VswitchId"],
			"gmt_created":  object["GmtCreated"],
			"gmt_modified": object["GmtModified"],
			"gmt_expired":  object["ExpireTime"],
			"lock_mode":    object["LockMode"],
			"lock_reason":  object["LockReason"],
		}
		// cpu,mem,cache
		instanceResp, err := selectDBService.DescribeSelectDBDbInstance(object["DBInstanceId"].(string))
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_selectdb_db_instance", AlibabaCloudSdkGoERROR)
		}
		result := instanceResp["DBClusterList"]
		cpuPrepaid := 0
		cpuPostpaid := 0
		memPrepaid := 0
		memPostpaid := 0
		cachePrepaid := 0
		cachePostpaid := 0

		clusterPrepaidCount := 0
		clusterPostpaidCount := 0

		for _, v := range result.([]interface{}) {
			item := v.(map[string]interface{})
			if item["ChargeType"].(string) == "Postpaid" {
				cpuP, _ := item["CpuCores"].(json.Number).Int64()
				cpuPostpaid += int(cpuP)
				memP, _ := item["Memory"].(json.Number).Int64()
				memPostpaid += int(memP)
				cacheP, _ := item["CacheStorageSizeGB"].(json.Number).Int64()
				cachePostpaid += int(cacheP)
				clusterPostpaidCount += 1
			}
			if item["ChargeType"].(string) == "Prepaid" {
				cpuP, _ := item["CpuCores"].(json.Number).Int64()
				cpuPrepaid += int(cpuP)
				memP, _ := item["Memory"].(json.Number).Int64()
				memPrepaid += int(memP)
				cacheP, _ := item["CacheStorageSizeGB"].(json.Number).Int64()
				cachePrepaid += int(cacheP)
				clusterPrepaidCount += 1
			}
		}
		mapping["cpu_prepaid"] = cpuPrepaid
		mapping["memory_prepaid"] = memPrepaid
		mapping["cache_size_prepaid"] = cachePrepaid
		mapping["cpu_postpaid"] = cpuPostpaid
		mapping["memory_postpaid"] = memPostpaid
		mapping["cache_size_postpaid"] = cachePostpaid

		mapping["cluster_count_prepaid"] = clusterPrepaidCount
		mapping["cluster_count_postpaid"] = clusterPostpaidCount

		mapping["engine_minor_version"] = instanceResp["EngineMinorVersion"]
		mapping["sub_domain"] = instanceResp["SubDomain"]

		id := fmt.Sprint(object["DBInstanceId"])
		mapping["id"] = id
		ids = append(ids, id)

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
