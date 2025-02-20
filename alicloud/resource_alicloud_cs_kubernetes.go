package alicloud

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	roacs "github.com/alibabacloud-go/cs-20151215/v5/client"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"

	"strconv"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/cs"
	aliyungoecs "github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"gopkg.in/yaml.v2"
)

const (
	KubernetesClusterNetworkTypeFlannel = "flannel"
	KubernetesClusterNetworkTypeTerway  = "terway"

	KubernetesClusterLoggingTypeSLS = "SLS"

	KubernetesClusterRRSASupportedVersion = "1.22.3-aliyun.1"
)

var (
	KubernetesClusterNodeCIDRMasksByDefault = 24
)

func resourceAlicloudCSKubernetes() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSKubernetesCreate,
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
				ValidateFunc:  StringLenBetween(1, 63),
				ConflictsWith: []string{"name_prefix"},
			},
			"name_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "Terraform-Creation",
				ValidateFunc:  StringLenBetween(0, 37),
				ConflictsWith: []string{"name"},
				Deprecated:    "Field 'name_prefix' has been deprecated from provider version 1.75.0.",
			},
			// master configurations
			"master_vswitch_ids": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				MinItems: 3,
				MaxItems: 5,
			},
			"master_instance_types": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems: 3,
				MaxItems: 5,
			},
			"master_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      40,
				ValidateFunc: IntBetween(40, 500),
			},
			"master_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  DiskCloudEfficiency,
			},
			"master_disk_performance_level": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateFunc:     StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
				DiffSuppressFunc: masterDiskPerformanceLevelDiffSuppressFunc,
			},
			"master_disk_snapshot_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"master_instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
				Default:      PostPaid,
			},
			"master_period_unit": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          Month,
				ValidateFunc:     StringInSlice([]string{"Week", "Month"}, false),
				DiffSuppressFunc: csKubernetesMasterPostPaidDiffSuppressFunc,
			},
			"master_period": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  1,
				// must be a valid period, expected [1-9], 12, 24, 36, 48 or 60,
				ValidateFunc: validation.Any(
					IntBetween(1, 9),
					IntInSlice([]int{12, 24, 36, 48, 60})),
				DiffSuppressFunc: csKubernetesMasterPostPaidDiffSuppressFunc,
			},
			"master_auto_renew": {
				Type:             schema.TypeBool,
				Default:          false,
				ForceNew:         true,
				Optional:         true,
				DiffSuppressFunc: csKubernetesMasterPostPaidDiffSuppressFunc,
			},
			"master_auto_renew_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				Default:          1,
				ValidateFunc:     IntInSlice([]int{1, 2, 3, 6, 12}),
				DiffSuppressFunc: csKubernetesMasterPostPaidDiffSuppressFunc,
			},
			// worker configurations
			"worker_vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				Removed:  "Field 'worker_vswitch_ids' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster worker nodes, by using field 'vswitch_ids' to replace it",
				MinItems: 1,
			},
			"worker_instance_types": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems: 1,
				MaxItems: 10,
				Removed:  "Field 'worker_instance_types' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster worker nodes, by using field 'instance_types' to replace it",
			},
			"worker_number": {
				Type:     schema.TypeInt,
				Optional: true,
				Removed:  "Field 'worker_number' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster worker nodes, by using field 'desired_size' to replace it",
			},
			"worker_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(20, 32768),
				Removed:      "Field 'worker_disk_size' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster worker nodes, by using field 'system_disk_size' to replace it",
			},
			"worker_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field 'worker_disk_category' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster worker nodes, by using field 'system_disk_category' to replace it",
			},
			"worker_disk_performance_level": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
				DiffSuppressFunc: workerDiskPerformanceLevelDiffSuppressFunc,
				Removed:          "Field 'worker_disk_performance_level' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster worker nodes, by using field 'system_disk_performance_level' to replace it",
			},
			"worker_disk_snapshot_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field 'worker_disk_snapshot_policy_id' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster worker nodes, by using field 'system_disk_snapshot_policy_id' to replace it",
			},
			"worker_data_disk_size": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     IntBetween(20, 32768),
				DiffSuppressFunc: workerDataDiskSizeSuppressFunc,
				Removed:          "Field 'worker_data_disk_size' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster worker nodes, by using field 'data_disks.size' to replace it",
			},
			"worker_data_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field 'worker_data_disk_category' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster worker nodes, by using field 'data_disks.category' to replace it",
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
							ValidateFunc: StringInSlice([]string{"all", "cloud", "ephemeral_ssd", "cloud_essd", "cloud_efficiency", "cloud_ssd", "local_disk"}, false),
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
				Removed: "Field 'worker_data_disks' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster worker nodes, by using field 'data_disks' to replace it",
			},
			"worker_instance_charge_type": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				Removed:  "Field 'worker_instance_charge_type' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster worker nodes, by using field 'instance_charge_type' to replace it",
			},
			"worker_period_unit": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     StringInSlice([]string{"Week", "Month"}, false),
				DiffSuppressFunc: csKubernetesWorkerPostPaidDiffSuppressFunc,
				Removed:          "Field 'worker_period_unit' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster worker nodes, by using field 'period_unit' to replace it",
			},
			"worker_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.Any(
					IntBetween(1, 9),
					IntInSlice([]int{12, 24, 36, 48, 60})),
				DiffSuppressFunc: csKubernetesWorkerPostPaidDiffSuppressFunc,
				Removed:          "Field 'worker_period' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster worker nodes, by using field 'period' to replace it",
			},
			"worker_auto_renew": {
				Type:             schema.TypeBool,
				Optional:         true,
				DiffSuppressFunc: csKubernetesWorkerPostPaidDiffSuppressFunc,
				Removed:          "Field 'worker_auto_renew' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster worker nodes, by using field 'auto_renew' to replace it",
			},
			"worker_auto_renew_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     IntInSlice([]int{1, 2, 3, 6, 12}),
				DiffSuppressFunc: csKubernetesWorkerPostPaidDiffSuppressFunc,
				Removed:          "Field 'worker_auto_renew_period' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster worker nodes, by using field 'auto_renew_period' to replace it",
			},
			"exclude_autoscaler_nodes": {
				Type:     schema.TypeBool,
				Optional: true,
				Removed:  "Field 'exclude_autoscaler_nodes' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster worker nodes.",
			},
			// global configurations
			// Terway network
			"pod_vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
			},
			// Flannel network
			"pod_cidr": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"service_cidr": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"node_cidr_mask": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      KubernetesClusterNodeCIDRMasksByDefault,
				ValidateFunc: IntBetween(24, 28),
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
			},
			"key_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"password", "kms_encrypted_password"},
				ForceNew:      true,
			},
			"kms_encrypted_password": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"password", "key_name"},
			},
			"kms_encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_password").(string) == ""
				},
				Elem: schema.TypeString,
			},
			"user_ca": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_ssh": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"load_balancer_spec": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"slb.s1.small", "slb.s2.small", "slb.s2.medium", "slb.s3.small", "slb.s3.medium", "slb.s3.large"}, false),
				Computed:     true,
				Deprecated:   "Field 'load_balancer_spec' has been deprecated from provider version 1.232.0. The spec will not take effect because the charge of the load balancer has been changed to PayByCLCU",
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"install_cloud_monitor": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
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
				ValidateFunc: StringInSlice([]string{"none", "static"}, false),
				Removed:      "Field 'cpu_policy' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster worker nodes, by using field 'cpu_policy' to replace it",
			},
			"proxy_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"iptables", "ipvs"}, false),
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
						"version": {
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
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"slb_internet_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
				ValidateFunc: StringInSlice([]string{"Windows", "Linux"}, false),
			},
			"platform": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"node_port_range": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Removed:  "Field 'node_port_range' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster worker nodes.",
			},
			"runtime": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "docker",
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "19.03.5",
						},
					},
				},
			},
			"cluster_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "cluster.local",
				ForceNew:    true,
				Description: "cluster local domain",
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
							ValidateFunc: StringInSlice([]string{"NoSchedule", "NoExecute", "PreferNoSchedule"}, false),
						},
					},
				},
				Removed: "Field 'taints' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster worker nodes, by using field 'taints' to replace it",
			},
			"rds_instances": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"custom_san": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// computed parameters
			"kube_config": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field 'kube_config' has been removed from provider version 1.212.0. New DataSource 'alicloud_cs_cluster_credential' manage your cluster's kube config.",
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
				Type:     schema.TypeString,
				Computed: true,
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
				ForceNew: true,
			},
			"is_enterprise_security_group": {
				Type:          schema.TypeBool,
				Optional:      true,
				ForceNew:      true,
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
			"master_nodes": {
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
			},
			"worker_nodes": {
				Type:     schema.TypeList,
				Optional: true,
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
				Removed: "Field 'worker_nodes' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster worker nodes.",
			},
			// remove parameters below
			// mix vswitch_ids between master and worker is not a good guidance to create cluster
			"worker_instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field 'worker_instance_type' has been removed from provider version 1.75.0. New field 'worker_instance_types' replaces it.",
			},
			"vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				MinItems: 3,
				MaxItems: 5,
				Removed:  "Field 'vswitch_ids' has been removed from provider version 1.75.0. New field 'master_vswitch_ids' and 'worker_vswitch_ids' replace it.",
			},
			// single instance type would cause extra troubles
			"master_instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field 'master_instance_type' has been removed from provider version 1.75.0. New field 'master_instance_types' replaces it.",
			},
			// force update is a high risk operation
			"force_update": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Removed:  "Field 'force_update' has been removed from provider version 1.75.0.",
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Removed:  "Field 'availability_zone' has been removed from provider version 1.212.0.",
			},
			// single az would be never supported.
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field 'vswitch_id' has been removed from provider version 1.75.0. New field 'master_vswitch_ids' and 'worker_vswitch_ids' replaces it.",
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
			"nodes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Removed:  "Field 'nodes' has been removed from provider version 1.9.4. New field 'master_nodes' replaces it.",
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
							ValidateFunc: StringInSlice([]string{KubernetesClusterLoggingTypeSLS}, false),
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
			"cluster_network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{KubernetesClusterNetworkTypeFlannel, KubernetesClusterNetworkTypeTerway}, false),
				Removed:      "Field 'cluster_network_type' has been removed from provider version 1.75.0. New field 'addons' replaces it.",
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field 'user_data' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster worker nodes, by using field 'user_data' to replace it",
			},
			"node_name_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringMatch(regexp.MustCompile(`^customized,[a-z0-9]([-a-z0-9\.])*,([5-9]|[1][0-2]),([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`), "Each node name consists of a prefix, an IP substring, and a suffix. For example, if the node IP address is 192.168.0.55, the prefix is aliyun.com, IP substring length is 5, and the suffix is test, the node name will be aliyun.com00055test."),
				ForceNew:     true,
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
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"retain_resources": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"delete_options": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"SLB", "ALB", "SLS_Data", "SLS_ControlPlane", "PrivateZone"}, false),
						},
						"delete_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"delete", "retain"}, false),
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudCSKubernetesCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	invoker := NewInvoker()

	var requestInfo *cs.Client
	var raw interface{}

	// prepare args and set default value
	args, err := buildKubernetesArgs(d, meta)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_kubernetes", "PrepareKubernetesClusterArgs", err)
	}

	if err = invoker.Run(func() error {
		raw, err = client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			args.RegionId = common.Region(client.RegionId)
			args.ClusterType = cs.DelicatedKubernetes
			return csClient.CreateDelicatedKubernetesCluster(args)
		})
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_kubernetes", "CreateKubernetesCluster", raw)
	}

	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["RegionId"] = common.Region(client.RegionId)
		requestMap["Params"] = args
		addDebug("CreateKubernetesCluster", raw, requestInfo, requestMap)
	}

	cluster, ok := raw.(*cs.ClusterCommonResponse)
	if ok != true {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_kubernetes", "ParseKubernetesClusterResponse", raw)
	}

	d.SetId(cluster.ClusterID)
	taskId := cluster.TaskId
	roaCsClient, err := client.NewRoaCsClient()
	if err == nil {
		csClient := CsClient{client: roaCsClient}
		stateConf := BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, csClient.DescribeTaskRefreshFunc(d, taskId, []string{"fail", "failed"}))
		if jobDetail, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, ResponseCodeMsg, d.Id(), "createCluster", jobDetail)
		}
	}

	// reset interval to 10s
	stateConf := BuildStateConf([]string{"initial"}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudCSKubernetesRead(d, meta)
}

func resourceAlicloudCSKubernetesUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	invoker := NewInvoker()
	// modifyCluster
	if !d.IsNewResource() && d.HasChanges("resource_group_id", "name", "name_prefix", "deletion_protection", "custom_san", "maintenance_window", "operation_policy", "enable_rrsa") {
		if err := modifyCluster(d, meta, &invoker); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifyCluster", AlibabaCloudSdkGoERROR)
		}
	}

	// modify cluster tag
	if d.HasChange("tags") {
		err := updateKubernetesClusterTag(d, meta)
		if err != nil {
			return WrapErrorf(err, ResponseCodeMsg, d.Id(), "ModifyClusterTags", AlibabaCloudSdkGoERROR)
		}
	}

	// update control plane config
	if d.HasChanges([]string{"control_plane_log_ttl", "control_plane_log_project", "control_plane_log_components"}...) {
		if err := updateControlPlaneLog(d, meta); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateControlPlaneLog", AlibabaCloudSdkGoERROR)
		}
	}

	// migrate cluster to pro from standard
	if d.HasChange("cluster_spec") {
		err := migrateCluster(d, meta)
		if err != nil {
			return WrapError(err)
		}
	}

	err := UpgradeAlicloudKubernetesCluster(d, meta)
	if err != nil {
		return WrapError(err)
	}

	d.Partial(false)
	return resourceAlicloudCSKubernetesRead(d, meta)
}

