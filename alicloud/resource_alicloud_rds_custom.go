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
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"amount": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"auto_pay": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"auto_renew": {
				Type:     schema.TypeBool,
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
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

	request["Amount"] = d.Get("amount")
	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}
	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	if v, ok := d.GetOk("data_disk"); ok {
		dataDiskMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Category"] = dataLoopTmp["category"]
			dataLoopMap["PerformanceLevel"] = dataLoopTmp["performance_level"]
			dataLoopMap["Size"] = dataLoopTmp["size"]
			dataDiskMapsArray = append(dataDiskMapsArray, dataLoopMap)
		}
		dataDiskMapsJson, err := json.Marshal(dataDiskMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["DataDisk"] = string(dataDiskMapsJson)
	}

	if v, ok := d.GetOk("deployment_set_id"); ok {
		request["DeploymentSetId"] = v
	}
	request["InstanceType"] = d.Get("instance_type")
	if v, ok := d.GetOkExists("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("security_group_ids"); ok {
		jsonPathResult6, err := jsonpath.Get("$", v)
		if err == nil && jsonPathResult6 != "" {
			request["SecurityGroupId"] = jsonPathResult6.([]interface{})[0]
		}
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	request["VSwitchId"] = d.Get("vswitch_id")
	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}
	if v, ok := d.GetOk("security_enhancement_strategy"); ok {
		request["SecurityEnhancementStrategy"] = v
	}
	if v, ok := d.GetOk("password"); ok {
		request["Password"] = v
	}
	if v, ok := d.GetOk("internet_charge_type"); ok {
		request["InternetChargeType"] = v
	}
	if v, ok := d.GetOkExists("internet_max_bandwidth_out"); ok {
		request["InternetMaxBandwidthOut"] = v
	}
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}
	if v, ok := d.GetOk("io_optimized"); ok {
		request["IoOptimized"] = v
	}
	if v, ok := d.GetOk("image_id"); ok {
		request["ImageId"] = v
	}
	if v, ok := d.GetOk("instance_charge_type"); ok {
		request["InstanceChargeType"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("key_pair_name"); ok {
		request["KeyPairName"] = v
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("system_disk"); !IsNil(v) {
		size3, _ := jsonpath.Get("$[0].size", d.Get("system_disk"))
		if size3 != nil && size3 != "" {
			objectDataLocalMap["Size"] = size3
		}
		category3, _ := jsonpath.Get("$[0].category", d.Get("system_disk"))
		if category3 != nil && category3 != "" {
			objectDataLocalMap["Category"] = category3
		}

		objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
		if err != nil {
			return WrapError(err)
		}
		request["SystemDisk"] = string(objectDataLocalMapJson)
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("create_mode"); ok {
		request["CreateMode"] = v
	}
	if v, ok := d.GetOk("host_name"); ok {
		request["HostName"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Rds", "2014-08-15", action, query, request, false)
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

	if objectRaw["DeploymentSetId"] != nil {
		d.Set("deployment_set_id", objectRaw["DeploymentSetId"])
	}
	if objectRaw["Description"] != nil {
		d.Set("description", objectRaw["Description"])
	}
	if objectRaw["InstanceType"] != nil {
		d.Set("instance_type", objectRaw["InstanceType"])
	}
	if objectRaw["RegionId"] != nil {
		d.Set("region_id", objectRaw["RegionId"])
	}
	if objectRaw["ResourceGroupId"] != nil {
		d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}
	if objectRaw["ZoneId"] != nil {
		d.Set("zone_id", objectRaw["ZoneId"])
	}

	vpcAttributes1RawObj, _ := jsonpath.Get("$.VpcAttributes", objectRaw)
	vpcAttributes1Raw := make(map[string]interface{})
	if vpcAttributes1RawObj != nil {
		vpcAttributes1Raw = vpcAttributes1RawObj.(map[string]interface{})
	}
	if vpcAttributes1Raw["VSwitchId"] != nil {
		d.Set("vswitch_id", vpcAttributes1Raw["VSwitchId"])
	}

	dataDisk1Raw, _ := jsonpath.Get("$.DataDisks.DataDisk", objectRaw)
	dataDiskMaps := make([]map[string]interface{}, 0)
	if dataDisk1Raw != nil {
		for _, dataDiskChild1Raw := range dataDisk1Raw.([]interface{}) {
			dataDiskMap := make(map[string]interface{})
			dataDiskChild1Raw := dataDiskChild1Raw.(map[string]interface{})
			dataDiskMap["category"] = dataDiskChild1Raw["Category"]
			dataDiskMap["performance_level"] = dataDiskChild1Raw["PerformanceLevel"]
			dataDiskMap["size"] = dataDiskChild1Raw["Size"]

			dataDiskMaps = append(dataDiskMaps, dataDiskMap)
		}
	}
	if dataDisk1Raw != nil {
		if err := d.Set("data_disk", dataDiskMaps); err != nil {
			return err
		}
	}
	securityGroupId1Raw, _ := jsonpath.Get("$.SecurityGroupIds.SecurityGroupId", objectRaw)
	d.Set("security_group_ids", securityGroupId1Raw)

	objectRaw, err = rdsServiceV2.DescribeListTagResources(d.Id())
	if err != nil {
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
					response, err = client.RpcPost("Rds", "2014-08-15", action, query, request, false)
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
				stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, rdsServiceV2.RdsCustomStateRefreshFunc(d.Id(), "Status", []string{}))
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
					response, err = client.RpcPost("Rds", "2014-08-15", action, query, request, false)
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
				stateConf := BuildStateConf([]string{}, []string{"Stopped"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, rdsServiceV2.RdsCustomStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}

	action := "ModifyResourceGroup"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBInstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["ResourceGroupId"] = d.Get("resource_group_id")
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
	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("direction"); ok {
		request["Direction"] = v
	}
	if !d.IsNewResource() && d.HasChange("instance_type") {
		update = true
	}
	request["InstanceType"] = d.Get("instance_type")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, query, request, false)
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
		response, err = client.RpcPost("Rds", "2014-08-15", action, query, request, false)
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 60*time.Second, rdsServiceV2.RdsCustomStateRefreshFunc(d.Id(), "InstanceId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
