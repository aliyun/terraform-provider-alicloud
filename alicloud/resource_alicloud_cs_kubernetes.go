package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"

	"strings"

	"errors"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
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

	KubernetesClusterLoggingTypeSLS = "SLS"
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
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateContainerVswitchId,
				},
				MaxItems: 3,
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
				MinItems: 1,
				MaxItems: 3,
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
				MinItems: 1,
				MaxItems: 3,
			},
			"worker_number": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "Field 'worker_number' has been deprecated from provider version 1.16.0. New field 'worker_numbers' replaces it.",
			},

			"worker_numbers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:    schema.TypeInt,
					Default: 3,
				},
				MinItems: 1,
				MaxItems: 3,
			},
			"password": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Sensitive:     true,
				ConflictsWith: []string{"key_name"},
			},
			"key_name": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				ConflictsWith: []string{"password"},
			},
			"pod_cidr": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"service_cidr": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"cluster_network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{KubernetesClusterNetworkTypeFlannel, KubernetesClusterNetworkTypeTerway}),
			},
			"node_cidr_mask": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      KubernetesClusterNodeCIDRMasksByDefault,
				ValidateFunc: validateIntegerInRange(24, 28),
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
			},
			"enable_ssh": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"image_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: imageIdSuppressFunc,
			},
			"master_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      40,
				ForceNew:     true,
				ValidateFunc: validateIntegerInRange(40, 500),
			},
			"master_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  DiskCloudEfficiency,
				ForceNew: true,
				ValidateFunc: validateAllowedStringValue([]string{
					string(DiskCloudEfficiency), string(DiskCloudSSD)}),
			},
			"worker_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      40,
				ForceNew:     true,
				ValidateFunc: validateIntegerInRange(20, 32768),
			},
			"worker_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  DiskCloudEfficiency,
				ValidateFunc: validateAllowedStringValue([]string{
					string(DiskCloudEfficiency), string(DiskCloudSSD)}),
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
			},
			"master_instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateInstanceChargeType,
				Default:      PostPaid,
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
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateInstanceChargeType,
				Default:      PostPaid,
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
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"is_outdated": {
				Type:     schema.TypeBool,
				Optional: true,
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
				Type:     schema.TypeString,
				Optional: true,
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
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
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
		return err
	} else if isMultiAZ {
		args, err := buildKubernetesMultiAZArgs(d, meta)
		if err != nil {
			return err
		}
		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				return csClient.CreateKubernetesMultiAZCluster(common.Region(client.RegionId), args)
			})
			if err != nil {
				return err
			}
			cluster, _ := raw.(cs.ClusterCreationResponse)
			d.SetId(cluster.ClusterID)
			return nil
		}); err != nil {
			return fmt.Errorf("Creating Kubernetes Cluster got an error: %#v", err)
		}
	} else {
		args, err := buildKubernetesArgs(d, meta)
		if err != nil {
			return err
		}
		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				return csClient.CreateKubernetesCluster(common.Region(client.RegionId), args)
			})
			if err != nil {
				return err
			}
			cluster, _ := raw.(cs.ClusterCreationResponse)
			d.SetId(cluster.ClusterID)
			return nil
		}); err != nil {
			return fmt.Errorf("Creating Kubernetes Cluster got an error: %#v", err)
		}
	}

	if err := invoker.Run(func() error {
		_, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return nil, csClient.WaitForClusterAsyn(d.Id(), cs.Running, 3600)
		})
		return err
	}); err != nil {
		return fmt.Errorf("Waitting for kubernetes cluster %#v got an error: %#v", cs.Running, err)
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

		args := &cs.KubernetesClusterResizeArgs{
			DisableRollback: true,
			TimeoutMins:     60,
			LoginPassword:   d.Get("password").(string),
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
		if err := invoker.Run(func() error {
			_, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				return nil, csClient.ResizeKubernetesCluster(d.Id(), args)
			})
			return err
		}); err != nil {
			return fmt.Errorf("Resize Cluster got an error: %#v", err)
		}

		if err := invoker.Run(func() error {
			_, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				return nil, csClient.WaitForClusterAsyn(d.Id(), cs.Running, 3600)
			})
			return err
		}); err != nil {
			return fmt.Errorf("Waitting for container Cluster %#v got an error: %#v", cs.Running, err)
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
		if err := invoker.Run(func() error {
			_, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				return nil, csClient.ModifyClusterName(d.Id(), clusterName)
			})
			if err != nil && !IsExceptedError(err, ErrorClusterNameAlreadyExist) {
				return err
			}
			return nil
		}); err != nil {
			return fmt.Errorf("Modify Cluster Name got an error: %#v", err)
		}
		d.SetPartial("name")
		d.SetPartial("name_prefix")
	}
	d.Partial(false)

	return resourceAlicloudCSKubernetesRead(d, meta)
}

func resourceAlicloudCSKubernetesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var cluster cs.KubernetesCluster
	invoker := NewInvoker()
	if err := invoker.Run(func() error {
		raw, e := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return csClient.DescribeKubernetesCluster(d.Id())
		})
		if e != nil {
			return fmt.Errorf("Describing kubernetes cluster %#v failed, error message: %#v. Please check cluster in the console,", d.Id(), e)
		}
		cluster, _ = raw.(cs.KubernetesCluster)
		return nil
	}); err != nil {
		if NotFoundError(err) || IsExceptedError(err, ErrorClusterNotFound) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", cluster.Name)
	if cluster.Parameters.ImageId != "" {
		d.Set("image_id", cluster.Parameters.ImageId)
	} else {
		d.Set("image_id", cluster.Parameters.MasterImageId)
	}
	d.Set("vpc_id", cluster.VPCID)
	d.Set("security_group_id", cluster.SecurityGroupID)
	d.Set("key_name", cluster.Parameters.KeyPair)
	if size, err := strconv.Atoi(cluster.Parameters.MasterSystemDiskSize); err != nil {
		return BuildWrapError("strconv.Atoi", d.Id(), ProviderERROR, err)
	} else {
		d.Set("master_disk_size", size)
	}
	d.Set("master_disk_category", cluster.Parameters.MasterSystemDiskCategory)
	if size, err := strconv.Atoi(cluster.Parameters.WorkerSystemDiskSize); err != nil {
		return BuildWrapError("strconv.Atoi", d.Id(), ProviderERROR, err)
	} else {
		d.Set("worker_disk_size", size)
	}
	d.Set("worker_disk_category", cluster.Parameters.WorkerSystemDiskCategory)
	d.Set("availability_zone", cluster.ZoneId)
	d.Set("slb_internet_enabled", cluster.Parameters.PublicSLB)

	if cluster.Parameters.MasterInstanceChargeType == string(PrePaid) {
		d.Set("master_instance_charge_type", string(PrePaid))
		if period, err := strconv.Atoi(cluster.Parameters.MasterPeriod); err != nil {
			return BuildWrapError("strconv.Atoi", d.Id(), ProviderERROR, err)
		} else {
			d.Set("master_period", period)
		}
		d.Set("master_period_unit", cluster.Parameters.MasterPeriodUnit)
		d.Set("master_auto_renew", cluster.Parameters.MasterAutoRenew)
		if period, err := strconv.Atoi(cluster.Parameters.MasterAutoRenewPeriod); err != nil {
			return BuildWrapError("strconv.Atoi", d.Id(), ProviderERROR, err)
		} else {
			d.Set("master_auto_renew_period", period)
		}
	} else {
		d.Set("master_instance_charge_type", string(PostPaid))
	}

	if cluster.Parameters.WorkerInstanceChargeType == string(PrePaid) {
		d.Set("worker_instance_charge_type", string(PrePaid))
		if period, err := strconv.Atoi(cluster.Parameters.WorkerPeriod); err != nil {
			return BuildWrapError("strconv.Atoi", d.Id(), ProviderERROR, err)
		} else {
			d.Set("worker_period", period)
		}
		d.Set("worker_period_unit", cluster.Parameters.WorkerPeriodUnit)
		d.Set("worker_auto_renew", cluster.Parameters.WorkerAutoRenew)
		if period, err := strconv.Atoi(cluster.Parameters.WorkerAutoRenewPeriod); err != nil {
			return BuildWrapError("strconv.Atoi", d.Id(), ProviderERROR, err)
		} else {
			d.Set("worker_auto_renew_period", period)
		}
	} else {
		d.Set("worker_instance_charge_type", string(PostPaid))
	}

	if cidrMask, err := strconv.Atoi(cluster.Parameters.NodeCIDRMask); err == nil {
		d.Set("node_cidr_mask", cidrMask)
	} else {
		return BuildWrapError("strconv.Atoi", d.Id(), ProviderERROR, err)
	}

	if cluster.Parameters.WorkerDataDisk {
		if size, err := strconv.Atoi(cluster.Parameters.WorkerDataDiskSize); err != nil {
			return BuildWrapError("strconv.Atoi", d.Id(), ProviderERROR, err)
		} else {
			d.Set("worker_data_disk_size", size)
		}
		d.Set("worker_data_disk_category", cluster.Parameters.WorkerDataDiskCategory)
	}

	if cluster.Parameters.LoggingType != "None" {
		logConfig := map[string]interface{}{}
		logConfig["type"] = cluster.Parameters.LoggingType
		if cluster.Parameters.SLSProjectName == "None" {
			logConfig["project"] = ""
		} else {
			logConfig["project"] = cluster.Parameters.SLSProjectName
		}
		if err := d.Set("log_config", []map[string]interface{}{logConfig}); err != nil {
			return err
		}
	}

	// Each k8s cluster contains 3 master nodes
	if cluster.MetaData.MultiAZ || cluster.MetaData.SubClass == "3az" {
		numOfNodeA, err := strconv.Atoi(cluster.Parameters.NumOfNodesA)
		if err != nil {
			return fmt.Errorf("error convert NumOfNodesA %s to int: %s", cluster.Parameters.NumOfNodesA, err.Error())
		}
		numOfNodeB, err := strconv.Atoi(cluster.Parameters.NumOfNodesB)
		if err != nil {
			return fmt.Errorf("error convert NumOfNodesB %s to int: %s", cluster.Parameters.NumOfNodesB, err.Error())
		}
		numOfNodeC, err := strconv.Atoi(cluster.Parameters.NumOfNodesC)
		if err != nil {
			return fmt.Errorf("error convert NumOfNodesC %s to int: %s", cluster.Parameters.NumOfNodesC, err.Error())
		}
		d.Set("worker_numbers", []int{numOfNodeA, numOfNodeB, numOfNodeC})
		d.Set("vswitch_ids", []string{cluster.Parameters.VSwitchIdA, cluster.Parameters.VSwitchIdB, cluster.Parameters.VSwitchIdC})
		d.Set("master_instance_types", []string{cluster.Parameters.MasterInstanceTypeA, cluster.Parameters.MasterInstanceTypeB, cluster.Parameters.MasterInstanceTypeC})
		d.Set("worker_instance_types", []string{cluster.Parameters.WorkerInstanceTypeA, cluster.Parameters.WorkerInstanceTypeB, cluster.Parameters.WorkerInstanceTypeC})
	} else {
		if numOfNode, err := strconv.Atoi(cluster.Parameters.NumOfNodes); err != nil {
			return BuildWrapError("strconv.Atoi", d.Id(), ProviderERROR, err)
		} else {
			d.Set("worker_numbers", []int{numOfNode})
		}
		d.Set("vswitch_ids", []string{cluster.Parameters.VSwitchID})
		d.Set("master_instance_types", []string{cluster.Parameters.MasterInstanceType})
		d.Set("worker_instance_types", []string{cluster.Parameters.WorkerInstanceType})
	}

	var masterNodes []map[string]interface{}
	var workerNodes []map[string]interface{}

	pageNumber := 1
	for {
		var result []cs.KubernetesNodeType
		var pagination *cs.PaginationResult

		if err := invoker.Run(func() error {
			raw, e := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				nodes, paginationResult, err := csClient.GetKubernetesClusterNodes(d.Id(), common.Pagination{PageNumber: pageNumber, PageSize: PageSizeLarge})
				return []interface{}{nodes, paginationResult}, err
			})
			if e != nil {
				return e
			}
			result, _ = raw.([]interface{})[0].([]cs.KubernetesNodeType)
			pagination, _ = raw.([]interface{})[1].(*cs.PaginationResult)
			return nil
		}); err != nil {
			return fmt.Errorf("[ERROR] GetKubernetesClusterNodes got an error: %#v.", err)
		}

		if pageNumber == 1 && (len(result) == 0 || result[0].InstanceId == "") {
			err := resource.Retry(5*time.Minute, func() *resource.RetryError {
				if err := invoker.Run(func() error {
					raw, e := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
						nodes, _, err := csClient.GetKubernetesClusterNodes(d.Id(), common.Pagination{PageNumber: pageNumber, PageSize: PageSizeLarge})
						return nodes, err
					})
					if e != nil {
						return e
					}
					tmp, _ := raw.([]cs.KubernetesNodeType)
					if len(tmp) > 0 && tmp[0].InstanceId != "" {
						result = tmp
					}
					return nil
				}); err != nil {
					return resource.NonRetryableError(fmt.Errorf("[ERROR] GetKubernetesClusterNodes got an error: %#v.", err))
				}
				for _, stableState := range cs.NodeStableClusterState {
					// If cluster is in NodeStableClusteState, node list will not change
					if cluster.State == stableState {
						return nil
					}
				}
				time.Sleep(5 * time.Second)
				return resource.RetryableError(fmt.Errorf("[ERROR] There is no any nodes in kubernetes cluster %s.", d.Id()))
			})
			if err != nil {
				return err
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
	reqSLB := slb.CreateDescribeLoadBalancersRequest()
	reqSLB.ServerId = masterNodes[0]["id"].(string)
	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeLoadBalancers(reqSLB)
	})
	if err != nil {
		return fmt.Errorf("[ERROR] DescribeLoadBalancers by server id %s got an error: %#v.", masterNodes[0]["id"].(string), err)
	}
	lbs, _ := raw.(*slb.DescribeLoadBalancersResponse)
	for _, lb := range lbs.LoadBalancers.LoadBalancer {
		if strings.ToLower(lb.AddressType) == strings.ToLower(string(Internet)) {
			d.Set("slb_internet", lb.LoadBalancerId)
			connection["api_server_internet"] = fmt.Sprintf("https://%s:6443", lb.Address)
			connection["master_public_ip"] = lb.Address
		} else {
			d.Set("slb_intranet", lb.LoadBalancerId)
			connection["api_server_intranet"] = fmt.Sprintf("https://%s:6443", lb.Address)
		}
	}
	connection["service_domain"] = fmt.Sprintf("*.%s.%s.alicontainer.com", d.Id(), cluster.RegionID)

	d.Set("connections", connection)
	req := vpc.CreateDescribeNatGatewaysRequest()
	req.VpcId = cluster.VPCID
	raw, err = client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeNatGateways(req)
	})
	if err != nil {
		return fmt.Errorf("[ERROR] DescribeNatGateways by VPC Id %s: %#v.", cluster.VPCID, err)
	}
	nat, _ := raw.(*vpc.DescribeNatGatewaysResponse)
	if nat != nil && len(nat.NatGateways.NatGateway) > 0 {
		d.Set("nat_gateway_id", nat.NatGateways.NatGateway[0].NatGatewayId)
	}

	if err := invoker.Run(func() error {
		raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return csClient.GetClusterCerts(d.Id())
		})
		if err != nil {
			return err
		}
		cert, _ := raw.(cs.ClusterCerts)
		if ce, ok := d.GetOk("client_cert"); ok && ce.(string) != "" {
			if err := writeToFile(ce.(string), cert.Cert); err != nil {
				return err
			}
		}
		if key, ok := d.GetOk("client_key"); ok && key.(string) != "" {
			if err := writeToFile(key.(string), cert.Key); err != nil {
				return err
			}
		}
		if ca, ok := d.GetOk("cluster_ca_cert"); ok && ca.(string) != "" {
			if err := writeToFile(ca.(string), cert.CA); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("Get Cluster %s Certs got an error: %#v.", d.Id(), err)
	}

	var config cs.ClusterConfig
	if file, ok := d.GetOk("kube_config"); ok && file.(string) != "" {
		if err := invoker.Run(func() error {
			raw, e := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				return csClient.GetClusterConfig(d.Id())
			})
			if e != nil {
				return e
			}
			config, _ = raw.(cs.ClusterConfig)
			return nil
		}); err != nil {
			return fmt.Errorf("GetClusterConfig got an error: %#v.", err)
		}
		if err := writeToFile(file.(string), config.Config); err != nil {
			return err
		}
	}

	return nil
}

func resourceAlicloudCSKubernetesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	invoker := NewInvoker()
	var cluster cs.ClusterType
	return resource.Retry(15*time.Minute, func() *resource.RetryError {
		if err := invoker.Run(func() error {
			_, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				return nil, csClient.DeleteCluster(d.Id())
			})
			return err
		}); err != nil {
			if NotFoundError(err) || IsExceptedError(err, ErrorClusterNotFound) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Delete Kubernetes Cluster timeout and get an error: %#v.", err))
		}

		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				return csClient.DescribeCluster(d.Id())
			})
			if err != nil {
				return err
			}
			cluster, _ = raw.(cs.ClusterType)
			return nil
		}); err != nil {
			if NotFoundError(err) || IsExceptedError(err, ErrorClusterNotFound) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Describing Kubernetes Cluster got an error: %#v", err))
		}
		if cluster.ClusterID == "" {
			return nil
		}

		if string(cluster.State) == string(Deleting) {
			time.Sleep(5 * time.Second)
		}

		return resource.RetryableError(fmt.Errorf("Delete Kubernetes Cluster timeout."))
	})
}

func isMultiAZClusterAndCheck(d *schema.ResourceData) (bool, error) {
	masterInstanceTypes := expandStringList(d.Get("master_instance_types").([]interface{}))
	workerInstanceTypes := expandStringList(d.Get("worker_instance_types").([]interface{}))
	vswitchIDs := expandStringList(d.Get("vswitch_ids").([]interface{}))
	workerNumbers := expandIntList(d.Get("worker_numbers").([]interface{}))

	if len(masterInstanceTypes) != len(workerInstanceTypes) {
		return false, errors.New("length of fields `master_instance_types`, `worker_instance_types` must match")
	}
	if len(masterInstanceTypes) == 1 {
		// single AZ
		return false, nil
	} else if len(masterInstanceTypes) == 3 {
		if len(vswitchIDs) != 3 || len(workerNumbers) != 3 {
			return true, errors.New("length of fields `vswitch_ids`, `worker_numbers` must be 3 for multiAZ clusters")

		} else {
			return true, nil
		}
	} else {
		return false, errors.New("length of fields `master_instance_types`, `worker_instance_types` should be either 3 (for MultiAZ) or 1 (for Single AZ)")
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
		vsw, err := vpcService.DescribeVswitch(vswitchID)
		if err != nil {
			return nil, err
		}
		vpcId = vsw.VpcId
		if zoneId != "" && zoneId != vsw.ZoneId {
			return nil, fmt.Errorf("The specified vswitch %s isn't in the zone %s.", vsw.VSwitchId, zoneId)
		}
		zoneId = vsw.ZoneId
	} else if !d.Get("new_nat_gateway").(bool) {
		return nil, fmt.Errorf("The automatic created VPC and VSwitch must set 'new_nat_gateway' to 'true'.")
	}

	loggingType, slsProjectName, err := parseKubernetesClusterLogConfig(d)
	if err != nil {
		return nil, err
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
		LoginPassword:            d.Get("password").(string),
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

	return creationArgs, nil
}

func buildKubernetesMultiAZArgs(d *schema.ResourceData, meta interface{}) (*cs.KubernetesMultiAZCreationArgs, error) {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	vpcService := VpcService{client}

	// Ensure instance_type is valid
	zoneId, validZones, err := ecsService.DescribeAvailableResources(d, meta, InstanceTypeResource)
	if err != nil {
		return nil, err
	}
	instanceTypes := expandStringList(d.Get("master_instance_types").([]interface{}))
	instanceTypes = append(instanceTypes, expandStringList(d.Get("worker_instance_types").([]interface{}))...)

	for _, instanceType := range instanceTypes {
		if err := ecsService.InstanceTypeValidation(instanceType, zoneId, validZones); err != nil {
			return nil, err
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

	vsw, err := vpcService.DescribeVswitch(vswitchIDs[0])
	if err != nil {
		return nil, err
	}

	loggingType, slsProjectName, err := parseKubernetesClusterLogConfig(d)
	if err != nil {
		return nil, err
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
		LoginPassword:            d.Get("password").(string),
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
			creationArgs.WorkerPeriodUnit = d.Get("Worker_period_unit").(string)
		}
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
			case KubernetesClusterLoggingTypeSLS:
				if config["project"].(string) == "" {
					return "", "", fmt.Errorf("SLS project name must be provided when choosing SLS as log_config.")
				}
				if config["project"].(string) == "None" {
					return "", "", fmt.Errorf("SLS project name must not be `None`.")
				}
				slsProjectName = config["project"].(string)
				break
			default:
				break
			}
		}
	}
	return loggingType, slsProjectName, nil
}
