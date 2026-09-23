// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEcsElasticityAssurance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEcsElasticityAssuranceCreate,
		Read:   resourceAliCloudEcsElasticityAssuranceRead,
		Update: resourceAliCloudEcsElasticityAssuranceUpdate,
		Delete: resourceAliCloudEcsElasticityAssuranceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"assurance_times": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"auto_renew_period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"auto_renew_period_unit": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"elasticity_assurance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_amount": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"instance_charge_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_type": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			},
			"period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Month", "Year"}, false),
			},
			"private_pool_options_match_criteria": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"private_pool_options_name": {
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
				ForceNew: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old != "" && new != "" && strings.HasPrefix(new, strings.Trim(old, "00Z"))
				},
			},
			"start_time_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"used_assurance_times": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"zone_ids": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAliCloudEcsElasticityAssuranceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateElasticityAssurance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("instance_type"); ok {
		instanceTypeMapsArray := v.([]interface{})
		request["InstanceType"] = instanceTypeMapsArray
	}

	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("private_pool_options_name"); ok {
		request["PrivatePoolOptions.Name"] = v
	}
	if v, ok := d.GetOk("assurance_times"); ok {
		request["AssuranceTimes"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("start_time"); ok {
		request["StartTime"] = v
	}
	if v, ok := d.GetOk("private_pool_options_match_criteria"); ok {
		request["PrivatePoolOptions.MatchCriteria"] = v
	}
	if v, ok := d.GetOkExists("auto_renew_period"); ok {
		request["AutoRenewPeriod"] = v
	}
	if v, ok := d.GetOkExists("instance_amount"); ok {
		request["InstanceAmount"] = v
	}
	if v, ok := d.GetOkExists("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("zone_ids"); ok {
		zoneIdMapsArray := v.([]interface{})
		request["ZoneId"] = zoneIdMapsArray
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_elasticity_assurance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["PrivatePoolOptionsId"]))

	ecsServiceV2 := EcsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Active", "Prepared"}, d.Timeout(schema.TimeoutCreate), 20*time.Second, ecsServiceV2.EcsElasticityAssuranceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEcsElasticityAssuranceUpdate(d, meta)
}

func resourceAliCloudEcsElasticityAssuranceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsServiceV2 := EcsServiceV2{client}

	objectRaw, err := ecsServiceV2.DescribeEcsElasticityAssurance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_elasticity_assurance DescribeEcsElasticityAssurance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("assurance_times", objectRaw["TotalAssuranceTimes"])
	d.Set("description", objectRaw["Description"])
	d.Set("end_time", objectRaw["EndTime"])
	d.Set("instance_charge_type", objectRaw["InstanceChargeType"])
	d.Set("private_pool_options_match_criteria", objectRaw["PrivatePoolOptionsMatchCriteria"])
	d.Set("private_pool_options_name", objectRaw["PrivatePoolOptionsName"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("start_time_type", objectRaw["StartTimeType"])
	d.Set("status", objectRaw["Status"])
	d.Set("used_assurance_times", objectRaw["UsedAssuranceTimes"])
	d.Set("elasticity_assurance_id", objectRaw["PrivatePoolOptionsId"])
	d.Set("start_time", objectRaw["StartTime"])

	if v, ok := objectRaw["Tags"].(map[string]interface{}); ok {
		d.Set("tags", tagsToMap(v["Tag"]))
	}

	if v, ok := objectRaw["AllocatedResources"]; ok {
		allocatedResources := v.(map[string]interface{})
		if v, ok := allocatedResources["AllocatedResource"]; ok && len(v.([]interface{})) > 0 {
			allocatedResourceMap := v.([]interface{})[0].(map[string]interface{})
			d.Set("instance_type", []string{fmt.Sprint(allocatedResourceMap["InstanceType"])})
			d.Set("instance_amount", allocatedResourceMap["TotalAmount"])
			d.Set("zone_ids", []string{fmt.Sprint(allocatedResourceMap["zoneId"])})
		}
	}

	objectRaw, err = ecsServiceV2.DescribeElasticityAssuranceDescribeElasticityAssuranceAutoRenewAttribute(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("auto_renew", formatBool(convertEcsElasticityAssuranceElasticityAssuranceRenewAttributesElasticityAssuranceRenewAttributeRenewalStatusResponse(objectRaw["RenewalStatus"])))
	d.Set("auto_renew_period", objectRaw["Period"])
	d.Set("auto_renew_period_unit", objectRaw["PeriodUnit"])
	d.Set("elasticity_assurance_id", objectRaw["PrivatePoolOptionsId"])

	return nil
}

func resourceAliCloudEcsElasticityAssuranceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "ModifyElasticityAssurance"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["PrivatePoolOptions.Id"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("private_pool_options_name") {
		update = true
		request["PrivatePoolOptions.Name"] = d.Get("private_pool_options_name")
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("instance_amount") {
		update = true
		request["InstanceAmount"] = d.Get("instance_amount")
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
		ecsServiceV2 := EcsServiceV2{client}

		if d.HasChange("instance_amount") {
			stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("instance_amount"))}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, ecsServiceV2.EcsElasticityAssuranceStateRefreshFunc(d.Id(), "InstanceAmount", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
	}
	update = false
	action = "ModifyElasticityAssuranceAutoRenewAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["PrivatePoolOptions.Id.1"] = d.Id()
	request["RegionId"] = client.RegionId

	if !d.IsNewResource() && d.HasChange("auto_renew") {
		update = true
	}
	isAutoRenew := d.Get("auto_renew").(bool)
	request["RenewalStatus"] = convertEcsElasticityAssuranceRenewalStatusRequest(isAutoRenew)

	if isAutoRenew && d.HasChange("auto_renew_period_unit") {
		update = true
	}
	request["PeriodUnit"] = d.Get("auto_renew_period_unit")

	if !d.IsNewResource() && d.HasChange("auto_renew_period") {
		update = true
	}
	request["Period"] = d.Get("auto_renew_period")

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
		if err := ecsServiceV2.SetResourceTags(d, "ELASTICITYASSURANCE"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudEcsElasticityAssuranceRead(d, meta)
}

func resourceAliCloudEcsElasticityAssuranceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "RefundInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)

	request["ProductCode"] = "ecsrep"
	request["ImmediatelyRelease"] = "1"
	var endpoint string
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"ServiceUnavailable"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{""}) {
				endpoint = connectivity.BssOpenAPIEndpointInternational
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

	ecsServiceV2 := EcsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ecsServiceV2.EcsElasticityAssuranceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func convertEcsElasticityAssuranceElasticityAssuranceRenewAttributesElasticityAssuranceRenewAttributeRenewalStatusResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "AutoRenewal":
		return "true"
	case "Normal":
		return "false"
	case "NotRenewal":
		return "false"
	}
	return source
}
func convertEcsElasticityAssuranceRenewalStatusRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "true":
		return "AutoRenewal"
	case "false":
		return "Normal"
	}
	return source
}
