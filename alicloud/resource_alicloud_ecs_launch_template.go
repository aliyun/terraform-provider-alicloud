package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudEcsLaunchTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEcsLaunchTemplateCreate,
		Read:   resourceAliCloudEcsLaunchTemplateRead,
		Update: resourceAliCloudEcsLaunchTemplateUpdate,
		Delete: resourceAliCloudEcsLaunchTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"auto_renew_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"auto_release_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_disks": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"delete_with_instance": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.All(StringLenBetween(2, 256), StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with \"http://\", \"https://\".")),
						},
						"encrypted": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.All(StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with \"http://\", \"https://\"."), StringMatch(regexp.MustCompile(`^[a-zA-Z\p{Han}][a-zA-Z\p{Han}_0-9\-\.\:]{1,127}$`), `It can contain A-Z, a-z, Chinese characters, numbers, periods (.), colons (:), underscores (_), and hyphens (-).`)),
						},
						"performance_level": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"device": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"deployment_set_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.All(StringLenBetween(2, 256), StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with \"http://\", \"https://\".")),
			},
			"enable_vm_os_config": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"host_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.All(StringDoesNotMatch(regexp.MustCompile(`(^\..*)|(^\-.*)|(.*\-$)|(.*\.$)`), "It cannot begin or end with period (.), hyphen (-).")),
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_owner_alias": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"system", "self", "others", "marketplace", ""}, false),
				Default:      "",
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"PostPaid", "PrePaid"}, false),
			},
			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.All(StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with \"http://\", \"https://\"."), StringMatch(regexp.MustCompile(`^[a-zA-Z\p{Han}][a-zA-Z\p{Han}_0-9\-\.\:\,\[\]]{1,127}$`), `It must begin with an English or a Chinese character. It can contain A-Z, a-z, Chinese characters, numbers, colons (:), underscores (_), periods (.), commas (,), brackets ([]), and hyphens (-).`)),
			},
			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile(`^ecs\..*`), "prefix must be 'ecs.'"),
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"PayByBandwidth", "PayByTraffic"}, false),
			},
			"internet_max_bandwidth_in": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"internet_max_bandwidth_out": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 100),
			},
			"io_optimized": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"none", "optimized"}, false),
			},
			"key_pair_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringLenBetween(2, 128),
			},
			"launch_template_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "Field 'name' has been deprecated from provider version 1.120.0. New field 'launch_template_name' instead.",
				ConflictsWith: []string{"launch_template_name"},
			},
			"network_interfaces": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.All(StringLenBetween(2, 256), StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with \"http://\", \"https://\".")),
						},
						"name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.All(StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with \"http://\", \"https://\"."), StringMatch(regexp.MustCompile(`^[a-zA-Z\p{Han}][a-zA-Z\p{Han}_0-9\-\.\:]{1,127}$`), `It can contain A-Z, a-z, Chinese characters, numbers, periods (.), colons (:), underscores (_), and hyphens (-).`)),
						},
						"primary_ip": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.SingleIP(),
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				MaxItems: 1,
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"classic", "vpc"}, false),
			},
			"password_inherit": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"period_unit": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"private_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ram_role_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_enhancement_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Active", "Deactive"}, false),
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"spot_duration": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"0", "1", "2", "3", "4", "5", "6"}, false),
				Default:      "1",
			},
			"spot_price_limit": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"spot_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"NoSpot", "SpotAsPriceGo", "SpotWithPriceLimit"}, false),
			},
			"system_disk": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ValidateFunc:  StringInSlice([]string{"all", "cloud", "ephemeral_ssd", "cloud_essd", "cloud_efficiency", "cloud_ssd", "local_disk"}, false),
							ConflictsWith: []string{"system_disk_category"},
						},
						"delete_with_instance": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"description": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ValidateFunc:  validation.All(StringLenBetween(2, 256), StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with \"http://\", \"https://\".")),
							ConflictsWith: []string{"system_disk_description"},
						},
						"iops": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ValidateFunc:  validation.All(StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with \"http://\", \"https://\"."), StringMatch(regexp.MustCompile(`^[a-zA-Z\p{Han}][a-zA-Z\p{Han}_0-9\-\.\:]{1,127}$`), `It can contain A-Z, a-z, Chinese characters, numbers, periods (.), colons (:), underscores (_), and hyphens (-).`)),
							ConflictsWith: []string{"system_disk_name"},
						},
						"performance_level": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
						},
						"size": {
							Type:          schema.TypeInt,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"system_disk_size"},
						},
						"encrypted": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"tags":          tagsSchema(),
			"template_tags": tagsSchema(),
			"template_resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"default_version_number": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"update_default_version_number"},
				ValidateFunc:  validation.IntAtLeast(1),
			},
			"update_default_version_number": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"default_version_number"},
			},
			"latest_version_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"user_data": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"userdata"},
			},
			"userdata": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'userdata' has been deprecated from provider version 1.120.0. New field 'user_data' instead.",
				ConflictsWith: []string{"user_data"},
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"http_endpoint": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"http_tokens": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"http_put_response_hop_limit": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"system_disk_category": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  StringInSlice([]string{"all", "cloud", "ephemeral_ssd", "cloud_essd", "cloud_efficiency", "cloud_ssd", "local_disk"}, false),
				Deprecated:    "Field 'system_disk_category' has been deprecated from provider version 1.120.0. New field 'system_disk' instead.",
				ConflictsWith: []string{"system_disk"},
			},
			"system_disk_description": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.All(StringLenBetween(2, 256), StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with \"http://\", \"https://\".")),
				Deprecated:    "Field 'system_disk_description' has been deprecated from provider version 1.120.0. New field 'system_disk' instead.",
				ConflictsWith: []string{"system_disk"},
			},
			"system_disk_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.All(StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with \"http://\", \"https://\"."), StringMatch(regexp.MustCompile(`^[a-zA-Z\p{Han}][a-zA-Z\p{Han}_0-9\-\.\:]{1,127}$`), `It can contain A-Z, a-z, Chinese characters, numbers, periods (.), colons (:), underscores (_), and hyphens (-).`)),
				Deprecated:    "Field 'system_disk_name' has been deprecated from provider version 1.120.0. New field 'system_disk' instead.",
				ConflictsWith: []string{"system_disk"},
			},
			"system_disk_size": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'system_disk_size' has been deprecated from provider version 1.120.0. New field 'system_disk' instead.",
				ConflictsWith: []string{"system_disk"},
			},
		},
	}
}

func resourceAliCloudEcsLaunchTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateLaunchTemplate"
	request := make(map[string]interface{})
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("auto_release_time"); ok {
		request["AutoReleaseTime"] = v
	}

	if v, ok := d.GetOk("data_disks"); ok {
		dataDisksMaps := make([]map[string]interface{}, 0)
		for _, dataDisks := range v.(*schema.Set).List() {
			dataDisksMap := make(map[string]interface{})
			dataDisksArg := dataDisks.(map[string]interface{})

			if category, ok := dataDisksArg["category"].(string); ok && category != "" {
				dataDisksMap["Category"] = category
			}

			if description, ok := dataDisksArg["description"].(string); ok && description != "" {
				dataDisksMap["Description"] = description
			}

			if name, ok := dataDisksArg["name"].(string); ok && name != "" {
				dataDisksMap["DiskName"] = name
			}

			if performanceLevel, ok := dataDisksArg["performance_level"].(string); ok && performanceLevel != "" {
				dataDisksMap["PerformanceLevel"] = performanceLevel
			}

			if snapshotId, ok := dataDisksArg["snapshot_id"].(string); ok && snapshotId != "" {
				dataDisksMap["SnapshotId"] = snapshotId
			}

			if device, ok := dataDisksArg["device"].(string); ok && device != "" {
				dataDisksMap["Device"] = device
			}

			dataDisksMap["DeleteWithInstance"] = requests.NewBoolean(dataDisksArg["delete_with_instance"].(bool))
			dataDisksMap["Encrypted"] = fmt.Sprintf("%v", dataDisksArg["encrypted"].(bool))
			dataDisksMap["Size"] = requests.NewInteger(dataDisksArg["size"].(int))

			dataDisksMaps = append(dataDisksMaps, dataDisksMap)
		}
		request["DataDisk"] = dataDisksMaps

	}

	if v, ok := d.GetOk("auto_renew"); ok {
		request["AutoRenew"] = v
	}

	if v, ok := d.GetOk("auto_renew_period"); ok {
		request["AutoRenewPeriod"] = v
	}

	if v, ok := d.GetOk("deployment_set_id"); ok {
		request["DeploymentSetId"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOkExists("enable_vm_os_config"); ok {
		request["EnableVmOsConfig"] = v
	}

	if v, ok := d.GetOk("host_name"); ok {
		request["HostName"] = v
	}

	if v, ok := d.GetOk("image_id"); ok {
		request["ImageId"] = v
	}

	if v, ok := d.GetOk("image_owner_alias"); ok {
		request["ImageOwnerAlias"] = v
	}

	if v, ok := d.GetOk("instance_charge_type"); ok {
		request["InstanceChargeType"] = v
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request["InstanceName"] = v
	}

	if v, ok := d.GetOk("instance_type"); ok {
		request["InstanceType"] = v
	}

	if v, ok := d.GetOk("internet_charge_type"); ok {
		request["InternetChargeType"] = v
	}

	if v, ok := d.GetOk("internet_max_bandwidth_in"); ok {
		request["InternetMaxBandwidthIn"] = v
	}

	if v, ok := d.GetOk("internet_max_bandwidth_out"); ok {
		request["InternetMaxBandwidthOut"] = v
	}

	if v, ok := d.GetOk("io_optimized"); ok {
		request["IoOptimized"] = v
	}

	if v, ok := d.GetOk("key_pair_name"); ok {
		request["KeyPairName"] = v
	}

	if v, ok := d.GetOk("launch_template_name"); ok {
		request["LaunchTemplateName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["LaunchTemplateName"] = v
	} else {
		return WrapError(Error(`[ERROR] Argument "name" or "launch_template_name" must be set one!`))
	}

	if v, ok := d.GetOk("network_interfaces"); ok {
		networkInterfacesMaps := make([]map[string]interface{}, 0)
		for _, networkInterfaces := range v.(*schema.Set).List() {
			networkInterfacesMap := make(map[string]interface{})
			networkInterfacesArg := networkInterfaces.(map[string]interface{})
			networkInterfacesMap["Description"] = networkInterfacesArg["description"]
			networkInterfacesMap["NetworkInterfaceName"] = networkInterfacesArg["name"]
			networkInterfacesMap["PrimaryIpAddress"] = networkInterfacesArg["primary_ip"]
			networkInterfacesMap["SecurityGroupId"] = networkInterfacesArg["security_group_id"]
			networkInterfacesMap["VSwitchId"] = networkInterfacesArg["vswitch_id"]
			networkInterfacesMaps = append(networkInterfacesMaps, networkInterfacesMap)
		}
		request["NetworkInterface"] = networkInterfacesMaps

	}

	if v, ok := d.GetOk("network_type"); ok {
		request["NetworkType"] = v
	}

	if v, ok := d.GetOkExists("password_inherit"); ok {
		request["PasswordInherit"] = v
	}

	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}

	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}

	if v, ok := d.GetOk("private_ip_address"); ok {
		request["PrivateIpAddress"] = v
	}

	if v, ok := d.GetOk("ram_role_name"); ok {
		request["RamRoleName"] = v
	}

	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("security_enhancement_strategy"); ok {
		request["SecurityEnhancementStrategy"] = v
	}

	if v, ok := d.GetOk("security_group_id"); ok {
		request["SecurityGroupId"] = v
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		request["SecurityGroupIds"] = v
	}

	if v, ok := d.GetOk("spot_duration"); ok {
		request["SpotDuration"] = v
	}

	if v, ok := d.GetOk("spot_price_limit"); ok {
		request["SpotPriceLimit"] = v
	}

	if v, ok := d.GetOk("spot_strategy"); ok {
		request["SpotStrategy"] = v
	}

	if v, ok := d.GetOk("system_disk"); ok {
		systemDiskMap := make(map[string]interface{})
		for _, systemDisk := range v.(*schema.Set).List() {
			systemDiskArg := systemDisk.(map[string]interface{})
			systemDiskMap["Category"] = systemDiskArg["category"]
			systemDiskMap["DeleteWithInstance"] = requests.NewBoolean(systemDiskArg["delete_with_instance"].(bool))
			systemDiskMap["Description"] = systemDiskArg["description"]
			systemDiskMap["Iops"] = systemDiskArg["iops"]
			systemDiskMap["DiskName"] = systemDiskArg["name"]
			systemDiskMap["PerformanceLevel"] = systemDiskArg["performance_level"]
			systemDiskMap["Size"] = requests.NewInteger(systemDiskArg["size"].(int))
			systemDiskMap["Encrypted"] = requests.NewBoolean(systemDiskArg["encrypted"].(bool))
		}
		request["SystemDisk"] = systemDiskMap

	} else {
		systemDiskMap := make(map[string]interface{})
		if v, ok := d.GetOk("system_disk_category"); ok {
			systemDiskMap["Category"] = v
		}
		if v, ok := d.GetOk("system_disk_description"); ok {
			systemDiskMap["Description"] = v
		}
		if v, ok := d.GetOk("system_disk_name"); ok {
			systemDiskMap["DiskName"] = v
		}
		if v, ok := d.GetOk("system_disk_size"); ok {
			systemDiskMap["Size"] = v
		}
		if len(systemDiskMap) > 0 {
			request["SystemDisk"] = systemDiskMap
		}

	}
	if v, ok := d.GetOk("tags"); ok {
		count := 1
		for key, value := range v.(map[string]interface{}) {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
	}
	if v, ok := d.GetOk("template_tags"); ok {
		count := 1
		for key, value := range v.(map[string]interface{}) {
			request[fmt.Sprintf("TemplateTag.%d.Key", count)] = key
			request[fmt.Sprintf("TemplateTag.%d.Value", count)] = value
			count++
		}
	}
	if v, ok := d.GetOk("template_resource_group_id"); ok {
		request["TemplateResourceGroupId"] = v
	}

	if v, ok := d.GetOk("user_data"); ok {
		request["UserData"] = v
	} else if v, ok := d.GetOk("userdata"); ok {
		request["UserData"] = v
	}

	if v, ok := d.GetOk("version_description"); ok {
		request["VersionDescription"] = v
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}

	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}

	vswitchId := Trim(d.Get("vswitch_id").(string))
	if vswitchId != "" {
		vpcService := VpcService{client}
		vsw, err := vpcService.DescribeVSwitchWithTeadsl(vswitchId)
		if err != nil {
			return WrapError(err)
		}
		if v, ok := request["VpcId"].(string); !ok || v == "" {
			request["VpcId"] = vsw["VpcId"]
		}
		request["VSwitchId"] = vswitchId
		if v, ok := request["ZoneId"].(string); !ok || v == "" {
			request["ZoneId"] = vsw["ZoneId"]
		}
	}
	if v, ok := d.GetOk("http_endpoint"); ok {
		request["HttpEndpoint"] = v
	}
	if v, ok := d.GetOk("http_tokens"); ok {
		request["HttpTokens"] = v
	}
	if v, ok := d.GetOk("http_put_response_hop_limit"); ok {
		request["HttpPutResponseHopLimit"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_launch_template", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["LaunchTemplateId"]))

	return resourceAliCloudEcsLaunchTemplateRead(d, meta)
}

func resourceAliCloudEcsLaunchTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeEcsLaunchTemplate(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_launch_template ecsService.DescribeEcsLaunchTemplate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("launch_template_name", object["LaunchTemplateName"])
	d.Set("name", object["LaunchTemplateName"])
	d.Set("default_version_number", object["DefaultVersionNumber"])
	d.Set("latest_version_number", object["LatestVersionNumber"])

	describeLaunchTemplateVersionsObject, err := ecsService.DescribeLaunchTemplateVersions(d.Id(), object["LatestVersionNumber"])
	if err != nil {
		return WrapError(err)
	}
	d.Set("auto_release_time", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["AutoReleaseTime"])

	dataDisk := make([]map[string]interface{}, 0)
	if dataDiskList, ok := describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["DataDisks"].(map[string]interface{})["DataDisk"].([]interface{}); ok {
		for _, v := range dataDiskList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"category":             m1["Category"],
					"delete_with_instance": m1["DeleteWithInstance"],
					"description":          m1["Description"],
					"encrypted":            m1["Encrypted"],
					"name":                 m1["DiskName"],
					"performance_level":    m1["PerformanceLevel"],
					"size":                 m1["Size"],
					"snapshot_id":          m1["SnapshotId"],
					"device":               m1["Device"],
				}
				dataDisk = append(dataDisk, temp1)

			}
		}
	}
	if err := d.Set("data_disks", dataDisk); err != nil {
		return WrapError(err)
	}
	d.Set("deployment_set_id", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["DeploymentSetId"])
	d.Set("description", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["Description"])
	d.Set("enable_vm_os_config", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["EnableVmOsConfig"])
	d.Set("host_name", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["HostName"])
	d.Set("image_id", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["ImageId"])
	d.Set("image_owner_alias", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["ImageOwnerAlias"])
	d.Set("instance_charge_type", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["InstanceChargeType"])
	d.Set("instance_name", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["InstanceName"])
	d.Set("instance_type", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["InstanceType"])
	d.Set("internet_charge_type", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["InternetChargeType"])
	d.Set("internet_max_bandwidth_in", formatInt(describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["InternetMaxBandwidthIn"]))
	d.Set("internet_max_bandwidth_out", formatInt(describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["InternetMaxBandwidthOut"]))
	d.Set("io_optimized", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["IoOptimized"])
	d.Set("key_pair_name", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["KeyPairName"])

	networkInterface := make([]map[string]interface{}, 0)
	if networkInterfaceList, ok := describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["NetworkInterfaces"].(map[string]interface{})["NetworkInterface"].([]interface{}); ok {
		for _, v := range networkInterfaceList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"description":       m1["Description"],
					"name":              m1["NetworkInterfaceName"],
					"primary_ip":        m1["PrimaryIpAddress"],
					"security_group_id": m1["SecurityGroupId"],
					"vswitch_id":        m1["VSwitchId"],
				}
				networkInterface = append(networkInterface, temp1)

			}
		}
	}
	if err := d.Set("network_interfaces", networkInterface); err != nil {
		return WrapError(err)
	}
	d.Set("auto_renew", formatBool(describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["AutoRenew"]))
	d.Set("auto_renew_period", formatInt(describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["AutoRenewPeriod"]))
	d.Set("network_type", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["NetworkType"])
	d.Set("password_inherit", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["PasswordInherit"])
	d.Set("period", formatInt(describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["Period"]))
	d.Set("period_unit", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["PeriodUnit"])
	d.Set("private_ip_address", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["PrivateIpAddress"])
	d.Set("ram_role_name", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["RamRoleName"])
	d.Set("security_enhancement_strategy", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SecurityEnhancementStrategy"])
	d.Set("security_group_id", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SecurityGroupId"])
	if describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SecurityGroupIds"] != nil {
		d.Set("security_group_ids", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SecurityGroupIds"].(map[string]interface{})["SecurityGroupId"])
	}
	d.Set("spot_duration", fmt.Sprint(formatInt(describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SpotDuration"])))
	d.Set("spot_price_limit", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SpotPriceLimit"])
	d.Set("spot_strategy", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SpotStrategy"])
	d.Set("system_disk_category", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SystemDisk.Category"])
	d.Set("system_disk_size", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SystemDisk.Size"])
	d.Set("system_disk_description", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SystemDisk.Description"])
	d.Set("system_disk_name", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SystemDisk.DiskName"])

	systemDiskSli := make([]map[string]interface{}, 0)
	systemDiskMap := make(map[string]interface{})
	systemDiskMap["category"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SystemDisk.Category"]
	systemDiskMap["delete_with_instance"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SystemDisk.DeleteWithInstance"]
	systemDiskMap["description"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SystemDisk.Description"]
	systemDiskMap["iops"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SystemDisk.Iops"]
	systemDiskMap["name"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SystemDisk.DiskName"]
	systemDiskMap["performance_level"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SystemDisk.PerformanceLevel"]
	systemDiskMap["size"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SystemDisk.Size"]
	systemDiskMap["encrypted"] = describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["SystemDisk.Encrypted"]
	systemDiskSli = append(systemDiskSli, systemDiskMap)
	d.Set("system_disk", systemDiskSli)

	d.Set("vpc_id", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["VpcId"])
	d.Set("resource_group_id", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["ResourceGroupId"])
	if object["Tags"] != nil {
		d.Set("template_tags", tagsToMap(object["Tags"].(map[string]interface{})["Tag"]))
	}

	d.Set("user_data", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["UserData"])
	d.Set("userdata", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["UserData"])
	d.Set("vswitch_id", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["VSwitchId"])
	d.Set("version_description", describeLaunchTemplateVersionsObject["VersionDescription"])
	d.Set("vpc_id", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["VpcId"])
	d.Set("zone_id", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["ZoneId"])
	d.Set("tags", tagsToMap(describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["Tags"].(map[string]interface{})["InstanceTag"]))
	d.Set("http_endpoint", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["HttpEndpoint"])
	d.Set("http_tokens", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["HttpTokens"])
	d.Set("http_put_response_hop_limit", describeLaunchTemplateVersionsObject["LaunchTemplateData"].(map[string]interface{})["HttpPutResponseHopLimit"])
	return nil
}

func resourceAliCloudEcsLaunchTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("template_tags") {
		if err := ecsService.SetResourceTemplateTags(d, "launchtemplate"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("template_tags")
	}
	versions, err := getLaunchTemplateVersions(d.Id(), meta)
	if err != nil {
		return WrapError(err)
	}
	// Remove one of the oldest and non-default version when the total number reach 30
	if len(versions) > 29 {
		var oldestVersion int
		for _, version := range versions {
			if !version.(map[string]interface{})["DefaultVersion"].(bool) && (oldestVersion == 0 || formatInt(version.(map[string]interface{})["VersionNumber"]) < oldestVersion) {
				oldestVersion = formatInt(version.(map[string]interface{})["VersionNumber"])
			}
		}

		err = deleteLaunchTemplateVersion(d.Id(), oldestVersion, meta)
		if err != nil {
			return WrapError(err)
		}
	}

	latestVersion := d.Get("latest_version_number")
	update := false
	request := map[string]interface{}{
		"LaunchTemplateId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("auto_release_time") {
		update = true
	}
	if v, ok := d.GetOk("auto_release_time"); ok {
		request["AutoReleaseTime"] = v
	}
	if d.HasChange("data_disks") {
		update = true
	}
	if v, ok := d.GetOk("data_disks"); ok {
		DataDisks := make([]map[string]interface{}, len(v.(*schema.Set).List()))
		for i, DataDisksValue := range v.(*schema.Set).List() {
			DataDisksMap := DataDisksValue.(map[string]interface{})
			DataDisks[i] = make(map[string]interface{})

			if category, ok := DataDisksMap["category"].(string); ok && category != "" {
				DataDisks[i]["Category"] = category
			}

			if description, ok := DataDisksMap["description"].(string); ok && description != "" {
				DataDisks[i]["Description"] = description
			}

			if name, ok := DataDisksMap["name"].(string); ok && name != "" {
				DataDisks[i]["DiskName"] = name
			}

			if performanceLevel, ok := DataDisksMap["performance_level"].(string); ok && performanceLevel != "" {
				DataDisks[i]["PerformanceLevel"] = performanceLevel
			}

			if snapshotId, ok := DataDisksMap["snapshot_id"].(string); ok && snapshotId != "" {
				DataDisks[i]["SnapshotId"] = snapshotId
			}

			if device, ok := DataDisksMap["device"].(string); ok && device != "" {
				DataDisks[i]["Device"] = device
			}

			DataDisks[i]["DeleteWithInstance"] = DataDisksMap["delete_with_instance"]
			DataDisks[i]["Encrypted"] = fmt.Sprintf("%v", DataDisksMap["encrypted"].(bool))
			DataDisks[i]["Size"] = DataDisksMap["size"]

		}
		request["DataDisk"] = DataDisks

	}
	if d.HasChange("deployment_set_id") {
		update = true
	}
	if v, ok := d.GetOk("deployment_set_id"); ok {
		request["DeploymentSetId"] = v
	}
	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if d.HasChange("enable_vm_os_config") {
		update = true
	}
	if v, ok := d.GetOkExists("enable_vm_os_config"); ok {
		request["EnableVmOsConfig"] = v
	}
	if d.HasChange("host_name") {
		update = true
	}
	if v, ok := d.GetOk("host_name"); ok {
		request["HostName"] = v
	}
	if d.HasChange("image_id") {
		update = true
	}
	if v, ok := d.GetOk("image_id"); ok {
		request["ImageId"] = v
	}
	if d.HasChange("image_owner_alias") {
		update = true
	}
	if v, ok := d.GetOk("image_owner_alias"); ok {
		request["ImageOwnerAlias"] = v
	}
	if d.HasChange("instance_charge_type") {
		update = true
	}
	if v, ok := d.GetOk("instance_charge_type"); ok {
		request["InstanceChargeType"] = v
	}
	if d.HasChange("instance_name") {
		update = true
	}
	if v, ok := d.GetOk("instance_name"); ok {
		request["InstanceName"] = v
	}
	if d.HasChange("instance_type") {
		update = true
	}
	if v, ok := d.GetOk("instance_type"); ok {
		request["InstanceType"] = v
	}
	if d.HasChange("internet_charge_type") {
		update = true
	}
	if v, ok := d.GetOk("internet_charge_type"); ok {
		request["InternetChargeType"] = v
	}
	if d.HasChange("internet_max_bandwidth_in") {
		update = true
	}
	if v, ok := d.GetOk("internet_max_bandwidth_in"); ok {
		request["InternetMaxBandwidthIn"] = v
	}
	if d.HasChange("internet_max_bandwidth_out") {
		update = true
	}
	if v, ok := d.GetOk("internet_max_bandwidth_out"); ok {
		request["InternetMaxBandwidthOut"] = v
	}
	if d.HasChange("io_optimized") {
		update = true
	}
	if v, ok := d.GetOk("io_optimized"); ok {
		request["IoOptimized"] = v
	}
	if d.HasChange("key_pair_name") {
		update = true
	}
	if v, ok := d.GetOk("key_pair_name"); ok {
		request["KeyPairName"] = v
	}

	if d.HasChange("launch_template_name") {
		update = true
		request["LaunchTemplateName"] = d.Get("launch_template_name")
	} else if d.HasChange("name") {
		update = true
		request["LaunchTemplateName"] = d.Get("name")
	} else {
		if v, ok := d.GetOk("launch_template_name"); ok {
			request["LaunchTemplateName"] = v
		} else if v, ok := d.GetOk("name"); ok {
			request["LaunchTemplateName"] = v
		}
	}
	if d.HasChange("network_interfaces") {
		update = true
	}
	if v, ok := d.GetOk("network_interfaces"); ok {
		NetworkInterfaces := make([]map[string]interface{}, len(v.(*schema.Set).List()))
		for i, NetworkInterfacesValue := range v.(*schema.Set).List() {
			NetworkInterfacesMap := NetworkInterfacesValue.(map[string]interface{})
			NetworkInterfaces[i] = make(map[string]interface{})
			NetworkInterfaces[i]["Description"] = NetworkInterfacesMap["description"]
			NetworkInterfaces[i]["NetworkInterfaceName"] = NetworkInterfacesMap["name"]
			NetworkInterfaces[i]["PrimaryIpAddress"] = NetworkInterfacesMap["primary_ip"]
			NetworkInterfaces[i]["SecurityGroupId"] = NetworkInterfacesMap["security_group_id"]
			NetworkInterfaces[i]["VSwitchId"] = NetworkInterfacesMap["vswitch_id"]
		}
		request["NetworkInterface"] = NetworkInterfaces

	}
	if d.HasChange("network_type") {
		update = true
	}
	if v, ok := d.GetOk("network_type"); ok {
		request["NetworkType"] = v
	}
	if d.HasChange("password_inherit") {
		update = true
	}
	if v, ok := d.GetOkExists("password_inherit"); ok {
		request["PasswordInherit"] = v
	}
	if d.HasChange("period") {
		update = true
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if d.HasChange("private_ip_address") {
		update = true
	}
	if v, ok := d.GetOk("private_ip_address"); ok {
		request["PrivateIpAddress"] = v
	}
	if d.HasChange("ram_role_name") {
		update = true
	}
	if v, ok := d.GetOk("ram_role_name"); ok {
		request["RamRoleName"] = v
	}
	if d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if d.HasChange("security_enhancement_strategy") {
		update = true
	}
	if v, ok := d.GetOk("security_enhancement_strategy"); ok {
		request["SecurityEnhancementStrategy"] = v
	}
	if d.HasChange("security_group_id") {
		update = true
	}
	if v, ok := d.GetOk("security_group_id"); ok {
		request["SecurityGroupId"] = v
	}
	if d.HasChange("security_group_ids") {
		update = true
	}
	if v, ok := d.GetOk("security_group_ids"); ok {
		request["SecurityGroupIds"] = v
	}
	if d.HasChange("spot_duration") {
		update = true
	}
	if v, ok := d.GetOk("spot_duration"); ok {
		request["SpotDuration"] = v
	}
	if d.HasChange("spot_price_limit") {
		update = true
	}
	if v, ok := d.GetOk("spot_price_limit"); ok {
		request["SpotPriceLimit"] = v
	}
	if d.HasChange("spot_strategy") {
		update = true
	}
	if v, ok := d.GetOk("spot_strategy"); ok {
		request["SpotStrategy"] = v
	}
	if d.HasChange("system_disk") {
		update = true

		systemDiskMap := make(map[string]interface{})
		for _, systemDisk := range d.Get("system_disk").(*schema.Set).List() {
			systemDiskArg := systemDisk.(map[string]interface{})
			systemDiskMap["Category"] = systemDiskArg["category"]
			systemDiskMap["DeleteWithInstance"] = requests.NewBoolean(systemDiskArg["delete_with_instance"].(bool))
			systemDiskMap["Description"] = systemDiskArg["description"]
			systemDiskMap["Iops"] = systemDiskArg["iops"]
			systemDiskMap["DiskName"] = systemDiskArg["name"]
			systemDiskMap["PerformanceLevel"] = systemDiskArg["performance_level"]
			systemDiskMap["Size"] = requests.NewInteger(systemDiskArg["size"].(int))
			systemDiskMap["Encrypted"] = requests.NewBoolean(systemDiskArg["encrypted"].(bool))
		}
		request["SystemDisk"] = systemDiskMap
	} else {
		systemDiskMap := make(map[string]interface{})
		if d.HasChange("system_disk_category") {
			update = true
		}
		if v, ok := d.GetOk("system_disk_category"); ok {
			systemDiskMap["Category"] = v
		}
		if d.HasChange("system_disk_description") {
			update = true
		}
		if v, ok := d.GetOk("system_disk_description"); ok {
			systemDiskMap["Description"] = v
		}
		if d.HasChange("system_disk_name") {
			update = true
		}
		if v, ok := d.GetOk("system_disk_name"); ok {
			systemDiskMap["DiskName"] = v
		}
		if d.HasChange("system_disk_size") {
			update = true
		}
		if v, ok := d.GetOk("system_disk_size"); ok {
			systemDiskMap["Size"] = v
		}
		diskMap := d.Get("system_disk").(*schema.Set).List()[0].(map[string]interface{})
		systemDiskMap["DeleteWithInstance"] = diskMap["delete_with_instance"]
		systemDiskMap["Iops"] = diskMap["iops"]
		systemDiskMap["PerformanceLevel"] = diskMap["performance_level"]
		request["SystemDisk"] = systemDiskMap
	}
	if d.HasChange("tags") {
		update = true
	}
	if v, ok := d.GetOk("tags"); ok {
		count := 1
		for key, value := range v.(map[string]interface{}) {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
	}
	if d.HasChange("user_data") {
		update = true
		request["UserData"] = d.Get("user_data")
	} else if d.HasChange("userdata") {
		update = true
		request["UserData"] = d.Get("userdata")
	} else {
		if v, ok := d.GetOk("user_data"); ok {
			request["UserData"] = v
		} else if v, ok := d.GetOk("userdata"); ok {
			request["UserData"] = v
		}
	}
	if d.HasChange("vswitch_id") {
		update = true
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	if d.HasChange("version_description") {
		update = true
	}
	if v, ok := d.GetOk("version_description"); ok {
		request["VersionDescription"] = v
	}
	if d.HasChange("vpc_id") {
		update = true
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}
	if d.HasChange("zone_id") {
		update = true
	}
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}
	if update {
		action := "CreateLaunchTemplateVersion"
		conn, err := client.NewEcsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		latestVersion = response["LaunchTemplateVersionNumber"]
	}

	if d.Get("update_default_version_number").(bool) || d.HasChange("default_version_number") {
		action := "ModifyLaunchTemplateDefaultVersion"
		request = map[string]interface{}{
			"LaunchTemplateId": d.Id(),
			"RegionId":         client.RegionId,
		}

		if d.Get("update_default_version_number").(bool) {
			request["DefaultVersionNumber"] = latestVersion
		} else if d.HasChange("default_version_number") {
			request["DefaultVersionNumber"] = d.Get("default_version_number")
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
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
	d.Partial(false)
	return resourceAliCloudEcsLaunchTemplateRead(d, meta)
}

func resourceAliCloudEcsLaunchTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteLaunchTemplate"
	var response map[string]interface{}
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"LaunchTemplateId": d.Id(),
	}

	if v, ok := d.GetOk("launch_template_name"); ok {
		request["LaunchTemplateName"] = v
	}
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return nil
}

func getLaunchTemplateVersions(id string, meta interface{}) ([]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)
	action := "DescribeLaunchTemplateVersions"
	var response map[string]interface{}
	conn, err := client.NewEcsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	request := map[string]interface{}{
		"LaunchTemplateId": id,
	}
	request["PageSize"] = requests.NewInteger(50)
	request["RegionId"] = client.RegionId

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return nil, WrapError(err)
	}

	addDebug(action, response, request)
	v, err := jsonpath.Get("$.LaunchTemplateVersionSets.LaunchTemplateVersionSet", response)
	if err != nil {
		return nil, WrapErrorf(err, FailedGetAttributeMsg, id, "$.LaunchTemplateVersionSets.LaunchTemplateVersionSet", response)
	}

	if len(v.([]interface{})) < 1 {
		return nil, WrapErrorf(Error(GetNotFoundMessage("ECS", id)), NotFoundWithResponse, response)
	}

	return v.([]interface{}), nil
}

func deleteLaunchTemplateVersion(id string, version int, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteLaunchTemplateVersion"
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"LaunchTemplateId": id,
	}
	request["DeleteVersion"] = &[]string{strconv.FormatInt(int64(version), 10)}
	request["RegionId"] = client.RegionId

	_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
