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

func resourceAliCloudEfloExperimentPlanTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEfloExperimentPlanTemplateCreate,
		Read:   resourceAliCloudEfloExperimentPlanTemplateRead,
		Update: resourceAliCloudEfloExperimentPlanTemplateUpdate,
		Delete: resourceAliCloudEfloExperimentPlanTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"privacy_level": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"template_description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"template_pipeline": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"env_params": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"gpu_driver_version": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"py_torch_version": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"nccl_version": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"gpu_per_worker": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"cpu_per_worker": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"share_memory": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"worker_num": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"memory_per_worker": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"cuda_version": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"workload_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"workload_id": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"setting_params": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
						},
						"pipeline_order": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"scene": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudEfloExperimentPlanTemplateCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateExperimentPlanTemplate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("template_pipeline"); ok {
		templatePipelineMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			localData1 := make(map[string]interface{})
			memoryPerWorker1, _ := jsonpath.Get("$[0].memory_per_worker", dataLoopTmp["env_params"])
			if memoryPerWorker1 != nil && memoryPerWorker1 != "" {
				localData1["MemoryPerWorker"] = memoryPerWorker1
			}
			pyTorchVersion1, _ := jsonpath.Get("$[0].py_torch_version", dataLoopTmp["env_params"])
			if pyTorchVersion1 != nil && pyTorchVersion1 != "" {
				localData1["PyTorchVersion"] = pyTorchVersion1
			}
			cudaVersion1, _ := jsonpath.Get("$[0].cuda_version", dataLoopTmp["env_params"])
			if cudaVersion1 != nil && cudaVersion1 != "" {
				localData1["CudaVersion"] = cudaVersion1
			}
			shareMemory1, _ := jsonpath.Get("$[0].share_memory", dataLoopTmp["env_params"])
			if shareMemory1 != nil && shareMemory1 != "" {
				localData1["ShareMemory"] = shareMemory1
			}
			cpuPerWorker1, _ := jsonpath.Get("$[0].cpu_per_worker", dataLoopTmp["env_params"])
			if cpuPerWorker1 != nil && cpuPerWorker1 != "" {
				localData1["CpuPerWorker"] = cpuPerWorker1
			}
			gpuDriverVersion1, _ := jsonpath.Get("$[0].gpu_driver_version", dataLoopTmp["env_params"])
			if gpuDriverVersion1 != nil && gpuDriverVersion1 != "" {
				localData1["GpuDriverVersion"] = gpuDriverVersion1
			}
			ncclVersion, _ := jsonpath.Get("$[0].nccl_version", dataLoopTmp["env_params"])
			if ncclVersion != nil && ncclVersion != "" {
				localData1["NCCLVersion"] = ncclVersion
			}
			gpuPerWorker1, _ := jsonpath.Get("$[0].gpu_per_worker", dataLoopTmp["env_params"])
			if gpuPerWorker1 != nil && gpuPerWorker1 != "" {
				localData1["GpuPerWorker"] = gpuPerWorker1
			}
			workerNum1, _ := jsonpath.Get("$[0].worker_num", dataLoopTmp["env_params"])
			if workerNum1 != nil && workerNum1 != "" {
				localData1["WorkerNum"] = workerNum1
			}
			dataLoopMap["EnvParams"] = localData1
			dataLoopMap["WorkloadId"] = dataLoopTmp["workload_id"]
			dataLoopMap["WorkloadName"] = dataLoopTmp["workload_name"]
			dataLoopMap["Scene"] = dataLoopTmp["scene"]
			dataLoopMap["PipelineOrder"] = dataLoopTmp["pipeline_order"]
			dataLoopMap["SettingParams"] = dataLoopTmp["setting_params"]
			templatePipelineMapsArray = append(templatePipelineMapsArray, dataLoopMap)
		}
		templatePipelineMapsJson, err := json.Marshal(templatePipelineMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["TemplatePipeline"] = string(templatePipelineMapsJson)
	}

	request["PrivacyLevel"] = d.Get("privacy_level")
	request["TemplateName"] = d.Get("template_name")
	if v, ok := d.GetOk("template_description"); ok {
		request["TemplateDescription"] = v
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_eflo_experiment_plan_template", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Data.TemplateCode", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudEfloExperimentPlanTemplateRead(d, meta)
}

func resourceAliCloudEfloExperimentPlanTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	efloServiceV2 := EfloServiceV2{client}

	objectRaw, err := efloServiceV2.DescribeEfloExperimentPlanTemplate(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_eflo_experiment_plan_template DescribeEfloExperimentPlanTemplate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("privacy_level", objectRaw["PrivacyLevel"])
	d.Set("template_description", objectRaw["TemplateDescription"])
	d.Set("template_name", objectRaw["TemplateName"])

	templatePipelineParamRaw := objectRaw["TemplatePipelineParam"]
	templatePipelineMaps := make([]map[string]interface{}, 0)
	if templatePipelineParamRaw != nil {
		for index, templatePipelineParamChildRaw := range templatePipelineParamRaw.([]interface{}) {
			templatePipelineMap := make(map[string]interface{})
			templatePipelineParamChildRaw := templatePipelineParamChildRaw.(map[string]interface{})
			templatePipelineMap["pipeline_order"] = templatePipelineParamChildRaw["PipelineOrder"]
			templatePipelineMap["scene"] = templatePipelineParamChildRaw["Scene"]
			templatePipelineMap["workload_id"] = templatePipelineParamChildRaw["WorkloadId"]
			templatePipelineMap["workload_name"] = templatePipelineParamChildRaw["WorkloadName"]

			getSettingParams := make(map[string]interface{})
			settingParamsField := fmt.Sprintf("%v.%v.%v", "template_pipeline", index, "setting_params")
			if v, ok := d.GetOk(settingParamsField); ok && v != nil {
				getSettingParams = v.(map[string]interface{})
			}

			if templatePipelineParamChildRaw["SettingParams"] != "" {
				settingParamsMap := make(map[string]string)
				settingParams := templatePipelineParamChildRaw["SettingParams"].(map[string]interface{})

				for k, v := range settingParams {
					// There is an openapi bug that it will return all of SettingParams even through the config does not specified by user.
					// This workaround is not prefect when user set the SettingParams
					if _, ok := getSettingParams[k]; !ok && len(getSettingParams) > 0 {
						continue
					}
					settingParamsMap[k] = fmt.Sprint(v)
				}
				templatePipelineMap["setting_params"] = settingParamsMap
			}

			envParamsMaps := make([]map[string]interface{}, 0)
			envParamsMap := make(map[string]interface{})
			envParamsRaw := make(map[string]interface{})
			if templatePipelineParamChildRaw["EnvParams"] != nil {
				envParamsRaw = templatePipelineParamChildRaw["EnvParams"].(map[string]interface{})
			}
			if len(envParamsRaw) > 0 {
				envParamsMap["cpu_per_worker"] = envParamsRaw["CpuPerWorker"]
				envParamsMap["cuda_version"] = envParamsRaw["CudaVersion"]
				envParamsMap["gpu_driver_version"] = envParamsRaw["GpuDriverVersion"]
				envParamsMap["gpu_per_worker"] = envParamsRaw["GpuPerWorker"]
				envParamsMap["memory_per_worker"] = envParamsRaw["MemoryPerWorker"]
				envParamsMap["nccl_version"] = envParamsRaw["NCCLVersion"]
				envParamsMap["py_torch_version"] = envParamsRaw["PyTorchVersion"]
				envParamsMap["share_memory"] = envParamsRaw["ShareMemory"]
				envParamsMap["worker_num"] = envParamsRaw["WorkerNum"]

				envParamsMaps = append(envParamsMaps, envParamsMap)
			}
			templatePipelineMap["env_params"] = envParamsMaps
			templatePipelineMaps = append(templatePipelineMaps, templatePipelineMap)
		}
	}
	if err := d.Set("template_pipeline", templatePipelineMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudEfloExperimentPlanTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateExperimentPlanTemplate"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["TemplateId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("template_pipeline") {
		update = true
	}
	if v, ok := d.GetOk("template_pipeline"); ok {
		templatePipelineMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			if !IsNil(dataLoopTmp["env_params"]) {
				localData1 := make(map[string]interface{})
				memoryPerWorker1, _ := jsonpath.Get("$[0].memory_per_worker", dataLoopTmp["env_params"])
				if memoryPerWorker1 != nil && memoryPerWorker1 != "" {
					localData1["MemoryPerWorker"] = memoryPerWorker1
				}
				pyTorchVersion1, _ := jsonpath.Get("$[0].py_torch_version", dataLoopTmp["env_params"])
				if pyTorchVersion1 != nil && pyTorchVersion1 != "" {
					localData1["PyTorchVersion"] = pyTorchVersion1
				}
				cudaVersion1, _ := jsonpath.Get("$[0].cuda_version", dataLoopTmp["env_params"])
				if cudaVersion1 != nil && cudaVersion1 != "" {
					localData1["CudaVersion"] = cudaVersion1
				}
				shareMemory1, _ := jsonpath.Get("$[0].share_memory", dataLoopTmp["env_params"])
				if shareMemory1 != nil && shareMemory1 != "" {
					localData1["ShareMemory"] = shareMemory1
				}
				cpuPerWorker1, _ := jsonpath.Get("$[0].cpu_per_worker", dataLoopTmp["env_params"])
				if cpuPerWorker1 != nil && cpuPerWorker1 != "" {
					localData1["CpuPerWorker"] = cpuPerWorker1
				}
				gpuDriverVersion1, _ := jsonpath.Get("$[0].gpu_driver_version", dataLoopTmp["env_params"])
				if gpuDriverVersion1 != nil && gpuDriverVersion1 != "" {
					localData1["GpuDriverVersion"] = gpuDriverVersion1
				}
				ncclVersion, _ := jsonpath.Get("$[0].nccl_version", dataLoopTmp["env_params"])
				if ncclVersion != nil && ncclVersion != "" {
					localData1["NCCLVersion"] = ncclVersion
				}
				gpuPerWorker1, _ := jsonpath.Get("$[0].gpu_per_worker", dataLoopTmp["env_params"])
				if gpuPerWorker1 != nil && gpuPerWorker1 != "" {
					localData1["GpuPerWorker"] = gpuPerWorker1
				}
				workerNum1, _ := jsonpath.Get("$[0].worker_num", dataLoopTmp["env_params"])
				if workerNum1 != nil && workerNum1 != "" {
					localData1["WorkerNum"] = workerNum1
				}
				dataLoopMap["EnvParams"] = localData1
			}
			dataLoopMap["WorkloadId"] = dataLoopTmp["workload_id"]
			dataLoopMap["WorkloadName"] = dataLoopTmp["workload_name"]
			dataLoopMap["Scene"] = dataLoopTmp["scene"]
			dataLoopMap["PipelineOrder"] = dataLoopTmp["pipeline_order"]
			dataLoopMap["SettingParams"] = dataLoopTmp["setting_params"]
			templatePipelineMapsArray = append(templatePipelineMapsArray, dataLoopMap)
		}
		templatePipelineMapsJson, err := json.Marshal(templatePipelineMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["TemplatePipeline"] = string(templatePipelineMapsJson)
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

	return resourceAliCloudEfloExperimentPlanTemplateRead(d, meta)
}

func resourceAliCloudEfloExperimentPlanTemplateDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteExperimentPlanTemplate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["TemplateId"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
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
		if IsExpectedErrors(err, []string{"NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
