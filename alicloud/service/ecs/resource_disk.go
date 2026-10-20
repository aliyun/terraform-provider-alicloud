package ecs

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceDiskCreate,
		Read:   resourceDiskRead,
		Update: resourceDiskUpdate,
		Delete: resourceDiskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"advanced_features": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bursting_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"category": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"cloud", "cloud_efficiency", "cloud_essd", "cloud_ssd", "cloud_auto", "cloud_essd_entry", "elastic_ephemeral_disk_standard", "elastic_ephemeral_disk_premium"}, false),
				Default:      "cloud_efficiency",
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delete_auto_snapshot": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"delete_with_instance": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disk_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_auto_snapshot": {
				Type:       schema.TypeBool,
				Optional:   true,
				Computed:   true,
				Deprecated: "This field is deprecated. The automatic snapshot policy feature is enabled by default for a cloud disk after it is created. To use the automatic snapshot policy, apply one to the cloud disk.",
			},
			"encrypt_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"multi_attach": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"performance_level": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("category").(string) != "cloud_essd"
				},
			},
			"provisioned_iops": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
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
			"storage_set_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"storage_set_partition_number": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"tags": tags.TagsSchema(),
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"offline", "online"}, false),
			},
			"zone_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"availability_zone"},
			},
			"availability_zone": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "Field `availability_zone` has been deprecated from provider version 1.122.0. New field `zone_id` instead",
				ConflictsWith: []string{"zone_id"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field `name` has been deprecated from provider version 1.122.0. New field `disk_name` instead.",
				ConflictsWith: []string{"disk_name"},
			},
		},
	}
}

func resourceDiskCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDisk"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := tags.ConvertTags(v.(map[string]interface{}))
		request = tags.ExpandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("encrypt_algorithm"); ok {
		request["EncryptAlgorithm"] = v
	}
	if v, ok := d.GetOkExists("encrypted"); ok {
		request["Encrypted"] = v
	}
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	if v, ok := d.GetOk("category"); ok {
		request["DiskCategory"] = v
	}
	if v, ok := d.GetOkExists("size"); ok {
		request["Size"] = v
	}
	if v, ok := d.GetOk("disk_name"); ok {
		request["DiskName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["DiskName"] = v
	}
	if v, ok := d.GetOk("performance_level"); ok {
		request["PerformanceLevel"] = v
	}
	if v, ok := d.GetOk("snapshot_id"); ok {
		request["SnapshotId"] = v
	}
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	} else if v, ok := d.GetOk("availability_zone"); ok {
		request["ZoneId"] = v
	}
	if v, ok := d.GetOk("kms_key_id"); ok {
		request["KMSKeyId"] = v
	}
	if v, ok := d.GetOk("multi_attach"); ok {
		request["MultiAttach"] = v
	}
	if v, ok := d.GetOkExists("provisioned_iops"); ok {
		request["ProvisionedIops"] = v
	}
	if v, ok := d.GetOkExists("bursting_enabled"); ok {
		request["BurstingEnabled"] = v
	}
	if v, ok := d.GetOk("storage_set_id"); ok {
		request["StorageSetId"] = v
	}
	if v, ok := d.GetOkExists("storage_set_partition_number"); ok {
		request["StorageSetPartitionNumber"] = v
	}
	if v, ok := d.GetOk("advanced_features"); ok {
		request["AdvancedFeatures"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"ServiceUnavailable", "UnknownError", "LastTokenProcessing"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ecs_disk", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DiskId"]))

	// Setting instanceId aims to creating the PrePaid disk and the job is async
	if fmt.Sprint(request["InstanceId"]) != "" {
		service := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{"Available", "In_use"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, service.DiskStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	return resourceDiskUpdate(d, meta)
}

func resourceDiskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	service := NewService(client)

	objectRaw, err := service.DescribeDisk(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_disk DescribeDisk Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if objectRaw["BurstingEnabled"] != nil {
		d.Set("bursting_enabled", objectRaw["BurstingEnabled"])
	}
	if objectRaw["Category"] != nil {
		d.Set("category", objectRaw["Category"])
	}
	if objectRaw["CreationTime"] != nil {
		d.Set("create_time", objectRaw["CreationTime"])
	}
	if objectRaw["DeleteAutoSnapshot"] != nil {
		d.Set("delete_auto_snapshot", objectRaw["DeleteAutoSnapshot"])
	}
	if objectRaw["DeleteWithInstance"] != nil {
		d.Set("delete_with_instance", objectRaw["DeleteWithInstance"])
	}
	if objectRaw["Description"] != nil {
		d.Set("description", objectRaw["Description"])
	}
	if objectRaw["DiskName"] != nil {
		d.Set("disk_name", objectRaw["DiskName"])
		d.Set("name", objectRaw["DiskName"])
	}
	if objectRaw["EnableAutoSnapshot"] != nil {
		d.Set("enable_auto_snapshot", objectRaw["EnableAutoSnapshot"])
	}
	if objectRaw["Encrypted"] != nil {
		d.Set("encrypted", objectRaw["Encrypted"])
	}
	if objectRaw["InstanceId"] != nil {
		d.Set("instance_id", objectRaw["InstanceId"])
	}
	if objectRaw["KMSKeyId"] != nil {
		d.Set("kms_key_id", objectRaw["KMSKeyId"])
	}
	if objectRaw["MultiAttach"] != nil {
		d.Set("multi_attach", objectRaw["MultiAttach"])
	}
	if objectRaw["DiskChargeType"] != nil {
		d.Set("payment_type", convertDiskPaymentTypeResponse(objectRaw["DiskChargeType"]))
	}
	if objectRaw["PerformanceLevel"] != nil {
		d.Set("performance_level", objectRaw["PerformanceLevel"])
	}
	if objectRaw["ProvisionedIops"] != nil {
		d.Set("provisioned_iops", objectRaw["ProvisionedIops"])
	}
	if objectRaw["RegionId"] != nil {
		d.Set("region_id", objectRaw["RegionId"])
	}
	if objectRaw["ResourceGroupId"] != nil {
		d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	}
	if objectRaw["Size"] != nil {
		d.Set("size", objectRaw["Size"])
	}
	if objectRaw["SourceSnapshotId"] != nil {
		d.Set("snapshot_id", objectRaw["SourceSnapshotId"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}
	if objectRaw["StorageSetId"] != nil {
		d.Set("storage_set_id", objectRaw["StorageSetId"])
	}
	if objectRaw["StorageSetPartitionNumber"] != nil {
		d.Set("storage_set_partition_number", objectRaw["StorageSetPartitionNumber"])
	}
	if objectRaw["ZoneId"] != nil {
		d.Set("zone_id", objectRaw["ZoneId"])
		d.Set("availability_zone", objectRaw["ZoneId"])
	}

	tagsMaps, _ := jsonpath.Get("$.Tags.Tag", objectRaw)
	d.Set("tags", tags.TagsToMap(tagsMaps))

	return nil
}

func resourceDiskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	service := NewService(client)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	// needBurstingPost defers the BurstingEnabled API call to after ModifyDiskSpec +
	// WaitForState when category is being changed to cloud_auto in the same apply.
	// BurstingEnabled is only valid on cloud_auto disks (OpenAPI returns an error if
	// specified on a disk that does not support performance burst).
	needBurstingPost := false
	d.Partial(true)

	action := "ModifyDiskAttribute"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DiskId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("delete_auto_snapshot") {
		update = true
		request["DeleteAutoSnapshot"] = d.Get("delete_auto_snapshot")
	}

	objectRaw, err := service.DescribeDisk(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}
	if fmt.Sprint(objectRaw["DeleteWithInstance"]) != fmt.Sprint(d.Get("delete_with_instance")) {
		update = true
		request["DeleteWithInstance"] = d.Get("delete_with_instance")
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("disk_name") {
		update = true
		request["DiskName"] = d.Get("disk_name")
	}

	if !d.IsNewResource() && d.HasChange("name") {
		update = true
		request["DiskName"] = d.Get("name")
	}

	enableAutoSnapshot, ok := d.GetOkExists("enable_auto_snapshot")
	if ok && fmt.Sprint(objectRaw["EnableAutoSnapshot"]) != fmt.Sprint(enableAutoSnapshot) {
		update = true
		request["EnableAutoSnapshot"] = enableAutoSnapshot
	}

	if !d.IsNewResource() && d.HasChange("bursting_enabled") {
		// When category is being changed to cloud_auto in the same apply, the disk is
		// still on its previous category (e.g. cloud_essd) at this point, so setting
		// BurstingEnabled here would fail. Defer it to after ModifyDiskSpec +
		// WaitForState so the disk is already cloud_auto. When category is unchanged
		// or transitions away from cloud_auto, setting it here is safe because the
		// disk is still cloud_auto at this point when disabling the burst feature.
		if d.HasChange("category") && d.Get("category").(string) == "cloud_auto" {
			needBurstingPost = true
		} else {
			update = true
			request["BurstingEnabled"] = d.Get("bursting_enabled")
		}
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	update = false
	action = "ResizeDisk"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DiskId"] = d.Id()

	request["ClientToken"] = helper.BuildClientToken(action)
	if !d.IsNewResource() && d.HasChange("size") {
		update = true
	}
	request["NewSize"] = d.Get("size")
	if v, ok := d.GetOk("type"); ok {
		request["Type"] = v
	}
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"LastTokenProcessing", "InternalError"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "JoinResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	request["ResourceType"] = "disk"
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"InternalError"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "ModifyDiskSpec"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DiskId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("performance_level") {
		update = true
		request["PerformanceLevel"] = d.Get("performance_level")
	}

	if !d.IsNewResource() && d.HasChange("category") {
		update = true
		request["DiskCategory"] = d.Get("category")
	}

	if !d.IsNewResource() && d.HasChange("provisioned_iops") {
		// ProvisionedIops is only valid on cloud_auto disks. When category is
		// transitioning to a non-cloud_auto value, ModifyDiskSpec rejects
		// ProvisionedIops (any value, including 0) with
		// ProvisionedIopsForDiskCategoryUnsupported; the disk's ProvisionedIops
		// is implicitly cleared by the category change, so skip the explicit
		// set in that case. When category is unchanged or transitions to
		// cloud_auto, the target category supports ProvisionedIops and the
		// direct set is safe.
		if !d.HasChange("category") || d.Get("category").(string) == "cloud_auto" {
			update = true
			request["ProvisionedIops"] = d.Get("provisioned_iops")
		}
	}

	if update {
		if _, ok := d.GetOkExists("dry_run"); ok {
			request["DryRun"] = d.Get("dry_run")
		}
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"ServiceUnavailable", "Throttling.ConcurrentLimitExceeded", "DiskInCoolingPeriod"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		// Category conversion (e.g. cloud_essd -> cloud_auto) is applied asynchronously
		// by ECS: after ModifyDiskSpec returns, DescribeDisks may report Status
		// "Modifying" during the transition, and even after Status returns to
		// "Available" the Category field can still report the previous value. Two
		// separate waits cannot exclude such interleavings (e.g. a response reading
		// Status=Available with the old category followed by one reading
		// Modifying with the new category), so on the category-change path both
		// fields are checked within one DescribeDisks response: the wait completes
		// only when Status is Available/In_use AND Category equals the requested
		// value, so the deferred BurstingEnabled update below, which is only valid
		// on cloud_auto disks, is not issued while the disk is still converting
		// (BurstingEnabledForModifyingDiskUnsupported). Conversion duration depends
		// on the disk (throughput, capacity, previous category) and comes with no
		// completion guarantee, so the wait uses the configured update timeout as-is
		// instead of a multiplied budget; set `timeouts.update` explicitly for large
		// disks whose conversion may take longer.
		if d.HasChange("category") {
			targetCategory := d.Get("category").(string)
			categoryConf := helper.BuildStateConf([]string{"Modifying", "CategoryConverting"}, []string{"Available", "In_use"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, func() (interface{}, string, error) {
				objectRaw, err := service.DescribeDisk(d.Id())
				return diskCategoryTransitionRefresh(objectRaw, err, targetCategory)
			})
			if _, err := categoryConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}
		} else {
			// With an empty Pending list the state machine keeps polling through any
			// unknown state (including "Modifying" on non-category updates) until the
			// target or the timeout.
			stateConf := helper.BuildStateConf([]string{}, []string{"Available", "In_use"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, service.DiskStateRefreshFunc(d.Id(), "Status", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}
		}
	}
	// Deferred BurstingEnabled set: category has just been changed to cloud_auto and
	// the Status + Category WaitForState above confirmed the disk reached the target
	// category, so BurstingEnabled is valid.
	if needBurstingPost {
		update = true
		action = "ModifyDiskAttribute"
		request = make(map[string]interface{})
		query = make(map[string]interface{})
		request["DiskId"] = d.Id()
		request["RegionId"] = client.RegionId
		request["BurstingEnabled"] = d.Get("bursting_enabled")
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "ModifyDiskChargeType"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DiskIds"] = convert.ConvertListToJsonString([]interface{}{d.Id()})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)
	request["AutoPay"] = true
	request["InstanceId"] = d.Get("instance_id")
	if !d.IsNewResource() && d.HasChange("payment_type") {
		update = true
		request["DiskChargeType"] = convertDiskPaymentTypeRequest(d.Get("payment_type").(string))
	}
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		stateConf := helper.BuildStateConf([]string{}, []string{fmt.Sprint(request["DiskChargeType"])}, d.Timeout(schema.TimeoutUpdate), 1*time.Second, service.DiskStateRefreshFunc(d.Id(), "DiskChargeType", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	if !d.IsNewResource() && d.HasChange("tags") {
		if err := service.SetResourceTags(d, "disk"); err != nil {
			return errs.WrapError(err)
		}
	}
	d.Partial(false)
	return resourceDiskRead(d, meta)
}

func resourceDiskDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	service := NewService(client)

	objectRaw, err := service.DescribeDisk(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_disk DescribeDisk Failed!!! %s", err)
			return nil
		}
		return errs.WrapError(err)
	}

	if fmt.Sprint(objectRaw["DiskChargeType"]) == "PrePaid" && fmt.Sprint(objectRaw["DeleteWithInstance"]) == "true" {
		log.Printf("[DEBUG] Resource alicloud_ecs_disk %s charge type is PrePaid and its attribute DeleteWithInstance is true, so it will remove from state.", d.Id())
		return nil
	}

	action := "DeleteDisk"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	request = make(map[string]interface{})
	request["DiskId"] = d.Id()

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"IncorrectDiskStatus.Initializing"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDiskId.NotFound"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}

func convertDiskPaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "PostPaid":
		return "PayAsYouGo"
	case "PrePaid":
		return "Subscription"
	}
	return source
}

func convertDiskPaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PostPaid"
	case "Subscription":
		return "PrePaid"
	}
	return source
}

// diskCategoryTransitionState maps one DescribeDisks response into the state of
// the combined Status/Category wait on the category-change path. The wait targets
// Status Available/In_use and only completes when Category already equals the
// requested value in the same response; "Modifying" and "CategoryConverting" keep
// polling, and any other status is returned as-is so the state machine fails fast
// instead of polling until the timeout.
func diskCategoryTransitionState(status, category, targetCategory string) string {
	switch {
	case status == "Modifying":
		return "Modifying"
	case status == "Available" || status == "In_use":
		if category == targetCategory {
			return status
		}
		return "CategoryConverting"
	default:
		return status
	}
}

// diskCategoryTransitionRefresh converts one DescribeDisk result into a refresh
// result for the combined Status/Category wait. A disk that DescribeDisk reports
// as NotFound must surface as a nil result: StateChangeConf only treats a nil
// result as "not found, keep polling", while a typed nil map boxed in interface{}
// is non-nil and its empty state would match neither Pending nor Target, failing
// the wait with an unexpected state.
func diskCategoryTransitionRefresh(objectRaw map[string]interface{}, err error, targetCategory string) (interface{}, string, error) {
	if err != nil {
		if errs.NotFoundError(err) {
			return nil, "", nil
		}
		return nil, "", errs.WrapError(err)
	}
	return objectRaw, diskCategoryTransitionState(fmt.Sprint(objectRaw["Status"]), fmt.Sprint(objectRaw["Category"]), targetCategory), nil
}
