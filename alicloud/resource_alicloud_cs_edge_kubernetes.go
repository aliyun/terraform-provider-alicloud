package alicloud

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	roacs "github.com/alibabacloud-go/cs-20151215/v7/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	EdgeKubernetesDefaultTimeoutInMinutes = 60
	EdgeProfile                           = "Edge"
)

func resourceAlicloudCSEdgeKubernetes() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSEdgeKubernetesCreate,
		Read:   resourceAlicloudCSEdgeKubernetesRead,
		Update: resourceAlicloudCSEdgeKubernetesUpdate,
		Delete: resourceAlicloudCSEdgeKubernetesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: "This resource has been deprecated since v1.276.0 and will be removed in the future. Please use 'alicloud_cs_managed_kubernetes' instead.",

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(EdgeKubernetesDefaultTimeoutInMinutes * time.Minute),
			Update: schema.DefaultTimeout(EdgeKubernetesDefaultTimeoutInMinutes * time.Minute),
			Delete: schema.DefaultTimeout(EdgeKubernetesDefaultTimeoutInMinutes * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type:    resourceAlicloudCSEdgeKubernetesV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceAlicloudCSEdgeKubernetesStateUpgradeV0,
			},
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
			"cluster_spec": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"ack.standard", "ack.pro.small"}, false),
			},
			// worker configurations
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
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems: 1,
				MaxItems: 10,
			},
			//cloud worker number
			"worker_number": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"worker_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      40,
				ValidateFunc: validation.IntBetween(20, 32768),
			},
			"worker_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  DiskCloudEfficiency,
				ValidateFunc: validation.StringInSlice([]string{
					string(DiskCloudEfficiency), string(DiskCloudSSD), string(DiskCloudESSD)}, false),
			},
			"worker_disk_performance_level": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
				DiffSuppressFunc: workerDiskPerformanceLevelDiffSuppressFunc,
			},
			"worker_disk_snapshot_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"proxy_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ipvs",
				ValidateFunc: validation.StringInSlice([]string{"iptables", "ipvs"}, false),
			},
			"worker_instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{string(common.PostPaid)}, false),
				Default:      PostPaid,
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
			},
			"worker_ram_role_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// global configurations
			"pod_cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"runtime": {
				Type:     schema.TypeList,
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
			},
			"node_cidr_mask": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      KubernetesClusterNodeCIDRMasksByDefault,
				ValidateFunc: validation.IntBetween(24, 28),
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
				ConflictsWith: []string{"key_name"},
			},
			"key_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"password"},
			},
			"install_cloud_monitor": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
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
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"slb.s1.small", "slb.s2.small", "slb.s2.medium", "slb.s3.small", "slb.s3.medium", "slb.s3.large"}, false),
			},
			"kube_config": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'kube_config' has been deprecated from provider version 1.187.0. Please use the attribute 'output_file' of new DataSource 'alicloud_cs_cluster_credential' to replace it.",
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
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"rds_instances": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// computed parameters start
			"certificate_authority": {
				Type:     schema.TypeList,
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
				Type:     schema.TypeList,
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
			"slb_internet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slb_id": {
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
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// computed params end

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
				Deprecated: "Field 'log_config' has been removed from provider version 1.103.0. New field 'addons' replaces it.",
			},
			// lintignore: S006
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
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

func resourceAlicloudCSEdgeKubernetesCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request, err := buildEdgeKubernetesRequest(d, meta)
	if err != nil {
		return WrapError(err)
	}
	nodepools := request.Nodepools
	if len(nodepools) != 1 || nodepools[0] == nil || nodepools[0].ScalingGroup == nil {
		return WrapError(Error("exactly one default Edge node pool must be configured"))
	}
	expectedWorkers := int(tea.Int64Value(nodepools[0].ScalingGroup.DesiredSize))

	roaClient, err := client.NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_edge_kubernetes", "InitializeClient", AlibabaCloudSdkGoERROR)
	}
	csClient := CsClient{client: roaClient}
	var response *roacs.CreateClusterResponse
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = csClient.client.CreateCluster(request)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug("/clusters", response, map[string]interface{}{"request": request})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_edge_kubernetes", "CreateCluster", AlibabaCloudSdkGoERROR)
	}
	if response == nil || response.Body == nil {
		return WrapErrorf(Error("CreateCluster returned an empty response"), DefaultErrorMsg, "alicloud_cs_edge_kubernetes", "CreateCluster", AlibabaCloudSdkGoERROR)
	}
	clusterId := tea.StringValue(response.Body.ClusterId)
	if clusterId == "" {
		return WrapErrorf(Error("CreateCluster returned an empty cluster ID"), DefaultErrorMsg, "alicloud_cs_edge_kubernetes", "CreateCluster", AlibabaCloudSdkGoERROR)
	}
	d.SetId(clusterId)
	taskId := tea.StringValue(response.Body.TaskId)
	if taskId != "" {
		stateConf := BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, csClient.DescribeTaskRefreshFunc(d, taskId, []string{"fail", "failed"}))
		if jobDetail, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, ResponseCodeMsg, d.Id(), "createCluster", jobDetail)
		}
	}

	csService := CsService{client}
	stateConf := BuildStateConf([]string{"initial"}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 10*time.Minute, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	// Confirm that ACK accepted the SDK v2 node-pool payload and retained the
	// requested capacity. Kubelet readiness is deliberately not part of the
	// provider success contract, matching the SDK v1 implementation.
	if err := waitForEdgeDefaultNodePool(csClient.client, d.Id(), "", expectedWorkers, 20*time.Minute); err != nil {
		logEdgeKubernetesDiagnostics(d, meta, "", err)
		return WrapError(err)
	}
	// CreateCluster's embedded Nodepool model does not expose
	// scaling_group.system_disk_snapshot_policy_id, so the legacy top-level
	// compatibility field is not applied to the actual SDK v2 node pool.
	// Persist the setting through ModifyClusterNodePool before the first Read.
	if snapshotPolicyId := d.Get("worker_disk_snapshot_policy_id").(string); snapshotPolicyId != "" {
		nodepoolId, err := scaleOutEdgeKubernetesNodePool(d, meta, expectedWorkers)
		if err != nil {
			logEdgeKubernetesDiagnostics(d, meta, nodepoolId, err)
			return WrapError(err)
		}
	}
	return resourceAlicloudCSEdgeKubernetesRead(d, meta)
}

