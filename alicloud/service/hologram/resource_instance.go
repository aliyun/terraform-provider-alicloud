// Package hologram. This file is generated automatically. Please do not modify it manually, thank you!
package hologram

import (
	"fmt"
	"log"
	"regexp"
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

func ResourceInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceInstanceCreate,
		Read:   resourceInstanceRead,
		Update: resourceInstanceUpdate,
		Delete: resourceInstanceDelete,
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
			"enable_ssl": {
				Type:     schema.TypeBool,
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
				ValidateFunc: validators.StringMatch(regexp.MustCompile("^.{2,64}$"), "The name of the resource"),
			},
			"instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"Standard", "Follower", "Warehouse", "Shared", "Serverless"}, false),
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
				ValidateFunc: validators.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"pricing_cycle": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"Month", "Hour"}, false),
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
				ValidateFunc: validators.StringInSlice([]string{"DOWNGRADE", "UPGRADE"}, false),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"Running", "Allocating", "Suspended", "Creating"}, false),
			},
			"storage_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"tags": tags.TagsSchema(),
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceInstanceCreate(d *schema.ResourceData, meta interface{}) error {

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
		request["vSwitchId"] = convert.ConvertListToCommaSeparate(convert.FilterEmptyStrings(jsonPathResult2.([]interface{})))
	}

	if v, ok := d.GetOkExists("gateway_count"); ok {
		request["gatewayCount"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["resourceGroupId"] = v
	}
	request["chargeType"] = convertInstancechargeTypeRequest(d.Get("payment_type").(string))
	request["instanceName"] = d.Get("instance_name")
	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["autoPay"] = v
	}
	if v, ok := d.GetOk("pricing_cycle"); ok {
		request["pricingCycle"] = v
	}
	jsonPathResult9, err := jsonpath.Get("$[*].vpc_id", d.Get("endpoints").(*schema.Set).List())
	if err == nil {
		request["vpcId"] = convert.ConvertListToCommaSeparate(convert.FilterEmptyStrings(jsonPathResult9.([]interface{})))
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
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RoaPost("Hologram", "2022-06-01", action, query, nil, body, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_hologram_instance", action, errs.AlibabaCloudSdkGoERROR)
	}

	code, _ := jsonpath.Get("$.Data.Success", response)
	if fmt.Sprint(code) != "true" {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_hologram_instance", action, errs.AlibabaCloudSdkGoERROR, response)
	}

	id, _ := jsonpath.Get("$.Data.InstanceId", response)
	d.SetId(fmt.Sprint(id))

	service := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, service.InstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceInstanceUpdate(d, meta)
}

func resourceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	service := NewService(client)

	objectRaw, err := service.DescribeInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_hologram_instance DescribeInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("cold_storage_size", objectRaw["ColdStorage"])
	d.Set("cpu", objectRaw["Cpu"])
	d.Set("create_time", objectRaw["CreationTime"])
	d.Set("enable_ssl", objectRaw["EnableSSL"])
	d.Set("gateway_count", objectRaw["GatewayCount"])
	d.Set("instance_name", objectRaw["InstanceName"])
	d.Set("instance_type", objectRaw["InstanceType"])
	d.Set("leader_instance_id", objectRaw["LeaderInstanceId"])
	d.Set("payment_type", convertInstanceInstanceInstanceChargeTypeResponse(objectRaw["InstanceChargeType"]))
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["InstanceStatus"])
	d.Set("storage_size", objectRaw["Disk"])
	d.Set("zone_id", objectRaw["ZoneId"])

	endpointsRaw := objectRaw["Endpoints"]
	endpointsMaps := make([]map[string]interface{}, 0)
	if endpointsRaw != nil {
		for _, endpointsChildRaw := range convert.ConvertToInterfaceArray(endpointsRaw) {
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
	d.Set("tags", tags.TagsToMap(tagsMaps))

	return nil
}

func resourceInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var header map[string]*string
	var query map[string]*string
	var body map[string]interface{}
	update := false
	d.Partial(true)

	service := NewService(client)
	objectRaw, _ := service.DescribeInstance(d.Id())

	if d.HasChange("status") {
		var err error
		target := d.Get("status").(string)

		currentStatus, err := jsonpath.Get("InstanceStatus", objectRaw)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, d.Id(), "InstanceStatus", objectRaw)
		}
		if fmt.Sprint(currentStatus) != target {
			if target == "Suspended" {
				instanceId := d.Id()
				action := fmt.Sprintf("/api/v1/instances/%s/stop", instanceId)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})
				request["instanceId"] = d.Id()

				body = request
				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RoaPost("Hologram", "2022-06-01", action, query, nil, body, true)
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
				service := NewService(client)
				stateConf := helper.BuildStateConf([]string{}, []string{"Suspended"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, service.InstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
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
				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RoaPost("Hologram", "2022-06-01", action, query, nil, body, true)
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
				service := NewService(client)
				stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, service.InstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
				}

			}
		}
	}

	if d.HasChange("enable_ssl") {
		var err error
		target := d.Get("enable_ssl").(bool)

		currentStatus, err := jsonpath.Get("EnableSSL", objectRaw)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, d.Id(), "EnableSSL", objectRaw)
		}
		currentStatusValue, formatErr := convert.FormatBool(currentStatus)
		if formatErr != nil {
			return formatErr
		}
		if currentStatusValue != target {
			if target == true {
				instanceId := d.Id()
				action := fmt.Sprintf("/api/v1/instances/%s/enableSSL", instanceId)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})

				body = request
				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RoaPost("Hologram", "2022-06-01", action, query, header, body, true)
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
			if target == false {
				instanceId := d.Id()
				action := fmt.Sprintf("/api/v1/instances/%s/disableSSL", instanceId)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})

				body = request
				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RoaPost("Hologram", "2022-06-01", action, query, header, body, true)
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
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RoaPost("Hologram", "2022-06-01", action, query, nil, body, true)
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
		service := NewService(client)

		stateConf := helper.BuildStateConf([]string{}, []string{fmt.Sprint(request["instanceName"])}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, service.InstanceStateRefreshFunc(d.Id(), "InstanceName", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}

		stateConf = helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, service.InstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
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
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RoaPost("Hologram", "2022-06-01", action, query, nil, body, true)
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
		service := NewService(client)

		if d.HasChange("cpu") {
			stateConf := helper.BuildStateConf([]string{}, []string{fmt.Sprint(request["cpu"])}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, service.InstanceStateRefreshFunc(d.Id(), "Cpu", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}
		}

		stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, service.InstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
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
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RoaPost("Hologram", "2022-06-01", action, query, nil, body, true)
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
		service := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("resource_group_id"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, service.InstanceStateRefreshFunc(d.Id(), "ResourceGroupId", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
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
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RoaPost("Hologram", "2022-06-01", action, query, nil, body, true)
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
		service := NewService(client)

		stateConf := helper.BuildStateConf([]string{}, []string{fmt.Sprint(request["gatewayCount"])}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, service.InstanceStateRefreshFunc(d.Id(), "GatewayCount", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}

		stateConf = helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, service.InstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
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
		request["vpcId"] = convert.ConvertListToCommaSeparate(convert.FilterEmptyStrings(jsonPathResult.([]interface{})))
	}

	if !d.IsNewResource() && d.HasChange("endpoints") {
		update = true
	}
	jsonPathResult1, err := jsonpath.Get("$[*].vswitch_id", d.Get("endpoints").(*schema.Set).List())
	if err == nil {
		request["vSwitchId"] = convert.ConvertListToCommaSeparate(convert.FilterEmptyStrings(jsonPathResult1.([]interface{})))
	}

	if d.HasChange("endpoints") {
		update = true
	}
	jsonPathResult2, err := jsonpath.Get("$[*].type", d.Get("endpoints").(*schema.Set).List())
	if err == nil {
		request["networkTypes"] = convert.ConvertListToCommaSeparate(convert.FilterEmptyStrings(jsonPathResult2.([]interface{})))
	}
	body = request
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RoaPost("Hologram", "2022-06-01", action, query, nil, body, true)
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
		service := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, service.InstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	if d.HasChange("tags") {
		service := NewService(client)
		if err := service.SetResourceTags(d, ""); err != nil {
			return errs.WrapError(err)
		}
	}
	d.Partial(false)
	return resourceInstanceRead(d, meta)
}

func resourceInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	enableDelete := true
	if v, ok := d.GetOk("payment_type"); ok {
		if helper.InArray(fmt.Sprint(v), []string{"Subscription"}) {
			return errs.WrapError(errs.Error("'Subscription' instance cannot be destroyed: please unsubscribe it in the console or wait for it to expire. To remove it from the state only, run 'terraform state rm'."))
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
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
			response, err = client.RoaPost("Hologram", "2022-06-01", action, query, nil, body, true)

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
			if errs.IsExpectedErrors(err, []string{"resource not exists failed"}) || errs.NotFoundError(err) {
				return nil
			}
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		service := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, service.InstanceStateRefreshFunc(d.Id(), "$.InstanceId", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}

	}
	return nil
}

func convertInstanceInstanceInstanceChargeTypeResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "PostPaid":
		return "PayAsYouGo"
	case "PrePaid":
		return "Subscription"
	}
	return source
}

func convertInstancechargeTypeRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "Subscription":
		return "PrePaid"
	case "PayAsYouGo":
		return "PostPaid"
	}
	return source
}
