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

func resourceAliCloudEnsGatewayQos() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEnsGatewayQosCreate,
		Read:   resourceAliCloudEnsGatewayQosRead,
		Update: resourceAliCloudEnsGatewayQosUpdate,
		Delete: resourceAliCloudEnsGatewayQosDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bandwidth_in": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"bandwidth_out": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ens_region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway_qos_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gateway_qos_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instances": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudEnsGatewayQosCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateEnsGatewayQos"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["GatewayQosType"] = d.Get("gateway_qos_type")
	if v, ok := d.GetOkExists("bandwidth_in"); ok {
		request["BandwidthIn"] = v
	}
	if v, ok := d.GetOk("instances"); ok && v.(*schema.Set).Len() > 0 {
		instancesMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["InstanceId"] = dataLoopTmp["instance_id"]
			dataLoopMap["InstanceType"] = dataLoopTmp["instance_type"]
			instancesMapsArray = append(instancesMapsArray, dataLoopMap)
		}
		request["Instances"] = instancesMapsArray
	}

	if v, ok := d.GetOk("gateway_qos_name"); ok {
		request["GatewayQosName"] = v
	}
	if v, ok := d.GetOkExists("bandwidth_out"); ok {
		request["BandwidthOut"] = v
	}
	request["NetworkId"] = d.Get("network_id")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ens_gateway_qos", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["GatewayQosId"]))

	return resourceAliCloudEnsGatewayQosRead(d, meta)
}

func resourceAliCloudEnsGatewayQosRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ensServiceV2 := EnsServiceV2{client}

	objectRaw, err := ensServiceV2.DescribeEnsGatewayQos(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ens_gateway_qos DescribeEnsGatewayQos Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("bandwidth_in", objectRaw["BandwidthIn"])
	d.Set("bandwidth_out", objectRaw["BandwidthOut"])
	d.Set("creation_time", objectRaw["CreationTime"])
	d.Set("ens_region_id", objectRaw["EnsRegionId"])
	d.Set("gateway_qos_name", objectRaw["GatewayQosName"])
	d.Set("gateway_qos_type", objectRaw["GatewayQosType"])
	d.Set("network_id", objectRaw["NetworkId"])
	d.Set("status", objectRaw["Status"])

	instanceRaw, _ := jsonpath.Get("$.Instances.Instance", objectRaw)
	instancesMaps := make([]map[string]interface{}, 0)
	if instanceRaw != nil {
		for _, instanceChildRaw := range convertToInterfaceArray(instanceRaw) {
			instancesMap := make(map[string]interface{})
			instanceChildRaw := instanceChildRaw.(map[string]interface{})
			instancesMap["instance_id"] = instanceChildRaw["InstanceId"]
			instancesMap["instance_type"] = instanceChildRaw["InstanceType"]
			instancesMap["status"] = instanceChildRaw["Status"]

			instancesMaps = append(instancesMaps, instancesMap)
		}
	}
	if err := d.Set("instances", instancesMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudEnsGatewayQosUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateEnsGatewayQos"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["GatewayQosId"] = d.Id()

	if d.HasChange("bandwidth_in") {
		update = true
		request["BandwidthIn"] = d.Get("bandwidth_in")
	}

	if d.HasChange("gateway_qos_name") {
		update = true
		request["GatewayQosName"] = d.Get("gateway_qos_name")
	}

	if d.HasChange("bandwidth_out") {
		update = true
		request["BandwidthOut"] = d.Get("bandwidth_out")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
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
		ensServiceV2 := EnsServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ensServiceV2.EnsGatewayQosStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("instances") {
		var err error
		oldEntry, newEntry := d.GetChange("instances")
		oldEntrySet := oldEntry.(*schema.Set)
		newEntrySet := newEntry.(*schema.Set)
		removed := oldEntrySet.Difference(newEntrySet)
		added := newEntrySet.Difference(oldEntrySet)

		if removed.Len() > 0 {
			action := "DeleteEnsGatewayQosInstances"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["GatewayQosId"] = d.Id()

			localData := removed.List()
			instancesMapsArray := make([]interface{}, 0)
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["InstanceId"] = dataLoopTmp["instance_id"]
				dataLoopMap["InstanceType"] = dataLoopTmp["instance_type"]
				instancesMapsArray = append(instancesMapsArray, dataLoopMap)
			}
			request["Instances"] = instancesMapsArray

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
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
			ensServiceV2 := EnsServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"[]"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ensServiceV2.EnsGatewayQosStateRefreshFunc(d.Id(), "", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}

		if added.Len() > 0 {
			action := "AddEnsGatewayQosInstances"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["GatewayQosId"] = d.Id()

			localData := added.List()
			instancesMapsArray := make([]interface{}, 0)
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["InstanceId"] = dataLoopTmp["instance_id"]
				dataLoopMap["InstanceType"] = dataLoopTmp["instance_type"]
				instancesMapsArray = append(instancesMapsArray, dataLoopMap)
			}
			request["Instances"] = instancesMapsArray

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
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

	}
	return resourceAliCloudEnsGatewayQosRead(d, meta)
}

func resourceAliCloudEnsGatewayQosDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteEnsGatewayQos"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["GatewayQosId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
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
		if IsExpectedErrors(err, []string{"GatewayQosNotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	ensServiceV2 := EnsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ensServiceV2.EnsGatewayQosStateRefreshFunc(d.Id(), "$.GatewayQosId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
