package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/mse"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudMseClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudMseClustersRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"cluster_alias_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"request_pars": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"DESTROY_FAILED", "DESTROY_ING", "DESTROY_SUCCESS", "INIT_FAILED", "INIT_ING", "INIT_SUCCESS", "INIT_TIME_OUT", "RESTART_FAILED", "RESTART_ING", "RESTART_SUCCESS", "SCALE_FAILED", "SCALE_ING", "SCALE_SUCCESS"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"acl_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"health_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"init_cost_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_models": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"health_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"internet_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"pod_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"role": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"single_tunnel_vip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vip": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"internet_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"intranet_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"intranet_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"intranet_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"pay_info": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pub_network_flow": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAlicloudMseClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := mse.CreateListClustersRequest()
	if v, ok := d.GetOk("cluster_alias_name"); ok {
		request.ClusterAliasName = v.(string)
	}
	request.RegionId = client.RegionId
	if v, ok := d.GetOk("request_pars"); ok {
		request.RequestPars = v.(string)
	}
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNum = requests.NewInteger(1)
	var objects []mse.ClusterForListModel
	var clusterNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		clusterNameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	status, statusOk := d.GetOk("status")
	var response *mse.ListClustersResponse
	for {
		raw, err := client.WithMseClient(func(mseClient *mse.Client) (interface{}, error) {
			return mseClient.ListClusters(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_mse_clusters", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ = raw.(*mse.ListClustersResponse)

		for _, item := range response.Data {
			if clusterNameRegex != nil {
				if !clusterNameRegex.MatchString(item.ClusterAliasName) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[item.InstanceId]; !ok {
					continue
				}
			}
			if statusOk && status != "" && status != item.InitStatus {
				continue
			}
			objects = append(objects, item)
		}
		if len(response.Data) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNum)
		if err != nil {
			return WrapError(err)
		}
		request.PageNum = page
	}
	ids := make([]string, len(objects))
	names := make([]string, len(objects))
	s := make([]map[string]interface{}, len(objects))
	for i, object := range objects {
		mapping := map[string]interface{}{
			"app_version":      object.AppVersion,
			"cluster_id":       object.ClusterId,
			"cluster_name":     object.ClusterAliasName,
			"cluster_type":     object.ClusterType,
			"id":               object.InstanceId,
			"instance_id":      object.InstanceId,
			"internet_address": object.InternetAddress,
			"internet_domain":  object.InternetDomain,
			"intranet_address": object.IntranetAddress,
			"intranet_domain":  object.IntranetDomain,
			"status":           object.InitStatus,
		}
		ids[i] = object.InstanceId
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			names[i] = object.ClusterAliasName
			s[i] = mapping
			continue
		}
		request := mse.CreateQueryClusterDetailRequest()
		request.RegionId = client.RegionId
		request.InstanceId = object.InstanceId
		raw, err := client.WithMseClient(func(mseClient *mse.Client) (interface{}, error) {
			return mseClient.QueryClusterDetail(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_mse_clusters", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		responseGet, _ := raw.(*mse.QueryClusterDetailResponse)
		mapping["acl_id"] = responseGet.Data.AclId
		mapping["cpu"] = responseGet.Data.Cpu
		mapping["health_status"] = responseGet.Data.HealthStatus
		mapping["init_cost_time"] = responseGet.Data.InitCostTime
		mapping["instance_count"] = responseGet.Data.InstanceCount

		instanceModels := make([]map[string]string, len(responseGet.Data.InstanceModels))
		for i, v := range responseGet.Data.InstanceModels {
			mapping1 := map[string]string{
				"health_status":     v.HealthStatus,
				"instance_type":     v.InstanceType,
				"internet_ip":       v.InternetIp,
				"ip":                v.Ip,
				"pod_name":          v.PodName,
				"role":              v.Role,
				"single_tunnel_vip": v.SingleTunnelVip,
				"vip":               v.Vip,
			}
			instanceModels[i] = mapping1
		}
		mapping["instance_models"] = instanceModels
		mapping["internet_port"] = responseGet.Data.InternetPort
		mapping["intranet_port"] = responseGet.Data.IntranetPort
		mapping["memory_capacity"] = responseGet.Data.MemoryCapacity
		mapping["pay_info"] = responseGet.Data.PayInfo
		mapping["pub_network_flow"] = responseGet.Data.PubNetworkFlow
		names[i] = object.ClusterAliasName
		s[i] = mapping
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("clusters", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
