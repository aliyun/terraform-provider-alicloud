package cs

import (
	"regexp"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/denverdino/aliyungo/common"
	aliyungocs "github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceCSEdgeKubernetesClusters() *schema.Resource {
	return &schema.Resource{
		Read:               dataSourceCSEdgeKubernetesClustersRead,
		DeprecationMessage: "This data source has been deprecated since v1.276.0 and will be removed in the future. Please use 'alicloud_cs_clusters' instead.",

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
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
								},
							},
						},
					},
				},
			},
		},
	}
}
func dataSourceCSEdgeKubernetesClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var allClusterTypes []aliyungocs.ClusterType

	var requestInfo *aliyungocs.Client
	invoker := helper.NewInvoker()
	var response interface{}
	if err := invoker.Run(func() error {
		raw, err := client.WithCsClient(func(csClient *aliyungocs.Client) (interface{}, error) {
			requestInfo = csClient
			return csClient.DescribeClusters("")
		})
		response = raw
		return err
	}); err != nil {
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cs_edge_kubernetes_clusters", "DescribeClusters", errs.DenverdinoAliyungo)
	}
	if logs.DebugOn() {
		logs.AddDebug("DescribeClusters", response, requestInfo)
	}
	allClusterTypes, _ = response.([]aliyungocs.ClusterType)

	var filteredClusterTypes []aliyungocs.ClusterType
	for _, v := range allClusterTypes {
		if v.ClusterType != aliyungocs.ClusterTypeManagedKubernetes {
			continue
		}
		if client.RegionId != string(v.RegionID) {
			continue
		}
		if nameRegex, ok := d.GetOk("name_regex"); ok {
			r, err := regexp.Compile(nameRegex.(string))
			if err != nil {
				return errs.WrapError(err)
			}
			if !r.MatchString(v.Name) {
				continue
			}
		}
		if ids, ok := d.GetOk("ids"); ok {
			var found bool
			for _, i := range convert.ExpandStringList(ids.([]interface{})) {
				if v.ClusterID == i {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		filteredClusterTypes = append(filteredClusterTypes, v)
	}

	var filteredKubernetesCluster []aliyungocs.KubernetesCluster

	for _, v := range filteredClusterTypes {
		var kubernetesCluster aliyungocs.KubernetesCluster

		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *aliyungocs.Client) (interface{}, error) {
				requestInfo = csClient
				return csClient.DescribeKubernetesCluster(v.ClusterID)
			})
			response = raw
			return err
		}); err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cs_edge_kubernetes_clusters", "DescribeKubernetesCluster", errs.DenverdinoAliyungo)
		}
		if logs.DebugOn() {
			requestMap := make(map[string]interface{})
			requestMap["Id"] = v.ClusterID
			logs.AddDebug("DescribeKubernetesCluster", response, requestInfo, requestMap)
		}
		kubernetesCluster = response.(aliyungocs.KubernetesCluster)

		if az, ok := d.GetOk("availability_zone"); ok && az != kubernetesCluster.ZoneId {
			continue
		}

		filteredKubernetesCluster = append(filteredKubernetesCluster, kubernetesCluster)
	}

	return csEdgeKubernetesClusterDescriptionAttributes(d, filteredKubernetesCluster, meta)
}

