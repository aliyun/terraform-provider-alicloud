package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"

	"strconv"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/cs"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

const (
	KubernetesClusterNetworkTypeFlannel = "flannel"
	KubernetesClusterNetworkTypeTerway  = "terway"

	KubernetesClusterLoggingTypeSLS       = "SLS"
	KubernetesClusterLoggingTypeLogtailDS = "logtail-ds"
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

		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validateContainerName,
				ConflictsWith: []string{"name_prefix"},
			},
			"name_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "Terraform-Creation",
				ValidateFunc:  validateContainerNamePrefix,
				ConflictsWith: []string{"name"},
			},
			"force_update": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"vswitch_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'vswitch_id' has been deprecated from provider version 1.16.0. New field 'vswitch_ids' replaces it.",
			},
			"vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateContainerVswitchId,
				},
				MaxItems:         3,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"new_nat_gateway": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"master_instance_type": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'master_instance_type' has been deprecated from provider version 1.16.0. New field 'master_instance_types' replaces it.",
			},
			"master_instance_types": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems:         1,
				MaxItems:         3,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"worker_instance_type": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'worker_instance_type' has been deprecated from provider version 1.16.0. New field 'worker_instance_types' replaces it.",
			},
			"worker_instance_types": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems:         1,
				MaxItems:         3,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"worker_number": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "Field 'worker_number' has been deprecated from provider version 1.16.0. New field 'worker_numbers' replaces it.",
			},
			"worker_numbers": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:    schema.TypeInt,
					Default: 3,
				},
				MinItems:         1,
				MaxItems:         3,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"password": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Sensitive:        true,
				ConflictsWith:    []string{"key_name", "kms_encrypted_password"},
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"key_name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Optional:         true,
				ConflictsWith:    []string{"password", "kms_encrypted_password"},
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"kms_encrypted_password": {
				Type:          schema.TypeString,
				ForceNew:      true,
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
				Type:             schema.TypeString,
				ForceNew:         true,
				Optional:         true,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"pod_cidr": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Optional:         true,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"service_cidr": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Optional:         true,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"cluster_network_type": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateFunc:     validateAllowedStringValue([]string{KubernetesClusterNetworkTypeFlannel, KubernetesClusterNetworkTypeTerway}),
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"node_cidr_mask": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				Default:          KubernetesClusterNodeCIDRMasksByDefault,
				ValidateFunc:     validateIntegerInRange(24, 28),
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"log_config": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							ValidateFunc: validateAllowedStringValue([]string{KubernetesClusterLoggingTypeSLS}),
							Required:     true,
						},
						"project": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"enable_ssh": {
				Type:             schema.TypeBool,
				Optional:         true,
				ForceNew:         true,
				Default:          false,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"image_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: imageIdSuppressFunc,
			},
			"master_disk_size": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          40,
				ForceNew:         true,
				ValidateFunc:     validateIntegerInRange(40, 500),
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"master_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  DiskCloudEfficiency,
				ForceNew: true,
				ValidateFunc: validateAllowedStringValue([]string{
					string(DiskCloudEfficiency), string(DiskCloudSSD)}),
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"worker_disk_size": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          40,
				ForceNew:         true,
				ValidateFunc:     validateIntegerInRange(20, 32768),
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"worker_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  DiskCloudEfficiency,
				ValidateFunc: validateAllowedStringValue([]string{
					string(DiskCloudEfficiency), string(DiskCloudSSD)}),
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"worker_data_disk_size": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				Default:          40,
				ValidateFunc:     validateIntegerInRange(20, 32768),
				DiffSuppressFunc: workerDataDiskSizeSuppressFunc,
			},
			"worker_data_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validateAllowedStringValue([]string{
					string(DiskCloudEfficiency), string(DiskCloudSSD)}),
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"master_instance_charge_type": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateFunc:     validateInstanceChargeType,
				Default:          PostPaid,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"master_period_unit": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          Month,
				ValidateFunc:     validateInstanceChargeTypePeriodUnit,
				DiffSuppressFunc: csKubernetesMasterPostPaidDiffSuppressFunc,
			},
			"master_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ValidateFunc:     validateInstanceChargeTypePeriod,
				DiffSuppressFunc: csKubernetesMasterPostPaidDiffSuppressFunc,
			},
			"master_auto_renew": {
				Type:             schema.TypeBool,
				Default:          false,
				Optional:         true,
				DiffSuppressFunc: csKubernetesMasterPostPaidDiffSuppressFunc,
			},
			"master_auto_renew_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ValidateFunc:     validateAllowedIntValue([]int{1, 2, 3, 6, 12}),
				DiffSuppressFunc: csKubernetesMasterPostPaidDiffSuppressFunc,
			},
			"worker_instance_charge_type": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateFunc:     validateInstanceChargeType,
				Default:          PostPaid,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"worker_period_unit": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          Month,
				ValidateFunc:     validateInstanceChargeTypePeriodUnit,
				DiffSuppressFunc: csKubernetesWorkerPostPaidDiffSuppressFunc,
			},
			"worker_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ValidateFunc:     validateInstanceChargeTypePeriod,
				DiffSuppressFunc: csKubernetesWorkerPostPaidDiffSuppressFunc,
			},
			"worker_auto_renew": {
				Type:             schema.TypeBool,
				Default:          false,
				Optional:         true,
				DiffSuppressFunc: csKubernetesWorkerPostPaidDiffSuppressFunc,
			},
			"worker_auto_renew_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ValidateFunc:     validateAllowedIntValue([]int{1, 2, 3, 6, 12}),
				DiffSuppressFunc: csKubernetesWorkerPostPaidDiffSuppressFunc,
			},
			"install_cloud_monitor": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          false,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"is_outdated": {
				Type:             schema.TypeBool,
				Optional:         true,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
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

			// 'version' is a reserved parameter and it just is used to test. No Recommendation to expose it.
			"version": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},

			"nodes": {
				Type:       schema.TypeList,
				Optional:   true,
				Elem:       &schema.Schema{Type: schema.TypeString},
				Deprecated: "Field 'nodes' has been deprecated from provider version 1.9.4. New field 'master_nodes' replaces it.",
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
			"slb_internet_enabled": {
				Type:             schema.TypeBool,
				Optional:         true,
				ForceNew:         true,
				Default:          true,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
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
		},
	}
}

func resourceAlicloudCSKubernetesCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	invoker := NewInvoker()

	if isMultiAZ, err := isMultiAZClusterAndCheck(d); err != nil {
		return WrapError(err)
	} else if isMultiAZ {
		args, err := buildKubernetesMultiAZArgs(d, meta)
		if err != nil {
			return WrapError(err)
		}
		var requestInfo *cs.Client
		var raw interface{}
		if err := invoker.Run(func() error {
			raw, err = client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				requestInfo = csClient
				return csClient.CreateKubernetesMultiAZCluster(common.Region(client.RegionId), args)
			})
			return err
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_kubernetes", "CreateKubernetesMultiAZCluster", DenverdinoAliyungo)
		}
		if debugOn() {
			requestMap := make(map[string]interface{})
			requestMap["RegionId"] = common.Region(client.RegionId)
			requestMap["Args"] = args
			addDebug("CreateKubernetesMultiAZCluster", raw, requestInfo, requestMap)
		}
		cluster, _ := raw.(cs.ClusterCreationResponse)
		d.SetId(cluster.ClusterID)
	} else {
		args, err := buildKubernetesArgs(d, meta)
		if err != nil {
			return WrapError(err)
		}
		var requestInfo *cs.Client
		var raw interface{}
		if err := invoker.Run(func() error {
			raw, err = client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				requestInfo = csClient
				return csClient.CreateKubernetesCluster(common.Region(client.RegionId), args)
			})
			return err
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_kubernetes", "CreateKubernetesCluster", DenverdinoAliyungo)
		}
		if debugOn() {
			requestMap := make(map[string]interface{})
			requestMap["RegionId"] = common.Region(client.RegionId)
			requestMap["Args"] = args
			addDebug("CreateKubernetesCluster", raw, requestInfo, requestMap)
		}
		cluster, _ := raw.(cs.ClusterCreationResponse)
		d.SetId(cluster.ClusterID)

	}

	var requestInfo *cs.Client

	raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
		requestInfo = csClient
		return nil, csClient.WaitForClusterAsyn(d.Id(), cs.Running, 3600)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_kubernetes", "WaitForClusterAsyn", DenverdinoAliyungo)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["ClusterId"] = d.Id()
		requestMap["Status"] = cs.Running
		requestMap["TimeOut"] = 3600
		addDebug("WaitForClusterAsyn", raw, requestInfo, requestMap)
	}
	return resourceAlicloudCSKubernetesUpdate(d, meta)
}

func resourceAlicloudCSKubernetesUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	d.Partial(true)
	invoker := NewInvoker()
	if d.HasChange("worker_numbers") && !d.IsNewResource() {

		workerNumbers := expandIntList(d.Get("worker_numbers").([]interface{}))
		workerInstanceTypes := expandStringList(d.Get("worker_instance_types").([]interface{}))

		password := d.Get("password").(string)
		if password == "" {
			if v := d.Get("kms_encrypted_password").(string); v != "" {
				kmsService := KmsService{client}
				decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
				if err != nil {
					return WrapError(err)
				}
				password = decryptResp.Plaintext
			}
		}

		args := &cs.KubernetesClusterResizeArgs{
			DisableRollback: true,
			TimeoutMins:     60,
			LoginPassword:   password,
		}

		if len(workerNumbers) == 1 {
			args.WorkerInstanceType = workerInstanceTypes[0]
			args.NumOfNodes = int64(workerNumbers[0])
		} else if len(workerNumbers) == 3 {
			args.WorkerInstanceTypeA = workerInstanceTypes[0]
			args.WorkerInstanceTypeB = workerInstanceTypes[1]
			args.WorkerInstanceTypeC = workerInstanceTypes[2]
			args.NumOfNodesA = int64(workerNumbers[0])
			args.NumOfNodesB = int64(workerNumbers[1])
			args.NumOfNodesC = int64(workerNumbers[2])
		}
		var requestInfo *cs.Client
		var resoponse interface{}
		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				requestInfo = csClient
				return nil, csClient.ResizeKubernetesCluster(d.Id(), args)
			})
			resoponse = raw
			return err
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ResizeKubernetesCluster", DenverdinoAliyungo)
		}
		if debugOn() {
			resizeRequestMap := make(map[string]interface{})
			resizeRequestMap["ClusterId"] = d.Id()
			resizeRequestMap["Args"] = args
			addDebug("ResizeKubernetesCluster", resoponse, requestInfo, resizeRequestMap)
		}
		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				requestInfo = csClient
				return nil, csClient.WaitForClusterAsyn(d.Id(), cs.Running, 3600)
			})
			resoponse = raw
			return err
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "WaitForClusterAsyn", DenverdinoAliyungo)
		}
		if debugOn() {
			waitRequestMap := make(map[string]interface{})
			waitRequestMap["ClusterId"] = d.Id()
			waitRequestMap["Status"] = cs.Running
			waitRequestMap["TimeOut"] = 3600
			addDebug("WaitForClusterAsyn", resoponse, requestInfo, waitRequestMap)
		}
		d.SetPartial("worker_number")
	}

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
		}); err != nil && !IsExceptedError(err, ErrorClusterNameAlreadyExist) {
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
	d.Partial(false)

	return resourceAlicloudCSKubernetesRead(d, meta)
}

func resourceAlicloudCSKubernetesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	invoker := NewInvoker()
	object, err := csService.DescribeCsKubernetes(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.Name)
	d.Set("vpc_id", object.VPCID)
	d.Set("security_group_id", object.SecurityGroupID)
	d.Set("availability_zone", object.ZoneId)

	var masterNodes []map[string]interface{}
	var workerNodes []map[string]interface{}

	pageNumber := 1
	for {
		var result []cs.KubernetesNodeType
		var pagination *cs.PaginationResult
		var requestInfo *cs.Client
		var response interface{}
		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				requestInfo = csClient
				nodes, paginationResult, err := csClient.GetKubernetesClusterNodes(d.Id(), common.Pagination{PageNumber: pageNumber, PageSize: PageSizeLarge})
				return []interface{}{nodes, paginationResult}, err
			})
			response = raw
			return err
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetKubernetesClusterNodes", DenverdinoAliyungo)
		}
		if debugOn() {
			requestMap := make(map[string]interface{})
			requestMap["ClusterId"] = d.Id()
			requestMap["Pagination"] = common.Pagination{PageNumber: pageNumber, PageSize: PageSizeLarge}
			addDebug("GetKubernetesClusterNodes", response, requestInfo, requestMap)
		}
		result, _ = response.([]interface{})[0].([]cs.KubernetesNodeType)
		pagination, _ = response.([]interface{})[1].(*cs.PaginationResult)

		if pageNumber == 1 && (len(result) == 0 || result[0].InstanceId == "") {
			err := resource.Retry(5*time.Minute, func() *resource.RetryError {
				if err := invoker.Run(func() error {
					raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
						requestInfo = csClient
						nodes, _, err := csClient.GetKubernetesClusterNodes(d.Id(), common.Pagination{PageNumber: pageNumber, PageSize: PageSizeLarge})
						return nodes, err
					})
					response = raw
					return err
				}); err != nil {
					return resource.NonRetryableError(err)
				}
				tmp, _ := response.([]cs.KubernetesNodeType)
				if len(tmp) > 0 && tmp[0].InstanceId != "" {
					result = tmp
				}

				for _, stableState := range cs.NodeStableClusterState {
					// If cluster is in NodeStableClusteState, node list will not change
					if object.State == stableState {
						if debugOn() {
							requestMap := make(map[string]interface{})
							requestMap["ClusterId"] = d.Id()
							requestMap["Pagination"] = common.Pagination{PageNumber: pageNumber, PageSize: PageSizeLarge}
							addDebug("GetKubernetesClusterNodes", response, requestInfo, requestMap)
						}
						return nil
					}
				}
				time.Sleep(5 * time.Second)
				return resource.RetryableError(Error("[ERROR] There is no any nodes in kubernetes cluster %s.", d.Id()))
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetKubernetesClusterNodes", DenverdinoAliyungo)
			}

		}

		for _, node := range result {
			mapping := map[string]interface{}{
				"id":         node.InstanceId,
				"name":       node.InstanceName,
				"private_ip": node.IpAddress[0],
			}
			if node.InstanceRole == "Master" {
				masterNodes = append(masterNodes, mapping)
			} else {
				workerNodes = append(workerNodes, mapping)
			}
		}

		if len(result) < pagination.PageSize {
			break
		}
		pageNumber += 1
	}
	d.Set("master_nodes", masterNodes)
	d.Set("worker_nodes", workerNodes)

	// Get slb information
	connection := make(map[string]string)
	request := slb.CreateDescribeLoadBalancersRequest()
	request.ServerId = masterNodes[0]["id"].(string)
	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeLoadBalancers(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	lbs, _ := raw.(*slb.DescribeLoadBalancersResponse)
	for _, lb := range lbs.LoadBalancers.LoadBalancer {
		if strings.ToLower(lb.AddressType) == strings.ToLower(string(Internet)) {
			d.Set("slb_internet", lb.LoadBalancerId)
			connection["api_server_internet"] = fmt.Sprintf("https://%s:6443", lb.Address)
			connection["master_public_ip"] = lb.Address
		} else {
			d.Set("slb_intranet", lb.LoadBalancerId)
			connection["api_server_intranet"] = fmt.Sprintf("https://%s:6443", lb.Address)

			reqVpc := vpc.CreateDescribeEipAddressesRequest()
			reqVpc.AssociatedInstanceId = lb.LoadBalancerId
			reqVpc.AssociatedInstanceType = "SlbInstance"
			raw, err = client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeEipAddresses(reqVpc)
			})
			eip, _ := raw.(*vpc.DescribeEipAddressesResponse)
			if eip != nil && len(eip.EipAddresses.EipAddress) > 0 {
				eipAddr := eip.EipAddresses.EipAddress[0].IpAddress
				connection["master_public_ip"] = eipAddr
				connection["api_server_internet"] = fmt.Sprintf("https://%s:6443", eipAddr)
			}
		}
	}
	connection["service_domain"] = fmt.Sprintf("*.%s.%s.alicontainer.com", d.Id(), object.RegionID)

	d.Set("connections", connection)
	natRequest := vpc.CreateDescribeNatGatewaysRequest()
	natRequest.VpcId = object.VPCID
	raw, err = client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeNatGateways(natRequest)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	nat, _ := raw.(*vpc.DescribeNatGatewaysResponse)
	if nat != nil && len(nat.NatGateways.NatGateway) > 0 {
		d.Set("nat_gateway_id", nat.NatGateways.NatGateway[0].NatGatewayId)
	}

	var requestInfo *cs.Client
	var response interface{}
	if err := invoker.Run(func() error {
		raw, err = client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			requestInfo = csClient
			return csClient.GetClusterCerts(d.Id())
		})
		response = raw
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetClusterCerts", DenverdinoAliyungo)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["Id"] = d.Id()
		addDebug("GetClusterCerts", response, requestInfo, requestMap)
	}
	cert, _ := response.(cs.ClusterCerts)
	if ce, ok := d.GetOk("client_cert"); ok && ce.(string) != "" {
		if err := writeToFile(ce.(string), cert.Cert); err != nil {
			return WrapError(err)
		}
	}
	if key, ok := d.GetOk("client_key"); ok && key.(string) != "" {
		if err := writeToFile(key.(string), cert.Key); err != nil {
			return WrapError(err)
		}
	}
	if ca, ok := d.GetOk("cluster_ca_cert"); ok && ca.(string) != "" {
		if err := writeToFile(ca.(string), cert.CA); err != nil {
			return WrapError(err)
		}
	}

	var config cs.ClusterConfig
	if file, ok := d.GetOk("kube_config"); ok && file.(string) != "" {
		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				requestInfo = csClient
				return csClient.GetClusterConfig(d.Id())
			})
			response = raw
			return err
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetClusterConfig", DenverdinoAliyungo)
		}
		if debugOn() {
			requestMap := make(map[string]interface{})
			requestMap["Id"] = d.Id()
			addDebug("GetClusterConfig", response, requestInfo, requestMap)
		}
		config, _ = response.(cs.ClusterConfig)

		if err := writeToFile(file.(string), config.Config); err != nil {
			return WrapError(err)
		}
	}

	return nil
}

func resourceAlicloudCSKubernetesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	invoker := NewInvoker()

	var requestInfo *cs.Client
	var response interface{}
	err := resource.Retry(30*time.Minute, func() *resource.RetryError {
		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				requestInfo = csClient
				return nil, csClient.DeleteCluster(d.Id())
			})
			response = raw
			return err
		}); err != nil {
			return resource.RetryableError(err)
		}
		if debugOn() {
			requestMap := make(map[string]interface{})
			requestMap["ClusterId"] = d.Id()
			addDebug("DeleteCluster", response, requestInfo, requestMap)
		}
		return nil
	})
	if err != nil {
		if NotFoundError(err) || IsExceptedError(err, ErrorClusterNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteCluster", DenverdinoAliyungo)
	}
	return WrapError(csService.WaitForCsKubernetes(d.Id(), Deleted, DefaultLongTimeout))
}

func isMultiAZClusterAndCheck(d *schema.ResourceData) (bool, error) {
	masterInstanceTypes := expandStringList(d.Get("master_instance_types").([]interface{}))
	workerInstanceTypes := expandStringList(d.Get("worker_instance_types").([]interface{}))
	vswitchIDs := expandStringList(d.Get("vswitch_ids").([]interface{}))
	workerNumbers := expandIntList(d.Get("worker_numbers").([]interface{}))

	if len(masterInstanceTypes) != len(workerInstanceTypes) {
		return false, WrapError(Error("length of fields `master_instance_types`, `worker_instance_types` must match"))
	}
	if len(masterInstanceTypes) == 1 {
		// single AZ
		return false, nil
	} else if len(masterInstanceTypes) == 3 {
		if len(vswitchIDs) != 3 || len(workerNumbers) != 3 {
			return true, WrapError(Error("length of fields `vswitch_ids`, `worker_numbers` must be 3 for multiAZ clusters"))

		} else {
			return true, nil
		}
	} else {
		return false, WrapError(Error("length of fields `master_instance_types`, `worker_instance_types` should be either 3 (for MultiAZ) or 1 (for Single AZ)"))
	}
}

