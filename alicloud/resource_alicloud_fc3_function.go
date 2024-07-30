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

func resourceAliCloudFc3Function() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudFc3FunctionCreate,
		Read:   resourceAliCloudFc3FunctionRead,
		Update: resourceAliCloudFc3FunctionUpdate,
		Delete: resourceAliCloudFc3FunctionDelete,
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
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(0, 86400),
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

func resourceAliCloudFc3FunctionCreate(d *schema.ResourceData, meta interface{}) error {

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
	request["functionName"] = d.Get("function_name")

	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOk("timeout"); ok && v.(int) > 0 {
		request["timeout"] = v
	}
	if v, ok := d.GetOk("memory_size"); ok && v.(int) > 0 {
		request["memorySize"] = v
	}
	request["runtime"] = d.Get("runtime")
	if v, ok := d.GetOk("environment_variables"); ok {
		request["environmentVariables"] = v
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("code"); !IsNil(v) {
		nodeNative, _ := jsonpath.Get("$[0].zip_file", d.Get("code"))
		if nodeNative != nil && nodeNative != "" {
			objectDataLocalMap["zipFile"] = nodeNative
		}
		nodeNative1, _ := jsonpath.Get("$[0].checksum", d.Get("code"))
		if nodeNative1 != nil && nodeNative1 != "" {
			objectDataLocalMap["checksum"] = nodeNative1
		}
		nodeNative2, _ := jsonpath.Get("$[0].oss_bucket_name", d.Get("code"))
		if nodeNative2 != nil && nodeNative2 != "" {
			objectDataLocalMap["ossBucketName"] = nodeNative2
		}
		nodeNative3, _ := jsonpath.Get("$[0].oss_object_name", d.Get("code"))
		if nodeNative3 != nil && nodeNative3 != "" {
			objectDataLocalMap["ossObjectName"] = nodeNative3
		}

		request["code"] = objectDataLocalMap
	}

	if v, ok := d.GetOk("cpu"); ok && v.(float64) > 0 {
		request["cpu"] = v
	}
	objectDataLocalMap1 := make(map[string]interface{})

	if v := d.Get("custom_runtime_config"); !IsNil(v) {
		nodeNative4, _ := jsonpath.Get("$[0].args", v)
		if nodeNative4 != nil && nodeNative4 != "" {
			objectDataLocalMap1["args"] = nodeNative4
		}
		nodeNative5, _ := jsonpath.Get("$[0].command", v)
		if nodeNative5 != nil && nodeNative5 != "" {
			objectDataLocalMap1["command"] = nodeNative5
		}
		healthCheckConfig := make(map[string]interface{})
		nodeNative6, _ := jsonpath.Get("$[0].health_check_config[0].failure_threshold", d.Get("custom_runtime_config"))
		if nodeNative6 != nil && nodeNative6 != "" && nodeNative6.(int) > 0 {
			healthCheckConfig["failureThreshold"] = nodeNative6
		}
		nodeNative7, _ := jsonpath.Get("$[0].health_check_config[0].http_get_url", d.Get("custom_runtime_config"))
		if nodeNative7 != nil && nodeNative7 != "" {
			healthCheckConfig["httpGetUrl"] = nodeNative7
		}
		nodeNative8, _ := jsonpath.Get("$[0].health_check_config[0].initial_delay_seconds", d.Get("custom_runtime_config"))
		if nodeNative8 != nil && nodeNative8 != "" {
			healthCheckConfig["initialDelaySeconds"] = nodeNative8
		}
		nodeNative9, _ := jsonpath.Get("$[0].health_check_config[0].period_seconds", d.Get("custom_runtime_config"))
		if nodeNative9 != nil && nodeNative9 != "" && nodeNative9.(int) > 0 {
			healthCheckConfig["periodSeconds"] = nodeNative9
		}
		nodeNative10, _ := jsonpath.Get("$[0].health_check_config[0].success_threshold", d.Get("custom_runtime_config"))
		if nodeNative10 != nil && nodeNative10 != "" && nodeNative10.(int) > 0 {
			healthCheckConfig["successThreshold"] = nodeNative10
		}
		nodeNative11, _ := jsonpath.Get("$[0].health_check_config[0].timeout_seconds", d.Get("custom_runtime_config"))
		if nodeNative11 != nil && nodeNative11 != "" && nodeNative11.(int) > 0 {
			healthCheckConfig["timeoutSeconds"] = nodeNative11
		}

		objectDataLocalMap1["healthCheckConfig"] = healthCheckConfig
		nodeNative12, _ := jsonpath.Get("$[0].port", d.Get("custom_runtime_config"))
		if nodeNative12 != nil && nodeNative12 != "" && nodeNative12.(int) > 0 {
			objectDataLocalMap1["port"] = nodeNative12
		}

		request["customRuntimeConfig"] = objectDataLocalMap1
	}

	if v, ok := d.GetOk("layers"); ok {
		layersMaps := v.([]interface{})
		request["layers"] = layersMaps
	}

	if v, ok := d.GetOk("disk_size"); ok && v.(int) > 0 {
		request["diskSize"] = v
	}
	if v, ok := d.GetOk("instance_concurrency"); ok && v.(int) > 0 {
		request["instanceConcurrency"] = v
	}
	objectDataLocalMap2 := make(map[string]interface{})

	if v := d.Get("instance_lifecycle_config"); !IsNil(v) {
		initializer := make(map[string]interface{})
		nodeNative13, _ := jsonpath.Get("$[0].initializer[0].handler", d.Get("instance_lifecycle_config"))
		if nodeNative13 != nil && nodeNative13 != "" {
			initializer["handler"] = nodeNative13
		}
		nodeNative14, _ := jsonpath.Get("$[0].initializer[0].timeout", d.Get("instance_lifecycle_config"))
		if nodeNative14 != nil && nodeNative14 != "" && nodeNative14.(int) > 0 {
			initializer["timeout"] = nodeNative14
		}

		objectDataLocalMap2["initializer"] = initializer
		preStop := make(map[string]interface{})
		nodeNative15, _ := jsonpath.Get("$[0].pre_stop[0].timeout", d.Get("instance_lifecycle_config"))
		if nodeNative15 != nil && nodeNative15 != "" && nodeNative15.(int) > 0 {
			preStop["timeout"] = nodeNative15
		}
		nodeNative16, _ := jsonpath.Get("$[0].pre_stop[0].handler", d.Get("instance_lifecycle_config"))
		if nodeNative16 != nil && nodeNative16 != "" {
			preStop["handler"] = nodeNative16
		}

		objectDataLocalMap2["preStop"] = preStop

		request["instanceLifecycleConfig"] = objectDataLocalMap2
	}

	if v, ok := d.GetOk("internet_access"); ok {
		request["internetAccess"] = v
	}
	objectDataLocalMap3 := make(map[string]interface{})

	if v := d.Get("log_config"); !IsNil(v) {
		nodeNative17, _ := jsonpath.Get("$[0].enable_instance_metrics", d.Get("log_config"))
		if nodeNative17 != nil && nodeNative17 != "" {
			objectDataLocalMap3["enableInstanceMetrics"] = nodeNative17
		}
		nodeNative18, _ := jsonpath.Get("$[0].enable_request_metrics", d.Get("log_config"))
		if nodeNative18 != nil && nodeNative18 != "" {
			objectDataLocalMap3["enableRequestMetrics"] = nodeNative18
		}
		nodeNative19, _ := jsonpath.Get("$[0].log_begin_rule", d.Get("log_config"))
		if nodeNative19 != nil && nodeNative19 != "" {
			objectDataLocalMap3["logBeginRule"] = nodeNative19
		}
		nodeNative20, _ := jsonpath.Get("$[0].logstore", d.Get("log_config"))
		if nodeNative20 != nil && nodeNative20 != "" {
			objectDataLocalMap3["logstore"] = nodeNative20
		}
		nodeNative21, _ := jsonpath.Get("$[0].project", d.Get("log_config"))
		if nodeNative21 != nil && nodeNative21 != "" {
			objectDataLocalMap3["project"] = nodeNative21
		}

		request["logConfig"] = objectDataLocalMap3
	}

	request["handler"] = d.Get("handler")
	if v, ok := d.GetOk("role"); ok {
		request["role"] = v
	}
	objectDataLocalMap4 := make(map[string]interface{})

	if v := d.Get("oss_mount_config"); !IsNil(v) {
		if v, ok := d.GetOk("oss_mount_config"); ok {
			localData1, err := jsonpath.Get("$[0].mount_points", v)
			if err != nil {
				localData1 = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop1 := range localData1.([]interface{}) {
				dataLoop1Tmp := make(map[string]interface{})
				if dataLoop1 != nil {
					dataLoop1Tmp = dataLoop1.(map[string]interface{})
				}
				dataLoop1Map := make(map[string]interface{})
				dataLoop1Map["bucketName"] = dataLoop1Tmp["bucket_name"]
				dataLoop1Map["bucketPath"] = dataLoop1Tmp["bucket_path"]
				dataLoop1Map["endpoint"] = dataLoop1Tmp["endpoint"]
				dataLoop1Map["mountDir"] = dataLoop1Tmp["mount_dir"]
				dataLoop1Map["readOnly"] = dataLoop1Tmp["read_only"]
				localMaps = append(localMaps, dataLoop1Map)
			}
			objectDataLocalMap4["mountPoints"] = localMaps
		}

		request["ossMountConfig"] = objectDataLocalMap4
	}

	objectDataLocalMap5 := make(map[string]interface{})

	if v := d.Get("vpc_config"); !IsNil(v) {
		nodeNative27, _ := jsonpath.Get("$[0].vpc_id", d.Get("vpc_config"))
		if nodeNative27 != nil && nodeNative27 != "" {
			objectDataLocalMap5["vpcId"] = nodeNative27
		}
		nodeNative28, _ := jsonpath.Get("$[0].vswitch_ids", v)
		if nodeNative28 != nil && nodeNative28 != "" {
			objectDataLocalMap5["vSwitchIds"] = nodeNative28
		}
		nodeNative29, _ := jsonpath.Get("$[0].security_group_id", d.Get("vpc_config"))
		if nodeNative29 != nil && nodeNative29 != "" {
			objectDataLocalMap5["securityGroupId"] = nodeNative29
		}

		request["vpcConfig"] = objectDataLocalMap5
	}

	objectDataLocalMap6 := make(map[string]interface{})

	if v := d.Get("custom_container_config"); !IsNil(v) {
		nodeNative30, _ := jsonpath.Get("$[0].image", d.Get("custom_container_config"))
		if nodeNative30 != nil && nodeNative30 != "" {
			objectDataLocalMap6["image"] = nodeNative30
		}
		nodeNative31, _ := jsonpath.Get("$[0].port", d.Get("custom_container_config"))
		if nodeNative31 != nil && nodeNative31 != "" {
			objectDataLocalMap6["port"] = nodeNative31
		}
		nodeNative32, _ := jsonpath.Get("$[0].entrypoint", v)
		if nodeNative32 != nil && nodeNative32 != "" {
			objectDataLocalMap6["entrypoint"] = nodeNative32
		}
		nodeNative33, _ := jsonpath.Get("$[0].command", v)
		if nodeNative33 != nil && nodeNative33 != "" {
			objectDataLocalMap6["command"] = nodeNative33
		}
		nodeNative34, _ := jsonpath.Get("$[0].acr_instance_id", d.Get("custom_container_config"))
		if nodeNative34 != nil && nodeNative34 != "" {
			objectDataLocalMap6["acrInstanceId"] = nodeNative34
		}
		nodeNative35, _ := jsonpath.Get("$[0].acceleration_type", d.Get("custom_container_config"))
		if nodeNative35 != nil && nodeNative35 != "" {
			objectDataLocalMap6["accelerationType"] = nodeNative35
		}
		healthCheckConfig1 := make(map[string]interface{})
		nodeNative36, _ := jsonpath.Get("$[0].health_check_config[0].failure_threshold", d.Get("custom_container_config"))
		if nodeNative36 != nil && nodeNative36 != "" && nodeNative36.(int) > 0 {
			healthCheckConfig1["failureThreshold"] = nodeNative36
		}
		nodeNative37, _ := jsonpath.Get("$[0].health_check_config[0].http_get_url", d.Get("custom_container_config"))
		if nodeNative37 != nil && nodeNative37 != "" {
			healthCheckConfig1["httpGetUrl"] = nodeNative37
		}
		nodeNative38, _ := jsonpath.Get("$[0].health_check_config[0].initial_delay_seconds", d.Get("custom_container_config"))
		if nodeNative38 != nil && nodeNative38 != "" {
			healthCheckConfig1["initialDelaySeconds"] = nodeNative38
		}
		nodeNative39, _ := jsonpath.Get("$[0].health_check_config[0].period_seconds", d.Get("custom_container_config"))
		if nodeNative39 != nil && nodeNative39 != "" && nodeNative39.(int) > 0 {
			healthCheckConfig1["periodSeconds"] = nodeNative39
		}
		nodeNative40, _ := jsonpath.Get("$[0].health_check_config[0].success_threshold", d.Get("custom_container_config"))
		if nodeNative40 != nil && nodeNative40 != "" && nodeNative40.(int) > 0 {
			healthCheckConfig1["successThreshold"] = nodeNative40
		}
		nodeNative41, _ := jsonpath.Get("$[0].health_check_config[0].timeout_seconds", d.Get("custom_container_config"))
		if nodeNative41 != nil && nodeNative41 != "" && nodeNative41.(int) > 0 {
			healthCheckConfig1["timeoutSeconds"] = nodeNative41
		}

		objectDataLocalMap6["healthCheckConfig"] = healthCheckConfig1

		request["customContainerConfig"] = objectDataLocalMap6
	}

	objectDataLocalMap7 := make(map[string]interface{})

	if v := d.Get("gpu_config"); !IsNil(v) {
		nodeNative42, _ := jsonpath.Get("$[0].gpu_type", d.Get("gpu_config"))
		if nodeNative42 != nil && nodeNative42 != "" {
			objectDataLocalMap7["gpuType"] = nodeNative42
		}
		nodeNative43, _ := jsonpath.Get("$[0].gpu_memory_size", d.Get("gpu_config"))
		if nodeNative43 != nil && nodeNative43 != "" {
			objectDataLocalMap7["gpuMemorySize"] = nodeNative43
		}

		request["gpuConfig"] = objectDataLocalMap7
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
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_fc3_function", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.body.functionName", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudFc3FunctionUpdate(d, meta)
}

func resourceAliCloudFc3FunctionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fc3ServiceV2 := Fc3ServiceV2{client}

	objectRaw, err := fc3ServiceV2.DescribeFc3Function(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_fc3_function DescribeFc3Function Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
	if objectRaw["handler"] != nil {
		d.Set("handler", objectRaw["handler"])
	}
	if objectRaw["instanceConcurrency"] != nil {
		d.Set("instance_concurrency", objectRaw["instanceConcurrency"])
	}
	if objectRaw["internetAccess"] != nil {
		d.Set("internet_access", objectRaw["internetAccess"])
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

func resourceAliCloudFc3FunctionUpdate(d *schema.ResourceData, meta interface{}) error {
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

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["description"] = v
	}
	if !d.IsNewResource() && d.HasChange("timeout") {
		update = true
	}
	if v, ok := d.GetOk("timeout"); (ok || d.HasChange("timeout")) && v.(int) > 0 {
		request["timeout"] = v
	}
	if !d.IsNewResource() && d.HasChange("memory_size") {
		update = true
	}
	if v, ok := d.GetOk("memory_size"); (ok || d.HasChange("memory_size")) && v.(int) > 0 {
		request["memorySize"] = v
	}
	if !d.IsNewResource() && d.HasChange("runtime") {
		update = true
	}
	request["runtime"] = d.Get("runtime")
	if !d.IsNewResource() && d.HasChange("environment_variables") {
		update = true
	}
	if v, ok := d.GetOk("environment_variables"); ok || d.HasChange("environment_variables") {
		request["environmentVariables"] = v
	}
	if d.HasChange("code") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("code"); !IsNil(v) || d.HasChange("code") {
		nodeNative, _ := jsonpath.Get("$[0].zip_file", v)
		if nodeNative != nil && (d.HasChange("code.0.zip_file") || nodeNative != "") {
			objectDataLocalMap["zipFile"] = nodeNative
		}
		nodeNative1, _ := jsonpath.Get("$[0].checksum", v)
		if nodeNative1 != nil && (d.HasChange("code.0.checksum") || nodeNative1 != "") {
			objectDataLocalMap["checksum"] = nodeNative1
		}
		nodeNative2, _ := jsonpath.Get("$[0].oss_bucket_name", v)
		if nodeNative2 != nil && (d.HasChange("code.0.oss_bucket_name") || nodeNative2 != "") {
			objectDataLocalMap["ossBucketName"] = nodeNative2
		}
		nodeNative3, _ := jsonpath.Get("$[0].oss_object_name", v)
		if nodeNative3 != nil && (d.HasChange("code.0.oss_object_name") || nodeNative3 != "") {
			objectDataLocalMap["ossObjectName"] = nodeNative3
		}

		request["code"] = objectDataLocalMap
	}

	if !d.IsNewResource() && d.HasChange("cpu") {
		update = true
	}
	if v, ok := d.GetOk("cpu"); (ok || d.HasChange("cpu")) && v.(float64) > 0 {
		request["cpu"] = v
	}
	if !d.IsNewResource() && d.HasChange("layers") {
		update = true
	}
	if v, ok := d.GetOk("layers"); ok || d.HasChange("layers") {
		layersMaps := v.([]interface{})
		request["layers"] = layersMaps
	}

	if d.HasChange("custom_runtime_config") {
		update = true
	}
	objectDataLocalMap1 := make(map[string]interface{})

	if v := d.Get("custom_runtime_config"); !IsNil(v) || d.HasChange("custom_runtime_config") {
		healthCheckConfig := make(map[string]interface{})
		nodeNative4, _ := jsonpath.Get("$[0].health_check_config[0].failure_threshold", v)
		if nodeNative4 != nil && (d.HasChange("custom_runtime_config.0.health_check_config.0.failure_threshold") || nodeNative4 != "") && nodeNative4.(int) > 0 {
			healthCheckConfig["failureThreshold"] = nodeNative4
		}
		nodeNative5, _ := jsonpath.Get("$[0].health_check_config[0].http_get_url", v)
		if nodeNative5 != nil && (d.HasChange("custom_runtime_config.0.health_check_config.0.http_get_url") || nodeNative5 != "") {
			healthCheckConfig["httpGetUrl"] = nodeNative5
		}
		nodeNative6, _ := jsonpath.Get("$[0].health_check_config[0].initial_delay_seconds", v)
		if nodeNative6 != nil && (d.HasChange("custom_runtime_config.0.health_check_config.0.initial_delay_seconds") || nodeNative6 != "") {
			healthCheckConfig["initialDelaySeconds"] = nodeNative6
		}
		nodeNative7, _ := jsonpath.Get("$[0].health_check_config[0].period_seconds", v)
		if nodeNative7 != nil && (d.HasChange("custom_runtime_config.0.health_check_config.0.period_seconds") || nodeNative7 != "") && nodeNative7.(int) > 0 {
			healthCheckConfig["periodSeconds"] = nodeNative7
		}
		nodeNative8, _ := jsonpath.Get("$[0].health_check_config[0].success_threshold", v)
		if nodeNative8 != nil && (d.HasChange("custom_runtime_config.0.health_check_config.0.success_threshold") || nodeNative8 != "") && nodeNative8.(int) > 0 {
			healthCheckConfig["successThreshold"] = nodeNative8
		}
		nodeNative9, _ := jsonpath.Get("$[0].health_check_config[0].timeout_seconds", v)
		if nodeNative9 != nil && (d.HasChange("custom_runtime_config.0.health_check_config.0.timeout_seconds") || nodeNative9 != "") && nodeNative9.(int) > 0 {
			healthCheckConfig["timeoutSeconds"] = nodeNative9
		}

		objectDataLocalMap1["healthCheckConfig"] = healthCheckConfig
		nodeNative10, _ := jsonpath.Get("$[0].port", v)
		if nodeNative10 != nil && (d.HasChange("custom_runtime_config.0.port") || nodeNative10 != "") && nodeNative10.(int) > 0 {
			objectDataLocalMap1["port"] = nodeNative10
		}
		nodeNative11, _ := jsonpath.Get("$[0].command", d.Get("custom_runtime_config"))
		if nodeNative11 != nil && (d.HasChange("custom_runtime_config.0.command") || nodeNative11 != "") {
			objectDataLocalMap1["command"] = nodeNative11
		}
		nodeNative12, _ := jsonpath.Get("$[0].args", d.Get("custom_runtime_config"))
		if nodeNative12 != nil && (d.HasChange("custom_runtime_config.0.args") || nodeNative12 != "") {
			objectDataLocalMap1["args"] = nodeNative12
		}

		request["customRuntimeConfig"] = objectDataLocalMap1
	}

	if !d.IsNewResource() && d.HasChange("disk_size") {
		update = true
	}
	if v, ok := d.GetOk("disk_size"); (ok || d.HasChange("disk_size")) && v.(int) > 0 {
		request["diskSize"] = v
	}
	if !d.IsNewResource() && d.HasChange("instance_concurrency") {
		update = true
	}
	if v, ok := d.GetOk("instance_concurrency"); (ok || d.HasChange("instance_concurrency")) && v.(int) > 0 {
		request["instanceConcurrency"] = v
	}
	if d.HasChange("instance_lifecycle_config") {
		update = true
	}
	objectDataLocalMap2 := make(map[string]interface{})

	if v := d.Get("instance_lifecycle_config"); !IsNil(v) || d.HasChange("instance_lifecycle_config") {
		initializer := make(map[string]interface{})
		nodeNative13, _ := jsonpath.Get("$[0].initializer[0].handler", v)
		if nodeNative13 != nil && (d.HasChange("instance_lifecycle_config.0.initializer.0.handler") || nodeNative13 != "") {
			initializer["handler"] = nodeNative13
		}
		nodeNative14, _ := jsonpath.Get("$[0].initializer[0].timeout", v)
		if nodeNative14 != nil && (d.HasChange("instance_lifecycle_config.0.initializer.0.timeout") || nodeNative14 != "") && nodeNative14.(int) > 0 {
			initializer["timeout"] = nodeNative14
		}

		objectDataLocalMap2["initializer"] = initializer
		preStop := make(map[string]interface{})
		nodeNative15, _ := jsonpath.Get("$[0].pre_stop[0].timeout", v)
		if nodeNative15 != nil && (d.HasChange("instance_lifecycle_config.0.pre_stop.0.timeout") || nodeNative15 != "") && nodeNative15.(int) > 0 {
			preStop["timeout"] = nodeNative15
		}
		nodeNative16, _ := jsonpath.Get("$[0].pre_stop[0].handler", v)
		if nodeNative16 != nil && (d.HasChange("instance_lifecycle_config.0.pre_stop.0.handler") || nodeNative16 != "") {
			preStop["handler"] = nodeNative16
		}

		objectDataLocalMap2["preStop"] = preStop

		request["instanceLifecycleConfig"] = objectDataLocalMap2
	}

	if !d.IsNewResource() && d.HasChange("internet_access") {
		update = true
	}
	if v, ok := d.GetOk("internet_access"); ok || d.HasChange("internet_access") {
		request["internetAccess"] = v
	}
	if d.HasChange("log_config") {
		update = true
	}
	objectDataLocalMap3 := make(map[string]interface{})

	if v := d.Get("log_config"); !IsNil(v) || d.HasChange("log_config") {
		nodeNative17, _ := jsonpath.Get("$[0].enable_instance_metrics", v)
		if nodeNative17 != nil && (d.HasChange("log_config.0.enable_instance_metrics") || nodeNative17 != "") {
			objectDataLocalMap3["enableInstanceMetrics"] = nodeNative17
		}
		nodeNative18, _ := jsonpath.Get("$[0].enable_request_metrics", v)
		if nodeNative18 != nil && (d.HasChange("log_config.0.enable_request_metrics") || nodeNative18 != "") {
			objectDataLocalMap3["enableRequestMetrics"] = nodeNative18
		}
		nodeNative19, _ := jsonpath.Get("$[0].log_begin_rule", v)
		if nodeNative19 != nil && (d.HasChange("log_config.0.log_begin_rule") || nodeNative19 != "") {
			objectDataLocalMap3["logBeginRule"] = nodeNative19
		}
		nodeNative20, _ := jsonpath.Get("$[0].logstore", v)
		if nodeNative20 != nil && (d.HasChange("log_config.0.logstore") || nodeNative20 != "") {
			objectDataLocalMap3["logstore"] = nodeNative20
		}
		nodeNative21, _ := jsonpath.Get("$[0].project", v)
		if nodeNative21 != nil && (d.HasChange("log_config.0.project") || nodeNative21 != "") {
			objectDataLocalMap3["project"] = nodeNative21
		}

		request["logConfig"] = objectDataLocalMap3
	}

	if !d.IsNewResource() && d.HasChange("handler") {
		update = true
	}
	request["handler"] = d.Get("handler")
	if !d.IsNewResource() && d.HasChange("role") {
		update = true
	}
	if v, ok := d.GetOk("role"); ok || d.HasChange("role") {
		request["role"] = v
	}
	if d.HasChange("oss_mount_config") {
		update = true
	}
	objectDataLocalMap4 := make(map[string]interface{})

	if v := d.Get("oss_mount_config"); !IsNil(v) || d.HasChange("oss_mount_config") {
		if v, ok := d.GetOk("oss_mount_config"); ok {
			localData1, err := jsonpath.Get("$[0].mount_points", v)
			if err != nil {
				localData1 = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop1 := range localData1.([]interface{}) {
				dataLoop1Tmp := make(map[string]interface{})
				if dataLoop1 != nil {
					dataLoop1Tmp = dataLoop1.(map[string]interface{})
				}
				dataLoop1Map := make(map[string]interface{})
				dataLoop1Map["bucketName"] = dataLoop1Tmp["bucket_name"]
				dataLoop1Map["bucketPath"] = dataLoop1Tmp["bucket_path"]
				dataLoop1Map["endpoint"] = dataLoop1Tmp["endpoint"]
				dataLoop1Map["mountDir"] = dataLoop1Tmp["mount_dir"]
				dataLoop1Map["readOnly"] = dataLoop1Tmp["read_only"]
				localMaps = append(localMaps, dataLoop1Map)
			}
			objectDataLocalMap4["mountPoints"] = localMaps
		}

		request["ossMountConfig"] = objectDataLocalMap4
	}

	if d.HasChange("custom_dns") {
		update = true
	}
	objectDataLocalMap5 := make(map[string]interface{})

	if v := d.Get("custom_dns"); !IsNil(v) || d.HasChange("custom_dns") {
		nodeNative27, _ := jsonpath.Get("$[0].name_servers", d.Get("custom_dns"))
		if nodeNative27 != nil && (d.HasChange("custom_dns.0.name_servers") || nodeNative27 != "") {
			objectDataLocalMap5["nameServers"] = nodeNative27
		}
		nodeNative28, _ := jsonpath.Get("$[0].searches", d.Get("custom_dns"))
		if nodeNative28 != nil && (d.HasChange("custom_dns.0.searches") || nodeNative28 != "") {
			objectDataLocalMap5["searches"] = nodeNative28
		}
		if v, ok := d.GetOk("custom_dns"); ok {
			localData2, err := jsonpath.Get("$[0].dns_options", v)
			if err != nil {
				localData2 = make([]interface{}, 0)
			}
			localMaps1 := make([]interface{}, 0)
			for _, dataLoop2 := range localData2.([]interface{}) {
				dataLoop2Tmp := make(map[string]interface{})
				if dataLoop2 != nil {
					dataLoop2Tmp = dataLoop2.(map[string]interface{})
				}
				dataLoop2Map := make(map[string]interface{})
				dataLoop2Map["name"] = dataLoop2Tmp["name"]
				dataLoop2Map["value"] = dataLoop2Tmp["value"]
				localMaps1 = append(localMaps1, dataLoop2Map)
			}
			objectDataLocalMap5["dnsOptions"] = localMaps1
		}

		request["customDNS"] = objectDataLocalMap5
	}

	if d.HasChange("vpc_config") {
		update = true
	}
	objectDataLocalMap6 := make(map[string]interface{})

	if v := d.Get("vpc_config"); !IsNil(v) || d.HasChange("vpc_config") {
		nodeNative31, _ := jsonpath.Get("$[0].vpc_id", v)
		if nodeNative31 != nil && (d.HasChange("vpc_config.0.vpc_id") || nodeNative31 != "") {
			objectDataLocalMap6["vpcId"] = nodeNative31
		}
		nodeNative32, _ := jsonpath.Get("$[0].vswitch_ids", d.Get("vpc_config"))
		if nodeNative32 != nil && (d.HasChange("vpc_config.0.vswitch_ids") || nodeNative32 != "") {
			objectDataLocalMap6["vSwitchIds"] = nodeNative32
		}
		nodeNative33, _ := jsonpath.Get("$[0].security_group_id", v)
		if nodeNative33 != nil && (d.HasChange("vpc_config.0.security_group_id") || nodeNative33 != "") {
			objectDataLocalMap6["securityGroupId"] = nodeNative33
		}

		request["vpcConfig"] = objectDataLocalMap6
	}

	if d.HasChange("custom_container_config") {
		update = true
	}
	objectDataLocalMap7 := make(map[string]interface{})

	if v := d.Get("custom_container_config"); !IsNil(v) || d.HasChange("custom_container_config") {
		nodeNative34, _ := jsonpath.Get("$[0].image", v)
		if nodeNative34 != nil && (d.HasChange("custom_container_config.0.image") || nodeNative34 != "") {
			objectDataLocalMap7["image"] = nodeNative34
		}
		nodeNative35, _ := jsonpath.Get("$[0].port", v)
		if nodeNative35 != nil && (d.HasChange("custom_container_config.0.port") || nodeNative35 != "") {
			objectDataLocalMap7["port"] = nodeNative35
		}
		nodeNative36, _ := jsonpath.Get("$[0].entrypoint", d.Get("custom_container_config"))
		if nodeNative36 != nil && (d.HasChange("custom_container_config.0.entrypoint") || nodeNative36 != "") {
			objectDataLocalMap7["entrypoint"] = nodeNative36
		}
		nodeNative37, _ := jsonpath.Get("$[0].command", d.Get("custom_container_config"))
		if nodeNative37 != nil && (d.HasChange("custom_container_config.0.command") || nodeNative37 != "") {
			objectDataLocalMap7["command"] = nodeNative37
		}
		healthCheckConfig1 := make(map[string]interface{})
		nodeNative38, _ := jsonpath.Get("$[0].health_check_config[0].failure_threshold", v)
		if nodeNative38 != nil && (d.HasChange("custom_container_config.0.health_check_config.0.failure_threshold") || nodeNative38 != "") && nodeNative38.(int) > 0 {
			healthCheckConfig1["failureThreshold"] = nodeNative38
		}
		nodeNative39, _ := jsonpath.Get("$[0].health_check_config[0].http_get_url", v)
		if nodeNative39 != nil && (d.HasChange("custom_container_config.0.health_check_config.0.http_get_url") || nodeNative39 != "") {
			healthCheckConfig1["httpGetUrl"] = nodeNative39
		}
		nodeNative40, _ := jsonpath.Get("$[0].health_check_config[0].initial_delay_seconds", v)
		if nodeNative40 != nil && (d.HasChange("custom_container_config.0.health_check_config.0.initial_delay_seconds") || nodeNative40 != "") {
			healthCheckConfig1["initialDelaySeconds"] = nodeNative40
		}
		nodeNative41, _ := jsonpath.Get("$[0].health_check_config[0].period_seconds", v)
		if nodeNative41 != nil && (d.HasChange("custom_container_config.0.health_check_config.0.period_seconds") || nodeNative41 != "") && nodeNative41.(int) > 0 {
			healthCheckConfig1["periodSeconds"] = nodeNative41
		}
		nodeNative42, _ := jsonpath.Get("$[0].health_check_config[0].success_threshold", v)
		if nodeNative42 != nil && (d.HasChange("custom_container_config.0.health_check_config.0.success_threshold") || nodeNative42 != "") && nodeNative42.(int) > 0 {
			healthCheckConfig1["successThreshold"] = nodeNative42
		}
		nodeNative43, _ := jsonpath.Get("$[0].health_check_config[0].timeout_seconds", v)
		if nodeNative43 != nil && (d.HasChange("custom_container_config.0.health_check_config.0.timeout_seconds") || nodeNative43 != "") && nodeNative43.(int) > 0 {
			healthCheckConfig1["timeoutSeconds"] = nodeNative43
		}

		objectDataLocalMap7["healthCheckConfig"] = healthCheckConfig1
		nodeNative44, _ := jsonpath.Get("$[0].acr_instance_id", v)
		if nodeNative44 != nil && (d.HasChange("custom_container_config.0.acr_instance_id") || nodeNative44 != "") {
			objectDataLocalMap7["acrInstanceId"] = nodeNative44
		}
		nodeNative45, _ := jsonpath.Get("$[0].acceleration_type", v)
		if nodeNative45 != nil && (d.HasChange("custom_container_config.0.acceleration_type") || nodeNative45 != "") {
			objectDataLocalMap7["accelerationType"] = nodeNative45
		}

		request["customContainerConfig"] = objectDataLocalMap7
	}

	if d.HasChange("gpu_config") {
		update = true
	}
	objectDataLocalMap8 := make(map[string]interface{})

	if v := d.Get("gpu_config"); !IsNil(v) || d.HasChange("gpu_config") {
		nodeNative46, _ := jsonpath.Get("$[0].gpu_type", v)
		if nodeNative46 != nil && (d.HasChange("gpu_config.0.gpu_type") || nodeNative46 != "") {
			objectDataLocalMap8["gpuType"] = nodeNative46
		}
		nodeNative47, _ := jsonpath.Get("$[0].gpu_memory_size", v)
		if nodeNative47 != nil && (d.HasChange("gpu_config.0.gpu_memory_size") || nodeNative47 != "") {
			objectDataLocalMap8["gpuMemorySize"] = nodeNative47
		}

		request["gpuConfig"] = objectDataLocalMap8
	}

	if d.HasChange("nas_config") {
		update = true
	}
	objectDataLocalMap9 := make(map[string]interface{})

	if v := d.Get("nas_config"); !IsNil(v) || d.HasChange("nas_config") {
		nodeNative48, _ := jsonpath.Get("$[0].group_id", v)
		if nodeNative48 != nil && (d.HasChange("nas_config.0.group_id") || nodeNative48 != "") {
			objectDataLocalMap9["groupId"] = nodeNative48
		}
		nodeNative49, _ := jsonpath.Get("$[0].user_id", v)
		if nodeNative49 != nil && (d.HasChange("nas_config.0.user_id") || nodeNative49 != "") {
			objectDataLocalMap9["userId"] = nodeNative49
		}
		if v, ok := d.GetOk("nas_config"); ok {
			localData3, err := jsonpath.Get("$[0].mount_points", v)
			if err != nil {
				localData3 = make([]interface{}, 0)
			}
			localMaps2 := make([]interface{}, 0)
			for _, dataLoop3 := range localData3.([]interface{}) {
				dataLoop3Tmp := make(map[string]interface{})
				if dataLoop3 != nil {
					dataLoop3Tmp = dataLoop3.(map[string]interface{})
				}
				dataLoop3Map := make(map[string]interface{})
				dataLoop3Map["serverAddr"] = dataLoop3Tmp["server_addr"]
				dataLoop3Map["mountDir"] = dataLoop3Tmp["mount_dir"]
				dataLoop3Map["enableTLS"] = dataLoop3Tmp["enable_tls"]
				localMaps2 = append(localMaps2, dataLoop3Map)
			}
			objectDataLocalMap9["mountPoints"] = localMaps2
		}

		request["nasConfig"] = objectDataLocalMap9
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
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudFc3FunctionRead(d, meta)
}

func resourceAliCloudFc3FunctionDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	functionName := d.Id()
	action := fmt.Sprintf("/2023-03-30/functions/%s", functionName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	conn, err := client.NewFcv2Client()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["functionName"] = d.Id()

	body = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2023-03-30"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)

		if err != nil {
			if IsExpectedErrors(err, []string{"429"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"FunctionNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
