// Package realtimecompute. This file is generated automatically. Please do not modify it manually, thank you!
package realtimecompute

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceRealtimeComputeNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceRealtimeComputeNamespaceCreate,
		Read:   resourceRealtimeComputeNamespaceRead,
		Update: resourceRealtimeComputeNamespaceUpdate,
		Delete: resourceRealtimeComputeNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ha": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"resource_spec": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"memory_gb": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceRealtimeComputeNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	realtimeComputeService := Service{client}
	var response map[string]interface{}
	action := "CreateNamespace"
	request := make(map[string]interface{})
	var err error

	request["Region"] = client.RegionId
	request["InstanceId"] = d.Get("instance_id")
	request["Namespace"] = d.Get("namespace")

	if v, ok := d.GetOk("ha"); ok {
		request["Ha"] = v
	}

	resourceSpec := d.Get("resource_spec")
	resourceSpecMap := map[string]interface{}{}
	for _, resourceSpecList := range convert.ConvertToInterfaceArray(resourceSpec) {
		if resourceSpecArg, ok := resourceSpecList.(map[string]interface{}); ok {
			resourceSpecMap["Cpu"] = resourceSpecArg["cpu"]
			resourceSpecMap["MemoryGB"] = resourceSpecArg["memory_gb"]
		}
	}

	resourceSpecJson, err := convert.ConvertMaptoJsonString(resourceSpecMap)
	if err != nil {
		return errs.WrapError(err)
	}

	request["ResourceSpec"] = resourceSpecJson

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("foasconsole", "2021-10-28", action, nil, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_realtime_compute_namespace", action, errs.AlibabaCloudSdkGoERROR)
	}

	stateConf := helper.BuildStateConf([]string{}, []string{"SUCCESS"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, realtimeComputeService.RealtimeComputeNamespaceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceRealtimeComputeNamespaceRead(d, meta)
}

func resourceRealtimeComputeNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	realtimeComputeService := Service{client}

	object, err := realtimeComputeService.DescribeRealtimeComputeNamespace(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("instance_id", parts[0])
	d.Set("namespace", object["Namespace"])
	d.Set("status", object["Status"])

	instance, err := realtimeComputeService.DescribeRealtimeComputeVvpInstance(parts[0])
	if err != nil {
		return errs.WrapError(err)
	}

	if fmt.Sprint(instance["ChargeType"]) == "POST" {
		if resourceSpec, ok := object["ResourceSpec"]; ok {
			resourceSpecMaps := make([]map[string]interface{}, 0)
			resourceSpecMap := make(map[string]interface{})

			if resourceSpecArg, ok := resourceSpec.(map[string]interface{}); ok {
				if cpu, ok := resourceSpecArg["Cpu"]; ok {
					resourceSpecMap["cpu"] = cpu
				}

				if memoryGB, ok := resourceSpecArg["MemoryGB"]; ok {
					resourceSpecMap["memory_gb"] = memoryGB
				}
			}

			resourceSpecMaps = append(resourceSpecMaps, resourceSpecMap)

			d.Set("resource_spec", resourceSpecMaps)
		}
	} else {
		if guaranteedResourceSpec, ok := object["GuaranteedResourceSpec"]; ok {
			guaranteedResourceSpecMaps := make([]map[string]interface{}, 0)
			guaranteedResourceSpecMap := make(map[string]interface{})

			if guaranteedResourceSpecArg, ok := guaranteedResourceSpec.(map[string]interface{}); ok {
				if cpu, ok := guaranteedResourceSpecArg["Cpu"]; ok {
					guaranteedResourceSpecMap["cpu"] = cpu
				}

				if memoryGB, ok := guaranteedResourceSpecArg["MemoryGB"]; ok {
					guaranteedResourceSpecMap["memory_gb"] = memoryGB
				}
			}

			guaranteedResourceSpecMaps = append(guaranteedResourceSpecMaps, guaranteedResourceSpecMap)

			d.Set("resource_spec", guaranteedResourceSpecMaps)
		}
	}

	return nil
}

func resourceRealtimeComputeNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	realtimeComputeService := Service{client}

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyNamespaceSpecV2"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Region"] = client.RegionId
	request["InstanceId"] = parts[0]
	request["Namespace"] = parts[1]

	if d.HasChange("resource_spec") {
		update = true
	}
	resourceSpec := d.Get("resource_spec")
	resourceSpecMap := map[string]interface{}{}
	for _, resourceSpecList := range convert.ConvertToInterfaceArray(resourceSpec) {
		if resourceSpecArg, ok := resourceSpecList.(map[string]interface{}); ok {
			resourceSpecMap["Cpu"] = resourceSpecArg["cpu"]
			resourceSpecMap["MemoryGB"] = resourceSpecArg["memory_gb"]
		}
	}

	resourceSpecJson, err := convert.ConvertMaptoJsonString(resourceSpecMap)
	if err != nil {
		return errs.WrapError(err)
	}

	zeroResourceSpecMap := map[string]interface{}{
		"Cpu":      0,
		"MemoryGB": 0,
	}

	zeroResourceSpecJson, err := convert.ConvertMaptoJsonString(zeroResourceSpecMap)
	if err != nil {
		return errs.WrapError(err)
	}

	instance, err := realtimeComputeService.DescribeRealtimeComputeVvpInstance(parts[0])
	if err != nil {
		return errs.WrapError(err)
	}

	if fmt.Sprint(instance["ChargeType"]) == "POST" {
		request["ElasticResourceSpec"] = resourceSpecJson
		request["GuaranteedResourceSpec"] = zeroResourceSpecJson
		request["Ha"] = false
	} else {
		request["GuaranteedResourceSpec"] = resourceSpecJson
		request["ElasticResourceSpec"] = zeroResourceSpecJson
	}

	if d.HasChange("ha") {
		update = true
	}
	request["Ha"] = d.Get("ha")

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("foasconsole", "2021-10-28", action, query, request, true)
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

		stateConf := helper.BuildStateConf([]string{}, []string{"SUCCESS"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, realtimeComputeService.RealtimeComputeNamespaceStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	return resourceRealtimeComputeNamespaceRead(d, meta)
}

func resourceRealtimeComputeNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	realtimeComputeService := Service{client}

	parts := strings.Split(d.Id(), ":")
	action := "DeleteNamespace"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["Namespace"] = parts[1]
	request["InstanceId"] = parts[0]
	request["Region"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("foasconsole", "2021-10-28", action, query, request, true)
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
		if errs.IsExpectedErrors(err, []string{"903042"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, realtimeComputeService.RealtimeComputeNamespaceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}