func buildKubernetesArgs(d *schema.ResourceData, meta interface{}) (*cs.KubernetesCreationArgs, error) {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	vpcService := VpcService{client}

	// Ensure instance_type is valid
	zoneId, validZones, err := ecsService.DescribeAvailableResources(d, meta, InstanceTypeResource)
	if err != nil {
		return nil, err
	}

	var masterInstanceType, workerInstanceType, vswitchID string
	var workerNumber int

	masterInstanceType = expandStringList(d.Get("master_instance_types").([]interface{}))[0]
	workerInstanceType = expandStringList(d.Get("worker_instance_types").([]interface{}))[0]

	if list := expandStringList(d.Get("vswitch_ids").([]interface{})); len(list) > 0 {
		vswitchID = list[0]
	} else {
		vswitchID = ""
	}

	if list := expandIntList(d.Get("worker_numbers").([]interface{})); len(list) > 0 {
		workerNumber = list[0]
	} else {
		workerNumber = 3
	}

	if err := ecsService.InstanceTypeValidation(masterInstanceType, zoneId, validZones); err != nil {
		return nil, err
	}

	if err := ecsService.InstanceTypeValidation(workerInstanceType, zoneId, validZones); err != nil {
		return nil, err
	}

	var clusterName string
	if v, ok := d.GetOk("name"); ok {
		clusterName = v.(string)
	} else {
		clusterName = resource.PrefixedUniqueId(d.Get("name_prefix").(string))
	}

	var vpcId string
	if vswitchID != "" {
		vsw, err := vpcService.DescribeVSwitch(vswitchID)
		if err != nil {
			return nil, err
		}
		vpcId = vsw.VpcId
		if zoneId != "" && zoneId != vsw.ZoneId {
			return nil, WrapError(Error("The specified vswitch %s isn't in the zone %s.", vsw.VSwitchId, zoneId))
		}
		zoneId = vsw.ZoneId
	} else if !d.Get("new_nat_gateway").(bool) {
		return nil, WrapError(Error("The automatic created VPC and VSwitch must set 'new_nat_gateway' to 'true'."))
	}

	loggingType, slsProjectName, err := parseKubernetesClusterLogConfig(d)
	if err != nil {
		return nil, WrapError(err)
	}

	password := d.Get("password").(string)
	if password == "" {
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return nil, WrapError(err)
			}
			password = decryptResp.Plaintext
		}
	}

	creationArgs := &cs.KubernetesCreationArgs{
		Name:                     clusterName,
		ClusterType:              "Kubernetes",
		DisableRollback:          true,
		TimeoutMins:              60,
		MasterInstanceType:       masterInstanceType,
		WorkerInstanceType:       workerInstanceType,
		VPCID:                    vpcId,
		VSwitchId:                vswitchID,
		LoginPassword:            password,
		KeyPair:                  d.Get("key_name").(string),
		ImageId:                  d.Get("image_id").(string),
		Network:                  d.Get("cluster_network_type").(string),
		NodeCIDRMask:             strconv.Itoa(d.Get("node_cidr_mask").(int)),
		LoggingType:              loggingType,
		SLSProjectName:           slsProjectName,
		NumOfNodes:               int64(workerNumber),
		MasterSystemDiskCategory: ecs.DiskCategory(d.Get("master_disk_category").(string)),
		MasterSystemDiskSize:     int64(d.Get("master_disk_size").(int)),
		WorkerSystemDiskCategory: ecs.DiskCategory(d.Get("worker_disk_category").(string)),
		WorkerSystemDiskSize:     int64(d.Get("worker_disk_size").(int)),
		SNatEntry:                d.Get("new_nat_gateway").(bool),
		KubernetesVersion:        d.Get("version").(string),
		ContainerCIDR:            d.Get("pod_cidr").(string),
		ServiceCIDR:              d.Get("service_cidr").(string),
		SSHFlags:                 d.Get("enable_ssh").(bool),
		PublicSLB:                d.Get("slb_internet_enabled").(bool),
		CloudMonitorFlags:        d.Get("install_cloud_monitor").(bool),
		ZoneId:                   zoneId,
	}

	if v, ok := d.GetOk("worker_data_disk_category"); ok {
		creationArgs.WorkerDataDiskCategory = v.(string)
		creationArgs.WorkerDataDisk = true
		creationArgs.WorkerDataDiskSize = int64(d.Get("worker_data_disk_size").(int))
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

	if v, ok := d.GetOk("worker_instance_charge_type"); ok {
		creationArgs.WorkerInstanceChargeType = v.(string)
		if creationArgs.WorkerInstanceChargeType == string(PrePaid) {
			creationArgs.WorkerAutoRenew = d.Get("worker_auto_renew").(bool)
			creationArgs.WorkerAutoRenewPeriod = d.Get("worker_auto_renew_period").(int)
			creationArgs.WorkerPeriod = d.Get("worker_period").(int)
			creationArgs.WorkerPeriodUnit = d.Get("worker_period_unit").(string)
		}
	}

	if userCa, ok := d.GetOk("user_ca"); ok {
		userCaContent, err := loadFileContent(userCa.(string))
		if err != nil {
			return nil, fmt.Errorf("reading user_ca file failed %s", err)
		}
		creationArgs.UserCA = string(userCaContent)
	}

	return creationArgs, nil
}

