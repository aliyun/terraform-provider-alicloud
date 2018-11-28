package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudVpnConnections() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVpnConnectionsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				MinItems: 1,
			},

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
			"connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
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

						"ipsec_config": {
							Type:     schema.TypeList,
							Optional: true,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipsec_enc_alg": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"ipsec_auth_alg": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"ipsec_pfs": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"ipsec_lifetime": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},

						"ike_config": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"psk": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"ike_version": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"ike_mode": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"ike_enc_alg": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"ike_auth_alg": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"ike_pfs": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"ike_lifetime": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
									},
									"ike_local_id": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"ike_remote_id": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
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

func dataSourceAlicloudVpnConnectionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := vpc.CreateDescribeVpnConnectionsRequest()
	args.RegionId = string(client.Region)
	args.PageSize = requests.NewInteger(PageSizeLarge)
	args.PageNumber = requests.NewInteger(1)
	var allVpnConns []vpc.VpnConnection

	if v, ok := d.GetOk("vpn_gateway_id"); ok && v.(string) != "" {
		args.VpnGatewayId = v.(string)
	}

	if v, ok := d.GetOk("customer_gateway_id"); ok && v.(string) != "" {
		args.CustomerGatewayId = v.(string)
	}

	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVpnConnections(args)
		})
		if err != nil {
			return err
		}
		resp, _ := raw.(*vpc.DescribeVpnConnectionsResponse)

		if resp == nil || len(resp.VpnConnections.VpnConnection) < 1 {
			break
		}

		allVpnConns = append(allVpnConns, resp.VpnConnections.VpnConnection...)

		if len(resp.VpnConnections.VpnConnection) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return err
		} else {
			args.PageNumber = page
		}
	}

	var filteredVpnConns []vpc.VpnConnection
	var reg *regexp.Regexp
	var ids []string
	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
		for _, item := range v.([]interface{}) {
			ids = append(ids, strings.Trim(item.(string), " "))
		}
	}
	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		if r, err := regexp.Compile(nameRegex.(string)); err == nil {
			reg = r
		}
	}

	for _, vpnConn := range allVpnConns {
		if reg != nil {
			if !reg.MatchString(vpnConn.Name) {
				continue
			}
		}
		if ids != nil && len(ids) != 0 {
			for _, id := range ids {
				if vpnConn.VpnConnectionId == id {
					filteredVpnConns = append(filteredVpnConns, vpnConn)
				}
			}
		} else {
			filteredVpnConns = append(filteredVpnConns, vpnConn)
		}
	}

	if len(filteredVpnConns) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}
	log.Printf("[DEBUG] alicloud_vpn_connections VPN connections filter: %#v", filteredVpnConns)
	return vpnConnectionsDecriptionAttributes(d, filteredVpnConns, meta)
}

func vpnConnectionsDecriptionAttributes(d *schema.ResourceData, vpnSetTypes []vpc.VpnConnection, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}
	var ids []string
	var s []map[string]interface{}
	for _, conn := range vpnSetTypes {
		mapping := map[string]interface{}{
			"id":                  conn.VpnConnectionId,
			"customer_gateway_id": conn.CustomerGatewayId,
			"vpn_gateway_id":      conn.VpnGatewayId,
			"name":                conn.Name,
			"local_subnet":        conn.LocalSubnet,
			"remote_subnet":       conn.RemoteSubnet,
			"create_time":         conn.CreateTime,
			"effect_immediately":  conn.EffectImmediately,
			"status":              conn.Status,
			"ike_config":          vpnGatewayService.ParseIkeConfig(conn.IkeConfig),
			"ipsec_config":        vpnGatewayService.ParseIpsecConfig(conn.IpsecConfig),
		}
		log.Printf("[DEBUG] alicloud_vpn - adding vpn connection: %v", mapping)
		ids = append(ids, conn.VpnConnectionId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("connections", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
