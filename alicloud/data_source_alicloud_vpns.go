package alicloud

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudVpns() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVpnsRead,

		Schema: map[string]*schema.Schema{
			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"internet_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
				ForceNew:     true,
			},

			"status": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"business_status": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"vpn_gateways": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpn_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"spec": {
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
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"business_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipsec_vpn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssl_vpn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssl_max_connections": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudVpnsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).vpcconn

	args := vpc.CreateDescribeVpnGatewaysRequest()
	args.RegionId = string(getRegion(d, meta))
	args.PageSize = requests.NewInteger(PageSizeLarge)

	var allVpns []vpc.VpnGateway

	for {
		resp, err := conn.DescribeVpnGateways(args)
		if err != nil {
			return err
		}

		if resp == nil || len(resp.VpnGateways.VpnGateway) < 1 {
			break
		}

		allVpns = append(allVpns, resp.VpnGateways.VpnGateway...)

		if len(resp.VpnGateways.VpnGateway) < PageSizeLarge {
			break
		}

		args.PageNumber = args.PageNumber + requests.NewInteger(1)
	}

	//var filteredVpcsTemp []vpc.Vpc
	//var route_tables []string
	var filteredVpnsTemp []vpc.VpnGateway

	for _, v := range allVpns {
		if vpnId, ok := d.GetOk("vpn_gateway_id"); ok && string(v.VpnGatewayId) != vpnId.(string) {
			continue
		}

		if vpcId, ok := d.GetOk("vpc_id"); ok && string(v.VpcId) != vpcId.(string) {
			continue
		}

		if status, ok := d.GetOk("status"); ok && string(v.Status) != status.(string) {
			continue
		}

		filteredVpnsTemp = append(filteredVpnsTemp, v)
	}

	var filteredVpns []vpc.VpnGateway

	if nameRegex, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(nameRegex.(string)); err == nil {
			for _, vpn := range filteredVpnsTemp {
				if r.MatchString(vpn.Name) {
					filteredVpns = append(filteredVpns, vpn)
				}
			}
		}
	} else {
		filteredVpns = filteredVpnsTemp[:]
	}

	if len(filteredVpns) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}
	log.Printf("[DEBUG] alicloud_vpns VPNs filter: %#v", filteredVpns)
	return vpnsDecriptionAttributes(d, filteredVpns, meta)
}

func vpnsDecriptionAttributes(d *schema.ResourceData, vpnSetTypes []vpc.VpnGateway, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, vpn := range vpnSetTypes {
		mapping := map[string]interface{}{
			"vpn_gateway_id":      vpn.VpnGatewayId,
			"vpc_id":              vpn.VpcId,
			"vswitch_id":          vpn.VSwitchId,
			"internet_ip":         vpn.InternetIp,
			"name":                vpn.Name,
			"spec":                vpn.Spec,
			"description":         vpn.Description,
			"status":              vpn.Status,
			"ipsec_vpn":           vpn.IpsecVpn,
			"ssl_vpn":             vpn.SslVpn,
			"ssl_max_connections": vpn.SslMaxConnections,
			"create_time":         vpn.CreateTime,
		}
		log.Printf("[DEBUG] alicloud_vpn - adding vpn: %v", mapping)
		ids = append(ids, vpn.VpnGatewayId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("vpn_gateways", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
