package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudGpdbInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudGpdbInstancesRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_instance_categories": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"db_instance_modes": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Creating", "DBInstanceClassChanging", "DBInstanceNetTypeChanging", "Deleting", "EngineVersionUpgrading", "GuardDBInstanceCreating", "GuardSwitching", "Importing", "ImportingFromOtherInstance", "Rebooting", "Restoring", "Running", "Transfering", "TransferingToOtherInstance"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_network_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connection_string": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_cores": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
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
						"ip_whitelist": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_group_attribute": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip_group_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"security_ip_list": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"instance_network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"maintain_end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"maintain_start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_node_num": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"seg_node_num": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"storage_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAlicloudGpdbInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDBInstances"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("db_instance_categories"); ok {
		request["DBInstanceCategories"] = v
	}
	if v, ok := d.GetOk("db_instance_modes"); ok {
		request["DBInstanceModes"] = v
	}
	if v, ok := d.GetOk("db_instance_statuses"); ok {
		request["DBInstanceStatuses"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["DBInstanceDescription"] = v
	}

	if v, ok := d.GetOk("instance_network_type"); ok {
		request["InstanceNetworkType"] = v
	}

	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	availabilityZone := d.Get("availability_zone")

	vSwitchId := d.Get("vswitch_id")

	if v, ok := d.GetOk("tags"); ok {
		tags := make([]map[string]interface{}, 0)
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, map[string]interface{}{
				"Key":   key,
				"Value": value.(string),
			})
		}
		request["Tag"] = tags
	}
	setPagingRequest(d, request, PageSizeLarge)
	var objects []map[string]interface{}
	var instanceNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		instanceNameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	status, statusOk := d.GetOk("status")
	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("gpdb", "2016-05-03", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_gpdb_instances", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Items.DBInstance", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Items.DBInstance", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["DBInstanceId"])]; !ok {
					continue
				}
			}
			if instanceNameRegex != nil && !instanceNameRegex.MatchString(fmt.Sprint(item["DBInstanceDescription"])) {
				continue
			}

			if availabilityZone != "" && availabilityZone != item["ZoneId"] {
				continue
			}

			if vSwitchId != "" && vSwitchId != item["VSwitchId"] {
				continue
			}
			if statusOk && status.(string) != "" && status.(string) != item["DBInstanceStatus"].(string) {
				continue
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"create_time":           object["CreateTime"],
			"id":                    fmt.Sprint(object["DBInstanceId"]),
			"db_instance_id":        fmt.Sprint(object["DBInstanceId"]),
			"description":           object["DBInstanceDescription"],
			"engine":                object["Engine"],
			"engine_version":        object["EngineVersion"],
			"instance_network_type": object["InstanceNetworkType"],
			"payment_type":          convertGpdbDbInstancePaymentTypeResponse(object["PayType"]),
			"status":                object["DBInstanceStatus"],
			"vswitch_id":            object["VSwitchId"],
			"vpc_id":                object["VpcId"],
			"zone_id":               object["ZoneId"],
			"region_id":             object["RegionId"],
			"availability_zone":     object["ZoneId"],
			"creation_time":         object["CreateTime"],
			"charge_type":           object["PayType"],
		}

		tags := make(map[string]interface{})
		t, _ := jsonpath.Get("$.Tags.Tag", object)
		if t != nil {
			for _, t := range t.([]interface{}) {
				key := t.(map[string]interface{})["Key"].(string)
				value := t.(map[string]interface{})["Value"].(string)
				if !ignoredTags(key, value) {
					tags[key] = value
				}
			}
		}
		mapping["tags"] = tags
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["DBInstanceDescription"])
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["DBInstanceId"])
		gpdbService := GpdbService{client}
		getResp, err := gpdbService.DescribeGpdbDbInstance(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["connection_string"] = getResp["ConnectionString"]
		mapping["cpu_cores"] = getResp["CpuCores"]
		mapping["db_instance_category"] = getResp["DBInstanceCategory"]
		mapping["db_instance_mode"] = getResp["DBInstanceMode"]
		mapping["maintain_end_time"] = getResp["MaintainEndTime"]
		mapping["maintain_start_time"] = getResp["MaintainStartTime"]
		mapping["master_node_num"] = getResp["MasterNodeNum"]
		mapping["memory_size"] = getResp["MemorySize"]
		mapping["seg_node_num"] = getResp["SegNodeNum"]
		if v, ok := getResp["StorageSize"]; ok && fmt.Sprint(v) != "0" {
			mapping["storage_size"] = formatInt(v)
		}
		mapping["storage_type"] = getResp["StorageType"]
		gpdbService = GpdbService{client}
		getResp1, err := gpdbService.DescribeDBInstanceIPArrayList(id)
		if err != nil {
			return WrapError(err)
		}
		if iPWhitelistMap, ok := getResp1["Items"].(map[string]interface{}); ok && iPWhitelistMap != nil {
			if dBInstanceIPArrayList, ok := iPWhitelistMap["DBInstanceIPArray"]; ok && dBInstanceIPArrayList != nil {
				iPWhitelistMaps := make([]map[string]interface{}, 0)
				for _, dBInstanceIPArrayListItem := range dBInstanceIPArrayList.([]interface{}) {
					if dBInstanceIPArrayListItemMap, ok := dBInstanceIPArrayListItem.(map[string]interface{}); ok {
						if dBInstanceIPArrayListItem.(map[string]interface{})["DBInstanceIPArrayAttribute"] == "hidden" {
							continue
						}
						dBInstanceIPArrayListMap := map[string]interface{}{}
						dBInstanceIPArrayListMap["ip_group_attribute"] = dBInstanceIPArrayListItemMap["DBInstanceIPArrayAttribute"]
						dBInstanceIPArrayListMap["ip_group_name"] = dBInstanceIPArrayListItemMap["DBInstanceIPArrayName"]
						dBInstanceIPArrayListMap["security_ip_list"] = dBInstanceIPArrayListItemMap["SecurityIPList"]
						iPWhitelistMaps = append(iPWhitelistMaps, dBInstanceIPArrayListMap)
					}
				}
				mapping["ip_whitelist"] = iPWhitelistMaps
			}
		}

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
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
