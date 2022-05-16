package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudEcsInvocations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEcsInvocationsRead,
		Schema: map[string]*schema.Schema{
			"command_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"content_encoding": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PlainText", "Base64"}, false),
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"invoke_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
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

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	var response map[string]interface{}
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecs_invocations", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Invocations.Invocation", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Invocations.Invocation", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["InvokeId"])]; !ok {
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
