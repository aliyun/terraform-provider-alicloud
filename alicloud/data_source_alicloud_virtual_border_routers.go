package alicloud

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudVirtualBorderRouters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVirtualBorderRoutersRead,

		Schema: map[string]*schema.Schema{
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"physical_connection_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"physical_connection_owner_uid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				MinItems: 1,
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
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vlan_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"route_table_id": {
							Type:     schema.TypeString,
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
						"physical_connection_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"physical_connection_owner_uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_point_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"circuit_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudVirtualBorderRoutersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := vpc.CreateDescribeVirtualBorderRoutersRequest()
	request.RegionId = string(client.Region)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var filters []vpc.DescribeVirtualBorderRoutersFilter
	for _, key := range []string{"status", "physical_connection_id", "physical_connection_owner_uid"} {
		if v, ok := d.GetOk(key); ok && v.(string) != "" {
			value := []string{v.(string)}
			filters = append(filters, vpc.DescribeVirtualBorderRoutersFilter{
				Key:   terraformToAPI(key),
				Value: &value,
			})
		}
	}

	request.Filter = &filters

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	var allVirtualBorderRouters []vpc.VirtualBorderRouterType
	invoker := NewInvoker()

	for {
		var response *vpc.DescribeVirtualBorderRoutersResponse
		if err := invoker.Run(func() error {
			raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeVirtualBorderRouters(request)
			})
			if err != nil {
				return err
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			response, _ = raw.(*vpc.DescribeVirtualBorderRoutersResponse)
			return nil
		}); err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_virtual_border_routers", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		if len(response.VirtualBorderRouterSet.VirtualBorderRouterType) < 1 {
			break
		}

		allVirtualBorderRouters = append(allVirtualBorderRouters, response.VirtualBorderRouterSet.VirtualBorderRouterType...)

		if len(response.VirtualBorderRouterSet.VirtualBorderRouterType) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}

	}

	var filteredVirtualBorderRouters []vpc.VirtualBorderRouterType
	var r *regexp.Regexp
	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}

	for _, v := range allVirtualBorderRouters {
		if len(idsMap) > 0 {
			if _, ok := idsMap[v.VbrId]; !ok {
				continue
			}
		}
		if r != nil && !r.MatchString(v.Name) {
			continue
		}
		filteredVirtualBorderRouters = append(filteredVirtualBorderRouters, v)
	}

	return vbrDescriptionAttributes(d, filteredVirtualBorderRouters, meta)
}

func vbrDescriptionAttributes(d *schema.ResourceData, vbrSetTypes []vpc.VirtualBorderRouterType, meta interface{}) error {
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, vbr := range vbrSetTypes {
		mapping := map[string]interface{}{
			"id":                            vbr.VbrId,
			"status":                        vbr.Status,
			"name":                          vbr.Name,
			"description":                   vbr.Description,
			"vlan_id":                       vbr.VlanId,
			"route_table_id":                vbr.RouteTableId,
			"vlan_interface_id":             vbr.VlanInterfaceId,
			"local_gateway_ip":              vbr.LocalGatewayIp,
			"peer_gateway_ip":               vbr.PeerGatewayIp,
			"peering_subnet_mask":           vbr.PeeringSubnetMask,
			"physical_connection_id":        vbr.PhysicalConnectionId,
			"physical_connection_owner_uid": vbr.PhysicalConnectionOwnerUid,
			"access_point_id":               vbr.AccessPointId,
			"creation_time":                 vbr.CreationTime,
			"circuit_code":                  vbr.CircuitCode,
		}
		ids = append(ids, vbr.VbrId)
		names = append(names, vbr.Name)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("vbrs", s); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
