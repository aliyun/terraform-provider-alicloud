// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
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
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"checksum": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"oss_object_name": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"oss_bucket_name": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
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
				Type:      schema.TypeList,
				Optional:  true,
				Sensitive: true,
				Elem:      &schema.Schema{Type: schema.TypeString},
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
			"role": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"runtime": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"python3.10", "python3.9", "python3", "nodejs20", "nodejs18", "nodejs16", "nodejs14", "java11", "java8", "php7.2", "dotnetcore3.1", "go1", "custom.debian10", "custom", "custom-container"}, false),
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
	conn, err := client.NewFcv2Client()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	if v, ok := d.GetOk("function_name"); ok {
		request["functionName"] = v
	}

	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("gpu_config"); !IsNil(v) {
		gpuMemorySize1, _ := jsonpath.Get("$[0].gpu_memory_size", v)
		if gpuMemorySize1 != nil && gpuMemorySize1 != "" {
			objectDataLocalMap["gpuMemorySize"] = gpuMemorySize1
		}
		gpuType1, _ := jsonpath.Get("$[0].gpu_type", v)
		if gpuType1 != nil && gpuType1 != "" {
			objectDataLocalMap["gpuType"] = gpuType1
		}

		request["gpuConfig"] = objectDataLocalMap
	}

	objectDataLocalMap1 := make(map[string]interface{})

	if v := d.Get("log_config"); !IsNil(v) {
		logBeginRule1, _ := jsonpath.Get("$[0].log_begin_rule", v)
		if logBeginRule1 != nil && logBeginRule1 != "" {
			objectDataLocalMap1["logBeginRule"] = logBeginRule1
		}
		enableInstanceMetrics1, _ := jsonpath.Get("$[0].enable_instance_metrics", v)
		if enableInstanceMetrics1 != nil && enableInstanceMetrics1 != "" {
			objectDataLocalMap1["enableInstanceMetrics"] = enableInstanceMetrics1
		}
		project1, _ := jsonpath.Get("$[0].project", v)
		if project1 != nil && project1 != "" {
			objectDataLocalMap1["project"] = project1
		}
		enableRequestMetrics1, _ := jsonpath.Get("$[0].enable_request_metrics", v)
		if enableRequestMetrics1 != nil && enableRequestMetrics1 != "" {
			objectDataLocalMap1["enableRequestMetrics"] = enableRequestMetrics1
		}
		logstore1, _ := jsonpath.Get("$[0].logstore", v)
		if logstore1 != nil && logstore1 != "" {
			objectDataLocalMap1["logstore"] = logstore1
		}

		request["logConfig"] = objectDataLocalMap1
	}

	objectDataLocalMap2 := make(map[string]interface{})

	if v := d.Get("custom_container_config"); !IsNil(v) {
		healthCheckConfig := make(map[string]interface{})
		httpGetUrl1, _ := jsonpath.Get("$[0].health_check_config[0].http_get_url", v)
		if httpGetUrl1 != nil && httpGetUrl1 != "" {
			healthCheckConfig["httpGetUrl"] = httpGetUrl1
		}
		failureThreshold1, _ := jsonpath.Get("$[0].health_check_config[0].failure_threshold", v)
		if failureThreshold1 != nil && failureThreshold1 != "" && failureThreshold1.(int) > 0 {
			healthCheckConfig["failureThreshold"] = failureThreshold1
		}
		successThreshold1, _ := jsonpath.Get("$[0].health_check_config[0].success_threshold", v)
		if successThreshold1 != nil && successThreshold1 != "" && successThreshold1.(int) > 0 {
			healthCheckConfig["successThreshold"] = successThreshold1
		}
		timeoutSeconds1, _ := jsonpath.Get("$[0].health_check_config[0].timeout_seconds", v)
		if timeoutSeconds1 != nil && timeoutSeconds1 != "" && timeoutSeconds1.(int) > 0 {
			healthCheckConfig["timeoutSeconds"] = timeoutSeconds1
		}
		initialDelaySeconds1, _ := jsonpath.Get("$[0].health_check_config[0].initial_delay_seconds", v)
		if initialDelaySeconds1 != nil && initialDelaySeconds1 != "" {
			healthCheckConfig["initialDelaySeconds"] = initialDelaySeconds1
		}
		periodSeconds1, _ := jsonpath.Get("$[0].health_check_config[0].period_seconds", v)
		if periodSeconds1 != nil && periodSeconds1 != "" && periodSeconds1.(int) > 0 {
			healthCheckConfig["periodSeconds"] = periodSeconds1
		}

		objectDataLocalMap2["healthCheckConfig"] = healthCheckConfig
		accelerationType1, _ := jsonpath.Get("$[0].acceleration_type", v)
		if accelerationType1 != nil && accelerationType1 != "" {
			objectDataLocalMap2["accelerationType"] = accelerationType1
		}
		command1, _ := jsonpath.Get("$[0].command", v)
		if command1 != nil && command1 != "" {
			objectDataLocalMap2["command"] = command1
		}
		image1, _ := jsonpath.Get("$[0].image", v)
		if image1 != nil && image1 != "" {
			objectDataLocalMap2["image"] = image1
		}
		port1, _ := jsonpath.Get("$[0].port", v)
		if port1 != nil && port1 != "" {
			objectDataLocalMap2["port"] = port1
		}
		acrInstanceId1, _ := jsonpath.Get("$[0].acr_instance_id", v)
		if acrInstanceId1 != nil && acrInstanceId1 != "" {
			objectDataLocalMap2["acrInstanceId"] = acrInstanceId1
		}
		entrypoint1, _ := jsonpath.Get("$[0].entrypoint", v)
		if entrypoint1 != nil && entrypoint1 != "" {
			objectDataLocalMap2["entrypoint"] = entrypoint1
		}

		request["customContainerConfig"] = objectDataLocalMap2
	}

	objectDataLocalMap3 := make(map[string]interface{})

	if v := d.Get("custom_runtime_config"); !IsNil(v) {
		healthCheckConfig1 := make(map[string]interface{})
		successThreshold3, _ := jsonpath.Get("$[0].health_check_config[0].success_threshold", v)
		if successThreshold3 != nil && successThreshold3 != "" && successThreshold3.(int) > 0 {
			healthCheckConfig1["successThreshold"] = successThreshold3
		}
		timeoutSeconds3, _ := jsonpath.Get("$[0].health_check_config[0].timeout_seconds", v)
		if timeoutSeconds3 != nil && timeoutSeconds3 != "" && timeoutSeconds3.(int) > 0 {
			healthCheckConfig1["timeoutSeconds"] = timeoutSeconds3
		}
		initialDelaySeconds3, _ := jsonpath.Get("$[0].health_check_config[0].initial_delay_seconds", v)
		if initialDelaySeconds3 != nil && initialDelaySeconds3 != "" {
			healthCheckConfig1["initialDelaySeconds"] = initialDelaySeconds3
		}
		httpGetUrl3, _ := jsonpath.Get("$[0].health_check_config[0].http_get_url", v)
		if httpGetUrl3 != nil && httpGetUrl3 != "" {
			healthCheckConfig1["httpGetUrl"] = httpGetUrl3
		}
		periodSeconds3, _ := jsonpath.Get("$[0].health_check_config[0].period_seconds", v)
		if periodSeconds3 != nil && periodSeconds3 != "" && periodSeconds3.(int) > 0 {
			healthCheckConfig1["periodSeconds"] = periodSeconds3
		}
		failureThreshold3, _ := jsonpath.Get("$[0].health_check_config[0].failure_threshold", v)
		if failureThreshold3 != nil && failureThreshold3 != "" && failureThreshold3.(int) > 0 {
			healthCheckConfig1["failureThreshold"] = failureThreshold3
		}

		objectDataLocalMap3["healthCheckConfig"] = healthCheckConfig1
		command3, _ := jsonpath.Get("$[0].command", v)
		if command3 != nil && command3 != "" {
			objectDataLocalMap3["command"] = command3
		}
		port3, _ := jsonpath.Get("$[0].port", v)
		if port3 != nil && port3 != "" && port3.(int) > 0 {
			objectDataLocalMap3["port"] = port3
		}
		args1, _ := jsonpath.Get("$[0].args", v)
		if args1 != nil && args1 != "" {
			objectDataLocalMap3["args"] = args1
		}

		request["customRuntimeConfig"] = objectDataLocalMap3
	}

	if v, ok := d.GetOk("layers"); ok {
		layersMapsArray := v.([]interface{})
		request["layers"] = layersMapsArray
	}

	if v, ok := d.GetOkExists("timeout"); ok && v.(int) > 0 {
		request["timeout"] = v
	}
	if v, ok := d.GetOkExists("instance_concurrency"); ok && v.(int) > 0 {
		request["instanceConcurrency"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request["Tags"] = tagsMap
	}

	objectDataLocalMap4 := make(map[string]interface{})

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

		objectDataLocalMap4["preStop"] = preStop
		initializer := make(map[string]interface{})
		handler3, _ := jsonpath.Get("$[0].initializer[0].handler", v)
		if handler3 != nil && handler3 != "" {
			initializer["handler"] = handler3
		}
		timeout3, _ := jsonpath.Get("$[0].initializer[0].timeout", v)
		if timeout3 != nil && timeout3 != "" && timeout3.(int) > 0 {
			initializer["timeout"] = timeout3
		}

		objectDataLocalMap4["initializer"] = initializer

		request["instanceLifecycleConfig"] = objectDataLocalMap4
	}

	if v, ok := d.GetOkExists("internet_access"); ok {
		request["internetAccess"] = v
	}
	objectDataLocalMap5 := make(map[string]interface{})

	if v := d.Get("oss_mount_config"); !IsNil(v) {
		if v, ok := d.GetOk("oss_mount_config"); ok {
			localData2, err := jsonpath.Get("$[0].mount_points", v)
			if err != nil {
				localData2 = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop2 := range localData2.([]interface{}) {
				dataLoop2Tmp := make(map[string]interface{})
				if dataLoop2 != nil {
					dataLoop2Tmp = dataLoop2.(map[string]interface{})
				}
				dataLoop2Map := make(map[string]interface{})
				dataLoop2Map["bucketName"] = dataLoop2Tmp["bucket_name"]
				dataLoop2Map["bucketPath"] = dataLoop2Tmp["bucket_path"]
				dataLoop2Map["mountDir"] = dataLoop2Tmp["mount_dir"]
				dataLoop2Map["readOnly"] = dataLoop2Tmp["read_only"]
				dataLoop2Map["endpoint"] = dataLoop2Tmp["endpoint"]
				localMaps = append(localMaps, dataLoop2Map)
			}
			objectDataLocalMap5["mountPoints"] = localMaps
		}

		request["ossMountConfig"] = objectDataLocalMap5
	}

	request["runtime"] = d.Get("runtime")
	if v, ok := d.GetOk("cpu"); ok && v.(float64) > 0 {
		request["cpu"] = v
	}
	if v, ok := d.GetOk("environment_variables"); ok {
		request["environmentVariables"] = v
	}
	objectDataLocalMap6 := make(map[string]interface{})

	if v := d.Get("code"); !IsNil(v) {
		ossBucketName1, _ := jsonpath.Get("$[0].oss_bucket_name", v)
		if ossBucketName1 != nil && ossBucketName1 != "" {
			objectDataLocalMap6["ossBucketName"] = ossBucketName1
		}
		zipFile1, _ := jsonpath.Get("$[0].zip_file", v)
		if zipFile1 != nil && zipFile1 != "" {
			objectDataLocalMap6["zipFile"] = zipFile1
		}
		checksum1, _ := jsonpath.Get("$[0].checksum", v)
		if checksum1 != nil && checksum1 != "" {
			objectDataLocalMap6["checksum"] = checksum1
		}
		ossObjectName1, _ := jsonpath.Get("$[0].oss_object_name", v)
		if ossObjectName1 != nil && ossObjectName1 != "" {
			objectDataLocalMap6["ossObjectName"] = ossObjectName1
		}

		request["code"] = objectDataLocalMap6
	}

	if v, ok := d.GetOk("role"); ok {
		request["role"] = v
	}
	if v, ok := d.GetOkExists("disk_size"); ok && v.(int) > 0 {
		request["diskSize"] = v
	}
	objectDataLocalMap7 := make(map[string]interface{})

	if v := d.Get("vpc_config"); !IsNil(v) {
		vpcId1, _ := jsonpath.Get("$[0].vpc_id", v)
		if vpcId1 != nil && vpcId1 != "" {
			objectDataLocalMap7["vpcId"] = vpcId1
		}
		securityGroupId1, _ := jsonpath.Get("$[0].security_group_id", v)
		if securityGroupId1 != nil && securityGroupId1 != "" {
			objectDataLocalMap7["securityGroupId"] = securityGroupId1
		}
		vSwitchIds1, _ := jsonpath.Get("$[0].vswitch_ids", v)
		if vSwitchIds1 != nil && vSwitchIds1 != "" {
			objectDataLocalMap7["vSwitchIds"] = vSwitchIds1
		}

		request["vpcConfig"] = objectDataLocalMap7
	}

	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	request["handler"] = d.Get("handler")
	if v, ok := d.GetOkExists("memory_size"); ok && v.(int) > 0 {
		request["memorySize"] = v
	}
	body = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2023-03-30"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
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

	id, _ := jsonpath.Get("$.body.functionName", response)
	d.SetId(fmt.Sprint(id))

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

	if objectRaw["codeSize"] != nil {
		d.Set("code_size", objectRaw["codeSize"])
	}
	if objectRaw["cpu"] != nil {
		d.Set("cpu", objectRaw["cpu"])
	}
	if objectRaw["createdTime"] != nil {
		d.Set("create_time", objectRaw["createdTime"])
	}
	if objectRaw["description"] != nil {
		d.Set("description", objectRaw["description"])
	}
	if objectRaw["diskSize"] != nil {
		d.Set("disk_size", objectRaw["diskSize"])
	}
	if objectRaw["environmentVariables"] != nil {
		d.Set("environment_variables", objectRaw["environmentVariables"])
	}
	if objectRaw["functionArn"] != nil {
		d.Set("function_arn", objectRaw["functionArn"])
	}
	if objectRaw["functionId"] != nil {
		d.Set("function_id", objectRaw["functionId"])
	}
	if objectRaw["handler"] != nil {
		d.Set("handler", objectRaw["handler"])
	}
	if objectRaw["instanceConcurrency"] != nil {
		d.Set("instance_concurrency", objectRaw["instanceConcurrency"])
	}
	if objectRaw["internetAccess"] != nil {
		d.Set("internet_access", objectRaw["internetAccess"])
	}
	if objectRaw["lastModifiedTime"] != nil {
		d.Set("last_modified_time", objectRaw["lastModifiedTime"])
	}
	if objectRaw["lastUpdateStatus"] != nil {
		d.Set("last_update_status", objectRaw["lastUpdateStatus"])
	}
	if objectRaw["lastUpdateStatusReason"] != nil {
		d.Set("last_update_status_reason", objectRaw["lastUpdateStatusReason"])
	}
	if objectRaw["lastUpdateStatusReasonCode"] != nil {
		d.Set("last_update_status_reason_code", objectRaw["lastUpdateStatusReasonCode"])
	}
	if objectRaw["memorySize"] != nil {
		d.Set("memory_size", objectRaw["memorySize"])
	}
	if objectRaw["role"] != nil {
		d.Set("role", objectRaw["role"])
	}
	if objectRaw["runtime"] != nil {
		d.Set("runtime", objectRaw["runtime"])
	}
	if objectRaw["state"] != nil {
		d.Set("state", objectRaw["state"])
	}
	if objectRaw["stateReason"] != nil {
		d.Set("state_reason", objectRaw["stateReason"])
	}
	if objectRaw["stateReasonCode"] != nil {
		d.Set("state_reason_code", objectRaw["stateReasonCode"])
	}
	if objectRaw["timeout"] != nil {
		d.Set("timeout", objectRaw["timeout"])
	}

	customContainerConfigMaps := make([]map[string]interface{}, 0)
	customContainerConfigMap := make(map[string]interface{})
	customContainerConfig1Raw := make(map[string]interface{})
	if objectRaw["customContainerConfig"] != nil {
		customContainerConfig1Raw = objectRaw["customContainerConfig"].(map[string]interface{})
	}
	if len(customContainerConfig1Raw) > 0 {
		customContainerConfigMap["acceleration_type"] = customContainerConfig1Raw["accelerationType"]
		customContainerConfigMap["acr_instance_id"] = customContainerConfig1Raw["acrInstanceId"]
		customContainerConfigMap["image"] = customContainerConfig1Raw["image"]
		customContainerConfigMap["port"] = customContainerConfig1Raw["port"]
		customContainerConfigMap["resolved_image_uri"] = customContainerConfig1Raw["resolvedImageUri"]

		accelerationInfoMaps := make([]map[string]interface{}, 0)
		accelerationInfoMap := make(map[string]interface{})
		accelerationInfo1Raw := make(map[string]interface{})
		if customContainerConfig1Raw["accelerationInfo"] != nil {
			accelerationInfo1Raw = customContainerConfig1Raw["accelerationInfo"].(map[string]interface{})
		}
		if len(accelerationInfo1Raw) > 0 {
			accelerationInfoMap["status"] = accelerationInfo1Raw["status"]

			accelerationInfoMaps = append(accelerationInfoMaps, accelerationInfoMap)
		}
		customContainerConfigMap["acceleration_info"] = accelerationInfoMaps
		command2Raw := make([]interface{}, 0)
		if customContainerConfig1Raw["command"] != nil {
			command2Raw = customContainerConfig1Raw["command"].([]interface{})
		}

		customContainerConfigMap["command"] = command2Raw
		entrypoint1Raw := make([]interface{}, 0)
		if customContainerConfig1Raw["entrypoint"] != nil {
			entrypoint1Raw = customContainerConfig1Raw["entrypoint"].([]interface{})
		}

		customContainerConfigMap["entrypoint"] = entrypoint1Raw
		healthCheckConfigMaps := make([]map[string]interface{}, 0)
		healthCheckConfigMap := make(map[string]interface{})
		healthCheckConfig2Raw := make(map[string]interface{})
		if customContainerConfig1Raw["healthCheckConfig"] != nil {
			healthCheckConfig2Raw = customContainerConfig1Raw["healthCheckConfig"].(map[string]interface{})
		}
		if len(healthCheckConfig2Raw) > 0 {
			healthCheckConfigMap["failure_threshold"] = healthCheckConfig2Raw["failureThreshold"]
			healthCheckConfigMap["http_get_url"] = healthCheckConfig2Raw["httpGetUrl"]
			healthCheckConfigMap["initial_delay_seconds"] = healthCheckConfig2Raw["initialDelaySeconds"]
			healthCheckConfigMap["period_seconds"] = healthCheckConfig2Raw["periodSeconds"]
			healthCheckConfigMap["success_threshold"] = healthCheckConfig2Raw["successThreshold"]
			healthCheckConfigMap["timeout_seconds"] = healthCheckConfig2Raw["timeoutSeconds"]

			healthCheckConfigMaps = append(healthCheckConfigMaps, healthCheckConfigMap)
		}
		customContainerConfigMap["health_check_config"] = healthCheckConfigMaps
		customContainerConfigMaps = append(customContainerConfigMaps, customContainerConfigMap)
	}
	if objectRaw["customContainerConfig"] != nil {
		if err := d.Set("custom_container_config", customContainerConfigMaps); err != nil {
			return err
		}
	}
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
	if objectRaw["customDNS"] != nil {
		if err := d.Set("custom_dns", customDnsMaps); err != nil {
			return err
		}
	}
	customRuntimeConfigMaps := make([]map[string]interface{}, 0)
	customRuntimeConfigMap := make(map[string]interface{})
	customRuntimeConfig1Raw := make(map[string]interface{})
	if objectRaw["customRuntimeConfig"] != nil {
		customRuntimeConfig1Raw = objectRaw["customRuntimeConfig"].(map[string]interface{})
	}
	if len(customRuntimeConfig1Raw) > 0 {
		customRuntimeConfigMap["port"] = customRuntimeConfig1Raw["port"]

		args1Raw := make([]interface{}, 0)
		if customRuntimeConfig1Raw["args"] != nil {
			args1Raw = customRuntimeConfig1Raw["args"].([]interface{})
		}

		customRuntimeConfigMap["args"] = args1Raw
		command3Raw := make([]interface{}, 0)
		if customRuntimeConfig1Raw["command"] != nil {
			command3Raw = customRuntimeConfig1Raw["command"].([]interface{})
		}

		customRuntimeConfigMap["command"] = command3Raw
		healthCheckConfigMaps := make([]map[string]interface{}, 0)
		healthCheckConfigMap := make(map[string]interface{})
		healthCheckConfig3Raw := make(map[string]interface{})
		if customRuntimeConfig1Raw["healthCheckConfig"] != nil {
			healthCheckConfig3Raw = customRuntimeConfig1Raw["healthCheckConfig"].(map[string]interface{})
		}
		if len(healthCheckConfig3Raw) > 0 {
			healthCheckConfigMap["failure_threshold"] = healthCheckConfig3Raw["failureThreshold"]
			healthCheckConfigMap["http_get_url"] = healthCheckConfig3Raw["httpGetUrl"]
			healthCheckConfigMap["initial_delay_seconds"] = healthCheckConfig3Raw["initialDelaySeconds"]
			healthCheckConfigMap["period_seconds"] = healthCheckConfig3Raw["periodSeconds"]
			healthCheckConfigMap["success_threshold"] = healthCheckConfig3Raw["successThreshold"]
			healthCheckConfigMap["timeout_seconds"] = healthCheckConfig3Raw["timeoutSeconds"]

			healthCheckConfigMaps = append(healthCheckConfigMaps, healthCheckConfigMap)
		}
		customRuntimeConfigMap["health_check_config"] = healthCheckConfigMaps
		customRuntimeConfigMaps = append(customRuntimeConfigMaps, customRuntimeConfigMap)
	}
	if objectRaw["customRuntimeConfig"] != nil {
		if err := d.Set("custom_runtime_config", customRuntimeConfigMaps); err != nil {
			return err
		}
	}
	gpuConfigMaps := make([]map[string]interface{}, 0)
	gpuConfigMap := make(map[string]interface{})
	gpuConfig1Raw := make(map[string]interface{})
	if objectRaw["gpuConfig"] != nil {
		gpuConfig1Raw = objectRaw["gpuConfig"].(map[string]interface{})
	}
	if len(gpuConfig1Raw) > 0 {
		gpuConfigMap["gpu_memory_size"] = gpuConfig1Raw["gpuMemorySize"]
		gpuConfigMap["gpu_type"] = gpuConfig1Raw["gpuType"]

		gpuConfigMaps = append(gpuConfigMaps, gpuConfigMap)
	}
	if objectRaw["gpuConfig"] != nil {
		if err := d.Set("gpu_config", gpuConfigMaps); err != nil {
			return err
		}
	}
	instanceLifecycleConfigMaps := make([]map[string]interface{}, 0)
	instanceLifecycleConfigMap := make(map[string]interface{})
	instanceLifecycleConfig1Raw := make(map[string]interface{})
	if objectRaw["instanceLifecycleConfig"] != nil {
		instanceLifecycleConfig1Raw = objectRaw["instanceLifecycleConfig"].(map[string]interface{})
	}
	if len(instanceLifecycleConfig1Raw) > 0 {

		initializerMaps := make([]map[string]interface{}, 0)
		initializerMap := make(map[string]interface{})
		initializer1Raw := make(map[string]interface{})
		if instanceLifecycleConfig1Raw["initializer"] != nil {
			initializer1Raw = instanceLifecycleConfig1Raw["initializer"].(map[string]interface{})
		}
		if len(initializer1Raw) > 0 {
			initializerMap["handler"] = initializer1Raw["handler"]
			initializerMap["timeout"] = initializer1Raw["timeout"]

			initializerMaps = append(initializerMaps, initializerMap)
		}
		instanceLifecycleConfigMap["initializer"] = initializerMaps
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
	if objectRaw["instanceLifecycleConfig"] != nil {
		if err := d.Set("instance_lifecycle_config", instanceLifecycleConfigMaps); err != nil {
			return err
		}
	}
	logConfigMaps := make([]map[string]interface{}, 0)
	logConfigMap := make(map[string]interface{})
	logConfig1Raw := make(map[string]interface{})
	if objectRaw["logConfig"] != nil {
		logConfig1Raw = objectRaw["logConfig"].(map[string]interface{})
	}
	if len(logConfig1Raw) > 0 {
		logConfigMap["enable_instance_metrics"] = logConfig1Raw["enableInstanceMetrics"]
		logConfigMap["enable_request_metrics"] = logConfig1Raw["enableRequestMetrics"]
		logConfigMap["log_begin_rule"] = logConfig1Raw["logBeginRule"]
		logConfigMap["logstore"] = logConfig1Raw["logstore"]
		logConfigMap["project"] = logConfig1Raw["project"]

		logConfigMaps = append(logConfigMaps, logConfigMap)
	}
	if objectRaw["logConfig"] != nil {
		if err := d.Set("log_config", logConfigMaps); err != nil {
			return err
		}
	}
	nasConfigMaps := make([]map[string]interface{}, 0)
	nasConfigMap := make(map[string]interface{})
	nasConfig1Raw := make(map[string]interface{})
	if objectRaw["nasConfig"] != nil {
		nasConfig1Raw = objectRaw["nasConfig"].(map[string]interface{})
	}
	if len(nasConfig1Raw) > 0 {
		nasConfigMap["group_id"] = nasConfig1Raw["groupId"]
		nasConfigMap["user_id"] = nasConfig1Raw["userId"]

		mountPoints2Raw := nasConfig1Raw["mountPoints"]
		mountPointsMaps := make([]map[string]interface{}, 0)
		if mountPoints2Raw != nil {
			for _, mountPointsChild2Raw := range mountPoints2Raw.([]interface{}) {
				mountPointsMap := make(map[string]interface{})
				mountPointsChild2Raw := mountPointsChild2Raw.(map[string]interface{})
				mountPointsMap["enable_tls"] = mountPointsChild2Raw["enableTLS"]
				mountPointsMap["mount_dir"] = mountPointsChild2Raw["mountDir"]
				mountPointsMap["server_addr"] = mountPointsChild2Raw["serverAddr"]

				mountPointsMaps = append(mountPointsMaps, mountPointsMap)
			}
		}
		nasConfigMap["mount_points"] = mountPointsMaps
		nasConfigMaps = append(nasConfigMaps, nasConfigMap)
	}
	if objectRaw["nasConfig"] != nil {
		if err := d.Set("nas_config", nasConfigMaps); err != nil {
			return err
		}
	}
	ossMountConfigMaps := make([]map[string]interface{}, 0)
	ossMountConfigMap := make(map[string]interface{})
	mountPoints3Raw, _ := jsonpath.Get("$.ossMountConfig.mountPoints", objectRaw)

	mountPointsMaps := make([]map[string]interface{}, 0)
	if mountPoints3Raw != nil {
		for _, mountPointsChild3Raw := range mountPoints3Raw.([]interface{}) {
			mountPointsMap := make(map[string]interface{})
			mountPointsChild3Raw := mountPointsChild3Raw.(map[string]interface{})
			mountPointsMap["bucket_name"] = mountPointsChild3Raw["bucketName"]
			mountPointsMap["bucket_path"] = mountPointsChild3Raw["bucketPath"]
			mountPointsMap["endpoint"] = mountPointsChild3Raw["endpoint"]
			mountPointsMap["mount_dir"] = mountPointsChild3Raw["mountDir"]
			mountPointsMap["read_only"] = mountPointsChild3Raw["readOnly"]

			mountPointsMaps = append(mountPointsMaps, mountPointsMap)
		}
	}
	ossMountConfigMap["mount_points"] = mountPointsMaps
	ossMountConfigMaps = append(ossMountConfigMaps, ossMountConfigMap)
	if mountPoints3Raw != nil {
		if err := d.Set("oss_mount_config", ossMountConfigMaps); err != nil {
			return err
		}
	}
	tagsMaps := objectRaw["tags"]
	d.Set("tags", tagsToMap(tagsMaps))
	tracingConfigMaps := make([]map[string]interface{}, 0)
	tracingConfigMap := make(map[string]interface{})
	tracingConfig1Raw := make(map[string]interface{})
	if objectRaw["tracingConfig"] != nil {
		tracingConfig1Raw = objectRaw["tracingConfig"].(map[string]interface{})
	}
	if len(tracingConfig1Raw) > 0 {
		tracingConfigMap["params"] = tracingConfig1Raw["params"]
		tracingConfigMap["type"] = tracingConfig1Raw["type"]

		tracingConfigMaps = append(tracingConfigMaps, tracingConfigMap)
	}
	if objectRaw["tracingConfig"] != nil {
		if err := d.Set("tracing_config", tracingConfigMaps); err != nil {
			return err
		}
	}
	vpcConfigMaps := make([]map[string]interface{}, 0)
	vpcConfigMap := make(map[string]interface{})
	vpcConfig1Raw := make(map[string]interface{})
	if objectRaw["vpcConfig"] != nil {
		vpcConfig1Raw = objectRaw["vpcConfig"].(map[string]interface{})
	}
	if len(vpcConfig1Raw) > 0 {
		vpcConfigMap["security_group_id"] = vpcConfig1Raw["securityGroupId"]
		vpcConfigMap["vpc_id"] = vpcConfig1Raw["vpcId"]

		vSwitchIds1Raw := make([]interface{}, 0)
		if vpcConfig1Raw["vSwitchIds"] != nil {
			vSwitchIds1Raw = vpcConfig1Raw["vSwitchIds"].([]interface{})
		}

		vpcConfigMap["vswitch_ids"] = vSwitchIds1Raw
		vpcConfigMaps = append(vpcConfigMaps, vpcConfigMap)
	}
	if objectRaw["vpcConfig"] != nil {
		if err := d.Set("vpc_config", vpcConfigMaps); err != nil {
			return err
		}
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

	functionName := d.Id()
	action := fmt.Sprintf("/2023-03-30/functions/%s", functionName)
	conn, err := client.NewFcv2Client()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["functionName"] = d.Id()

	if !d.IsNewResource() && d.HasChange("gpu_config") {
		update = true
		objectDataLocalMap := make(map[string]interface{})

		if v := d.Get("gpu_config"); v != nil {
			gpuMemorySize1, _ := jsonpath.Get("$[0].gpu_memory_size", v)
			if gpuMemorySize1 != nil && (d.HasChange("gpu_config.0.gpu_memory_size") || gpuMemorySize1 != "") {
				objectDataLocalMap["gpuMemorySize"] = gpuMemorySize1
			}
			gpuType1, _ := jsonpath.Get("$[0].gpu_type", v)
			if gpuType1 != nil && (d.HasChange("gpu_config.0.gpu_type") || gpuType1 != "") {
				objectDataLocalMap["gpuType"] = gpuType1
			}

			request["gpuConfig"] = objectDataLocalMap
		}
	}

	if !d.IsNewResource() && d.HasChange("log_config") {
		update = true
		objectDataLocalMap1 := make(map[string]interface{})

		if v := d.Get("log_config"); v != nil {
			logBeginRule1, _ := jsonpath.Get("$[0].log_begin_rule", v)
			if logBeginRule1 != nil && (d.HasChange("log_config.0.log_begin_rule") || logBeginRule1 != "") {
				objectDataLocalMap1["logBeginRule"] = logBeginRule1
			}
			project1, _ := jsonpath.Get("$[0].project", v)
			if project1 != nil && (d.HasChange("log_config.0.project") || project1 != "") {
				objectDataLocalMap1["project"] = project1
			}
			enableInstanceMetrics1, _ := jsonpath.Get("$[0].enable_instance_metrics", v)
			if enableInstanceMetrics1 != nil && (d.HasChange("log_config.0.enable_instance_metrics") || enableInstanceMetrics1 != "") {
				objectDataLocalMap1["enableInstanceMetrics"] = enableInstanceMetrics1
			}
			enableRequestMetrics1, _ := jsonpath.Get("$[0].enable_request_metrics", v)
			if enableRequestMetrics1 != nil && (d.HasChange("log_config.0.enable_request_metrics") || enableRequestMetrics1 != "") {
				objectDataLocalMap1["enableRequestMetrics"] = enableRequestMetrics1
			}
			logstore1, _ := jsonpath.Get("$[0].logstore", v)
			if logstore1 != nil && (d.HasChange("log_config.0.logstore") || logstore1 != "") {
				objectDataLocalMap1["logstore"] = logstore1
			}

			request["logConfig"] = objectDataLocalMap1
		}
	}

	if d.HasChange("nas_config") {
		update = true
		objectDataLocalMap2 := make(map[string]interface{})

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
				objectDataLocalMap2["mountPoints"] = localMaps
			}

			userId1, _ := jsonpath.Get("$[0].user_id", v)
			if userId1 != nil && (d.HasChange("nas_config.0.user_id") || userId1 != "") {
				objectDataLocalMap2["userId"] = userId1
			}
			groupId1, _ := jsonpath.Get("$[0].group_id", v)
			if groupId1 != nil && (d.HasChange("nas_config.0.group_id") || groupId1 != "") {
				objectDataLocalMap2["groupId"] = groupId1
			}

			request["nasConfig"] = objectDataLocalMap2
		}
	}

	if !d.IsNewResource() && d.HasChange("instance_concurrency") {
		update = true
		request["instanceConcurrency"] = d.Get("instance_concurrency")
	}

	if d.HasChange("custom_runtime_config") {
		update = true
		objectDataLocalMap3 := make(map[string]interface{})

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

			objectDataLocalMap3["healthCheckConfig"] = healthCheckConfig
			args1, _ := jsonpath.Get("$[0].args", d.Get("custom_runtime_config"))
			if args1 != nil && (d.HasChange("custom_runtime_config.0.args") || args1 != "") {
				objectDataLocalMap3["args"] = args1
			}
			command1, _ := jsonpath.Get("$[0].command", d.Get("custom_runtime_config"))
			if command1 != nil && (d.HasChange("custom_runtime_config.0.command") || command1 != "") {
				objectDataLocalMap3["command"] = command1
			}
			port1, _ := jsonpath.Get("$[0].port", v)
			if port1 != nil && (d.HasChange("custom_runtime_config.0.port") || port1 != "") && port1.(int) > 0 {
				objectDataLocalMap3["port"] = port1
			}

			request["customRuntimeConfig"] = objectDataLocalMap3
		}
	}

	if !d.IsNewResource() && d.HasChange("custom_container_config") {
		update = true
		objectDataLocalMap4 := make(map[string]interface{})

		if v := d.Get("custom_container_config"); v != nil {
			accelerationType1, _ := jsonpath.Get("$[0].acceleration_type", v)
			if accelerationType1 != nil && (d.HasChange("custom_container_config.0.acceleration_type") || accelerationType1 != "") {
				objectDataLocalMap4["accelerationType"] = accelerationType1
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

			objectDataLocalMap4["healthCheckConfig"] = healthCheckConfig1
			entrypoint1, _ := jsonpath.Get("$[0].entrypoint", d.Get("custom_container_config"))
			if entrypoint1 != nil && (d.HasChange("custom_container_config.0.entrypoint") || entrypoint1 != "") {
				objectDataLocalMap4["entrypoint"] = entrypoint1
			}
			command3, _ := jsonpath.Get("$[0].command", d.Get("custom_container_config"))
			if command3 != nil && (d.HasChange("custom_container_config.0.command") || command3 != "") {
				objectDataLocalMap4["command"] = command3
			}
			image1, _ := jsonpath.Get("$[0].image", v)
			if image1 != nil && (d.HasChange("custom_container_config.0.image") || image1 != "") {
				objectDataLocalMap4["image"] = image1
			}
			port3, _ := jsonpath.Get("$[0].port", v)
			if port3 != nil && (d.HasChange("custom_container_config.0.port") || port3 != "") {
				objectDataLocalMap4["port"] = port3
			}
			acrInstanceId1, _ := jsonpath.Get("$[0].acr_instance_id", v)
			if acrInstanceId1 != nil && (d.HasChange("custom_container_config.0.acr_instance_id") || acrInstanceId1 != "") {
				objectDataLocalMap4["acrInstanceId"] = acrInstanceId1
			}

			request["customContainerConfig"] = objectDataLocalMap4
		}
	}

	if d.HasChange("custom_dns") {
		update = true
		objectDataLocalMap5 := make(map[string]interface{})

		if v := d.Get("custom_dns"); v != nil {
			searches1, _ := jsonpath.Get("$[0].searches", d.Get("custom_dns"))
			if searches1 != nil && (d.HasChange("custom_dns.0.searches") || searches1 != "") {
				objectDataLocalMap5["searches"] = searches1
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
				objectDataLocalMap5["dnsOptions"] = localMaps1
			}

			nameServers1, _ := jsonpath.Get("$[0].name_servers", d.Get("custom_dns"))
			if nameServers1 != nil && (d.HasChange("custom_dns.0.name_servers") || nameServers1 != "") {
				objectDataLocalMap5["nameServers"] = nameServers1
			}

			request["customDNS"] = objectDataLocalMap5
		}
	}

	if d.HasChange("instance_lifecycle_config") {
		update = true
		objectDataLocalMap6 := make(map[string]interface{})

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

			objectDataLocalMap6["preStop"] = preStop
			initializer := make(map[string]interface{})
			handler3, _ := jsonpath.Get("$[0].initializer[0].handler", v)
			if handler3 != nil && (d.HasChange("instance_lifecycle_config.0.initializer.0.handler") || handler3 != "") {
				initializer["handler"] = handler3
			}
			timeout3, _ := jsonpath.Get("$[0].initializer[0].timeout", v)
			if timeout3 != nil && (d.HasChange("instance_lifecycle_config.0.initializer.0.timeout") || timeout3 != "") && timeout3.(int) > 0 {
				initializer["timeout"] = timeout3
			}

			objectDataLocalMap6["initializer"] = initializer

			request["instanceLifecycleConfig"] = objectDataLocalMap6
		}
	}

	if !d.IsNewResource() && d.HasChange("internet_access") {
		update = true
		request["internetAccess"] = d.Get("internet_access")
	}

	if !d.IsNewResource() && d.HasChange("oss_mount_config") {
		update = true
		objectDataLocalMap7 := make(map[string]interface{})

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
				objectDataLocalMap7["mountPoints"] = localMaps2
			}

			request["ossMountConfig"] = objectDataLocalMap7
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
		objectDataLocalMap8 := make(map[string]interface{})

		if v := d.Get("code"); v != nil {
			ossBucketName1, _ := jsonpath.Get("$[0].oss_bucket_name", v)
			if ossBucketName1 != nil && (d.HasChange("code.0.oss_bucket_name") || ossBucketName1 != "") {
				objectDataLocalMap8["ossBucketName"] = ossBucketName1
			}
			zipFile1, _ := jsonpath.Get("$[0].zip_file", v)
			if zipFile1 != nil && (d.HasChange("code.0.zip_file") || zipFile1 != "") {
				objectDataLocalMap8["zipFile"] = zipFile1
			}
			ossObjectName1, _ := jsonpath.Get("$[0].oss_object_name", v)
			if ossObjectName1 != nil && (d.HasChange("code.0.oss_object_name") || ossObjectName1 != "") {
				objectDataLocalMap8["ossObjectName"] = ossObjectName1
			}
			checksum1, _ := jsonpath.Get("$[0].checksum", v)
			if checksum1 != nil && (d.HasChange("code.0.checksum") || checksum1 != "") {
				objectDataLocalMap8["checksum"] = checksum1
			}

			request["code"] = objectDataLocalMap8
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
		objectDataLocalMap9 := make(map[string]interface{})

		if v := d.Get("vpc_config"); v != nil {
			vpcId1, _ := jsonpath.Get("$[0].vpc_id", v)
			if vpcId1 != nil && (d.HasChange("vpc_config.0.vpc_id") || vpcId1 != "") {
				objectDataLocalMap9["vpcId"] = vpcId1
			}
			securityGroupId1, _ := jsonpath.Get("$[0].security_group_id", v)
			if securityGroupId1 != nil && (d.HasChange("vpc_config.0.security_group_id") || securityGroupId1 != "") {
				objectDataLocalMap9["securityGroupId"] = securityGroupId1
			}
			vSwitchIds1, _ := jsonpath.Get("$[0].vswitch_ids", d.Get("vpc_config"))
			if vSwitchIds1 != nil && (d.HasChange("vpc_config.0.vswitch_ids") || vSwitchIds1 != "") {
				objectDataLocalMap9["vSwitchIds"] = vSwitchIds1
			}

			request["vpcConfig"] = objectDataLocalMap9
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
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2023-03-30"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
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

	if d.HasChange("tags") {
		fcv3ServiceV2 := Fcv3ServiceV2{client}
		if err := fcv3ServiceV2.SetResourceTags(d, "function"); err != nil {
			return WrapError(err)
		}
	}
	return resourceAliCloudFcv3FunctionRead(d, meta)
}

func resourceAliCloudFcv3FunctionDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	functionName := d.Id()
	action := fmt.Sprintf("/2023-03-30/functions/%s", functionName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	conn, err := client.NewFcv2Client()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["functionName"] = d.Id()

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2023-03-30"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), query, nil, nil, &runtime)

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
