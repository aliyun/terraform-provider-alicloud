package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEciContainerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEciContainerGroupCreate,
		Read:   resourceAlicloudEciContainerGroupRead,
		Update: resourceAlicloudEciContainerGroupUpdate,
		Delete: resourceAlicloudEciContainerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"container_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if old != "" && new != "" && old != new {
						return strings.Contains(new, old)
					}
					return false
				},
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"cpu": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"memory": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"ram_role_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"restart_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Always", "Never", "OnFailure"}, false),
			},
			"auto_match_image_cache": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"plain_http_registry": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"insecure_registry": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auto_create_eip": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"eip_bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"eip_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
			"containers": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
						"image": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cpu": {
							Type:     schema.TypeFloat,
							Optional: true,
							Default:  0,
						},
						"gpu": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Default:  0,
						},
						"memory": {
							Type:     schema.TypeFloat,
							Optional: true,
							Default:  0,
						},
						"image_pull_policy": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "IfNotPresent",
							ValidateFunc: StringInSlice([]string{"Always", "IfNotPresent", "Never"}, false),
						},
						"working_dir": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"commands": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"args": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"ports": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"environment_vars": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"field_ref": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"field_path": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
						"volume_mounts": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mount_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"read_only": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
								},
							},
						},
						"liveness_probe": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"initial_delay_seconds": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"period_seconds": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"timeout_seconds": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"success_threshold": {
										Type:     schema.TypeInt,
										ForceNew: true,
										Optional: true,
									},
									"failure_threshold": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"exec": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"commands": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"tcp_socket": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
									},
									"http_get": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"scheme": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"port": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"path": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
						"readiness_probe": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"initial_delay_seconds": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"period_seconds": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"timeout_seconds": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"success_threshold": {
										Type:     schema.TypeInt,
										ForceNew: true,
										Optional: true,
									},
									"failure_threshold": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"exec": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"commands": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"tcp_socket": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
									},
									"http_get": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"scheme": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"port": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"path": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
						"ready": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"restart_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"security_context": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"capability": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"add": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"run_as_user": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"privileged": {
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						"lifecycle_pre_stop_handler_exec": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"init_containers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"cpu": {
							Type:     schema.TypeFloat,
							Optional: true,
							Default:  0,
						},
						"gpu": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Default:  0,
						},
						"memory": {
							Type:     schema.TypeFloat,
							Optional: true,
							Default:  0,
						},
						"image": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"image_pull_policy": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "IfNotPresent",
							ValidateFunc: StringInSlice([]string{"Always", "IfNotPresent", "Never"}, false),
						},
						"working_dir": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"commands": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"args": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"ports": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"environment_vars": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"field_ref": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"field_path": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
						"volume_mounts": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mount_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"read_only": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
								},
							},
						},
						"ready": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"restart_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"security_context": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"capability": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"add": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"run_as_user": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"dns_policy": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"dns_config": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name_servers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"searches": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"options": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"eci_security_context": {
				Type:     schema.TypeSet,
				Removed:  "Field 'eci_security_context' has been removed from provider version ?",
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sysctls": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
			"security_context": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sysctl": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
			"host_aliases": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"hostnames": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"volumes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"disk_volume_disk_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"disk_volume_fs_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"flex_volume_driver": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"flex_volume_fs_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"flex_volume_options": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return old != "" && new != ""
							},
						},
						"nfs_volume_path": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"nfs_volume_server": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"nfs_volume_read_only": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: true,
						},
						"config_file_volume_config_file_to_paths": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"content": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"path": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
			"image_registry_credential": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"password": {
							Type:     schema.TypeString,
							Required: true,
						},
						"server": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"acr_registry_info": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"domains": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"internet_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"intranet_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"termination_grace_period_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"spot_price_limit": {
				Type:     schema.TypeFloat,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"spot_strategy": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"NoSpot", "SpotAsPriceGo", "SpotWithPriceLimit"}, false),
			}},
	}
}

func resourceAlicloudEciContainerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eciService := EciService{client}
	var response map[string]interface{}
	action := "CreateContainerGroup"
	request := make(map[string]interface{})
	var err error
	request["ContainerGroupName"] = d.Get("container_group_name")
	Containers := make([]map[string]interface{}, len(d.Get("containers").([]interface{})))
	for i, ContainersValue := range d.Get("containers").([]interface{}) {
		ContainersMap := ContainersValue.(map[string]interface{})
		Containers[i] = make(map[string]interface{})
		Containers[i]["Arg"] = ContainersMap["args"]
		Containers[i]["Command"] = ContainersMap["commands"]
		Containers[i]["Cpu"] = ContainersMap["cpu"]
		EnvironmentVars := make([]map[string]interface{}, len(ContainersMap["environment_vars"].([]interface{})))
		for i, EnvironmentVarsValue := range ContainersMap["environment_vars"].([]interface{}) {
			EnvironmentVarsMap := EnvironmentVarsValue.(map[string]interface{})
			EnvironmentVars[i] = make(map[string]interface{})
			EnvironmentVars[i]["Key"] = EnvironmentVarsMap["key"]
			EnvironmentVars[i]["Value"] = EnvironmentVarsMap["value"]
			for _, FieldRefValue := range EnvironmentVarsMap["field_ref"].([]interface{}) {
				FieldRef := map[string]interface{}{}
				FieldRefMap := FieldRefValue.(map[string]interface{})
				FieldRef["FieldPath"] = FieldRefMap["field_path"]
				EnvironmentVars[i]["FieldRef"] = FieldRef
			}
		}
		Containers[i]["EnvironmentVar"] = EnvironmentVars

		Containers[i]["Gpu"] = ContainersMap["gpu"]
		Containers[i]["Image"] = ContainersMap["image"]
		Containers[i]["ImagePullPolicy"] = ContainersMap["image_pull_policy"]
		Containers[i]["Memory"] = ContainersMap["memory"]
		Containers[i]["Name"] = ContainersMap["name"]
		Ports := make([]map[string]interface{}, len(ContainersMap["ports"].([]interface{})))
		for i, PortsValue := range ContainersMap["ports"].([]interface{}) {
			PortsMap := PortsValue.(map[string]interface{})
			Ports[i] = make(map[string]interface{})
			Ports[i]["Port"] = PortsMap["port"]
			Ports[i]["Protocol"] = PortsMap["protocol"]
		}
		Containers[i]["Port"] = Ports

		VolumeMounts := make([]map[string]interface{}, len(ContainersMap["volume_mounts"].([]interface{})))
		for i, VolumeMountsValue := range ContainersMap["volume_mounts"].([]interface{}) {
			VolumeMountsMap := VolumeMountsValue.(map[string]interface{})
			VolumeMounts[i] = make(map[string]interface{})
			VolumeMounts[i]["MountPath"] = VolumeMountsMap["mount_path"]
			VolumeMounts[i]["Name"] = VolumeMountsMap["name"]
			VolumeMounts[i]["ReadOnly"] = VolumeMountsMap["read_only"]
		}
		Containers[i]["VolumeMount"] = VolumeMounts

		Containers[i]["WorkingDir"] = ContainersMap["working_dir"]

		LivenessProbe := map[string]interface{}{}
		for _, LivenessProbesValue := range ContainersMap["liveness_probe"].([]interface{}) {

			LivenessProbesMap := LivenessProbesValue.(map[string]interface{})
			HttpGetValue := map[string]interface{}{}
			for _, HttpGetValues := range LivenessProbesMap["http_get"].([]interface{}) {
				HttpGetValueMap := HttpGetValues.(map[string]interface{})
				HttpGetValue["Scheme"] = HttpGetValueMap["scheme"]
				HttpGetValue["Port"] = HttpGetValueMap["port"]
				HttpGetValue["Path"] = HttpGetValueMap["path"]
			}
			ExecsValue := map[string]interface{}{}
			for _, ExecsValues := range LivenessProbesMap["exec"].([]interface{}) {
				ExecsValuesValueMap := ExecsValues.(map[string]interface{})
				ExecsValue["Command"] = ExecsValuesValueMap["commands"]
			}
			TcpSocketValue := map[string]interface{}{}
			for _, TcpSocketsValues := range LivenessProbesMap["tcp_socket"].([]interface{}) {
				TcpSocketsValueMap := TcpSocketsValues.(map[string]interface{})
				TcpSocketValue["Port"] = TcpSocketsValueMap["port"]
			}

			LivenessProbe["PeriodSeconds"] = LivenessProbesMap["period_seconds"]
			LivenessProbe["InitialDelaySeconds"] = LivenessProbesMap["initial_delay_seconds"]
			LivenessProbe["SuccessThreshold"] = LivenessProbesMap["success_threshold"]
			LivenessProbe["FailureThreshold"] = LivenessProbesMap["failure_threshold"]
			LivenessProbe["TimeoutSeconds"] = LivenessProbesMap["timeout_seconds"]
			LivenessProbe["HttpGet"] = HttpGetValue
			LivenessProbe["Exec"] = ExecsValue
			LivenessProbe["TcpSocket"] = TcpSocketValue
		}
		Containers[i]["LivenessProbe"] = LivenessProbe

		ReadinessProbe := map[string]interface{}{}
		for _, ReadinessProbesValue := range ContainersMap["readiness_probe"].([]interface{}) {

			ReadinessProbesMap := ReadinessProbesValue.(map[string]interface{})
			HttpGetValue := map[string]interface{}{}
			for _, HttpGetValues := range ReadinessProbesMap["http_get"].([]interface{}) {
				HttpGetValueMap := HttpGetValues.(map[string]interface{})
				HttpGetValue["Scheme"] = HttpGetValueMap["scheme"]
				HttpGetValue["Port"] = HttpGetValueMap["port"]
				HttpGetValue["Path"] = HttpGetValueMap["path"]

			}
			ExecsValue := map[string]interface{}{}
			for _, ExecsValues := range ReadinessProbesMap["exec"].([]interface{}) {
				ExecsValuesValueMap := ExecsValues.(map[string]interface{})
				ExecsValue["Command"] = ExecsValuesValueMap["commands"]
			}
			TcpSocketValue := map[string]interface{}{}
			for _, TcpSocketsValues := range ReadinessProbesMap["tcp_socket"].([]interface{}) {
				TcpSocketsValueMap := TcpSocketsValues.(map[string]interface{})
				TcpSocketValue["Port"] = TcpSocketsValueMap["port"]
			}

			ReadinessProbe["PeriodSeconds"] = ReadinessProbesMap["period_seconds"]
			ReadinessProbe["InitialDelaySeconds"] = ReadinessProbesMap["initial_delay_seconds"]
			ReadinessProbe["SuccessThreshold"] = ReadinessProbesMap["success_threshold"]
			ReadinessProbe["FailureThreshold"] = ReadinessProbesMap["failure_threshold"]
			ReadinessProbe["TimeoutSeconds"] = ReadinessProbesMap["timeout_seconds"]
			ReadinessProbe["HttpGet"] = HttpGetValue
			ReadinessProbe["Exec"] = ExecsValue
			ReadinessProbe["TcpSocket"] = TcpSocketValue
		}
		Containers[i]["ReadinessProbe"] = ReadinessProbe
		SecurityContext := map[string]interface{}{}
		for _, SecurityContextValue := range ContainersMap["security_context"].([]interface{}) {
			SecurityContextMap := SecurityContextValue.(map[string]interface{})
			if SecurityContextMap["capability"] != nil {
				for _, v := range SecurityContextMap["capability"].([]interface{}) {
					CapabilityValue := map[string]interface{}{}
					CapabilityValue["Add"] = v.(map[string]interface{})["add"]
					SecurityContext["Capability"] = CapabilityValue
				}
			}
			SecurityContext["RunAsUser"] = SecurityContextMap["run_as_user"]
			Containers[i]["SecurityContextPrivileged"] = SecurityContextMap["privileged"]
		}
		Containers[i]["SecurityContext"] = SecurityContext

		Containers[i]["LifecyclePreStopHandlerExec"] = ContainersMap["lifecycle_pre_stop_handler_exec"]

	}
	request["Container"] = Containers

	if v, ok := d.GetOk("acr_registry_info"); ok {
		AcrRegistryInfos := make([]map[string]interface{}, len(v.(*schema.Set).List()))
		for i, AcrRegistryInfosValue := range v.(*schema.Set).List() {
			AcrRegistryInfosMap := AcrRegistryInfosValue.(map[string]interface{})
			AcrRegistryInfos[i] = make(map[string]interface{})
			AcrRegistryInfos[i]["InstanceName"] = AcrRegistryInfosMap["instance_name"]
			AcrRegistryInfos[i]["InstanceId"] = AcrRegistryInfosMap["instance_id"]
			AcrRegistryInfos[i]["Domain"] = AcrRegistryInfosMap["domains"]
			AcrRegistryInfos[i]["RegionId"] = AcrRegistryInfosMap["region_id"]
		}
		request["AcrRegistryInfo"] = AcrRegistryInfos
	}

	if v, ok := d.GetOk("cpu"); ok {
		request["Cpu"] = v
	}

	if v, ok := d.GetOk("dns_policy"); ok {
		request["DnsPolicy"] = v
	}

	if v, ok := d.GetOk("dns_config"); ok {
		if v != nil {
			dnsConfigMap := make(map[string]interface{})
			for _, dnsConfig := range v.(*schema.Set).List() {
				dnsConfigArg := dnsConfig.(map[string]interface{})
				dnsConfigMap["NameServer"] = dnsConfigArg["name_servers"]
				if dnsConfigArg["options"] != nil {
					optionsMaps := make([]map[string]interface{}, 0)
					for _, options := range dnsConfigArg["options"].([]interface{}) {
						optionsMap := make(map[string]interface{})
						optionsArg := options.(map[string]interface{})
						optionsMap["Name"] = optionsArg["name"]
						optionsMap["Value"] = optionsArg["value"]
						optionsMaps = append(optionsMaps, optionsMap)
					}
					dnsConfigMap["Option"] = optionsMaps
				}
				dnsConfigMap["Search"] = dnsConfigArg["searches"]
			}
			request["DnsConfig"] = dnsConfigMap
		}
	}
	if v, ok := d.GetOk("eci_security_context"); ok {
		if v != nil {
			eciSecurityContextMap := make(map[string]interface{})
			for _, eciSecurityContext := range v.(*schema.Set).List() {
				eciSecurityContextArg := eciSecurityContext.(map[string]interface{})
				if eciSecurityContextArg["sysctls"] != nil {
					sysctlsMaps := make([]map[string]interface{}, 0)
					for _, sysctls := range eciSecurityContextArg["sysctls"].([]interface{}) {
						sysctlsMap := make(map[string]interface{})
						sysctlsArg := sysctls.(map[string]interface{})
						sysctlsMap["Name"] = sysctlsArg["name"]
						sysctlsMap["Value"] = sysctlsArg["value"]
						sysctlsMaps = append(sysctlsMaps, sysctlsMap)
					}
					eciSecurityContextMap["Sysctls"] = sysctlsMaps
				}
			}
			request["EciSecurityContext"] = eciSecurityContextMap
		}
	}
	if v, ok := d.GetOk("security_context"); ok {
		if v != nil {
			eciSecurityContextMap := make(map[string]interface{})
			for _, eciSecurityContext := range v.(*schema.Set).List() {
				eciSecurityContextArg := eciSecurityContext.(map[string]interface{})
				if eciSecurityContextArg["sysctl"] != nil {
					sysctlsMaps := make([]map[string]interface{}, 0)
					for _, sysctls := range eciSecurityContextArg["sysctl"].([]interface{}) {
						sysctlsMap := make(map[string]interface{})
						sysctlsArg := sysctls.(map[string]interface{})
						sysctlsMap["Name"] = sysctlsArg["name"]
						sysctlsMap["Value"] = sysctlsArg["value"]
						sysctlsMaps = append(sysctlsMaps, sysctlsMap)
					}
					eciSecurityContextMap["Sysctl"] = sysctlsMaps
				}
			}
			request["SecurityContext"] = eciSecurityContextMap
		}
	}
	if v, ok := d.GetOk("host_aliases"); ok {
		HostAliases := make([]map[string]interface{}, len(v.([]interface{})))
		for i, HostAliasesValue := range v.([]interface{}) {
			HostAliasesMap := HostAliasesValue.(map[string]interface{})
			HostAliases[i] = make(map[string]interface{})
			HostAliases[i]["Hostname"] = HostAliasesMap["hostnames"]
			HostAliases[i]["Ip"] = HostAliasesMap["ip"]
		}
		request["HostAliase"] = HostAliases

	}

	if v, ok := d.GetOk("init_containers"); ok {
		InitContainers := make([]map[string]interface{}, len(v.([]interface{})))
		for i, InitContainersValue := range v.([]interface{}) {
			InitContainersMap := InitContainersValue.(map[string]interface{})
			InitContainers[i] = make(map[string]interface{})
			InitContainers[i]["Arg"] = InitContainersMap["args"]
			InitContainers[i]["Command"] = InitContainersMap["commands"]
			InitContainers[i]["Cpu"] = InitContainersMap["cpu"]
			EnvironmentVars := make([]map[string]interface{}, len(InitContainersMap["environment_vars"].([]interface{})))
			for i, EnvironmentVarsValue := range InitContainersMap["environment_vars"].([]interface{}) {
				EnvironmentVarsMap := EnvironmentVarsValue.(map[string]interface{})
				EnvironmentVars[i] = make(map[string]interface{})
				EnvironmentVars[i]["Key"] = EnvironmentVarsMap["key"]
				EnvironmentVars[i]["Value"] = EnvironmentVarsMap["value"]
				for _, FieldRefValue := range EnvironmentVarsMap["field_ref"].([]interface{}) {
					FieldRef := map[string]interface{}{}
					FieldRefMap := FieldRefValue.(map[string]interface{})
					FieldRef["FieldPath"] = FieldRefMap["field_path"]
					EnvironmentVars[i]["FieldRef"] = FieldRef
				}
			}
			InitContainers[i]["EnvironmentVar"] = EnvironmentVars

			InitContainers[i]["Gpu"] = InitContainersMap["gpu"]
			InitContainers[i]["Image"] = InitContainersMap["image"]
			InitContainers[i]["ImagePullPolicy"] = InitContainersMap["image_pull_policy"]
			InitContainers[i]["Memory"] = InitContainersMap["memory"]
			InitContainers[i]["Name"] = InitContainersMap["name"]
			Ports := make([]map[string]interface{}, len(InitContainersMap["ports"].([]interface{})))
			for i, PortsValue := range InitContainersMap["ports"].([]interface{}) {
				PortsMap := PortsValue.(map[string]interface{})
				Ports[i] = make(map[string]interface{})
				Ports[i]["Port"] = PortsMap["port"]
				Ports[i]["Protocol"] = PortsMap["protocol"]
			}
			InitContainers[i]["Port"] = Ports

			VolumeMounts := make([]map[string]interface{}, len(InitContainersMap["volume_mounts"].([]interface{})))
			for i, VolumeMountsValue := range InitContainersMap["volume_mounts"].([]interface{}) {
				VolumeMountsMap := VolumeMountsValue.(map[string]interface{})
				VolumeMounts[i] = make(map[string]interface{})
				VolumeMounts[i]["MountPath"] = VolumeMountsMap["mount_path"]
				VolumeMounts[i]["Name"] = VolumeMountsMap["name"]
				VolumeMounts[i]["ReadOnly"] = VolumeMountsMap["read_only"]
			}
			InitContainers[i]["VolumeMount"] = VolumeMounts

			InitContainers[i]["WorkingDir"] = InitContainersMap["working_dir"]

			SecurityContext := map[string]interface{}{}
			for _, SecurityContextValue := range InitContainersMap["security_context"].([]interface{}) {
				SecurityContextMap := SecurityContextValue.(map[string]interface{})
				if SecurityContextMap["capability"] != nil {
					for _, v := range SecurityContextMap["capability"].([]interface{}) {
						CapabilityValue := map[string]interface{}{}
						CapabilityValue["Add"] = v.(map[string]interface{})["add"]
						SecurityContext["Capability"] = CapabilityValue
					}
				}
				SecurityContext["RunAsUser"] = SecurityContextMap["run_as_user"]
			}
			InitContainers[i]["SecurityContext"] = SecurityContext
		}
		request["InitContainer"] = InitContainers

	}

	if v, ok := d.GetOk("instance_type"); ok {
		request["InstanceType"] = v
	}

	if v, ok := d.GetOk("memory"); ok {
		request["Memory"] = v
	}

	if v, ok := d.GetOk("ram_role_name"); ok {
		request["RamRoleName"] = v
	}

	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("restart_policy"); ok {
		request["RestartPolicy"] = v
	}

	request["SecurityGroupId"] = d.Get("security_group_id")
	if v, ok := d.GetOk("tags"); ok {
		count := 1
		for key, value := range v.(map[string]interface{}) {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
	}
	if v, ok := d.GetOk("volumes"); ok {
		Volumes := make([]map[string]interface{}, len(v.([]interface{})))
		for i, VolumesValue := range v.([]interface{}) {
			VolumesMap := VolumesValue.(map[string]interface{})
			Volumes[i] = make(map[string]interface{})
			ConfigFileVolumeConfigFileToPaths := make([]map[string]interface{}, len(VolumesMap["config_file_volume_config_file_to_paths"].([]interface{})))
			for j, ConfigFileVolumeConfigFileToPathsValue := range VolumesMap["config_file_volume_config_file_to_paths"].([]interface{}) {
				ConfigFileVolumeConfigFileToPathsMap := ConfigFileVolumeConfigFileToPathsValue.(map[string]interface{})
				ConfigFileVolumeConfigFileToPaths[j] = make(map[string]interface{})
				ConfigFileVolumeConfigFileToPaths[j]["Content"] = ConfigFileVolumeConfigFileToPathsMap["content"]
				ConfigFileVolumeConfigFileToPaths[j]["Path"] = ConfigFileVolumeConfigFileToPathsMap["path"]
			}
			Volumes[i]["ConfigFileVolume.ConfigFileToPath"] = ConfigFileVolumeConfigFileToPaths

			Volumes[i]["DiskVolume.DiskId"] = VolumesMap["disk_volume_disk_id"]
			Volumes[i]["DiskVolume.FsType"] = VolumesMap["disk_volume_fs_type"]
			Volumes[i]["FlexVolume.Driver"] = VolumesMap["flex_volume_driver"]
			Volumes[i]["FlexVolume.FsType"] = VolumesMap["flex_volume_fs_type"]
			Volumes[i]["FlexVolume.Options"] = VolumesMap["flex_volume_options"]
			Volumes[i]["NFSVolume.Path"] = VolumesMap["nfs_volume_path"]
			Volumes[i]["NFSVolume.Server"] = VolumesMap["nfs_volume_server"]
			Volumes[i]["NFSVolume.ReadOnly"] = VolumesMap["nfs_volume_read_only"]
			Volumes[i]["Name"] = VolumesMap["name"]
			Volumes[i]["Type"] = VolumesMap["type"]
		}
		request["Volume"] = Volumes

	}

	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}

	vswitchId := Trim(d.Get("vswitch_id").(string))
	request["VSwitchId"] = vswitchId

	if v, ok := d.GetOk("image_registry_credential"); ok {
		imageRegisryCredentialMaps := make([]map[string]interface{}, 0)
		for _, raw := range v.(*schema.Set).List() {
			obj := raw.(map[string]interface{})
			imageRegisryCredentialMaps = append(imageRegisryCredentialMaps, map[string]interface{}{
				"Password": obj["password"],
				"Server":   obj["server"],
				"UserName": obj["user_name"],
			})
		}
		request["ImageRegistryCredential"] = imageRegisryCredentialMaps
	}
	request["AutoMatchImageCache"] = d.Get("auto_match_image_cache")
	if v, ok := d.GetOkExists("auto_create_eip"); ok {
		request["AutoCreateEip"] = v
	}
	if v, ok := d.GetOkExists("eip_bandwidth"); ok {
		request["EipBandwidth"] = v
	}
	if v, ok := d.GetOk("eip_instance_id"); ok {
		request["EipInstanceId"] = v
	}

	if v, ok := d.GetOk("plain_http_registry"); ok {
		request["PlainHttpRegistry"] = v
	}

	if v, ok := d.GetOk("insecure_registry"); ok {
		request["InsecureRegistry"] = v
	}

	if v, ok := d.GetOk("termination_grace_period_seconds"); ok {
		request["TerminationGracePeriodSeconds"] = v
	}

	if v, ok := d.GetOk("spot_strategy"); ok {
		request["SpotStrategy"] = v
	}
	if v, ok := d.GetOk("spot_price_limit"); ok {
		request["SpotPriceLimit"] = v
	}

	request["ClientToken"] = buildClientToken("CreateContainerGroup")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Eci", "2018-08-08", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_eci_container_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ContainerGroupId"]))

	stateConf := BuildStateConf([]string{}, []string{"Running", "Succeeded"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, eciService.EciContainerGroupStateRefreshFunc(d.Id(), []string{"Failed", "ScheduleFailed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudEciContainerGroupRead(d, meta)
}

func resourceAlicloudEciContainerGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eciService := EciService{client}
	object, err := eciService.DescribeEciContainerGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_eci_openapi_container_group eciService.DescribeEciContainerGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("container_group_name", object["ContainerGroupName"])
	d.Set("internet_ip", object["InternetIp"])
	d.Set("intranet_ip", object["IntranetIp"])

	containers := make([]map[string]interface{}, 0)
	if containersList, ok := object["Containers"].([]interface{}); ok {
		for _, v := range containersList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"args":              m1["Args"],
					"commands":          m1["Commands"],
					"cpu":               m1["Cpu"],
					"gpu":               m1["Gpu"],
					"image":             m1["Image"],
					"image_pull_policy": m1["ImagePullPolicy"],
					"memory":            m1["Memory"],
					"name":              m1["Name"],
					"ready":             m1["Ready"],
					"restart_count":     m1["RestartCount"],
					"working_dir":       m1["WorkingDir"],
				}
				if m1["EnvironmentVars"] != nil {
					environmentVarsMaps := make([]map[string]interface{}, 0)
					for _, environmentVarsValue := range m1["EnvironmentVars"].([]interface{}) {
						environmentVars := environmentVarsValue.(map[string]interface{})
						environmentVarsMap := map[string]interface{}{
							"key":   environmentVars["Key"],
							"value": environmentVars["Value"],
						}
						if environmentVars["ValueFrom"] != nil {
							fieldRefMaps := make([]map[string]interface{}, 0)
							fieldRefValue := environmentVars["ValueFrom"].(map[string]interface{})["FieldRef"]
							if fieldRefValue != nil && fieldRefValue.(map[string]interface{})["FieldPath"] != nil {
								fieldRefMap := map[string]interface{}{
									"field_path": fieldRefValue.(map[string]interface{})["FieldPath"],
								}
								fieldRefMaps = append(fieldRefMaps, fieldRefMap)
								environmentVarsMap["field_ref"] = fieldRefMaps
							}
						}
						environmentVarsMaps = append(environmentVarsMaps, environmentVarsMap)
					}
					temp1["environment_vars"] = environmentVarsMaps
				}
				if m1["Ports"] != nil {
					portsMaps := make([]map[string]interface{}, 0)
					for _, portsValue := range m1["Ports"].([]interface{}) {
						ports := portsValue.(map[string]interface{})
						portsMap := map[string]interface{}{
							"port":     ports["Port"],
							"protocol": ports["Protocol"],
						}
						portsMaps = append(portsMaps, portsMap)
					}
					temp1["ports"] = portsMaps
				}
				if m1["VolumeMounts"] != nil {
					volumeMountsMaps := make([]map[string]interface{}, 0)
					for _, volumeMountsValue := range m1["VolumeMounts"].([]interface{}) {
						volumeMounts := volumeMountsValue.(map[string]interface{})
						volumeMountsMap := map[string]interface{}{
							"mount_path": volumeMounts["MountPath"],
							"name":       volumeMounts["Name"],
							"read_only":  volumeMounts["ReadOnly"],
						}
						volumeMountsMaps = append(volumeMountsMaps, volumeMountsMap)
					}
					temp1["volume_mounts"] = volumeMountsMaps
				}
				if m1["ReadinessProbe"] != nil {
					ReadinessProbesMaps := make([]map[string]interface{}, 0)
					ReadinessProbes := m1["ReadinessProbe"].(map[string]interface{})
					ReadinessProbesMap := map[string]interface{}{
						"initial_delay_seconds": ReadinessProbes["InitialDelaySeconds"],
						"timeout_seconds":       ReadinessProbes["TimeoutSeconds"],
						"period_seconds":        ReadinessProbes["PeriodSeconds"],
						"success_threshold":     ReadinessProbes["SuccessThreshold"],
						"failure_threshold":     ReadinessProbes["FailureThreshold"],
					}

					Execs := make([]map[string]interface{}, 0)
					if ExecsList, ok := ReadinessProbes["Execs"]; ok && ExecsList != nil {
						ExecsMap := map[string]interface{}{
							"commands": ExecsList.([]interface{}),
						}
						Execs = append(Execs, ExecsMap)
						ReadinessProbesMap["exec"] = Execs
					}

					HttpGets := make([]map[string]interface{}, 0)
					if HttpGetList, ok := ReadinessProbes["HttpGet"]; ok && len(HttpGetList.(map[string]interface{})) > 0 {
						HttpGetsMap := map[string]interface{}{
							"scheme": HttpGetList.(map[string]interface{})["Scheme"],
							"port":   HttpGetList.(map[string]interface{})["Port"],
							"path":   HttpGetList.(map[string]interface{})["Path"],
						}
						HttpGets = append(HttpGets, HttpGetsMap)
						ReadinessProbesMap["http_get"] = HttpGets
					}

					TcpSockets := make([]map[string]interface{}, 0)
					if TcpSocketList, ok := ReadinessProbes["TcpSocket"]; ok && len(TcpSocketList.(map[string]interface{})) > 0 {
						TcpSocketsMap := map[string]interface{}{
							"port": TcpSocketList.(map[string]interface{})["Port"],
						}
						TcpSockets = append(TcpSockets, TcpSocketsMap)
						ReadinessProbesMap["tcp_socket"] = TcpSockets
					}

					ReadinessProbesMaps = append(ReadinessProbesMaps, ReadinessProbesMap)
					temp1["readiness_probe"] = ReadinessProbesMaps
				}
				if m1["LivenessProbe"] != nil {
					LivenessProbesMaps := make([]map[string]interface{}, 0)
					LivenessProbes := m1["LivenessProbe"].(map[string]interface{})
					LivenessProbesMap := map[string]interface{}{
						"initial_delay_seconds": LivenessProbes["InitialDelaySeconds"],
						"timeout_seconds":       LivenessProbes["TimeoutSeconds"],
						"period_seconds":        LivenessProbes["PeriodSeconds"],
						"success_threshold":     LivenessProbes["SuccessThreshold"],
						"failure_threshold":     LivenessProbes["FailureThreshold"],
					}
					Execs := make([]map[string]interface{}, 0)
					if ExecsList, ok := LivenessProbes["Execs"]; ok && ExecsList != nil {
						ExecsMap := map[string]interface{}{
							"commands": ExecsList.([]interface{}),
						}
						Execs = append(Execs, ExecsMap)
						LivenessProbesMap["exec"] = Execs
					}

					HttpGets := make([]map[string]interface{}, 0)
					if HttpGetList, ok := LivenessProbes["HttpGet"]; ok && len(HttpGetList.(map[string]interface{})) > 0 {
						HttpGetsMap := map[string]interface{}{
							"scheme": HttpGetList.(map[string]interface{})["Scheme"],
							"port":   HttpGetList.(map[string]interface{})["Port"],
							"path":   HttpGetList.(map[string]interface{})["Path"],
						}
						HttpGets = append(HttpGets, HttpGetsMap)
						LivenessProbesMap["http_get"] = HttpGets
					}

					TcpSockets := make([]map[string]interface{}, 0)
					if TcpSocketList, ok := LivenessProbes["TcpSocket"]; ok && len(TcpSocketList.(map[string]interface{})) > 0 {
						TcpSocketsMap := map[string]interface{}{
							"port": TcpSocketList.(map[string]interface{})["Port"],
						}
						TcpSockets = append(TcpSockets, TcpSocketsMap)
						LivenessProbesMap["tcp_socket"] = TcpSockets
					}
					LivenessProbesMaps = append(LivenessProbesMaps, LivenessProbesMap)
					temp1["liveness_probe"] = LivenessProbesMaps
				}

				if m1["SecurityContext"] != nil {
					SecurityContextMaps := make([]map[string]interface{}, 0)
					SecurityContextValue := m1["SecurityContext"].(map[string]interface{})
					SecurityContextMap := map[string]interface{}{
						"run_as_user": SecurityContextValue["RunAsUser"],
						"privileged":  getSecurityContextPrivileged(d, temp1["name"]),
					}
					if SecurityContextValue["Capability"] != nil && SecurityContextValue["Capability"].(map[string]interface{})["Adds"] != nil {
						Capabilities := make([]map[string]interface{}, 0)
						Capability := map[string]interface{}{
							"add": SecurityContextValue["Capability"].(map[string]interface{})["Adds"],
						}
						Capabilities = append(Capabilities, Capability)
						SecurityContextMap["capability"] = Capabilities
					}
					SecurityContextMaps = append(SecurityContextMaps, SecurityContextMap)
					temp1["security_context"] = SecurityContextMaps
				}

				exec := getLifecyclePreStopHandlerExec(d, temp1["name"])
				if exec != "" {
					temp1["lifecycle_pre_stop_handler_exec"] = exec
				}

				containers = append(containers, temp1)
			}
		}
	}
	if err := d.Set("containers", containers); err != nil {
		return WrapError(err)
	}
	d.Set("cpu", object["Cpu"])

	dnsConfigSli := make([]map[string]interface{}, 0)
	if len(object["DnsConfig"].(map[string]interface{})) > 0 {
		dnsConfig := object["DnsConfig"]
		dnsConfigMap := make(map[string]interface{})
		dnsConfigMap["name_servers"] = dnsConfig.(map[string]interface{})["NameServers"]

		optionsSli := make([]map[string]interface{}, 0)
		if dnsConfig.(map[string]interface{})["Options"] != nil && len(dnsConfig.(map[string]interface{})["Options"].([]interface{})) > 0 {
			for _, options := range dnsConfig.(map[string]interface{})["Options"].([]interface{}) {
				optionsMap := make(map[string]interface{})
				optionsMap["name"] = options.(map[string]interface{})["Name"]
				optionsMap["value"] = options.(map[string]interface{})["Value"]
				optionsSli = append(optionsSli, optionsMap)
			}
		}
		dnsConfigMap["options"] = optionsSli
		dnsConfigMap["searches"] = dnsConfig.(map[string]interface{})["Searches"]
		dnsConfigSli = append(dnsConfigSli, dnsConfigMap)
	}
	d.Set("dns_config", dnsConfigSli)

	securityContextSli := make([]map[string]interface{}, 0)
	if len(object["EciSecurityContext"].(map[string]interface{})) > 0 {
		eciSecurityContext := object["EciSecurityContext"]
		eciSecurityContextMap := make(map[string]interface{})

		sysctlsSli := make([]map[string]interface{}, 0)
		if len(eciSecurityContext.(map[string]interface{})["Sysctls"].([]interface{})) > 0 {
			for _, sysctls := range eciSecurityContext.(map[string]interface{})["Sysctls"].([]interface{}) {
				sysctlsMap := make(map[string]interface{})
				sysctlsMap["name"] = sysctls.(map[string]interface{})["Name"]
				sysctlsMap["value"] = sysctls.(map[string]interface{})["Value"]
				sysctlsSli = append(sysctlsSli, sysctlsMap)
			}
		}
		eciSecurityContextMap["sysctl"] = sysctlsSli
		securityContextSli = append(securityContextSli, eciSecurityContextMap)
	}
	d.Set("security_context", securityContextSli)

	hostAliases := make([]map[string]interface{}, 0)
	if hostAliasesList, ok := object["HostAliases"].([]interface{}); ok {
		for _, v := range hostAliasesList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"hostnames": m1["Hostnames"],
					"ip":        m1["Ip"],
				}
				hostAliases = append(hostAliases, temp1)

			}
		}
	}
	if err := d.Set("host_aliases", hostAliases); err != nil {
		return WrapError(err)
	}

	initContainers := make([]map[string]interface{}, 0)
	if initContainersList, ok := object["InitContainers"].([]interface{}); ok {
		for _, v := range initContainersList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"args":              m1["Args"],
					"commands":          m1["Command"],
					"cpu":               m1["Cpu"],
					"gpu":               m1["Gpu"],
					"image":             m1["Image"],
					"image_pull_policy": m1["ImagePullPolicy"],
					"memory":            m1["Memory"],
					"name":              m1["Name"],
					"ready":             m1["Ready"],
					"restart_count":     m1["RestartCount"],
					"working_dir":       m1["WorkingDir"],
				}
				if m1["EnvironmentVars"] != nil {
					environmentVarsMaps := make([]map[string]interface{}, 0)
					for _, environmentVarsValue := range m1["EnvironmentVars"].([]interface{}) {
						environmentVars := environmentVarsValue.(map[string]interface{})
						environmentVarsMap := map[string]interface{}{
							"key":   environmentVars["Key"],
							"value": environmentVars["Value"],
						}
						if environmentVars["ValueFrom"] != nil {
							fieldRefMaps := make([]map[string]interface{}, 0)
							fieldRefValue := environmentVars["ValueFrom"].(map[string]interface{})["FieldRef"]
							if fieldRefValue != nil && fieldRefValue.(map[string]interface{})["FieldPath"] != nil {
								fieldRefMap := map[string]interface{}{
									"field_path": fieldRefValue.(map[string]interface{})["FieldPath"],
								}
								fieldRefMaps = append(fieldRefMaps, fieldRefMap)
								environmentVarsMap["field_ref"] = fieldRefMaps
							}
						}
						environmentVarsMaps = append(environmentVarsMaps, environmentVarsMap)
					}
					temp1["environment_vars"] = environmentVarsMaps
				}
				if m1["Ports"] != nil {
					portsMaps := make([]map[string]interface{}, 0)
					for _, portsValue := range m1["Ports"].([]interface{}) {
						ports := portsValue.(map[string]interface{})
						portsMap := map[string]interface{}{
							"port":     ports["Port"],
							"protocol": ports["Protocol"],
						}
						portsMaps = append(portsMaps, portsMap)
					}
					temp1["ports"] = portsMaps
				}
				if m1["VolumeMounts"] != nil {
					volumeMountsMaps := make([]map[string]interface{}, 0)
					for _, volumeMountsValue := range m1["VolumeMounts"].([]interface{}) {
						volumeMounts := volumeMountsValue.(map[string]interface{})
						volumeMountsMap := map[string]interface{}{
							"mount_path": volumeMounts["MountPath"],
							"name":       volumeMounts["Name"],
							"read_only":  volumeMounts["ReadOnly"],
						}
						volumeMountsMaps = append(volumeMountsMaps, volumeMountsMap)
					}
					temp1["volume_mounts"] = volumeMountsMaps
				}
				if m1["SecurityContext"] != nil {
					SecurityContextMaps := make([]map[string]interface{}, 0)
					SecurityContextValue := m1["SecurityContext"].(map[string]interface{})
					SecurityContextMap := map[string]interface{}{
						"run_as_user": SecurityContextValue["RunAsUser"],
					}
					if SecurityContextValue["Capability"] != nil && SecurityContextValue["Capability"].(map[string]interface{})["Adds"] != nil {
						Capabilities := make([]map[string]interface{}, 0)
						Capability := map[string]interface{}{
							"add": SecurityContextValue["Capability"].(map[string]interface{})["Adds"],
						}
						Capabilities = append(Capabilities, Capability)
						SecurityContextMap["capability"] = Capabilities
					}
					SecurityContextMaps = append(SecurityContextMaps, SecurityContextMap)
					temp1["security_context"] = SecurityContextMaps
				}
				initContainers = append(initContainers, temp1)

			}
		}
	}
	if err := d.Set("init_containers", initContainers); err != nil {
		return WrapError(err)
	}
	d.Set("instance_type", object["InstanceType"])
	d.Set("memory", object["Memory"])
	d.Set("ram_role_name", object["RamRoleName"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("restart_policy", object["RestartPolicy"])
	d.Set("security_group_id", object["SecurityGroupId"])
	d.Set("status", object["Status"])
	d.Set("tags", tagsToMap(object["Tags"]))
	d.Set("vswitch_id", object["VSwitchId"])

	volumes := make([]map[string]interface{}, 0)
	if volumesList, ok := object["Volumes"].([]interface{}); ok {
		for _, v := range volumesList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"disk_volume_disk_id":  m1["DiskVolumeDiskId"],
					"disk_volume_fs_type":  m1["DiskVolumeFsType"],
					"flex_volume_driver":   m1["FlexVolumeDriver"],
					"flex_volume_fs_type":  m1["FlexVolumeFsType"],
					"flex_volume_options":  m1["FlexVolumeOptions"],
					"nfs_volume_path":      m1["NFSVolumePath"],
					"nfs_volume_read_only": m1["NFSVolumeReadOnly"],
					"nfs_volume_server":    m1["NFSVolumeServer"],
					"name":                 m1["Name"],
					"type":                 m1["Type"],
				}
				if m1["ConfigFileVolumeConfigFileToPaths"] != nil {
					configFileVolumeConfigFileToPathsMaps := make([]map[string]interface{}, 0)
					for _, configFileVolumeConfigFileToPathsValue := range m1["ConfigFileVolumeConfigFileToPaths"].([]interface{}) {
						configFileVolumeConfigFileToPaths := configFileVolumeConfigFileToPathsValue.(map[string]interface{})
						configFileVolumeConfigFileToPathsMap := map[string]interface{}{
							"path": configFileVolumeConfigFileToPaths["Path"],
						}
						content := getConfigFileContent(d, temp1["name"], configFileVolumeConfigFileToPathsMap["path"])
						configFileVolumeConfigFileToPathsMap["content"] = content
						configFileVolumeConfigFileToPathsMaps = append(configFileVolumeConfigFileToPathsMaps, configFileVolumeConfigFileToPathsMap)
					}
					temp1["config_file_volume_config_file_to_paths"] = configFileVolumeConfigFileToPathsMaps
				}
				volumes = append(volumes, temp1)

			}
		}
	}
	if err := d.Set("volumes", volumes); err != nil {
		return WrapError(err)
	}
	d.Set("zone_id", object["ZoneId"])
	d.Set("spot_strategy", object["SpotStrategy"])
	d.Set("spot_price_limit", object["SpotPriceLimit"])
	d.Set("dns_policy", object["DnsPolicy"])
	return nil
}

