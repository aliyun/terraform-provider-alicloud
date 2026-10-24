package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlicloudNasLifecyclePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudNasLifecyclePolicyCreate,
		Read:   resourceAlicloudNasLifecyclePolicyRead,
		Update: resourceAlicloudNasLifecyclePolicyUpdate,
		Delete: resourceAlicloudNasLifecyclePolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"file_system_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"lifecycle_policy_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"lifecycle_rule_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"DEFAULT_ATIME_14", "DEFAULT_ATIME_30", "DEFAULT_ATIME_60", "DEFAULT_ATIME_90"}, false),
			},
			"paths": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				MaxItems:      10,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"path"},
			},
			"path": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"paths"},
			},
			"storage_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"InfrequentAccess"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lifecycle_policy_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "Auto",
				ValidateFunc: validation.StringInSlice([]string{"Auto", "OnDemand"}, false),
			},
			"retrieve_rules": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attribute": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"RetrieveType"}, false),
						},
						"threshold": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"AfterVisit", "All"}, false),
						},
					},
				},
			},
			"transit_rules": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attribute": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"Atime"}, false),
						},
						"threshold": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"lifecycle_policy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudNasLifecyclePolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateLifecyclePolicy"
	request := make(map[string]interface{})
	var err error
	request["FileSystemId"] = d.Get("file_system_id")
	request["LifecyclePolicyName"] = d.Get("lifecycle_policy_name")
	request["LifecycleRuleName"] = d.Get("lifecycle_rule_name")
	request["StorageType"] = d.Get("storage_type")
	if v, ok := d.GetOk("paths"); ok && len(v.([]interface{})) > 0 {
		request["Paths"] = v
	}
	if v, ok := d.GetOk("path"); ok && v.(string) != "" {
		request["Path"] = v
	}
	if len(d.Get("paths").([]interface{})) == 0 && d.Get("path").(string) == "" {
		return WrapError(fmt.Errorf("one of `path` or `paths` must be specified when creating alicloud_nas_lifecycle_policy"))
	}
	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request["Description"] = v
	}
	if v, ok := d.GetOk("lifecycle_policy_type"); ok && v.(string) != "" {
		request["LifecyclePolicyType"] = v
	}
	if v, ok := d.GetOk("retrieve_rules"); ok {
		for i, rule := range v.([]interface{}) {
			ruleArg := rule.(map[string]interface{})
			request[fmt.Sprintf("RetrieveRules.%d.Attribute", i+1)] = ruleArg["attribute"]
			request[fmt.Sprintf("RetrieveRules.%d.Threshold", i+1)] = ruleArg["threshold"]
		}
	}
	if v, ok := d.GetOk("transit_rules"); ok {
		for i, rule := range v.([]interface{}) {
			ruleArg := rule.(map[string]interface{})
			request[fmt.Sprintf("TransitRules.%d.Attribute", i+1)] = ruleArg["attribute"]
			request[fmt.Sprintf("TransitRules.%d.Threshold", i+1)] = ruleArg["threshold"]
		}
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nas_lifecycle_policy", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["FileSystemId"], ":", request["LifecyclePolicyName"]))

	return resourceAlicloudNasLifecyclePolicyRead(d, meta)
}
func resourceAlicloudNasLifecyclePolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	object, err := nasService.DescribeNasLifecyclePolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nas_lifecycle_policy nasService.DescribeNasLifecyclePolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("file_system_id", object["FileSystemId"])
	d.Set("lifecycle_policy_name", object["LifecyclePolicyName"])
	d.Set("lifecycle_rule_name", object["LifecycleRuleName"])
	d.Set("storage_type", object["StorageType"])
	d.Set("description", object["Description"])
	d.Set("lifecycle_policy_type", object["LifecyclePolicyType"])
	d.Set("lifecycle_policy_id", object["LifecyclePolicyId"])
	// `path` and `paths` are mutually exclusive. Only populate `path` when the
	// API does not return a non-empty `Paths`, so that a paths-based config does
	// not drift on an echoed `Path` value and a path-based config is read back.
	pathsVal := object["Paths"]
	d.Set("paths", pathsVal)
	if pathsList, ok := pathsVal.([]interface{}); !ok || len(pathsList) == 0 {
		d.Set("path", object["Path"])
	} else {
		d.Set("path", "")
	}

	if retrieveRulesRaw, ok := object["RetrieveRules"].([]interface{}); ok && retrieveRulesRaw != nil {
		retrieveRules := make([]map[string]interface{}, 0, len(retrieveRulesRaw))
		for _, rule := range retrieveRulesRaw {
			if ruleMap, ok := rule.(map[string]interface{}); ok {
				retrieveRules = append(retrieveRules, map[string]interface{}{
					"attribute": ruleMap["Attribute"],
					"threshold": ruleMap["Threshold"],
				})
			}
		}
		if err := d.Set("retrieve_rules", retrieveRules); err != nil {
			return WrapError(err)
		}
	}
	if transitRulesRaw, ok := object["TransitRules"].([]interface{}); ok && transitRulesRaw != nil {
		transitRules := make([]map[string]interface{}, 0, len(transitRulesRaw))
		for _, rule := range transitRulesRaw {
			if ruleMap, ok := rule.(map[string]interface{}); ok {
				transitRules = append(transitRules, map[string]interface{}{
					"attribute": ruleMap["Attribute"],
					"threshold": ruleMap["Threshold"],
				})
			}
		}
		if err := d.Set("transit_rules", transitRules); err != nil {
			return WrapError(err)
		}
	}
	return nil
}
func resourceAlicloudNasLifecyclePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"FileSystemId":        parts[0],
		"LifecyclePolicyName": parts[1],
	}
	if v, ok := d.GetOk("lifecycle_policy_id"); ok && v.(string) != "" {
		request["LifecyclePolicyId"] = v
	}
	request["LifecycleRuleName"] = d.Get("lifecycle_rule_name")
	if d.HasChange("lifecycle_rule_name") {
		update = true
	}
	if d.HasChange("description") {
		request["Description"] = d.Get("description")
		update = true
	}
	if update {
		action := "ModifyLifecyclePolicy"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("NAS", "2017-06-26", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudNasLifecyclePolicyRead(d, meta)
}
func resourceAlicloudNasLifecyclePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteLifecyclePolicy"
	var response map[string]interface{}
	request := map[string]interface{}{
		"FileSystemId":        parts[0],
		"LifecyclePolicyName": parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidFileSystem.NotFound", "InvalidLifecyclePolicy.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
