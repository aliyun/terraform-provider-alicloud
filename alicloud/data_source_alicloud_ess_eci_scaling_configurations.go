package alicloud

import (
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
				ValidateFunc: validation.ValidateRegexp,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
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
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scaling_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": tagsSchema(),
						"restart_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"container_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spot_price_limit": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"egress_bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"auto_create_eip": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"lifecycle_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eip_bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ram_role_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ingress_bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spot_strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"containers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(1),
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
										Type:         schema.TypeInt,
										ValidateFunc: validation.IntInSlice([]int{1}),
										Optional:     true,
									},
									"liveness_probe_timeout_seconds": {
										Type:         schema.TypeInt,
										ValidateFunc: validation.IntAtLeast(1),
										Optional:     true,
									},
									"readiness_probe_exec_commands": {
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Optional: true,
									},
									"readiness_probe_period_seconds": {
										Type:         schema.TypeInt,
										ValidateFunc: validation.IntAtLeast(1),
										Optional:     true,
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
										Type:         schema.TypeInt,
										ValidateFunc: validation.IntInSlice([]int{1}),
										Optional:     true,
									},
									"readiness_probe_timeout_seconds": {
										Type:         schema.TypeInt,
										ValidateFunc: validation.IntAtLeast(1),
										Optional:     true,
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
											},
										},
									},
									"disk_volume_disk_id": {
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

						"init_containers": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEssEciScalingConfigurationsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DescribeEciScalingConfigurations"
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	if scalingGroupId, ok := d.GetOk("scaling_group_id"); ok && scalingGroupId.(string) != "" {
		request["ScalingGroupId"] = scalingGroupId.(string)
	}

	var configurationsResponse []interface{}
	var response map[string]interface{}
	conn, err := client.NewEssClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-28"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_ess_scaling_configurations", "DescribeEciScalingConfigurations", AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.ScalingConfigurations", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ScalingConfigurations", response)
		}
		result, _ := resp.([]interface{})
		if len(result) < 1 {
			break
		}
		configurationsResponse = append(configurationsResponse, result...)
		if len(result) < PageSizeLarge {
			break
		}

		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	var filteredScalingConfigurations = make([]map[string]interface{}, 0)

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
		for _, configuration := range configurationsResponse {
			if okNameRegex && nameRegex != "" {
				r, err := regexp.Compile(nameRegex.(string))
				if err != nil {
					return WrapError(err)
				}
				if r != nil && !r.MatchString(configuration.(map[string]interface{})["ScalingConfigurationName"].(string)) {
					continue
				}
			}
			if okIds && len(idsMap) > 0 {
				if _, ok := idsMap[configuration.(map[string]interface{})["ScalingConfigurationId"].(string)]; !ok {
					continue
				}
			}
			filteredScalingConfigurations = append(filteredScalingConfigurations, configuration.(map[string]interface{}))
		}
	} else {
		for _, configuration := range configurationsResponse {
			filteredScalingConfigurations = append(filteredScalingConfigurations, configuration.(map[string]interface{}))
		}
	}

	return eciScalingConfigurationsDescriptionAttribute(d, filteredScalingConfigurations, meta)
}

