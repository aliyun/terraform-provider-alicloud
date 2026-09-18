package alicloud

import (
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudEssEciScalingConfigurations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEssEciScalingConfigurationsRead,
		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"configurations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scaling_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scaling_configuration_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"container_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"restart_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cost_optimization": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable_sls": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"instance_family_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_snapshot_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ram_role_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"termination_grace_period_seconds": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"auto_match_image_cache": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ipv6_address_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cpu_options_core": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cpu_options_threads_per_core": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"active_deadline_seconds": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"spot_strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spot_price_limit": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"auto_create_eip": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"eip_bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ingress_bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"egress_bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ephemeral_storage": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"load_balancer_weight": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"tags": tagsSchemaComputed(),
						"instance_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lifecycle_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"acr_registry_infos": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domains": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"instance_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"image_registry_credentials": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"password": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"server": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"username": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"dns_config_options": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"security_context_sysctls": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"containers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"security_context_capability_adds": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"lifecycle_pre_stop_handler_execs": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"security_context_read_only_root_file_system": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"tty": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"stdin": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"security_context_run_as_user": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ports": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"protocol": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"port": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"environment_vars": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"field_ref_field_path": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"working_dir": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"args": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"cpu": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"gpu": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"memory": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"image": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"image_pull_policy": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"volume_mounts": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mount_path": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"mount_propagation": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"sub_path": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"read_only": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"commands": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"liveness_probe_exec_commands": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"liveness_probe_period_seconds": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"liveness_probe_http_get_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"liveness_probe_failure_threshold": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"liveness_probe_initial_delay_seconds": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"liveness_probe_http_get_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"liveness_probe_http_get_scheme": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"liveness_probe_tcp_socket_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"liveness_probe_success_threshold": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"liveness_probe_timeout_seconds": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"readiness_probe_exec_commands": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"readiness_probe_period_seconds": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"readiness_probe_http_get_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"readiness_probe_failure_threshold": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"readiness_probe_initial_delay_seconds": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"readiness_probe_http_get_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"readiness_probe_http_get_scheme": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"readiness_probe_tcp_socket_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"readiness_probe_success_threshold": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"readiness_probe_timeout_seconds": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"init_containers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"security_context_capability_adds": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"security_context_read_only_root_file_system": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"security_context_run_as_user": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ports": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"protocol": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"port": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"environment_vars": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"field_ref_field_path": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"working_dir": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"args": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"cpu": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"gpu": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"memory": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"image": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"image_pull_policy": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"volume_mounts": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mount_path": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"mount_propagation": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"sub_path": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"read_only": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"commands": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"volumes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config_file_volume_config_file_to_paths": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"content": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"path": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"mode": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"disk_volume_disk_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"host_path_volume_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"host_path_volume_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"config_file_volume_default_mode": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"empty_dir_volume_medium": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"empty_dir_volume_size_limit": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"disk_volume_fs_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"disk_volume_disk_size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"flex_volume_driver": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"flex_volume_fs_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"flex_volume_options": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"nfs_volume_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"nfs_volume_read_only": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"nfs_volume_server": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"host_aliases": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hostnames": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
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

func dataSourceAlicloudEssEciScalingConfigurationsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"PageSize":   requests.NewInteger(PageSizeLarge),
		"PageNumber": requests.NewInteger(1),
		"RegionId":   client.RegionId,
	}

	if v, ok := d.GetOk("scaling_group_id"); ok {
		request["ScalingGroupId"] = v.(string)
	}

	var allConfigs []interface{}
	for {
		response, err = client.RpcPost("Ess", "2014-08-28", "DescribeEciScalingConfigurations", nil, request, true)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ess_eci_scaling_configurations", "DescribeEciScalingConfigurations", AlibabaCloudSdkGoERROR)
		}

		v, err := jsonpath.Get("$.ScalingConfigurations", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, "$.ScalingConfigurations", response)
		}

		addDebug("DescribeEciScalingConfigurations", response, request, request)
		if len(v.([]interface{})) < 1 {
			break
		}

		allConfigs = append(allConfigs, v.([]interface{})...)

		if len(v.([]interface{})) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(requests.Integer(fmt.Sprint(request["PageNumber"]))); err != nil {
			return WrapError(err)
		} else {
			request["PageNumber"] = page
		}
	}

	var filteredConfigsTemp []interface{}

	nameRegex, okNameRegex := d.GetOk("name_regex")
	idsMap := make(map[string]string)
	ids, okIds := d.GetOk("ids")
	if okIds {
		for _, i := range ids.([]interface{}) {
			if i == nil {
				continue
			}
			idsMap[i.(string)] = i.(string)
		}
	}
	if okNameRegex || okIds {
		for _, config := range allConfigs {
			var object map[string]interface{}
			object = config.(map[string]interface{})
			if okNameRegex && nameRegex != "" {
				r, err := regexp.Compile(nameRegex.(string))
				if err != nil {
					return WrapError(err)
				}
				if r != nil && !r.MatchString(object["ScalingConfigurationName"].(string)) {
					continue
				}
			}
			if okIds && len(idsMap) > 0 {
				if _, ok := idsMap[object["ScalingConfigurationId"].(string)]; !ok {
					continue
				}
			}
			filteredConfigsTemp = append(filteredConfigsTemp, config)
		}
	} else {
		filteredConfigsTemp = allConfigs
	}
	return eciScalingConfigurationsDescriptionAttribute(d, filteredConfigsTemp, meta, client)
}

