package alicloud

import (
	"fmt"
	"time"

	"strings"

	"errors"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/cs"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	KubernetesClusterNetworkTypeFlannel = "flannel"
	KubernetesClusterNetworkTypeTerway  = "terway"

	KubernetesClusterLoggingTypeSLS = "SLS"
)

var (
	KubernetesClusterNodeCIDRMasks          = []string{"24", "25", "26", "27", "28"}
	KubernetesClusterNodeCIDRMasksByDefault = "24"
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
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'vswitch_id' has been deprecated from provider version 1.16.0. New field 'vswitch_ids' replaces it.",
			},
			"vswitch_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MaxItems: 3,
			},
			"new_nat_gateway": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"master_instance_type": &schema.Schema{
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'master_instance_type' has been deprecated from provider version 1.16.0. New field 'master_instance_types' replaces it.",
			},
			"master_instance_types": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems: 1,
				MaxItems: 3,
			},
			"worker_instance_type": &schema.Schema{
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'worker_instance_type' has been deprecated from provider version 1.16.0. New field 'worker_instance_types' replaces it.",
			},
			"worker_instance_types": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems: 1,
				MaxItems: 3,
			},
			"worker_number": &schema.Schema{
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "Field 'worker_number' has been deprecated from provider version 1.16.0. New field 'worker_numbers' replaces it.",
			},

			"worker_numbers": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:    schema.TypeInt,
					Default: 3,
				},
				MinItems: 1,
				MaxItems: 3,
			},
			"password": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Sensitive:     true,
				ConflictsWith: []string{"key_name"},
			},
			"key_name": &schema.Schema{
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				ConflictsWith: []string{"password"},
			},
			"pod_cidr": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"service_cidr": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"cluster_network_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{KubernetesClusterNetworkTypeFlannel, KubernetesClusterNetworkTypeTerway}),
			},
			"node_cidr_mask": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      KubernetesClusterNodeCIDRMasksByDefault,
				ValidateFunc: validateAllowedStringValue(KubernetesClusterNodeCIDRMasks),
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
			"enable_ssh": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"master_disk_size": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      40,
				ForceNew:     true,
				ValidateFunc: validateIntegerInRange(40, 500),
			},
			"master_disk_category": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  DiskCloudEfficiency,
				ForceNew: true,
				ValidateFunc: validateAllowedStringValue([]string{
					string(DiskCloudEfficiency), string(DiskCloudSSD)}),
			},
			"worker_disk_size": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      40,
				ForceNew:     true,
				ValidateFunc: validateIntegerInRange(20, 32768),
			},
			"worker_disk_category": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  DiskCloudEfficiency,
				ValidateFunc: validateAllowedStringValue([]string{
					string(DiskCloudEfficiency), string(DiskCloudSSD)}),
			},
			"worker_data_disk_size": &schema.Schema{
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				Default:          40,
				ValidateFunc:     validateIntegerInRange(20, 32768),
				DiffSuppressFunc: workerDataDiskSizeSuppressFunc,
			},
			"worker_data_disk_category": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validateAllowedStringValue([]string{
					string(DiskCloudEfficiency), string(DiskCloudSSD)}),
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
	invoker := NewInvoker()

	if isMultiAZ, err := isMultiAZClusterAndCheck(d); err != nil {
		return err
	} else if isMultiAZ {
		args, err := buildKubernetesMultiAZArgs(d, meta)
		if err != nil {
			return err
		}
		if err := invoker.Run(func() error {
			cluster, err := conn.CreateKubernetesMultiAZCluster(getRegion(d, meta), args)
			if err != nil {
				return err
			}
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
			cluster, err := conn.CreateKubernetesCluster(getRegion(d, meta), args)
			if err != nil {
				return err
			}
			d.SetId(cluster.ClusterID)
			return nil
		}); err != nil {
			return fmt.Errorf("Creating Kubernetes Cluster got an error: %#v", err)
		}
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
			return conn.ResizeKubernetesCluster(d.Id(), args)
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

	var cluster cs.KubernetesCluster
	invoker := NewInvoker()
	if err := invoker.Run(func() error {
		c, e := client.csconn.DescribeKubernetesCluster(d.Id())
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
	d.Set("vpc_id", cluster.VPCID)
	d.Set("security_group_id", cluster.SecurityGroupID)
	d.Set("key_name", cluster.Parameters.KeyPair)
	d.Set("master_disk_size", cluster.Parameters.MasterSystemDiskSize)
	d.Set("master_disk_category", cluster.Parameters.MasterSystemDiskCategory)
	d.Set("worker_disk_size", cluster.Parameters.WorkerSystemDiskSize)
	d.Set("worker_disk_category", cluster.Parameters.WorkerSystemDiskCategory)
	d.Set("availability_zone", cluster.ZoneId)
	d.Set("node_cidr_mask", cluster.Parameters.NodeCIDRMask)

	if cluster.Parameters.WorkerDataDisk {
		d.Set("worker_data_disk_size", cluster.Parameters.WorkerDataDiskSize)
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
		numOfNode, err := strconv.Atoi(cluster.Parameters.NumOfNodes)
		if err != nil {
			return fmt.Errorf("error convert NumOfNodes %s to int: %s", cluster.Parameters.NumOfNodes, err.Error())
		}
		d.Set("worker_numbers", []int{numOfNode})
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
	lbs, err := client.slbconn.DescribeLoadBalancers(reqSLB)
	if err != nil {
		return fmt.Errorf("[ERROR] DescribeLoadBalancers by server id %s got an error: %#v.", masterNodes[0]["id"].(string), err)
	}
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
	client := meta.(*AliyunClient)

	// Ensure instance_type is valid
	zoneId, validZones, err := meta.(*AliyunClient).DescribeAvailableResources(d, meta, InstanceTypeResource)
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

	if err := meta.(*AliyunClient).InstanceTypeValidation(masterInstanceType, zoneId, validZones); err != nil {
		return nil, err
	}

	if err := meta.(*AliyunClient).InstanceTypeValidation(workerInstanceType, zoneId, validZones); err != nil {
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
		vsw, err := client.DescribeVswitch(vswitchID)
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
		Network:                  d.Get("cluster_network_type").(string),
		NodeCIDRMask:             d.Get("node_cidr_mask").(string),
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
		CloudMonitorFlags:        d.Get("install_cloud_monitor").(bool),
		ZoneId:                   zoneId,
	}

	if v, ok := d.GetOk("worker_data_disk_category"); ok {
		creationArgs.WorkerDataDiskCategory = v.(string)
		creationArgs.WorkerDataDisk = true
		creationArgs.WorkerDataDiskSize = int64(d.Get("worker_data_disk_size").(int))
	}

	return creationArgs, nil
}

func buildKubernetesMultiAZArgs(d *schema.ResourceData, meta interface{}) (*cs.KubernetesMultiAZCreationArgs, error) {
	client := meta.(*AliyunClient)

	// Ensure instance_type is valid
	zoneId, validZones, err := client.DescribeAvailableResources(d, meta, InstanceTypeResource)
	if err != nil {
		return nil, err
	}
	instanceTypes := expandStringList(d.Get("master_instance_types").([]interface{}))
	instanceTypes = append(instanceTypes, expandStringList(d.Get("worker_instance_types").([]interface{}))...)

	for _, instanceType := range instanceTypes {
		if err := meta.(*AliyunClient).InstanceTypeValidation(instanceType, zoneId, validZones); err != nil {
			return nil, err
		}
	}

	var clusterName string
	if v, ok := d.GetOk("name"); ok {
		clusterName = v.(string)
	} else {
		return nil, errors.New("The 'name' is required for Kubernetes MultiAZ.")
	}

	masterInstanceTypes := expandStringList(d.Get("master_instance_types").([]interface{}))
	workerInstanceTypes := expandStringList(d.Get("worker_instance_types").([]interface{}))
	vswitchIDs := expandStringList(d.Get("vswitch_ids").([]interface{}))
	workerNumbers := expandIntList(d.Get("worker_numbers").([]interface{}))

	vsw, err := client.DescribeVswitch(vswitchIDs[0])
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
		NodeCIDRMask:             d.Get("node_cidr_mask").(string),
		LoggingType:              loggingType,
		SLSProjectName:           slsProjectName,
		MasterSystemDiskCategory: ecs.DiskCategory(d.Get("master_disk_category").(string)),
		MasterSystemDiskSize:     int64(d.Get("master_disk_size").(int)),
		WorkerSystemDiskCategory: ecs.DiskCategory(d.Get("worker_disk_category").(string)),
		WorkerSystemDiskSize:     int64(d.Get("worker_disk_size").(int)),
		ContainerCIDR:            d.Get("pod_cidr").(string),
		ServiceCIDR:              d.Get("service_cidr").(string),
		SSHFlags:                 d.Get("enable_ssh").(bool),
		CloudMonitorFlags:        d.Get("install_cloud_monitor").(bool),
		KubernetesVersion:        d.Get("version").(string),
	}

	if v, ok := d.GetOk("worker_data_disk_category"); ok {
		creationArgs.WorkerDataDiskCategory = v.(string)
		creationArgs.WorkerDataDisk = true
		creationArgs.WorkerDataDiskSize = int64(d.Get("worker_data_disk_size").(int))
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
