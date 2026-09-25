package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCmsEntityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsEntityGroupCreate,
		Read:   resourceAliCloudCmsEntityGroupRead,
		Update: resourceAliCloudCmsEntityGroupUpdate,
		Delete: resourceAliCloudCmsEntityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"entity_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entity_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"entity_rules": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"tags": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"op": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"tag_key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"tag_values": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"labels": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"op": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"tag_key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"tag_values": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"ip_match_rule": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_field_key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ip_cidr": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"instance_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"field_rules": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field_key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"op": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"field_values": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"entity_types": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"annotations": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"op": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"tag_key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"tag_values": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudCmsEntityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/entity-groups"
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})

	body["entityGroupName"] = d.Get("entity_group_name")
	body["description"] = d.Get("description")
	body["workspace"] = d.Get("workspace")

	if v, ok := d.GetOk("resource_group_id"); ok && v.(string) != "" {
		entityRules := expandCmsEntityRules(d.Get("entity_rules"))
		if entityRules == nil {
			entityRules = make(map[string]interface{})
		}
		entityRules["resourceGroupId"] = v.(string)
		body["entityRules"] = entityRules
	} else if er := expandCmsEntityRules(d.Get("entity_rules")); er != nil {
		body["entityRules"] = er
	}

	wait := incrementalWait(3*time.Second, 0*time.Second)
	var err error
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("Cms", "2024-03-30", action, query, nil, body, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, body)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_entity_group", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.entityGroupId", response)
	if id == nil || fmt.Sprint(id) == "" {
		return WrapError(Error("CreateEntityGroup failed: empty entityGroupId in response"))
	}
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudCmsEntityGroupRead(d, meta)
}

func resourceAliCloudCmsEntityGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	object, err := cmsServiceV2.DescribeCmsEntityGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_entity_group DescribeCmsEntityGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	entityGroupRaw, _ := jsonpath.Get("$.entityGroup", object)
	entityGroup := make(map[string]interface{})
	if entityGroupRaw != nil {
		if eg, ok := entityGroupRaw.(map[string]interface{}); ok {
			entityGroup = eg
		}
	}

	d.Set("entity_group_id", entityGroup["entityGroupId"])
	d.Set("entity_group_name", entityGroup["entityGroupName"])
	d.Set("description", entityGroup["description"])
	d.Set("workspace", entityGroup["workspace"])

	entityRulesMaps := flattenCmsEntityRules(entityGroup["entityRules"])
	if err := d.Set("entity_rules", entityRulesMaps); err != nil {
		return WrapError(err)
	}

	if entityRules, ok := entityGroup["entityRules"].(map[string]interface{}); ok {
		if rgId, ok := entityRules["resourceGroupId"].(string); ok && rgId != "" {
			d.Set("resource_group_id", rgId)
		}
	}

	return nil
}

func resourceAliCloudCmsEntityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	entityGroupId := d.Id()
	action := fmt.Sprintf("/entity-groups/%s", entityGroupId)
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})

	body["entityGroupName"] = d.Get("entity_group_name")
	body["description"] = d.Get("description")
	body["workspace"] = d.Get("workspace")

	if v, ok := d.GetOk("resource_group_id"); ok && v.(string) != "" {
		entityRules := expandCmsEntityRules(d.Get("entity_rules"))
		if entityRules == nil {
			entityRules = make(map[string]interface{})
		}
		entityRules["resourceGroupId"] = v.(string)
		body["entityRules"] = entityRules
	} else if er := expandCmsEntityRules(d.Get("entity_rules")); er != nil {
		body["entityRules"] = er
	}

	wait := incrementalWait(3*time.Second, 0*time.Second)
	var err error
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.RoaPut("Cms", "2024-03-30", action, query, nil, body, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, body)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return resourceAliCloudCmsEntityGroupRead(d, meta)
}

func resourceAliCloudCmsEntityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	entityGroupId := d.Id()
	action := fmt.Sprintf("/entity-groups/%s", entityGroupId)
	var response map[string]interface{}
	query := make(map[string]*string)

	if v, ok := d.GetOk("workspace"); ok && v.(string) != "" {
		query["workspace"] = StringPointer(v.(string))
	}

	wait := incrementalWait(3*time.Second, 0*time.Second)
	var err error
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("Cms", "2024-03-30", action, query, nil, nil, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, query)
	if err != nil {
		if IsExpectedErrors(err, []string{"404", "500"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

// expandCmsEntityRules converts the schema entity_rules list to the API entityRules map.
func expandCmsEntityRules(v interface{}) map[string]interface{} {
	rules, ok := v.([]interface{})
	if !ok || len(rules) == 0 {
		return nil
	}
	rule, ok := rules[0].(map[string]interface{})
	if !ok {
		return nil
	}
	entityRules := make(map[string]interface{})

	if rgId, ok := rule["resource_group_id"].(string); ok && rgId != "" {
		entityRules["resourceGroupId"] = rgId
	}

	if tags := expandCmsEntityTagRules(rule["tags"]); len(tags) > 0 {
		entityRules["tags"] = tags
	}
	if labels := expandCmsEntityTagRules(rule["labels"]); len(labels) > 0 {
		entityRules["labels"] = labels
	}
	if annotations := expandCmsEntityTagRules(rule["annotations"]); len(annotations) > 0 {
		entityRules["annotations"] = annotations
	}

	if ipr, ok := rule["ip_match_rule"].([]interface{}); ok && len(ipr) > 0 {
		if iprm, ok := ipr[0].(map[string]interface{}); ok {
			ipMatchRule := make(map[string]interface{})
			if v := iprm["ip_field_key"].(string); v != "" {
				ipMatchRule["ipFieldKey"] = v
			}
			if v := iprm["ip_cidr"].(string); v != "" {
				ipMatchRule["ipCIDR"] = v
			}
			if len(ipMatchRule) > 0 {
				entityRules["ipMatchRule"] = ipMatchRule
			}
		}
	}

	if iids, ok := rule["instance_ids"].([]interface{}); ok && len(iids) > 0 {
		ids := make([]string, 0)
		for _, id := range iids {
			ids = append(ids, id.(string))
		}
		entityRules["instanceIds"] = ids
	}

	if entityTypes, ok := rule["entity_types"].([]interface{}); ok && len(entityTypes) > 0 {
		et := make([]string, 0)
		for _, t := range entityTypes {
			et = append(et, t.(string))
		}
		entityRules["entityTypes"] = et
	}

	if fieldRules := expandCmsEntityFieldRules(rule["field_rules"]); len(fieldRules) > 0 {
		entityRules["fieldRules"] = fieldRules
	}

	return entityRules
}

// expandCmsEntityTagRules converts a list of tag/label/annotation schema blocks to the API list.
func expandCmsEntityTagRules(v interface{}) []map[string]interface{} {
	items, ok := v.([]interface{})
	if !ok || len(items) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, 0)
	for _, raw := range items {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		m := make(map[string]interface{})
		if v := item["op"].(string); v != "" {
			m["op"] = v
		}
		if v := item["tag_key"].(string); v != "" {
			m["tagKey"] = v
		}
		if tvs, ok := item["tag_values"].([]interface{}); ok && len(tvs) > 0 {
			vals := make([]string, 0)
			for _, tv := range tvs {
				vals = append(vals, tv.(string))
			}
			m["tagValues"] = vals
		}
		result = append(result, m)
	}
	return result
}

// expandCmsEntityFieldRules converts a list of field_rules schema blocks to the API list.
func expandCmsEntityFieldRules(v interface{}) []map[string]interface{} {
	items, ok := v.([]interface{})
	if !ok || len(items) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, 0)
	for _, raw := range items {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		m := make(map[string]interface{})
		if v := item["field_key"].(string); v != "" {
			m["fieldKey"] = v
		}
		if v := item["op"].(string); v != "" {
			m["op"] = v
		}
		if fvs, ok := item["field_values"].([]interface{}); ok && len(fvs) > 0 {
			vals := make([]string, 0)
			for _, fv := range fvs {
				vals = append(vals, fv.(string))
			}
			m["fieldValues"] = vals
		}
		result = append(result, m)
	}
	return result
}

// flattenCmsEntityRules converts the API entityRules map to the schema entity_rules list.
func flattenCmsEntityRules(v interface{}) []map[string]interface{} {
	entityRules, ok := v.(map[string]interface{})
	if !ok || len(entityRules) == 0 {
		return nil
	}
	rule := make(map[string]interface{})

	if rgId, ok := entityRules["resourceGroupId"].(string); ok {
		rule["resource_group_id"] = rgId
	}

	if tagsRaw, ok := entityRules["tags"].([]interface{}); ok {
		rule["tags"] = flattenCmsEntityTagBlocks(tagsRaw)
	}
	if labelsRaw, ok := entityRules["labels"].([]interface{}); ok {
		rule["labels"] = flattenCmsEntityTagBlocks(labelsRaw)
	}
	if annotationsRaw, ok := entityRules["annotations"].([]interface{}); ok {
		rule["annotations"] = flattenCmsEntityTagBlocks(annotationsRaw)
	}

	if ipMatchRuleRaw, ok := entityRules["ipMatchRule"].(map[string]interface{}); ok {
		ipMatchRule := make(map[string]interface{})
		if v, ok := ipMatchRuleRaw["ipFieldKey"].(string); ok {
			ipMatchRule["ip_field_key"] = v
		}
		if v, ok := ipMatchRuleRaw["ipCIDR"].(string); ok {
			ipMatchRule["ip_cidr"] = v
		}
		rule["ip_match_rule"] = []map[string]interface{}{ipMatchRule}
	}

	if instanceIdsRaw, ok := entityRules["instanceIds"].([]interface{}); ok {
		ids := make([]string, 0)
		for _, id := range instanceIdsRaw {
			ids = append(ids, fmt.Sprint(id))
		}
		rule["instance_ids"] = ids
	}

	if entityTypesRaw, ok := entityRules["entityTypes"].([]interface{}); ok {
		et := make([]string, 0)
		for _, t := range entityTypesRaw {
			et = append(et, fmt.Sprint(t))
		}
		rule["entity_types"] = et
	}

	if fieldRulesRaw, ok := entityRules["fieldRules"].([]interface{}); ok {
		rule["field_rules"] = flattenCmsEntityFieldBlocks(fieldRulesRaw)
	}

	return []map[string]interface{}{rule}
}

func flattenCmsEntityTagBlocks(items []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	for _, raw := range items {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		m := make(map[string]interface{})
		if v, ok := item["op"].(string); ok {
			m["op"] = v
		}
		if v, ok := item["tagKey"].(string); ok {
			m["tag_key"] = v
		}
		if tvs, ok := item["tagValues"].([]interface{}); ok {
			vals := make([]string, 0)
			for _, tv := range tvs {
				vals = append(vals, fmt.Sprint(tv))
			}
			m["tag_values"] = vals
		}
		result = append(result, m)
	}
	return result
}

func flattenCmsEntityFieldBlocks(items []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	for _, raw := range items {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		m := make(map[string]interface{})
		if v, ok := item["fieldKey"].(string); ok {
			m["field_key"] = v
		}
		if v, ok := item["op"].(string); ok {
			m["op"] = v
		}
		if fvs, ok := item["fieldValues"].([]interface{}); ok {
			vals := make([]string, 0)
			for _, fv := range fvs {
				vals = append(vals, fmt.Sprint(fv))
			}
			m["field_values"] = vals
		}
		result = append(result, m)
	}
	return result
}
