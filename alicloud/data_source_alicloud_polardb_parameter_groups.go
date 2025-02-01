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

func dataSourceAlicloudPolarDBParameterGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudPolarDBParameterGroupsRead,
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
			"db_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"MySQL"}, false),
			},
			"db_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"5.6", "5.7", "8.0"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter_group_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter_group_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter_counts": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"force_restart": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudPolarDBParameterGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeParameterGroups"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("db_type"); ok {
		request["DBType"] = v
	}
	if v, ok := d.GetOk("db_version"); ok {
		request["DBVersion"] = v
	}
	var objects []map[string]interface{}
	var parameterGroupNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		parameterGroupNameRegex = r
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
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, true)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_polardb_parameter_groups", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.ParameterGroups", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ParameterGroups", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if parameterGroupNameRegex != nil && !parameterGroupNameRegex.MatchString(fmt.Sprint(item["ParameterGroupName"])) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["ParameterGroupId"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                   fmt.Sprint(object["ParameterGroupId"]),
			"parameter_group_id":   fmt.Sprint(object["ParameterGroupId"]),
			"parameter_group_name": object["ParameterGroupName"],
			"db_type":              object["DBType"],
			"db_version":           object["DBVersion"],
			"parameter_group_desc": object["ParameterGroupDesc"],
			"parameter_group_type": object["ParameterGroupType"],
			"parameter_counts":     object["ParameterCounts"],
			"force_restart":        object["ForceRestart"],
			"create_time":          object["CreateTime"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["ParameterGroupName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("groups", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
