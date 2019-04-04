package alicloud

import (
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudSslVpnClientCerts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSslVpnClientCertsRead,

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
			"ssl_vpn_client_certs": {
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
						"ssl_vpn_client_cert_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
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
	client := meta.(*connectivity.AliyunClient)
	args := vpc.CreateDescribeSslVpnClientCertsRequest()
	args.RegionId = client.RegionId
	args.PageSize = requests.NewInteger(PageSizeLarge)
	args.PageNumber = requests.NewInteger(1)
	var allSslVpnClientCerts []vpc.SslVpnClientCertKey

	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeSslVpnClientCerts(args)
		})
		if err != nil {
			return err
		}
		resp, _ := raw.(*vpc.DescribeSslVpnClientCertsResponse)
		if resp == nil || len(resp.SslVpnClientCertKeys.SslVpnClientCertKey) < 1 {
			break
		}
		allSslVpnClientCerts = append(allSslVpnClientCerts, resp.SslVpnClientCertKeys.SslVpnClientCertKey...)
		if len(resp.SslVpnClientCertKeys.SslVpnClientCertKey) < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return err
		} else {
			args.PageNumber = page
		}
	}

	var filteredSslVpnClientCerts []vpc.SslVpnClientCertKey
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

	for _, sslVpnClientCertKey := range allSslVpnClientCerts {
		if reg != nil {
			if !reg.MatchString(sslVpnClientCertKey.Name) {
				continue
			}
		}

		if ids != nil && len(ids) != 0 {
			for _, id := range ids {
				if sslVpnClientCertKey.SslVpnClientCertId == id {
					filteredSslVpnClientCerts = append(filteredSslVpnClientCerts, sslVpnClientCertKey)
				}
			}
		} else {
			filteredSslVpnClientCerts = append(filteredSslVpnClientCerts, sslVpnClientCertKey)
		}
	}

	return sslVpnClientCertsDecriptionAttributes(d, filteredSslVpnClientCerts, meta)
}

func sslVpnClientCertsDecriptionAttributes(d *schema.ResourceData, vpnSetTypes []vpc.SslVpnClientCertKey, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, vpn := range vpnSetTypes {
		mapping := map[string]interface{}{
			"ssl_vpn_server_id":      vpn.SslVpnServerId,
			"ssl_vpn_client_cert_id": vpn.SslVpnClientCertId,
			"region_id":              vpn.RegionId,
			"name":                   vpn.Name,
			"end_time":               vpn.EndTime,
			"create_time":            vpn.CreateTime,
			"status":                 vpn.Status,
		}
		ids = append(ids, vpn.SslVpnClientCertId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ssl_vpn_client_certs", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
