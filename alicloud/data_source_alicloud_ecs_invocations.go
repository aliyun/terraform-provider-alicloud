package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudEcsInvocations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEcsInvocationsRead,
		Schema: map[string]*schema.Schema{
			"command_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"content_encoding": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PlainText", "Base64"}, false),
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"invoke_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Running", "Finished", "Failed", "PartialFailed", "Stopped"}, false),
			},
			"output_file": {
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
				Default:  PageSizeLarge,
			},
			"invocations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"command_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"command_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"command_content": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"command_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"frequency": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"invoke_instances": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"creation_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"update_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"finish_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"invocation_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"repeats": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"output": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"dropped": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"stop_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"exit_code": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"start_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"error_info": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"timed": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"error_code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_invoke_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"invocation_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"invoke_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"invocation_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameters": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repeat_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timed": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEcsInvocationsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeInvocations"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("command_id"); ok {
		request["CommandId"] = v
	}
	if v, ok := d.GetOk("content_encoding"); ok {
		request["ContentEncoding"] = v
	}
	if v, ok := d.GetOk("invoke_status"); ok {
		request["InvokeStatus"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("page_number"); ok && v.(int) > 0 {
		request["PageNumber"] = v.(int)
	} else {
		request["PageNumber"] = 1
	}
	if v, ok := d.GetOk("page_size"); ok && v.(int) > 0 {
		request["PageSize"] = v.(int)
	} else {
		request["PageSize"] = PageSizeLarge
	}
	var objects []map[string]interface{}

	// The command output is returned by DescribeInvocations only when
	// IncludeOutput is set to true and the InvokeId or InstanceId parameter
	// is specified. Therefore, when ids are set, query each invocation by
	// its InvokeId with IncludeOutput enabled so that the output of every
	// invoke instance can be returned. Otherwise, list all invocations by page.
	queries := make([]map[string]interface{}, 0)
	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
		seen := make(map[string]bool)
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			invokeId := vv.(string)
			if seen[invokeId] {
				continue
			}
			seen[invokeId] = true
			query := make(map[string]interface{}, len(request)+2)
			for key, value := range request {
				query[key] = value
			}
			query["InvokeId"] = invokeId
			query["IncludeOutput"] = "true"
			query["PageNumber"] = 1
			queries = append(queries, query)
		}
	} else {
		queries = append(queries, request)
	}
	var response map[string]interface{}
	var err error
	for _, query := range queries {
		for {
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, query, true)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, query)

			if err != nil {
				return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecs_invocations", action, AlibabaCloudSdkGoERROR)
			}

			resp, err := jsonpath.Get("$.Invocations.Invocation", response)
			if err != nil {
				return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Invocations.Invocation", response)
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				objects = append(objects, v.(map[string]interface{}))
			}
			if len(result) < PageSizeLarge {
				break
			}
			query["PageNumber"] = query["PageNumber"].(int) + 1
		}
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"command_id":        object["CommandId"],
			"frequency":         object["Frequency"],
			"repeat_mode":       object["RepeatMode"],
			"timed":             object["Timed"],
			"invocation_id":     object["InvokeId"],
			"id":                object["InvokeId"],
			"username":          object["Username"],
			"parameters":        object["Parameters"],
			"command_type":      object["CommandType"],
			"command_name":      object["CommandName"],
			"invocation_status": object["InvocationStatus"],
			"create_time":       object["CreationTime"],
			"invoke_status":     object["InvokeStatus"],
			"command_content":   object["CommandContent"],
		}
		instanceIdItems := make([]map[string]interface{}, 0)
		if invokeInstances, ok := object["InvokeInstances"]; ok && invokeInstances != nil {
			if invokeInstance, ok := invokeInstances.(map[string]interface{})["InvokeInstance"]; ok && invokeInstance != nil {
				for _, invokeInstanceItem := range invokeInstance.([]interface{}) {
					item := invokeInstanceItem.(map[string]interface{})
					instanceIdItems = append(instanceIdItems, map[string]interface{}{
						"creation_time":          item["CreationTime"],
						"update_time":            item["UpdateTime"],
						"finish_time":            item["FinishTime"],
						"invocation_status":      item["InvocationStatus"],
						"repeats":                formatInt(item["Repeats"]),
						"instance_id":            item["InstanceId"],
						"output":                 item["Output"],
						"dropped":                formatInt(item["Dropped"]),
						"stop_time":              item["StopTime"],
						"exit_code":              formatInt(item["ExitCode"]),
						"start_time":             item["StartTime"],
						"error_info":             item["ErrorInfo"],
						"timed":                  item["Timed"],
						"error_code":             item["ErrorCode"],
						"instance_invoke_status": item["InstanceInvokeStatus"],
					})
				}
				mapping["invoke_instances"] = instanceIdItems
			}
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("invocations", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
