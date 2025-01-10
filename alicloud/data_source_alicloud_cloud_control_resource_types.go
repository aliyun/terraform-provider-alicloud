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

func dataSourceAliCloudCloudControlResourceTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCloudControlResourceTypeRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"product": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_only_properties": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"delete_only_properties": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"filter_properties": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"get_only_properties": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"get_response_properties": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"handlers": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"delete": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"permissions": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"create": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"permissions": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"get": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"permissions": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"permissions": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"update": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"permissions": {
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
						"info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"delivery_scope": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"charge_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"title": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"list_only_properties": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"list_response_properties": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"primary_identifier": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"properties": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_properties": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"read_only_properties": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"required": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sensitive_info_properties": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"update_only_properties": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"update_type_properties": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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
				ForceNew: true,
			},
		},
	}
}

func dataSourceAliCloudCloudControlResourceTypeRead(d *schema.ResourceData, meta interface{}) error {
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
	action := fmt.Sprintf("/api/v1/providers/%s/products/%s/resourceTypes", "aliyun", d.Get("product").(string))
	conn, err := client.NewCloudcontrolClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query := make(map[string]*string)
	body := make(map[string]interface{})
	request["product"] = d.Get("product")
	request["provider"] = "aliyun"

	body = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	query["MaxResults"] = StringPointer("50")
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2022-08-30"), nil, StringPointer("GET"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)

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

		resp, _ := jsonpath.Get("$.body.resourceTypes[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["resourceType"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		nextToken, _ := jsonpath.Get("$.body.nextToken", response)
		if nextToken != nil && nextToken != "" {
			query["nextToken"] = StringPointer(fmt.Sprint(nextToken))
		} else {
			break
		}
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["resourceType"]

		mapping["primary_identifier"] = objectRaw["primaryIdentifier"]
		mapping["product"] = objectRaw["product"]
		mapping["properties"] = convertObjectToJsonString(objectRaw["properties"])
		mapping["resource_type"] = objectRaw["resourceType"]

		createOnlyProperties1Raw := make([]interface{}, 0)
		if objectRaw["createOnlyProperties"] != nil {
			createOnlyProperties1Raw = objectRaw["createOnlyProperties"].([]interface{})
		}

		mapping["create_only_properties"] = createOnlyProperties1Raw
		deleteOnlyProperties1Raw := make([]interface{}, 0)
		if objectRaw["deleteOnlyProperties"] != nil {
			deleteOnlyProperties1Raw = objectRaw["deleteOnlyProperties"].([]interface{})
		}

		mapping["delete_only_properties"] = deleteOnlyProperties1Raw
		filterProperties1Raw := make([]interface{}, 0)
		if objectRaw["filterProperties"] != nil {
			filterProperties1Raw = objectRaw["filterProperties"].([]interface{})
		}

		mapping["filter_properties"] = filterProperties1Raw
		getOnlyProperties1Raw := make([]interface{}, 0)
		if objectRaw["getOnlyProperties"] != nil {
			getOnlyProperties1Raw = objectRaw["getOnlyProperties"].([]interface{})
		}

		mapping["get_only_properties"] = getOnlyProperties1Raw
		getResponseProperties1Raw := make([]interface{}, 0)
		if objectRaw["getResponseProperties"] != nil {
			getResponseProperties1Raw = objectRaw["getResponseProperties"].([]interface{})
		}

		mapping["get_response_properties"] = getResponseProperties1Raw
		handlersMaps := make([]map[string]interface{}, 0)
		handlersMap := make(map[string]interface{})
		handlers1Raw := make(map[string]interface{})
		if objectRaw["handlers"] != nil {
			handlers1Raw = objectRaw["handlers"].(map[string]interface{})
		}
		if len(handlers1Raw) > 0 {

			createMaps := make([]map[string]interface{}, 0)
			createMap := make(map[string]interface{})
			permissions5Raw, _ := jsonpath.Get("$.handlers.create.permissions", objectRaw)

			createMap["permissions"] = permissions5Raw
			createMaps = append(createMaps, createMap)
			handlersMap["create"] = createMaps
			deleteMaps := make([]map[string]interface{}, 0)
			deleteMap := make(map[string]interface{})
			permissions6Raw, _ := jsonpath.Get("$.handlers.delete.permissions", objectRaw)

			deleteMap["permissions"] = permissions6Raw
			deleteMaps = append(deleteMaps, deleteMap)
			handlersMap["delete"] = deleteMaps
			getMaps := make([]map[string]interface{}, 0)
			getMap := make(map[string]interface{})
			permissions7Raw, _ := jsonpath.Get("$.handlers.get.permissions", objectRaw)

			getMap["permissions"] = permissions7Raw
			getMaps = append(getMaps, getMap)
			handlersMap["get"] = getMaps
			listMaps := make([]map[string]interface{}, 0)
			listMap := make(map[string]interface{})
			permissions8Raw, _ := jsonpath.Get("$.handlers.list.permissions", objectRaw)

			listMap["permissions"] = permissions8Raw
			listMaps = append(listMaps, listMap)
			handlersMap["list"] = listMaps
			updateMaps := make([]map[string]interface{}, 0)
			updateMap := make(map[string]interface{})
			permissions9Raw, _ := jsonpath.Get("$.handlers.update.permissions", objectRaw)

			updateMap["permissions"] = permissions9Raw
			updateMaps = append(updateMaps, updateMap)
			handlersMap["update"] = updateMaps
			handlersMaps = append(handlersMaps, handlersMap)
		}
		mapping["handlers"] = handlersMaps
		infoMaps := make([]map[string]interface{}, 0)
		infoMap := make(map[string]interface{})
		info1Raw := make(map[string]interface{})
		if objectRaw["info"] != nil {
			info1Raw = objectRaw["info"].(map[string]interface{})
		}
		if len(info1Raw) > 0 {
			infoMap["charge_type"] = info1Raw["chargeType"]
			infoMap["delivery_scope"] = info1Raw["deliveryScope"]
			infoMap["description"] = info1Raw["description"]
			infoMap["title"] = info1Raw["title"]

			infoMaps = append(infoMaps, infoMap)
		}
		mapping["info"] = infoMaps
		listOnlyProperties1Raw := make([]interface{}, 0)
		if objectRaw["listOnlyProperties"] != nil {
			listOnlyProperties1Raw = objectRaw["listOnlyProperties"].([]interface{})
		}

		mapping["list_only_properties"] = listOnlyProperties1Raw
		listResponseProperties1Raw := make([]interface{}, 0)
		if objectRaw["listResponseProperties"] != nil {
			listResponseProperties1Raw = objectRaw["listResponseProperties"].([]interface{})
		}

		mapping["list_response_properties"] = listResponseProperties1Raw
		publicProperties1Raw := make([]interface{}, 0)
		if objectRaw["publicProperties"] != nil {
			publicProperties1Raw = objectRaw["publicProperties"].([]interface{})
		}

		mapping["public_properties"] = publicProperties1Raw
		readOnlyProperties1Raw := make([]interface{}, 0)
		if objectRaw["readOnlyProperties"] != nil {
			readOnlyProperties1Raw = objectRaw["readOnlyProperties"].([]interface{})
		}

		mapping["read_only_properties"] = readOnlyProperties1Raw
		required1Raw := make([]interface{}, 0)
		if objectRaw["required"] != nil {
			required1Raw = objectRaw["required"].([]interface{})
		}

		mapping["required"] = required1Raw
		sensitiveInfoProperties1Raw := make([]interface{}, 0)
		if objectRaw["sensitiveInfoProperties"] != nil {
			sensitiveInfoProperties1Raw = objectRaw["sensitiveInfoProperties"].([]interface{})
		}

		mapping["sensitive_info_properties"] = sensitiveInfoProperties1Raw
		updateOnlyProperties1Raw := make([]interface{}, 0)
		if objectRaw["updateOnlyProperties"] != nil {
			updateOnlyProperties1Raw = objectRaw["updateOnlyProperties"].([]interface{})
		}

		mapping["update_only_properties"] = updateOnlyProperties1Raw
		updateTypeProperties1Raw := make([]interface{}, 0)
		if objectRaw["updateTypeProperties"] != nil {
			updateTypeProperties1Raw = objectRaw["updateTypeProperties"].([]interface{})
		}

		mapping["update_type_properties"] = updateTypeProperties1Raw

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["ResourceTypeName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("types", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
