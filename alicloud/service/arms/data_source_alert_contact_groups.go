package arms

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceAlertContactGroups() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceAlertContactGroupsRead,
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
			"alert_contact_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"contact_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"contact_name": {
				Type:     schema.TypeString,
				Optional: true,
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
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alert_contact_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alert_contact_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"contact_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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

func DataSourceAlertContactGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "SearchAlertContactGroup"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("alert_contact_group_name"); ok {
		request["ContactGroupName"] = v
	}
	if v, ok := d.GetOk("contact_id"); ok {
		request["ContactId"] = v
	}
	if v, ok := d.GetOk("contact_name"); ok {
		request["ContactName"] = v
	}
	request["IsDetail"] = true
	request["RegionId"] = client.RegionId
	var objects []map[string]interface{}
	var alertContactGroupNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
		}
		alertContactGroupNameRegex = r
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
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
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
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_arms_alert_contact_groups", action, errs.AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.ContactGroups", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.ContactGroups", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if alertContactGroupNameRegex != nil && !alertContactGroupNameRegex.MatchString(fmt.Sprint(item["ContactGroupName"])) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["ContactGroupId"])]; !ok {
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
			"id":                       fmt.Sprint(object["ContactGroupId"]),
			"alert_contact_group_id":   fmt.Sprint(object["ContactGroupId"]),
			"alert_contact_group_name": object["ContactGroupName"],
			"create_time":              fmt.Sprint(object["CreateTime"]),
		}
		contactIdsItems := make([]string, 0)
		if contacts, ok := object["Contacts"]; ok && contacts != nil {
			for _, contactsItem := range contacts.([]interface{}) {
				if contactId, ok := contactsItem.(map[string]interface{})["ContactId"]; ok && contactId != nil {
					contactIdsItems = append(contactIdsItems, fmt.Sprint(contactId))
				}
			}
		}
		mapping["contact_ids"] = contactIdsItems
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["ContactGroupName"])
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
