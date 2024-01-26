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

	roacs "github.com/alibabacloud-go/cs-20151215/v4/client"
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
				DiffSuppressFunc: csNodepoolInstancePostPaidDiffSuppressFunc,
				ValidateFunc:     IntInSlice([]int{0, 1, 2, 3, 6, 12}),
			},
			"cis_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
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
				ValidateFunc: StringInSlice([]string{"static", "none"}, true),
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
						"snapshot_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"category": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"cloud_efficiency", "cloud_ssd", "cloud_essd", "cloud_auto", "cloud", "cloud_essd_xc0", "cloud_essd_xc1", "all", "ephemeral_ssd", "local_disk"}, true),
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
						"provisioned_iops": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
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
				ValidateFunc: StringInSlice([]string{"AliyunLinux", "AliyunLinux3", "AliyunLinux3Arm64", "AliyunLinuxUEFI", "CentOS", "Windows", "WindowsCore", "AliyunLinux Qboot", "ContainerOS", "AliyunLinuxSecurity"}, true),
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
				ValidateFunc: StringInSlice([]string{"PrePaid", "PostPaid"}, true),
			},
			"instance_types": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"PayByBandwidth", "PayByTraffic"}, true),
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
				ValidateFunc: StringInSlice([]string{"PRIORITY", "COST_OPTIMIZED", "BALANCE"}, true),
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
				Computed:     true,
				ExactlyOneOf: []string{"node_pool_name", "name"},
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
				DiffSuppressFunc: csNodepoolInstancePostPaidDiffSuppressFunc,
				ValidateFunc:     IntInSlice([]int{0, 1, 2, 3, 6, 12}),
			},
			"period_unit": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: csNodepoolInstancePostPaidDiffSuppressFunc,
				ValidateFunc:     StringInSlice([]string{"Month"}, true),
			},
			"platform": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Deprecated:   "Field 'platform' has been deprecated from provider version 1.145.0. Operating system release, using `image_type` instead.",
				ValidateFunc: StringInSlice([]string{"CentOS", "AliyunLinux", "Windows", "WindowsCore"}, true),
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
				ConflictsWith: []string{"instances", "node_count", "desired_size"},
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
							ValidateFunc: StringInSlice([]string{"cpu", "gpu", "gpushare", "spot"}, true),
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
							ValidateFunc:  StringInSlice([]string{"PayByBandwidth", "PayByTraffic"}, true),
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
				DiffSuppressFunc: csNodepoolScalingPolicyDiffSuppressFunc,
				ValidateFunc:     StringInSlice([]string{"release", "recycle"}, true),
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
				ValidateFunc: StringInSlice([]string{"cloud_efficiency", "cloud_ssd", "cloud_essd", "cloud_auto"}, true),
			},
			"system_disk_encrypt_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"aes-256"}, true),
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
	request["ClusterId"] = d.Get("cluster_id")

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
		nodeNative2, _ := jsonpath.Get("$", v)
		if nodeNative2 != nil && nodeNative2 != "" {
			objectDataLocalMap1["security_group_ids"] = nodeNative2
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
				return WrapError(err)
			}
			localMaps := make([]map[string]interface{}, 0)
			for _, dataLoop := range localData.([]interface{}) {
				dataLoopTmp := dataLoop.(map[string]interface{})
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
				return WrapError(err)
			}
			localMaps1 := make([]map[string]interface{}, 0)
			for _, dataLoop1 := range localData1.([]interface{}) {
				dataLoop1Tmp := dataLoop1.(map[string]interface{})
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
		nodeNative41, _ := jsonpath.Get("$", v)
		if nodeNative41 != nil && nodeNative41 != "" {
			objectDataLocalMap1["vswitch_ids"] = nodeNative41
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
		nodeNative47, _ := jsonpath.Get("$[0].private_pool_options_match_criteria", d.Get("private_pool_options"))
		if nodeNative47 != nil && nodeNative47 != "" {
			private_pool_options["match_criteria"] = nodeNative47
		}
		nodeNative48, _ := jsonpath.Get("$[0].private_pool_options_id", d.Get("private_pool_options"))
		if nodeNative48 != nil && nodeNative48 != "" {
			private_pool_options["id"] = nodeNative48
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
		nodeNative54, _ := jsonpath.Get("$", v)
		if nodeNative54 != nil && nodeNative54 != "" {
			objectDataLocalMap1["system_disk_categories"] = nodeNative54
		}
	}
	if v, ok := d.GetOk("instance_types"); ok {
		nodeNative55, _ := jsonpath.Get("$", v)
		if nodeNative55 != nil && nodeNative55 != "" {
			objectDataLocalMap1["instance_types"] = nodeNative55
		}
	}
	if v, ok := d.GetOk("rds_instances"); ok {
		nodeNative56, _ := jsonpath.Get("$", v)
		if nodeNative56 != nil && nodeNative56 != "" {
			objectDataLocalMap1["rds_instances"] = nodeNative56
		}
	}
	if v, ok := d.GetOk("system_disk_kms_key"); ok {
		objectDataLocalMap1["system_disk_kms_key_id"] = v
	}
	if v, ok := d.GetOk("system_disk_snapshot_policy_id"); ok {
		objectDataLocalMap1["worker_system_disk_snapshot_policy_id"] = v
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
				return WrapError(err)
			}
			localMaps3 := make([]map[string]interface{}, 0)
			for _, dataLoop3 := range localData3.([]interface{}) {
				dataLoop3Tmp := dataLoop3.(map[string]interface{})
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
				return WrapError(err)
			}
			localMaps4 := make([]map[string]interface{}, 0)
			for _, dataLoop4 := range localData4.([]interface{}) {
				dataLoop4Tmp := dataLoop4.(map[string]interface{})
				dataLoop4Map := make(map[string]interface{})
				dataLoop4Map["key"] = dataLoop4Tmp["key"]
				dataLoop4Map["value"] = dataLoop4Tmp["value"]
				localMaps4 = append(localMaps4, dataLoop4Map)
			}
			objectDataLocalMap2["labels"] = localMaps4
		}
	}
	request["kubernetes_config"] = objectDataLocalMap2
	objectDataLocalMap3 := make(map[string]interface{})
	if v := d.Get("scaling_config"); !IsNil(v) {
		nodeNative71, _ := jsonpath.Get("$[0].type", d.Get("scaling_config"))
		if nodeNative71 != nil && nodeNative71 != "" {
			objectDataLocalMap3["type"] = nodeNative71
		}
		nodeNative72, _ := jsonpath.Get("$[0].max_size", d.Get("scaling_config"))
		if nodeNative72 != nil && nodeNative72 != "" {
			objectDataLocalMap3["max_instances"] = nodeNative72
		}
		nodeNative73, _ := jsonpath.Get("$[0].min_size", d.Get("scaling_config"))
		if nodeNative73 != nil && nodeNative73 != "" {
			objectDataLocalMap3["min_instances"] = nodeNative73
		}
		nodeNative74, _ := jsonpath.Get("$[0].is_bond_eip", d.Get("scaling_config"))
		if nodeNative74 != nil && nodeNative74 != "" {
			objectDataLocalMap3["is_bond_eip"] = nodeNative74
		}
		nodeNative75, _ := jsonpath.Get("$[0].enable", d.Get("scaling_config"))
		if nodeNative75 != nil && nodeNative75 != "" {
			objectDataLocalMap3["enable"] = nodeNative75
		}
		nodeNative76, _ := jsonpath.Get("$[0].eip_internet_charge_type", d.Get("scaling_config"))
		if nodeNative76 != nil && nodeNative76 != "" {
			objectDataLocalMap3["eip_internet_charge_type"] = nodeNative76
		}
		nodeNative77, _ := jsonpath.Get("$[0].eip_bandwidth", d.Get("scaling_config"))
		if nodeNative77 != nil && nodeNative77 != "" && nodeNative77.(int) > 0 {
			objectDataLocalMap3["eip_bandwidth"] = nodeNative77
		}
		request["auto_scaling"] = objectDataLocalMap3
	}

	objectDataLocalMap4 := make(map[string]interface{})
	if v := d.Get("management"); !IsNil(v) {
		nodeNative78, _ := jsonpath.Get("$[0].enable", d.Get("management"))
		if nodeNative78 != nil && nodeNative78 != "" {
			objectDataLocalMap4["enable"] = nodeNative78
		}
		nodeNative79, _ := jsonpath.Get("$[0].auto_repair", d.Get("management"))
		if nodeNative79 != nil && nodeNative79 != "" {
			objectDataLocalMap4["auto_repair"] = nodeNative79
		}
		auto_repair_policy := make(map[string]interface{})
		nodeNative80, _ := jsonpath.Get("$[0].auto_repair_policy[0].restart_node", d.Get("management"))
		if nodeNative80 != nil && nodeNative80 != "" {
			auto_repair_policy["restart_node"] = nodeNative80
		}
		objectDataLocalMap4["auto_repair_policy"] = auto_repair_policy
		nodeNative81, _ := jsonpath.Get("$[0].auto_vul_fix", d.Get("management"))
		if nodeNative81 != nil && nodeNative81 != "" {
			objectDataLocalMap4["auto_vul_fix"] = nodeNative81
		}
		auto_vul_fix_policy := make(map[string]interface{})
		nodeNative82, _ := jsonpath.Get("$[0].auto_vul_fix_policy[0].restart_node", d.Get("management"))
		if nodeNative82 != nil && nodeNative82 != "" {
			auto_vul_fix_policy["restart_node"] = nodeNative82
		}
		nodeNative83, _ := jsonpath.Get("$[0].auto_vul_fix_policy[0].vul_level", d.Get("management"))
		if nodeNative83 != nil && nodeNative83 != "" {
			auto_vul_fix_policy["vul_level"] = nodeNative83
		}
		objectDataLocalMap4["auto_vul_fix_policy"] = auto_vul_fix_policy
		nodeNative84, _ := jsonpath.Get("$[0].auto_upgrade", d.Get("management"))
		if nodeNative84 != nil && nodeNative84 != "" {
			objectDataLocalMap4["auto_upgrade"] = nodeNative84
		}
		auto_upgrade_policy := make(map[string]interface{})
		nodeNative85, _ := jsonpath.Get("$[0].auto_upgrade_policy[0].auto_upgrade_kubelet", d.Get("management"))
		if nodeNative85 != nil && nodeNative85 != "" {
			auto_upgrade_policy["auto_upgrade_kubelet"] = nodeNative85
		}
		objectDataLocalMap4["auto_upgrade_policy"] = auto_upgrade_policy
		upgrade_config := make(map[string]interface{})
		nodeNative86, _ := jsonpath.Get("$[0].surge", d.Get("management"))
		if nodeNative86 != nil && nodeNative86 != "" {
			upgrade_config["surge"] = nodeNative86
		}
		nodeNative87, _ := jsonpath.Get("$[0].surge_percentage", d.Get("management"))
		if nodeNative87 != nil && nodeNative87 != "" {
			upgrade_config["surge_percentage"] = nodeNative87
		}
		nodeNative88, _ := jsonpath.Get("$[0].max_unavailable", d.Get("management"))
		if nodeNative88 != nil && nodeNative88 != "" && nodeNative88.(int) > 0 {
			upgrade_config["max_unavailable"] = nodeNative88
		}
		objectDataLocalMap4["upgrade_config"] = upgrade_config
		request["management"] = objectDataLocalMap4
	}

	objectDataLocalMap5 := make(map[string]interface{})
	if v := d.Get("tee_config"); !IsNil(v) {
		nodeNative89, _ := jsonpath.Get("$[0].tee_enable", d.Get("tee_config"))
		if nodeNative89 != nil && nodeNative89 != "" {
			objectDataLocalMap5["tee_enable"] = nodeNative89
		}
		request["tee_config"] = objectDataLocalMap5
	}

	objectDataLocalMap6 := make(map[string]interface{})
	if v := d.Get("kubelet_configuration"); !IsNil(v) {
		kubelet_configuration := make(map[string]interface{})
		nodeNative90, _ := jsonpath.Get("$[0].registry_pull_qps", d.Get("kubelet_configuration"))
		if nodeNative90 != nil && nodeNative90 != "" {
			intVal, _ := strconv.ParseInt(nodeNative90.(string), 10, 64)
			kubelet_configuration["registryPullQPS"] = intVal
		}
		nodeNative91, _ := jsonpath.Get("$[0].registry_burst", d.Get("kubelet_configuration"))
		if nodeNative91 != nil && nodeNative91 != "" {
			intVal, _ := strconv.ParseInt(nodeNative91.(string), 10, 64)
			kubelet_configuration["registryBurst"] = intVal
		}
		nodeNative92, _ := jsonpath.Get("$[0].event_record_qps", d.Get("kubelet_configuration"))
		if nodeNative92 != nil && nodeNative92 != "" {
			intVal, _ := strconv.ParseInt(nodeNative92.(string), 10, 64)
			kubelet_configuration["eventRecordQPS"] = intVal
		}
		nodeNative93, _ := jsonpath.Get("$[0].event_burst", d.Get("kubelet_configuration"))
		if nodeNative93 != nil && nodeNative93 != "" {
			intVal, _ := strconv.ParseInt(nodeNative93.(string), 10, 64)
			kubelet_configuration["eventBurst"] = intVal
		}
		nodeNative94, _ := jsonpath.Get("$[0].kube_api_qps", d.Get("kubelet_configuration"))
		if nodeNative94 != nil && nodeNative94 != "" {
			intVal, _ := strconv.ParseInt(nodeNative94.(string), 10, 64)
			kubelet_configuration["kubeAPIQPS"] = intVal
		}
		nodeNative95, _ := jsonpath.Get("$[0].serialize_image_pulls", d.Get("kubelet_configuration"))
		if nodeNative95 != nil && nodeNative95 != "" {
			boolVal, _ := strconv.ParseBool(nodeNative95.(string))
			kubelet_configuration["serializeImagePulls"] = boolVal
		}
		nodeNative96, _ := jsonpath.Get("$[0].cpu_manager_policy", d.Get("kubelet_configuration"))
		if nodeNative96 != nil && nodeNative96 != "" {
			kubelet_configuration["cpuManagerPolicy"] = nodeNative96
		}
		nodeNative97, _ := jsonpath.Get("$[0].allowed_unsafe_sysctls", v)
		if nodeNative97 != nil && nodeNative97 != "" {
			kubelet_configuration["allowedUnsafeSysctls"] = nodeNative97
		}
		nodeNative98, _ := jsonpath.Get("$[0].feature_gates", d.Get("kubelet_configuration"))
		if nodeNative98 != nil && nodeNative98 != "" {
			kubelet_configuration["featureGates"] = nodeNative98
		}
		nodeNative99, _ := jsonpath.Get("$[0].container_log_max_files", d.Get("kubelet_configuration"))
		if nodeNative99 != nil && nodeNative99 != "" {
			intVal, _ := strconv.ParseInt(nodeNative99.(string), 10, 64)
			kubelet_configuration["containerLogMaxFiles"] = intVal
		}
		nodeNative100, _ := jsonpath.Get("$[0].container_log_max_size", d.Get("kubelet_configuration"))
		if nodeNative100 != nil && nodeNative100 != "" {
			kubelet_configuration["containerLogMaxSize"] = nodeNative100
		}
		nodeNative101, _ := jsonpath.Get("$[0].max_pods", d.Get("kubelet_configuration"))
		if nodeNative101 != nil && nodeNative101 != "" {
			intVal, _ := strconv.ParseInt(nodeNative101.(string), 10, 64)
			kubelet_configuration["maxPods"] = intVal
		}
		nodeNative102, _ := jsonpath.Get("$[0].read_only_port", d.Get("kubelet_configuration"))
		if nodeNative102 != nil && nodeNative102 != "" {
			intVal, _ := strconv.ParseInt(nodeNative102.(string), 10, 64)
			kubelet_configuration["readOnlyPort"] = intVal
		}
		nodeNative103, _ := jsonpath.Get("$[0].kube_reserved", d.Get("kubelet_configuration"))
		if nodeNative103 != nil && nodeNative103 != "" {
			kubelet_configuration["kubeReserved"] = nodeNative103
		}
		nodeNative104, _ := jsonpath.Get("$[0].system_reserved", d.Get("kubelet_configuration"))
		if nodeNative104 != nil && nodeNative104 != "" {
			kubelet_configuration["systemReserved"] = nodeNative104
		}
		nodeNative105, _ := jsonpath.Get("$[0].eviction_soft_grace_period", d.Get("kubelet_configuration"))
		if nodeNative105 != nil && nodeNative105 != "" {
			kubelet_configuration["evictionSoftGracePeriod"] = nodeNative105
		}
		nodeNative106, _ := jsonpath.Get("$[0].eviction_soft", d.Get("kubelet_configuration"))
		if nodeNative106 != nil && nodeNative106 != "" {
			kubelet_configuration["evictionSoft"] = nodeNative106
		}
		nodeNative107, _ := jsonpath.Get("$[0].eviction_hard", d.Get("kubelet_configuration"))
		if nodeNative107 != nil && nodeNative107 != "" {
			kubelet_configuration["evictionHard"] = nodeNative107
		}
		nodeNative108, _ := jsonpath.Get("$[0].kube_api_burst", d.Get("kubelet_configuration"))
		if nodeNative108 != nil && nodeNative108 != "" {
			intVal, _ := strconv.ParseInt(nodeNative108.(string), 10, 64)
			kubelet_configuration["kubeAPIBurst"] = intVal
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
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_kubernetes_node_pool", action, AlibabaCloudSdkGoERROR)
	}

	nodepool_id, _ := jsonpath.Get("$.body.nodepool_id", response)
	d.SetId(fmt.Sprintf("%v:%v", ClusterId, nodepool_id))

	ackServiceV2 := AckServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ackServiceV2.DescribeAsyncAckNodepoolStateRefreshFunc(d, response, "$.state", []string{"fail", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	if v, ok := d.GetOk("instances"); ok && v != nil {
		attachExistingInstance(d, meta)
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
			dataDisksMap["auto_snapshot_policy_id"] = data_disksChild1Raw["auto_snapshot_policy_id"]
			dataDisksMap["bursting_enabled"] = data_disksChild1Raw["bursting_enabled"]
			dataDisksMap["category"] = data_disksChild1Raw["category"]
			dataDisksMap["device"] = data_disksChild1Raw["device"]
			dataDisksMap["encrypted"] = data_disksChild1Raw["encrypted"]
			dataDisksMap["kms_key_id"] = data_disksChild1Raw["kms_key_id"]
			dataDisksMap["name"] = data_disksChild1Raw["disk_name"]
			dataDisksMap["performance_level"] = data_disksChild1Raw["performance_level"]
			dataDisksMap["provisioned_iops"] = data_disksChild1Raw["provisioned_iops"]
			dataDisksMap["size"] = data_disksChild1Raw["size"]
			dataDisksMap["snapshot_id"] = data_disksChild1Raw["snapshot_id"]

			dataDisksMaps = append(dataDisksMaps, dataDisksMap)
		}
	}
	d.Set("data_disks", dataDisksMaps)
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
	d.Set("kubelet_configuration", kubeletConfigurationMaps)
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
	d.Set("labels", labelsMaps)
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
	d.Set("management", managementMaps)
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
	d.Set("private_pool_options", privatePoolOptionsMaps)
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
	d.Set("scaling_config", scalingConfigMaps)
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
	d.Set("spot_price_limit", spotPriceLimitMaps)
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
	d.Set("taints", taintsMaps)
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
	d.Set("tee_config", teeConfigMaps)
	vswitch_ids1Raw, _ := jsonpath.Get("$.scaling_group.vswitch_ids", objectRaw)
	d.Set("vswitch_ids", vswitch_ids1Raw)

	parts := strings.Split(d.Id(), ":")
	d.Set("cluster_id", parts[0])
	d.Set("node_pool_id", parts[1])

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
	request["ClusterId"] = parts[0]
	request["NodepoolId"] = parts[1]
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
					return WrapError(err)
				}
				localMaps := make([]map[string]interface{}, 0)
				for _, dataLoop := range localData.([]interface{}) {
					dataLoopTmp := dataLoop.(map[string]interface{})
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
					return WrapError(err)
				}
				localMaps1 := make([]map[string]interface{}, 0)
				for _, dataLoop1 := range localData1.([]interface{}) {
					dataLoop1Tmp := dataLoop1.(map[string]interface{})
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
		nodeNative36, _ := jsonpath.Get("$", d.Get("vswitch_ids"))
		if nodeNative36 != nil && nodeNative36 != "" {
			objectDataLocalMap1["vswitch_ids"] = nodeNative36
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
			nodeNative41, _ := jsonpath.Get("$[0].private_pool_options_match_criteria", v)
			if nodeNative41 != nil && nodeNative41 != "" {
				private_pool_options["match_criteria"] = nodeNative41
			}
			nodeNative42, _ := jsonpath.Get("$[0].private_pool_options_id", v)
			if nodeNative42 != nil && nodeNative42 != "" {
				private_pool_options["id"] = nodeNative42
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
		nodeNative46, _ := jsonpath.Get("$", d.Get("system_disk_categories"))
		if nodeNative46 != nil && nodeNative46 != "" {
			objectDataLocalMap1["system_disk_categories"] = nodeNative46
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
		nodeNative49, _ := jsonpath.Get("$", d.Get("instance_types"))
		if nodeNative49 != nil && nodeNative49 != "" {
			objectDataLocalMap1["instance_types"] = nodeNative49
		}
	}
	if d.HasChange("rds_instances") {
		update = true
		nodeNative50, _ := jsonpath.Get("$", d.Get("rds_instances"))
		if nodeNative50 != nil && nodeNative50 != "" {
			objectDataLocalMap1["rds_instances"] = nodeNative50
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
					return WrapError(err)
				}
				localMaps3 := make([]map[string]interface{}, 0)
				for _, dataLoop3 := range localData3.([]interface{}) {
					dataLoop3Tmp := dataLoop3.(map[string]interface{})
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
					return WrapError(err)
				}
				localMaps4 := make([]map[string]interface{}, 0)
				for _, dataLoop4 := range localData4.([]interface{}) {
					dataLoop4Tmp := dataLoop4.(map[string]interface{})
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
	request["kubernetes_config"] = objectDataLocalMap2
	if d.HasChange("scaling_config") {
		update = true
		objectDataLocalMap3 := make(map[string]interface{})
		if v := d.Get("scaling_config"); v != nil {
			nodeNative64, _ := jsonpath.Get("$[0].type", v)
			if nodeNative64 != nil && nodeNative64 != "" {
				objectDataLocalMap3["type"] = nodeNative64
			}
			nodeNative65, _ := jsonpath.Get("$[0].enable", v)
			if nodeNative65 != nil && nodeNative65 != "" {
				objectDataLocalMap3["enable"] = nodeNative65
			}
			nodeNative66, _ := jsonpath.Get("$[0].max_size", v)
			if nodeNative66 != nil && nodeNative66 != "" {
				objectDataLocalMap3["max_instances"] = nodeNative66
			}
			nodeNative67, _ := jsonpath.Get("$[0].min_size", v)
			if nodeNative67 != nil && nodeNative67 != "" {
				objectDataLocalMap3["min_instances"] = nodeNative67
			}
			nodeNative68, _ := jsonpath.Get("$[0].eip_bandwidth", v)
			if nodeNative68 != nil && nodeNative68 != "" && nodeNative68.(int) > 0 {
				objectDataLocalMap3["eip_bandwidth"] = nodeNative68
			}
			nodeNative69, _ := jsonpath.Get("$[0].eip_internet_charge_type", v)
			if nodeNative69 != nil && nodeNative69 != "" {
				objectDataLocalMap3["eip_internet_charge_type"] = nodeNative69
			}
			nodeNative70, _ := jsonpath.Get("$[0].is_bond_eip", v)
			if nodeNative70 != nil && nodeNative70 != "" {
				objectDataLocalMap3["is_bond_eip"] = nodeNative70
			}
			request["auto_scaling"] = objectDataLocalMap3
		}
	}

	if d.HasChange("management") {
		update = true
		objectDataLocalMap4 := make(map[string]interface{})
		if v := d.Get("management"); v != nil {
			nodeNative71, _ := jsonpath.Get("$[0].enable", v)
			if nodeNative71 != nil && nodeNative71 != "" {
				objectDataLocalMap4["enable"] = nodeNative71
			}
			nodeNative72, _ := jsonpath.Get("$[0].auto_repair", v)
			if nodeNative72 != nil && nodeNative72 != "" {
				objectDataLocalMap4["auto_repair"] = nodeNative72
			}
			auto_repair_policy := make(map[string]interface{})
			nodeNative73, _ := jsonpath.Get("$[0].auto_repair_policy[0].restart_node", v)
			if nodeNative73 != nil && nodeNative73 != "" {
				auto_repair_policy["restart_node"] = nodeNative73
			}
			objectDataLocalMap4["auto_repair_policy"] = auto_repair_policy
			nodeNative74, _ := jsonpath.Get("$[0].auto_vul_fix", v)
			if nodeNative74 != nil && nodeNative74 != "" {
				objectDataLocalMap4["auto_vul_fix"] = nodeNative74
			}
			auto_vul_fix_policy := make(map[string]interface{})
			nodeNative75, _ := jsonpath.Get("$[0].auto_vul_fix_policy[0].restart_node", v)
			if nodeNative75 != nil && nodeNative75 != "" {
				auto_vul_fix_policy["restart_node"] = nodeNative75
			}
			nodeNative76, _ := jsonpath.Get("$[0].auto_vul_fix_policy[0].vul_level", v)
			if nodeNative76 != nil && nodeNative76 != "" {
				auto_vul_fix_policy["vul_level"] = nodeNative76
			}
			objectDataLocalMap4["auto_vul_fix_policy"] = auto_vul_fix_policy
			nodeNative77, _ := jsonpath.Get("$[0].auto_upgrade", v)
			if nodeNative77 != nil && nodeNative77 != "" {
				objectDataLocalMap4["auto_upgrade"] = nodeNative77
			}
			auto_upgrade_policy := make(map[string]interface{})
			nodeNative78, _ := jsonpath.Get("$[0].auto_upgrade_policy[0].auto_upgrade_kubelet", v)
			if nodeNative78 != nil && nodeNative78 != "" {
				auto_upgrade_policy["auto_upgrade_kubelet"] = nodeNative78
			}
			objectDataLocalMap4["auto_upgrade_policy"] = auto_upgrade_policy
			upgrade_config := make(map[string]interface{})
			nodeNative79, _ := jsonpath.Get("$[0].surge", v)
			if nodeNative79 != nil && nodeNative79 != "" {
				upgrade_config["surge"] = nodeNative79
			}
			nodeNative80, _ := jsonpath.Get("$[0].surge_percentage", v)
			if nodeNative80 != nil && nodeNative80 != "" {
				upgrade_config["surge_percentage"] = nodeNative80
			}
			nodeNative81, _ := jsonpath.Get("$[0].max_unavailable", v)
			if nodeNative81 != nil && nodeNative81 != "" && nodeNative81.(int) > 0 {
				upgrade_config["max_unavailable"] = nodeNative81
			}
			objectDataLocalMap4["upgrade_config"] = upgrade_config
			request["management"] = objectDataLocalMap4
		}
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
			removeNodePoolNodes(d, meta, parts, nil, nil)
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
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		ackServiceV2 := AckServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ackServiceV2.DescribeAsyncAckNodepoolStateRefreshFunc(d, response, "$.state", []string{"fail", "failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id(), response)
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
	request["ClusterId"] = parts[0]
	request["NodepoolId"] = parts[1]
	if d.HasChange("kubelet_configuration") {
		update = true
		objectDataLocalMap := make(map[string]interface{})
		if v := d.Get("kubelet_configuration"); v != nil {
			nodeNative, _ := jsonpath.Get("$[0].registry_burst", v)
			if nodeNative != nil && nodeNative != "" {
				intVal, _ := strconv.ParseInt(nodeNative.(string), 10, 64)
				objectDataLocalMap["registryBurst"] = intVal
			}
			nodeNative1, _ := jsonpath.Get("$[0].registry_pull_qps", v)
			if nodeNative1 != nil && nodeNative1 != "" {
				intVal, _ := strconv.ParseInt(nodeNative1.(string), 10, 64)
				objectDataLocalMap["registryPullQPS"] = intVal
			}
			nodeNative2, _ := jsonpath.Get("$[0].event_record_qps", v)
			if nodeNative2 != nil && nodeNative2 != "" {
				intVal, _ := strconv.ParseInt(nodeNative2.(string), 10, 64)
				objectDataLocalMap["eventRecordQPS"] = intVal
			}
			nodeNative3, _ := jsonpath.Get("$[0].event_burst", v)
			if nodeNative3 != nil && nodeNative3 != "" {
				intVal, _ := strconv.ParseInt(nodeNative3.(string), 10, 64)
				objectDataLocalMap["eventBurst"] = intVal
			}
			nodeNative4, _ := jsonpath.Get("$[0].kube_api_qps", v)
			if nodeNative4 != nil && nodeNative4 != "" {
				intVal, _ := strconv.ParseInt(nodeNative4.(string), 10, 64)
				objectDataLocalMap["kubeAPIQPS"] = intVal
			}
			nodeNative5, _ := jsonpath.Get("$[0].serialize_image_pulls", v)
			if nodeNative5 != nil && nodeNative5 != "" {
				boolVal, _ := strconv.ParseBool(nodeNative5.(string))
				objectDataLocalMap["serializeImagePulls"] = boolVal
			}
			nodeNative6, _ := jsonpath.Get("$[0].cpu_manager_policy", v)
			if nodeNative6 != nil && nodeNative6 != "" {
				objectDataLocalMap["cpuManagerPolicy"] = nodeNative6
			}
			nodeNative7, _ := jsonpath.Get("$[0].eviction_hard", v)
			if nodeNative7 != nil && nodeNative7 != "" {
				objectDataLocalMap["evictionHard"] = nodeNative7
			}
			nodeNative8, _ := jsonpath.Get("$[0].eviction_soft", v)
			if nodeNative8 != nil && nodeNative8 != "" {
				objectDataLocalMap["evictionSoft"] = nodeNative8
			}
			nodeNative9, _ := jsonpath.Get("$[0].eviction_soft_grace_period", v)
			if nodeNative9 != nil && nodeNative9 != "" {
				objectDataLocalMap["evictionSoftGracePeriod"] = nodeNative9
			}
			nodeNative10, _ := jsonpath.Get("$[0].system_reserved", v)
			if nodeNative10 != nil && nodeNative10 != "" {
				objectDataLocalMap["systemReserved"] = nodeNative10
			}
			nodeNative11, _ := jsonpath.Get("$[0].kube_reserved", v)
			if nodeNative11 != nil && nodeNative11 != "" {
				objectDataLocalMap["kubeReserved"] = nodeNative11
			}
			nodeNative12, _ := jsonpath.Get("$[0].read_only_port", v)
			if nodeNative12 != nil && nodeNative12 != "" {
				intVal, _ := strconv.ParseInt(nodeNative12.(string), 10, 64)
				objectDataLocalMap["readOnlyPort"] = intVal
			}
			nodeNative13, _ := jsonpath.Get("$[0].max_pods", v)
			if nodeNative13 != nil && nodeNative13 != "" {
				intVal, _ := strconv.ParseInt(nodeNative13.(string), 10, 64)
				objectDataLocalMap["maxPods"] = intVal
			}
			nodeNative14, _ := jsonpath.Get("$[0].container_log_max_size", v)
			if nodeNative14 != nil && nodeNative14 != "" {
				objectDataLocalMap["containerLogMaxSize"] = nodeNative14
			}
			nodeNative15, _ := jsonpath.Get("$[0].container_log_max_files", v)
			if nodeNative15 != nil && nodeNative15 != "" {
				intVal, _ := strconv.ParseInt(nodeNative15.(string), 10, 64)
				objectDataLocalMap["containerLogMaxFiles"] = intVal
			}
			nodeNative16, _ := jsonpath.Get("$[0].feature_gates", v)
			if nodeNative16 != nil && nodeNative16 != "" {
				objectDataLocalMap["featureGates"] = nodeNative16
			}
			nodeNative17, _ := jsonpath.Get("$[0].allowed_unsafe_sysctls", d.Get("kubelet_configuration"))
			if nodeNative17 != nil && nodeNative17 != "" {
				objectDataLocalMap["allowedUnsafeSysctls"] = nodeNative17
			}
			nodeNative18, _ := jsonpath.Get("$[0].kube_api_burst", v)
			if nodeNative18 != nil && nodeNative18 != "" {
				intVal, _ := strconv.ParseInt(nodeNative18.(string), 10, 64)
				objectDataLocalMap["kubeAPIBurst"] = intVal
			}
			request["kubelet_config"] = objectDataLocalMap
		}
	}

	objectDataLocalMap1 = make(map[string]interface{})
	if v := d.Get("rolling_policy"); v != nil {
		nodeNative19, _ := jsonpath.Get("$[0].max_parallelism", v)
		if nodeNative19 != nil && nodeNative19 != "" {
			objectDataLocalMap1["max_parallelism"] = nodeNative19
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
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		ackServiceV2 := AckServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ackServiceV2.DescribeAsyncAckNodepoolStateRefreshFunc(d, response, "$.state", []string{"fail", "failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id(), response)
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

		if len(newValue) > len(oldValue) {
			attachExistingInstance(d, meta)
		} else {
			removeNodePoolNodes(d, meta, parts, oldValue, newValue)
		}
	}

	csService := CsService{client}
	_ = resource.Retry(10*time.Minute, func() *resource.RetryError {
		log.Printf("[DEBUG] Start retry fetch node pool info: %s", d.Id())
		nodePoolDetail, err := csService.DescribeCsKubernetesNodePool(d.Id())
		if err != nil {
			return resource.NonRetryableError(err)
		}

		if nodePoolDetail.TotalNodes != d.Get("node_count").(int) && nodePoolDetail.TotalNodes != d.Get("desired_size").(int) {
			time.Sleep(20 * time.Second)
			return resource.RetryableError(Error("[ERROR] The number of nodes is inconsistent %s", d.Id()))
		}

		return resource.NonRetryableError(Error("[DEBUG] The number of nodes is the same"))
	})
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
	body := make(map[string]interface{})
	conn, err := client.NewAckClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["NodepoolId"] = parts[1]
	request["ClusterId"] = parts[0]

	query["force"] = tea.String("true")
	if v, ok := d.GetOk("force_delete"); ok {
		query["force"] = StringPointer(strconv.FormatBool(v.(bool)))
	}
	body = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2015-12-15"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)

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

	ackServiceV2 := AckServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ackServiceV2.DescribeAsyncAckNodepoolStateRefreshFunc(d, response, "$.state", []string{"fail", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
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
		var attachNodeList []string
		if oldNodes != nil && newNodes != nil {
			attachNodeList = difference(expandStringList(oldNodes), expandStringList(newNodes))
		}
		if len(newNodes) == 0 {
			attachNodeList = expandStringList(oldNodes)
		}
		for _, v := range ret {
			for _, name := range attachNodeList {
				if name == v.InstanceId {
					removeInstanceList = append(removeInstanceList, v.NodeName)
				}
			}
		}
		removeNodesName = removeInstanceList
	}

	removeNodesArgs := &cs.DeleteKubernetesClusterNodesRequest{
		Nodes:       removeNodesName,
		ReleaseNode: true,
		DrainNode:   false,
	}
	if err := invoker.Run(func() error {
		var err error
		response, err = client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			resp, err := csClient.DeleteKubernetesClusterNodes(parseId[0], removeNodesArgs)
			return resp, err
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

func attachExistingInstance(d *schema.ResourceData, meta interface{}) error {
	csService := CsService{meta.(*connectivity.AliyunClient)}
	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "InitializeClient", err)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	clusterId := parts[0]
	nodePoolId := parts[1]

	args := &roacs.AttachInstancesRequest{
		NodepoolId:       tea.String(nodePoolId),
		FormatDisk:       tea.Bool(false),
		KeepInstanceName: tea.Bool(true),
	}

	if v, ok := d.GetOk("password"); ok {
		args.Password = tea.String(v.(string))
	}

	if v, ok := d.GetOk("key_name"); ok {
		args.KeyPair = tea.String(v.(string))
	}

	if v, ok := d.GetOk("format_disk"); ok {
		args.FormatDisk = tea.Bool(v.(bool))
	}

	if v, ok := d.GetOk("keep_instance_name"); ok {
		args.KeepInstanceName = tea.Bool(v.(bool))
	}

	if v, ok := d.GetOk("image_id"); ok {
		args.ImageId = tea.String(v.(string))
	}

	if v, ok := d.GetOk("instances"); ok {
		args.Instances = tea.StringSlice(expandStringList(v.([]interface{})))
	}

	_, err = client.AttachInstances(tea.String(clusterId), args)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "AttachInstances", AliyunTablestoreGoSdk)
	}

	stateConf := BuildStateConf([]string{"scaling"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, csService.CsKubernetesNodePoolStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
