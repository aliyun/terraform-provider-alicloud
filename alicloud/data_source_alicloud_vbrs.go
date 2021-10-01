package alicloud

import (
	"regexp"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudVbrs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVbrsRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"vbr_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"physical_connection_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"access_point_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			// Computed values
			"vbrs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"activation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"termination_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"recovery_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_point_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"physical_connection_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"physical_connection_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"physical_connection_owner_uid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"physical_connection_business_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ecc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vlan_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vlan_interface_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_gateway_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"peer_gateway_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"peering_subnet_mask": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_ipv6": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"local_ipv6_gateway_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"peer_ipv6_gateway_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"peering_ipv6_subnet_mask": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"min_rx_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"min_tx_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"detect_multiplier": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func toStringArray(v []interface{}) []string {
	var arr []string
	for _, vv := range v {
		if vv == nil {
			continue
		}
		arr = append(arr, Trim(vv.(string)))
	}
	return arr
}

func dataSourceAlicloudVbrsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := vpc.CreateDescribeVirtualBorderRoutersRequest()
	request.RegionId = string(client.Region)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	var filters []vpc.DescribeVirtualBorderRoutersFilter
	if v, ok := d.GetOk("vbr_ids"); ok {
		values := toStringArray(v.([]interface{}))
		filter := vpc.DescribeVirtualBorderRoutersFilter{
			Key:   "VbrId",
			Value: &values,
		}
		filters = append(filters, filter)
	}
	if v, ok := d.GetOk("physical_connection_ids"); ok {
		values := toStringArray(v.([]interface{}))
		filter := vpc.DescribeVirtualBorderRoutersFilter{
			Key:   "PhysicalConnectionId",
			Value: &values,
		}
		filters = append(filters, filter)
	}
	if v, ok := d.GetOk("access_point_ids"); ok {
		values := toStringArray(v.([]interface{}))
		filter := vpc.DescribeVirtualBorderRoutersFilter{
			Key:   "AccessPointId",
			Value: &values,
		}
		filters = append(filters, filter)
	}
	request.Filter = &filters

	var allVbrs []vpc.VirtualBorderRouterType
	invoker := NewInvoker()
	for {
		var rawResponse interface{}
		err := invoker.Run(func() error {
			iraw, ierr := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeVirtualBorderRouters(request)
			})
			addDebug(request.GetActionName(), iraw, request.RpcRequest, request)
			rawResponse = iraw
			return ierr
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_vbrs", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		response, _ := rawResponse.(*vpc.DescribeVirtualBorderRoutersResponse)
		if len(response.VirtualBorderRouterSet.VirtualBorderRouterType) < 1 {
			break
		}

		allVbrs = append(allVbrs, response.VirtualBorderRouterSet.VirtualBorderRouterType...)

		if len(response.VirtualBorderRouterSet.VirtualBorderRouterType) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}

	var filteredVbrs []vpc.VirtualBorderRouterType
	var r *regexp.Regexp
	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}

	for _, v := range allVbrs {
		if r != nil && !r.MatchString(v.Name) {
			continue
		}

		filteredVbrs = append(filteredVbrs, v)
	}

	return vbrsDescriptionAttributes(d, filteredVbrs, meta)
}

func vbrsDescriptionAttributes(d *schema.ResourceData, vbrs []vpc.VirtualBorderRouterType, meta interface{}) error {
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, vbr := range vbrs {
		uid, _ := strconv.Atoi(vbr.PhysicalConnectionOwnerUid)

		mapping := map[string]interface{}{
			"id":                                  vbr.VbrId,
			"status":                              vbr.Status,
			"route_table_id":                      vbr.RouteTableId,
			"name":                                vbr.Name,
			"description":                         vbr.Description,
			"type":                                vbr.Type,
			"creation_time":                       vbr.CreationTime,
			"activation_time":                     vbr.ActivationTime,
			"termination_time":                    vbr.TerminationTime,
			"recovery_time":                       vbr.RecoveryTime,
			"access_point_id":                     vbr.AccessPointId,
			"physical_connection_id":              vbr.PhysicalConnectionId,
			"physical_connection_status":          vbr.PhysicalConnectionStatus,
			"physical_connection_owner_uid":       uid,
			"physical_connection_business_status": vbr.PhysicalConnectionBusinessStatus,
			"ecc_id":                              vbr.EccId,
			"vlan_id":                             vbr.VlanId,
			"vlan_interface_id":                   vbr.VlanInterfaceId,
			"local_gateway_ip":                    vbr.LocalGatewayIp,
			"peer_gateway_ip":                     vbr.PeerGatewayIp,
			"peering_subnet_mask":                 vbr.PeeringSubnetMask,
			"enable_ipv6":                         vbr.EnableIpv6,
			"local_ipv6_gateway_ip":               vbr.LocalIpv6GatewayIp,
			"peer_ipv6_gateway_ip":                vbr.PeerIpv6GatewayIp,
			"peering_ipv6_subnet_mask":            vbr.PeeringIpv6SubnetMask,
			"min_rx_interval":                     vbr.MinRxInterval,
			"min_tx_interval":                     vbr.MinTxInterval,
			"detect_multiplier":                   vbr.DetectMultiplier,
		}
		ids = append(ids, vbr.VbrId)
		names = append(names, vbr.Name)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("vbrs", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
