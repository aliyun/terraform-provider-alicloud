package alicloud

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudCSKubernetesClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCSKubernetesClustersRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"clusters": {
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
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slb_internet_enabled": {
							Type:     schema.TypeBool,
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
						"vswitch_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"master_instance_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"worker_instance_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"worker_numbers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"key_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pod_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_cidr_mask": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"log_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"project": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"master_disk_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"worker_disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"worker_disk_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"worker_data_disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"worker_data_disk_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_instance_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_period_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"master_auto_renew": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"master_auto_renew_period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"worker_instance_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"worker_period_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"worker_period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"worker_auto_renew": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"worker_auto_renew_period": {
							Type:     schema.TypeInt,
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
					},
				},
			},
		},
	}
}
func dataSourceAlicloudCSKubernetesClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var allClusterTypes []cs.ClusterType

	invoker := NewInvoker()
	if err := invoker.Run(func() error {
		raw, e := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return csClient.DescribeClusters("")
		})
		if e != nil {
			return fmt.Errorf("Describing cluster failed, error message: %#v.", e)
		}
		allClusterTypes, _ = raw.([]cs.ClusterType)
		return nil
	}); err != nil {
		return err
	}

	var filteredClusterTypes []cs.ClusterType
	for _, v := range allClusterTypes {
		if v.ClusterType != cs.ClusterTypeKubernetes {
			continue
		}
		if client.RegionId != string(v.RegionID) {
			continue
		}
		if nameRegex, ok := d.GetOk("name_regex"); ok {
			r := regexp.MustCompile(nameRegex.(string))
			if !r.MatchString(v.Name) {
				continue
			}
		}
		if ids, ok := d.GetOk("ids"); ok {
			findId := func(id string, ids []string) (ret bool) {
				for _, i := range ids {
					if id == i {
						ret = true
					}
				}
				return
			}
			if !findId(v.ClusterID, expandStringList(ids.([]interface{}))) {
				continue
			}
		}
		filteredClusterTypes = append(filteredClusterTypes, v)
	}

	var filteredKubernetesCluster []cs.KubernetesCluster

	for _, v := range filteredClusterTypes {
		var kubernetesCluster cs.KubernetesCluster

		invoker := NewInvoker()
		if err := invoker.Run(func() error {
			raw, e := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				return csClient.DescribeKubernetesCluster(v.ClusterID)
			})
			if e != nil {
				return fmt.Errorf("Describing kubernetes cluster %#v failed, error message: %#v. Please check cluster in the console,", v.ClusterID, e)
			}
			kubernetesCluster = raw.(cs.KubernetesCluster)
			return nil
		}); err != nil {
			return err
		}

		if az, ok := d.GetOk("availability_zone"); ok && az != kubernetesCluster.ZoneId {
			continue
		}

		filteredKubernetesCluster = append(filteredKubernetesCluster, kubernetesCluster)
	}

	return csKubernetesClusterDescriptionAttributes(d, filteredKubernetesCluster, meta)
}

