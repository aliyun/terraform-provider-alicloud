package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCmsDigitalEmployee() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsDigitalEmployeeCreate,
		Read:   resourceAliCloudCmsDigitalEmployeeRead,
		Update: resourceAliCloudCmsDigitalEmployeeUpdate,
		Delete: resourceAliCloudCmsDigitalEmployeeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"digital_employee_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_arn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default_rule": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attributes": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"knowledges": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bailian": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"workspace_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"index_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"region": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"attributes": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"sandbox_network_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_fqdns": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"allow_cidrs": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"enable_acl": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"tool_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"aliyun": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"deny_policy": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"auto_pass_policy": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"statements": {
										Type:       schema.TypeList,
										Optional:   true,
										Deprecated: "`statements` is deprecated and is only read back for compatibility; use `deny_policy` and `auto_pass_policy` instead.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"decision": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"product": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"api_version": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"actions": {
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
				},
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"employee_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCmsDigitalEmployeeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/digital-employee"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	request["name"] = d.Get("digital_employee_name")
	request["roleArn"] = d.Get("role_arn")
	if v, ok := d.GetOk("default_rule"); ok {
		request["defaultRule"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOk("display_name"); ok {
		request["displayName"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["resourceGroupId"] = v
	}
	if v, ok := d.GetOk("attributes"); ok {
		request["attributes"] = v
	}
	if v, ok := d.GetOk("knowledges"); ok {
		request["knowledges"] = expandCmsDigitalEmployeeKnowledges(v.([]interface{}))
	}
	if v, ok := d.GetOk("sandbox_network_policy"); ok {
		request["sandboxNetworkPolicy"] = expandCmsDigitalEmployeeSandboxNetworkPolicy(v.([]interface{}))
	}
	if v, ok := d.GetOk("tool_policy"); ok {
		request["toolPolicy"] = expandCmsDigitalEmployeeToolPolicy(v.([]interface{}))
	}
	if v, ok := d.GetOk("tags"); ok {
		request["tags"] = expandCmsDigitalEmployeeTags(v.(*schema.Set))
	}
	body := request
	wait := incrementalWait(3*time.Second, 5*time.Second)
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
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_digital_employee", action, AlibabaCloudSdkGoERROR)
	}

	name, _ := response["name"].(string)
	if name == "" {
		name = d.Get("digital_employee_name").(string)
	}
	d.SetId(name)

	return resourceAliCloudCmsDigitalEmployeeRead(d, meta)
}

func resourceAliCloudCmsDigitalEmployeeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	object, err := cmsServiceV2.DescribeCmsDigitalEmployee(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_digital_employee DescribeCmsDigitalEmployee Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("digital_employee_name", object["name"])
	d.Set("role_arn", object["roleArn"])
	d.Set("default_rule", object["defaultRule"])
	d.Set("description", object["description"])
	d.Set("display_name", object["displayName"])
	d.Set("resource_group_id", object["resourceGroupId"])
	d.Set("create_time", object["createTime"])
	d.Set("update_time", object["updateTime"])
	d.Set("employee_type", object["employeeType"])
	d.Set("region_id", object["regionId"])
	d.Set("resource_type", "ALIYUN::CMS::DIGITALEMPLOYEE")
	if attributes, ok := object["attributes"].(map[string]interface{}); ok {
		converted := make(map[string]interface{})
		for k, vv := range attributes {
			converted[k] = fmt.Sprintf("%v", vv)
		}
		d.Set("attributes", converted)
	}
	if knowledges, ok := object["knowledges"].(map[string]interface{}); ok {
		d.Set("knowledges", flattenCmsDigitalEmployeeKnowledges(knowledges))
	}
	if sandboxNetworkPolicy, ok := object["sandboxNetworkPolicy"].(map[string]interface{}); ok {
		d.Set("sandbox_network_policy", flattenCmsDigitalEmployeeSandboxNetworkPolicy(sandboxNetworkPolicy))
	}
	if toolPolicy, ok := object["toolPolicy"].(map[string]interface{}); ok {
		d.Set("tool_policy", flattenCmsDigitalEmployeeToolPolicy(toolPolicy))
	}
	if tags, ok := object["tags"].([]interface{}); ok && len(tags) > 0 {
		d.Set("tags", flattenCmsDigitalEmployeeTags(tags))
	}

	return nil
}

func resourceAliCloudCmsDigitalEmployeeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/digital-employee/%s", d.Id())
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	update := false

	if d.HasChange("role_arn") {
		update = true
	}
	request["roleArn"] = d.Get("role_arn")
	if d.HasChange("default_rule") {
		update = true
	}
	if v, ok := d.GetOk("default_rule"); ok {
		request["defaultRule"] = v
	}
	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if d.HasChange("display_name") {
		update = true
	}
	if v, ok := d.GetOk("display_name"); ok {
		request["displayName"] = v
	}
	if d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["resourceGroupId"] = v
	}
	if d.HasChange("attributes") {
		update = true
	}
	if v, ok := d.GetOk("attributes"); ok {
		request["attributes"] = v
	}
	if d.HasChange("knowledges") {
		update = true
	}
	if v, ok := d.GetOk("knowledges"); ok {
		request["knowledges"] = expandCmsDigitalEmployeeKnowledges(v.([]interface{}))
	}
	if d.HasChange("sandbox_network_policy") {
		update = true
	}
	if v, ok := d.GetOk("sandbox_network_policy"); ok {
		request["sandboxNetworkPolicy"] = expandCmsDigitalEmployeeSandboxNetworkPolicy(v.([]interface{}))
	}
	if d.HasChange("tool_policy") {
		update = true
	}
	if v, ok := d.GetOk("tool_policy"); ok {
		request["toolPolicy"] = expandCmsDigitalEmployeeToolPolicy(v.([]interface{}))
	}
	if d.HasChange("tags") {
		update = true
	}
	if v, ok := d.GetOk("tags"); ok {
		request["tags"] = expandCmsDigitalEmployeeTags(v.(*schema.Set))
	}

	body := request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPatch("Cms", "2024-03-30", action, query, nil, body, true)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudCmsDigitalEmployeeRead(d, meta)
}

func resourceAliCloudCmsDigitalEmployeeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/digital-employee/%s", d.Id())
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
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
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDigitalEmployee.NotFound", "404"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func expandCmsDigitalEmployeeTags(s *schema.Set) []map[string]interface{} {
	if s == nil || s.Len() == 0 {
		return nil
	}
	tags := make([]map[string]interface{}, 0, s.Len())
	for _, raw := range s.List() {
		if m, ok := raw.(map[string]interface{}); ok {
			item := map[string]interface{}{
				"key":   m["key"],
				"value": m["value"],
			}
			tags = append(tags, item)
		}
	}
	return tags
}

func flattenCmsDigitalEmployeeTags(tags []interface{}) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(tags))
	for _, raw := range tags {
		if m, ok := raw.(map[string]interface{}); ok {
			item := map[string]interface{}{
				"key":   m["key"],
				"value": m["value"],
			}
			out = append(out, item)
		}
	}
	return out
}

func expandCmsDigitalEmployeeKnowledges(list []interface{}) map[string]interface{} {
	if len(list) == 0 {
		return nil
	}
	m, ok := list[0].(map[string]interface{})
	if !ok {
		return nil
	}
	out := map[string]interface{}{}
	if bailian, ok := m["bailian"].([]interface{}); ok && len(bailian) > 0 {
		items := make([]map[string]interface{}, 0, len(bailian))
		for _, raw := range bailian {
			if b, ok := raw.(map[string]interface{}); ok {
				items = append(items, map[string]interface{}{
					"workspaceId": b["workspace_id"],
					"indexId":     b["index_id"],
					"region":      b["region"],
					"attributes":  b["attributes"],
				})
			}
		}
		out["bailian"] = items
	}
	return out
}

func flattenCmsDigitalEmployeeKnowledges(knowledges map[string]interface{}) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, 1)
	item := map[string]interface{}{}
	if bailian, ok := knowledges["bailian"].([]interface{}); ok {
		items := make([]map[string]interface{}, 0, len(bailian))
		for _, raw := range bailian {
			if b, ok := raw.(map[string]interface{}); ok {
				items = append(items, map[string]interface{}{
					"workspace_id": b["workspaceId"],
					"index_id":     b["indexId"],
					"region":       b["region"],
					"attributes":   b["attributes"],
				})
			}
		}
		item["bailian"] = items
	}
	out = append(out, item)
	return out
}

