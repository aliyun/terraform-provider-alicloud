// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudRealtimeComputeMembers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudRealtimeComputeMemberRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"members": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"member": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudRealtimeComputeMemberRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	namespace := d.Get("namespace")
	workspace := d.Get("resource_id")
	action := fmt.Sprintf("/gateway/v2/namespaces/%s/members", namespace)
	request := make(map[string]*string)
	header := make(map[string]*string)
	request["pageSize"] = StringPointer(strconv.Itoa(PageSizeLarge))
	request["pageIndex"] = StringPointer("1")
	header["workspace"] = StringPointer(workspace.(string))

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
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RoaGet("ververica", "2022-07-18", action, request, header, nil)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_realtime_compute_members", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.data", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.data", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprintf("%v:%v:%v", workspace, namespace, item["member"])]; !ok {
					continue
				}
			}

			objects = append(objects, item)
		}

		if len(result) < PageSizeLarge {
			break
		}

		pageIndex, err := strconv.Atoi(*request["pageIndex"])
		if err != nil {
			return WrapError(err)
		}

		request["pageIndex"] = StringPointer(strconv.Itoa(pageIndex + 1))
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":          fmt.Sprintf("%v:%v:%v", workspace, namespace, object["member"]),
			"resource_id": fmt.Sprint(workspace),
			"namespace":   fmt.Sprint(namespace),
			"member":      fmt.Sprint(object["member"]),
			"role":        object["role"],
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("members", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
