// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"time"
)

func resourceAliCloudEfloResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEfloResourceCreate,
		Read:   resourceAliCloudEfloResourceRead,
		Update: resourceAliCloudEfloResourceUpdate,
		Delete: resourceAliCloudEfloResourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cluster_desc": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"machine_types": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"memory_info": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"bond_num": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"node_count": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"cpu_info": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"network_info": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"gpu_info": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"disk_info": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"network_mode": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"user_access_param": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint": {
							Type:     schema.TypeString,
							Required: true,
						},
						"access_id": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"workspace_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"access_key": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								if old != "" && new != "" && old != new {
									return true
								}
								return false
							},
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudEfloResourceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateResource"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		request["ClusterId"] = v
	}
	request["RegionId"] = client.RegionId

	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("machine_types"); v != nil {
		name1, _ := jsonpath.Get("$[0].name", v)
		if name1 != nil && name1 != "" {
			objectDataLocalMap["Name"] = name1
		}
		networkInfo1, _ := jsonpath.Get("$[0].network_info", v)
		if networkInfo1 != nil && networkInfo1 != "" {
			objectDataLocalMap["NetworkInfo"] = networkInfo1
		}
		gpuInfo1, _ := jsonpath.Get("$[0].gpu_info", v)
		if gpuInfo1 != nil && gpuInfo1 != "" {
			objectDataLocalMap["GpuInfo"] = gpuInfo1
		}
		diskInfo1, _ := jsonpath.Get("$[0].disk_info", v)
		if diskInfo1 != nil && diskInfo1 != "" {
			objectDataLocalMap["DiskInfo"] = diskInfo1
		}
		nodeCount1, _ := jsonpath.Get("$[0].node_count", v)
		if nodeCount1 != nil && nodeCount1 != "" {
			objectDataLocalMap["NodeCount"] = nodeCount1
		}
		type1, _ := jsonpath.Get("$[0].type", v)
		if type1 != nil && type1 != "" {
			objectDataLocalMap["Type"] = type1
		}
		networkMode1, _ := jsonpath.Get("$[0].network_mode", v)
		if networkMode1 != nil && networkMode1 != "" {
			objectDataLocalMap["NetworkMode"] = networkMode1
		}
		cpuInfo1, _ := jsonpath.Get("$[0].cpu_info", v)
		if cpuInfo1 != nil && cpuInfo1 != "" {
			objectDataLocalMap["CpuInfo"] = cpuInfo1
		}
		bondNum1, _ := jsonpath.Get("$[0].bond_num", v)
		if bondNum1 != nil && bondNum1 != "" {
			objectDataLocalMap["BondNum"] = bondNum1
		}
		memoryInfo1, _ := jsonpath.Get("$[0].memory_info", v)
		if memoryInfo1 != nil && memoryInfo1 != "" {
			objectDataLocalMap["MemoryInfo"] = memoryInfo1
		}

		objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
		if err != nil {
			return WrapError(err)
		}
		request["MachineTypes"] = string(objectDataLocalMapJson)
	}

	objectDataLocalMap1 := make(map[string]interface{})

	if v := d.Get("user_access_param"); v != nil {
		accessId1, _ := jsonpath.Get("$[0].access_id", v)
		if accessId1 != nil && accessId1 != "" {
			objectDataLocalMap1["AccessId"] = accessId1
		}
		workspaceId1, _ := jsonpath.Get("$[0].workspace_id", v)
		if workspaceId1 != nil && workspaceId1 != "" {
			objectDataLocalMap1["WorkspaceId"] = workspaceId1
		}
		accessKey1, _ := jsonpath.Get("$[0].access_key", v)
		if accessKey1 != nil && accessKey1 != "" {
			objectDataLocalMap1["AccessKey"] = accessKey1
		}
		endpoint1, _ := jsonpath.Get("$[0].endpoint", v)
		if endpoint1 != nil && endpoint1 != "" {
			objectDataLocalMap1["Endpoint"] = endpoint1
		}

		objectDataLocalMap1Json, err := json.Marshal(objectDataLocalMap1)
		if err != nil {
			return WrapError(err)
		}
		request["UserAccessParam"] = string(objectDataLocalMap1Json)
	}

	request["ClusterName"] = d.Get("cluster_name")
	if v, ok := d.GetOk("cluster_desc"); ok {
		request["ClusterDesc"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("eflo-cnp", "2023-08-28", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_eflo_resource", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["ClusterId"]))

	return resourceAliCloudEfloResourceRead(d, meta)
}

func resourceAliCloudEfloResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	efloServiceV2 := EfloServiceV2{client}

	objectRaw, err := efloServiceV2.DescribeEfloResource(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_eflo_resource DescribeEfloResource Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cluster_desc", objectRaw["ClusterDesc"])
	d.Set("cluster_name", objectRaw["ClusterName"])
	d.Set("cluster_id", objectRaw["ClusterId"])
	d.Set("resource_id", objectRaw["ResourceId"])

	machineTypesMaps := make([]map[string]interface{}, 0)
	machineTypesMap := make(map[string]interface{})
	machineTypeRaw := make(map[string]interface{})
	if objectRaw["MachineType"] != nil {
		machineTypeRaw = objectRaw["MachineType"].(map[string]interface{})
	}
	if len(machineTypeRaw) > 0 {
		machineTypesMap["bond_num"] = machineTypeRaw["BondNum"]
		machineTypesMap["cpu_info"] = machineTypeRaw["CpuInfo"]
		machineTypesMap["disk_info"] = machineTypeRaw["DiskInfo"]
		machineTypesMap["gpu_info"] = machineTypeRaw["GpuInfo"]
		machineTypesMap["memory_info"] = machineTypeRaw["MemoryInfo"]
		machineTypesMap["name"] = machineTypeRaw["Name"]
		machineTypesMap["network_info"] = machineTypeRaw["NetworkInfo"]
		machineTypesMap["network_mode"] = machineTypeRaw["NetworkMode"]
		machineTypesMap["node_count"] = machineTypeRaw["NodeCount"]
		machineTypesMap["type"] = machineTypeRaw["Type"]

		machineTypesMaps = append(machineTypesMaps, machineTypesMap)
	}
	if err := d.Set("machine_types", machineTypesMaps); err != nil {
		return err
	}
	userAccessParamMaps := make([]map[string]interface{}, 0)
	userAccessParamMap := make(map[string]interface{})
	userAccessParamRaw := make(map[string]interface{})
	if objectRaw["UserAccessParam"] != nil {
		userAccessParamRaw = objectRaw["UserAccessParam"].(map[string]interface{})
	}
	if len(userAccessParamRaw) > 0 {
		userAccessParamMap["access_id"] = userAccessParamRaw["AccessId"]
		userAccessParamMap["access_key"] = userAccessParamRaw["AccessKey"]
		userAccessParamMap["endpoint"] = userAccessParamRaw["Endpoint"]
		userAccessParamMap["workspace_id"] = userAccessParamRaw["WorkspaceId"]

		userAccessParamMaps = append(userAccessParamMaps, userAccessParamMap)
	}
	if err := d.Set("user_access_param", userAccessParamMaps); err != nil {
		return err
	}

	d.Set("cluster_id", d.Id())

	return nil
}

func resourceAliCloudEfloResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ValidateResource"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ClusterId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("user_access_param") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("user_access_param"); v != nil {
		accessId1, _ := jsonpath.Get("$[0].access_id", v)
		if accessId1 != nil && (d.HasChange("user_access_param.0.access_id") || accessId1 != "") {
			objectDataLocalMap["AccessId"] = accessId1
		}
		accessKey1, _ := jsonpath.Get("$[0].access_key", v)
		if accessKey1 != nil && (d.HasChange("user_access_param.0.access_key") || accessKey1 != "") {
			objectDataLocalMap["AccessKey"] = accessKey1
		}
		workspaceId1, _ := jsonpath.Get("$[0].workspace_id", v)
		if workspaceId1 != nil && (d.HasChange("user_access_param.0.workspace_id") || workspaceId1 != "") {
			objectDataLocalMap["WorkspaceId"] = workspaceId1
		}
		endpoint1, _ := jsonpath.Get("$[0].endpoint", v)
		if endpoint1 != nil && (d.HasChange("user_access_param.0.endpoint") || endpoint1 != "") {
			objectDataLocalMap["Endpoint"] = endpoint1
		}

		objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
		if err != nil {
			return WrapError(err)
		}
		request["UserAccessParam"] = string(objectDataLocalMapJson)
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("eflo-cnp", "2023-08-28", action, query, request, true)
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

	return resourceAliCloudEfloResourceRead(d, meta)
}

func resourceAliCloudEfloResourceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Resource. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
