package alicloud

import (
	"fmt"
	"time"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/cs"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

const (
	ManagedKubernetesClusterNetworkTypeFlannel      = "flannel"
	ManagedKubernetesClusterNetworkTypeTerway       = "terway"
	ManagedKubernetesCreationDefaultTimeoutInMinute = 60
)

func resourceAlicloudCSManagedKubernetes() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSManagedKubernetesCreate,
		Read:   resourceAlicloudCSManagedKubernetesRead,
		Update: resourceAlicloudCSManagedKubernetesUpdate,
		Delete: resourceAlicloudCSManagedKubernetesDelete,
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
			"vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateContainerVswitchId,
				},
				MinItems:         1,
				MaxItems:         5,
				DiffSuppressFunc: csForceUpdate,
			},
			"new_nat_gateway": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"worker_instance_types": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems:         1,
				MaxItems:         5,
				DiffSuppressFunc: csForceUpdate,
			},
			"worker_number": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3,
			},
			"worker_numbers": {
				Type:       schema.TypeList,
				Deprecated: "Field 'worker_numbers' has been deprecated from provider version 1.53.0. New field 'worker_number' replaces it.",
				Elem: &schema.Schema{
					Type:    schema.TypeInt,
					Default: 3,
				},
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
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
				ValidateFunc: validateAllowedStringValue([]string{ManagedKubernetesClusterNetworkTypeFlannel, ManagedKubernetesClusterNetworkTypeTerway}),
			},
			"image_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: imageIdSuppressFunc,
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
			"slb_internet_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"install_cloud_monitor": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"force_update": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
			"security_group_id": {
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

func resourceAlicloudCSManagedKubernetesCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	invoker := NewInvoker()

	args, err := buildManagedKubernetesArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}
	var requestInfo *cs.Client
	var response interface{}
	if err := invoker.Run(func() error {
		raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			requestInfo = csClient
			return csClient.CreateKubernetesCluster(common.Region(client.RegionId), args)
		})
		response = raw
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_managed_kubernetes", "CreateKubernetesCluster", DenverdinoAliyungo)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["RegionId"] = common.Region(client.RegionId)
		requestMap["Args"] = args
		addDebug("CreateKubernetesCluster", response, requestInfo, requestMap)
	}
	cluster, _ := response.(cs.ClusterCreationResponse)
	d.SetId(cluster.ClusterID)

	if err := invoker.Run(func() error {
		raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			requestInfo = csClient
			return nil, csClient.WaitForClusterAsyn(d.Id(), cs.Running, 3600)
		})
		response = raw
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "WaitForClusterAsyn", DenverdinoAliyungo)
	}
	if debugOn() {
		waitRequestMap := make(map[string]interface{})
		waitRequestMap["ClusterId"] = d.Id()
		waitRequestMap["Status"] = cs.Running
		waitRequestMap["TimeOut"] = 3600
		addDebug("WaitForClusterAsyn", response, requestInfo, waitRequestMap)
	}

	return resourceAlicloudCSManagedKubernetesRead(d, meta)
}

func resourceAlicloudCSManagedKubernetesUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	d.Partial(true)
	invoker := NewInvoker()
	if d.HasChange("worker_number") || d.HasChange("worker_numbers") {
		var scaleSize int
		if d.HasChange("worker_number") {
			o, n := d.GetChange("worker_number")
			oldNumber := o.(int)
			newNumber := n.(int)
			scaleSize = newNumber - oldNumber
		} else {
			var oSum, nSum int
			o, n := d.GetChange("worker_numbers")
			oInts := o.([]interface{})
			for _, v := range oInts {
				oSum += v.(int)
			}
			nInts := n.([]interface{})
			for _, v := range nInts {
				nSum += v.(int)
			}
			scaleSize = nSum - oSum
		}
		if scaleSize < 0 {
			return fmt.Errorf("cannot scale down cluster")
		}

		workerInstanceTypes := deduplicateInstanceTypes(expandStringList(d.Get("worker_instance_types").([]interface{})))

		// When cluster was created using keypair, LoginPassword will be ignored.
		// When cluster was created using password, LoginPassword is required to resize.
		args := &cs.KubernetesClusterScaleArgs{
			LoginPassword:            d.Get("password").(string),
			KeyPair:                  d.Get("key_name").(string),
			WorkerInstanceTypes:      workerInstanceTypes,
			WorkerSystemDiskCategory: ecs.DiskCategory(d.Get("worker_disk_category").(string)),
			WorkerSystemDiskSize:     int64(d.Get("worker_disk_size").(int)),
			Count:                    scaleSize,
		}
		if _, ok := d.GetOk("worker_data_disk_category"); ok {
			args.WorkerDataDisk = true
		}

		var requestInfo *cs.Client
		var response interface{}
		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				requestInfo = csClient
				return nil, csClient.ScaleKubernetesCluster(d.Id(), args)
			})
			response = raw
			return err
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ScaleKubernetesCluster", DenverdinoAliyungo)
		}
		if debugOn() {
			requestMap := make(map[string]interface{})
			requestMap["ClusterId"] = d.Id()
			requestMap["Args"] = args
			addDebug("ScaleKubernetesCluster", response, requestInfo, requestMap)
		}

		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				requestInfo = csClient
				return nil, csClient.WaitForClusterAsyn(d.Id(), cs.Running, 3600)
			})
			response = raw
			return err
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "WaitForClusterAsyn", DenverdinoAliyungo)
		}
		if debugOn() {
			waitRequestMap := make(map[string]interface{})
			waitRequestMap["ClusterId"] = d.Id()
			waitRequestMap["Status"] = cs.Running
			waitRequestMap["TimeOut"] = 3600
			addDebug("WaitForClusterAsyn", response, requestInfo, waitRequestMap)
		}
		d.SetPartial("worker_number")
	}

	if d.HasChange("name") || d.HasChange("name_prefix") {
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

	return resourceAlicloudCSManagedKubernetesRead(d, meta)
}

func resourceAlicloudCSManagedKubernetesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	invoker := NewInvoker()
	object, err := csService.DescribeCsManagedKubernetes(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
	}

	d.Set("name", object.Name)
	d.Set("vpc_id", object.VPCID)
	d.Set("security_group_id", object.SecurityGroupID)
	d.Set("availability_zone", object.ZoneId)

	var workerNodes []map[string]interface{}
	var response interface{}
	pageNumber := 1
	for {
		var result []cs.KubernetesNodeType
		var pagination *cs.PaginationResult

		var requestInfo *cs.Client

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
			var requestInfo *cs.Client

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
				return resource.RetryableError(Error("[ERROR] There is no any nodes in ManagedKubernetes cluster %s.", d.Id()))
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
			workerNodes = append(workerNodes, mapping)
		}

		if len(result) < pagination.PageSize {
			break
		}
		pageNumber += 1
	}
	d.Set("worker_nodes", workerNodes)

	var requestInfo *cs.Client

	if err := invoker.Run(func() error {
		raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
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
		requestMap["ClusterId"] = d.Id()
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
		var requestInfo *cs.Client

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
			requestMap["ClusterId"] = d.Id()
			addDebug("GetClusterConfig", response, requestInfo, requestMap)
		}
		config, _ = response.(cs.ClusterConfig)

		if err := writeToFile(file.(string), config.Config); err != nil {
			return WrapError(err)
		}
	}
	return nil
}

func resourceAlicloudCSManagedKubernetesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	var requestInfo *cs.Client
	invoker := NewInvoker()
	var response interface{}

	err := resource.Retry(30*time.Minute, func() *resource.RetryError {
		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				return nil, csClient.DeleteCluster(d.Id())
			})
			response = raw
			return err
		}); err != nil {
			return resource.NonRetryableError(err)
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
	return WrapError(csService.WaitForCSManagedKubernetes(d.Id(), Deleted, DefaultLongTimeout))
}

func buildManagedKubernetesArgs(d *schema.ResourceData, meta interface{}) (*cs.KubernetesCreationArgs, error) {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	rawWorkerInstanceTypes := expandStringList(d.Get("worker_instance_types").([]interface{}))

	workerInstanceTypes := deduplicateInstanceTypes(rawWorkerInstanceTypes)

	vswitchIds := expandStringList(d.Get("vswitch_ids").([]interface{}))

	workerNumber := d.Get("worker_number").(int)

	var clusterName string
	if v, ok := d.GetOk("name"); ok {
		clusterName = v.(string)
	} else {
		clusterName = resource.PrefixedUniqueId(d.Get("name_prefix").(string))
	}

	var vpcId, zoneId string
	for _, vswId := range vswitchIds {
		vsw, err := vpcService.DescribeVSwitch(vswId)
		if err != nil {
			return nil, err
		}
		if vpcId == "" {
			vpcId = vsw.VpcId
			zoneId = vsw.ZoneId
		} else {
			if vsw.VpcId != vpcId {
				return nil, fmt.Errorf("all specified vswitch should be in the same vpc %s.", vswitchIds)
			}
		}
	}

	creationArgs := &cs.KubernetesCreationArgs{
		Name:                     clusterName,
		ClusterType:              "ManagedKubernetes",
		DisableRollback:          true,
		TimeoutMins:              ManagedKubernetesCreationDefaultTimeoutInMinute,
		RegionId:                 client.RegionId,
		WorkerInstanceTypes:      workerInstanceTypes,
		VPCID:                    vpcId,
		VSwitchIds:               vswitchIds,
		LoginPassword:            d.Get("password").(string),
		KeyPair:                  d.Get("key_name").(string),
		ImageId:                  d.Get("image_id").(string),
		Network:                  d.Get("cluster_network_type").(string),
		NumOfNodes:               int64(workerNumber),
		WorkerSystemDiskCategory: ecs.DiskCategory(d.Get("worker_disk_category").(string)),
		WorkerSystemDiskSize:     int64(d.Get("worker_disk_size").(int)),
		SNatEntry:                d.Get("new_nat_gateway").(bool),
		KubernetesVersion:        d.Get("version").(string),
		ContainerCIDR:            d.Get("pod_cidr").(string),
		ServiceCIDR:              d.Get("service_cidr").(string),
		CloudMonitorFlags:        d.Get("install_cloud_monitor").(bool),
		PublicSLB:                d.Get("slb_internet_enabled").(bool),
		ZoneId:                   zoneId,
	}

	if v, ok := d.GetOk("worker_data_disk_category"); ok {
		creationArgs.WorkerDataDiskCategory = v.(string)
		creationArgs.WorkerDataDisk = true
		creationArgs.WorkerDataDiskSize = int64(d.Get("worker_data_disk_size").(int))
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

func deduplicateInstanceTypes(instanceTypes []string) []string {
	var ret []string
	instanceTypesMaps := make(map[string]bool)
	for _, v := range instanceTypes {
		instanceTypesMaps[v] = true
	}
	for k := range instanceTypesMaps {
		ret = append(ret, k)
	}
	return ret
}
