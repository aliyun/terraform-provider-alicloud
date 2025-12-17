package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"

	"github.com/alibabacloud-go/tea/tea"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	roacs "github.com/alibabacloud-go/cs-20151215/v5/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCSManagedKubernetes() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSManagedKubernetesCreate,
		Read:   resourceAlicloudCSManagedKubernetesRead,
		Update: resourceAlicloudCSManagedKubernetesUpdate,
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
			},
			"profile": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
				Optional: true,
			},
			"worker_vswitch_ids": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				Optional:   true,
				Deprecated: "Field 'worker_vswitch_ids' has been deprecated from provider version 1.241.0. Please use 'vswitch_ids' to managed control plane vswtiches",
			},
			"worker_instance_types": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MaxItems: 10,
				Removed:  "Field 'worker_instance_types' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'instance_types' to replace it.",
			},
			"worker_number": {
				Type:     schema.TypeInt,
				Optional: true,
				Removed:  "Field 'worker_number' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes., by using field 'desired_size' to replace it.",
			},
			"worker_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(20, 32768),
				Removed:      "Field 'worker_disk_size' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'system_disk_size' to replace it.",
			},
			"worker_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field 'worker_disk_category' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'system_disk_category' to replace it.",
			},
			"worker_disk_performance_level": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
				DiffSuppressFunc: workerDiskPerformanceLevelDiffSuppressFunc,
				Removed:          "Field 'worker_disk_performance_level' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'system_disk_performance_level' to replace it",
			},
			"worker_disk_snapshot_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field 'worker_disk_snapshot_policy_id' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'system_disk_snapshot_policy_id' to replace it",
			},
			"worker_data_disk_size": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     IntBetween(20, 32768),
				DiffSuppressFunc: workerDataDiskSizeSuppressFunc,
				Removed:          "Field 'worker_data_disk_size' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'data_disks.size' to replace it",
			},
			"worker_data_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field 'worker_data_disk_category' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'data_disks.category' to replace it",
			},
			"worker_instance_charge_type": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				Removed:  "Field 'worker_instance_charge_type' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'instance_charge_type' to replace it",
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
				Removed: "Field 'worker_data_disks' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'data_disks' to replace it",
			},
			"worker_period_unit": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     StringInSlice([]string{"Week", "Month"}, false),
				DiffSuppressFunc: csKubernetesWorkerPostPaidDiffSuppressFunc,
				Removed:          "Field 'worker_period_unit' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'period_unit' to replace it",
			},
			"worker_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: csKubernetesWorkerPostPaidDiffSuppressFunc,
				Removed:          "Field 'worker_period' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'period' to replace it",
			},
			"worker_auto_renew": {
				Type:             schema.TypeBool,
				Optional:         true,
				DiffSuppressFunc: csKubernetesWorkerPostPaidDiffSuppressFunc,
				Removed:          "Field 'worker_auto_renew' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'auto_renew' to replace it",
			},
			"worker_auto_renew_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     IntInSlice([]int{1, 2, 3, 6, 12}),
				DiffSuppressFunc: csKubernetesWorkerPostPaidDiffSuppressFunc,
				Removed:          "Field 'worker_auto_renew_period' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'auto_renew_period' to replace it",
			},
			"exclude_autoscaler_nodes": {
				Type:     schema.TypeBool,
				Optional: true,
				Removed:  "Field 'exclude_autoscaler_nodes' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes",
			},
			// global configurations
			"zone_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems: 1,
				MaxItems: 5,
			},
			"vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				MinItems:     1,
				MaxItems:     5,
				ExactlyOneOf: []string{"worker_vswitch_ids", "vswitch_ids", "zone_ids"},
			},
			"pod_vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
			},
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
			"enable_ssh": {
				Type:     schema.TypeBool,
				Optional: true,
				Removed:  "Field 'enable_ssh' has been removed from provider version 1.212.0.",
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
				Removed:       "Field 'password' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'password' to replace it",
			},
			"key_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"password", "kms_encrypted_password"},
				Removed:       "Field 'key_name' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'key_name' to replace it",
			},
			"kms_encrypted_password": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"password", "key_name"},
				Removed:       "Field 'kms_encrypted_password' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'kms_encrypted_password' to replace it",
			},
			"kms_encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_password").(string) == ""
				},
				Elem:    schema.TypeString,
				Removed: "Field 'kms_encryption_context' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'kms_encryption_context' to replace it",
			},
			"user_ca": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field 'image_id' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'image_id' to replace it",
			},
			"install_cloud_monitor": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Removed:  "Field 'install_cloud_monitor' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'install_cloud_monitor' to replace it",
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
				Removed:      "Field 'cpu_policy' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'cpu_policy' to replace it",
			},
			"proxy_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "ipvs",
				ValidateFunc: StringInSlice([]string{"iptables", "ipvs"}, false),
			},
			"ip_stack": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"ipv4", "dual"}, false),
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
			"slb_internet_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"load_balancer_spec": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"slb.s1.small", "slb.s2.small", "slb.s2.medium", "slb.s3.small", "slb.s3.medium", "slb.s3.large"}, false),
				Computed:     true,
				Deprecated:   "Field 'load_balancer_spec' has been deprecated from provider version 1.232.0. The spec will not take effect because the charge of the load balancer has been changed to PayByCLCU",
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"enable_rrsa": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"os_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Windows", "Linux"}, false),
				Removed:      "Field 'os_type' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes.",
			},
			"platform": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Removed:  "Field 'platform' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'platform' to replace it.",
			},
			"node_port_range": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Removed:  "Field 'node_port_range' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes.",
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
				Removed: "Field 'runtime' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'runtime_name' and 'runtime_version' to replace it.",
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
				Removed: "Field 'taints' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'taints' to replace it.",
			},
			"rds_instances": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Removed: "Field 'rds_instances' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'rds_instances' to replace it.",
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field 'user_data' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'user_data' to replace it.",
			},
			"node_name_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile(`^customized,[a-z0-9]([-a-z0-9\.])*,([5-9]|[1][0-2]),([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`), "Each node name consists of a prefix, an IP substring, and a suffix. For example, if the node IP address is 192.168.0.55, the prefix is aliyun.com, IP substring length is 5, and the suffix is test, the node name will be aliyun.com00055test."),
				Removed:      "Field 'node_name_mode' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes, by using field 'node_name_mode' to replace it.",
			},
			"worker_nodes": {
				Type:     schema.TypeList,
				Optional: true,
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
				Removed: "Field 'worker_nodes' has been removed from provider version 1.212.0. Please use resource 'alicloud_cs_kubernetes_node_pool' to manage cluster nodes.",
			},
			"custom_san": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encryption_provider_key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The ID of the Key Management Service (KMS) key that is used to encrypt Kubernetes Secrets.",
			},
			// computed parameters
			"kube_config": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field 'kube_config' has been removed from provider version 1.212.0. Please use the attribute 'output_file' of new DataSource 'alicloud_cs_cluster_credential' to replace it",
			},
			"client_cert": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'client_cert' has been deprecated from provider version 1.248.0. From version 1.248.0, new DataSource 'alicloud_cs_cluster_credential' is recommended to manage cluster's kubeconfig, you can also save the 'certificate_authority.client_cert' attribute content of new DataSource 'alicloud_cs_cluster_credential' to an appropriate path(like ~/.kube/client-cert.pem) for replace it.",
			},
			"client_key": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'client_key' has been deprecated from provider version 1.248.0. From version 1.248.0, new DataSource 'alicloud_cs_cluster_credential' is recommended to manage cluster's kubeconfig, you can also save the 'certificate_authority.client_key' attribute content of new DataSource 'alicloud_cs_cluster_credential' to an appropriate path(like ~/.kube/client-key.pem) for replace it.",
			},
			"cluster_ca_cert": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'cluster_ca_cert' has been deprecated from provider version 1.248.0. From version 1.248.0, new DataSource 'alicloud_cs_cluster_credential' is recommended to manage cluster's kubeconfig, you can also save the 'certificate_authority.cluster_cert' attribute content of new DataSource 'alicloud_cs_cluster_credential' to an appropriate path(like ~/.kube/cluster-ca-cert.pem) for replace it.",
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
				Deprecated: "Field 'certificate_authority' has been deprecated from provider version 1.248.0. Please use the attribute 'certificate_authority' of new DataSource 'alicloud_cs_cluster_credential' to replace it.",
			},
			"skip_set_certificate_authority": {
				Type:     schema.TypeBool,
				Optional: true,
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
				Removed:  "Field 'availability_zone' has been removed from provider version 1.212.0.",
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
				ValidateFunc: StringInSlice([]string{KubernetesClusterNetworkTypeFlannel, KubernetesClusterNetworkTypeTerway}, false),
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
			},
			"cluster_spec": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"ack.standard", "ack.pro.small"}, false),
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
							Optional: true,
							Computed: true,
						},
						"maintenance_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"duration": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"weekly_period": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"operation_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_auto_upgrade": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"channel": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"control_plane_log_ttl": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"control_plane_log_project": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"control_plane_log_components": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 0,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"audit_log_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"sls_project_name": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
					},
				},
			},
			"auto_mode": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							ForceNew: true,
							Optional: true,
						},
					},
				},
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
			"rrsa_metadata": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"rrsa_oidc_issuer_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ram_oidc_provider_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ram_oidc_provider_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudCSManagedKubernetesCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	roa, _ := client.NewRoaCsClient()
	csClient := CsClient{roa}

	var clusterName string
	if v, ok := d.GetOk("name"); ok {
		clusterName = v.(string)
	} else {
		clusterName = resource.PrefixedUniqueId(d.Get("name_prefix").(string))
	}

	tags := make([]*roacs.Tag, 0)
	tagsMap, ok := d.Get("tags").(map[string]interface{})
	if ok {
		for key, value := range tagsMap {
			if value != nil {
				if v, ok := value.(string); ok {
					tags = append(tags, &roacs.Tag{
						Key:   tea.String(key),
						Value: tea.String(v),
					})
				}
			}
		}
	}
	addons := make([]*roacs.Addon, 0)
	if v, ok := d.GetOk("addons"); ok {
		all, ok := v.([]interface{})
		if ok {
			for _, a := range all {
				addon, ok := a.(map[string]interface{})
				if ok {
					addons = append(addons, &roacs.Addon{
						Name:     tea.String(addon["name"].(string)),
						Config:   tea.String(addon["config"].(string)),
						Version:  tea.String(addon["version"].(string)),
						Disabled: tea.Bool(addon["disabled"].(bool)),
					})
				}
			}
		}
	}

	vpcService := VpcService{client}
	var vSwitchIds []string
	if v, ok := d.GetOk("vswitch_ids"); ok {
		vSwitchIds = expandStringList(v.([]interface{}))
	} else {
		if v, ok := d.GetOk("worker_vswitch_ids"); ok {
			vSwitchIds = expandStringList(v.([]interface{}))
		}
	}
	var vpcId string
	if len(vSwitchIds) > 0 {
		vsw, err := vpcService.DescribeVSwitch(vSwitchIds[0])
		if err != nil {
			return err
		}
		vpcId = vsw.VpcId
	}

	request := &roacs.CreateClusterRequest{
		Name:        tea.String(clusterName),
		RegionId:    tea.String(client.RegionId),
		ClusterType: tea.String("ManagedKubernetes"),
		Profile:     tea.String("Default"),
		Tags:        tags,
		Addons:      addons,
		Vpcid:       tea.String(vpcId),
		VswitchIds:  tea.StringSlice(vSwitchIds),
	}
	if v, ok := d.GetOk("profile"); ok {
		request.SetProfile(v.(string))
	}

	if v, ok := d.GetOk("version"); ok {
		request.SetKubernetesVersion(v.(string))
	}

	if v, ok := d.GetOkExists("deletion_protection"); ok {
		request.SetDeletionProtection(v.(bool))
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request.SetResourceGroupId(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.SetVpcid(v.(string))
	}

	if v, ok := d.GetOk("zone_ids"); ok {
		request.SetZoneIds(tea.StringSlice(expandStringList(v.([]interface{}))))
	}

	if v, ok := d.GetOk("new_nat_gateway"); ok {
		request.SetSnatEntry(v.(bool))
	}

	if v, ok := d.GetOk("slb_internet_enabled"); ok {
		request.SetEndpointPublicAccess(v.(bool))
	}

	if v, ok := d.GetOk("load_balancer_spec"); ok {
		request.SetLoadBalancerSpec(v.(string))
	}

	if v, ok := d.GetOk("is_enterprise_security_group"); ok {
		request.SetIsEnterpriseSecurityGroup(v.(bool))
	}

	if v, ok := d.GetOk("security_group_id"); ok {
		request.SetSecurityGroupId(v.(string))
	}

	if v, ok := d.GetOk("service_cidr"); ok {
		request.SetServiceCidr(v.(string))
	}

	if v, ok := d.GetOk("proxy_mode"); ok {
		request.SetProxyMode(v.(string))
	}

	if v, ok := d.GetOk("ip_stack"); ok {
		request.SetIpStack(v.(string))
	}

	if v, ok := d.GetOk("timezone"); ok {
		request.SetTimezone(v.(string))
	}

	if v, ok := d.GetOk("pod_vswitch_ids"); ok {
		request.SetPodVswitchIds(tea.StringSlice(expandStringList(v.([]interface{}))))
	}

	if v, ok := d.GetOk("pod_cidr"); ok {
		request.SetContainerCidr(v.(string))
	}

	if v, ok := d.GetOk("node_cidr_mask"); ok {
		request.SetNodeCidrMask(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("cluster_spec"); ok {
		request.SetClusterSpec(v.(string))
	}

	if v, ok := d.GetOk("cluster_domain"); ok {
		request.SetClusterDomain(v.(string))
	}

	if v, ok := d.GetOk("service_account_issuer"); ok {
		request.SetServiceAccountIssuer(v.(string))
	}

	if v, ok := d.GetOk("api_audiences"); ok {
		if list := expandStringList(v.([]interface{})); len(list) > 0 {
			request.SetApiAudiences(strings.Join(list, ","))
		}
	}

	if v, ok := d.GetOk("enable_rrsa"); ok {
		request.SetEnableRrsa(v.(bool))
	}
	if v, ok := d.GetOk("custom_san"); ok {
		request.SetCustomSan(v.(string))
	}

	if v, ok := d.GetOk("encryption_provider_key"); ok {
		request.SetEncryptionProviderKey(v.(string))
	}

	// Configure control plane log. Effective only in the professional managed cluster
	if v, ok := d.GetOk("control_plane_log_components"); ok {
		request.SetControlplaneLogComponents(tea.StringSlice(expandStringList(v.([]interface{}))))
		// ttl default is 30 days
		request.SetControlplaneLogTtl("30")
	}
	if v, ok := d.GetOk("control_plane_log_ttl"); ok {
		request.SetControlplaneLogTtl(v.(string))
	}
	if v, ok := d.GetOk("control_plane_log_project"); ok {
		request.SetControlplaneLogProject(v.(string))
	}

	if v, ok := d.GetOk("maintenance_window"); ok {
		request.SetMaintenanceWindow(expandMaintenanceWindowConfigRoa(v.([]interface{})))
	}
	if v, ok := d.GetOk("operation_policy"); ok {
		request.OperationPolicy = &roacs.CreateClusterRequestOperationPolicy{}
		m := v.([]interface{})[0].(map[string]interface{})
		if vv, ok := m["cluster_auto_upgrade"]; ok {
			policy := vv.([]interface{})[0].(map[string]interface{})
			request.OperationPolicy.ClusterAutoUpgrade = &roacs.CreateClusterRequestOperationPolicyClusterAutoUpgrade{
				Enabled: tea.Bool(policy["enabled"].(bool)),
				Channel: tea.String(policy["channel"].(string)),
			}
		}
	}
	if v, ok := d.GetOk("audit_log_config"); ok {
		m := v.([]interface{})[0].(map[string]interface{})
		if vv, ok := m["enabled"]; ok {
			request.AuditLogConfig = &roacs.CreateClusterRequestAuditLogConfig{
				Enabled: tea.Bool(vv.(bool)),
			}
		}
		if vv, ok := m["sls_project_name"]; ok {
			request.AuditLogConfig.SlsProjectName = tea.String(vv.(string))
		}
	}

	if v, ok := d.GetOk("auto_mode"); ok {
		m := v.([]interface{})[0].(map[string]interface{})
		if vv, ok := m["enabled"]; ok {
			request.AutoMode = &roacs.CreateClusterRequestAutoMode{
				Enable: tea.Bool(vv.(bool)),
			}
		}
	}

	var err error
	var resp *roacs.CreateClusterResponse
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = csClient.client.CreateCluster(request)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_managed_kubernetes", "CreateManagedKubernetesCluster", AlibabaCloudSdkGoERROR)
	}
	d.SetId(tea.StringValue(resp.Body.ClusterId))
	taskId := tea.StringValue(resp.Body.TaskId)
	roaCsClient, err := client.NewRoaCsClient()
	if err == nil {
		csClient := CsClient{client: roaCsClient}
		stateConf := BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, csClient.DescribeTaskRefreshFunc(d, taskId, []string{"fail", "failed"}))
		if jobDetail, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, ResponseCodeMsg, d.Id(), "createCluster", jobDetail)
		}
	}

	csService := CsService{client}
	stateConf := BuildStateConf([]string{"initial"}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudCSManagedKubernetesRead(d, meta)
}

func resourceAlicloudCSManagedKubernetesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rosClient, err := client.NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "InitializeClient", err)
	}
	csClient := CsClient{rosClient}

	object, err := csClient.DescribeClusterDetail(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cs_managed_kubernetes DescribeClusterDetail Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	// compat for default value
	if spec := d.Get("load_balancer_spec").(string); spec != "" {
		d.Set("load_balancer_spec", spec)
	}

	if object.Name != nil {
		d.Set("name", object.Name)
	}

	if object.Profile != nil {
		d.Set("profile", object.Profile)
	}

	if object.VpcId != nil {
		d.Set("vpc_id", object.VpcId)
	}

	if object.VswitchIds != nil {
		d.Set("vswitch_ids", tea.StringSliceValue(object.VswitchIds))
	}

	// compat for old value
	if v := d.Get("worker_vswitch_ids"); v != nil {
		d.Set("worker_vswitch_ids", v)
	}

	if object.SecurityGroupId != nil {
		d.Set("security_group_id", object.SecurityGroupId)
	}

	if object.DeletionProtection != nil {
		d.Set("deletion_protection", object.DeletionProtection)
	}

	if object.CurrentVersion != nil {
		d.Set("version", object.CurrentVersion)
	}

	if object.ResourceGroupId != nil {
		d.Set("resource_group_id", object.ResourceGroupId)
	}

	if object.ClusterSpec != nil {
		d.Set("cluster_spec", object.ClusterSpec)
	}

	if object.Timezone != nil {
		d.Set("timezone", object.Timezone)
	}

	if object.WorkerRamRoleName != nil {
		d.Set("worker_ram_role_name", object.WorkerRamRoleName)
	}

	d.Set("cluster_domain", "cluster.local")
	if object.ClusterDomain != nil {
		d.Set("cluster_domain", object.ClusterDomain)
	}

	if err := d.Set("tags", flattenTags(object.Tags)); err != nil {
		return WrapError(err)
	}

	slbId, err := getApiServerSlbID(d, meta)
	if err != nil {
		log.Printf(DefaultErrorMsg, d.Id(), "DescribeClusterResources", err.Error())
	}
	d.Set("slb_id", slbId)

	if object.ServiceCidr != nil {
		d.Set("service_cidr", object.ServiceCidr)
	} else {
		if v, ok := object.Parameters["ServiceCIDR"]; ok {
			d.Set("service_cidr", v)
		}
	}
	if object.ProxyMode != nil {
		d.Set("proxy_mode", object.ProxyMode)
	} else {
		if v, ok := object.Parameters["ProxyMode"]; ok {
			d.Set("proxy_mode", v)
		}
	}

	if object.IpStack != nil {
		d.Set("ip_stack", object.IpStack)
	}

	if object.ContainerCidr != nil {
		d.Set("pod_cidr", object.ContainerCidr)
	} else {
		if v, ok := object.Parameters["ContainerCIDR"]; ok {
			d.Set("pod_cidr", v)
		}
	}

	if object.NodeCidrMask != nil {
		d.Set("node_cidr_mask", formatInt(tea.StringValue(object.NodeCidrMask)))
	} else {
		// node_cidr_mask
		capabilities := fetchClusterCapabilities(tea.StringValue(object.MetaData))
		if v, ok := capabilities["NodeCIDRMask"]; ok {
			d.Set("node_cidr_mask", formatInt(v))
		}
	}

	metadata := fetchClusterMetaDataMap(tea.StringValue(object.MetaData))
	if v, ok := metadata["ExtraCertSAN"]; ok && v != nil {
		l := expandStringList(v.([]interface{}))
		d.Set("custom_san", strings.Join(l, ","))
	}
	// rrsa metadata only for managed, ignore attributes error
	if data, err := flattenRRSAMetadata(tea.StringValue(object.MetaData)); err != nil {
		return WrapError(err)
	} else {
		d.Set("rrsa_metadata", data)
		if len(data) > 0 {
			d.Set("enable_rrsa", data[0]["enabled"].(bool))
		}
	}

	if object.MaintenanceWindow != nil {
		d.Set("maintenance_window", flattenMaintenanceWindowConfigRoa(object.MaintenanceWindow))
	}

	if object.OperationPolicy != nil {
		m := make([]map[string]interface{}, 0)
		if object.OperationPolicy.ClusterAutoUpgrade != nil {
			m = append(m, map[string]interface{}{
				"cluster_auto_upgrade": []map[string]interface{}{
					{
						"enabled": tea.BoolValue(object.OperationPolicy.ClusterAutoUpgrade.Enabled),
						"channel": tea.StringValue(object.OperationPolicy.ClusterAutoUpgrade.Channel),
					},
				},
			})
		}
		d.Set("operation_policy", m)
	}

	if object.AutoMode != nil {
		m := make(map[string]interface{})
		if object.AutoMode.Enable != nil {
			m["enabled"] = tea.BoolValue(object.AutoMode.Enable)
		}
		d.Set("auto_mode", []map[string]interface{}{m})
	}

	// Get slb information and set connect
	connection := make(map[string]string)
	masterURL := tea.StringValue(object.MasterUrl)
	endPoint := make(map[string]string)
	_ = json.Unmarshal([]byte(masterURL), &endPoint)
	connection["api_server_internet"] = endPoint["api_server_endpoint"]
	connection["api_server_intranet"] = endPoint["intranet_api_server_endpoint"]
	if endPoint["api_server_endpoint"] != "" {
		connection["master_public_ip"] = strings.Split(strings.Split(endPoint["api_server_endpoint"], ":")[1], "/")[2]
	}
	connection["service_domain"] = fmt.Sprintf("*.%s.%s.alicontainer.com", d.Id(), tea.StringValue(object.RegionId))

	d.Set("connections", connection)
	d.Set("slb_internet", connection["master_public_ip"])
	if endPoint["intranet_api_server_endpoint"] != "" {
		d.Set("slb_intranet", strings.Split(strings.Split(endPoint["intranet_api_server_endpoint"], ":")[1], "/")[2])
	}

	// set nat gateway
	natRequest := vpc.CreateDescribeNatGatewaysRequest()
	natRequest.VpcId = tea.StringValue(object.VpcId)
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
	if tea.StringValue(object.State) == "failed" || tea.StringValue(object.State) == "delete_failed" || tea.StringValue(object.State) == "deleting" {
		return nil
	}

	if err = setCerts(d, meta, d.Get("skip_set_certificate_authority").(bool)); err != nil {
		return WrapError(err)
	}

	if err = checkControlPlaneLogEnable(d, meta); err != nil {
		return WrapError(err)
	}

	if err = getClusterAuditProject(d, meta); err != nil {
		return WrapError(err)
	}

	return nil

}

func resourceAlicloudCSManagedKubernetesUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	invoker := NewInvoker()
	// modifyCluster
	if !d.IsNewResource() && d.HasChanges("resource_group_id", "name", "name_prefix", "deletion_protection", "maintenance_window", "operation_policy",
		"custom_san", "vswitch_ids", "timezone", "security_group_id", "enable_rrsa") {
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

	if d.HasChange("audit_log_config") {
		if err := updateClusterAuditConfig(d, meta); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateClusterAuditConfig", AlibabaCloudSdkGoERROR)
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
	return resourceAlicloudCSManagedKubernetesRead(d, meta)
}

func UpgradeAlicloudKubernetesCluster(d *schema.ResourceData, meta interface{}) error {
	if !d.HasChange("version") {
		return nil
	}

	clusterId := d.Id()
	version := d.Get("version").(string)
	action := "UpgradeCluster"
	c := meta.(*connectivity.AliyunClient)
	rosCsClient, err := c.NewRoaCsClient()
	if err != nil {
		return err
	}
	args := &roacs.UpgradeClusterRequest{
		NextVersion: tea.String(version),
	}
	// upgrade cluster
	var resp *roacs.UpgradeClusterResponse
	err = resource.Retry(UpgradeClusterTimeout, func() *resource.RetryError {
		resp, err = rosCsClient.UpgradeCluster(tea.String(clusterId), args)
		if NeedRetry(err) || resp == nil {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return WrapErrorf(err, ResponseCodeMsg, d.Id(), action, err)
	}

	taskId := tea.StringValue(resp.Body.TaskId)
	if taskId == "" {
		return WrapErrorf(err, ResponseCodeMsg, d.Id(), action, resp)
	}

	csClient := CsClient{client: rosCsClient}
	stateConf := BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, csClient.DescribeTaskRefreshFunc(d, taskId, []string{"fail", "failed"}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		// try to cancel task
		wait := incrementalWait(3*time.Second, 3*time.Second)
		_ = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, _err := rosCsClient.CancelTask(tea.String(taskId))
			if _err != nil {
				if NeedRetry(_err) {
					wait()
					return resource.RetryableError(_err)
				}
				log.Printf("[WARN] %s ACK Cluster cancel upgrade error: %#v", clusterId, err)
			}
			return nil
		})
		// output error message
		return WrapErrorf(err, ResponseCodeMsg, d.Id(), action, jobDetail)
	}
	// ensure cluster state is running
	csService := CsService{client: c}
	stateConf = BuildStateConf([]string{}, []string{"running"}, UpgradeClusterTimeout, 10*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(clusterId, []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	d.SetPartial("version")
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

// versionCompare check version,
// if cueVersion is newer than neededVersion return 1
// if curVersion is equal neededVersion return 0
// if curVersion is older than neededVersion return -1
// example: neededVersion = 1.20.11-aliyun.1, curVersion = 1.22.3-aliyun.1, it will return 1
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

func updateControlPlaneLog(d *schema.ResourceData, meta interface{}) error {
	request := &roacs.UpdateControlPlaneLogRequest{}
	client := meta.(*connectivity.AliyunClient)
	csClient, err := client.NewRoaCsClient()
	if err != nil {
		return err
	}
	csService := CsService{client}
	if d.HasChange("control_plane_log_ttl") {
		if v, ok := d.GetOk("control_plane_log_ttl"); ok {
			request.LogTtl = tea.String(v.(string))
		}
	}
	if d.HasChange("control_plane_log_project") {
		if v, ok := d.GetOk("control_plane_log_project"); ok {
			request.LogProject = tea.String(v.(string))
		}
	}
	if d.HasChange("control_plane_log_components") {
		if v, ok := d.GetOk("control_plane_log_components"); ok {
			list := v.([]interface{})
			components := make([]*string, len(list))
			for i, c := range list {
				components[i] = tea.String(c.(string))
			}
			request.Components = components
		}
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = csClient.UpdateControlPlaneLog(tea.String(d.Id()), request)
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
		return err
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

func checkControlPlaneLogEnable(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return err
	}
	var response *roacs.CheckControlPlaneLogEnableResponse
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.CheckControlPlaneLogEnable(tea.String(d.Id()))
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
		return err
	}
	if response.Body != nil {
		if response.Body.LogTtl != nil {
			d.Set("control_plane_log_ttl", *response.Body.LogTtl)
		}
		if response.Body.LogProject != nil {
			d.Set("control_plane_log_project", *response.Body.LogProject)
		}
		components := make([]string, len(response.Body.Components))
		for i, c := range response.Body.Components {
			components[i] = *c
		}
		d.Set("control_plane_log_components", components)
	}

	return nil
}

func updateClusterAuditConfig(d *schema.ResourceData, meta interface{}) error {
	request := &roacs.UpdateClusterAuditLogConfigRequest{}
	client := meta.(*connectivity.AliyunClient)
	csClient, err := client.NewRoaCsClient()
	if err != nil {
		return err
	}

	if d.HasChange("audit_log_config") {
		v, ok := d.GetOk("audit_log_config")
		if ok && len(v.([]interface{})) > 0 {
			m := v.([]interface{})[0].(map[string]interface{})
			if vv, ok := m["enabled"].(bool); ok {
				request.Disable = tea.Bool(!vv)
			}
			if vv, ok := m["sls_project_name"].(string); ok {
				request.SlsProjectName = tea.String(vv)
			}
		}
	}
	csService := CsService{client}
	var response *roacs.UpdateClusterAuditLogConfigResponse
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = csClient.UpdateClusterAuditLogConfig(tea.String(d.Id()), request)
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
		return err
	}
	taskId := tea.StringValue(response.Body.TaskId)
	c := CsClient{client: csClient}
	stateConf := BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, c.DescribeTaskRefreshFunc(d, taskId, []string{"fail", "failed"}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, ResponseCodeMsg, d.Id(), "UpdateClusterAuditLogConfig", jobDetail)
	}

	stateConf = BuildStateConf([]string{"updating"}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func getClusterAuditProject(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return err
	}
	var response *roacs.GetClusterAuditProjectResponse
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.GetClusterAuditProject(tea.String(d.Id()))
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
		return err
	}
	if response.Body != nil {
		m := make(map[string]interface{})
		if response.Body.AuditEnabled != nil {
			m["enabled"] = tea.BoolValue(response.Body.AuditEnabled)
		}
		if response.Body.SlsProjectName != nil {
			m["sls_project_name"] = tea.StringValue(response.Body.SlsProjectName)
		}
		d.Set("audit_log_config", []map[string]interface{}{m})
	}

	return nil
}
