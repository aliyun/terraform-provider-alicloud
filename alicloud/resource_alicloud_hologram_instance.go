// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudHologramInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudHologramInstanceCreate,
		Read:   resourceAliCloudHologramInstanceRead,
		Update: resourceAliCloudHologramInstanceUpdate,
		Delete: resourceAliCloudHologramInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_pay": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cold_storage_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"cpu": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"duration": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"endpoints": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"vpc_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alternative_endpoints": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"gateway_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"initial_databases": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^.{2,64}$"), "The name of the resource"),
			},
			"instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Standard", "Follower", "Warehouse", "Shared", "Serverless"}, false),
			},
			"leader_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"pricing_cycle": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Month", "Hour"}, false),
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
			"scale_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"DOWNGRADE", "UPGRADE"}, false),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Running", "Allocating", "Suspended", "Creating"}, false),
			},
			"storage_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"tags": tagsSchema(),
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudHologramInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/api/v1/instances/create")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOkExists("duration"); ok {
		request["duration"] = v
	}
	if v, ok := d.GetOk("leader_instance_id"); ok {
		request["leaderInstanceId"] = v
	}
	jsonPathResult2, err := jsonpath.Get("$[*].vswitch_id", d.Get("endpoints").(*schema.Set).List())
	if err == nil {
		request["vSwitchId"] = convertListToCommaSeparate(filterEmptyStrings(jsonPathResult2.([]interface{})))
	}

	if v, ok := d.GetOkExists("gateway_count"); ok {
		request["gatewayCount"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["resourceGroupId"] = v
	}
	request["chargeType"] = convertHologramInstancechargeTypeRequest(d.Get("payment_type").(string))
	request["instanceName"] = d.Get("instance_name")
	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["autoPay"] = v
	}
	if v, ok := d.GetOk("pricing_cycle"); ok {
		request["pricingCycle"] = v
	}
	jsonPathResult9, err := jsonpath.Get("$[*].vpc_id", d.Get("endpoints").(*schema.Set).List())
	if err == nil {
		request["vpcId"] = convertListToCommaSeparate(filterEmptyStrings(jsonPathResult9.([]interface{})))
	}

	request["regionId"] = client.RegionId
	if v, ok := d.GetOkExists("cpu"); ok {
		request["cpu"] = v
	}
	if v, ok := d.GetOk("initial_databases"); ok {
		request["initialDatabases"] = v
	}
	request["instanceType"] = d.Get("instance_type")
	if v, ok := d.GetOkExists("storage_size"); ok {
		request["storageSize"] = v
	}
	request["zoneId"] = d.Get("zone_id")
	if v, ok := d.GetOkExists("cold_storage_size"); ok {
		request["coldStorageSize"] = v
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("Hologram", "2022-06-01", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_hologram_instance", action, AlibabaCloudSdkGoERROR)
	}

	code, _ := jsonpath.Get("$.Data.Success", response)
	if fmt.Sprint(code) != "true" {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_hologram_instance", action, AlibabaCloudSdkGoERROR, response)
	}

	id, _ := jsonpath.Get("$.Data.InstanceId", response)
	d.SetId(fmt.Sprint(id))

	hologramServiceV2 := HologramServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, hologramServiceV2.HologramInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudHologramInstanceUpdate(d, meta)
}

func resourceAliCloudHologramInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hologramServiceV2 := HologramServiceV2{client}

	objectRaw, err := hologramServiceV2.DescribeHologramInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_hologram_instance DescribeHologramInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cold_storage_size", objectRaw["ColdStorage"])
	d.Set("cpu", objectRaw["Cpu"])
	d.Set("create_time", objectRaw["CreationTime"])
	d.Set("gateway_count", objectRaw["GatewayCount"])
	d.Set("instance_name", objectRaw["InstanceName"])
	d.Set("instance_type", objectRaw["InstanceType"])
	d.Set("leader_instance_id", objectRaw["LeaderInstanceId"])
	d.Set("payment_type", convertHologramInstanceInstanceInstanceChargeTypeResponse(objectRaw["InstanceChargeType"]))
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["InstanceStatus"])
	d.Set("storage_size", objectRaw["Disk"])
	d.Set("zone_id", objectRaw["ZoneId"])

	endpointsRaw := objectRaw["Endpoints"]
	endpointsMaps := make([]map[string]interface{}, 0)
	if endpointsRaw != nil {
		for _, endpointsChildRaw := range endpointsRaw.([]interface{}) {
			endpointsMap := make(map[string]interface{})
			endpointsChildRaw := endpointsChildRaw.(map[string]interface{})
			endpointsMap["alternative_endpoints"] = endpointsChildRaw["AlternativeEndpoints"]
			endpointsMap["enabled"] = endpointsChildRaw["Enabled"]
			endpointsMap["endpoint"] = endpointsChildRaw["Endpoint"]
			endpointsMap["type"] = endpointsChildRaw["Type"]
			endpointsMap["vswitch_id"] = endpointsChildRaw["VSwitchId"]
			endpointsMap["vpc_id"] = endpointsChildRaw["VpcId"]
			endpointsMap["vpc_instance_id"] = endpointsChildRaw["VpcInstanceId"]

			endpointsMaps = append(endpointsMaps, endpointsMap)
		}
	}
	if err := d.Set("endpoints", endpointsMaps); err != nil {
		return err
	}
	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudHologramInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	d.Partial(true)

	if d.HasChange("status") {
		var err error
		hologramServiceV2 := HologramServiceV2{client}
		object, err := hologramServiceV2.DescribeHologramInstance(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["InstanceStatus"].(string) != target {
			if target == "Suspended" {
				instanceId := d.Id()
				action := fmt.Sprintf("/api/v1/instances/%s/stop", instanceId)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})
				request["instanceId"] = d.Id()

				body = request
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RoaPost("Hologram", "2022-06-01", action, query, nil, body, true)
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
				hologramServiceV2 := HologramServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Suspended"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, hologramServiceV2.HologramInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			if target == "Running" {
				instanceId := d.Id()
				action := fmt.Sprintf("/api/v1/instances/%s/resume", instanceId)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})
				request["instanceId"] = d.Id()

				body = request
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RoaPost("Hologram", "2022-06-01", action, query, nil, body, true)
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
				hologramServiceV2 := HologramServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, hologramServiceV2.HologramInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}

	var err error
	instanceId := d.Id()
	action := fmt.Sprintf("/api/v1/instances/%s/instanceName", instanceId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["instanceId"] = d.Id()
	if !d.IsNewResource() && d.HasChange("instance_name") {
		update = true
	}
	request["instanceName"] = d.Get("instance_name")
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("Hologram", "2022-06-01", action, query, nil, body, true)
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
		hologramServiceV2 := HologramServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, hologramServiceV2.HologramInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	instanceId = d.Id()
	action = fmt.Sprintf("/api/v1/instances/%s/scale", instanceId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["instanceId"] = d.Id()

	if v, ok := d.GetOk("scale_type"); ok {
		request["scaleType"] = v
	}

	if !d.IsNewResource() && d.HasChange("cpu") {
		update = true

		if v, ok := d.GetOkExists("cpu"); ok {
			request["cpu"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("storage_size") {
		update = true

		if v, ok := d.GetOkExists("storage_size"); ok {
			request["storageSize"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("cold_storage_size") {
		update = true

		if v, ok := d.GetOkExists("cold_storage_size"); ok {
			request["coldStorageSize"] = v
		}
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("Hologram", "2022-06-01", action, query, nil, body, true)
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
		hologramServiceV2 := HologramServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, hologramServiceV2.HologramInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = fmt.Sprintf("/api/v1/tag/changeResourceGroup")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["instanceId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["newResourceGroupId"] = v
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("Hologram", "2022-06-01", action, query, nil, body, true)
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
		hologramServiceV2 := HologramServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("resource_group_id"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, hologramServiceV2.HologramInstanceStateRefreshFunc(d.Id(), "ResourceGroupId", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	instanceId = d.Id()
	action = fmt.Sprintf("/api/v1/instances/%s/updateGatewayCount", instanceId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["instanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("gateway_count") {
		update = true
		if v, ok := d.GetOkExists("gateway_count"); ok {
			request["gatewayCount"] = v
		}
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("Hologram", "2022-06-01", action, query, nil, body, true)
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
		hologramServiceV2 := HologramServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, hologramServiceV2.HologramInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	instanceId = d.Id()
	action = fmt.Sprintf("/api/v1/instances/%s/network", instanceId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["instanceId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("endpoints") {
		update = true
	}
	jsonPathResult, err := jsonpath.Get("$[*].vpc_id", d.Get("endpoints").(*schema.Set).List())
	if err == nil {
		request["vpcId"] = convertListToCommaSeparate(filterEmptyStrings(jsonPathResult.([]interface{})))
	}

	if !d.IsNewResource() && d.HasChange("endpoints") {
		update = true
	}
	jsonPathResult1, err := jsonpath.Get("$[*].vswitch_id", d.Get("endpoints").(*schema.Set).List())
	if err == nil {
		request["vSwitchId"] = convertListToCommaSeparate(filterEmptyStrings(jsonPathResult1.([]interface{})))
	}

	if d.HasChange("endpoints") {
		update = true
	}
	jsonPathResult2, err := jsonpath.Get("$[*].type", d.Get("endpoints").(*schema.Set).List())
	if err == nil {
		request["networkTypes"] = convertListToCommaSeparate(filterEmptyStrings(jsonPathResult2.([]interface{})))
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("Hologram", "2022-06-01", action, query, nil, body, true)
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
		hologramServiceV2 := HologramServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, hologramServiceV2.HologramInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("tags") {
		hologramServiceV2 := HologramServiceV2{client}
		if err := hologramServiceV2.SetResourceTags(d, ""); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudHologramInstanceRead(d, meta)
}

func resourceAliCloudHologramInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	enableDelete := true
	if v, ok := d.GetOk("payment_type"); ok {
		if InArray(fmt.Sprint(v), []string{"Subscription"}) {
			enableDelete = false
			log.Printf("[WARN] Cannot destroy resource alicloud_hologram_instance which payment_type valued Subscription. Terraform will remove this resource from the state file, however resources may remain.")
		}
	}
	if enableDelete {
		client := meta.(*connectivity.AliyunClient)
		instanceId := d.Id()
		action := fmt.Sprintf("/api/v1/instances/%s/delete", instanceId)
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]*string)
		body := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		request["instanceId"] = d.Id()
		request["RegionId"] = client.RegionId

		body = request
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			response, err = client.RoaPost("Hologram", "2022-06-01", action, query, nil, body, true)

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
			if IsExpectedErrors(err, []string{"resource not exists failed"}) || NotFoundError(err) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		hologramServiceV2 := HologramServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, hologramServiceV2.HologramInstanceStateRefreshFunc(d.Id(), "$.InstanceId", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}
	return nil
}

func convertHologramInstanceInstanceInstanceChargeTypeResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "PostPaid":
		return "PayAsYouGo"
	case "PrePaid":
		return "Subscription"
	}
	return source
}

func convertHologramInstancechargeTypeRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "Subscription":
		return "PrePaid"
	case "PayAsYouGo":
		return "PostPaid"
	}
	return source
}
