package alicloud

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCSManagedKubernetes() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSManagedKubernetesCreate,
		Read:   resourceAlicloudCSKubernetesRead,
		Update: resourceAlicloudCSKubernetesUpdate,
		Delete: resourceAlicloudCSKubernetesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringLenBetween(1, 63),
				ConflictsWith: []string{"name_prefix"},
			},
			"name_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "Terraform-Creation",
				ValidateFunc:  validation.StringLenBetween(0, 37),
				ConflictsWith: []string{"name"},
			},
			// worker configurations，TODO: name issue
			"worker_vswitch_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				MinItems: 1,
			},
			"worker_instance_types": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MaxItems:   10,
				Deprecated: "Field 'worker_instance_types' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'instance_types' to replace it.",
			},
			"worker_number": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "Field 'worker_number' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes., by using field 'desired_size' to replace it.",
			},
			"worker_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(20, 32768),
				Deprecated:   "Field 'worker_disk_size' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'system_disk_size' to replace it.",
			},
			"worker_disk_category": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'worker_disk_category' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'system_disk_category' to replace it.",
			},
			"worker_disk_performance_level": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
				DiffSuppressFunc: workerDiskPerformanceLevelDiffSuppressFunc,
				Deprecated:       "Field 'worker_disk_performance_level' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'system_disk_performance_level' to replace it",
			},
			"worker_disk_snapshot_policy_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'worker_disk_snapshot_policy_id' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'system_disk_snapshot_policy_id' to replace it",
			},
			"worker_data_disk_size": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     validation.IntBetween(20, 32768),
				DiffSuppressFunc: workerDataDiskSizeSuppressFunc,
				Deprecated:       "Field 'worker_data_disk_size' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'data_disks.size' to replace it",
			},
			"worker_data_disk_category": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'worker_data_disk_category' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'data_disks.category' to replace it",
			},
			"worker_instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
				Deprecated:   "Field 'worker_instance_charge_type' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'instance_charge_type' to replace it",
			},
			"worker_data_disks": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeString,
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
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"device": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"encrypted": {
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
				Deprecated: "Field 'worker_data_disks' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'data_disks' to replace it",
			},
			"worker_period_unit": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{"Week", "Month"}, false),
				DiffSuppressFunc: csKubernetesWorkerPostPaidDiffSuppressFunc,
				Deprecated:       "Field 'worker_period_unit' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'period_unit' to replace it",
			},
			"worker_period": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: validation.Any(
					validation.IntBetween(1, 9),
					validation.IntInSlice([]int{12, 24, 36, 48, 60})),
				DiffSuppressFunc: csKubernetesWorkerPostPaidDiffSuppressFunc,
				Deprecated:       "Field 'worker_period' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'period' to replace it",
			},
			"worker_auto_renew": {
				Type:             schema.TypeBool,
				Optional:         true,
				DiffSuppressFunc: csKubernetesWorkerPostPaidDiffSuppressFunc,
				Deprecated:       "Field 'worker_auto_renew' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'auto_renew' to replace it",
			},
			"worker_auto_renew_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 6, 12}),
				DiffSuppressFunc: csKubernetesWorkerPostPaidDiffSuppressFunc,
				Deprecated:       "Field 'worker_auto_renew_period' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'auto_renew_period' to replace it",
			},
			"exclude_autoscaler_nodes": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Field 'worker_auto_renew_period' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes",
			},
			// global configurations
			"pod_vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				MaxItems: 10,
			},
			"pod_cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"node_cidr_mask": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      KubernetesClusterNodeCIDRMasksByDefault,
				ValidateFunc: validation.IntBetween(24, 28),
			},
			"enable_ssh": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Field 'enable_ssh' has been deprecated from provider version 1.177.0.",
			},
			"new_nat_gateway": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"password": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"key_name", "kms_encrypted_password"},
				Deprecated:    "Field 'password' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'password' to replace it",
			},
			"key_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"password", "kms_encrypted_password"},
				Deprecated:    "Field 'key_name' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'key_name' to replace it",
			},
			"kms_encrypted_password": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"password", "key_name"},
				Deprecated:    "Field 'kms_encrypted_password' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'kms_encrypted_password' to replace it",
			},
			"kms_encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_password").(string) == ""
				},
				Elem:       schema.TypeString,
				Deprecated: "Field 'kms_encryption_context' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'kms_encryption_context' to replace it",
			},
			"user_ca": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'image_id' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'image_id' to replace it",
			},
			"install_cloud_monitor": {
				Type:       schema.TypeBool,
				Optional:   true,
				Default:    true,
				Deprecated: "Field 'install_cloud_monitor' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'install_cloud_monitor' to replace it",
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// cpu policy options of kubelet
			"cpu_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"none", "static"}, false),
				Deprecated:   "Field 'cpu_policy' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'cpu_policy' to replace it",
			},
			"proxy_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ipvs",
				ValidateFunc: validation.StringInSlice([]string{"iptables", "ipvs"}, false),
			},
			"addons": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"config": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"disabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"slb_internet_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"load_balancer_spec": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"slb.s1.small", "slb.s2.small", "slb.s2.medium", "slb.s3.small", "slb.s3.medium", "slb.s3.large"}, false),
				Default:      "slb.s1.small",
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"enable_rrsa": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"os_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Linux",
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Windows", "Linux"}, false),
				Deprecated:   "Field 'os_type' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes.",
			},
			"platform": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				ForceNew:   true,
				Deprecated: "Field 'platform' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'platform' to replace it.",
			},
			"node_port_range": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Field 'platform' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes.",
			},
			"cluster_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "cluster.local",
				ForceNew:    true,
				Description: "cluster local domain ",
			},
			"runtime": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Deprecated: "Field 'runtime' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'runtime_name' and 'runtime_version' to replace it.",
			},
			"taints": {
				Type:     schema.TypeList,
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
						"effect": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"NoSchedule", "NoExecute", "PreferNoSchedule"}, false),
						},
					},
				},
				Deprecated: "Field 'taints' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'taints' to replace it.",
			},
			"rds_instances": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Deprecated: "Field 'rds_instances' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'rds_instances' to replace it.",
			},
			"user_data": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'user_data' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'user_data' to replace it.",
			},
			"node_name_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^customized,[a-z0-9]([-a-z0-9\.])*,([5-9]|[1][0-2]),([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`), "Each node name consists of a prefix, an IP substring, and a suffix. For example, if the node IP address is 192.168.0.55, the prefix is aliyun.com, IP substring length is 5, and the suffix is test, the node name will be aliyun.com00055test."),
				Deprecated:   "Field 'node_name_mode' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'node_name_mode' to replace it.",
			},
			"worker_nodes": {
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
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Deprecated: "Field 'worker_nodes' has been deprecated from provider version 1.177.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes.",
			},
			"custom_san": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"encryption_provider_key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "disk encryption key, only in ack-pro",
			},
			// computed parameters
			"kube_config": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_cert": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_ca_cert": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"certificate_authority": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_cert": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_cert": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"connections": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_server_internet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"api_server_intranet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"slb_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Field 'slb_id' has been deprecated from provider version 1.9.2. New field 'slb_internet' replaces it.",
			},
			"slb_internet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slb_intranet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_enterprise_security_group": {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"security_group_id"},
			},
			"nat_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// remove parameters below
			// mix vswitch_ids between master and worker is not a good guidance to create cluster
			"vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				MinItems: 3,
				MaxItems: 5,
				Removed:  "Field 'vswitch_ids' has been deprecated from provider version 1.75.0. New field 'master_vswitch_ids' and 'worker_vswitch_ids' replace it.",
			},
			// force update is a high risk operation
			"force_update": {
				Type:     schema.TypeBool,
				Optional: true,
				Removed:  "Field 'force_update' has been removed from provider version 1.75.0.",
			},
			// worker_numbers in array is a hell of management
			"worker_numbers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:    schema.TypeInt,
					Default: 3,
				},
				MinItems: 1,
				MaxItems: 3,
				Removed:  "Field 'worker_numbers' has been removed from provider version 1.75.0. New field 'worker_number' replaces it.",
			},
			"cluster_network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{KubernetesClusterNetworkTypeFlannel, KubernetesClusterNetworkTypeTerway}, false),
				Removed:      "Field 'cluster_network_type' has been removed from provider version 1.75.0. New field 'addons' replaces it.",
			},
			// too hard to use this config
			"log_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{KubernetesClusterLoggingTypeSLS}, false),
							Required:     true,
						},
						"project": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Removed: "Field 'log_config' has been removed from provider version 1.75.0. New field 'addons' replaces it.",
			},
			"worker_instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field 'worker_instance_type' has been removed from provider version 1.75.0. New field 'worker_instance_types' replaces it.",
			},
			"worker_ram_role_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_account_issuer": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"api_audiences": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ForceNew: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_spec": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"ack.standard", "ack.pro.small"}, false),
			},
			"maintenance_window": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"maintenance_time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"duration": {
							Type:     schema.TypeString,
							Required: true,
						},
						"weekly_period": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"control_plane_log_ttl": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"control_plane_log_project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"control_plane_log_components": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"retain_resources": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceAlicloudCSManagedKubernetesCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	invoker := NewInvoker()
	csService := CsService{client}
	args, err := buildKubernetesArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}
	var requestInfo *cs.Client
	var response interface{}
	if err := invoker.Run(func() error {
		raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			requestInfo = csClient
			args.RegionId = common.Region(client.RegionId)
			args.ClusterType = cs.ManagedKubernetes
			return csClient.CreateManagedKubernetesCluster(&cs.ManagedKubernetesClusterCreationRequest{
				ClusterArgs: args.ClusterArgs,
				WorkerArgs:  args.WorkerArgs,
			})
		})
		response = raw
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_managed_kubernetes", "CreateKubernetesCluster", response)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["RegionId"] = common.Region(client.RegionId)
		requestMap["Args"] = args
		addDebug("CreateKubernetesCluster", response, requestInfo, requestMap)
	}
	cluster, _ := response.(*cs.ClusterCommonResponse)
	d.SetId(cluster.ClusterID)

	stateConf := BuildStateConf([]string{"initial"}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudCSKubernetesRead(d, meta)
}

