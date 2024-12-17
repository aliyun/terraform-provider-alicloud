package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/denverdino/aliyungo/common"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"encoding/base64"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudInstanceCreate,
		Read:   resourceAliCloudInstanceRead,
		Update: resourceAliCloudInstanceUpdate,
		Delete: resourceAliCloudInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"image_id": {
				Type:         schema.TypeString,
				AtLeastOneOf: []string{"image_id", "launch_template_id", "launch_template_name"},
				Optional:     true,
				Computed:     true,
			},
			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile(`^ecs\..*`), "prefix must be 'ecs.'"),
				AtLeastOneOf: []string{"instance_type", "launch_template_id", "launch_template_name"},
				Computed:     true,
			},
			"credit_specification": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: StringInSlice([]string{
					string(CreditSpecificationStandard),
					string(CreditSpecificationUnlimited),
				}, false),
			},
			"security_groups": {
				Type:         schema.TypeSet,
				Elem:         &schema.Schema{Type: schema.TypeString},
				Computed:     true,
				Optional:     true,
				AtLeastOneOf: []string{"security_groups", "launch_template_id", "launch_template_name"},
			},
			"allocate_public_ip": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Field 'allocate_public_ip' has been deprecated from provider version 1.6.1. Setting 'internet_max_bandwidth_out' larger than 0 will allocate public ip for instance.",
			},
			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringLenBetween(2, 128),
				Computed:     true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringLenBetween(2, 256),
				Computed:     true,
			},
			"internet_charge_type": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     StringInSlice([]string{"PayByBandwidth", "PayByTraffic"}, false),
				DiffSuppressFunc: ecsInternetDiffSuppressFunc,
				Computed:         true,
			},
			"internet_max_bandwidth_in": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: ecsInternetDiffSuppressFunc,
				Deprecated:       "The attribute is invalid and no any affect for the instance. So it has been deprecated since version v1.121.2.",
			},
			"internet_max_bandwidth_out": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"kms_encrypted_password": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: kmsDiffSuppressFunc,
			},
			"kms_encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_password") == ""
				},
				Elem: schema.TypeString,
			},
			"password_inherit": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"io_optimized": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Attribute io_optimized has been deprecated on instance resource. All the launched alicloud instances will be IO optimized. Suggest to remove it from your template.",
				Removed:    "Attribute 'io_optimized' has been removed from provider version 1.213.1.",
			},
			"is_outdated": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"system_disk_category": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"all", "cloud", "ephemeral_ssd", "cloud_essd", "cloud_efficiency", "cloud_ssd", "local_disk", "cloud_auto", "cloud_essd_entry"}, false),
			},
			"system_disk_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringLenBetween(2, 128),
			},
			"system_disk_description": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringLenBetween(2, 256),
			},
			"system_disk_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"system_disk_performance_level": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: ecsSystemDiskPerformanceLevelSuppressFunc,
				ValidateFunc:     StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
			},
			"system_disk_auto_snapshot_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"system_disk_storage_cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"system_disk_encrypted": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"system_disk_kms_key_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"system_disk_encrypt_algorithm": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"system_disk_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_interface_traffic_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Standard", "HighPerformance"}, false),
			},
			"network_card_index": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"queue_pair_number": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"data_disks": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 16,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: StringLenBetween(2, 128),
						},
						"size": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"category": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"all", "cloud", "ephemeral_ssd", "cloud_essd", "cloud_efficiency", "cloud_ssd", "local_disk", "cloud_auto", "cloud_essd_entry"}, false),
							Default:      DiskCloudEfficiency,
							ForceNew:     true,
						},
						"encrypted": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: true,
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"auto_snapshot_policy_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"delete_with_instance": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  true,
						},
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: StringLenBetween(2, 256),
						},
						"performance_level": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
						},
						"device": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"network_interfaces": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_interface_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"network_interface_traffic_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"Standard", "HighPerformance"}, false),
						},
						"network_card_index": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"queue_pair_number": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"security_group_ids": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			//subnet_id and vswitch_id both exists, cause compatible old version, and aws habit.
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true, //add this schema cause subnet_id not used enter parameter, will different, so will be ForceNew
				Removed:  "Field 'subnet_id' has been removed from version 1.210.0",
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.Any(
					IntBetween(1, 9),
					IntInSlice([]int{12, 24, 36, 48, 60})),
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"period_unit": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          Month,
				ValidateFunc:     StringInSlice([]string{"Week", "Month"}, false),
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"renewal_status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  RenewNormal,
				ValidateFunc: StringInSlice([]string{
					string(RenewAutoRenewal),
					string(RenewNormal),
					string(RenewNotRenewal)}, false),
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"auto_renew_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ValidateFunc:     IntInSlice([]int{1, 2, 3, 6, 12}),
				DiffSuppressFunc: ecsNotAutoRenewDiffSuppressFunc,
			},
			"include_data_disks": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          true,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Running", "Stopped"}, false),
				Computed:     true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if (d.Get("user_data") != "" && new == "") && (d.Get("launch_template_name").(string) != "" ||
						d.Get("launch_template_id").(string) != "") {
						return true
					}
					return false
				},
			},
			"role_name": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				DiffSuppressFunc: vpcTypeResourceDiffSuppressFunc,
			},
			"key_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"spot_strategy": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateFunc:     StringInSlice([]string{"NoSpot", "SpotAsPriceGo", "SpotWithPriceLimit"}, false),
				DiffSuppressFunc: ecsSpotStrategyDiffSuppressFunc,
			},
			"spot_price_limit": {
				Type:             schema.TypeFloat,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				DiffSuppressFunc: ecsSpotPriceLimitDiffSuppressFunc,
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"force_delete": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          false,
				Description:      descriptions["A behavior mark used to delete 'PrePaid' ECS instance forcibly."],
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"security_enhancement_strategy": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				ValidateFunc: StringInSlice([]string{
					string(ActiveSecurityEnhancementStrategy),
					string(DeactiveSecurityEnhancementStrategy),
				}, false),
			},
			"tags":        tagsSchemaWithIgnore(),
			"volume_tags": tagsSchemaComputed(),
			"auto_release_time": {
				Type:     schema.TypeString,
				Optional: true,
				//ValidateFunc: ValidateRFC3339TimeString(true),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if (d.Get("auto_release_time") != "" && new == "") && (d.Get("launch_template_name").(string) != "" ||
						d.Get("launch_template_id").(string) != "") {
						return true
					}
					if d.Get("instance_charge_type").(string) == "PrePaid" {
						return true
					}
					return old != "" && new != "" && strings.HasPrefix(new, strings.Trim(old, ":00Z"))
				},
			},
			"hpc_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"secondary_private_ips": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:      true,
				ConflictsWith: []string{"secondary_private_ip_address_count"},
			},
			"secondary_private_ip_address_count": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"secondary_private_ips"},
			},
			"deployment_set_id": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if (d.Get("deployment_set_id") != "" && new == "") && (d.Get("launch_template_name").(string) != "" ||
						d.Get("launch_template_id").(string) != "") {
						return true
					}
					return false
				},
			},
			"deployment_set_group_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operator_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"upgrade", "downgrade"}, false),
			},
			"stopped_mode": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"StopCharging", "KeepCharging", "Not-applicable"}, false),
			},
			"maintenance_time": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"maintenance_action": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Stop", "AutoRecover", "AutoRedeploy"}, false),
			},
			"maintenance_notify": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"spot_duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(0, 6),
			},
			"http_tokens": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"optional", "required"}, false),
			},
			"http_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"enabled", "disabled"}, false),
			},
			"http_put_response_hop_limit": {
				Type:         schema.TypeInt,
				ForceNew:     true,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(1, 64),
			},
			"ipv6_address_count": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  IntBetween(1, 10),
				ConflictsWith: []string{"ipv6_addresses"},
			},
			"ipv6_addresses": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				MaxItems:      10,
				ForceNew:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"ipv6_address_count"},
			},
			"network_interface_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"os_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dedicated_host_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"spot_strategy", "spot_price_limit"},
			},
			"launch_template_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"launch_template_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"launch_template_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enable_jumbo_frame": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expired_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_options": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"login_as_non_root": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	var response map[string]interface{}
	action := "RunInstances"
	request := make(map[string]interface{})
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("instance_type"); ok {
		request["InstanceType"] = v
	}

	if v, ok := d.GetOk("image_id"); ok {
		request["ImageId"] = v
	}

	if v, ok := d.GetOk("availability_zone"); ok {
		request["ZoneId"] = v
	}

	if v, ok := d.GetOk("system_disk_name"); ok {
		request["SystemDisk.DiskName"] = v
	}

	if v, ok := d.GetOk("system_disk_description"); ok {
		request["SystemDisk.Description"] = v
	}

	if v, ok := d.GetOk("system_disk_performance_level"); ok {
		request["SystemDisk.PerformanceLevel"] = v
	}

	if v, ok := d.GetOk("system_disk_category"); ok {
		request["SystemDisk.Category"] = v
	}

	request["SystemDisk.Size"] = 40
	if v, ok := d.GetOk("system_disk_size"); ok {
		request["SystemDisk.Size"] = v
	}

	if v, ok := d.GetOk("system_disk_auto_snapshot_policy_id"); ok {
		request["SystemDisk.AutoSnapshotPolicyId"] = v
	}

	if v, ok := d.GetOk("system_disk_storage_cluster_id"); ok {
		request["SystemDisk.StorageClusterId"] = v
	}

	if v, ok := d.GetOkExists("system_disk_encrypted"); ok {
		request["SystemDisk.Encrypted"] = v
	}

	if v, ok := d.GetOk("system_disk_kms_key_id"); ok {
		request["SystemDisk.KMSKeyId"] = v
	}

	if v, ok := d.GetOk("system_disk_encrypt_algorithm"); ok {
		request["SystemDisk.EncryptAlgorithm"] = v
	}

	request["InstanceName"] = "ECS-Instance"
	if v, ok := d.GetOk("instance_name"); ok {
		request["InstanceName"] = v
	}

	if v, ok := d.GetOk("credit_specification"); ok {
		request["CreditSpecification"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("launch_template_name"); ok {
		request["LaunchTemplateName"] = v
	}
	if v, ok := d.GetOk("launch_template_id"); ok {
		request["LaunchTemplateId"] = v
	}
	if v, ok := d.GetOk("launch_template_version"); ok {
		request["LaunchTemplateVersion"] = v
	}

	request["InternetChargeType"] = "PayByTraffic"
	if v, ok := d.GetOk("internet_charge_type"); ok {
		request["InternetChargeType"] = v
	}

	request["InternetMaxBandwidthOut"] = 0
	if v, ok := d.GetOk("internet_max_bandwidth_out"); ok {
		request["InternetMaxBandwidthOut"] = v
	}

	if v, ok := d.GetOk("internet_max_bandwidth_in"); ok {
		request["InternetMaxBandwidthIn"] = v
	}

	if v, ok := d.GetOk("host_name"); ok {
		request["HostName"] = v
	}

	if v, ok := d.GetOk("password"); ok {
		request["Password"] = v
	}

	if v, ok := d.GetOk("kms_encrypted_password"); ok {
		kmsService := KmsService{client}
		decryptResp, err := kmsService.Decrypt(v.(string), d.Get("kms_encryption_context").(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request["Password"] = decryptResp
	}

	if v, ok := d.GetOkExists("password_inherit"); ok {
		request["PasswordInherit"] = v
	}

	vswitchValue := d.Get("vswitch_id")
	if vswitchValue == "" {
		vswitchValue = d.Get("subnet_id")
	}

	request["InstanceChargeType"] = "PostPaid"
	if v, ok := d.GetOk("instance_charge_type"); ok {
		request["InstanceChargeType"] = v
	}

	if request["InstanceChargeType"] == string(PrePaid) {
		if v, ok := d.GetOk("period"); ok {
			request["Period"] = v
		}
		if v, ok := d.GetOk("period_unit"); ok {
			request["PeriodUnit"] = v
		}
		if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
			request["AutoRenew"] = true
		}
		if v, ok := d.GetOk("auto_renew_period"); ok {
			request["AutoRenewPeriod"] = v
		}
	} else {
		if v, ok := d.GetOk("spot_strategy"); ok {
			request["SpotStrategy"] = v
		}
		if v, ok := d.GetOk("spot_price_limit"); ok {
			request["SpotPriceLimit"] = v
		}
	}

	if v, ok := d.GetOk("user_data"); ok {
		_, base64DecodeError := base64.StdEncoding.DecodeString(v.(string))
		if base64DecodeError == nil {
			request["UserData"] = v
		} else {
			request["UserData"] = base64.StdEncoding.EncodeToString([]byte(v.(string)))
		}
	}

	if v, ok := d.GetOk("role_name"); ok {
		request["RamRoleName"] = v
	}

	if v, ok := d.GetOk("key_name"); ok {
		request["KeyPairName"] = v
	}

	if v, ok := d.GetOk("security_enhancement_strategy"); ok {
		request["SecurityEnhancementStrategy"] = v
	}

	if v, ok := d.GetOk("auto_release_time"); ok && v.(string) != "" {
		request["AutoReleaseTime"] = v
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	if v, ok := d.GetOkExists("deletion_protection"); ok {
		request["DeletionProtection"] = v
	}

	if v, ok := d.GetOk("tags"); ok {
		count := 1
		for key, value := range v.(map[string]interface{}) {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
	}

	if v, ok := d.GetOk("data_disks"); ok {
		disksMaps := make([]map[string]interface{}, 0)
		disks := v.([]interface{})
		for _, rew := range disks {

			disksMap := make(map[string]interface{})
			item := rew.(map[string]interface{})

			disksMap["DeleteWithInstance"] = item["delete_with_instance"].(bool)
			disksMap["Encrypted"] = item["encrypted"].(bool)
			disksMap["Size"] = item["size"].(int)

			if category, ok := item["category"].(string); ok && category != "" {
				disksMap["Category"] = category
			}

			if name, ok := item["name"].(string); ok && name != "" {
				disksMap["DiskName"] = name
			}

			if kmsKeyId, ok := item["kms_key_id"].(string); ok && kmsKeyId != "" {
				disksMap["KMSKeyId"] = kmsKeyId
			}

			if snapshotId, ok := item["snapshot_id"].(string); ok && snapshotId != "" {
				disksMap["SnapshotId"] = snapshotId
			}

			if description, ok := item["description"].(string); ok && description != "" {
				disksMap["Description"] = description
			}

			if autoSnapshotPolicyId, ok := item["auto_snapshot_policy_id"].(string); ok && autoSnapshotPolicyId != "" {
				disksMap["AutoSnapshotPolicyId"] = autoSnapshotPolicyId
			}

			if device, ok := item["device"].(string); ok && device != "" {
				disksMap["Device"] = device
			}

			if performanceLevel, ok := item["performance_level"].(string); ok && performanceLevel != "" && disksMap["Category"] == string(DiskCloudESSD) {
				disksMap["PerformanceLevel"] = performanceLevel
			}

			if disksMap["Category"] == string(DiskEphemeralSSD) {
				disksMap["DeleteWithInstance"] = ""
			}

			disksMaps = append(disksMaps, disksMap)
		}
		request["DataDisk"] = disksMaps
	}

	networkInterfacesMaps := make([]map[string]interface{}, 0)

	_, networkInterfaceTrafficModeOk := d.GetOk("network_interface_traffic_mode")
	_, networkCardIndexOk := d.GetOkExists("network_card_index")
	_, queuePairNumberOk := d.GetOkExists("queue_pair_number")

	if networkInterfaceTrafficModeOk || networkCardIndexOk || queuePairNumberOk {
		primaryNetworkInterfacesMap := make(map[string]interface{})
		primaryNetworkInterfacesMap["InstanceType"] = "Primary"

		if v, ok := d.GetOk("security_groups"); ok {
			// At present, the classic network instance does not support multi sg in runInstances
			sgs := expandStringList(v.(*schema.Set).List())
			if d.Get("vswitch_id").(string) == "" && len(sgs) > 0 {
				primaryNetworkInterfacesMap["SecurityGroupId"] = sgs[0]
			} else {
				primaryNetworkInterfacesMap["SecurityGroupIds"] = sgs
			}
		}

		if vswitchValue != "" {
			primaryNetworkInterfacesMap["VSwitchId"] = vswitchValue

			if v, ok := d.GetOk("private_ip"); ok {
				primaryNetworkInterfacesMap["PrimaryIpAddress"] = v
			}
		}

		if v, ok := d.GetOk("ipv6_addresses"); ok {
			primaryNetworkInterfacesMap["Ipv6Address"] = v.(*schema.Set).List()
		}

		if v, ok := d.GetOkExists("ipv6_address_count"); ok {
			primaryNetworkInterfacesMap["Ipv6AddressCount"] = v
		}

		if v, ok := d.GetOk("network_interface_traffic_mode"); ok {
			primaryNetworkInterfacesMap["NetworkInterfaceTrafficMode"] = v
		}

		if v, ok := d.GetOkExists("network_card_index"); ok {
			primaryNetworkInterfacesMap["NetworkCardIndex"] = v
		}

		if v, ok := d.GetOkExists("queue_pair_number"); ok {
			primaryNetworkInterfacesMap["QueuePairNumber"] = v
		}

		networkInterfacesMaps = append(networkInterfacesMaps, primaryNetworkInterfacesMap)
	} else {
		if vswitchValue != "" {
			request["VSwitchId"] = vswitchValue

			if v, ok := d.GetOk("private_ip"); ok {
				request["PrivateIpAddress"] = v
			}
		}

		if v, ok := d.GetOk("security_groups"); ok {
			// At present, the classic network instance does not support multi sg in runInstances
			sgs := expandStringList(v.(*schema.Set).List())
			if d.Get("vswitch_id").(string) == "" && len(sgs) > 0 {
				request["SecurityGroupId"] = sgs[0]
			} else {
				request["SecurityGroupIds"] = sgs
			}
		}

		if v, ok := d.GetOk("ipv6_addresses"); ok {
			request["Ipv6Address"] = v.(*schema.Set).List()
		}

		if v, ok := d.GetOkExists("ipv6_address_count"); ok {
			request["Ipv6AddressCount"] = v
		}
	}

	if v, ok := d.GetOk("network_interfaces"); ok {
		for _, networkInterfaces := range v.([]interface{}) {
			secondaryNetworkInterfacesMap := make(map[string]interface{})
			secondaryNetworkInterfacesArg := networkInterfaces.(map[string]interface{})

			if networkInterfaceId, ok := secondaryNetworkInterfacesArg["network_interface_id"]; ok && fmt.Sprint(networkInterfaceId) != "" {
				secondaryNetworkInterfacesMap["NetworkInterfaceId"] = networkInterfaceId
			} else {
				secondaryNetworkInterfacesMap["InstanceType"] = "Secondary"

				if vSwitchId, ok := secondaryNetworkInterfacesArg["vswitch_id"]; ok {
					secondaryNetworkInterfacesMap["VSwitchId"] = vSwitchId
				}

				if networkInterfaceTrafficMode, ok := secondaryNetworkInterfacesArg["network_interface_traffic_mode"]; ok {
					secondaryNetworkInterfacesMap["NetworkInterfaceTrafficMode"] = networkInterfaceTrafficMode
				}

				isSupported, err := ecsService.isSupportedNetworkCardIndex(fmt.Sprint(request["InstanceType"]))
				if err != nil {
					return WrapError(err)
				}

				if networkCardIndex, ok := secondaryNetworkInterfacesArg["network_card_index"]; ok && isSupported {
					secondaryNetworkInterfacesMap["NetworkCardIndex"] = networkCardIndex
				}

				if queuePairNumber, ok := secondaryNetworkInterfacesArg["queue_pair_number"]; ok && fmt.Sprint(queuePairNumber) != "0" {
					secondaryNetworkInterfacesMap["QueuePairNumber"] = queuePairNumber
				}

				if securityGroupIds, ok := secondaryNetworkInterfacesArg["security_group_ids"]; ok {
					secondaryNetworkInterfacesMap["SecurityGroupIds"] = securityGroupIds
				}
			}

			networkInterfacesMaps = append(networkInterfacesMaps, secondaryNetworkInterfacesMap)
		}
	}

	if len(networkInterfacesMaps) > 0 {
		request["NetworkInterface"] = networkInterfacesMaps
	}

	networkOptionsMap := make(map[string]interface{})

	if v, ok := d.GetOkExists("enable_jumbo_frame"); ok {
		networkOptionsMap["EnableJumboFrame"] = v
	}

	if len(networkOptionsMap) > 0 {
		request["NetworkOptions"] = networkOptionsMap
	}

	if v, ok := d.GetOk("hpc_cluster_id"); ok {
		request["HpcClusterId"] = v
	}

	if v, ok := d.GetOk("deployment_set_id"); ok {
		request["DeploymentSetId"] = v
	}

	if v, ok := d.GetOk("http_tokens"); ok {
		request["HttpTokens"] = v
	}

	if v, ok := d.GetOk("http_endpoint"); ok {
		request["HttpEndpoint"] = v
	}

	if v, ok := d.GetOk("http_put_response_hop_limit"); ok {
		request["HttpPutResponseHopLimit"] = v
	}

	request["IoOptimized"] = "optimized"
	if d.Get("is_outdated").(bool) {
		request["IoOptimized"] = "none"
	}

	if v, ok := d.GetOkExists("spot_duration"); ok {
		request["SpotDuration"] = v
	}

	if v, ok := d.GetOk("dedicated_host_id"); ok {
		request["DedicatedHostId"] = v
	}

	if v, ok := d.GetOk("image_options"); ok {
		for _, raw := range v.(*schema.Set).List() {
			imageOptionsArg := raw.(map[string]interface{})
			if v, ok := imageOptionsArg["login_as_non_root"]; ok {
				request["ImageOptions.LoginAsNonRoot"] = v
			}
		}
	}

	wait := incrementalWait(1*time.Second, 1*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"IncorrectVSwitchStatus"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_instance", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(response["InstanceIdSets"].(map[string]interface{})["InstanceIdSet"].([]interface{})[0]))

	stateConf := BuildStateConf([]string{"Pending", "Starting", "Stopped"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, ecsService.InstanceStateRefreshFunc(d.Id(), []string{"Stopping"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudInstanceUpdate(d, meta)
}

func resourceAliCloudInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	instance, err := ecsService.DescribeInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_instance ecsService.DescribeInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}

		return WrapError(err)
	}

	var disk ecs.Disk
	disk, err = ecsService.DescribeInstanceSystemDisk(d.Id(), instance.ResourceGroupId, d.Get("system_disk_id").(string))
	if err != nil {
		// if old resource happenes an not found error, there may system has been detached
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[WARNING] describing instance %s system disk failed. Error: %v", d.Id(), err)
		} else {
			return WrapError(err)
		}
	} else {
		d.Set("system_disk_category", disk.Category)
		d.Set("system_disk_name", disk.DiskName)
		d.Set("system_disk_description", disk.Description)
		d.Set("system_disk_size", disk.Size)
		d.Set("system_disk_auto_snapshot_policy_id", disk.AutoSnapshotPolicyId)
		d.Set("system_disk_storage_cluster_id", disk.StorageClusterId)
		d.Set("system_disk_encrypted", disk.Encrypted)
		d.Set("system_disk_kms_key_id", disk.KMSKeyId)
		d.Set("system_disk_id", disk.DiskId)
		d.Set("volume_tags", ecsService.tagsToMap(disk.Tags.Tag))
		d.Set("system_disk_performance_level", disk.PerformanceLevel)
	}
	d.Set("instance_name", instance.InstanceName)
	d.Set("resource_group_id", instance.ResourceGroupId)
	d.Set("description", instance.Description)
	d.Set("status", instance.Status)
	d.Set("availability_zone", instance.ZoneId)
	d.Set("host_name", instance.HostName)
	d.Set("image_id", instance.ImageId)
	d.Set("instance_type", instance.InstanceType)
	d.Set("password", d.Get("password").(string))
	d.Set("internet_max_bandwidth_out", instance.InternetMaxBandwidthOut)
	d.Set("internet_max_bandwidth_in", instance.InternetMaxBandwidthIn)
	d.Set("instance_charge_type", instance.InstanceChargeType)
	d.Set("key_name", instance.KeyPairName)
	d.Set("spot_strategy", instance.SpotStrategy)
	d.Set("spot_price_limit", instance.SpotPriceLimit)
	d.Set("internet_charge_type", instance.InternetChargeType)
	d.Set("deletion_protection", instance.DeletionProtection)
	d.Set("credit_specification", instance.CreditSpecification)
	// There is an api bug that it returns non-RFC3339 timestamp, like "2023-07-26T23:50Z"
	if _, err := time.Parse(time.RFC3339, instance.AutoReleaseTime); err != nil && strings.HasSuffix(instance.AutoReleaseTime, "Z") {
		d.Set("auto_release_time", fmt.Sprintf("%s:00Z", strings.TrimSuffix(instance.AutoReleaseTime, "Z")))
	} else {
		d.Set("auto_release_time", instance.AutoReleaseTime)
	}
	d.Set("tags", ecsService.tagsToMap(instance.Tags.Tag))
	d.Set("hpc_cluster_id", instance.HpcClusterId)
	d.Set("deployment_set_id", instance.DeploymentSetId)
	d.Set("deployment_set_group_no", instance.DeploymentSetGroupNo)
	d.Set("stopped_mode", instance.StoppedMode)
	if len(instance.PublicIpAddress.IpAddress) > 0 {
		d.Set("public_ip", instance.PublicIpAddress.IpAddress[0])
	} else {
		d.Set("public_ip", "")
	}

	d.Set("subnet_id", instance.VpcAttributes.VSwitchId)
	d.Set("vpc_id", instance.VpcAttributes.VpcId)
	d.Set("vswitch_id", instance.VpcAttributes.VSwitchId)
	d.Set("spot_duration", instance.SpotDuration)
	d.Set("http_tokens", instance.MetadataOptions.HttpTokens)
	d.Set("http_endpoint", instance.MetadataOptions.HttpEndpoint)
	d.Set("http_put_response_hop_limit", instance.MetadataOptions.HttpPutResponseHopLimit)
	d.Set("cpu", instance.Cpu)
	d.Set("memory", instance.Memory)
	d.Set("os_name", instance.OSName)
	d.Set("os_type", instance.OSType)
	d.Set("dedicated_host_id", instance.DedicatedHostAttribute.DedicatedHostId)
	d.Set("create_time", instance.CreationTime)
	d.Set("start_time", instance.StartTime)
	d.Set("expired_time", instance.ExpiredTime)

	imageOptionsMaps := make([]map[string]interface{}, 0)
	imageOptionsMap := make(map[string]interface{})
	imageOptionsMap["login_as_non_root"] = instance.ImageOptions.LoginAsNonRoot
	imageOptionsMaps = append(imageOptionsMaps, imageOptionsMap)
	d.Set("image_options", imageOptionsMaps)

	if len(instance.VpcAttributes.PrivateIpAddress.IpAddress) > 0 {
		d.Set("private_ip", instance.VpcAttributes.PrivateIpAddress.IpAddress[0])
	} else {
		d.Set("private_ip", strings.Join(instance.InnerIpAddress.IpAddress, ","))
	}

	sgs := make([]string, 0, len(instance.SecurityGroupIds.SecurityGroupId))
	for _, sg := range instance.SecurityGroupIds.SecurityGroupId {
		sgs = append(sgs, sg)
	}
	if err := d.Set("security_groups", sgs); err != nil {
		return WrapError(err)
	}

	dataRequest := ecs.CreateDescribeUserDataRequest()
	dataRequest.RegionId = client.RegionId
	dataRequest.InstanceId = d.Id()
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeUserData(dataRequest)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), dataRequest.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(dataRequest.GetActionName(), raw, dataRequest.RpcRequest, dataRequest)
	userDataResponse, _ := raw.(*ecs.DescribeUserDataResponse)
	if userDataResponse.UserData != "" {
		if v, ok := d.GetOk("user_data"); ok {
			_, base64DecodeError := base64.StdEncoding.DecodeString(v.(string))
			if base64DecodeError == nil {
				d.Set("user_data", userDataResponse.UserData)
			} else {
				d.Set("user_data", userDataHashSum(userDataResponse.UserData))
			}
		} else {
			d.Set("user_data", userDataResponse.UserData)
		}
	} else {
		d.Set("user_data", "")
	}

	if len(instance.VpcAttributes.VSwitchId) > 0 && (!d.IsNewResource() || d.HasChange("role_name")) {
		request := ecs.CreateDescribeInstanceRamRoleRequest()
		request.RegionId = client.RegionId
		request.InstanceIds = convertListToJsonString([]interface{}{d.Id()})
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstanceRamRole(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ecs.DescribeInstanceRamRoleResponse)
		if len(response.InstanceRamRoleSets.InstanceRamRoleSet) >= 1 {
			d.Set("role_name", response.InstanceRamRoleSets.InstanceRamRoleSet[0].RamRoleName)
		}
	}

	if instance.InstanceChargeType == string(PrePaid) {
		request := ecs.CreateDescribeInstanceAutoRenewAttributeRequest()
		request.RegionId = client.RegionId
		request.InstanceId = d.Id()

		var raw interface{}
		var err error
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.DescribeInstanceAutoRenewAttribute(request)
			})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		response, _ := raw.(*ecs.DescribeInstanceAutoRenewAttributeResponse)
		periodUnit := d.Get("period_unit").(string)
		if periodUnit == "" {
			periodUnit = "Month"
		}
		if len(response.InstanceRenewAttributes.InstanceRenewAttribute) > 0 {
			renew := response.InstanceRenewAttributes.InstanceRenewAttribute[0]
			d.Set("renewal_status", renew.RenewalStatus)
			d.Set("auto_renew_period", renew.Duration)
			if renew.RenewalStatus == "AutoRenewal" {
				periodUnit = renew.PeriodUnit
			}
			if periodUnit == "Year" {
				periodUnit = "Month"
				d.Set("auto_renew_period", renew.Duration*12)
			}
		}
		//period, err := computePeriodByUnit(instance.CreationTime, instance.ExpiredTime, d.Get("period").(int), periodUnit)
		//if err != nil {
		//	return WrapError(err)
		//}
		//thisPeriod := d.Get("period").(int)
		//if thisPeriod != 0 && thisPeriod != period {
		//	d.Set("period", thisPeriod)
		//} else {
		//	d.Set("period", period)
		//}
		d.Set("period_unit", periodUnit)
	}
	networkInterfaceId := ""
	networkInterfaceMaps := make([]map[string]interface{}, 0)
	for _, obj := range instance.NetworkInterfaces.NetworkInterface {
		if obj.Type == "Primary" {
			networkInterfaceId = obj.NetworkInterfaceId
			object, err := ecsService.DescribeEcsNetworkInterface(obj.NetworkInterfaceId)
			if err != nil {
				return WrapError(err)
			}

			d.Set("primary_ip_address", obj.PrimaryIpAddress)
			d.Set("network_interface_traffic_mode", object["NetworkInterfaceTrafficMode"])
			d.Set("queue_pair_number", object["QueuePairNumber"])

			if attachment, ok := object["Attachment"]; ok {
				attachmentArg := attachment.(map[string]interface{})

				if networkCardIndex, ok := attachmentArg["NetworkCardIndex"]; ok {
					d.Set("network_card_index", networkCardIndex)
				}
			}
		} else {
			networkInterfaceMap := make(map[string]interface{})
			networkInterfaceMap["network_interface_id"] = obj.NetworkInterfaceId

			object, err := ecsService.DescribeEcsNetworkInterface(obj.NetworkInterfaceId)
			if err != nil {
				return WrapError(err)
			}

			networkInterfaceMap["vswitch_id"] = object["VSwitchId"]
			networkInterfaceMap["network_interface_traffic_mode"] = object["NetworkInterfaceTrafficMode"]
			networkInterfaceMap["queue_pair_number"] = object["QueuePairNumber"]
			networkInterfaceMap["security_group_ids"] = object["SecurityGroupIds"].(map[string]interface{})["SecurityGroupId"]

			if attachment, ok := object["Attachment"]; ok {
				attachmentArg := attachment.(map[string]interface{})

				if networkCardIndex, ok := attachmentArg["NetworkCardIndex"]; ok {
					networkInterfaceMap["network_card_index"] = networkCardIndex
				}
			}

			networkInterfaceMaps = append(networkInterfaceMaps, networkInterfaceMap)
		}
	}
	d.Set("network_interfaces", networkInterfaceMaps)

	if len(networkInterfaceId) != 0 {
		object, err := ecsService.DescribeEcsNetworkInterface(networkInterfaceId)
		if err != nil {
			return WrapError(err)
		}
		secondaryPrivateIpsSli := make([]interface{}, 0, len(object["PrivateIpSets"].(map[string]interface{})["PrivateIpSet"].([]interface{})))
		for _, v := range object["PrivateIpSets"].(map[string]interface{})["PrivateIpSet"].([]interface{}) {
			if !v.(map[string]interface{})["Primary"].(bool) {
				secondaryPrivateIpsSli = append(secondaryPrivateIpsSli, v.(map[string]interface{})["PrivateIpAddress"])
			}
		}
		ipv6SetList := make([]interface{}, 0)
		for _, v := range object["Ipv6Sets"].(map[string]interface{})["Ipv6Set"].([]interface{}) {
			ipv6Set := v.(map[string]interface{})
			ipv6SetList = append(ipv6SetList, ipv6Set["Ipv6Address"])
		}

		d.Set("network_interface_id", networkInterfaceId)
		d.Set("ipv6_addresses", ipv6SetList)
		d.Set("ipv6_address_count", len(ipv6SetList))
		d.Set("secondary_private_ips", secondaryPrivateIpsSli)
		d.Set("secondary_private_ip_address_count", len(secondaryPrivateIpsSli))
	}

	maintenanceAttribute, err := ecsService.DescribeInstanceMaintenanceAttribute(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if v, ok := maintenanceAttribute["MaintenanceWindows"]; ok {
		maintenanceWindowsMaps := make([]map[string]interface{}, 0)
		maintenanceWindowsList := v.(map[string]interface{})["MaintenanceWindow"].([]interface{})
		maintenanceWindowsMap := make(map[string]interface{})
		for _, maintenanceWindowsItem := range maintenanceWindowsList {
			if maintenanceWindowsItemArg, ok := maintenanceWindowsItem.(map[string]interface{}); ok {
				maintenanceWindowsMap["start_time"] = maintenanceWindowsItemArg["StartTime"]
				maintenanceWindowsMap["end_time"] = maintenanceWindowsItemArg["EndTime"]
				maintenanceWindowsMaps = append(maintenanceWindowsMaps, maintenanceWindowsMap)
			}
		}
		d.Set("maintenance_time", maintenanceWindowsMaps)
	}

	if v, ok := maintenanceAttribute["ActionOnMaintenance"]; ok {
		d.Set("maintenance_action", v.(map[string]interface{})["Value"])
	}

	d.Set("maintenance_notify", maintenanceAttribute["NotifyOnMaintenance"])

	instanceAttribute, err := ecsService.DescribeInstanceAttribute(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("enable_jumbo_frame", instanceAttribute.EnableJumboFrame)

	return nil
}

func resourceAliCloudInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)

	if !d.IsNewResource() {
		if err := setTags(client, TagResourceInstance, d); err != nil {
			return WrapError(err)
		} else {
			d.SetPartial("tags")
		}
	}

	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		action := "JoinResourceGroup"
		request := map[string]interface{}{
			"ResourceType":    "instance",
			"ResourceId":      d.Id(),
			"RegionId":        client.RegionId,
			"ResourceGroupId": d.Get("resource_group_id"),
		}
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		d.SetPartial("resource_group_id")
	}

	if err := setVolumeTags(client, TagResourceDisk, d); err != nil {
		return WrapError(err)
	} else {
		d.SetPartial("volume_tags")
	}

	if !d.IsNewResource() && !d.HasChange("vpc_id") && d.HasChange("security_groups") {
		if !d.IsNewResource() || d.Get("vswitch_id").(string) == "" {
			o, n := d.GetChange("security_groups")
			os := o.(*schema.Set)
			ns := n.(*schema.Set)

			rl := expandStringList(os.Difference(ns).List())
			al := expandStringList(ns.Difference(os).List())

			if len(rl) > 0 {
				err := ecsService.LeaveSecurityGroups(d.Id(), rl)
				if err != nil {
					return WrapError(err)
				}
			}

			if len(al) > 0 {
				err := ecsService.JoinSecurityGroups(d.Id(), al)
				if err != nil {
					return WrapError(err)
				}
			}

			d.SetPartial("security_groups")
		}
	}

	if !d.IsNewResource() && (d.HasChange("system_disk_size") || d.HasChange("system_disk_auto_snapshot_policy_id") || d.HasChange("system_disk_name") || d.HasChange("system_disk_description")) {
		disk, err := ecsService.DescribeEcsSystemDisk(d.Id())
		if err != nil {
			return WrapError(err)
		}
		if d.HasChange("system_disk_size") {
			instance, errDesc := ecsService.DescribeInstance(d.Id())
			if errDesc != nil {
				return WrapError(errDesc)
			}

			request := ecs.CreateResizeDiskRequest()
			request.NewSize = requests.NewInteger(d.Get("system_disk_size").(int))
			if instance.Status == string(Stopped) {
				request.Type = "offline"
			} else {
				request.Type = "online"
			}
			request.DiskId = disk["DiskId"].(string)
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ResizeDisk(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			d.SetPartial("system_disk_size")
		}

		if d.HasChange("system_disk_auto_snapshot_policy_id") {
			action := "ApplyAutoSnapshotPolicy"
			var response map[string]interface{}
			request := map[string]interface{}{
				"RegionId": client.RegionId,
			}
			request["autoSnapshotPolicyId"] = d.Get("system_disk_auto_snapshot_policy_id")
			request["diskIds"] = convertListToJsonString([]interface{}{disk["DiskId"]})
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
				return nil
			})
			addDebug(action, response, request)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			d.SetPartial("system_disk_auto_snapshot_policy_id")
		}

		if d.HasChange("system_disk_name") || d.HasChange("system_disk_description") {
			var response map[string]interface{}
			modifyDiskAttributeReq := map[string]interface{}{
				"DiskId": disk["DiskId"],
			}
			modifyDiskAttributeReq["DiskName"] = d.Get("system_disk_name")
			modifyDiskAttributeReq["Description"] = d.Get("system_disk_description")
			action := "ModifyDiskAttribute"
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, modifyDiskAttributeReq, &util.RuntimeOptions{})
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, modifyDiskAttributeReq)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			d.SetPartial("system_disk_name")
			d.SetPartial("system_disk_description")
		}
	}

	run := false
	imageUpdate, err := modifyInstanceImage(d, meta, run)
	if err != nil {
		return WrapError(err)
	}

	vpcUpdate, err := modifyVpcAttribute(d, meta, run)
	if err != nil {
		return WrapError(err)
	}

	passwordUpdate, err := modifyInstanceAttribute(d, meta)
	if err != nil {
		return WrapError(err)
	}

	if !d.IsNewResource() && d.HasChange("auto_release_time") {
		request := ecs.CreateModifyInstanceAutoReleaseTimeRequest()
		request.InstanceId = d.Id()
		request.RegionId = client.RegionId
		request.AutoReleaseTime = d.Get("auto_release_time").(string)
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifyInstanceAutoReleaseTime(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("auto_release_time")
	}

	typeUpdate, err := modifyInstanceType(d, meta, run)
	if err != nil {
		return WrapError(err)
	}
	target, targetExist := d.GetOk("status")
	statusUpdate := d.HasChange("status")
	if d.IsNewResource() && targetExist && target.(string) == string(Running) {
		statusUpdate = false
	}
	if imageUpdate || vpcUpdate || passwordUpdate || typeUpdate || statusUpdate {
		run = true
		instance, errDesc := ecsService.DescribeInstance(d.Id())
		if errDesc != nil {
			return WrapError(errDesc)
		}
		if statusUpdate && targetExist && target == string(Stopped) || instance.Status == string(Running) {
			stopRequest := ecs.CreateStopInstanceRequest()
			stopRequest.RegionId = client.RegionId
			stopRequest.InstanceId = d.Id()
			stopRequest.ForceStop = requests.NewBoolean(false)
			if v, ok := d.GetOk("stopped_mode"); ok {
				stopRequest.StoppedMode = v.(string)
			}
			err := resource.Retry(5*time.Minute, func() *resource.RetryError {
				raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
					return ecsClient.StopInstance(stopRequest)
				})
				if err != nil {
					if IsExpectedErrors(err, []string{"IncorrectInstanceStatus"}) {
						time.Sleep(time.Second)
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(stopRequest.GetActionName(), raw)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), stopRequest.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			stateConf := BuildStateConf([]string{"Pending", "Running", "Stopping"}, []string{"Stopped"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ecsService.InstanceStateRefreshFunc(d.Id(), []string{}))

			if _, err = stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
		if _, err := modifyInstanceImage(d, meta, run); err != nil {
			return WrapError(err)
		}

		if _, err := modifyVpcAttribute(d, meta, run); err != nil {
			return WrapError(err)
		}

		if _, err := modifyInstanceType(d, meta, run); err != nil {
			return WrapError(err)
		}

		if targetExist && target == string(Running) {
			startRequest := ecs.CreateStartInstanceRequest()
			startRequest.InstanceId = d.Id()

			err := resource.Retry(5*time.Minute, func() *resource.RetryError {
				raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
					return ecsClient.StartInstance(startRequest)
				})
				if err != nil {
					if IsExpectedErrors(err, []string{"IncorrectInstanceStatus"}) {
						time.Sleep(time.Second)
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(startRequest.GetActionName(), raw)
				return nil
			})

			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), startRequest.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			// Start instance sometimes costs more than 8 minutes when os type is centos.
			stateConf := &resource.StateChangeConf{
				Pending:    []string{"Pending", "Starting", "Stopped"},
				Target:     []string{"Running"},
				Refresh:    ecsService.InstanceStateRefreshFunc(d.Id(), []string{}),
				Timeout:    d.Timeout(schema.TimeoutUpdate),
				Delay:      5 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			if _, err = stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
		if d.HasChange("status") {
			d.SetPartial("status")
		}
	}

	if err := modifyInstanceNetworkSpec(d, meta); err != nil {
		return WrapError(err)
	}

	if d.HasChange("force_delete") {
		d.SetPartial("force_delete")
	}

	if err := modifyInstanceChargeType(d, meta, false); err != nil {
		return WrapError(err)
	}

	// Only PrePaid instance can support modifying renewal attribute
	if !d.IsNewResource() && d.Get("instance_charge_type").(string) == string(PrePaid) &&
		(d.HasChange("renewal_status") || d.HasChange("auto_renew_period")) {
		status := d.Get("renewal_status").(string)
		request := ecs.CreateModifyInstanceAutoRenewAttributeRequest()
		request.InstanceId = d.Id()
		request.RenewalStatus = status

		if status == string(RenewAutoRenewal) {
			request.PeriodUnit = d.Get("period_unit").(string)
			request.Duration = requests.NewInteger(d.Get("auto_renew_period").(int))
		}

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifyInstanceAutoRenewAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("renewal_status")
		d.SetPartial("auto_renew_period")
	}

	if d.HasChange("secondary_private_ips") {
		client := meta.(*connectivity.AliyunClient)
		ecsService := EcsService{client}
		var response map[string]interface{}
		instance, err := ecsService.DescribeInstance(d.Id())
		if err != nil {
			return WrapError(err)
		}
		networkInterfaceId := ""
		for _, obj := range instance.NetworkInterfaces.NetworkInterface {
			if obj.Type == "Primary" {
				networkInterfaceId = obj.NetworkInterfaceId
				break
			}
		}
		oraw, nraw := d.GetChange("secondary_private_ips")
		remove := oraw.(*schema.Set).Difference(nraw.(*schema.Set)).List()
		create := nraw.(*schema.Set).Difference(oraw.(*schema.Set)).List()
		if len(remove) > 0 {
			action := "UnassignPrivateIpAddresses"
			request := map[string]interface{}{
				"RegionId":           client.RegionId,
				"NetworkInterfaceId": networkInterfaceId,
				"ClientToken":        buildClientToken(action),
			}

			for index, val := range remove {
				request[fmt.Sprintf("PrivateIpAddress.%d", index+1)] = val
			}

			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					if NeedRetry(err) || IsExpectedErrors(err, []string{"OperationConflict", "Operation.Conflict"}) {
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
			addDebug(action, response, request)
			d.SetPartial("secondary_private_ips")
		}
		if len(create) > 0 {
			action := "AssignPrivateIpAddresses"
			request := map[string]interface{}{
				"RegionId":           client.RegionId,
				"NetworkInterfaceId": networkInterfaceId,
				"ClientToken":        buildClientToken(action),
			}
			for index, val := range create {
				request[fmt.Sprintf("PrivateIpAddress.%d", index+1)] = val
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					if NeedRetry(err) || IsExpectedErrors(err, []string{"OperationConflict", "Operation.Conflict"}) {
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
			d.SetPartial("secondary_private_ips")
		}

	}

	if d.HasChange("secondary_private_ip_address_count") {
		client := meta.(*connectivity.AliyunClient)
		ecsService := EcsService{client}
		var response map[string]interface{}
		instance, err := ecsService.DescribeInstance(d.Id())
		if err != nil {
			return WrapError(err)
		}
		// query for the Primary NetworkInterfaceId
		networkInterfaceId := ""
		for _, obj := range instance.NetworkInterfaces.NetworkInterface {
			if obj.Type == "Primary" {
				networkInterfaceId = obj.NetworkInterfaceId
				break
			}
		}
		privateIpList := expandStringList(d.Get("secondary_private_ips").(*schema.Set).List())
		oldIpsCount, newIpsCount := d.GetChange("secondary_private_ip_address_count")
		if oldIpsCount != nil && newIpsCount != nil && newIpsCount != len(privateIpList) {
			diff := newIpsCount.(int) - oldIpsCount.(int)
			if diff > 0 {
				action := "AssignPrivateIpAddresses"
				request := map[string]interface{}{
					"RegionId":                       client.RegionId,
					"NetworkInterfaceId":             networkInterfaceId,
					"ClientToken":                    buildClientToken(action),
					"SecondaryPrivateIpAddressCount": diff,
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
					if err != nil {
						if NeedRetry(err) || IsExpectedErrors(err, []string{"OperationConflict"}) {
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
				d.SetPartial("secondary_private_ip_address_count")
			}
			if diff < 0 {
				diff *= -1
				action := "UnassignPrivateIpAddresses"
				request := map[string]interface{}{
					"RegionId":           client.RegionId,
					"NetworkInterfaceId": networkInterfaceId,
					"ClientToken":        buildClientToken(action),
				}
				for index, val := range privateIpList[:diff] {
					request[fmt.Sprintf("PrivateIpAddress.%d", index+1)] = val
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
					if err != nil {
						if NeedRetry(err) || IsExpectedErrors(err, []string{"OperationConflict"}) {
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
				addDebug(action, response, request)
				d.SetPartial("secondary_private_ip_address_count")
			}
		}
	}
	if !d.IsNewResource() && d.HasChange("deployment_set_id") {
		action := "ModifyInstanceDeployment"
		var response map[string]interface{}
		request := map[string]interface{}{
			"RegionId":    client.RegionId,
			"InstanceId":  d.Id(),
			"ClientToken": buildClientToken(action),
		}
		if v, ok := d.GetOk("deployment_set_id"); ok {
			request["DeploymentSetId"] = v
		}
		if v := d.Get("deployment_set_id"); len(v.(string)) == 0 {
			oldDeploymentSetId, _ := d.GetChange("deployment_set_id")
			request["DeploymentSetId"] = oldDeploymentSetId
			request["RemoveFromDeploymentSet"] = true
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
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("deployment_set_id")
	}

	if d.HasChange("maintenance_time") || d.HasChange("maintenance_action") || d.HasChange("maintenance_notify") {
		var response map[string]interface{}
		action := "ModifyInstanceMaintenanceAttributes"
		request := map[string]interface{}{
			"RegionId":   client.RegionId,
			"InstanceId": []string{d.Id()},
		}

		maintenanceWindowsMaps := make([]map[string]interface{}, 0)
		for _, maintenanceWindows := range d.Get("maintenance_time").(*schema.Set).List() {
			maintenanceWindowsMap := make(map[string]interface{})
			maintenanceWindowsArg := maintenanceWindows.(map[string]interface{})

			if v, ok := maintenanceWindowsArg["start_time"].(string); ok && v != "" {
				maintenanceWindowsMap["StartTime"] = v
			}
			if v, ok := maintenanceWindowsArg["end_time"].(string); ok && v != "" {
				maintenanceWindowsMap["EndTime"] = v
			}

			maintenanceWindowsMaps = append(maintenanceWindowsMaps, maintenanceWindowsMap)
		}
		request["MaintenanceWindow"] = maintenanceWindowsMaps

		if v, ok := d.GetOk("maintenance_action"); ok {
			request["ActionOnMaintenance"] = v
		}
		if v, ok := d.GetOkExists("maintenance_notify"); ok {
			request["NotifyOnMaintenance"] = v
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
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("maintenance_time")
		d.SetPartial("maintenance_action")
		d.SetPartial("maintenance_notify")
	}

	if d.HasChange("http_endpoint") || d.HasChange("http_tokens") {
		var response map[string]interface{}
		action := "ModifyInstanceMetadataOptions"
		request := map[string]interface{}{
			"RegionId":   client.RegionId,
			"InstanceId": d.Id(),
		}

		if v, ok := d.GetOk("http_endpoint"); ok {
			request["HttpEndpoint"] = v
		} else {
			request["HttpEndpoint"] = "enabled"
		}
		if v, ok := d.GetOk("http_tokens"); ok {
			request["HttpTokens"] = v
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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

		d.SetPartial("http_endpoint")
		d.SetPartial("http_tokens")
	}

	d.Partial(false)
	return resourceAliCloudInstanceRead(d, meta)
}

func resourceAliCloudInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	if d.Get("instance_charge_type").(string) == string(PrePaid) {
		force := d.Get("force_delete").(bool)
		if !force {
			return WrapError(Error("Please convert 'PrePaid' instance to 'PostPaid' or set 'force_delete' as true before deleting 'PrePaid' instance."))
		} else if err := modifyInstanceChargeType(d, meta, force); err != nil {
			return WrapError(err)
		}
	}
	stopRequest := ecs.CreateStopInstanceRequest()
	stopRequest.InstanceId = d.Id()
	stopRequest.ForceStop = requests.NewBoolean(true)

	deleteRequest := ecs.CreateDeleteInstanceRequest()
	deleteRequest.InstanceId = d.Id()
	deleteRequest.Force = requests.NewBoolean(true)

	wait := incrementalWait(1*time.Second, 1*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteInstance(deleteRequest)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectInstanceStatus", "DependencyViolation.RouteEntry", "IncorrectInstanceStatus.Initializing"}) {
				return resource.RetryableError(err)
			}
			if IsExpectedErrors(err, []string{Throttling, "LastTokenProcessing"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(deleteRequest.GetActionName(), raw)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, EcsNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), deleteRequest.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"Pending", "Running", "Stopped", "Stopping"}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Second, ecsService.InstanceStateRefreshFunc(d.Id(), []string{}))

	if _, err = stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func modifyInstanceChargeType(d *schema.ResourceData, meta interface{}, forceDelete bool) error {
	if d.IsNewResource() {
		d.Partial(false)
		return nil
	}

	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	chargeType := d.Get("instance_charge_type").(string)
	if d.HasChange("instance_charge_type") || forceDelete {
		if forceDelete {
			chargeType = string(PostPaid)
		}
		request := ecs.CreateModifyInstanceChargeTypeRequest()
		request.InstanceIds = convertListToJsonString(append(make([]interface{}, 0, 1), d.Id()))
		request.IncludeDataDisks = requests.NewBoolean(d.Get("include_data_disks").(bool))
		request.AutoPay = requests.NewBoolean(true)
		request.DryRun = requests.NewBoolean(d.Get("dry_run").(bool))
		request.ClientToken = fmt.Sprintf("terraform-modify-instance-charge-type-%s", d.Id())
		if chargeType == string(PrePaid) {
			if v, ok := d.GetOk("period"); ok {
				request.Period = requests.NewInteger(v.(int))
			}
			request.PeriodUnit = d.Get("period_unit").(string)
		}
		request.InstanceChargeType = chargeType
		if err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ModifyInstanceChargeType(request)
			})
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"InternalError"}) {
					time.Sleep(3 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		// Wait for instance charge type has been changed
		if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			if instance, err := ecsService.DescribeInstance(d.Id()); err != nil {
				return resource.NonRetryableError(err)
			} else if instance.InstanceChargeType == chargeType {
				return nil
			}
			return resource.RetryableError(Error("Waitting for instance %s to be %s timeout.", d.Id(), chargeType))
		}); err != nil {
			return WrapError(err)
		}

		d.SetPartial("instance_charge_type")
		return nil
	}

	return nil
}

func modifyInstanceImage(d *schema.ResourceData, meta interface{}, run bool) (bool, error) {
	if d.IsNewResource() {
		d.Partial(false)
		return false, nil
	}
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	update := false
	if d.HasChange("image_id") {
		update = true
		if !run {
			return update, nil
		}
		instance, err := ecsService.DescribeInstance(d.Id())
		if err != nil {
			return update, WrapError(err)
		}
		keyPairName := instance.KeyPairName
		request := ecs.CreateReplaceSystemDiskRequest()
		request.InstanceId = d.Id()
		request.ImageId = d.Get("image_id").(string)
		request.SystemDiskSize = requests.NewInteger(d.Get("system_disk_size").(int))
		if v, ok := d.GetOk("system_disk_encrypted"); ok {
			request.Encrypted = requests.NewBoolean(v.(bool))
		}
		if v, ok := d.GetOk("system_disk_encrypt_algorithm"); ok {
			request.EncryptAlgorithm = v.(string)
		}
		if v, ok := d.GetOk("security_enhancement_strategy"); ok {
			request.SecurityEnhancementStrategy = v.(string)
		}
		if v, ok := d.GetOk("password"); ok {
			request.Password = v.(string)
		}
		if v, ok := d.GetOkExists("password_inherit"); ok {
			request.PasswordInherit = requests.NewBoolean(v.(bool))
		}
		if v, ok := d.GetOk("system_disk_kms_key_id"); ok {
			request.KMSKeyId = v.(string)
		}
		request.ClientToken = buildClientToken(request.GetActionName())
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ReplaceSystemDisk(request)
		})
		if err != nil {
			return update, WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		response, _ := raw.(*ecs.ReplaceSystemDiskResponse)
		d.Set("system_disk_id", response.DiskId)
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		// Ensure instance's image has been replaced successfully.
		timeout := DefaultTimeoutMedium
		for {
			instance, errDesc := ecsService.DescribeInstance(d.Id())
			if errDesc != nil {
				return update, WrapError(errDesc)
			}
			disk, err := ecsService.DescribeInstanceSystemDisk(d.Id(), instance.ResourceGroupId, d.Get("system_disk_id").(string))
			if err != nil {
				return update, WrapError(err)
			}

			if instance.ImageId == d.Get("image_id") && disk.Size == d.Get("system_disk_size").(int) {
				break
			}
			time.Sleep(DefaultIntervalShort * time.Second)

			timeout = timeout - DefaultIntervalShort
			if timeout <= 0 {
				return update, WrapError(GetTimeErrorFromString(fmt.Sprintf("Replacing instance %s system disk timeout.", d.Id())))
			}
		}

		// After updating image, it need to re-attach key pair
		if keyPairName != "" {
			if err := ecsService.AttachKeyPair(keyPairName, []interface{}{d.Id()}); err != nil {
				return update, WrapError(err)
			}
		}
	}
	return update, nil
}

func modifyInstanceAttribute(d *schema.ResourceData, meta interface{}) (bool, error) {
	if d.IsNewResource() {
		d.Partial(false)
		return false, nil
	}

	update := false
	reboot := false
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "ModifyInstanceAttribute"
	request := make(map[string]interface{})
	conn, err := client.NewEcsClient()
	if err != nil {
		return reboot, WrapError(err)
	}

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	request["InstanceId"] = d.Id()

	if d.HasChange("instance_name") {
		d.SetPartial("instance_name")
		request["InstanceName"] = d.Get("instance_name").(string)
		update = true
	}

	if d.HasChange("description") {
		d.SetPartial("description")
		request["Description"] = d.Get("description").(string)
		update = true
	}

	if d.HasChange("user_data") {
		d.SetPartial("user_data")
		v := d.Get("user_data")
		_, base64DecodeError := base64.StdEncoding.DecodeString(v.(string))
		if base64DecodeError == nil {
			request["UserData"] = v.(string)
		} else {
			request["UserData"] = base64.StdEncoding.EncodeToString([]byte(v.(string)))
		}

		update = true
		reboot = true
	}

	if d.HasChange("host_name") {
		d.SetPartial("host_name")
		request["HostName"] = d.Get("host_name").(string)
		update = true
		reboot = true
	}

	if d.HasChange("password") || d.HasChange("kms_encrypted_password") {
		if v := d.Get("password").(string); v != "" {
			d.SetPartial("password")
			request["Password"] = v
			update = true
			reboot = true
		}
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{meta.(*connectivity.AliyunClient)}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return reboot, WrapError(err)
			}
			request["Password"] = decryptResp
			d.SetPartial("kms_encrypted_password")
			d.SetPartial("kms_encryption_context")
			update = true
			reboot = true
		}
	}

	if d.HasChange("deletion_protection") {
		d.SetPartial("deletion_protection")
		request["DeletionProtection"] = requests.NewBoolean(d.Get("deletion_protection").(bool))
		update = true
	}

	if d.HasChange("credit_specification") {
		d.SetPartial("credit_specification")
		request["CreditSpecification"] = d.Get("credit_specification").(string)
		update = true
	}

	if d.HasChange("enable_jumbo_frame") {
		d.SetPartial("enable_jumbo_frame")

		if v, ok := d.GetOkExists("enable_jumbo_frame"); ok {
			request["EnableJumboFrame"] = v
		}

		update = true
	}

	if update {
		wait := incrementalWait(1*time.Minute, 1*time.Minute)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"InvalidChargeType.ValueNotSupported"}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return reboot, WrapErrorf(err, DefaultErrorMsg, "alicloud_instance", action, AlibabaCloudSdkGoERROR)
		}
	}

	return reboot, nil
}

func modifyVpcAttribute(d *schema.ResourceData, meta interface{}, run bool) (bool, error) {
	if d.IsNewResource() {
		d.Partial(false)
		return false, nil
	}

	update := false
	request := ecs.CreateModifyInstanceVpcAttributeRequest()
	request.InstanceId = d.Id()
	request.VSwitchId = d.Get("vswitch_id").(string)

	if d.HasChange("vswitch_id") {
		update = true
		if d.Get("vswitch_id").(string) == "" {
			return update, WrapError(Error("Field 'vswitch_id' is required when modifying the instance VPC attribute."))
		}
	}

	if d.HasChange("subnet_id") {
		update = true
		if d.Get("subnet_id").(string) == "" {
			return update, WrapError(Error("Field 'subnet_id' is required when modifying the instance VPC attribute."))
		}
		request.VSwitchId = d.Get("subnet_id").(string)
	}

	if request.VSwitchId != "" && d.HasChange("private_ip") {
		request.PrivateIpAddress = d.Get("private_ip").(string)
		update = true
	}

	if d.HasChange("vpc_id") {
		update = true

		if v, ok := d.GetOk("vpc_id"); ok {
			request.VpcId = v.(string)
		}

		if v, ok := d.GetOk("security_groups"); ok {
			securityGroupIds := expandStringList(v.(*schema.Set).List())
			if len(securityGroupIds) > 0 {
				request.SecurityGroupId = &securityGroupIds
			}
		}
	}

	if !run {
		return update, nil
	}

	if update {
		client := meta.(*connectivity.AliyunClient)
		err := resource.Retry(1*time.Minute, func() *resource.RetryError {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ModifyInstanceVpcAttribute(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationConflict"}) {
					time.Sleep(1 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})

		if err != nil {
			return update, WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		ecsService := EcsService{client}
		if err := ecsService.WaitForVpcAttributesChanged(d.Id(), request.VSwitchId, request.PrivateIpAddress); err != nil {
			return update, WrapError(err)
		}

		d.SetPartial("vswitch_id")
		d.SetPartial("subnet_id")
		d.SetPartial("private_ip")
		d.SetPartial("vpc_id")
		d.SetPartial("security_groups")
	}

	return update, nil
}

func modifyInstanceType(d *schema.ResourceData, meta interface{}, run bool) (bool, error) {
	if d.IsNewResource() {
		d.Partial(false)
		return false, nil
	}
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	update := false
	if d.HasChange("instance_type") {
		update = true
		if !run {
			return update, nil
		}
		// Ensure instance_type is valid
		//zoneId, validZones, err := ecsService.DescribeAvailableResources(d, meta, InstanceTypeResource)
		//if err != nil {
		//	return update, WrapError(err)
		//}
		//if err = ecsService.InstanceTypeValidation(d.Get("instance_type").(string), zoneId, validZones); err != nil {
		//	return update, WrapError(err)
		//}

		// There should use the old instance charge type to decide API method because of instance_charge_type will be updated at last step
		oldCharge, _ := d.GetChange("instance_charge_type")
		if oldCharge.(string) == string(PrePaid) {
			request := ecs.CreateModifyPrepayInstanceSpecRequest()
			request.InstanceId = d.Id()
			request.InstanceType = d.Get("instance_type").(string)
			if v, ok := d.GetOk("operator_type"); ok {
				request.OperatorType = v.(string)
			}

			err := resource.Retry(6*time.Minute, func() *resource.RetryError {
				raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
					return ecsClient.ModifyPrepayInstanceSpec(request)
				})
				if err != nil {
					if IsExpectedErrors(err, []string{Throttling}) {
						time.Sleep(5 * time.Second)
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			})
			if err != nil {
				return update, WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
		} else {
			//An instance that was successfully modified once cannot be modified again within 5 minutes.
			request := ecs.CreateModifyInstanceSpecRequest()
			request.InstanceId = d.Id()
			request.InstanceType = d.Get("instance_type").(string)
			request.ClientToken = buildClientToken(request.GetActionName())

			err := resource.Retry(6*time.Minute, func() *resource.RetryError {
				args := *request
				raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
					return ecsClient.ModifyInstanceSpec(&args)
				})
				if err != nil {
					if IsExpectedErrors(err, []string{Throttling}) {
						time.Sleep(10 * time.Second)
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			})
			if err != nil {
				return update, WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
		}

		// Ensure instance's type has been replaced successfully.
		timeout := DefaultTimeoutMedium
		for {
			instance, err := ecsService.DescribeInstance(d.Id())

			if err != nil {
				return update, WrapError(err)
			}

			if instance.InstanceType == d.Get("instance_type").(string) {
				break
			}

			timeout = timeout - DefaultIntervalShort
			if timeout <= 0 {
				return update, WrapErrorf(err, WaitTimeoutMsg, d.Id(), GetFunc(1), timeout, instance.InstanceType, d.Get("instance_type"), ProviderERROR)
			}

			time.Sleep(DefaultIntervalShort * time.Second)
		}
		d.SetPartial("instance_type")
	}
	return update, nil
}

func modifyInstanceNetworkSpec(d *schema.ResourceData, meta interface{}) error {
	if d.IsNewResource() {
		d.Partial(false)
		return nil
	}

	allocate := false
	update := false
	request := ecs.CreateModifyInstanceNetworkSpecRequest()
	request.InstanceId = d.Id()
	request.ClientToken = buildClientToken(request.GetActionName())

	if d.HasChange("internet_charge_type") {
		request.NetworkChargeType = d.Get("internet_charge_type").(string)
		update = true
		d.SetPartial("internet_charge_type")
	}

	if d.HasChange("internet_max_bandwidth_out") {
		o, n := d.GetChange("internet_max_bandwidth_out")
		if o.(int) <= 0 && n.(int) > 0 {
			allocate = true
		}
		request.InternetMaxBandwidthOut = requests.NewInteger(n.(int))
		update = true
		d.SetPartial("internet_max_bandwidth_out")
	}

	if d.HasChange("internet_max_bandwidth_in") {
		request.InternetMaxBandwidthIn = requests.NewInteger(d.Get("internet_max_bandwidth_in").(int))
		update = true
		d.SetPartial("internet_max_bandwidth_in")
	}

	//An instance that was successfully modified once cannot be modified again within 5 minutes.
	wait := incrementalWait(2*time.Second, 2*time.Second)
	client := meta.(*connectivity.AliyunClient)
	if update {
		if err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ModifyInstanceNetworkSpec(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{Throttling, "LastOrderProcessing", "LastRequestProcessing", "LastTokenProcessing"}) {
					wait()
					return resource.RetryableError(err)
				}
				if IsExpectedErrors(err, []string{"InternalError"}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		ecsService := EcsService{client: client}

		deadline := time.Now().Add(DefaultTimeout * time.Second)
		for {
			instance, err := ecsService.DescribeInstance(d.Id())
			if err != nil {
				return WrapError(err)
			}

			if instance.InternetMaxBandwidthOut == d.Get("internet_max_bandwidth_out").(int) &&
				instance.InternetChargeType == d.Get("internet_charge_type").(string) {
				break
			}

			if time.Now().After(deadline) {
				return WrapError(Error(`wait for internet update timeout! expect internet_charge_type value %s, get %s
					expect internet_max_bandwidth_out value %d, get %d,`,
					d.Get("internet_charge_type").(string), instance.InternetChargeType, d.Get("internet_max_bandwidth_out").(int),
					instance.InternetMaxBandwidthOut))
			}
			time.Sleep(1 * time.Second)
		}

		if allocate {
			request := ecs.CreateAllocatePublicIpAddressRequest()
			request.InstanceId = d.Id()
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.AllocatePublicIpAddress(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}
	}
	return nil
}
