package alicloud

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudVpnCustomerGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVpnCgwsRead,

		Schema: map[string]*schema.Schema{
			"customer_gateway_id": {
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

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"customer_gateways": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"customer_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudVpnCgwsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).vpcconn

	args := vpc.CreateDescribeCustomerGatewaysRequest()
	args.RegionId = string(getRegion(d, meta))
	args.PageSize = requests.NewInteger(PageSizeLarge)

	var allVpns []vpc.CustomerGateway

	for {
		resp, err := conn.DescribeCustomerGateways(args)
		if err != nil {
			return err
		}

		if resp == nil || len(resp.CustomerGateways.CustomerGateway) < 1 {
			break
		}

		allVpns = append(allVpns, resp.CustomerGateways.CustomerGateway...)

		if len(resp.CustomerGateways.CustomerGateway) < PageSizeLarge {
			break
		}

		args.PageNumber = args.PageNumber + requests.NewInteger(1)
	}

	//var filteredVpcsTemp []vpc.Vpc
	//var route_tables []string
	var filteredVpnsTemp []vpc.CustomerGateway

	for _, v := range allVpns {
		if vpnId, ok := d.GetOk("customer_gateway_id"); ok && string(v.CustomerGatewayId) != vpnId.(string) {
			continue
		}

		filteredVpnsTemp = append(filteredVpnsTemp, v)
	}

	log.Printf("[DEBUG] alicloud_vpns - VPNs found: %#v", allVpns)

	var filteredVpns []vpc.CustomerGateway

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
	return vpnCgwsDecriptionAttributes(d, filteredVpns, meta)
}

func vpnCgwsDecriptionAttributes(d *schema.ResourceData, vpnSetTypes []vpc.CustomerGateway, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, vpn := range vpnSetTypes {
		mapping := map[string]interface{}{
			"customer_gateway_id": vpn.CustomerGatewayId,
			"name":                vpn.Name,
			"ip_address":          vpn.IpAddress,
			"description":         vpn.Description,
			"create_time":         vpn.CreateTime,
		}
		log.Printf("[DEBUG] alicloud_vpn - adding vpn cgw: %v", mapping)
		ids = append(ids, vpn.CustomerGatewayId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("customer_gateways", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
