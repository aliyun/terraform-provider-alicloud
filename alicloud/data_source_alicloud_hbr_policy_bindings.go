package alicloud

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudHbrPolicyBindings() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudHbrPolicyBindingsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"UDM_ECS", "NAS", "OSS", "File", "ECS_FILE", "OTS"}, false),
			},
			"data_source_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bindings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_binding_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_source_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"source": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"include": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"exclude": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"speed_limit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_binding_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cross_account_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cross_account_user_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cross_account_role_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_by_tag": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"hit_tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"operator": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"advanced_options": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"oss_detail": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ignore_archive_object": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"inventory_cleanup_policy": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"inventory_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"udm_detail": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"disk_id_list": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"exclude_disk_id_list": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"destination_kms_key_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"app_consistent": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"snapshot_group": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"ram_role_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"pre_script_path": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"post_script_path": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"enable_fs_freeze": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"timeout_in_seconds": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"enable_writers": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
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

func dataSourceAlicloudHbrPolicyBindingsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribePolicyBindings"
	request := make(map[string]interface{})
	query := make(map[string]interface{})

	if v, ok := d.GetOk("policy_id"); ok {
		request["PolicyId"] = v
	}
	if v, ok := d.GetOk("source_type"); ok {
		request["SourceType"] = v
	}
	dataSourceIdsMap := make(map[string]string)
	if v, ok := d.GetOk("data_source_ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			dataSourceIdsMap[vv.(string)] = vv.(string)
		}
		// The API matches DataSourceIds only when SourceType is set at the same time;
		// otherwise the filter silently returns an empty list. Membership is always
		// verified client-side below, so only narrow server-side when SourceType exists.
		if _, ok := request["SourceType"]; ok && len(dataSourceIdsMap) > 0 {
			dataSourceIdsJson, err := json.Marshal(v.([]interface{}))
			if err != nil {
				return WrapError(err)
			}
			request["DataSourceIds"] = string(dataSourceIdsJson)
		}
	}
	request["MaxResults"] = PageSizeLarge

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

	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("hbr", "2017-09-08", action, query, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_hbr_policy_bindings", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.PolicyBindings[*]", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.PolicyBindings[*]", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				id := fmt.Sprintf("%v:%v:%v", item["PolicyId"], item["SourceType"], item["DataSourceId"])
				if _, ok := idsMap[id]; !ok {
					continue
				}
			}
			if len(dataSourceIdsMap) > 0 {
				if _, ok := dataSourceIdsMap[fmt.Sprint(item["DataSourceId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                         fmt.Sprintf("%v:%v:%v", object["PolicyId"], object["SourceType"], object["DataSourceId"]),
			"policy_binding_id":          object["PolicyBindingId"],
			"policy_id":                  object["PolicyId"],
			"source_type":                object["SourceType"],
			"data_source_id":             object["DataSourceId"],
			"disabled":                   object["Disabled"],
			"source":                     object["Source"],
			"include":                    object["Include"],
			"exclude":                    object["Exclude"],
			"speed_limit":                object["SpeedLimit"],
			"policy_binding_description": object["PolicyBindingDescription"],
			"create_time":                fmt.Sprint(object["CreatedTime"]),
			"cross_account_type":         object["CrossAccountType"],
			"cross_account_user_id":      object["CrossAccountUserId"],
			"cross_account_role_name":    object["CrossAccountRoleName"],
			"created_by_tag":             object["CreatedByTag"],
		}

		hitTagsMaps := make([]map[string]interface{}, 0)
		if hitTagsRaw, ok := object["HitTags"].([]interface{}); ok {
			for _, v := range hitTagsRaw {
				hitTag := v.(map[string]interface{})
				hitTagsMaps = append(hitTagsMaps, map[string]interface{}{
					"key":      hitTag["Key"],
					"value":    hitTag["Value"],
					"operator": hitTag["Operator"],
				})
			}
		}
		mapping["hit_tags"] = hitTagsMaps

		advancedOptionsMaps := make([]map[string]interface{}, 0)
		advancedOptionsMap := make(map[string]interface{})
		advancedOptionsRaw := make(map[string]interface{})
		if object["AdvancedOptions"] != nil {
			advancedOptionsRaw = object["AdvancedOptions"].(map[string]interface{})
		}
		if len(advancedOptionsRaw) > 0 {
			ossDetailMaps := make([]map[string]interface{}, 0)
			ossDetailMap := make(map[string]interface{})
			ossDetailRaw := make(map[string]interface{})
			if advancedOptionsRaw["OssDetail"] != nil {
				ossDetailRaw = advancedOptionsRaw["OssDetail"].(map[string]interface{})
			}
			if len(ossDetailRaw) > 0 {
				ossDetailMap["ignore_archive_object"] = ossDetailRaw["IgnoreArchiveObject"]
				ossDetailMap["inventory_cleanup_policy"] = ossDetailRaw["InventoryCleanupPolicy"]
				ossDetailMap["inventory_id"] = ossDetailRaw["InventoryId"]
				ossDetailMaps = append(ossDetailMaps, ossDetailMap)
			}
			advancedOptionsMap["oss_detail"] = ossDetailMaps
			udmDetailMaps := make([]map[string]interface{}, 0)
			udmDetailMap := make(map[string]interface{})
			udmDetailRaw := make(map[string]interface{})
			if advancedOptionsRaw["UdmDetail"] != nil {
				udmDetailRaw = advancedOptionsRaw["UdmDetail"].(map[string]interface{})
			}
			if len(udmDetailRaw) > 0 {
				udmDetailMap["destination_kms_key_id"] = udmDetailRaw["DestinationKmsKeyId"]
				diskIdListRaw := make([]interface{}, 0)
				if udmDetailRaw["DiskIdList"] != nil {
					diskIdListRaw = convertToInterfaceArray(udmDetailRaw["DiskIdList"])
				}
				udmDetailMap["disk_id_list"] = diskIdListRaw
				excludeDiskIdListRaw := make([]interface{}, 0)
				if udmDetailRaw["ExcludeDiskIdList"] != nil {
					excludeDiskIdListRaw = convertToInterfaceArray(udmDetailRaw["ExcludeDiskIdList"])
				}
				udmDetailMap["exclude_disk_id_list"] = excludeDiskIdListRaw
				udmDetailMap["app_consistent"] = udmDetailRaw["AppConsistent"]
				udmDetailMap["snapshot_group"] = udmDetailRaw["SnapshotGroup"]
				udmDetailMap["ram_role_name"] = udmDetailRaw["RamRoleName"]
				udmDetailMap["pre_script_path"] = udmDetailRaw["PreScriptPath"]
				udmDetailMap["post_script_path"] = udmDetailRaw["PostScriptPath"]
				udmDetailMap["enable_fs_freeze"] = udmDetailRaw["EnableFsFreeze"]
				udmDetailMap["timeout_in_seconds"] = formatInt(udmDetailRaw["TimeoutInSeconds"])
				udmDetailMap["enable_writers"] = udmDetailRaw["EnableWriters"]
				udmDetailMaps = append(udmDetailMaps, udmDetailMap)
			}
			advancedOptionsMap["udm_detail"] = udmDetailMaps
			advancedOptionsMaps = append(advancedOptionsMaps, advancedOptionsMap)
		}
		mapping["advanced_options"] = advancedOptionsMaps

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("bindings", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
