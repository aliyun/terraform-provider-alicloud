// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudRealtimeComputeVvpInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRealtimeComputeVvpInstanceCreate,
		Read:   resourceAliCloudRealtimeComputeVvpInstanceRead,
		Update: resourceAliCloudRealtimeComputeVvpInstanceUpdate,
		Delete: resourceAliCloudRealtimeComputeVvpInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(35 * time.Minute),
			Update: schema.DefaultTimeout(35 * time.Minute),
			Delete: schema.DefaultTimeout(35 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"duration": {
				Type:     schema.TypeInt,
				Optional: true,
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
				ValidateFunc: StringInSlice([]string{"Month"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_spec": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"memory_gb": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"oss": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
			"tags": tagsSchema(),
			"vswitch_ids": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vvp_instance_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudRealtimeComputeVvpInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["CreateInstanceRequest.VpcId"] = d.Get("vpc_id")
	objectDataLocalMap := make(map[string]interface{})
	if v, ok := d.GetOk("vpc_id"); ok {
		objectDataLocalMap["VpcId"] = v
	}
	if v := d.Get("storage"); !IsNil(v) {
		storage := make(map[string]interface{})
		oss := make(map[string]interface{})
		nodeNative1, _ := jsonpath.Get("$[0].oss[0].bucket", v)
		if nodeNative1 != "" {
			oss["Bucket"] = nodeNative1
		}
		storage["Oss"] = oss
		objectDataLocalMap["Storage"] = storage
	}
	if v, ok := d.GetOk("zone_id"); ok {
		objectDataLocalMap["ZoneId"] = v
	}
	objectDataLocalMap["Region"] = client.RegionId
	if v, ok := d.GetOk("vswitch_ids"); ok {
		nodeNative4, _ := jsonpath.Get("$", v)
		if nodeNative4 != "" {
			objectDataLocalMap["VSwitchIds"] = nodeNative4
		}
	}
	if v := d.Get("resource_spec"); !IsNil(v) {
		resourceSpec := make(map[string]interface{})
		nodeNative5, _ := jsonpath.Get("$[0].cpu", v)
		if nodeNative5 != "" {
			resourceSpec["Cpu"] = nodeNative5
		}
		nodeNative6, _ := jsonpath.Get("$[0].memory_gb", v)
		if nodeNative6 != "" {
			resourceSpec["MemoryGB"] = nodeNative6
		}
		objectDataLocalMap["ResourceSpec"] = resourceSpec
	}
	if v, ok := d.GetOk("vvp_instance_name"); ok {
		objectDataLocalMap["InstanceName"] = v
	}
	if v, ok := d.GetOk("pricing_cycle"); ok {
		objectDataLocalMap["PricingCycle"] = v
	}
	if v, ok := d.GetOk("duration"); ok {
		objectDataLocalMap["Duration"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		objectDataLocalMap["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("payment_type"); ok {
		objectDataLocalMap["ChargeType"] = convertRealtimeComputeCreateInstanceRequestChargeTypeRequest(v)
	}
	request["CreateInstanceRequest"] = objectDataLocalMap
	request["CreateInstanceRequest.ZoneId"] = d.Get("zone_id")
	request["CreateInstanceRequest.InstanceName"] = d.Get("vvp_instance_name")
	if v, ok := d.GetOk("pricing_cycle"); ok {
		request["CreateInstanceRequest.PricingCycle"] = v
	}
	if v, ok := d.GetOk("duration"); ok {
		request["CreateInstanceRequest.Duration"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["CreateInstanceRequest.ResourceGroupId"] = v
	}
	request["CreateInstanceRequest.ChargeType"] = convertRealtimeComputeCreateInstanceRequestChargeTypeRequest(d.Get("payment_type").(string))
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("foasconsole", "2019-06-01", action, query, request, true)

		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"998001", "Invoke ABM failed"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_realtime_compute_vvp_instance", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.OrderInfo.InstanceId", response)
	d.SetId(fmt.Sprint(id))

	realtimeComputeServiceV2 := RealtimeComputeServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, realtimeComputeServiceV2.RealtimeComputeVvpInstanceStateRefreshFunc(d.Id(), "ClusterStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudRealtimeComputeVvpInstanceRead(d, meta)
}

func resourceAliCloudRealtimeComputeVvpInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	realtimeComputeServiceV2 := RealtimeComputeServiceV2{client}

	objectRaw, err := realtimeComputeServiceV2.DescribeRealtimeComputeVvpInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_realtime_compute_vvp_instance DescribeRealtimeComputeVvpInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["ResourceCreateTime"])
	d.Set("payment_type", convertRealtimeComputeInstancesChargeTypeResponse(objectRaw["ChargeType"]))
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["ClusterStatus"])
	d.Set("vpc_id", objectRaw["VpcId"])
	d.Set("vvp_instance_name", objectRaw["InstanceName"])
	d.Set("resource_id", objectRaw["ResourceId"])
	resourceSpecMaps := make([]map[string]interface{}, 0)
	resourceSpecMap := make(map[string]interface{})
	resourceSpec1Raw := make(map[string]interface{})
	if objectRaw["ResourceSpec"] != nil {
		resourceSpec1Raw = objectRaw["ResourceSpec"].(map[string]interface{})
	}
	if len(resourceSpec1Raw) > 0 {
		resourceSpecMap["cpu"] = resourceSpec1Raw["Cpu"]
		resourceSpecMap["memory_gb"] = resourceSpec1Raw["MemoryGB"]
		resourceSpecMaps = append(resourceSpecMaps, resourceSpecMap)
	}
	d.Set("resource_spec", resourceSpecMaps)
	storageMaps := make([]map[string]interface{}, 0)
	storageMap := make(map[string]interface{})
	storage1Raw := make(map[string]interface{})
	if objectRaw["Storage"] != nil {
		storage1Raw = objectRaw["Storage"].(map[string]interface{})
	}
	if len(storage1Raw) > 0 {
		storageMap["oss"] = storage1Raw["Oss"]
		ossMaps := make([]map[string]interface{}, 0)
		ossMap := make(map[string]interface{})
		oss3Raw := make(map[string]interface{})
		if storage1Raw["Oss"] != nil {
			oss3Raw = storage1Raw["Oss"].(map[string]interface{})
		}
		if len(oss3Raw) > 0 {
			ossMap["bucket"] = oss3Raw["Bucket"]
			ossMaps = append(ossMaps, ossMap)
		}
		storageMap["oss"] = ossMaps
		storageMaps = append(storageMaps, storageMap)
	}
	d.Set("storage", storageMaps)
	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))
	vSwitchIds1Raw := make([]interface{}, 0)
	if objectRaw["VSwitchIds"] != nil {
		vSwitchIds1Raw = objectRaw["VSwitchIds"].([]interface{})
	}

	d.Set("vswitch_ids", vSwitchIds1Raw)

	return nil
}

func resourceAliCloudRealtimeComputeVvpInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	action := "ModifyPrepayInstanceSpec"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ModifyPrepayInstanceSpecRequest.InstanceId"] = d.Id()
	objectDataLocalMap := make(map[string]interface{})
	if d.HasChange("resource_spec") {
		update = true
		if v := d.Get("resource_spec"); !IsNil(v) {
			resourceSpec := make(map[string]interface{})
			nodeNative, _ := jsonpath.Get("$[0].cpu", v)
			if nodeNative != "" {
				resourceSpec["Cpu"] = nodeNative
			}
			nodeNative1, _ := jsonpath.Get("$[0].memory_gb", v)
			if nodeNative1 != "" {
				resourceSpec["MemoryGB"] = nodeNative1
			}
			objectDataLocalMap["ResourceSpec"] = resourceSpec
		}
	}
	if d.HasChange("region_id") {
		update = true
	}
	objectDataLocalMap["Region"] = client.RegionId
	request["ModifyPrepayInstanceSpecRequest"] = objectDataLocalMap
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("foasconsole", "2019-06-01", action, query, request, true)

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
		realtimeComputeServiceV2 := RealtimeComputeServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutUpdate), 2*time.Minute, realtimeComputeServiceV2.RealtimeComputeVvpInstanceStateRefreshFunc(d.Id(), "ClusterStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "ChangeResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ResourceId"] = d.Id()
	query["ResourceRegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["ResourceGroupId"] = d.Get("resource_group_id")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("foasconsole", "2019-06-01", action, query, request, true)

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
		realtimeComputeServiceV2 := RealtimeComputeServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("resource_group_id"))}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, realtimeComputeServiceV2.RealtimeComputeVvpInstanceStateRefreshFunc(d.Id(), "ResourceGroupId", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("resource_group_id")
	}

	if d.HasChange("tags") {
		realtimeComputeServiceV2 := RealtimeComputeServiceV2{client}
		if err := realtimeComputeServiceV2.SetResourceTags(d, "vvpinstance"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	d.Partial(false)
	return resourceAliCloudRealtimeComputeVvpInstanceRead(d, meta)
}

func resourceAliCloudRealtimeComputeVvpInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	if v, ok := d.GetOk("payment_type"); ok {
		if v == "Subscription" {
			log.Printf("[WARN] Cannot destroy resource alicloud_realtime_compute_vvp_instance which payment_type valued Subscription. Terraform will remove this resource from the state file, however resources may remain.")
			return nil
		}
	}
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["DeleteInstanceRequest.InstanceId"] = d.Id()
	request["DeleteInstanceRequest.Region"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("foasconsole", "2019-06-01", action, query, request, true)

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

	realtimeComputeServiceV2 := RealtimeComputeServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 9*time.Minute, realtimeComputeServiceV2.RealtimeComputeVvpInstanceStateRefreshFunc(d.Id(), "InstanceId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func convertRealtimeComputeInstancesChargeTypeResponse(source interface{}) interface{} {
	switch source {
	case "POST":
		return "PayAsYouGo"
	case "PRE":
		return "Subscription"
	}
	return source
}
func convertRealtimeComputeCreateInstanceRequestChargeTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "POST"
	case "Subscription":
		return "PRE"
	}
	return source
}