func csKubernetesClusterDescriptionAttributes(d *schema.ResourceData, clusterTypes []cs.KubernetesCluster, meta interface{}) error {
	var ids, names []string
	var s []map[string]interface{}
	for _, ct := range clusterTypes {
		mapping := map[string]interface{}{
			"id":   ct.ClusterID,
			"name": ct.Name,
		}

		if detailedEnabled, ok := d.GetOk("enable_details"); ok && !detailedEnabled.(bool) {
			ids = append(ids, ct.ClusterID)
			names = append(names, ct.Name)
			s = append(s, mapping)
			continue
		}

		mapping["vpc_id"] = ct.VPCID
		mapping["security_group_id"] = ct.SecurityGroupID
		mapping["availability_zone"] = ct.ZoneId
		mapping["key_name"] = ct.Parameters.KeyPair
		mapping["master_disk_category"] = ct.Parameters.MasterSystemDiskCategory
		mapping["worker_disk_category"] = ct.Parameters.WorkerSystemDiskCategory
		if ct.Parameters.PublicSLB != nil {
			mapping["slb_internet_enabled"] = *ct.Parameters.PublicSLB
		}

		if ct.Parameters.ImageId != "" {
			mapping["image_id"] = ct.Parameters.ImageId
		} else {
			mapping["image_id"] = ct.Parameters.MasterImageId
		}

		if size, err := strconv.Atoi(ct.Parameters.MasterSystemDiskSize); err != nil {
			return BuildWrapError("strconv.Atoi", d.Id(), ProviderERROR, err, "")
		} else {
			mapping["master_disk_size"] = size
		}

		if size, err := strconv.Atoi(ct.Parameters.WorkerSystemDiskSize); err != nil {
			return BuildWrapError("strconv.Atoi", d.Id(), ProviderERROR, err, "")
		} else {
			mapping["worker_disk_size"] = size
		}

		if ct.Parameters.MasterInstanceChargeType == string(PrePaid) {
			mapping["master_instance_charge_type"] = string(PrePaid)
			if period, err := strconv.Atoi(ct.Parameters.MasterPeriod); err != nil {
				return BuildWrapError("strconv.Atoi", d.Id(), ProviderERROR, err, "")
			} else {
				mapping["master_period"] = period
			}
			mapping["master_period_unit"] = ct.Parameters.MasterPeriodUnit
			if ct.Parameters.MasterAutoRenew != nil {
				mapping["master_auto_renew"] = *ct.Parameters.MasterAutoRenew
			}
			if period, err := strconv.Atoi(ct.Parameters.MasterAutoRenewPeriod); err != nil {
				return BuildWrapError("strconv.Atoi", d.Id(), ProviderERROR, err, "")
			} else {
				mapping["master_auto_renew_period"] = period
			}
		} else {
			mapping["master_instance_charge_type"] = string(PostPaid)
		}

		if ct.Parameters.WorkerInstanceChargeType == string(PrePaid) {
			mapping["worker_instance_charge_type"] = string(PrePaid)
			if period, err := strconv.Atoi(ct.Parameters.WorkerPeriod); err != nil {
				return BuildWrapError("strconv.Atoi", d.Id(), ProviderERROR, err, "")
			} else {
				mapping["worker_period"] = period
			}
			mapping["worker_period_unit"] = ct.Parameters.WorkerPeriodUnit
			if ct.Parameters.WorkerAutoRenew != nil {
				mapping["worker_auto_renew"] = *ct.Parameters.WorkerAutoRenew
			}
			if period, err := strconv.Atoi(ct.Parameters.WorkerAutoRenewPeriod); err != nil {
				return BuildWrapError("strconv.Atoi", d.Id(), ProviderERROR, err, "")
			} else {
				mapping["worker_auto_renew_period"] = period
			}
		} else {
			mapping["worker_instance_charge_type"] = string(PostPaid)
		}

		if cidrMask, err := strconv.Atoi(ct.Parameters.NodeCIDRMask); err == nil {
			mapping["node_cidr_mask"] = cidrMask
		} else {
			return BuildWrapError("strconv.Atoi", d.Id(), ProviderERROR, err, "")
		}

		if ct.Parameters.WorkerDataDisk != nil && *ct.Parameters.WorkerDataDisk {
			if size, err := strconv.Atoi(ct.Parameters.WorkerDataDiskSize); err != nil {
				return BuildWrapError("strconv.Atoi", d.Id(), ProviderERROR, err, "")
			} else {
				mapping["worker_data_disk_size"] = size
			}
			mapping["worker_data_disk_category"] = ct.Parameters.WorkerDataDiskCategory
		}

		if ct.Parameters.LoggingType != "None" {
			logConfig := map[string]interface{}{}
			logConfig["type"] = ct.Parameters.LoggingType
			if ct.Parameters.SLSProjectName == "None" {
				logConfig["project"] = ""
			} else {
				logConfig["project"] = ct.Parameters.SLSProjectName
			}
			mapping["log_config"] = []map[string]interface{}{logConfig}
		}

		// Each k8s cluster contains 3 master nodes
		if ct.MetaData.MultiAZ || ct.MetaData.SubClass == "3az" {
			numOfNodeA, err := strconv.Atoi(ct.Parameters.NumOfNodesA)
			if err != nil {
				return fmt.Errorf("error convert NumOfNodesA %s to int: %s", ct.Parameters.NumOfNodesA, err.Error())
			}
			numOfNodeB, err := strconv.Atoi(ct.Parameters.NumOfNodesB)
			if err != nil {
				return fmt.Errorf("error convert NumOfNodesB %s to int: %s", ct.Parameters.NumOfNodesB, err.Error())
			}
			numOfNodeC, err := strconv.Atoi(ct.Parameters.NumOfNodesC)
			if err != nil {
				return fmt.Errorf("error convert NumOfNodesC %s to int: %s", ct.Parameters.NumOfNodesC, err.Error())
			}
			mapping["worker_numbers"] = []int{numOfNodeA, numOfNodeB, numOfNodeC}
			mapping["vswitch_ids"] = []string{ct.Parameters.VSwitchIdA, ct.Parameters.VSwitchIdB, ct.Parameters.VSwitchIdC}
			mapping["master_instance_types"] = []string{ct.Parameters.MasterInstanceTypeA, ct.Parameters.MasterInstanceTypeB, ct.Parameters.MasterInstanceTypeC}
			mapping["worker_instance_types"] = []string{ct.Parameters.WorkerInstanceTypeA, ct.Parameters.WorkerInstanceTypeB, ct.Parameters.WorkerInstanceTypeC}
		} else {
			if numOfNode, err := strconv.Atoi(ct.Parameters.NumOfNodes); err != nil {
				return BuildWrapError("strconv.Atoi", d.Id(), ProviderERROR, err, "")
			} else {
				mapping["worker_numbers"] = []int{numOfNode}
			}
			mapping["vswitch_ids"] = []string{ct.Parameters.VSwitchID}
			mapping["master_instance_types"] = []string{ct.Parameters.MasterInstanceType}
			mapping["worker_instance_types"] = []string{ct.Parameters.WorkerInstanceType}
		}

		var masterNodes []map[string]interface{}
		var workerNodes []map[string]interface{}

		invoker := NewInvoker()
		client := meta.(*connectivity.AliyunClient)
		pageNumber := 1
		for {
			var result []cs.KubernetesNodeType
			var pagination *cs.PaginationResult

			if err := invoker.Run(func() error {
				raw, e := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
					nodes, paginationResult, err := csClient.GetKubernetesClusterNodes(ct.ClusterID, common.Pagination{PageNumber: pageNumber, PageSize: PageSizeLarge})
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
							nodes, _, err := csClient.GetKubernetesClusterNodes(ct.ClusterID, common.Pagination{PageNumber: pageNumber, PageSize: PageSizeLarge})
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
						if ct.State == stableState {
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
				subMapping := map[string]interface{}{
					"id":         node.InstanceId,
					"name":       node.InstanceName,
					"private_ip": node.IpAddress[0],
				}
				if node.InstanceRole == "Master" {
					masterNodes = append(masterNodes, subMapping)
				} else {
					workerNodes = append(workerNodes, subMapping)
				}
			}

			if len(result) < pagination.PageSize {
				break
			}
			pageNumber += 1
		}
		mapping["master_nodes"] = masterNodes
		mapping["worker_nodes"] = workerNodes

		if err := invoker.Run(func() error {
			rawEndpoints, e := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				endpoints, err := csClient.GetClusterEndpoints(ct.ClusterID)
				return endpoints, err
			})
			if e != nil {
				return e
			}
			connection := make(map[string]string)
			if endpoints, ok := rawEndpoints.(cs.ClusterEndpoints); ok && endpoints.ApiServerEndpoint != "" {
				connection["api_server_internet"] = endpoints.ApiServerEndpoint
				connection["master_public_ip"] = strings.TrimSuffix(strings.TrimPrefix(endpoints.ApiServerEndpoint, "https://"), ":6443")
			}
			if endpoints, ok := rawEndpoints.(cs.ClusterEndpoints); ok && endpoints.IntranetApiServerEndpoint != "" {
				connection["api_server_intranet"] = endpoints.IntranetApiServerEndpoint
			}
			connection["service_domain"] = fmt.Sprintf("*.%s.%s.alicontainer.com", ct.ClusterID, ct.RegionID)

			mapping["connections"] = connection
			return nil
		}); err != nil {
			return fmt.Errorf("[ERROR] GetKubernetesClusterNodes got an error: %#v.", err)
		}

		req := vpc.CreateDescribeNatGatewaysRequest()
		req.VpcId = ct.VPCID
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeNatGateways(req)
		})
		if err != nil {
			return fmt.Errorf("[ERROR] DescribeNatGateways by VPC Id %s: %#v.", ct.VPCID, err)
		}
		nat, _ := raw.(*vpc.DescribeNatGatewaysResponse)
		if nat != nil && len(nat.NatGateways.NatGateway) > 0 {
			mapping["nat_gateway_id"] = nat.NatGateways.NatGateway[0].NatGatewayId
		}

		ids = append(ids, ct.ClusterID)
		names = append(names, ct.Name)
		s = append(s, mapping)
	}

	d.Set("ids", ids)
	d.Set("names", names)
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("clusters", s); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
