package alicloud

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEssEciScalingConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssEciScalingConfigurationCreate,
		Read:   resourceAliyunEssEciScalingConfigurationRead,
		Update: resourceAliyunEssEciScalingConfigurationUpdate,
		Delete: resourceAliyunEssEciScalingConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"force_delete": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"scaling_group_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"scaling_configuration_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringMatch(regexp.MustCompile(`^[\\u4e00-\\u9fa5a-zA-Z0-9][\\u4e00-\\u9fa5a-zA-Z0-9\-_.]{1,63}$`), "It must be 2 to 64 characters in length and can contain letters, digits, underscores (_), hyphens (-), and periods (.). It must start with a letter or a digit."),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"container_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"restart_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cpu": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"memory": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dns_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cost_optimization": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_sls": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"instance_family_level": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ram_role_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"termination_grace_period_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"auto_match_image_cache": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ipv6_address_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"cpu_options_core": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntAtLeast(1),
			},
			"cpu_options_threads_per_core": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntAtLeast(1),
			},
			"active_deadline_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"spot_strategy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"spot_price_limit": {
				Type:     schema.TypeFloat,
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
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ingress_bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"egress_bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ephemeral_storage": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"load_balancer_weight": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"tags": tagsSchema(),
			"acr_registry_infos": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domains": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
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
						},
					},
				},
			},
			"image_registry_credentials": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"password": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"server": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"username": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"dns_config_options": {
				Type:     schema.TypeSet,
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
			"security_context_sysctls": {
				Type:     schema.TypeSet,
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
			"containers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_context_capability_adds": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
						"lifecycle_pre_stop_handler_execs": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
						"security_context_read_only_root_file_system": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"tty": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"stdin": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"security_context_run_as_user": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"ports": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"environment_vars": {
							Type:     schema.TypeSet,
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
									"field_ref_field_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"working_dir": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"args": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cpu": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"gpu": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"memory": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"image": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"image_pull_policy": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"volume_mounts": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mount_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"mount_propagation": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"sub_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"read_only": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"commands": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"liveness_probe_exec_commands": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
						"liveness_probe_period_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"liveness_probe_http_get_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"liveness_probe_failure_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"liveness_probe_initial_delay_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"liveness_probe_http_get_port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"liveness_probe_http_get_scheme": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS"}, false),
						},
						"liveness_probe_tcp_socket_port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"liveness_probe_success_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"liveness_probe_timeout_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"readiness_probe_exec_commands": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
						"readiness_probe_period_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"readiness_probe_http_get_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"readiness_probe_failure_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"readiness_probe_initial_delay_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"readiness_probe_http_get_port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"readiness_probe_http_get_scheme": {
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS"}, false),
							Optional:     true,
						},
						"readiness_probe_tcp_socket_port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"readiness_probe_success_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"readiness_probe_timeout_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"instance_types": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					MaxItems: 5,
				},
			},
			"init_containers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_context_capability_adds": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
						"security_context_read_only_root_file_system": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"security_context_run_as_user": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"ports": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"environment_vars": {
							Type:     schema.TypeSet,
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
									"field_ref_field_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"working_dir": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"args": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cpu": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"gpu": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"memory": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"image": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"image_pull_policy": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"volume_mounts": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mount_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"mount_propagation": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"sub_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"read_only": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
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
			"volumes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_file_volume_config_file_to_paths": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"content": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"mode": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"disk_volume_disk_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"host_path_volume_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"host_path_volume_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"config_file_volume_default_mode": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"empty_dir_volume_medium": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"empty_dir_volume_size_limit": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"disk_volume_fs_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"disk_volume_disk_size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"flex_volume_driver": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"flex_volume_fs_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"flex_volume_options": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"nfs_volume_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"nfs_volume_read_only": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"nfs_volume_server": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"host_aliases": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostnames": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliyunEssEciScalingConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	// Ensure instance_type is generation three
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateEciScalingConfiguration"
	request := make(map[string]interface{})
	var err error
	request["ScalingGroupId"] = d.Get("scaling_group_id")
	request["ScalingConfigurationName"] = d.Get("scaling_configuration_name")
	request["Description"] = d.Get("description")
	request["SecurityGroupId"] = d.Get("security_group_id")
	request["ContainerGroupName"] = d.Get("container_group_name")
	request["RestartPolicy"] = d.Get("restart_policy")
	request["Cpu"] = d.Get("cpu")
	request["Memory"] = d.Get("memory")
	request["ResourceGroupId"] = d.Get("resource_group_id")
	request["DnsPolicy"] = d.Get("dns_policy")
	request["EnableSls"] = d.Get("enable_sls")
	request["RamRoleName"] = d.Get("ram_role_name")
	request["AutoCreateEip"] = d.Get("auto_create_eip")
	request["EipBandwidth"] = d.Get("eip_bandwidth")
	request["HostName"] = d.Get("host_name")
	request["IngressBandwidth"] = d.Get("ingress_bandwidth")
	request["EgressBandwidth"] = d.Get("egress_bandwidth")
	request["SpotStrategy"] = d.Get("spot_strategy")
	request["ImageSnapshotId"] = d.Get("image_snapshot_id")
	request["TerminationGracePeriodSeconds"] = d.Get("termination_grace_period_seconds")
	request["AutoMatchImageCache"] = d.Get("auto_match_image_cache")
	request["Ipv6AddressCount"] = d.Get("ipv6_address_count")
	if v, ok := d.GetOkExists("cpu_options_core"); ok {
		request["CpuOptionsCore"] = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("instance_family_level"); ok {
		request["InstanceFamilyLevel"] = v
	}
	if v, ok := d.GetOkExists("cpu_options_threads_per_core"); ok {
		request["CpuOptionsThreadsPerCore"] = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOkExists("cost_optimization"); ok {
		request["CostOptimization"] = v.(bool)
	}
	request["EphemeralStorage"] = d.Get("ephemeral_storage")
	request["LoadBalancerWeight"] = d.Get("load_balancer_weight")

	if v, ok := d.GetOk("active_deadline_seconds"); ok {
		request["ActiveDeadlineSeconds"] = v
	}
	if v, ok := d.GetOk("spot_price_limit"); ok {
		request["SpotPriceLimit"] = strconv.FormatFloat(v.(float64), 'f', 2, 64)
	}

	if v, ok := d.GetOk("tags"); ok {
		count := 1
		for key, value := range v.(map[string]interface{}) {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
	}

	if v, ok := d.GetOk("image_registry_credentials"); ok {
		imageRegisryCredentialMaps := make([]map[string]interface{}, 0)
		for _, raw := range v.(*schema.Set).List() {
			obj := raw.(map[string]interface{})
			imageRegisryCredentialMaps = append(imageRegisryCredentialMaps, map[string]interface{}{
				"Password": obj["password"],
				"Server":   obj["server"],
				"UserName": obj["username"],
			})
		}
		request["ImageRegistryCredential"] = imageRegisryCredentialMaps
	}

	if v, ok := d.GetOk("dns_config_options"); ok {
		dnsConfigOptionMaps := make([]map[string]interface{}, 0)
		for _, raw := range v.(*schema.Set).List() {
			obj := raw.(map[string]interface{})
			dnsConfigOptionMaps = append(dnsConfigOptionMaps, map[string]interface{}{
				"Name":  obj["name"],
				"Value": obj["value"],
			})
		}
		request["DnsConfigOption"] = dnsConfigOptionMaps
	}

	if v, ok := d.GetOk("security_context_sysctls"); ok {
		securityContextSysctlMaps := make([]map[string]interface{}, 0)
		for _, raw := range v.(*schema.Set).List() {
			obj := raw.(map[string]interface{})
			securityContextSysctlMaps = append(securityContextSysctlMaps, map[string]interface{}{
				"Name":  obj["name"],
				"Value": obj["value"],
			})
		}
		request["SecurityContextSysctl"] = securityContextSysctlMaps
	}

	if v, ok := d.GetOk("acr_registry_infos"); ok {
		acrRegistryInfoMaps := make([]map[string]interface{}, 0)
		for _, raw := range v.(*schema.Set).List() {
			obj := raw.(map[string]interface{})

			acrRegistryInfoMaps = append(acrRegistryInfoMaps, map[string]interface{}{
				"Domain":       expandStringList(obj["domains"].(*schema.Set).List()),
				"InstanceName": obj["instance_name"],
				"InstanceId":   obj["instance_id"],
				"RegionId":     obj["region_id"],
			})
		}
		request["AcrRegistryInfo"] = acrRegistryInfoMaps
	}

	Containers := make([]map[string]interface{}, len(d.Get("containers").([]interface{})))
	for i, v := range d.Get("containers").([]interface{}) {
		ContainersMap := v.(map[string]interface{})
		Containers[i] = make(map[string]interface{})
		Ports := make([]map[string]interface{}, len(ContainersMap["ports"].(*schema.Set).List()))
		for i, PortsValue := range ContainersMap["ports"].(*schema.Set).List() {
			PortsMap := PortsValue.(map[string]interface{})
			Ports[i] = make(map[string]interface{})
			Ports[i]["Port"] = PortsMap["port"]
			Ports[i]["Protocol"] = PortsMap["protocol"]
		}
		Containers[i]["Port"] = Ports
		EnvironmentVars := make([]map[string]interface{}, len(ContainersMap["environment_vars"].(*schema.Set).List()))
		for i, EnvironmentVarsValue := range ContainersMap["environment_vars"].(*schema.Set).List() {
			EnvironmentVarsMap := EnvironmentVarsValue.(map[string]interface{})
			EnvironmentVars[i] = make(map[string]interface{})
			EnvironmentVars[i]["Key"] = EnvironmentVarsMap["key"]
			EnvironmentVars[i]["Value"] = EnvironmentVarsMap["value"]
			EnvironmentVars[i]["FieldRefFieldPath"] = EnvironmentVarsMap["field_ref_field_path"]
		}
		Containers[i]["EnvironmentVar"] = EnvironmentVars
		Containers[i]["WorkingDir"] = ContainersMap["working_dir"]
		Containers[i]["Arg"] = ContainersMap["args"]
		Containers[i]["Cpu"] = ContainersMap["cpu"]
		Containers[i]["Gpu"] = ContainersMap["gpu"]
		Containers[i]["Memory"] = ContainersMap["memory"]
		Containers[i]["Name"] = ContainersMap["name"]
		Containers[i]["Image"] = ContainersMap["image"]
		Containers[i]["ImagePullPolicy"] = ContainersMap["image_pull_policy"]
		Containers[i]["SecurityContext.Capability.Add"] = ContainersMap["security_context_capability_adds"]
		Containers[i]["LifecyclePreStopHandlerExec"] = ContainersMap["lifecycle_pre_stop_handler_execs"]
		Containers[i]["SecurityContext.ReadOnlyRootFilesystem"] = ContainersMap["security_context_read_only_root_file_system"]
		Containers[i]["Tty"] = ContainersMap["tty"]
		Containers[i]["Stdin"] = ContainersMap["stdin"]
		Containers[i]["SecurityContext.RunAsUser"] = ContainersMap["security_context_run_as_user"]

		Containers[i]["ReadinessProbe.Exec.Command"] = ContainersMap["readiness_probe_exec_commands"]
		if ContainersMap["readiness_probe_period_seconds"] != nil && ContainersMap["readiness_probe_period_seconds"] != 0 {
			Containers[i]["ReadinessProbe.PeriodSeconds"] = ContainersMap["readiness_probe_period_seconds"]
		}
		Containers[i]["ReadinessProbe.HttpGet.Path"] = ContainersMap["readiness_probe_http_get_path"]
		if ContainersMap["readiness_probe_failure_threshold"] != 0 {
			Containers[i]["ReadinessProbe.FailureThreshold"] = ContainersMap["readiness_probe_failure_threshold"]
		}
		if ContainersMap["readiness_probe_initial_delay_seconds"] != 0 {
			Containers[i]["ReadinessProbe.InitialDelaySeconds"] = ContainersMap["readiness_probe_initial_delay_seconds"]
		}
		if ContainersMap["readiness_probe_http_get_port"] != 0 {
			Containers[i]["ReadinessProbe.HttpGet.Port"] = ContainersMap["readiness_probe_http_get_port"]
		}
		Containers[i]["ReadinessProbe.HttpGet.Scheme"] = ContainersMap["readiness_probe_http_get_scheme"]
		if ContainersMap["readiness_probe_tcp_socket_port"] != 0 {
			Containers[i]["ReadinessProbe.TcpSocket.Port"] = ContainersMap["readiness_probe_tcp_socket_port"]
		}
		if ContainersMap["readiness_probe_success_threshold"] != nil && ContainersMap["readiness_probe_success_threshold"] != 0 {
			Containers[i]["ReadinessProbe.SuccessThreshold"] = ContainersMap["readiness_probe_success_threshold"]
		}
		if ContainersMap["readiness_probe_timeout_seconds"] != nil && ContainersMap["readiness_probe_timeout_seconds"] != 0 {
			Containers[i]["ReadinessProbe.TimeoutSeconds"] = ContainersMap["readiness_probe_timeout_seconds"]
		}

		Containers[i]["LivenessProbe.Exec.Command"] = ContainersMap["liveness_probe_exec_commands"]
		if ContainersMap["liveness_probe_period_seconds"] != nil && ContainersMap["liveness_probe_period_seconds"] != 0 {
			Containers[i]["LivenessProbe.PeriodSeconds"] = ContainersMap["liveness_probe_period_seconds"]
		}
		Containers[i]["LivenessProbe.HttpGet.Path"] = ContainersMap["liveness_probe_http_get_path"]
		if ContainersMap["liveness_probe_failure_threshold"] != 0 {
			Containers[i]["LivenessProbe.FailureThreshold"] = ContainersMap["liveness_probe_failure_threshold"]
		}
		if ContainersMap["liveness_probe_initial_delay_seconds"] != 0 {
			Containers[i]["LivenessProbe.InitialDelaySeconds"] = ContainersMap["liveness_probe_initial_delay_seconds"]
		}
		if ContainersMap["liveness_probe_http_get_port"] != 0 {
			Containers[i]["LivenessProbe.HttpGet.Port"] = ContainersMap["liveness_probe_http_get_port"]
		}
		Containers[i]["LivenessProbe.HttpGet.Scheme"] = ContainersMap["liveness_probe_http_get_scheme"]
		if ContainersMap["liveness_probe_tcp_socket_port"] != 0 {
			Containers[i]["LivenessProbe.TcpSocket.Port"] = ContainersMap["liveness_probe_tcp_socket_port"]
		}
		if ContainersMap["liveness_probe_success_threshold"] != nil && ContainersMap["liveness_probe_success_threshold"] != 0 {
			Containers[i]["LivenessProbe.SuccessThreshold"] = ContainersMap["liveness_probe_success_threshold"]
		}
		if ContainersMap["liveness_probe_timeout_seconds"] != nil && ContainersMap["liveness_probe_timeout_seconds"] != 0 {
			Containers[i]["LivenessProbe.TimeoutSeconds"] = ContainersMap["liveness_probe_timeout_seconds"]
		}

		VolumeMounts := make([]map[string]interface{}, len(ContainersMap["volume_mounts"].(*schema.Set).List()))
		for i, VolumeMountsValue := range ContainersMap["volume_mounts"].(*schema.Set).List() {
			VolumeMountsMap := VolumeMountsValue.(map[string]interface{})
			VolumeMounts[i] = make(map[string]interface{})
			VolumeMounts[i]["MountPath"] = VolumeMountsMap["mount_path"]
			VolumeMounts[i]["Name"] = VolumeMountsMap["name"]
			VolumeMounts[i]["ReadOnly"] = VolumeMountsMap["read_only"]
			VolumeMounts[i]["SubPath"] = VolumeMountsMap["sub_path"]
			if VolumeMountsMap["mount_propagation"] != nil && VolumeMountsMap["mount_propagation"] != "" {
				VolumeMounts[i]["MountPropagation"] = VolumeMountsMap["mount_propagation"]
			}
		}
		Containers[i]["VolumeMount"] = VolumeMounts
		Containers[i]["Command"] = ContainersMap["commands"]
	}
	request["Container"] = Containers
	// instance_types
	types := make([]string, 0, int(5))
	instanceTypes := d.Get("instance_types").([]interface{})
	if instanceTypes != nil && len(instanceTypes) > 0 {
		types = expandStringList(instanceTypes)
		request["InstanceType"] = types
	}
	if _, ok := d.GetOk("init_containers"); ok {
		InitContainers := make([]map[string]interface{}, len(d.Get("init_containers").([]interface{})))
		for i, v := range d.Get("init_containers").([]interface{}) {
			InitContainersMap := v.(map[string]interface{})
			InitContainers[i] = make(map[string]interface{})
			Ports := make([]map[string]interface{}, len(InitContainersMap["ports"].(*schema.Set).List()))
			for i, PortsValue := range InitContainersMap["ports"].(*schema.Set).List() {
				PortsMap := PortsValue.(map[string]interface{})
				Ports[i] = make(map[string]interface{})
				Ports[i]["Port"] = PortsMap["port"]
				Ports[i]["Protocol"] = PortsMap["protocol"]
			}
			InitContainers[i]["InitContainerPort"] = Ports
			EnvironmentVars := make([]map[string]interface{}, len(InitContainersMap["environment_vars"].(*schema.Set).List()))
			for i, EnvironmentVarsValue := range InitContainersMap["environment_vars"].(*schema.Set).List() {
				EnvironmentVarsMap := EnvironmentVarsValue.(map[string]interface{})
				EnvironmentVars[i] = make(map[string]interface{})
				EnvironmentVars[i]["Key"] = EnvironmentVarsMap["key"]
				EnvironmentVars[i]["Value"] = EnvironmentVarsMap["value"]
				EnvironmentVars[i]["FieldRefFieldPath"] = EnvironmentVarsMap["field_ref_field_path"]
			}
			InitContainers[i]["InitContainerEnvironmentVar"] = EnvironmentVars
			InitContainers[i]["WorkingDir"] = InitContainersMap["working_dir"]
			InitContainers[i]["Arg"] = InitContainersMap["args"]
			InitContainers[i]["Cpu"] = InitContainersMap["cpu"]
			InitContainers[i]["Gpu"] = InitContainersMap["gpu"]
			InitContainers[i]["Memory"] = InitContainersMap["memory"]
			InitContainers[i]["Name"] = InitContainersMap["name"]
			InitContainers[i]["Image"] = InitContainersMap["image"]
			InitContainers[i]["ImagePullPolicy"] = InitContainersMap["image_pull_policy"]
			InitContainers[i]["SecurityContext.Capability.Add"] = InitContainersMap["security_context_capability_adds"]
			InitContainers[i]["SecurityContext.ReadOnlyRootFilesystem"] = InitContainersMap["security_context_read_only_root_file_system"]
			InitContainers[i]["SecurityContext.RunAsUser"] = InitContainersMap["security_context_run_as_user"]
			VolumeMounts := make([]map[string]interface{}, len(InitContainersMap["volume_mounts"].(*schema.Set).List()))
			for i, VolumeMountsValue := range InitContainersMap["volume_mounts"].(*schema.Set).List() {
				VolumeMountsMap := VolumeMountsValue.(map[string]interface{})
				VolumeMounts[i] = make(map[string]interface{})
				VolumeMounts[i]["MountPath"] = VolumeMountsMap["mount_path"]
				VolumeMounts[i]["Name"] = VolumeMountsMap["name"]
				VolumeMounts[i]["ReadOnly"] = VolumeMountsMap["read_only"]
				VolumeMounts[i]["SubPath"] = VolumeMountsMap["sub_path"]
				if VolumeMountsMap["mount_propagation"] != nil && VolumeMountsMap["mount_propagation"] != "" {
					VolumeMounts[i]["MountPropagation"] = VolumeMountsMap["mount_propagation"]
				}
			}
			InitContainers[i]["InitContainerVolumeMount"] = VolumeMounts
			InitContainers[i]["Command"] = InitContainersMap["commands"]
		}
		request["InitContainer"] = InitContainers
	}

	if v, ok := d.GetOk("volumes"); ok {
		Volumes := make([]map[string]interface{}, len(v.(*schema.Set).List()))
		for i, v := range v.(*schema.Set).List() {
			VolumesMap := v.(map[string]interface{})
			Volumes[i] = make(map[string]interface{})
			ConfigFileVolumeConfigFileToPaths := make([]map[string]interface{}, len(VolumesMap["config_file_volume_config_file_to_paths"].(*schema.Set).List()))
			for i, ConfigFileVolumeConfigFileToPathsValue := range VolumesMap["config_file_volume_config_file_to_paths"].(*schema.Set).List() {
				ConfigFileVolumeConfigFileToPathsMap := ConfigFileVolumeConfigFileToPathsValue.(map[string]interface{})
				ConfigFileVolumeConfigFileToPaths[i] = make(map[string]interface{})
				ConfigFileVolumeConfigFileToPaths[i]["Content"] = ConfigFileVolumeConfigFileToPathsMap["content"]
				ConfigFileVolumeConfigFileToPaths[i]["Path"] = ConfigFileVolumeConfigFileToPathsMap["path"]
				if ConfigFileVolumeConfigFileToPathsMap["mode"] != nil && ConfigFileVolumeConfigFileToPathsMap["mode"] != 0 {
					ConfigFileVolumeConfigFileToPaths[i]["Mode"] = ConfigFileVolumeConfigFileToPathsMap["mode"]
				}
			}
			Volumes[i]["ConfigFileVolumeConfigFileToPath"] = ConfigFileVolumeConfigFileToPaths
			Volumes[i]["DiskVolume.DiskId"] = VolumesMap["disk_volume_disk_id"]
			Volumes[i]["HostPathVolume.Type"] = VolumesMap["host_path_volume_type"]
			Volumes[i]["HostPathVolume.Path"] = VolumesMap["host_path_volume_path"]
			Volumes[i]["EmptyDirVolume.SizeLimit"] = VolumesMap["empty_dir_volume_size_limit"]
			if VolumesMap["config_file_volume_default_mode"] != nil && VolumesMap["config_file_volume_default_mode"] != 0 {
				Volumes[i]["ConfigFileVolumeDefaultMode"] = VolumesMap["config_file_volume_default_mode"]
			}
			Volumes[i]["EmptyDirVolume.Medium"] = VolumesMap["empty_dir_volume_medium"]

			Volumes[i]["DiskVolume.FsType"] = VolumesMap["disk_volume_fs_type"]
			Volumes[i]["DiskVolume.DiskSize"] = VolumesMap["disk_volume_disk_size"]
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

	if v, ok := d.GetOk("host_aliases"); ok {
		HostAliases := make([]map[string]interface{}, len(v.(*schema.Set).List()))
		for i, HostAliasesValue := range v.(*schema.Set).List() {
			HostAliasesMap := HostAliasesValue.(map[string]interface{})
			HostAliases[i] = make(map[string]interface{})
			HostAliases[i]["Hostname"] = HostAliasesMap["hostnames"]
			HostAliases[i]["Ip"] = HostAliasesMap["ip"]
		}
		request["HostAliase"] = HostAliases
	}
	response, err = client.RpcPost("Ess", "2014-08-28", action, nil, request, true)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ess_eci_scaling_configuration", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ScalingConfigurationId"]))

	return resourceAliyunEssEciScalingConfigurationUpdate(d, meta)
}

func resourceAliyunEssEciScalingConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	o, err := essService.DescribeEssEciScalingConfiguration(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("active", d.Get("active"))
	d.Set("force_delete", d.Get("force_delete"))
	d.Set("scaling_group_id", o["ScalingGroupId"])
	d.Set("scaling_configuration_name", o["ScalingConfigurationName"])
	d.Set("description", o["Description"])
	d.Set("security_group_id", o["SecurityGroupId"])
	d.Set("container_group_name", o["ContainerGroupName"])
	d.Set("restart_policy", o["RestartPolicy"])
	d.Set("cpu", o["Cpu"])
	d.Set("memory", o["Memory"])
	d.Set("resource_group_id", o["ResourceGroupId"])
	d.Set("dns_policy", o["DnsPolicy"])
	d.Set("enable_sls", o["SlsEnable"])
	d.Set("cost_optimization", o["CostOptimization"])
	d.Set("image_snapshot_id", o["ImageSnapshotId"])
	if o["InstanceFamilyLevel"] != nil && o["InstanceFamilyLevel"] != "" {
		d.Set("instance_family_level", o["InstanceFamilyLevel"])
	}
	d.Set("ram_role_name", o["RamRoleName"])
	d.Set("termination_grace_period_seconds", o["TerminationGracePeriodSeconds"])
	d.Set("auto_match_image_cache", o["AutoMatchImageCache"])
	d.Set("ipv6_address_count", o["Ipv6AddressCount"])
	if o["CpuOptionsCore"] != nil && o["CpuOptionsCore"] != 0 {
		d.Set("cpu_options_core", o["CpuOptionsCore"])
	}
	if o["CpuOptionsThreadsPerCore"] != nil && o["CpuOptionsThreadsPerCore"] != 0 {
		d.Set("cpu_options_threads_per_core", o["CpuOptionsThreadsPerCore"])

	}
	if o["ActiveDeadlineSeconds"] != nil && o["ActiveDeadlineSeconds"] != 0 {
		d.Set("active_deadline_seconds", o["ActiveDeadlineSeconds"])
	}
	d.Set("auto_create_eip", o["AutoCreateEip"])
	d.Set("eip_bandwidth", o["EipBandwidth"])
	d.Set("host_name", o["HostName"])
	d.Set("ingress_bandwidth", o["IngressBandwidth"])
	d.Set("egress_bandwidth", o["EgressBandwidth"])
	d.Set("ephemeral_storage", o["EphemeralStorage"])
	d.Set("load_balancer_weight", o["LoadBalancerWeight"])
	d.Set("tags", o["Tags"])
	d.Set("instance_types", o["InstanceTypes"])
	d.Set("spot_strategy", o["SpotStrategy"])
	if o["spot_price_limit"] != nil {
		d.Set("spot_price_limit", strconv.FormatFloat(o["SpotPriceLimit"].(float64), 'f', 2, 64))
	}

	credentials := make([]map[string]interface{}, 0)
	if credentialList, ok := o["ImageRegistryCredentials"].([]interface{}); ok {
		for _, v := range credentialList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"password": m1["Password"],
					"server":   m1["Server"],
					"username": m1["UserName"],
				}
				credentials = append(credentials, temp1)
			}
		}
	}
	d.Set("image_registry_credentials", credentials)

	options := make([]map[string]interface{}, 0)
	if optionList, ok := o["DnsConfigOptions"].([]interface{}); ok {
		for _, v := range optionList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"name":  m1["Name"],
					"value": m1["Value"],
				}
				options = append(options, temp1)
			}
		}
	}
	d.Set("dns_config_options", options)

	sysctls := make([]map[string]interface{}, 0)
	if sysctlList, ok := o["SecurityContextSysCtls"].([]interface{}); ok {
		for _, v := range sysctlList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"name":  m1["Name"],
					"value": m1["Value"],
				}
				sysctls = append(sysctls, temp1)
			}
		}
	}
	d.Set("security_context_sysctls", sysctls)

	acrRegistryInfos := make([]map[string]interface{}, 0)
	if acrRegistryInfoList, ok := o["AcrRegistryInfos"].([]interface{}); ok {
		for _, v := range acrRegistryInfoList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"domains":       m1["Domains"],
					"instance_name": m1["InstanceName"],
					"instance_id":   m1["InstanceId"],
					"region_id":     m1["RegionId"],
				}
				acrRegistryInfos = append(acrRegistryInfos, temp1)
			}
		}
	}
	d.Set("acr_registry_infos", acrRegistryInfos)

	containers := make([]map[string]interface{}, 0)
	if containersList, ok := o["Containers"].([]interface{}); ok {
		for _, v := range containersList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"security_context_capability_adds":            m1["SecurityContextCapabilityAdds"],
					"lifecycle_pre_stop_handler_execs":            m1["LifecyclePreStopHandlerExecs"],
					"security_context_read_only_root_file_system": m1["SecurityContextReadOnlyRootFilesystem"],
					"tty":                          m1["Tty"],
					"stdin":                        m1["Stdin"],
					"security_context_run_as_user": m1["SecurityContextRunAsUser"],
					"working_dir":                  m1["WorkingDir"],
					"args":                         m1["Args"],
					"cpu":                          m1["Cpu"],
					"gpu":                          m1["Gpu"],
					"memory":                       m1["Memory"],
					"name":                         m1["Name"],
					"image":                        m1["Image"],
					"image_pull_policy":            m1["ImagePullPolicy"],
					"commands":                     m1["Commands"],

					"readiness_probe_exec_commands":         m1["ReadinessProbeExecCommands"],
					"readiness_probe_http_get_path":         m1["ReadinessProbeHttpGetPath"],
					"readiness_probe_failure_threshold":     m1["ReadinessProbeFailureThreshold"],
					"readiness_probe_initial_delay_seconds": m1["ReadinessProbeInitialDelaySeconds"],
					"readiness_probe_http_get_port":         m1["ReadinessProbeHttpGetPort"],
					"readiness_probe_http_get_scheme":       m1["ReadinessProbeHttpGetScheme"],
					"readiness_probe_tcp_socket_port":       m1["ReadinessProbeTcpSocketPort"],

					"liveness_probe_exec_commands":         m1["LivenessProbeExecCommands"],
					"liveness_probe_http_get_path":         m1["LivenessProbeHttpGetPath"],
					"liveness_probe_failure_threshold":     m1["LivenessProbeFailureThreshold"],
					"liveness_probe_initial_delay_seconds": m1["LivenessProbeInitialDelaySeconds"],
					"liveness_probe_http_get_port":         m1["LivenessProbeHttpGetPort"],
					"liveness_probe_http_get_scheme":       m1["LivenessProbeHttpGetScheme"],
					"liveness_probe_tcp_socket_port":       m1["LivenessProbeTcpSocketPort"],
				}
				if m1["LivenessProbeTimeoutSeconds"] != nil && m1["LivenessProbeTimeoutSeconds"] != 0 {
					temp1["liveness_probe_timeout_seconds"] = m1["LivenessProbeTimeoutSeconds"]
				}
				if m1["LivenessProbeSuccessThreshold"] != nil && m1["LivenessProbeSuccessThreshold"] != 0 {
					temp1["liveness_probe_success_threshold"] = m1["LivenessProbeSuccessThreshold"]
				}
				if m1["LivenessProbePeriodSeconds"] != nil && m1["LivenessProbePeriodSeconds"] != 0 {
					temp1["liveness_probe_period_seconds"] = m1["LivenessProbePeriodSeconds"]
				}
				if m1["ReadinessProbeTimeoutSeconds"] != nil && m1["ReadinessProbeTimeoutSeconds"] != 0 {
					temp1["readiness_probe_timeout_seconds"] = m1["ReadinessProbeTimeoutSeconds"]
				}
				if m1["ReadinessProbeSuccessThreshold"] != nil && m1["ReadinessProbeSuccessThreshold"] != 0 {
					temp1["readiness_probe_success_threshold"] = m1["ReadinessProbeSuccessThreshold"]
				}

				if m1["ReadinessProbePeriodSeconds"] != nil && m1["ReadinessProbePeriodSeconds"] != 0 {
					temp1["readiness_probe_period_seconds"] = m1["ReadinessProbePeriodSeconds"]
				}
				if m1["EnvironmentVars"] != nil {
					environmentVarsMaps := make([]map[string]interface{}, 0)
					for _, environmentVarsValue := range m1["EnvironmentVars"].([]interface{}) {
						environmentVars := environmentVarsValue.(map[string]interface{})
						environmentVarsMap := map[string]interface{}{
							"key":                  environmentVars["Key"],
							"value":                environmentVars["Value"],
							"field_ref_field_path": environmentVars["FieldRefFieldPath"],
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
							"sub_path":   volumeMounts["SubPath"],
						}
						if volumeMounts["MountPropagation"] != nil && volumeMounts["MountPropagation"] != "" {
							volumeMountsMap["mount_propagation"] = volumeMounts["MountPropagation"]
						}
						volumeMountsMaps = append(volumeMountsMaps, volumeMountsMap)
					}
					temp1["volume_mounts"] = volumeMountsMaps
				}
				containers = append(containers, temp1)
			}
		}
	}

	if err := d.Set("containers", containers); err != nil {
		return WrapError(err)
	}

	initContainers := make([]map[string]interface{}, 0)
	if initContainersList, ok := o["InitContainers"].([]interface{}); ok {
		for _, v := range initContainersList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"security_context_capability_adds":            m1["SecurityContextCapabilityAdds"],
					"security_context_read_only_root_file_system": m1["SecurityContextReadOnlyRootFilesystem"],
					"security_context_run_as_user":                m1["SecurityContextRunAsUser"],
					"working_dir":                                 m1["WorkingDir"],
					"args":                                        m1["InitContainerArgs"],
					"cpu":                                         m1["Cpu"],
					"gpu":                                         m1["Gpu"],
					"memory":                                      m1["Memory"],
					"image":                                       m1["Image"],
					"image_pull_policy":                           m1["ImagePullPolicy"],
					"name":                                        m1["Name"],
					"commands":                                    m1["InitContainerCommands"],
				}
				if m1["InitContainerEnvironmentVars"] != nil {
					environmentVarsMaps := make([]map[string]interface{}, 0)
					for _, environmentVarsValue := range m1["InitContainerEnvironmentVars"].([]interface{}) {
						environmentVars := environmentVarsValue.(map[string]interface{})
						environmentVarsMap := map[string]interface{}{
							"key":                  environmentVars["Key"],
							"value":                environmentVars["Value"],
							"field_ref_field_path": environmentVars["FieldRefFieldPath"],
						}
						environmentVarsMaps = append(environmentVarsMaps, environmentVarsMap)
					}
					temp1["environment_vars"] = environmentVarsMaps
				}
				if m1["InitContainerPorts"] != nil {
					portsMaps := make([]map[string]interface{}, 0)
					for _, portsValue := range m1["InitContainerPorts"].([]interface{}) {
						ports := portsValue.(map[string]interface{})
						portsMap := map[string]interface{}{
							"port":     ports["Port"],
							"protocol": ports["Protocol"],
						}
						portsMaps = append(portsMaps, portsMap)
					}
					temp1["ports"] = portsMaps
				}
				if m1["InitContainerVolumeMounts"] != nil {
					volumeMountsMaps := make([]map[string]interface{}, 0)
					for _, volumeMountsValue := range m1["InitContainerVolumeMounts"].([]interface{}) {
						volumeMounts := volumeMountsValue.(map[string]interface{})
						volumeMountsMap := map[string]interface{}{
							"mount_path": volumeMounts["MountPath"],
							"name":       volumeMounts["Name"],
							"read_only":  volumeMounts["ReadOnly"],
							"sub_path":   volumeMounts["SubPath"],
						}
						if volumeMounts["MountPropagation"] != nil && volumeMounts["MountPropagation"] != "" {
							volumeMountsMap["mount_propagation"] = volumeMounts["MountPropagation"]
						}
						volumeMountsMaps = append(volumeMountsMaps, volumeMountsMap)
					}
					temp1["volume_mounts"] = volumeMountsMaps
				}
				initContainers = append(initContainers, temp1)

			}
		}
	}
	if err := d.Set("init_containers", initContainers); err != nil {
		return WrapError(err)
	}

	volumes := make([]map[string]interface{}, 0)
	if volumesList, ok := o["Volumes"].([]interface{}); ok {
		for _, v := range volumesList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"disk_volume_disk_id":     m1["DiskVolumeDiskId"],
					"disk_volume_fs_type":     m1["DiskVolumeFsType"],
					"disk_volume_disk_size":   m1["DiskVolumeDiskSize"],
					"flex_volume_driver":      m1["FlexVolumeDriver"],
					"flex_volume_fs_type":     m1["FlexVolumeFsType"],
					"flex_volume_options":     m1["FlexVolumeOptions"],
					"nfs_volume_path":         m1["NFSVolumePath"],
					"nfs_volume_read_only":    m1["NFSVolumeReadOnly"],
					"nfs_volume_server":       m1["NFSVolumeServer"],
					"name":                    m1["Name"],
					"type":                    m1["Type"],
					"empty_dir_volume_medium": m1["EmptyDirVolumeMedium"],
				}
				if m1["HostPathVolumeType"] != nil && m1["HostPathVolumeType"] != "" {
					temp1["host_path_volume_type"] = m1["HostPathVolumeType"]
				}
				if m1["EmptyDirVolumeSizeLimit"] != nil && m1["EmptyDirVolumeSizeLimit"] != "" {
					temp1["empty_dir_volume_size_limit"] = m1["EmptyDirVolumeSizeLimit"]
				}
				if m1["HostPathVolumePath"] != nil && m1["HostPathVolumePath"] != "" {
					temp1["host_path_volume_path"] = m1["HostPathVolumePath"]
				}
				if m1["ConfigFileVolumeDefaultMode"] != nil && m1["ConfigFileVolumeDefaultMode"] != 0 {
					temp1["config_file_volume_default_mode"] = m1["ConfigFileVolumeDefaultMode"]
				}
				if m1["ConfigFileVolumeConfigFileToPaths"] != nil {
					configFileVolumeConfigFileToPathsMaps := make([]map[string]interface{}, 0)
					for _, configFileVolumeConfigFileToPathsValue := range m1["ConfigFileVolumeConfigFileToPaths"].([]interface{}) {
						configFileVolumeConfigFileToPaths := configFileVolumeConfigFileToPathsValue.(map[string]interface{})
						configFileVolumeConfigFileToPathsMap := map[string]interface{}{
							"content": configFileVolumeConfigFileToPaths["Content"],
							"path":    configFileVolumeConfigFileToPaths["Path"],
						}
						if configFileVolumeConfigFileToPaths["Mode"] != nil && configFileVolumeConfigFileToPaths["Mode"] != 0 {
							configFileVolumeConfigFileToPathsMap["mode"] = configFileVolumeConfigFileToPaths["Mode"]
						}
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

	hostAliases := make([]map[string]interface{}, 0)
	if hostAliasesList, ok := o["HostAliases"].([]interface{}); ok {
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

	return nil
}

func resourceAliyunEssEciScalingConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	d.Partial(true)

	sgId := d.Get("scaling_group_id").(string)
	group, err := essService.DescribeEssScalingGroup(sgId)
	if err != nil {
		return WrapError(err)
	}

	if d.HasChange("active") {
		c, err := essService.DescribeEssEciScalingConfiguration(d.Id())
		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapError(err)
		}
		if d.Get("active").(bool) {
			if c["LifecycleState"] == string(Inactive) {
				modifyGroupRequest := ess.CreateModifyScalingGroupRequest()
				modifyGroupRequest.ScalingGroupId = sgId
				modifyGroupRequest.ActiveScalingConfigurationId = d.Id()
				_, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
					return essClient.ModifyScalingGroup(modifyGroupRequest)
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyGroupRequest.GetActionName(), AlibabaCloudSdkGoERROR)
				}
			}
			if group.LifecycleState == string(Inactive) {
				enableGroupRequest := ess.CreateEnableScalingGroupRequest()
				enableGroupRequest.RegionId = client.RegionId
				enableGroupRequest.ScalingGroupId = sgId
				_, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
					return essClient.EnableScalingGroup(enableGroupRequest)
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), enableGroupRequest.GetActionName(), AlibabaCloudSdkGoERROR)
				}
			}
		}
		d.SetPartial("active")
	}

	request := map[string]interface{}{
		"ScalingConfigurationId": d.Id(),
		"RegionId":               client.RegionId,
	}
	update := false

	if d.HasChange("scaling_configuration_name") {
		update = true
		request["ScalingConfigurationName"] = d.Get("scaling_configuration_name")
	}
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	if d.HasChange("security_group_id") {
		update = true
		request["SecurityGroupId"] = d.Get("security_group_id")
	}
	if d.HasChange("container_group_name") {
		update = true
		request["ContainerGroupName"] = d.Get("container_group_name")
	}
	if d.HasChange("restart_policy") {
		update = true
		request["RestartPolicy"] = d.Get("restart_policy")
	}
	if d.HasChange("cpu") {
		update = true
		request["Cpu"] = d.Get("cpu")
	}
	if d.HasChange("memory") {
		update = true
		request["Memory"] = d.Get("memory")
	}
	if d.HasChange("resource_group_id") {
		update = true
		request["ResourceGroupId"] = d.Get("resource_group_id")
	}
	if d.HasChange("dns_policy") {
		update = true
		request["DnsPolicy"] = d.Get("dns_policy")
	}
	if d.HasChange("cost_optimization") {
		update = true
		request["CostOptimization"] = d.Get("cost_optimization")
	}
	if d.HasChange("enable_sls") {
		update = true
		request["EnableSls"] = d.Get("enable_sls")
	}
	if d.HasChange("image_snapshot_id") {
		update = true
		request["ImageSnapshotId"] = d.Get("image_snapshot_id")
	}
	if d.HasChange("instance_family_level") {
		update = true
		request["InstanceFamilyLevel"] = d.Get("instance_family_level")
	}
	if d.HasChange("ram_role_name") {
		update = true
		request["RamRoleName"] = d.Get("ram_role_name")
	}
	if d.HasChange("termination_grace_period_seconds") {
		update = true
		request["TerminationGracePeriodSeconds"] = d.Get("termination_grace_period_seconds")
	}
	if d.HasChange("auto_match_image_cache") {
		update = true
		request["AutoMatchImageCache"] = d.Get("auto_match_image_cache")
	}
	if d.HasChange("ipv6_address_count") {
		update = true
		request["Ipv6AddressCount"] = d.Get("ipv6_address_count")
	}
	if d.HasChange("cpu_options_core") {
		if v, ok := d.GetOkExists("cpu_options_core"); ok {
			if v != 0 {
				update = true
				request["CpuOptionsCore"] = v
			}
		}
	}
	if d.HasChange("cpu_options_threads_per_core") {
		if v, ok := d.GetOkExists("cpu_options_threads_per_core"); ok {
			if v != 0 {
				update = true
				request["CpuOptionsThreadsPerCore"] = v
			}
		}
	}
	if d.HasChange("active_deadline_seconds") {
		update = true
		v := d.Get("active_deadline_seconds")
		request["ActiveDeadlineSeconds"] = v
	}
	if d.HasChange("spot_strategy") {
		update = true
		request["SpotStrategy"] = d.Get("spot_strategy")
	}
	if d.HasChange("spot_price_limit") {
		update = true
		request["SpotPriceLimit"] = strconv.FormatFloat(d.Get("spot_price_limit").(float64), 'f', 2, 64)
	}
	if d.HasChange("auto_create_eip") {
		update = true
		request["AutoCreateEip"] = d.Get("auto_create_eip")
	}
	if d.HasChange("eip_bandwidth") {
		update = true
		request["EipBandwidth"] = d.Get("eip_bandwidth")
	}
	if d.HasChange("host_name") {
		update = true
		request["HostName"] = d.Get("host_name")
	}
	if d.HasChange("ingress_bandwidth") {
		update = true
		request["IngressBandwidth"] = d.Get("ingress_bandwidth")
	}
	if d.HasChange("egress_bandwidth") {
		update = true
		request["EgressBandwidth"] = d.Get("egress_bandwidth")
	}
	if d.HasChange("ephemeral_storage") {
		update = true
		request["EphemeralStorage"] = d.Get("ephemeral_storage")
	}
	if d.HasChange("load_balancer_weight") {
		update = true
		request["LoadBalancerWeight"] = d.Get("load_balancer_weight")
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

	if d.HasChange("image_registry_credentials") {
		update = true
		if v, ok := d.GetOk("image_registry_credentials"); ok {
			imageRegisryCredentialMaps := make([]map[string]interface{}, 0)
			for _, raw := range v.(*schema.Set).List() {
				obj := raw.(map[string]interface{})
				imageRegisryCredentialMaps = append(imageRegisryCredentialMaps, map[string]interface{}{
					"Password": obj["password"],
					"Server":   obj["server"],
					"UserName": obj["username"],
				})
			}
			request["ImageRegistryCredential"] = imageRegisryCredentialMaps
		}
	}

	if d.HasChange("dns_config_options") {
		update = true
		if v, ok := d.GetOk("dns_config_options"); ok {
			dnsConfigOptionMaps := make([]map[string]interface{}, 0)
			for _, raw := range v.(*schema.Set).List() {
				obj := raw.(map[string]interface{})
				dnsConfigOptionMaps = append(dnsConfigOptionMaps, map[string]interface{}{
					"Name":  obj["name"],
					"Value": obj["value"],
				})
			}
			request["DnsConfigOption"] = dnsConfigOptionMaps
		}
	}

	if d.HasChange("security_context_sysctls") {
		update = true
		if v, ok := d.GetOk("security_context_sysctls"); ok {
			securityContextSysctlMaps := make([]map[string]interface{}, 0)
			for _, raw := range v.(*schema.Set).List() {
				obj := raw.(map[string]interface{})
				securityContextSysctlMaps = append(securityContextSysctlMaps, map[string]interface{}{
					"Name":  obj["name"],
					"Value": obj["value"],
				})
			}
			request["SecurityContextSysctl"] = securityContextSysctlMaps
		}
	}

	if d.HasChange("acr_registry_infos") {
		update = true
		if v, ok := d.GetOk("acr_registry_infos"); ok {
			acrRegistryInfoMaps := make([]map[string]interface{}, 0)
			for _, raw := range v.(*schema.Set).List() {
				obj := raw.(map[string]interface{})

				acrRegistryInfoMaps = append(acrRegistryInfoMaps, map[string]interface{}{
					"Domain":       expandStringList(obj["domains"].(*schema.Set).List()),
					"InstanceName": obj["instance_name"],
					"InstanceId":   obj["instance_id"],
					"RegionId":     obj["region_id"],
				})
			}
			request["AcrRegistryInfo"] = acrRegistryInfoMaps
		}
	}

	if d.HasChange("containers") {
		update = true
		Containers := make([]map[string]interface{}, len(d.Get("containers").([]interface{})))
		for i, ContainersValue := range d.Get("containers").([]interface{}) {
			ContainersMap := ContainersValue.(map[string]interface{})
			Containers[i] = make(map[string]interface{})
			Containers[i]["SecurityContext.Capability.Add"] = ContainersMap["security_context_capability_adds"]
			Containers[i]["LifecyclePreStopHandlerExec"] = ContainersMap["lifecycle_pre_stop_handler_execs"]
			Containers[i]["SecurityContext.ReadOnlyRootFilesystem"] = ContainersMap["security_context_read_only_root_file_system"]
			Containers[i]["Tty"] = ContainersMap["tty"]
			Containers[i]["Stdin"] = ContainersMap["stdin"]
			Containers[i]["SecurityContext.RunAsUser"] = ContainersMap["security_context_run_as_user"]
			Containers[i]["WorkingDir"] = ContainersMap["working_dir"]
			Containers[i]["Arg"] = ContainersMap["args"]
			Containers[i]["Cpu"] = ContainersMap["cpu"]
			Containers[i]["Gpu"] = ContainersMap["gpu"]
			Containers[i]["Memory"] = ContainersMap["memory"]
			Containers[i]["Name"] = ContainersMap["name"]
			Containers[i]["Image"] = ContainersMap["image"]
			Containers[i]["ImagePullPolicy"] = ContainersMap["image_pull_policy"]
			Containers[i]["Command"] = ContainersMap["commands"]

			Containers[i]["ReadinessProbe.Exec.Command"] = ContainersMap["readiness_probe_exec_commands"]
			if ContainersMap["readiness_probe_period_seconds"] != nil && ContainersMap["readiness_probe_period_seconds"] != 0 {
				Containers[i]["ReadinessProbe.PeriodSeconds"] = ContainersMap["readiness_probe_period_seconds"]
			}
			Containers[i]["ReadinessProbe.HttpGet.Path"] = ContainersMap["readiness_probe_http_get_path"]
			if ContainersMap["readiness_probe_failure_threshold"] != 0 {
				Containers[i]["ReadinessProbe.FailureThreshold"] = ContainersMap["readiness_probe_failure_threshold"]
			}
			if ContainersMap["readiness_probe_initial_delay_seconds"] != 0 {
				Containers[i]["ReadinessProbe.InitialDelaySeconds"] = ContainersMap["readiness_probe_initial_delay_seconds"]
			}
			if ContainersMap["readiness_probe_http_get_port"] != 0 {
				Containers[i]["ReadinessProbe.HttpGet.Port"] = ContainersMap["readiness_probe_http_get_port"]
			}
			Containers[i]["ReadinessProbe.HttpGet.Scheme"] = ContainersMap["readiness_probe_http_get_scheme"]
			if ContainersMap["readiness_probe_tcp_socket_port"] != 0 {
				Containers[i]["ReadinessProbe.TcpSocket.Port"] = ContainersMap["readiness_probe_tcp_socket_port"]
			}
			if ContainersMap["readiness_probe_success_threshold"] != nil && ContainersMap["readiness_probe_success_threshold"] != 0 {
				Containers[i]["ReadinessProbe.SuccessThreshold"] = ContainersMap["readiness_probe_success_threshold"]
			}
			if ContainersMap["readiness_probe_timeout_seconds"] != nil && ContainersMap["readiness_probe_timeout_seconds"] != 0 {
				Containers[i]["ReadinessProbe.TimeoutSeconds"] = ContainersMap["readiness_probe_timeout_seconds"]
			}

			Containers[i]["LivenessProbe.Exec.Command"] = ContainersMap["liveness_probe_exec_commands"]
			if ContainersMap["liveness_probe_period_seconds"] != nil && ContainersMap["liveness_probe_period_seconds"] != 0 {
				Containers[i]["LivenessProbe.PeriodSeconds"] = ContainersMap["liveness_probe_period_seconds"]
			}
			Containers[i]["LivenessProbe.HttpGet.Path"] = ContainersMap["liveness_probe_http_get_path"]
			if ContainersMap["liveness_probe_failure_threshold"] != 0 {
				Containers[i]["LivenessProbe.FailureThreshold"] = ContainersMap["liveness_probe_failure_threshold"]
			}
			if ContainersMap["liveness_probe_initial_delay_seconds"] != 0 {
				Containers[i]["LivenessProbe.InitialDelaySeconds"] = ContainersMap["liveness_probe_initial_delay_seconds"]
			}
			if ContainersMap["liveness_probe_http_get_port"] != 0 {
				Containers[i]["LivenessProbe.HttpGet.Port"] = ContainersMap["liveness_probe_http_get_port"]
			}
			Containers[i]["LivenessProbe.HttpGet.Scheme"] = ContainersMap["liveness_probe_http_get_scheme"]
			if ContainersMap["liveness_probe_tcp_socket_port"] != 0 {
				Containers[i]["LivenessProbe.TcpSocket.Port"] = ContainersMap["liveness_probe_tcp_socket_port"]
			}
			if ContainersMap["liveness_probe_success_threshold"] != nil && ContainersMap["liveness_probe_success_threshold"] != 0 {
				Containers[i]["LivenessProbe.SuccessThreshold"] = ContainersMap["liveness_probe_success_threshold"]
			}
			if ContainersMap["liveness_probe_timeout_seconds"] != nil && ContainersMap["liveness_probe_timeout_seconds"] != 0 {
				Containers[i]["LivenessProbe.TimeoutSeconds"] = ContainersMap["liveness_probe_timeout_seconds"]
			}

			EnvironmentVars := make([]map[string]interface{}, len(ContainersMap["environment_vars"].(*schema.Set).List()))
			for i, EnvironmentVarsValue := range ContainersMap["environment_vars"].(*schema.Set).List() {
				EnvironmentVarsMap := EnvironmentVarsValue.(map[string]interface{})
				EnvironmentVars[i] = make(map[string]interface{})
				EnvironmentVars[i]["Key"] = EnvironmentVarsMap["key"]
				EnvironmentVars[i]["Value"] = EnvironmentVarsMap["value"]
				EnvironmentVars[i]["FieldRef.FieldPath"] = EnvironmentVarsMap["field_ref_field_path"]
			}
			Containers[i]["EnvironmentVar"] = EnvironmentVars

			Ports := make([]map[string]interface{}, len(ContainersMap["ports"].(*schema.Set).List()))
			for i, PortsValue := range ContainersMap["ports"].(*schema.Set).List() {
				PortsMap := PortsValue.(map[string]interface{})
				Ports[i] = make(map[string]interface{})
				Ports[i]["Port"] = PortsMap["port"]
				Ports[i]["Protocol"] = PortsMap["protocol"]
			}
			Containers[i]["Port"] = Ports

			VolumeMounts := make([]map[string]interface{}, len(ContainersMap["volume_mounts"].(*schema.Set).List()))
			for i, VolumeMountsValue := range ContainersMap["volume_mounts"].(*schema.Set).List() {
				VolumeMountsMap := VolumeMountsValue.(map[string]interface{})
				VolumeMounts[i] = make(map[string]interface{})
				VolumeMounts[i]["MountPath"] = VolumeMountsMap["mount_path"]
				VolumeMounts[i]["Name"] = VolumeMountsMap["name"]
				VolumeMounts[i]["ReadOnly"] = VolumeMountsMap["read_only"]
				VolumeMounts[i]["SubPath"] = VolumeMountsMap["sub_path"]
				if VolumeMountsMap["mount_propagation"] != nil && VolumeMountsMap["mount_propagation"] != "" {
					VolumeMounts[i]["MountPropagation"] = VolumeMountsMap["mount_propagation"]
				}
			}
			Containers[i]["VolumeMount"] = VolumeMounts
		}
		request["Container"] = Containers
	}

	if d.HasChange("init_containers") {
		update = true
		InitContainers := make([]map[string]interface{}, len(d.Get("init_containers").([]interface{})))
		for i, InitContainersValue := range d.Get("init_containers").([]interface{}) {
			InitContainersMap := InitContainersValue.(map[string]interface{})
			InitContainers[i] = make(map[string]interface{})
			InitContainers[i]["WorkingDir"] = InitContainersMap["working_dir"]
			InitContainers[i]["Arg"] = InitContainersMap["args"]
			InitContainers[i]["Cpu"] = InitContainersMap["cpu"]
			InitContainers[i]["Gpu"] = InitContainersMap["gpu"]
			InitContainers[i]["Memory"] = InitContainersMap["memory"]
			InitContainers[i]["Name"] = InitContainersMap["name"]
			InitContainers[i]["Image"] = InitContainersMap["image"]
			InitContainers[i]["ImagePullPolicy"] = InitContainersMap["image_pull_policy"]
			InitContainers[i]["Command"] = InitContainersMap["commands"]
			InitContainers[i]["SecurityContext.Capability.Add"] = InitContainersMap["security_context_capability_adds"]
			InitContainers[i]["SecurityContext.ReadOnlyRootFilesystem"] = InitContainersMap["security_context_read_only_root_file_system"]
			InitContainers[i]["SecurityContext.RunAsUser"] = InitContainersMap["security_context_run_as_user"]

			EnvironmentVars := make([]map[string]interface{}, len(InitContainersMap["environment_vars"].(*schema.Set).List()))
			for i, EnvironmentVarsValue := range InitContainersMap["environment_vars"].(*schema.Set).List() {
				EnvironmentVarsMap := EnvironmentVarsValue.(map[string]interface{})
				EnvironmentVars[i] = make(map[string]interface{})
				EnvironmentVars[i]["Key"] = EnvironmentVarsMap["key"]
				EnvironmentVars[i]["Value"] = EnvironmentVarsMap["value"]
				EnvironmentVars[i]["FieldRef.FieldPath"] = EnvironmentVarsMap["field_ref_field_path"]
			}
			InitContainers[i]["InitContainerEnvironmentVar"] = EnvironmentVars

			Ports := make([]map[string]interface{}, len(InitContainersMap["ports"].(*schema.Set).List()))
			for i, PortsValue := range InitContainersMap["ports"].(*schema.Set).List() {
				PortsMap := PortsValue.(map[string]interface{})
				Ports[i] = make(map[string]interface{})
				Ports[i]["Port"] = PortsMap["port"]
				Ports[i]["Protocol"] = PortsMap["protocol"]
			}
			InitContainers[i]["InitContainerPort"] = Ports

			VolumeMounts := make([]map[string]interface{}, len(InitContainersMap["volume_mounts"].(*schema.Set).List()))
			for i, VolumeMountsValue := range InitContainersMap["volume_mounts"].(*schema.Set).List() {
				VolumeMountsMap := VolumeMountsValue.(map[string]interface{})
				VolumeMounts[i] = make(map[string]interface{})
				VolumeMounts[i]["MountPath"] = VolumeMountsMap["mount_path"]
				VolumeMounts[i]["Name"] = VolumeMountsMap["name"]
				VolumeMounts[i]["ReadOnly"] = VolumeMountsMap["read_only"]
				VolumeMounts[i]["SubPath"] = VolumeMountsMap["sub_path"]
				if VolumeMountsMap["mount_propagation"] != nil && VolumeMountsMap["mount_propagation"] != "" {
					VolumeMounts[i]["MountPropagation"] = VolumeMountsMap["mount_propagation"]
				}
			}
			InitContainers[i]["InitContainerVolumeMount"] = VolumeMounts
		}
		request["InitContainer"] = InitContainers
	}

	if d.HasChange("volumes") {
		update = true
		Volumes := make([]map[string]interface{}, len(d.Get("volumes").(*schema.Set).List()))
		for i, VolumesValue := range d.Get("volumes").(*schema.Set).List() {
			VolumesMap := VolumesValue.(map[string]interface{})
			Volumes[i] = make(map[string]interface{})
			ConfigFileVolumeConfigFileToPaths := make([]map[string]interface{}, len(VolumesMap["config_file_volume_config_file_to_paths"].(*schema.Set).List()))
			for i, ConfigFileVolumeConfigFileToPathsValue := range VolumesMap["config_file_volume_config_file_to_paths"].(*schema.Set).List() {
				ConfigFileVolumeConfigFileToPathsMap := ConfigFileVolumeConfigFileToPathsValue.(map[string]interface{})
				ConfigFileVolumeConfigFileToPaths[i] = make(map[string]interface{})
				ConfigFileVolumeConfigFileToPaths[i]["Content"] = ConfigFileVolumeConfigFileToPathsMap["content"]
				ConfigFileVolumeConfigFileToPaths[i]["Path"] = ConfigFileVolumeConfigFileToPathsMap["path"]
				if ConfigFileVolumeConfigFileToPathsMap["mode"] != nil && ConfigFileVolumeConfigFileToPathsMap["mode"] != 0 {
					ConfigFileVolumeConfigFileToPaths[i]["Mode"] = ConfigFileVolumeConfigFileToPathsMap["mode"]
				}
			}
			Volumes[i]["ConfigFileVolumeConfigFileToPath"] = ConfigFileVolumeConfigFileToPaths
			Volumes[i]["DiskVolume.DiskId"] = VolumesMap["disk_volume_disk_id"]
			Volumes[i]["HostPathVolume.Type"] = VolumesMap["host_path_volume_type"]
			Volumes[i]["EmptyDirVolume.SizeLimit"] = VolumesMap["empty_dir_volume_size_limit"]
			Volumes[i]["EmptyDirVolume.Medium"] = VolumesMap["empty_dir_volume_medium"]
			Volumes[i]["HostPathVolume.Path"] = VolumesMap["host_path_volume_path"]
			if VolumesMap["config_file_volume_default_mode"] != nil && VolumesMap["config_file_volume_default_mode"] != 0 {
				Volumes[i]["ConfigFileVolumeDefaultMode"] = VolumesMap["config_file_volume_default_mode"]
			}
			Volumes[i]["DiskVolume.FsType"] = VolumesMap["disk_volume_fs_type"]
			Volumes[i]["DiskVolume.DiskSize"] = VolumesMap["disk_volume_disk_size"]
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

	if d.HasChange("host_aliases") {
		update = true
		aliases := make([]map[string]interface{}, len(d.Get("host_aliases").(*schema.Set).List()))
		for i, value := range d.Get("host_aliases").(*schema.Set).List() {
			aliasMap := value.(map[string]interface{})
			aliases[i] = make(map[string]interface{})
			aliases[i]["Hostname"] = aliasMap["hostnames"]
			aliases[i]["Ip"] = aliasMap["ip"]
		}
		request["HostAliase"] = aliases
	}

	if d.HasChange("instance_types") {
		update = true
		instanceTypes := d.Get("instance_types").([]interface{})
		types := make([]string, 0, int(5))
		if instanceTypes != nil && len(instanceTypes) > 0 {
			types = expandStringList(instanceTypes)
		}
		request["InstanceType"] = types
	}
	if update {
		_, err = client.RpcPost("Ess", "2014-08-28", "ModifyEciScalingConfiguration", nil, request, false)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_ess_eci_scaling_configuration", "ModifyEciScalingConfiguration", AlibabaCloudSdkGoERROR)
		}
	}
	d.Partial(false)

	return resourceAliyunEssEciScalingConfigurationRead(d, meta)
}

func resourceAliyunEssEciScalingConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	var err error

	o, err := essService.DescribeEssEciScalingConfiguration(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}

	if d.Get("force_delete").(bool) {
		request := ess.CreateDeleteScalingGroupRequest()
		request.ScalingGroupId = o["ScalingGroupId"].(string)
		request.ForceDelete = requests.NewBoolean(true)
		request.RegionId = client.RegionId
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DeleteScalingGroup(request)
		})

		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidScalingGroupId.NotFound"}) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return WrapError(essService.WaitForEssScalingGroup(d.Id(), Deleted, DefaultTimeout))
	} else {
		request := map[string]interface{}{
			"ScalingConfigurationId": d.Id(),
			"RegionId":               client.RegionId,
		}
		_, err = client.RpcPost("Ess", "2014-08-28", "DeleteEciScalingConfiguration", nil, request, false)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteEciScalingConfiguration", AlibabaCloudSdkGoERROR)
		}
		return nil
	}

}
