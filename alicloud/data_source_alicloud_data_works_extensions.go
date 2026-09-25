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

func dataSourceAlicloudDataWorksExtensions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDataWorksExtensionsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"extension_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"extensions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extension_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extension_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extension_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"help_doc_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_testing": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"detail_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter_setting": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"option_setting": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bind_event_list": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"event_category_list": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDataWorksExtensionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

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

	var objects []map[string]interface{}

	if v, ok := d.GetOk("extension_code"); ok && v.(string) != "" {
		action := "GetExtension"
		request := map[string]interface{}{
			"ExtensionCode": v.(string),
		}
		var response map[string]interface{}
		var err error
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_data_works_extensions", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Data", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data", response)
		}
		if obj, ok := resp.(map[string]interface{}); ok && len(obj) > 0 {
			objects = append(objects, obj)
		}
	} else {
		action := "ListExtensions"
		request := make(map[string]interface{})
		request["PageSize"] = PageSizeLarge
		request["PageNumber"] = 1
		var response map[string]interface{}
		var err error
		for {
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, true)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, request)
			if err != nil {
				return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_data_works_extensions", action, AlibabaCloudSdkGoERROR)
			}
			resp, err := jsonpath.Get("$.Data.Extensions", response)
			if err != nil {
				return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.Extensions", response)
			}
			result, _ := resp.([]interface{})
			for _, item := range result {
				obj, ok := item.(map[string]interface{})
				if !ok {
					continue
				}
				if len(idsMap) > 0 {
					if _, ok := idsMap[fmt.Sprint(obj["ExtensionCode"])]; !ok {
						continue
					}
				}
				if nameRegex != nil {
					if name := fmt.Sprint(obj["ExtensionName"]); !nameRegex.MatchString(name) {
						continue
					}
				}
				objects = append(objects, obj)
			}
			if len(result) < PageSizeLarge {
				break
			}
			request["PageNumber"] = request["PageNumber"].(int) + 1
		}
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                  fmt.Sprint(object["ExtensionCode"]),
			"extension_code":      fmt.Sprint(object["ExtensionCode"]),
			"extension_name":      fmt.Sprint(object["ExtensionName"]),
			"extension_desc":      fmt.Sprint(object["ExtensionDesc"]),
			"status":              fmt.Sprint(object["Status"]),
			"help_doc_url":        fmt.Sprint(object["HelpDocUrl"]),
			"project_testing":     fmt.Sprint(object["ProjectTesting"]),
			"detail_url":          fmt.Sprint(object["DetailUrl"]),
			"parameter_setting":   convertObjectToJsonString(object["ParameterSetting"]),
			"option_setting":      convertObjectToJsonString(object["OptionSetting"]),
			"bind_event_list":     convertObjectToJsonString(object["BindEventList"]),
			"event_category_list": convertObjectToJsonString(object["EventCategoryList"]),
			"region_id":           fmt.Sprint(object["RegionId"]),
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("extensions", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
