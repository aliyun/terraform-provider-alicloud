package alicloud

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudCmsAlertActions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCmsAlertActionsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alert_actions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alert_action_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alert_action_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"webhook_param": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"method":       {Type: schema.TypeString, Computed: true},
									"url":          {Type: schema.TypeString, Computed: true},
									"content_type": {Type: schema.TypeString, Computed: true},
									"headers":      {Type: schema.TypeMap, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
								},
							},
						},
						"mns_param": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mns_type":  {Type: schema.TypeString, Computed: true},
									"name":      {Type: schema.TypeString, Computed: true},
									"region_id": {Type: schema.TypeString, Computed: true},
								},
							},
						},
						"sls_param": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"logstore":  {Type: schema.TypeString, Computed: true},
									"project":   {Type: schema.TypeString, Computed: true},
									"region_id": {Type: schema.TypeString, Computed: true},
								},
							},
						},
						"ess_param": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ess_group_id": {Type: schema.TypeString, Computed: true},
									"ess_rule_id":  {Type: schema.TypeString, Computed: true},
									"region_id":    {Type: schema.TypeString, Computed: true},
								},
							},
						},
						"fc_param": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"function":  {Type: schema.TypeString, Computed: true},
									"region_id": {Type: schema.TypeString, Computed: true},
									"service":   {Type: schema.TypeString, Computed: true},
								},
							},
						},
						"pager_duty_param": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {Type: schema.TypeString, Computed: true},
									"url": {Type: schema.TypeString, Computed: true},
								},
							},
						},
						"fc3_param": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region_id": {Type: schema.TypeString, Computed: true},
									"function":  {Type: schema.TypeString, Computed: true},
									"qualifier": {Type: schema.TypeString, Computed: true},
								},
							},
						},
						"eb_param": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region_id":      {Type: schema.TypeString, Computed: true},
									"event_bus_name": {Type: schema.TypeString, Computed: true},
									"subject":        {Type: schema.TypeString, Computed: true},
									"eb_source":      {Type: schema.TypeString, Computed: true},
								},
							},
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudCmsAlertActionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/alertActions"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	query["RegionId"] = StringPointer(client.RegionId)

	if v, ok := d.GetOk("type"); ok {
		query["type"] = StringPointer(v.(string))
	}

	var allAlertActions []map[string]interface{}
	pageNumber := 1
	pageSize := 100
	for {
		query["pageNumber"] = StringPointer(fmt.Sprintf("%d", pageNumber))
		query["pageSize"] = StringPointer(fmt.Sprintf("%d", pageSize))

		var ids []string
		if v, ok := d.GetOk("ids"); ok {
			idList := v.([]interface{})
			for _, id := range idList {
				ids = append(ids, id.(string))
			}
		}
		if len(ids) > 0 {
			idsJSON, _ := json.Marshal(ids)
			query["alertActionIds"] = StringPointer(string(idsJSON))
		}

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)
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
			return WrapErrorf(err, DefaultErrorMsg, "data.alicloud_cms_alert_actions", action, AlibabaCloudSdkGoERROR)
		}

		v, err := jsonpath.Get("$.alertActions", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, "data.alicloud_cms_alert_actions", "$.alertActions", response)
		}

		if v == nil {
			break
		}

		alertActionsRaw := v.([]interface{})
		for _, item := range alertActionsRaw {
			allAlertActions = append(allAlertActions, item.(map[string]interface{}))
		}

		total, _ := jsonpath.Get("$.total", response)
		var totalInt int
		if n, ok := total.(json.Number); ok {
			if i, err := n.Int64(); err == nil {
				totalInt = int(i)
			}
		} else if f, ok := total.(float64); ok {
			totalInt = int(f)
		}
		if pageNumber*pageSize >= totalInt {
			break
		}
		pageNumber++
	}

	var filteredActions []map[string]interface{}
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex, err = regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
	}

	for _, item := range allAlertActions {
		if nameRegex != nil {
			name, _ := item["alertActionName"].(string)
			if !nameRegex.MatchString(name) {
				continue
			}
		}
		filteredActions = append(filteredActions, item)
	}

	var ids []string
	for _, item := range filteredActions {
		ids = append(ids, fmt.Sprint(item["alertActionId"]))
	}
	sort.Strings(ids)
	d.SetId(dataResourceIdHash(ids))

	maps := make([]map[string]interface{}, 0)
	for _, item := range filteredActions {
		m := make(map[string]interface{})
		m["alert_action_id"] = item["alertActionId"]
		m["alert_action_name"] = item["alertActionName"]
		m["type"] = item["type"]
		m["region_id"] = item["regionId"]
		m["webhook_param"] = flattenCmsAlertActionWebhookParam(item["webhookParam"])
		m["mns_param"] = flattenCmsAlertActionMnsParam(item["mnsParam"])
		m["sls_param"] = flattenCmsAlertActionSlsParam(item["slsParam"])
		m["ess_param"] = flattenCmsAlertActionEssParam(item["essParam"])
		m["fc_param"] = flattenCmsAlertActionFcParam(item["fcParam"])
		m["pager_duty_param"] = flattenCmsAlertActionPagerDutyParam(item["pagerDutyParam"])
		m["fc3_param"] = flattenCmsAlertActionFc3Param(item["fc3Param"])
		m["eb_param"] = flattenCmsAlertActionEbParam(item["ebParam"])
		maps = append(maps, m)
	}

	if err := d.Set("alert_actions", maps); err != nil {
		return err
	}
	if err := d.Set("ids", ids); err != nil {
		return err
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), maps); err != nil {
			return err
		}
	}

	return nil
}

