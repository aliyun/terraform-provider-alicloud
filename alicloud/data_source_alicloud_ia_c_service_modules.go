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
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudIaCServiceModules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudIaCServiceModuleRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"module_name": {
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
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"modules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"latest_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"module_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"module_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"group_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"project_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"project_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tag_value": {
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

func dataSourceAliCloudIaCServiceModuleRead(d *schema.ResourceData, meta interface{}) error {
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
	action := fmt.Sprintf("/modules")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)

	if v, ok := d.GetOk("group_id"); ok {
		query["groupId"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("module_name"); ok {
		query["keyword"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("project_id"); ok {
		query["projectId"] = StringPointer(v.(string))
	}

	pageSize := PageSizeLarge
	if v, ok := d.GetOk("page_size"); ok && v.(int) > 0 {
		pageSize = v.(int)
	}
	pageNumber := 1
	if v, ok := d.GetOk("page_number"); ok && v.(int) > 0 {
		pageNumber = v.(int)
	}
	query["pageSize"] = StringPointer(strconv.Itoa(pageSize))
	query["pageNumber"] = StringPointer(strconv.Itoa(pageNumber))
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RoaGet("IaCService", "2021-08-06", action, query, nil, nil)

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

		resp, _ := jsonpath.Get("$.modules[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["moduleId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < pageSize {
			break
		}
		pageNumber++
		query["pageNumber"] = StringPointer(strconv.Itoa(pageNumber))
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["moduleId"]

		mapping["create_time"] = objectRaw["createTime"]
		mapping["description"] = objectRaw["description"]
		mapping["latest_version"] = objectRaw["latestVersion"]
		mapping["module_name"] = objectRaw["name"]
		mapping["source"] = objectRaw["source"]
		mapping["status"] = objectRaw["status"]
		mapping["module_id"] = objectRaw["moduleId"]

		groupInfoRaw := objectRaw["groupInfo"]
		groupInfoMaps := make([]map[string]interface{}, 0)
		if groupInfoRaw != nil {
			groupInfoMap := make(map[string]interface{})
			groupInfoChildRaw := groupInfoRaw.(map[string]interface{})
			groupInfoMap["group_id"] = groupInfoChildRaw["groupId"]
			groupInfoMap["group_name"] = groupInfoChildRaw["groupName"]
			groupInfoMap["project_id"] = groupInfoChildRaw["projectId"]
			groupInfoMap["project_name"] = groupInfoChildRaw["projectName"]

			groupInfoMaps = append(groupInfoMaps, groupInfoMap)
		}
		mapping["group_info"] = groupInfoMaps
		tagsRaw := objectRaw["tags"]
		tagsMaps := make([]map[string]interface{}, 0)
		if tagsRaw != nil {
			for _, tagsChildRaw := range convertToInterfaceArray(tagsRaw) {
				tagsMap := make(map[string]interface{})
				tagsChildRaw := tagsChildRaw.(map[string]interface{})
				tagsMap["tag_key"] = tagsChildRaw["tagKey"]
				tagsMap["tag_value"] = tagsChildRaw["tagValue"]

				tagsMaps = append(tagsMaps, tagsMap)
			}
		}
		mapping["tags"] = tagsMaps

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
	if err := d.Set("modules", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
