// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"create_extra_param": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_disk": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
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
					},
				},
			},
			"deployment_set_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"size": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"tags": tagsSchema(),
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
	// request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}
	if v, ok := d.GetOk("data_disk"); ok {
		dataDiskMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Size"] = dataLoopTmp["size"]
			dataLoopMap["PerformanceLevel"] = dataLoopTmp["performance_level"]
			dataLoopMap["Category"] = dataLoopTmp["category"]
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
	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}
	if v, ok := d.GetOk("key_pair_name"); ok {
		request["KeyPairName"] = v
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("system_disk"); !IsNil(v) {
		category3, _ := jsonpath.Get("$[0].category", v)
		if category3 != nil && category3 != "" {
			objectDataLocalMap["Category"] = category3
		}
		size3, _ := jsonpath.Get("$[0].size", v)
		if size3 != nil && size3 != "" {
			objectDataLocalMap["Size"] = size3
		}

		objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
		if err != nil {
			return WrapError(err)
		}
		request["SystemDisk"] = string(objectDataLocalMapJson)
	}

	if v, ok := d.GetOk("io_optimized"); ok {
		request["IoOptimized"] = v
	}
	if v, ok := d.GetOk("deployment_set_id"); ok {
		request["DeploymentSetId"] = v
	}
	if v, ok := d.GetOk("security_group_ids"); ok {
		jsonPathResult7, err := jsonpath.Get("$", v)
		if err == nil && jsonPathResult7 != "" {
			request["SecurityGroupId"] = jsonPathResult7.([]interface{})[0]
		}
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
	if v, ok := d.GetOk("create_extra_param"); ok {
		request["CreateExtraParam"] = v
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
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
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
	d.Set("instance_type", objectRaw["InstanceType"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["Status"])
	d.Set("zone_id", objectRaw["ZoneId"])

	vpcAttributesRawObj, _ := jsonpath.Get("$.VpcAttributes", objectRaw)
	vpcAttributesRaw := make(map[string]interface{})
	if vpcAttributesRawObj != nil {
		vpcAttributesRaw = vpcAttributesRawObj.(map[string]interface{})
	}
	d.Set("vswitch_id", vpcAttributesRaw["VSwitchId"])

	dataDiskRaw, _ := jsonpath.Get("$.DataDisks.DataDisk", objectRaw)
	dataDiskMaps := make([]map[string]interface{}, 0)
	if dataDiskRaw != nil {
		for _, dataDiskChildRaw := range dataDiskRaw.([]interface{}) {
			dataDiskMap := make(map[string]interface{})
			dataDiskChildRaw := dataDiskChildRaw.(map[string]interface{})
			dataDiskMap["category"] = dataDiskChildRaw["Category"]
			dataDiskMap["performance_level"] = dataDiskChildRaw["PerformanceLevel"]
			dataDiskMap["size"] = dataDiskChildRaw["Size"]

			dataDiskMaps = append(dataDiskMaps, dataDiskMap)
		}
	}
	if err := d.Set("data_disk", dataDiskMaps); err != nil {
		return err
	}
	securityGroupIdRaw, _ := jsonpath.Get("$.SecurityGroupIds.SecurityGroupId", objectRaw)
	d.Set("security_group_ids", securityGroupIdRaw)

	objectRaw, err = rdsServiceV2.DescribeCustomListTagResources(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	tagsMaps, _ := jsonpath.Get("$.TagResources.TagResource", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudRdsCustomUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	if d.HasChange("status") {
		var err error
		rdsServiceV2 := RdsServiceV2{client}
		object, err := rdsServiceV2.DescribeRdsCustom(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "Running" {
				action := "StartRCInstance"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["InstanceId"] = d.Id()
				request["RegionId"] = client.RegionId
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Rds", "2014-08-15", action, query, request, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"undefined"}) || NeedRetry(err) {
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
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Rds", "2014-08-15", action, query, request, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"undefined"}) || NeedRetry(err) {
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
	request["ClientToken"] = buildClientToken(action)
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["ResourceGroupId"] = d.Get("resource_group_id")
	request["ResourceType"] = "Custom"
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
	action = "ModifyRCInstance"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("instance_type") {
		update = true
	}
	request["InstanceType"] = d.Get("instance_type")
	if v, ok := d.GetOk("auto_pay"); ok {
		request["AutoPay"] = v
	}
	if v, ok := d.GetOk("direction"); ok {
		request["Direction"] = v
	}
	if v, ok := d.GetOk("dry_run"); ok {
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
		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("instance_type"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsServiceV2.RdsCustomStateRefreshFunc(d.Id(), "InstanceType", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("tags") {
		rdsServiceV2 := RdsServiceV2{client}
		if err := rdsServiceV2.SetResourceTags(d, "Custom"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
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
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
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
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound", "InvalidDBInstance.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	rdsServiceV2 := RdsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 60*time.Second, rdsServiceV2.RdsCustomStateRefreshFunc(d.Id(), "$.InstanceId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
