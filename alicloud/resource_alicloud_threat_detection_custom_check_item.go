package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudThreatDetectionCustomCheckItem() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudThreatDetectionCustomCheckItemCreate,
		Read:   resourceAlicloudThreatDetectionCustomCheckItemRead,
		Update: resourceAlicloudThreatDetectionCustomCheckItemUpdate,
		Delete: resourceAlicloudThreatDetectionCustomCheckItemDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"check_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the custom check item.",
			},
			"check_show_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name of the check item.",
			},
			"section_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "The list of section IDs associated with the custom check item.",
			},
			"vendor": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The vendor to which the check item belongs.",
			},
			"instance_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The asset type to which the check item belongs.",
			},
			"instance_sub_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The asset subtype to which the check item belongs.",
			},
			"risk_level": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The risk level of the check item.",
			},
			"status": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"RELEASE", "EDIT"}, false),
				Description:  "The status of the check item. Valid values: `RELEASE`, `EDIT`.",
			},
			"check_rule": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The check rule of the check item.",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The remarks of the check item.",
			},
			"description": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "The description information of the check item.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The type of the description information.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The value of the description information.",
						},
					},
				},
			},
			"assist_info": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "The description help information of the check item.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The type of the description help information.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The value of the description help information.",
						},
					},
				},
			},
			"solution": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "The solution of the check item.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The type of the solution.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The value of the solution.",
						},
					},
				},
			},
			"check_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the check item.",
			},
		},
	}
}

func resourceAlicloudThreatDetectionCustomCheckItemCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateCheckItem"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if err := setCustomCheckItemRequest(d, request, false); err != nil {
		return err
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_custom_check_item", action, AlibabaCloudSdkGoERROR)
	}

	id, err := jsonpath.Get("$.Data.CheckId", response)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_custom_check_item", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(id))

	return resourceAlicloudThreatDetectionCustomCheckItemRead(d, meta)
}

func resourceAlicloudThreatDetectionCustomCheckItemRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionServiceV2 := ThreatDetectionServiceV2{client}

	objectRaw, err := threatDetectionServiceV2.DescribeThreatDetectionCustomCheckItem(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_custom_check_item DescribeThreatDetectionCustomCheckItem Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("check_id", objectRaw["CheckId"])
	d.Set("check_show_name", objectRaw["CheckShowName"])
	d.Set("vendor", objectRaw["Vendor"])
	d.Set("instance_type", objectRaw["InstanceType"])
	d.Set("instance_sub_type", objectRaw["InstanceSubType"])
	d.Set("risk_level", objectRaw["RiskLevel"])
	d.Set("status", objectRaw["Status"])
	d.Set("check_rule", objectRaw["CheckRule"])
	d.Set("remark", objectRaw["Remark"])
	d.Set("check_type", objectRaw["CheckType"])

	// section_ids is not returned by ListCheckItems, the config value is preserved in state.

	if err := d.Set("description", flattenCustomCheckItemStruct(objectRaw["Description"])); err != nil {
		return WrapError(err)
	}
	if err := d.Set("assist_info", flattenCustomCheckItemStruct(objectRaw["AssistInfo"])); err != nil {
		return WrapError(err)
	}
	if err := d.Set("solution", flattenCustomCheckItemStruct(objectRaw["Solution"])); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAlicloudThreatDetectionCustomCheckItemUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "UpdateCheckItem"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	checkId, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return WrapError(err)
	}
	request["CheckId"] = checkId

	if err := setCustomCheckItemRequest(d, request, true); err != nil {
		return err
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)
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

	return resourceAlicloudThreatDetectionCustomCheckItemRead(d, meta)
}

