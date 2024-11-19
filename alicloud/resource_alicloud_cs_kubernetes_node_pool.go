// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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
	util "github.com/alibabacloud-go/tea-utils/service"
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
							ValidateFunc: StringInSlice([]string{"cloud_efficiency", "cloud_ssd", "cloud_essd", "cloud_auto", "cloud", "cloud_essd_xc0", "cloud_essd_xc1", "all", "ephemeral_ssd", "local_disk"}, false),
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
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"instances", "node_count"},
				ValidateFunc:  IntAtLeast(0),
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
				ValidateFunc: StringInSlice([]string{"AliyunLinux", "AliyunLinux3", "AliyunLinux3Arm64", "AliyunLinuxUEFI", "CentOS", "Windows", "WindowsCore", "AliyunLinux Qboot", "ContainerOS", "AliyunLinuxSecurity", "Ubuntu"}, false),
			},
			"install_cloud_monitor": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      PostPaid,
				ValidateFunc: StringInSlice([]string{"PrePaid", "PostPaid"}, false),
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
						"event_burst": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cpu_manager_policy": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"allowed_unsafe_sysctls": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"eviction_soft_grace_period": {
							Type:     schema.TypeMap,
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
						"system_reserved": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"eviction_soft": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"feature_gates": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeBool},
						},
						"event_record_qps": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"registry_burst": {
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
						"read_only_port": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"registry_pull_qps": {
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
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 1000),
			},
			"on_demand_percentage_above_base_capacity": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 100),
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
				Type:     schema.TypeList,
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
				Type:             schema.TypeList,
				Optional:         true,
				DiffSuppressFunc: csNodepoolSpotInstanceSettingDiffSuppressFunc,
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
				ValidateFunc: StringInSlice([]string{"cloud_efficiency", "cloud_ssd", "cloud_essd", "cloud_auto"}, false),
			},
			"system_disk_encrypt_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
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
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(20, 500),
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
	conn, err := client.NewAckClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	objectDataLocalMap := make(map[string]interface{})

	if v, ok := d.GetOk("resource_group_id"); ok {
		objectDataLocalMap["resource_group_id"] = v
	}

	objectDataLocalMap["name"] = d.Get("name")
	if v, ok := d.GetOk("node_pool_name"); ok {
		objectDataLocalMap["name"] = v
	}

	request["nodepool_info"] = objectDataLocalMap
	objectDataLocalMap1 := make(map[string]interface{})

	if v, ok := d.GetOk("node_count"); ok {
		request["count"] = v
	}
	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIds, _ := jsonpath.Get("$", v)
		if securityGroupIds != nil && securityGroupIds != "" {
			objectDataLocalMap1["security_group_ids"] = securityGroupIds
		}
	}

	if v, ok := d.GetOk("period"); ok {
		objectDataLocalMap1["period"] = v
	}

	if v, ok := d.GetOk("platform"); ok {
		objectDataLocalMap1["platform"] = v
	}

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
				dataLoopMap["category"] = dataLoopTmp["category"]
				if dataLoopMap["category"] == "cloud_auto" {
					dataLoopMap["bursting_enabled"] = dataLoopTmp["bursting_enabled"]
				}
				dataLoopMap["performance_level"] = dataLoopTmp["performance_level"]
				dataLoopMap["auto_snapshot_policy_id"] = dataLoopTmp["auto_snapshot_policy_id"]
				if dataLoopTmp["provisioned_iops"].(int) > 0 {
					dataLoopMap["provisioned_iops"] = dataLoopTmp["provisioned_iops"]
				}
				dataLoopMap["encrypted"] = dataLoopTmp["encrypted"]
				if dataLoopTmp["size"].(int) > 0 {
					dataLoopMap["size"] = dataLoopTmp["size"]
				}
				dataLoopMap["kms_key_id"] = dataLoopTmp["kms_key_id"]
				dataLoopMap["disk_name"] = dataLoopTmp["name"]
				dataLoopMap["device"] = dataLoopTmp["device"]
				dataLoopMap["snapshot_id"] = dataLoopTmp["snapshot_id"]

				if autoFormatRaw, ok := dataLoopTmp["auto_format"]; ok && autoFormatRaw != "" {
					autoFormat, _ := strconv.ParseBool(autoFormatRaw.(string))
					dataLoopMap["auto_format"] = autoFormat

				}
				dataLoopMap["file_system"] = dataLoopTmp["file_system"]
				dataLoopMap["mount_target"] = dataLoopTmp["mount_target"]
				localMaps = append(localMaps, dataLoopMap)
			}
			objectDataLocalMap1["data_disks"] = localMaps
		}

	}

	if v, ok := d.GetOk("deployment_set_id"); ok {
		objectDataLocalMap1["deploymentset_id"] = v
	}

	if v, ok := d.GetOk("compensate_with_on_demand"); ok {
		objectDataLocalMap1["compensate_with_on_demand"] = v
	}

	if v, ok := d.GetOk("auto_renew"); ok {
		objectDataLocalMap1["auto_renew"] = v
	}

	if v, ok := d.GetOk("auto_renew_period"); ok {
		objectDataLocalMap1["auto_renew_period"] = v
	}

	if v, ok := d.GetOk("desired_size"); ok {
		objectDataLocalMap1["desired_size"] = v
	}

	if v, ok := d.GetOk("image_id"); ok {
		objectDataLocalMap1["image_id"] = v
	}

	if v, ok := d.GetOk("image_type"); ok {
		objectDataLocalMap1["image_type"] = v
	}

	if v, ok := d.GetOk("instance_charge_type"); ok {
		objectDataLocalMap1["instance_charge_type"] = v
	}

	if v, ok := d.GetOk("internet_charge_type"); ok {
		objectDataLocalMap1["internet_charge_type"] = v
	}

	if v, ok := d.GetOk("internet_max_bandwidth_out"); ok {
		objectDataLocalMap1["internet_max_bandwidth_out"] = v
	}

	if v, ok := d.GetOk("key_name"); ok {
		objectDataLocalMap1["key_pair"] = v
	}

	if v, ok := d.GetOk("multi_az_policy"); ok {
		objectDataLocalMap1["multi_az_policy"] = v
	}

	if v, ok := d.GetOk("on_demand_base_capacity"); ok {
		objectDataLocalMap1["on_demand_base_capacity"] = v
	}

	if v, ok := d.GetOk("on_demand_percentage_above_base_capacity"); ok {
		objectDataLocalMap1["on_demand_percentage_above_base_capacity"] = v
	}

	if v, ok := d.GetOk("period_unit"); ok {
		objectDataLocalMap1["period_unit"] = v
	}

	if v, ok := d.GetOk("scaling_policy"); ok {
		objectDataLocalMap1["scaling_policy"] = v
	}

	if v, ok := d.GetOk("security_group_id"); ok {
		objectDataLocalMap1["security_group_id"] = v
	}

	if v, ok := d.GetOk("spot_instance_pools"); ok {
		objectDataLocalMap1["spot_instance_pools"] = v
	}

	if v, ok := d.GetOk("spot_instance_remedy"); ok {
		objectDataLocalMap1["spot_instance_remedy"] = v
	}

	if v := d.Get("spot_price_limit"); !IsNil(v) {
		if v, ok := d.GetOk("spot_price_limit"); ok {
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
				dataLoop1Map["instance_type"] = dataLoop1Tmp["instance_type"]
				dataLoop1Map["price_limit"] = dataLoop1Tmp["price_limit"]
				localMaps1 = append(localMaps1, dataLoop1Map)
			}
			objectDataLocalMap1["spot_price_limit"] = localMaps1
		}

	}

	if v, ok := d.GetOk("spot_strategy"); ok {
		objectDataLocalMap1["spot_strategy"] = v
	}

	if v, ok := d.GetOk("system_disk_bursting_enabled"); ok {
		objectDataLocalMap1["system_disk_bursting_enabled"] = v
	}

	if v, ok := d.GetOk("system_disk_category"); ok {
		objectDataLocalMap1["system_disk_category"] = v
	}

	if v, ok := d.GetOk("system_disk_performance_level"); ok {
		objectDataLocalMap1["system_disk_performance_level"] = v
	}

	if v, ok := d.GetOk("vswitch_ids"); ok {
		vswitchIds, _ := jsonpath.Get("$", v)
		if vswitchIds != nil && vswitchIds != "" {
			objectDataLocalMap1["vswitch_ids"] = vswitchIds
		}
	}

	if v := d.Get("tags"); !IsNil(v) {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		objectDataLocalMap1["tags"] = tagsMap
	}

	if v, ok := d.GetOk("system_disk_size"); ok {
		objectDataLocalMap1["system_disk_size"] = v
	}

	if v, ok := d.GetOk("system_disk_provisioned_iops"); ok {
		objectDataLocalMap1["system_disk_provisioned_iops"] = v
	}

	if v, ok := d.GetOk("password"); ok {
		objectDataLocalMap1["login_password"] = v
	}

	if v := d.Get("private_pool_options"); !IsNil(v) {
		private_pool_options := make(map[string]interface{})
		privatePoolOptionsMatchCriteria, _ := jsonpath.Get("$[0].private_pool_options_match_criteria", d.Get("private_pool_options"))
		if privatePoolOptionsMatchCriteria != nil && privatePoolOptionsMatchCriteria != "" {
			private_pool_options["match_criteria"] = privatePoolOptionsMatchCriteria
		}
		privatePoolOptionsId, _ := jsonpath.Get("$[0].private_pool_options_id", d.Get("private_pool_options"))
		if privatePoolOptionsId != nil && privatePoolOptionsId != "" {
			private_pool_options["id"] = privatePoolOptionsId
		}

		objectDataLocalMap1["private_pool_options"] = private_pool_options
	}

	password := d.Get("password").(string)
	if password == "" {
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}
			objectDataLocalMap1["login_password"] = decryptResp
		}
	}
	if v, ok := d.GetOk("cis_enabled"); ok {
		objectDataLocalMap1["cis_enabled"] = v
	}

	if v, ok := d.GetOk("soc_enabled"); ok {
		objectDataLocalMap1["soc_enabled"] = v
	}

	if v, ok := d.GetOk("system_disk_encrypt_algorithm"); ok {
		objectDataLocalMap1["system_disk_encrypt_algorithm"] = v
	}

	if v, ok := d.GetOk("login_as_non_root"); ok {
		objectDataLocalMap1["login_as_non_root"] = v
	}

	if v, ok := d.GetOk("system_disk_encrypted"); ok {
		objectDataLocalMap1["system_disk_encrypted"] = v
	}

	if v, ok := d.GetOk("system_disk_categories"); ok {
		systemDiskCategories, _ := jsonpath.Get("$", v)
		if systemDiskCategories != nil && systemDiskCategories != "" {
			objectDataLocalMap1["system_disk_categories"] = systemDiskCategories
		}
	}

	if v, ok := d.GetOk("instance_types"); ok {
		instanceTypes, _ := jsonpath.Get("$", v)
		if instanceTypes != nil && instanceTypes != "" {
			objectDataLocalMap1["instance_types"] = instanceTypes
		}
	}

	if v, ok := d.GetOk("rds_instances"); ok {
		rdsInstances, _ := jsonpath.Get("$", v)
		if rdsInstances != nil && rdsInstances != "" {
			objectDataLocalMap1["rds_instances"] = rdsInstances
		}
	}

	if v, ok := d.GetOk("system_disk_kms_key"); ok {
		objectDataLocalMap1["system_disk_kms_key_id"] = v
	}

	if v, ok := d.GetOk("system_disk_snapshot_policy_id"); ok {
		objectDataLocalMap1["worker_system_disk_snapshot_policy_id"] = v
	}

	if v, ok := d.GetOk("security_hardening_os"); ok {
		objectDataLocalMap1["security_hardening_os"] = v
	}

	request["scaling_group"] = objectDataLocalMap1
	objectDataLocalMap2 := make(map[string]interface{})

	if v, ok := d.GetOk("cpu_policy"); ok {
		objectDataLocalMap2["cpu_policy"] = v
	}

	if v, ok := d.GetOk("install_cloud_monitor"); ok {
		objectDataLocalMap2["cms_enabled"] = v
	}

	if v, ok := d.GetOk("runtime_version"); ok {
		objectDataLocalMap2["runtime_version"] = v
	}

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

	if v := d.Get("taints"); !IsNil(v) {
		if v, ok := d.GetOk("taints"); ok {
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
				dataLoop3Map["key"] = dataLoop3Tmp["key"]
				dataLoop3Map["effect"] = dataLoop3Tmp["effect"]
				dataLoop3Map["value"] = dataLoop3Tmp["value"]
				localMaps3 = append(localMaps3, dataLoop3Map)
			}
			objectDataLocalMap2["taints"] = localMaps3
		}

	}

	if v, ok := d.GetOk("node_name_mode"); ok {
		objectDataLocalMap2["node_name_mode"] = v
	}

	if v, ok := d.GetOk("unschedulable"); ok {
		objectDataLocalMap2["unschedulable"] = v
	}

	if v, ok := d.GetOk("runtime_name"); ok {
		objectDataLocalMap2["runtime"] = v
	}

	if v := d.Get("labels"); !IsNil(v) {
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
				dataLoop4Map["key"] = dataLoop4Tmp["key"]
				dataLoop4Map["value"] = dataLoop4Tmp["value"]
				localMaps4 = append(localMaps4, dataLoop4Map)
			}
			objectDataLocalMap2["labels"] = localMaps4
		}

	}

	if v, ok := d.GetOk("pre_user_data"); ok {
		objectDataLocalMap2["pre_user_data"] = v
	}

	request["kubernetes_config"] = objectDataLocalMap2
	objectDataLocalMap3 := make(map[string]interface{})

	if v := d.Get("scaling_config"); !IsNil(v) {
		type1, _ := jsonpath.Get("$[0].type", d.Get("scaling_config"))
		if type1 != nil && type1 != "" {
			objectDataLocalMap3["type"] = type1
		}
		maxSize, _ := jsonpath.Get("$[0].max_size", d.Get("scaling_config"))
		if maxSize != nil && maxSize != "" {
			objectDataLocalMap3["max_instances"] = maxSize
		}
		minSize, _ := jsonpath.Get("$[0].min_size", d.Get("scaling_config"))
		if minSize != nil && minSize != "" {
			objectDataLocalMap3["min_instances"] = minSize
		}
		isBondEip, _ := jsonpath.Get("$[0].is_bond_eip", d.Get("scaling_config"))
		if isBondEip != nil && isBondEip != "" {
			objectDataLocalMap3["is_bond_eip"] = isBondEip
		}
		enable1, _ := jsonpath.Get("$[0].enable", d.Get("scaling_config"))
		if enable1 != nil && enable1 != "" {
			objectDataLocalMap3["enable"] = enable1
		}
		eipInternetChargeType, _ := jsonpath.Get("$[0].eip_internet_charge_type", d.Get("scaling_config"))
		if eipInternetChargeType != nil && eipInternetChargeType != "" {
			objectDataLocalMap3["eip_internet_charge_type"] = eipInternetChargeType
		}
		eipBandwidth, _ := jsonpath.Get("$[0].eip_bandwidth", d.Get("scaling_config"))
		if eipBandwidth != nil && eipBandwidth != "" && eipBandwidth.(int) > 0 {
			objectDataLocalMap3["eip_bandwidth"] = eipBandwidth
		}

		request["auto_scaling"] = objectDataLocalMap3
	}

	objectDataLocalMap4 := make(map[string]interface{})

	if v := d.Get("management"); !IsNil(v) {
		enable3, _ := jsonpath.Get("$[0].enable", d.Get("management"))
		if enable3 != nil && enable3 != "" {
			objectDataLocalMap4["enable"] = enable3
		}
		autoRepair, _ := jsonpath.Get("$[0].auto_repair", d.Get("management"))
		if autoRepair != nil && autoRepair != "" {
			objectDataLocalMap4["auto_repair"] = autoRepair
		}
		auto_repair_policy := make(map[string]interface{})
		restartNode, _ := jsonpath.Get("$[0].auto_repair_policy[0].restart_node", d.Get("management"))
		if restartNode != nil && restartNode != "" {
			auto_repair_policy["restart_node"] = restartNode
		}

		objectDataLocalMap4["auto_repair_policy"] = auto_repair_policy
		autoVulFix, _ := jsonpath.Get("$[0].auto_vul_fix", d.Get("management"))
		if autoVulFix != nil && autoVulFix != "" {
			objectDataLocalMap4["auto_vul_fix"] = autoVulFix
		}
		auto_vul_fix_policy := make(map[string]interface{})
		restartNode1, _ := jsonpath.Get("$[0].auto_vul_fix_policy[0].restart_node", d.Get("management"))
		if restartNode1 != nil && restartNode1 != "" {
			auto_vul_fix_policy["restart_node"] = restartNode1
		}
		vulLevel, _ := jsonpath.Get("$[0].auto_vul_fix_policy[0].vul_level", d.Get("management"))
		if vulLevel != nil && vulLevel != "" {
			auto_vul_fix_policy["vul_level"] = vulLevel
		}

		objectDataLocalMap4["auto_vul_fix_policy"] = auto_vul_fix_policy
		autoUpgrade, _ := jsonpath.Get("$[0].auto_upgrade", d.Get("management"))
		if autoUpgrade != nil && autoUpgrade != "" {
			objectDataLocalMap4["auto_upgrade"] = autoUpgrade
		}
		auto_upgrade_policy := make(map[string]interface{})
		autoUpgradeKubelet, _ := jsonpath.Get("$[0].auto_upgrade_policy[0].auto_upgrade_kubelet", d.Get("management"))
		if autoUpgradeKubelet != nil && autoUpgradeKubelet != "" {
			auto_upgrade_policy["auto_upgrade_kubelet"] = autoUpgradeKubelet
		}

		objectDataLocalMap4["auto_upgrade_policy"] = auto_upgrade_policy
		upgrade_config := make(map[string]interface{})
		surge1, _ := jsonpath.Get("$[0].surge", d.Get("management"))
		if surge1 != nil && surge1 != "" {
			upgrade_config["surge"] = surge1
		}
		surgePercentage, _ := jsonpath.Get("$[0].surge_percentage", d.Get("management"))
		if surgePercentage != nil && surgePercentage != "" {
			upgrade_config["surge_percentage"] = surgePercentage
		}
		maxUnavailable, _ := jsonpath.Get("$[0].max_unavailable", d.Get("management"))
		if maxUnavailable != nil && maxUnavailable != "" && maxUnavailable.(int) > 0 {
			upgrade_config["max_unavailable"] = maxUnavailable
		}

		objectDataLocalMap4["upgrade_config"] = upgrade_config

		request["management"] = objectDataLocalMap4
	}

	objectDataLocalMap5 := make(map[string]interface{})

	if v := d.Get("tee_config"); !IsNil(v) {
		teeEnable, _ := jsonpath.Get("$[0].tee_enable", d.Get("tee_config"))
		if teeEnable != nil && teeEnable != "" {
			objectDataLocalMap5["tee_enable"] = teeEnable
		}

		request["tee_config"] = objectDataLocalMap5
	}

	objectDataLocalMap6 := make(map[string]interface{})

	if v := d.Get("kubelet_configuration"); !IsNil(v) {
		kubelet_configuration := make(map[string]interface{})
		registryPullQpsRaw, _ := jsonpath.Get("$[0].registry_pull_qps", d.Get("kubelet_configuration"))
		if registryPullQpsRaw != nil && registryPullQpsRaw != "" {
			registryPullQps, _ := strconv.ParseInt(registryPullQpsRaw.(string), 10, 64)
			kubelet_configuration["registryPullQPS"] = registryPullQps
		}
		registryBurst1Raw, _ := jsonpath.Get("$[0].registry_burst", d.Get("kubelet_configuration"))
		if registryBurst1Raw != nil && registryBurst1Raw != "" {
			registryBurst1, _ := strconv.ParseInt(registryBurst1Raw.(string), 10, 64)
			kubelet_configuration["registryBurst"] = registryBurst1
		}
		eventRecordQpsRaw, _ := jsonpath.Get("$[0].event_record_qps", d.Get("kubelet_configuration"))
		if eventRecordQpsRaw != nil && eventRecordQpsRaw != "" {
			eventRecordQps, _ := strconv.ParseInt(eventRecordQpsRaw.(string), 10, 64)
			kubelet_configuration["eventRecordQPS"] = eventRecordQps
		}
		eventBurst1Raw, _ := jsonpath.Get("$[0].event_burst", d.Get("kubelet_configuration"))
		if eventBurst1Raw != nil && eventBurst1Raw != "" {
			eventBurst1, _ := strconv.ParseInt(eventBurst1Raw.(string), 10, 64)
			kubelet_configuration["eventBurst"] = eventBurst1
		}
		kubeApiQpsRaw, _ := jsonpath.Get("$[0].kube_api_qps", d.Get("kubelet_configuration"))
		if kubeApiQpsRaw != nil && kubeApiQpsRaw != "" {
			kubeApiQps, _ := strconv.ParseInt(kubeApiQpsRaw.(string), 10, 64)
			kubelet_configuration["kubeAPIQPS"] = kubeApiQps
		}
		serializeImagePulls1Raw, _ := jsonpath.Get("$[0].serialize_image_pulls", d.Get("kubelet_configuration"))
		if serializeImagePulls1Raw != nil && serializeImagePulls1Raw != "" {
			serializeImagePulls1, _ := strconv.ParseBool(serializeImagePulls1Raw.(string))
			kubelet_configuration["serializeImagePulls"] = serializeImagePulls1
		}
		cpuManagerPolicy1, _ := jsonpath.Get("$[0].cpu_manager_policy", d.Get("kubelet_configuration"))
		if cpuManagerPolicy1 != nil && cpuManagerPolicy1 != "" {
			kubelet_configuration["cpuManagerPolicy"] = cpuManagerPolicy1
		}
		allowedUnsafeSysctls1, _ := jsonpath.Get("$[0].allowed_unsafe_sysctls", v)
		if allowedUnsafeSysctls1 != nil && allowedUnsafeSysctls1 != "" {
			kubelet_configuration["allowedUnsafeSysctls"] = allowedUnsafeSysctls1
		}
		featureGates1, _ := jsonpath.Get("$[0].feature_gates", d.Get("kubelet_configuration"))
		if featureGates1 != nil && featureGates1 != "" {
			kubelet_configuration["featureGates"] = featureGates1
		}
		containerLogMaxFiles1Raw, _ := jsonpath.Get("$[0].container_log_max_files", d.Get("kubelet_configuration"))
		if containerLogMaxFiles1Raw != nil && containerLogMaxFiles1Raw != "" {
			containerLogMaxFiles1, _ := strconv.ParseInt(containerLogMaxFiles1Raw.(string), 10, 64)
			kubelet_configuration["containerLogMaxFiles"] = containerLogMaxFiles1
		}
		containerLogMaxSize1, _ := jsonpath.Get("$[0].container_log_max_size", d.Get("kubelet_configuration"))
		if containerLogMaxSize1 != nil && containerLogMaxSize1 != "" {
			kubelet_configuration["containerLogMaxSize"] = containerLogMaxSize1
		}
		maxPods1Raw, _ := jsonpath.Get("$[0].max_pods", d.Get("kubelet_configuration"))
		if maxPods1Raw != nil && maxPods1Raw != "" {
			maxPods1, _ := strconv.ParseInt(maxPods1Raw.(string), 10, 64)
			kubelet_configuration["maxPods"] = maxPods1
		}
		readOnlyPort1Raw, _ := jsonpath.Get("$[0].read_only_port", d.Get("kubelet_configuration"))
		if readOnlyPort1Raw != nil && readOnlyPort1Raw != "" {
			readOnlyPort1, _ := strconv.ParseInt(readOnlyPort1Raw.(string), 10, 64)
			kubelet_configuration["readOnlyPort"] = readOnlyPort1
		}
		kubeReserved1, _ := jsonpath.Get("$[0].kube_reserved", d.Get("kubelet_configuration"))
		if kubeReserved1 != nil && kubeReserved1 != "" {
			kubelet_configuration["kubeReserved"] = kubeReserved1
		}
		systemReserved1, _ := jsonpath.Get("$[0].system_reserved", d.Get("kubelet_configuration"))
		if systemReserved1 != nil && systemReserved1 != "" {
			kubelet_configuration["systemReserved"] = systemReserved1
		}
		evictionSoftGracePeriod1, _ := jsonpath.Get("$[0].eviction_soft_grace_period", d.Get("kubelet_configuration"))
		if evictionSoftGracePeriod1 != nil && evictionSoftGracePeriod1 != "" {
			kubelet_configuration["evictionSoftGracePeriod"] = evictionSoftGracePeriod1
		}
		evictionSoft1, _ := jsonpath.Get("$[0].eviction_soft", d.Get("kubelet_configuration"))
		if evictionSoft1 != nil && evictionSoft1 != "" {
			kubelet_configuration["evictionSoft"] = evictionSoft1
		}
		evictionHard1, _ := jsonpath.Get("$[0].eviction_hard", d.Get("kubelet_configuration"))
		if evictionHard1 != nil && evictionHard1 != "" {
			kubelet_configuration["evictionHard"] = evictionHard1
		}
		kubeApiBurstRaw, _ := jsonpath.Get("$[0].kube_api_burst", d.Get("kubelet_configuration"))
		if kubeApiBurstRaw != nil && kubeApiBurstRaw != "" {
			kubeApiBurst, _ := strconv.ParseInt(kubeApiBurstRaw.(string), 10, 64)
			kubelet_configuration["kubeAPIBurst"] = kubeApiBurst
		}

		objectDataLocalMap6["kubelet_configuration"] = kubelet_configuration

		request["node_config"] = objectDataLocalMap6
	}

	body = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2015-12-15"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
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

	nodepool_idVar, _ := jsonpath.Get("$.body.nodepool_id", response)
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

	kubernetes_config1RawObj, _ := jsonpath.Get("$.kubernetes_config", objectRaw)
	kubernetes_config1Raw := make(map[string]interface{})
	if kubernetes_config1RawObj != nil {
		kubernetes_config1Raw = kubernetes_config1RawObj.(map[string]interface{})
	}
	d.Set("cpu_policy", kubernetes_config1Raw["cpu_policy"])
	d.Set("install_cloud_monitor", kubernetes_config1Raw["cms_enabled"])
	d.Set("node_name_mode", kubernetes_config1Raw["node_name_mode"])
	d.Set("pre_user_data", kubernetes_config1Raw["pre_user_data"])
	d.Set("runtime_name", kubernetes_config1Raw["runtime"])
	d.Set("runtime_version", kubernetes_config1Raw["runtime_version"])
	d.Set("unschedulable", kubernetes_config1Raw["unschedulable"])
	d.Set("user_data", kubernetes_config1Raw["user_data"])

	nodepool_info1RawObj, _ := jsonpath.Get("$.nodepool_info", objectRaw)
	nodepool_info1Raw := make(map[string]interface{})
	if nodepool_info1RawObj != nil {
		nodepool_info1Raw = nodepool_info1RawObj.(map[string]interface{})
	}
	d.Set("node_pool_name", nodepool_info1Raw["name"])
	d.Set("resource_group_id", nodepool_info1Raw["resource_group_id"])
	d.Set("node_pool_id", nodepool_info1Raw["nodepool_id"])

	scaling_group1RawObj, _ := jsonpath.Get("$.scaling_group", objectRaw)
	scaling_group1Raw := make(map[string]interface{})
	if scaling_group1RawObj != nil {
		scaling_group1Raw = scaling_group1RawObj.(map[string]interface{})
	}
	d.Set("auto_renew", scaling_group1Raw["auto_renew"])
	d.Set("auto_renew_period", scaling_group1Raw["auto_renew_period"])
	d.Set("cis_enabled", scaling_group1Raw["cis_enabled"])
	d.Set("compensate_with_on_demand", scaling_group1Raw["compensate_with_on_demand"])
	d.Set("deployment_set_id", scaling_group1Raw["deploymentset_id"])
	d.Set("desired_size", scaling_group1Raw["desired_size"])
	d.Set("image_id", scaling_group1Raw["image_id"])
	d.Set("image_type", scaling_group1Raw["image_type"])
	d.Set("instance_charge_type", scaling_group1Raw["instance_charge_type"])
	d.Set("internet_charge_type", scaling_group1Raw["internet_charge_type"])
	d.Set("internet_max_bandwidth_out", scaling_group1Raw["internet_max_bandwidth_out"])
	d.Set("key_name", scaling_group1Raw["key_pair"])
	d.Set("login_as_non_root", scaling_group1Raw["login_as_non_root"])
	d.Set("multi_az_policy", scaling_group1Raw["multi_az_policy"])
	d.Set("on_demand_base_capacity", scaling_group1Raw["on_demand_base_capacity"])
	d.Set("on_demand_percentage_above_base_capacity", scaling_group1Raw["on_demand_percentage_above_base_capacity"])
	if passwd, ok := d.GetOk("password"); ok && passwd.(string) != "" {
		d.Set("password", passwd)
	}
	d.Set("period", scaling_group1Raw["period"])
	d.Set("period_unit", scaling_group1Raw["period_unit"])
	d.Set("platform", scaling_group1Raw["platform"])
	d.Set("scaling_group_id", scaling_group1Raw["scaling_group_id"])
	d.Set("scaling_policy", scaling_group1Raw["scaling_policy"])
	d.Set("security_group_id", scaling_group1Raw["security_group_id"])
	d.Set("security_hardening_os", scaling_group1Raw["security_hardening_os"])
	d.Set("soc_enabled", scaling_group1Raw["soc_enabled"])
	d.Set("spot_instance_pools", scaling_group1Raw["spot_instance_pools"])
	d.Set("spot_instance_remedy", scaling_group1Raw["spot_instance_remedy"])
	d.Set("spot_strategy", scaling_group1Raw["spot_strategy"])
	d.Set("system_disk_bursting_enabled", scaling_group1Raw["system_disk_bursting_enabled"])
	d.Set("system_disk_category", scaling_group1Raw["system_disk_category"])
	d.Set("system_disk_encrypt_algorithm", scaling_group1Raw["system_disk_encrypt_algorithm"])
	d.Set("system_disk_encrypted", scaling_group1Raw["system_disk_encrypted"])
	d.Set("system_disk_kms_key", scaling_group1Raw["system_disk_kms_key_id"])
	d.Set("system_disk_performance_level", scaling_group1Raw["system_disk_performance_level"])
	d.Set("system_disk_provisioned_iops", scaling_group1Raw["system_disk_provisioned_iops"])
	d.Set("system_disk_size", scaling_group1Raw["system_disk_size"])
	d.Set("system_disk_snapshot_policy_id", scaling_group1Raw["worker_system_disk_snapshot_policy_id"])
	status1RawObj, _ := jsonpath.Get("$.status", objectRaw)
	status1Raw := make(map[string]interface{})
	if status1RawObj != nil {
		status1Raw = status1RawObj.(map[string]interface{})
	}
	d.Set("node_count", status1Raw["total_nodes"])

	data_disks1Raw, _ := jsonpath.Get("$.scaling_group.data_disks", objectRaw)
	dataDisksMaps := make([]map[string]interface{}, 0)
	if data_disks1Raw != nil {
		for _, data_disksChild1Raw := range data_disks1Raw.([]interface{}) {
			dataDisksMap := make(map[string]interface{})
			data_disksChild1Raw := data_disksChild1Raw.(map[string]interface{})
			if v, ok := data_disksChild1Raw["auto_format"].(bool); ok {
				dataDisksMap["auto_format"] = strconv.FormatBool(v)
			}

			dataDisksMap["auto_snapshot_policy_id"] = data_disksChild1Raw["auto_snapshot_policy_id"]
			dataDisksMap["bursting_enabled"] = data_disksChild1Raw["bursting_enabled"]
			dataDisksMap["category"] = data_disksChild1Raw["category"]
			dataDisksMap["device"] = data_disksChild1Raw["device"]
			dataDisksMap["encrypted"] = data_disksChild1Raw["encrypted"]
			dataDisksMap["file_system"] = data_disksChild1Raw["file_system"]
			dataDisksMap["kms_key_id"] = data_disksChild1Raw["kms_key_id"]
			dataDisksMap["mount_target"] = data_disksChild1Raw["mount_target"]
			dataDisksMap["name"] = data_disksChild1Raw["disk_name"]
			dataDisksMap["performance_level"] = data_disksChild1Raw["performance_level"]
			dataDisksMap["provisioned_iops"] = data_disksChild1Raw["provisioned_iops"]
			dataDisksMap["size"] = data_disksChild1Raw["size"]
			dataDisksMap["snapshot_id"] = data_disksChild1Raw["snapshot_id"]

			dataDisksMaps = append(dataDisksMaps, dataDisksMap)
		}
	}
	if data_disks1Raw != nil {
		if err := d.Set("data_disks", dataDisksMaps); err != nil {
			return err
		}
	}
	instance_types1Raw, _ := jsonpath.Get("$.scaling_group.instance_types", objectRaw)
	d.Set("instance_types", instance_types1Raw)
	kubeletConfigurationMaps := make([]map[string]interface{}, 0)
	kubeletConfigurationMap := make(map[string]interface{})
	kubelet_configuration1RawObj, _ := jsonpath.Get("$.node_config.kubelet_configuration", objectRaw)
	kubelet_configuration1Raw := make(map[string]interface{})
	if kubelet_configuration1RawObj != nil {
		kubelet_configuration1Raw = kubelet_configuration1RawObj.(map[string]interface{})
	}
	if len(kubelet_configuration1Raw) > 0 {
		if v, ok := kubelet_configuration1Raw["containerLogMaxFiles"].(json.Number); ok {
			kubeletConfigurationMap["container_log_max_files"] = v.String()
		}

		kubeletConfigurationMap["container_log_max_size"] = kubelet_configuration1Raw["containerLogMaxSize"]
		kubeletConfigurationMap["cpu_manager_policy"] = kubelet_configuration1Raw["cpuManagerPolicy"]
		if v, ok := kubelet_configuration1Raw["eventBurst"].(json.Number); ok {
			kubeletConfigurationMap["event_burst"] = v.String()
		}

		if v, ok := kubelet_configuration1Raw["eventRecordQPS"].(json.Number); ok {
			kubeletConfigurationMap["event_record_qps"] = v.String()
		}

		kubeletConfigurationMap["eviction_hard"] = kubelet_configuration1Raw["evictionHard"]
		kubeletConfigurationMap["eviction_soft"] = kubelet_configuration1Raw["evictionSoft"]
		kubeletConfigurationMap["eviction_soft_grace_period"] = kubelet_configuration1Raw["evictionSoftGracePeriod"]
		kubeletConfigurationMap["feature_gates"] = kubelet_configuration1Raw["featureGates"]
		if v, ok := kubelet_configuration1Raw["kubeAPIBurst"].(json.Number); ok {
			kubeletConfigurationMap["kube_api_burst"] = v.String()
		}

		if v, ok := kubelet_configuration1Raw["kubeAPIQPS"].(json.Number); ok {
			kubeletConfigurationMap["kube_api_qps"] = v.String()
		}

		kubeletConfigurationMap["kube_reserved"] = kubelet_configuration1Raw["kubeReserved"]
		if v, ok := kubelet_configuration1Raw["maxPods"].(json.Number); ok {
			kubeletConfigurationMap["max_pods"] = v.String()
		}

		if v, ok := kubelet_configuration1Raw["readOnlyPort"].(json.Number); ok {
			kubeletConfigurationMap["read_only_port"] = v.String()
		}

		if v, ok := kubelet_configuration1Raw["registryBurst"].(json.Number); ok {
			kubeletConfigurationMap["registry_burst"] = v.String()
		}

		if v, ok := kubelet_configuration1Raw["registryPullQPS"].(json.Number); ok {
			kubeletConfigurationMap["registry_pull_qps"] = v.String()
		}

		if v, ok := kubelet_configuration1Raw["serializeImagePulls"].(bool); ok {
			kubeletConfigurationMap["serialize_image_pulls"] = strconv.FormatBool(v)
		}

		kubeletConfigurationMap["system_reserved"] = kubelet_configuration1Raw["systemReserved"]

		allowedUnsafeSysctls1Raw, _ := jsonpath.Get("$.node_config.kubelet_configuration.allowedUnsafeSysctls", objectRaw)
		kubeletConfigurationMap["allowed_unsafe_sysctls"] = allowedUnsafeSysctls1Raw
		kubeletConfigurationMaps = append(kubeletConfigurationMaps, kubeletConfigurationMap)
	}
	if kubelet_configuration1RawObj != nil {
		if err := d.Set("kubelet_configuration", kubeletConfigurationMaps); err != nil {
			return err
		}
	}
	labels1Raw, _ := jsonpath.Get("$.kubernetes_config.labels", objectRaw)
	labelsMaps := make([]map[string]interface{}, 0)
	if labels1Raw != nil {
		for _, labelsChild1Raw := range labels1Raw.([]interface{}) {
			labelsMap := make(map[string]interface{})
			labelsChild1Raw := labelsChild1Raw.(map[string]interface{})
			labelsMap["key"] = labelsChild1Raw["key"]
			labelsMap["value"] = labelsChild1Raw["value"]

			labelsMaps = append(labelsMaps, labelsMap)
		}
	}
	if labels1Raw != nil {
		if err := d.Set("labels", labelsMaps); err != nil {
			return err
		}
	}
	managementMaps := make([]map[string]interface{}, 0)
	managementMap := make(map[string]interface{})
	management1Raw := make(map[string]interface{})
	if objectRaw["management"] != nil {
		management1Raw = objectRaw["management"].(map[string]interface{})
	}
	if len(management1Raw) > 0 {
		managementMap["auto_repair"] = management1Raw["auto_repair"]
		managementMap["auto_upgrade"] = management1Raw["auto_upgrade"]
		managementMap["auto_vul_fix"] = management1Raw["auto_vul_fix"]
		managementMap["enable"] = management1Raw["enable"]

		upgrade_config1RawObj, _ := jsonpath.Get("$.management.upgrade_config", objectRaw)
		upgrade_config1Raw := make(map[string]interface{})
		if upgrade_config1RawObj != nil {
			upgrade_config1Raw = upgrade_config1RawObj.(map[string]interface{})
		}
		if len(upgrade_config1Raw) > 0 {
			managementMap["max_unavailable"] = upgrade_config1Raw["max_unavailable"]
			managementMap["surge"] = upgrade_config1Raw["surge"]
			managementMap["surge_percentage"] = upgrade_config1Raw["surge_percentage"]
		}
		autoRepairPolicyMaps := make([]map[string]interface{}, 0)
		autoRepairPolicyMap := make(map[string]interface{})
		auto_repair_policy1Raw := make(map[string]interface{})
		if management1Raw["auto_repair_policy"] != nil {
			auto_repair_policy1Raw = management1Raw["auto_repair_policy"].(map[string]interface{})
		}
		if len(auto_repair_policy1Raw) > 0 {
			autoRepairPolicyMap["restart_node"] = auto_repair_policy1Raw["restart_node"]

			autoRepairPolicyMaps = append(autoRepairPolicyMaps, autoRepairPolicyMap)
		}
		managementMap["auto_repair_policy"] = autoRepairPolicyMaps
		autoUpgradePolicyMaps := make([]map[string]interface{}, 0)
		autoUpgradePolicyMap := make(map[string]interface{})
		auto_upgrade_policy1Raw := make(map[string]interface{})
		if management1Raw["auto_upgrade_policy"] != nil {
			auto_upgrade_policy1Raw = management1Raw["auto_upgrade_policy"].(map[string]interface{})
		}
		if len(auto_upgrade_policy1Raw) > 0 {
			autoUpgradePolicyMap["auto_upgrade_kubelet"] = auto_upgrade_policy1Raw["auto_upgrade_kubelet"]

			autoUpgradePolicyMaps = append(autoUpgradePolicyMaps, autoUpgradePolicyMap)
		}
		managementMap["auto_upgrade_policy"] = autoUpgradePolicyMaps
		autoVulFixPolicyMaps := make([]map[string]interface{}, 0)
		autoVulFixPolicyMap := make(map[string]interface{})
		auto_vul_fix_policy1Raw := make(map[string]interface{})
		if management1Raw["auto_vul_fix_policy"] != nil {
			auto_vul_fix_policy1Raw = management1Raw["auto_vul_fix_policy"].(map[string]interface{})
		}
		if len(auto_vul_fix_policy1Raw) > 0 {
			autoVulFixPolicyMap["restart_node"] = auto_vul_fix_policy1Raw["restart_node"]
			autoVulFixPolicyMap["vul_level"] = auto_vul_fix_policy1Raw["vul_level"]

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
	private_pool_options1RawObj, _ := jsonpath.Get("$.scaling_group.private_pool_options", objectRaw)
	private_pool_options1Raw := make(map[string]interface{})
	if private_pool_options1RawObj != nil {
		private_pool_options1Raw = private_pool_options1RawObj.(map[string]interface{})
	}
	if len(private_pool_options1Raw) > 0 {
		privatePoolOptionsMap["private_pool_options_id"] = private_pool_options1Raw["id"]
		privatePoolOptionsMap["private_pool_options_match_criteria"] = private_pool_options1Raw["match_criteria"]

		privatePoolOptionsMaps = append(privatePoolOptionsMaps, privatePoolOptionsMap)
	}
	if private_pool_options1RawObj != nil {
		if err := d.Set("private_pool_options", privatePoolOptionsMaps); err != nil {
			return err
		}
	}
	rds_instances1Raw, _ := jsonpath.Get("$.scaling_group.rds_instances", objectRaw)
	d.Set("rds_instances", rds_instances1Raw)
	scalingConfigMaps := make([]map[string]interface{}, 0)
	scalingConfigMap := make(map[string]interface{})
	auto_scaling1Raw := make(map[string]interface{})
	if objectRaw["auto_scaling"] != nil {
		auto_scaling1Raw = objectRaw["auto_scaling"].(map[string]interface{})
	}
	if len(auto_scaling1Raw) > 0 {
		scalingConfigMap["eip_bandwidth"] = auto_scaling1Raw["eip_bandwidth"]
		scalingConfigMap["eip_internet_charge_type"] = auto_scaling1Raw["eip_internet_charge_type"]
		scalingConfigMap["enable"] = auto_scaling1Raw["enable"]
		scalingConfigMap["is_bond_eip"] = auto_scaling1Raw["is_bond_eip"]
		scalingConfigMap["max_size"] = auto_scaling1Raw["max_instances"]
		scalingConfigMap["min_size"] = auto_scaling1Raw["min_instances"]
		scalingConfigMap["type"] = auto_scaling1Raw["type"]

		scalingConfigMaps = append(scalingConfigMaps, scalingConfigMap)
	}
	if objectRaw["auto_scaling"] != nil {
		if err := d.Set("scaling_config", scalingConfigMaps); err != nil {
			return err
		}
	}
	security_group_ids1Raw, _ := jsonpath.Get("$.scaling_group.security_group_ids", objectRaw)
	d.Set("security_group_ids", security_group_ids1Raw)
	spot_price_limit1Raw, _ := jsonpath.Get("$.scaling_group.spot_price_limit", objectRaw)
	spotPriceLimitMaps := make([]map[string]interface{}, 0)
	if spot_price_limit1Raw != nil {
		for _, spot_price_limitChild1Raw := range spot_price_limit1Raw.([]interface{}) {
			spotPriceLimitMap := make(map[string]interface{})
			spot_price_limitChild1Raw := spot_price_limitChild1Raw.(map[string]interface{})
			spotPriceLimitMap["instance_type"] = spot_price_limitChild1Raw["instance_type"]
			spotPriceLimitMap["price_limit"] = spot_price_limitChild1Raw["price_limit"]

			spotPriceLimitMaps = append(spotPriceLimitMaps, spotPriceLimitMap)
		}
	}
	if spot_price_limit1Raw != nil {
		if err := d.Set("spot_price_limit", spotPriceLimitMaps); err != nil {
			return err
		}
	}
	system_disk_categories1Raw, _ := jsonpath.Get("$.scaling_group.system_disk_categories", objectRaw)
	d.Set("system_disk_categories", system_disk_categories1Raw)
	tagsMaps, _ := jsonpath.Get("$.scaling_group.tags", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))
	taints1Raw, _ := jsonpath.Get("$.kubernetes_config.taints", objectRaw)
	taintsMaps := make([]map[string]interface{}, 0)
	if taints1Raw != nil {
		for _, taintsChild1Raw := range taints1Raw.([]interface{}) {
			taintsMap := make(map[string]interface{})
			taintsChild1Raw := taintsChild1Raw.(map[string]interface{})
			taintsMap["effect"] = taintsChild1Raw["effect"]
			taintsMap["key"] = taintsChild1Raw["key"]
			taintsMap["value"] = taintsChild1Raw["value"]

			taintsMaps = append(taintsMaps, taintsMap)
		}
	}
	if taints1Raw != nil {
		if err := d.Set("taints", taintsMaps); err != nil {
			return err
		}
	}
	teeConfigMaps := make([]map[string]interface{}, 0)
	teeConfigMap := make(map[string]interface{})
	tee_config1Raw := make(map[string]interface{})
	if objectRaw["tee_config"] != nil {
		tee_config1Raw = objectRaw["tee_config"].(map[string]interface{})
	}
	if len(tee_config1Raw) > 0 {
		teeConfigMap["tee_enable"] = tee_config1Raw["tee_enable"]

		teeConfigMaps = append(teeConfigMaps, teeConfigMap)
	}
	if objectRaw["tee_config"] != nil {
		if err := d.Set("tee_config", teeConfigMaps); err != nil {
			return err
		}
	}
	vswitch_ids1Raw, _ := jsonpath.Get("$.scaling_group.vswitch_ids", objectRaw)
	d.Set("vswitch_ids", vswitch_ids1Raw)

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

	parts := strings.Split(d.Id(), ":")
	ClusterId := parts[0]
	NodepoolId := parts[1]
	action := fmt.Sprintf("/clusters/%s/nodepools/%s", ClusterId, NodepoolId)
	conn, err := client.NewAckClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	objectDataLocalMap := make(map[string]interface{})

	if d.HasChange("resource_group_id") {
		update = true
		objectDataLocalMap["resource_group_id"] = d.Get("resource_group_id")
	}

	if d.HasChange("name") {
		update = true
		objectDataLocalMap["name"] = d.Get("name")
	}

	if d.HasChange("node_pool_name") {
		update = true
		objectDataLocalMap["name"] = d.Get("node_pool_name")
	}

	request["nodepool_info"] = objectDataLocalMap
	objectDataLocalMap1 := make(map[string]interface{})

	if d.HasChange("period") {
		update = true
		objectDataLocalMap1["period"] = d.Get("period")
	}

	if d.HasChange("platform") {
		update = true
		objectDataLocalMap1["platform"] = d.Get("platform")
	}

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
					dataLoopMap["category"] = dataLoopTmp["category"]
					if dataLoopMap["category"] == "cloud_auto" {
						dataLoopMap["bursting_enabled"] = dataLoopTmp["bursting_enabled"]
					}
					dataLoopMap["performance_level"] = dataLoopTmp["performance_level"]
					dataLoopMap["auto_snapshot_policy_id"] = dataLoopTmp["auto_snapshot_policy_id"]
					if dataLoopTmp["provisioned_iops"].(int) > 0 {
						dataLoopMap["provisioned_iops"] = dataLoopTmp["provisioned_iops"]
					}
					dataLoopMap["encrypted"] = dataLoopTmp["encrypted"]
					if dataLoopTmp["size"].(int) > 0 {
						dataLoopMap["size"] = dataLoopTmp["size"]
					}
					dataLoopMap["kms_key_id"] = dataLoopTmp["kms_key_id"]
					dataLoopMap["device"] = dataLoopTmp["device"]
					dataLoopMap["snapshot_id"] = dataLoopTmp["snapshot_id"]
					dataLoopMap["disk_name"] = dataLoopTmp["name"]

					if autoFormatRaw, ok := dataLoopTmp["auto_format"]; ok && autoFormatRaw != "" {
						autoFormat, _ := strconv.ParseBool(autoFormatRaw.(string))
						dataLoopMap["auto_format"] = autoFormat

					}
					dataLoopMap["file_system"] = dataLoopTmp["file_system"]
					dataLoopMap["mount_target"] = dataLoopTmp["mount_target"]
					localMaps = append(localMaps, dataLoopMap)
				}
				objectDataLocalMap1["data_disks"] = localMaps
			}

		}
	}

	if d.HasChange("compensate_with_on_demand") {
		update = true
		objectDataLocalMap1["compensate_with_on_demand"] = d.Get("compensate_with_on_demand")
	}

	if d.HasChange("auto_renew") {
		update = true
		objectDataLocalMap1["auto_renew"] = d.Get("auto_renew")
	}

	if d.HasChange("auto_renew_period") {
		update = true
		objectDataLocalMap1["auto_renew_period"] = d.Get("auto_renew_period")
	}

	if d.HasChange("desired_size") {
		update = true
		objectDataLocalMap1["desired_size"] = d.Get("desired_size")
	}

	if d.HasChange("image_id") {
		update = true
		objectDataLocalMap1["image_id"] = d.Get("image_id")
	}

	if d.HasChange("instance_charge_type") {
		update = true
		objectDataLocalMap1["instance_charge_type"] = d.Get("instance_charge_type")
	}

	if d.HasChange("internet_charge_type") {
		update = true
		objectDataLocalMap1["internet_charge_type"] = d.Get("internet_charge_type")
	}

	if d.HasChange("internet_max_bandwidth_out") {
		update = true
		objectDataLocalMap1["internet_max_bandwidth_out"] = d.Get("internet_max_bandwidth_out")
	}

	if d.HasChange("key_name") {
		update = true
		objectDataLocalMap1["key_pair"] = d.Get("key_name")
	}

	if d.HasChange("multi_az_policy") {
		update = true
		objectDataLocalMap1["multi_az_policy"] = d.Get("multi_az_policy")
	}

	if d.HasChange("on_demand_base_capacity") {
		update = true
		objectDataLocalMap1["on_demand_base_capacity"] = d.Get("on_demand_base_capacity")
	}

	if d.HasChange("on_demand_percentage_above_base_capacity") {
		update = true
		objectDataLocalMap1["on_demand_percentage_above_base_capacity"] = d.Get("on_demand_percentage_above_base_capacity")
	}

	if d.HasChange("period_unit") {
		update = true
		objectDataLocalMap1["period_unit"] = d.Get("period_unit")
	}

	if d.HasChange("scaling_policy") {
		update = true
		objectDataLocalMap1["scaling_policy"] = d.Get("scaling_policy")
	}

	if d.HasChange("spot_instance_pools") {
		update = true
		objectDataLocalMap1["spot_instance_pools"] = d.Get("spot_instance_pools")
	}

	if d.HasChange("spot_instance_remedy") {
		update = true
		objectDataLocalMap1["spot_instance_remedy"] = d.Get("spot_instance_remedy")
	}

	if d.HasChange("spot_price_limit") {
		update = true
		if v := d.Get("spot_price_limit"); v != nil {
			if v, ok := d.GetOk("spot_price_limit"); ok {
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
					dataLoop1Map["instance_type"] = dataLoop1Tmp["instance_type"]
					dataLoop1Map["price_limit"] = dataLoop1Tmp["price_limit"]
					localMaps1 = append(localMaps1, dataLoop1Map)
				}
				objectDataLocalMap1["spot_price_limit"] = localMaps1
			}

		}
	}

	if d.HasChange("spot_strategy") {
		update = true
		objectDataLocalMap1["spot_strategy"] = d.Get("spot_strategy")
	}

	if d.HasChange("system_disk_category") {
		update = true
		objectDataLocalMap1["system_disk_category"] = d.Get("system_disk_category")
	}

	if d.HasChange("system_disk_performance_level") {
		update = true
		objectDataLocalMap1["system_disk_performance_level"] = d.Get("system_disk_performance_level")
	}

	if d.HasChange("vswitch_ids") {
		update = true
		vswitchIds, _ := jsonpath.Get("$", d.Get("vswitch_ids"))
		if vswitchIds != nil && vswitchIds != "" {
			objectDataLocalMap1["vswitch_ids"] = vswitchIds
		}
	}

	if d.HasChange("tags") {
		update = true
		if v := d.Get("tags"); v != nil {
			tagsMap := ConvertTags(v.(map[string]interface{}))
			objectDataLocalMap1["tags"] = tagsMap
		}
	}

	if d.HasChange("system_disk_size") {
		update = true
		objectDataLocalMap1["system_disk_size"] = d.Get("system_disk_size")
	}

	if d.HasChange("password") {
		update = true
		objectDataLocalMap1["login_password"] = d.Get("password")
	}

	password := d.Get("password").(string)
	if password == "" {
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}
			objectDataLocalMap1["login_password"] = decryptResp
		}
	}
	if d.HasChange("private_pool_options") {
		update = true
		if v := d.Get("private_pool_options"); v != nil {
			private_pool_options := make(map[string]interface{})
			privatePoolOptionsMatchCriteria, _ := jsonpath.Get("$[0].private_pool_options_match_criteria", v)
			if privatePoolOptionsMatchCriteria != nil && privatePoolOptionsMatchCriteria != "" {
				private_pool_options["match_criteria"] = privatePoolOptionsMatchCriteria
			}
			privatePoolOptionsId, _ := jsonpath.Get("$[0].private_pool_options_id", v)
			if privatePoolOptionsId != nil && privatePoolOptionsId != "" {
				private_pool_options["id"] = privatePoolOptionsId
			}

			objectDataLocalMap1["private_pool_options"] = private_pool_options
		}
	}

	if d.HasChange("system_disk_provisioned_iops") {
		update = true
		objectDataLocalMap1["system_disk_provisioned_iops"] = d.Get("system_disk_provisioned_iops")
	}

	if d.HasChange("system_disk_bursting_enabled") {
		update = true
		objectDataLocalMap1["system_disk_bursting_enabled"] = d.Get("system_disk_bursting_enabled")
	}

	if d.HasChange("system_disk_encrypted") {
		update = true
		objectDataLocalMap1["system_disk_encrypted"] = d.Get("system_disk_encrypted")
	}

	if d.HasChange("system_disk_categories") {
		update = true
		systemDiskCategories, _ := jsonpath.Get("$", d.Get("system_disk_categories"))
		if systemDiskCategories != nil && systemDiskCategories != "" {
			objectDataLocalMap1["system_disk_categories"] = systemDiskCategories
		}
	}

	if d.HasChange("system_disk_encrypt_algorithm") {
		update = true
		objectDataLocalMap1["system_disk_encrypt_algorithm"] = d.Get("system_disk_encrypt_algorithm")
	}

	if d.HasChange("image_type") {
		update = true
		objectDataLocalMap1["image_type"] = d.Get("image_type")
	}

	if d.HasChange("instance_types") {
		update = true
		instanceTypes, _ := jsonpath.Get("$", d.Get("instance_types"))
		if instanceTypes != nil && instanceTypes != "" {
			objectDataLocalMap1["instance_types"] = instanceTypes
		}
	}

	if d.HasChange("rds_instances") {
		update = true
		rdsInstances, _ := jsonpath.Get("$", d.Get("rds_instances"))
		if rdsInstances != nil && rdsInstances != "" {
			objectDataLocalMap1["rds_instances"] = rdsInstances
		}
	}

	if d.HasChange("system_disk_kms_key") {
		update = true
		objectDataLocalMap1["system_disk_kms_key_id"] = d.Get("system_disk_kms_key")
	}

	if d.HasChange("system_disk_snapshot_policy_id") {
		update = true
		objectDataLocalMap1["worker_system_disk_snapshot_policy_id"] = d.Get("system_disk_snapshot_policy_id")
	}

	request["scaling_group"] = objectDataLocalMap1
	objectDataLocalMap2 := make(map[string]interface{})

	if d.HasChange("cpu_policy") {
		update = true
		objectDataLocalMap2["cpu_policy"] = d.Get("cpu_policy")
	}

	if d.HasChange("install_cloud_monitor") {
		update = true
		objectDataLocalMap2["cms_enabled"] = d.Get("install_cloud_monitor")
	}

	if d.HasChange("runtime_version") {
		update = true
		objectDataLocalMap2["runtime_version"] = d.Get("runtime_version")
	}

	if d.HasChange("user_data") {
		update = true
		objectDataLocalMap2["user_data"] = d.Get("user_data")
		if v := d.Get("user_data").(string); v != "" {
			_, base64DecodeError := base64.StdEncoding.DecodeString(v)
			if base64DecodeError == nil {
				objectDataLocalMap2["user_data"] = tea.String(v)
			} else {
				objectDataLocalMap2["user_data"] = tea.String(base64.StdEncoding.EncodeToString([]byte(v)))
			}
		}
	}

	if d.HasChange("taints") {
		update = true
		if v := d.Get("taints"); v != nil {
			if v, ok := d.GetOk("taints"); ok {
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
					dataLoop3Map["key"] = dataLoop3Tmp["key"]
					dataLoop3Map["effect"] = dataLoop3Tmp["effect"]
					dataLoop3Map["value"] = dataLoop3Tmp["value"]
					localMaps3 = append(localMaps3, dataLoop3Map)
				}
				objectDataLocalMap2["taints"] = localMaps3
			}

		}
	}

	if d.HasChange("runtime_name") {
		update = true
		objectDataLocalMap2["runtime"] = d.Get("runtime_name")
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
					dataLoop4Map["key"] = dataLoop4Tmp["key"]
					dataLoop4Map["value"] = dataLoop4Tmp["value"]
					localMaps4 = append(localMaps4, dataLoop4Map)
				}
				objectDataLocalMap2["labels"] = localMaps4
			}

		}
	}

	if d.HasChange("unschedulable") {
		update = true
		objectDataLocalMap2["unschedulable"] = d.Get("unschedulable")
	}

	if d.HasChange("pre_user_data") {
		update = true
		objectDataLocalMap2["pre_user_data"] = d.Get("pre_user_data")
	}

	request["kubernetes_config"] = objectDataLocalMap2
	if d.HasChange("scaling_config") {
		update = true
		objectDataLocalMap3 := make(map[string]interface{})

		if v := d.Get("scaling_config"); !IsNil(v) {
			type1, _ := jsonpath.Get("$[0].type", v)
			if type1 != nil && type1 != "" {
				objectDataLocalMap3["type"] = type1
			}
			enable1, _ := jsonpath.Get("$[0].enable", v)
			if enable1 != nil && enable1 != "" {
				objectDataLocalMap3["enable"] = enable1
			}
			maxSize, _ := jsonpath.Get("$[0].max_size", v)
			if maxSize != nil && maxSize != "" {
				objectDataLocalMap3["max_instances"] = maxSize
			}
			minSize, _ := jsonpath.Get("$[0].min_size", v)
			if minSize != nil && minSize != "" {
				objectDataLocalMap3["min_instances"] = minSize
			}
			eipBandwidth, _ := jsonpath.Get("$[0].eip_bandwidth", v)
			if eipBandwidth != nil && eipBandwidth != "" && eipBandwidth.(int) > 0 {
				objectDataLocalMap3["eip_bandwidth"] = eipBandwidth
			}
			eipInternetChargeType, _ := jsonpath.Get("$[0].eip_internet_charge_type", v)
			if eipInternetChargeType != nil && eipInternetChargeType != "" {
				objectDataLocalMap3["eip_internet_charge_type"] = eipInternetChargeType
			}
			isBondEip, _ := jsonpath.Get("$[0].is_bond_eip", v)
			if isBondEip != nil && isBondEip != "" {
				objectDataLocalMap3["is_bond_eip"] = isBondEip
			}

			request["auto_scaling"] = objectDataLocalMap3
		}
	}

	if d.HasChange("management") {
		update = true
		objectDataLocalMap4 := make(map[string]interface{})

		if v := d.Get("management"); !IsNil(v) {
			enable3, _ := jsonpath.Get("$[0].enable", v)
			if enable3 != nil && enable3 != "" {
				objectDataLocalMap4["enable"] = enable3
			}
			autoRepair, _ := jsonpath.Get("$[0].auto_repair", v)
			if autoRepair != nil && autoRepair != "" {
				objectDataLocalMap4["auto_repair"] = autoRepair
			}
			auto_repair_policy := make(map[string]interface{})
			restartNode, _ := jsonpath.Get("$[0].auto_repair_policy[0].restart_node", v)
			if restartNode != nil && restartNode != "" {
				auto_repair_policy["restart_node"] = restartNode
			}

			objectDataLocalMap4["auto_repair_policy"] = auto_repair_policy
			autoVulFix, _ := jsonpath.Get("$[0].auto_vul_fix", v)
			if autoVulFix != nil && autoVulFix != "" {
				objectDataLocalMap4["auto_vul_fix"] = autoVulFix
			}
			auto_vul_fix_policy := make(map[string]interface{})
			restartNode1, _ := jsonpath.Get("$[0].auto_vul_fix_policy[0].restart_node", v)
			if restartNode1 != nil && restartNode1 != "" {
				auto_vul_fix_policy["restart_node"] = restartNode1
			}
			vulLevel, _ := jsonpath.Get("$[0].auto_vul_fix_policy[0].vul_level", v)
			if vulLevel != nil && vulLevel != "" {
				auto_vul_fix_policy["vul_level"] = vulLevel
			}

			objectDataLocalMap4["auto_vul_fix_policy"] = auto_vul_fix_policy
			autoUpgrade, _ := jsonpath.Get("$[0].auto_upgrade", v)
			if autoUpgrade != nil && autoUpgrade != "" {
				objectDataLocalMap4["auto_upgrade"] = autoUpgrade
			}
			auto_upgrade_policy := make(map[string]interface{})
			autoUpgradeKubelet, _ := jsonpath.Get("$[0].auto_upgrade_policy[0].auto_upgrade_kubelet", v)
			if autoUpgradeKubelet != nil && autoUpgradeKubelet != "" {
				auto_upgrade_policy["auto_upgrade_kubelet"] = autoUpgradeKubelet
			}

			objectDataLocalMap4["auto_upgrade_policy"] = auto_upgrade_policy
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

			request["management"] = objectDataLocalMap4
		}
	}

	if v, ok := d.GetOkExists("update_nodes"); ok {
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
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2015-12-15"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
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
	conn, err = client.NewAckClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("kubelet_configuration") {
		update = true
		objectDataLocalMap := make(map[string]interface{})

		if v := d.Get("kubelet_configuration"); !IsNil(v) {
			registryBurst1Raw, _ := jsonpath.Get("$[0].registry_burst", v)
			if registryBurst1Raw != nil && registryBurst1Raw != "" {
				registryBurst1, _ := strconv.ParseInt(registryBurst1Raw.(string), 10, 64)
				objectDataLocalMap["registryBurst"] = registryBurst1
			}
			registryPullQpsRaw, _ := jsonpath.Get("$[0].registry_pull_qps", v)
			if registryPullQpsRaw != nil && registryPullQpsRaw != "" {
				registryPullQps, _ := strconv.ParseInt(registryPullQpsRaw.(string), 10, 64)
				objectDataLocalMap["registryPullQPS"] = registryPullQps
			}
			eventRecordQpsRaw, _ := jsonpath.Get("$[0].event_record_qps", v)
			if eventRecordQpsRaw != nil && eventRecordQpsRaw != "" {
				eventRecordQps, _ := strconv.ParseInt(eventRecordQpsRaw.(string), 10, 64)
				objectDataLocalMap["eventRecordQPS"] = eventRecordQps
			}
			eventBurst1Raw, _ := jsonpath.Get("$[0].event_burst", v)
			if eventBurst1Raw != nil && eventBurst1Raw != "" {
				eventBurst1, _ := strconv.ParseInt(eventBurst1Raw.(string), 10, 64)
				objectDataLocalMap["eventBurst"] = eventBurst1
			}
			kubeApiQpsRaw, _ := jsonpath.Get("$[0].kube_api_qps", v)
			if kubeApiQpsRaw != nil && kubeApiQpsRaw != "" {
				kubeApiQps, _ := strconv.ParseInt(kubeApiQpsRaw.(string), 10, 64)
				objectDataLocalMap["kubeAPIQPS"] = kubeApiQps
			}
			serializeImagePulls1Raw, _ := jsonpath.Get("$[0].serialize_image_pulls", v)
			if serializeImagePulls1Raw != nil && serializeImagePulls1Raw != "" {
				serializeImagePulls1, _ := strconv.ParseBool(serializeImagePulls1Raw.(string))
				objectDataLocalMap["serializeImagePulls"] = serializeImagePulls1
			}
			cpuManagerPolicy1, _ := jsonpath.Get("$[0].cpu_manager_policy", v)
			if cpuManagerPolicy1 != nil && cpuManagerPolicy1 != "" {
				objectDataLocalMap["cpuManagerPolicy"] = cpuManagerPolicy1
			}
			evictionHard1, _ := jsonpath.Get("$[0].eviction_hard", v)
			if evictionHard1 != nil && evictionHard1 != "" {
				objectDataLocalMap["evictionHard"] = evictionHard1
			}
			evictionSoft1, _ := jsonpath.Get("$[0].eviction_soft", v)
			if evictionSoft1 != nil && evictionSoft1 != "" {
				objectDataLocalMap["evictionSoft"] = evictionSoft1
			}
			evictionSoftGracePeriod1, _ := jsonpath.Get("$[0].eviction_soft_grace_period", v)
			if evictionSoftGracePeriod1 != nil && evictionSoftGracePeriod1 != "" {
				objectDataLocalMap["evictionSoftGracePeriod"] = evictionSoftGracePeriod1
			}
			systemReserved1, _ := jsonpath.Get("$[0].system_reserved", v)
			if systemReserved1 != nil && systemReserved1 != "" {
				objectDataLocalMap["systemReserved"] = systemReserved1
			}
			kubeReserved1, _ := jsonpath.Get("$[0].kube_reserved", v)
			if kubeReserved1 != nil && kubeReserved1 != "" {
				objectDataLocalMap["kubeReserved"] = kubeReserved1
			}
			readOnlyPort1Raw, _ := jsonpath.Get("$[0].read_only_port", v)
			if readOnlyPort1Raw != nil && readOnlyPort1Raw != "" {
				readOnlyPort1, _ := strconv.ParseInt(readOnlyPort1Raw.(string), 10, 64)
				objectDataLocalMap["readOnlyPort"] = readOnlyPort1
			}
			maxPods1Raw, _ := jsonpath.Get("$[0].max_pods", v)
			if maxPods1Raw != nil && maxPods1Raw != "" {
				maxPods1, _ := strconv.ParseInt(maxPods1Raw.(string), 10, 64)
				objectDataLocalMap["maxPods"] = maxPods1
			}
			containerLogMaxSize1, _ := jsonpath.Get("$[0].container_log_max_size", v)
			if containerLogMaxSize1 != nil && containerLogMaxSize1 != "" {
				objectDataLocalMap["containerLogMaxSize"] = containerLogMaxSize1
			}
			containerLogMaxFiles1Raw, _ := jsonpath.Get("$[0].container_log_max_files", v)
			if containerLogMaxFiles1Raw != nil && containerLogMaxFiles1Raw != "" {
				containerLogMaxFiles1, _ := strconv.ParseInt(containerLogMaxFiles1Raw.(string), 10, 64)
				objectDataLocalMap["containerLogMaxFiles"] = containerLogMaxFiles1
			}
			featureGates1, _ := jsonpath.Get("$[0].feature_gates", v)
			if featureGates1 != nil && featureGates1 != "" {
				objectDataLocalMap["featureGates"] = featureGates1
			}
			allowedUnsafeSysctls1, _ := jsonpath.Get("$[0].allowed_unsafe_sysctls", d.Get("kubelet_configuration"))
			if allowedUnsafeSysctls1 != nil && allowedUnsafeSysctls1 != "" {
				objectDataLocalMap["allowedUnsafeSysctls"] = allowedUnsafeSysctls1
			}
			kubeApiBurstRaw, _ := jsonpath.Get("$[0].kube_api_burst", v)
			if kubeApiBurstRaw != nil && kubeApiBurstRaw != "" {
				kubeApiBurst, _ := strconv.ParseInt(kubeApiBurstRaw.(string), 10, 64)
				objectDataLocalMap["kubeAPIBurst"] = kubeApiBurst
			}

			request["kubelet_config"] = objectDataLocalMap
		}
	}

	objectDataLocalMap1 = make(map[string]interface{})

	if v := d.Get("rolling_policy"); !IsNil(v) {
		maxParallelism, _ := jsonpath.Get("$[0].max_parallelism", v)
		if maxParallelism != nil && maxParallelism != "" {
			objectDataLocalMap1["max_parallelism"] = maxParallelism
		}

		request["rolling_policy"] = objectDataLocalMap1
	}

	body = request
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2015-12-15"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
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
	conn, err := client.NewAckClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	query["force"] = tea.String("true")
	if v, ok := d.GetOk("force_delete"); ok {
		query["force"] = StringPointer(strconv.FormatBool(v.(bool)))
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2015-12-15"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), query, nil, nil, &runtime)

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
