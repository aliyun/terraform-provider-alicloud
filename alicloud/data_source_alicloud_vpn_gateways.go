package alicloud

import (
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudVpnGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVpnsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				MinItems: 1,
			},

			"vpc_id": {
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
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{string(Init), string(Provisioning), string(Active), string(Updating), string(Deleting)}),
			},

			"business_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{string(Normal), string(FinancialLocked)}),
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
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
						"vpc_id": {
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
						"specification": {
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
						"instance_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_ipsec": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_ssl": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssl_connections": {
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
	client := meta.(*connectivity.AliyunClient)

	args := vpc.CreateDescribeVpnGatewaysRequest()
	args.RegionId = client.RegionId
	args.PageSize = requests.NewInteger(PageSizeLarge)
	args.PageNumber = requests.NewInteger(1)

	var allVpns []vpc.VpnGateway

	if v, ok := d.GetOk("vpc_id"); ok && v.(string) != "" {
		args.VpcId = v.(string)
	}

	if v, ok := d.GetOk("status"); ok && v.(string) != "" {
		args.Status = strings.ToLower(v.(string))
	}

	if v, ok := d.GetOk("business_status"); ok && v.(string) != "" {
		args.BusinessStatus = v.(string)
	}

	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVpnGateways(args)
		})
		if err != nil {
			return err
		}
		resp, _ := raw.(*vpc.DescribeVpnGatewaysResponse)

		if resp == nil || len(resp.VpnGateways.VpnGateway) < 1 {
			break
		}

		allVpns = append(allVpns, resp.VpnGateways.VpnGateway...)

		if len(resp.VpnGateways.VpnGateway) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return err
		} else {
			args.PageNumber = page
		}
	}

	var filteredVpns []vpc.VpnGateway
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

	for _, vpn := range allVpns {
		if reg != nil {
			if !reg.MatchString(vpn.Name) {
				continue
			}
		}
		if ids != nil && len(ids) != 0 {
			for _, id := range ids {
				if vpn.VpnGatewayId == id {
					filteredVpns = append(filteredVpns, vpn)
				}
			}
		} else {
			filteredVpns = append(filteredVpns, vpn)
		}

	}

	return vpnsDecriptionAttributes(d, filteredVpns, meta)
}

func convertStatus(lower string) string {
	upStr := strings.ToUpper(lower)

	wholeStr := string(upStr[0]) + lower[1:]
	return wholeStr
}

func vpnsDecriptionAttributes(d *schema.ResourceData, vpnSetTypes []vpc.VpnGateway, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, vpn := range vpnSetTypes {
		mapping := map[string]interface{}{
			"id":                   vpn.VpnGatewayId,
			"vpc_id":               vpn.VpcId,
			"internet_ip":          vpn.InternetIp,
			"create_time":          vpn.CreateTime,
			"end_time":             vpn.EndTime,
			"specification":        vpn.Spec,
			"name":                 vpn.Name,
			"description":          vpn.Description,
			"status":               convertStatus(vpn.Status),
			"business_status":      vpn.BusinessStatus,
			"instance_charge_type": vpn.ChargeType,
			"enable_ipsec":         vpn.IpsecVpn,
			"enable_ssl":           vpn.SslVpn,
			"ssl_connections":      vpn.SslMaxConnections,
		}

		ids = append(ids, vpn.VpnGatewayId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("gateways", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
