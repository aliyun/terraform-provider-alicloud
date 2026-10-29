// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"sort"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudRdsCustom() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRdsCustomCreate,
		Read:   resourceAliCloudRdsCustomRead,
		Update: resourceAliCloudRdsCustomUpdate,
		Delete: resourceAliCloudRdsCustomDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(7 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"amount": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"auto_pay": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"auto_use_coupon": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"create_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_disk": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"device": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"encrypted": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: StringInSlice([]string{"true", "false"}, false),
						},
						"performance_level": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"deployment_set_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"direction": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[\\w.,;/@-]+$"), "Instance configuration type, value range:> This parameter does not need to be uploaded, and the system can automatically determine whether to upgrade or downgrade. If you want to upload, please follow the following logic rules.-**Up** (default): upgrade the instance specification. Please ensure that your account balance is sufficient.-**Down**: Downgrade instance specifications. When the instance type set to InstanceType is lower than the current instance type, set Direction = down."),
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_jumbo_frame": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"force_stop": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_charge_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"internet_charge_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"internet_max_bandwidth_out": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"io_optimized": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_pair_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"period_unit": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"promotion_code": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"reboot": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"reboot_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"reboot_when_finished": {
				Type:     schema.TypeBool,
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
			"security_enhancement_strategy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"spot_strategy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"support_case": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"system_disk": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"performance_level": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"size": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"system_disk_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"user_data_in_base64": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"create_extra_param": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudRdsCustomCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "RunRCInstances"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}
	if v, ok := d.GetOk("data_disk"); ok {
		dataDiskMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Size"] = dataLoopTmp["size"]
			dataLoopMap["PerformanceLevel"] = dataLoopTmp["performance_level"]
			dataLoopMap["Category"] = dataLoopTmp["category"]
			if device, ok := dataLoopTmp["device"].(string); ok && device != "" {
				dataLoopMap["Device"] = device
			}
			if encrypted, ok := dataLoopTmp["encrypted"].(string); ok && encrypted != "" {
				dataLoopMap["Encrypted"] = encrypted
			}
			if snapshotId, ok := dataLoopTmp["snapshot_id"].(string); ok && snapshotId != "" {
				dataLoopMap["SnapshotId"] = snapshotId
			}
			dataDiskMapsArray = append(dataDiskMapsArray, dataLoopMap)
		}
		dataDiskMapsJson, err := json.Marshal(dataDiskMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["DataDisk"] = string(dataDiskMapsJson)
	}

	request["VSwitchId"] = d.Get("vswitch_id")
	if v, ok := d.GetOk("password"); ok {
		request["Password"] = v
	}
	if v, ok := d.GetOk("private_ip_address"); ok {
		request["PrivateIpAddress"] = v
	}
	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}
	if v, ok := d.GetOk("key_pair_name"); ok {
		request["KeyPairName"] = v
	}
	systemDisk := make(map[string]interface{})

	if v := d.Get("system_disk"); !IsNil(v) {
		category3, _ := jsonpath.Get("$[0].category", v)
		if category3 != nil && category3 != "" {
			systemDisk["Category"] = category3
		}
		size3, _ := jsonpath.Get("$[0].size", v)
		if size3 != nil && size3 != "" {
			systemDisk["Size"] = size3
		}
		performanceLevel3, _ := jsonpath.Get("$[0].performance_level", v)
		if performanceLevel3 != nil && performanceLevel3 != "" {
			systemDisk["PerformanceLevel"] = performanceLevel3
		}

		systemDiskJson, err := json.Marshal(systemDisk)
		if err != nil {
			return WrapError(err)
		}
		request["SystemDisk"] = string(systemDiskJson)
	}

	if v, ok := d.GetOkExists("enable_jumbo_frame"); ok {
		networkOptions := map[string]interface{}{
			"EnableJumboFrame": v,
		}
		networkOptionsJson, err := json.Marshal(networkOptions)
		if err != nil {
			return WrapError(err)
		}
		request["NetworkOptions"] = string(networkOptionsJson)
	}

	if v, ok := d.GetOk("io_optimized"); ok {
		request["IoOptimized"] = v
	}
	if v, ok := d.GetOk("deployment_set_id"); ok {
		request["DeploymentSetId"] = v
	}
	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsMapsArray := convertToInterfaceArray(v)

		securityGroupIdsMapsJson, err := json.Marshal(securityGroupIdsMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["SecurityGroupIds"] = string(securityGroupIdsMapsJson)
	}

	request["InstanceType"] = d.Get("instance_type")
	if v, ok := d.GetOk("create_mode"); ok {
		request["CreateMode"] = v
	}
	if v, ok := d.GetOk("support_case"); ok {
		request["SupportCase"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("internet_charge_type"); ok {
		request["InternetChargeType"] = v
	}
	if v, ok := d.GetOk("spot_strategy"); ok {
		request["SpotStrategy"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("instance_name"); ok {
		request["InstanceName"] = v
	}
	if v, ok := d.GetOk("image_id"); ok {
		request["ImageId"] = v
	}
	if v, ok := d.GetOk("instance_charge_type"); ok {
		request["InstanceChargeType"] = v
	}
	if v, ok := d.GetOkExists("amount"); ok {
		request["Amount"] = v
	}
	if v, ok := d.GetOk("host_name"); ok {
		request["HostName"] = v
	}
	if v, ok := d.GetOkExists("internet_max_bandwidth_out"); ok {
		request["InternetMaxBandwidthOut"] = v
	}
	if v, ok := d.GetOk("security_enhancement_strategy"); ok {
		request["SecurityEnhancementStrategy"] = v
	}
	if v, ok := d.GetOkExists("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}
	if v, ok := d.GetOk("create_extra_param"); ok {
		request["CreateExtraParam"] = v
	}
	if v, ok := d.GetOkExists("auto_use_coupon"); ok {
		request["AutoUseCoupon"] = v
	}
	if v, ok := d.GetOk("promotion_code"); ok {
		request["PromotionCode"] = v
	}
	if v, ok := d.GetOkExists("deletion_protection"); ok {
		request["DeletionProtection"] = v
	}
	if v, ok := d.GetOk("user_data"); ok {
		request["UserData"] = v
	}
	if v, ok := d.GetOkExists("user_data_in_base64"); ok {
		request["UserDataInBase64"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Rds", "2014-08-15", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"InsufficientCapacity"}) || NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_rds_custom", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.InstanceIdSets.InstanceIdSet[0]", response)
	d.SetId(fmt.Sprint(id))

	rdsServiceV2 := RdsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 40*time.Second, rdsServiceV2.RdsCustomStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudRdsCustomUpdate(d, meta)
}

func resourceAliCloudRdsCustomRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsServiceV2 := RdsServiceV2{client}

	objectRaw, err := rdsServiceV2.DescribeRdsCustom(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_rds_custom DescribeRdsCustom Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("deployment_set_id", objectRaw["DeploymentSetId"])
	d.Set("description", objectRaw["Description"])
	d.Set("instance_name", objectRaw["InstanceName"])
	d.Set("instance_type", objectRaw["InstanceType"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["Status"])
	d.Set("zone_id", objectRaw["ZoneId"])

	if v, ok := objectRaw["DeletionProtection"].(bool); ok {
		d.Set("deletion_protection", v)
	}
	if v, ok := objectRaw["EnableJumboFrame"].(bool); ok {
		d.Set("enable_jumbo_frame", v)
	}

	// DescribeRCInstanceAttribute always returns UserData as a Base64 string.
	// Decode it back to the configured representation so refresh and import keep
	// reporting the value the user configured.
	if userDataRaw, ok := objectRaw["UserData"].(string); ok && userDataRaw != "" {
		if rdsCustomRawConfigHasAttribute(d, "user_data_in_base64") && d.Get("user_data_in_base64").(bool) {
			d.Set("user_data", userDataRaw)
		} else if decoded, err := base64.StdEncoding.DecodeString(userDataRaw); err == nil {
			d.Set("user_data", string(decoded))
		} else {
			d.Set("user_data", userDataRaw)
		}
	}

	vpcAttributesRawObj, _ := jsonpath.Get("$.VpcAttributes", objectRaw)
	vpcAttributesRaw := make(map[string]interface{})
	if vpcAttributesRawObj != nil {
		vpcAttributesRaw = vpcAttributesRawObj.(map[string]interface{})
	}
	d.Set("vswitch_id", vpcAttributesRaw["VSwitchId"])

	privateIpAddressRaw, _ := jsonpath.Get("$.VpcAttributes.PrivateIpAddress.IpAddress", objectRaw)
	if arr, ok := privateIpAddressRaw.([]interface{}); ok && len(arr) > 0 {
		d.Set("private_ip_address", arr[0])
	}

	dataDiskRaw, _ := jsonpath.Get("$.DataDisks.DataDisk", objectRaw)
	dataDiskMaps := make([]map[string]interface{}, 0)
	if dataDiskRaw != nil {
		for _, dataDiskChildRaw := range convertToInterfaceArray(dataDiskRaw) {
			dataDiskMap := make(map[string]interface{})
			dataDiskChildRaw := dataDiskChildRaw.(map[string]interface{})
			dataDiskMap["category"] = dataDiskChildRaw["Category"]
			dataDiskMap["performance_level"] = dataDiskChildRaw["PerformanceLevel"]
			dataDiskMap["size"] = dataDiskChildRaw["Size"]
			if v, ok := dataDiskChildRaw["Device"].(string); ok {
				dataDiskMap["device"] = v
			}
			if v, ok := dataDiskChildRaw["Encrypted"].(string); ok {
				dataDiskMap["encrypted"] = v
			}
			if v, ok := dataDiskChildRaw["SnapshotId"].(string); ok {
				dataDiskMap["snapshot_id"] = v
			}

			dataDiskMaps = append(dataDiskMaps, dataDiskMap)
		}
	}
	dataDiskMaps = orderRdsCustomDataDisks(d.Get("data_disk"), dataDiskMaps)
	if err := d.Set("data_disk", dataDiskMaps); err != nil {
		return err
	}
	securityGroupIdRaw, _ := jsonpath.Get("$.SecurityGroupIds.SecurityGroupId", objectRaw)
	d.Set("security_group_ids", securityGroupIdRaw)
	systemDiskMaps := make([]map[string]interface{}, 0)
	systemDiskMap := make(map[string]interface{})
	systemDiskRaw := make(map[string]interface{})
	if objectRaw["SystemDisk"] != nil {
		systemDiskRaw = objectRaw["SystemDisk"].(map[string]interface{})
	}
	if len(systemDiskRaw) > 0 {
		systemDiskMap["category"] = systemDiskRaw["SystemDiskCategory"]
		systemDiskMap["size"] = systemDiskRaw["SystemDiskSize"]
		if v, ok := systemDiskRaw["SystemDiskPerformanceLevel"].(string); ok {
			systemDiskMap["performance_level"] = v
		}

		systemDiskMaps = append(systemDiskMaps, systemDiskMap)
	}
	if err := d.Set("system_disk", systemDiskMaps); err != nil {
		return err
	}

	objectRaw, err = rdsServiceV2.DescribeCustomListTagResources(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	tagsMaps, _ := jsonpath.Get("$.TagResources.TagResource", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	objectRaw, err = rdsServiceV2.DescribeCustomDescribeRCDisks(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("system_disk_id", objectRaw["DiskId"])

	return nil
}

func resourceAliCloudRdsCustomUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	rdsServiceV2 := RdsServiceV2{client}
	objectRaw, _ := rdsServiceV2.DescribeRdsCustom(d.Id())

	if d.HasChange("status") {
		var err error
		target := d.Get("status").(string)

		currentStatus, err := jsonpath.Get("Status", objectRaw)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "Status", objectRaw)
		}
		if fmt.Sprint(currentStatus) != target {
			if target == "Running" {
				action := "StartRCInstance"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["InstanceId"] = d.Id()
				request["RegionId"] = client.RegionId

				wait := incrementalWait(5*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("Rds", "2014-08-15", action, query, request, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"undefined"}) || NeedRetry(err) {
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
				rdsServiceV2 := RdsServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, rdsServiceV2.RdsCustomStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			if target == "Stopped" {
				action := "StopRCInstance"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["InstanceId"] = d.Id()
				request["RegionId"] = client.RegionId
				if v, ok := d.GetOkExists("force_stop"); ok {
					request["ForceStop"] = v
				}

				wait := incrementalWait(5*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("Rds", "2014-08-15", action, query, request, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"undefined"}) || NeedRetry(err) {
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
				rdsServiceV2 := RdsServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Stopped"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, rdsServiceV2.RdsCustomStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}

	var err error
	action := "ModifyResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBInstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["ResourceGroupId"] = d.Get("resource_group_id")
	request["ResourceType"] = "Custom"

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, query, request, true)
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
	update = false
	action = "ModifyRCInstance"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("instance_type") {
		update = true
	}
	request["InstanceType"] = d.Get("instance_type")
	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}
	if v, ok := d.GetOk("direction"); ok {
		request["Direction"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if rdsCustomRawConfigHasAttribute(d, "auto_use_coupon") {
		request["AutoUseCoupon"] = d.Get("auto_use_coupon").(bool)
	}
	if v, ok := d.GetOk("promotion_code"); ok {
		request["PromotionCode"] = v
	}
	rebootWhenFinished := true
	if rdsCustomRawConfigHasAttribute(d, "reboot_when_finished") {
		rebootWhenFinished = d.Get("reboot_when_finished").(bool)
		request["RebootWhenFinished"] = rebootWhenFinished
	}
	rebootTime, _ := d.Get("reboot_time").(string)
	if rebootTime != "" {
		request["RebootTime"] = rebootTime
	}

	if update {
		if rebootTime != "" {
			if err := validateRdsCustomRebootTime(rebootTime, rebootWhenFinished); err != nil {
				return WrapError(err)
			}
		} else if !rebootWhenFinished && fmt.Sprint(objectRaw["Status"]) == "Running" {
			return WrapError(errors.New("reboot_time is required when reboot_when_finished is false and the instance is Running"))
		}
		wait := incrementalWait(20*time.Second, 20*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"InvalidOrderTask.NotSupport"}) || NeedRetry(err) {
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
		rdsServiceV2 := RdsServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("instance_type"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsServiceV2.RdsCustomStateRefreshFunc(d.Id(), "InstanceType", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "ModifyRCInstanceDescription"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["InstanceDescription"] = d.Get("description")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, query, request, true)
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
	update = false
	action = "ModifyRCInstanceAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("host_name"); ok {
		request["HostName"] = v
	}
	if !d.IsNewResource() && d.HasChange("security_group_ids") {
		update = true
		if v, ok := d.GetOk("security_group_ids"); ok || d.HasChange("security_group_ids") {
			securityGroupIdsMapsArray := convertToInterfaceArray(v)

			securityGroupIdsMapsJson, err := json.Marshal(securityGroupIdsMapsArray)
			if err != nil {
				return WrapError(err)
			}
			request["SecurityGroupIds"] = string(securityGroupIdsMapsJson)
		}
	}

	if !d.IsNewResource() && d.HasChange("instance_name") {
		update = true
		request["InstanceName"] = d.Get("instance_name")
	}

	if v, ok := d.GetOk("password"); ok {
		request["Password"] = v
	}

	if !d.IsNewResource() && (d.HasChange("deletion_protection") || d.HasChange("enable_jumbo_frame") || d.HasChange("reboot")) {
		update = true
	}
	if !d.IsNewResource() && d.HasChange("deletion_protection") && rdsCustomRawConfigHasAttribute(d, "deletion_protection") {
		request["DeletionProtection"] = d.Get("deletion_protection").(bool)
	}
	if !d.IsNewResource() && d.HasChange("enable_jumbo_frame") && rdsCustomRawConfigHasAttribute(d, "enable_jumbo_frame") {
		request["EnableJumboFrame"] = d.Get("enable_jumbo_frame").(bool)
	}
	if !d.IsNewResource() && d.HasChange("reboot") && d.Get("reboot").(bool) {
		request["Reboot"] = true
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, query, request, true)
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

	if d.HasChange("tags") {
		rdsServiceV2 := RdsServiceV2{client}
		if err := rdsServiceV2.SetResourceTags(d, "Custom"); err != nil {
			return WrapError(err)
		}
	}
	return resourceAliCloudRdsCustomRead(d, meta)
}

func resourceAliCloudRdsCustomDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteRCInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOkExists("force"); ok {
		request["Force"] = v
	}
	wait := incrementalWait(5*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Rds", "2014-08-15", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectInstanceStatus"}) || NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound", "InvalidDBInstance.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	rdsServiceV2 := RdsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 60*time.Second, rdsServiceV2.RdsCustomStateRefreshFunc(d.Id(), "$.InstanceId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

// rdsCustomRawConfigHasAttribute reports whether the attribute is present in the
// configuration (explicit null is treated as absent). It distinguishes an
// explicit false from an unset or removed attribute, which GetOkExists cannot.
func rdsCustomRawConfigHasAttribute(d *schema.ResourceData, name string) bool {
	raw := d.GetRawConfig()
	if raw.IsNull() || !raw.Type().IsObjectType() || !raw.Type().HasAttribute(name) {
		return false
	}
	return !raw.GetAttr(name).IsNull()
}

// validateRdsCustomRebootTime mirrors the ModifyRCInstance contract: when
// RebootWhenFinished is false the restart time is mandatory, must be expressed
// in UTC as yyyy-MM-ddTHH:mmZ and must fall within the next 48 hours.
func validateRdsCustomRebootTime(rebootTime string, rebootWhenFinished bool) error {
	if rebootWhenFinished {
		return nil
	}
	if rebootTime == "" {
		return errors.New("reboot_time is required when reboot_when_finished is false")
	}
	parsed, err := time.Parse("2006-01-02T15:04Z", rebootTime)
	if err != nil {
		return fmt.Errorf("reboot_time must follow the ISO 8601 UTC format yyyy-MM-ddTHH:mmZ: %s", err)
	}
	now := time.Now().UTC()
	if parsed.Before(now) || parsed.After(now.Add(48*time.Hour)) {
		return errors.New("reboot_time must be within the next 48 hours in UTC")
	}
	return nil
}

// orderRdsCustomDataDisks keeps the order the disks are tracked in the
// configuration/state. DescribeRCInstanceAttribute does not echo the submitted
// order, so members are joined by device (or by category/size/performance level
// when no device is configured). Without configured members (import) the
// canonical device order is used, which is the order the instance was created
// with.
func orderRdsCustomDataDisks(prior interface{}, disks []map[string]interface{}) []map[string]interface{} {
	priorList := convertToInterfaceArray(prior)
	if len(priorList) == 0 {
		sort.SliceStable(disks, func(i, j int) bool {
			left, _ := disks[i]["device"].(string)
			right, _ := disks[j]["device"].(string)
			if left == "" || right == "" {
				return left != "" && right == ""
			}
			return left < right
		})
		return disks
	}
	used := make([]bool, len(disks))
	ordered := make([]map[string]interface{}, 0, len(disks))
	for _, item := range priorList {
		priorItem, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		for i, disk := range disks {
			if used[i] || !rdsCustomDataDiskMatches(priorItem, disk) {
				continue
			}
			used[i] = true
			ordered = append(ordered, disk)
			break
		}
	}
	for i, disk := range disks {
		if !used[i] {
			ordered = append(ordered, disk)
		}
	}
	return ordered
}

func rdsCustomDataDiskMatches(prior, remote map[string]interface{}) bool {
	if device, _ := prior["device"].(string); device != "" {
		if remoteDevice, _ := remote["device"].(string); remoteDevice != "" {
			return device == remoteDevice
		}
	}
	return fmt.Sprint(prior["category"]) == fmt.Sprint(remote["category"]) &&
		fmt.Sprint(prior["size"]) == fmt.Sprint(remote["size"]) &&
		fmt.Sprint(prior["performance_level"]) == fmt.Sprint(remote["performance_level"])
}