func resourceAlicloudCSEdgeKubernetesDelete(d *schema.ResourceData, meta interface{}) error {
	// Keep the SDK v1 lifecycle contract: DeleteCluster owns normal node-pool
	// cleanup. Pre-deleting a default pool changes imported-cluster semantics
	// and can turn a recoverable node-pool error into a failed destroy.
	return resourceAlicloudCSKubernetesDelete(d, meta)
}

func resourceAlicloudCSEdgeKubernetesRead(d *schema.ResourceData, meta interface{}) error {
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
	d.Set("deletion_protection", object.DeletionProtection)

	if object.ClusterType == cs.ManagedKubernetes {
		d.Set("cluster_spec", object.ClusterSpec)
	}

	// compat for default value
	if spec := d.Get("load_balancer_spec").(string); spec != "" {
		d.Set("load_balancer_spec", spec)
	}

	if err := d.Set("tags", flattenTagsConfig(object.Tags)); err != nil {
		return WrapError(err)
	}

	//request.Parameters
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
		// Node pools created through the SDK v2 nodepools request keep this
		// setting in kubernetes_config.cms_enabled. The cluster-level legacy
		// CloudMonitorFlags parameter remains false even when CMS is enabled,
		// so prefer the default node pool as the source of truth.
		cloudMonitorReadFromNodePool := false
		workerNumberReadFromNodePool := false
		if roaClient, clientErr := client.NewRoaCsClient(); clientErr != nil {
			log.Printf("[WARN] reading Edge default node pool CloudMonitor setting failed to initialize ACK client: %s", clientErr)
		} else if nodepool, nodepoolErr := findEdgeDefaultNodePool(roaClient, d.Id()); nodepoolErr != nil {
			log.Printf("[WARN] reading Edge default node pool CloudMonitor setting failed: %s", nodepoolErr)
		} else {
			if nodepool.KubernetesConfig != nil && nodepool.KubernetesConfig.CmsEnabled != nil {
				d.Set("install_cloud_monitor", tea.BoolValue(nodepool.KubernetesConfig.CmsEnabled))
				cloudMonitorReadFromNodePool = true
			}
			if nodepool.ScalingGroup != nil && nodepool.ScalingGroup.DesiredSize != nil {
				// worker_number is desired capacity, not the count of
				// kubelets that currently happen to be registered.
				d.Set("worker_number", int(tea.Int64Value(nodepool.ScalingGroup.DesiredSize)))
				workerNumberReadFromNodePool = true
			}
			nodepoolId := tea.StringValue(nodepool.NodepoolInfo.NodepoolId)
			snapshotPolicyConfigured := d.Get("worker_disk_snapshot_policy_id").(string) != ""
			if detail, detailErr := describeEdgeNodePoolDetail(roaClient, d.Id(), nodepoolId); detailErr != nil {
				if snapshotPolicyConfigured {
					return WrapError(detailErr)
				}
				log.Printf("[WARN] reading Edge default node pool detail failed: %s", detailErr)
			} else if detail.ScalingGroup == nil {
				if snapshotPolicyConfigured {
					return Error("DescribeClusterNodePoolDetail returned no scaling group for Edge cluster %s node pool %s", d.Id(), nodepoolId)
				}
			} else {
				// The list API returns a summary and omits the snapshot policy.
				// Use the detail API so state reflects the actual backend value
				// and the ESSD acceptance test verifies the SDK v2 mapping.
				if err := d.Set("worker_disk_snapshot_policy_id", tea.StringValue(detail.ScalingGroup.SystemDiskSnapshotPolicyId)); err != nil {
					return WrapError(err)
				}
			}
		}
		if !cloudMonitorReadFromNodePool {
			if v, ok := object.Parameters["CloudMonitorFlags"]; ok {
				d.Set("install_cloud_monitor", Interface2Bool(v))
			}
		}
		// only works with default-nodepool
		workerNodes := fetchWorkerNodes(d, meta)
		d.Set("worker_nodes", workerNodes)
		if !workerNumberReadFromNodePool {
			d.Set("worker_number", len(workerNodes))
		}
	}
	if object.ClusterType == "Kubernetes" {
		if v, ok := object.Parameters["CloudMonitorFlags"]; ok {
			d.Set("install_cloud_monitor", Interface2Bool(v))
		}
		if v, ok := object.Parameters["MasterKeyPair"]; ok {
			d.Set("key_name", Interface2String(v))
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

	d.Set("connections", []interface{}{connection})
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

	if err = setCerts(d, meta, d.Get("skip_set_certificate_authority").(bool), true); err != nil {
		return WrapError(err)
	}

	return nil
}

func describeEdgeWorkerNodes(client *roacs.Client, clusterId, nodepoolId string) ([]*roacs.DescribeClusterNodesResponseBodyNodes, error) {
	nodes := make([]*roacs.DescribeClusterNodesResponseBodyNodes, 0)
	for pageNumber := 1; ; pageNumber++ {
		request := &roacs.DescribeClusterNodesRequest{
			NodepoolId: tea.String(nodepoolId),
			PageNumber: tea.String(strconv.Itoa(pageNumber)),
			PageSize:   tea.String("100"),
			State:      tea.String("all"),
		}
		response, err := client.DescribeClusterNodes(tea.String(clusterId), request)
		if err != nil {
			return nil, WrapErrorf(err, DefaultErrorMsg, clusterId, "DescribeClusterNodes", AlibabaCloudSdkGoERROR)
		}
		if response == nil || response.Body == nil {
			return nil, Error("DescribeClusterNodes returned an empty response for Edge cluster %s node pool %s", clusterId, nodepoolId)
		}
		nodes = append(nodes, response.Body.Nodes...)
		if len(response.Body.Nodes) < 100 {
			break
		}
	}
	return nodes, nil
}

func edgeWorkerNodeSummary(nodes []*roacs.DescribeClusterNodesResponseBodyNodes) string {
	items := make([]map[string]interface{}, 0, len(nodes))
	for _, node := range nodes {
		if node == nil {
			continue
		}
		ips := make([]string, 0, len(node.IpAddress))
		for _, ip := range node.IpAddress {
			ips = append(ips, tea.StringValue(ip))
		}
		items = append(items, map[string]interface{}{
			"instance_id":     tea.StringValue(node.InstanceId),
			"instance_status": tea.StringValue(node.InstanceStatus),
			"node_name":       tea.StringValue(node.NodeName),
			"node_status":     tea.StringValue(node.NodeStatus),
			"state":           tea.StringValue(node.State),
			"error_message":   tea.StringValue(node.ErrorMessage),
			"ip_addresses":    ips,
		})
		if len(items) == 10 {
			break
		}
	}
	payload, err := json.Marshal(items)
	if err != nil {
		return fmt.Sprintf("%v", items)
	}
	return string(payload)
}

// logEdgeKubernetesDiagnostics writes only remote control-plane observations
// to the ACC log. In particular, node error_message exposes bootstrap and
// private API-server connectivity failures without requiring local access.
func logEdgeKubernetesDiagnostics(d *schema.ResourceData, meta interface{}, nodepoolId string, cause error) {
	log.Printf("[WARN] Edge Kubernetes diagnostics: cluster=%s nodepool=%s cause=%s", d.Id(), nodepoolId, cause)
	client := meta.(*connectivity.AliyunClient)
	if cluster, err := (&CsService{client}).DescribeCsKubernetes(d.Id()); err != nil {
		log.Printf("[WARN] Edge Kubernetes diagnostics: DescribeCsKubernetes failed: %s", err)
	} else {
		log.Printf("[WARN] Edge Kubernetes cluster: state=%s vpc=%s vswitches=%s security_group=%s ingress_slb=%s master_url=%s", cluster.State, cluster.VpcId, cluster.VSwitchIds, cluster.SecurityGroupId, cluster.IngressLoadbalancerId, cluster.MasterURL)
	}

	roaClient, err := client.NewRoaCsClient()
	if err != nil {
		log.Printf("[WARN] Edge Kubernetes diagnostics: ACK client initialization failed: %s", err)
		return
	}
	response, err := roaClient.DescribeClusterNodePools(tea.String(d.Id()), &roacs.DescribeClusterNodePoolsRequest{})
	if err != nil {
		log.Printf("[WARN] Edge Kubernetes diagnostics: DescribeClusterNodePools failed: %s", err)
	} else if response == nil || response.Body == nil {
		log.Printf("[WARN] Edge Kubernetes diagnostics: DescribeClusterNodePools returned an empty response")
	} else {
		for _, nodepool := range response.Body.Nodepools {
			if nodepool == nil || nodepool.NodepoolInfo == nil {
				continue
			}
			currentId := tea.StringValue(nodepool.NodepoolInfo.NodepoolId)
			if nodepoolId == "" && tea.StringValue(nodepool.NodepoolInfo.Type) == defaultNodePoolType && (tea.BoolValue(nodepool.NodepoolInfo.IsDefault) || tea.StringValue(nodepool.NodepoolInfo.Name) == "default-nodepool") {
				nodepoolId = currentId
			}
			if currentId != nodepoolId {
				continue
			}
			payload, marshalErr := json.Marshal(edgeNodePoolDiagnostic(nodepool))
			if marshalErr != nil {
				log.Printf("[WARN] Edge Kubernetes diagnostics: node pool summary failed: %s", marshalErr)
			} else {
				log.Printf("[WARN] Edge Kubernetes node pool: %s", string(payload))
			}
		}
	}

	if nodepoolId == "" {
		return
	}
	nodes, err := describeEdgeWorkerNodes(roaClient, d.Id(), nodepoolId)
	if err != nil {
		log.Printf("[WARN] Edge Kubernetes diagnostics: DescribeClusterNodes failed: %s", err)
		return
	}
	log.Printf("[WARN] Edge Kubernetes workers: %s", edgeWorkerNodeSummary(nodes))
}

func edgeNodePoolDiagnostic(nodepool *roacs.DescribeClusterNodePoolsResponseBodyNodepools) map[string]interface{} {
	result := map[string]interface{}{}
	if nodepool.NodepoolInfo != nil {
		result["id"] = tea.StringValue(nodepool.NodepoolInfo.NodepoolId)
		result["name"] = tea.StringValue(nodepool.NodepoolInfo.Name)
		result["type"] = tea.StringValue(nodepool.NodepoolInfo.Type)
		result["is_default"] = tea.BoolValue(nodepool.NodepoolInfo.IsDefault)
	}
	if nodepool.ScalingGroup != nil {
		result["desired_size"] = tea.Int64Value(nodepool.ScalingGroup.DesiredSize)
		result["security_group_id"] = tea.StringValue(nodepool.ScalingGroup.SecurityGroupId)
		result["security_group_ids"] = nodepool.ScalingGroup.SecurityGroupIds
		result["vswitch_ids"] = nodepool.ScalingGroup.VswitchIds
		result["instance_types"] = nodepool.ScalingGroup.InstanceTypes
		result["system_disk_snapshot_policy_id"] = tea.StringValue(nodepool.ScalingGroup.SystemDiskSnapshotPolicyId)
	}
	if nodepool.Status != nil {
		result["state"] = tea.StringValue(nodepool.Status.State)
		result["total_nodes"] = tea.Int64Value(nodepool.Status.TotalNodes)
		result["healthy_nodes"] = tea.Int64Value(nodepool.Status.HealthyNodes)
		result["initial_nodes"] = tea.Int64Value(nodepool.Status.InitialNodes)
		result["failed_nodes"] = tea.Int64Value(nodepool.Status.FailedNodes)
		result["offline_nodes"] = tea.Int64Value(nodepool.Status.OfflineNodes)
	}
	return result
}

func waitForEdgeDefaultNodePool(client *roacs.Client, clusterId, expectedNodepoolId string, expectedSize int, timeout time.Duration) error {
	return resource.Retry(timeout, func() *resource.RetryError {
		nodepool, err := findEdgeDefaultNodePool(client, clusterId)
		if err != nil {
			return resource.RetryableError(err)
		}
		actualNodepoolId := tea.StringValue(nodepool.NodepoolInfo.NodepoolId)
		if expectedNodepoolId != "" && actualNodepoolId != expectedNodepoolId {
			return resource.RetryableError(Error("Edge cluster %s default ESS node pool ID is %s, expected %s", clusterId, actualNodepoolId, expectedNodepoolId))
		}
		if nodepool.ScalingGroup == nil {
			return resource.RetryableError(Error("Edge cluster %s default ESS node pool %s has no scaling group", clusterId, expectedNodepoolId))
		}
		actualSize := int(tea.Int64Value(nodepool.ScalingGroup.DesiredSize))
		if actualSize != expectedSize {
			return resource.RetryableError(Error("Edge cluster %s default ESS node pool %s desired size is %d, expected %d", clusterId, expectedNodepoolId, actualSize, expectedSize))
		}
		// The ACK task and desired size are the resource-level completion
		// contract. Node readiness can remain non-active because of account
		// networking policies, which the SDK v1 resource did not treat as an
		// additional provider failure.
		return nil
	})
}

func findEdgeDefaultNodePool(csClient *roacs.Client, clusterId string) (*roacs.DescribeClusterNodePoolsResponseBodyNodepools, error) {
	var response *roacs.DescribeClusterNodePoolsResponse
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = csClient.DescribeClusterNodePools(tea.String(clusterId), &roacs.DescribeClusterNodePoolsRequest{})
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
		return nil, WrapErrorf(err, DefaultErrorMsg, clusterId, "DescribeClusterNodePools", AlibabaCloudSdkGoERROR)
	}
	if response == nil || response.Body == nil {
		return nil, Error("DescribeClusterNodePools returned an empty response for Edge cluster %s", clusterId)
	}

	defaultPools := make([]*roacs.DescribeClusterNodePoolsResponseBodyNodepools, 0, 1)
	namedPools := make([]*roacs.DescribeClusterNodePoolsResponseBodyNodepools, 0, 1)
	for _, nodepool := range response.Body.Nodepools {
		if nodepool == nil || nodepool.NodepoolInfo == nil || tea.StringValue(nodepool.NodepoolInfo.Type) != defaultNodePoolType {
			continue
		}
		if tea.BoolValue(nodepool.NodepoolInfo.IsDefault) {
			defaultPools = append(defaultPools, nodepool)
		}
		if tea.StringValue(nodepool.NodepoolInfo.Name) == "default-nodepool" {
			namedPools = append(namedPools, nodepool)
		}
	}

	var nodepool *roacs.DescribeClusterNodePoolsResponseBodyNodepools
	switch len(defaultPools) {
	case 1:
		nodepool = defaultPools[0]
	case 0:
		switch len(namedPools) {
		case 1:
			nodepool = namedPools[0]
		case 0:
			return nil, Error("no default ESS node pool was found for Edge cluster %s", clusterId)
		default:
			return nil, Error("found %d ESS node pools named default-nodepool for Edge cluster %s", len(namedPools), clusterId)
		}
	default:
		return nil, Error("found %d default ESS node pools for Edge cluster %s", len(defaultPools), clusterId)
	}

	if tea.StringValue(nodepool.NodepoolInfo.NodepoolId) == "" {
		return nil, Error("the default ESS node pool for Edge cluster %s returned an empty node pool ID", clusterId)
	}
	return nodepool, nil
}

