package alicloud

import (
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudSslVpnServers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSslVpnServersRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				MinItems: 1,
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
	client := meta.(*connectivity.AliyunClient)
	args := vpc.CreateDescribeSslVpnServersRequest()
	args.RegionId = client.RegionId
	args.PageSize = requests.NewInteger(PageSizeLarge)
	args.PageNumber = requests.NewInteger(1)
	var allSslVpnServers []vpc.SslVpnServer

	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeSslVpnServers(args)
		})
		if err != nil {
			return err
		}
		resp, _ := raw.(*vpc.DescribeSslVpnServersResponse)
		if resp == nil || len(resp.SslVpnServers.SslVpnServer) < 1 {
			break
		}
		allSslVpnServers = append(allSslVpnServers, resp.SslVpnServers.SslVpnServer...)
		if len(resp.SslVpnServers.SslVpnServer) < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return err
		} else {
			args.PageNumber = page
		}
	}

	var filteredSslVpnServers []vpc.SslVpnServer
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

	for _, sslVpnServer := range allSslVpnServers {
		if reg != nil {
			if !reg.MatchString(sslVpnServer.Name) {
				continue
			}
		}

		if ids != nil && len(ids) != 0 {
			for _, id := range ids {
				if sslVpnServer.SslVpnServerId == id {
					filteredSslVpnServers = append(filteredSslVpnServers, sslVpnServer)
				}
			}
		} else {
			filteredSslVpnServers = append(filteredSslVpnServers, sslVpnServer)
		}
	}

	return sslVpnServersDecriptionAttributes(d, filteredSslVpnServers, meta)
}

func sslVpnServersDecriptionAttributes(d *schema.ResourceData, vpnSetTypes []vpc.SslVpnServer, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, vpn := range vpnSetTypes {
		mapping := map[string]interface{}{
			"ssl_vpn_server_id": vpn.SslVpnServerId,
			"vpn_gateway_id":    vpn.VpnGatewayId,
			"region_id":         vpn.RegionId,
			"local_subnet":      vpn.LocalSubnet,
			"name":              vpn.Name,
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
