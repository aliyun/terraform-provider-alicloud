package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudFcv3Functions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudFcv3FunctionRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"prefix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"functions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cpu": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"custom_container_config": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resolved_image_uri": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"entrypoint": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"acr_instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"acceleration_info": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"status": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"command": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"acceleration_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"health_check_config": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"initial_delay_seconds": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"timeout_seconds": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"http_get_url": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"period_seconds": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"failure_threshold": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"success_threshold": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"image": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"custom_dns": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"searches": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"dns_options": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"name_servers": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"custom_runtime_config": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"args": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"command": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"health_check_config": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"initial_delay_seconds": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"timeout_seconds": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"http_get_url": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"period_seconds": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"failure_threshold": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"success_threshold": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"environment_variables": {
							Type:     schema.TypeMap,
							Computed: true,
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
							Type:     schema.TypeString,
							Computed: true,
						},
						"gpu_config": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"gpu_memory_size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"gpu_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"handler": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"idle_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_concurrency": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_isolation_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_lifecycle_config": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pre_stop": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"timeout": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"handler": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"initializer": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"command": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"timeout": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"handler": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"internet_access": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"invocation_restriction": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"last_modified_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"disable": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"reason": {
										Type:     schema.TypeString,
										Computed: true,
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
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"log_config": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"project": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"log_begin_rule": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"logstore": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enable_instance_metrics": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"enable_request_metrics": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"memory_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"nas_config": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mount_points": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enable_tls": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"server_addr": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"mount_dir": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"user_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"group_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"oss_mount_config": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mount_points": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"read_only": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"bucket_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"endpoint": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"bucket_path": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"mount_dir": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"runtime": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"session_affinity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"session_affinity_config": {
							Type:     schema.TypeString,
							Computed: true,
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
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Computed: true,
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
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vswitch_ids": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"security_group_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func dataSourceAliCloudFcv3FunctionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	action := fmt.Sprintf("/2023-03-30/functions")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	if v, ok := d.GetOk("resource_group_id"); ok {
		query["resourceGroupId"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("prefix"); ok {
		query["prefix"] = StringPointer(v.(string))
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	query["limit"] = StringPointer(fmt.Sprint(PageSizeSmall))
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaGet("FC", "2023-03-30", action, query, nil, request)

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

		resp, _ := jsonpath.Get("$.functions[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["functionName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["functionName"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if nextToken, ok := response["nextToken"].(string); ok && nextToken != "" {
			query["nextToken"] = StringPointer(fmt.Sprint(nextToken))
		} else {
			break
		}
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["functionName"]

		mapping["code_size"] = objectRaw["codeSize"]
		mapping["cpu"] = objectRaw["cpu"]
		mapping["create_time"] = objectRaw["createdTime"]
		mapping["description"] = objectRaw["description"]
		mapping["disk_size"] = objectRaw["diskSize"]
		mapping["environment_variables"] = objectRaw["environmentVariables"]
		mapping["function_arn"] = objectRaw["functionArn"]
		mapping["function_id"] = objectRaw["functionId"]
		mapping["handler"] = objectRaw["handler"]
		mapping["idle_timeout"] = objectRaw["idleTimeout"]
		mapping["instance_concurrency"] = objectRaw["instanceConcurrency"]
		mapping["instance_isolation_mode"] = objectRaw["instanceIsolationMode"]
		mapping["internet_access"] = objectRaw["internetAccess"]
		mapping["last_modified_time"] = objectRaw["lastModifiedTime"]
		mapping["last_update_status"] = objectRaw["lastUpdateStatus"]
		mapping["last_update_status_reason"] = objectRaw["lastUpdateStatusReason"]
		mapping["last_update_status_reason_code"] = objectRaw["lastUpdateStatusReasonCode"]
		mapping["memory_size"] = objectRaw["memorySize"]
		mapping["resource_group_id"] = objectRaw["resourceGroupId"]
		mapping["role"] = objectRaw["role"]
		mapping["runtime"] = objectRaw["runtime"]
		mapping["session_affinity"] = objectRaw["sessionAffinity"]
		mapping["session_affinity_config"] = objectRaw["sessionAffinityConfig"]
		mapping["state"] = objectRaw["state"]
		mapping["state_reason"] = objectRaw["stateReason"]
		mapping["state_reason_code"] = objectRaw["stateReasonCode"]
		mapping["timeout"] = objectRaw["timeout"]
		mapping["function_name"] = objectRaw["functionName"]

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
				commandRaw = convertToInterfaceArray(customContainerConfigRaw["command"])
			}

			customContainerConfigMap["command"] = commandRaw
			entrypointRaw := make([]interface{}, 0)
			if customContainerConfigRaw["entrypoint"] != nil {
				entrypointRaw = convertToInterfaceArray(customContainerConfigRaw["entrypoint"])
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
		mapping["custom_container_config"] = customContainerConfigMaps
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
				for _, dnsOptionsChildRaw := range convertToInterfaceArray(dnsOptionsRaw) {
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
				nameServersRaw = convertToInterfaceArray(customDNSRaw["nameServers"])
			}

			customDnsMap["name_servers"] = nameServersRaw
			searchesRaw := make([]interface{}, 0)
			if customDNSRaw["searches"] != nil {
				searchesRaw = convertToInterfaceArray(customDNSRaw["searches"])
			}

			customDnsMap["searches"] = searchesRaw
			customDnsMaps = append(customDnsMaps, customDnsMap)
		}
		mapping["custom_dns"] = customDnsMaps
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
				argsRaw = convertToInterfaceArray(customRuntimeConfigRaw["args"])
			}

			customRuntimeConfigMap["args"] = argsRaw
			commandRaw := make([]interface{}, 0)
			if customRuntimeConfigRaw["command"] != nil {
				commandRaw = convertToInterfaceArray(customRuntimeConfigRaw["command"])
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
		mapping["custom_runtime_config"] = customRuntimeConfigMaps
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
		mapping["gpu_config"] = gpuConfigMaps
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
					commandRaw = convertToInterfaceArray(initializerRaw["command"])
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
		mapping["instance_lifecycle_config"] = instanceLifecycleConfigMaps
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
		mapping["invocation_restriction"] = invocationRestrictionMaps

		mapping["layers"] = objectRaw["layers"]
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
		mapping["log_config"] = logConfigMaps
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
				for _, mountPointsChildRaw := range convertToInterfaceArray(mountPointsRaw) {
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
		mapping["nas_config"] = nasConfigMaps
		ossMountConfigMaps := make([]map[string]interface{}, 0)
		ossMountConfigMap := make(map[string]interface{})
		mountPointsRaw, _ := jsonpath.Get("$.ossMountConfig.mountPoints", objectRaw)

		mountPointsMaps := make([]map[string]interface{}, 0)
		if mountPointsRaw != nil {
			for _, mountPointsChildRaw := range convertToInterfaceArray(mountPointsRaw) {
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
		mapping["oss_mount_config"] = ossMountConfigMaps
		tagsMaps := objectRaw["tags"]
		mapping["tags"] = tagsToMap(tagsMaps)
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
		mapping["tracing_config"] = tracingConfigMaps
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
				vSwitchIdsRaw = convertToInterfaceArray(vpcConfigRaw["vSwitchIds"])
			}

			vpcConfigMap["vswitch_ids"] = vSwitchIdsRaw
			vpcConfigMaps = append(vpcConfigMaps, vpcConfigMap)
		}
		mapping["vpc_config"] = vpcConfigMaps

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["FunctionName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("functions", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
