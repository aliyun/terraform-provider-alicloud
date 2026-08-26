// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudMilvusInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudMilvusInstanceRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_backup": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"components": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cu_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cu_num": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"data_disk": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"storage_class": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"size": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"performance_level": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"enabled": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"disk_size_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"pay_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"replica": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"configuration": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"encrypted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"expire_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ha": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"multi_zone_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"order_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"running_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"security_group_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"vswitch_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vsw_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"zone_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAliCloudMilvusInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
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

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	// ListInstancesV2
	action := fmt.Sprintf("/webapi/instance/list")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	query["RegionId"] = StringPointer(client.RegionId)
	query["instanceId"] = StringPointer(d.Get("instance_id").(string))
	if v, ok := d.GetOk("instance_id"); ok {
		query["instanceId"] = StringPointer(v.(string))
	}

	request["instancename"] = d.Get("instance_name")
	if v, ok := d.GetOk("resource_group_id"); ok {
		query["resourceGroupId"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsMap := v.(map[string]interface{})
		i := 0
		for key, value := range tagsMap {
			query[fmt.Sprintf("tag[%d].key", i)] = StringPointer(key)
			query[fmt.Sprintf("tag[%d].value", i)] = StringPointer(fmt.Sprint(value))
			i++
		}
	}

	query["PageSize"] = StringPointer(strconv.Itoa(PageSizeLarge))
	query["PageNumber"] = StringPointer("1")
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RoaGet("milvus", "2023-10-12", action, query, nil, nil)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.instances[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["instanceName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["instanceId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < PageSizeLarge {
			break
		}
		pageNum, _ := strconv.Atoi(*query["PageNumber"])
		query["PageNumber"] = StringPointer(strconv.Itoa(pageNum + 1))
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["instanceId"]

		mapping["auto_backup"] = objectRaw["autoBackup"]
		mapping["configuration"] = objectRaw["configuration"]
		mapping["create_time"] = objectRaw["createTime"]
		mapping["db_version"] = objectRaw["dbVersion"]
		mapping["encrypted"] = objectRaw["encrypted"]
		mapping["expire_time"] = objectRaw["expireTime"]
		mapping["ha"] = objectRaw["ha"]
		mapping["instance_name"] = objectRaw["instanceName"]
		mapping["kms_key_id"] = objectRaw["kmsKeyId"]
		mapping["multi_zone_mode"] = objectRaw["multiZoneMode"]
		mapping["order_id"] = objectRaw["orderId"]
		mapping["payment_type"] = objectRaw["paymentType"]
		mapping["region_id"] = objectRaw["regionId"]
		mapping["resource_group_id"] = objectRaw["resourceGroupId"]
		mapping["running_time"] = objectRaw["runningTime"]
		mapping["status"] = objectRaw["status"]
		mapping["vpc_id"] = objectRaw["vpcId"]
		mapping["zone_id"] = objectRaw["zoneId"]
		mapping["instance_id"] = objectRaw["instanceId"]

		componentsRaw := objectRaw["components"]
		componentsMaps := make([]map[string]interface{}, 0)
		if componentsRaw != nil {
			for _, componentsChildRaw := range convertToInterfaceArray(componentsRaw) {
				componentsMap := make(map[string]interface{})
				componentsChildRaw := componentsChildRaw.(map[string]interface{})
				componentsMap["cu_num"] = componentsChildRaw["cuNum"]
				componentsMap["cu_type"] = componentsChildRaw["cuType"]
				componentsMap["disk_size_type"] = componentsChildRaw["diskSizeType"]
				componentsMap["replica"] = componentsChildRaw["replica"]
				componentsMap["type"] = componentsChildRaw["type"]

				dataDiskMaps := make([]map[string]interface{}, 0)
				dataDiskMap := make(map[string]interface{})
				dataDiskRaw := make(map[string]interface{})
				if componentsChildRaw["dataDisk"] != nil {
					dataDiskRaw = componentsChildRaw["dataDisk"].(map[string]interface{})
				}
				if len(dataDiskRaw) > 0 {
					dataDiskMap["enabled"] = dataDiskRaw["Enabled"]
					dataDiskMap["performance_level"] = dataDiskRaw["PerformanceLevel"]
					dataDiskMap["size"] = dataDiskRaw["Size"]
					dataDiskMap["storage_class"] = dataDiskRaw["StorageClass"]

					dataDiskMaps = append(dataDiskMaps, dataDiskMap)
				}
				componentsMap["data_disk"] = dataDiskMaps
				componentsMaps = append(componentsMaps, componentsMap)
			}
		}
		mapping["components"] = componentsMaps
		securityGroupIdsRaw := make([]interface{}, 0)
		if objectRaw["securityGroupIds"] != nil {
			securityGroupIdsRaw = convertToInterfaceArray(objectRaw["securityGroupIds"])
		}

		mapping["security_group_ids"] = securityGroupIdsRaw
		tagsMaps := objectRaw["tags"]
		mapping["tags"] = tagsToMap(tagsMaps)
		vSwitchIdsRaw := objectRaw["vSwitchIds"]
		vSwitchIdsMaps := make([]map[string]interface{}, 0)
		if vSwitchIdsRaw != nil {
			for _, vSwitchIdsChildRaw := range convertToInterfaceArray(vSwitchIdsRaw) {
				vSwitchIdsMap := make(map[string]interface{})
				vSwitchIdsChildRaw := vSwitchIdsChildRaw.(map[string]interface{})
				vSwitchIdsMap["vsw_id"] = vSwitchIdsChildRaw["vswId"]
				vSwitchIdsMap["zone_id"] = vSwitchIdsChildRaw["zoneId"]

				vSwitchIdsMaps = append(vSwitchIdsMaps, vSwitchIdsMap)
			}
		}
		mapping["vswitch_ids"] = vSwitchIdsMaps

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["instanceName"])
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
