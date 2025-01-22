package alicloud

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudMaxComputeProjects() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudMaxComputeProjectRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
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
							MaxItems: 1,
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
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"timezone": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"sql_metering_max": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type_system": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"table_lifecycle": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
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
									"retention_days": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"encryption": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
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
									"enable_decimal2": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"security_properties": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
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
										MaxItems: 1,
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
						"type": {
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
	action := fmt.Sprintf("/api/v1/projects")
	conn, err := client.NewOdpsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query := make(map[string]*string)
	query["maxItem"] = StringPointer(strconv.Itoa(PageSizeLarge))

	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2022-01-04"), nil, StringPointer("GET"), StringPointer("AK"), StringPointer(action), query, nil, nil, &runtime)

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

		resp, _ := jsonpath.Get("$.body.data.projects[*]", response)
		marker, _ := jsonpath.Get("$.body.data.marker", response)
		for _, v := range resp.([]interface{}) {
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

		if nextToken, ok := marker.(string); ok && nextToken != "" {
			query["marker"] = StringPointer(nextToken)
		} else {
			break
		}
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)

	maxComputeServiceV2 := MaxComputeServiceV2{client}
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["comment"] = objectRaw["comment"]
		mapping["create_time"] = objectRaw["createdTime"]
		mapping["default_quota"] = objectRaw["defaultQuota"]
		mapping["owner"] = objectRaw["owner"]
		mapping["status"] = objectRaw["status"]
		mapping["type"] = objectRaw["type"]
		mapping["project_name"] = objectRaw["name"]

		objectDetail, _ := maxComputeServiceV2.DescribeMaxComputeProject(objectRaw["name"].(string))
		objectRaw, _ := jsonpath.Get("$", objectDetail)

		ipWhiteListMaps := make([]map[string]interface{}, 0)
		ipWhiteListMap := make(map[string]interface{})
		ipWhiteList2Raw := make(map[string]interface{})
		if objectRaw.(map[string]interface{})["ipWhiteList"] != nil {
			ipWhiteList2Raw = objectRaw.(map[string]interface{})["ipWhiteList"].(map[string]interface{})
		}
		if len(ipWhiteList2Raw) > 0 {
			ipWhiteListMap["ip_list"] = ipWhiteList2Raw["ipList"]
			ipWhiteListMap["vpc_ip_list"] = ipWhiteList2Raw["vpcIpList"]

			ipWhiteListMaps = append(ipWhiteListMaps, ipWhiteListMap)
		}
		mapping["ip_white_list"] = ipWhiteListMaps
		propertiesMaps := make([]map[string]interface{}, 0)
		propertiesMap := make(map[string]interface{})
		properties2Raw := make(map[string]interface{})
		if objectRaw.(map[string]interface{})["properties"] != nil {
			properties2Raw = objectRaw.(map[string]interface{})["properties"].(map[string]interface{})
		}
		if len(properties2Raw) > 0 {
			propertiesMap["allow_full_scan"] = properties2Raw["allowFullScan"]
			propertiesMap["enable_decimal2"] = properties2Raw["enableDecimal2"]
			propertiesMap["retention_days"] = properties2Raw["retentionDays"]
			propertiesMap["sql_metering_max"] = properties2Raw["sqlMeteringMax"]
			propertiesMap["timezone"] = properties2Raw["timezone"]
			propertiesMap["type_system"] = properties2Raw["typeSystem"]

			encryptionMaps := make([]map[string]interface{}, 0)
			encryptionMap := make(map[string]interface{})
			encryption2Raw := make(map[string]interface{})
			if properties2Raw["encryption"] != nil {
				encryption2Raw = properties2Raw["encryption"].(map[string]interface{})
			}
			if len(encryption2Raw) > 0 {
				encryptionMap["algorithm"] = encryption2Raw["algorithm"]
				encryptionMap["enable"] = encryption2Raw["enable"]
				encryptionMap["key"] = encryption2Raw["key"]

				encryptionMaps = append(encryptionMaps, encryptionMap)
			}
			propertiesMap["encryption"] = encryptionMaps
			tableLifecycleMaps := make([]map[string]interface{}, 0)
			tableLifecycleMap := make(map[string]interface{})
			tableLifecycle2Raw := make(map[string]interface{})
			if properties2Raw["tableLifecycle"] != nil {
				tableLifecycle2Raw = properties2Raw["tableLifecycle"].(map[string]interface{})
			}
			if len(tableLifecycle2Raw) > 0 {
				tableLifecycleMap["type"] = tableLifecycle2Raw["type"]
				tableLifecycleMap["value"] = tableLifecycle2Raw["value"]

				tableLifecycleMaps = append(tableLifecycleMaps, tableLifecycleMap)
			}
			propertiesMap["table_lifecycle"] = tableLifecycleMaps
			propertiesMaps = append(propertiesMaps, propertiesMap)
		}
		mapping["properties"] = propertiesMaps
		securityPropertiesMaps := make([]map[string]interface{}, 0)
		securityPropertiesMap := make(map[string]interface{})
		securityProperties2Raw := make(map[string]interface{})
		if objectRaw.(map[string]interface{})["securityProperties"] != nil {
			securityProperties2Raw = objectRaw.(map[string]interface{})["securityProperties"].(map[string]interface{})
		}
		if len(securityProperties2Raw) > 0 {
			securityPropertiesMap["enable_download_privilege"] = securityProperties2Raw["enableDownloadPrivilege"]
			securityPropertiesMap["label_security"] = securityProperties2Raw["labelSecurity"]
			securityPropertiesMap["object_creator_has_access_permission"] = securityProperties2Raw["objectCreatorHasAccessPermission"]
			securityPropertiesMap["object_creator_has_grant_permission"] = securityProperties2Raw["objectCreatorHasGrantPermission"]
			securityPropertiesMap["using_acl"] = securityProperties2Raw["usingAcl"]
			securityPropertiesMap["using_policy"] = securityProperties2Raw["usingPolicy"]

			projectProtectionMaps := make([]map[string]interface{}, 0)
			projectProtectionMap := make(map[string]interface{})
			projectProtection2Raw := make(map[string]interface{})
			if securityProperties2Raw["projectProtection"] != nil {
				projectProtection2Raw = securityProperties2Raw["projectProtection"].(map[string]interface{})
			}
			if len(projectProtection2Raw) > 0 {
				projectProtectionMap["exception_policy"] = projectProtection2Raw["exceptionPolicy"]
				projectProtectionMap["protected"] = projectProtection2Raw["protected"]

				projectProtectionMaps = append(projectProtectionMaps, projectProtectionMap)
			}
			securityPropertiesMap["project_protection"] = projectProtectionMaps
			securityPropertiesMaps = append(securityPropertiesMaps, securityPropertiesMap)
		}
		mapping["security_properties"] = securityPropertiesMaps

		ids = append(ids, fmt.Sprint(mapping["project_name"]))
		names = append(names, fmt.Sprint(mapping["project_name"]))
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