func UpgradeAlicloudKubernetesCluster(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if d.HasChange("version") {
		nextVersion := d.Get("version").(string)
		args := &cs.UpgradeClusterArgs{
			Version: nextVersion,
		}

		csService := CsService{client}
		err := csService.UpgradeCluster(d.Id(), args)

		if err != nil {
			return WrapError(err)
		}

		d.SetPartial("version")
	}
	return nil
}

func migrateAlicloudManagedKubernetesCluster(d *schema.ResourceData, meta interface{}) error {
	action := "MigrateCluster"
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}

	migrateClusterRequest := map[string]string{
		"type": "ManagedKubernetes",
		"spec": d.Get("cluster_spec").(string),
	}
	conn, err := meta.(*connectivity.AliyunClient).NewTeaRoaCommonClient(connectivity.OpenAckService)
	if err != nil {
		return WrapError(err)
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := conn.DoRequestWithAction(StringPointer(action), StringPointer("2015-12-15"), nil, StringPointer("POST"), StringPointer("AK"), String(fmt.Sprintf("/clusters/%s/migrate", d.Id())), nil, nil, migrateClusterRequest, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"QPS Limit Exceeded"}) || NeedRetry(err) {
				return resource.RetryableError(err)
			}
			addDebug(action, response, nil)
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, nil)
		return nil
	})

	stateConf := BuildStateConf([]string{"migrating"}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 20*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return err
	}

	d.SetPartial("cluster_spec")

	return nil
}

