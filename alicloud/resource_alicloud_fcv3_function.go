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

func resourceAliCloudFcv3Function() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudFcv3FunctionCreate,
		Read:   resourceAliCloudFcv3FunctionRead,
		Update: resourceAliCloudFcv3FunctionUpdate,
		Delete: resourceAliCloudFcv3FunctionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
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
						"checksum": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
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
			"code_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cpu": {
				Type:         schema.TypeFloat,
				Optional:     true,
				Computed:     true,
				ValidateFunc: FloatBetween(0.05, 16),
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
						"resolved_image_uri": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"entrypoint": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"acr_instance_id": {
							Type:       schema.TypeString,
							Optional:   true,
							Deprecated: "Field 'acr_instance_id' has been deprecated from provider version 1.228.0. ACR Enterprise version Image Repository ID, which must be entered when using ACR Enterprise version image. (Obsolete)",
						},
						"acceleration_info": {
							Type:       schema.TypeList,
							Computed:   true,
							Deprecated: "Field 'acceleration_info' has been deprecated from provider version 1.228.0. Image Acceleration Information (Obsolete)",
							MaxItems:   1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": {
										Type:       schema.TypeString,
										Computed:   true,
										Deprecated: "Field 'status' has been deprecated from provider version 1.228.0. Image Acceleration Status (Deprecated)",
									},
								},
							},
						},
						"command": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"acceleration_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Deprecated:   "Field 'acceleration_type' has been deprecated from provider version 1.228.0. Whether to enable Image acceleration. Default: The Default value, indicating that image acceleration is enabled. None: indicates that image acceleration is disabled. (Obsolete)",
							ValidateFunc: StringInSlice([]string{"Default", "None"}, false),
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"health_check_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"initial_delay_seconds": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: IntBetween(0, 120),
									},
									"timeout_seconds": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: IntBetween(0, 3),
									},
									"http_get_url": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"period_seconds": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: IntBetween(0, 120),
									},
									"failure_threshold": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: IntBetween(0, 120),
									},
									"success_threshold": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: IntBetween(0, 120),
									},
								},
							},
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
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(0, 65535),
						},
						"health_check_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"initial_delay_seconds": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: IntBetween(0, 120),
									},
									"timeout_seconds": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: IntBetween(0, 3),
									},
									"http_get_url": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"period_seconds": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: IntBetween(0, 120),
									},
									"failure_threshold": {
										Type:         schema.TypeInt,
										Optional:     true,
										Computed:     true,
										ValidateFunc: IntBetween(0, 120),
									},
									"success_threshold": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: IntBetween(0, 120),
									},
								},
							},
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntAtLeast(512),
			},
			"environment_variables": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"function_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"function_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"function_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[0-9a-zA-Z_-]+$"), "The function name. Consists of uppercase and lowercase letters, digits (0 to 9), underscores (_), and dashes (-). It must begin with an English letter (a ~ z), (A ~ Z), or an underscore (_). Case sensitive. The length is 1~128 characters."),
			},
			"gpu_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gpu_memory_size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"gpu_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"fc.gpu.tesla.1", "fc.gpu.ampere.1", "fc.gpu.ada.1", "g1"}, false),
						},
					},
				},
			},
			"handler": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_concurrency": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(0, 200),
			},
			"instance_isolation_mode": {
				Type:     schema.TypeString,
				Optional: true,
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
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: IntBetween(0, 900),
									},
									"handler": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"initializer": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"command": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"timeout": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: IntBetween(0, 600),
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
			"internet_access": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"invocation_restriction": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"last_modified_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"reason": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"last_modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_update_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_update_status_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_update_status_reason_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"layers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"log_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"log_begin_rule": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"None", "DefaultRegex"}, false),
						},
						"logstore": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enable_instance_metrics": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"enable_request_metrics": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"memory_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(64, 32768),
			},
			"nas_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount_points": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable_tls": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"server_addr": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"mount_dir": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"user_id": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"oss_mount_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount_points": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"read_only": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"bucket_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"endpoint": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"bucket_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"mount_dir": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"role": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"runtime": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"python3.10", "python3.9", "python3", "nodejs20", "nodejs18", "nodejs16", "nodejs14", "java11", "java8", "php7.2", "dotnetcore3.1", "go1", "custom.debian10", "custom", "custom-container", "python3.12", "custom.debian11", "custom.debian12"}, false),
			},
			"session_affinity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"session_affinity_config": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state_reason_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(0, 86400),
			},
			"tracing_config": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"params": {
							Type:     schema.TypeMap,
							Computed: true,
						},
					},
				},
			},
			"vpc_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"vswitch_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudFcv3FunctionCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/2023-03-30/functions")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("function_name"); ok {
		request["functionName"] = v
	}

	dataList := make(map[string]interface{})

	if v := d.Get("gpu_config"); !IsNil(v) {
		gpuMemorySize1, _ := jsonpath.Get("$[0].gpu_memory_size", v)
		if gpuMemorySize1 != nil && gpuMemorySize1 != "" {
			dataList["gpuMemorySize"] = gpuMemorySize1
		}
		gpuType1, _ := jsonpath.Get("$[0].gpu_type", v)
		if gpuType1 != nil && gpuType1 != "" {
			dataList["gpuType"] = gpuType1
		}

		request["gpuConfig"] = dataList
	}

	dataList1 := make(map[string]interface{})

	if v := d.Get("log_config"); !IsNil(v) {
		logBeginRule1, _ := jsonpath.Get("$[0].log_begin_rule", v)
		if logBeginRule1 != nil && logBeginRule1 != "" {
			dataList1["logBeginRule"] = logBeginRule1
		}
		project1, _ := jsonpath.Get("$[0].project", v)
		if project1 != nil && project1 != "" {
			dataList1["project"] = project1
		}
		enableInstanceMetrics1, _ := jsonpath.Get("$[0].enable_instance_metrics", v)
		if enableInstanceMetrics1 != nil && enableInstanceMetrics1 != "" {
			dataList1["enableInstanceMetrics"] = enableInstanceMetrics1
		}
		enableRequestMetrics1, _ := jsonpath.Get("$[0].enable_request_metrics", v)
		if enableRequestMetrics1 != nil && enableRequestMetrics1 != "" {
			dataList1["enableRequestMetrics"] = enableRequestMetrics1
		}
		logstore1, _ := jsonpath.Get("$[0].logstore", v)
		if logstore1 != nil && logstore1 != "" {
			dataList1["logstore"] = logstore1
		}

		request["logConfig"] = dataList1
	}

	if v, ok := d.GetOkExists("instance_concurrency"); ok && v.(int) > 0 {
		request["instanceConcurrency"] = v
	}
	if v, ok := d.GetOk("instance_isolation_mode"); ok {
		request["instanceIsolationMode"] = v
	}
	dataList2 := make(map[string]interface{})

	if v := d.Get("custom_runtime_config"); !IsNil(v) {
		healthCheckConfig := make(map[string]interface{})
		timeoutSeconds1, _ := jsonpath.Get("$[0].health_check_config[0].timeout_seconds", v)
		if timeoutSeconds1 != nil && timeoutSeconds1 != "" && timeoutSeconds1.(int) > 0 {
			healthCheckConfig["timeoutSeconds"] = timeoutSeconds1
		}
		httpGetUrl1, _ := jsonpath.Get("$[0].health_check_config[0].http_get_url", v)
		if httpGetUrl1 != nil && httpGetUrl1 != "" {
			healthCheckConfig["httpGetUrl"] = httpGetUrl1
		}
		successThreshold1, _ := jsonpath.Get("$[0].health_check_config[0].success_threshold", v)
		if successThreshold1 != nil && successThreshold1 != "" && successThreshold1.(int) > 0 {
			healthCheckConfig["successThreshold"] = successThreshold1
		}
		initialDelaySeconds1, _ := jsonpath.Get("$[0].health_check_config[0].initial_delay_seconds", v)
		if initialDelaySeconds1 != nil && initialDelaySeconds1 != "" {
			healthCheckConfig["initialDelaySeconds"] = initialDelaySeconds1
		}
		periodSeconds1, _ := jsonpath.Get("$[0].health_check_config[0].period_seconds", v)
		if periodSeconds1 != nil && periodSeconds1 != "" && periodSeconds1.(int) > 0 {
			healthCheckConfig["periodSeconds"] = periodSeconds1
		}
		failureThreshold1, _ := jsonpath.Get("$[0].health_check_config[0].failure_threshold", v)
		if failureThreshold1 != nil && failureThreshold1 != "" && failureThreshold1.(int) > 0 {
			healthCheckConfig["failureThreshold"] = failureThreshold1
		}

		dataList2["healthCheckConfig"] = healthCheckConfig
		args1, _ := jsonpath.Get("$[0].args", v)
		if args1 != nil && args1 != "" {
			dataList2["args"] = args1
		}
		command1, _ := jsonpath.Get("$[0].command", v)
		if command1 != nil && command1 != "" {
			dataList2["command"] = command1
		}
		port1, _ := jsonpath.Get("$[0].port", v)
		if port1 != nil && port1 != "" && port1.(int) > 0 {
			dataList2["port"] = port1
		}

		request["customRuntimeConfig"] = dataList2
	}

	dataList3 := make(map[string]interface{})

	if v := d.Get("custom_container_config"); !IsNil(v) {
		accelerationType1, _ := jsonpath.Get("$[0].acceleration_type", v)
		if accelerationType1 != nil && accelerationType1 != "" {
			dataList3["accelerationType"] = accelerationType1
		}
		healthCheckConfig1 := make(map[string]interface{})
		failureThreshold3, _ := jsonpath.Get("$[0].health_check_config[0].failure_threshold", v)
		if failureThreshold3 != nil && failureThreshold3 != "" && failureThreshold3.(int) > 0 {
			healthCheckConfig1["failureThreshold"] = failureThreshold3
		}
		timeoutSeconds3, _ := jsonpath.Get("$[0].health_check_config[0].timeout_seconds", v)
		if timeoutSeconds3 != nil && timeoutSeconds3 != "" && timeoutSeconds3.(int) > 0 {
			healthCheckConfig1["timeoutSeconds"] = timeoutSeconds3
		}
		initialDelaySeconds3, _ := jsonpath.Get("$[0].health_check_config[0].initial_delay_seconds", v)
		if initialDelaySeconds3 != nil && initialDelaySeconds3 != "" {
			healthCheckConfig1["initialDelaySeconds"] = initialDelaySeconds3
		}
		periodSeconds3, _ := jsonpath.Get("$[0].health_check_config[0].period_seconds", v)
		if periodSeconds3 != nil && periodSeconds3 != "" && periodSeconds3.(int) > 0 {
			healthCheckConfig1["periodSeconds"] = periodSeconds3
		}
		httpGetUrl3, _ := jsonpath.Get("$[0].health_check_config[0].http_get_url", v)
		if httpGetUrl3 != nil && httpGetUrl3 != "" {
			healthCheckConfig1["httpGetUrl"] = httpGetUrl3
		}
		successThreshold3, _ := jsonpath.Get("$[0].health_check_config[0].success_threshold", v)
		if successThreshold3 != nil && successThreshold3 != "" && successThreshold3.(int) > 0 {
			healthCheckConfig1["successThreshold"] = successThreshold3
		}

		dataList3["healthCheckConfig"] = healthCheckConfig1
		entrypoint1, _ := jsonpath.Get("$[0].entrypoint", v)
		if entrypoint1 != nil && entrypoint1 != "" {
			dataList3["entrypoint"] = entrypoint1
		}
		command3, _ := jsonpath.Get("$[0].command", v)
		if command3 != nil && command3 != "" {
			dataList3["command"] = command3
		}
		image1, _ := jsonpath.Get("$[0].image", v)
		if image1 != nil && image1 != "" {
			dataList3["image"] = image1
		}
		port3, _ := jsonpath.Get("$[0].port", v)
		if port3 != nil && port3 != "" {
			dataList3["port"] = port3
		}
		acrInstanceId1, _ := jsonpath.Get("$[0].acr_instance_id", v)
		if acrInstanceId1 != nil && acrInstanceId1 != "" {
			dataList3["acrInstanceId"] = acrInstanceId1
		}

		request["customContainerConfig"] = dataList3
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["resourceGroupId"] = v
	}
	dataList4 := make(map[string]interface{})

	if v := d.Get("instance_lifecycle_config"); !IsNil(v) {
		preStop := make(map[string]interface{})
		timeout1, _ := jsonpath.Get("$[0].pre_stop[0].timeout", v)
		if timeout1 != nil && timeout1 != "" && timeout1.(int) > 0 {
			preStop["timeout"] = timeout1
		}
		handler1, _ := jsonpath.Get("$[0].pre_stop[0].handler", v)
		if handler1 != nil && handler1 != "" {
			preStop["handler"] = handler1
		}

		dataList4["preStop"] = preStop
		initializer := make(map[string]interface{})
		command5, _ := jsonpath.Get("$[0].initializer[0].command", v)
		if command5 != nil && command5 != "" {
			initializer["command"] = command5
		}
		handler3, _ := jsonpath.Get("$[0].initializer[0].handler", v)
		if handler3 != nil && handler3 != "" {
			initializer["handler"] = handler3
		}
		timeout3, _ := jsonpath.Get("$[0].initializer[0].timeout", v)
		if timeout3 != nil && timeout3 != "" && timeout3.(int) > 0 {
			initializer["timeout"] = timeout3
		}

		dataList4["initializer"] = initializer

		request["instanceLifecycleConfig"] = dataList4
	}

	if v, ok := d.GetOkExists("internet_access"); ok {
		request["internetAccess"] = v
	}
	dataList5 := make(map[string]interface{})

	if v := d.Get("oss_mount_config"); !IsNil(v) {
		if v, ok := d.GetOk("oss_mount_config"); ok {
			localData, err := jsonpath.Get("$[0].mount_points", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range localData.([]interface{}) {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["bucketName"] = dataLoopTmp["bucket_name"]
				dataLoopMap["readOnly"] = dataLoopTmp["read_only"]
				dataLoopMap["bucketPath"] = dataLoopTmp["bucket_path"]
				dataLoopMap["mountDir"] = dataLoopTmp["mount_dir"]
				dataLoopMap["endpoint"] = dataLoopTmp["endpoint"]
				localMaps = append(localMaps, dataLoopMap)
			}
			dataList5["mountPoints"] = localMaps
		}

		request["ossMountConfig"] = dataList5
	}

	request["runtime"] = d.Get("runtime")
	if v, ok := d.GetOk("environment_variables"); ok {
		request["environmentVariables"] = v
	}
	dataList6 := make(map[string]interface{})

	if v := d.Get("code"); !IsNil(v) {
		ossBucketName1, _ := jsonpath.Get("$[0].oss_bucket_name", v)
		if ossBucketName1 != nil && ossBucketName1 != "" {
			dataList6["ossBucketName"] = ossBucketName1
		}
		zipFile1, _ := jsonpath.Get("$[0].zip_file", v)
		if zipFile1 != nil && zipFile1 != "" {
			dataList6["zipFile"] = zipFile1
		}
		ossObjectName1, _ := jsonpath.Get("$[0].oss_object_name", v)
		if ossObjectName1 != nil && ossObjectName1 != "" {
			dataList6["ossObjectName"] = ossObjectName1
		}
		checksum1, _ := jsonpath.Get("$[0].checksum", v)
		if checksum1 != nil && checksum1 != "" {
			dataList6["checksum"] = checksum1
		}

		request["code"] = dataList6
	}

	if v, ok := d.GetOk("role"); ok {
		request["role"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request["Tags"] = tagsMap
	}

	if v, ok := d.GetOk("layers"); ok {
		layersMapsArray := v.([]interface{})
		request["layers"] = layersMapsArray
	}

	if v, ok := d.GetOkExists("timeout"); ok && v.(int) > 0 {
		request["timeout"] = v
	}
	if v, ok := d.GetOk("session_affinity"); ok {
		request["sessionAffinity"] = v
	}
	if v, ok := d.GetOk("session_affinity_config"); ok {
		request["sessionAffinityConfig"] = v
	}
	if v, ok := d.GetOk("cpu"); ok && v.(float64) > 0 {
		request["cpu"] = v
	}
	if v, ok := d.GetOkExists("disk_size"); ok && v.(int) > 0 {
		request["diskSize"] = v
	}
	dataList7 := make(map[string]interface{})

	if v := d.Get("vpc_config"); !IsNil(v) {
		vpcId1, _ := jsonpath.Get("$[0].vpc_id", v)
		if vpcId1 != nil && vpcId1 != "" {
			dataList7["vpcId"] = vpcId1
		}
		securityGroupId1, _ := jsonpath.Get("$[0].security_group_id", v)
		if securityGroupId1 != nil && securityGroupId1 != "" {
			dataList7["securityGroupId"] = securityGroupId1
		}
		vSwitchIds1, _ := jsonpath.Get("$[0].vswitch_ids", v)
		if vSwitchIds1 != nil && vSwitchIds1 != "" {
			dataList7["vSwitchIds"] = vSwitchIds1
		}

		request["vpcConfig"] = dataList7
	}

	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	request["handler"] = d.Get("handler")
	if v, ok := d.GetOkExists("memory_size"); ok && v.(int) > 0 {
		request["memorySize"] = v
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("FC", "2023-03-30", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_fcv3_function", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["functionName"]))

	return resourceAliCloudFcv3FunctionUpdate(d, meta)
}

func resourceAliCloudFcv3FunctionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fcv3ServiceV2 := Fcv3ServiceV2{client}

	objectRaw, err := fcv3ServiceV2.DescribeFcv3Function(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_fcv3_function DescribeFcv3Function Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("code_size", objectRaw["codeSize"])
	d.Set("cpu", objectRaw["cpu"])
	d.Set("create_time", objectRaw["createdTime"])
	d.Set("description", objectRaw["description"])
	d.Set("disk_size", objectRaw["diskSize"])
	d.Set("environment_variables", objectRaw["environmentVariables"])
	d.Set("function_arn", objectRaw["functionArn"])
	d.Set("function_id", objectRaw["functionId"])
	d.Set("handler", objectRaw["handler"])
	d.Set("instance_concurrency", objectRaw["instanceConcurrency"])
	d.Set("instance_isolation_mode", objectRaw["instanceIsolationMode"])
	d.Set("internet_access", objectRaw["internetAccess"])
	d.Set("last_modified_time", objectRaw["lastModifiedTime"])
	d.Set("last_update_status", objectRaw["lastUpdateStatus"])
	d.Set("last_update_status_reason", objectRaw["lastUpdateStatusReason"])
	d.Set("last_update_status_reason_code", objectRaw["lastUpdateStatusReasonCode"])
	d.Set("memory_size", objectRaw["memorySize"])
	d.Set("resource_group_id", objectRaw["resourceGroupId"])
	d.Set("role", objectRaw["role"])
	d.Set("runtime", objectRaw["runtime"])
	d.Set("session_affinity", objectRaw["sessionAffinity"])
	d.Set("session_affinity_config", objectRaw["sessionAffinityConfig"])
	d.Set("state", objectRaw["state"])
	d.Set("state_reason", objectRaw["stateReason"])
	d.Set("state_reason_code", objectRaw["stateReasonCode"])
	d.Set("timeout", objectRaw["timeout"])

	customContainerConfigMaps := make([]map[string]interface{}, 0)
	customContainerConfigMap := make(map[string]interface{})
	customContainerConfigRaw := make(map[string]interface{})
	if objectRaw["customContainerConfig"] != nil {
		customContainerConfigRaw = objectRaw["customContainerConfig"].(map[string]interface{})
	}
	if len(customContainerConfigRaw) > 0 {
		customContainerConfigMap["acceleration_type"] = customContainerConfigRaw["accelerationType"]
		customContainerConfigMap["acr_instance_id"] = customContainerConfigRaw["acrInstanceId"]
		customContainerConfigMap["image"] = customContainerConfigRaw["image"]
		customContainerConfigMap["port"] = customContainerConfigRaw["port"]
		customContainerConfigMap["resolved_image_uri"] = customContainerConfigRaw["resolvedImageUri"]

		accelerationInfoMaps := make([]map[string]interface{}, 0)
		accelerationInfoMap := make(map[string]interface{})
		accelerationInfoRaw := make(map[string]interface{})
		if customContainerConfigRaw["accelerationInfo"] != nil {
			accelerationInfoRaw = customContainerConfigRaw["accelerationInfo"].(map[string]interface{})
		}
		if len(accelerationInfoRaw) > 0 {
			accelerationInfoMap["status"] = accelerationInfoRaw["status"]

			accelerationInfoMaps = append(accelerationInfoMaps, accelerationInfoMap)
		}
		customContainerConfigMap["acceleration_info"] = accelerationInfoMaps
		commandRaw := make([]interface{}, 0)
		if customContainerConfigRaw["command"] != nil {
			commandRaw = customContainerConfigRaw["command"].([]interface{})
		}

		customContainerConfigMap["command"] = commandRaw
		entrypointRaw := make([]interface{}, 0)
		if customContainerConfigRaw["entrypoint"] != nil {
			entrypointRaw = customContainerConfigRaw["entrypoint"].([]interface{})
		}

		customContainerConfigMap["entrypoint"] = entrypointRaw
		healthCheckConfigMaps := make([]map[string]interface{}, 0)
		healthCheckConfigMap := make(map[string]interface{})
		healthCheckConfigRaw := make(map[string]interface{})
		if customContainerConfigRaw["healthCheckConfig"] != nil {
			healthCheckConfigRaw = customContainerConfigRaw["healthCheckConfig"].(map[string]interface{})
		}
		if len(healthCheckConfigRaw) > 0 {
			healthCheckConfigMap["failure_threshold"] = healthCheckConfigRaw["failureThreshold"]
			healthCheckConfigMap["http_get_url"] = healthCheckConfigRaw["httpGetUrl"]
			healthCheckConfigMap["initial_delay_seconds"] = healthCheckConfigRaw["initialDelaySeconds"]
			healthCheckConfigMap["period_seconds"] = healthCheckConfigRaw["periodSeconds"]
			healthCheckConfigMap["success_threshold"] = healthCheckConfigRaw["successThreshold"]
			healthCheckConfigMap["timeout_seconds"] = healthCheckConfigRaw["timeoutSeconds"]

			healthCheckConfigMaps = append(healthCheckConfigMaps, healthCheckConfigMap)
		}
		customContainerConfigMap["health_check_config"] = healthCheckConfigMaps
		customContainerConfigMaps = append(customContainerConfigMaps, customContainerConfigMap)
	}
	if err := d.Set("custom_container_config", customContainerConfigMaps); err != nil {
		return err
	}
	customDnsMaps := make([]map[string]interface{}, 0)
	customDnsMap := make(map[string]interface{})
	customDNSRaw := make(map[string]interface{})
	if objectRaw["customDNS"] != nil {
		customDNSRaw = objectRaw["customDNS"].(map[string]interface{})
	}
	if len(customDNSRaw) > 0 {

		dnsOptionsRaw := customDNSRaw["dnsOptions"]
		dnsOptionsMaps := make([]map[string]interface{}, 0)
		if dnsOptionsRaw != nil {
			for _, dnsOptionsChildRaw := range dnsOptionsRaw.([]interface{}) {
				dnsOptionsMap := make(map[string]interface{})
				dnsOptionsChildRaw := dnsOptionsChildRaw.(map[string]interface{})
				dnsOptionsMap["name"] = dnsOptionsChildRaw["name"]
				dnsOptionsMap["value"] = dnsOptionsChildRaw["value"]

				dnsOptionsMaps = append(dnsOptionsMaps, dnsOptionsMap)
			}
		}
		customDnsMap["dns_options"] = dnsOptionsMaps
		nameServersRaw := make([]interface{}, 0)
		if customDNSRaw["nameServers"] != nil {
			nameServersRaw = customDNSRaw["nameServers"].([]interface{})
		}

		customDnsMap["name_servers"] = nameServersRaw
		searchesRaw := make([]interface{}, 0)
		if customDNSRaw["searches"] != nil {
			searchesRaw = customDNSRaw["searches"].([]interface{})
		}

		customDnsMap["searches"] = searchesRaw
		customDnsMaps = append(customDnsMaps, customDnsMap)
	}
	if err := d.Set("custom_dns", customDnsMaps); err != nil {
		return err
	}
	customRuntimeConfigMaps := make([]map[string]interface{}, 0)
	customRuntimeConfigMap := make(map[string]interface{})
	customRuntimeConfigRaw := make(map[string]interface{})
	if objectRaw["customRuntimeConfig"] != nil {
		customRuntimeConfigRaw = objectRaw["customRuntimeConfig"].(map[string]interface{})
	}
	if len(customRuntimeConfigRaw) > 0 {
		customRuntimeConfigMap["port"] = customRuntimeConfigRaw["port"]

		argsRaw := make([]interface{}, 0)
		if customRuntimeConfigRaw["args"] != nil {
			argsRaw = customRuntimeConfigRaw["args"].([]interface{})
		}

		customRuntimeConfigMap["args"] = argsRaw
		commandRaw := make([]interface{}, 0)
		if customRuntimeConfigRaw["command"] != nil {
			commandRaw = customRuntimeConfigRaw["command"].([]interface{})
		}

		customRuntimeConfigMap["command"] = commandRaw
		healthCheckConfigMaps := make([]map[string]interface{}, 0)
		healthCheckConfigMap := make(map[string]interface{})
		healthCheckConfigRaw := make(map[string]interface{})
		if customRuntimeConfigRaw["healthCheckConfig"] != nil {
			healthCheckConfigRaw = customRuntimeConfigRaw["healthCheckConfig"].(map[string]interface{})
		}
		if len(healthCheckConfigRaw) > 0 {
			healthCheckConfigMap["failure_threshold"] = healthCheckConfigRaw["failureThreshold"]
			healthCheckConfigMap["http_get_url"] = healthCheckConfigRaw["httpGetUrl"]
			healthCheckConfigMap["initial_delay_seconds"] = healthCheckConfigRaw["initialDelaySeconds"]
			healthCheckConfigMap["period_seconds"] = healthCheckConfigRaw["periodSeconds"]
			healthCheckConfigMap["success_threshold"] = healthCheckConfigRaw["successThreshold"]
			healthCheckConfigMap["timeout_seconds"] = healthCheckConfigRaw["timeoutSeconds"]

			healthCheckConfigMaps = append(healthCheckConfigMaps, healthCheckConfigMap)
		}
		customRuntimeConfigMap["health_check_config"] = healthCheckConfigMaps
		customRuntimeConfigMaps = append(customRuntimeConfigMaps, customRuntimeConfigMap)
	}
	if err := d.Set("custom_runtime_config", customRuntimeConfigMaps); err != nil {
		return err
	}
	gpuConfigMaps := make([]map[string]interface{}, 0)
	gpuConfigMap := make(map[string]interface{})
	gpuConfigRaw := make(map[string]interface{})
	if objectRaw["gpuConfig"] != nil {
		gpuConfigRaw = objectRaw["gpuConfig"].(map[string]interface{})
	}
	if len(gpuConfigRaw) > 0 {
		gpuConfigMap["gpu_memory_size"] = gpuConfigRaw["gpuMemorySize"]
		gpuConfigMap["gpu_type"] = gpuConfigRaw["gpuType"]

		gpuConfigMaps = append(gpuConfigMaps, gpuConfigMap)
	}
	if err := d.Set("gpu_config", gpuConfigMaps); err != nil {
		return err
	}
	instanceLifecycleConfigMaps := make([]map[string]interface{}, 0)
	instanceLifecycleConfigMap := make(map[string]interface{})
	instanceLifecycleConfigRaw := make(map[string]interface{})
	if objectRaw["instanceLifecycleConfig"] != nil {
		instanceLifecycleConfigRaw = objectRaw["instanceLifecycleConfig"].(map[string]interface{})
	}
	if len(instanceLifecycleConfigRaw) > 0 {

		initializerMaps := make([]map[string]interface{}, 0)
		initializerMap := make(map[string]interface{})
		initializerRaw := make(map[string]interface{})
		if instanceLifecycleConfigRaw["initializer"] != nil {
			initializerRaw = instanceLifecycleConfigRaw["initializer"].(map[string]interface{})
		}
		if len(initializerRaw) > 0 {
			initializerMap["handler"] = initializerRaw["handler"]
			initializerMap["timeout"] = initializerRaw["timeout"]

			commandRaw := make([]interface{}, 0)
			if initializerRaw["command"] != nil {
				commandRaw = initializerRaw["command"].([]interface{})
			}

			initializerMap["command"] = commandRaw
			initializerMaps = append(initializerMaps, initializerMap)
		}
		instanceLifecycleConfigMap["initializer"] = initializerMaps
		preStopMaps := make([]map[string]interface{}, 0)
		preStopMap := make(map[string]interface{})
		preStopRaw := make(map[string]interface{})
		if instanceLifecycleConfigRaw["preStop"] != nil {
			preStopRaw = instanceLifecycleConfigRaw["preStop"].(map[string]interface{})
		}
		if len(preStopRaw) > 0 {
			preStopMap["handler"] = preStopRaw["handler"]
			preStopMap["timeout"] = preStopRaw["timeout"]

			preStopMaps = append(preStopMaps, preStopMap)
		}
		instanceLifecycleConfigMap["pre_stop"] = preStopMaps
		instanceLifecycleConfigMaps = append(instanceLifecycleConfigMaps, instanceLifecycleConfigMap)
	}
	if err := d.Set("instance_lifecycle_config", instanceLifecycleConfigMaps); err != nil {
		return err
	}
	invocationRestrictionMaps := make([]map[string]interface{}, 0)
	invocationRestrictionMap := make(map[string]interface{})
	invocationRestrictionRaw := make(map[string]interface{})
	if objectRaw["invocationRestriction"] != nil {
		invocationRestrictionRaw = objectRaw["invocationRestriction"].(map[string]interface{})
	}
	if len(invocationRestrictionRaw) > 0 {
		invocationRestrictionMap["disable"] = invocationRestrictionRaw["disable"]
		invocationRestrictionMap["last_modified_time"] = invocationRestrictionRaw["lastModifiedTime"]
		invocationRestrictionMap["reason"] = invocationRestrictionRaw["reason"]

		invocationRestrictionMaps = append(invocationRestrictionMaps, invocationRestrictionMap)
	}
	if err := d.Set("invocation_restriction", invocationRestrictionMaps); err != nil {
		return err
	}
	logConfigMaps := make([]map[string]interface{}, 0)
	logConfigMap := make(map[string]interface{})
	logConfigRaw := make(map[string]interface{})
	if objectRaw["logConfig"] != nil {
		logConfigRaw = objectRaw["logConfig"].(map[string]interface{})
	}
	if len(logConfigRaw) > 0 {
		logConfigMap["enable_instance_metrics"] = logConfigRaw["enableInstanceMetrics"]
		logConfigMap["enable_request_metrics"] = logConfigRaw["enableRequestMetrics"]
		logConfigMap["log_begin_rule"] = logConfigRaw["logBeginRule"]
		logConfigMap["logstore"] = logConfigRaw["logstore"]
		logConfigMap["project"] = logConfigRaw["project"]

		logConfigMaps = append(logConfigMaps, logConfigMap)
	}
	if err := d.Set("log_config", logConfigMaps); err != nil {
		return err
	}
	nasConfigMaps := make([]map[string]interface{}, 0)
	nasConfigMap := make(map[string]interface{})
	nasConfigRaw := make(map[string]interface{})
	if objectRaw["nasConfig"] != nil {
		nasConfigRaw = objectRaw["nasConfig"].(map[string]interface{})
	}
	if len(nasConfigRaw) > 0 {
		nasConfigMap["group_id"] = nasConfigRaw["groupId"]
		nasConfigMap["user_id"] = nasConfigRaw["userId"]

		mountPointsRaw := nasConfigRaw["mountPoints"]
		mountPointsMaps := make([]map[string]interface{}, 0)
		if mountPointsRaw != nil {
			for _, mountPointsChildRaw := range mountPointsRaw.([]interface{}) {
				mountPointsMap := make(map[string]interface{})
				mountPointsChildRaw := mountPointsChildRaw.(map[string]interface{})
				mountPointsMap["enable_tls"] = mountPointsChildRaw["enableTLS"]
				mountPointsMap["mount_dir"] = mountPointsChildRaw["mountDir"]
				mountPointsMap["server_addr"] = mountPointsChildRaw["serverAddr"]

				mountPointsMaps = append(mountPointsMaps, mountPointsMap)
			}
		}
		nasConfigMap["mount_points"] = mountPointsMaps
		nasConfigMaps = append(nasConfigMaps, nasConfigMap)
	}
	if err := d.Set("nas_config", nasConfigMaps); err != nil {
		return err
	}
	ossMountConfigMaps := make([]map[string]interface{}, 0)
	ossMountConfigMap := make(map[string]interface{})
	mountPointsRaw, _ := jsonpath.Get("$.ossMountConfig.mountPoints", objectRaw)

	mountPointsMaps := make([]map[string]interface{}, 0)
	if mountPointsRaw != nil {
		for _, mountPointsChildRaw := range mountPointsRaw.([]interface{}) {
			mountPointsMap := make(map[string]interface{})
			mountPointsChildRaw := mountPointsChildRaw.(map[string]interface{})
			mountPointsMap["bucket_name"] = mountPointsChildRaw["bucketName"]
			mountPointsMap["bucket_path"] = mountPointsChildRaw["bucketPath"]
			mountPointsMap["endpoint"] = mountPointsChildRaw["endpoint"]
			mountPointsMap["mount_dir"] = mountPointsChildRaw["mountDir"]
			mountPointsMap["read_only"] = mountPointsChildRaw["readOnly"]

			mountPointsMaps = append(mountPointsMaps, mountPointsMap)
		}
	}
	ossMountConfigMap["mount_points"] = mountPointsMaps
	ossMountConfigMaps = append(ossMountConfigMaps, ossMountConfigMap)
	if err := d.Set("oss_mount_config", ossMountConfigMaps); err != nil {
		return err
	}
	tagsMaps := objectRaw["tags"]
	d.Set("tags", tagsToMap(tagsMaps))
	tracingConfigMaps := make([]map[string]interface{}, 0)
	tracingConfigMap := make(map[string]interface{})
	tracingConfigRaw := make(map[string]interface{})
	if objectRaw["tracingConfig"] != nil {
		tracingConfigRaw = objectRaw["tracingConfig"].(map[string]interface{})
	}
	if len(tracingConfigRaw) > 0 {
		tracingConfigMap["params"] = tracingConfigRaw["params"]
		tracingConfigMap["type"] = tracingConfigRaw["type"]

		tracingConfigMaps = append(tracingConfigMaps, tracingConfigMap)
	}
	if err := d.Set("tracing_config", tracingConfigMaps); err != nil {
		return err
	}
	vpcConfigMaps := make([]map[string]interface{}, 0)
	vpcConfigMap := make(map[string]interface{})
	vpcConfigRaw := make(map[string]interface{})
	if objectRaw["vpcConfig"] != nil {
		vpcConfigRaw = objectRaw["vpcConfig"].(map[string]interface{})
	}
	if len(vpcConfigRaw) > 0 {
		vpcConfigMap["security_group_id"] = vpcConfigRaw["securityGroupId"]
		vpcConfigMap["vpc_id"] = vpcConfigRaw["vpcId"]

		vSwitchIdsRaw := make([]interface{}, 0)
		if vpcConfigRaw["vSwitchIds"] != nil {
			vSwitchIdsRaw = vpcConfigRaw["vSwitchIds"].([]interface{})
		}

		vpcConfigMap["vswitch_ids"] = vSwitchIdsRaw
		vpcConfigMaps = append(vpcConfigMaps, vpcConfigMap)
	}
	if err := d.Set("vpc_config", vpcConfigMaps); err != nil {
		return err
	}

	d.Set("function_name", d.Id())

	return nil
}

func resourceAliCloudFcv3FunctionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	d.Partial(true)

	if d.HasChange("invocation_restriction.0.disable") {
		var err error
		fcv3ServiceV2 := Fcv3ServiceV2{client}
		object, err := fcv3ServiceV2.DescribeFcv3Function(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("invocation_restriction.0.disable").(bool)
		disable, _ := jsonpath.Get("$.invocationRestriction.disable", object)
		if formatBool(disable) != target {
			if target == true {
				functionName := d.Id()
				action := fmt.Sprintf("/2023-03-30/functions/%s/invoke/disable", functionName)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})
				request["functionName"] = d.Id()

				if v, ok := d.GetOk("invocation_restriction"); ok {
					invocationRestrictionReasonJsonPath, err := jsonpath.Get("$[0].reason", v)
					if err == nil && invocationRestrictionReasonJsonPath != "" {
						request["reason"] = invocationRestrictionReasonJsonPath
					}
				}
				body = request
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RoaPost("FC", "2023-03-30", action, query, nil, body, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"500"}) || NeedRetry(err) {
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
			if target == false {
				functionName := d.Id()
				action := fmt.Sprintf("/2023-03-30/functions/%s/invoke/enable", functionName)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})
				request["functionName"] = d.Id()

				body = request
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RoaPost("FC", "2023-03-30", action, query, nil, body, true)
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
	}

	var err error
	functionName := d.Id()
	action := fmt.Sprintf("/2023-03-30/functions/%s", functionName)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["functionName"] = d.Id()

	if !d.IsNewResource() && d.HasChange("gpu_config") {
		update = true
		dataList := make(map[string]interface{})

		if v := d.Get("gpu_config"); v != nil {
			gpuMemorySize1, _ := jsonpath.Get("$[0].gpu_memory_size", v)
			if gpuMemorySize1 != nil && (d.HasChange("gpu_config.0.gpu_memory_size") || gpuMemorySize1 != "") {
				dataList["gpuMemorySize"] = gpuMemorySize1
			}
			gpuType1, _ := jsonpath.Get("$[0].gpu_type", v)
			if gpuType1 != nil && (d.HasChange("gpu_config.0.gpu_type") || gpuType1 != "") {
				dataList["gpuType"] = gpuType1
			}

			request["gpuConfig"] = dataList
		}
	}

	if !d.IsNewResource() && d.HasChange("log_config") {
		update = true
		dataList1 := make(map[string]interface{})

		if v := d.Get("log_config"); v != nil {
			logBeginRule1, _ := jsonpath.Get("$[0].log_begin_rule", v)
			if logBeginRule1 != nil && (d.HasChange("log_config.0.log_begin_rule") || logBeginRule1 != "") {
				dataList1["logBeginRule"] = logBeginRule1
			}
			project1, _ := jsonpath.Get("$[0].project", v)
			if project1 != nil && (d.HasChange("log_config.0.project") || project1 != "") {
				dataList1["project"] = project1
			}
			enableInstanceMetrics1, _ := jsonpath.Get("$[0].enable_instance_metrics", v)
			if enableInstanceMetrics1 != nil && (d.HasChange("log_config.0.enable_instance_metrics") || enableInstanceMetrics1 != "") {
				dataList1["enableInstanceMetrics"] = enableInstanceMetrics1
			}
			enableRequestMetrics1, _ := jsonpath.Get("$[0].enable_request_metrics", v)
			if enableRequestMetrics1 != nil && (d.HasChange("log_config.0.enable_request_metrics") || enableRequestMetrics1 != "") {
				dataList1["enableRequestMetrics"] = enableRequestMetrics1
			}
			logstore1, _ := jsonpath.Get("$[0].logstore", v)
			if logstore1 != nil && (d.HasChange("log_config.0.logstore") || logstore1 != "") {
				dataList1["logstore"] = logstore1
			}

			request["logConfig"] = dataList1
		}
	}

	if d.HasChange("nas_config") {
		update = true
		dataList2 := make(map[string]interface{})

		if v := d.Get("nas_config"); v != nil {
			if v, ok := d.GetOk("nas_config"); ok {
				localData, err := jsonpath.Get("$[0].mount_points", v)
				if err != nil {
					localData = make([]interface{}, 0)
				}
				localMaps := make([]interface{}, 0)
				for _, dataLoop := range localData.([]interface{}) {
					dataLoopTmp := make(map[string]interface{})
					if dataLoop != nil {
						dataLoopTmp = dataLoop.(map[string]interface{})
					}
					dataLoopMap := make(map[string]interface{})
					dataLoopMap["enableTLS"] = dataLoopTmp["enable_tls"]
					dataLoopMap["serverAddr"] = dataLoopTmp["server_addr"]
					dataLoopMap["mountDir"] = dataLoopTmp["mount_dir"]
					localMaps = append(localMaps, dataLoopMap)
				}
				dataList2["mountPoints"] = localMaps
			}

			userId1, _ := jsonpath.Get("$[0].user_id", v)
			if userId1 != nil && (d.HasChange("nas_config.0.user_id") || userId1 != "") {
				dataList2["userId"] = userId1
			}
			groupId1, _ := jsonpath.Get("$[0].group_id", v)
			if groupId1 != nil && (d.HasChange("nas_config.0.group_id") || groupId1 != "") {
				dataList2["groupId"] = groupId1
			}

			request["nasConfig"] = dataList2
		}
	}

	if !d.IsNewResource() && d.HasChange("instance_concurrency") {
		update = true
		request["instanceConcurrency"] = d.Get("instance_concurrency")
	}

	if !d.IsNewResource() && d.HasChange("instance_isolation_mode") {
		update = true
		request["instanceIsolationMode"] = d.Get("instance_isolation_mode")
	}

	if !d.IsNewResource() && d.HasChange("custom_runtime_config") {
		update = true
		dataList3 := make(map[string]interface{})

		if v := d.Get("custom_runtime_config"); v != nil {
			healthCheckConfig := make(map[string]interface{})
			timeoutSeconds1, _ := jsonpath.Get("$[0].health_check_config[0].timeout_seconds", v)
			if timeoutSeconds1 != nil && (d.HasChange("custom_runtime_config.0.health_check_config.0.timeout_seconds") || timeoutSeconds1 != "") && timeoutSeconds1.(int) > 0 {
				healthCheckConfig["timeoutSeconds"] = timeoutSeconds1
			}
			httpGetUrl1, _ := jsonpath.Get("$[0].health_check_config[0].http_get_url", v)
			if httpGetUrl1 != nil && (d.HasChange("custom_runtime_config.0.health_check_config.0.http_get_url") || httpGetUrl1 != "") {
				healthCheckConfig["httpGetUrl"] = httpGetUrl1
			}
			successThreshold1, _ := jsonpath.Get("$[0].health_check_config[0].success_threshold", v)
			if successThreshold1 != nil && (d.HasChange("custom_runtime_config.0.health_check_config.0.success_threshold") || successThreshold1 != "") && successThreshold1.(int) > 0 {
				healthCheckConfig["successThreshold"] = successThreshold1
			}
			initialDelaySeconds1, _ := jsonpath.Get("$[0].health_check_config[0].initial_delay_seconds", v)
			if initialDelaySeconds1 != nil && (d.HasChange("custom_runtime_config.0.health_check_config.0.initial_delay_seconds") || initialDelaySeconds1 != "") {
				healthCheckConfig["initialDelaySeconds"] = initialDelaySeconds1
			}
			periodSeconds1, _ := jsonpath.Get("$[0].health_check_config[0].period_seconds", v)
			if periodSeconds1 != nil && (d.HasChange("custom_runtime_config.0.health_check_config.0.period_seconds") || periodSeconds1 != "") && periodSeconds1.(int) > 0 {
				healthCheckConfig["periodSeconds"] = periodSeconds1
			}
			failureThreshold1, _ := jsonpath.Get("$[0].health_check_config[0].failure_threshold", v)
			if failureThreshold1 != nil && (d.HasChange("custom_runtime_config.0.health_check_config.0.failure_threshold") || failureThreshold1 != "") && failureThreshold1.(int) > 0 {
				healthCheckConfig["failureThreshold"] = failureThreshold1
			}

			dataList3["healthCheckConfig"] = healthCheckConfig
			args1, _ := jsonpath.Get("$[0].args", d.Get("custom_runtime_config"))
			if args1 != nil && (d.HasChange("custom_runtime_config.0.args") || args1 != "") {
				dataList3["args"] = args1
			}
			command1, _ := jsonpath.Get("$[0].command", d.Get("custom_runtime_config"))
			if command1 != nil && (d.HasChange("custom_runtime_config.0.command") || command1 != "") {
				dataList3["command"] = command1
			}
			port1, _ := jsonpath.Get("$[0].port", v)
			if port1 != nil && (d.HasChange("custom_runtime_config.0.port") || port1 != "") && port1.(int) > 0 {
				dataList3["port"] = port1
			}

			request["customRuntimeConfig"] = dataList3
		}
	}

	if !d.IsNewResource() && d.HasChange("custom_container_config") {
		update = true
		dataList4 := make(map[string]interface{})

		if v := d.Get("custom_container_config"); v != nil {
			accelerationType1, _ := jsonpath.Get("$[0].acceleration_type", v)
			if accelerationType1 != nil && (d.HasChange("custom_container_config.0.acceleration_type") || accelerationType1 != "") {
				dataList4["accelerationType"] = accelerationType1
			}
			healthCheckConfig1 := make(map[string]interface{})
			failureThreshold3, _ := jsonpath.Get("$[0].health_check_config[0].failure_threshold", v)
			if failureThreshold3 != nil && (d.HasChange("custom_container_config.0.health_check_config.0.failure_threshold") || failureThreshold3 != "") && failureThreshold3.(int) > 0 {
				healthCheckConfig1["failureThreshold"] = failureThreshold3
			}
			timeoutSeconds3, _ := jsonpath.Get("$[0].health_check_config[0].timeout_seconds", v)
			if timeoutSeconds3 != nil && (d.HasChange("custom_container_config.0.health_check_config.0.timeout_seconds") || timeoutSeconds3 != "") && timeoutSeconds3.(int) > 0 {
				healthCheckConfig1["timeoutSeconds"] = timeoutSeconds3
			}
			initialDelaySeconds3, _ := jsonpath.Get("$[0].health_check_config[0].initial_delay_seconds", v)
			if initialDelaySeconds3 != nil && (d.HasChange("custom_container_config.0.health_check_config.0.initial_delay_seconds") || initialDelaySeconds3 != "") {
				healthCheckConfig1["initialDelaySeconds"] = initialDelaySeconds3
			}
			periodSeconds3, _ := jsonpath.Get("$[0].health_check_config[0].period_seconds", v)
			if periodSeconds3 != nil && (d.HasChange("custom_container_config.0.health_check_config.0.period_seconds") || periodSeconds3 != "") && periodSeconds3.(int) > 0 {
				healthCheckConfig1["periodSeconds"] = periodSeconds3
			}
			httpGetUrl3, _ := jsonpath.Get("$[0].health_check_config[0].http_get_url", v)
			if httpGetUrl3 != nil && (d.HasChange("custom_container_config.0.health_check_config.0.http_get_url") || httpGetUrl3 != "") {
				healthCheckConfig1["httpGetUrl"] = httpGetUrl3
			}
			successThreshold3, _ := jsonpath.Get("$[0].health_check_config[0].success_threshold", v)
			if successThreshold3 != nil && (d.HasChange("custom_container_config.0.health_check_config.0.success_threshold") || successThreshold3 != "") && successThreshold3.(int) > 0 {
				healthCheckConfig1["successThreshold"] = successThreshold3
			}

			dataList4["healthCheckConfig"] = healthCheckConfig1
			entrypoint1, _ := jsonpath.Get("$[0].entrypoint", d.Get("custom_container_config"))
			if entrypoint1 != nil && (d.HasChange("custom_container_config.0.entrypoint") || entrypoint1 != "") {
				dataList4["entrypoint"] = entrypoint1
			}
			command3, _ := jsonpath.Get("$[0].command", d.Get("custom_container_config"))
			if command3 != nil && (d.HasChange("custom_container_config.0.command") || command3 != "") {
				dataList4["command"] = command3
			}
			image1, _ := jsonpath.Get("$[0].image", v)
			if image1 != nil && (d.HasChange("custom_container_config.0.image") || image1 != "") {
				dataList4["image"] = image1
			}
			port3, _ := jsonpath.Get("$[0].port", v)
			if port3 != nil && (d.HasChange("custom_container_config.0.port") || port3 != "") {
				dataList4["port"] = port3
			}
			acrInstanceId1, _ := jsonpath.Get("$[0].acr_instance_id", v)
			if acrInstanceId1 != nil && (d.HasChange("custom_container_config.0.acr_instance_id") || acrInstanceId1 != "") {
				dataList4["acrInstanceId"] = acrInstanceId1
			}

			request["customContainerConfig"] = dataList4
		}
	}

	if d.HasChange("custom_dns") {
		update = true
		dataList5 := make(map[string]interface{})

		if v := d.Get("custom_dns"); v != nil {
			searches1, _ := jsonpath.Get("$[0].searches", d.Get("custom_dns"))
			if searches1 != nil && (d.HasChange("custom_dns.0.searches") || searches1 != "") {
				dataList5["searches"] = searches1
			}
			if v, ok := d.GetOk("custom_dns"); ok {
				localData1, err := jsonpath.Get("$[0].dns_options", v)
				if err != nil {
					localData1 = make([]interface{}, 0)
				}
				localMaps1 := make([]interface{}, 0)
				for _, dataLoop1 := range localData1.([]interface{}) {
					dataLoop1Tmp := make(map[string]interface{})
					if dataLoop1 != nil {
						dataLoop1Tmp = dataLoop1.(map[string]interface{})
					}
					dataLoop1Map := make(map[string]interface{})
					dataLoop1Map["value"] = dataLoop1Tmp["value"]
					dataLoop1Map["name"] = dataLoop1Tmp["name"]
					localMaps1 = append(localMaps1, dataLoop1Map)
				}
				dataList5["dnsOptions"] = localMaps1
			}

			nameServers1, _ := jsonpath.Get("$[0].name_servers", d.Get("custom_dns"))
			if nameServers1 != nil && (d.HasChange("custom_dns.0.name_servers") || nameServers1 != "") {
				dataList5["nameServers"] = nameServers1
			}

			request["customDNS"] = dataList5
		}
	}

	if !d.IsNewResource() && d.HasChange("instance_lifecycle_config") {
		update = true
		dataList6 := make(map[string]interface{})

		if v := d.Get("instance_lifecycle_config"); v != nil {
			preStop := make(map[string]interface{})
			timeout1, _ := jsonpath.Get("$[0].pre_stop[0].timeout", v)
			if timeout1 != nil && (d.HasChange("instance_lifecycle_config.0.pre_stop.0.timeout") || timeout1 != "") && timeout1.(int) > 0 {
				preStop["timeout"] = timeout1
			}
			handler1, _ := jsonpath.Get("$[0].pre_stop[0].handler", v)
			if handler1 != nil && (d.HasChange("instance_lifecycle_config.0.pre_stop.0.handler") || handler1 != "") {
				preStop["handler"] = handler1
			}

			dataList6["preStop"] = preStop
			initializer := make(map[string]interface{})
			command5, _ := jsonpath.Get("$[0].initializer[0].command", d.Get("instance_lifecycle_config"))
			if command5 != nil && (d.HasChange("instance_lifecycle_config.0.initializer.0.command") || command5 != "") {
				initializer["command"] = command5
			}
			handler3, _ := jsonpath.Get("$[0].initializer[0].handler", v)
			if handler3 != nil && (d.HasChange("instance_lifecycle_config.0.initializer.0.handler") || handler3 != "") {
				initializer["handler"] = handler3
			}
			timeout3, _ := jsonpath.Get("$[0].initializer[0].timeout", v)
			if timeout3 != nil && (d.HasChange("instance_lifecycle_config.0.initializer.0.timeout") || timeout3 != "") && timeout3.(int) > 0 {
				initializer["timeout"] = timeout3
			}

			dataList6["initializer"] = initializer

			request["instanceLifecycleConfig"] = dataList6
		}
	}

	if !d.IsNewResource() && d.HasChange("internet_access") {
		update = true
		request["internetAccess"] = d.Get("internet_access")
	}

	if !d.IsNewResource() && d.HasChange("oss_mount_config") {
		update = true
		dataList7 := make(map[string]interface{})

		if v := d.Get("oss_mount_config"); v != nil {
			if v, ok := d.GetOk("oss_mount_config"); ok {
				localData2, err := jsonpath.Get("$[0].mount_points", v)
				if err != nil {
					localData2 = make([]interface{}, 0)
				}
				localMaps2 := make([]interface{}, 0)
				for _, dataLoop2 := range localData2.([]interface{}) {
					dataLoop2Tmp := make(map[string]interface{})
					if dataLoop2 != nil {
						dataLoop2Tmp = dataLoop2.(map[string]interface{})
					}
					dataLoop2Map := make(map[string]interface{})
					dataLoop2Map["bucketName"] = dataLoop2Tmp["bucket_name"]
					dataLoop2Map["readOnly"] = dataLoop2Tmp["read_only"]
					dataLoop2Map["bucketPath"] = dataLoop2Tmp["bucket_path"]
					dataLoop2Map["mountDir"] = dataLoop2Tmp["mount_dir"]
					dataLoop2Map["endpoint"] = dataLoop2Tmp["endpoint"]
					localMaps2 = append(localMaps2, dataLoop2Map)
				}
				dataList7["mountPoints"] = localMaps2
			}

			request["ossMountConfig"] = dataList7
		}
	}

	if !d.IsNewResource() && d.HasChange("runtime") {
		update = true
	}
	request["runtime"] = d.Get("runtime")
	if !d.IsNewResource() && d.HasChange("environment_variables") {
		update = true
		request["environmentVariables"] = d.Get("environment_variables")
	}

	if !d.IsNewResource() && d.HasChange("code") {
		update = true
		dataList8 := make(map[string]interface{})

		if v := d.Get("code"); v != nil {
			ossBucketName1, _ := jsonpath.Get("$[0].oss_bucket_name", v)
			if ossBucketName1 != nil && (d.HasChange("code.0.oss_bucket_name") || ossBucketName1 != "") {
				dataList8["ossBucketName"] = ossBucketName1
			}
			zipFile1, _ := jsonpath.Get("$[0].zip_file", v)
			if zipFile1 != nil && (d.HasChange("code.0.zip_file") || zipFile1 != "") {
				dataList8["zipFile"] = zipFile1
			}
			ossObjectName1, _ := jsonpath.Get("$[0].oss_object_name", v)
			if ossObjectName1 != nil && (d.HasChange("code.0.oss_object_name") || ossObjectName1 != "") {
				dataList8["ossObjectName"] = ossObjectName1
			}
			checksum1, _ := jsonpath.Get("$[0].checksum", v)
			if checksum1 != nil && (d.HasChange("code.0.checksum") || checksum1 != "") {
				dataList8["checksum"] = checksum1
			}

			request["code"] = dataList8
		}
	}

	if !d.IsNewResource() && d.HasChange("role") {
		update = true
		request["role"] = d.Get("role")
	}

	if !d.IsNewResource() && d.HasChange("layers") {
		update = true
		if v, ok := d.GetOk("layers"); ok || d.HasChange("layers") {
			layersMapsArray := v.([]interface{})
			request["layers"] = layersMapsArray
		}
	}

	if !d.IsNewResource() && d.HasChange("timeout") {
		update = true
		request["timeout"] = d.Get("timeout")
	}

	if !d.IsNewResource() && d.HasChange("session_affinity") {
		update = true
		request["sessionAffinity"] = d.Get("session_affinity")
	}

	if !d.IsNewResource() && d.HasChange("session_affinity_config") {
		update = true
		request["sessionAffinityConfig"] = d.Get("session_affinity_config")
	}

	if !d.IsNewResource() && d.HasChange("cpu") {
		update = true
		request["cpu"] = d.Get("cpu")
	}

	if !d.IsNewResource() && d.HasChange("disk_size") {
		update = true
		request["diskSize"] = d.Get("disk_size")
	}

	if !d.IsNewResource() && d.HasChange("vpc_config") {
		update = true
		dataList9 := make(map[string]interface{})

		if v := d.Get("vpc_config"); v != nil {
			vpcId1, _ := jsonpath.Get("$[0].vpc_id", v)
			if vpcId1 != nil && (d.HasChange("vpc_config.0.vpc_id") || vpcId1 != "") {
				dataList9["vpcId"] = vpcId1
			}
			securityGroupId1, _ := jsonpath.Get("$[0].security_group_id", v)
			if securityGroupId1 != nil && (d.HasChange("vpc_config.0.security_group_id") || securityGroupId1 != "") {
				dataList9["securityGroupId"] = securityGroupId1
			}
			vSwitchIds1, _ := jsonpath.Get("$[0].vswitch_ids", d.Get("vpc_config"))
			if vSwitchIds1 != nil && (d.HasChange("vpc_config.0.vswitch_ids") || vSwitchIds1 != "") {
				dataList9["vSwitchIds"] = vSwitchIds1
			}

			request["vpcConfig"] = dataList9
		}
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["description"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("handler") {
		update = true
	}
	request["handler"] = d.Get("handler")
	if !d.IsNewResource() && d.HasChange("memory_size") {
		update = true
		request["memorySize"] = d.Get("memory_size")
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("FC", "2023-03-30", action, query, nil, body, true)
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

	if !d.IsNewResource() && d.HasChange("tags") {
		fcv3ServiceV2 := Fcv3ServiceV2{client}
		if err := fcv3ServiceV2.SetResourceTags(d, "function"); err != nil {
			return WrapError(err)
		}
	}
	update = false
	action = fmt.Sprintf("/2023-03-30/resource-groups")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["resourceId"] = d.Id()

	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["newResourceGroupId"] = d.Get("resource_group_id")
	}

	request["resourceType"] = "function"
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("FC", "2023-03-30", action, query, nil, body, true)
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

	d.Partial(false)
	return resourceAliCloudFcv3FunctionRead(d, meta)
}

func resourceAliCloudFcv3FunctionDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	functionName := d.Id()
	action := fmt.Sprintf("/2023-03-30/functions/%s", functionName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	request["functionName"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("FC", "2023-03-30", action, query, nil, nil, true)

		if err != nil {
			if IsExpectedErrors(err, []string{"429"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"FunctionNotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
