package selectdb

import (
	"encoding/json"
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDbClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDbClustersRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_cluster_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_cluster_description": {
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
						"create_time": {
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
						"cpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cache_size": {
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
						"params": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"optional": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"comment": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"param_category": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"default_value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"is_dynamic": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"is_user_modifiable": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"param_change_logs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"old_value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"new_value": {
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
									"config_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"is_applied": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceDbClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	selectDBService := NewService(client)

	instanceMap := make(map[string]string)
	idsMap := make(map[string]string)

	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			parts, err := helper.ParseResourceId(vv.(string), 2)
			if err != nil {
				return errs.WrapError(err)
			}
			//instanceid:clusterid clusterid
			idsMap[vv.(string)] = parts[1]
			instanceMap[parts[0]] = parts[0]
		}
	}

	instanceResult := make(map[string]interface{})
	instanceClusterResult := make(map[string]interface{})
	for _, instanceId := range instanceMap {
		instanceResp, err := selectDBService.DescribeDbInstance(instanceId)
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_selectdb_db_clusters", errs.AlibabaCloudSdkGoERROR)
		}
		resp := instanceResp["DBClusterList"]
		instanceResult[instanceId] = instanceResp
		instanceClusterResult[instanceId] = resp.([]interface{})
	}

	var objects []map[string]interface{}

	if len(idsMap) > 0 {
		for pairId, pairClusterId := range idsMap {
			parts, err := helper.ParseResourceId(pairId, 2)
			if err != nil {
				return errs.WrapError(err)
			}
			result := instanceClusterResult[parts[0]].([]interface{})
			for _, v := range result {
				item := v.(map[string]interface{})
				if item["DbClusterId"].(string) == pairClusterId {
					objects = append(objects, item)
				}
			}
		}
	} else {
		item := make(map[string]interface{})
		objects = append(objects, item)
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		instanceResp := instanceResult[object["DbInstanceName"].(string)].(map[string]interface{})
		cpuP, _ := object["CpuCores"].(json.Number).Int64()
		memoryP, _ := object["Memory"].(json.Number).Int64()
		cacheP, _ := object["CacheStorageSizeGB"].(json.Number).Int64()

		mapping := map[string]interface{}{
			"status":                 object["Status"].(string),
			"create_time":            object["CreatedTime"].(string),
			"db_cluster_description": object["DbClusterName"].(string),
			"payment_type":           convert.ConvertChargeTypeToPaymentType(object["ChargeType"]),
			"cpu":                    cpuP,
			"memory":                 memoryP,
			"cache_size":             cacheP,
			"db_instance_id":         object["DbInstanceName"].(string),
			"db_cluster_id":          object["DbClusterId"].(string),
			"db_cluster_class":       object["DbClusterClass"].(string),
			"engine":                 fmt.Sprint(instanceResp["Engine"]),
			"engine_version":         fmt.Sprint(instanceResp["EngineVersion"]),
			"vpc_id":                 fmt.Sprint(instanceResp["VpcId"]),
			"zone_id":                fmt.Sprint(instanceResp["ZoneId"]),
			"region_id":              fmt.Sprint(instanceResp["RegionId"]),
		}

		id := fmt.Sprint(object["DbInstanceName"]) + ":" + fmt.Sprint(object["DbClusterId"])
		mapping["id"] = id

		selectDBService := NewService(client)

		configArrayList, err := selectDBService.DescribeDbClusterConfig(id)
		if err != nil {
			return errs.WrapError(err)
		}
		configArray := make([]map[string]interface{}, 0)
		for _, v := range configArrayList {
			m1 := v.(map[string]interface{})
			temp1 := map[string]interface{}{
				"comment":            m1["Comment"].(string),
				"default_value":      m1["DefaultValue"].(string),
				"optional":           m1["Optional"].(string),
				"param_category":     m1["ParamCategory"].(string),
				"value":              m1["Value"].(string),
				"is_user_modifiable": m1["IsUserModifiable"],
				"is_dynamic":         m1["IsDynamic"],
				"name":               m1["Name"].(string),
			}
			// config with default value will not be updated
			if m1["DefaultValue"].(string) != m1["Value"].(string) {
				configArray = append(configArray, temp1)
			}
		}
		mapping["params"] = configArray

		configChangeArrayList, err := selectDBService.DescribeDbClusterConfigChangeLog(id)
		if err != nil {
			return errs.WrapError(err)
		}
		configChangeArray := make([]map[string]interface{}, 0)
		for _, v := range configChangeArrayList {
			m1 := v.(map[string]interface{})
			temp1 := map[string]interface{}{
				"name":         m1["Name"].(string),
				"old_value":    m1["OldValue"].(string),
				"new_value":    m1["NewValue"].(string),
				"is_applied":   m1["isApplied"],
				"gmt_created":  m1["GmtCreated"].(string),
				"gmt_modified": m1["GmtModified"].(string),
				"config_id":    m1["ConfigId"],
			}
			configChangeArray = append(configChangeArray, temp1)
		}
		mapping["param_change_logs"] = configChangeArray

		ids = append(ids, id)
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("clusters", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
