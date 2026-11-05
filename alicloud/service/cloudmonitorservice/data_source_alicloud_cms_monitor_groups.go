package cloudmonitorservice

import (
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlicloudCmsMonitorGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCmsMonitorGroupsRead,
		Schema: map[string]*schema.Schema{
			"dynamic_tag_rule_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"include_template_history": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"keyword": {
				Type:     schema.TypeString,
				Optional: true,
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
			"monitor_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"select_contact_groups": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": tags.TagsSchema(),
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"custom", "ehpc_cluster", "kubernetes"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bind_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"contact_groups": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"dynamic_tag_rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gmt_create": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"gmt_modified": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"monitor_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"template_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCmsMonitorGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeMonitorGroups"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("dynamic_tag_rule_id"); ok {
		request["DynamicTagRuleId"] = v
	}
	if v, ok := d.GetOkExists("include_template_history"); ok {
		request["IncludeTemplateHistory"] = v
	}
	if v, ok := d.GetOk("keyword"); ok {
		request["Keyword"] = v
	}
	if v, ok := d.GetOk("monitor_group_name"); ok {
		request["GroupName"] = v
	}
	if v, ok := d.GetOkExists("select_contact_groups"); ok {
		request["SelectContactGroups"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tags := make([]map[string]interface{}, 0)
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, map[string]interface{}{
				"Key":   key,
				"Value": value.(string),
			})
		}
		request["Tag"] = tags
	}
	if v, ok := d.GetOk("type"); ok {
		request["Type"] = v
	}
	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var monitorGroupNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
		}
		monitorGroupNameRegex = r
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
	for {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cms_monitor_groups", action, errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(action, response, request)

		resp, err := jsonpath.Get("$.Resources.Resource", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Resources.Resource", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if monitorGroupNameRegex != nil {
				if !monitorGroupNameRegex.MatchString(fmt.Sprint(item["GroupName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["GroupId"])]; !ok {
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
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {

		gmt_createVal, err := convert.FormatInt(object["GmtCreate"])
		if err != nil {
			return err
		}
		gmt_modifiedVal, err := convert.FormatInt(object["GmtModified"])
		if err != nil {
			return err
		}
		resourceGroupId, err := convert.FormatInt(object["GroupId"])
		if err != nil {
			return err
		}
		groupId, err := convert.FormatInt(object["GroupId"])
		if err != nil {
			return err
		}
		mapping := map[string]interface{}{
			"bind_url":            object["BindUrl"],
			"dynamic_tag_rule_id": object["DynamicTagRuleId"],
			"gmt_create":          gmt_createVal,
			"gmt_modified":        gmt_modifiedVal,
			"id":                  fmt.Sprint(resourceGroupId),
			"group_id":            fmt.Sprint(groupId),
			"monitor_group_name":  object["GroupName"],
			"service_id":          object["ServiceId"],
			"template_ids":        object["TemplateIds"].(map[string]interface{})["TemplateId"],
			"type":                object["Type"],
		}
		if v := object["ContactGroups"].(map[string]interface{})["ContactGroup"].([]interface{}); len(v) > 0 {
			var contactGroups []interface{}
			if v, ok := object["ContactGroups"].(map[string]interface{})["ContactGroup"]; ok {
				for _, contactGroup := range v.([]interface{}) {
					contactGroups = append(contactGroups, contactGroup.(map[string]interface{})["Name"])
				}
			}
			mapping["contact_groups"] = contactGroups
		}

		tagsMap := make(map[string]interface{})
		t, _ := jsonpath.Get("$.Tags.Tag", object)
		if t != nil {
			for _, t := range t.([]interface{}) {
				key := t.(map[string]interface{})["Key"].(string)
				value := t.(map[string]interface{})["Value"].(string)
				if !tags.IgnoredTags(key, value) {
					tagsMap[key] = value
				}
			}
		}
		mapping["tags"] = tagsMap
		ids = append(ids, fmt.Sprint(object["GroupId"]))
		names = append(names, object["GroupName"])
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
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