func flattenCmsAlertActionWebhookParam(v interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if v == nil {
		return result
	}
	m, ok := v.(map[string]interface{})
	if !ok || len(m) == 0 {
		return result
	}
	item := make(map[string]interface{})
	item["method"] = m["method"]
	item["url"] = m["url"]
	item["content_type"] = m["contentType"]
	item["headers"] = m["headers"]
	result = append(result, item)
	return result
}

func flattenCmsAlertActionMnsParam(v interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if v == nil {
		return result
	}
	m, ok := v.(map[string]interface{})
	if !ok || len(m) == 0 {
		return result
	}
	item := make(map[string]interface{})
	item["mns_type"] = m["mnsType"]
	item["name"] = m["name"]
	item["region_id"] = m["regionId"]
	result = append(result, item)
	return result
}

func flattenCmsAlertActionSlsParam(v interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if v == nil {
		return result
	}
	m, ok := v.(map[string]interface{})
	if !ok || len(m) == 0 {
		return result
	}
	item := make(map[string]interface{})
	item["logstore"] = m["logstore"]
	item["project"] = m["project"]
	item["region_id"] = m["regionId"]
	result = append(result, item)
	return result
}

func flattenCmsAlertActionEssParam(v interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if v == nil {
		return result
	}
	m, ok := v.(map[string]interface{})
	if !ok || len(m) == 0 {
		return result
	}
	item := make(map[string]interface{})
	item["ess_group_id"] = m["essGroupId"]
	item["ess_rule_id"] = m["essRuleId"]
	item["region_id"] = m["regionId"]
	result = append(result, item)
	return result
}

func flattenCmsAlertActionFcParam(v interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if v == nil {
		return result
	}
	m, ok := v.(map[string]interface{})
	if !ok || len(m) == 0 {
		return result
	}
	item := make(map[string]interface{})
	item["function"] = m["function"]
	item["region_id"] = m["regionId"]
	item["service"] = m["service"]
	result = append(result, item)
	return result
}

func flattenCmsAlertActionPagerDutyParam(v interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if v == nil {
		return result
	}
	m, ok := v.(map[string]interface{})
	if !ok || len(m) == 0 {
		return result
	}
	item := make(map[string]interface{})
	item["key"] = m["key"]
	item["url"] = m["url"]
	result = append(result, item)
	return result
}

func flattenCmsAlertActionFc3Param(v interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if v == nil {
		return result
	}
	m, ok := v.(map[string]interface{})
	if !ok || len(m) == 0 {
		return result
	}
	item := make(map[string]interface{})
	item["region_id"] = m["regionId"]
	item["function"] = m["function"]
	item["qualifier"] = m["qualifier"]
	result = append(result, item)
	return result
}

func flattenCmsAlertActionEbParam(v interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if v == nil {
		return result
	}
	m, ok := v.(map[string]interface{})
	if !ok || len(m) == 0 {
		return result
	}
	item := make(map[string]interface{})
	item["region_id"] = m["regionId"]
	item["event_bus_name"] = m["eventBusName"]
	item["subject"] = m["subject"]
	item["eb_source"] = m["ebSource"]
	result = append(result, item)
	return result
}