func migrateCluster(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	oldValue, newValue := d.GetChange("cluster_spec")
	o, ok := oldValue.(string)
	if ok != true {
		return WrapErrorf(fmt.Errorf("cluster_spec old value can not be parsed"), "parseError %d", oldValue)
	}
	n, ok := newValue.(string)
	if ok != true {
		return WrapErrorf(fmt.Errorf("cluster_pec new value can not be parsed"), "parseError %d", newValue)
	}

	// The field `cluster_spec` of some ack.standard managed cluster is "" since some historical reasons.
	// The logic here is to confirm whether the cluster is the above.
	// The interface error should not block the main process and errors will be output to the log.
	clusterInfo, err := csService.DescribeCsManagedKubernetes(d.Id())
	if err != nil || clusterInfo == nil {
		log.Printf("[DEBUG] Failed to DescribeCsManagedKubernetes cluster %s when migrate", d.Id())
	} else {
		o = clusterInfo.ClusterSpec
	}

	if (o == "ack.standard" || o == "") && strings.Contains(n, "pro") {
		err := migrateAlicloudManagedKubernetesCluster(d, meta)
		if err != nil {
			return WrapErrorf(err, ResponseCodeMsg, d.Id(), "MigrateCluster", AlibabaCloudSdkGoERROR)
		}
	}

	return nil
}

