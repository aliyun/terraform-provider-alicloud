package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/cs"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
			"name": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validateContainerName,
				ConflictsWith: []string{"name_prefix"},
			},
			"name_prefix": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "Terraform-Creation",
				ValidateFunc:  validateContainerNamePrefix,
				ConflictsWith: []string{"name"},
			},
			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"vswitch_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"new_nat_gateway": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"master_instance_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateInstanceType,
			},
			"worker_instance_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateInstanceType,
			},
			"worker_number": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3,
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"pod_cidr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_cidr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_ssh": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"master_disk_size": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      40,
				ValidateFunc: validateIntegerInRange(40, 500),
			},
			"master_disk_category": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  ecs.DiskCategoryCloudEfficiency,
				ValidateFunc: validateAllowedStringValue([]string{
					string(ecs.DiskCategoryCloudEfficiency), string(ecs.DiskCategoryCloudSSD)}),
			},
			"worker_disk_size": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      40,
				ValidateFunc: validateIntegerInRange(20, 32768),
			},
			"worker_disk_category": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  ecs.DiskCategoryCloudEfficiency,
				ValidateFunc: validateAllowedStringValue([]string{
					string(ecs.DiskCategoryCloudEfficiency), string(ecs.DiskCategoryCloudSSD)}),
			},
			"install_cloud_monitor": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"is_outdated": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},

			"kube_config": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_cert": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_ca_cert": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			// 'version' is a reserved parameter and it just is used to test. No Recommendation to expose it.
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"nodes": &schema.Schema{
				Type:       schema.TypeList,
				Optional:   true,
				Elem:       &schema.Schema{Type: schema.TypeString},
				Deprecated: "Field 'nodes' has been deprecated from provider version 1.9.4. New field 'master_nodes' replaces it.",
			},
			"master_nodes": &schema.Schema{
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
			"worker_nodes": &schema.Schema{
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
			"connections": &schema.Schema{
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
			"slb_id": &schema.Schema{
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Field 'slb_id' has been deprecated from provider version 1.9.2. New field 'slb_internet' replaces it.",
			},
			"slb_internet": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"slb_intranet": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"nat_gateway_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCSKubernetesCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	conn := client.csconn

	args, err := buildKunernetesArgs(d, meta)
	if err != nil {
		return err
	}
	invoker := NewInvoker()
	if err := invoker.Run(func() error {
		cluster, err := conn.CreateKubernetesCluster(getRegion(d, meta), args)
		if err != nil {
			return err
		}
		d.SetId(cluster.ClusterID)
		return nil
	}); err != nil {
		return fmt.Errorf("Creating Kubernetes Cluster got an error: %#v", err)
	}

	if err := invoker.Run(func() error {
		return conn.WaitForClusterAsyn(d.Id(), cs.Running, 3600)
	}); err != nil {
		return fmt.Errorf("Waitting for kubernetes cluster %#v got an error: %#v", cs.Running, err)
	}

	return resourceAlicloudCSKubernetesUpdate(d, meta)
}

func resourceAlicloudCSKubernetesUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).csconn
	d.Partial(true)
	invoker := NewInvoker()
	if d.HasChange("worker_number") && !d.IsNewResource() {
		// Ensure instance_type is generation three
		args, err := buildKunernetesArgs(d, meta)
		if err != nil {
			return err
		}
		if err := invoker.Run(func() error {
			return conn.ResizeKubernetes(d.Id(), args)
		}); err != nil {
			return fmt.Errorf("Resize Cluster got an error: %#v", err)
		}

		if err := invoker.Run(func() error {
			return conn.WaitForClusterAsyn(d.Id(), cs.Running, 3600)
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
			if err := conn.ModifyClusterName(d.Id(), clusterName); err != nil && !IsExceptedError(err, ErrorClusterNameAlreadyExist) {
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
	client := meta.(*AliyunClient)

	var cluster cs.ClusterType
	invoker := NewInvoker()
	if err := invoker.Run(func() error {
		c, e := client.csconn.DescribeCluster(d.Id())
		if e != nil {
			return e
		}
		cluster = c
		return nil
	}); err != nil {
		if NotFoundError(err) || IsExceptedError(err, ErrorClusterNotFound) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", cluster.Name)
	// Each k8s cluster contains 3 master nodes
	d.Set("worker_number", cluster.Size-KubernetesMasterNumber)
	d.Set("vswitch_id", cluster.VSwitchID)
	d.Set("vpc_id", cluster.VPCID)
	d.Set("security_group_id", cluster.SecurityGroupID)

	var masterNodes []map[string]interface{}
	var workerNodes []map[string]interface{}
	var master, worker cs.KubernetesNodeType
	var workerId string

	pageNumber := 1
	for {
		var result []cs.KubernetesNodeType
		var pagination *cs.PaginationResult

		if err := invoker.Run(func() error {
			r, p, e := client.csconn.GetKubernetesClusterNodes(d.Id(), common.Pagination{PageNumber: pageNumber, PageSize: PageSizeLarge})
			if e != nil {
				return e
			}
			result = r
			pagination = p
			return nil
		}); err != nil {
			return fmt.Errorf("[ERROR] GetKubernetesClusterNodes got an error: %#v.", err)
		}

		if pageNumber == 1 && (len(result) == 0 || result[0].InstanceId == "") {
			err := resource.Retry(5*time.Minute, func() *resource.RetryError {
				if err := invoker.Run(func() error {
					tmp, _, e := client.csconn.GetKubernetesClusterNodes(d.Id(), common.Pagination{PageNumber: pageNumber, PageSize: PageSizeLarge})
					if e != nil {
						return e
					}
					if len(tmp) > 0 && tmp[0].InstanceId != "" {
						result = tmp
					}
					return nil
				}); err != nil {
					return resource.NonRetryableError(fmt.Errorf("[ERROR] GetKubernetesClusterNodes got an error: %#v.", err))
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
				master = node
				masterNodes = append(masterNodes, mapping)
			} else {
				if workerId == "" {
					workerId = node.InstanceId
				}
				worker = node
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

	d.Set("master_instance_type", master.InstanceType)
	if disks, err := client.DescribeDisksByType(master.InstanceId, DiskTypeSystem); err != nil {
		return fmt.Errorf("[ERROR] DescribeDisks By Id %s: %#v.", master.InstanceId, err)
	} else if len(disks) > 0 {
		d.Set("master_disk_size", disks[0].Size)
		d.Set("master_disk_category", disks[0].Category)
		d.Set("availability_zone", disks[0].ZoneId)
	}

	d.Set("worker_instance_type", worker.InstanceType)
	// worker.InstanceId will be empty in sometimes

	if disks, err := client.DescribeDisksByType(workerId, DiskTypeSystem); err != nil {
		return fmt.Errorf("[ERROR] DescribeDisks By Id %s: %#v.", workerId, err)
	} else if len(disks) > 0 {
		d.Set("worker_disk_size", disks[0].Size)
		d.Set("worker_disk_category", disks[0].Category)
	}

	if cluster.SecurityGroupID == "" {
		if inst, err := client.DescribeInstanceAttribute(workerId); err != nil {
			return fmt.Errorf("[ERROR] DescribeInstanceAttribute %s got an error: %#v.", workerId, err)
		} else {
			d.Set("security_group_id", inst.SecurityGroupIds.SecurityGroupId[0])
		}
	}

	// Get slb information
	connection := make(map[string]string)
	lbs, err := client.slbconn.DescribeLoadBalancers(&slb.DescribeLoadBalancersArgs{
		RegionId: getRegion(d, meta),
		ServerId: master.InstanceId,
	})
	if err != nil {
		return fmt.Errorf("[ERROR] DescribeLoadBalancers by server id %s got an error: %#v.", workerId, err)
	}
	for _, lb := range lbs {
		if lb.AddressType == slb.InternetAddressType {
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
	if nat, err := client.vpcconn.DescribeNatGateways(req); err != nil {
		return fmt.Errorf("[ERROR] DescribeNatGateways by VPC Id %s: %#v.", cluster.VPCID, err)
	} else if nat != nil && len(nat.NatGateways.NatGateway) > 0 {
		d.Set("nat_gateway_id", nat.NatGateways.NatGateway[0].NatGatewayId)
	}

	if err := invoker.Run(func() error {
		cert, err := client.csconn.GetClusterCerts(d.Id())
		if err != nil {
			return err
		}
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
			c, e := client.csconn.GetClusterConfig(d.Id())
			if e != nil {
				return e
			}
			config = c
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
	conn := meta.(*AliyunClient).csconn
	invoker := NewInvoker()
	var cluster cs.ClusterType
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if err := invoker.Run(func() error {
			return conn.DeleteCluster(d.Id())
		}); err != nil {
			if NotFoundError(err) || IsExceptedError(err, ErrorClusterNotFound) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Delete Kubernetes Cluster timeout and get an error: %#v.", err))
		}

		if err := invoker.Run(func() error {
			resp, err := conn.DescribeCluster(d.Id())
			if err != nil {
				return err
			}
			cluster = resp
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

func buildKunernetesArgs(d *schema.ResourceData, meta interface{}) (*cs.KubernetesCreationArgs, error) {
	client := meta.(*AliyunClient)

	// Ensure instance_type is valid
	zoneId, validZones, err := meta.(*AliyunClient).DescribeAvailableResources(d, meta, InstanceTypeResource)
	if err != nil {
		return nil, err
	}
	if err := meta.(*AliyunClient).InstanceTypeValidation(d.Get("master_instance_type").(string), zoneId, validZones); err != nil {
		return nil, err
	}

	if err := meta.(*AliyunClient).InstanceTypeValidation(d.Get("worker_instance_type").(string), zoneId, validZones); err != nil {
		return nil, err
	}

	var clusterName string
	if v, ok := d.GetOk("name"); ok {
		clusterName = v.(string)
	} else {
		clusterName = resource.PrefixedUniqueId(d.Get("name_prefix").(string))
	}

	stackArgs := &cs.KubernetesStackArgs{
		MasterInstanceType:       d.Get("master_instance_type").(string),
		WorkerInstanceType:       d.Get("worker_instance_type").(string),
		Password:                 d.Get("password").(string),
		NumOfNodes:               int64(d.Get("worker_number").(int)),
		MasterSystemDiskCategory: ecs.DiskCategory(d.Get("master_disk_category").(string)),
		MasterSystemDiskSize:     int64(d.Get("master_disk_size").(int)),
		WorkerSystemDiskCategory: ecs.DiskCategory(d.Get("worker_disk_category").(string)),
		WorkerSystemDiskSize:     int64(d.Get("worker_disk_size").(int)),
		SNatEntry:                d.Get("new_nat_gateway").(bool),
		KubernetesVersion:        d.Get("version").(string),
		ContainerCIDR:            d.Get("pod_cidr").(string),
		ServiceCIDR:              d.Get("service_cidr").(string),
		SSHFlags:                 d.Get("enable_ssh").(bool),
		ImageID:                  KubernetesImageId,
		CloudMonitorFlags:        d.Get("install_cloud_monitor").(bool),
		ZoneId:                   zoneId,
	}

	if v, ok := d.GetOk("vswitch_id"); ok && len(Trim(v.(string))) > 0 {
		stackArgs.VSwitchID = Trim(v.(string))
		vsw, err := client.DescribeVswitch(stackArgs.VSwitchID)
		if err != nil {
			return nil, err
		}
		stackArgs.VPCID = vsw.VpcId
		if stackArgs.ZoneId != "" && vsw.ZoneId != vsw.ZoneId {
			return nil, fmt.Errorf("The specified vswitch %s isn't in the zone %s.", vsw.VSwitchId, stackArgs.ZoneId)
		}
		stackArgs.ZoneId = vsw.ZoneId
	} else if !stackArgs.SNatEntry {
		return nil, fmt.Errorf("The automatic created VPC and VSwitch must set 'new_nat_gateway' to 'true'.")
	}

	return &cs.KubernetesCreationArgs{
		Name:              clusterName,
		ClusterType:       "Kubernetes",
		DisableRollback:   true,
		TimeoutMins:       60,
		KubernetesVersion: stackArgs.KubernetesVersion,
		StackParams:       *stackArgs,
	}, nil
}
