package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCmsTransformer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsTransformerCreate,
		Read:   resourceAliCloudCmsTransformerRead,
		Update: resourceAliCloudCmsTransformerUpdate,
		Delete: resourceAliCloudCmsTransformerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"transformer_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"transformer_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"quit_after_match": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"sort_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"actions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filter_setting": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"conditions": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"field": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"op": {
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
									"expression": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"relation": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"label_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"mapping": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"reg_exp": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"source": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"target": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"variable": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"filter_setting": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"conditions": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"op": {
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
						"expression": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"relation": {
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
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCmsTransformerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "/transformers"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error

	workspace := d.Get("workspace").(string)
	query["workspace"] = StringPointer(workspace)
	query["regionId"] = StringPointer(client.RegionId)

	request = buildCmsTransformerRequestBody(d)
	body = request

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_transformer", action, AlibabaCloudSdkGoERROR)
	}

	transformerId, _ := jsonpath.Get("$.transformerId", response)
	if transformerId == nil || transformerId == "" {
		return WrapError(fmt.Errorf("failed to create CMS Transformer: transformerId is empty in response"))
	}

	d.SetId(fmt.Sprintf("%v:%v", transformerId, workspace))

	return resourceAliCloudCmsTransformerUpdate(d, meta)
}

func resourceAliCloudCmsTransformerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	objectRaw, err := cmsServiceV2.DescribeCmsTransformer(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_transformer DescribeCmsTransformer Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	object := objectRaw

	d.Set("transformer_id", object["transformerId"])
	d.Set("transformer_name", object["transformerName"])
	d.Set("workspace", object["workspace"])
	d.Set("description", object["description"])
	d.Set("quit_after_match", object["quitAfterMatch"])
	d.Set("sort_id", object["sortId"])
	d.Set("enable", object["enable"])
	d.Set("create_time", object["createTime"])
	d.Set("update_time", object["updateTime"])
	d.Set("user_id", object["userId"])
	d.Set("region_id", object["regionId"])

	if actionsRaw, ok := object["actions"]; ok && actionsRaw != nil {
		actionsMaps := make([]map[string]interface{}, 0)
		for _, actionItem := range convertToInterfaceArray(actionsRaw) {
			actionMap := make(map[string]interface{})
			actionItemMap, ok := actionItem.(map[string]interface{})
			if !ok {
				actionsMaps = append(actionsMaps, actionMap)
				continue
			}
			actionMap["label_key"] = actionItemMap["labelKey"]
			actionMap["reg_exp"] = actionItemMap["regExp"]
			actionMap["source"] = actionItemMap["source"]
			actionMap["target"] = actionItemMap["target"]
			actionMap["type"] = actionItemMap["type"]
			actionMap["value"] = actionItemMap["value"]
			actionMap["variable"] = actionItemMap["variable"]
			actionMap["mapping"] = actionItemMap["mapping"]

			filterSettingMaps := make([]map[string]interface{}, 0)
			if fsRaw, ok := actionItemMap["filterSetting"]; ok && fsRaw != nil {
				if fsMap, ok := fsRaw.(map[string]interface{}); ok && len(fsMap) > 0 {
					filterSettingMap := buildCmsTransformerFilterSettingRead(fsMap)
					filterSettingMaps = append(filterSettingMaps, filterSettingMap)
				}
			}
			actionMap["filter_setting"] = filterSettingMaps

			actionsMaps = append(actionsMaps, actionMap)
		}
		d.Set("actions", actionsMaps)
	} else {
		d.Set("actions", make([]map[string]interface{}, 0))
	}

	filterSettingMaps := make([]map[string]interface{}, 0)
	if fsRaw, ok := object["filterSetting"]; ok && fsRaw != nil {
		if fsMap, ok := fsRaw.(map[string]interface{}); ok && len(fsMap) > 0 {
			filterSettingMap := buildCmsTransformerFilterSettingRead(fsMap)
			filterSettingMaps = append(filterSettingMaps, filterSettingMap)
		}
	}
	d.Set("filter_setting", filterSettingMaps)

	return nil
}

func resourceAliCloudCmsTransformerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length 2, got %d", d.Id(), len(parts)))
	}
	transformerId := parts[0]
	workspace := parts[1]

	// Handle enable/disable sub-operations when the enable field changes.
	if d.HasChange("enable") && !d.IsNewResource() {
		target := d.Get("enable").(bool)
		if target {
			if err := cmsTransformerToggleState(client, transformerId, workspace, "enable", d.Timeout(schema.TimeoutUpdate)); err != nil {
				return err
			}
		} else {
			if err := cmsTransformerToggleState(client, transformerId, workspace, "disable", d.Timeout(schema.TimeoutUpdate)); err != nil {
				return err
			}
		}
	}

	// General update via PUT /transformers/{transformerId}.
	action := fmt.Sprintf("/transformers/%s", transformerId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = buildCmsTransformerRequestBody(d)
	query["workspace"] = StringPointer(workspace)
	query["regionId"] = StringPointer(client.RegionId)
	body = request

	wait := incrementalWait(3*time.Second, 5*time.Second)
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
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return resourceAliCloudCmsTransformerRead(d, meta)
}

func resourceAliCloudCmsTransformerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length 2, got %d", d.Id(), len(parts)))
	}
	transformerId := parts[0]
	workspace := parts[1]

	action := fmt.Sprintf("/transformers/%s", transformerId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	query["workspace"] = StringPointer(workspace)
	query["regionId"] = StringPointer(client.RegionId)

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
		if IsExpectedErrors(err, []string{"NotFound", "ResourceNotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

// cmsTransformerToggleState calls the enable or disable sub-operation.
func cmsTransformerToggleState(client *connectivity.AliyunClient, transformerId, workspace, op string, timeout time.Duration) error {
	action := fmt.Sprintf("/transformers/%s/%s", transformerId, op)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["workspace"] = StringPointer(workspace)
	query["regionId"] = StringPointer(client.RegionId)
	body = request

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(timeout, func() *resource.RetryError {
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
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, transformerId, action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

// buildCmsTransformerRequestBody constructs the TransformerForModify request body
// from the Terraform schema.
func buildCmsTransformerRequestBody(d *schema.ResourceData) map[string]interface{} {
	request := make(map[string]interface{})

	if v, ok := d.GetOk("transformer_name"); ok {
		request["transformerName"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOkExists("quit_after_match"); ok {
		request["quitAfterMatch"] = v
	}
	if v, ok := d.GetOkExists("sort_id"); ok {
		request["sortId"] = v
	}

	// Top-level filter_setting
	if v, ok := d.GetOk("filter_setting"); ok {
		if fsList, ok := v.([]interface{}); ok && len(fsList) > 0 {
			filterSetting := buildCmsTransformerFilterSettingRequest(fsList[0])
			if len(filterSetting) > 0 {
				request["filterSetting"] = filterSetting
			}
		}
	}

	// Actions array
	if v, ok := d.GetOk("actions"); ok {
		actions := make([]interface{}, 0)
		for _, actionItem := range convertToInterfaceArray(v) {
			actionMap := make(map[string]interface{})
			actionItemMap, ok := actionItem.(map[string]interface{})
			if !ok {
				actions = append(actions, actionMap)
				continue
			}
			if val, ok := actionItemMap["label_key"].(string); ok && val != "" {
				actionMap["labelKey"] = val
			}
			if val, ok := actionItemMap["reg_exp"].(string); ok && val != "" {
				actionMap["regExp"] = val
			}
			if val, ok := actionItemMap["source"].(string); ok && val != "" {
				actionMap["source"] = val
			}
			if val, ok := actionItemMap["target"].(string); ok && val != "" {
				actionMap["target"] = val
			}
			if val, ok := actionItemMap["type"].(string); ok && val != "" {
				actionMap["type"] = val
			}
			if val, ok := actionItemMap["value"].(string); ok && val != "" {
				actionMap["value"] = val
			}
			if val, ok := actionItemMap["variable"].(string); ok && val != "" {
				actionMap["variable"] = val
			}
			if mapping, ok := actionItemMap["mapping"].(map[string]interface{}); ok && len(mapping) > 0 {
				actionMap["mapping"] = mapping
			}

			if fsList, ok := actionItemMap["filter_setting"].([]interface{}); ok && len(fsList) > 0 {
				filterSetting := buildCmsTransformerFilterSettingRequest(fsList[0])
				if len(filterSetting) > 0 {
					actionMap["filterSetting"] = filterSetting
				}
			}

			actions = append(actions, actionMap)
		}
		request["actions"] = actions
	}

	return request
}

// buildCmsTransformerFilterSettingRequest converts a schema filter_setting block
// into the API request shape.
func buildCmsTransformerFilterSettingRequest(raw interface{}) map[string]interface{} {
	filterSetting := make(map[string]interface{})
	filterSettingMap, ok := raw.(map[string]interface{})
	if !ok {
		return filterSetting
	}

	conditionsRaw, ok := filterSettingMap["conditions"].([]interface{})
	if ok {
		conditions := make([]interface{}, 0)
		for _, conditionItem := range convertToInterfaceArray(conditionsRaw) {
			conditionMap := make(map[string]interface{})
			conditionItemMap, ok := conditionItem.(map[string]interface{})
			if !ok {
				conditions = append(conditions, conditionMap)
				continue
			}
			if val, ok := conditionItemMap["field"].(string); ok {
				conditionMap["field"] = val
			}
			if val, ok := conditionItemMap["op"].(string); ok {
				conditionMap["op"] = val
			}
			if val, ok := conditionItemMap["value"].(string); ok {
				conditionMap["value"] = val
			}
			conditions = append(conditions, conditionMap)
		}
		filterSetting["conditions"] = conditions
	}

	if val, ok := filterSettingMap["expression"].(string); ok && val != "" {
		filterSetting["expression"] = val
	}
	if val, ok := filterSettingMap["relation"].(string); ok && val != "" {
		filterSetting["relation"] = val
	}

	return filterSetting
}

// buildCmsTransformerFilterSettingRead converts an API filter_setting map
// into the Terraform schema shape.
func buildCmsTransformerFilterSettingRead(fsMap map[string]interface{}) map[string]interface{} {
	filterSettingMap := make(map[string]interface{})

	conditionsMaps := make([]map[string]interface{}, 0)
	if conditionsRaw, ok := fsMap["conditions"]; ok && conditionsRaw != nil {
		for _, conditionItem := range convertToInterfaceArray(conditionsRaw) {
			conditionMap := make(map[string]interface{})
			conditionItemMap, ok := conditionItem.(map[string]interface{})
			if !ok {
				conditionsMaps = append(conditionsMaps, conditionMap)
				continue
			}
			conditionMap["field"] = conditionItemMap["field"]
			conditionMap["op"] = conditionItemMap["op"]
			conditionMap["value"] = conditionItemMap["value"]
			conditionsMaps = append(conditionsMaps, conditionMap)
		}
	}
	filterSettingMap["conditions"] = conditionsMaps
	filterSettingMap["expression"] = fsMap["expression"]
	filterSettingMap["relation"] = fsMap["relation"]

	return filterSettingMap
}