func resourceAlicloudThreatDetectionCustomCheckItemDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteCheckItem"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	checkId, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return WrapError(err)
	}
	request["CheckIds.1"] = checkId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)
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
		if NotFoundError(err) {
			return nil
		}
		if IsExpectedErrors(err, []string{"CspmDeleteCheckCustomItemError"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

// setCustomCheckItemRequest populates the request map with the schema fields shared by
// CreateCheckItem and UpdateCheckItem. When forUpdate is true, only fields with a value
// or that have changed are sent.
func setCustomCheckItemRequest(d *schema.ResourceData, request map[string]interface{}, forUpdate bool) error {
	if v, ok := d.GetOk("check_show_name"); ok && (!forUpdate || d.HasChange("check_show_name")) {
		request["CheckShowName"] = v
	} else if !forUpdate {
		request["CheckShowName"] = d.Get("check_show_name")
	}
	if v, ok := d.GetOk("vendor"); ok && (!forUpdate || d.HasChange("vendor")) {
		request["Vendor"] = v
	} else if !forUpdate {
		request["Vendor"] = d.Get("vendor")
	}
	if v, ok := d.GetOk("instance_type"); ok && (!forUpdate || d.HasChange("instance_type")) {
		request["InstanceType"] = v
	} else if !forUpdate {
		request["InstanceType"] = d.Get("instance_type")
	}
	if v, ok := d.GetOk("instance_sub_type"); ok && (!forUpdate || d.HasChange("instance_sub_type")) {
		request["InstanceSubType"] = v
	} else if !forUpdate {
		request["InstanceSubType"] = d.Get("instance_sub_type")
	}
	if v, ok := d.GetOk("risk_level"); ok && (!forUpdate || d.HasChange("risk_level")) {
		request["RiskLevel"] = v
	} else if !forUpdate {
		request["RiskLevel"] = d.Get("risk_level")
	}
	if v, ok := d.GetOk("status"); ok && (!forUpdate || d.HasChange("status")) {
		request["Status"] = v
	} else if !forUpdate {
		request["Status"] = d.Get("status")
	}
	if v, ok := d.GetOk("check_rule"); ok && (!forUpdate || d.HasChange("check_rule")) {
		request["CheckRule"] = v
	} else if !forUpdate {
		request["CheckRule"] = d.Get("check_rule")
	}
	if v, ok := d.GetOk("section_ids"); ok && (!forUpdate || d.HasChange("section_ids")) {
		request["SectionIds"] = convertToInterfaceArray(v)
	} else if !forUpdate {
		request["SectionIds"] = convertToInterfaceArray(d.Get("section_ids"))
	}
	if v, ok := d.GetOk("remark"); ok && (!forUpdate || d.HasChange("remark")) {
		request["Remark"] = v
	}

	if m := buildCustomCheckItemStructMap(d, "description"); len(m) > 0 {
		b, err := json.Marshal(m)
		if err != nil {
			return WrapError(err)
		}
		request["Description"] = string(b)
	}
	if m := buildCustomCheckItemStructMap(d, "assist_info"); len(m) > 0 {
		b, err := json.Marshal(m)
		if err != nil {
			return WrapError(err)
		}
		request["AssistInfo"] = string(b)
	}
	if m := buildCustomCheckItemStructMap(d, "solution"); len(m) > 0 {
		b, err := json.Marshal(m)
		if err != nil {
			return WrapError(err)
		}
		request["Solution"] = string(b)
	}
	return nil
}

// buildCustomCheckItemStructMap builds the map for description/assist_info/solution schema blocks.
func buildCustomCheckItemStructMap(d *schema.ResourceData, key string) map[string]interface{} {
	m := make(map[string]interface{})
	v := d.Get(key)
	items := convertToInterfaceArray(v)
	if len(items) == 0 {
		return m
	}
	item, ok := items[0].(map[string]interface{})
	if !ok {
		return m
	}
	if t, ok := item["type"].(string); ok && t != "" {
		m["Type"] = t
	}
	if val, ok := item["value"].(string); ok && val != "" {
		m["Value"] = val
	}
	return m
}

// flattenCustomCheckItemStruct flattens a description/assist_info/solution response object into
// the schema list form.
func flattenCustomCheckItemStruct(raw interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if raw == nil {
		return result
	}
	m, ok := raw.(map[string]interface{})
	if !ok || len(m) == 0 {
		return result
	}
	return append(result, map[string]interface{}{
		"type":  m["Type"],
		"value": m["Value"],
	})
}
