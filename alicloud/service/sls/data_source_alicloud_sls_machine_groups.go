// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package sls

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAliCloudSlsMachineGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudSlsMachineGroupRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_name": {
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
			},
		},
	}
}

func dataSourceAliCloudSlsMachineGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []string

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
	action := fmt.Sprintf("/machinegroups")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	query["groupName"] = helper.StringPointer(d.Get("group_name").(string))
	hostMap["project"] = helper.StringPointer(d.Get("project_name").(string))
	if v, ok := d.GetOk("group_name"); ok {
		query["groupName"] = helper.StringPointer(v.(string))
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	query["size"] = helper.StringPointer(strconv.Itoa(helper.PageSizeLarge))
	query["offset"] = helper.StringPointer("0")
	for {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.Do("Sls", helper.RoaParam("GET", "2020-12-30", "ListMachineGroup", action), query, nil, nil, hostMap, true)

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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.machinegroups", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(string)
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(*hostMap["project"], ":", item)]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < helper.PageSizeLarge {
			break
		}
		offsetFromQuery, err := convert.FormatInt(*query["offset"])
		if err != nil {
			return err
		}
		offset := offsetFromQuery + helper.PageSizeLarge
		query["offset"] = helper.StringPointer(strconv.Itoa(offset))
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = fmt.Sprint(*hostMap["project"], ":", objectRaw)

		mapping["group_name"] = objectRaw

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw)
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("groups", s); err != nil {
		return errs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}
	return nil
}