func modifyCluster(d *schema.ResourceData, meta interface{}, invoker *Invoker) error {
	updated := false
	request := &roacs.ModifyClusterRequest{}
	client := meta.(*connectivity.AliyunClient)
	csClient, err := client.NewRoaCsClient()
	csService := CsService{client}

	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		request.SetResourceGroupId(d.Get("resource_group_id").(string))
		updated = true
	}

	if !d.IsNewResource() && d.HasChange("name") {
		var clusterName string
		if v, ok := d.GetOk("name"); ok {
			clusterName = v.(string)
			request.SetClusterName(clusterName)
			updated = true
		}
	}

	// modify cluster deletion protection
	if !d.IsNewResource() && d.HasChange("deletion_protection") {
		v := d.Get("deletion_protection")
		request.SetDeletionProtection(v.(bool))
		updated = true
	}

	// modify cluster maintenance window
	if !d.IsNewResource() && d.HasChange("maintenance_window") {
		if v := d.Get("maintenance_window").([]interface{}); len(v) > 0 {
			request.MaintenanceWindow = expandMaintenanceWindowConfigRoa(v)
			updated = true
		}
		d.SetPartial("maintenance_window")
	}

	// modify cluster maintenance window
	if !d.IsNewResource() && d.HasChange("operation_policy") {
		if v := d.Get("operation_policy").([]interface{}); len(v) > 0 {
			request.OperationPolicy = &roacs.ModifyClusterRequestOperationPolicy{}
			if vv := d.Get("operation_policy.0.cluster_auto_upgrade").([]interface{}); len(vv) > 0 {
				policy := vv[0].(map[string]interface{})
				request.OperationPolicy.ClusterAutoUpgrade = &roacs.ModifyClusterRequestOperationPolicyClusterAutoUpgrade{
					Enabled: tea.Bool(policy["enabled"].(bool)),
					Channel: tea.String(policy["channel"].(string)),
				}
			}
			updated = true
		}
	}

	// modify cluster rrsa policy
	if d.HasChange("enable_rrsa") {
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
		request.SetEnableRrsa(enableRRSA)
		updated = true
		d.SetPartial("enable_rrsa")
	}

	if d.HasChange("custom_san") {
		customSan := d.Get("custom_san").(string)
		request.SetApiServerCustomCertSans(
			&roacs.ModifyClusterRequestApiServerCustomCertSans{
				SubjectAlternativeNames: tea.StringSlice(strings.Split(customSan, ",")),
				Action:                  tea.String("overwrite"),
			},
		)
		updated = true
	}

	if d.HasChange("vswitch_ids") {
		vSwitchIds := expandStringList(d.Get("vswitch_ids").([]interface{}))
		request.SetVswitchIds(tea.StringSlice(vSwitchIds))
		updated = true
	}

	if updated == false {
		return nil
	}

	if err := invoker.Run(func() error {
		_, err := csClient.ModifyCluster(tea.String(d.Id()), request)
		return err
	}); err != nil && !IsExpectedErrors(err, []string{"ClusterNameAlreadyExist", "ErrorModifyDeletionProtectionFailed"}) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifyCluster", AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"updating"}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func resourceAlicloudCSKubernetesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}

	object, err := csService.DescribeCsKubernetes(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cs_kubernetes DescribeCsKubernetes Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.Name)
	d.Set("vpc_id", object.VpcId)
	d.Set("security_group_id", object.SecurityGroupId)
	d.Set("version", object.CurrentVersion)
	d.Set("worker_ram_role_name", object.WorkerRamRoleName)
	d.Set("resource_group_id", object.ResourceGroupId)
	d.Set("cluster_spec", object.ClusterSpec)
	d.Set("deletion_protection", object.DeletionProtection)

	slbId, err := getApiServerSlbID(d, meta)
	if err != nil {
		log.Printf(DefaultErrorMsg, d.Id(), "DescribeClusterResources", err.Error())
	}
	d.Set("slb_id", slbId)

	// compat for default value
	if spec := d.Get("load_balancer_spec").(string); spec != "" {
		d.Set("load_balancer_spec", spec)
	}

	if err := d.Set("tags", flattenTagsConfig(object.Tags)); err != nil {
		return WrapError(err)
	}
	if d.Get("os_type") == "" {
		d.Set("os_type", "Linux")
	}

	if d.Get("cluster_domain") == "" {
		d.Set("cluster_domain", "cluster.local")
	}

	d.Set("maintenance_window", flattenMaintenanceWindowConfig(&object.MaintenanceWindow))

	//request.Parameters
	if v, ok := object.Parameters["MasterVSwitchIds"]; ok {
		d.Set("master_vswitch_ids", strings.Split(Interface2String(v), ","))
	}
	if v, ok := object.Parameters["MasterSystemDiskCategory"]; ok {
		category := Interface2String(v)
		d.Set("master_disk_category", category)
		if category == string(DiskCloudESSD) {
			if v, ok := object.Parameters["MasterSystemDiskPerformanceLevel"]; ok && v != nil {
				d.Set("master_disk_performance_level", Interface2String(v))
			}
		}
	}
	if v, ok := object.Parameters["MasterSystemDiskSize"]; ok {
		d.Set("master_disk_size", formatInt(v))
	}
	if v, ok := object.Parameters["MasterSnapshotPolicyId"]; ok {
		d.Set("master_disk_snapshot_policy_id", Interface2String(v))
	}
	if v, ok := object.Parameters["MasterInstanceChargeType"]; ok {
		chargeType := Interface2String(v)
		d.Set("master_instance_charge_type", chargeType)
		if chargeType == string(PrePaid) {
			if v, ok := object.Parameters["MasterAutoRenew"]; ok {
				d.Set("master_auto_renew", Interface2Bool(v))
			}
			if v, ok := object.Parameters["MasterAutoRenewPeriod"]; ok {
				d.Set("master_auto_renew_period", formatInt(v))
			}
			if v, ok := object.Parameters["MasterPeriod"]; ok {
				d.Set("master_period", formatInt(v))
			}
			if v, ok := object.Parameters["MasterPeriodUnit"]; ok {
				d.Set("master_period_unit", Interface2String(v))
			}
		}
	}
	if v, ok := object.Parameters["MasterInstanceTypes"]; ok {
		d.Set("master_instance_types", strings.Split(Interface2String(v), ","))
	}
	if object.ClusterType != "Kubernetes" {
		if v, ok := object.Parameters["WorkerVSwitchIds"]; ok {
			d.Set("worker_vswitch_ids", strings.Split(Interface2String(v), ","))
		}
	}
	if object.Profile == EdgeProfile {
		if v, ok := object.Parameters["WorkerInstanceChargeType"]; ok {
			d.Set("worker_instance_charge_type", Interface2String(v))
		}
		if v, ok := object.Parameters["WorkerInstanceTypes"]; ok {
			d.Set("worker_instance_types", strings.Split(Interface2String(v), ","))
		}
		if v, ok := object.Parameters["WorkerSystemDiskCategory"]; ok {
			d.Set("worker_disk_category", Interface2String(v))
		}
		if v, ok := object.Parameters["WorkerSystemDiskSize"]; ok {
			d.Set("worker_disk_size", formatInt(v))
		}
		if v, ok := object.Parameters["CloudMonitorFlags"]; ok {
			d.Set("install_cloud_monitor", Interface2Bool(v))
		}
		// only works with default-nodepool
		workerNodes := fetchWorkerNodes(d, meta)
		d.Set("worker_nodes", workerNodes)
		d.Set("worker_number", len(workerNodes))
	}
	if object.ClusterType == "Kubernetes" {
		if v, ok := object.Parameters["CloudMonitorFlags"]; ok {
			d.Set("install_cloud_monitor", Interface2Bool(v))
		}
		if v, ok := object.Parameters["SSHFlags"]; ok {
			d.Set("enable_ssh", Interface2Bool(v))
		}
		if v, ok := object.Parameters["NodeNameMode"]; ok {
			d.Set("node_name_mode", Interface2String(v))
		}
		if v, ok := object.Parameters["MasterImageId"]; ok {
			d.Set("image_id", Interface2String(v))
		}
		if v, ok := object.Parameters["MasterKeyPair"]; ok {
			d.Set("key_name", Interface2String(v))
		}
		if v, ok := object.Parameters["IsEnterpriseSecurityGroup"]; ok {
			d.Set("is_enterprise_security_group", Interface2Bool(v))
		}
		d.Set("master_nodes", fetchMasterNodes(d, meta))
	}
	if v, ok := object.Parameters["PodVswitchIds"]; ok {
		l := make([]string, 0)
		err := json.Unmarshal([]byte(Interface2String(v)), &l)
		if err == nil && len(l) > 0 {
			d.Set("pod_vswitch_ids", l)
		}
	}
	if v, ok := object.Parameters["ProxyMode"]; ok {
		d.Set("proxy_mode", Interface2String(v))
	}
	if v, ok := object.Parameters["ServiceCIDR"]; ok {
		d.Set("service_cidr", Interface2String(v))
	}
	if v, ok := object.Parameters["ContainerCIDR"]; ok {
		d.Set("pod_cidr", Interface2String(v))
	}
	//if v, ok := object.Parameters["SNatEntry"]; ok {
	//	d.Set("new_nat_gateway", Interface2String(v))
	//}

	// Cluster Metadata
	metadata := object.GetMetaData()
	if v, ok := metadata["ExtraCertSAN"]; ok && v != nil {
		l := expandStringList(v.([]interface{}))
		d.Set("custom_san", strings.Join(l, ","))
	}
	//if v, ok := metadata["Timezone"]; ok {
	//	d.Set("timezone", Interface2String(v))
	//}
	// Cluster capabilities
	capabilities := fetchClusterCapabilities(object.MetaData)
	if v, ok := capabilities["PublicSLB"]; ok {
		d.Set("slb_internet_enabled", Interface2Bool(v))
	}
	if v, ok := capabilities["NodeCIDRMask"]; ok {
		d.Set("node_cidr_mask", formatInt(v))
	}

	// Get slb information and set connect
	connection := make(map[string]string)
	masterURL := object.MasterURL
	endPoint := make(map[string]string)
	_ = json.Unmarshal([]byte(masterURL), &endPoint)
	connection["api_server_internet"] = endPoint["api_server_endpoint"]
	connection["api_server_intranet"] = endPoint["intranet_api_server_endpoint"]
	if endPoint["api_server_endpoint"] != "" {
		connection["master_public_ip"] = strings.Split(strings.Split(endPoint["api_server_endpoint"], ":")[1], "/")[2]
	}
	if object.Profile != EdgeProfile {
		connection["service_domain"] = fmt.Sprintf("*.%s.%s.alicontainer.com", d.Id(), object.RegionId)
	}

	d.Set("connections", connection)
	d.Set("slb_internet", connection["master_public_ip"])
	if endPoint["intranet_api_server_endpoint"] != "" {
		d.Set("slb_intranet", strings.Split(strings.Split(endPoint["intranet_api_server_endpoint"], ":")[1], "/")[2])
	}

	// set nat gateway
	natRequest := vpc.CreateDescribeNatGatewaysRequest()
	natRequest.VpcId = object.VpcId
	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeNatGateways(natRequest)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), natRequest.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(natRequest.GetActionName(), raw, natRequest.RpcRequest, natRequest)
	nat, _ := raw.(*vpc.DescribeNatGatewaysResponse)
	if nat != nil && len(nat.NatGateways.NatGateway) > 0 {
		d.Set("nat_gateway_id", nat.NatGateways.NatGateway[0].NatGatewayId)
	}

	// get cluster conn certs
	// If the cluster is failed, there is no need to get cluster certs
	if object.State == "failed" || object.State == "delete_failed" || object.State == "deleting" {
		return nil
	}

	if err = setCerts(d, meta, true); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAlicloudCSKubernetesDelete(d *schema.ResourceData, meta interface{}) error {
	csService := CsService{meta.(*connectivity.AliyunClient)}
	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "InitializeClient", err)
	}

	args := &roacs.DeleteClusterRequest{}
	if v := d.Get("retain_resources"); len(v.([]interface{})) > 0 {
		args.RetainResources = tea.StringSlice(expandStringList(v.([]interface{})))
	}
	if v, ok := d.GetOk("delete_options"); ok && len(v.([]interface{})) > 0 {
		for _, vv := range v.([]interface{}) {
			if options, ok := vv.(map[string]interface{}); ok {
				args.DeleteOptions = append(args.DeleteOptions, &roacs.DeleteClusterRequestDeleteOptions{
					DeleteMode:   tea.String(options["delete_mode"].(string)),
					ResourceType: tea.String(options["resource_type"].(string)),
				})
			}
		}
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	var resp *roacs.DeleteClusterResponse
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = client.DeleteCluster(tea.String(d.Id()), args)
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
		if IsExpectedErrors(err, []string{"ErrorClusterNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteCluster", AliyunTablestoreGoSdk)
	}

	taskId := tea.StringValue(resp.Body.TaskId)
	csClient := CsClient{client: client}
	stateConf := BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, csClient.DescribeTaskRefreshFunc(d, taskId, []string{"fail", "failed"}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, ResponseCodeMsg, d.Id(), "deleteCluster", jobDetail)
	}

	stateConf = BuildStateConf([]string{"running", "deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func buildKubernetesArgs(d *schema.ResourceData, meta interface{}) (*cs.DelicatedKubernetesClusterCreationRequest, error) {
	client := meta.(*connectivity.AliyunClient)

	vpcService := VpcService{client}

	var vswitchID string
	list := make([]string, 0)
	if v, ok := d.GetOk("master_vswitch_ids"); ok {
		list = append(list, expandStringList(v.([]interface{}))...)
	}
	if v, ok := d.GetOk("worker_vswitch_ids"); ok {
		list = append(list, expandStringList(v.([]interface{}))...)
	}
	if len(list) > 0 {
		vswitchID = list[0]
	} else {
		vswitchID = ""
	}

	var vpcId string
	if vswitchID != "" {
		vsw, err := vpcService.DescribeVSwitch(vswitchID)
		if err != nil {
			return nil, err
		}
		vpcId = vsw.VpcId
	}

	var clusterName string
	if v, ok := d.GetOk("name"); ok {
		clusterName = v.(string)
	} else {
		clusterName = resource.PrefixedUniqueId(d.Get("name_prefix").(string))
	}

	addons := make([]cs.Addon, 0)
	if v, ok := d.GetOk("addons"); ok {
		all, ok := v.([]interface{})
		if ok {
			for _, a := range all {
				addon, ok := a.(map[string]interface{})
				if ok {
					addons = append(addons, cs.Addon{
						Name:     addon["name"].(string),
						Config:   addon["config"].(string),
						Version:  addon["version"].(string),
						Disabled: addon["disabled"].(bool),
					})
				}
			}
		}
	}

	var apiAudiences string
	if d.Get("api_audiences") != nil {
		if list := expandStringList(d.Get("api_audiences").([]interface{})); len(list) > 0 {
			apiAudiences = strings.Join(list, ",")
		}
	}

	creationArgs := &cs.DelicatedKubernetesClusterCreationRequest{
		ClusterArgs: cs.ClusterArgs{
			DisableRollback:    true,
			Name:               clusterName,
			DeletionProtection: d.Get("deletion_protection").(bool),
			VpcId:              vpcId,
			// the params below is ok to be empty
			KubernetesVersion:         d.Get("version").(string),
			NodeCidrMask:              strconv.Itoa(d.Get("node_cidr_mask").(int)),
			KeyPair:                   d.Get("key_name").(string),
			ServiceCidr:               d.Get("service_cidr").(string),
			CloudMonitorFlags:         d.Get("install_cloud_monitor").(bool),
			SecurityGroupId:           d.Get("security_group_id").(string),
			IsEnterpriseSecurityGroup: d.Get("is_enterprise_security_group").(bool),
			EndpointPublicAccess:      d.Get("slb_internet_enabled").(bool),
			SnatEntry:                 d.Get("new_nat_gateway").(bool),
			Addons:                    addons,
			ApiAudiences:              apiAudiences,
		},
	}

	if enableRRSA, ok := d.GetOk("enable_rrsa"); ok {
		creationArgs.EnableRRSA = enableRRSA.(bool)
	}

	if lbSpec, ok := d.GetOk("load_balancer_spec"); ok {
		creationArgs.LoadBalancerSpec = lbSpec.(string)
	}

	if osType, ok := d.GetOk("os_type"); ok {
		creationArgs.OsType = osType.(string)
	}

	if platform, ok := d.GetOk("platform"); ok {
		creationArgs.Platform = platform.(string)
	}

	if timezone, ok := d.GetOk("timezone"); ok {
		creationArgs.Timezone = timezone.(string)
	}

	if clusterDomain, ok := d.GetOk("cluster_domain"); ok {
		creationArgs.ClusterDomain = clusterDomain.(string)
	}

	if customSan, ok := d.GetOk("custom_san"); ok {
		creationArgs.CustomSAN = customSan.(string)
	}

	if imageId, ok := d.GetOk("image_id"); ok {
		creationArgs.ClusterArgs.ImageId = imageId.(string)
	}
	if nodeNameMode, ok := d.GetOk("node_name_mode"); ok {
		creationArgs.ClusterArgs.NodeNameMode = nodeNameMode.(string)
	}
	if saIssuer, ok := d.GetOk("service_account_issuer"); ok {
		creationArgs.ClusterArgs.ServiceAccountIssuer = saIssuer.(string)
	}
	if resourceGroupId, ok := d.GetOk("resource_group_id"); ok {
		creationArgs.ClusterArgs.ResourceGroupId = resourceGroupId.(string)
	}

	if v := d.Get("user_data").(string); v != "" {
		_, base64DecodeError := base64.StdEncoding.DecodeString(v)
		if base64DecodeError == nil {
			creationArgs.UserData = v
		} else {
			creationArgs.UserData = base64.StdEncoding.EncodeToString([]byte(v))
		}
	}

	if _, ok := d.GetOk("pod_vswitch_ids"); ok {
		creationArgs.PodVswitchIds = expandStringList(d.Get("pod_vswitch_ids").([]interface{}))
	} else {
		creationArgs.ContainerCidr = d.Get("pod_cidr").(string)
	}

	if password := d.Get("password").(string); password == "" {
		if v, ok := d.GetOk("kms_encrypted_password"); ok && v != "" {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(v.(string), d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return nil, WrapError(err)
			}
			password = decryptResp
		}
		creationArgs.LoginPassword = password
	} else {
		creationArgs.LoginPassword = password
	}

	if tags, err := ConvertCsTags(d); err == nil {
		creationArgs.Tags = tags
	}
	// CA default is empty
	if userCa, ok := d.GetOk("user_ca"); ok {
		userCaContent, err := loadFileContent(userCa.(string))
		if err != nil {
			return nil, fmt.Errorf("reading user_ca file failed %s", err)
		}
		creationArgs.UserCa = string(userCaContent)
	}

	// set proxy mode and default is ipvs
	if proxyMode := d.Get("proxy_mode").(string); proxyMode != "" {
		creationArgs.ProxyMode = cs.ProxyMode(proxyMode)
	} else {
		creationArgs.ProxyMode = cs.ProxyMode(cs.IPVS)
	}

	// dedicated kubernetes must provide master_vswitch_ids
	if _, ok := d.GetOk("master_vswitch_ids"); ok {
		creationArgs.MasterArgs = cs.MasterArgs{
			MasterCount:              len(d.Get("master_vswitch_ids").([]interface{})),
			MasterVSwitchIds:         expandStringList(d.Get("master_vswitch_ids").([]interface{})),
			MasterInstanceTypes:      expandStringList(d.Get("master_instance_types").([]interface{})),
			MasterSystemDiskCategory: aliyungoecs.DiskCategory(d.Get("master_disk_category").(string)),
			MasterSystemDiskSize:     int64(d.Get("master_disk_size").(int)),
			// TODO support other params
		}
	}

	if v, ok := d.GetOk("master_disk_snapshot_policy_id"); ok && v != "" {
		creationArgs.MasterArgs.MasterSnapshotPolicyId = v.(string)
	}

	if v, ok := d.GetOk("master_disk_performance_level"); ok && v != "" {
		creationArgs.MasterArgs.MasterSystemDiskPerformanceLevel = v.(string)
	}

	if v, ok := d.GetOk("master_instance_charge_type"); ok {
		creationArgs.MasterInstanceChargeType = v.(string)
		if creationArgs.MasterInstanceChargeType == string(PrePaid) {
			creationArgs.MasterAutoRenew = d.Get("master_auto_renew").(bool)
			creationArgs.MasterAutoRenewPeriod = d.Get("master_auto_renew_period").(int)
			creationArgs.MasterPeriod = d.Get("master_period").(int)
			creationArgs.MasterPeriodUnit = d.Get("master_period_unit").(string)
		}
	}

	var workerDiskSize int64
	if d.Get("worker_disk_size") != nil {
		workerDiskSize = int64(d.Get("worker_disk_size").(int))
	}

	if v, ok := d.GetOk("worker_vswitch_ids"); ok {
		creationArgs.WorkerArgs.WorkerVSwitchIds = expandStringList(v.([]interface{}))
	}
	if v, ok := d.GetOk("worker_instance_types"); ok {
		creationArgs.WorkerArgs.WorkerInstanceTypes = expandStringList(v.([]interface{}))
	}
	if v, ok := d.GetOk("worker_number"); ok {
		creationArgs.WorkerArgs.NumOfNodes = int64(v.(int))
	}
	if v, ok := d.GetOk("worker_disk_category"); ok {
		creationArgs.WorkerArgs.WorkerSystemDiskCategory = aliyungoecs.DiskCategory(v.(string))
	}
	if v, ok := d.GetOk("worker_disk_snapshot_policy_id"); ok && v != "" {
		creationArgs.WorkerArgs.WorkerSnapshotPolicyId = v.(string)
	}
	if v, ok := d.GetOk("worker_disk_performance_level"); ok && v != "" {
		creationArgs.WorkerArgs.WorkerSystemDiskPerformanceLevel = v.(string)
	}

	if dds, ok := d.GetOk("worker_data_disks"); ok {
		disks := dds.([]interface{})
		createDataDisks := make([]cs.DataDisk, 0, len(disks))
		for _, e := range disks {
			pack := e.(map[string]interface{})
			dataDisk := cs.DataDisk{
				Size:                 pack["size"].(string),
				DiskName:             pack["name"].(string),
				Category:             pack["category"].(string),
				Device:               pack["device"].(string),
				AutoSnapshotPolicyId: pack["auto_snapshot_policy_id"].(string),
				KMSKeyId:             pack["kms_key_id"].(string),
				Encrypted:            pack["encrypted"].(string),
				PerformanceLevel:     pack["performance_level"].(string),
			}
			createDataDisks = append(createDataDisks, dataDisk)
		}
		creationArgs.WorkerDataDisks = createDataDisks
	}
	if workerDiskSize != 0 {
		creationArgs.WorkerArgs.WorkerSystemDiskSize = workerDiskSize
	}

	if v, ok := d.GetOk("worker_instance_charge_type"); ok {
		creationArgs.WorkerInstanceChargeType = v.(string)
		if creationArgs.WorkerInstanceChargeType == string(PrePaid) {
			creationArgs.WorkerAutoRenew = d.Get("worker_auto_renew").(bool)
			creationArgs.WorkerAutoRenewPeriod = d.Get("worker_auto_renew_period").(int)
			creationArgs.WorkerPeriod = d.Get("worker_period").(int)
			creationArgs.WorkerPeriodUnit = d.Get("worker_period_unit").(string)
		}
	}

	if v, ok := d.GetOk("cluster_spec"); ok {
		creationArgs.ClusterSpec = v.(string)
	}

	if encryptionProviderKey, ok := d.GetOk("encryption_provider_key"); ok {
		creationArgs.EncryptionProviderKey = encryptionProviderKey.(string)
	}

	if rdsInstances, ok := d.GetOk("rds_instances"); ok {
		creationArgs.RdsInstances = expandStringList(rdsInstances.([]interface{}))
	}

	if nodePortRange, ok := d.GetOk("node_port_range"); ok {
		creationArgs.NodePortRange = nodePortRange.(string)
	}

	if runtime, ok := d.GetOk("runtime"); ok {
		if v := runtime.(map[string]interface{}); len(v) > 0 {
			creationArgs.Runtime = expandKubernetesRuntimeConfig(v)
		}
	}

	if taints, ok := d.GetOk("taints"); ok {
		if v := taints.([]interface{}); len(v) > 0 {
			creationArgs.Taints = expandKubernetesTaintsConfig(v)
		}
	}

	// Cluster maintenance window. Effective only in the professional managed cluster
	if v, ok := d.GetOk("maintenance_window"); ok {
		creationArgs.MaintenanceWindow = expandMaintenanceWindowConfig(v.([]interface{}))
	}

	// Configure control plane log. Effective only in the professional managed cluster
	if v, ok := d.GetOk("control_plane_log_components"); ok {
		creationArgs.ControlplaneComponents = expandStringList(v.([]interface{}))
		// ttl default is 30 days
		creationArgs.ControlplaneLogTTL = "30"
	}
	if v, ok := d.GetOk("control_plane_log_ttl"); ok {
		creationArgs.ControlplaneLogTTL = v.(string)
	}
	if v, ok := d.GetOk("control_plane_log_project"); ok {
		creationArgs.ControlplaneLogProject = v.(string)
	}

	return creationArgs, nil
}

func expandKubernetesTaintsConfig(l []interface{}) []cs.Taint {
	config := []cs.Taint{}

	for _, v := range l {
		if m, ok := v.(map[string]interface{}); ok {
			config = append(config, cs.Taint{
				Key:    m["key"].(string),
				Value:  m["value"].(string),
				Effect: cs.Effect(m["effect"].(string)),
			})
		}
	}

	return config
}

func expandKubernetesRuntimeConfig(l map[string]interface{}) cs.Runtime {
	config := cs.Runtime{}

	if v, ok := l["name"]; ok && v != "" {
		config.Name = v.(string)
	}
	if v, ok := l["version"]; ok && v != "" {
		config.Version = v.(string)
	}

	return config
}

func flattenAlicloudCSCertificate(certificate *roacs.DescribeClusterUserKubeconfigResponseBody) map[string]string {
	if certificate == nil {
		return map[string]string{}
	}

	kubeConfig := make(map[string]interface{})
	_ = yaml.Unmarshal([]byte(tea.StringValue(certificate.Config)), &kubeConfig)

	m := make(map[string]string)
	m["cluster_cert"] = kubeConfig["clusters"].([]interface{})[0].(map[interface{}]interface{})["cluster"].(map[interface{}]interface{})["certificate-authority-data"].(string)
	m["client_cert"] = kubeConfig["users"].([]interface{})[0].(map[interface{}]interface{})["user"].(map[interface{}]interface{})["client-certificate-data"].(string)
	m["client_key"] = kubeConfig["users"].([]interface{})[0].(map[interface{}]interface{})["user"].(map[interface{}]interface{})["client-key-data"].(string)

	return m
}

// ACK pro maintenance window
func expandMaintenanceWindowConfig(l []interface{}) (config cs.MaintenanceWindow) {
	if len(l) == 0 || l[0] == nil {
		return
	}

	m := l[0].(map[string]interface{})

	if v, ok := m["enable"]; ok {
		config.Enable = v.(bool)
	}
	if v, ok := m["maintenance_time"]; ok && v != "" {
		config.MaintenanceTime = cs.MaintenanceTime(v.(string))
	}
	if v, ok := m["duration"]; ok && v != "" {
		config.Duration = v.(string)
	}
	if v, ok := m["weekly_period"]; ok && v != "" {
		config.WeeklyPeriod = cs.WeeklyPeriod(v.(string))
	}

	return
}

func expandMaintenanceWindowConfigRoa(l []interface{}) *roacs.MaintenanceWindow {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})
	config := &roacs.MaintenanceWindow{}
	if v, ok := m["enable"]; ok {
		config.SetEnable(v.(bool))
	}
	if v, ok := m["maintenance_time"]; ok {
		config.SetMaintenanceTime(v.(string))
	}
	if v, ok := m["duration"]; ok {
		config.SetDuration(v.(string))
	}
	if v, ok := m["weekly_period"]; ok {
		config.SetWeeklyPeriod(v.(string))
	}

	return config
}

