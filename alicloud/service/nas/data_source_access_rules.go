package nas

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlicloudAccessRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAccessRulesRead,

		Schema: map[string]*schema.Schema{
			"source_cidr_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_access": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rw_access": {
				Type:     schema.TypeString,
				Optional: true,
			},
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
			// Computed values
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_cidr_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"access_rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_access": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rw_access": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func dataSourceAlicloudAccessRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeAccessRules"
	request := make(map[string]interface{})
	request["AccessGroupName"] = d.Get("access_group_name")
	request["RegionId"] = client.Region
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
		response, err = client.RpcPost("NAS", "2017-06-26", action, nil, request, true)
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_nas_access_rules", action, errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(action, response, request)

		resp, err := jsonpath.Get("$.AccessRules.AccessRule", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.AccessRules.AccessRule", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if v, ok := d.GetOk("source_cidr_ip"); ok && v.(string) != "" && item["SourceCidrIp"].(string) != v.(string) {
				continue
			}
			if v, ok := d.GetOk("user_access"); ok && v.(string) != "" && item["UserAccess"].(string) != v.(string) {
				continue
			}
			if v, ok := d.GetOk("rw_access"); ok && v.(string) != "" && item["RWAccess"].(string) != v.(string) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["AccessRuleId"])]; !ok {
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

		priorityVal, err := convert.FormatInt(object["Priority"])
		if err != nil {
			return err
		}
		mapping := map[string]interface{}{
			"source_cidr_ip": object["SourceCidrIp"],
			"priority":       priorityVal,
			"access_rule_id": object["AccessRuleId"],
			"user_access":    object["UserAccess"],
			"rw_access":      object["RWAccess"],
		}
		ids = append(ids, fmt.Sprint(object["AccessRuleId"]))
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("rules", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}
	return nil
}
