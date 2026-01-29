// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudThreatDetectionCheckItemConfigs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudThreatDetectionCheckItemConfigRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"task_sources": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"configs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"check_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"check_show_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"check_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"custom_configs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"default_value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type_define": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"show_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"description": {
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
						"estimated_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_sub_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"risk_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"section_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"vendor": {
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

func dataSourceAliCloudThreatDetectionCheckItemConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

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

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "ListCheckItem"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("task_sources"); ok {
		taskSourcesMapsArray := convertToInterfaceArray(v)

		request["TaskSources"] = taskSourcesMapsArray
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)

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

		resp, _ := jsonpath.Get("$.CheckItems[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint()]; !ok {
					continue
				}
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
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["check_show_name"] = objectRaw["CheckShowName"]
		mapping["check_type"] = objectRaw["CheckType"]
		mapping["risk_level"] = objectRaw["RiskLevel"]
		mapping["check_id"] = objectRaw["CheckId"]
		mapping["estimated_count"] = objectRaw["EstimatedCount"]
		mapping["instance_sub_type"] = objectRaw["InstanceSubType"]
		mapping["instance_type"] = objectRaw["InstanceType"]
		mapping["vendor"] = objectRaw["Vendor"]

		sectionIdsRaw := make([]interface{}, 0)
		if objectRaw["SectionIds"] != nil {
			sectionIdsRaw = convertToInterfaceArray(objectRaw["SectionIds"])
		}

		mapping["section_ids"] = sectionIdsRaw
		customConfigsRaw := objectRaw["CustomConfigs"]
		customConfigsMaps := make([]map[string]interface{}, 0)
		if customConfigsRaw != nil {
			for _, customConfigsChildRaw := range convertToInterfaceArray(customConfigsRaw) {
				customConfigsMap := make(map[string]interface{})
				customConfigsChildRaw := customConfigsChildRaw.(map[string]interface{})
				customConfigsMap["default_value"] = customConfigsChildRaw["DefaultValue"]
				customConfigsMap["name"] = customConfigsChildRaw["Name"]
				customConfigsMap["show_name"] = customConfigsChildRaw["ShowName"]
				customConfigsMap["type_define"] = customConfigsChildRaw["TypeDefine"]
				customConfigsMap["value"] = customConfigsChildRaw["Value"]

				customConfigsMaps = append(customConfigsMaps, customConfigsMap)
			}
		}
		mapping["custom_configs"] = customConfigsMaps
		descriptionMaps := make([]map[string]interface{}, 0)
		descriptionMap := make(map[string]interface{})
		descriptionRaw := make(map[string]interface{})
		if objectRaw["Description"] != nil {
			descriptionRaw = objectRaw["Description"].(map[string]interface{})
		}
		if len(descriptionRaw) > 0 {
			descriptionMap["type"] = descriptionRaw["Type"]
			descriptionMap["value"] = descriptionRaw["Value"]

			descriptionMaps = append(descriptionMaps, descriptionMap)
		}
		mapping["description"] = descriptionMaps

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw[""])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("configs", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
