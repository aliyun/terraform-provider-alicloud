package alicloud

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	roacs "github.com/alibabacloud-go/cs-20151215/v5/client"
	"github.com/denverdino/aliyungo/cs"

	"github.com/PaesslerAG/jsonpath"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudAckNodepool() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAckNodepoolCreate,
		Read:   resourceAliCloudAckNodepoolRead,
		Update: resourceAliCloudAckNodepoolUpdate,
		Delete: resourceAliCloudAckNodepoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_renew": {
				Type:             schema.TypeBool,
				Optional:         true,
				DiffSuppressFunc: csNodepoolInstancePostPaidDiffSuppressFunc,
			},
			"auto_renew_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     IntInSlice([]int{0, 1, 2, 3, 6, 12}),
				DiffSuppressFunc: csNodepoolInstancePostPaidDiffSuppressFunc,
			},
			"cis_enabled": {
				Type:       schema.TypeBool,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Field 'cis_enabled' has been deprecated from provider version 1.223.1. Whether enable worker node to support cis security reinforcement, its valid value `true` or `false`. Default to `false` and apply to AliyunLinux series. Use `security_hardening_os` instead.",
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"compensate_with_on_demand": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cpu_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"static", "none"}, false),
			},
			"data_disks": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bursting_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"category": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"cloud_efficiency", "cloud_ssd", "cloud_essd", "cloud_auto", "cloud", "cloud_essd_xc0", "cloud_essd_xc1", "all", "ephemeral_ssd", "local_disk", "cloud_essd_entry", "elastic_ephemeral_disk_premium", "elastic_ephemeral_disk_standard"}, false),
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"performance_level": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"encrypted": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"size": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: IntBetween(40, 32767),
						},
						"device": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"auto_snapshot_policy_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"mount_target": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"auto_format": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"provisioned_iops": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"file_system": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"deployment_set_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"desired_size": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"instances", "node_count"},
			},
			"eflo_node_group": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"force_delete": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"image_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"AliyunLinux", "AliyunLinux3", "AliyunLinux3Arm64", "AliyunLinuxUEFI", "CentOS", "Windows", "WindowsCore", "AliyunLinux Qboot", "ContainerOS", "AliyunLinuxSecurity", "Ubuntu", "AliyunLinux3ContainerOptimized", "Custom"}, false),
			},
			"install_cloud_monitor": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"PrePaid", "PostPaid"}, false),
			},
			"instance_patterns": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cores": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"cpu_architectures": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"max_cpu_cores": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"memory": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"instance_categories": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"instance_family_level": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"max_memory_size": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"min_cpu_cores": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"min_memory_size": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"excluded_instance_types": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"instance_type_families": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"instance_types": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"PayByBandwidth", "PayByTraffic"}, false),
			},
			"internet_max_bandwidth_out": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 100),
			},
			"key_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"password", "kms_encrypted_password"},
			},
			"kubelet_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu_manager_policy": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"allowed_unsafe_sysctls": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"topology_manager_policy": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"pod_pids_limit": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cpu_cfs_quota": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"serialize_image_pulls": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"kube_api_burst": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cluster_dns": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"system_reserved": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"feature_gates": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeBool},
						},
						"registry_burst": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"read_only_port": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"registry_pull_qps": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"reserved_memory": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"limits": {
										Type:     schema.TypeMap,
										Optional: true,
									},
									"numa_node": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"container_log_monitor_interval": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"container_log_max_workers": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"event_burst": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"image_gc_high_threshold_percent": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"eviction_soft_grace_period": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"image_gc_low_threshold_percent": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"memory_manager_policy": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cpu_cfs_quota_period": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"eviction_soft": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"event_record_qps": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"kube_reserved": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"max_pods": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"container_log_max_files": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"eviction_hard": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"tracing": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sampling_rate_per_million": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"endpoint": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"container_log_max_size": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"kube_api_qps": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"labels": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"login_as_non_root": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"management": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_upgrade_policy": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auto_upgrade_kubelet": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"auto_repair": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"auto_upgrade": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"surge_percentage": {
							Type:       schema.TypeInt,
							Optional:   true,
							Deprecated: "Field 'surge_percentage' has been deprecated from provider version 1.219.0. Proportion of additional nodes. You have to specify one of surge, surge_percentage.",
						},
						"surge": {
							Type:       schema.TypeInt,
							Optional:   true,
							Deprecated: "Field 'surge' has been deprecated from provider version 1.219.0. Number of additional nodes. You have to specify one of surge, surge_percentage.",
						},
						"auto_vul_fix_policy": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"restart_node": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"vul_level": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"enable": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"auto_repair_policy": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"restart_node": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"auto_vul_fix": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"max_unavailable": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: IntBetween(0, 1000),
						},
					},
				},
			},
			"multi_az_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"PRIORITY", "COST_OPTIMIZED", "BALANCE"}, false),
			},
			"node_name_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"node_pool_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_pool_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"node_pool_name", "name"},
				Computed:     true,
			},
			"on_demand_base_capacity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"on_demand_percentage_above_base_capacity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"key_name", "kms_encrypted_password"},
			},
			"period": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     IntInSlice([]int{0, 1, 2, 3, 6, 12}),
				DiffSuppressFunc: csNodepoolInstancePostPaidDiffSuppressFunc,
			},
			"period_unit": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     StringInSlice([]string{"Month"}, false),
				DiffSuppressFunc: csNodepoolInstancePostPaidDiffSuppressFunc,
			},
			"platform": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Deprecated:   "Field 'platform' has been deprecated from provider version 1.145.0. Operating system release, using `image_type` instead.",
				ValidateFunc: StringInSlice([]string{"CentOS", "AliyunLinux", "Windows", "WindowsCore"}, false),
			},
			"pre_user_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_pool_options": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"private_pool_options_match_criteria": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"private_pool_options_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"ram_role_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"rds_instances": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"rolling_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_parallelism": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"runtime_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"runtime_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scaling_config": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"instances", "node_count"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: IntBetween(0, 1000),
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"cpu", "gpu", "gpushare", "spot"}, false),
						},
						"eip_bandwidth": {
							Type:          schema.TypeInt,
							Optional:      true,
							ConflictsWith: []string{"internet_charge_type"},
							ValidateFunc:  IntBetween(0, 500),
						},
						"is_bond_eip": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"internet_charge_type"},
						},
						"enable": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"max_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: IntBetween(0, 1000),
						},
						"eip_internet_charge_type": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"internet_charge_type"},
							ValidateFunc:  StringInSlice([]string{"PayByBandwidth", "PayByTraffic"}, false),
						},
					},
				},
			},
			"scaling_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scaling_policy": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     StringInSlice([]string{"release", "recycle"}, false),
				DiffSuppressFunc: csNodepoolScalingPolicyDiffSuppressFunc,
			},
			"security_group_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				ForceNew:   true,
				Deprecated: "Field 'security_group_id' has been deprecated from provider version 1.145.0. The security group ID of the node pool. This field has been replaced by `security_group_ids`, please use the `security_group_ids` field instead.",
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"security_hardening_os": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"soc_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"spot_instance_pools": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 10),
			},
			"spot_instance_remedy": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"spot_price_limit": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"price_limit": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"spot_strategy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"system_disk_bursting_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"system_disk_categories": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"system_disk_category": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"cloud_efficiency", "cloud_ssd", "cloud_essd", "cloud_auto", "cloud_essd_entry", "cloud"}, false),
			},
			"system_disk_encrypt_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"aes-256"}, false),
			},
			"system_disk_encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"system_disk_kms_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"system_disk_performance_level": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: csNodepoolDiskPerformanceLevelDiffSuppressFunc,
			},
			"system_disk_provisioned_iops": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"system_disk_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"system_disk_snapshot_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
			"taints": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"effect": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"tee_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tee_enable": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"unschedulable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"update_nodes": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					_, base64DecodeError := base64.StdEncoding.DecodeString(old)
					if base64DecodeError == nil {
						return new == old
					}
					return new == base64.StdEncoding.EncodeToString([]byte(old))
				},
			},
			"vswitch_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"auto_mode": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'name' has been deprecated since provider version 1.219.0. New field 'node_pool_name' instead.",
			},
			"kms_encrypted_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MaxItems:      100,
				ConflictsWith: []string{"node_count", "scaling_config", "desired_size"},
			},
			"keep_instance_name": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"format_disk": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"node_count": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"instances", "desired_size"},
				Deprecated:    "Field 'node_count' has been deprecated from provider version 1.158.0. New field 'desired_size' instead.",
			},
			"rollout_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_unavailable": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
				Removed: "Field 'rollout_policy' has been removed from provider version 1.184.0. Please use new field 'rolling_policy' instead it to ensure the config takes effect",
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
				Removed:  "Field 'vpc_id' has been removed from provider version 1.218.0.",
			},
			"kms_encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_password").(string) == ""
				},
				Elem: schema.TypeString,
			},
		},
	}
}

func resourceAliCloudAckNodepoolCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	ClusterId := d.Get("cluster_id")
	action := fmt.Sprintf("/clusters/%s/nodepools", ClusterId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("data_disks"); !IsNil(v) {
		if v, ok := d.GetOk("data_disks"); ok {
			localData, err := jsonpath.Get("$", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range localData.([]interface{}) {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["disk_name"] = dataLoopTmp["name"]
				dataLoopMap["auto_snapshot_policy_id"] = dataLoopTmp["auto_snapshot_policy_id"]
				dataLoopMap["file_system"] = dataLoopTmp["file_system"]
				dataLoopMap["kms_key_id"] = dataLoopTmp["kms_key_id"]

				if autoFormatRaw, ok := dataLoopTmp["auto_format"]; ok && autoFormatRaw != "" {
					autoFormat, _ := strconv.ParseBool(autoFormatRaw.(string))
					dataLoopMap["auto_format"] = autoFormat

				}
				dataLoopMap["category"] = dataLoopTmp["category"]
				if dataLoopTmp["size"].(int) > 0 {
					dataLoopMap["size"] = dataLoopTmp["size"]
				}
				dataLoopMap["device"] = dataLoopTmp["device"]
				if dataLoopTmp["provisioned_iops"].(int) > 0 {
					dataLoopMap["provisioned_iops"] = dataLoopTmp["provisioned_iops"]
				}
				dataLoopMap["snapshot_id"] = dataLoopTmp["snapshot_id"]
				dataLoopMap["performance_level"] = dataLoopTmp["performance_level"]
				dataLoopMap["encrypted"] = dataLoopTmp["encrypted"]
				if dataLoopMap["category"] == "cloud_auto" {
					dataLoopMap["bursting_enabled"] = dataLoopTmp["bursting_enabled"]
				}
				dataLoopMap["mount_target"] = dataLoopTmp["mount_target"]
				localMaps = append(localMaps, dataLoopMap)
			}
			objectDataLocalMap["data_disks"] = localMaps
		}

	}

	if v := d.Get("instance_patterns"); !IsNil(v) {
		if v, ok := d.GetOk("instance_patterns"); ok {
			localData1, err := jsonpath.Get("$", v)
			if err != nil {
				localData1 = make([]interface{}, 0)
			}
			localMaps1 := make([]interface{}, 0)
			for _, dataLoop1 := range localData1.([]interface{}) {
				dataLoop1Tmp := make(map[string]interface{})
				if dataLoop1 != nil {
					dataLoop1Tmp = dataLoop1.(map[string]interface{})
				}
				dataLoop1Map := make(map[string]interface{})
				dataLoop1Map["instance_family_level"] = dataLoop1Tmp["instance_family_level"]
				dataLoop1Map["max_cpu_cores"] = dataLoop1Tmp["max_cpu_cores"]
				dataLoop1Map["cpu_architectures"] = dataLoop1Tmp["cpu_architectures"]
				dataLoop1Map["instance_type_families"] = dataLoop1Tmp["instance_type_families"]
				dataLoop1Map["instance_categories"] = dataLoop1Tmp["instance_categories"]
				dataLoop1Map["min_cpu_cores"] = dataLoop1Tmp["min_cpu_cores"]
				dataLoop1Map["max_memory_size"] = dataLoop1Tmp["max_memory_size"]
				dataLoop1Map["min_memory_size"] = dataLoop1Tmp["min_memory_size"]
				dataLoop1Map["excluded_instance_types"] = dataLoop1Tmp["excluded_instance_types"]
				dataLoop1Map["memory"] = dataLoop1Tmp["memory"]
				dataLoop1Map["cores"] = dataLoop1Tmp["cores"]
				localMaps1 = append(localMaps1, dataLoop1Map)
			}
			objectDataLocalMap["instance_patterns"] = localMaps1
		}

	}

	if v, ok := d.GetOk("security_group_id"); ok {
		objectDataLocalMap["security_group_id"] = v
	}

	if v, ok := d.GetOk("soc_enabled"); ok {
		objectDataLocalMap["soc_enabled"] = v
	}

	if v, ok := d.GetOk("node_count"); ok {
		request["count"] = v
	}
	if v, ok := d.GetOk("security_group_ids"); ok {
		objectDataLocalMap["security_group_ids"] = v.(*schema.Set).List()
	}

	if v, ok := d.GetOk("key_name"); ok {
		objectDataLocalMap["key_pair"] = v
	}

	if v, ok := d.GetOk("system_disk_category"); ok {
		objectDataLocalMap["system_disk_category"] = v
	}

	password := d.Get("password").(string)
	if password == "" {
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}
			objectDataLocalMap["login_password"] = decryptResp
		}
	}
	if v, ok := d.GetOk("cis_enabled"); ok {
		objectDataLocalMap["cis_enabled"] = v
	}

	if v, ok := d.GetOk("security_hardening_os"); ok {
		objectDataLocalMap["security_hardening_os"] = v
	}

	if v, ok := d.GetOk("system_disk_kms_key"); ok {
		objectDataLocalMap["system_disk_kms_key_id"] = v
	}

	if v, ok := d.GetOk("system_disk_encrypted"); ok {
		objectDataLocalMap["system_disk_encrypted"] = v
	}

	if v, ok := d.GetOk("instance_charge_type"); ok {
		objectDataLocalMap["instance_charge_type"] = v
	}

	if v, ok := d.GetOk("instance_types"); ok {
		instanceTypes, _ := jsonpath.Get("$", v)
		if instanceTypes != nil && instanceTypes != "" {
			objectDataLocalMap["instance_types"] = instanceTypes
		}
	}

	if v, ok := d.GetOk("platform"); ok {
		objectDataLocalMap["platform"] = v
	}

	if v, ok := d.GetOk("system_disk_snapshot_policy_id"); ok {
		objectDataLocalMap["worker_system_disk_snapshot_policy_id"] = v
	}

	if v, ok := d.GetOk("compensate_with_on_demand"); ok {
		objectDataLocalMap["compensate_with_on_demand"] = v
	}

	if v, ok := d.GetOk("password"); ok {
		objectDataLocalMap["login_password"] = v
	}

	if v, ok := d.GetOk("on_demand_base_capacity"); ok {
		if v != nil && v != "" {
			onDemandBaseCapacity, _ := strconv.ParseInt(v.(string), 10, 64)
			objectDataLocalMap["on_demand_base_capacity"] = onDemandBaseCapacity
		}
	}

	if v, ok := d.GetOk("system_disk_provisioned_iops"); ok {
		objectDataLocalMap["system_disk_provisioned_iops"] = v
	}

	if v := d.Get("tags"); !IsNil(v) {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		objectDataLocalMap["tags"] = tagsMap
	}

	if v, ok := d.GetOk("multi_az_policy"); ok {
		objectDataLocalMap["multi_az_policy"] = v
	}

	if v, ok := d.GetOk("desired_size"); ok {
		if v != nil && v != "" {
			desiredSize, _ := strconv.ParseInt(v.(string), 10, 64)
			objectDataLocalMap["desired_size"] = desiredSize
		}
	}

	if v, ok := d.GetOk("vswitch_ids"); ok {
		vswitchIds, _ := jsonpath.Get("$", v)
		if vswitchIds != nil && vswitchIds != "" {
			objectDataLocalMap["vswitch_ids"] = vswitchIds
		}
	}

	if v := d.Get("private_pool_options"); !IsNil(v) {
		private_pool_options := make(map[string]interface{})
		privatePoolOptionsId, _ := jsonpath.Get("$[0].private_pool_options_id", v)
		if privatePoolOptionsId != nil && privatePoolOptionsId != "" {
			private_pool_options["id"] = privatePoolOptionsId
		}
		privatePoolOptionsMatchCriteria, _ := jsonpath.Get("$[0].private_pool_options_match_criteria", v)
		if privatePoolOptionsMatchCriteria != nil && privatePoolOptionsMatchCriteria != "" {
			private_pool_options["match_criteria"] = privatePoolOptionsMatchCriteria
		}

		objectDataLocalMap["private_pool_options"] = private_pool_options
	}

	if v, ok := d.GetOk("system_disk_size"); ok {
		objectDataLocalMap["system_disk_size"] = v
	}

	if v, ok := d.GetOk("system_disk_performance_level"); ok {
		objectDataLocalMap["system_disk_performance_level"] = v
	}

	if v, ok := d.GetOk("on_demand_percentage_above_base_capacity"); ok {
		if v != nil && v != "" {
			onDemandPercentageAboveBaseCapacity, _ := strconv.ParseInt(v.(string), 10, 64)
			objectDataLocalMap["on_demand_percentage_above_base_capacity"] = onDemandPercentageAboveBaseCapacity
		}
	}

	if v, ok := d.GetOk("internet_charge_type"); ok {
		objectDataLocalMap["internet_charge_type"] = v
	}

	if v, ok := d.GetOk("spot_strategy"); ok {
		objectDataLocalMap["spot_strategy"] = v
	}

	if v, ok := d.GetOk("auto_renew"); ok {
		objectDataLocalMap["auto_renew"] = v
	}

	if v, ok := d.GetOk("image_id"); ok {
		objectDataLocalMap["image_id"] = v
	}

	if v := d.Get("spot_price_limit"); !IsNil(v) {
		if v, ok := d.GetOk("spot_price_limit"); ok {
			localData3, err := jsonpath.Get("$", v)
			if err != nil {
				localData3 = make([]interface{}, 0)
			}
			localMaps3 := make([]interface{}, 0)
			for _, dataLoop3 := range localData3.([]interface{}) {
				dataLoop3Tmp := make(map[string]interface{})
				if dataLoop3 != nil {
					dataLoop3Tmp = dataLoop3.(map[string]interface{})
				}
				dataLoop3Map := make(map[string]interface{})
				dataLoop3Map["price_limit"] = dataLoop3Tmp["price_limit"]
				dataLoop3Map["instance_type"] = dataLoop3Tmp["instance_type"]
				localMaps3 = append(localMaps3, dataLoop3Map)
			}
			objectDataLocalMap["spot_price_limit"] = localMaps3
		}

	}

	if v, ok := d.GetOk("system_disk_categories"); ok {
		systemDiskCategories, _ := jsonpath.Get("$", v)
		if systemDiskCategories != nil && systemDiskCategories != "" {
			objectDataLocalMap["system_disk_categories"] = systemDiskCategories
		}
	}

	if v, ok := d.GetOk("period"); ok {
		objectDataLocalMap["period"] = v
	}

	if v, ok := d.GetOk("image_type"); ok {
		objectDataLocalMap["image_type"] = v
	}

	if v, ok := d.GetOk("period_unit"); ok {
		objectDataLocalMap["period_unit"] = v
	}

	if v, ok := d.GetOk("rds_instances"); ok {
		rdsInstances, _ := jsonpath.Get("$", v)
		if rdsInstances != nil && rdsInstances != "" {
			objectDataLocalMap["rds_instances"] = rdsInstances
		}
	}

	if v, ok := d.GetOk("login_as_non_root"); ok {
		objectDataLocalMap["login_as_non_root"] = v
	}

	if v, ok := d.GetOk("auto_renew_period"); ok {
		objectDataLocalMap["auto_renew_period"] = v
	}

	if v, ok := d.GetOk("deployment_set_id"); ok {
		objectDataLocalMap["deploymentset_id"] = v
	}

	if v, ok := d.GetOk("spot_instance_remedy"); ok {
		objectDataLocalMap["spot_instance_remedy"] = v
	}

	if v, ok := d.GetOk("system_disk_encrypt_algorithm"); ok {
		objectDataLocalMap["system_disk_encrypt_algorithm"] = v
	}

	if v, ok := d.GetOk("system_disk_bursting_enabled"); ok {
		objectDataLocalMap["system_disk_bursting_enabled"] = v
	}

	if v, ok := d.GetOk("scaling_policy"); ok {
		objectDataLocalMap["scaling_policy"] = v
	}

	if v, ok := d.GetOk("spot_instance_pools"); ok {
		objectDataLocalMap["spot_instance_pools"] = v
	}

	if v, ok := d.GetOk("ram_role_name"); ok {
		objectDataLocalMap["ram_role_name"] = v
	}

	if v, ok := d.GetOk("internet_max_bandwidth_out"); ok {
		objectDataLocalMap["internet_max_bandwidth_out"] = v
	}

	request["scaling_group"] = objectDataLocalMap

	objectDataLocalMap1 := make(map[string]interface{})

	if v := d.Get("kubelet_configuration"); !IsNil(v) {
		kubelet_configuration := make(map[string]interface{})
		memoryManagerPolicy1, _ := jsonpath.Get("$[0].memory_manager_policy", v)
		if memoryManagerPolicy1 != nil && memoryManagerPolicy1 != "" {
			kubelet_configuration["memoryManagerPolicy"] = memoryManagerPolicy1
		}
		if v, ok := d.GetOk("kubelet_configuration"); ok {
			localData4, err := jsonpath.Get("$[0].reserved_memory", v)
			if err != nil {
				localData4 = make([]interface{}, 0)
			}
			localMaps4 := make([]interface{}, 0)
			for _, dataLoop4 := range localData4.([]interface{}) {
				dataLoop4Tmp := make(map[string]interface{})
				if dataLoop4 != nil {
					dataLoop4Tmp = dataLoop4.(map[string]interface{})
				}
				dataLoop4Map := make(map[string]interface{})
				dataLoop4Map["numaNode"] = dataLoop4Tmp["numa_node"]
				dataLoop4Map["limits"] = dataLoop4Tmp["limits"]
				localMaps4 = append(localMaps4, dataLoop4Map)
			}
			kubelet_configuration["reservedMemory"] = localMaps4
		}

		allowedUnsafeSysctls1, _ := jsonpath.Get("$[0].allowed_unsafe_sysctls", v)
		if allowedUnsafeSysctls1 != nil && allowedUnsafeSysctls1 != "" {
			kubelet_configuration["allowedUnsafeSysctls"] = allowedUnsafeSysctls1
		}
		evictionSoft1, _ := jsonpath.Get("$[0].eviction_soft", v)
		if evictionSoft1 != nil && evictionSoft1 != "" {
			kubelet_configuration["evictionSoft"] = evictionSoft1
		}
		containerLogMaxWorkers1Raw, _ := jsonpath.Get("$[0].container_log_max_workers", v)
		if containerLogMaxWorkers1Raw != nil && containerLogMaxWorkers1Raw != "" {
			containerLogMaxWorkers1, _ := strconv.ParseInt(containerLogMaxWorkers1Raw.(string), 10, 64)
			kubelet_configuration["containerLogMaxWorkers"] = containerLogMaxWorkers1
		}
		kubeApiBurstRaw, _ := jsonpath.Get("$[0].kube_api_burst", v)
		if kubeApiBurstRaw != nil && kubeApiBurstRaw != "" {
			kubeApiBurst, _ := strconv.ParseInt(kubeApiBurstRaw.(string), 10, 64)
			kubelet_configuration["kubeAPIBurst"] = kubeApiBurst
		}
		registryPullQpsRaw, _ := jsonpath.Get("$[0].registry_pull_qps", v)
		if registryPullQpsRaw != nil && registryPullQpsRaw != "" {
			registryPullQps, _ := strconv.ParseInt(registryPullQpsRaw.(string), 10, 64)
			kubelet_configuration["registryPullQPS"] = registryPullQps
		}
		readOnlyPort1Raw, _ := jsonpath.Get("$[0].read_only_port", v)
		if readOnlyPort1Raw != nil && readOnlyPort1Raw != "" {
			readOnlyPort1, _ := strconv.ParseInt(readOnlyPort1Raw.(string), 10, 64)
			kubelet_configuration["readOnlyPort"] = readOnlyPort1
		}
		cpuCfsQuotaPeriod, _ := jsonpath.Get("$[0].cpu_cfs_quota_period", v)
		if cpuCfsQuotaPeriod != nil && cpuCfsQuotaPeriod != "" {
			kubelet_configuration["cpuCFSQuotaPeriod"] = cpuCfsQuotaPeriod
		}
		evictionSoftGracePeriod1, _ := jsonpath.Get("$[0].eviction_soft_grace_period", v)
		if evictionSoftGracePeriod1 != nil && evictionSoftGracePeriod1 != "" {
			kubelet_configuration["evictionSoftGracePeriod"] = evictionSoftGracePeriod1
		}
		kubeReserved1, _ := jsonpath.Get("$[0].kube_reserved", v)
		if kubeReserved1 != nil && kubeReserved1 != "" {
			kubelet_configuration["kubeReserved"] = kubeReserved1
		}
		cpuCfsQuotaRaw, _ := jsonpath.Get("$[0].cpu_cfs_quota", v)
		if cpuCfsQuotaRaw != nil && cpuCfsQuotaRaw != "" {
			cpuCfsQuota, _ := strconv.ParseBool(cpuCfsQuotaRaw.(string))
			kubelet_configuration["cpuCFSQuota"] = cpuCfsQuota
		}
		topologyManagerPolicy1, _ := jsonpath.Get("$[0].topology_manager_policy", v)
		if topologyManagerPolicy1 != nil && topologyManagerPolicy1 != "" {
			kubelet_configuration["topologyManagerPolicy"] = topologyManagerPolicy1
		}
		featureGates1, _ := jsonpath.Get("$[0].feature_gates", v)
		if featureGates1 != nil && featureGates1 != "" {
			kubelet_configuration["featureGates"] = featureGates1
		}
		podPidsLimit1Raw, _ := jsonpath.Get("$[0].pod_pids_limit", v)
		if podPidsLimit1Raw != nil && podPidsLimit1Raw != "" {
			podPidsLimit1, _ := strconv.ParseInt(podPidsLimit1Raw.(string), 10, 64)
			kubelet_configuration["podPidsLimit"] = podPidsLimit1
		}
		clusterDns, _ := jsonpath.Get("$[0].cluster_dns", v)
		if clusterDns != nil && clusterDns != "" {
			kubelet_configuration["clusterDNS"] = clusterDns
		}
		cpuManagerPolicy1, _ := jsonpath.Get("$[0].cpu_manager_policy", v)
		if cpuManagerPolicy1 != nil && cpuManagerPolicy1 != "" {
			kubelet_configuration["cpuManagerPolicy"] = cpuManagerPolicy1
		}
		systemReserved1, _ := jsonpath.Get("$[0].system_reserved", v)
		if systemReserved1 != nil && systemReserved1 != "" {
			kubelet_configuration["systemReserved"] = systemReserved1
		}
		containerLogMonitorInterval1, _ := jsonpath.Get("$[0].container_log_monitor_interval", v)
		if containerLogMonitorInterval1 != nil && containerLogMonitorInterval1 != "" {
			kubelet_configuration["containerLogMonitorInterval"] = containerLogMonitorInterval1
		}
		imageGcLowThresholdPercentRaw, _ := jsonpath.Get("$[0].image_gc_low_threshold_percent", v)
		if imageGcLowThresholdPercentRaw != nil && imageGcLowThresholdPercentRaw != "" {
			imageGcLowThresholdPercent, _ := strconv.ParseInt(imageGcLowThresholdPercentRaw.(string), 10, 64)
			kubelet_configuration["imageGCLowThresholdPercent"] = imageGcLowThresholdPercent
		}
		eventBurst1Raw, _ := jsonpath.Get("$[0].event_burst", v)
		if eventBurst1Raw != nil && eventBurst1Raw != "" {
			eventBurst1, _ := strconv.ParseInt(eventBurst1Raw.(string), 10, 64)
			kubelet_configuration["eventBurst"] = eventBurst1
		}
		maxPods1Raw, _ := jsonpath.Get("$[0].max_pods", v)
		if maxPods1Raw != nil && maxPods1Raw != "" {
			maxPods1, _ := strconv.ParseInt(maxPods1Raw.(string), 10, 64)
			kubelet_configuration["maxPods"] = maxPods1
		}
		imageGcHighThresholdPercentRaw, _ := jsonpath.Get("$[0].image_gc_high_threshold_percent", v)
		if imageGcHighThresholdPercentRaw != nil && imageGcHighThresholdPercentRaw != "" {
			imageGcHighThresholdPercent, _ := strconv.ParseInt(imageGcHighThresholdPercentRaw.(string), 10, 64)
			kubelet_configuration["imageGCHighThresholdPercent"] = imageGcHighThresholdPercent
		}
		tracing := make(map[string]interface{})
		samplingRatePerMillion1Raw, _ := jsonpath.Get("$[0].tracing[0].sampling_rate_per_million", v)
		if samplingRatePerMillion1Raw != nil && samplingRatePerMillion1Raw != "" {
			samplingRatePerMillion1, _ := strconv.ParseInt(samplingRatePerMillion1Raw.(string), 10, 64)
			tracing["samplingRatePerMillion"] = samplingRatePerMillion1
		}
		endpoint1, _ := jsonpath.Get("$[0].tracing[0].endpoint", v)
		if endpoint1 != nil && endpoint1 != "" {
			tracing["endpoint"] = endpoint1
		}

		kubelet_configuration["tracing"] = tracing
		eventRecordQpsRaw, _ := jsonpath.Get("$[0].event_record_qps", v)
		if eventRecordQpsRaw != nil && eventRecordQpsRaw != "" {
			eventRecordQps, _ := strconv.ParseInt(eventRecordQpsRaw.(string), 10, 64)
			kubelet_configuration["eventRecordQPS"] = eventRecordQps
		}
		containerLogMaxFiles1Raw, _ := jsonpath.Get("$[0].container_log_max_files", v)
		if containerLogMaxFiles1Raw != nil && containerLogMaxFiles1Raw != "" {
			containerLogMaxFiles1, _ := strconv.ParseInt(containerLogMaxFiles1Raw.(string), 10, 64)
			kubelet_configuration["containerLogMaxFiles"] = containerLogMaxFiles1
		}
		evictionHard1, _ := jsonpath.Get("$[0].eviction_hard", v)
		if evictionHard1 != nil && evictionHard1 != "" {
			kubelet_configuration["evictionHard"] = evictionHard1
		}
		registryBurst1Raw, _ := jsonpath.Get("$[0].registry_burst", v)
		if registryBurst1Raw != nil && registryBurst1Raw != "" {
			registryBurst1, _ := strconv.ParseInt(registryBurst1Raw.(string), 10, 64)
			kubelet_configuration["registryBurst"] = registryBurst1
		}
		serializeImagePulls1Raw, _ := jsonpath.Get("$[0].serialize_image_pulls", v)
		if serializeImagePulls1Raw != nil && serializeImagePulls1Raw != "" {
			serializeImagePulls1, _ := strconv.ParseBool(serializeImagePulls1Raw.(string))
			kubelet_configuration["serializeImagePulls"] = serializeImagePulls1
		}
		containerLogMaxSize1, _ := jsonpath.Get("$[0].container_log_max_size", v)
		if containerLogMaxSize1 != nil && containerLogMaxSize1 != "" {
			kubelet_configuration["containerLogMaxSize"] = containerLogMaxSize1
		}
		kubeApiQpsRaw, _ := jsonpath.Get("$[0].kube_api_qps", v)
		if kubeApiQpsRaw != nil && kubeApiQpsRaw != "" {
			kubeApiQps, _ := strconv.ParseInt(kubeApiQpsRaw.(string), 10, 64)
			kubelet_configuration["kubeAPIQPS"] = kubeApiQps
		}

		objectDataLocalMap1["kubelet_configuration"] = kubelet_configuration

		request["node_config"] = objectDataLocalMap1
	}

	objectDataLocalMap2 := make(map[string]interface{})

	if v := d.Get("scaling_config"); !IsNil(v) {
		type1, _ := jsonpath.Get("$[0].type", v)
		if type1 != nil && type1 != "" {
			objectDataLocalMap2["type"] = type1
		}
		eipInternetChargeType, _ := jsonpath.Get("$[0].eip_internet_charge_type", v)
		if eipInternetChargeType != nil && eipInternetChargeType != "" {
			objectDataLocalMap2["eip_internet_charge_type"] = eipInternetChargeType
		}
		eipBandwidth, _ := jsonpath.Get("$[0].eip_bandwidth", v)
		if eipBandwidth != nil && eipBandwidth != "" && eipBandwidth.(int) > 0 {
			objectDataLocalMap2["eip_bandwidth"] = eipBandwidth
		}
		maxSize, _ := jsonpath.Get("$[0].max_size", v)
		if maxSize != nil && maxSize != "" {
			objectDataLocalMap2["max_instances"] = maxSize
		}
		minSize, _ := jsonpath.Get("$[0].min_size", v)
		if minSize != nil && minSize != "" {
			objectDataLocalMap2["min_instances"] = minSize
		}
		isBondEip, _ := jsonpath.Get("$[0].is_bond_eip", v)
		if isBondEip != nil && isBondEip != "" {
			objectDataLocalMap2["is_bond_eip"] = isBondEip
		}
		enable1, _ := jsonpath.Get("$[0].enable", v)
		if enable1 != nil && enable1 != "" {
			objectDataLocalMap2["enable"] = enable1
		}

		request["auto_scaling"] = objectDataLocalMap2
	}

	objectDataLocalMap3 := make(map[string]interface{})

	if v, ok := d.GetOk("runtime_version"); ok {
		objectDataLocalMap3["runtime_version"] = v
	}

	if v, ok := d.GetOk("runtime_name"); ok {
		objectDataLocalMap3["runtime"] = v
	}

	if v := d.Get("labels"); !IsNil(v) {
		if v, ok := d.GetOk("labels"); ok {
			localData5, err := jsonpath.Get("$", v)
			if err != nil {
				localData5 = make([]interface{}, 0)
			}
			localMaps5 := make([]interface{}, 0)
			for _, dataLoop5 := range localData5.([]interface{}) {
				dataLoop5Tmp := make(map[string]interface{})
				if dataLoop5 != nil {
					dataLoop5Tmp = dataLoop5.(map[string]interface{})
				}
				dataLoop5Map := make(map[string]interface{})
				dataLoop5Map["value"] = dataLoop5Tmp["value"]
				dataLoop5Map["key"] = dataLoop5Tmp["key"]
				localMaps5 = append(localMaps5, dataLoop5Map)
			}
			objectDataLocalMap3["labels"] = localMaps5
		}

	}

	if v, ok := d.GetOk("user_data"); ok {
		objectDataLocalMap3["user_data"] = v
		if v := d.Get("user_data").(string); v != "" {
			_, base64DecodeError := base64.StdEncoding.DecodeString(v)
			if base64DecodeError == nil {
				objectDataLocalMap3["user_data"] = tea.String(v)
			} else {
				objectDataLocalMap3["user_data"] = tea.String(base64.StdEncoding.EncodeToString([]byte(v)))
			}
		}
	}

	if v, ok := d.GetOk("unschedulable"); ok {
		objectDataLocalMap3["unschedulable"] = v
	}

	if v, ok := d.GetOk("cpu_policy"); ok {
		objectDataLocalMap3["cpu_policy"] = v
	}

	if v := d.Get("taints"); !IsNil(v) {
		if v, ok := d.GetOk("taints"); ok {
			localData6, err := jsonpath.Get("$", v)
			if err != nil {
				localData6 = make([]interface{}, 0)
			}
			localMaps6 := make([]interface{}, 0)
			for _, dataLoop6 := range localData6.([]interface{}) {
				dataLoop6Tmp := make(map[string]interface{})
				if dataLoop6 != nil {
					dataLoop6Tmp = dataLoop6.(map[string]interface{})
				}
				dataLoop6Map := make(map[string]interface{})
				dataLoop6Map["value"] = dataLoop6Tmp["value"]
				dataLoop6Map["key"] = dataLoop6Tmp["key"]
				dataLoop6Map["effect"] = dataLoop6Tmp["effect"]
				localMaps6 = append(localMaps6, dataLoop6Map)
			}
			objectDataLocalMap3["taints"] = localMaps6
		}

	}

	if v, ok := d.GetOk("install_cloud_monitor"); ok {
		objectDataLocalMap3["cms_enabled"] = v
	}

	if v, ok := d.GetOk("pre_user_data"); ok {
		objectDataLocalMap3["pre_user_data"] = v
	}

	if v, ok := d.GetOk("node_name_mode"); ok {
		objectDataLocalMap3["node_name_mode"] = v
	}

	request["kubernetes_config"] = objectDataLocalMap3

	objectDataLocalMap4 := make(map[string]interface{})

	if v := d.Get("management"); !IsNil(v) {
		upgrade_config := make(map[string]interface{})
		surge1, _ := jsonpath.Get("$[0].surge", v)
		if surge1 != nil && surge1 != "" {
			upgrade_config["surge"] = surge1
		}
		surgePercentage, _ := jsonpath.Get("$[0].surge_percentage", v)
		if surgePercentage != nil && surgePercentage != "" {
			upgrade_config["surge_percentage"] = surgePercentage
		}
		maxUnavailable, _ := jsonpath.Get("$[0].max_unavailable", v)
		if maxUnavailable != nil && maxUnavailable != "" && maxUnavailable.(int) > 0 {
			upgrade_config["max_unavailable"] = maxUnavailable
		}

		objectDataLocalMap4["upgrade_config"] = upgrade_config
		auto_upgrade_policy := make(map[string]interface{})
		autoUpgradeKubelet, _ := jsonpath.Get("$[0].auto_upgrade_policy[0].auto_upgrade_kubelet", v)
		if autoUpgradeKubelet != nil && autoUpgradeKubelet != "" {
			auto_upgrade_policy["auto_upgrade_kubelet"] = autoUpgradeKubelet
		}

		objectDataLocalMap4["auto_upgrade_policy"] = auto_upgrade_policy
		autoUpgrade, _ := jsonpath.Get("$[0].auto_upgrade", v)
		if autoUpgrade != nil && autoUpgrade != "" {
			objectDataLocalMap4["auto_upgrade"] = autoUpgrade
		}
		autoRepair, _ := jsonpath.Get("$[0].auto_repair", v)
		if autoRepair != nil && autoRepair != "" {
			objectDataLocalMap4["auto_repair"] = autoRepair
		}
		enable3, _ := jsonpath.Get("$[0].enable", v)
		if enable3 != nil && enable3 != "" {
			objectDataLocalMap4["enable"] = enable3
		}
		autoVulFix, _ := jsonpath.Get("$[0].auto_vul_fix", v)
		if autoVulFix != nil && autoVulFix != "" {
			objectDataLocalMap4["auto_vul_fix"] = autoVulFix
		}
		auto_vul_fix_policy := make(map[string]interface{})
		vulLevel, _ := jsonpath.Get("$[0].auto_vul_fix_policy[0].vul_level", v)
		if vulLevel != nil && vulLevel != "" {
			auto_vul_fix_policy["vul_level"] = vulLevel
		}
		restartNode, _ := jsonpath.Get("$[0].auto_vul_fix_policy[0].restart_node", v)
		if restartNode != nil && restartNode != "" {
			auto_vul_fix_policy["restart_node"] = restartNode
		}

		objectDataLocalMap4["auto_vul_fix_policy"] = auto_vul_fix_policy
		auto_repair_policy := make(map[string]interface{})
		restartNode1, _ := jsonpath.Get("$[0].auto_repair_policy[0].restart_node", v)
		if restartNode1 != nil && restartNode1 != "" {
			auto_repair_policy["restart_node"] = restartNode1
		}

		objectDataLocalMap4["auto_repair_policy"] = auto_repair_policy

		request["management"] = objectDataLocalMap4
	}

	objectDataLocalMap5 := make(map[string]interface{})

	if v, ok := d.GetOk("resource_group_id"); ok {
		objectDataLocalMap5["resource_group_id"] = v
	}

	objectDataLocalMap5["name"] = d.Get("name")
	if v, ok := d.GetOk("node_pool_name"); ok {
		objectDataLocalMap5["name"] = v
	}

	if v, ok := d.GetOk("type"); ok {
		objectDataLocalMap5["type"] = v
	}

	request["nodepool_info"] = objectDataLocalMap5

	objectDataLocalMap6 := make(map[string]interface{})

	if v := d.Get("tee_config"); !IsNil(v) {
		teeEnable, _ := jsonpath.Get("$[0].tee_enable", v)
		if teeEnable != nil && teeEnable != "" {
			objectDataLocalMap6["tee_enable"] = teeEnable
		}

		request["tee_config"] = objectDataLocalMap6
	}

	objectDataLocalMap7 := make(map[string]interface{})

	if v := d.Get("eflo_node_group"); !IsNil(v) {
		clusterId, _ := jsonpath.Get("$[0].cluster_id", v)
		if clusterId != nil && clusterId != "" {
			objectDataLocalMap7["cluster_id"] = clusterId
		}
		groupId, _ := jsonpath.Get("$[0].group_id", v)
		if groupId != nil && groupId != "" {
			objectDataLocalMap7["group_id"] = groupId
		}

		request["eflo_node_group"] = objectDataLocalMap7
	}

	objectDataLocalMap8 := make(map[string]interface{})

	if v := d.Get("auto_mode"); !IsNil(v) {
		enable5, _ := jsonpath.Get("$[0].enable", v)
		if enable5 != nil && enable5 != "" {
			objectDataLocalMap8["enable"] = enable5
		}

		request["auto_mode"] = objectDataLocalMap8
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("CS", "2015-12-15", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_kubernetes_node_pool", action, AlibabaCloudSdkGoERROR)
	}

	nodepool_idVar, _ := jsonpath.Get("$.nodepool_id", response)
	d.SetId(fmt.Sprintf("%v:%v", ClusterId, nodepool_idVar))

	ackServiceV2 := AckServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ackServiceV2.DescribeAsyncAckNodepoolStateRefreshFunc(d, response, "$.state", []string{"fail", "failed"}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	if v, ok := d.GetOk("instances"); ok && v != nil {
		if err := attachExistingInstance(d, meta, expandStringList(v.([]interface{}))); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_kubernetes_node_pool", action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudAckNodepoolRead(d, meta)
}

func resourceAliCloudAckNodepoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ackServiceV2 := AckServiceV2{client}

	objectRaw, err := ackServiceV2.DescribeAckNodepool(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cs_kubernetes_node_pool DescribeAckNodepool Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	kubernetes_configRawObj, _ := jsonpath.Get("$.kubernetes_config", objectRaw)
	kubernetes_configRaw := make(map[string]interface{})
	if kubernetes_configRawObj != nil {
		kubernetes_configRaw = kubernetes_configRawObj.(map[string]interface{})
	}
	d.Set("cpu_policy", kubernetes_configRaw["cpu_policy"])
	d.Set("install_cloud_monitor", kubernetes_configRaw["cms_enabled"])
	d.Set("node_name_mode", kubernetes_configRaw["node_name_mode"])
	d.Set("pre_user_data", kubernetes_configRaw["pre_user_data"])
	d.Set("runtime_name", kubernetes_configRaw["runtime"])
	d.Set("runtime_version", kubernetes_configRaw["runtime_version"])
	d.Set("unschedulable", kubernetes_configRaw["unschedulable"])
	d.Set("user_data", kubernetes_configRaw["user_data"])

	nodepool_infoRawObj, _ := jsonpath.Get("$.nodepool_info", objectRaw)
	nodepool_infoRaw := make(map[string]interface{})
	if nodepool_infoRawObj != nil {
		nodepool_infoRaw = nodepool_infoRawObj.(map[string]interface{})
	}
	d.Set("node_pool_name", nodepool_infoRaw["name"])
	d.Set("resource_group_id", nodepool_infoRaw["resource_group_id"])
	d.Set("type", nodepool_infoRaw["type"])
	d.Set("node_pool_id", nodepool_infoRaw["nodepool_id"])

	scaling_groupRawObj, _ := jsonpath.Get("$.scaling_group", objectRaw)
	scaling_groupRaw := make(map[string]interface{})
	if scaling_groupRawObj != nil {
		scaling_groupRaw = scaling_groupRawObj.(map[string]interface{})
	}
	d.Set("auto_renew", scaling_groupRaw["auto_renew"])
	d.Set("auto_renew_period", scaling_groupRaw["auto_renew_period"])
	d.Set("cis_enabled", scaling_groupRaw["cis_enabled"])
	d.Set("compensate_with_on_demand", scaling_groupRaw["compensate_with_on_demand"])
	d.Set("deployment_set_id", scaling_groupRaw["deploymentset_id"])
	if v, ok := scaling_groupRaw["desired_size"].(json.Number); ok {
		d.Set("desired_size", v.String())
	}

	d.Set("image_id", scaling_groupRaw["image_id"])
	d.Set("image_type", scaling_groupRaw["image_type"])
	d.Set("instance_charge_type", scaling_groupRaw["instance_charge_type"])
	d.Set("internet_charge_type", scaling_groupRaw["internet_charge_type"])
	d.Set("internet_max_bandwidth_out", scaling_groupRaw["internet_max_bandwidth_out"])
	d.Set("key_name", scaling_groupRaw["key_pair"])
	d.Set("login_as_non_root", scaling_groupRaw["login_as_non_root"])
	d.Set("multi_az_policy", scaling_groupRaw["multi_az_policy"])
	if v, ok := scaling_groupRaw["on_demand_base_capacity"].(json.Number); ok {
		d.Set("on_demand_base_capacity", v.String())
	}

	if v, ok := scaling_groupRaw["on_demand_percentage_above_base_capacity"].(json.Number); ok {
		d.Set("on_demand_percentage_above_base_capacity", v.String())
	}
	if passwd, ok := d.GetOk("password"); ok && passwd.(string) != "" {
		d.Set("password", passwd)
	}
	d.Set("period", scaling_groupRaw["period"])
	d.Set("period_unit", scaling_groupRaw["period_unit"])
	d.Set("platform", scaling_groupRaw["platform"])
	d.Set("ram_role_name", scaling_groupRaw["ram_role_name"])
	d.Set("scaling_group_id", scaling_groupRaw["scaling_group_id"])
	d.Set("scaling_policy", scaling_groupRaw["scaling_policy"])
	d.Set("security_group_id", scaling_groupRaw["security_group_id"])
	d.Set("security_hardening_os", scaling_groupRaw["security_hardening_os"])
	d.Set("soc_enabled", scaling_groupRaw["soc_enabled"])
	d.Set("spot_instance_pools", scaling_groupRaw["spot_instance_pools"])
	d.Set("spot_instance_remedy", scaling_groupRaw["spot_instance_remedy"])
	d.Set("spot_strategy", scaling_groupRaw["spot_strategy"])
	d.Set("system_disk_bursting_enabled", scaling_groupRaw["system_disk_bursting_enabled"])
	d.Set("system_disk_category", scaling_groupRaw["system_disk_category"])
	d.Set("system_disk_encrypt_algorithm", scaling_groupRaw["system_disk_encrypt_algorithm"])
	d.Set("system_disk_encrypted", scaling_groupRaw["system_disk_encrypted"])
	d.Set("system_disk_kms_key", scaling_groupRaw["system_disk_kms_key_id"])
	d.Set("system_disk_performance_level", scaling_groupRaw["system_disk_performance_level"])
	d.Set("system_disk_provisioned_iops", scaling_groupRaw["system_disk_provisioned_iops"])
	d.Set("system_disk_size", scaling_groupRaw["system_disk_size"])
	d.Set("system_disk_snapshot_policy_id", scaling_groupRaw["worker_system_disk_snapshot_policy_id"])
	status1RawObj, _ := jsonpath.Get("$.status", objectRaw)
	status1Raw := make(map[string]interface{})
	if status1RawObj != nil {
		status1Raw = status1RawObj.(map[string]interface{})
	}
	d.Set("node_count", status1Raw["total_nodes"])

	data_disksRaw, _ := jsonpath.Get("$.scaling_group.data_disks", objectRaw)
	dataDisksMaps := make([]map[string]interface{}, 0)
	if data_disksRaw != nil {
		for _, data_disksChildRaw := range data_disksRaw.([]interface{}) {
			dataDisksMap := make(map[string]interface{})
			data_disksChildRaw := data_disksChildRaw.(map[string]interface{})
			if v, ok := data_disksChildRaw["auto_format"].(bool); ok {
				dataDisksMap["auto_format"] = strconv.FormatBool(v)
			}

			dataDisksMap["auto_snapshot_policy_id"] = data_disksChildRaw["auto_snapshot_policy_id"]
			dataDisksMap["bursting_enabled"] = data_disksChildRaw["bursting_enabled"]
			dataDisksMap["category"] = data_disksChildRaw["category"]
			dataDisksMap["device"] = data_disksChildRaw["device"]
			dataDisksMap["encrypted"] = data_disksChildRaw["encrypted"]
			dataDisksMap["file_system"] = data_disksChildRaw["file_system"]
			dataDisksMap["kms_key_id"] = data_disksChildRaw["kms_key_id"]
			dataDisksMap["mount_target"] = data_disksChildRaw["mount_target"]
			dataDisksMap["name"] = data_disksChildRaw["disk_name"]
			dataDisksMap["performance_level"] = data_disksChildRaw["performance_level"]
			dataDisksMap["provisioned_iops"] = data_disksChildRaw["provisioned_iops"]
			dataDisksMap["size"] = data_disksChildRaw["size"]
			dataDisksMap["snapshot_id"] = data_disksChildRaw["snapshot_id"]

			dataDisksMaps = append(dataDisksMaps, dataDisksMap)
		}
	}
	if err := d.Set("data_disks", dataDisksMaps); err != nil {
		return err
	}
	instance_patternsRaw, _ := jsonpath.Get("$.scaling_group.instance_patterns", objectRaw)
	instancePatternsMaps := make([]map[string]interface{}, 0)
	if instance_patternsRaw != nil {
		for _, instance_patternsChildRaw := range instance_patternsRaw.([]interface{}) {
			instancePatternsMap := make(map[string]interface{})
			instance_patternsChildRaw := instance_patternsChildRaw.(map[string]interface{})
			instancePatternsMap["cores"] = instance_patternsChildRaw["cores"]
			instancePatternsMap["instance_family_level"] = instance_patternsChildRaw["instance_family_level"]
			instancePatternsMap["max_cpu_cores"] = instance_patternsChildRaw["max_cpu_cores"]
			instancePatternsMap["max_memory_size"] = instance_patternsChildRaw["max_memory_size"]
			instancePatternsMap["memory"] = instance_patternsChildRaw["memory"]
			instancePatternsMap["min_cpu_cores"] = instance_patternsChildRaw["min_cpu_cores"]
			instancePatternsMap["min_memory_size"] = instance_patternsChildRaw["min_memory_size"]

			cpu_architecturesRaw, _ := jsonpath.Get("$.cpu_architectures", instance_patternsChildRaw)
			instancePatternsMap["cpu_architectures"] = cpu_architecturesRaw
			excluded_instance_typesRaw, _ := jsonpath.Get("$.excluded_instance_types", instance_patternsChildRaw)
			instancePatternsMap["excluded_instance_types"] = excluded_instance_typesRaw
			instance_categoriesRaw, _ := jsonpath.Get("$.instance_categories", instance_patternsChildRaw)
			instancePatternsMap["instance_categories"] = instance_categoriesRaw
			instance_type_familiesRaw, _ := jsonpath.Get("$.instance_type_families", instance_patternsChildRaw)
			instancePatternsMap["instance_type_families"] = instance_type_familiesRaw
			instancePatternsMaps = append(instancePatternsMaps, instancePatternsMap)
		}
	}
	if err := d.Set("instance_patterns", instancePatternsMaps); err != nil {
		return err
	}
	instance_typesRaw, _ := jsonpath.Get("$.scaling_group.instance_types", objectRaw)
	d.Set("instance_types", instance_typesRaw)
	kubeletConfigurationMaps := make([]map[string]interface{}, 0)
	kubeletConfigurationMap := make(map[string]interface{})
	kubelet_configurationRawObj, _ := jsonpath.Get("$.node_config.kubelet_configuration", objectRaw)
	kubelet_configurationRaw := make(map[string]interface{})
	if kubelet_configurationRawObj != nil {
		kubelet_configurationRaw = kubelet_configurationRawObj.(map[string]interface{})
	}
	if len(kubelet_configurationRaw) > 0 {
		if v, ok := kubelet_configurationRaw["containerLogMaxFiles"].(json.Number); ok {
			kubeletConfigurationMap["container_log_max_files"] = v.String()
		}

		kubeletConfigurationMap["container_log_max_size"] = kubelet_configurationRaw["containerLogMaxSize"]
		if v, ok := kubelet_configurationRaw["containerLogMaxWorkers"].(json.Number); ok {
			kubeletConfigurationMap["container_log_max_workers"] = v.String()
		}

		kubeletConfigurationMap["container_log_monitor_interval"] = kubelet_configurationRaw["containerLogMonitorInterval"]
		if v, ok := kubelet_configurationRaw["cpuCFSQuota"].(bool); ok {
			kubeletConfigurationMap["cpu_cfs_quota"] = strconv.FormatBool(v)
		}

		kubeletConfigurationMap["cpu_cfs_quota_period"] = kubelet_configurationRaw["cpuCFSQuotaPeriod"]
		kubeletConfigurationMap["cpu_manager_policy"] = kubelet_configurationRaw["cpuManagerPolicy"]
		if v, ok := kubelet_configurationRaw["eventBurst"].(json.Number); ok {
			kubeletConfigurationMap["event_burst"] = v.String()
		}

		if v, ok := kubelet_configurationRaw["eventRecordQPS"].(json.Number); ok {
			kubeletConfigurationMap["event_record_qps"] = v.String()
		}

		kubeletConfigurationMap["eviction_hard"] = kubelet_configurationRaw["evictionHard"]
		kubeletConfigurationMap["eviction_soft"] = kubelet_configurationRaw["evictionSoft"]
		kubeletConfigurationMap["eviction_soft_grace_period"] = kubelet_configurationRaw["evictionSoftGracePeriod"]
		kubeletConfigurationMap["feature_gates"] = kubelet_configurationRaw["featureGates"]
		if v, ok := kubelet_configurationRaw["imageGCHighThresholdPercent"].(json.Number); ok {
			kubeletConfigurationMap["image_gc_high_threshold_percent"] = v.String()
		}

		if v, ok := kubelet_configurationRaw["imageGCLowThresholdPercent"].(json.Number); ok {
			kubeletConfigurationMap["image_gc_low_threshold_percent"] = v.String()
		}

		if v, ok := kubelet_configurationRaw["kubeAPIBurst"].(json.Number); ok {
			kubeletConfigurationMap["kube_api_burst"] = v.String()
		}

		if v, ok := kubelet_configurationRaw["kubeAPIQPS"].(json.Number); ok {
			kubeletConfigurationMap["kube_api_qps"] = v.String()
		}

		kubeletConfigurationMap["kube_reserved"] = kubelet_configurationRaw["kubeReserved"]
		if v, ok := kubelet_configurationRaw["maxPods"].(json.Number); ok {
			kubeletConfigurationMap["max_pods"] = v.String()
		}

		kubeletConfigurationMap["memory_manager_policy"] = kubelet_configurationRaw["memoryManagerPolicy"]
		if v, ok := kubelet_configurationRaw["podPidsLimit"].(json.Number); ok {
			kubeletConfigurationMap["pod_pids_limit"] = v.String()
		}

		if v, ok := kubelet_configurationRaw["readOnlyPort"].(json.Number); ok {
			kubeletConfigurationMap["read_only_port"] = v.String()
		}

		if v, ok := kubelet_configurationRaw["registryBurst"].(json.Number); ok {
			kubeletConfigurationMap["registry_burst"] = v.String()
		}

		if v, ok := kubelet_configurationRaw["registryPullQPS"].(json.Number); ok {
			kubeletConfigurationMap["registry_pull_qps"] = v.String()
		}

		if v, ok := kubelet_configurationRaw["serializeImagePulls"].(bool); ok {
			kubeletConfigurationMap["serialize_image_pulls"] = strconv.FormatBool(v)
		}

		kubeletConfigurationMap["system_reserved"] = kubelet_configurationRaw["systemReserved"]
		kubeletConfigurationMap["topology_manager_policy"] = kubelet_configurationRaw["topologyManagerPolicy"]

		allowedUnsafeSysctlsRaw, _ := jsonpath.Get("$.node_config.kubelet_configuration.allowedUnsafeSysctls", objectRaw)
		kubeletConfigurationMap["allowed_unsafe_sysctls"] = allowedUnsafeSysctlsRaw
		clusterDNSRaw, _ := jsonpath.Get("$.node_config.kubelet_configuration.clusterDNS", objectRaw)
		kubeletConfigurationMap["cluster_dns"] = clusterDNSRaw
		reservedMemoryRaw, _ := jsonpath.Get("$.node_config.kubelet_configuration.reservedMemory", objectRaw)
		reservedMemoryMaps := make([]map[string]interface{}, 0)
		if reservedMemoryRaw != nil {
			for _, reservedMemoryChildRaw := range reservedMemoryRaw.([]interface{}) {
				reservedMemoryMap := make(map[string]interface{})
				reservedMemoryChildRaw := reservedMemoryChildRaw.(map[string]interface{})
				reservedMemoryMap["limits"] = reservedMemoryChildRaw["limits"]
				reservedMemoryMap["numa_node"] = reservedMemoryChildRaw["numaNode"]

				reservedMemoryMaps = append(reservedMemoryMaps, reservedMemoryMap)
			}
		}
		kubeletConfigurationMap["reserved_memory"] = reservedMemoryMaps
		tracingMaps := make([]map[string]interface{}, 0)
		tracingMap := make(map[string]interface{})
		tracingRawObj, _ := jsonpath.Get("$.node_config.kubelet_configuration.tracing", objectRaw)
		tracingRaw := make(map[string]interface{})
		if tracingRawObj != nil {
			tracingRaw = tracingRawObj.(map[string]interface{})
		}
		if len(tracingRaw) > 0 {
			tracingMap["endpoint"] = tracingRaw["endpoint"]
			if v, ok := tracingRaw["samplingRatePerMillion"].(json.Number); ok {
				tracingMap["sampling_rate_per_million"] = v.String()
			}

			tracingMaps = append(tracingMaps, tracingMap)
		}
		kubeletConfigurationMap["tracing"] = tracingMaps
		kubeletConfigurationMaps = append(kubeletConfigurationMaps, kubeletConfigurationMap)
	}
	if err := d.Set("kubelet_configuration", kubeletConfigurationMaps); err != nil {
		return err
	}
	labelsRaw, _ := jsonpath.Get("$.kubernetes_config.labels", objectRaw)
	labelsMaps := make([]map[string]interface{}, 0)
	if labelsRaw != nil {
		for _, labelsChildRaw := range labelsRaw.([]interface{}) {
			labelsMap := make(map[string]interface{})
			labelsChildRaw := labelsChildRaw.(map[string]interface{})
			labelsMap["key"] = labelsChildRaw["key"]
			labelsMap["value"] = labelsChildRaw["value"]

			labelsMaps = append(labelsMaps, labelsMap)
		}
	}
	if err := d.Set("labels", labelsMaps); err != nil {
		return err
	}
	managementMaps := make([]map[string]interface{}, 0)
	managementMap := make(map[string]interface{})
	managementRaw := make(map[string]interface{})
	if objectRaw["management"] != nil {
		managementRaw = objectRaw["management"].(map[string]interface{})
	}
	if len(managementRaw) > 0 {
		managementMap["auto_repair"] = managementRaw["auto_repair"]
		managementMap["auto_upgrade"] = managementRaw["auto_upgrade"]
		managementMap["auto_vul_fix"] = managementRaw["auto_vul_fix"]
		managementMap["enable"] = managementRaw["enable"]

		upgrade_configRawObj, _ := jsonpath.Get("$.management.upgrade_config", objectRaw)
		upgrade_configRaw := make(map[string]interface{})
		if upgrade_configRawObj != nil {
			upgrade_configRaw = upgrade_configRawObj.(map[string]interface{})
		}
		if len(upgrade_configRaw) > 0 {
			managementMap["max_unavailable"] = upgrade_configRaw["max_unavailable"]
			managementMap["surge"] = upgrade_configRaw["surge"]
			managementMap["surge_percentage"] = upgrade_configRaw["surge_percentage"]
		}
		autoRepairPolicyMaps := make([]map[string]interface{}, 0)
		autoRepairPolicyMap := make(map[string]interface{})
		auto_repair_policyRaw := make(map[string]interface{})
		if managementRaw["auto_repair_policy"] != nil {
			auto_repair_policyRaw = managementRaw["auto_repair_policy"].(map[string]interface{})
		}
		if len(auto_repair_policyRaw) > 0 {
			autoRepairPolicyMap["restart_node"] = auto_repair_policyRaw["restart_node"]

			autoRepairPolicyMaps = append(autoRepairPolicyMaps, autoRepairPolicyMap)
		}
		managementMap["auto_repair_policy"] = autoRepairPolicyMaps
		autoUpgradePolicyMaps := make([]map[string]interface{}, 0)
		autoUpgradePolicyMap := make(map[string]interface{})
		auto_upgrade_policyRaw := make(map[string]interface{})
		if managementRaw["auto_upgrade_policy"] != nil {
			auto_upgrade_policyRaw = managementRaw["auto_upgrade_policy"].(map[string]interface{})
		}
		if len(auto_upgrade_policyRaw) > 0 {
			autoUpgradePolicyMap["auto_upgrade_kubelet"] = auto_upgrade_policyRaw["auto_upgrade_kubelet"]

			autoUpgradePolicyMaps = append(autoUpgradePolicyMaps, autoUpgradePolicyMap)
		}
		managementMap["auto_upgrade_policy"] = autoUpgradePolicyMaps
		autoVulFixPolicyMaps := make([]map[string]interface{}, 0)
		autoVulFixPolicyMap := make(map[string]interface{})
		auto_vul_fix_policyRaw := make(map[string]interface{})
		if managementRaw["auto_vul_fix_policy"] != nil {
			auto_vul_fix_policyRaw = managementRaw["auto_vul_fix_policy"].(map[string]interface{})
		}
		if len(auto_vul_fix_policyRaw) > 0 {
			autoVulFixPolicyMap["restart_node"] = auto_vul_fix_policyRaw["restart_node"]
			autoVulFixPolicyMap["vul_level"] = auto_vul_fix_policyRaw["vul_level"]

			autoVulFixPolicyMaps = append(autoVulFixPolicyMaps, autoVulFixPolicyMap)
		}
		managementMap["auto_vul_fix_policy"] = autoVulFixPolicyMaps
		managementMaps = append(managementMaps, managementMap)
	}
	if err := d.Set("management", managementMaps); err != nil {
		return err
	}
	privatePoolOptionsMaps := make([]map[string]interface{}, 0)
	privatePoolOptionsMap := make(map[string]interface{})
	private_pool_optionsRawObj, _ := jsonpath.Get("$.scaling_group.private_pool_options", objectRaw)
	private_pool_optionsRaw := make(map[string]interface{})
	if private_pool_optionsRawObj != nil {
		private_pool_optionsRaw = private_pool_optionsRawObj.(map[string]interface{})
	}
	if len(private_pool_optionsRaw) > 0 {
		privatePoolOptionsMap["private_pool_options_id"] = private_pool_optionsRaw["id"]
		privatePoolOptionsMap["private_pool_options_match_criteria"] = private_pool_optionsRaw["match_criteria"]

		privatePoolOptionsMaps = append(privatePoolOptionsMaps, privatePoolOptionsMap)
	}
	if err := d.Set("private_pool_options", privatePoolOptionsMaps); err != nil {
		return err
	}
	rds_instancesRaw, _ := jsonpath.Get("$.scaling_group.rds_instances", objectRaw)
	d.Set("rds_instances", rds_instancesRaw)
	scalingConfigMaps := make([]map[string]interface{}, 0)
	scalingConfigMap := make(map[string]interface{})
	auto_scalingRaw := make(map[string]interface{})
	if objectRaw["auto_scaling"] != nil {
		auto_scalingRaw = objectRaw["auto_scaling"].(map[string]interface{})
	}
	if len(auto_scalingRaw) > 0 {
		scalingConfigMap["eip_bandwidth"] = auto_scalingRaw["eip_bandwidth"]
		scalingConfigMap["eip_internet_charge_type"] = auto_scalingRaw["eip_internet_charge_type"]
		scalingConfigMap["enable"] = auto_scalingRaw["enable"]
		scalingConfigMap["is_bond_eip"] = auto_scalingRaw["is_bond_eip"]
		scalingConfigMap["max_size"] = auto_scalingRaw["max_instances"]
		scalingConfigMap["min_size"] = auto_scalingRaw["min_instances"]
		scalingConfigMap["type"] = auto_scalingRaw["type"]

		scalingConfigMaps = append(scalingConfigMaps, scalingConfigMap)
	}
	if err := d.Set("scaling_config", scalingConfigMaps); err != nil {
		return err
	}
	security_group_idsRaw, _ := jsonpath.Get("$.scaling_group.security_group_ids", objectRaw)
	d.Set("security_group_ids", security_group_idsRaw)
	spot_price_limitRaw, _ := jsonpath.Get("$.scaling_group.spot_price_limit", objectRaw)
	spotPriceLimitMaps := make([]map[string]interface{}, 0)
	if spot_price_limitRaw != nil {
		for _, spot_price_limitChildRaw := range spot_price_limitRaw.([]interface{}) {
			spotPriceLimitMap := make(map[string]interface{})
			spot_price_limitChildRaw := spot_price_limitChildRaw.(map[string]interface{})
			spotPriceLimitMap["instance_type"] = spot_price_limitChildRaw["instance_type"]
			spotPriceLimitMap["price_limit"] = spot_price_limitChildRaw["price_limit"]

			spotPriceLimitMaps = append(spotPriceLimitMaps, spotPriceLimitMap)
		}
	}
	if err := d.Set("spot_price_limit", spotPriceLimitMaps); err != nil {
		return err
	}
	system_disk_categoriesRaw, _ := jsonpath.Get("$.scaling_group.system_disk_categories", objectRaw)
	d.Set("system_disk_categories", system_disk_categoriesRaw)
	tagsMaps, _ := jsonpath.Get("$.scaling_group.tags", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))
	taintsRaw, _ := jsonpath.Get("$.kubernetes_config.taints", objectRaw)
	taintsMaps := make([]map[string]interface{}, 0)
	if taintsRaw != nil {
		for _, taintsChildRaw := range taintsRaw.([]interface{}) {
			taintsMap := make(map[string]interface{})
			taintsChildRaw := taintsChildRaw.(map[string]interface{})
			taintsMap["effect"] = taintsChildRaw["effect"]
			taintsMap["key"] = taintsChildRaw["key"]
			taintsMap["value"] = taintsChildRaw["value"]

			taintsMaps = append(taintsMaps, taintsMap)
		}
	}
	if err := d.Set("taints", taintsMaps); err != nil {
		return err
	}
	teeConfigMaps := make([]map[string]interface{}, 0)
	teeConfigMap := make(map[string]interface{})
	tee_configRaw := make(map[string]interface{})
	if objectRaw["tee_config"] != nil {
		tee_configRaw = objectRaw["tee_config"].(map[string]interface{})
	}
	if len(tee_configRaw) > 0 {
		teeConfigMap["tee_enable"] = tee_configRaw["tee_enable"]

		teeConfigMaps = append(teeConfigMaps, teeConfigMap)
	}
	if err := d.Set("tee_config", teeConfigMaps); err != nil {
		return err
	}
	vswitch_idsRaw, _ := jsonpath.Get("$.scaling_group.vswitch_ids", objectRaw)
	d.Set("vswitch_ids", vswitch_idsRaw)
	auto_modeMaps := make([]map[string]interface{}, 0)
	auto_modeMap := make(map[string]interface{})
	auto_modeRaw := make(map[string]interface{})
	if objectRaw["auto_mode"] != nil {
		auto_modeRaw = objectRaw["auto_mode"].(map[string]interface{})
	}
	if len(auto_modeRaw) > 0 {
		auto_modeMap["enable"] = auto_modeRaw["enable"]

		auto_modeMaps = append(auto_modeMaps, auto_modeMap)
	}
	if err := d.Set("auto_mode", auto_modeMaps); err != nil {
		return err
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("cluster_id", parts[0])

	d.Set("name", d.Get("node_pool_name"))
	return nil
}

func resourceAliCloudAckNodepoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	parts := strings.Split(d.Id(), ":")
	ClusterId := parts[0]
	NodepoolId := parts[1]
	action := fmt.Sprintf("/clusters/%s/nodepools/%s", ClusterId, NodepoolId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	objectDataLocalMap := make(map[string]interface{})

	if d.HasChange("data_disks") {
		update = true
		if v := d.Get("data_disks"); v != nil {
			if v, ok := d.GetOk("data_disks"); ok {
				localData, err := jsonpath.Get("$", v)
				if err != nil {
					localData = make([]interface{}, 0)
				}
				localMaps := make([]interface{}, 0)
				for _, dataLoop := range localData.([]interface{}) {
					dataLoopTmp := make(map[string]interface{})
					if dataLoop != nil {
						dataLoopTmp = dataLoop.(map[string]interface{})
					}
					dataLoopMap := make(map[string]interface{})
					dataLoopMap["disk_name"] = dataLoopTmp["name"]
					dataLoopMap["auto_snapshot_policy_id"] = dataLoopTmp["auto_snapshot_policy_id"]
					dataLoopMap["file_system"] = dataLoopTmp["file_system"]
					dataLoopMap["kms_key_id"] = dataLoopTmp["kms_key_id"]

					if autoFormatRaw, ok := dataLoopTmp["auto_format"]; ok && autoFormatRaw != "" {
						autoFormat, _ := strconv.ParseBool(autoFormatRaw.(string))
						dataLoopMap["auto_format"] = autoFormat

					}
					dataLoopMap["category"] = dataLoopTmp["category"]
					if dataLoopTmp["size"].(int) > 0 {
						dataLoopMap["size"] = dataLoopTmp["size"]
					}
					dataLoopMap["device"] = dataLoopTmp["device"]
					if dataLoopTmp["provisioned_iops"].(int) > 0 {
						dataLoopMap["provisioned_iops"] = dataLoopTmp["provisioned_iops"]
					}
					dataLoopMap["snapshot_id"] = dataLoopTmp["snapshot_id"]
					dataLoopMap["performance_level"] = dataLoopTmp["performance_level"]
					dataLoopMap["encrypted"] = dataLoopTmp["encrypted"]
					if dataLoopMap["category"] == "cloud_auto" {
						dataLoopMap["bursting_enabled"] = dataLoopTmp["bursting_enabled"]
					}
					dataLoopMap["mount_target"] = dataLoopTmp["mount_target"]
					localMaps = append(localMaps, dataLoopMap)
				}
				objectDataLocalMap["data_disks"] = localMaps
			}

		}
	}

	if d.HasChange("instance_patterns") {
		update = true
		if v := d.Get("instance_patterns"); v != nil {
			if v, ok := d.GetOk("instance_patterns"); ok {
				localData1, err := jsonpath.Get("$", v)
				if err != nil {
					localData1 = make([]interface{}, 0)
				}
				localMaps1 := make([]interface{}, 0)
				for _, dataLoop1 := range localData1.([]interface{}) {
					dataLoop1Tmp := make(map[string]interface{})
					if dataLoop1 != nil {
						dataLoop1Tmp = dataLoop1.(map[string]interface{})
					}
					dataLoop1Map := make(map[string]interface{})
					dataLoop1Map["instance_family_level"] = dataLoop1Tmp["instance_family_level"]
					dataLoop1Map["max_cpu_cores"] = dataLoop1Tmp["max_cpu_cores"]
					dataLoop1Map["cpu_architectures"] = dataLoop1Tmp["cpu_architectures"]
					dataLoop1Map["instance_type_families"] = dataLoop1Tmp["instance_type_families"]
					dataLoop1Map["instance_categories"] = dataLoop1Tmp["instance_categories"]
					dataLoop1Map["min_cpu_cores"] = dataLoop1Tmp["min_cpu_cores"]
					dataLoop1Map["max_memory_size"] = dataLoop1Tmp["max_memory_size"]
					dataLoop1Map["min_memory_size"] = dataLoop1Tmp["min_memory_size"]
					dataLoop1Map["excluded_instance_types"] = dataLoop1Tmp["excluded_instance_types"]
					dataLoop1Map["memory"] = dataLoop1Tmp["memory"]
					dataLoop1Map["cores"] = dataLoop1Tmp["cores"]
					localMaps1 = append(localMaps1, dataLoop1Map)
				}
				objectDataLocalMap["instance_patterns"] = localMaps1
			}

		}
	}

	if d.HasChange("key_name") {
		update = true
		if v, ok := d.GetOk("key_name"); ok {
			objectDataLocalMap["key_pair"] = v
		}
	}

	if d.HasChange("system_disk_category") {
		update = true
		if v, ok := d.GetOk("system_disk_category"); ok {
			objectDataLocalMap["system_disk_category"] = v
		}
	}

	if d.HasChange("system_disk_kms_key") {
		update = true
		if v, ok := d.GetOk("system_disk_kms_key"); ok {
			objectDataLocalMap["system_disk_kms_key_id"] = v
		}
	}

	if d.HasChange("system_disk_encrypted") {
		update = true
		objectDataLocalMap["system_disk_encrypted"] = d.Get("system_disk_encrypted")
	}

	if d.HasChange("instance_charge_type") {
		update = true
		if v, ok := d.GetOk("instance_charge_type"); ok {
			objectDataLocalMap["instance_charge_type"] = v
		}
	}

	if d.HasChange("instance_types") {
		update = true
		instanceTypes, _ := jsonpath.Get("$", d.Get("instance_types"))
		if instanceTypes != nil && instanceTypes != "" {
			objectDataLocalMap["instance_types"] = instanceTypes
		}
	}

	if d.HasChange("platform") {
		update = true
		if v, ok := d.GetOk("platform"); ok {
			objectDataLocalMap["platform"] = v
		}
	}

	if d.HasChange("system_disk_snapshot_policy_id") {
		update = true
		if v, ok := d.GetOk("system_disk_snapshot_policy_id"); ok {
			objectDataLocalMap["worker_system_disk_snapshot_policy_id"] = v
		}
	}

	if d.HasChange("compensate_with_on_demand") {
		update = true
		if v, ok := d.GetOk("compensate_with_on_demand"); ok {
			objectDataLocalMap["compensate_with_on_demand"] = v
		}
	}

	if d.HasChange("password") {
		update = true
		if v, ok := d.GetOk("password"); ok {
			objectDataLocalMap["login_password"] = v
		}
	}

	if d.HasChange("on_demand_base_capacity") {
		update = true
		if d.Get("on_demand_base_capacity") != nil && d.Get("on_demand_base_capacity") != "" {
			update = true
			onDemandBaseCapacity, _ := strconv.ParseInt(d.Get("on_demand_base_capacity").(string), 10, 64)
			objectDataLocalMap["on_demand_base_capacity"] = onDemandBaseCapacity
		}
	}

	if d.HasChange("system_disk_provisioned_iops") {
		update = true
		if v, ok := d.GetOk("system_disk_provisioned_iops"); ok {
			objectDataLocalMap["system_disk_provisioned_iops"] = v
		}
	}

	if d.HasChange("tags") {
		update = true
		if v := d.Get("tags"); v != nil {
			tagsMap := ConvertTags(v.(map[string]interface{}))
			objectDataLocalMap["tags"] = tagsMap
		}
	}

	if d.HasChange("multi_az_policy") {
		update = true
		if v, ok := d.GetOk("multi_az_policy"); ok {
			objectDataLocalMap["multi_az_policy"] = v
		}
	}

	if d.HasChange("desired_size") {
		update = true
		if d.Get("desired_size") != nil && d.Get("desired_size") != "" {
			update = true
			desiredSize, _ := strconv.ParseInt(d.Get("desired_size").(string), 10, 64)
			objectDataLocalMap["desired_size"] = desiredSize
		}
	}

	if d.HasChange("vswitch_ids") {
		update = true
		vswitchIds, _ := jsonpath.Get("$", d.Get("vswitch_ids"))
		if vswitchIds != nil && vswitchIds != "" {
			objectDataLocalMap["vswitch_ids"] = vswitchIds
		}
	}

	password := d.Get("password").(string)
	if password == "" {
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}
			objectDataLocalMap["login_password"] = decryptResp
		}
	}
	if d.HasChange("private_pool_options") {
		update = true
		if v := d.Get("private_pool_options"); v != nil {
			private_pool_options := make(map[string]interface{})
			privatePoolOptionsId, _ := jsonpath.Get("$[0].private_pool_options_id", v)
			if privatePoolOptionsId != nil && privatePoolOptionsId != "" {
				private_pool_options["id"] = privatePoolOptionsId
			}
			privatePoolOptionsMatchCriteria, _ := jsonpath.Get("$[0].private_pool_options_match_criteria", v)
			if privatePoolOptionsMatchCriteria != nil && privatePoolOptionsMatchCriteria != "" {
				private_pool_options["match_criteria"] = privatePoolOptionsMatchCriteria
			}

			objectDataLocalMap["private_pool_options"] = private_pool_options
		}
	}

	if d.HasChange("system_disk_size") {
		update = true
		if v, ok := d.GetOk("system_disk_size"); ok {
			objectDataLocalMap["system_disk_size"] = v
		}
	}

	if d.HasChange("system_disk_performance_level") {
		update = true
		if v, ok := d.GetOk("system_disk_performance_level"); ok {
			objectDataLocalMap["system_disk_performance_level"] = v
		}
	}

	if d.HasChange("on_demand_percentage_above_base_capacity") {
		update = true
		if d.Get("on_demand_percentage_above_base_capacity") != nil && d.Get("on_demand_percentage_above_base_capacity") != "" {
			update = true
			onDemandPercentageAboveBaseCapacity, _ := strconv.ParseInt(d.Get("on_demand_percentage_above_base_capacity").(string), 10, 64)
			objectDataLocalMap["on_demand_percentage_above_base_capacity"] = onDemandPercentageAboveBaseCapacity
		}
	}

	if d.HasChange("internet_charge_type") {
		update = true
		if v, ok := d.GetOk("internet_charge_type"); ok {
			objectDataLocalMap["internet_charge_type"] = v
		}
	}

	if d.HasChange("spot_strategy") {
		update = true
		if v, ok := d.GetOk("spot_strategy"); ok {
			objectDataLocalMap["spot_strategy"] = v
		}
	}

	if d.HasChange("auto_renew") {
		update = true
		if v, ok := d.GetOk("auto_renew"); ok {
			objectDataLocalMap["auto_renew"] = v
		}
	}

	if d.HasChange("image_id") {
		update = true
		if v, ok := d.GetOk("image_id"); ok {
			objectDataLocalMap["image_id"] = v
		}
	}

	if d.HasChange("spot_price_limit") {
		update = true
		if v := d.Get("spot_price_limit"); v != nil {
			if v, ok := d.GetOk("spot_price_limit"); ok {
				localData3, err := jsonpath.Get("$", v)
				if err != nil {
					localData3 = make([]interface{}, 0)
				}
				localMaps3 := make([]interface{}, 0)
				for _, dataLoop3 := range localData3.([]interface{}) {
					dataLoop3Tmp := make(map[string]interface{})
					if dataLoop3 != nil {
						dataLoop3Tmp = dataLoop3.(map[string]interface{})
					}
					dataLoop3Map := make(map[string]interface{})
					dataLoop3Map["price_limit"] = dataLoop3Tmp["price_limit"]
					dataLoop3Map["instance_type"] = dataLoop3Tmp["instance_type"]
					localMaps3 = append(localMaps3, dataLoop3Map)
				}
				objectDataLocalMap["spot_price_limit"] = localMaps3
			}

		}
	}

	if d.HasChange("system_disk_categories") {
		update = true
		systemDiskCategories, _ := jsonpath.Get("$", d.Get("system_disk_categories"))
		if systemDiskCategories != nil && systemDiskCategories != "" {
			objectDataLocalMap["system_disk_categories"] = systemDiskCategories
		}
	}

	if d.HasChange("period") {
		update = true
		if v, ok := d.GetOk("period"); ok {
			objectDataLocalMap["period"] = v
		}
	}

	if d.HasChange("image_type") {
		update = true
		if v, ok := d.GetOk("image_type"); ok {
			objectDataLocalMap["image_type"] = v
		}
	}

	if d.HasChange("period_unit") {
		update = true
		if v, ok := d.GetOk("period_unit"); ok {
			objectDataLocalMap["period_unit"] = v
		}
	}

	if d.HasChange("rds_instances") {
		update = true
		rdsInstances, _ := jsonpath.Get("$", d.Get("rds_instances"))
		if rdsInstances != nil && rdsInstances != "" {
			objectDataLocalMap["rds_instances"] = rdsInstances
		}
	}

	if d.HasChange("auto_renew_period") {
		update = true
		if v, ok := d.GetOk("auto_renew_period"); ok {
			objectDataLocalMap["auto_renew_period"] = v
		}
	}

	if d.HasChange("spot_instance_remedy") {
		update = true
		if v, ok := d.GetOk("spot_instance_remedy"); ok {
			objectDataLocalMap["spot_instance_remedy"] = v
		}
	}

	if d.HasChange("system_disk_encrypt_algorithm") {
		update = true
		if v, ok := d.GetOk("system_disk_encrypt_algorithm"); ok {
			objectDataLocalMap["system_disk_encrypt_algorithm"] = v
		}
	}

	if d.HasChange("system_disk_bursting_enabled") {
		update = true
		if v, ok := d.GetOk("system_disk_bursting_enabled"); ok {
			objectDataLocalMap["system_disk_bursting_enabled"] = v
		}
	}

	if d.HasChange("scaling_policy") {
		update = true
		if v, ok := d.GetOk("scaling_policy"); ok {
			objectDataLocalMap["scaling_policy"] = v
		}
	}

	if d.HasChange("spot_instance_pools") {
		update = true
		if v, ok := d.GetOk("spot_instance_pools"); ok {
			objectDataLocalMap["spot_instance_pools"] = v
		}
	}

	if d.HasChange("internet_max_bandwidth_out") {
		update = true
		if v, ok := d.GetOk("internet_max_bandwidth_out"); ok {
			objectDataLocalMap["internet_max_bandwidth_out"] = v
		}
	}

	request["scaling_group"] = objectDataLocalMap

	if !d.IsNewResource() && d.HasChange("scaling_config") {
		update = true
		objectDataLocalMap1 := make(map[string]interface{})

		if v := d.Get("scaling_config"); v != nil {
			type1, _ := jsonpath.Get("$[0].type", v)
			if type1 != nil && type1 != "" {
				objectDataLocalMap1["type"] = type1
			}
			eipInternetChargeType, _ := jsonpath.Get("$[0].eip_internet_charge_type", v)
			if eipInternetChargeType != nil && eipInternetChargeType != "" {
				objectDataLocalMap1["eip_internet_charge_type"] = eipInternetChargeType
			}
			eipBandwidth, _ := jsonpath.Get("$[0].eip_bandwidth", v)
			if eipBandwidth != nil && eipBandwidth != "" && eipBandwidth.(int) > 0 {
				objectDataLocalMap1["eip_bandwidth"] = eipBandwidth
			}
			maxSize, _ := jsonpath.Get("$[0].max_size", v)
			if maxSize != nil && maxSize != "" {
				objectDataLocalMap1["max_instances"] = maxSize
			}
			minSize, _ := jsonpath.Get("$[0].min_size", v)
			if minSize != nil && minSize != "" {
				objectDataLocalMap1["min_instances"] = minSize
			}
			isBondEip, _ := jsonpath.Get("$[0].is_bond_eip", v)
			if isBondEip != nil && isBondEip != "" {
				objectDataLocalMap1["is_bond_eip"] = isBondEip
			}
			enable1, _ := jsonpath.Get("$[0].enable", v)
			if enable1 != nil && enable1 != "" {
				objectDataLocalMap1["enable"] = enable1
			}

			request["auto_scaling"] = objectDataLocalMap1
		}
	}

	objectDataLocalMap2 := make(map[string]interface{})

	if d.HasChange("runtime_version") {
		update = true
		if v, ok := d.GetOk("runtime_version"); ok {
			objectDataLocalMap2["runtime_version"] = v
		}
	}

	if d.HasChange("runtime_name") {
		update = true
		if v, ok := d.GetOk("runtime_name"); ok {
			objectDataLocalMap2["runtime"] = v
		}
	}

	if d.HasChange("labels") {
		update = true
		if v := d.Get("labels"); v != nil {
			if v, ok := d.GetOk("labels"); ok {
				localData4, err := jsonpath.Get("$", v)
				if err != nil {
					localData4 = make([]interface{}, 0)
				}
				localMaps4 := make([]interface{}, 0)
				for _, dataLoop4 := range localData4.([]interface{}) {
					dataLoop4Tmp := make(map[string]interface{})
					if dataLoop4 != nil {
						dataLoop4Tmp = dataLoop4.(map[string]interface{})
					}
					dataLoop4Map := make(map[string]interface{})
					dataLoop4Map["value"] = dataLoop4Tmp["value"]
					dataLoop4Map["key"] = dataLoop4Tmp["key"]
					localMaps4 = append(localMaps4, dataLoop4Map)
				}
				objectDataLocalMap2["labels"] = localMaps4
			}

		}
	}

	if d.HasChange("user_data") {
		update = true
		if v, ok := d.GetOk("user_data"); ok {
			objectDataLocalMap2["user_data"] = v
			if v := d.Get("user_data").(string); v != "" {
				_, base64DecodeError := base64.StdEncoding.DecodeString(v)
				if base64DecodeError == nil {
					objectDataLocalMap2["user_data"] = tea.String(v)
				} else {
					objectDataLocalMap2["user_data"] = tea.String(base64.StdEncoding.EncodeToString([]byte(v)))
				}
			}
		}
	}

	if d.HasChange("unschedulable") {
		update = true
		if v, ok := d.GetOk("unschedulable"); ok {
			objectDataLocalMap2["unschedulable"] = v
		}
	}

	if d.HasChange("cpu_policy") {
		update = true
		if v, ok := d.GetOk("cpu_policy"); ok {
			objectDataLocalMap2["cpu_policy"] = v
		}
	}

	if d.HasChange("taints") {
		update = true
		if v := d.Get("taints"); v != nil {
			if v, ok := d.GetOk("taints"); ok {
				localData5, err := jsonpath.Get("$", v)
				if err != nil {
					localData5 = make([]interface{}, 0)
				}
				localMaps5 := make([]interface{}, 0)
				for _, dataLoop5 := range localData5.([]interface{}) {
					dataLoop5Tmp := make(map[string]interface{})
					if dataLoop5 != nil {
						dataLoop5Tmp = dataLoop5.(map[string]interface{})
					}
					dataLoop5Map := make(map[string]interface{})
					dataLoop5Map["value"] = dataLoop5Tmp["value"]
					dataLoop5Map["key"] = dataLoop5Tmp["key"]
					dataLoop5Map["effect"] = dataLoop5Tmp["effect"]
					localMaps5 = append(localMaps5, dataLoop5Map)
				}
				objectDataLocalMap2["taints"] = localMaps5
			}

		}
	}

	if d.HasChange("install_cloud_monitor") {
		update = true
		if v, ok := d.GetOk("install_cloud_monitor"); ok {
			objectDataLocalMap2["cms_enabled"] = v
		}
	}

	if d.HasChange("pre_user_data") {
		update = true
		if v, ok := d.GetOk("pre_user_data"); ok {
			objectDataLocalMap2["pre_user_data"] = v
		}
	}

	request["kubernetes_config"] = objectDataLocalMap2

	if d.HasChange("management") {
		update = true
		objectDataLocalMap3 := make(map[string]interface{})

		if v := d.Get("management"); v != nil {
			upgrade_config := make(map[string]interface{})
			surge1, _ := jsonpath.Get("$[0].surge", v)
			if surge1 != nil && surge1 != "" {
				upgrade_config["surge"] = surge1
			}
			surgePercentage, _ := jsonpath.Get("$[0].surge_percentage", v)
			if surgePercentage != nil && surgePercentage != "" {
				upgrade_config["surge_percentage"] = surgePercentage
			}
			maxUnavailable, _ := jsonpath.Get("$[0].max_unavailable", v)
			if maxUnavailable != nil && maxUnavailable != "" && maxUnavailable.(int) > 0 {
				upgrade_config["max_unavailable"] = maxUnavailable
			}

			objectDataLocalMap3["upgrade_config"] = upgrade_config
			auto_upgrade_policy := make(map[string]interface{})
			autoUpgradeKubelet, _ := jsonpath.Get("$[0].auto_upgrade_policy[0].auto_upgrade_kubelet", v)
			if autoUpgradeKubelet != nil && autoUpgradeKubelet != "" {
				auto_upgrade_policy["auto_upgrade_kubelet"] = autoUpgradeKubelet
			}

			objectDataLocalMap3["auto_upgrade_policy"] = auto_upgrade_policy
			autoUpgrade, _ := jsonpath.Get("$[0].auto_upgrade", v)
			if autoUpgrade != nil && autoUpgrade != "" {
				objectDataLocalMap3["auto_upgrade"] = autoUpgrade
			}
			autoRepair, _ := jsonpath.Get("$[0].auto_repair", v)
			if autoRepair != nil && autoRepair != "" {
				objectDataLocalMap3["auto_repair"] = autoRepair
			}
			enable3, _ := jsonpath.Get("$[0].enable", v)
			if enable3 != nil && enable3 != "" {
				objectDataLocalMap3["enable"] = enable3
			}
			autoVulFix, _ := jsonpath.Get("$[0].auto_vul_fix", v)
			if autoVulFix != nil && autoVulFix != "" {
				objectDataLocalMap3["auto_vul_fix"] = autoVulFix
			}
			auto_vul_fix_policy := make(map[string]interface{})
			vulLevel, _ := jsonpath.Get("$[0].auto_vul_fix_policy[0].vul_level", v)
			if vulLevel != nil && vulLevel != "" {
				auto_vul_fix_policy["vul_level"] = vulLevel
			}
			restartNode, _ := jsonpath.Get("$[0].auto_vul_fix_policy[0].restart_node", v)
			if restartNode != nil && restartNode != "" {
				auto_vul_fix_policy["restart_node"] = restartNode
			}

			objectDataLocalMap3["auto_vul_fix_policy"] = auto_vul_fix_policy
			auto_repair_policy := make(map[string]interface{})
			restartNode1, _ := jsonpath.Get("$[0].auto_repair_policy[0].restart_node", v)
			if restartNode1 != nil && restartNode1 != "" {
				auto_repair_policy["restart_node"] = restartNode1
			}

			objectDataLocalMap3["auto_repair_policy"] = auto_repair_policy

			request["management"] = objectDataLocalMap3
		}
	}

	objectDataLocalMap4 := make(map[string]interface{})

	if d.HasChange("resource_group_id") {
		update = true
		if v, ok := d.GetOk("resource_group_id"); ok {
			objectDataLocalMap4["resource_group_id"] = v
		}
	}

	if d.HasChange("name") {
		update = true
		objectDataLocalMap4["name"] = d.Get("name")
	}

	if d.HasChange("node_pool_name") {
		update = true
		if v, ok := d.GetOk("node_pool_name"); ok {
			objectDataLocalMap4["name"] = v
		}
	}

	request["nodepool_info"] = objectDataLocalMap4

	if v, ok := d.GetOk("update_nodes"); ok {
		request["update_nodes"] = v
	}
	if _, exist := d.GetOk("desired_size"); !exist && d.HasChange("node_count") {
		oldV, newV := d.GetChange("node_count")
		oldValue, ok := oldV.(int)
		if ok != true {
			return WrapErrorf(fmt.Errorf("node_count old value can not be parsed"), "parseError %d", oldValue)
		}
		newValue, ok := newV.(int)
		if ok != true {
			return WrapErrorf(fmt.Errorf("node_count new value can not be parsed"), "parseError %d", newValue)
		}
		if newValue < oldValue {
			if err = removeNodePoolNodes(d, meta, parts, nil, nil); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "RemoveNodePoolNodes", AlibabaCloudSdkGoERROR)
			}
			// The removal of a node is logically independent.
			// The removal of a node should not involve parameter changes.
			return resourceAliCloudAckNodepoolRead(d, meta)
		}
		update = true
		request["count"] = int64(newValue) - int64(oldValue)
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("CS", "2015-12-15", action, query, nil, body, true)
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
		ackServiceV2 := AckServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ackServiceV2.DescribeAsyncAckNodepoolStateRefreshFunc(d, response, "$.state", []string{"fail", "failed"}))
		if jobDetail, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
		}
	}
	update = false
	parts = strings.Split(d.Id(), ":")
	ClusterId = parts[0]
	NodepoolId = parts[1]
	action = fmt.Sprintf("/clusters/%s/nodepools/%s/node_config", ClusterId, NodepoolId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("kubelet_configuration") {
		update = true
		objectDataLocalMap := make(map[string]interface{})

		if v := d.Get("kubelet_configuration"); v != nil {
			memoryManagerPolicy1, _ := jsonpath.Get("$[0].memory_manager_policy", v)
			if memoryManagerPolicy1 != nil && memoryManagerPolicy1 != "" {
				objectDataLocalMap["memoryManagerPolicy"] = memoryManagerPolicy1
			}
			registryPullQpsRaw, _ := jsonpath.Get("$[0].registry_pull_qps", v)
			if registryPullQpsRaw != nil && registryPullQpsRaw != "" {
				registryPullQps, _ := strconv.ParseInt(registryPullQpsRaw.(string), 10, 64)
				objectDataLocalMap["registryPullQPS"] = registryPullQps
			}
			tracing := make(map[string]interface{})
			samplingRatePerMillion1Raw, _ := jsonpath.Get("$[0].tracing[0].sampling_rate_per_million", v)
			if samplingRatePerMillion1Raw != nil && samplingRatePerMillion1Raw != "" {
				samplingRatePerMillion1, _ := strconv.ParseInt(samplingRatePerMillion1Raw.(string), 10, 64)
				tracing["samplingRatePerMillion"] = samplingRatePerMillion1
			}
			endpoint1, _ := jsonpath.Get("$[0].tracing[0].endpoint", v)
			if endpoint1 != nil && endpoint1 != "" {
				tracing["endpoint"] = endpoint1
			}

			objectDataLocalMap["tracing"] = tracing
			eventRecordQpsRaw, _ := jsonpath.Get("$[0].event_record_qps", v)
			if eventRecordQpsRaw != nil && eventRecordQpsRaw != "" {
				eventRecordQps, _ := strconv.ParseInt(eventRecordQpsRaw.(string), 10, 64)
				objectDataLocalMap["eventRecordQPS"] = eventRecordQps
			}
			containerLogMaxFiles1Raw, _ := jsonpath.Get("$[0].container_log_max_files", v)
			if containerLogMaxFiles1Raw != nil && containerLogMaxFiles1Raw != "" {
				containerLogMaxFiles1, _ := strconv.ParseInt(containerLogMaxFiles1Raw.(string), 10, 64)
				objectDataLocalMap["containerLogMaxFiles"] = containerLogMaxFiles1
			}
			if v, ok := d.GetOk("kubelet_configuration"); ok {
				localData, err := jsonpath.Get("$[0].reserved_memory", v)
				if err != nil {
					localData = make([]interface{}, 0)
				}
				localMaps := make([]interface{}, 0)
				for _, dataLoop := range localData.([]interface{}) {
					dataLoopTmp := make(map[string]interface{})
					if dataLoop != nil {
						dataLoopTmp = dataLoop.(map[string]interface{})
					}
					dataLoopMap := make(map[string]interface{})
					dataLoopMap["limits"] = dataLoopTmp["limits"]
					dataLoopMap["numaNode"] = dataLoopTmp["numa_node"]
					localMaps = append(localMaps, dataLoopMap)
				}
				objectDataLocalMap["reservedMemory"] = localMaps
			}

			allowedUnsafeSysctls1, _ := jsonpath.Get("$[0].allowed_unsafe_sysctls", d.Get("kubelet_configuration"))
			if allowedUnsafeSysctls1 != nil && allowedUnsafeSysctls1 != "" {
				objectDataLocalMap["allowedUnsafeSysctls"] = allowedUnsafeSysctls1
			}
			readOnlyPort1Raw, _ := jsonpath.Get("$[0].read_only_port", v)
			if readOnlyPort1Raw != nil && readOnlyPort1Raw != "" {
				readOnlyPort1, _ := strconv.ParseInt(readOnlyPort1Raw.(string), 10, 64)
				objectDataLocalMap["readOnlyPort"] = readOnlyPort1
			}
			cpuCfsQuotaPeriod, _ := jsonpath.Get("$[0].cpu_cfs_quota_period", v)
			if cpuCfsQuotaPeriod != nil && cpuCfsQuotaPeriod != "" {
				objectDataLocalMap["cpuCFSQuotaPeriod"] = cpuCfsQuotaPeriod
			}
			evictionSoftGracePeriod1, _ := jsonpath.Get("$[0].eviction_soft_grace_period", v)
			if evictionSoftGracePeriod1 != nil && evictionSoftGracePeriod1 != "" {
				objectDataLocalMap["evictionSoftGracePeriod"] = evictionSoftGracePeriod1
			}
			evictionSoft1, _ := jsonpath.Get("$[0].eviction_soft", v)
			if evictionSoft1 != nil && evictionSoft1 != "" {
				objectDataLocalMap["evictionSoft"] = evictionSoft1
			}
			containerLogMonitorInterval1, _ := jsonpath.Get("$[0].container_log_monitor_interval", v)
			if containerLogMonitorInterval1 != nil && containerLogMonitorInterval1 != "" {
				objectDataLocalMap["containerLogMonitorInterval"] = containerLogMonitorInterval1
			}
			evictionHard1, _ := jsonpath.Get("$[0].eviction_hard", v)
			if evictionHard1 != nil && evictionHard1 != "" {
				objectDataLocalMap["evictionHard"] = evictionHard1
			}
			kubeReserved1, _ := jsonpath.Get("$[0].kube_reserved", v)
			if kubeReserved1 != nil && kubeReserved1 != "" {
				objectDataLocalMap["kubeReserved"] = kubeReserved1
			}
			registryBurst1Raw, _ := jsonpath.Get("$[0].registry_burst", v)
			if registryBurst1Raw != nil && registryBurst1Raw != "" {
				registryBurst1, _ := strconv.ParseInt(registryBurst1Raw.(string), 10, 64)
				objectDataLocalMap["registryBurst"] = registryBurst1
			}
			imageGcLowThresholdPercentRaw, _ := jsonpath.Get("$[0].image_gc_low_threshold_percent", v)
			if imageGcLowThresholdPercentRaw != nil && imageGcLowThresholdPercentRaw != "" {
				imageGcLowThresholdPercent, _ := strconv.ParseInt(imageGcLowThresholdPercentRaw.(string), 10, 64)
				objectDataLocalMap["imageGCLowThresholdPercent"] = imageGcLowThresholdPercent
			}
			eventBurst1Raw, _ := jsonpath.Get("$[0].event_burst", v)
			if eventBurst1Raw != nil && eventBurst1Raw != "" {
				eventBurst1, _ := strconv.ParseInt(eventBurst1Raw.(string), 10, 64)
				objectDataLocalMap["eventBurst"] = eventBurst1
			}
			serializeImagePulls1Raw, _ := jsonpath.Get("$[0].serialize_image_pulls", v)
			if serializeImagePulls1Raw != nil && serializeImagePulls1Raw != "" {
				serializeImagePulls1, _ := strconv.ParseBool(serializeImagePulls1Raw.(string))
				objectDataLocalMap["serializeImagePulls"] = serializeImagePulls1
			}
			cpuCfsQuotaRaw, _ := jsonpath.Get("$[0].cpu_cfs_quota", v)
			if cpuCfsQuotaRaw != nil && cpuCfsQuotaRaw != "" {
				cpuCfsQuota, _ := strconv.ParseBool(cpuCfsQuotaRaw.(string))
				objectDataLocalMap["cpuCFSQuota"] = cpuCfsQuota
			}
			containerLogMaxSize1, _ := jsonpath.Get("$[0].container_log_max_size", v)
			if containerLogMaxSize1 != nil && containerLogMaxSize1 != "" {
				objectDataLocalMap["containerLogMaxSize"] = containerLogMaxSize1
			}
			topologyManagerPolicy1, _ := jsonpath.Get("$[0].topology_manager_policy", v)
			if topologyManagerPolicy1 != nil && topologyManagerPolicy1 != "" {
				objectDataLocalMap["topologyManagerPolicy"] = topologyManagerPolicy1
			}
			maxPods1Raw, _ := jsonpath.Get("$[0].max_pods", v)
			if maxPods1Raw != nil && maxPods1Raw != "" {
				maxPods1, _ := strconv.ParseInt(maxPods1Raw.(string), 10, 64)
				objectDataLocalMap["maxPods"] = maxPods1
			}
			featureGates1, _ := jsonpath.Get("$[0].feature_gates", v)
			if featureGates1 != nil && featureGates1 != "" {
				objectDataLocalMap["featureGates"] = featureGates1
			}
			podPidsLimit1Raw, _ := jsonpath.Get("$[0].pod_pids_limit", v)
			if podPidsLimit1Raw != nil && podPidsLimit1Raw != "" {
				podPidsLimit1, _ := strconv.ParseInt(podPidsLimit1Raw.(string), 10, 64)
				objectDataLocalMap["podPidsLimit"] = podPidsLimit1
			}
			kubeApiQpsRaw, _ := jsonpath.Get("$[0].kube_api_qps", v)
			if kubeApiQpsRaw != nil && kubeApiQpsRaw != "" {
				kubeApiQps, _ := strconv.ParseInt(kubeApiQpsRaw.(string), 10, 64)
				objectDataLocalMap["kubeAPIQPS"] = kubeApiQps
			}
			imageGcHighThresholdPercentRaw, _ := jsonpath.Get("$[0].image_gc_high_threshold_percent", v)
			if imageGcHighThresholdPercentRaw != nil && imageGcHighThresholdPercentRaw != "" {
				imageGcHighThresholdPercent, _ := strconv.ParseInt(imageGcHighThresholdPercentRaw.(string), 10, 64)
				objectDataLocalMap["imageGCHighThresholdPercent"] = imageGcHighThresholdPercent
			}
			clusterDns, _ := jsonpath.Get("$[0].cluster_dns", d.Get("kubelet_configuration"))
			if clusterDns != nil && clusterDns != "" {
				objectDataLocalMap["clusterDNS"] = clusterDns
			}
			containerLogMaxWorkers1Raw, _ := jsonpath.Get("$[0].container_log_max_workers", v)
			if containerLogMaxWorkers1Raw != nil && containerLogMaxWorkers1Raw != "" {
				containerLogMaxWorkers1, _ := strconv.ParseInt(containerLogMaxWorkers1Raw.(string), 10, 64)
				objectDataLocalMap["containerLogMaxWorkers"] = containerLogMaxWorkers1
			}
			kubeApiBurstRaw, _ := jsonpath.Get("$[0].kube_api_burst", v)
			if kubeApiBurstRaw != nil && kubeApiBurstRaw != "" {
				kubeApiBurst, _ := strconv.ParseInt(kubeApiBurstRaw.(string), 10, 64)
				objectDataLocalMap["kubeAPIBurst"] = kubeApiBurst
			}
			cpuManagerPolicy1, _ := jsonpath.Get("$[0].cpu_manager_policy", v)
			if cpuManagerPolicy1 != nil && cpuManagerPolicy1 != "" {
				objectDataLocalMap["cpuManagerPolicy"] = cpuManagerPolicy1
			}
			systemReserved1, _ := jsonpath.Get("$[0].system_reserved", v)
			if systemReserved1 != nil && systemReserved1 != "" {
				objectDataLocalMap["systemReserved"] = systemReserved1
			}

			request["kubelet_config"] = objectDataLocalMap
		}
	}

	objectDataLocalMap1 := make(map[string]interface{})

	if v := d.Get("rolling_policy"); v != nil {
		maxParallelism, _ := jsonpath.Get("$[0].max_parallelism", v)
		if maxParallelism != nil && maxParallelism != "" {
			objectDataLocalMap1["max_parallelism"] = maxParallelism
		}

		request["rolling_policy"] = objectDataLocalMap1
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("CS", "2015-12-15", action, query, nil, body, true)
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
		ackServiceV2 := AckServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ackServiceV2.DescribeAsyncAckNodepoolStateRefreshFunc(d, response, "$.state", []string{"fail", "failed"}))
		if jobDetail, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
		}
	}

	// attach or remove existing node
	if d.HasChange("instances") {
		rawOldValue, rawNewValue := d.GetChange("instances")
		oldValue, ok := rawOldValue.([]interface{})
		if ok != true {
			return WrapErrorf(fmt.Errorf("instances old value can not be parsed"), "parseError %d", oldValue)
		}
		newValue, ok := rawNewValue.([]interface{})
		if ok != true {
			return WrapErrorf(fmt.Errorf("instances new value can not be parsed"), "parseError %d", oldValue)
		}
		attach, remove := diffInstances(expandStringList(oldValue), expandStringList(newValue))
		if len(attach) > 0 {
			if err = attachExistingInstance(d, meta, attach); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "AttachInstances", AlibabaCloudSdkGoERROR)
			}
		}
		if len(remove) > 0 {
			if err = removeNodePoolNodes(d, meta, parts, oldValue, newValue); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "RemoveNodePoolNodes", AlibabaCloudSdkGoERROR)
			}
		}
	}
	d.Partial(false)
	return resourceAliCloudAckNodepoolRead(d, meta)
}

func resourceAliCloudAckNodepoolDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	NodepoolId := parts[1]
	ClusterId := parts[0]
	action := fmt.Sprintf("/clusters/%s/nodepools/%s", ClusterId, NodepoolId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	query["force"] = tea.String("true")
	if v, ok := d.GetOk("force_delete"); ok {
		query["force"] = StringPointer(strconv.FormatBool(v.(bool)))
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("CS", "2015-12-15", action, query, nil, nil, true)

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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	ackServiceV2 := AckServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ackServiceV2.DescribeAsyncAckNodepoolStateRefreshFunc(d, response, "$.state", []string{"fail", "failed"}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return nil
}

func removeNodePoolNodes(d *schema.ResourceData, meta interface{}, parseId []string, oldNodes []interface{}, newNodes []interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	invoker := NewInvoker()

	var response interface{}
	// list all nodes of the nodepool
	if err := invoker.Run(func() error {
		var err error
		response, err = client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			nodes, _, err := csClient.GetKubernetesClusterNodes(parseId[0], common.Pagination{PageNumber: 1, PageSize: PageSizeLarge}, parseId[1])
			return nodes, err
		})
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetKubernetesClusterNodes", DenverdinoAliyungo)
	}

	ret := response.([]cs.KubernetesNodeType)
	// fetch the NodeName of all nodes
	var allNodeName []string
	for _, value := range ret {
		allNodeName = append(allNodeName, value.NodeName)
	}

	removeNodesName := allNodeName

	// remove automatically created nodes
	if d.HasChange("node_count") {
		o, n := d.GetChange("node_count")
		count := o.(int) - n.(int)
		removeNodesName = allNodeName[:count]
	}

	// remove manually added nodes
	if d.HasChange("instances") {
		var removeInstanceList []string
		var removeInstances []string
		if oldNodes != nil && newNodes != nil {
			_, removeInstances = diffInstances(expandStringList(oldNodes), expandStringList(newNodes))
		}
		for _, v := range ret {
			for _, name := range removeInstances {
				if name == v.InstanceId {
					removeInstanceList = append(removeInstanceList, v.NodeName)
				}
			}
		}
		removeNodesName = removeInstanceList
	}

	if len(removeNodesName) == 0 {
		return nil
	}

	removeNodesArgs := &cs.DeleteKubernetesClusterNodesRequest{
		Nodes:       removeNodesName,
		ReleaseNode: true,
		DrainNode:   false,
	}
	if err := invoker.Run(func() error {
		var err error
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				resp, err := csClient.DeleteKubernetesClusterNodes(parseId[0], removeNodesArgs)
				return resp, err
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

		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteKubernetesClusterNodes", DenverdinoAliyungo)
	}

	stateConf := BuildStateConf([]string{"removing"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, csService.CsKubernetesNodePoolStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	d.SetPartial("node_count")

	return nil
}

const defaultNodePoolType = "ess"

func ConvertCsTags(d *schema.ResourceData) ([]cs.Tag, error) {
	tags := make([]cs.Tag, 0)
	tagsMap, ok := d.Get("tags").(map[string]interface{})
	if ok {
		for key, value := range tagsMap {
			if value != nil {
				if v, ok := value.(string); ok {
					tags = append(tags, cs.Tag{
						Key:   key,
						Value: v,
					})
				}
			}
		}
	}

	return tags, nil
}

func flattenTagsConfig(config []cs.Tag) map[string]string {
	m := make(map[string]string, len(config))
	if len(config) < 0 {
		return m
	}

	for _, tag := range config {
		if tag.Key != DefaultClusterTag {
			m[tag.Key] = tag.Value
		}
	}

	return m
}

func attachExistingInstance(d *schema.ResourceData, meta interface{}, attachInstances []string) error {
	action := "AttachInstancesToNodePool"
	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "InitializeClient", err)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	clusterId := parts[0]
	nodePoolId := parts[1]

	args := &roacs.AttachInstancesToNodePoolRequest{
		FormatDisk:       tea.Bool(false),
		KeepInstanceName: tea.Bool(true),
	}

	if v, ok := d.GetOk("password"); ok {
		args.Password = tea.String(v.(string))
	}

	if v, ok := d.GetOk("format_disk"); ok {
		args.FormatDisk = tea.Bool(v.(bool))
	}

	if v, ok := d.GetOk("keep_instance_name"); ok {
		args.KeepInstanceName = tea.Bool(v.(bool))
	}

	args.Instances = tea.StringSlice(attachInstances)

	var resp *roacs.AttachInstancesToNodePoolResponse
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = client.AttachInstancesToNodePool(tea.String(clusterId), tea.String(nodePoolId), args)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	taskId := tea.StringValue(resp.Body.TaskId)
	if taskId == "" {
		return WrapErrorf(err, ResponseCodeMsg, d.Id(), action, resp)
	}

	csClient := CsClient{client: client}
	stateConf := BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, csClient.DescribeTaskRefreshFunc(d, taskId, []string{"fail", "failed"}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, ResponseCodeMsg, d.Id(), action, jobDetail)
	}

	return nil
}

func diffInstances(old []string, new []string) (attach []string, remove []string) {
	for i, _ := range new {
		found := false
		for j, _ := range old {
			if new[i] == old[j] {
				found = true
			}
		}
		if found == false {
			attach = append(attach, new[i])
		}
	}

	for i, _ := range old {
		found := false
		for j, _ := range new {
			if old[i] == new[j] {
				found = true
			}
		}
		if found == false {
			remove = append(remove, old[i])
		}
	}

	return
}