func flattenMaintenanceWindowConfig(config *cs.MaintenanceWindow) (m []map[string]interface{}) {
	if config == nil {
		return []map[string]interface{}{}
	}

	m = append(m, map[string]interface{}{
		"enable":           config.Enable,
		"maintenance_time": config.MaintenanceTime,
		"duration":         config.Duration,
		"weekly_period":    config.WeeklyPeriod,
	})

	return
}

func flattenMaintenanceWindowConfigRoa(config *roacs.MaintenanceWindow) (m []map[string]interface{}) {
	if config == nil {
		return []map[string]interface{}{}
	}

	m = append(m, map[string]interface{}{
		"enable":           config.Enable,
		"maintenance_time": config.MaintenanceTime,
		"duration":         config.Duration,
		"weekly_period":    config.WeeklyPeriod,
	})

	return
}

// getApiServerSlbID gets cluster's API server SLB ID.
func getApiServerSlbID(d *schema.ResourceData, meta interface{}) (string, error) {
	rosClient, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return "", err
	}
	var clusterResources *roacs.DescribeClusterResourcesResponse
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := &roacs.DescribeClusterResourcesRequest{}
		clusterResources, err = rosClient.DescribeClusterResources(tea.String(d.Id()), request)
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
		return "", err
	}

	for _, clusterResource := range clusterResources.Body {
		if tea.StringValue(clusterResource.ResourceType) == "SLB" || tea.StringValue(clusterResource.ResourceType) == "ALIYUN::SLB::LoadBalancer" {
			return tea.StringValue(clusterResource.InstanceId), nil
		}
	}

	return "", fmt.Errorf("cannot found api server SLB information for cluster: %s", d.Id())
}

