// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudSlsLogtailConfigs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudSlsLogtailConfigRead,
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
			"logstore_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"logtail_config_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"offset": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
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
				ForceNew: true,
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
	action := fmt.Sprintf("/configs")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	query["configName"] = StringPointer(d.Get("logtail_config_name").(string))
	hostMap["project"] = StringPointer(d.Get("project_name").(string))
	if v, ok := d.GetOk("logstore_name"); ok {
		query["logstoreName"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("logtail_config_name"); ok {
		query["configName"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOkExists("offset"); ok {
		query["offset"] = StringPointer(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOkExists("size"); ok {
		query["size"] = StringPointer(strconv.Itoa(v.(int)))
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.Do("Sls", roaParam("GET", "2020-12-30", "ListConfig", action), query, nil, nil, hostMap, false)

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

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("configs", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
