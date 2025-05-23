// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudFcv3Triggers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudFcv3TriggerRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"function_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"triggers": {
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
						"http_trigger": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url_intranet": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"url_internet": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"invocation_role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_modified_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"qualifier": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_config": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_type": {
							Type:     schema.TypeString,
							Computed: true,
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

func dataSourceAliCloudFcv3TriggerRead(d *schema.ResourceData, meta interface{}) error {
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

	var functionName string
	if v, ok := d.GetOk("function_name"); ok {
		functionName = v.(string)
	}

	action := fmt.Sprintf("/2023-03-30/functions/%s/triggers", functionName)
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	request["functionName"] = d.Get("function_name")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	request["MaxResults"] = PageSizeLarge
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaGet("FC", "2023-03-30", action, query, nil, nil)

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

		resp, _ := jsonpath.Get("$.triggers[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["triggerName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(functionName, ":", item["triggerName"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = fmt.Sprint(functionName, ":", objectRaw["triggerName"])

		mapping["create_time"] = objectRaw["createdTime"]
		mapping["description"] = objectRaw["description"]
		mapping["invocation_role"] = objectRaw["invocationRole"]
		mapping["last_modified_time"] = objectRaw["lastModifiedTime"]
		mapping["qualifier"] = objectRaw["qualifier"]
		mapping["source_arn"] = objectRaw["sourceArn"]
		mapping["status"] = objectRaw["status"]
		mapping["target_arn"] = objectRaw["targetArn"]
		mapping["trigger_config"] = objectRaw["triggerConfig"]
		mapping["trigger_id"] = objectRaw["triggerId"]
		mapping["trigger_type"] = objectRaw["triggerType"]
		mapping["trigger_name"] = objectRaw["triggerName"]

		httpTriggerMaps := make([]map[string]interface{}, 0)
		httpTriggerMap := make(map[string]interface{})
		httpTriggerRaw := make(map[string]interface{})
		if objectRaw["httpTrigger"] != nil {
			httpTriggerRaw = objectRaw["httpTrigger"].(map[string]interface{})
		}
		if len(httpTriggerRaw) > 0 {
			httpTriggerMap["url_internet"] = httpTriggerRaw["urlInternet"]
			httpTriggerMap["url_intranet"] = httpTriggerRaw["urlIntranet"]

			httpTriggerMaps = append(httpTriggerMaps, httpTriggerMap)
		}
		mapping["http_trigger"] = httpTriggerMaps

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["TriggerName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("triggers", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
