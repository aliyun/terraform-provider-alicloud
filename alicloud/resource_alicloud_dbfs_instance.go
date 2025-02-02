// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudDbfsDbfsInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDbfsDbfsInstanceCreate,
		Read:   resourceAliCloudDbfsDbfsInstanceRead,
		Update: resourceAliCloudDbfsDbfsInstanceUpdate,
		Delete: resourceAliCloudDbfsDbfsInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"advanced_features": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"category": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"standard", "enterprise"}, false),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delete_snapshot": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_raid": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"encryption": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"fs_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"fs_name", "instance_name"},
			},
			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"dbfs.small", "dbfs.medium", "dbfs.large "}, false),
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"performance_level": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
			},
			"raid_stripe_unit_number": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: IntBetween(2, 8),
			},
			"size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: IntBetween(20, 262144),
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"used_scene": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"MySQL 5.7", "PostgreSQL ", "MongoDB", "DataCube"}, false),
			},
			"tags": tagsSchema(),
			"ecs_list": {
				Type:       schema.TypeSet,
				Optional:   true,
				Deprecated: "Field 'ecs_list' has been deprecated from provider version 1.156.0 and it will be removed in the future version. Please use the new resource 'alicloud_dbfs_instance_attachment' to attach ECS and DBFS.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ecs_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'instance_name' has been deprecated since provider version 1.212.0. New field 'fs_name' instead.",
			},
		},
	}
}

func resourceAliCloudDbfsDbfsInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDbfs"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("instance_name"); ok {
		request["FsName"] = v
	}

	if v, ok := d.GetOk("fs_name"); ok {
		request["FsName"] = v
	}
	request["Category"] = d.Get("category")
	request["ZoneId"] = d.Get("zone_id")
	request["SizeG"] = d.Get("size")
	if v, ok := d.GetOk("snapshot_id"); ok {
		request["SnapshotId"] = v
	}
	if v, ok := d.GetOkExists("delete_snapshot"); ok {
		request["DeleteSnapshot"] = v
	}
	if v, ok := d.GetOk("performance_level"); ok {
		request["PerformanceLevel"] = v
	}
	if v, ok := d.GetOkExists("enable_raid"); ok {
		request["EnableRaid"] = v
	}
	if v, ok := d.GetOk("raid_stripe_unit_number"); ok {
		request["RaidStripeUnitNumber"] = v
	}
	if v, ok := d.GetOk("kms_key_id"); ok {
		request["KMSKeyId"] = v
	}
	if v, ok := d.GetOkExists("encryption"); ok {
		request["Encryption"] = v
	}
	if v, ok := d.GetOk("used_scene"); ok {
		request["UsedScene"] = v
	}
	if v, ok := d.GetOk("instance_type"); ok {
		request["InstanceType"] = v
	}
	if v, ok := d.GetOk("advanced_features"); ok {
		request["AdvancedFeatures"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("DBFS", "2020-04-18", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dbfs_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["FsId"]))

	dbfsServiceV2 := DbfsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"unattached"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, dbfsServiceV2.DbfsDbfsInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudDbfsDbfsInstanceRead(d, meta)
}

func resourceAliCloudDbfsDbfsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dbfsServiceV2 := DbfsServiceV2{client}

	objectRaw, err := dbfsServiceV2.DescribeDbfsDbfsInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dbfs_instance DescribeDbfsDbfsInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("advanced_features", objectRaw["AdvancedFeatures"])
	d.Set("category", objectRaw["Category"])
	d.Set("create_time", objectRaw["CreatedTime"])
	d.Set("enable_raid", objectRaw["EnableRaid"])
	d.Set("encryption", objectRaw["Encryption"])
	d.Set("fs_name", objectRaw["FsName"])
	d.Set("instance_type", objectRaw["InstanceType"])
	d.Set("kms_key_id", objectRaw["KMSKeyId"])
	d.Set("performance_level", objectRaw["PerformanceLevel"])
	d.Set("raid_stripe_unit_number", objectRaw["RaidStrip"])
	d.Set("size", objectRaw["SizeG"])
	d.Set("snapshot_id", objectRaw["SnapshotId"])
	d.Set("status", objectRaw["Status"])
	d.Set("used_scene", objectRaw["UsedScene"])
	d.Set("zone_id", objectRaw["ZoneId"])
	d.Set("tags", tagsToMap(objectRaw["Tags"]))

	d.Set("instance_name", d.Get("fs_name"))
	if ecsListList, ok := objectRaw["EcsList"]; ok && ecsListList != nil {
		ecsListMaps := make([]map[string]interface{}, 0)
		for _, ecsListListItem := range ecsListList.([]interface{}) {
			if ecsListListItemMap, ok := ecsListListItem.(map[string]interface{}); ok {
				ecsListListItemMap["ecs_id"] = ecsListListItemMap["EcsId"]
				ecsListMaps = append(ecsListMaps, ecsListListItemMap)
			}
			d.Set("ecs_list", ecsListMaps)
		}
	}
	return nil
}

func resourceAliCloudDbfsDbfsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	var query map[string]interface{}
	update := false
	d.Partial(true)
	action := "RenameDbfs"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["FsId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("instance_name") {
		update = true
		request["FsName"] = d.Get("instance_name")
	}

	if d.HasChange("fs_name") {
		update = true
		request["FsName"] = d.Get("fs_name")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("DBFS", "2020-04-18", action, query, request, false)

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
		d.SetPartial("fs_name")
	}
	update = false
	action = "ResizeDbfs"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["FsId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("size") {
		update = true
	}
	request["NewSizeG"] = d.Get("size")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("DBFS", "2020-04-18", action, query, request, false)

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
		dbfsServiceV2 := DbfsServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("size"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, dbfsServiceV2.DbfsDbfsInstanceStateRefreshFunc(d.Id(), "SizeG", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("size")
	}
	update = false
	action = "UpdateDbfs"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["FsId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("used_scene") {
		update = true
		request["UsedScene"] = d.Get("used_scene")
	}

	if d.HasChange("instance_type") {
		update = true
		request["InstanceType"] = d.Get("instance_type")
	}

	if d.HasChange("advanced_features") {
		update = true
		request["AdvancedFeatures"] = d.Get("advanced_features")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("DBFS", "2020-04-18", action, query, request, false)

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
		d.SetPartial("used_scene")
		d.SetPartial("instance_type")
		d.SetPartial("advanced_features")
	}
	update = false
	action = "ModifyPerformanceLevel"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["FsId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("performance_level") {
		update = true
		request["PerformanceLevel"] = d.Get("performance_level")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("DBFS", "2020-04-18", action, query, request, false)

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
		d.SetPartial("performance_level")
	}

	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		remove := oraw.(map[string]interface{})
		create := nraw.(map[string]interface{})

		if len(remove) > 0 {

			deleteTagsBatchReq := map[string]interface{}{
				"DbfsList": "[\"" + d.Id() + "\"]",
				"RegionId": client.RegionId,
			}

			tagsMaps := make([]map[string]interface{}, 0)
			for key, value := range remove {
				tagsMap := map[string]interface{}{}
				tagsMap["TagKey"] = key
				tagsMap["TagValue"] = value
				tagsMaps = append(tagsMaps, tagsMap)
			}
			deleteTagsBatchReq["Tags"], _ = convertListMapToJsonString(tagsMaps)

			action := "DeleteTagsBatch"
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("DBFS", "2020-04-18", action, nil, deleteTagsBatchReq, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, deleteTagsBatchReq)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}

		if len(create) > 0 {

			addTagsBatchReq := map[string]interface{}{
				"DbfsList": "[\"" + d.Id() + "\"]",
				"RegionId": client.RegionId,
			}

			tagsMaps := make([]map[string]interface{}, 0)
			for key, value := range create {
				tagsMap := map[string]interface{}{}
				tagsMap["TagKey"] = key
				tagsMap["TagValue"] = value
				tagsMaps = append(tagsMaps, tagsMap)
			}
			addTagsBatchReq["Tags"], _ = convertListMapToJsonString(tagsMaps)

			action := "AddTagsBatch"
			request["ClientToken"] = buildClientToken("AddTagsBatch")
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("DBFS", "2020-04-18", action, nil, addTagsBatchReq, true)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, addTagsBatchReq)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}

		d.SetPartial("tags")
	}

	if d.HasChange("ecs_list") {
		oldEcsList, newEcsList := d.GetChange("ecs_list")
		oldEcsListSet := oldEcsList.(*schema.Set)
		newEcsListSet := newEcsList.(*schema.Set)
		removed := oldEcsListSet.Difference(newEcsListSet)
		added := newEcsListSet.Difference(oldEcsListSet)

		if removed.Len() > 0 {
			detachdbfsrequest := map[string]interface{}{
				"FsId": d.Id(),
			}
			detachdbfsrequest["RegionId"] = client.RegionId
			detachdbfsrequest["ECSInstanceId"] = d.Get("ecs_instance_id")
			for _, ecsArg := range removed.List() {
				ecsMap := ecsArg.(map[string]interface{})
				detachdbfsrequest["ECSInstanceId"] = ecsMap["ecs_id"]

				action := "DetachDbfs"
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("DBFS", "2020-04-18", action, nil, detachdbfsrequest, false)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, detachdbfsrequest)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}

				dbfsService := DbfsService{client}
				stateConf := BuildStateConf([]string{}, []string{"unattached", "attached"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, dbfsService.DbfsInstanceStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
		}

		if added.Len() > 0 {
			attachdbfsrequest := map[string]interface{}{
				"FsId": d.Id(),
			}
			attachdbfsrequest["RegionId"] = client.RegionId

			action := "AttachDbfs"
			for _, ecsArg := range added.List() {
				ecsMap := ecsArg.(map[string]interface{})
				attachdbfsrequest["ECSInstanceId"] = ecsMap["ecs_id"]
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("DBFS", "2020-04-18", action, nil, attachdbfsrequest, false)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, attachdbfsrequest)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				dbfsService := DbfsService{client}
				stateConf := BuildStateConf([]string{}, []string{"attached"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, dbfsService.DbfsInstanceStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
		}

		d.SetPartial("ecs_list")
	}
	d.Partial(false)
	return resourceAliCloudDbfsDbfsInstanceRead(d, meta)
}

func resourceAliCloudDbfsDbfsInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDbfs"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["FsId"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("DBFS", "2020-04-18", action, query, request, false)

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
		if IsExpectedErrors(err, []string{"EntityNotExist.DBFS"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	dbfsServiceV2 := DbfsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 1*time.Minute, dbfsServiceV2.DbfsDbfsInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func convertDbfsInstancePaymentTypeResponse(source string) string {
	switch source {
	case "postpaid":
		return "PayAsYouGo"
	}
	return source
}

func convertDbfsDBFSInfoPayTypeResponse(source interface{}) interface{} {
	switch source {
	case "postpaid":
		return "PayAsYouGo"
	}
	return source
}
