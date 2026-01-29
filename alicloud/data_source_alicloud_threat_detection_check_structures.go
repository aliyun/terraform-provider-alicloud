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

func dataSourceAliCloudThreatDetectionCheckStructures() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudThreatDetectionCheckStructureRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"current_page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"task_sources": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"structures": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"standard_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"standards": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"show_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"requirements": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"sections": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:     schema.TypeInt,
																Computed: true,
															},
															"show_name": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
												"total_check_count": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"id": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"show_name": {
													Type:     schema.TypeString,
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
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAliCloudThreatDetectionCheckStructureRead(d *schema.ResourceData, meta interface{}) error {
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
	action := "GetCheckStructure"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOkExists("current_page"); ok {
		request["CurrentPage"] = v
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("task_sources"); ok {
		taskSourcesMapsArray := convertToInterfaceArray(v)

		request["TaskSources"] = taskSourcesMapsArray
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
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

	resp, _ := jsonpath.Get("$.CheckStructureResponse[*]", response)

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

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["standard_type"] = objectRaw["StandardType"]

		standardsRaw := objectRaw["Standards"]
		standardsMaps := make([]map[string]interface{}, 0)
		if standardsRaw != nil {
			for _, standardsChildRaw := range convertToInterfaceArray(standardsRaw) {
				standardsMap := make(map[string]interface{})
				standardsChildRaw := standardsChildRaw.(map[string]interface{})
				standardsMap["id"] = standardsChildRaw["Id"]
				standardsMap["show_name"] = standardsChildRaw["ShowName"]
				standardsMap["type"] = standardsChildRaw["Type"]

				requirementsRaw := standardsChildRaw["Requirements"]
				requirementsMaps := make([]map[string]interface{}, 0)
				if requirementsRaw != nil {
					for _, requirementsChildRaw := range convertToInterfaceArray(requirementsRaw) {
						requirementsMap := make(map[string]interface{})
						requirementsChildRaw := requirementsChildRaw.(map[string]interface{})
						requirementsMap["id"] = requirementsChildRaw["Id"]
						requirementsMap["show_name"] = requirementsChildRaw["ShowName"]
						requirementsMap["total_check_count"] = requirementsChildRaw["TotalCheckCount"]

						sectionsRaw := requirementsChildRaw["Sections"]
						sectionsMaps := make([]map[string]interface{}, 0)
						if sectionsRaw != nil {
							for _, sectionsChildRaw := range convertToInterfaceArray(sectionsRaw) {
								sectionsMap := make(map[string]interface{})
								sectionsChildRaw := sectionsChildRaw.(map[string]interface{})
								sectionsMap["id"] = sectionsChildRaw["Id"]
								sectionsMap["show_name"] = sectionsChildRaw["ShowName"]

								sectionsMaps = append(sectionsMaps, sectionsMap)
							}
						}
						requirementsMap["sections"] = sectionsMaps
						requirementsMaps = append(requirementsMaps, requirementsMap)
					}
				}
				standardsMap["requirements"] = requirementsMaps
				standardsMaps = append(standardsMaps, standardsMap)
			}
		}
		mapping["standards"] = standardsMaps

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw[""])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("structures", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