func eciScalingConfigurationsDescriptionAttribute(d *schema.ResourceData, configs []interface{}, meta interface{}, client *connectivity.AliyunClient) error {
	var ids []string
	var names []string
	var s = make([]map[string]interface{}, 0)

	for _, config := range configs {
		var object map[string]interface{}
		object = config.(map[string]interface{})

		tagsMap := make(map[string]interface{})
		if object["Tags"] != nil {
			if tags, ok := object["Tags"].(map[string]interface{}); ok {
				// The API may wrap the tag list under a "Tag" key, or return a plain key/value map.
				if tagList, ok := tags["Tag"]; ok {
					tagsMap = tagsToMap(tagList)
				} else {
					tagsMap = tagsToMap(tags)
				}
			} else {
				// The API may also return the tags directly as a list of {Key, Value}.
				tagsMap = tagsToMap(object["Tags"])
			}
		}

		mapping := map[string]interface{}{
			"id":                               object["ScalingConfigurationId"],
			"scaling_group_id":                 object["ScalingGroupId"],
			"scaling_configuration_name":       object["ScalingConfigurationName"],
			"description":                      object["Description"],
			"security_group_id":                object["SecurityGroupId"],
			"container_group_name":             object["ContainerGroupName"],
			"restart_policy":                   object["RestartPolicy"],
			"cpu":                              object["Cpu"],
			"memory":                           object["Memory"],
			"resource_group_id":                object["ResourceGroupId"],
			"dns_policy":                       object["DnsPolicy"],
			"enable_sls":                       object["SlsEnable"],
			"cost_optimization":                object["CostOptimization"],
			"image_snapshot_id":                object["ImageSnapshotId"],
			"instance_family_level":            object["InstanceFamilyLevel"],
			"ram_role_name":                    object["RamRoleName"],
			"termination_grace_period_seconds": object["TerminationGracePeriodSeconds"],
			"auto_match_image_cache":           object["AutoMatchImageCache"],
			"ipv6_address_count":               object["Ipv6AddressCount"],
			"cpu_options_core":                 object["CpuOptionsCore"],
			"cpu_options_threads_per_core":     object["CpuOptionsThreadsPerCore"],
			"active_deadline_seconds":          object["ActiveDeadlineSeconds"],
			"spot_strategy":                    object["SpotStrategy"],
			"auto_create_eip":                  object["AutoCreateEip"],
			"eip_bandwidth":                    object["EipBandwidth"],
			"host_name":                        object["HostName"],
			"ingress_bandwidth":                object["IngressBandwidth"],
			"egress_bandwidth":                 object["EgressBandwidth"],
			"ephemeral_storage":                object["EphemeralStorage"],
			"load_balancer_weight":             object["LoadBalancerWeight"],
			"tags":                             tagsMap,
			"instance_types":                   object["InstanceTypes"],
			"creation_time":                    object["CreationTime"],
			"lifecycle_state":                  object["LifecycleState"],
		}

		if object["SpotPriceLimit"] != nil {
			mapping["spot_price_limit"] = object["SpotPriceLimit"]
		}

		credentials := make([]map[string]interface{}, 0)
		if object["ImageRegistryCredentials"] != nil && len(object["ImageRegistryCredentials"].([]interface{})) != 0 {
			if credentialList, ok := object["ImageRegistryCredentials"].([]interface{}); ok {
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
		}
		mapping["image_registry_credentials"] = credentials

		options := make([]map[string]interface{}, 0)
		if optionList, ok := object["DnsConfigOptions"].([]interface{}); ok {
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
		mapping["dns_config_options"] = options

		sysctls := make([]map[string]interface{}, 0)
		if sysctlList, ok := object["SecurityContextSysCtls"].([]interface{}); ok {
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
		mapping["security_context_sysctls"] = sysctls

		acrRegistryInfos := make([]map[string]interface{}, 0)
		if object["AcrRegistryInfos"] != nil && len(object["AcrRegistryInfos"].([]interface{})) != 0 {
			if acrRegistryInfoList, ok := object["AcrRegistryInfos"].([]interface{}); ok {
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
		}
		mapping["acr_registry_infos"] = acrRegistryInfos

		containers := make([]map[string]interface{}, 0)
		if containersList, ok := object["Containers"].([]interface{}); ok {
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
		mapping["containers"] = containers

		initContainers := make([]map[string]interface{}, 0)
		if initContainersList, ok := object["InitContainers"].([]interface{}); ok {
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
		mapping["init_containers"] = initContainers

		volumes := make([]map[string]interface{}, 0)
		if object["Volumes"] != nil && len(object["Volumes"].([]interface{})) != 0 {
			if volumesList, ok := object["Volumes"].([]interface{}); ok {
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
		}
		mapping["volumes"] = volumes

		hostAliases := make([]map[string]interface{}, 0)
		if object["HostAliases"] != nil && len(object["HostAliases"].([]interface{})) != 0 {
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
		}
		mapping["host_aliases"] = hostAliases

		ids = append(ids, object["ScalingConfigurationId"].(string))
		names = append(names, object["ScalingConfigurationName"].(string))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("configurations", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok {
		writeToFile(output.(string), s)
	}

	return nil
}