func updateKubernetesClusterTag(d *schema.ResourceData, meta interface{}) error {
	action := "ModifyClusterTags"
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}

	var modifyClusterTagsRequest []cs.Tag
	if tags, err := ConvertCsTags(d); err == nil {
		modifyClusterTagsRequest = tags
	}
	d.SetPartial("tags")
	conn, err := meta.(*connectivity.AliyunClient).NewTeaRoaCommonClient(connectivity.OpenAckService)
	if err != nil {
		return WrapError(err)
	}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := conn.DoRequestWithAction(StringPointer(action), StringPointer("2015-12-15"), nil, StringPointer("POST"), StringPointer("AK"), String(fmt.Sprintf("/clusters/%s/tags", d.Id())), nil, nil, modifyClusterTagsRequest, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"QPS Limit Exceeded"}) || NeedRetry(err) {
				return resource.RetryableError(err)
			}
			addDebug(action, response, nil)
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, nil)
		return nil
	})

	stateConf := BuildStateConf([]string{"updating"}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func updateKubernetesClusterRRSA(d *schema.ResourceData, meta interface{}, invoker *Invoker) error {
	enableRRSA := false
	if v, ok := d.GetOk("enable_rrsa"); ok {
		enableRRSA = v.(bool)
	}
	// it's not allowed to disable rrsa
	if !enableRRSA {
		return fmt.Errorf("It's not supported to disable RRSA! " +
			"If your cluster has enabled this function, please manually modify your tf file and add the rrsa configuration to the file.")
	}

	// version check
	clusterVersion := d.Get("version").(string)
	if res, err := versionCompare(KubernetesClusterRRSASupportedVersion, clusterVersion); res < 0 || err != nil {
		return fmt.Errorf("RRSA is not supported in current version: %s", clusterVersion)
	}

	action := "ModifyClusterRRSA"
	client := meta.(*connectivity.AliyunClient)

	var requestInfo cs.ModifyClusterArgs
	requestInfo.EnableRRSA = enableRRSA

	var response interface{}
	if err := invoker.Run(func() error {
		_, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return nil, csClient.ModifyCluster(d.Id(), &requestInfo)
		})
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, DenverdinoAliyungo)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["ClusterId"] = d.Id()
		requestMap["enable_rrsa"] = requestInfo.EnableRRSA
		addDebug(action, response, requestInfo, requestMap)
	}
	d.SetPartial("enable_rrsa")
	return nil
}

