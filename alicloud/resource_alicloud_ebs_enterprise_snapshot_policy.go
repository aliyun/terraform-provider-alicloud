// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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

func resourceAliCloudEbsEnterpriseSnapshotPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEbsEnterpriseSnapshotPolicyCreate,
		Read:   resourceAliCloudEbsEnterpriseSnapshotPolicyRead,
		Update: resourceAliCloudEbsEnterpriseSnapshotPolicyUpdate,
		Delete: resourceAliCloudEbsEnterpriseSnapshotPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cross_region_copy_info": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"regions": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"retain_days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: IntAtLeast(1),
									},
									"region_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_snapshot_policy_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"retain_rule": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"number": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: IntAtLeast(1),
						},
						"time_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: IntAtLeast(1),
						},
						"time_unit": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"schedule": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cron_expression": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"special_retain_rules": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"rules": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"special_period_unit": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"time_interval": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"time_unit": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"storage_rule": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_immediate_access": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"tags": tagsSchema(),
			"target_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudEbsEnterpriseSnapshotPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateEnterpriseSnapshotPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["Name"] = d.Get("enterprise_snapshot_policy_name")
	request["TargetType"] = d.Get("target_type")
	if v, ok := d.GetOk("desc"); ok {
		request["Desc"] = v
	}
	objectDataLocalMap := make(map[string]interface{})
	if v := d.Get("schedule"); !IsNil(v) {
		nodeNative, _ := jsonpath.Get("$[0].cron_expression", d.Get("schedule"))
		if nodeNative != "" {
			objectDataLocalMap["CronExpression"] = nodeNative
		}
		request["Schedule"], err = convertMaptoJsonString(objectDataLocalMap)
		if err != nil {
			return WrapError(err)
		}
	}

	objectDataLocalMap1 := make(map[string]interface{})
	if v := d.Get("retain_rule"); !IsNil(v) {
		nodeNative1, _ := jsonpath.Get("$[0].number", d.Get("retain_rule"))
		if nodeNative1 != "" && nodeNative1.(int) > 0 {
			objectDataLocalMap1["Number"] = nodeNative1
		}
		nodeNative2, _ := jsonpath.Get("$[0].time_interval", d.Get("retain_rule"))
		if nodeNative2 != "" && nodeNative2.(int) > 0 {
			objectDataLocalMap1["TimeInterval"] = nodeNative2
		}
		nodeNative3, _ := jsonpath.Get("$[0].time_unit", d.Get("retain_rule"))
		if nodeNative3 != "" {
			objectDataLocalMap1["TimeUnit"] = nodeNative3
		}
		request["RetainRule"], err = convertMaptoJsonString(objectDataLocalMap1)
		if err != nil {
			return WrapError(err)
		}
	}

	objectDataLocalMap2 := make(map[string]interface{})
	if v := d.Get("storage_rule"); !IsNil(v) {
		nodeNative4, _ := jsonpath.Get("$[0].enable_immediate_access", d.Get("storage_rule"))
		if nodeNative4 != "" {
			objectDataLocalMap2["EnableImmediateAccess"] = nodeNative4
		}
		request["StorageRule"], err = convertMaptoJsonString(objectDataLocalMap2)
		if err != nil {
			return WrapError(err)
		}
	}

	objectDataLocalMap3 := make(map[string]interface{})
	if v := d.Get("cross_region_copy_info"); !IsNil(v) {
		nodeNative5, _ := jsonpath.Get("$[0].enabled", d.Get("cross_region_copy_info"))
		if nodeNative5 != "" {
			objectDataLocalMap3["Enabled"] = nodeNative5
		}
		if v, ok := d.GetOk("cross_region_copy_info"); ok {
			localData, err := jsonpath.Get("$[0].regions", v)
			if err != nil {
				return WrapError(err)
			}
			localMaps := make([]map[string]interface{}, 0)
			for _, dataLoop := range localData.([]interface{}) {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["RegionId"] = dataLoopTmp["region_id"]
				if dataLoopTmp["retain_days"].(int) > 0 {
					dataLoopMap["RetainDays"] = dataLoopTmp["retain_days"]
				}
				localMaps = append(localMaps, dataLoopMap)
			}
			objectDataLocalMap3["Regions"] = localMaps
		}
		request["CrossRegionCopyInfo"], err = convertMaptoJsonString(objectDataLocalMap3)
		if err != nil {
			return WrapError(err)
		}
	}

	if v, ok := d.GetOk("status"); ok {
		request["State"] = v
	}
	if _, ok := d.GetOk("tags"); ok {
		added, _ := parsingTags(d)
		count := 1
		for key, value := range added {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
	}

	objectDataLocalMap4 := make(map[string]interface{})
	if v := d.Get("special_retain_rules"); !IsNil(v) {
		nodeNative10, _ := jsonpath.Get("$[0].enabled", d.Get("special_retain_rules"))
		if nodeNative10 != "" {
			objectDataLocalMap4["Enabled"] = nodeNative10
		}
		if v, ok := d.GetOk("special_retain_rules"); ok {
			localData2, err := jsonpath.Get("$[0].rules", v)
			if err != nil {
				return WrapError(err)
			}
			localMaps1 := make([]map[string]interface{}, 0)
			for _, dataLoop2 := range localData2.([]interface{}) {
				dataLoop2Tmp := dataLoop2.(map[string]interface{})
				dataLoop2Map := make(map[string]interface{})
				dataLoop2Map["SpecialPeriodUnit"] = dataLoop2Tmp["special_period_unit"]
				dataLoop2Map["TimeInterval"] = dataLoop2Tmp["time_interval"]
				dataLoop2Map["TimeUnit"] = dataLoop2Tmp["time_unit"]
				localMaps1 = append(localMaps1, dataLoop2Map)
			}
			objectDataLocalMap4["Rules"] = localMaps1
		}
		request["SpecialRetainRules"], err = convertMaptoJsonString(objectDataLocalMap4)
		if err != nil {
			return WrapError(err)
		}
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ebs_enterprise_snapshot_policy", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["PolicyId"]))

	return resourceAliCloudEbsEnterpriseSnapshotPolicyRead(d, meta)
}

func resourceAliCloudEbsEnterpriseSnapshotPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ebsServiceV2 := EbsServiceV2{client}

	objectRaw, err := ebsServiceV2.DescribeEbsEnterpriseSnapshotPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ebs_enterprise_snapshot_policy DescribeEbsEnterpriseSnapshotPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("desc", objectRaw["Desc"])
	d.Set("enterprise_snapshot_policy_name", objectRaw["Name"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["State"])
	d.Set("target_type", objectRaw["TargetType"])

	crossRegionCopyInfoMaps := make([]map[string]interface{}, 0)
	crossRegionCopyInfoMap := make(map[string]interface{})
	crossRegionCopyInfo1Raw := make(map[string]interface{})
	if objectRaw["CrossRegionCopyInfo"] != nil {
		crossRegionCopyInfo1Raw = objectRaw["CrossRegionCopyInfo"].(map[string]interface{})
	}
	if len(crossRegionCopyInfo1Raw) > 0 {
		crossRegionCopyInfoMap["enabled"] = crossRegionCopyInfo1Raw["Enabled"]

		regions1Raw := crossRegionCopyInfo1Raw["Regions"]
		regionsMaps := make([]map[string]interface{}, 0)
		if regions1Raw != nil {
			for _, regionsChild1Raw := range regions1Raw.([]interface{}) {
				regionsMap := make(map[string]interface{})
				regionsChild1Raw := regionsChild1Raw.(map[string]interface{})
				regionsMap["region_id"] = regionsChild1Raw["RegionId"]
				regionsMap["retain_days"] = regionsChild1Raw["RetainDays"]

				regionsMaps = append(regionsMaps, regionsMap)
			}
		}
		crossRegionCopyInfoMap["regions"] = regionsMaps
		crossRegionCopyInfoMaps = append(crossRegionCopyInfoMaps, crossRegionCopyInfoMap)
	}
	d.Set("cross_region_copy_info", crossRegionCopyInfoMaps)
	retainRuleMaps := make([]map[string]interface{}, 0)
	retainRuleMap := make(map[string]interface{})
	retainRule1Raw := make(map[string]interface{})
	if objectRaw["RetainRule"] != nil {
		retainRule1Raw = objectRaw["RetainRule"].(map[string]interface{})
	}
	if len(retainRule1Raw) > 0 {
		retainRuleMap["number"] = retainRule1Raw["Number"]
		retainRuleMap["time_interval"] = retainRule1Raw["TimeInterval"]
		retainRuleMap["time_unit"] = retainRule1Raw["TimeUnit"]

		retainRuleMaps = append(retainRuleMaps, retainRuleMap)
	}
	d.Set("retain_rule", retainRuleMaps)
	scheduleMaps := make([]map[string]interface{}, 0)
	scheduleMap := make(map[string]interface{})
	schedule1Raw := make(map[string]interface{})
	if objectRaw["Schedule"] != nil {
		schedule1Raw = objectRaw["Schedule"].(map[string]interface{})
	}
	if len(schedule1Raw) > 0 {
		scheduleMap["cron_expression"] = schedule1Raw["CronExpression"]

		scheduleMaps = append(scheduleMaps, scheduleMap)
	}
	d.Set("schedule", scheduleMaps)
	specialRetainRulesMaps := make([]map[string]interface{}, 0)
	specialRetainRulesMap := make(map[string]interface{})
	specialRetainRules1Raw := make(map[string]interface{})
	if objectRaw["SpecialRetainRules"] != nil {
		specialRetainRules1Raw = objectRaw["SpecialRetainRules"].(map[string]interface{})
	}
	if len(specialRetainRules1Raw) > 0 {
		specialRetainRulesMap["enabled"] = specialRetainRules1Raw["Enabled"]

		rules1Raw := specialRetainRules1Raw["Rules"]
		rulesMaps := make([]map[string]interface{}, 0)
		if rules1Raw != nil {
			for _, rulesChild1Raw := range rules1Raw.([]interface{}) {
				rulesMap := make(map[string]interface{})
				rulesChild1Raw := rulesChild1Raw.(map[string]interface{})
				rulesMap["special_period_unit"] = rulesChild1Raw["SpecialPeriodUnit"]
				rulesMap["time_interval"] = rulesChild1Raw["TimeInterval"]
				rulesMap["time_unit"] = rulesChild1Raw["TimeUnit"]

				rulesMaps = append(rulesMaps, rulesMap)
			}
		}
		specialRetainRulesMap["rules"] = rulesMaps
		specialRetainRulesMaps = append(specialRetainRulesMaps, specialRetainRulesMap)
	}
	d.Set("special_retain_rules", specialRetainRulesMaps)
	storageRuleMaps := make([]map[string]interface{}, 0)
	storageRuleMap := make(map[string]interface{})
	storageRule1Raw := make(map[string]interface{})
	if objectRaw["StorageRule"] != nil {
		storageRule1Raw = objectRaw["StorageRule"].(map[string]interface{})
	}
	if len(storageRule1Raw) > 0 {
		storageRuleMap["enable_immediate_access"] = storageRule1Raw["EnableImmediateAccess"]

		storageRuleMaps = append(storageRuleMaps, storageRuleMap)
	}
	d.Set("storage_rule", storageRuleMaps)
	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudEbsEnterpriseSnapshotPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	action := "UpdateEnterpriseSnapshotPolicy"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["PolicyId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("enterprise_snapshot_policy_name") {
		update = true
	}
	request["Name"] = d.Get("enterprise_snapshot_policy_name")
	if d.HasChange("desc") {
		update = true
		request["Desc"] = d.Get("desc")
	}

	if d.HasChange("schedule") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})
	if v := d.Get("schedule"); !IsNil(v) {
		nodeNative, _ := jsonpath.Get("$[0].cron_expression", v)
		if nodeNative != "" {
			objectDataLocalMap["CronExpression"] = nodeNative
		}
		request["Schedule"], err = convertMaptoJsonString(objectDataLocalMap)
		if err != nil {
			return WrapError(err)
		}
	}

	if d.HasChange("retain_rule") {
		update = true
	}
	objectDataLocalMap1 := make(map[string]interface{})
	if v := d.Get("retain_rule"); !IsNil(v) {
		nodeNative1, _ := jsonpath.Get("$[0].number", v)
		if nodeNative1 != "" && nodeNative1.(int) > 0 {
			objectDataLocalMap1["Number"] = nodeNative1
		}
		nodeNative2, _ := jsonpath.Get("$[0].time_interval", v)
		if nodeNative2 != "" && nodeNative2.(int) > 0 {
			objectDataLocalMap1["TimeInterval"] = nodeNative2
		}
		nodeNative3, _ := jsonpath.Get("$[0].time_unit", v)
		if nodeNative3 != "" {
			objectDataLocalMap1["TimeUnit"] = nodeNative3
		}
		request["RetainRule"], err = convertMaptoJsonString(objectDataLocalMap1)
		if err != nil {
			return WrapError(err)
		}
	}

	if d.HasChange("storage_rule") {
		update = true
		objectDataLocalMap2 := make(map[string]interface{})
		if v := d.Get("storage_rule"); !IsNil(v) {
			nodeNative4, _ := jsonpath.Get("$[0].enable_immediate_access", v)
			if nodeNative4 != "" {
				objectDataLocalMap2["EnableImmediateAccess"] = nodeNative4
			}
			request["StorageRule"], err = convertMaptoJsonString(objectDataLocalMap2)
			if err != nil {
				return WrapError(err)
			}
		}
	}

	if d.HasChange("cross_region_copy_info") {
		update = true
		objectDataLocalMap3 := make(map[string]interface{})
		if v := d.Get("cross_region_copy_info"); !IsNil(v) {
			nodeNative5, _ := jsonpath.Get("$[0].enabled", v)
			if nodeNative5 != "" {
				objectDataLocalMap3["Enabled"] = nodeNative5
			}
			if v, ok := d.GetOk("cross_region_copy_info"); ok {
				localData, err := jsonpath.Get("$[0].regions", v)
				if err != nil {
					return WrapError(err)
				}
				localMaps := make([]map[string]interface{}, 0)
				for _, dataLoop := range localData.([]interface{}) {
					dataLoopTmp := dataLoop.(map[string]interface{})
					dataLoopMap := make(map[string]interface{})
					dataLoopMap["RegionId"] = dataLoopTmp["region_id"]
					if dataLoopTmp["retain_days"].(int) > 0 {
						dataLoopMap["RetainDays"] = dataLoopTmp["retain_days"]
					}
					localMaps = append(localMaps, dataLoopMap)
				}
				objectDataLocalMap3["Regions"] = localMaps
			}
			request["CrossRegionCopyInfo"], err = convertMaptoJsonString(objectDataLocalMap3)
			if err != nil {
				return WrapError(err)
			}
		}
	}

	if d.HasChange("status") {
		update = true
		request["State"] = d.Get("status")
	}

	if d.HasChange("special_retain_rules") {
		update = true
		objectDataLocalMap4 := make(map[string]interface{})
		if v := d.Get("special_retain_rules"); !IsNil(v) {
			nodeNative8, _ := jsonpath.Get("$[0].enabled", v)
			if nodeNative8 != "" {
				objectDataLocalMap4["Enabled"] = nodeNative8
			}
			if v, ok := d.GetOk("special_retain_rules"); ok {
				localData1, err := jsonpath.Get("$[0].rules", v)
				if err != nil {
					return WrapError(err)
				}
				localMaps1 := make([]map[string]interface{}, 0)
				for _, dataLoop1 := range localData1.([]interface{}) {
					dataLoop1Tmp := dataLoop1.(map[string]interface{})
					dataLoop1Map := make(map[string]interface{})
					dataLoop1Map["SpecialPeriodUnit"] = dataLoop1Tmp["special_period_unit"]
					dataLoop1Map["TimeInterval"] = dataLoop1Tmp["time_interval"]
					dataLoop1Map["TimeUnit"] = dataLoop1Tmp["time_unit"]
					localMaps1 = append(localMaps1, dataLoop1Map)
				}
				objectDataLocalMap4["Rules"] = localMaps1
			}
			request["SpecialRetainRules"], err = convertMaptoJsonString(objectDataLocalMap4)
			if err != nil {
				return WrapError(err)
			}
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
			request["ClientToken"] = buildClientToken(action)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("enterprise_snapshot_policy_name")
		d.SetPartial("desc")
		d.SetPartial("status")
	}
	update = false
	action = "ChangeResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if _, ok := d.GetOk("resource_group_id"); ok && d.HasChange("resource_group_id") {
		update = true
		request["NewResourceGroupId"] = d.Get("resource_group_id")
	}

	request["ResourceType"] = "EnterpriseSnapshotPolicy"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
			request["ClientToken"] = buildClientToken(action)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("resource_group_id")
	}

	if d.HasChange("tags") {
		ebsServiceV2 := EbsServiceV2{client}
		if err := ebsServiceV2.SetResourceTags(d, "EnterpriseSnapshotPolicy"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	d.Partial(false)
	return resourceAliCloudEbsEnterpriseSnapshotPolicyRead(d, meta)
}

func resourceAliCloudEbsEnterpriseSnapshotPolicyDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteEnterpriseSnapshotPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["PolicyId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["ClientToken"] = buildClientToken(action)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