func eciScalingConfigurationsDescriptionAttribute(d *schema.ResourceData, scalingConfigurations []map[string]interface{}, meta interface{}) error {
	var ids []string
	var names []string
	var s = make([]map[string]interface{}, 0)

	for _, scalingConfiguration := range scalingConfigurations {
		var mapping = map[string]interface{}{
			"id":                   scalingConfiguration["ScalingConfigurationId"],
			"name":                 scalingConfiguration["ScalingConfigurationName"],
			"scaling_group_id":     scalingConfiguration["ScalingGroupId"],
			"restart_policy":       scalingConfiguration["RestartPolicy"],
			"security_group_id":    scalingConfiguration["SecurityGroupId"],
			"container_group_name": scalingConfiguration["ContainerGroupName"],
			"resource_group_id":    scalingConfiguration["ResourceGroupId"],
			"description":          scalingConfiguration["Description"],
			"dns_policy":           scalingConfiguration["DnsPolicy"],
			"egress_bandwidth":     scalingConfiguration["EgressBandwidth"],
			"spot_price_limit":     scalingConfiguration["SpotPriceLimit"],
			"auto_create_eip":      scalingConfiguration["AutoCreateEip"],
			"memory":               scalingConfiguration["Memory"],
			"lifecycle_state":      scalingConfiguration["LifecycleState"],
			"creation_time":        scalingConfiguration["CreationTime"],
			"eip_bandwidth":        scalingConfiguration["EipBandwidth"],
			"ram_role_name":        scalingConfiguration["RamRoleName"],
			"ingress_bandwidth":    scalingConfiguration["IngressBandwidth"],
			"host_name":            scalingConfiguration["HostName"],
			"spot_strategy":        scalingConfiguration["SpotStrategy"],
			"cpu":                  scalingConfiguration["Cpu"],
		}
		if containersList, ok := scalingConfiguration["Containers"].([]interface{}); ok {
			mapping["containers"] = flattenContainerMappings(containersList)
		}

		tags := make(map[string]interface{}, 0)
		if tagList, ok := scalingConfiguration["Tags"].([]interface{}); ok {
			for _, v := range tagList {
				for key, value := range v.(map[string]interface{}) {
					tags[key] = value
				}
			}
			mapping["tags"] = tags
		}

		if initContainersList, ok := scalingConfiguration["InitContainers"].([]interface{}); ok {
			mapping["init_containers"] = flattenInitContainerMappings(initContainersList)
		}

		if volumeList, ok := scalingConfiguration["Volumes"].([]interface{}); ok {
			mapping["volumes"] = flattenVloumeMappings(volumeList)
		}

		hostAliases := make([]map[string]interface{}, 0)
		if hostAliasesList, ok := scalingConfiguration["HostAliases"].([]interface{}); ok {
			for _, v := range hostAliasesList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"hostnames": m1["Hostnames"],
						"ip":        m1["Ip"],
					}
					hostAliases = append(hostAliases, temp1)
				}
			}
			mapping["host_aliases"] = hostAliases
		}

		credentials := make([]map[string]interface{}, 0)
		if credentialList, ok := scalingConfiguration["ImageRegistryCredentials"].([]interface{}); ok {
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
			mapping["image_registry_credentials"] = credentials
		}

		acrRegistryInfos := make([]map[string]interface{}, 0)
		if acrRegistryInfoList, ok := scalingConfiguration["AcrRegistryInfos"].([]interface{}); ok {
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
			mapping["acr_registry_infos"] = acrRegistryInfos
		}

		ids = append(ids, scalingConfiguration["ScalingConfigurationId"].(string))
		names = append(names, scalingConfiguration["ScalingConfigurationName"].(string))
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("configurations", s); err != nil {
		return WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

func flattenInitContainerMappings(initContainerList []interface{}) []map[string]interface{} {
	initContainers := make([]map[string]interface{}, 0)

	for _, v := range initContainerList {
		if m1, ok := v.(map[string]interface{}); ok {
			temp1 := map[string]interface{}{
				"working_dir":       m1["WorkingDir"],
				"args":              m1["InitContainerArgs"],
				"cpu":               m1["Cpu"],
				"gpu":               m1["Gpu"],
				"memory":            m1["Memory"],
				"image":             m1["Image"],
				"image_pull_policy": m1["ImagePullPolicy"],
				"name":              m1["Name"],
				"commands":          m1["InitContainerCommands"],
			}
			if m1["InitContainerEnvironmentVars"] != nil {
				environmentVarsMaps := make([]map[string]interface{}, 0)
				for _, environmentVarsValue := range m1["InitContainerEnvironmentVars"].([]interface{}) {
					environmentVars := environmentVarsValue.(map[string]interface{})
					environmentVarsMap := map[string]interface{}{
						"key":   environmentVars["Key"],
						"value": environmentVars["Value"],
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
					}
					volumeMountsMaps = append(volumeMountsMaps, volumeMountsMap)
				}
				temp1["volume_mounts"] = volumeMountsMaps
			}
			initContainers = append(initContainers, temp1)

		}
	}
	return initContainers
}

func flattenVloumeMappings(volumesList []interface{}) []map[string]interface{} {
	volumes := make([]map[string]interface{}, 0)
	for _, v := range volumesList {
		if m1, ok := v.(map[string]interface{}); ok {
			temp1 := map[string]interface{}{
				"disk_volume_disk_id":   m1["DiskVolumeDiskId"],
				"disk_volume_fs_type":   m1["DiskVolumeFsType"],
				"disk_volume_disk_size": m1["DiskVolumeDiskSize"],
				"flex_volume_driver":    m1["FlexVolumeDriver"],
				"flex_volume_fs_type":   m1["FlexVolumeFsType"],
				"flex_volume_options":   m1["FlexVolumeOptions"],
				"nfs_volume_path":       m1["NFSVolumePath"],
				"nfs_volume_read_only":  m1["NFSVolumeReadOnly"],
				"nfs_volume_server":     m1["NFSVolumeServer"],
				"name":                  m1["Name"],
				"type":                  m1["Type"],
			}
			if m1["ConfigFileVolumeConfigFileToPaths"] != nil {
				configFileVolumeConfigFileToPathsMaps := make([]map[string]interface{}, 0)
				for _, configFileVolumeConfigFileToPathsValue := range m1["ConfigFileVolumeConfigFileToPaths"].([]interface{}) {
					configFileVolumeConfigFileToPaths := configFileVolumeConfigFileToPathsValue.(map[string]interface{})
					configFileVolumeConfigFileToPathsMap := map[string]interface{}{
						"content": configFileVolumeConfigFileToPaths["Content"],
						"path":    configFileVolumeConfigFileToPaths["Path"],
					}
					configFileVolumeConfigFileToPathsMaps = append(configFileVolumeConfigFileToPathsMaps, configFileVolumeConfigFileToPathsMap)
				}
				temp1["config_file_volume_config_file_to_paths"] = configFileVolumeConfigFileToPathsMaps
			}
			volumes = append(volumes, temp1)

		}
	}
	return volumes
}

func flattenContainerMappings(containerList []interface{}) []map[string]interface{} {
	containers := make([]map[string]interface{}, 0)
	for _, v := range containerList {
		if m1, ok := v.(map[string]interface{}); ok {
			temp1 := map[string]interface{}{
				"working_dir":       m1["WorkingDir"],
				"args":              m1["Args"],
				"cpu":               m1["Cpu"],
				"gpu":               m1["Gpu"],
				"memory":            m1["Memory"],
				"name":              m1["Name"],
				"image":             m1["Image"],
				"image_pull_policy": m1["ImagePullPolicy"],
				"commands":          m1["Commands"],

				"readiness_probe_exec_commands":         m1["ReadinessProbeExecCommands"],
				"readiness_probe_period_seconds":        m1["ReadinessProbePeriodSeconds"],
				"readiness_probe_http_get_path":         m1["ReadinessProbeHttpGetPath"],
				"readiness_probe_failure_threshold":     m1["ReadinessProbeFailureThreshold"],
				"readiness_probe_initial_delay_seconds": m1["ReadinessProbeInitialDelaySeconds"],
				"readiness_probe_http_get_port":         m1["ReadinessProbeHttpGetPort"],
				"readiness_probe_http_get_scheme":       m1["ReadinessProbeHttpGetScheme"],
				"readiness_probe_tcp_socket_port":       m1["ReadinessProbeTcpSocketPort"],
				"readiness_probe_success_threshold":     m1["ReadinessProbeSuccessThreshold"],
				"readiness_probe_timeout_seconds":       m1["ReadinessProbeTimeoutSeconds"],

				"liveness_probe_exec_commands":         m1["LivenessProbeExecCommands"],
				"liveness_probe_period_seconds":        m1["LivenessProbePeriodSeconds"],
				"liveness_probe_http_get_path":         m1["LivenessProbeHttpGetPath"],
				"liveness_probe_failure_threshold":     m1["LivenessProbeFailureThreshold"],
				"liveness_probe_initial_delay_seconds": m1["LivenessProbeInitialDelaySeconds"],
				"liveness_probe_http_get_port":         m1["LivenessProbeHttpGetPort"],
				"liveness_probe_http_get_scheme":       m1["LivenessProbeHttpGetScheme"],
				"liveness_probe_tcp_socket_port":       m1["LivenessProbeTcpSocketPort"],
				"liveness_probe_success_threshold":     m1["LivenessProbeSuccessThreshold"],
				"liveness_probe_timeout_seconds":       m1["LivenessProbeTimeoutSeconds"],
			}
			if m1["EnvironmentVars"] != nil {
				environmentVarsMaps := make([]map[string]interface{}, 0)
				for _, environmentVarsValue := range m1["EnvironmentVars"].([]interface{}) {
					environmentVars := environmentVarsValue.(map[string]interface{})
					environmentVarsMap := map[string]interface{}{
						"key":   environmentVars["Key"],
						"value": environmentVars["Value"],
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
			containers = append(containers, temp1)
		}
	}
	return containers
}