//versionCompare check version,
//if cueVersion is newer than neededVersion return 1
//if curVersion is equal neededVersion return 0
//if curVersion is older than neededVersion return -1
//example: neededVersion = 1.20.11-aliyun.1, curVersion = 1.22.3-aliyun.1, it will return 1
func versionCompare(neededVersion, curVersion string) (int, error) {
	if neededVersion == "" || curVersion == "" {
		if neededVersion == "" && curVersion == "" {
			return 0, nil
		} else {
			if neededVersion == "" {
				return 1, nil
			} else {
				return -1, nil
			}
		}
	}

	// 取出版本号
	regx := regexp.MustCompile(`[0-9]+\.[0-9]+\.[0-9]+`)
	neededVersion = regx.FindString(neededVersion)
	curVersion = regx.FindString(curVersion)

	currentVersions := strings.Split(neededVersion, ".")
	newVersions := strings.Split(curVersion, ".")

	compare := 0

	for index, val := range currentVersions {
		newVal := newVersions[index]
		v1, err1 := strconv.Atoi(val)
		v2, err2 := strconv.Atoi(newVal)

		if err1 != nil || err2 != nil {
			return -2, fmt.Errorf("NotSupport, current cluster version is not support: %s", curVersion)
		}

		if v1 > v2 {
			compare = -1
		} else if v1 == v2 {
			compare = 0
		} else {
			compare = 1
		}

		if compare != 0 {
			break
		}
	}

	return compare, nil
}
