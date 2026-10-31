package ehpc

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlicloudEhpcJobTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEhpcJobTemplatesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"array_request": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"clock_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"command_line": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"job_template_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"job_template_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mem": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"package_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"queue": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"re_runable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"runas_user": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stderr_redirect_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stdout_redirect_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"task": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"thread": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"variables": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEhpcJobTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListJobTemplates"
	request := make(map[string]interface{})

	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1
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
	var err error
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcGet("EHPC", "2018-04-12", action, request, nil)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_ehpc_job_templates", action, errs.AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Templates.JobTemplates", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Templates.JobTemplates", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["Id"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {

		gpuVal, err := convert.FormatInt(object["Gpu"])
		if err != nil {
			return err
		}

		nodeVal, err := convert.FormatInt(object["Node"])
		if err != nil {
			return err
		}

		priorityVal, err := convert.FormatInt(object["Priority"])
		if err != nil {
			return err
		}

		taskVal, err := convert.FormatInt(object["Task"])
		if err != nil {
			return err
		}
		threadVal, err := convert.FormatInt(object["Thread"])
		if err != nil {
			return err
		}
		mapping := map[string]interface{}{
			"array_request":        object["ArrayRequest"],
			"clock_time":           object["ClockTime"],
			"command_line":         object["CommandLine"],
			"gpu":                  gpuVal,
			"id":                   fmt.Sprint(object["Id"]),
			"job_template_id":      fmt.Sprint(object["Id"]),
			"job_template_name":    object["Name"],
			"mem":                  object["Mem"],
			"node":                 nodeVal,
			"package_path":         object["PackagePath"],
			"priority":             priorityVal,
			"queue":                object["Queue"],
			"re_runable":           object["ReRunable"],
			"runas_user":           object["RunasUser"],
			"stderr_redirect_path": object["StderrRedirectPath"],
			"stdout_redirect_path": object["StdoutRedirectPath"],
			"task":                 taskVal,
			"thread":               threadVal,
			"variables":            object["Variables"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("templates", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
