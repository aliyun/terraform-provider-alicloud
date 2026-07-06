package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudRdsCustomDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRdsCustomDiskCreate,
		Read:   resourceAliCloudRdsCustomDiskRead,
		Update: resourceAliCloudRdsCustomDiskUpdate,
		Delete: resourceAliCloudRdsCustomDiskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_pay": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"bursting_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
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
			"disk_category": {
				Type:     schema.TypeString,
				Required: true,
			},
			"disk_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"instance_charge_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"performance_level": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
				Required: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"period": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "Field `period` has been deprecated from provider version 1.283.0.",
			},
			"period_unit": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field `period_unit` has been deprecated from provider version 1.283.0.",
			},
		},
	}
}

func resourceAliCloudRdsCustomDiskCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateRCDisk"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	if v, ok := d.GetOk("performance_level"); ok {
		request["PerformanceLevel"] = v
	}
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
	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}
	if v, ok := d.GetOk("instance_charge_type"); ok {
		request["InstanceChargeType"] = v
	}
	if v, ok := d.GetOk("snapshot_id"); ok {
		request["SnapshotId"] = v
	}
	request["DiskCategory"] = d.Get("disk_category")
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	request["Size"] = d.Get("size")
	if v, ok := d.GetOk("disk_name"); ok {
		request["DiskName"] = v
	}
	request["ZoneId"] = d.Get("zone_id")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Rds", "2014-08-15", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_rds_custom_disk", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DiskId"]))

	rdsServiceV2 := RdsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available", "In_use"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, rdsServiceV2.RdsCustomDiskStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudRdsCustomDiskUpdate(d, meta)
}

func resourceAliCloudRdsCustomDiskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsServiceV2 := RdsServiceV2{client}

	objectRaw, err := rdsServiceV2.DescribeRdsCustomDisk(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_rds_custom_disk DescribeRdsCustomDisk Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("bursting_enabled", objectRaw["BurstingEnabled"])
	d.Set("create_time", objectRaw["CreationTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("disk_category", objectRaw["Category"])
	d.Set("disk_name", objectRaw["DiskName"])
	d.Set("instance_charge_type", objectRaw["DiskChargeType"])
	d.Set("performance_level", objectRaw["PerformanceLevel"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("size", objectRaw["Size"])
	d.Set("snapshot_id", objectRaw["SourceSnapshotId"])
	d.Set("status", objectRaw["Status"])
	d.Set("zone_id", objectRaw["ZoneId"])
	d.Set("delete_with_instance", objectRaw["DeleteWithInstance"])
	d.Set("instance_id", objectRaw["InstanceId"])

	tagsMaps, _ := jsonpath.Get("$.Tag", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudRdsCustomDiskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ResizeRCInstanceDisk"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DiskId"] = d.Id()
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	if !d.IsNewResource() && d.HasChange("size") {
		update = true
	}
	request["NewSize"] = d.Get("size")
	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}
	if v, ok := d.GetOk("type"); ok {
		request["Type"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"InvalidOrderTask.NotSupport"}) || NeedRetry(err) {
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
		rdsServiceV2 := RdsServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("size"))}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, rdsServiceV2.RdsCustomDiskStateRefreshFunc(d.Id(), "Size", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "ModifyResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBInstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["ResourceType"] = "CustomDisk"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, query, request, true)
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
	action = "ModifyRCDiskSpec"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DiskId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("performance_level") {
		update = true

		if v, ok := d.GetOk("performance_level"); ok {
			request["PerformanceLevel"] = v
		}
	}

	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}
	if !d.IsNewResource() && d.HasChange("disk_category") {
		update = true
	}
	request["DiskCategory"] = d.Get("disk_category")
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"InvalidOrderTask.NotSupport"}) || NeedRetry(err) {
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

		rdsServiceV2 := RdsServiceV2{client}
		if d.HasChange("disk_category") {
			stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("disk_category"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsServiceV2.RdsCustomDiskStateRefreshFunc(d.Id(), "Category", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		stateConf := BuildStateConf([]string{}, []string{"Available", "In_use"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsServiceV2.RdsCustomDiskStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}
	update = false
	action = "ModifyRCDiskAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DiskId"] = d.Id()
	request["RegionId"] = client.RegionId

	rdsServiceV2 := RdsServiceV2{client}
	objectRaw, err := rdsServiceV2.DescribeRdsCustomDisk(d.Id())
	if err != nil {
		return WrapError(err)
	}
	deleteWithInstance, ok := objectRaw["DeleteWithInstance"]
	if ok && (fmt.Sprint(deleteWithInstance) != fmt.Sprint(d.Get("delete_with_instance"))) {
		update = true

		if v, ok := d.GetOkExists("delete_with_instance"); ok {
			request["DeleteWithInstance"] = v
		}
	}
	burstingEnabled, ok := objectRaw["BurstingEnabled"]
	if ok && (fmt.Sprint(burstingEnabled) != fmt.Sprint(d.Get("bursting_enabled"))) {
		update = true

		if v, ok := d.GetOkExists("bursting_enabled"); ok {
			request["BurstingEnabled"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("disk_name") {
		update = true
		request["DiskName"] = d.Get("disk_name")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, query, request, true)
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
		if err := rdsServiceV2.SetResourceTags(d, "CustomDisk"); err != nil {
			return WrapError(err)
		}
	}
	return resourceAliCloudRdsCustomDiskRead(d, meta)
}

func resourceAliCloudRdsCustomDiskDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	rdsServiceV2 := RdsServiceV2{client}
	objectRaw, err := rdsServiceV2.DescribeRdsCustomDisk(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_rds_custom_disk DescribeRCDisks Failed!!! %s", err)
			return nil
		}
		return WrapError(err)
	}

	if fmt.Sprint(objectRaw["DiskChargeType"]) == "Prepaid" {
		log.Printf("[DEBUG] Resource alicloud_rds_custom_disk %s charge type is PrePaid, so it will remove from state.", d.Id())
		return nil
	}

	action := "DeleteRCDisk"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	request = make(map[string]interface{})
	request["DiskId"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Rds", "2014-08-15", action, query, request, true)

		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectDBInstanceState"}) || NeedRetry(err) {
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

	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutCreate), 5*time.Second, rdsServiceV2.RdsCustomDiskStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
