package alicloud

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/cs"
	aliyungoecs "github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	EdgeKubernetesDefaultTimeoutInMinutes = 60
	EdgeProfile                           = "Edge"
)

func resourceAlicloudCSEdgeKubernetes() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSEdgeKubernetesCreate,
		Read:   resourceAlicloudCSKubernetesRead, // TODO Refactor read from k8s resources
		Update: resourceAlicloudCSEdgeKubernetesUpdate,
		Delete: resourceAlicloudCSKubernetesDelete, // TODO Refactor delete from k8s resources
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: "This resource has been deprecated since v1.276.0 and will be removed in the future. Please use 'alicloud_cs_managed_kubernetes' instead.",

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(EdgeKubernetesDefaultTimeoutInMinutes * time.Minute),
			Update: schema.DefaultTimeout(EdgeKubernetesDefaultTimeoutInMinutes * time.Minute),
			Delete: schema.DefaultTimeout(EdgeKubernetesDefaultTimeoutInMinutes * time.Minute),
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
	invoker := NewInvoker()
	csService := CsService{client}
	args, err := buildEdgeKubernetesArgs(d, meta)
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
			args.Profile = EdgeProfile
			return csClient.CreateManagedKubernetesCluster(&cs.ManagedKubernetesClusterCreationRequest{
				ClusterArgs: args.ClusterArgs,
				WorkerArgs:  args.WorkerArgs,
			})
		})
		response = raw
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_edge_kubernetes", "CreateKubernetesCluster", response)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["RegionId"] = common.Region(client.RegionId)
		requestMap["Args"] = args
		addDebug("CreateKubernetesCluster", response, requestInfo, requestMap)
	}
	cluster, _ := response.(*cs.ClusterCommonResponse)
	d.SetId(cluster.ClusterID)

	stateConf := BuildStateConf([]string{"initial"}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 10*time.Minute, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudCSKubernetesRead(d, meta)
}

func resourceAlicloudCSEdgeKubernetesUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	d.Partial(true)
	invoker := NewInvoker()
	//scale up cloud worker nodes
	var resp interface{}
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
			password := d.Get("password").(string)
			keyPair := d.Get("key_name").(string)

			args := &cs.ScaleOutKubernetesClusterRequest{
				KeyPair:             keyPair,
				LoginPassword:       password,
				Count:               int64(newValue) - int64(oldValue),
				WorkerVSwitchIds:    expandStringList(d.Get("worker_vswitch_ids").([]interface{})),
				WorkerInstanceTypes: expandStringList(d.Get("worker_instance_types").([]interface{})),
			}

			if userData, ok := d.GetOk("user_data"); ok {
				_, base64DecodeError := base64.StdEncoding.DecodeString(userData.(string))
				if base64DecodeError == nil {
					args.UserData = userData.(string)
				} else {
					args.UserData = base64.StdEncoding.EncodeToString([]byte(userData.(string)))
				}
			}

			if imageID, ok := d.GetOk("image_id"); ok {
				args.ImageId = imageID.(string)

			}

			if v, ok := d.GetOk("worker_disk_category"); ok {
				args.WorkerSystemDiskCategory = aliyungoecs.DiskCategory(v.(string))
			}

			if v, ok := d.GetOk("worker_disk_size"); ok {
				args.WorkerSystemDiskSize = int64(v.(int))
			}

			if v, ok := d.GetOk("worker_disk_snapshot_policy_id"); ok {
				args.WorkerSnapshotPolicyId = v.(string)
			}

			if v, ok := d.GetOk("worker_disk_performance_level"); ok {
				args.WorkerSystemDiskPerformanceLevel = v.(string)
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
				args.WorkerDataDisks = createDataDisks

			}

			if d.HasChange("tags") && !d.IsNewResource() {
				if tags, err := ConvertCsTags(d); err == nil {
					args.Tags = tags
				}
				d.SetPartial("tags")
			}

			if err := invoker.Run(func() error {
				var err error
				resp, err = client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
					resp, err := csClient.ScaleOutKubernetesCluster(d.Id(), args)
					return resp, err
				})
				return err
			}); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ScaleOutCloudWorkers", DenverdinoAliyungo)
			}
			if debugOn() {
				resizeRequestMap := make(map[string]interface{})
				resizeRequestMap["ClusterId"] = d.Id()
				resizeRequestMap["Args"] = args
				addDebug("ResizeKubernetesCluster", resp, resizeRequestMap)
			}
			stateConf := BuildStateConf([]string{"scaling"}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))

			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
			d.SetPartial("worker_data_disks")
			d.SetPartial("worker_number")
			d.SetPartial("worker_disk_category")
			d.SetPartial("worker_disk_size")
			d.SetPartial("worker_disk_snapshot_policy_id")
			d.SetPartial("worker_disk_performance_level")
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
		d.SetPartial("name")
		d.SetPartial("name_prefix")
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
		d.SetPartial("deletion_protection")
	}

	// modify cluster tag
	if d.HasChange("tags") {
		err := updateKubernetesClusterTag(d, meta)
		if err != nil {
			return WrapErrorf(err, ResponseCodeMsg, d.Id(), "ModifyClusterTags", AlibabaCloudSdkGoERROR)
		}
	}
	d.SetPartial("tags")

	// upgrade cluster version
	err := UpgradeAlicloudKubernetesCluster(d, meta)
	if err != nil {
		return WrapError(err)
	}
	d.Partial(false)
	return resourceAlicloudCSKubernetesRead(d, meta)
}

func buildEdgeKubernetesArgs(d *schema.ResourceData, meta interface{}) (*cs.DelicatedKubernetesClusterCreationRequest, error) {
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

	if rdsInstances, ok := d.GetOk("rds_instances"); ok {
		creationArgs.RdsInstances = expandStringList(rdsInstances.([]interface{}))
	}

	if nodePortRange, ok := d.GetOk("node_port_range"); ok {
		creationArgs.NodePortRange = nodePortRange.(string)
	}

	if runtime, ok := d.GetOk("runtime"); ok {
		if raw := runtime.([]interface{}); len(raw) > 0 {
			if v := raw[0].(map[string]interface{}); len(v) > 0 {
				creationArgs.Runtime = expandKubernetesRuntimeConfig(v)
			}
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
