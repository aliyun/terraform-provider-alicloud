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

func dataSourceAliCloudMongoDBInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudMongoDBInstancesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"replicate", "sharding"}, false),
			},
			"instance_class": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
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
						"id": {
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
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"replication": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lock_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expiration_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"mongos": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"node_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"class": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"shards": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"node_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"class": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"storage": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"restore_ranges": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"restore_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"restore_begin_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"restore_end_time": {
										Type:     schema.TypeString,
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

func dataSourceAliCloudMongoDBInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDBInstances"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	if v, ok := d.GetOk("instance_type"); ok {
		request["DBInstanceType"] = v
	}

	if v, ok := d.GetOk("instance_class"); ok {
		request["DBInstanceClass"] = v
	}

	if v, ok := d.GetOk("availability_zone"); ok {
		request["ZoneId"] = v
	}

	if v, ok := d.GetOk("status"); ok {
		request["DBInstanceStatus"] = v
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	var objects []map[string]interface{}
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var instanceNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		instanceNameRegex = r
	}

	var response map[string]interface{}
	var err error

	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_mongodb_instances", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.DBInstances.DBInstance", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.DBInstances.DBInstance", response)
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

			objects = append(objects, item)
		}

		if len(result) < request["PageSize"].(int) {
			break
		}

		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                fmt.Sprint(object["DBInstanceId"]),
			"engine":            object["Engine"],
			"engine_version":    object["EngineVersion"],
			"instance_type":     object["DBInstanceType"],
			"instance_class":    object["DBInstanceClass"],
			"storage":           object["DBInstanceStorage"],
			"network_type":      object["NetworkType"],
			"availability_zone": object["ZoneId"],
			"name":              object["DBInstanceDescription"],
			"charge_type":       object["ChargeType"],
			"replication":       object["ReplicationFactor"],
			"lock_mode":         object["LockMode"],
			"region_id":         object["RegionId"],
			"creation_time":     object["CreationTime"],
			"expiration_time":   object["ExpireTime"],
			"status":            object["DBInstanceStatus"],
		}

		if tagsList, ok := object["Tags"]; ok {
			tags := tagsList.(map[string]interface{})
			if tagMaps, ok := tags["Tag"]; ok {
				mapping["tags"] = tagsToMap(tagMaps)
			}
		}

		if MongosListMap, ok := object["MongosList"].(map[string]interface{}); ok && MongosListMap != nil {
			if MongosList, ok := MongosListMap["MongosAttribute"]; ok && MongosList != nil {
				MongosListMaps := make([]map[string]interface{}, 0)
				for _, MongosLists := range MongosList.([]interface{}) {
					MongosListItemMap := make(map[string]interface{})
					MongosListArg := MongosLists.(map[string]interface{})

					if nodeId, ok := MongosListArg["NodeId"]; ok {
						MongosListItemMap["node_id"] = nodeId
					}

					if class, ok := MongosListArg["NodeClass"]; ok {
						MongosListItemMap["class"] = class
					}

					if description, ok := MongosListArg["NodeDescription"]; ok {
						MongosListItemMap["description"] = description
					}

					MongosListMaps = append(MongosListMaps, MongosListItemMap)
				}

				mapping["mongos"] = MongosListMaps
			}
		}

		if shardListMap, ok := object["ShardList"].(map[string]interface{}); ok && shardListMap != nil {
			if shardList, ok := shardListMap["ShardAttribute"]; ok && shardList != nil {
				shardListMaps := make([]map[string]interface{}, 0)
				for _, shardLists := range shardList.([]interface{}) {
					shardListItemMap := make(map[string]interface{})
					shardListArg := shardLists.(map[string]interface{})

					if nodeId, ok := shardListArg["NodeId"]; ok {
						shardListItemMap["node_id"] = nodeId
					}

					if class, ok := shardListArg["NodeClass"]; ok {
						shardListItemMap["class"] = class
					}

					if storage, ok := shardListArg["NodeStorage"]; ok {
						shardListItemMap["storage"] = storage
					}

					if description, ok := shardListArg["NodeDescription"]; ok {
						shardListItemMap["description"] = description
					}

					shardListMaps = append(shardListMaps, shardListItemMap)
				}

				mapping["shards"] = shardListMaps
			}
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["DBInstanceDescription"])

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(object["DBInstanceId"])
		ddsService := MongoDBService{client}

		recoverTimeDetail, err := ddsService.DescribeMongoDBRecoverTime(d, id)
		if err != nil {
			return WrapError(err)
		}

		if len(recoverTimeDetail) > 0 {
			recoverTimeMaps := make([]map[string]interface{}, 0)
			for _, recoverTime := range recoverTimeDetail {
				recoverTimeArg := recoverTime.(map[string]interface{})
				recoverTimeMap := map[string]interface{}{}

				if restoreType, ok := recoverTimeArg["RestoreType"]; ok {
					recoverTimeMap["restore_type"] = restoreType
				}

				if restoreBeginTime, ok := recoverTimeArg["RestoreBeginTime"]; ok {
					recoverTimeMap["restore_begin_time"] = restoreBeginTime
				}

				if restoreEndTime, ok := recoverTimeArg["RestoreEndTime"]; ok {
					recoverTimeMap["restore_end_time"] = restoreEndTime
				}

				recoverTimeMaps = append(recoverTimeMaps, recoverTimeMap)
			}

			mapping["restore_ranges"] = recoverTimeMaps
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
