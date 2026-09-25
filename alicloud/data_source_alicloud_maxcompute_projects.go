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

func dataSourceAliCloudMaxComputeProjects() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudMaxComputeProjectRead,
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
			"projects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"comment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cost_storage": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_quota": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_white_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_ip_list": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip_list": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"properties": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sql_metering_max": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"timezone": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tunnel_quota": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"table_lifecycle": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"type_system": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enable_data_masking": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"enable_tunnel_quota_route": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"retention_days": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"encryption": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enable": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"algorithm": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"key": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"allow_full_scan": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"enable_dr": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"enable_decimal2": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_properties": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"using_policy": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"label_security": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"object_creator_has_grant_permission": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"object_creator_has_access_permission": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"using_acl": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"enable_download_privilege": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"project_protection": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"protected": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"exception_policy": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
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
						"three_tier_model": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"trusted_projects": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"type": {
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
			"product_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAliCloudMaxComputeProjectRead(d *schema.ResourceData, meta interface{}) error {
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
	// ListProjects
	action := fmt.Sprintf("/api/v1/projects")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)

	query["maxItem"] = StringPointer(strconv.Itoa(PageSizeLarge))
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RoaGet("MaxCompute", "2022-01-04", action, query, nil, nil)

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

		resp, _ := jsonpath.Get("$.data.projects[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["name"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if nextToken, ok := response["data.marker"].(string); ok && nextToken != "" {
			query["marker"] = StringPointer(nextToken)
		} else {
			break
		}
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["name"]

		mapping["comment"] = objectRaw["comment"]
		mapping["cost_storage"] = objectRaw["costStorage"]
		mapping["create_time"] = objectRaw["createdTime"]
		mapping["default_quota"] = objectRaw["defaultQuota"]
		mapping["owner"] = objectRaw["owner"]
		mapping["region_id"] = objectRaw["regionId"]
		mapping["status"] = objectRaw["status"]
		mapping["three_tier_model"] = objectRaw["threeTierModel"]
		mapping["type"] = objectRaw["type"]
		mapping["project_name"] = objectRaw["name"]

		ipWhiteListMaps := make([]map[string]interface{}, 0)
		ipWhiteListMap := make(map[string]interface{})
		ipWhiteListRaw := make(map[string]interface{})
		if objectRaw["ipWhiteList"] != nil {
			ipWhiteListRaw = objectRaw["ipWhiteList"].(map[string]interface{})
		}
		if len(ipWhiteListRaw) > 0 {
			ipWhiteListMap["ip_list"] = ipWhiteListRaw["ipList"]
			ipWhiteListMap["vpc_ip_list"] = ipWhiteListRaw["vpcIpList"]

			ipWhiteListMaps = append(ipWhiteListMaps, ipWhiteListMap)
		}
		mapping["ip_white_list"] = ipWhiteListMaps
		propertiesMaps := make([]map[string]interface{}, 0)
		propertiesMap := make(map[string]interface{})
		propertiesRaw := make(map[string]interface{})
		if objectRaw["properties"] != nil {
			propertiesRaw = objectRaw["properties"].(map[string]interface{})
		}
		if len(propertiesRaw) > 0 {
			propertiesMap["allow_full_scan"] = propertiesRaw["allowFullScan"]
			propertiesMap["enable_decimal2"] = propertiesRaw["enableDecimal2"]
			propertiesMap["enable_tunnel_quota_route"] = propertiesRaw["enableTunnelQuotaRoute"]
			propertiesMap["retention_days"] = propertiesRaw["retentionDays"]
			propertiesMap["sql_metering_max"] = propertiesRaw["sqlMeteringMax"]
			propertiesMap["timezone"] = propertiesRaw["timezone"]
			propertiesMap["tunnel_quota"] = propertiesRaw["tunnelQuota"]
			propertiesMap["type_system"] = propertiesRaw["typeSystem"]

			encryptionMaps := make([]map[string]interface{}, 0)
			encryptionMap := make(map[string]interface{})
			encryptionRaw := make(map[string]interface{})
			if propertiesRaw["encryption"] != nil {
				encryptionRaw = propertiesRaw["encryption"].(map[string]interface{})
			}
			if len(encryptionRaw) > 0 {
				encryptionMap["algorithm"] = encryptionRaw["algorithm"]
				encryptionMap["enable"] = encryptionRaw["enable"]
				encryptionMap["key"] = encryptionRaw["key"]

				encryptionMaps = append(encryptionMaps, encryptionMap)
			}
			propertiesMap["encryption"] = encryptionMaps
			tableLifecycleMaps := make([]map[string]interface{}, 0)
			tableLifecycleMap := make(map[string]interface{})
			tableLifecycleRaw := make(map[string]interface{})
			if propertiesRaw["tableLifecycle"] != nil {
				tableLifecycleRaw = propertiesRaw["tableLifecycle"].(map[string]interface{})
			}
			if len(tableLifecycleRaw) > 0 {
				tableLifecycleMap["type"] = tableLifecycleRaw["type"]
				tableLifecycleMap["value"] = tableLifecycleRaw["value"]

				tableLifecycleMaps = append(tableLifecycleMaps, tableLifecycleMap)
			}
			propertiesMap["table_lifecycle"] = tableLifecycleMaps
			propertiesMaps = append(propertiesMaps, propertiesMap)
		}
		mapping["properties"] = propertiesMaps
		securityPropertiesMaps := make([]map[string]interface{}, 0)
		securityPropertiesMap := make(map[string]interface{})
		securityPropertiesRaw := make(map[string]interface{})
		if objectRaw["securityProperties"] != nil {
			securityPropertiesRaw = objectRaw["securityProperties"].(map[string]interface{})
		}
		if len(securityPropertiesRaw) > 0 {
			securityPropertiesMap["enable_download_privilege"] = securityPropertiesRaw["enableDownloadPrivilege"]
			securityPropertiesMap["label_security"] = securityPropertiesRaw["labelSecurity"]
			securityPropertiesMap["object_creator_has_access_permission"] = securityPropertiesRaw["objectCreatorHasAccessPermission"]
			securityPropertiesMap["object_creator_has_grant_permission"] = securityPropertiesRaw["objectCreatorHasGrantPermission"]
			securityPropertiesMap["using_acl"] = securityPropertiesRaw["usingAcl"]
			securityPropertiesMap["using_policy"] = securityPropertiesRaw["usingPolicy"]

			projectProtectionMaps := make([]map[string]interface{}, 0)
			projectProtectionMap := make(map[string]interface{})
			projectProtectionRaw := make(map[string]interface{})
			if securityPropertiesRaw["projectProtection"] != nil {
				projectProtectionRaw = securityPropertiesRaw["projectProtection"].(map[string]interface{})
			}
			if len(projectProtectionRaw) > 0 {
				projectProtectionMap["exception_policy"] = projectProtectionRaw["exceptionPolicy"]
				projectProtectionMap["protected"] = projectProtectionRaw["protected"]

				projectProtectionMaps = append(projectProtectionMaps, projectProtectionMap)
			}
			securityPropertiesMap["project_protection"] = projectProtectionMaps
			securityPropertiesMaps = append(securityPropertiesMaps, securityPropertiesMap)
		}
		mapping["security_properties"] = securityPropertiesMaps

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			names = append(names, objectRaw["name"])
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(objectRaw["name"])
		mapping, err = dataSourceAliCloudMaxComputeProjectReadDescription(d, id, mapping, meta)
		if err != nil {
			return WrapError(err)
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["name"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("projects", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

func dataSourceAliCloudMaxComputeProjectReadDescription(d *schema.ResourceData, id string, object map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	maxComputeServiceV2 := MaxComputeServiceV2{client}
	getResp, err := maxComputeServiceV2.DescribeMaxComputeProject(id)
	if err != nil {
		return nil, WrapError(err)
	}

	// Merge additional fields from Get API response to mapping
	// Reuse the response mapping template from Resource's read function
	mapping := object
	objectRaw := getResp

	mapping["comment"] = objectRaw["comment"]
	mapping["cost_storage"] = objectRaw["costStorage"]
	mapping["create_time"] = objectRaw["createdTime"]
	mapping["default_quota"] = objectRaw["defaultQuota"]
	mapping["owner"] = objectRaw["owner"]
	mapping["region_id"] = objectRaw["regionId"]
	mapping["status"] = objectRaw["status"]
	mapping["three_tier_model"] = objectRaw["threeTierModel"]
	mapping["type"] = objectRaw["type"]
	mapping["project_name"] = objectRaw["name"]

	ipWhiteListMaps := make([]map[string]interface{}, 0)
	ipWhiteListMap := make(map[string]interface{})
	ipWhiteListRaw := make(map[string]interface{})
	if objectRaw["ipWhiteList"] != nil {
		ipWhiteListRaw = objectRaw["ipWhiteList"].(map[string]interface{})
	}
	if len(ipWhiteListRaw) > 0 {
		ipWhiteListMap["ip_list"] = ipWhiteListRaw["ipList"]
		ipWhiteListMap["vpc_ip_list"] = ipWhiteListRaw["vpcIpList"]

		ipWhiteListMaps = append(ipWhiteListMaps, ipWhiteListMap)
	}
	mapping["ip_white_list"] = ipWhiteListMaps
	propertiesMaps := make([]map[string]interface{}, 0)
	propertiesMap := make(map[string]interface{})
	propertiesRaw := make(map[string]interface{})
	if objectRaw["properties"] != nil {
		propertiesRaw = objectRaw["properties"].(map[string]interface{})
	}
	if len(propertiesRaw) > 0 {
		propertiesMap["allow_full_scan"] = propertiesRaw["allowFullScan"]
		propertiesMap["enable_data_masking"] = propertiesRaw["enableDataMasking"]
		propertiesMap["enable_decimal2"] = propertiesRaw["enableDecimal2"]
		propertiesMap["enable_dr"] = propertiesRaw["enableDr"]
		propertiesMap["enable_tunnel_quota_route"] = propertiesRaw["enableTunnelQuotaRoute"]
		propertiesMap["retention_days"] = propertiesRaw["retentionDays"]
		propertiesMap["sql_metering_max"] = propertiesRaw["sqlMeteringMax"]
		propertiesMap["timezone"] = propertiesRaw["timezone"]
		propertiesMap["tunnel_quota"] = propertiesRaw["tunnelQuota"]
		propertiesMap["type_system"] = propertiesRaw["typeSystem"]

		encryptionMaps := make([]map[string]interface{}, 0)
		encryptionMap := make(map[string]interface{})
		encryptionRaw := make(map[string]interface{})
		if propertiesRaw["encryption"] != nil {
			encryptionRaw = propertiesRaw["encryption"].(map[string]interface{})
		}
		if len(encryptionRaw) > 0 {
			encryptionMap["algorithm"] = encryptionRaw["algorithm"]
			encryptionMap["enable"] = encryptionRaw["enable"]
			encryptionMap["key"] = encryptionRaw["key"]

			encryptionMaps = append(encryptionMaps, encryptionMap)
		}
		propertiesMap["encryption"] = encryptionMaps
		tableLifecycleMaps := make([]map[string]interface{}, 0)
		tableLifecycleMap := make(map[string]interface{})
		tableLifecycleRaw := make(map[string]interface{})
		if propertiesRaw["tableLifecycle"] != nil {
			tableLifecycleRaw = propertiesRaw["tableLifecycle"].(map[string]interface{})
		}
		if len(tableLifecycleRaw) > 0 {
			tableLifecycleMap["type"] = tableLifecycleRaw["type"]
			tableLifecycleMap["value"] = tableLifecycleRaw["value"]

			tableLifecycleMaps = append(tableLifecycleMaps, tableLifecycleMap)
		}
		propertiesMap["table_lifecycle"] = tableLifecycleMaps
		propertiesMaps = append(propertiesMaps, propertiesMap)
	}
	mapping["properties"] = propertiesMaps
	securityPropertiesMaps := make([]map[string]interface{}, 0)
	securityPropertiesMap := make(map[string]interface{})
	securityPropertiesRaw := make(map[string]interface{})
	if objectRaw["securityProperties"] != nil {
		securityPropertiesRaw = objectRaw["securityProperties"].(map[string]interface{})
	}
	if len(securityPropertiesRaw) > 0 {
		securityPropertiesMap["enable_download_privilege"] = securityPropertiesRaw["enableDownloadPrivilege"]
		securityPropertiesMap["label_security"] = securityPropertiesRaw["labelSecurity"]
		securityPropertiesMap["object_creator_has_access_permission"] = securityPropertiesRaw["objectCreatorHasAccessPermission"]
		securityPropertiesMap["object_creator_has_grant_permission"] = securityPropertiesRaw["objectCreatorHasGrantPermission"]
		securityPropertiesMap["using_acl"] = securityPropertiesRaw["usingAcl"]
		securityPropertiesMap["using_policy"] = securityPropertiesRaw["usingPolicy"]

		projectProtectionMaps := make([]map[string]interface{}, 0)
		projectProtectionMap := make(map[string]interface{})
		projectProtectionRaw := make(map[string]interface{})
		if securityPropertiesRaw["projectProtection"] != nil {
			projectProtectionRaw = securityPropertiesRaw["projectProtection"].(map[string]interface{})
		}
		if len(projectProtectionRaw) > 0 {
			projectProtectionMap["exception_policy"] = projectProtectionRaw["exceptionPolicy"]
			projectProtectionMap["protected"] = projectProtectionRaw["protected"]

			projectProtectionMaps = append(projectProtectionMaps, projectProtectionMap)
		}
		securityPropertiesMap["project_protection"] = projectProtectionMaps
		securityPropertiesMaps = append(securityPropertiesMaps, securityPropertiesMap)
	}
	mapping["security_properties"] = securityPropertiesMaps

	objectRaw = getResp

	tagResourcesRaw := make(map[string]interface{})
	if objectRaw["TagResources"] != nil {
		tagResourcesRaw = objectRaw["TagResources"].(map[string]interface{})
	}
	if len(tagResourcesRaw) > 0 {

		tagsMaps, _ := jsonpath.Get("$.TagResources.TagResource", objectRaw)
		mapping["tags"] = tagsToMap(tagsMaps)

		objectRaw = getResp

		dataRaw := make([]interface{}, 0)
		if objectRaw["data"] != nil {
			dataRaw = convertToInterfaceArray(objectRaw["data"])
		}

		mapping["trusted_projects"] = dataRaw
	}

	return mapping, nil
}
