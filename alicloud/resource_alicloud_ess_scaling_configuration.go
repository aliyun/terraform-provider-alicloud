package alicloud

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEssScalingConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssScalingConfigurationCreate,
		Read:   resourceAliyunEssScalingConfigurationRead,
		Update: resourceAliyunEssScalingConfigurationUpdate,
		Delete: resourceAliyunEssScalingConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"scaling_group_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"system_disk_encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"image_options_login_as_non_root": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"system_disk_encrypt_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"AES-256", "SM4-128"}, false),
			},
			"system_disk_kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_type": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  StringMatch(regexp.MustCompile(`^ecs\..*`), "prefix must be 'ecs.'"),
				ConflictsWith: []string{"instance_types", "instance_type_override"},
			},
			"instance_types": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:      true,
				ConflictsWith: []string{"instance_type", "instance_type_override"},
				MaxItems:      int(MaxScalingConfigurationInstanceTypes),
			},
			"io_optimized": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Attribute io_optimized has been deprecated on instance resource. All the launched alicloud instances will be IO optimized. Suggest to remove it from your template.",
			},
			"is_outdated": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"security_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"security_group_ids"},
			},
			"security_group_ids": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ConflictsWith: []string{"security_group_id"},
				Optional:      true,
				MaxItems:      16,
			},
			"network_interfaces": {
				Optional:      true,
				Type:          schema.TypeSet,
				ConflictsWith: []string{"security_group_id"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"Primary", "Secondary"}, false),
						},
						"network_interface_traffic_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"Standard", "HighPerformance"}, false),
						},
						"ipv6_address_count": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"security_group_ids": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
					},
				},
			},
			"custom_priorities": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"scaling_configuration_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringLenBetween(2, 256),
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      PayByBandwidth,
				ValidateFunc: StringInSlice([]string{"PayByBandwidth", "PayByTraffic"}, false),
			},
			"internet_max_bandwidth_in": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"spot_duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 1),
			},
			"internet_max_bandwidth_out": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 1024),
			},
			"credit_specification": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: StringInSlice([]string{
					string(CreditSpecificationStandard),
					string(CreditSpecificationUnlimited),
				}, false),
			},
			"system_disk_category": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      DiskCloudEfficiency,
				ValidateFunc: StringInSlice([]string{"cloud", "ephemeral_ssd", "cloud_ssd", "cloud_essd", "cloud_efficiency", "cloud_essd_xc1"}, false),
			},
			"security_enhancement_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Active", "Deactive"}, false),
			},
			"system_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(20, 500),
			},
			"system_disk_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"system_disk_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"system_disk_auto_snapshot_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"system_disk_provisioned_iops": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"data_disk": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"provisioned_iops": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"category": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"all", "cloud", "ephemeral_ssd", "cloud_essd", "cloud_efficiency", "cloud_ssd", "local_disk"}, false),
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"device": {
							Type:       schema.TypeString,
							Optional:   true,
							Deprecated: "Attribute device has been deprecated on disk attachment resource. Suggest to remove it from your template.",
						},
						"delete_with_instance": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"encrypted": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"auto_snapshot_policy_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"performance_level": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"instance_ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
				Deprecated: "Field 'instance_ids' has been deprecated from provider version 1.6.0. New resource 'alicloud_ess_attachment' replaces it.",
			},

			"substitute": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"role_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"key_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"force_delete": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ESS-Instance",
				ValidateFunc: StringLenBetween(2, 128),
			},

			"override": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("password_inherit").(bool)
				},
			},
			"password_inherit": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"kms_encrypted_password": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("password_inherit").(bool) || d.Get("password").(string) != ""
				},
			},
			"kms_encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_password").(string) == ""
				},
				Elem: schema.TypeString,
			},
			"system_disk_performance_level": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"spot_strategy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"spot_price_limit": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"price_limit": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
					},
				},
			},
			"instance_type_override": {
				Optional:      true,
				Type:          schema.TypeSet,
				ConflictsWith: []string{"instance_type", "instance_types"},
				MaxItems:      int(MaxScalingConfigurationInstanceTypes),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"weighted_capacity": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"instance_pattern_info": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_family_level": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"EntryLevel", "EnterpriseLevel", "CreditEntryLevel"}, false),
						},
						"burstable_performance": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"Exclude", "Include", "Required"}, false),
						},
						"excluded_instance_types": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
						"architectures": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: StringInSlice([]string{"X86", "Heterogeneous", "BareMental", "Arm", "SuperComputeCluster"}, false),
							},
							Optional: true,
						},
						"cores": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"memory": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"max_price": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliyunEssScalingConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateScalingConfiguration"
	request := make(map[string]interface{})
	var err error
	securityGroupId := d.Get("security_group_id").(string)
	securityGroupIds := d.Get("security_group_ids").([]interface{})
	if securityGroupId == "" && (securityGroupIds == nil || len(securityGroupIds) == 0) {
		return WrapError(Error("security_group_id or security_group_ids must be assigned"))
	}
	instanceType := d.Get("instance_type").(string)
	instanceTypes := d.Get("instance_types").([]interface{})
	instanceTypeOverrides := d.Get("instance_type_override").(*schema.Set).List()
	if instanceType == "" && (instanceTypes == nil || len(instanceTypes) == 0) && (instanceTypeOverrides == nil || len(instanceTypeOverrides) == 0) {
		return WrapError(Error("instance_type or instance_types or instance_type_override must be assigned"))
	}
	request["ImageId"] = d.Get("image_id")
	request["ScalingGroupId"] = d.Get("scaling_group_id")
	request["PasswordInherit"] = d.Get("password_inherit")
	request["SystemDisk.Encrypted"] = d.Get("system_disk_encrypted")
	request["ImageOptions.LoginAsNonRoot"] = d.Get("image_options_login_as_non_root")
	request["DeletionProtection"] = d.Get("deletion_protection")
	request["SystemDisk.KMSKeyId"] = d.Get("system_disk_kms_key_id")
	if securityGroupId != "" {
		request["SecurityGroupId"] = d.Get("security_group_id")
	}
	password := d.Get("password").(string)
	kmsPassword := d.Get("kms_encrypted_password").(string)
	if password != "" {
		request["Password"] = d.Get("password")
	} else if kmsPassword != "" {
		kmsService := KmsService{client}
		decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request["Password"] = decryptResp
	}
	if securityGroupIds != nil && len(securityGroupIds) > 0 {
		request["SecurityGroupIds"] = securityGroupIds
	}

	types := make([]string, 0, int(MaxScalingConfigurationInstanceTypes))
	if instanceTypes != nil && len(instanceTypes) > 0 {
		types = expandStringList(instanceTypes)
	}

	if instanceType != "" {
		types = append(types, instanceType)
	}
	request["InstanceTypes"] = types

	if instanceTypeOverrides != nil && len(instanceTypeOverrides) != 0 {
		instanceTypeOverridesMaps := make([]map[string]interface{}, 0)
		for _, rew := range instanceTypeOverrides {
			instanceTypeOverridesMap := make(map[string]interface{})
			item := rew.(map[string]interface{})
			if instanceType, ok := item["instance_type"].(string); ok && instanceType != "" {
				instanceTypeOverridesMap["InstanceType"] = instanceType
			}
			instanceTypeOverridesMap["WeightedCapacity"] = item["weighted_capacity"].(int)

			instanceTypeOverridesMaps = append(instanceTypeOverridesMaps, instanceTypeOverridesMap)
		}
		request["InstanceTypeOverride"] = instanceTypeOverridesMaps
	}

	if v := d.Get("scaling_configuration_name").(string); v != "" {
		request["ScalingConfigurationName"] = d.Get("scaling_configuration_name")
	}
	if v, ok := d.GetOk("instance_description"); ok {
		request["InstanceDescription"] = v
	}

	if v, ok := d.GetOk("system_disk_encrypt_algorithm"); ok {
		request["SystemDisk.EncryptAlgorithm"] = v
	}

	if v := d.Get("image_name").(string); v != "" {
		request["ImageName"] = d.Get("image_name")
	}

	if v := d.Get("internet_charge_type").(string); v != "" {
		request["InternetChargeType"] = d.Get("internet_charge_type")
	}

	if v, ok := d.GetOkExists("spot_duration"); ok {
		request["SpotDuration"] = v
	}

	if v, ok := d.GetOk("system_disk_provisioned_iops"); ok {
		request["SystemDisk.ProvisionedIops"] = v
	}

	if v := d.Get("system_disk_category").(string); v != "" {
		request["SystemDisk.Category"] = d.Get("system_disk_category")
	}
	if v, ok := d.GetOk("security_enhancement_strategy"); ok && v.(string) != "" {
		request["SecurityEnhancementStrategy"] = v
	}

	if v := d.Get("internet_max_bandwidth_in").(int); v != 0 {
		request["InternetMaxBandwidthIn"] = d.Get("internet_max_bandwidth_in")
	}

	request["InternetMaxBandwidthOut"] = d.Get("internet_max_bandwidth_out")

	if v := d.Get("credit_specification").(string); v != "" {
		request["CreditSpecification"] = d.Get("credit_specification")
	}

	if v := d.Get("system_disk_size").(int); v != 0 {
		request["SystemDisk.Size"] = d.Get("system_disk_size")
	}

	if v := d.Get("system_disk_name").(string); v != "" {
		request["SystemDisk.DiskName"] = d.Get("system_disk_name")
	}

	if v := d.Get("system_disk_description").(string); v != "" {
		request["SystemDisk.Description"] = d.Get("system_disk_description")
	}

	if v := d.Get("system_disk_auto_snapshot_policy_id").(string); v != "" {
		request["SystemDisk.AutoSnapshotPolicyId"] = d.Get("system_disk_auto_snapshot_policy_id")
	}

	if v := d.Get("system_disk_performance_level").(string); v != "" {
		request["SystemDisk.PerformanceLevel"] = d.Get("system_disk_performance_level")
	}

	if v := d.Get("resource_group_id").(string); v != "" {
		request["ResourceGroupId"] = d.Get("resource_group_id")
	}

	if v, ok := d.GetOk("role_name"); ok && v.(string) != "" {
		request["RamRoleName"] = v
	}

	if v, ok := d.GetOk("key_name"); ok && v.(string) != "" {
		request["KeyPairName"] = v
	}

	if v, ok := d.GetOk("instance_name"); ok && v.(string) != "" {
		request["InstanceName"] = v
	}

	if v, ok := d.GetOk("host_name"); ok && v.(string) != "" {
		request["HostName"] = v
	}

	if v, ok := d.GetOk("spot_strategy"); ok && v.(string) != "" {
		request["SpotStrategy"] = v
	}

	if v, ok := d.GetOk("user_data"); ok {
		_, base64DecodeError := base64.StdEncoding.DecodeString(v.(string))
		if base64DecodeError == nil {
			request["UserData"] = v
		} else {
			request["UserData"] = base64.StdEncoding.EncodeToString([]byte(v.(string)))
		}
	}

	if v, ok := d.GetOk("tags"); ok {
		vv, okk := convertMaptoJsonString(v.(map[string]interface{}))
		if okk == nil {
			request["Tags"] = vv
		}
	}

	if v, ok := d.GetOk("data_disk"); ok {
		disksMaps := make([]map[string]interface{}, 0)
		disks := v.([]interface{})

		for _, rew := range disks {
			disksMap := make(map[string]interface{})
			item := rew.(map[string]interface{})
			disksMap["Size"] = item["size"].(int)
			disksMap["ProvisionedIops"] = item["provisioned_iops"].(int)

			if category, ok := item["category"].(string); ok && category != "" {
				disksMap["Category"] = category
			}
			if snapshotId, ok := item["snapshot_id"].(string); ok && snapshotId != "" {
				disksMap["SnapshotId"] = snapshotId
			}
			disksMap["DeleteWithInstance"] = item["delete_with_instance"].(bool)
			if device, ok := item["device"].(string); ok && device != "" {
				disksMap["Device"] = device
			}
			disksMap["Encrypted"] = item["encrypted"].(bool)
			if kmsKeyId, ok := item["kms_key_id"].(string); ok && kmsKeyId != "" {
				disksMap["KMSKeyId"] = kmsKeyId
			}
			if name, ok := item["name"].(string); ok && name != "" {
				disksMap["DiskName"] = name
			}
			if description, ok := item["description"].(string); ok && description != "" {
				disksMap["Description"] = description
			}
			if autoSnapshotPolicyId, ok := item["auto_snapshot_policy_id"].(string); ok && autoSnapshotPolicyId != "" {
				disksMap["AutoSnapshotPolicyId"] = autoSnapshotPolicyId
			}
			if performanceLevel, ok := item["performance_level"].(string); ok && performanceLevel != "" {
				disksMap["PerformanceLevel"] = performanceLevel
			}
			disksMaps = append(disksMaps, disksMap)
		}
		request["DataDisk"] = disksMaps
	}
	if v, ok := d.GetOk("spot_price_limit"); ok {
		spotPriceLimitsMaps := make([]map[string]interface{}, 0)
		spotPriceLimits := v.(*schema.Set).List()
		for _, rew := range spotPriceLimits {
			spotPriceLimitsMap := make(map[string]interface{})
			item := rew.(map[string]interface{})
			if instanceType, ok := item["instance_type"].(string); ok && instanceType != "" {
				spotPriceLimitsMap["InstanceType"] = instanceType
			}
			spotPriceLimitsMap["PriceLimit"] = strconv.FormatFloat(item["price_limit"].(float64), 'f', 2, 64)
			spotPriceLimitsMaps = append(spotPriceLimitsMaps, spotPriceLimitsMap)
		}
		request["SpotPriceLimit"] = spotPriceLimitsMaps
	}
	if v, ok := d.GetOk("instance_pattern_info"); ok {
		instancePatternInfosMaps := make([]map[string]interface{}, 0)
		instancePatternInfos := v.(*schema.Set).List()
		for _, rew := range instancePatternInfos {
			instancePatternInfosMap := make(map[string]interface{})
			item := rew.(map[string]interface{})

			if instanceFamilyLevel, ok := item["instance_family_level"].(string); ok && instanceFamilyLevel != "" {
				instancePatternInfosMap["InstanceFamilyLevel"] = instanceFamilyLevel
			}
			if burstablePerformance, ok := item["burstable_performance"].(string); ok && burstablePerformance != "" {
				instancePatternInfosMap["BurstablePerformance"] = burstablePerformance
			}
			instancePatternInfosMap["Memory"] = strconv.FormatFloat(item["memory"].(float64), 'f', 2, 64)
			if item["max_price"].(float64) != 0 {
				instancePatternInfosMap["MaxPrice"] = strconv.FormatFloat(item["max_price"].(float64), 'f', 2, 64)
			}
			instancePatternInfosMap["Cores"] = item["cores"].(int)
			instancePatternInfosMap["Architectures"] = item["architectures"]
			instancePatternInfosMap["ExcludedInstanceTypes"] = item["excluded_instance_types"]

			instancePatternInfosMaps = append(instancePatternInfosMaps, instancePatternInfosMap)
		}
		request["InstancePatternInfo"] = instancePatternInfosMaps
	}
	if v, ok := d.GetOk("network_interfaces"); ok {
		networkInterfacesMaps := make([]map[string]interface{}, 0)
		networkInterfaces := v.(*schema.Set).List()
		for _, rew := range networkInterfaces {
			networkInterfacesMap := make(map[string]interface{})
			item := rew.(map[string]interface{})

			if instanceType, ok := item["instance_type"].(string); ok && instanceType != "" {
				networkInterfacesMap["InstanceType"] = instanceType
			}
			if networkInterfaceTrafficMode, ok := item["network_interface_traffic_mode"].(string); ok && networkInterfaceTrafficMode != "" {
				networkInterfacesMap["NetworkInterfaceTrafficMode"] = networkInterfaceTrafficMode
			}
			networkInterfacesMap["Ipv6AddressCount"] = item["ipv6_address_count"].(int)
			networkInterfacesMap["SecurityGroupIds"] = item["security_group_ids"]
			networkInterfacesMaps = append(networkInterfacesMaps, networkInterfacesMap)
		}
		request["NetworkInterfaces"] = networkInterfacesMaps
	}

	if v, ok := d.GetOk("custom_priorities"); ok {
		customPrioritiesMaps := make([]map[string]interface{}, 0)
		customPriorities := v.(*schema.Set).List()
		for _, rew := range customPriorities {
			customPrioritiesMap := make(map[string]interface{})
			item := rew.(map[string]interface{})

			if instanceType, ok := item["instance_type"].(string); ok && instanceType != "" {
				customPrioritiesMap["InstanceType"] = instanceType
			}

			if vswitchId, ok := item["vswitch_id"].(string); ok && vswitchId != "" {
				customPrioritiesMap["VswitchId"] = vswitchId
			}
			customPrioritiesMaps = append(customPrioritiesMaps, customPrioritiesMap)
		}
		request["CustomPriorities"] = customPrioritiesMaps
	}
	request["IoOptimized"] = string(IOOptimized)

	if d.Get("is_outdated").(bool) == true {
		request["IoOptimized"] = string(NoneOptimized)
	}
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Ess", "2014-08-28", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{Throttling, "IncorrectScalingGroupStatus"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ess_scaling_configuration", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ScalingConfigurationId"]))
	return resourceAliyunEssScalingConfigurationUpdate(d, meta)
}

func resourceAliyunEssScalingConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	d.Partial(true)
	if strings.Contains(d.Id(), COLON_SEPARATED) {
		d.SetId(strings.Split(d.Id(), COLON_SEPARATED)[1])
	}

	if d.HasChange("active") {
		c, err := essService.DescribeEssScalingConfiguration(d.Id())
		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapError(err)
		}

		if d.Get("active").(bool) {
			if c.LifecycleState == string(Inactive) {

				err := essService.ActiveEssScalingConfiguration(c.ScalingGroupId, d.Id())
				if err != nil {
					return WrapError(err)
				}
			}
		} else {
			if c.LifecycleState == string(Active) {
				_, err := activeSubstituteScalingConfiguration(d, meta)
				if err != nil {
					return WrapError(err)
				}
			}
		}
		d.SetPartial("active")
	}

	if err := enableEssScalingConfiguration(d, meta); err != nil {
		return WrapError(err)
	}
	if err := modifyEssScalingConfiguration(d, meta); err != nil {
		return WrapError(err)
	}

	d.Partial(false)

	return resourceAliyunEssScalingConfigurationRead(d, meta)
}

func modifyEssScalingConfiguration(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	action := "ModifyScalingConfiguration"
	request := map[string]interface{}{
		"ScalingConfigurationId": d.Id(),
	}
	update := false

	if d.HasChange("override") {
		request["Override"] = d.Get("override")
		update = true
	}
	if d.HasChange("password_inherit") {
		request["PasswordInherit"] = d.Get("password_inherit")
		update = true
	}
	if d.HasChange("image_id") || d.Get("override").(bool) {
		request["ImageId"] = d.Get("image_id")
		update = true
	}
	if d.HasChange("image_name") || d.Get("override").(bool) {
		request["ImageName"] = d.Get("image_name")
		update = true
	}
	if d.HasChange("instance_description") {
		request["InstanceDescription"] = d.Get("instance_description")
		update = true
	}

	if d.HasChange("system_disk_encrypt_algorithm") {
		request["SystemDisk.EncryptAlgorithm"] = d.Get("system_disk_encrypt_algorithm")
		update = true
	}

	if d.HasChange("scaling_configuration_name") {
		request["ScalingConfigurationName"] = d.Get("scaling_configuration_name")
		update = true
	}

	if d.HasChange("internet_charge_type") {
		request["InternetChargeType"] = d.Get("internet_charge_type")
		update = true
	}

	if d.HasChange("spot_duration") {
		if v, ok := d.GetOkExists("spot_duration"); ok {
			request["SpotDuration"] = requests.NewInteger(v.(int))
			update = true
		}
	}

	if d.HasChange("system_disk_provisioned_iops") {
		if v, ok := d.GetOkExists("system_disk_provisioned_iops"); ok {
			request["SystemDisk.ProvisionedIops"] = requests.NewInteger(v.(int))
			update = true
		}
	}

	if d.HasChange("system_disk_category") {
		request["SystemDisk.Category"] = d.Get("system_disk_category")
		update = true
	}

	if d.HasChange("internet_max_bandwidth_in") {
		if v := d.Get("internet_max_bandwidth_in").(int); v != 0 {
			request["InternetMaxBandwidthIn"] = v
			update = true
		}
	}

	if d.HasChange("internet_max_bandwidth_out") {
		request["InternetMaxBandwidthOut"] = d.Get("internet_max_bandwidth_out")
		update = true
	}

	if d.HasChange("credit_specification") {
		request["CreditSpecification"] = d.Get("credit_specification")
		update = true
	}

	if d.HasChange("system_disk_size") {
		request["SystemDisk.Size"] = d.Get("system_disk_size")
		update = true
	}

	if d.HasChange("system_disk_name") {
		request["SystemDisk.DiskName"] = d.Get("system_disk_name")
		update = true
	}

	if d.HasChange("system_disk_description") {
		request["SystemDisk.Description"] = d.Get("system_disk_description")
		update = true
	}

	if d.HasChange("system_disk_auto_snapshot_policy_id") {
		request["SystemDisk.AutoSnapshotPolicyId"] = d.Get("system_disk_auto_snapshot_policy_id")
		update = true
	}

	if d.HasChange("system_disk_performance_level") {
		request["SystemDisk.PerformanceLevel"] = d.Get("system_disk_performance_level")
		update = true
	}
	if d.HasChange("resource_group_id") {
		request["ResourceGroupId"] = d.Get("resource_group_id")
		update = true
	}
	if d.HasChange("role_name") {
		request["RamRoleName"] = d.Get("role_name")
		update = true
	}

	if d.HasChange("key_name") {
		request["KeyPairName"] = d.Get("key_name")
		update = true
	}

	if d.HasChange("instance_name") {
		request["InstanceName"] = d.Get("instance_name")
		update = true
	}

	if d.HasChange("host_name") {
		request["HostName"] = d.Get("host_name")
		update = true
	}

	if d.HasChange("spot_strategy") {
		request["SpotStrategy"] = d.Get("spot_strategy")
		update = true
	}

	hasChangeInstanceType := d.HasChange("instance_type")
	hasChangeInstanceTypes := d.HasChange("instance_types")
	hasChangeInstanceTypeOverrides := d.HasChange("instance_type_override")
	typeOverride := d.Get("instance_type_override")

	if (hasChangeInstanceType || hasChangeInstanceTypes) && (!hasChangeInstanceTypeOverrides || len(typeOverride.(*schema.Set).List()) == 0) {
		instanceType := d.Get("instance_type").(string)
		instanceTypes := d.Get("instance_types").([]interface{})
		if instanceType == "" && (instanceTypes == nil || len(instanceTypes) == 0) {
			return fmt.Errorf("instance_type or instance_types must be assigned")
		}
		types := make([]string, 0, int(MaxScalingConfigurationInstanceTypes))
		if instanceTypes != nil && len(instanceTypes) > 0 {
			types = expandStringList(instanceTypes)
		}
		if instanceType != "" {
			types = append(types, instanceType)
		}
		request["InstanceTypes"] = types
		update = true
	}

	if hasChangeInstanceTypeOverrides && d.Get("override").(bool) {
		v, ok := d.GetOk("instance_type_override")
		if ok {
			instanceTypeOverrides := make([]ess.ModifyScalingConfigurationInstanceTypeOverride, 0)
			for _, e := range v.(*schema.Set).List() {
				pack := e.(map[string]interface{})
				l := ess.ModifyScalingConfigurationInstanceTypeOverride{
					InstanceType:     pack["instance_type"].(string),
					WeightedCapacity: strconv.Itoa(pack["weighted_capacity"].(int)),
				}
				instanceTypeOverrides = append(instanceTypeOverrides, l)
			}
			request["InstanceTypeOverride"] = instanceTypeOverrides
		}
		update = true
	}

	hasChangeSecurityGroupId := d.HasChange("security_group_id")
	hasChangeSecurityGroupIds := d.HasChange("security_group_ids")
	hasChangeNetworkInterfaces := d.HasChange("network_interfaces")
	if hasChangeSecurityGroupId || hasChangeSecurityGroupIds || hasChangeNetworkInterfaces || d.Get("override").(bool) {
		securityGroupId := d.Get("security_group_id").(string)
		securityGroupIds := d.Get("security_group_ids").([]interface{})
		v, ok := d.GetOk("network_interfaces")
		networkInterfaces := make([]map[string]interface{}, 0)
		if ok {
			for _, e := range v.(*schema.Set).List() {
				pack := e.(map[string]interface{})
				securityGroupIdsFormat := pack["security_group_ids"]
				securityGroupIds := toStringArray(securityGroupIdsFormat)
				networkInterface := map[string]interface{}{
					"InstanceType":                pack["instance_type"].(string),
					"NetworkInterfaceTrafficMode": pack["network_interface_traffic_mode"].(string),
					"Ipv6AddressCount":            strconv.Itoa(pack["ipv6_address_count"].(int)),
					"SecurityGroupIds":            &securityGroupIds,
				}
				networkInterfaces = append(networkInterfaces, networkInterface)
			}
			request["NetworkInterfaces"] = networkInterfaces
		}
		if securityGroupId == "" && len(networkInterfaces) == 0 && (securityGroupIds == nil || len(securityGroupIds) == 0) {
			return fmt.Errorf("securityGroupId or securityGroupIds must be assigned")
		}
		if len(networkInterfaces) > 0 {
			request["NetworkInterfaces"] = networkInterfaces
		}
		if securityGroupIds != nil && len(securityGroupIds) > 0 {
			sgs := expandStringList(securityGroupIds)
			request["SecurityGroupIds"] = sgs
		}
		if securityGroupId != "" {
			request["SecurityGroupId"] = securityGroupId
		}
		update = true
	}

	if d.HasChange("custom_priorities") {
		v, ok := d.GetOk("custom_priorities")
		if ok {
			customPriorities := make([]map[string]interface{}, 0)
			for _, e := range v.(*schema.Set).List() {
				pack := e.(map[string]interface{})
				customPrioritiy := map[string]interface{}{
					"InstanceType": pack["instance_type"].(string),
					"VswitchId":    pack["vswitch_id"].(string),
				}
				customPriorities = append(customPriorities, customPrioritiy)
			}
			request["CustomPriorities"] = customPriorities
			update = true
		}
	}

	if d.HasChange("user_data") {
		if v, ok := d.GetOk("user_data"); ok && v.(string) != "" {
			_, base64DecodeError := base64.StdEncoding.DecodeString(v.(string))
			if base64DecodeError == nil {
				request["UserData"] = v
			} else {
				request["UserData"] = base64.StdEncoding.EncodeToString([]byte(v.(string)))
			}
		}
		if v, ok := d.GetOk("user_data"); ok && v.(string) == "" {
			request["UserData"] = ""
		}
		update = true
	}
	if d.HasChange("spot_price_limit") {
		v, ok := d.GetOk("spot_price_limit")
		if ok {
			spotPriceLimits := make([]ess.ModifyScalingConfigurationSpotPriceLimit, 0)
			for _, e := range v.(*schema.Set).List() {
				pack := e.(map[string]interface{})
				l := ess.ModifyScalingConfigurationSpotPriceLimit{
					InstanceType: pack["instance_type"].(string),
					PriceLimit:   strconv.FormatFloat(pack["price_limit"].(float64), 'f', 2, 64),
				}
				spotPriceLimits = append(spotPriceLimits, l)
			}
			request["SpotPriceLimit"] = spotPriceLimits
		}
		update = true
	}

	if d.HasChange("instance_pattern_info") {
		v, ok := d.GetOk("instance_pattern_info")
		if ok {
			instancePatternInfos := make([]map[string]interface{}, 0)
			for _, e := range v.(*schema.Set).List() {
				pack := e.(map[string]interface{})
				excludedInstancezTypesFormat := pack["excluded_instance_types"]
				excludedInstancezTypes := toStringArray(excludedInstancezTypesFormat)
				architecturesFormat := pack["architectures"]
				architectures := toStringArray(architecturesFormat)
				instancePatternInfo := map[string]interface{}{
					"InstanceFamilyLevel":  pack["instance_family_level"].(string),
					"Memory":               strconv.FormatFloat(pack["memory"].(float64), 'f', 2, 64),
					"Cores":                strconv.Itoa(pack["cores"].(int)),
					"BurstablePerformance": pack["burstable_performance"].(string),
					"ExcludedInstanceType": &excludedInstancezTypes,
					"Architecture":         &architectures,
				}
				if pack["max_price"].(float64) != 0 {
					instancePatternInfo["MaxPrice"] = strconv.FormatFloat(pack["max_price"].(float64), 'f', 2, 64)
				}
				instancePatternInfos = append(instancePatternInfos, instancePatternInfo)
			}
			request["InstancePatternInfo"] = instancePatternInfos
		}
		update = true
	}
	if d.HasChange("tags") {
		if v, ok := d.GetOk("tags"); ok {
			tags := "{"
			for key, value := range v.(map[string]interface{}) {
				tags += "\"" + key + "\"" + ":" + "\"" + value.(string) + "\"" + ","
			}
			request["Tags"] = strings.TrimSuffix(tags, ",") + "}"
		}
		update = true
	}

	if d.HasChange("data_disk") {
		dds, ok := d.GetOk("data_disk")
		if ok {
			disks := dds.([]interface{})
			dataDisks := make([]map[string]interface{}, 0, len(disks))
			for _, e := range disks {
				pack := e.(map[string]interface{})
				dataDisks = append(dataDisks, map[string]interface{}{
					"Size":                 pack["size"],
					"ProvisionedIops":      pack["provisioned_iops"],
					"Category":             pack["category"],
					"SnapshotId":           pack["snapshot_id"],
					"DeleteWithInstance":   pack["delete_with_instance"],
					"Device":               pack["device"],
					"Encrypted":            pack["encrypted"],
					"KMSKeyId":             pack["kms_key_id"],
					"DiskName":             pack["name"],
					"Description":          pack["description"],
					"AutoSnapshotPolicyId": pack["auto_snapshot_policy_id"],
					"PerformanceLevel":     pack["performance_level"],
				})
			}
			request["DataDisk"] = dataDisks
		}
		update = true
	}

	if d.HasChange("system_disk_encrypted") {
		request["SystemDisk.Encrypted"] = d.Get("system_disk_encrypted")
		update = true
	}

	if d.HasChange("image_options_login_as_non_root") {
		request["ImageOptions.LoginAsNonRoot"] = d.Get("image_options_login_as_non_root")
		update = true
	}

	if d.HasChange("deletion_protection") {
		request["DeletionProtection"] = d.Get("deletion_protection")
		update = true
	}

	if d.HasChange("system_disk_kms_key_id") {
		request["SystemDisk.KMSKeyId"] = d.Get("system_disk_kms_key_id")
		update = true
	}

	if update {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := client.RpcPost("Ess", "2014-08-28", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_ess_scaling_configuration", "ModifyScalingConfiguration", AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}

func enableEssScalingConfiguration(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}

	if d.HasChange("enable") {
		sgId := d.Get("scaling_group_id").(string)
		group, err := essService.DescribeEssScalingGroup(sgId)
		if err != nil {
			return WrapError(err)
		}

		if d.Get("enable").(bool) {
			if group.LifecycleState == string(Inactive) {

				object, err := essService.DescribeEssScalingConfifurations(sgId)

				if err != nil {
					return WrapError(err)
				}
				activeConfig := ""
				var csIds []string
				for _, c := range object {
					csIds = append(csIds, c.ScalingConfigurationId)
					if c.LifecycleState == string(Active) {
						activeConfig = c.ScalingConfigurationId
					}
				}

				if activeConfig == "" {
					return WrapError(Error("Please active a scaling configuration before enabling scaling group %s. Its all scaling configuration are %s.",
						sgId, strings.Join(csIds, ",")))
				}

				request := ess.CreateEnableScalingGroupRequest()
				request.RegionId = client.RegionId
				request.ScalingGroupId = sgId
				request.ActiveScalingConfigurationId = activeConfig

				raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
					return essClient.EnableScalingGroup(request)
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				if err := essService.WaitForEssScalingGroup(sgId, Active, DefaultTimeout); err != nil {
					return WrapError(err)
				}

				d.SetPartial("scaling_configuration_id")
			}
		} else {
			if group.LifecycleState == string(Active) {
				request := ess.CreateDisableScalingGroupRequest()
				request.RegionId = client.RegionId
				request.ScalingGroupId = sgId
				raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
					return essClient.DisableScalingGroup(request)
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				if err := essService.WaitForEssScalingGroup(sgId, Inactive, DefaultTimeout); err != nil {
					return WrapError(err)
				}
			}
		}
		d.SetPartial("enable")
	}

	return nil
}

func resourceAliyunEssScalingConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	if strings.Contains(d.Id(), COLON_SEPARATED) {
		d.SetId(strings.Split(d.Id(), COLON_SEPARATED)[1])
	}
	response, err := essService.DescribeEssScalingConfigurationByCommonApi(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("scaling_group_id", response["ScalingGroupId"])
	d.Set("image_id", response["ImageId"])
	d.Set("image_name", response["ImageName"])
	d.Set("scaling_configuration_name", response["ScalingConfigurationName"])
	d.Set("instance_description", response["InstanceDescription"])
	if response["SystemDisk.EncryptAlgorithm"] != nil {
		d.Set("system_disk_encrypt_algorithm", response["SystemDisk.EncryptAlgorithm"])
	}
	d.Set("internet_charge_type", response["InternetChargeType"])
	if response["SpotDuration"] != nil {
		d.Set("spot_duration", response["SpotDuration"])
	}
	if response["SystemDisk.ProvisionedIops"] != nil {
		d.Set("system_disk_provisioned_iops", response["SystemDisk.ProvisionedIops"])
	}
	d.Set("system_disk_category", response["SystemDiskCategory"])
	d.Set("security_enhancement_strategy", response["SecurityEnhancementStrategy"])
	if response["InternetMaxBandwidthIn"] != nil {
		d.Set("internet_max_bandwidth_in", response["InternetMaxBandwidthIn"])
	}
	d.Set("internet_max_bandwidth_out", response["InternetMaxBandwidthOut"])

	d.Set("credit_specification", response["CreditSpecification"])
	if response["SystemDiskSize"] != nil && response["SystemDiskSize"] != 0 {
		d.Set("system_disk_size", response["SystemDiskSize"])
	}
	d.Set("system_disk_name", response["SystemDiskName"])
	d.Set("system_disk_description", response["SystemDiskDescription"])
	d.Set("system_disk_auto_snapshot_policy_id", response["SystemDiskAutoSnapshotPolicyId"])
	d.Set("system_disk_performance_level", response["SystemDiskPerformanceLevel"])
	d.Set("system_disk_kms_key_id", response["SystemDisk.KMSKeyId"])
	d.Set("system_disk_encrypted", response["SystemDisk.Encrypted"])
	d.Set("image_options_login_as_non_root", response["ImageOptions.LoginAsNonRoot"])
	d.Set("deletion_protection", response["DeletionProtection"])
	d.Set("role_name", response["RamRoleName"])
	d.Set("key_name", response["KeyPairName"])
	d.Set("force_delete", d.Get("force_delete").(bool))

	d.Set("instance_name", response["InstanceName"])
	d.Set("override", d.Get("override").(bool))

	d.Set("password_inherit", response["PasswordInherit"])
	d.Set("resource_group_id", response["ResourceGroupId"])

	d.Set("host_name", response["HostName"])
	d.Set("spot_strategy", response["SpotStrategy"])

	if sg, ok := d.GetOk("security_group_id"); ok && sg.(string) != "" {
		d.Set("security_group_id", response["SecurityGroupId"])
	}
	if instanceType, ok := d.GetOk("instance_type"); ok && instanceType.(string) != "" {
		d.Set("instance_type", response["InstanceType"])
	}
	userData := d.Get("user_data")
	if userData.(string) != "" {
		_, base64DecodeError := base64.StdEncoding.DecodeString(userData.(string))
		if base64DecodeError == nil {
			d.Set("user_data", response["UserData"])
		} else {
			if response["UserData"] != nil {
				d.Set("user_data", userDataHashSum(response["UserData"].(string)))
			}
		}
	} else {
		if response["UserData"] != nil {
			d.Set("user_data", userDataHashSum(response["UserData"].(string)))
		}
	}
	d.Set("active", response["LifecycleState"].(string) == string(Active))

	if instanceTypes, ok := d.GetOk("instance_types"); ok && len(instanceTypes.([]interface{})) > 0 {
		v := response["InstanceTypes"].(map[string]interface{})
		d.Set("instance_types", v["InstanceType"].([]interface{}))
	}

	if instanceTypeOverride, ok := d.GetOk("instance_type_override"); ok && len(instanceTypeOverride.(*schema.Set).List()) > 0 {
		if v := response["InstanceTypeOverrides"]; v != nil {
			result := make([]map[string]interface{}, 0)
			for _, i := range v.(map[string][]ess.ModifyScalingConfigurationInstanceTypeOverride)["InstanceTypeOverride"] {
				r := i
				n, _ := strconv.Atoi(r.WeightedCapacity)
				l := map[string]interface{}{
					"instance_type":     r.InstanceType,
					"weighted_capacity": n,
				}
				result = append(result, l)
			}
			err := d.Set("instance_type_override", result)
			if err != nil {
				return WrapError(err)
			}
		}
	}
	if sgs, ok := d.GetOk("security_group_ids"); ok && len(sgs.([]interface{})) > 0 {
		v := response["SecurityGroupIds"].(map[string]interface{})
		d.Set("security_group_ids", v["SecurityGroupId"].([]interface{}))
	}

	if v := response["DataDisks"]; v != nil {
		result := make([]map[string]interface{}, 0)
		for _, i := range v.(map[string]interface{})["DataDisk"].([]interface{}) {
			r := i.(map[string]interface{})
			l := map[string]interface{}{
				"size":                    r["Size"],
				"provisioned_iops":        r["ProvisionedIops"],
				"category":                r["Category"],
				"snapshot_id":             r["SnapshotId"],
				"device":                  r["Device"],
				"delete_with_instance":    r["DeleteWithInstance"],
				"encrypted":               r["Encrypted"],
				"kms_key_id":              r["KMSKeyId"],
				"disk_name":               r["DiskName"],
				"description":             r["Description"],
				"auto_snapshot_policy_id": r["AutoSnapshotPolicyId"],
				"performance_level":       r["PerformanceLevel"],
			}
			result = append(result, l)
		}
		d.Set("data_disk", result)
	}
	if v := response["Tags"]; v != nil {
		d.Set("tags", tagsToMap(response["Tags"].(map[string]interface{})["Tag"]))
	}
	if v := response["SpotPriceLimit"]; v != nil {
		result := make([]map[string]interface{}, 0)
		for _, i := range v.(map[string]interface{})["SpotPriceModel"].([]interface{}) {
			r := i.(map[string]interface{})
			f, _ := r["PriceLimit"].(json.Number).Float64()
			p, _ := strconv.ParseFloat(strconv.FormatFloat(f, 'f', 2, 64), 64)
			l := map[string]interface{}{
				"instance_type": r["InstanceType"],
				"price_limit":   p,
			}
			result = append(result, l)
		}
		d.Set("spot_price_limit", result)
	}

	if v := response["InstancePatternInfos"]; v != nil {
		result := make([]map[string]interface{}, 0)
		for _, i := range v.(map[string]interface{})["InstancePatternInfo"].([]interface{}) {
			r := i.(map[string]interface{})
			var arr []string
			var arr1 []string
			if r["ExcludedInstanceTypes"] != nil {
				i := r["ExcludedInstanceTypes"].(map[string]interface{})
				arr = toStringArray(i["ExcludedInstanceType"])
			}
			if r["Architectures"] != nil {
				j := r["Architectures"].(map[string]interface{})
				arr1 = toStringArray(j["Architecture"])
			}
			l := map[string]interface{}{
				"instance_family_level":   r["InstanceFamilyLevel"],
				"memory":                  r["Memory"],
				"cores":                   r["Cores"],
				"excluded_instance_types": &arr,
				"architectures":           &arr1,
				"burstable_performance":   r["BurstablePerformance"],
			}
			if r["MaxPrice"] != nil {
				f, _ := r["MaxPrice"].(json.Number).Float64()
				maxPrice, _ := strconv.ParseFloat(strconv.FormatFloat(f, 'f', 2, 64), 64)
				l["max_price"] = maxPrice
			}
			result = append(result, l)
		}
		err := d.Set("instance_pattern_info", result)
		if err != nil {
			return WrapError(err)
		}
	}
	if v := response["NetworkInterfaces"]; v != nil {
		result := make([]map[string]interface{}, 0)
		for _, i := range v.(map[string]interface{})["NetworkInterface"].([]interface{}) {
			r := i.(map[string]interface{})
			var arr1 []string
			if r["SecurityGroupIds"] != nil {
				j := r["SecurityGroupIds"].(map[string]interface{})
				arr1 = toStringArray(j["SecurityGroupId"])
			}
			l := map[string]interface{}{
				"instance_type":                  r["InstanceType"],
				"network_interface_traffic_mode": r["NetworkInterfaceTrafficMode"],
				"ipv6_address_count":             r["Ipv6AddressCount"],
				"security_group_ids":             &arr1,
			}
			result = append(result, l)
		}
		err := d.Set("network_interfaces", result)
		if err != nil {
			return WrapError(err)
		}
	}

	if v := response["CustomPriorities"]; v != nil {
		result := make([]map[string]interface{}, 0)
		for _, i := range v.(map[string]interface{})["CustomPriority"].([]interface{}) {
			r := i.(map[string]interface{})
			l := map[string]interface{}{
				"instance_type": r["InstanceType"],
				"vswitch_id":    r["VswitchId"],
			}
			result = append(result, l)
		}
		err := d.Set("custom_priorities", result)
		if err != nil {
			return WrapError(err)
		}
	}
	return nil
}

func resourceAliyunEssScalingConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}

	if strings.Contains(d.Id(), COLON_SEPARATED) {
		d.SetId(strings.Split(d.Id(), COLON_SEPARATED)[1])
	}

	object, err := essService.DescribeEssScalingConfiguration(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}

	request := ess.CreateDescribeScalingConfigurationsRequest()
	request.RegionId = client.RegionId
	request.ScalingGroupId = object.ScalingGroupId

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScalingConfigurations(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ess.DescribeScalingConfigurationsResponse)
	if len(response.ScalingConfigurations.ScalingConfiguration) < 1 {
		return nil
	} else if len(response.ScalingConfigurations.ScalingConfiguration) == 1 {
		if d.Get("force_delete").(bool) {
			request := ess.CreateDeleteScalingGroupRequest()
			request.ScalingGroupId = object.ScalingGroupId
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
		}
		return WrapError(Error("Current scaling configuration %s is the last configuration for the scaling group %s. Please launch a new "+
			"active scaling configuration or set 'force_delete' to 'true' to delete it with deleting its scaling group.", d.Id(), object.ScalingGroupId))
	}

	deleteScalingConfigurationRequest := ess.CreateDeleteScalingConfigurationRequest()
	deleteScalingConfigurationRequest.ScalingConfigurationId = d.Id()

	rawDeleteScalingConfiguration, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DeleteScalingConfiguration(deleteScalingConfigurationRequest)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidScalingGroupId.NotFound", "InvalidScalingConfigurationId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), rawDeleteScalingConfiguration, request.RpcRequest, request)

	return WrapError(essService.WaitForScalingConfiguration(d.Id(), Deleted, DefaultTimeout))
}

