// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudFcv2Function() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudFcv2FunctionCreate,
		Read:   resourceAliCloudFcv2FunctionRead,
		Update: resourceAliCloudFcv2FunctionUpdate,
		Delete: resourceAliCloudFcv2FunctionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"ca_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"code": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zip_file": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"oss_object_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"oss_bucket_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"code_checksum": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cpu": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_container_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"args": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"command": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"acceleration_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"web_server_mode": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"image": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"custom_dns": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"searches": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"dns_options": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"name_servers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"custom_health_check_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"initial_delay_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"timeout_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"http_get_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"period_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"failure_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"success_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"custom_runtime_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"args": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"command": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disk_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"environment_variables": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"function_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"gpu_memory_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"handler": {
				Type:     schema.TypeString,
				Required: true,
			},
			"initialization_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"initializer": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_concurrency": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"instance_lifecycle_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pre_stop": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"timeout": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"handler": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"pre_freeze": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"timeout": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"handler": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"layers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"memory_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"runtime": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"function_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudFcv2FunctionCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	serviceName := d.Get("service_name")
	action := fmt.Sprintf("/2021-04-06/services/%s/functions", serviceName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	conn, err := client.NewFcv2Client()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	objectDataLocalMap := make(map[string]interface{})
	if v, ok := d.GetOk("description"); ok {
		objectDataLocalMap["description"] = v
	}
	if v, ok := d.GetOk("function_name"); ok {
		objectDataLocalMap["functionName"] = v
	}
	if v, ok := d.GetOk("handler"); ok {
		objectDataLocalMap["handler"] = v
	}
	if v, ok := d.GetOk("initialization_timeout"); ok {
		objectDataLocalMap["initializationTimeout"] = v
	}
	if v, ok := d.GetOk("initializer"); ok {
		objectDataLocalMap["initializer"] = v
	}
	if v, ok := d.GetOk("memory_size"); ok {
		objectDataLocalMap["memorySize"] = v
	}
	if v, ok := d.GetOk("runtime"); ok {
		objectDataLocalMap["runtime"] = v
	}
	if v, ok := d.GetOk("timeout"); ok {
		objectDataLocalMap["timeout"] = v
	}
	if v := d.Get("code"); !IsNil(v) {
		code := make(map[string]interface{})
		nodeNative8, _ := jsonpath.Get("$[0].oss_bucket_name", v)
		if nodeNative8 != "" {
			code["ossBucketName"] = nodeNative8
		}
		nodeNative9, _ := jsonpath.Get("$[0].oss_object_name", v)
		if nodeNative9 != "" {
			code["ossObjectName"] = nodeNative9
		}
		nodeNative10, _ := jsonpath.Get("$[0].zip_file", v)
		if nodeNative10 != "" {
			code["zipFile"] = nodeNative10
		}
		objectDataLocalMap["code"] = code
	}
	if v := d.Get("custom_container_config"); !IsNil(v) {
		customContainerConfig := make(map[string]interface{})
		nodeNative11, _ := jsonpath.Get("$[0].args", v)
		if nodeNative11 != "" {
			customContainerConfig["args"] = nodeNative11
		}
		nodeNative12, _ := jsonpath.Get("$[0].command", v)
		if nodeNative12 != "" {
			customContainerConfig["command"] = nodeNative12
		}
		nodeNative13, _ := jsonpath.Get("$[0].image", v)
		if nodeNative13 != "" {
			customContainerConfig["image"] = nodeNative13
		}
		nodeNative14, _ := jsonpath.Get("$[0].acceleration_type", v)
		if nodeNative14 != "" {
			customContainerConfig["accelerationType"] = nodeNative14
		}
		nodeNative15, _ := jsonpath.Get("$[0].web_server_mode", v)
		if nodeNative15 != "" {
			customContainerConfig["webServerMode"] = nodeNative15
		}
		objectDataLocalMap["customContainerConfig"] = customContainerConfig
	}
	if v, ok := d.GetOk("ca_port"); ok {
		objectDataLocalMap["caPort"] = v
	}
	if v, ok := d.GetOk("instance_concurrency"); ok {
		objectDataLocalMap["instanceConcurrency"] = v
	}
	if v, ok := d.GetOk("instance_type"); ok {
		objectDataLocalMap["instanceType"] = v
	}
	if v := d.Get("custom_health_check_config"); !IsNil(v) {
		customHealthCheckConfig := make(map[string]interface{})
		nodeNative19, _ := jsonpath.Get("$[0].http_get_url", v)
		if nodeNative19 != "" {
			customHealthCheckConfig["httpGetUrl"] = nodeNative19
		}
		nodeNative20, _ := jsonpath.Get("$[0].initial_delay_seconds", v)
		if nodeNative20 != "" {
			customHealthCheckConfig["initialDelaySeconds"] = nodeNative20
		}
		nodeNative21, _ := jsonpath.Get("$[0].period_seconds", v)
		if nodeNative21 != "" {
			customHealthCheckConfig["periodSeconds"] = nodeNative21
		}
		nodeNative22, _ := jsonpath.Get("$[0].timeout_seconds", v)
		if nodeNative22 != "" {
			customHealthCheckConfig["timeoutSeconds"] = nodeNative22
		}
		nodeNative23, _ := jsonpath.Get("$[0].failure_threshold", v)
		if nodeNative23 != "" {
			customHealthCheckConfig["failureThreshold"] = nodeNative23
		}
		nodeNative24, _ := jsonpath.Get("$[0].success_threshold", v)
		if nodeNative24 != "" {
			customHealthCheckConfig["successThreshold"] = nodeNative24
		}
		objectDataLocalMap["customHealthCheckConfig"] = customHealthCheckConfig
	}
	if v, ok := d.GetOk("cpu"); ok {
		objectDataLocalMap["diskSize"] = d.Get("disk_size")
		objectDataLocalMap["cpu"] = v
	}
	if v, ok := d.GetOk("disk_size"); ok {
		objectDataLocalMap["cpu"] = d.Get("cpu")
		objectDataLocalMap["diskSize"] = v
	}
	if v, ok := d.GetOk("gpu_memory_size"); ok {
		objectDataLocalMap["gpuMemorySize"] = v
	}
	if v, ok := d.GetOk("layers"); ok {
		nodeNative28, _ := jsonpath.Get("$", v)
		if nodeNative28 != "" {
			objectDataLocalMap["layers"] = nodeNative28
		}
	}
	if v := d.Get("custom_runtime_config"); !IsNil(v) {
		customRuntimeConfig := make(map[string]interface{})
		nodeNative29, _ := jsonpath.Get("$[0].command", v)
		if nodeNative29 != "" {
			customRuntimeConfig["command"] = nodeNative29
		}
		nodeNative30, _ := jsonpath.Get("$[0].args", v)
		if nodeNative30 != "" {
			customRuntimeConfig["args"] = nodeNative30
		}
		objectDataLocalMap["customRuntimeConfig"] = customRuntimeConfig
	}
	if v := d.Get("custom_dns"); !IsNil(v) {
		customDNS := make(map[string]interface{})
		if v, ok := d.GetOk("custom_dns"); ok {
			localData, err := jsonpath.Get("$[0].dns_options", v)
			if err != nil {
				return WrapError(err)
			}
			localMaps := make([]map[string]interface{}, 0)
			for _, dataLoop := range localData.([]interface{}) {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["name"] = dataLoopTmp["name"]
				dataLoopMap["value"] = dataLoopTmp["value"]
				localMaps = append(localMaps, dataLoopMap)
			}
			customDNS["dnsOptions"] = localMaps
		}
		nodeNative33, _ := jsonpath.Get("$[0].name_servers", v)
		if nodeNative33 != "" {
			customDNS["nameServers"] = nodeNative33
		}
		nodeNative34, _ := jsonpath.Get("$[0].searches", v)
		if nodeNative34 != "" {
			customDNS["searches"] = nodeNative34
		}
		objectDataLocalMap["customDNS"] = customDNS
	}
	if v, ok := d.GetOk("environment_variables"); ok {
		objectDataLocalMap["environmentVariables"] = v
	}
	if v := d.Get("instance_lifecycle_config"); !IsNil(v) {
		instanceLifecycleConfig := make(map[string]interface{})
		preFreeze := make(map[string]interface{})
		nodeNative36, _ := jsonpath.Get("$[0].pre_freeze[0].handler", v)
		if nodeNative36 != "" {
			preFreeze["handler"] = nodeNative36
		}
		nodeNative37, _ := jsonpath.Get("$[0].pre_freeze[0].timeout", v)
		if nodeNative37 != "" {
			preFreeze["timeout"] = nodeNative37
		}
		instanceLifecycleConfig["preFreeze"] = preFreeze
		preStop := make(map[string]interface{})
		nodeNative38, _ := jsonpath.Get("$[0].pre_stop[0].handler", v)
		if nodeNative38 != "" {
			preStop["handler"] = nodeNative38
		}
		nodeNative39, _ := jsonpath.Get("$[0].pre_stop[0].timeout", v)
		if nodeNative39 != "" {
			preStop["timeout"] = nodeNative39
		}
		instanceLifecycleConfig["preStop"] = preStop
		objectDataLocalMap["instanceLifecycleConfig"] = instanceLifecycleConfig
	}
	request["function"] = objectDataLocalMap

	body = objectDataLocalMap
	headerParams := make(map[string]*string)
	b, err := json.Marshal(body)
	headerParams["content-md5"] = tea.String(MD5(b))
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2021-04-06"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), query, headerParams, body, &util.RuntimeOptions{})

		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrentUpdateError"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_fcv2_function", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", serviceName, objectDataLocalMap["functionName"]))

	return resourceAliCloudFcv2FunctionRead(d, meta)
}