func getLifecyclePreStopHandlerExec(d *schema.ResourceData, name interface{}) (result interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("getLifecyclePreStopHandlerExec Recovered from panic: %v", r)
			result = ""
		}
	}()
	for _, srcContainer := range d.Get("containers").([]interface{}) {
		c := srcContainer.(map[string]interface{})
		if c["name"].(string) == name.(string) {
			return c["lifecycle_pre_stop_handler_exec"]
		}
	}
	return ""
}

func getSecurityContextPrivileged(d *schema.ResourceData, name interface{}) (result interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("getSecurityContextPrivileged Recovered from panic: %v", r)
			result = nil
		}
	}()
	for _, srcContainer := range d.Get("containers").([]interface{}) {
		c := srcContainer.(map[string]interface{})
		if c["name"].(string) == name.(string) {
			if c["security_context"] != nil {
				scs := c["security_context"].([]interface{})
				for _, tscs := range scs {
					mapData, ok := tscs.(map[string]interface{})
					if ok {
						return mapData["privileged"]
					}
				}
			}
		}
	}
	return nil
}

func getConfigFileContent(d *schema.ResourceData, name interface{}, path interface{}) (result interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("getConfigFileContent Recovered from panic: %v", r)
			result = ""
		}
	}()
	for _, v := range d.Get("volumes").([]interface{}) {
		volume := v.(map[string]interface{})
		if volume["name"].(string) == name.(string) {
			configs := volume["config_file_volume_config_file_to_paths"].([]interface{})
			for _, c := range configs {
				config := c.(map[string]interface{})
				if config["path"].(string) == path.(string) {
					return config["content"]
				}
			}
		}
	}
	return ""
}

func resourceAlicloudEciContainerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eciService := EciService{client}
	var response map[string]interface{}
	var err error
	update := false
	request := map[string]interface{}{
		"ContainerGroupId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("cpu") {
		update = true
		request["Cpu"] = d.Get("cpu")
	}
	if d.HasChange("dns_config") {
		update = true
		if d.Get("dns_config") != nil {
			dnsConfigMap := make(map[string]interface{})
			for _, dnsConfig := range d.Get("dns_config").(*schema.Set).List() {
				dnsConfigArg := dnsConfig.(map[string]interface{})
				dnsConfigMap["NameServer"] = dnsConfigArg["name_servers"]
				if dnsConfigArg["options"] != nil {
					optionsMaps := make([]map[string]interface{}, 0)
					for _, options := range dnsConfigArg["options"].([]interface{}) {
						optionsMap := make(map[string]interface{})
						optionsArg := options.(map[string]interface{})
						optionsMap["Name"] = optionsArg["name"]
						optionsMap["Value"] = optionsArg["value"]
						optionsMaps = append(optionsMaps, optionsMap)
					}
					dnsConfigMap["Option"] = optionsMaps
				}
				dnsConfigMap["Search"] = dnsConfigArg["searches"]
			}
			request["DnsConfig"] = dnsConfigMap
		}
	}

	if d.HasChange("memory") {
		update = true
		request["Memory"] = d.Get("memory")
	}
	if d.HasChange("restart_policy") {
		update = true
		request["RestartPolicy"] = d.Get("restart_policy")
	}
	if d.HasChange("tags") {
		update = true
		count := 1
		for key, value := range d.Get("tags").(map[string]interface{}) {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
	}
	if d.HasChange("volumes") {
		update = true
		Volumes := make([]map[string]interface{}, len(d.Get("volumes").([]interface{})))
		for i, VolumesValue := range d.Get("volumes").([]interface{}) {
			VolumesMap := VolumesValue.(map[string]interface{})
			Volumes[i] = make(map[string]interface{})
			ConfigFileVolumeConfigFileToPaths := make([]map[string]interface{}, len(VolumesMap["config_file_volume_config_file_to_paths"].([]interface{})))
			for i, ConfigFileVolumeConfigFileToPathsValue := range VolumesMap["config_file_volume_config_file_to_paths"].([]interface{}) {
				ConfigFileVolumeConfigFileToPathsMap := ConfigFileVolumeConfigFileToPathsValue.(map[string]interface{})
				ConfigFileVolumeConfigFileToPaths[i] = make(map[string]interface{})
				ConfigFileVolumeConfigFileToPaths[i]["Content"] = ConfigFileVolumeConfigFileToPathsMap["content"]
				ConfigFileVolumeConfigFileToPaths[i]["Path"] = ConfigFileVolumeConfigFileToPathsMap["path"]
			}
			Volumes[i]["ConfigFileVolume.ConfigFileToPath"] = ConfigFileVolumeConfigFileToPaths

			Volumes[i]["NFSVolume.Path"] = VolumesMap["nfs_volume_path"]
			Volumes[i]["NFSVolume.Server"] = VolumesMap["nfs_volume_server"]
			Volumes[i]["Name"] = VolumesMap["name"]
			Volumes[i]["Type"] = VolumesMap["type"]
		}
		request["Volume"] = Volumes

	}
	if d.HasChange("image_registry_credential") {
		update = true
		if v, ok := d.GetOk("image_registry_credential"); ok {
			imageRegisryCredentialMaps := make([]map[string]interface{}, 0)
			for _, raw := range v.(*schema.Set).List() {
				obj := raw.(map[string]interface{})
				imageRegisryCredentialMaps = append(imageRegisryCredentialMaps, map[string]interface{}{
					"Password": obj["password"],
					"Server":   obj["server"],
					"UserName": obj["user_name"],
				})
			}
			request["ImageRegistryCredential"] = imageRegisryCredentialMaps
		}
	}

	if d.HasChange("resource_group_id") {
		update = true
		if v, ok := d.GetOk("resource_group_id"); ok {
			request["ResourceGroupId"] = v
		}
	}
	if d.HasChange("containers") || d.HasChange("init_containers") {
		update = true
	}

	request["UpdateType"] = "IncrementalUpdate"

	if update {

		Containers := make([]map[string]interface{}, len(d.Get("containers").([]interface{})))
		for i, ContainersValue := range d.Get("containers").([]interface{}) {
			ContainersMap := ContainersValue.(map[string]interface{})
			Containers[i] = make(map[string]interface{})
			Containers[i]["Arg"] = ContainersMap["args"]
			Containers[i]["Command"] = ContainersMap["commands"]
			Containers[i]["Cpu"] = ContainersMap["cpu"]
			EnvironmentVars := make([]map[string]interface{}, len(ContainersMap["environment_vars"].([]interface{})))
			for i, EnvironmentVarsValue := range ContainersMap["environment_vars"].([]interface{}) {
				EnvironmentVarsMap := EnvironmentVarsValue.(map[string]interface{})
				EnvironmentVars[i] = make(map[string]interface{})
				EnvironmentVars[i]["Key"] = EnvironmentVarsMap["key"]
				EnvironmentVars[i]["Value"] = EnvironmentVarsMap["value"]
				for _, FieldRefValue := range EnvironmentVarsMap["field_ref"].([]interface{}) {
					FieldRef := map[string]interface{}{}
					FieldRefMap := FieldRefValue.(map[string]interface{})
					FieldRef["FieldPath"] = FieldRefMap["field_path"]
					EnvironmentVars[i]["FieldRef"] = FieldRef
				}
			}
			Containers[i]["EnvironmentVar"] = EnvironmentVars

			Containers[i]["Gpu"] = ContainersMap["gpu"]
			Containers[i]["Image"] = ContainersMap["image"]
			Containers[i]["ImagePullPolicy"] = ContainersMap["image_pull_policy"]
			Containers[i]["Memory"] = ContainersMap["memory"]
			Containers[i]["Name"] = ContainersMap["name"]
			Ports := make([]map[string]interface{}, len(ContainersMap["ports"].([]interface{})))
			for i, PortsValue := range ContainersMap["ports"].([]interface{}) {
				PortsMap := PortsValue.(map[string]interface{})
				Ports[i] = make(map[string]interface{})
				Ports[i]["Port"] = PortsMap["port"]
				Ports[i]["Protocol"] = PortsMap["protocol"]
			}
			Containers[i]["Port"] = Ports

			VolumeMounts := make([]map[string]interface{}, len(ContainersMap["volume_mounts"].([]interface{})))
			for i, VolumeMountsValue := range ContainersMap["volume_mounts"].([]interface{}) {
				VolumeMountsMap := VolumeMountsValue.(map[string]interface{})
				VolumeMounts[i] = make(map[string]interface{})
				VolumeMounts[i]["MountPath"] = VolumeMountsMap["mount_path"]
				VolumeMounts[i]["Name"] = VolumeMountsMap["name"]
				VolumeMounts[i]["ReadOnly"] = VolumeMountsMap["read_only"]
			}
			Containers[i]["VolumeMount"] = VolumeMounts

			Containers[i]["WorkingDir"] = ContainersMap["working_dir"]

			LivenessProbe := map[string]interface{}{}
			for _, LivenessProbesValue := range ContainersMap["liveness_probe"].([]interface{}) {

				LivenessProbesMap := LivenessProbesValue.(map[string]interface{})
				HttpGetValue := map[string]interface{}{}
				for _, HttpGetValues := range LivenessProbesMap["http_get"].([]interface{}) {
					HttpGetValueMap := HttpGetValues.(map[string]interface{})
					HttpGetValue["Scheme"] = HttpGetValueMap["scheme"]
					HttpGetValue["Port"] = HttpGetValueMap["port"]
					HttpGetValue["Path"] = HttpGetValueMap["path"]

				}
				ExecsValue := map[string]interface{}{}
				for _, ExecsValues := range LivenessProbesMap["exec"].([]interface{}) {
					ExecsValuesValueMap := ExecsValues.(map[string]interface{})
					ExecsValue["Command"] = ExecsValuesValueMap["commands"]
				}
				TcpSocketValue := map[string]interface{}{}
				for _, TcpSocketsValues := range LivenessProbesMap["tcp_socket"].([]interface{}) {
					TcpSocketsValueMap := TcpSocketsValues.(map[string]interface{})
					TcpSocketValue["Port"] = TcpSocketsValueMap["port"]
				}

				LivenessProbe["PeriodSeconds"] = LivenessProbesMap["period_seconds"]
				LivenessProbe["InitialDelaySeconds"] = LivenessProbesMap["initial_delay_seconds"]
				LivenessProbe["SuccessThreshold"] = LivenessProbesMap["success_threshold"]
				LivenessProbe["FailureThreshold"] = LivenessProbesMap["failure_threshold"]
				LivenessProbe["TimeoutSeconds"] = LivenessProbesMap["timeout_seconds"]
				LivenessProbe["HttpGet"] = HttpGetValue
				LivenessProbe["Exec"] = ExecsValue
				LivenessProbe["TcpSocket"] = TcpSocketValue
			}
			Containers[i]["LivenessProbe"] = LivenessProbe

			ReadinessProbe := map[string]interface{}{}
			for _, ReadinessProbesValue := range ContainersMap["readiness_probe"].([]interface{}) {

				ReadinessProbesMap := ReadinessProbesValue.(map[string]interface{})
				HttpGetValue := map[string]interface{}{}
				for _, HttpGetValues := range ReadinessProbesMap["http_get"].([]interface{}) {
					HttpGetValueMap := HttpGetValues.(map[string]interface{})
					HttpGetValue["Scheme"] = HttpGetValueMap["scheme"]
					HttpGetValue["Port"] = HttpGetValueMap["port"]
					HttpGetValue["Path"] = HttpGetValueMap["path"]

				}
				ExecsValue := map[string]interface{}{}
				for _, ExecsValues := range ReadinessProbesMap["exec"].([]interface{}) {
					ExecsValuesValueMap := ExecsValues.(map[string]interface{})
					ExecsValue["Command"] = ExecsValuesValueMap["commands"]
				}
				TcpSocketValue := map[string]interface{}{}
				for _, TcpSocketsValues := range ReadinessProbesMap["tcp_socket"].([]interface{}) {
					TcpSocketsValueMap := TcpSocketsValues.(map[string]interface{})
					TcpSocketValue["Port"] = TcpSocketsValueMap["port"]
				}

				ReadinessProbe["PeriodSeconds"] = ReadinessProbesMap["period_seconds"]
				ReadinessProbe["InitialDelaySeconds"] = ReadinessProbesMap["initial_delay_seconds"]
				ReadinessProbe["SuccessThreshold"] = ReadinessProbesMap["success_threshold"]
				ReadinessProbe["FailureThreshold"] = ReadinessProbesMap["failure_threshold"]
				ReadinessProbe["TimeoutSeconds"] = ReadinessProbesMap["timeout_seconds"]
				ReadinessProbe["HttpGet"] = HttpGetValue
				ReadinessProbe["Exec"] = ExecsValue
				ReadinessProbe["TcpSocket"] = TcpSocketValue
			}
			Containers[i]["ReadinessProbe"] = ReadinessProbe

			SecurityContext := map[string]interface{}{}
			for _, SecurityContextValue := range ContainersMap["security_context"].([]interface{}) {
				SecurityContextMap := SecurityContextValue.(map[string]interface{})
				if SecurityContextMap["capability"] != nil {
					for _, v := range SecurityContextMap["capability"].([]interface{}) {
						CapabilityValue := map[string]interface{}{}
						CapabilityValue["Add"] = v.(map[string]interface{})["add"]
						SecurityContext["Capability"] = CapabilityValue
					}
				}
				SecurityContext["RunAsUser"] = SecurityContextMap["run_as_user"]
			}
			Containers[i]["SecurityContext"] = SecurityContext

			Containers[i]["LifecyclePreStopHandlerExec"] = ContainersMap["lifecycle_pre_stop_handler_exec"]
		}
		request["Container"] = Containers

		InitContainers := make([]map[string]interface{}, len(d.Get("init_containers").([]interface{})))
		for i, InitContainersValue := range d.Get("init_containers").([]interface{}) {
			InitContainersMap := InitContainersValue.(map[string]interface{})
			InitContainers[i] = make(map[string]interface{})
			InitContainers[i]["Arg"] = InitContainersMap["args"]
			InitContainers[i]["Command"] = InitContainersMap["commands"]
			InitContainers[i]["Cpu"] = InitContainersMap["cpu"]
			EnvironmentVars := make([]map[string]interface{}, len(InitContainersMap["environment_vars"].([]interface{})))
			for i, EnvironmentVarsValue := range InitContainersMap["environment_vars"].([]interface{}) {
				EnvironmentVarsMap := EnvironmentVarsValue.(map[string]interface{})
				EnvironmentVars[i] = make(map[string]interface{})
				EnvironmentVars[i]["Key"] = EnvironmentVarsMap["key"]
				EnvironmentVars[i]["Value"] = EnvironmentVarsMap["value"]
				for _, FieldRefValue := range EnvironmentVarsMap["field_ref"].([]interface{}) {
					FieldRef := map[string]interface{}{}
					FieldRefMap := FieldRefValue.(map[string]interface{})
					FieldRef["FieldPath"] = FieldRefMap["field_path"]
					EnvironmentVars[i]["FieldRef"] = FieldRef
				}
			}
			InitContainers[i]["EnvironmentVar"] = EnvironmentVars

			InitContainers[i]["Gpu"] = InitContainersMap["gpu"]
			InitContainers[i]["Image"] = InitContainersMap["image"]
			InitContainers[i]["ImagePullPolicy"] = InitContainersMap["image_pull_policy"]
			InitContainers[i]["Memory"] = InitContainersMap["memory"]
			InitContainers[i]["Name"] = InitContainersMap["name"]
			Ports := make([]map[string]interface{}, len(InitContainersMap["ports"].([]interface{})))
			for i, PortsValue := range InitContainersMap["ports"].([]interface{}) {
				PortsMap := PortsValue.(map[string]interface{})
				Ports[i] = make(map[string]interface{})
				Ports[i]["Port"] = PortsMap["port"]
				Ports[i]["Protocol"] = PortsMap["protocol"]
			}
			InitContainers[i]["Port"] = Ports

			VolumeMounts := make([]map[string]interface{}, len(InitContainersMap["volume_mounts"].([]interface{})))
			for i, VolumeMountsValue := range InitContainersMap["volume_mounts"].([]interface{}) {
				VolumeMountsMap := VolumeMountsValue.(map[string]interface{})
				VolumeMounts[i] = make(map[string]interface{})
				VolumeMounts[i]["MountPath"] = VolumeMountsMap["mount_path"]
				VolumeMounts[i]["Name"] = VolumeMountsMap["name"]
				VolumeMounts[i]["ReadOnly"] = VolumeMountsMap["read_only"]
			}
			InitContainers[i]["VolumeMount"] = VolumeMounts

			InitContainers[i]["WorkingDir"] = InitContainersMap["working_dir"]

			SecurityContext := map[string]interface{}{}
			for _, SecurityContextValue := range InitContainersMap["security_context"].([]interface{}) {
				SecurityContextMap := SecurityContextValue.(map[string]interface{})
				CapabilityValue := map[string]interface{}{}
				for _, CapabilityValues := range SecurityContextMap["capability"].([]interface{}) {
					CapabilityValueMap := CapabilityValues.(map[string]interface{})
					CapabilityValue["Add"] = CapabilityValueMap["add"]
				}
				SecurityContext["Capability"] = CapabilityValue
				SecurityContext["RunAsUser"] = SecurityContextMap["run_as_user"]
			}
			InitContainers[i]["SecurityContext"] = SecurityContext
		}
		request["InitContainer"] = InitContainers

		action := "UpdateContainerGroup"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Eci", "2018-08-08", action, nil, request, true)
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

		stateConf := BuildStateConf([]string{}, []string{"Running", "Succeeded"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, eciService.EciContainerGroupStateRefreshFunc(d.Id(), []string{"Failed", "ScheduleFailed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAlicloudEciContainerGroupRead(d, meta)
}

func resourceAlicloudEciContainerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteContainerGroup"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"ContainerGroupId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Eci", "2018-08-08", action, nil, request, true)
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

	return nil
}
