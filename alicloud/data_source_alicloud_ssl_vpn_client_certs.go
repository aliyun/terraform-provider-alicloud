package alicloud

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudSslVpnClientCerts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSslVpnClientCertsRead,

		Schema: map[string]*schema.Schema{
			"ssl_vpn_server_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"ssl_vpn_client_cert_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"ssl_vpn_client_cert_keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssl_vpn_client_cert_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssl_vpn_server_id": {
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
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudSslVpnClientCertsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).vpcconn

	args := vpc.CreateDescribeSslVpnClientCertsRequest()
	args.RegionId = string(getRegion(d, meta))
	args.PageSize = requests.NewInteger(PageSizeLarge)

	var allVpns []vpc.SslVpnClientCertKey

	for {
		resp, err := conn.DescribeSslVpnClientCerts(args)
		if err != nil {
			return err
		}

		if resp == nil || len(resp.SslVpnClientCertKeys.SslVpnClientCertKey) < 1 {
			break
		}

		allVpns = append(allVpns, resp.SslVpnClientCertKeys.SslVpnClientCertKey...)

		if len(resp.SslVpnClientCertKeys.SslVpnClientCertKey) < PageSizeLarge {
			break
		}

		args.PageNumber = args.PageNumber + requests.NewInteger(1)
	}

	//var filteredVpcsTemp []vpc.Vpc
	//var route_tables []string
	var filteredVpnsTemp []vpc.SslVpnClientCertKey

	for _, v := range allVpns {
		if sslVpnServerId, ok := d.GetOk("ssl_vpn_server_id"); ok && string(v.SslVpnServerId) != sslVpnServerId.(string) {
			continue
		}

		if sslVpnClientCertId, ok := d.GetOk("ssl_vpn_client_cert_id"); ok && string(v.SslVpnClientCertId) != sslVpnClientCertId.(string) {
			continue
		}

		filteredVpnsTemp = append(filteredVpnsTemp, v)
	}

	var filteredVpns []vpc.SslVpnClientCertKey

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
	return sslVpnClientCertsDecriptionAttributes(d, filteredVpns, meta)
}

func sslVpnClientCertsDecriptionAttributes(d *schema.ResourceData, vpnSetTypes []vpc.SslVpnClientCertKey, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, vpn := range vpnSetTypes {
		mapping := map[string]interface{}{
			"region_id":              vpn.RegionId,
			"ssl_vpn_client_cert_id": vpn.SslVpnClientCertId,
			"name":              vpn.Name,
			"ssl_vpn_server_id": vpn.SslVpnServerId,
			"create_time":       vpn.CreateTime,
			"end_time":          vpn.EndTime,
			"status":            vpn.Status,
		}
		log.Printf("[DEBUG] alicloud_vpn - adding vpn: %v", mapping)
		ids = append(ids, vpn.SslVpnClientCertId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ssl_vpn_client_cert_keys", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