func resourceAliCloudFcv2FunctionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fcv2ServiceV2 := Fcv2ServiceV2{client}

	objectRaw, err := fcv2ServiceV2.DescribeFcv2Function(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_fcv2_function DescribeFcv2Function Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("ca_port", objectRaw["caPort"])
	d.Set("code_checksum", objectRaw["codeChecksum"])
	d.Set("cpu", objectRaw["cpu"])
	d.Set("create_time", objectRaw["createdTime"])
	d.Set("description", objectRaw["description"])
	d.Set("disk_size", objectRaw["diskSize"])
	d.Set("environment_variables", objectRaw["environmentVariables"])
	d.Set("gpu_memory_size", objectRaw["gpuMemorySize"])
	d.Set("handler", objectRaw["handler"])
	d.Set("initialization_timeout", objectRaw["initializationTimeout"])
	d.Set("initializer", objectRaw["initializer"])
	d.Set("instance_concurrency", objectRaw["instanceConcurrency"])
	d.Set("instance_type", objectRaw["instanceType"])
	d.Set("memory_size", objectRaw["memorySize"])
	d.Set("runtime", objectRaw["runtime"])
	d.Set("timeout", objectRaw["timeout"])
	d.Set("function_name", objectRaw["functionName"])
	customContainerConfigMaps := make([]map[string]interface{}, 0)
	customContainerConfigMap := make(map[string]interface{})
	customContainerConfig1Raw := make(map[string]interface{})
	if objectRaw["customContainerConfig"] != nil {
		customContainerConfig1Raw = objectRaw["customContainerConfig"].(map[string]interface{})
	}
	if len(customContainerConfig1Raw) > 0 {
		customContainerConfigMap["acceleration_type"] = customContainerConfig1Raw["accelerationType"]
		customContainerConfigMap["args"] = customContainerConfig1Raw["args"]
		customContainerConfigMap["command"] = customContainerConfig1Raw["command"]
		customContainerConfigMap["image"] = customContainerConfig1Raw["image"]
		customContainerConfigMap["web_server_mode"] = customContainerConfig1Raw["webServerMode"]
		customContainerConfigMaps = append(customContainerConfigMaps, customContainerConfigMap)
	}
	d.Set("custom_container_config", customContainerConfigMaps)
	customDnsMaps := make([]map[string]interface{}, 0)
	customDnsMap := make(map[string]interface{})
	customDNS1Raw := make(map[string]interface{})
	if objectRaw["customDNS"] != nil {
		customDNS1Raw = objectRaw["customDNS"].(map[string]interface{})
	}
	if len(customDNS1Raw) > 0 {
		dnsOptions1Raw := customDNS1Raw["dnsOptions"]
		dnsOptionsMaps := make([]map[string]interface{}, 0)
		if dnsOptions1Raw != nil {
			for _, dnsOptionsChild1Raw := range dnsOptions1Raw.([]interface{}) {
				dnsOptionsMap := make(map[string]interface{})
				dnsOptionsChild1Raw := dnsOptionsChild1Raw.(map[string]interface{})
				dnsOptionsMap["name"] = dnsOptionsChild1Raw["name"]
				dnsOptionsMap["value"] = dnsOptionsChild1Raw["value"]
				dnsOptionsMaps = append(dnsOptionsMaps, dnsOptionsMap)
			}
		}
		customDnsMap["dns_options"] = dnsOptionsMaps
		nameServers1Raw := make([]interface{}, 0)
		if customDNS1Raw["nameServers"] != nil {
			nameServers1Raw = customDNS1Raw["nameServers"].([]interface{})
		}

		customDnsMap["name_servers"] = nameServers1Raw
		searches1Raw := make([]interface{}, 0)
		if customDNS1Raw["searches"] != nil {
			searches1Raw = customDNS1Raw["searches"].([]interface{})
		}

		customDnsMap["searches"] = searches1Raw
		customDnsMaps = append(customDnsMaps, customDnsMap)
	}
	d.Set("custom_dns", customDnsMaps)
	customHealthCheckConfigMaps := make([]map[string]interface{}, 0)
	customHealthCheckConfigMap := make(map[string]interface{})
	customHealthCheckConfig1Raw := make(map[string]interface{})
	if objectRaw["customHealthCheckConfig"] != nil {
		customHealthCheckConfig1Raw = objectRaw["customHealthCheckConfig"].(map[string]interface{})
	}
	if len(customHealthCheckConfig1Raw) > 0 {
		customHealthCheckConfigMap["failure_threshold"] = customHealthCheckConfig1Raw["failureThreshold"]
		customHealthCheckConfigMap["http_get_url"] = customHealthCheckConfig1Raw["httpGetUrl"]
		customHealthCheckConfigMap["initial_delay_seconds"] = customHealthCheckConfig1Raw["initialDelaySeconds"]
		customHealthCheckConfigMap["period_seconds"] = customHealthCheckConfig1Raw["periodSeconds"]
		customHealthCheckConfigMap["success_threshold"] = customHealthCheckConfig1Raw["successThreshold"]
		customHealthCheckConfigMap["timeout_seconds"] = customHealthCheckConfig1Raw["timeoutSeconds"]
		customHealthCheckConfigMaps = append(customHealthCheckConfigMaps, customHealthCheckConfigMap)
	}
	d.Set("custom_health_check_config", customHealthCheckConfigMaps)
	customRuntimeConfigMaps := make([]map[string]interface{}, 0)
	customRuntimeConfigMap := make(map[string]interface{})
	customRuntimeConfig1Raw := make(map[string]interface{})
	if objectRaw["customRuntimeConfig"] != nil {
		customRuntimeConfig1Raw = objectRaw["customRuntimeConfig"].(map[string]interface{})
	}
	if len(customRuntimeConfig1Raw) > 0 {
		args5Raw := make([]interface{}, 0)
		if customRuntimeConfig1Raw["args"] != nil {
			args5Raw = customRuntimeConfig1Raw["args"].([]interface{})
		}

		customRuntimeConfigMap["args"] = args5Raw
		command5Raw := make([]interface{}, 0)
		if customRuntimeConfig1Raw["command"] != nil {
			command5Raw = customRuntimeConfig1Raw["command"].([]interface{})
		}

		customRuntimeConfigMap["command"] = command5Raw
		customRuntimeConfigMaps = append(customRuntimeConfigMaps, customRuntimeConfigMap)
	}
	d.Set("custom_runtime_config", customRuntimeConfigMaps)
	instanceLifecycleConfigMaps := make([]map[string]interface{}, 0)
	instanceLifecycleConfigMap := make(map[string]interface{})
	instanceLifecycleConfig1Raw := make(map[string]interface{})
	if objectRaw["instanceLifecycleConfig"] != nil {
		instanceLifecycleConfig1Raw = objectRaw["instanceLifecycleConfig"].(map[string]interface{})
	}
	if len(instanceLifecycleConfig1Raw) > 0 {
		preFreezeMaps := make([]map[string]interface{}, 0)
		preFreezeMap := make(map[string]interface{})
		preFreeze1Raw := make(map[string]interface{})
		if instanceLifecycleConfig1Raw["preFreeze"] != nil {
			preFreeze1Raw = instanceLifecycleConfig1Raw["preFreeze"].(map[string]interface{})
		}
		if len(preFreeze1Raw) > 0 {
			preFreezeMap["handler"] = preFreeze1Raw["handler"]
			preFreezeMap["timeout"] = preFreeze1Raw["timeout"]
			preFreezeMaps = append(preFreezeMaps, preFreezeMap)
		}
		instanceLifecycleConfigMap["pre_freeze"] = preFreezeMaps
		preStopMaps := make([]map[string]interface{}, 0)
		preStopMap := make(map[string]interface{})
		preStop1Raw := make(map[string]interface{})
		if instanceLifecycleConfig1Raw["preStop"] != nil {
			preStop1Raw = instanceLifecycleConfig1Raw["preStop"].(map[string]interface{})
		}
		if len(preStop1Raw) > 0 {
			preStopMap["handler"] = preStop1Raw["handler"]
			preStopMap["timeout"] = preStop1Raw["timeout"]
			preStopMaps = append(preStopMaps, preStopMap)
		}
		instanceLifecycleConfigMap["pre_stop"] = preStopMaps
		instanceLifecycleConfigMaps = append(instanceLifecycleConfigMaps, instanceLifecycleConfigMap)
	}
	d.Set("instance_lifecycle_config", instanceLifecycleConfigMaps)
	layersArnV21Raw := make([]interface{}, 0)
	if objectRaw["layersArnV2"] != nil {
		layersArnV21Raw = objectRaw["layersArnV2"].([]interface{})
	}

	d.Set("layers", layersArnV21Raw)

	codeMaps := make([]map[string]interface{}, 0)
	if v, ok := d.GetOk("code"); ok {
		codeMap := make(map[string]interface{})
		oldConfig := v.([]interface{})
		if len(oldConfig) > 0 {
			val := oldConfig[0].(map[string]interface{})
			codeMap["zip_file"] = val["zip_file"]
			codeMap["oss_object_name"] = val["oss_object_name"]
			codeMap["oss_bucket_name"] = val["oss_bucket_name"]
		}
		codeMaps = append(codeMaps, codeMap)
	}
	d.Set("code", codeMaps)

	parts := strings.Split(d.Id(), ":")
	d.Set("service_name", parts[0])
	if accountId, err := client.AccountId(); err != nil {
		log.Print(WrapError(err))
	} else {
		d.Set("function_arn", fmt.Sprintf("acs:fc:%s:%s:services/%s.LATEST/functions/%s", client.RegionId, accountId, parts[0], parts[1]))
	}

	return nil
}

