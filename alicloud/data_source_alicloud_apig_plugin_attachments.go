// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudApigPluginAttachments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudApigPluginAttachmentRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"attach_resource_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attach_resource_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"environment_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"plugin_info": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"plugin_config": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"plugin_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"gateway_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"attachments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attach_resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attach_resource_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"attach_resource_names": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"attach_resource_parent_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"attach_resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"environment_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plugin_attachment_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plugin_class_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"execute_priority": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"execute_stage": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"direction": {
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
						"plugin_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"plugin_config": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"plugin_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"gateway_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
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

func dataSourceAliCloudApigPluginAttachmentRead(d *schema.ResourceData, meta interface{}) error {
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
	var query map[string]*string
	// ListPluginAttachments
	action := fmt.Sprintf("/v1/plugin-attachments")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)

	if v, ok := d.GetOk("attach_resource_id"); ok {
		query["attachResourceId"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("attach_resource_type"); ok {
		query["attachResourceType"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("environment_id"); ok {
		query["environmentId"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("plugin_info.0.gateway_id"); ok {
		query["gatewayId"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("plugin_info.0.plugin_id"); ok {
		query["pluginId"] = StringPointer(v.(string))
	}

	query["pageSize"] = StringPointer(strconv.Itoa(PageSizeLarge))
	query["pageNumber"] = StringPointer("1")
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

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
			if NotFoundError(err) {
				break
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.data.items[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["pluginAttachmentId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < PageSizeLarge {
			break
		}
		pageNum, _ := strconv.Atoi(*query["pageNumber"])
		query["pageNumber"] = StringPointer(strconv.Itoa(pageNum + 1))
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["pluginAttachmentId"]
		mapping["attach_resource_type"] = objectRaw["attachResourceType"]
		mapping["enable"] = objectRaw["enable"]
		mapping["plugin_attachment_id"] = objectRaw["pluginAttachmentId"]

		environmentInfoRawObj, _ := jsonpath.Get("$.environmentInfo", objectRaw)
		environmentInfoRaw := make(map[string]interface{})
		if environmentInfoRawObj != nil {
			environmentInfoRaw = environmentInfoRawObj.(map[string]interface{})
		}
		mapping["environment_id"] = environmentInfoRaw["environmentId"]

		resourceInfosRawObj, _ := jsonpath.Get("$.resourceInfos[*]", objectRaw)
		resourceInfosRaw := make([]interface{}, 0)
		if resourceInfosRawObj != nil {
			resourceInfosRaw = convertToInterfaceArray(resourceInfosRawObj)
		}

		attachResourceIds := make([]string, 0)
		attachResourceNames := make([]string, 0)
		for _, item := range resourceInfosRaw {
			if info, ok := item.(map[string]interface{}); ok {
				if v, ok := info["resourceId"]; ok && fmt.Sprint(v) != "" {
					attachResourceIds = append(attachResourceIds, fmt.Sprint(v))
				}
				if v, ok := info["resourceName"]; ok && fmt.Sprint(v) != "" {
					attachResourceNames = append(attachResourceNames, fmt.Sprint(v))
				}
			}
		}
		mapping["attach_resource_ids"] = attachResourceIds
		mapping["attach_resource_names"] = attachResourceNames
		if len(attachResourceIds) > 0 {
			mapping["attach_resource_id"] = attachResourceIds[0]
		}

		parentResourceInfoRawObj, _ := jsonpath.Get("$.parentResourceInfo", objectRaw)
		parentResourceInfoRaw := make(map[string]interface{})
		if parentResourceInfoRawObj != nil {
			parentResourceInfoRaw = parentResourceInfoRawObj.(map[string]interface{})
		}
		attachResourceParentIds := make([]string, 0)
		if apiInfoRaw, ok := parentResourceInfoRaw["apiInfo"].(map[string]interface{}); ok {
			if v, ok := apiInfoRaw["httpApiId"]; ok && fmt.Sprint(v) != "" {
				attachResourceParentIds = append(attachResourceParentIds, fmt.Sprint(v))
			}
		}
		mapping["attach_resource_parent_ids"] = attachResourceParentIds

		pluginClassInfoMaps := make([]map[string]interface{}, 0)
		pluginClassInfoMap := make(map[string]interface{})
		pluginClassInfoRaw := make(map[string]interface{})
		if objectRaw["pluginClassInfo"] != nil {
			pluginClassInfoRaw = objectRaw["pluginClassInfo"].(map[string]interface{})
		}
		if len(pluginClassInfoRaw) > 0 {
			pluginClassInfoMap["direction"] = pluginClassInfoRaw["mode"]
			pluginClassInfoMap["execute_priority"] = pluginClassInfoRaw["executePriority"]
			pluginClassInfoMap["execute_stage"] = pluginClassInfoRaw["executeStage"]
			pluginClassInfoMap["name"] = pluginClassInfoRaw["name"]
			pluginClassInfoMap["type"] = pluginClassInfoRaw["type"]

			pluginClassInfoMaps = append(pluginClassInfoMaps, pluginClassInfoMap)
		}
		mapping["plugin_class_info"] = pluginClassInfoMaps

		pluginInfoMaps := make([]map[string]interface{}, 0)
		pluginInfoMap := make(map[string]interface{})

		pluginInfoMap["plugin_config"] = objectRaw["pluginConfig"]
		pluginInfoMap["plugin_id"] = objectRaw["pluginId"]

		gatewayInfoRawObj, _ := jsonpath.Get("$.environmentInfo.gatewayInfo", objectRaw)
		gatewayInfoRaw := make(map[string]interface{})
		if gatewayInfoRawObj != nil {
			gatewayInfoRaw = gatewayInfoRawObj.(map[string]interface{})
		}
		pluginInfoMap["gateway_id"] = gatewayInfoRaw["gatewayId"]

		pluginInfoMaps = append(pluginInfoMaps, pluginInfoMap)
		mapping["plugin_info"] = pluginInfoMaps

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("attachments", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
