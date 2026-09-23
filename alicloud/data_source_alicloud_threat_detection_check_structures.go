package alicloud

import (
	"encoding/json"
	"fmt"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type checkStructureResponse struct {
	CheckStructureResponse []checkStructureItem `json:"CheckStructureResponse"`
}

type checkStructureItem struct {
	StandardType string              `json:"StandardType"`
	Standards    []checkStructureStd `json:"Standards"`
}

type checkStructureStd struct {
	Id           int                 `json:"Id"`
	ShowName     string              `json:"ShowName"`
	Type         string              `json:"Type"`
	Requirements []checkStructureReq `json:"Requirements"`
}

type checkStructureReq struct {
	Id              int                 `json:"Id"`
	ShowName        string              `json:"ShowName"`
	TotalCheckCount int                 `json:"TotalCheckCount"`
	Sections        []checkStructureSec `json:"Sections"`
}

type checkStructureSec struct {
	Id       int    `json:"Id"`
	ShowName string `json:"ShowName"`
}

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
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntAtLeast(1),
			},
			"page_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntAtLeast(1),
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

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	action := "GetCheckStructure"
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	query := map[string]interface{}{}
	var response map[string]interface{}
	var err error

	if v := d.Get("current_page").(int); v > 0 {
		request["CurrentPage"] = v
	}
	if v := d.Get("page_size").(int); v > 0 {
		request["PageSize"] = v
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("task_sources"); ok {
		request["TaskSources"] = convertToInterfaceArray(v)
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

	respBytes, err := json.Marshal(response)
	if err != nil {
		return WrapError(err)
	}
	var parsed checkStructureResponse
	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		return WrapError(err)
	}

	objects := make([]checkStructureItem, 0, len(parsed.CheckStructureResponse))
	for _, item := range parsed.CheckStructureResponse {
		matched := make([]checkStructureStd, 0, len(item.Standards))
		for _, std := range item.Standards {
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(std.Id)]; !ok {
					continue
				}
			}
			matched = append(matched, std)
		}
		if len(matched) == 0 {
			continue
		}
		objects = append(objects, checkStructureItem{
			StandardType: item.StandardType,
			Standards:    matched,
		})
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0, len(objects))
	for _, object := range objects {
		standardsMaps := make([]map[string]interface{}, 0, len(object.Standards))
		for _, std := range object.Standards {
			requirementsMaps := make([]map[string]interface{}, 0, len(std.Requirements))
			for _, req := range std.Requirements {
				sectionsMaps := make([]map[string]interface{}, 0, len(req.Sections))
				for _, sec := range req.Sections {
					sectionsMaps = append(sectionsMaps, map[string]interface{}{
						"id":        sec.Id,
						"show_name": sec.ShowName,
					})
				}
				requirementsMaps = append(requirementsMaps, map[string]interface{}{
					"id":                req.Id,
					"show_name":         req.ShowName,
					"total_check_count": req.TotalCheckCount,
					"sections":          sectionsMaps,
				})
			}
			standardsMaps = append(standardsMaps, map[string]interface{}{
				"id":           std.Id,
				"show_name":    std.ShowName,
				"type":         std.Type,
				"requirements": requirementsMaps,
			})
			ids = append(ids, fmt.Sprint(std.Id))
		}

		s = append(s, map[string]interface{}{
			"standard_type": object.StandardType,
			"standards":     standardsMaps,
		})
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