func describeEdgeNodePoolDetail(csClient *roacs.Client, clusterId, nodepoolId string) (*roacs.DescribeClusterNodePoolDetailResponseBody, error) {
	var response *roacs.DescribeClusterNodePoolDetailResponse
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = csClient.DescribeClusterNodePoolDetail(tea.String(clusterId), tea.String(nodepoolId))
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"ErrorNodePoolNotFound"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(fmt.Sprintf("/clusters/%s/nodepools/%s", clusterId, nodepoolId), response, map[string]interface{}{
		"cluster_id":  clusterId,
		"nodepool_id": nodepoolId,
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, clusterId, "DescribeClusterNodePoolDetail", AlibabaCloudSdkGoERROR)
	}
	if response == nil || response.Body == nil {
		return nil, Error("DescribeClusterNodePoolDetail returned an empty response for Edge cluster %s node pool %s", clusterId, nodepoolId)
	}
	return response.Body, nil
}

func scaleOutEdgeKubernetesNodePool(d *schema.ResourceData, meta interface{}, desiredSize int) (string, error) {
	client := meta.(*connectivity.AliyunClient)
	roaClient, err := client.NewRoaCsClient()
	if err != nil {
		return "", WrapErrorf(err, DefaultErrorMsg, d.Id(), "InitializeClient", AlibabaCloudSdkGoERROR)
	}
	csClient := CsClient{client: roaClient}
	nodepool, err := findEdgeDefaultNodePool(csClient.client, d.Id())
	if err != nil {
		return "", err
	}
	nodepoolId := tea.StringValue(nodepool.NodepoolInfo.NodepoolId)

	scalingGroup := &roacs.ModifyClusterNodePoolRequestScalingGroup{
		DesiredSize:        tea.Int64(int64(desiredSize)),
		VswitchIds:         tea.StringSlice(expandStringList(d.Get("worker_vswitch_ids").([]interface{}))),
		InstanceTypes:      tea.StringSlice(expandStringList(d.Get("worker_instance_types").([]interface{}))),
		InstanceChargeType: tea.String(d.Get("worker_instance_charge_type").(string)),
	}
	if v, ok := d.GetOk("password"); ok && v.(string) != "" {
		scalingGroup.SetLoginPassword(v.(string))
	}
	if v, ok := d.GetOk("key_name"); ok && v.(string) != "" {
		scalingGroup.SetKeyPair(v.(string))
	}
	if v, ok := d.GetOk("rds_instances"); ok {
		scalingGroup.SetRdsInstances(tea.StringSlice(expandStringList(v.([]interface{}))))
	}
	if v, ok := d.GetOk("security_group_id"); ok && v.(string) != "" {
		scalingGroup.SetSecurityGroupIds(tea.StringSlice([]string{v.(string)}))
	}
	tags := make([]*roacs.Tag, 0)
	if tagsMap, ok := d.Get("tags").(map[string]interface{}); ok {
		for key, value := range tagsMap {
			if tagValue, ok := value.(string); ok {
				tags = append(tags, &roacs.Tag{Key: tea.String(key), Value: tea.String(tagValue)})
			}
		}
	}
	if len(tags) > 0 {
		scalingGroup.SetTags(tags)
	}
	if v, ok := d.GetOk("worker_disk_category"); ok && v.(string) != "" {
		scalingGroup.SetSystemDiskCategory(v.(string))
	}
	if v, ok := d.GetOk("worker_disk_size"); ok {
		scalingGroup.SetSystemDiskSize(int64(v.(int)))
	}
	if v, ok := d.GetOk("worker_disk_performance_level"); ok && v.(string) != "" {
		scalingGroup.SetSystemDiskPerformanceLevel(v.(string))
	}
	if v, ok := d.GetOk("worker_disk_snapshot_policy_id"); ok && v.(string) != "" {
		scalingGroup.SetSystemDiskSnapshotPolicyId(v.(string))
	}
	dataDisks, err := expandEdgeKubernetesDataDisks(d)
	if err != nil {
		return nodepoolId, err
	}
	if dataDisks != nil {
		scalingGroup.SetDataDisks(dataDisks)
	}

	kubernetesConfig := &roacs.ModifyClusterNodePoolRequestKubernetesConfig{
		CmsEnabled: tea.Bool(d.Get("install_cloud_monitor").(bool)),
	}
	if userData, ok := d.GetOk("user_data"); ok && userData.(string) != "" {
		encoded := userData.(string)
		if _, err := base64.StdEncoding.DecodeString(encoded); err != nil {
			encoded = base64.StdEncoding.EncodeToString([]byte(encoded))
		}
		kubernetesConfig.SetUserData(encoded)
	}
	if runtime, ok := d.GetOk("runtime"); ok {
		if raw := runtime.([]interface{}); len(raw) > 0 {
			if values := raw[0].(map[string]interface{}); len(values) > 0 {
				if name, ok := values["name"].(string); ok && name != "" {
					kubernetesConfig.SetRuntime(name)
				}
				if version, ok := values["version"].(string); ok && version != "" {
					kubernetesConfig.SetRuntimeVersion(version)
				}
			}
		}
	}

	request := &roacs.ModifyClusterNodePoolRequest{
		ScalingGroup:     scalingGroup,
		KubernetesConfig: kubernetesConfig,
	}
	action := fmt.Sprintf("/clusters/%s/nodepools/%s", d.Id(), nodepoolId)
	var response *roacs.ModifyClusterNodePoolResponse
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = csClient.client.ModifyClusterNodePool(tea.String(d.Id()), tea.String(nodepoolId), request)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, map[string]interface{}{"request": request})
	if err != nil {
		return nodepoolId, WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	if response == nil || response.Body == nil {
		return nodepoolId, WrapErrorf(Error("ModifyClusterNodePool returned an empty response"), ResponseCodeMsg, d.Id(), action, response)
	}
	taskId := tea.StringValue(response.Body.TaskId)
	if taskId == "" {
		return nodepoolId, WrapErrorf(Error("ModifyClusterNodePool returned an empty task ID"), ResponseCodeMsg, d.Id(), action, response)
	}

	stateConf := BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, csClient.DescribeTaskRefreshFunc(d, taskId, []string{"fail", "failed"}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return nodepoolId, WrapErrorf(err, ResponseCodeMsg, d.Id(), action, jobDetail)
	}
	if err := waitForEdgeDefaultNodePool(csClient.client, d.Id(), nodepoolId, desiredSize, 10*time.Minute); err != nil {
		return nodepoolId, err
	}
	return nodepoolId, nil
}

func resourceAlicloudCSEdgeKubernetesUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	d.Partial(true)
	invoker := NewInvoker()
	//scale up cloud worker nodes
	if d.HasChanges("worker_number") {
		oldV, newV := d.GetChange("worker_number")
		oldValue, ok := oldV.(int)
		if !ok {
			return WrapErrorf(fmt.Errorf("worker_number old value can not be parsed"), "parseError %d", oldValue)
		}
		newValue, ok := newV.(int)
		if !ok {
			return WrapErrorf(fmt.Errorf("worker_number new value can not be parsed"), "parseError %d", oldValue)
		}

		// Edge cluster node support remove nodes.
		if newValue < oldValue {
			return WrapErrorf(fmt.Errorf("worker_number can not be less than before"), "scaleOutCloudWorkersFailed %d:%d", newValue, oldValue)
		}

		// scale out cluster.
		if newValue > oldValue {
			nodepoolId, err := scaleOutEdgeKubernetesNodePool(d, meta, newValue)
			if err != nil {
				logEdgeKubernetesDiagnostics(d, meta, nodepoolId, err)
				return WrapError(err)
			}
		}
	}
	// modify cluster name
	if !d.IsNewResource() && (d.HasChange("name") || d.HasChange("name_prefix")) {
		var clusterName string
		if v, ok := d.GetOk("name"); ok {
			clusterName = v.(string)
		} else {
			clusterName = resource.PrefixedUniqueId(d.Get("name_prefix").(string))
		}
		var requestInfo *cs.Client
		var response interface{}
		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				requestInfo = csClient
				return nil, csClient.ModifyClusterName(d.Id(), clusterName)
			})
			response = raw
			return err
		}); err != nil && !IsExpectedErrors(err, []string{"ErrorClusterNameAlreadyExist"}) {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifyClusterName", DenverdinoAliyungo)
		}
		if debugOn() {
			requestMap := make(map[string]interface{})
			requestMap["ClusterId"] = d.Id()
			requestMap["ClusterName"] = clusterName
			addDebug("ModifyClusterName", response, requestInfo, requestMap)
		}
	}

	// modify cluster deletion protection
	if !d.IsNewResource() && d.HasChange("deletion_protection") {
		var requestInfo cs.ModifyClusterArgs
		if v, ok := d.GetOk("deletion_protection"); ok {
			requestInfo.DeletionProtection = v.(bool)
		}

		var response interface{}
		if err := invoker.Run(func() error {
			_, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				return nil, csClient.ModifyCluster(d.Id(), &requestInfo)
			})
			return err
		}); err != nil && !IsExpectedErrors(err, []string{"ErrorModifyDeletionProtectionFailed"}) {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifyCluster", DenverdinoAliyungo)
		}
		if debugOn() {
			requestMap := make(map[string]interface{})
			requestMap["ClusterId"] = d.Id()
			requestMap["deletion_protection"] = requestInfo.DeletionProtection
			addDebug("ModifyCluster", response, requestInfo, requestMap)
		}
	}

	// modify cluster tag
	if d.HasChange("tags") {
		err := updateKubernetesClusterTag(d, meta)
		if err != nil {
			return WrapErrorf(err, ResponseCodeMsg, d.Id(), "ModifyClusterTags", AlibabaCloudSdkGoERROR)
		}
	}

	// upgrade cluster version
	err := UpgradeAlicloudKubernetesCluster(d, meta)
	if err != nil {
		return WrapError(err)
	}
	d.Partial(false)
	return resourceAlicloudCSEdgeKubernetesRead(d, meta)
}

