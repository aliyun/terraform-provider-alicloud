// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudRealtimeComputeVvpNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRealtimeComputeVvpNamespaceCreate,
		Read:   resourceAliCloudRealtimeComputeVvpNamespaceRead,
		Update: resourceAliCloudRealtimeComputeVvpNamespaceUpdate,
		Delete: resourceAliCloudRealtimeComputeVvpNamespaceDelete,
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

func resourceAliCloudRealtimeComputeVvpNamespaceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateNamespace"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("namespace"); ok {
		request["Namespace"] = v
	}
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	request["Region"] = client.RegionId

	dataList := make(map[string]interface{})

	if v := d.Get("resource_spec"); !IsNil(v) {
		cpu1, _ := jsonpath.Get("$[0].cpu", v)
		if cpu1 != nil && cpu1 != "" {
			dataList["Cpu"] = cpu1
		}
		memoryGb, _ := jsonpath.Get("$[0].memory_gb", v)
		if memoryGb != nil && memoryGb != "" {
			dataList["MemoryGB"] = memoryGb
		}

		dataListJson, err := json.Marshal(dataList)
		if err != nil {
			return WrapError(err)
		}
		request["ResourceSpec"] = string(dataListJson)
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("foasconsole", "2021-10-28", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_realtime_compute_vvp_namespace", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["InstanceId"], request["Namespace"]))

	realtimeComputeServiceV2 := RealtimeComputeServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"SUCCESS"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, realtimeComputeServiceV2.RealtimeComputeVvpNamespaceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudRealtimeComputeVvpNamespaceRead(d, meta)
}

func resourceAliCloudRealtimeComputeVvpNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	realtimeComputeServiceV2 := RealtimeComputeServiceV2{client}

	objectRaw, err := realtimeComputeServiceV2.DescribeRealtimeComputeVvpNamespace(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_realtime_compute_vvp_namespace DescribeRealtimeComputeVvpNamespace Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("status", objectRaw["Status"])
	d.Set("namespace", objectRaw["Namespace"])

	resourceSpecMaps := make([]map[string]interface{}, 0)
	resourceSpecMap := make(map[string]interface{})
	resourceSpecRaw := make(map[string]interface{})
	if objectRaw["ResourceSpec"] != nil {
		resourceSpecRaw = objectRaw["ResourceSpec"].(map[string]interface{})
	}
	if len(resourceSpecRaw) > 0 {
		resourceSpecMap["cpu"] = resourceSpecRaw["Cpu"]
		resourceSpecMap["memory_gb"] = resourceSpecRaw["MemoryGB"]

		resourceSpecMaps = append(resourceSpecMaps, resourceSpecMap)
	}
	if err := d.Set("resource_spec", resourceSpecMaps); err != nil {
		return err
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("instance_id", parts[0])

	return nil
}

func resourceAliCloudRealtimeComputeVvpNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyNamespaceSpec"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Namespace"] = parts[1]
	request["InstanceId"] = parts[0]
	request["Region"] = client.RegionId
	if d.HasChange("resource_spec") {
		update = true
	}
	dataList := make(map[string]interface{})

	if v := d.Get("resource_spec"); v != nil {
		cpu1, _ := jsonpath.Get("$[0].cpu", v)
		if cpu1 != nil && (d.HasChange("resource_spec.0.cpu") || cpu1 != "") {
			dataList["Cpu"] = cpu1
		}
		memoryGb, _ := jsonpath.Get("$[0].memory_gb", v)
		if memoryGb != nil && (d.HasChange("resource_spec.0.memory_gb") || memoryGb != "") {
			dataList["MemoryGB"] = memoryGb
		}

		dataListJson, err := json.Marshal(dataList)
		if err != nil {
			return WrapError(err)
		}
		request["ResourceSpec"] = string(dataListJson)
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("foasconsole", "2021-10-28", action, query, request, true)
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
		realtimeComputeServiceV2 := RealtimeComputeServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"SUCCESS"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, realtimeComputeServiceV2.RealtimeComputeVvpNamespaceStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudRealtimeComputeVvpNamespaceRead(d, meta)
}

func resourceAliCloudRealtimeComputeVvpNamespaceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
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

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("foasconsole", "2021-10-28", action, query, request, true)
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
		if IsExpectedErrors(err, []string{"903042"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	realtimeComputeServiceV2 := RealtimeComputeServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 30*time.Second, realtimeComputeServiceV2.RealtimeComputeVvpNamespaceStateRefreshFunc(d.Id(), "Namespace", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
