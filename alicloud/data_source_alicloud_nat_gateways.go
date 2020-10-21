package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudNatGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudNatGatewaysRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				Computed: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			// Computed values
			"gateways": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec": {
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
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snat_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"forward_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_lists": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}
func dataSourceAlicloudNatGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := vpc.CreateDescribeNatGatewaysRequest()
	request.RegionId = string(client.Region)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	request.VpcId = d.Get("vpc_id").(string)
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	var allNatGateways []vpc.NatGateway
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(Trim(v.(string))); err == nil {
			nameRegex = r
		} else {
			return WrapError(err)
		}
	}
	invoker := NewInvoker()
	for {
		var raw interface{}
		if err := invoker.Run(func() error {
			response, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeNatGateways(request)
			})
			raw = response
			return err
		}); err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_nat_gateways", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeNatGatewaysResponse)
		if len(response.NatGateways.NatGateway) < 1 {
			break
		}

		for _, gateways := range response.NatGateways.NatGateway {
			if nameRegex != nil {
				if !nameRegex.MatchString(gateways.Name) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[gateways.NatGatewayId]; !ok {
					continue
				}
			}
			allNatGateways = append(allNatGateways, gateways)
		}

		if len(response.NatGateways.NatGateway) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	return NatGatewaysDecriptionAttributes(d, allNatGateways, meta)
}

func NatGatewaysDecriptionAttributes(d *schema.ResourceData, gateways []vpc.NatGateway, meta interface{}) error {
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, gateway := range gateways {
		var ips []string
		for _, ip := range gateway.IpLists.IpList {
			ips = append(ips, ip.IpAddress)
		}
		mapping := map[string]interface{}{
			"id":               gateway.NatGatewayId,
			"spec":             gateway.Spec,
			"status":           gateway.Status,
			"name":             gateway.Name,
			"description":      gateway.Description,
			"creation_time":    gateway.CreationTime,
			"snat_table_id":    gateway.SnatTableIds.SnatTableId[0],
			"forward_table_id": gateway.ForwardTableIds.ForwardTableId[0],
			"vpc_id":           gateway.VpcId,
			"ip_lists":         ips,
		}
		names = append(names, gateway.Name)
		ids = append(ids, gateway.NatGatewayId)
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("gateways", s); err != nil {
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
