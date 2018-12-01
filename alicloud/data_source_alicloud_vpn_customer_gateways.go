package alicloud

import (
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudVpnCustomerGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVpnCgwsRead,

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
			"gateways": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
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
	client := meta.(*connectivity.AliyunClient)
	args := vpc.CreateDescribeCustomerGatewaysRequest()
	args.RegionId = client.RegionId
	args.PageSize = requests.NewInteger(PageSizeLarge)
	args.PageNumber = requests.NewInteger(1)
	var allCgws []vpc.CustomerGateway

	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeCustomerGateways(args)
		})
		if err != nil {
			return err
		}
		resp, _ := raw.(*vpc.DescribeCustomerGatewaysResponse)
		if resp == nil || len(resp.CustomerGateways.CustomerGateway) < 1 {
			break
		}
		allCgws = append(allCgws, resp.CustomerGateways.CustomerGateway...)
		if len(resp.CustomerGateways.CustomerGateway) < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return err
		} else {
			args.PageNumber = page
		}
	}

	var filteredCgws []vpc.CustomerGateway
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

	for _, cgw := range allCgws {
		if reg != nil {
			if !reg.MatchString(cgw.Name) {
				continue
			}
		}

		if ids != nil && len(ids) != 0 {
			for _, id := range ids {
				if cgw.CustomerGatewayId == id {
					filteredCgws = append(filteredCgws, cgw)
				}
			}
		} else {
			filteredCgws = append(filteredCgws, cgw)
		}
	}

	return vpnCgwsDecriptionAttributes(d, filteredCgws, meta)
}

func vpnCgwsDecriptionAttributes(d *schema.ResourceData, vpnSetTypes []vpc.CustomerGateway, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, vpn := range vpnSetTypes {
		mapping := map[string]interface{}{
			"id":          vpn.CustomerGatewayId,
			"name":        vpn.Name,
			"ip_address":  vpn.IpAddress,
			"description": vpn.Description,
			"create_time": vpn.CreateTime,
		}
		ids = append(ids, vpn.CustomerGatewayId)
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
