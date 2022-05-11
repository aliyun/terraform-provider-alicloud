package alicloud

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
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
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[\\u4e00-\\u9fa5a-zA-Z0-9][\\u4e00-\\u9fa5a-zA-Z0-9\-_.]{1,63}$`), "It must be 2 to 64 characters in length and can contain letters, digits, underscores (_), hyphens (-), and periods (.). It must start with a letter or a digit."),
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
			"enable_sls": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ram_role_name": {
				Type:     schema.TypeString,
				Optional: true,
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
			"tags": tagsSchema(),
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
			"containers": {
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
		},
	}
}

func resourceAliyunEssEciScalingConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	// Ensure instance_type is generation three
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateEciScalingConfiguration"
	request := make(map[string]interface{})
	conn, err := client.NewEssClient()
	if err != nil {
		return WrapError(err)
	}
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

	Containers := make([]map[string]interface{}, len(d.Get("containers").(*schema.Set).List()))
	for i, v := range d.Get("containers").(*schema.Set).List() {
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
		VolumeMounts := make([]map[string]interface{}, len(ContainersMap["volume_mounts"].(*schema.Set).List()))
		for i, VolumeMountsValue := range ContainersMap["volume_mounts"].(*schema.Set).List() {
			VolumeMountsMap := VolumeMountsValue.(map[string]interface{})
			VolumeMounts[i] = make(map[string]interface{})
			VolumeMounts[i]["MountPath"] = VolumeMountsMap["mount_path"]
			VolumeMounts[i]["Name"] = VolumeMountsMap["name"]
			VolumeMounts[i]["ReadOnly"] = VolumeMountsMap["read_only"]
		}
		Containers[i]["VolumeMount"] = VolumeMounts
		Containers[i]["Command"] = ContainersMap["commands"]
	}
	request["Container"] = Containers

	if v, ok := d.GetOk("init_containers"); ok {
		InitContainers := make([]map[string]interface{}, len(v.(*schema.Set).List()))
		for i, v := range v.(*schema.Set).List() {
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
			VolumeMounts := make([]map[string]interface{}, len(InitContainersMap["volume_mounts"].(*schema.Set).List()))
			for i, VolumeMountsValue := range InitContainersMap["volume_mounts"].(*schema.Set).List() {
				VolumeMountsMap := VolumeMountsValue.(map[string]interface{})
				VolumeMounts[i] = make(map[string]interface{})
				VolumeMounts[i]["MountPath"] = VolumeMountsMap["mount_path"]
				VolumeMounts[i]["Name"] = VolumeMountsMap["name"]
				VolumeMounts[i]["ReadOnly"] = VolumeMountsMap["read_only"]
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
			}
			Volumes[i]["ConfigFileVolumeConfigFileToPath"] = ConfigFileVolumeConfigFileToPaths
			Volumes[i]["DiskVolume.DiskId"] = VolumesMap["disk_volume_disk_id"]
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

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-28"), StringPointer("AK"), nil, request, &runtime)
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
	d.Set("ram_role_name", o["RamRoleName"])
	d.Set("auto_create_eip", o["AutoCreateEip"])
	d.Set("eip_bandwidth", o["EipBandwidth"])
	d.Set("host_name", o["HostName"])
	d.Set("ingress_bandwidth", o["IngressBandwidth"])
	d.Set("egress_bandwidth", o["EgressBandwidth"])
	d.Set("tags", o["Tags"])
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

	containers := make([]map[string]interface{}, 0)
	if containersList, ok := o["Containers"].([]interface{}); ok {
		for _, v := range containersList {
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
	}

	if err := d.Set("containers", containers); err != nil {
		return WrapError(err)
	}

	initContainers := make([]map[string]interface{}, 0)
	if initContainersList, ok := o["InitContainers"].([]interface{}); ok {
		for _, v := range initContainersList {
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
	}
	if err := d.Set("init_containers", initContainers); err != nil {
		return WrapError(err)
	}

	volumes := make([]map[string]interface{}, 0)
	if volumesList, ok := o["Volumes"].([]interface{}); ok {
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
	if d.HasChange("enable_sls") {
		update = true
		request["EnableSls"] = d.Get("enable_sls")
	}
	if d.HasChange("ram_role_name") {
		update = true
		request["RamRoleName"] = d.Get("ram_role_name")
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

	if d.HasChange("containers") {
		update = true
		Containers := make([]map[string]interface{}, len(d.Get("containers").(*schema.Set).List()))
		for i, ContainersValue := range d.Get("containers").(*schema.Set).List() {
			ContainersMap := ContainersValue.(map[string]interface{})
			Containers[i] = make(map[string]interface{})
			Containers[i]["WorkingDir"] = ContainersMap["working_dir"]
			Containers[i]["Arg"] = ContainersMap["args"]
			Containers[i]["Cpu"] = ContainersMap["cpu"]
			Containers[i]["Gpu"] = ContainersMap["gpu"]
			Containers[i]["Memory"] = ContainersMap["memory"]
			Containers[i]["Name"] = ContainersMap["name"]
			Containers[i]["Image"] = ContainersMap["image"]
			Containers[i]["ImagePullPolicy"] = ContainersMap["image_pull_policy"]
			Containers[i]["Command"] = ContainersMap["commands"]

			EnvironmentVars := make([]map[string]interface{}, len(ContainersMap["environment_vars"].(*schema.Set).List()))
			for i, EnvironmentVarsValue := range ContainersMap["environment_vars"].(*schema.Set).List() {
				EnvironmentVarsMap := EnvironmentVarsValue.(map[string]interface{})
				EnvironmentVars[i] = make(map[string]interface{})
				EnvironmentVars[i]["Key"] = EnvironmentVarsMap["key"]
				EnvironmentVars[i]["Value"] = EnvironmentVarsMap["value"]
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
			}
			Containers[i]["VolumeMount"] = VolumeMounts
		}
		request["Container"] = Containers
	}

	if d.HasChange("init_containers") {
		update = true
		InitContainers := make([]map[string]interface{}, len(d.Get("init_containers").(*schema.Set).List()))
		for i, InitContainersValue := range d.Get("init_containers").(*schema.Set).List() {
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

			EnvironmentVars := make([]map[string]interface{}, len(InitContainersMap["environment_vars"].(*schema.Set).List()))
			for i, EnvironmentVarsValue := range InitContainersMap["environment_vars"].(*schema.Set).List() {
				EnvironmentVarsMap := EnvironmentVarsValue.(map[string]interface{})
				EnvironmentVars[i] = make(map[string]interface{})
				EnvironmentVars[i]["Key"] = EnvironmentVarsMap["key"]
				EnvironmentVars[i]["Value"] = EnvironmentVarsMap["value"]
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
			}
			Volumes[i]["ConfigFileVolumeConfigFileToPath"] = ConfigFileVolumeConfigFileToPaths
			Volumes[i]["DiskVolume.DiskId"] = VolumesMap["disk_volume_disk_id"]
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

	if update {
		conn, err := client.NewEssClient()
		if err != nil {
			return WrapError(err)
		}
		_, err = conn.DoRequest(StringPointer("ModifyEciScalingConfiguration"), nil, StringPointer("POST"), StringPointer("2014-08-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	conn, err := client.NewEssClient()
	if err != nil {
		return WrapError(err)
	}

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
		_, err = conn.DoRequest(StringPointer("DeleteEciScalingConfiguration"), nil, StringPointer("POST"), StringPointer("2014-08-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteEciScalingConfiguration", AlibabaCloudSdkGoERROR)
		}
		return nil
	}

}
