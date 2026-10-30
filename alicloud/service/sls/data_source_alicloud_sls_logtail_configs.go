// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package sls

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAliCloudSlsLogtailConfigs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudSlsLogtailConfigRead,
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
			"logstore_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"logtail_config_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"offset": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"configs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"logtail_config_name": {
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

func dataSourceAliCloudSlsLogtailConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []string
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
	action := fmt.Sprintf("/configs")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	query["configName"] = helper.StringPointer(d.Get("logtail_config_name").(string))
	hostMap["project"] = helper.StringPointer(d.Get("project_name").(string))
	if v, ok := d.GetOk("logstore_name"); ok {
		query["logstoreName"] = helper.StringPointer(v.(string))
	}

	if v, ok := d.GetOk("logtail_config_name"); ok {
		query["configName"] = helper.StringPointer(v.(string))
	}

	if v, ok := d.GetOkExists("offset"); ok {
		query["offset"] = helper.StringPointer(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOkExists("size"); ok {
		query["size"] = helper.StringPointer(strconv.Itoa(v.(int)))
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		response, err = client.Do("Sls", helper.RoaParam("GET", "2020-12-30", "ListConfig", action), query, nil, nil, hostMap, false)

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

	resp, _ := jsonpath.Get("$.configs", response)

	result, _ := resp.([]interface{})
	for _, v := range result {
		item := fmt.Sprint(v)
		if nameRegex != nil && !nameRegex.MatchString(item) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(*hostMap["project"], ":", item)]; !ok {
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

		mapping["id"] = fmt.Sprint(*hostMap["project"], ":", response["configName"], ":", objectRaw)

		mapping["logtail_config_name"] = objectRaw

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw)
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("configs", s); err != nil {
		return errs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}
	return nil
}
