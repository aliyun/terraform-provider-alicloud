package alicloud

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudSslVpnServers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSslVpnServersRead,

		Schema: map[string]*schema.Schema{
			"ssl_vpn_server_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"ssl_vpn_servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssl_vpn_server_id": {
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
						"client_ip_pool": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cipher": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"proto": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"compress": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"connections": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_connections": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"internet_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudSslVpnServersRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).vpcconn

	args := vpc.CreateDescribeSslVpnServersRequest()
	args.RegionId = string(getRegion(d, meta))
	args.PageSize = requests.NewInteger(PageSizeLarge)

	var allVpns []vpc.SslVpnServer

	for {
		resp, err := conn.DescribeSslVpnServers(args)
		if err != nil {
			return err
		}

		if resp == nil || len(resp.SslVpnServers.SslVpnServer) < 1 {
			break
		}

		allVpns = append(allVpns, resp.SslVpnServers.SslVpnServer...)

		if len(resp.SslVpnServers.SslVpnServer) < PageSizeLarge {
			break
		}

		args.PageNumber = args.PageNumber + requests.NewInteger(1)
	}

	var filteredVpnsTemp []vpc.SslVpnServer

	for _, v := range allVpns {
		if vpnId, ok := d.GetOk("vpn_gateway_id"); ok && string(v.VpnGatewayId) != vpnId.(string) {
			continue
		}

		if sslVpnServerId, ok := d.GetOk("ssl_vpn_server_id"); ok && string(v.SslVpnServerId) != sslVpnServerId.(string) {
			continue
		}

		filteredVpnsTemp = append(filteredVpnsTemp, v)
	}

	var filteredVpns []vpc.SslVpnServer

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
	return sslVpnServersDecriptionAttributes(d, filteredVpns, meta)
}

func sslVpnServersDecriptionAttributes(d *schema.ResourceData, vpnSetTypes []vpc.SslVpnServer, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, vpn := range vpnSetTypes {
		mapping := map[string]interface{}{
			"region_id":         vpn.RegionId,
			"ssl_vpn_server_id": vpn.SslVpnServerId,
			"vpn_gateway_id":    vpn.VpnGatewayId,
			"name":              vpn.Name,
			"local_subnet":      vpn.LocalSubnet,
			"client_ip_pool":    vpn.ClientIpPool,
			"create_time":       vpn.CreateTime,
			"cipher":            vpn.Cipher,
			"proto":             vpn.Proto,
			"port":              vpn.Port,
			"compress":          vpn.Compress,
			"connections":       vpn.Connections,
			"max_connections":   vpn.MaxConnections,
			"internet_ip":       vpn.InternetIp,
		}
		log.Printf("[DEBUG] alicloud_vpn - adding ssl vpn server: %v", mapping)
		ids = append(ids, vpn.SslVpnServerId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ssl_vpn_servers", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
