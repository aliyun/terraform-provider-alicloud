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

func dataSourceAlicloudThreatDetectionComponents() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudThreatDetectionComponentsRead,
		Schema: map[string]*schema.Schema{
			"component_name": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"name_regex": {
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"lang": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"role_for": {
				Optional: true,
				Type:     schema.TypeInt,
			},
			"ids": {
				Optional: true,
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
			"components": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"component_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"component_alias": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"component_description": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"component_logo": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"component_extension": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"create_time": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"update_time": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"component_actions": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"component_action_name": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"component_action_description": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"input_configs": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"field_name": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"field_type": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"field_description": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"field_display_config": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"default_value": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"required": {
													Computed: true,
													Type:     schema.TypeBool,
												},
											},
										},
									},
									"output_configs": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"field_name": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"field_type": {
													Computed: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
								},
							},
						},
						"component_asset_configs": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field_name": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"field_type": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"field_description": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"required": {
										Computed: true,
										Type:     schema.TypeBool,
									},
									"encrypted": {
										Computed: true,
										Type:     schema.TypeBool,
									},
									"default_value": {
										Computed: true,
										Type:     schema.TypeString,
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

func dataSourceAlicloudThreatDetectionComponentsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	var componentNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		componentNameRegex = r
	}
	if v, ok := d.GetOk("component_name"); ok {
		request["ComponentName"] = v
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOkExists("role_for"); ok {
		request["RoleFor"] = v
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

	var err error
	var objects []interface{}
	var response map[string]interface{}

	request["PageNumber"] = 1
	request["PageSize"] = PageSizeMedium
	for {
		action := "ListComponents"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := client.RpcPost("sophonsoar", "2025-09-03", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			response = resp
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_threat_detection_components", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Components", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Components", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ComponentName"])]; !ok {
					continue
				}
			}
			if componentNameRegex != nil && !componentNameRegex.MatchString(fmt.Sprint(item["ComponentName"])) {
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
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                      fmt.Sprint(object["ComponentName"]),
			"component_name":          object["ComponentName"],
			"component_alias":         object["ComponentAlias"],
			"component_description":   object["ComponentDescription"],
			"component_logo":          object["ComponentLogo"],
			"component_extension":     object["ComponentExtension"],
			"create_time":             formatInt(object["CreateTime"]),
			"update_time":             formatInt(object["UpdateTime"]),
			"component_actions":       expandThreatDetectionComponentActions(object["ComponentActions"]),
			"component_asset_configs": expandThreatDetectionComponentAssetConfigs(object["ComponentAssetConfigs"]),
		}

		ids = append(ids, fmt.Sprint(object["ComponentName"]))

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("components", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return WrapError(err)
		}
	}
	return nil
}

func expandThreatDetectionComponentActions(raw interface{}) []map[string]interface{} {
	items, _ := raw.([]interface{})
	s := make([]map[string]interface{}, 0, len(items))
	for _, v := range items {
		item, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		s = append(s, map[string]interface{}{
			"component_action_name":        item["ComponentActionName"],
			"component_action_description": item["ComponentActionDescription"],
			"input_configs":                expandThreatDetectionComponentFieldConfigs(item["InputConfigs"], true),
			"output_configs":               expandThreatDetectionComponentFieldConfigs(item["OutputConfigs"], false),
		})
	}
	return s
}

func expandThreatDetectionComponentFieldConfigs(raw interface{}, withDisplayAndDefault bool) []map[string]interface{} {
	items, _ := raw.([]interface{})
	s := make([]map[string]interface{}, 0, len(items))
	for _, v := range items {
		item, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		entry := map[string]interface{}{
			"field_name": item["FieldName"],
			"field_type": item["FieldType"],
		}
		if withDisplayAndDefault {
			entry["field_description"] = item["FieldDescription"]
			entry["field_display_config"] = item["FieldDisplayConfig"]
			entry["default_value"] = item["DefaultValue"]
			entry["required"] = item["Required"]
		}
		s = append(s, entry)
	}
	return s
}

func expandThreatDetectionComponentAssetConfigs(raw interface{}) []map[string]interface{} {
	items, _ := raw.([]interface{})
	s := make([]map[string]interface{}, 0, len(items))
	for _, v := range items {
		item, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		s = append(s, map[string]interface{}{
			"field_name":        item["FieldName"],
			"field_type":        item["FieldType"],
			"field_description": item["FieldDescription"],
			"required":          item["Required"],
			"encrypted":         item["Encrypted"],
			"default_value":     item["DefaultValue"],
		})
	}
	return s
}