func activeSubstituteScalingConfiguration(d *schema.ResourceData, meta interface{}) (configures []ess.ScalingConfigurationInDescribeScalingConfigurations, err error) {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	substituteId, ok := d.GetOk("substitute")

	c, err := essService.DescribeEssScalingConfiguration(d.Id())
	if err != nil {
		err = WrapError(err)
		return
	}

	request := ess.CreateDescribeScalingConfigurationsRequest()
	request.RegionId = client.RegionId
	request.ScalingGroupId = c.ScalingGroupId

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScalingConfigurations(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ess.DescribeScalingConfigurationsResponse)
	if len(response.ScalingConfigurations.ScalingConfiguration) < 1 {
		return
	}

	if !ok || substituteId.(string) == "" {

		if len(response.ScalingConfigurations.ScalingConfiguration) == 1 {
			return configures, WrapError(Error("Current scaling configuration %s is the last configuration for the scaling group %s, and it can't be inactive.", d.Id(), c.ScalingGroupId))
		}

		var configs []string
		for _, cc := range response.ScalingConfigurations.ScalingConfiguration {
			if cc.ScalingConfigurationId != d.Id() {
				configs = append(configs, cc.ScalingConfigurationId)
			}
		}

		return configures, WrapError(Error("Before inactivating current scaling configuration, you must select a substitute for scaling group from: %s.", strings.Join(configs, ",")))

	}

	err = essService.ActiveEssScalingConfiguration(c.ScalingGroupId, substituteId.(string))
	if err != nil {
		return configures, WrapError(Error("Inactive scaling configuration %s err: %#v. Substitute scaling configuration ID: %s",
			d.Id(), err, substituteId.(string)))
	}

	return response.ScalingConfigurations.ScalingConfiguration, nil
}

func toStringArray(a interface{}) []string {
	var paramSlice []string
	strArr, _ := a.([]interface{})
	for _, param := range strArr {
		paramSlice = append(paramSlice, param.(string))
	}
	return paramSlice
}
