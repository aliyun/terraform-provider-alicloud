// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"time"
)

func resourceAliCloudEcsAutoSnapshotPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEcsAutoSnapshotPolicyCreate,
		Read:   resourceAliCloudEcsAutoSnapshotPolicyRead,
		Update: resourceAliCloudEcsAutoSnapshotPolicyUpdate,
		Delete: resourceAliCloudEcsAutoSnapshotPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_snapshot_policy_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
			},
			"copied_snapshots_retention_days": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"copy_encryption_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kms_key_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"encrypted": {
							Type:     schema.TypeBool,
							Optional: true,
							//Computed: true,
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_cross_region_copy": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"repeat_weekdays": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				//Computed: true,
			},
			"retention_days": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"target_copy_regions": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"time_points": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field `name` has been deprecated from provider version 1.236.0. New field `auto_snapshot_policy_name` instead.",
			},
		},
	}
}

func resourceAliCloudEcsAutoSnapshotPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAutoSnapshotPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["regionId"] = client.RegionId

	if v, ok := d.GetOk("name"); ok || d.HasChange("name") {
		request["autoSnapshotPolicyName"] = v
	}

	if v, ok := d.GetOk("auto_snapshot_policy_name"); ok {
		request["autoSnapshotPolicyName"] = v
	}
	request["retentionDays"] = d.Get("retention_days")
	if v, ok := d.GetOkExists("enable_cross_region_copy"); ok {
		request["EnableCrossRegionCopy"] = v
	}
	if v, ok := d.GetOkExists("copied_snapshots_retention_days"); ok {
		request["CopiedSnapshotsRetentionDays"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("copy_encryption_configuration"); ok {
		jsonPathResult5, err := jsonpath.Get("$[0].encrypted", v)
		if err == nil && jsonPathResult5 != "" {
			request["CopyEncryptionConfiguration.Encrypted"] = jsonPathResult5
		}
	}
	if v, ok := d.GetOk("copy_encryption_configuration"); ok {
		jsonPathResult6, err := jsonpath.Get("$[0].kms_key_id", v)
		if err == nil && jsonPathResult6 != "" {
			request["CopyEncryptionConfiguration.KMSKeyId"] = jsonPathResult6
		}
	}
	jsonPathResult7, err := jsonpath.Get("$", d.Get("time_points"))
	if err == nil {
		request["timePoints"] = convertListToJsonString(jsonPathResult7.([]interface{}))
	}

	jsonPathResult8, err := jsonpath.Get("$", d.Get("repeat_weekdays"))
	if err == nil {
		request["repeatWeekdays"] = convertListToJsonString(jsonPathResult8.([]interface{}))
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("target_copy_regions"); ok {
		jsonPathResult9, err := jsonpath.Get("$", v)
		if err == nil && jsonPathResult9 != "" {
			request["TargetCopyRegions"] = convertListToJsonString(jsonPathResult9.(*schema.Set).List())
		}
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_auto_snapshot_policy", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AutoSnapshotPolicyId"]))

	return resourceAliCloudEcsAutoSnapshotPolicyRead(d, meta)
}

func resourceAliCloudEcsAutoSnapshotPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsServiceV2 := EcsServiceV2{client}

	objectRaw, err := ecsServiceV2.DescribeEcsAutoSnapshotPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_auto_snapshot_policy DescribeEcsAutoSnapshotPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["AutoSnapshotPolicyName"] != nil {
		d.Set("auto_snapshot_policy_name", objectRaw["AutoSnapshotPolicyName"])
	}
	if objectRaw["CopiedSnapshotsRetentionDays"] != nil {
		d.Set("copied_snapshots_retention_days", objectRaw["CopiedSnapshotsRetentionDays"])
	}
	if objectRaw["CreationTime"] != nil {
		d.Set("create_time", objectRaw["CreationTime"])
	}
	if objectRaw["EnableCrossRegionCopy"] != nil {
		d.Set("enable_cross_region_copy", objectRaw["EnableCrossRegionCopy"])
	}
	if objectRaw["RegionId"] != nil {
		d.Set("region_id", objectRaw["RegionId"])
	}
	if objectRaw["ResourceGroupId"] != nil {
		d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	}
	if objectRaw["RetentionDays"] != nil {
		d.Set("retention_days", objectRaw["RetentionDays"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}

	copyEncryptionConfigurationMaps := make([]map[string]interface{}, 0)
	copyEncryptionConfigurationMap := make(map[string]interface{})
	copyEncryptionConfiguration1Raw := make(map[string]interface{})
	if objectRaw["CopyEncryptionConfiguration"] != nil {
		copyEncryptionConfiguration1Raw = objectRaw["CopyEncryptionConfiguration"].(map[string]interface{})
	}
	if len(copyEncryptionConfiguration1Raw) > 0 {
		if copyEncryptionConfiguration1Raw["Encrypted"] != nil {
			copyEncryptionConfigurationMap["encrypted"] = copyEncryptionConfiguration1Raw["Encrypted"]
		}

		if copyEncryptionConfiguration1Raw["KMSKeyId"] != nil {
			copyEncryptionConfigurationMap["kms_key_id"] = copyEncryptionConfiguration1Raw["KMSKeyId"]
		}

		copyEncryptionConfigurationMaps = append(copyEncryptionConfigurationMaps, copyEncryptionConfigurationMap)
	}

	if objectRaw["CopyEncryptionConfiguration"] != nil {
		if err := d.Set("copy_encryption_configuration", copyEncryptionConfigurationMaps); err != nil {
			return err
		}
	}

	tagsMaps, _ := jsonpath.Get("$.Tags.Tag", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	if objectRaw["RepeatWeekdays"] != nil {
		if repeatWeekdays, err := convertJsonStringToList(objectRaw["RepeatWeekdays"].(string)); err != nil {
			return WrapError(err)
		} else {
			d.Set("repeat_weekdays", repeatWeekdays)
		}
	}

	if objectRaw["TargetCopyRegions"] != nil {
		if targetCopyRegions, err := convertJsonStringToList(objectRaw["TargetCopyRegions"].(string)); err != nil {
			return WrapError(err)
		} else {
			d.Set("target_copy_regions", targetCopyRegions)
		}
	}

	if objectRaw["TimePoints"] != nil {
		if timePoints, err := convertJsonStringToList(objectRaw["TimePoints"].(string)); err != nil {
			return WrapError(err)
		} else {
			d.Set("time_points", timePoints)
		}
	}

	d.Set("name", d.Get("auto_snapshot_policy_name"))
	return nil
}

func resourceAliCloudEcsAutoSnapshotPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	action := "ModifyAutoSnapshotPolicyEx"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["autoSnapshotPolicyId"] = d.Id()
	request["regionId"] = client.RegionId
	if d.HasChange("retention_days") {
		update = true
	}
	request["retentionDays"] = d.Get("retention_days")
	if d.HasChange("repeat_weekdays") {
		update = true
	}
	jsonPathResult1, err := jsonpath.Get("$", d.Get("repeat_weekdays"))
	if err == nil {
		request["repeatWeekdays"] = convertListToJsonString(jsonPathResult1.([]interface{}))
	}

	if d.HasChange("time_points") {
		update = true
	}
	jsonPathResult2, err := jsonpath.Get("$", d.Get("time_points"))
	if err == nil {
		request["timePoints"] = convertListToJsonString(jsonPathResult2.([]interface{}))
	}

	if d.HasChange("copied_snapshots_retention_days") {
		update = true
		request["CopiedSnapshotsRetentionDays"] = d.Get("copied_snapshots_retention_days")
	}

	if d.HasChange("enable_cross_region_copy") {
		update = true
	}
	if v, ok := d.GetOkExists("enable_cross_region_copy"); ok {
		request["EnableCrossRegionCopy"] = v
	}

	if d.HasChange("target_copy_regions") {
		update = true
	}
	if v, ok := d.GetOk("target_copy_regions"); ok {
		jsonPathResult9, err := jsonpath.Get("$", v)
		if err == nil && jsonPathResult9 != "" {
			request["TargetCopyRegions"] = convertListToJsonString(jsonPathResult9.(*schema.Set).List())
		}
	}

	if d.HasChange("name") {
		update = true
		request["autoSnapshotPolicyName"] = d.Get("name")
	}

	if d.HasChange("auto_snapshot_policy_name") {
		update = true
		request["autoSnapshotPolicyName"] = d.Get("auto_snapshot_policy_name")
	}

	if d.HasChange("copy_encryption_configuration.0.encrypted") {
		update = true
	}
	if v, ok := d.GetOk("copy_encryption_configuration"); ok {
		jsonPathResult5, err := jsonpath.Get("$[0].encrypted", v)
		if err == nil && jsonPathResult5 != "" {
			request["CopyEncryptionConfiguration.Encrypted"] = jsonPathResult5
		}
	}

	if d.HasChange("copy_encryption_configuration.0.kms_key_id") {
		update = true
	}
	if v, ok := d.GetOk("copy_encryption_configuration"); ok {
		jsonPathResult6, err := jsonpath.Get("$[0].kms_key_id", v)
		if err == nil && jsonPathResult6 != "" {
			request["CopyEncryptionConfiguration.KMSKeyId"] = jsonPathResult6
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
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
	update = false
	action = "JoinResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ResourceType"] = "snapshotpolicy"
	if _, ok := d.GetOk("resource_group_id"); ok && d.HasChange("resource_group_id") {
		update = true
		request["ResourceGroupId"] = d.Get("resource_group_id")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
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

	if d.HasChange("tags") {
		ecsServiceV2 := EcsServiceV2{client}
		if err := ecsServiceV2.SetResourceTags(d, "snapshotpolicy"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudEcsAutoSnapshotPolicyRead(d, meta)
}

func resourceAliCloudEcsAutoSnapshotPolicyDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAutoSnapshotPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["autoSnapshotPolicyId"] = d.Id()
	request["regionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)

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
