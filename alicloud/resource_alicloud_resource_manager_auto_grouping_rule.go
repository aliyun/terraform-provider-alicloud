// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"time"
)

func resourceAliCloudResourceManagerAutoGroupingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudResourceManagerAutoGroupingRuleCreate,
		Read:   resourceAliCloudResourceManagerAutoGroupingRuleRead,
		Update: resourceAliCloudResourceManagerAutoGroupingRuleUpdate,
		Delete: resourceAliCloudResourceManagerAutoGroupingRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"exclude_region_ids_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"exclude_resource_group_ids_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"exclude_resource_ids_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"exclude_resource_types_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_ids_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_ids_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_ids_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_types_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule_contents": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_grouping_scope_condition": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								equal, _ := compareJsonTemplateAreEquivalent(old, new)
								return equal
							},
						},
						"target_resource_group_condition": {
							Type:     schema.TypeString,
							Required: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								equal, _ := compareJsonTemplateAreEquivalent(old, new)
								return equal
							},
						},
					},
				},
			},
			"rule_desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"associated_transfer", "custom_condition"}, false),
			},
		},
	}
}

func resourceAliCloudResourceManagerAutoGroupingRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAutoGroupingRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("exclude_resource_types_scope"); ok {
		request["ExcludeResourceTypesScope"] = v
	}
	if v, ok := d.GetOk("resource_ids_scope"); ok {
		request["ResourceIdsScope"] = v
	}
	if v, ok := d.GetOk("exclude_resource_ids_scope"); ok {
		request["ExcludeResourceIdsScope"] = v
	}
	if v, ok := d.GetOk("rule_desc"); ok {
		request["RuleDesc"] = v
	}

	if v, ok := d.GetOk("rule_contents"); ok {
		ruleContentsMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["TargetResourceGroupCondition"] = dataLoopTmp["target_resource_group_condition"]
			dataLoopMap["AutoGroupingScopeCondition"] = dataLoopTmp["auto_grouping_scope_condition"]
			ruleContentsMapsArray = append(ruleContentsMapsArray, dataLoopMap)
		}
		request["RuleContents"] = ruleContentsMapsArray
	}

	if v, ok := d.GetOk("resource_types_scope"); ok {
		request["ResourceTypesScope"] = v
	}
	request["RuleName"] = d.Get("rule_name")
	if v, ok := d.GetOk("exclude_resource_group_ids_scope"); ok {
		request["ExcludeResourceGroupIdsScope"] = v
	}
	if v, ok := d.GetOk("region_ids_scope"); ok {
		request["RegionIdsScope"] = v
	}
	if v, ok := d.GetOk("resource_group_ids_scope"); ok {
		request["ResourceGroupIdsScope"] = v
	}
	if v, ok := d.GetOk("exclude_region_ids_scope"); ok {
		request["ExcludeRegionIdsScope"] = v
	}
	request["RuleType"] = d.Get("rule_type")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_auto_grouping_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["RuleId"]))

	return resourceAliCloudResourceManagerAutoGroupingRuleRead(d, meta)
}