func fetchClusterCapabilities(meta string) map[string]interface{} {
	metadata := make(map[string]interface{}, 0)
	capabilities := make(map[string]interface{}, 0)
	if meta != "" {
		err := json.Unmarshal([]byte(meta), &metadata)
		if err != nil {
			log.Printf("[DEBUG] Failed to unmarshal metadata due to %++v", err)
		}
	}
	if v, ok := metadata["Capabilities"]; ok {
		if IsEmpty(v) {
			return capabilities
		}
		if m, ok := v.(map[string]interface{}); ok {
			return m
		}
	}
	return capabilities
}

type RRSAMetadata struct {
	Enabled      bool   `json:"enabled"`
	IssuerURL    string `json:"issuer"`
	ProviderName string `json:"oidc_name"`
	ProviderArn  string `json:"oidc_arn"`
}

func flattenRRSAMetadata(meta string) ([]map[string]interface{}, error) {
	meta = strings.TrimSpace(meta)
	if meta == "" {
		return nil, errors.New("invalid metadata")
	}
	metadata := struct {
		RRSAMetadata RRSAMetadata `json:"RRSAConfig"`
	}{}

	err := json.Unmarshal([]byte(meta), &metadata)
	if err != nil {
		log.Printf("[DEBUG] Failed to unmarshal metadata due to %++v", err)
		return nil, err
	}

	data := metadata.RRSAMetadata
	attributes := map[string]interface{}{
		"enabled":                data.Enabled,
		"rrsa_oidc_issuer_url":   "",
		"ram_oidc_provider_name": "",
		"ram_oidc_provider_arn":  "",
	}
	if !data.Enabled {
		return []map[string]interface{}{attributes}, nil
	}

	issuer := data.IssuerURL
	if strings.Contains(issuer, ",") {
		issuer = strings.Split(issuer, ",")[0]
	}
	attributes["rrsa_oidc_issuer_url"] = issuer
	attributes["ram_oidc_provider_name"] = data.ProviderName
	attributes["ram_oidc_provider_arn"] = data.ProviderArn

	return []map[string]interface{}{attributes}, nil
}

