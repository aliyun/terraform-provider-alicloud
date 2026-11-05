// Package apig. This file is generated automatically. Please do not modify it manually, thank you!
package apig

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourcePluginClasses() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePluginClassesRead,
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
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"classes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plugin_class_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plugin_class_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"document": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"wasm_language": {
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
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourcePluginClassesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
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
	// ListPluginClasses
	action := fmt.Sprintf("/v1/plugin-classes")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)

	if v, ok := d.GetOk("type"); ok {
		query["type"] = helper.StringPointer(v.(string))
	}

	query["pageSize"] = helper.StringPointer(strconv.Itoa(helper.PageSizeLarge))
	query["pageNumber"] = helper.StringPointer("1")
	for {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
			response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_apig_plugin_classes", action, errs.AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.data.items", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["pluginClassId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < helper.PageSizeLarge {
			break
		}
		pageNum, _ := strconv.Atoi(*query["pageNumber"])
		query["pageNumber"] = helper.StringPointer(strconv.Itoa(pageNum + 1))
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["pluginClassId"]
		mapping["plugin_class_id"] = objectRaw["pluginClassId"]
		mapping["plugin_class_name"] = objectRaw["name"]
		mapping["alias"] = objectRaw["alias"]
		mapping["description"] = objectRaw["description"]
		mapping["type"] = objectRaw["type"]
		mapping["version"] = objectRaw["version"]
		mapping["status"] = objectRaw["publishStatus"]

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			names = append(names, objectRaw["name"])
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(objectRaw["pluginClassId"])
		mapping, err = dataSourcePluginClassesReadDescription(id, mapping, meta)
		if err != nil {
			return errs.WrapError(err)
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["name"])
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("classes", s); err != nil {
		return errs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := helper.WriteToFile(output.(string), s); err != nil {
			return errs.WrapError(err)
		}
	}
	return nil
}

func dataSourcePluginClassesReadDescription(id string, object map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	apigServiceV2 := NewService(client)
	objectRaw, err := apigServiceV2.DescribeApigPluginClass(id)
	if err != nil {
		return nil, errs.WrapError(err)
	}

	object["alias"] = objectRaw["alias"]
	object["description"] = objectRaw["description"]
	object["document"] = objectRaw["document"]
	object["plugin_class_name"] = objectRaw["name"]
	object["status"] = objectRaw["publishStatus"]
	object["type"] = objectRaw["type"]
	object["wasm_language"] = objectRaw["wasmLanguage"]

	return object, nil
}
