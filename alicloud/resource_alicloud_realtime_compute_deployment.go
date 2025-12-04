// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudRealtimeComputeDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRealtimeComputeDeploymentCreate,
		Read:   resourceAliCloudRealtimeComputeDeploymentRead,
		Update: resourceAliCloudRealtimeComputeDeploymentUpdate,
		Delete: resourceAliCloudRealtimeComputeDeploymentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"artifact": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: StringInSlice([]string{"SQLSCRIPT", "JAR", "PYTHON"}, false),
						},
						"jar_artifact": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"jar_uri": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"main_args": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"entry_class": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"additional_dependencies": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"sql_artifact": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sql_script": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"additional_dependencies": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"python_artifact": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"entry_module": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"main_args": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"additional_python_archives": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"additional_python_libraries": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"additional_dependencies": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"python_artifact_uri": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"batch_resource_setting": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"basic_resource_setting": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"taskmanager_resource_setting_spec": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"memory": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"cpu": {
													Type:     schema.TypeFloat,
													Optional: true,
												},
											},
										},
									},
									"parallelism": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"jobmanager_resource_setting_spec": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"memory": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"cpu": {
													Type:     schema.TypeFloat,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
						"max_slot": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"deployment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deployment_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"deployment_target": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: StringInSlice([]string{"SESSION", "PER_JOB"}, false),
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"execution_mode": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"STREAMING", "BATCH"}, false),
			},
			"flink_conf": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"local_variables": {
				Type:     schema.TypeSet,
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
			"logging": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"log4j2_configuration_template": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"logging_profile": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"log4j_loggers": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"logger_name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"logger_level": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"log_reserve_policy": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"open_history": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"expiration_days": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"streaming_resource_setting": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"basic_resource_setting": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"taskmanager_resource_setting_spec": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"memory": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"cpu": {
													Type:     schema.TypeFloat,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
									"parallelism": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"jobmanager_resource_setting_spec": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"memory": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"cpu": {
													Type:     schema.TypeFloat,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"resource_setting_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"BASIC", "EXPERT"}, false),
						},
						"expert_resource_setting": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_plan": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"jobmanager_resource_setting_spec": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"memory": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"cpu": {
													Type:     schema.TypeFloat,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudRealtimeComputeDeploymentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	namespace := d.Get("namespace")
	action := fmt.Sprintf("/api/v2/namespaces/%s/deployments", namespace)
	var request map[string]interface{}
	var response map[string]interface{}
	header := make(map[string]*string)
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	header["workspace"] = StringPointer(d.Get("resource_id").(string))

	streamingResourceSetting := make(map[string]interface{})

	if v := d.Get("streaming_resource_setting"); !IsNil(v) {
		expertResourceSetting := make(map[string]interface{})
		jobmanagerResourceSettingSpec := make(map[string]interface{})
		cpu1, _ := jsonpath.Get("$[0].expert_resource_setting[0].jobmanager_resource_setting_spec[0].cpu", d.Get("streaming_resource_setting"))
		if cpu1 != nil && cpu1 != "" {
			jobmanagerResourceSettingSpec["cpu"] = cpu1
		}
		memory1, _ := jsonpath.Get("$[0].expert_resource_setting[0].jobmanager_resource_setting_spec[0].memory", d.Get("streaming_resource_setting"))
		if memory1 != nil && memory1 != "" {
			jobmanagerResourceSettingSpec["memory"] = memory1
		}

		if len(jobmanagerResourceSettingSpec) > 0 {
			expertResourceSetting["jobmanagerResourceSettingSpec"] = jobmanagerResourceSettingSpec
		}
		resourcePlan1, _ := jsonpath.Get("$[0].expert_resource_setting[0].resource_plan", d.Get("streaming_resource_setting"))
		if resourcePlan1 != nil && resourcePlan1 != "" {
			expertResourceSetting["resourcePlan"] = resourcePlan1
		}

		if len(expertResourceSetting) > 0 {
			streamingResourceSetting["expertResourceSetting"] = expertResourceSetting
		}
		resourceSettingMode1, _ := jsonpath.Get("$[0].resource_setting_mode", v)
		if resourceSettingMode1 != nil && resourceSettingMode1 != "" {
			streamingResourceSetting["resourceSettingMode"] = resourceSettingMode1
		}
		basicResourceSetting := make(map[string]interface{})
		taskmanagerResourceSettingSpec := make(map[string]interface{})
		memory3, _ := jsonpath.Get("$[0].basic_resource_setting[0].taskmanager_resource_setting_spec[0].memory", d.Get("streaming_resource_setting"))
		if memory3 != nil && memory3 != "" {
			taskmanagerResourceSettingSpec["memory"] = memory3
		}
		cpu3, _ := jsonpath.Get("$[0].basic_resource_setting[0].taskmanager_resource_setting_spec[0].cpu", d.Get("streaming_resource_setting"))
		if cpu3 != nil && cpu3 != "" {
			taskmanagerResourceSettingSpec["cpu"] = cpu3
		}

		if len(taskmanagerResourceSettingSpec) > 0 {
			basicResourceSetting["taskmanagerResourceSettingSpec"] = taskmanagerResourceSettingSpec
		}
		jobmanagerResourceSettingSpec1 := make(map[string]interface{})
		memory5, _ := jsonpath.Get("$[0].basic_resource_setting[0].jobmanager_resource_setting_spec[0].memory", d.Get("streaming_resource_setting"))
		if memory5 != nil && memory5 != "" {
			jobmanagerResourceSettingSpec1["memory"] = memory5
		}
		cpu5, _ := jsonpath.Get("$[0].basic_resource_setting[0].jobmanager_resource_setting_spec[0].cpu", d.Get("streaming_resource_setting"))
		if cpu5 != nil && cpu5 != "" {
			jobmanagerResourceSettingSpec1["cpu"] = cpu5
		}

		if len(jobmanagerResourceSettingSpec1) > 0 {
			basicResourceSetting["jobmanagerResourceSettingSpec"] = jobmanagerResourceSettingSpec1
		}
		parallelism1, _ := jsonpath.Get("$[0].basic_resource_setting[0].parallelism", d.Get("streaming_resource_setting"))
		if parallelism1 != nil && parallelism1 != "" {
			basicResourceSetting["parallelism"] = parallelism1
		}

		if len(basicResourceSetting) > 0 {
			streamingResourceSetting["basicResourceSetting"] = basicResourceSetting
		}

		request["streamingResourceSetting"] = streamingResourceSetting
	}

	if v, ok := d.GetOk("flink_conf"); ok {
		request["flinkConf"] = v
	}
	logging := make(map[string]interface{})

	if v := d.Get("logging"); !IsNil(v) {
		if v, ok := d.GetOk("logging"); ok {
			localData, err := jsonpath.Get("$[0].log4j_loggers", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range convertToInterfaceArray(localData) {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["loggerLevel"] = dataLoopTmp["logger_level"]
				dataLoopMap["loggerName"] = dataLoopTmp["logger_name"]
				localMaps = append(localMaps, dataLoopMap)
			}
			logging["log4jLoggers"] = localMaps
		}

		logReservePolicy := make(map[string]interface{})
		openHistory1, _ := jsonpath.Get("$[0].log_reserve_policy[0].open_history", d.Get("logging"))
		if openHistory1 != nil && openHistory1 != "" {
			logReservePolicy["openHistory"] = openHistory1
		}
		expirationDays1, _ := jsonpath.Get("$[0].log_reserve_policy[0].expiration_days", d.Get("logging"))
		if expirationDays1 != nil && expirationDays1 != "" {
			logReservePolicy["expirationDays"] = expirationDays1
		}

		if len(logReservePolicy) > 0 {
			logging["logReservePolicy"] = logReservePolicy
		}
		loggingProfile1, _ := jsonpath.Get("$[0].logging_profile", v)
		if loggingProfile1 != nil && loggingProfile1 != "" {
			logging["loggingProfile"] = loggingProfile1
		}
		log4J2ConfigurationTemplate, _ := jsonpath.Get("$[0].log4j2_configuration_template", v)
		if log4J2ConfigurationTemplate != nil && log4J2ConfigurationTemplate != "" {
			logging["log4j2ConfigurationTemplate"] = log4J2ConfigurationTemplate
		}

		request["logging"] = logging
	}

	batchResourceSetting := make(map[string]interface{})

	if v := d.Get("batch_resource_setting"); !IsNil(v) {
		basicResourceSetting1 := make(map[string]interface{})
		taskmanagerResourceSettingSpec1 := make(map[string]interface{})
		cpu7, _ := jsonpath.Get("$[0].basic_resource_setting[0].taskmanager_resource_setting_spec[0].cpu", d.Get("batch_resource_setting"))
		if cpu7 != nil && cpu7 != "" {
			taskmanagerResourceSettingSpec1["cpu"] = cpu7
		}
		memory7, _ := jsonpath.Get("$[0].basic_resource_setting[0].taskmanager_resource_setting_spec[0].memory", d.Get("batch_resource_setting"))
		if memory7 != nil && memory7 != "" {
			taskmanagerResourceSettingSpec1["memory"] = memory7
		}

		if len(taskmanagerResourceSettingSpec1) > 0 {
			basicResourceSetting1["taskmanagerResourceSettingSpec"] = taskmanagerResourceSettingSpec1
		}
		parallelism3, _ := jsonpath.Get("$[0].basic_resource_setting[0].parallelism", d.Get("batch_resource_setting"))
		if parallelism3 != nil && parallelism3 != "" {
			basicResourceSetting1["parallelism"] = parallelism3
		}
		jobmanagerResourceSettingSpec2 := make(map[string]interface{})
		cpu9, _ := jsonpath.Get("$[0].basic_resource_setting[0].jobmanager_resource_setting_spec[0].cpu", d.Get("batch_resource_setting"))
		if cpu9 != nil && cpu9 != "" {
			jobmanagerResourceSettingSpec2["cpu"] = cpu9
		}
		memory9, _ := jsonpath.Get("$[0].basic_resource_setting[0].jobmanager_resource_setting_spec[0].memory", d.Get("batch_resource_setting"))
		if memory9 != nil && memory9 != "" {
			jobmanagerResourceSettingSpec2["memory"] = memory9
		}

		if len(jobmanagerResourceSettingSpec2) > 0 {
			basicResourceSetting1["jobmanagerResourceSettingSpec"] = jobmanagerResourceSettingSpec2
		}

		if len(basicResourceSetting1) > 0 {
			batchResourceSetting["basicResourceSetting"] = basicResourceSetting1
		}
		maxSlot1, _ := jsonpath.Get("$[0].max_slot", v)
		if maxSlot1 != nil && maxSlot1 != "" {
			batchResourceSetting["maxSlot"] = maxSlot1
		}

		request["batchResourceSetting"] = batchResourceSetting
	}

	if v, ok := d.GetOk("labels"); ok {
		request["labels"] = v
	}
	if v, ok := d.GetOk("local_variables"); ok {
		localVariablesMapsArray := make([]interface{}, 0)
		for _, dataLoop1 := range convertToInterfaceArray(v) {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["value"] = dataLoop1Tmp["value"]
			dataLoop1Map["name"] = dataLoop1Tmp["name"]
			localVariablesMapsArray = append(localVariablesMapsArray, dataLoop1Map)
		}
		request["localVariables"] = localVariablesMapsArray
	}

	if v, ok := d.GetOk("engine_version"); ok {
		request["engineVersion"] = v
	}
	deploymentTarget := make(map[string]interface{})

	if v := d.Get("deployment_target"); v != nil {
		name3, _ := jsonpath.Get("$[0].name", v)
		if name3 != nil && name3 != "" {
			deploymentTarget["name"] = name3
		}
		mode1, _ := jsonpath.Get("$[0].mode", v)
		if mode1 != nil && mode1 != "" {
			deploymentTarget["mode"] = mode1
		}

		request["deploymentTarget"] = deploymentTarget
	}

	request["executionMode"] = d.Get("execution_mode")
	artifact := make(map[string]interface{})

	if v := d.Get("artifact"); v != nil {
		pythonArtifact := make(map[string]interface{})
		additionalDependencies1, _ := jsonpath.Get("$[0].python_artifact[0].additional_dependencies", d.Get("artifact"))
		if additionalDependencies1 != nil && additionalDependencies1 != "" {
			pythonArtifact["additionalDependencies"] = additionalDependencies1
		}
		additionalPythonArchives1, _ := jsonpath.Get("$[0].python_artifact[0].additional_python_archives", d.Get("artifact"))
		if additionalPythonArchives1 != nil && additionalPythonArchives1 != "" {
			pythonArtifact["additionalPythonArchives"] = additionalPythonArchives1
		}
		pythonArtifactUri1, _ := jsonpath.Get("$[0].python_artifact[0].python_artifact_uri", d.Get("artifact"))
		if pythonArtifactUri1 != nil && pythonArtifactUri1 != "" {
			pythonArtifact["pythonArtifactUri"] = pythonArtifactUri1
		}
		mainArgs1, _ := jsonpath.Get("$[0].python_artifact[0].main_args", d.Get("artifact"))
		if mainArgs1 != nil && mainArgs1 != "" {
			pythonArtifact["mainArgs"] = mainArgs1
		}
		additionalPythonLibraries1, _ := jsonpath.Get("$[0].python_artifact[0].additional_python_libraries", d.Get("artifact"))
		if additionalPythonLibraries1 != nil && additionalPythonLibraries1 != "" {
			pythonArtifact["additionalPythonLibraries"] = additionalPythonLibraries1
		}
		entryModule1, _ := jsonpath.Get("$[0].python_artifact[0].entry_module", d.Get("artifact"))
		if entryModule1 != nil && entryModule1 != "" {
			pythonArtifact["entryModule"] = entryModule1
		}

		if len(pythonArtifact) > 0 {
			artifact["pythonArtifact"] = pythonArtifact
		}
		jarArtifact := make(map[string]interface{})
		jarUri1, _ := jsonpath.Get("$[0].jar_artifact[0].jar_uri", d.Get("artifact"))
		if jarUri1 != nil && jarUri1 != "" {
			jarArtifact["jarUri"] = jarUri1
		}
		additionalDependencies3, _ := jsonpath.Get("$[0].jar_artifact[0].additional_dependencies", d.Get("artifact"))
		if additionalDependencies3 != nil && additionalDependencies3 != "" {
			jarArtifact["additionalDependencies"] = additionalDependencies3
		}
		mainArgs3, _ := jsonpath.Get("$[0].jar_artifact[0].main_args", d.Get("artifact"))
		if mainArgs3 != nil && mainArgs3 != "" {
			jarArtifact["mainArgs"] = mainArgs3
		}
		entryClass1, _ := jsonpath.Get("$[0].jar_artifact[0].entry_class", d.Get("artifact"))
		if entryClass1 != nil && entryClass1 != "" {
			jarArtifact["entryClass"] = entryClass1
		}

		if len(jarArtifact) > 0 {
			artifact["jarArtifact"] = jarArtifact
		}
		sqlArtifact := make(map[string]interface{})
		sqlScript1, _ := jsonpath.Get("$[0].sql_artifact[0].sql_script", d.Get("artifact"))
		if sqlScript1 != nil && sqlScript1 != "" {
			sqlArtifact["sqlScript"] = sqlScript1
		}
		additionalDependencies5, _ := jsonpath.Get("$[0].sql_artifact[0].additional_dependencies", d.Get("artifact"))
		if additionalDependencies5 != nil && additionalDependencies5 != "" {
			sqlArtifact["additionalDependencies"] = additionalDependencies5
		}

		if len(sqlArtifact) > 0 {
			artifact["sqlArtifact"] = sqlArtifact
		}
		kind1, _ := jsonpath.Get("$[0].kind", v)
		if kind1 != nil && kind1 != "" {
			artifact["kind"] = kind1
		}

		request["artifact"] = artifact
	}

	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	request["name"] = d.Get("deployment_name")
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("ververica", "2022-07-18", action, query, header, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_realtime_compute_deployment", action, AlibabaCloudSdkGoERROR)
	}

	dataworkspaceVar, _ := jsonpath.Get("$.data.workspace", response)
	datanamespaceVar, _ := jsonpath.Get("$.data.namespace", response)
	datadeploymentIdVar, _ := jsonpath.Get("$.data.deploymentId", response)
	d.SetId(fmt.Sprintf("%v:%v:%v", dataworkspaceVar, datanamespaceVar, datadeploymentIdVar))

	return resourceAliCloudRealtimeComputeDeploymentRead(d, meta)
}

func resourceAliCloudRealtimeComputeDeploymentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	realtimeComputeServiceV2 := RealtimeComputeServiceV2{client}

	objectRaw, err := realtimeComputeServiceV2.DescribeRealtimeComputeDeployment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_realtime_compute_deployment DescribeRealtimeComputeDeployment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("deployment_name", objectRaw["name"])
	d.Set("description", objectRaw["description"])
	d.Set("engine_version", objectRaw["engineVersion"])
	d.Set("execution_mode", objectRaw["executionMode"])
	d.Set("flink_conf", objectRaw["flinkConf"])
	d.Set("labels", objectRaw["labels"])
	d.Set("deployment_id", objectRaw["deploymentId"])
	d.Set("namespace", objectRaw["namespace"])
	d.Set("resource_id", objectRaw["workspace"])

	artifactMaps := make([]map[string]interface{}, 0)
	artifactMap := make(map[string]interface{})
	artifactRaw := make(map[string]interface{})
	if objectRaw["artifact"] != nil {
		artifactRaw = objectRaw["artifact"].(map[string]interface{})
	}
	if len(artifactRaw) > 0 {
		artifactMap["kind"] = artifactRaw["kind"]

		jarArtifactMaps := make([]map[string]interface{}, 0)
		jarArtifactMap := make(map[string]interface{})
		jarArtifactRaw := make(map[string]interface{})
		if artifactRaw["jarArtifact"] != nil {
			jarArtifactRaw = artifactRaw["jarArtifact"].(map[string]interface{})
		}
		if len(jarArtifactRaw) > 0 {
			jarArtifactMap["entry_class"] = jarArtifactRaw["entryClass"]
			jarArtifactMap["jar_uri"] = jarArtifactRaw["jarUri"]
			jarArtifactMap["main_args"] = jarArtifactRaw["mainArgs"]

			additionalDependenciesRaw := make([]interface{}, 0)
			if jarArtifactRaw["additionalDependencies"] != nil {
				additionalDependenciesRaw = convertToInterfaceArray(jarArtifactRaw["additionalDependencies"])
			}

			jarArtifactMap["additional_dependencies"] = additionalDependenciesRaw
			jarArtifactMaps = append(jarArtifactMaps, jarArtifactMap)
		}
		artifactMap["jar_artifact"] = jarArtifactMaps
		pythonArtifactMaps := make([]map[string]interface{}, 0)
		pythonArtifactMap := make(map[string]interface{})
		pythonArtifactRaw := make(map[string]interface{})
		if artifactRaw["pythonArtifact"] != nil {
			pythonArtifactRaw = artifactRaw["pythonArtifact"].(map[string]interface{})
		}
		if len(pythonArtifactRaw) > 0 {
			pythonArtifactMap["entry_module"] = pythonArtifactRaw["entryModule"]
			pythonArtifactMap["main_args"] = pythonArtifactRaw["mainArgs"]
			pythonArtifactMap["python_artifact_uri"] = pythonArtifactRaw["pythonArtifactUri"]

			additionalDependenciesRaw := make([]interface{}, 0)
			if pythonArtifactRaw["additionalDependencies"] != nil {
				additionalDependenciesRaw = convertToInterfaceArray(pythonArtifactRaw["additionalDependencies"])
			}

			pythonArtifactMap["additional_dependencies"] = additionalDependenciesRaw
			additionalPythonArchivesRaw := make([]interface{}, 0)
			if pythonArtifactRaw["additionalPythonArchives"] != nil {
				additionalPythonArchivesRaw = convertToInterfaceArray(pythonArtifactRaw["additionalPythonArchives"])
			}

			pythonArtifactMap["additional_python_archives"] = additionalPythonArchivesRaw
			additionalPythonLibrariesRaw := make([]interface{}, 0)
			if pythonArtifactRaw["additionalPythonLibraries"] != nil {
				additionalPythonLibrariesRaw = convertToInterfaceArray(pythonArtifactRaw["additionalPythonLibraries"])
			}

			pythonArtifactMap["additional_python_libraries"] = additionalPythonLibrariesRaw
			pythonArtifactMaps = append(pythonArtifactMaps, pythonArtifactMap)
		}
		artifactMap["python_artifact"] = pythonArtifactMaps
		sqlArtifactMaps := make([]map[string]interface{}, 0)
		sqlArtifactMap := make(map[string]interface{})
		sqlArtifactRaw := make(map[string]interface{})
		if artifactRaw["sqlArtifact"] != nil {
			sqlArtifactRaw = artifactRaw["sqlArtifact"].(map[string]interface{})
		}
		if len(sqlArtifactRaw) > 0 {
			sqlArtifactMap["sql_script"] = sqlArtifactRaw["sqlScript"]

			additionalDependenciesRaw := make([]interface{}, 0)
			if sqlArtifactRaw["additionalDependencies"] != nil {
				additionalDependenciesRaw = convertToInterfaceArray(sqlArtifactRaw["additionalDependencies"])
			}

			sqlArtifactMap["additional_dependencies"] = additionalDependenciesRaw
			sqlArtifactMaps = append(sqlArtifactMaps, sqlArtifactMap)
		}
		artifactMap["sql_artifact"] = sqlArtifactMaps
		artifactMaps = append(artifactMaps, artifactMap)
	}
	if err := d.Set("artifact", artifactMaps); err != nil {
		return err
	}
	batchResourceSettingMaps := make([]map[string]interface{}, 0)
	batchResourceSettingMap := make(map[string]interface{})
	batchResourceSettingRaw := make(map[string]interface{})
	if objectRaw["batchResourceSetting"] != nil {
		batchResourceSettingRaw = objectRaw["batchResourceSetting"].(map[string]interface{})
	}
	if len(batchResourceSettingRaw) > 0 {
		batchResourceSettingMap["max_slot"] = batchResourceSettingRaw["maxSlot"]

		basicResourceSettingMaps := make([]map[string]interface{}, 0)
		basicResourceSettingMap := make(map[string]interface{})
		basicResourceSettingRaw := make(map[string]interface{})
		if batchResourceSettingRaw["basicResourceSetting"] != nil {
			basicResourceSettingRaw = batchResourceSettingRaw["basicResourceSetting"].(map[string]interface{})
		}
		if len(basicResourceSettingRaw) > 0 {
			basicResourceSettingMap["parallelism"] = basicResourceSettingRaw["parallelism"]

			jobmanagerResourceSettingSpecMaps := make([]map[string]interface{}, 0)
			jobmanagerResourceSettingSpecMap := make(map[string]interface{})
			jobmanagerResourceSettingSpecRaw := make(map[string]interface{})
			if basicResourceSettingRaw["jobmanagerResourceSettingSpec"] != nil {
				jobmanagerResourceSettingSpecRaw = basicResourceSettingRaw["jobmanagerResourceSettingSpec"].(map[string]interface{})
			}
			if len(jobmanagerResourceSettingSpecRaw) > 0 {
				jobmanagerResourceSettingSpecMap["cpu"] = jobmanagerResourceSettingSpecRaw["cpu"]
				jobmanagerResourceSettingSpecMap["memory"] = jobmanagerResourceSettingSpecRaw["memory"]

				jobmanagerResourceSettingSpecMaps = append(jobmanagerResourceSettingSpecMaps, jobmanagerResourceSettingSpecMap)
			}
			basicResourceSettingMap["jobmanager_resource_setting_spec"] = jobmanagerResourceSettingSpecMaps
			taskmanagerResourceSettingSpecMaps := make([]map[string]interface{}, 0)
			taskmanagerResourceSettingSpecMap := make(map[string]interface{})
			taskmanagerResourceSettingSpecRaw := make(map[string]interface{})
			if basicResourceSettingRaw["taskmanagerResourceSettingSpec"] != nil {
				taskmanagerResourceSettingSpecRaw = basicResourceSettingRaw["taskmanagerResourceSettingSpec"].(map[string]interface{})
			}
			if len(taskmanagerResourceSettingSpecRaw) > 0 {
				taskmanagerResourceSettingSpecMap["cpu"] = taskmanagerResourceSettingSpecRaw["cpu"]
				taskmanagerResourceSettingSpecMap["memory"] = taskmanagerResourceSettingSpecRaw["memory"]

				taskmanagerResourceSettingSpecMaps = append(taskmanagerResourceSettingSpecMaps, taskmanagerResourceSettingSpecMap)
			}
			basicResourceSettingMap["taskmanager_resource_setting_spec"] = taskmanagerResourceSettingSpecMaps
			basicResourceSettingMaps = append(basicResourceSettingMaps, basicResourceSettingMap)
		}
		batchResourceSettingMap["basic_resource_setting"] = basicResourceSettingMaps
		batchResourceSettingMaps = append(batchResourceSettingMaps, batchResourceSettingMap)
	}
	if err := d.Set("batch_resource_setting", batchResourceSettingMaps); err != nil {
		return err
	}
	deploymentTargetMaps := make([]map[string]interface{}, 0)
	deploymentTargetMap := make(map[string]interface{})
	deploymentTargetRaw := make(map[string]interface{})
	if objectRaw["deploymentTarget"] != nil {
		deploymentTargetRaw = objectRaw["deploymentTarget"].(map[string]interface{})
	}
	if len(deploymentTargetRaw) > 0 {
		deploymentTargetMap["mode"] = deploymentTargetRaw["mode"]
		deploymentTargetMap["name"] = deploymentTargetRaw["name"]

		deploymentTargetMaps = append(deploymentTargetMaps, deploymentTargetMap)
	}
	if err := d.Set("deployment_target", deploymentTargetMaps); err != nil {
		return err
	}
	localVariablesRaw := objectRaw["localVariables"]
	localVariablesMaps := make([]map[string]interface{}, 0)
	if localVariablesRaw != nil {
		for _, localVariablesChildRaw := range convertToInterfaceArray(localVariablesRaw) {
			localVariablesMap := make(map[string]interface{})
			localVariablesChildRaw := localVariablesChildRaw.(map[string]interface{})
			localVariablesMap["name"] = localVariablesChildRaw["name"]
			localVariablesMap["value"] = localVariablesChildRaw["value"]

			localVariablesMaps = append(localVariablesMaps, localVariablesMap)
		}
	}
	if err := d.Set("local_variables", localVariablesMaps); err != nil {
		return err
	}
	loggingMaps := make([]map[string]interface{}, 0)
	loggingMap := make(map[string]interface{})
	loggingRaw := make(map[string]interface{})
	if objectRaw["logging"] != nil {
		loggingRaw = objectRaw["logging"].(map[string]interface{})
	}
	if len(loggingRaw) > 0 {
		loggingMap["log4j2_configuration_template"] = loggingRaw["log4j2ConfigurationTemplate"]
		loggingMap["logging_profile"] = loggingRaw["loggingProfile"]

		log4jLoggersRaw := loggingRaw["log4jLoggers"]
		log4JLoggersMaps := make([]map[string]interface{}, 0)
		if log4jLoggersRaw != nil {
			for _, log4jLoggersChildRaw := range convertToInterfaceArray(log4jLoggersRaw) {
				log4JLoggersMap := make(map[string]interface{})
				log4jLoggersChildRaw := log4jLoggersChildRaw.(map[string]interface{})
				log4JLoggersMap["logger_level"] = log4jLoggersChildRaw["loggerLevel"]
				log4JLoggersMap["logger_name"] = log4jLoggersChildRaw["loggerName"]

				log4JLoggersMaps = append(log4JLoggersMaps, log4JLoggersMap)
			}
		}
		loggingMap["log4j_loggers"] = log4JLoggersMaps
		logReservePolicyMaps := make([]map[string]interface{}, 0)
		logReservePolicyMap := make(map[string]interface{})
		logReservePolicyRaw := make(map[string]interface{})
		if loggingRaw["logReservePolicy"] != nil {
			logReservePolicyRaw = loggingRaw["logReservePolicy"].(map[string]interface{})
		}
		if len(logReservePolicyRaw) > 0 {
			logReservePolicyMap["expiration_days"] = logReservePolicyRaw["expirationDays"]
			logReservePolicyMap["open_history"] = logReservePolicyRaw["openHistory"]

			logReservePolicyMaps = append(logReservePolicyMaps, logReservePolicyMap)
		}
		loggingMap["log_reserve_policy"] = logReservePolicyMaps
		loggingMaps = append(loggingMaps, loggingMap)
	}
	if err := d.Set("logging", loggingMaps); err != nil {
		return err
	}
	streamingResourceSettingMaps := make([]map[string]interface{}, 0)
	streamingResourceSettingMap := make(map[string]interface{})
	streamingResourceSettingRaw := make(map[string]interface{})
	if objectRaw["streamingResourceSetting"] != nil {
		streamingResourceSettingRaw = objectRaw["streamingResourceSetting"].(map[string]interface{})
	}
	if len(streamingResourceSettingRaw) > 0 {
		streamingResourceSettingMap["resource_setting_mode"] = streamingResourceSettingRaw["resourceSettingMode"]

		basicResourceSettingMaps := make([]map[string]interface{}, 0)
		basicResourceSettingMap := make(map[string]interface{})
		basicResourceSettingRaw := make(map[string]interface{})
		if streamingResourceSettingRaw["basicResourceSetting"] != nil {
			basicResourceSettingRaw = streamingResourceSettingRaw["basicResourceSetting"].(map[string]interface{})
		}
		if len(basicResourceSettingRaw) > 0 {
			basicResourceSettingMap["parallelism"] = basicResourceSettingRaw["parallelism"]

			jobmanagerResourceSettingSpecMaps := make([]map[string]interface{}, 0)
			jobmanagerResourceSettingSpecMap := make(map[string]interface{})
			jobmanagerResourceSettingSpecRaw := make(map[string]interface{})
			if basicResourceSettingRaw["jobmanagerResourceSettingSpec"] != nil {
				jobmanagerResourceSettingSpecRaw = basicResourceSettingRaw["jobmanagerResourceSettingSpec"].(map[string]interface{})
			}
			if len(jobmanagerResourceSettingSpecRaw) > 0 {
				jobmanagerResourceSettingSpecMap["cpu"] = jobmanagerResourceSettingSpecRaw["cpu"]
				jobmanagerResourceSettingSpecMap["memory"] = jobmanagerResourceSettingSpecRaw["memory"]

				jobmanagerResourceSettingSpecMaps = append(jobmanagerResourceSettingSpecMaps, jobmanagerResourceSettingSpecMap)
			}
			basicResourceSettingMap["jobmanager_resource_setting_spec"] = jobmanagerResourceSettingSpecMaps
			taskmanagerResourceSettingSpecMaps := make([]map[string]interface{}, 0)
			taskmanagerResourceSettingSpecMap := make(map[string]interface{})
			taskmanagerResourceSettingSpecRaw := make(map[string]interface{})
			if basicResourceSettingRaw["taskmanagerResourceSettingSpec"] != nil {
				taskmanagerResourceSettingSpecRaw = basicResourceSettingRaw["taskmanagerResourceSettingSpec"].(map[string]interface{})
			}
			if len(taskmanagerResourceSettingSpecRaw) > 0 {
				taskmanagerResourceSettingSpecMap["cpu"] = taskmanagerResourceSettingSpecRaw["cpu"]
				taskmanagerResourceSettingSpecMap["memory"] = taskmanagerResourceSettingSpecRaw["memory"]

				taskmanagerResourceSettingSpecMaps = append(taskmanagerResourceSettingSpecMaps, taskmanagerResourceSettingSpecMap)
			}
			basicResourceSettingMap["taskmanager_resource_setting_spec"] = taskmanagerResourceSettingSpecMaps
			basicResourceSettingMaps = append(basicResourceSettingMaps, basicResourceSettingMap)
		}
		streamingResourceSettingMap["basic_resource_setting"] = basicResourceSettingMaps
		expertResourceSettingMaps := make([]map[string]interface{}, 0)
		expertResourceSettingMap := make(map[string]interface{})
		expertResourceSettingRaw := make(map[string]interface{})
		if streamingResourceSettingRaw["expertResourceSetting"] != nil {
			expertResourceSettingRaw = streamingResourceSettingRaw["expertResourceSetting"].(map[string]interface{})
		}
		if len(expertResourceSettingRaw) > 0 {
			expertResourceSettingMap["resource_plan"] = expertResourceSettingRaw["resourcePlan"]

			jobmanagerResourceSettingSpecMaps := make([]map[string]interface{}, 0)
			jobmanagerResourceSettingSpecMap := make(map[string]interface{})
			jobmanagerResourceSettingSpecRaw := make(map[string]interface{})
			if expertResourceSettingRaw["jobmanagerResourceSettingSpec"] != nil {
				jobmanagerResourceSettingSpecRaw = expertResourceSettingRaw["jobmanagerResourceSettingSpec"].(map[string]interface{})
			}
			if len(jobmanagerResourceSettingSpecRaw) > 0 {
				jobmanagerResourceSettingSpecMap["cpu"] = jobmanagerResourceSettingSpecRaw["cpu"]
				jobmanagerResourceSettingSpecMap["memory"] = jobmanagerResourceSettingSpecRaw["memory"]

				jobmanagerResourceSettingSpecMaps = append(jobmanagerResourceSettingSpecMaps, jobmanagerResourceSettingSpecMap)
			}
			expertResourceSettingMap["jobmanager_resource_setting_spec"] = jobmanagerResourceSettingSpecMaps
			expertResourceSettingMaps = append(expertResourceSettingMaps, expertResourceSettingMap)
		}
		streamingResourceSettingMap["expert_resource_setting"] = expertResourceSettingMaps
		streamingResourceSettingMaps = append(streamingResourceSettingMaps, streamingResourceSettingMap)
	}
	if err := d.Set("streaming_resource_setting", streamingResourceSettingMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudRealtimeComputeDeploymentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var header map[string]*string
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	namespace := parts[1]
	deploymentId := parts[2]
	action := fmt.Sprintf("/api/v2/namespaces/%s/deployments/%s", namespace, deploymentId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	header = make(map[string]*string)
	body = make(map[string]interface{})
	header["workspace"] = StringPointer(parts[0])

	if d.HasChange("streaming_resource_setting") {
		update = true
	}
	streamingResourceSetting := make(map[string]interface{})

	if v := d.Get("streaming_resource_setting"); !IsNil(v) || d.HasChange("streaming_resource_setting") {
		expertResourceSetting := make(map[string]interface{})
		resourcePlan1, _ := jsonpath.Get("$[0].expert_resource_setting[0].resource_plan", d.Get("streaming_resource_setting"))
		if resourcePlan1 != nil && resourcePlan1 != "" {
			expertResourceSetting["resourcePlan"] = resourcePlan1
		}
		jobmanagerResourceSettingSpec := make(map[string]interface{})
		memory1, _ := jsonpath.Get("$[0].expert_resource_setting[0].jobmanager_resource_setting_spec[0].memory", d.Get("streaming_resource_setting"))
		if memory1 != nil && memory1 != "" {
			jobmanagerResourceSettingSpec["memory"] = memory1
		}
		cpu1, _ := jsonpath.Get("$[0].expert_resource_setting[0].jobmanager_resource_setting_spec[0].cpu", d.Get("streaming_resource_setting"))
		if cpu1 != nil && cpu1 != "" {
			jobmanagerResourceSettingSpec["cpu"] = cpu1
		}

		if len(jobmanagerResourceSettingSpec) > 0 {
			expertResourceSetting["jobmanagerResourceSettingSpec"] = jobmanagerResourceSettingSpec
		}

		if len(expertResourceSetting) > 0 {
			streamingResourceSetting["expertResourceSetting"] = expertResourceSetting
		}
		resourceSettingMode1, _ := jsonpath.Get("$[0].resource_setting_mode", v)
		if resourceSettingMode1 != nil && resourceSettingMode1 != "" {
			streamingResourceSetting["resourceSettingMode"] = resourceSettingMode1
		}
		basicResourceSetting := make(map[string]interface{})
		taskmanagerResourceSettingSpec := make(map[string]interface{})
		memory3, _ := jsonpath.Get("$[0].basic_resource_setting[0].taskmanager_resource_setting_spec[0].memory", d.Get("streaming_resource_setting"))
		if memory3 != nil && memory3 != "" {
			taskmanagerResourceSettingSpec["memory"] = memory3
		}
		cpu3, _ := jsonpath.Get("$[0].basic_resource_setting[0].taskmanager_resource_setting_spec[0].cpu", d.Get("streaming_resource_setting"))
		if cpu3 != nil && cpu3 != "" {
			taskmanagerResourceSettingSpec["cpu"] = cpu3
		}

		if len(taskmanagerResourceSettingSpec) > 0 {
			basicResourceSetting["taskmanagerResourceSettingSpec"] = taskmanagerResourceSettingSpec
		}
		jobmanagerResourceSettingSpec1 := make(map[string]interface{})
		cpu5, _ := jsonpath.Get("$[0].basic_resource_setting[0].jobmanager_resource_setting_spec[0].cpu", d.Get("streaming_resource_setting"))
		if cpu5 != nil && cpu5 != "" {
			jobmanagerResourceSettingSpec1["cpu"] = cpu5
		}
		memory5, _ := jsonpath.Get("$[0].basic_resource_setting[0].jobmanager_resource_setting_spec[0].memory", d.Get("streaming_resource_setting"))
		if memory5 != nil && memory5 != "" {
			jobmanagerResourceSettingSpec1["memory"] = memory5
		}

		if len(jobmanagerResourceSettingSpec1) > 0 {
			basicResourceSetting["jobmanagerResourceSettingSpec"] = jobmanagerResourceSettingSpec1
		}
		parallelism1, _ := jsonpath.Get("$[0].basic_resource_setting[0].parallelism", d.Get("streaming_resource_setting"))
		if parallelism1 != nil && parallelism1 != "" {
			basicResourceSetting["parallelism"] = parallelism1
		}

		if len(basicResourceSetting) > 0 {
			streamingResourceSetting["basicResourceSetting"] = basicResourceSetting
		}

		request["streamingResourceSetting"] = streamingResourceSetting
	}

	if d.HasChange("artifact") {
		update = true
	}
	artifact := make(map[string]interface{})

	if v := d.Get("artifact"); v != nil {
		pythonArtifact := make(map[string]interface{})
		additionalDependencies1, _ := jsonpath.Get("$[0].python_artifact[0].additional_dependencies", d.Get("artifact"))
		if additionalDependencies1 != nil && additionalDependencies1 != "" {
			pythonArtifact["additionalDependencies"] = additionalDependencies1
		}
		additionalPythonArchives1, _ := jsonpath.Get("$[0].python_artifact[0].additional_python_archives", d.Get("artifact"))
		if additionalPythonArchives1 != nil && additionalPythonArchives1 != "" {
			pythonArtifact["additionalPythonArchives"] = additionalPythonArchives1
		}
		pythonArtifactUri1, _ := jsonpath.Get("$[0].python_artifact[0].python_artifact_uri", d.Get("artifact"))
		if pythonArtifactUri1 != nil && pythonArtifactUri1 != "" {
			pythonArtifact["pythonArtifactUri"] = pythonArtifactUri1
		}
		mainArgs1, _ := jsonpath.Get("$[0].python_artifact[0].main_args", d.Get("artifact"))
		if mainArgs1 != nil && mainArgs1 != "" {
			pythonArtifact["mainArgs"] = mainArgs1
		}
		additionalPythonLibraries1, _ := jsonpath.Get("$[0].python_artifact[0].additional_python_libraries", d.Get("artifact"))
		if additionalPythonLibraries1 != nil && additionalPythonLibraries1 != "" {
			pythonArtifact["additionalPythonLibraries"] = additionalPythonLibraries1
		}
		entryModule1, _ := jsonpath.Get("$[0].python_artifact[0].entry_module", d.Get("artifact"))
		if entryModule1 != nil && entryModule1 != "" {
			pythonArtifact["entryModule"] = entryModule1
		}

		if len(pythonArtifact) > 0 {
			artifact["pythonArtifact"] = pythonArtifact
		}
		jarArtifact := make(map[string]interface{})
		jarUri1, _ := jsonpath.Get("$[0].jar_artifact[0].jar_uri", d.Get("artifact"))
		if jarUri1 != nil && jarUri1 != "" {
			jarArtifact["jarUri"] = jarUri1
		}
		additionalDependencies3, _ := jsonpath.Get("$[0].jar_artifact[0].additional_dependencies", d.Get("artifact"))
		if additionalDependencies3 != nil && additionalDependencies3 != "" {
			jarArtifact["additionalDependencies"] = additionalDependencies3
		}
		mainArgs3, _ := jsonpath.Get("$[0].jar_artifact[0].main_args", d.Get("artifact"))
		if mainArgs3 != nil && mainArgs3 != "" {
			jarArtifact["mainArgs"] = mainArgs3
		}
		entryClass1, _ := jsonpath.Get("$[0].jar_artifact[0].entry_class", d.Get("artifact"))
		if entryClass1 != nil && entryClass1 != "" {
			jarArtifact["entryClass"] = entryClass1
		}

		if len(jarArtifact) > 0 {
			artifact["jarArtifact"] = jarArtifact
		}
		sqlArtifact := make(map[string]interface{})
		sqlScript1, _ := jsonpath.Get("$[0].sql_artifact[0].sql_script", d.Get("artifact"))
		if sqlScript1 != nil && sqlScript1 != "" {
			sqlArtifact["sqlScript"] = sqlScript1
		}
		additionalDependencies5, _ := jsonpath.Get("$[0].sql_artifact[0].additional_dependencies", d.Get("artifact"))
		if additionalDependencies5 != nil && additionalDependencies5 != "" {
			sqlArtifact["additionalDependencies"] = additionalDependencies5
		}

		if len(sqlArtifact) > 0 {
			artifact["sqlArtifact"] = sqlArtifact
		}
		kind1, _ := jsonpath.Get("$[0].kind", v)
		if kind1 != nil && kind1 != "" {
			artifact["kind"] = kind1
		}

		request["artifact"] = artifact
	}

	if d.HasChange("flink_conf") {
		update = true
	}
	if v, ok := d.GetOk("flink_conf"); ok || d.HasChange("flink_conf") {
		request["flinkConf"] = v
	}
	if d.HasChange("logging") {
		update = true
	}
	logging := make(map[string]interface{})

	if v := d.Get("logging"); !IsNil(v) || d.HasChange("logging") {
		logReservePolicy := make(map[string]interface{})
		expirationDays1, _ := jsonpath.Get("$[0].log_reserve_policy[0].expiration_days", d.Get("logging"))
		if expirationDays1 != nil && expirationDays1 != "" {
			logReservePolicy["expirationDays"] = expirationDays1
		}
		openHistory1, _ := jsonpath.Get("$[0].log_reserve_policy[0].open_history", d.Get("logging"))
		if openHistory1 != nil && openHistory1 != "" {
			logReservePolicy["openHistory"] = openHistory1
		}

		if len(logReservePolicy) > 0 {
			logging["logReservePolicy"] = logReservePolicy
		}
		if v, ok := d.GetOk("logging"); ok {
			localData, err := jsonpath.Get("$[0].log4j_loggers", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range convertToInterfaceArray(localData) {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["loggerLevel"] = dataLoopTmp["logger_level"]
				dataLoopMap["loggerName"] = dataLoopTmp["logger_name"]
				localMaps = append(localMaps, dataLoopMap)
			}
			logging["log4jLoggers"] = localMaps
		}

		loggingProfile1, _ := jsonpath.Get("$[0].logging_profile", v)
		if loggingProfile1 != nil && loggingProfile1 != "" {
			logging["loggingProfile"] = loggingProfile1
		}
		log4J2ConfigurationTemplate, _ := jsonpath.Get("$[0].log4j2_configuration_template", v)
		if log4J2ConfigurationTemplate != nil && log4J2ConfigurationTemplate != "" {
			logging["log4j2ConfigurationTemplate"] = log4J2ConfigurationTemplate
		}

		request["logging"] = logging
	}

	if d.HasChange("batch_resource_setting") {
		update = true
	}
	batchResourceSetting := make(map[string]interface{})

	if v := d.Get("batch_resource_setting"); !IsNil(v) || d.HasChange("batch_resource_setting") {
		basicResourceSetting1 := make(map[string]interface{})
		taskmanagerResourceSettingSpec1 := make(map[string]interface{})
		cpu7, _ := jsonpath.Get("$[0].basic_resource_setting[0].taskmanager_resource_setting_spec[0].cpu", d.Get("batch_resource_setting"))
		if cpu7 != nil && cpu7 != "" {
			taskmanagerResourceSettingSpec1["cpu"] = cpu7
		}
		memory7, _ := jsonpath.Get("$[0].basic_resource_setting[0].taskmanager_resource_setting_spec[0].memory", d.Get("batch_resource_setting"))
		if memory7 != nil && memory7 != "" {
			taskmanagerResourceSettingSpec1["memory"] = memory7
		}

		if len(taskmanagerResourceSettingSpec1) > 0 {
			basicResourceSetting1["taskmanagerResourceSettingSpec"] = taskmanagerResourceSettingSpec1
		}
		parallelism3, _ := jsonpath.Get("$[0].basic_resource_setting[0].parallelism", d.Get("batch_resource_setting"))
		if parallelism3 != nil && parallelism3 != "" {
			basicResourceSetting1["parallelism"] = parallelism3
		}
		jobmanagerResourceSettingSpec2 := make(map[string]interface{})
		cpu9, _ := jsonpath.Get("$[0].basic_resource_setting[0].jobmanager_resource_setting_spec[0].cpu", d.Get("batch_resource_setting"))
		if cpu9 != nil && cpu9 != "" {
			jobmanagerResourceSettingSpec2["cpu"] = cpu9
		}
		memory9, _ := jsonpath.Get("$[0].basic_resource_setting[0].jobmanager_resource_setting_spec[0].memory", d.Get("batch_resource_setting"))
		if memory9 != nil && memory9 != "" {
			jobmanagerResourceSettingSpec2["memory"] = memory9
		}

		if len(jobmanagerResourceSettingSpec2) > 0 {
			basicResourceSetting1["jobmanagerResourceSettingSpec"] = jobmanagerResourceSettingSpec2
		}

		if len(basicResourceSetting1) > 0 {
			batchResourceSetting["basicResourceSetting"] = basicResourceSetting1
		}
		maxSlot1, _ := jsonpath.Get("$[0].max_slot", v)
		if maxSlot1 != nil && maxSlot1 != "" {
			batchResourceSetting["maxSlot"] = maxSlot1
		}

		request["batchResourceSetting"] = batchResourceSetting
	}

	if d.HasChange("labels") {
		update = true
	}
	if v, ok := d.GetOk("labels"); ok || d.HasChange("labels") {
		request["labels"] = v
	}
	if d.HasChange("deployment_target") {
		update = true
	}
	deploymentTarget := make(map[string]interface{})

	if v := d.Get("deployment_target"); v != nil {
		mode1, _ := jsonpath.Get("$[0].mode", v)
		if mode1 != nil && mode1 != "" {
			deploymentTarget["mode"] = mode1
		}
		name1, _ := jsonpath.Get("$[0].name", v)
		if name1 != nil && name1 != "" {
			deploymentTarget["name"] = name1
		}

		request["deploymentTarget"] = deploymentTarget
	}

	if d.HasChange("local_variables") {
		update = true
	}
	if v, ok := d.GetOk("local_variables"); ok || d.HasChange("local_variables") {
		localVariablesMapsArray := make([]interface{}, 0)
		for _, dataLoop1 := range convertToInterfaceArray(v) {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["value"] = dataLoop1Tmp["value"]
			dataLoop1Map["name"] = dataLoop1Tmp["name"]
			localVariablesMapsArray = append(localVariablesMapsArray, dataLoop1Map)
		}
		request["localVariables"] = localVariablesMapsArray
	}

	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["description"] = v
	}
	if d.HasChange("engine_version") {
		update = true
	}
	if v, ok := d.GetOk("engine_version"); ok || d.HasChange("engine_version") {
		request["engineVersion"] = v
	}
	if d.HasChange("deployment_name") {
		update = true
	}
	request["name"] = d.Get("deployment_name")
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("ververica", "2022-07-18", action, query, header, body, true)
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

	return resourceAliCloudRealtimeComputeDeploymentRead(d, meta)
}

func resourceAliCloudRealtimeComputeDeploymentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	namespace := parts[1]
	deploymentId := parts[2]
	action := fmt.Sprintf("/api/v2/namespaces/%s/deployments/%s", namespace, deploymentId)
	var request map[string]interface{}
	var response map[string]interface{}
	header := make(map[string]*string)
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	header["workspace"] = StringPointer(parts[0])

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("ververica", "2022-07-18", action, query, header, nil, true)
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
		if IsExpectedErrors(err, []string{"990301"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