func fetchMasterNodes(d *schema.ResourceData, meta interface{}) []map[string]interface{} {
	csClient, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return nil
	}
	var nodes []*roacs.DescribeClusterNodesResponseBodyNodes
	num := 1
	size := 100
	for {
		request := &roacs.DescribeClusterNodesRequest{
			PageNumber: tea.String(strconv.Itoa(num)),
			PageSize:   tea.String(strconv.Itoa(size)),
		}
		var response *roacs.DescribeClusterNodesResponse
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = csClient.DescribeClusterNodes(tea.String(d.Id()), request)
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
			return nil
		}
		if response.Body == nil || response.Body.Nodes == nil {
			break
		}
		if len(response.Body.Nodes) > 0 {
			nodes = append(nodes, response.Body.Nodes...)
		}
		if len(response.Body.Nodes) == 0 || (response.Body.Page != nil && len(nodes) >= int(*response.Body.Page.TotalCount)) {
			break
		}
		num = num + 1
	}
	var masterNodes []map[string]interface{}
	for _, node := range nodes {
		if *node.InstanceRole != "Master" {
			continue
		}
		masterNodes = append(masterNodes, map[string]interface{}{
			"id":         node.InstanceId,
			"name":       node.NodeName,
			"private_ip": node.IpAddress[0],
		})
	}

	return masterNodes
}

