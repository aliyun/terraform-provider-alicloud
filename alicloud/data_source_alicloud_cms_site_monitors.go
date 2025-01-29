// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudCloudMonitorServiceSiteMonitors() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCloudMonitorServiceSiteMonitorRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"task_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"task_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"monitors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"interval": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"task_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"task_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"task_type": {
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

func dataSourceAliCloudCloudMonitorServiceSiteMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	var query map[string]interface{}
	action := "DescribeSiteMonitorList"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	if v, ok := d.GetOk("task_id"); ok {
		request["TaskId"] = v
	}
	request["TaskType"] = d.Get("task_type")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, query, request, true)

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

	resp, _ := jsonpath.Get("$.SiteMonitors.SiteMonitor[*]", response)

	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["Name"])) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["TaskId"])]; !ok {
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

		mapping["address"] = objectRaw["Address"]
		mapping["create_time"] = objectRaw["CreateTime"]
		mapping["interval"] = objectRaw["Interval"]
		mapping["task_name"] = objectRaw["TaskName"]
		mapping["task_type"] = objectRaw["TaskType"]
		mapping["task_id"] = objectRaw["TaskId"]

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["AlertName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("monitors", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