func expandCmsDigitalEmployeeSandboxNetworkPolicy(list []interface{}) map[string]interface{} {
	if len(list) == 0 {
		return nil
	}
	m, ok := list[0].(map[string]interface{})
	if !ok {
		return nil
	}
	out := map[string]interface{}{}
	if v, ok := m["allow_fqdns"].([]interface{}); ok {
		out["allowFqdns"] = toStringSlice(v)
	}
	if v, ok := m["allow_cidrs"].([]interface{}); ok {
		out["allowCidrs"] = toStringSlice(v)
	}
	out["enableAcl"] = m["enable_acl"]
	return out
}

func flattenCmsDigitalEmployeeSandboxNetworkPolicy(p map[string]interface{}) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, 1)
	item := map[string]interface{}{}
	if v, ok := p["allowFqdns"].([]interface{}); ok {
		item["allow_fqdns"] = toStringSlice(v)
	}
	if v, ok := p["allowCidrs"].([]interface{}); ok {
		item["allow_cidrs"] = toStringSlice(v)
	}
	item["enable_acl"] = p["enableAcl"]
	out = append(out, item)
	return out
}

func expandCmsDigitalEmployeeToolPolicy(list []interface{}) map[string]interface{} {
	if len(list) == 0 {
		return nil
	}
	m, ok := list[0].(map[string]interface{})
	if !ok {
		return nil
	}
	out := map[string]interface{}{}
	if aliyunList, ok := m["aliyun"].([]interface{}); ok && len(aliyunList) > 0 {
		if a, ok := aliyunList[0].(map[string]interface{}); ok {
			aliyunOut := map[string]interface{}{
				"enable": a["enable"],
			}
			if v, ok := a["deny_policy"].([]interface{}); ok {
				aliyunOut["denyPolicy"] = toStringSlice(v)
			}
			if v, ok := a["auto_pass_policy"].([]interface{}); ok {
				aliyunOut["autoPassPolicy"] = toStringSlice(v)
			}
			if v, ok := a["statements"].([]interface{}); ok && len(v) > 0 {
				items := make([]map[string]interface{}, 0, len(v))
				for _, raw := range v {
					if s, ok := raw.(map[string]interface{}); ok {
						stmt := map[string]interface{}{
							"decision":   s["decision"],
							"product":    s["product"],
							"apiVersion": s["api_version"],
						}
						if actions, ok := s["actions"].([]interface{}); ok {
							stmt["actions"] = toStringSlice(actions)
						}
						items = append(items, stmt)
					}
				}
				aliyunOut["statements"] = items
			}
			out["aliyun"] = aliyunOut
		}
	}
	return out
}

func flattenCmsDigitalEmployeeToolPolicy(toolPolicy map[string]interface{}) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, 1)
	item := map[string]interface{}{}
	if aliyun, ok := toolPolicy["aliyun"].(map[string]interface{}); ok {
		aliyunItem := map[string]interface{}{
			"enable": aliyun["enable"],
		}
		if v, ok := aliyun["denyPolicy"].([]interface{}); ok {
			aliyunItem["deny_policy"] = toStringSlice(v)
		}
		if v, ok := aliyun["autoPassPolicy"].([]interface{}); ok {
			aliyunItem["auto_pass_policy"] = toStringSlice(v)
		}
		if v, ok := aliyun["statements"].([]interface{}); ok {
			items := make([]map[string]interface{}, 0, len(v))
			for _, raw := range v {
				if s, ok := raw.(map[string]interface{}); ok {
					stmt := map[string]interface{}{
						"decision":    s["decision"],
						"product":     s["product"],
						"api_version": s["apiVersion"],
					}
					if actions, ok := s["actions"].([]interface{}); ok {
						stmt["actions"] = toStringSlice(actions)
					}
					items = append(items, stmt)
				}
			}
			aliyunItem["statements"] = items
		}
		item["aliyun"] = []map[string]interface{}{aliyunItem}
	}
	out = append(out, item)
	return out
}

func toStringSlice(list []interface{}) []string {
	out := make([]string, 0, len(list))
	for _, v := range list {
		out = append(out, fmt.Sprintf("%v", v))
	}
	return out
}
