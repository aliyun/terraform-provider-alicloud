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

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

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
			"instance_type": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringMatch(regexp.MustCompile(`^ecs\..*`), "prefix must be 'ecs.'"),
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
			"scaling_configuration_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      PayByBandwidth,
				ValidateFunc: validation.StringInSlice([]string{"PayByBandwidth", "PayByTraffic"}, false),
			},
			"internet_max_bandwidth_in": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"internet_max_bandwidth_out": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 1024),
			},
			"credit_specification": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(CreditSpecificationStandard),
					string(CreditSpecificationUnlimited),
				}, false),
			},
			"system_disk_category": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      DiskCloudEfficiency,
				ValidateFunc: validation.StringInSlice([]string{"cloud", "ephemeral_ssd", "cloud_ssd", "cloud_essd", "cloud_efficiency"}, false),
			},
			"system_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(20, 500),
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
			"data_disk": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"category": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"all", "cloud", "ephemeral_ssd", "cloud_essd", "cloud_efficiency", "cloud_ssd", "local_disk"}, false),
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
				ValidateFunc: validation.StringLenBetween(2, 128),
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
							ValidateFunc: validation.StringInSlice([]string{"EntryLevel", "EnterpriseLevel", "CreditEntryLevel"}, false),
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
	conn, err := client.NewEssClient()
	if err != nil {
		return WrapError(err)
	}
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

	if v := d.Get("image_name").(string); v != "" {
		request["ImageName"] = d.Get("image_name")
	}

	if v := d.Get("internet_charge_type").(string); v != "" {
		request["InternetChargeType"] = d.Get("internet_charge_type")
	}

	if v := d.Get("system_disk_category").(string); v != "" {
		request["SystemDisk.Category"] = d.Get("system_disk_category")
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
			instancePatternInfosMap["Memory"] = strconv.FormatFloat(item["memory"].(float64), 'f', 2, 64)
			instancePatternInfosMap["MaxPrice"] = strconv.FormatFloat(item["max_price"].(float64), 'f', 2, 64)
			instancePatternInfosMap["Cores"] = item["cores"].(int)

			instancePatternInfosMaps = append(instancePatternInfosMaps, instancePatternInfosMap)
		}
		request["InstancePatternInfo"] = instancePatternInfosMaps
	}
	request["IoOptimized"] = string(IOOptimized)

	if d.Get("is_outdated").(bool) == true {
		request["IoOptimized"] = string(NoneOptimized)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-28"), StringPointer("AK"), nil, request, &runtime)
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
	conn, err := client.NewEssClient()
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
	if d.HasChange("scaling_configuration_name") {
		request["ScalingConfigurationName"] = d.Get("scaling_configuration_name")
		update = true
	}

	if d.HasChange("internet_charge_type") {
		request["InternetChargeType"] = d.Get("internet_charge_type")
		update = true
	}

	if d.HasChange("system_disk_category") {
		request["SystemDisk.Category"] = d.Get("system_disk_category")
		update = true
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

	if (hasChangeInstanceType || hasChangeInstanceTypes || d.Get("override").(bool)) && (!hasChangeInstanceTypeOverrides || len(typeOverride.(*schema.Set).List()) == 0) {
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
	if hasChangeSecurityGroupId || hasChangeSecurityGroupIds || d.Get("override").(bool) {
		securityGroupId := d.Get("security_group_id").(string)
		securityGroupIds := d.Get("security_group_ids").([]interface{})
		if securityGroupId == "" && (securityGroupIds == nil || len(securityGroupIds) == 0) {
			return fmt.Errorf("securityGroupId or securityGroupIds must be assigned")
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

	if d.HasChange("user_data") {
		if v, ok := d.GetOk("user_data"); ok && v.(string) != "" {
			_, base64DecodeError := base64.StdEncoding.DecodeString(v.(string))
			if base64DecodeError == nil {
				request["UserData"] = v
			} else {
				request["UserData"] = base64.StdEncoding.EncodeToString([]byte(v.(string)))
			}
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
			instancePatternInfos := make([]ess.ModifyScalingConfigurationInstancePatternInfo, 0)
			for _, e := range v.(*schema.Set).List() {
				pack := e.(map[string]interface{})
				l := ess.ModifyScalingConfigurationInstancePatternInfo{
					InstanceFamilyLevel: pack["instance_family_level"].(string),
					Memory:              strconv.FormatFloat(pack["memory"].(float64), 'f', 2, 64),
					MaxPrice:            strconv.FormatFloat(pack["max_price"].(float64), 'f', 2, 64),
					Cores:               strconv.Itoa(pack["cores"].(int)),
				}
				instancePatternInfos = append(instancePatternInfos, l)
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

	if update {
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-28"), StringPointer("AK"), nil, request, &runtime)
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
	d.Set("internet_charge_type", response["InternetChargeType"])
	d.Set("system_disk_category", response["SystemDiskCategory"])
	d.Set("internet_max_bandwidth_in", response["InternetMaxBandwidthIn"])
	d.Set("internet_max_bandwidth_out", response["InternetMaxBandwidthOut"])

	d.Set("credit_specification", response["CreditSpecification"])
	d.Set("system_disk_size", response["SystemDiskSize"])
	d.Set("system_disk_name", response["SystemDiskName"])
	d.Set("system_disk_description", response["SystemDiskDescription"])
	d.Set("system_disk_auto_snapshot_policy_id", response["SystemDiskAutoSnapshotPolicyId"])
	d.Set("system_disk_performance_level", response["SystemDiskPerformanceLevel"])
	d.Set("system_disk_encrypted", response["SystemDisk.Encrypted"])

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
			d.Set("user_data", userDataHashSum(response["UserData"].(string)))
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
			f, _ := r["MaxPrice"].(json.Number).Float64()
			maxPrice, _ := strconv.ParseFloat(strconv.FormatFloat(f, 'f', 2, 64), 64)
			l := map[string]interface{}{
				"instance_family_level": r["InstanceFamilyLevel"],
				"memory":                r["Memory"],
				"cores":                 r["Cores"],
				"max_price":             maxPrice,
			}
			result = append(result, l)
		}
		err := d.Set("instance_pattern_info", result)
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
