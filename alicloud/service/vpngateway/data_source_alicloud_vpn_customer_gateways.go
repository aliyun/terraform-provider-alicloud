package vpngateway

import (
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
)

func dataSourceAlicloudVpnCustomerGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVpnCgwsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				MinItems: 1,
			},

			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
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
							Type:     schema.TypeString,
							Computed: true,
						},
						"asn": {
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
	request := vpc.CreateDescribeCustomerGatewaysRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(helper.PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var allCgws []vpc.CustomerGateway

	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeCustomerGateways(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "aliclioud_vpn_customer_gateways", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeCustomerGatewaysResponse)
		if len(response.CustomerGateways.CustomerGateway) < 1 {
			break
		}
		allCgws = append(allCgws, response.CustomerGateways.CustomerGateway...)
		if len(response.CustomerGateways.CustomerGateway) < helper.PageSizeLarge {
			break
		}
		page, err := helper.GetNextpageNumber(request.PageNumber)
		if err != nil {
			return errs.WrapError(err)
		}
		request.PageNumber = page
	}

	var filteredCgws []vpc.CustomerGateway
	var reg *regexp.Regexp
	var ids []string
	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
		for _, item := range v.([]interface{}) {
			if item == nil {
				continue
			}
			ids = append(ids, strings.Trim(item.(string), " "))
		}
	}

	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		if r, err := regexp.Compile(nameRegex.(string)); err == nil {
			reg = r
		} else {
			return errs.WrapError(err)
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
	var names []string
	var s []map[string]interface{}
	for _, vpn := range vpnSetTypes {
		mapping := map[string]interface{}{
			"id":          vpn.CustomerGatewayId,
			"name":        vpn.Name,
			"ip_address":  vpn.IpAddress,
			"description": vpn.Description,
			"create_time": TimestampToStr(vpn.CreateTime),
			"asn":         vpn.Asn,
		}
		ids = append(ids, vpn.CustomerGatewayId)
		names = append(names, vpn.Name)
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("gateways", s); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}
	return nil
}
