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

func dataSourceAlicloudHbrPolicyBindings() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudHbrPolicyBindingsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"UDM_ECS", "NAS", "OSS", "File", "ECS_FILE", "OTS"}, false),
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_bindings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
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
						"exclude": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"include": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source": {
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
						"cross_account_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cross_account_role_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cross_account_user_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
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
												"exclude_disk_id_list": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"destination_kms_key_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"disk_id_list": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
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
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("policy_id"); ok {
		request["PolicyId"] = v
	}
	if v, ok := d.GetOk("source_type"); ok {
		request["SourceType"] = v
	}
	request["MaxResults"] = PageSizeLarge

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

	var objects []map[string]interface{}
	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("hbr", "2017-09-08", action, nil, request, true)
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
			policyId := fmt.Sprint(item["PolicyId"])
			sourceType := fmt.Sprint(item["SourceType"])
			dataSourceId := fmt.Sprint(item["DataSourceId"])
			compositeId := fmt.Sprintf("%v:%v:%v", policyId, sourceType, dataSourceId)
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["PolicyBindingDescription"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[compositeId]; !ok {
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
		policyId := fmt.Sprint(object["PolicyId"])
		sourceType := fmt.Sprint(object["SourceType"])
		dataSourceId := fmt.Sprint(object["DataSourceId"])
		compositeId := fmt.Sprintf("%v:%v:%v", policyId, sourceType, dataSourceId)

		advancedOptionsMaps := make([]map[string]interface{}, 0)
		advancedOptionsRaw := make(map[string]interface{})
		if object["AdvancedOptions"] != nil {
			if raw, ok := object["AdvancedOptions"].(map[string]interface{}); ok {
				advancedOptionsRaw = raw
			}
		}
		if len(advancedOptionsRaw) > 0 {
			advancedOptionsMap := make(map[string]interface{})

			ossDetailMaps := make([]map[string]interface{}, 0)
			ossDetailRaw := make(map[string]interface{})
			if advancedOptionsRaw["OssDetail"] != nil {
				if raw, ok := advancedOptionsRaw["OssDetail"].(map[string]interface{}); ok {
					ossDetailRaw = raw
				}
			}
			if len(ossDetailRaw) > 0 {
				ossDetailMap := make(map[string]interface{})
				ossDetailMap["ignore_archive_object"] = ossDetailRaw["IgnoreArchiveObject"]
				ossDetailMap["inventory_cleanup_policy"] = ossDetailRaw["InventoryCleanupPolicy"]
				ossDetailMap["inventory_id"] = ossDetailRaw["InventoryId"]
				ossDetailMaps = append(ossDetailMaps, ossDetailMap)
			}
			advancedOptionsMap["oss_detail"] = ossDetailMaps

			udmDetailMaps := make([]map[string]interface{}, 0)
			udmDetailRaw := make(map[string]interface{})
			if advancedOptionsRaw["UdmDetail"] != nil {
				if raw, ok := advancedOptionsRaw["UdmDetail"].(map[string]interface{}); ok {
					udmDetailRaw = raw
				}
			}
			if len(udmDetailRaw) > 0 {
				udmDetailMap := make(map[string]interface{})
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
				udmDetailMaps = append(udmDetailMaps, udmDetailMap)
			}
			advancedOptionsMap["udm_detail"] = udmDetailMaps

			advancedOptionsMaps = append(advancedOptionsMaps, advancedOptionsMap)
		}

		mapping := map[string]interface{}{
			"id":                         compositeId,
			"policy_id":                  policyId,
			"source_type":                sourceType,
			"data_source_id":             dataSourceId,
			"disabled":                   object["Disabled"],
			"exclude":                    object["Exclude"],
			"include":                    object["Include"],
			"source":                     object["Source"],
			"speed_limit":                object["SpeedLimit"],
			"policy_binding_description": object["PolicyBindingDescription"],
			"cross_account_type":         object["CrossAccountType"],
			"cross_account_role_name":    object["CrossAccountRoleName"],
			"cross_account_user_id":      object["CrossAccountUserId"],
			"create_time":                fmt.Sprint(object["CreatedTime"]),
			"advanced_options":           advancedOptionsMaps,
		}
		ids = append(ids, compositeId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("policy_bindings", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