func buildEdgeKubernetesRequest(d *schema.ResourceData, meta interface{}) (*roacs.CreateClusterRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	workerVSwitchIds := make([]string, 0)
	if v, ok := d.GetOk("worker_vswitch_ids"); ok {
		workerVSwitchIds = append(workerVSwitchIds, expandStringList(v.([]interface{}))...)
	}
	var vpcId string
	if len(workerVSwitchIds) > 0 {
		vsw, err := vpcService.DescribeVSwitch(workerVSwitchIds[0])
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

	tags := make([]*roacs.Tag, 0)
	nodeTags := make([]*roacs.NodepoolScalingGroupTags, 0)
	if tagsMap, ok := d.Get("tags").(map[string]interface{}); ok {
		for key, value := range tagsMap {
			if value == nil {
				continue
			}
			if tagValue, ok := value.(string); ok {
				tags = append(tags, &roacs.Tag{Key: tea.String(key), Value: tea.String(tagValue)})
				nodeTags = append(nodeTags, &roacs.NodepoolScalingGroupTags{Key: tea.String(key), Value: tea.String(tagValue)})
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

	request := &roacs.CreateClusterRequest{
		Name:                      tea.String(clusterName),
		RegionId:                  tea.String(client.RegionId),
		ClusterType:               tea.String("ManagedKubernetes"),
		Profile:                   tea.String(EdgeProfile),
		Vpcid:                     tea.String(vpcId),
		VswitchIds:                tea.StringSlice(workerVSwitchIds),
		Tags:                      tags,
		Addons:                    addons,
		DeletionProtection:        tea.Bool(d.Get("deletion_protection").(bool)),
		DisableRollback:           tea.Bool(true),
		NodeCidrMask:              tea.String(strconv.Itoa(d.Get("node_cidr_mask").(int))),
		SnatEntry:                 tea.Bool(d.Get("new_nat_gateway").(bool)),
		EndpointPublicAccess:      tea.Bool(d.Get("slb_internet_enabled").(bool)),
		IsEnterpriseSecurityGroup: tea.Bool(d.Get("is_enterprise_security_group").(bool)),
	}
	if v, ok := d.GetOk("version"); ok && v.(string) != "" {
		request.SetKubernetesVersion(v.(string))
	}
	if v, ok := d.GetOk("cluster_spec"); ok && v.(string) != "" {
		request.SetClusterSpec(v.(string))
	}
	if v, ok := d.GetOk("service_cidr"); ok && v.(string) != "" {
		request.SetServiceCidr(v.(string))
	}
	if v, ok := d.GetOk("pod_cidr"); ok && v.(string) != "" {
		request.SetContainerCidr(v.(string))
	}
	if v, ok := d.GetOk("proxy_mode"); ok && v.(string) != "" {
		request.SetProxyMode(v.(string))
	}
	if v, ok := d.GetOk("load_balancer_spec"); ok && v.(string) != "" {
		request.SetLoadBalancerSpec(v.(string))
	}
	if v, ok := d.GetOk("resource_group_id"); ok && v.(string) != "" {
		request.SetResourceGroupId(v.(string))
	}
	if v, ok := d.GetOk("security_group_id"); ok && v.(string) != "" {
		request.SetSecurityGroupId(v.(string))
	}
	if v, ok := d.GetOk("worker_disk_snapshot_policy_id"); ok && v.(string) != "" {
		// CreateCluster's embedded Nodepool model does not expose this field,
		// while the request keeps the legacy compatibility field. Send it so
		// SDK v2 does not silently drop an existing resource argument.
		request.SetWorkerSystemDiskSnapshotPolicyId(v.(string))
	}

	workerInstanceTypes := make([]string, 0)
	if v, ok := d.GetOk("worker_instance_types"); ok {
		workerInstanceTypes = expandStringList(v.([]interface{}))
	}
	scalingGroup := &roacs.NodepoolScalingGroup{
		DesiredSize:        tea.Int64(int64(d.Get("worker_number").(int))),
		VswitchIds:         tea.StringSlice(workerVSwitchIds),
		InstanceTypes:      tea.StringSlice(workerInstanceTypes),
		InstanceChargeType: tea.String(d.Get("worker_instance_charge_type").(string)),
		SystemDiskCategory: tea.String(d.Get("worker_disk_category").(string)),
		SystemDiskSize:     tea.Int64(int64(d.Get("worker_disk_size").(int))),
		// The legacy Edge resource has no system-disk encryption option. Omit
		// the field and let the ECS account-level default encryption policy apply.
		Tags: nodeTags,
	}
	if v, ok := d.GetOk("worker_disk_performance_level"); ok && v.(string) != "" {
		scalingGroup.SetSystemDiskPerformanceLevel(v.(string))
	}
	if v, ok := d.GetOk("password"); ok && v.(string) != "" {
		scalingGroup.SetLoginPassword(v.(string))
	}
	if v, ok := d.GetOk("key_name"); ok && v.(string) != "" {
		scalingGroup.SetKeyPair(v.(string))
	}
	if v, ok := d.GetOk("rds_instances"); ok {
		scalingGroup.SetRdsInstances(tea.StringSlice(expandStringList(v.([]interface{}))))
	}
	if v, ok := d.GetOk("security_group_id"); ok && v.(string) != "" {
		// CreateClusterNodePool treats an explicitly supplied security group as
		// a custom node-pool group and does not add access rules for it. Preserve
		// an explicitly configured group, but use the SDK v2 plural field. The
		// cluster-created default group is added after cluster creation, when its
		// ID is available.
		scalingGroup.SetSecurityGroupIds(tea.StringSlice([]string{v.(string)}))
	}
	dataDisks, err := expandEdgeKubernetesDataDisks(d)
	if err != nil {
		return nil, err
	}
	if len(dataDisks) > 0 {
		scalingGroup.SetDataDisks(dataDisks)
	}

	kubernetesConfig := &roacs.NodepoolKubernetesConfig{
		CmsEnabled: tea.Bool(d.Get("install_cloud_monitor").(bool)),
	}
	runtimeName := ""
	runtimeVersion := ""
	if runtime, ok := d.GetOk("runtime"); ok {
		if raw := runtime.([]interface{}); len(raw) > 0 {
			if v := raw[0].(map[string]interface{}); len(v) > 0 {
				if name, ok := v["name"].(string); ok && name != "" {
					runtimeName = name
				}
				if version, ok := v["version"].(string); ok && version != "" {
					runtimeVersion = version
				}
			}
		}
	}
	// The legacy Edge resource made runtime optional and ACK selected a
	// compatible default. The SDK v2 nodepools request requires both values,
	// so resolve any omitted value from the selected Kubernetes version to
	// preserve that behavior for existing configurations.
	kubernetesVersion := d.Get("version").(string)
	if kubernetesVersion == "" || runtimeName == "" || runtimeVersion == "" {
		selectedName, selectedRuntimeVersion, selectedKubernetesVersion, err := resolveEdgeKubernetesRuntime(client, kubernetesVersion, runtimeName, runtimeVersion)
		if err != nil {
			return nil, err
		}
		if kubernetesVersion == "" {
			// The embedded node-pool request requires an explicit runtime.
			// Pin the Kubernetes version selected by metadata as well, so the
			// runtime can never be chosen for a different backend default.
			request.SetKubernetesVersion(selectedKubernetesVersion)
		}
		runtimeName = selectedName
		runtimeVersion = selectedRuntimeVersion
	}
	kubernetesConfig.SetRuntime(runtimeName)
	kubernetesConfig.SetRuntimeVersion(runtimeVersion)
	if userData, ok := d.GetOk("user_data"); ok && userData.(string) != "" {
		encoded := userData.(string)
		if _, err := base64.StdEncoding.DecodeString(encoded); err != nil {
			encoded = base64.StdEncoding.EncodeToString([]byte(encoded))
		}
		kubernetesConfig.SetUserData(encoded)
	}

	nodepoolInfo := &roacs.NodepoolNodepoolInfo{
		Name: tea.String("default-nodepool"),
		Type: tea.String(defaultNodePoolType),
	}
	if v, ok := d.GetOk("resource_group_id"); ok && v.(string) != "" {
		nodepoolInfo.SetResourceGroupId(v.(string))
	}
	request.SetNodepools([]*roacs.Nodepool{{
		NodepoolInfo:     nodepoolInfo,
		ScalingGroup:     scalingGroup,
		KubernetesConfig: kubernetesConfig,
	}})

	return request, nil
}

func resolveEdgeKubernetesRuntime(client *connectivity.AliyunClient, kubernetesVersion, preferredName, preferredVersion string) (string, string, string, error) {
	roaClient, err := client.NewRoaCsClient()
	if err != nil {
		return "", "", "", WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_edge_kubernetes", "InitializeClient", AlibabaCloudSdkGoERROR)
	}
	request := &roacs.DescribeKubernetesVersionMetadataRequest{
		Region:      tea.String(client.RegionId),
		ClusterType: tea.String("ManagedKubernetes"),
		Profile:     tea.String(EdgeProfile),
	}
	if kubernetesVersion != "" {
		request.SetKubernetesVersion(kubernetesVersion)
	}

	var response *roacs.DescribeKubernetesVersionMetadataResponse
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = roaClient.DescribeKubernetesVersionMetadata(request)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug("DescribeKubernetesVersionMetadata", response, request)
	if err != nil {
		return "", "", "", WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_edge_kubernetes", "DescribeKubernetesVersionMetadata", AlibabaCloudSdkGoERROR)
	}
	if response == nil || len(response.Body) == 0 {
		return "", "", "", Error("no ACK Edge runtime metadata was returned for Kubernetes version %q", kubernetesVersion)
	}

	fallbackName := ""
	fallbackRuntimeVersion := ""
	fallbackKubernetesVersion := ""
	for _, metadata := range response.Body {
		if metadata == nil || (kubernetesVersion != "" && tea.StringValue(metadata.Version) != kubernetesVersion) {
			continue
		}
		metadataKubernetesVersion := tea.StringValue(metadata.Version)
		for _, runtime := range metadata.Runtimes {
			if runtime == nil {
				continue
			}
			name := tea.StringValue(runtime.Name)
			version := tea.StringValue(runtime.Version)
			if name == "" || version == "" {
				continue
			}
			if preferredName != "" && !strings.EqualFold(name, preferredName) {
				continue
			}
			if preferredVersion != "" && version != preferredVersion {
				continue
			}
			if fallbackName == "" {
				fallbackName, fallbackRuntimeVersion = name, version
				fallbackKubernetesVersion = metadataKubernetesVersion
			}
			if preferredName != "" || strings.EqualFold(name, "containerd") {
				return name, version, metadataKubernetesVersion, nil
			}
		}
	}
	if fallbackName != "" {
		return fallbackName, fallbackRuntimeVersion, fallbackKubernetesVersion, nil
	}
	return "", "", "", Error("ACK Edge returned no runtime metadata matching name %q, runtime version %q, Kubernetes version %q", preferredName, preferredVersion, kubernetesVersion)
}

func expandEdgeKubernetesDataDisks(d *schema.ResourceData) ([]*roacs.DataDisk, error) {
	raw, ok := d.GetOk("worker_data_disks")
	if !ok {
		return nil, nil
	}
	disks := raw.([]interface{})
	result := make([]*roacs.DataDisk, 0, len(disks))
	for i, item := range disks {
		pack := item.(map[string]interface{})
		dataDisk := &roacs.DataDisk{}
		if size := strings.TrimSpace(pack["size"].(string)); size != "" {
			parsed, err := strconv.ParseInt(size, 10, 64)
			if err != nil {
				return nil, Error("worker_data_disks.%d.size %q is invalid: %s", i, size, err)
			}
			dataDisk.SetSize(parsed)
		}
		if value := pack["category"].(string); value != "" {
			dataDisk.SetCategory(value)
		}
		if value := pack["snapshot_id"].(string); value != "" {
			dataDisk.SetSnapshotId(value)
		}
		if value := pack["name"].(string); value != "" {
			dataDisk.SetDiskName(value)
		}
		if value := pack["device"].(string); value != "" {
			dataDisk.SetDevice(value)
		}
		if value := pack["kms_key_id"].(string); value != "" {
			dataDisk.SetKmsKeyId(value)
		}
		// Preserve explicit false as well as true. Omitting false changes the
		// user's configuration into the account-level encryption default.
		if value := strings.TrimSpace(pack["encrypted"].(string)); value != "" {
			dataDisk.SetEncrypted(value)
		}
		if value := pack["auto_snapshot_policy_id"].(string); value != "" {
			dataDisk.SetAutoSnapshotPolicyId(value)
		}
		if value := pack["performance_level"].(string); value != "" {
			dataDisk.SetPerformanceLevel(value)
		}
		result = append(result, dataDisk)
	}
	return result, nil
}

func resourceAlicloudCSEdgeKubernetesV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_spec": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"worker_vswitch_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"worker_instance_types": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"worker_number": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"worker_disk_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"worker_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"worker_disk_performance_level": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"worker_disk_snapshot_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"proxy_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"worker_instance_charge_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"worker_data_disks": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"category": {
							Type:     schema.TypeString,
							Optional: true,
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
			},
			"worker_ram_role_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pod_cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// lintignore: S022
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
			},
			"node_cidr_mask": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"new_nat_gateway": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"key_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"install_cloud_monitor": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
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
						},
					},
				},
			},
			"slb_internet_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"load_balancer_spec": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
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
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
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
			// lintignore: S022
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
			"skip_set_certificate_authority": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			// lintignore: S022
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
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"nat_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
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
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"log_config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"project": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			// lintignore: S006
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"retain_resources": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAlicloudCSEdgeKubernetesStateUpgradeV0(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	for _, field := range []string{"runtime", "certificate_authority", "connections"} {
		if v, ok := rawState[field]; ok && v != nil {
			switch val := v.(type) {
			case map[string]interface{}:
				if len(val) > 0 {
					rawState[field] = []interface{}{val}
				} else {
					rawState[field] = []interface{}{}
				}
			}
		}
	}
	return rawState, nil
}