func fetchWorkerNodes(d *schema.ResourceData, meta interface{}) []map[string]interface{} {
	csClient, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return nil
	}
	var response *roacs.DescribeClusterNodePoolsResponse
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = csClient.DescribeClusterNodePools(tea.String(d.Id()), &roacs.DescribeClusterNodePoolsRequest{})
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
		return nil
	}
	nodepoolId := ""
	if response.Body != nil && response.Body.Nodepools != nil {
		for _, nodepool := range response.Body.Nodepools {
			if *nodepool.NodepoolInfo.Type == defaultNodePoolType && *nodepool.NodepoolInfo.Name == "default-nodepool" {
				nodepoolId = *nodepool.NodepoolInfo.NodepoolId
				break
			}
		}
	}
	if nodepoolId == "" {
		return nil
	}

	var nodes []*roacs.DescribeClusterNodesResponseBodyNodes
	num := 1
	size := 100
	for {
		request := &roacs.DescribeClusterNodesRequest{
			PageNumber: tea.String(strconv.Itoa(num)),
			PageSize:   tea.String(strconv.Itoa(size)),
			NodepoolId: tea.String(nodepoolId),
		}
		var response *roacs.DescribeClusterNodesResponse
		wait = incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = csClient.DescribeClusterNodes(tea.String(d.Id()), request)
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
			return nil
		}
		if response.Body == nil || response.Body.Nodes == nil {
			break
		}
		if len(response.Body.Nodes) > 0 {
			nodes = append(nodes, response.Body.Nodes...)
		}
		if len(response.Body.Nodes) == 0 || (response.Body.Page != nil && len(nodes) >= int(*response.Body.Page.TotalCount)) {
			break
		}
		num = num + 1
	}
	var workerNodes []map[string]interface{}
	for _, node := range nodes {
		workerNodes = append(workerNodes, map[string]interface{}{
			"id":         node.InstanceId,
			"name":       node.NodeName,
			"private_ip": node.IpAddress[0],
		})
	}

	return workerNodes
}

func flattenTags(config []*roacs.Tag) map[string]string {
	m := make(map[string]string, len(config))
	if len(config) < 0 {
		return m
	}

	for _, tag := range config {
		key := tea.StringValue(tag.Key)
		value := tea.StringValue(tag.Value)
		if key != DefaultClusterTag && key != CsPlayerAccountIdTag {
			m[key] = value
		}
	}

	return m
}

func fetchClusterMetaDataMap(meta string) map[string]interface{} {
	metadata := make(map[string]interface{}, 0)
	if meta != "" {
		err := json.Unmarshal([]byte(meta), &metadata)
		if err != nil {
			log.Printf("[DEBUG] Failed to unmarshal metadata due to %++v", err)
		}
	}

	return metadata
}