func resourceAliCloudFcv2FunctionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	serviceName := parts[0]
	functionName := parts[1]
	action := fmt.Sprintf("/2021-04-06/services/%s/functions/%s", serviceName, functionName)
	conn, err := client.NewFcv2Client()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	objectDataLocalMap := make(map[string]interface{})
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			objectDataLocalMap["description"] = v
		}
	}
	if d.HasChange("handler") {
		update = true
		if v, ok := d.GetOk("handler"); ok {
			objectDataLocalMap["handler"] = v
		}
	}
	if d.HasChange("memory_size") {
		update = true
		if v, ok := d.GetOk("memory_size"); ok {
			objectDataLocalMap["memorySize"] = v
		}
	}
	if d.HasChange("runtime") {
		update = true
	}
	if v, ok := d.GetOk("runtime"); ok {
		objectDataLocalMap["runtime"] = v
	}
	if d.HasChange("timeout") {
		update = true
		if v, ok := d.GetOk("timeout"); ok {
			objectDataLocalMap["timeout"] = v
		}
	}
	if d.HasChange("initialization_timeout") {
		update = true
		if v, ok := d.GetOk("initialization_timeout"); ok {
			objectDataLocalMap["initializationTimeout"] = v
		}
	}
	if d.HasChange("initializer") {
		update = true
		if v, ok := d.GetOk("initializer"); ok {
			objectDataLocalMap["initializer"] = v
		}
	}
	if d.HasChange("code") {
		update = true
		if v := d.Get("code"); !IsNil(v) {
			code := make(map[string]interface{})
			nodeNative7, _ := jsonpath.Get("$[0].oss_bucket_name", v)
			if nodeNative7 != "" {
				code["ossBucketName"] = nodeNative7
			}
			nodeNative8, _ := jsonpath.Get("$[0].oss_object_name", v)
			if nodeNative8 != "" {
				code["ossObjectName"] = nodeNative8
			}
			nodeNative9, _ := jsonpath.Get("$[0].zip_file", v)
			if nodeNative9 != "" {
				code["zipFile"] = nodeNative9
			}
			objectDataLocalMap["code"] = code
		}
	}
	if d.HasChange("custom_container_config") {
		update = true
		if v := d.Get("custom_container_config"); !IsNil(v) {
			customContainerConfig := make(map[string]interface{})
			nodeNative10, _ := jsonpath.Get("$[0].args", v)
			if nodeNative10 != "" {
				customContainerConfig["args"] = nodeNative10
			}
			nodeNative11, _ := jsonpath.Get("$[0].command", v)
			if nodeNative11 != "" {
				customContainerConfig["command"] = nodeNative11
			}
			nodeNative12, _ := jsonpath.Get("$[0].image", v)
			if nodeNative12 != "" {
				customContainerConfig["image"] = nodeNative12
			}
			nodeNative13, _ := jsonpath.Get("$[0].acceleration_type", v)
			if nodeNative13 != "" {
				customContainerConfig["accelerationType"] = nodeNative13
			}
			nodeNative14, _ := jsonpath.Get("$[0].web_server_mode", v)
			if nodeNative14 != "" {
				customContainerConfig["webServerMode"] = nodeNative14
			}
			objectDataLocalMap["customContainerConfig"] = customContainerConfig
		}
	}
	if d.HasChange("ca_port") {
		update = true
		if v, ok := d.GetOk("ca_port"); ok {
			objectDataLocalMap["caPort"] = v
		}
	}
	if d.HasChange("instance_concurrency") {
		update = true
		if v, ok := d.GetOk("instance_concurrency"); ok {
			objectDataLocalMap["InstanceConcurrency"] = v
		}
	}
	if d.HasChange("instance_type") {
		update = true
		if v, ok := d.GetOk("instance_type"); ok {
			objectDataLocalMap["instanceType"] = v
		}
	}
	if d.HasChange("custom_health_check_config") {
		update = true
		if v := d.Get("custom_health_check_config"); !IsNil(v) {
			customHealthCheckConfig := make(map[string]interface{})
			nodeNative18, _ := jsonpath.Get("$[0].http_get_url", v)
			if nodeNative18 != "" {
				customHealthCheckConfig["httpGetUrl"] = nodeNative18
			}
			nodeNative19, _ := jsonpath.Get("$[0].initial_delay_seconds", v)
			if nodeNative19 != "" {
				customHealthCheckConfig["initialDelaySeconds"] = nodeNative19
			}
			nodeNative20, _ := jsonpath.Get("$[0].period_seconds", v)
			if nodeNative20 != "" {
				customHealthCheckConfig["periodSeconds"] = nodeNative20
			}
			nodeNative21, _ := jsonpath.Get("$[0].timeout_seconds", v)
			if nodeNative21 != "" {
				customHealthCheckConfig["timeoutSeconds"] = nodeNative21
			}
			nodeNative22, _ := jsonpath.Get("$[0].failure_threshold", v)
			if nodeNative22 != "" {
				customHealthCheckConfig["failureThreshold"] = nodeNative22
			}
			nodeNative23, _ := jsonpath.Get("$[0].success_threshold", v)
			if nodeNative23 != "" {
				customHealthCheckConfig["successThreshold"] = nodeNative23
			}
			objectDataLocalMap["customHealthCheckConfig"] = customHealthCheckConfig
		}
	}
	if d.HasChange("cpu") {
		update = true
		if v, ok := d.GetOk("cpu"); ok {
			objectDataLocalMap["diskSize"] = d.Get("disk_size")
			objectDataLocalMap["cpu"] = v
		}
	}
	if d.HasChange("disk_size") {
		update = true
		if v, ok := d.GetOk("disk_size"); ok {
			objectDataLocalMap["cpu"] = d.Get("cpu")
			objectDataLocalMap["diskSize"] = v
		}
	}
	if d.HasChange("gpu_memory_size") {
		update = true
	}
	if v, ok := d.GetOk("gpu_memory_size"); ok {
		objectDataLocalMap["gpuMemorySize"] = v
	}
	if d.HasChange("layers") {
		update = true
		objectDataLocalMap["layers"] = d.Get("layers")
	}
	if d.HasChange("custom_dns") {
		update = true
		if v := d.Get("custom_dns"); !IsNil(v) {
			customDNS := make(map[string]interface{})
			nodeNative28, _ := jsonpath.Get("$[0].name_servers", v)
			if nodeNative28 != "" {
				customDNS["nameServers"] = nodeNative28
			}
			if v, ok := d.GetOk("custom_dns"); ok {
				localData, err := jsonpath.Get("$[0].dns_options", v)
				if err != nil {
					return WrapError(err)
				}
				localMaps := make([]map[string]interface{}, 0)
				for _, dataLoop := range localData.([]interface{}) {
					dataLoopTmp := dataLoop.(map[string]interface{})
					dataLoopMap := make(map[string]interface{})
					dataLoopMap["name"] = dataLoopTmp["name"]
					dataLoopMap["value"] = dataLoopTmp["value"]
					localMaps = append(localMaps, dataLoopMap)
				}
				customDNS["dnsOptions"] = localMaps
			}
			nodeNative31, _ := jsonpath.Get("$[0].searches", v)
			if nodeNative31 != "" {
				customDNS["searches"] = nodeNative31
			}
			objectDataLocalMap["customDNS"] = customDNS
		}
	}
	if d.HasChange("custom_runtime_config") {
		update = true
		if v := d.Get("custom_runtime_config"); !IsNil(v) {
			customRuntimeConfig := make(map[string]interface{})
			nodeNative32, _ := jsonpath.Get("$[0].command", v)
			if nodeNative32 != "" {
				customRuntimeConfig["command"] = nodeNative32
			}
			nodeNative33, _ := jsonpath.Get("$[0].args", v)
			if nodeNative33 != "" {
				customRuntimeConfig["args"] = nodeNative33
			}
			objectDataLocalMap["customRuntimeConfig"] = customRuntimeConfig
		}
	}
	if d.HasChange("environment_variables") {
		update = true
		objectDataLocalMap["environmentVariables"] = d.Get("environment_variables")
	}
	if d.HasChange("instance_lifecycle_config") {
		update = true
		if v := d.Get("instance_lifecycle_config"); !IsNil(v) {
			instanceLifecycleConfig := make(map[string]interface{})
			preStop := make(map[string]interface{})
			nodeNative35, _ := jsonpath.Get("$[0].pre_stop[0].handler", v)
			if nodeNative35 != "" {
				preStop["handler"] = nodeNative35
			}
			nodeNative36, _ := jsonpath.Get("$[0].pre_stop[0].timeout", v)
			if nodeNative36 != "" {
				preStop["timeout"] = nodeNative36
			}
			instanceLifecycleConfig["preStop"] = preStop
			preFreeze := make(map[string]interface{})
			nodeNative37, _ := jsonpath.Get("$[0].pre_freeze[0].handler", v)
			if nodeNative37 != "" {
				preFreeze["handler"] = nodeNative37
			}
			nodeNative38, _ := jsonpath.Get("$[0].pre_freeze[0].timeout", v)
			if nodeNative38 != "" {
				preFreeze["timeout"] = nodeNative38
			}
			instanceLifecycleConfig["preFreeze"] = preFreeze
			objectDataLocalMap["instanceLifecycleConfig"] = instanceLifecycleConfig
		}
	}
	request["functionUpdateFields"] = objectDataLocalMap
	if d.HasChange("code_checksum") {
		update = true
		request["X-Fc-Code-Checksum"] = d.Get("code_checksum")
	}

	body = objectDataLocalMap
	headerParams := make(map[string]*string)
	b, err := json.Marshal(body)
	headerParams["content-md5"] = tea.String(MD5(b))
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2021-04-06"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), query, headerParams, body, &util.RuntimeOptions{})

			if err != nil {
				if IsExpectedErrors(err, []string{"ConcurrentUpdateError"}) || NeedRetry(err) {
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
	}

	return resourceAliCloudFcv2FunctionRead(d, meta)
}

func resourceAliCloudFcv2FunctionDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	serviceName := parts[0]
	functionName := parts[1]
	action := fmt.Sprintf("/2021-04-06/services/%s/functions/%s", serviceName, functionName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	conn, err := client.NewFcv2Client()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2021-04-06"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), query, nil, nil, &util.RuntimeOptions{})

		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrentUpdateError"}) || NeedRetry(err) {
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

	return nil
}
