package alicloud

import (
	"fmt"
	"log"
	"regexp"

	"encoding/json"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudVpnConnections() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVpnConnectionsRead,

		Schema: map[string]*schema.Schema{
			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"customer_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"vpn_connection_id": {
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
			"vpn_connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpn_connection_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"customer_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpn_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"effect_immediately": {
							Type:     schema.TypeBool,
							Computed: true,
						},

						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"ike_config": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"ipsec_config": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudVpnConnectionsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).vpcconn

	args := vpc.CreateDescribeVpnConnectionsRequest()
	args.RegionId = string(getRegion(d, meta))
	args.PageSize = requests.NewInteger(PageSizeLarge)

	var allVpns []vpc.VpnConnection

	for {
		resp, err := conn.DescribeVpnConnections(args)
		if err != nil {
			return err
		}

		if resp == nil || len(resp.VpnConnections.VpnConnection) < 1 {
			break
		}

		allVpns = append(allVpns, resp.VpnConnections.VpnConnection...)

		if len(resp.VpnConnections.VpnConnection) < PageSizeLarge {
			break
		}

		args.PageNumber = args.PageNumber + requests.NewInteger(1)
	}

	//var filteredVpcsTemp []vpc.Vpc
	//var route_tables []string
	var filteredVpnsTemp []vpc.VpnConnection

	for _, v := range allVpns {
		if vpnId, ok := d.GetOk("vpn_gateway_id"); ok && string(v.VpnGatewayId) != vpnId.(string) {
			continue
		}

		if cgwId, ok := d.GetOk("customer_gateway_id"); ok && string(v.CustomerGatewayId) != cgwId.(string) {
			continue
		}

		if vpnConnId, ok := d.GetOk("vpn_connection_id"); ok && string(v.VpnConnectionId) != vpnConnId.(string) {
			continue
		}

		filteredVpnsTemp = append(filteredVpnsTemp, v)
	}

	var filteredVpns []vpc.VpnConnection

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
	return vpnConnectionsDecriptionAttributes(d, filteredVpns, meta)
}

func vpnConnectionsDecriptionAttributes(d *schema.ResourceData, vpnSetTypes []vpc.VpnConnection, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, vpn := range vpnSetTypes {
		ikeConfig, _ := json.Marshal(vpn.IkeConfig)
		ipsecConfig, _ := json.Marshal(vpn.IpsecConfig)
		mapping := map[string]interface{}{
			"vpn_connection_id":   vpn.VpnConnectionId,
			"customer_gateway_id": vpn.CustomerGatewayId,
			"vpn_gateway_id":      vpn.VpnGatewayId,
			"name":                vpn.Name,
			"local_subnet":        vpn.LocalSubnet,
			"remote_subnet":       vpn.RemoteSubnet,
			"create_time":         vpn.CreateTime,
			"effect_immediately":  vpn.EffectImmediately,
			"status":              vpn.Status,
			"ike_config":          string(ikeConfig),
			"ipsec_config":        string(ipsecConfig),
		}
		log.Printf("[DEBUG] alicloud_vpn - adding vpn connection: %v", mapping)
		ids = append(ids, vpn.VpnConnectionId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("vpn_connections", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
