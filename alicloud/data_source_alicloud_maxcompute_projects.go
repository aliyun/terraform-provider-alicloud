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

func dataSourceAlicloudMaxcomputeProjects() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudMaxcomputeProjectsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name_regex": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.ValidateRegexp,
			},
			"names": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"projects": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"comment": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"default_quota": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"ip_white_list": {
							Computed: true,
							Type:     schema.TypeList,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_list": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"vpc_ip_list": {
										Computed: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"owner": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"project_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"properties": {
							Computed: true,
							Type:     schema.TypeList,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allow_full_scan": {
										Computed: true,
										Type:     schema.TypeBool,
									},
									"enable_decimal2": {
										Computed: true,
										Type:     schema.TypeBool,
									},
									"encryption": {
										Computed: true,
										Type:     schema.TypeList,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"algorithm": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"enable": {
													Computed: true,
													Type:     schema.TypeBool,
												},
												"key": {
													Computed: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"retention_days": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"sql_metering_max": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"table_lifecycle": {
										Computed: true,
										Type:     schema.TypeList,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"value": {
													Computed: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"timezone": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"type_system": {
										Computed: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"security_properties": {
							Computed: true,
							Type:     schema.TypeList,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable_download_privilege": {
										Computed: true,
										Type:     schema.TypeBool,
									},
									"label_security": {
										Computed: true,
										Type:     schema.TypeBool,
									},
									"object_creator_has_access_permission": {
										Computed: true,
										Type:     schema.TypeBool,
									},
									"object_creator_has_grant_permission": {
										Computed: true,
										Type:     schema.TypeBool,
									},
									"project_protection": {
										Computed: true,
										Type:     schema.TypeList,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"exception_policy": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"protected": {
													Computed: true,
													Type:     schema.TypeBool,
												},
											},
										},
									},
									"using_acl": {
										Computed: true,
										Type:     schema.TypeBool,
									},
									"using_policy": {
										Computed: true,
										Type:     schema.TypeBool,
									},
								},
							},
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"type": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudMaxcomputeProjectsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]*string)

	request["maxItem"] = StringPointer(strconv.Itoa(PageSizeLarge))

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var projectNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		projectNameRegex = r
	}

	conn, err := client.NewOdpsClient()
	if err != nil {
		return WrapError(err)
	}
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "/api/v1/projects"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer("2022-01-04"), nil, StringPointer("GET"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
			if err != nil {
				return resource.NonRetryableError(err)
			}
			response = resp
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_maxcompute_projects", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.body.data.projects", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.data.projects", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["name"])]; !ok {
					continue
				}
			}

			if projectNameRegex != nil && !projectNameRegex.MatchString(fmt.Sprint(item["name"])) {
				continue
			}
			objects = append(objects, item)
		}
		if nextToken, ok := response["marker"].(string); ok && nextToken != "" {
			request["marker"] = StringPointer(nextToken)
		} else {
			break
		}
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":            fmt.Sprint(object["name"]),
			"default_quota": object["defaultQuota"],
			"owner":         object["owner"],
			"project_name":  object["name"],
			"status":        object["status"],
			"type":          object["type"],
			"comment":       object["comment"],
		}
		ipWhiteList, err := jsonpath.Get("$.ipWhiteList", object)
		if err == nil {
			ipWhiteList60Maps := make([]map[string]interface{}, 0)
			ipWhiteList60Map := make(map[string]interface{})
			ipWhiteList60Raw := ipWhiteList.(map[string]interface{})
			ipWhiteList60Map["ip_list"] = ipWhiteList60Raw["ipList"]
			ipWhiteList60Map["vpc_ip_list"] = ipWhiteList60Raw["vpcIpList"]
			ipWhiteList60Maps = append(ipWhiteList60Maps, ipWhiteList60Map)
			mapping["ip_white_list"] = ipWhiteList60Maps
		}
		securityPropertiesMaps := make([]map[string]interface{}, 0)
		securityPropertiesMap := make(map[string]interface{})
		securityPropertiesRaw := object["securityProperties"].(map[string]interface{})
		securityPropertiesMap["enable_download_privilege"] = securityPropertiesRaw["enableDownloadPrivilege"]
		securityPropertiesMap["using_acl"] = securityPropertiesRaw["usingAcl"]
		securityPropertiesMap["using_policy"] = securityPropertiesRaw["usingPolicy"]
		securityPropertiesMap["label_security"] = securityPropertiesRaw["labelSecurity"]
		securityPropertiesMap["object_creator_has_access_permission"] = securityPropertiesRaw["objectCreatorHasAccessPermission"]
		securityPropertiesMap["object_creator_has_grant_permission"] = securityPropertiesRaw["objectCreatorHasGrantPermission"]
		projectProtectionMaps := make([]map[string]interface{}, 0)
		projectProtectionMap := make(map[string]interface{})
		projectProtectionRaw := securityPropertiesRaw["projectProtection"].(map[string]interface{})
		projectProtectionMap["protected"] = projectProtectionRaw["protected"]
		exceptionPolicy, err := jsonpath.Get("$.exceptionPolicy", projectProtectionRaw)
		if err == nil {
			projectProtectionMap["exceptionPolicy"] = exceptionPolicy
		}
		projectProtectionMaps = append(projectProtectionMaps, projectProtectionMap)
		securityPropertiesMap["project_protection"] = projectProtectionMaps
		securityPropertiesMaps = append(securityPropertiesMaps, securityPropertiesMap)
		mapping["security_properties"] = securityPropertiesMaps
		propertiesMaps := make([]map[string]interface{}, 0)
		propertiesMap := make(map[string]interface{})
		propertiesRaw := object["properties"].(map[string]interface{})
		propertiesMap["timezone"] = propertiesRaw["timezone"]
		propertiesMap["allow_full_scan"] = propertiesRaw["allowFullScan"]
		propertiesMap["enable_decimal2"] = propertiesRaw["enableDecimal2"]
		propertiesMap["retention_days"] = propertiesRaw["retentionDays"]
		propertiesMap["sql_metering_max"] = propertiesRaw["sqlMeteringMax"]
		propertiesMap["type_system"] = propertiesRaw["typeSystem"]
		tableLifecycleMaps := make([]map[string]interface{}, 0)
		tableLifecycleMap := make(map[string]interface{})
		tableLifecycleRaw := propertiesRaw["tableLifecycle"].(map[string]interface{})
		tableLifecycleMap["type"] = tableLifecycleRaw["type"]
		tableLifecycleMap["value"] = tableLifecycleRaw["value"]
		tableLifecycleMaps = append(tableLifecycleMaps, tableLifecycleMap)
		propertiesMap["table_lifecycle"] = tableLifecycleMaps
		encryptionMaps := make([]map[string]interface{}, 0)
		encryptionMap := make(map[string]interface{})
		encryptionRaw := propertiesRaw["encryption"].(map[string]interface{})
		encryptionMap["enable"] = encryptionRaw["enable"]
		encryptionMap["algorithm"] = encryptionRaw["algorithm"]
		encryptionMap["key"] = encryptionRaw["key"]
		encryptionMaps = append(encryptionMaps, encryptionMap)
		propertiesMap["encryption"] = encryptionMaps
		propertiesMaps = append(propertiesMaps, propertiesMap)
		mapping["properties"] = propertiesMaps
		ids = append(ids, fmt.Sprint(object["name"]))
		names = append(names, object["name"])
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
