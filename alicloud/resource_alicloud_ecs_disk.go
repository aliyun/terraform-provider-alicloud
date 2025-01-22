// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEcsDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEcsDiskCreate,
		Read:   resourceAliCloudEcsDiskRead,
		Update: resourceAliCloudEcsDiskUpdate,
		Delete: resourceAliCloudEcsDiskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(15 * time.Minute),
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
				ValidateFunc: StringInSlice([]string{"cloud", "cloud_efficiency", "cloud_essd", "cloud_ssd", "cloud_auto", "cloud_essd_entry", "elastic_ephemeral_disk_standard", "elastic_ephemeral_disk_premium"}, false),
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
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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
				ValidateFunc: StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"performance_level": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
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
			"tags": tagsSchema(),
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"offline", "online"}, false),
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
			"dedicated_block_storage_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field `dedicated_block_storage_cluster_id` is unavailable and it has been removed since 1.208.0.",
			},
		},
	}
}

func resourceAliCloudEcsDiskCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDisk"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
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

	if v, ok := d.GetOk("dedicated_block_storage_cluster_id"); ok {
		request["DedicatedBlockStorageClusterId"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), query, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"ServiceUnavailable", "UnknownError", "LastTokenProcessing"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_disk", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DiskId"]))

	// Setting instanceId aims to creating the PrePaid disk and the job is async
	if fmt.Sprint(request["InstanceId"]) != "" {
		ecsServiceV2 := EcsServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available", "In_use"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, ecsServiceV2.EcsDiskStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudEcsDiskUpdate(d, meta)
}

func resourceAliCloudEcsDiskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsServiceV2 := EcsServiceV2{client}

	objectRaw, err := ecsServiceV2.DescribeEcsDisk(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_disk DescribeEcsDisk Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
		d.Set("payment_type", convertEcsDiskPaymentTypeResponse(objectRaw["DiskChargeType"]))
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
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudEcsDiskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsServiceV2 := EcsServiceV2{client}
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	action := "ModifyDiskAttribute"
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DiskId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("delete_auto_snapshot") {
		update = true
		request["DeleteAutoSnapshot"] = d.Get("delete_auto_snapshot")
	}

	objectRaw, err := ecsServiceV2.DescribeEcsDisk(d.Id())
	if err != nil {
		return WrapError(err)
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
		update = true
		request["BurstingEnabled"] = d.Get("bursting_enabled")
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), query, request, &runtime)
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
	action = "ResizeDisk"
	conn, err = client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DiskId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("size") {
		update = true
	}
	request["NewSize"] = d.Get("size")
	if v, ok := d.GetOk("type"); ok {
		request["Type"] = v
	}
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), query, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"InternalError"}) || NeedRetry(err) {
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
	conn, err = client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
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
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), query, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"InternalError"}) || NeedRetry(err) {
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
	action = "ModifyDiskSpec"
	conn, err = client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
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
		update = true
		request["ProvisionedIops"] = d.Get("provisioned_iops")
	}

	if update {
		if _, ok := d.GetOkExists("dry_run"); ok {
			request["DryRun"] = d.Get("dry_run")
		}
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), query, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"ServiceUnavailable", "Throttling.ConcurrentLimitExceeded"}) || NeedRetry(err) {
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

		stateConf := BuildStateConf([]string{}, []string{"Available", "In_use"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, ecsServiceV2.EcsDiskStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "ModifyDiskChargeType"
	conn, err = client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DiskIds"] = convertListToJsonString([]interface{}{d.Id()})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	request["AutoPay"] = true
	request["InstanceId"] = d.Get("instance_id")
	if !d.IsNewResource() && d.HasChange("payment_type") {
		update = true
		request["DiskChargeType"] = convertEcsDiskPaymentTypeRequest(d.Get("payment_type").(string))
	}
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), query, request, &runtime)
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

		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(request["DiskChargeType"])}, d.Timeout(schema.TimeoutUpdate), 1*time.Second, ecsServiceV2.EcsDiskStateRefreshFunc(d.Id(), "DiskChargeType", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if !d.IsNewResource() && d.HasChange("tags") {
		if err := ecsServiceV2.SetResourceTags(d, "disk"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	d.Partial(false)
	return resourceAliCloudEcsDiskRead(d, meta)
}

func resourceAliCloudEcsDiskDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsServiceV2 := EcsServiceV2{client}

	objectRaw, err := ecsServiceV2.DescribeEcsDisk(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_disk DescribeEcsDisk Failed!!! %s", err)
			return nil
		}
		return WrapError(err)
	}

	if fmt.Sprint(objectRaw["DiskChargeType"]) == "PrePaid" && fmt.Sprint(objectRaw["DeleteWithInstance"]) == "true" {
		log.Printf("[DEBUG] Resource alicloud_ecs_disk %s charge type is PrePaid and its attribute DeleteWithInstance is true, so it will remove from state.", d.Id())
		return nil
	}

	action := "DeleteDisk"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["DiskId"] = d.Id()

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), query, request, &runtime)

		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectDiskStatus.Initializing"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDiskId.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func convertEcsDiskPaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "PostPaid":
		return "PayAsYouGo"
	case "PrePaid":
		return "Subscription"
	}
	return source
}

func convertEcsDiskPaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PostPaid"
	case "Subscription":
		return "PrePaid"
	}
	return source
}