func resourceAliCloudResourceManagerAutoGroupingRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerServiceV2 := ResourceManagerServiceV2{client}

	objectRaw, err := resourceManagerServiceV2.DescribeResourceManagerAutoGroupingRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_auto_grouping_rule DescribeResourceManagerAutoGroupingRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("exclude_region_ids_scope", objectRaw["ExcludeRegionIdsScope"])
	d.Set("exclude_resource_group_ids_scope", objectRaw["ExcludeResourceGroupIdsScope"])
	d.Set("exclude_resource_ids_scope", objectRaw["ExcludeResourceIdsScope"])
	d.Set("exclude_resource_types_scope", objectRaw["ExcludeResourceTypesScope"])
	d.Set("region_ids_scope", objectRaw["RegionIdsScope"])
	d.Set("resource_group_ids_scope", objectRaw["ResourceGroupIdsScope"])
	d.Set("resource_ids_scope", objectRaw["ResourceIdsScope"])
	d.Set("resource_types_scope", objectRaw["ResourceTypesScope"])
	d.Set("rule_desc", objectRaw["RuleDesc"])
	d.Set("rule_name", objectRaw["RuleName"])
	d.Set("rule_type", objectRaw["RuleType"])

	ruleContentsRaw := objectRaw["RuleContents"]
	ruleContentsMaps := make([]map[string]interface{}, 0)
	if ruleContentsRaw != nil {
		for _, ruleContentsChildRaw := range ruleContentsRaw.([]interface{}) {
			ruleContentsMap := make(map[string]interface{})
			ruleContentsChildRaw := ruleContentsChildRaw.(map[string]interface{})
			ruleContentsMap["auto_grouping_scope_condition"] = ruleContentsChildRaw["AutoGroupingScopeCondition"]
			ruleContentsMap["target_resource_group_condition"] = ruleContentsChildRaw["TargetResourceGroupCondition"]

			ruleContentsMaps = append(ruleContentsMaps, ruleContentsMap)
		}
	}
	if err := d.Set("rule_contents", ruleContentsMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudResourceManagerAutoGroupingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateAutoGroupingRule"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RuleId"] = d.Id()

	if d.HasChange("resource_ids_scope") || d.HasChange("resource_types_scope") ||
		d.HasChange("region_ids_scope") || d.HasChange("resource_group_ids_scope") {
		update = true

	}

	if d.HasChange("exclude_resource_types_scope") {
		update = true
	}
	if v, ok := d.GetOk("exclude_resource_types_scope"); ok {
		request["ExcludeResourceTypesScope"] = v
	}

	if d.HasChange("resource_ids_scope") {
		update = true
	}
	if v, ok := d.GetOk("resource_ids_scope"); ok {
		request["ResourceIdsScope"] = v
	}

	if d.HasChange("exclude_resource_ids_scope") {
		update = true
	}
	if v, ok := d.GetOk("exclude_resource_ids_scope"); ok {
		request["ExcludeResourceIdsScope"] = v
	}

	if d.HasChange("rule_desc") {
		update = true
	}
	if v, ok := d.GetOk("rule_desc"); ok {
		request["RuleDesc"] = v
	}

	if d.HasChange("rule_contents") {
		update = true
	}
	if v, ok := d.GetOk("rule_contents"); ok {
		ruleContentsMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["TargetResourceGroupCondition"] = dataLoopTmp["target_resource_group_condition"]
			dataLoopMap["AutoGroupingScopeCondition"] = dataLoopTmp["auto_grouping_scope_condition"]
			ruleContentsMapsArray = append(ruleContentsMapsArray, dataLoopMap)
		}
		request["RuleContents"] = ruleContentsMapsArray
	}

	if d.HasChange("resource_types_scope") {
		update = true
	}
	if v, ok := d.GetOk("resource_types_scope"); ok {
		request["ResourceTypesScope"] = v
	}

	if d.HasChange("rule_name") {
		update = true
	}
	request["RuleName"] = d.Get("rule_name")
	if d.HasChange("exclude_resource_group_ids_scope") {
		update = true
	}
	if v, ok := d.GetOk("exclude_resource_group_ids_scope"); ok {
		request["ExcludeResourceGroupIdsScope"] = v
	}

	if d.HasChange("region_ids_scope") {
		update = true
	}
	if v, ok := d.GetOk("region_ids_scope"); ok {
		request["RegionIdsScope"] = v
	}

	if d.HasChange("resource_group_ids_scope") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_ids_scope"); ok {
		request["ResourceGroupIdsScope"] = v
	}

	if d.HasChange("exclude_region_ids_scope") {
		update = true
	}
	if v, ok := d.GetOk("exclude_region_ids_scope"); ok {
		request["ExcludeRegionIdsScope"] = v
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ResourceManager", "2020-03-31", action, query, request, true)
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

	return resourceAliCloudResourceManagerAutoGroupingRuleRead(d, meta)
}

func resourceAliCloudResourceManagerAutoGroupingRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAutoGroupingRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RuleId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, query, request, true)

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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