func csEdgeKubernetesClusterDescriptionAttributes(d *schema.ResourceData, clusterTypes []aliyungocs.KubernetesCluster, meta interface{}) error {
	detailEnabled := false
	if v, ok := d.GetOk("enable_details"); ok {
		detailEnabled = v.(bool)
	}

	var ids, names []string
	var s []map[string]interface{}
	for _, ct := range clusterTypes {
		mapping := map[string]interface{}{
			"id":   ct.ClusterID,
			"name": ct.Name,
		}

		if !detailEnabled {
			ids = append(ids, ct.ClusterID)
			names = append(names, ct.Name)
			s = append(s, mapping)
			continue
		}

		mapping["vpc_id"] = ct.VPCID
		mapping["security_group_id"] = ct.SecurityGroupID
		mapping["availability_zone"] = ct.ZoneId

		var workerNodes []map[string]interface{}

		invoker := helper.NewInvoker()
		client := meta.(*connectivity.AliyunClient)
		pageNumber := 1
		for {
			var result []aliyungocs.KubernetesNodeType
			var pagination *aliyungocs.PaginationResult
			var requestInfo *aliyungocs.Client
			var response interface{}

			if err := invoker.Run(func() error {
				raw, err := client.WithCsClient(func(csClient *aliyungocs.Client) (interface{}, error) {
					requestInfo = csClient
					nodes, paginationResult, err := csClient.GetKubernetesClusterNodes(ct.ClusterID, common.Pagination{PageNumber: pageNumber, PageSize: helper.PageSizeLarge}, "")
					return []interface{}{nodes, paginationResult}, err
				})
				response = raw
				return err
			}); err != nil {
				return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cs_edge_kubernetes_clusters", "GetKubernetesClusterNodes", errs.DenverdinoAliyungo)
			}
			if logs.DebugOn() {
				requestMap := make(map[string]interface{})
				requestMap["Id"] = ct.ClusterID
				requestMap["Pagination"] = common.Pagination{PageNumber: pageNumber, PageSize: helper.PageSizeLarge}
				logs.AddDebug("GetKubernetesClusterNodes", response, requestInfo, requestMap)
			}
			result, _ = response.([]interface{})[0].([]aliyungocs.KubernetesNodeType)
			pagination, _ = response.([]interface{})[1].(*aliyungocs.PaginationResult)

			if pageNumber == 1 && (len(result) == 0 || result[0].InstanceId == "") {
				err := retry.Retry(5*time.Minute, func() *retry.RetryError {
					if err := invoker.Run(func() error {
						raw, err := client.WithCsClient(func(csClient *aliyungocs.Client) (interface{}, error) {
							requestInfo = csClient
							nodes, _, err := csClient.GetKubernetesClusterNodes(ct.ClusterID, common.Pagination{PageNumber: pageNumber, PageSize: helper.PageSizeLarge}, "")
							return nodes, err
						})
						response = raw
						return err
					}); err != nil {
						return retry.NonRetryableError(err)
					}
					if logs.DebugOn() {
						requestMap := make(map[string]interface{})
						requestMap["Id"] = ct.ClusterID
						requestMap["Pagination"] = common.Pagination{PageNumber: pageNumber, PageSize: helper.PageSizeLarge}
						logs.AddDebug("GetKubernetesClusterNodes", response, requestInfo, requestMap)
					}
					tmp, _ := response.([]aliyungocs.KubernetesNodeType)
					if len(tmp) > 0 && tmp[0].InstanceId != "" {
						result = tmp
					}
					for _, stableState := range aliyungocs.NodeStableClusterState {
						// If cluster is in NodeStableClusteState, node list will not change
						if ct.State == stableState {
							return nil
						}
					}
					time.Sleep(5 * time.Second)
					return retry.RetryableError(errs.Error("there is no any nodes in kubernetes cluster %s", d.Id()))
				})
				if err != nil {
					return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cs_managed_kubernetes_clusters", "GetKubernetesClusterNodes", errs.DenverdinoAliyungo)
				}

			}

			for _, node := range result {
				subMapping := map[string]interface{}{
					"id":         node.InstanceId,
					"name":       node.InstanceName,
					"private_ip": CsNodePrivateIp(node.IpAddress),
				}
				workerNodes = append(workerNodes, subMapping)
			}

			if len(result) < pagination.PageSize {
				break
			}
			pageNumber += 1
		}
		mapping["worker_nodes"] = workerNodes

		var response interface{}
		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *aliyungocs.Client) (interface{}, error) {
				endpoints, err := csClient.GetClusterEndpoints(ct.ClusterID)
				return endpoints, err
			})
			response = raw
			return err
		}); err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cs_edge_kubernetes_clusters", "GetClusterEndpoints", errs.DenverdinoAliyungo)
		}
		connection := make(map[string]string)
		if endpoints, ok := response.(aliyungocs.ClusterEndpoints); ok && endpoints.ApiServerEndpoint != "" {
			connection["api_server_internet"] = endpoints.ApiServerEndpoint
		}
		if endpoints, ok := response.(aliyungocs.ClusterEndpoints); ok && endpoints.IntranetApiServerEndpoint != "" {
			connection["api_server_intranet"] = endpoints.IntranetApiServerEndpoint
		}
		mapping["connections"] = []interface{}{connection}

		request := vpc.CreateDescribeNatGatewaysRequest()
		request.VpcId = ct.VPCID

		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeNatGateways(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, ct.VPCID, "DescribeNatGateways", errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug("DescribeNatGateways", raw, request.RpcRequest, request)
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
	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("clusters", s); err != nil {
		return errs.WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