func buildKubernetesMultiAZArgs(d *schema.ResourceData, meta interface{}) (*cs.KubernetesMultiAZCreationArgs, error) {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	vpcService := VpcService{client}

	// Ensure instance_type is valid
	zoneId, validZones, err := ecsService.DescribeAvailableResources(d, meta, InstanceTypeResource)
	if err != nil {
		return nil, WrapError(err)
	}
	instanceTypes := expandStringList(d.Get("master_instance_types").([]interface{}))
	instanceTypes = append(instanceTypes, expandStringList(d.Get("worker_instance_types").([]interface{}))...)

	for _, instanceType := range instanceTypes {
		if err := ecsService.InstanceTypeValidation(instanceType, zoneId, validZones); err != nil {
			return nil, WrapError(err)
		}
	}

	var clusterName string
	if v, ok := d.GetOk("name"); ok {
		clusterName = v.(string)
	} else {
		clusterName = resource.PrefixedUniqueId(d.Get("name_prefix").(string))
	}

	masterInstanceTypes := expandStringList(d.Get("master_instance_types").([]interface{}))
	workerInstanceTypes := expandStringList(d.Get("worker_instance_types").([]interface{}))
	vswitchIDs := expandStringList(d.Get("vswitch_ids").([]interface{}))
	workerNumbers := expandIntList(d.Get("worker_numbers").([]interface{}))

	vsw, err := vpcService.DescribeVSwitch(vswitchIDs[0])
	if err != nil {
		return nil, WrapError(err)
	}

	loggingType, slsProjectName, err := parseKubernetesClusterLogConfig(d)
	if err != nil {
		return nil, WrapError(err)
	}

	password := d.Get("password").(string)
	if password == "" {
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return nil, WrapError(err)
			}
			password = decryptResp.Plaintext
		}
	}

	creationArgs := &cs.KubernetesMultiAZCreationArgs{
		Name:                     clusterName,
		ClusterType:              "Kubernetes",
		DisableRollback:          true,
		TimeoutMins:              60,
		MultiAZ:                  true,
		MasterInstanceTypeA:      masterInstanceTypes[0],
		MasterInstanceTypeB:      masterInstanceTypes[1],
		MasterInstanceTypeC:      masterInstanceTypes[2],
		WorkerInstanceTypeA:      workerInstanceTypes[0],
		WorkerInstanceTypeB:      workerInstanceTypes[1],
		WorkerInstanceTypeC:      workerInstanceTypes[2],
		LoginPassword:            password,
		KeyPair:                  d.Get("key_name").(string),
		VSwitchIdA:               vswitchIDs[0],
		VSwitchIdB:               vswitchIDs[1],
		VSwitchIdC:               vswitchIDs[2],
		NumOfNodesA:              int64(workerNumbers[0]),
		NumOfNodesB:              int64(workerNumbers[1]),
		NumOfNodesC:              int64(workerNumbers[2]),
		VPCID:                    vsw.VpcId,
		Network:                  d.Get("cluster_network_type").(string),
		NodeCIDRMask:             strconv.Itoa(d.Get("node_cidr_mask").(int)),
		LoggingType:              loggingType,
		SLSProjectName:           slsProjectName,
		ImageId:                  d.Get("image_id").(string),
		MasterSystemDiskCategory: ecs.DiskCategory(d.Get("master_disk_category").(string)),
		MasterSystemDiskSize:     int64(d.Get("master_disk_size").(int)),
		WorkerSystemDiskCategory: ecs.DiskCategory(d.Get("worker_disk_category").(string)),
		WorkerSystemDiskSize:     int64(d.Get("worker_disk_size").(int)),
		ContainerCIDR:            d.Get("pod_cidr").(string),
		ServiceCIDR:              d.Get("service_cidr").(string),
		SSHFlags:                 d.Get("enable_ssh").(bool),
		PublicSLB:                d.Get("slb_internet_enabled").(bool),
		CloudMonitorFlags:        d.Get("install_cloud_monitor").(bool),
		KubernetesVersion:        d.Get("version").(string),
	}

	if v, ok := d.GetOk("worker_data_disk_category"); ok {
		creationArgs.WorkerDataDiskCategory = v.(string)
		creationArgs.WorkerDataDisk = true
		creationArgs.WorkerDataDiskSize = int64(d.Get("worker_data_disk_size").(int))
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

	if v, ok := d.GetOk("worker_instance_charge_type"); ok {
		creationArgs.WorkerInstanceChargeType = v.(string)
		if creationArgs.WorkerInstanceChargeType == string(PrePaid) {
			creationArgs.WorkerAutoRenew = d.Get("worker_auto_renew").(bool)
			creationArgs.WorkerAutoRenewPeriod = d.Get("worker_auto_renew_period").(int)
			creationArgs.WorkerPeriod = d.Get("worker_period").(int)
			creationArgs.WorkerPeriodUnit = d.Get("worker_period_unit").(string)
		}
	}

	if userCa, ok := d.GetOk("user_ca"); ok {
		userCaContent, err := loadFileContent(userCa.(string))
		if err != nil {
			return nil, fmt.Errorf("reading user_ca file failed %s", err)
		}
		creationArgs.UserCA = string(userCaContent)
	}

	return creationArgs, nil
}

func parseKubernetesClusterLogConfig(d *schema.ResourceData) (string, string, error) {
	var loggingType, slsProjectName string

	if v, ok := d.GetOk("log_config"); ok {
		configs := v.([]interface{})
		config, ok := configs[0].(map[string]interface{})
		if ok && config != nil {
			loggingType = config["type"].(string)
			switch loggingType {
			case KubernetesClusterLoggingTypeSLS, KubernetesClusterLoggingTypeLogtailDS:
				if config["project"].(string) == "" {
					return "", "", WrapError(Error("SLS project name must be provided when choosing SLS as log_config."))
				}
				if config["project"].(string) == "None" {
					return "", "", WrapError(Error("SLS project name must not be `None`."))
				}
				slsProjectName = config["project"].(string)
				//rename log controller name
				loggingType = KubernetesClusterLoggingTypeLogtailDS
				break
			default:
				break
			}
		}
	}
	return loggingType, slsProjectName, nil
}
