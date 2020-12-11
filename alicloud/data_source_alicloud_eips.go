package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudEips() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEipsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"in_use": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Field 'in_use' has been deprecated from provider version 1.8.0.",
			},
			"ip_addresses": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				MinItems: 1,
			},
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Computed values
			"eips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}
func dataSourceAlicloudEipsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	request := vpc.CreateDescribeEipAddressesRequest()
	request.RegionId = string(client.Region)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	request.ResourceGroupId = d.Get("resource_group_id").(string)

	idsMap := make(map[string]string)
	ipsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}
	if v, ok := d.GetOk("ip_addresses"); ok && len(v.([]interface{})) > 0 {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			ipsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	var allEips []vpc.EipAddress

	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeEipAddresses(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_eips", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeEipAddressesResponse)
		if len(response.EipAddresses.EipAddress) < 1 {
			break
		}

		for _, e := range response.EipAddresses.EipAddress {
			if len(idsMap) > 0 {
				if _, ok := idsMap[e.AllocationId]; !ok {
					continue
				}
			}
			if len(ipsMap) > 0 {
				if _, ok := ipsMap[e.IpAddress]; !ok {
					continue
				}
			}
			if value, ok := d.GetOk("tags"); ok && len(value.(map[string]interface{})) > 0 {
				tags, err := vpcService.DescribeTags(e.AllocationId, value.(map[string]interface{}), TagResourceEip)
				if err != nil {
					return WrapError(err)
				}
				if len(tags) < 1 {
					continue
				}
			}
			allEips = append(allEips, e)
		}

		if len(response.EipAddresses.EipAddress) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}

	return eipsDecriptionAttributes(d, allEips, meta)
}

func eipsDecriptionAttributes(d *schema.ResourceData, eipSetTypes []vpc.EipAddress, meta interface{}) error {
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, eip := range eipSetTypes {
		mapping := map[string]interface{}{
			"id":                   eip.AllocationId,
			"status":               eip.Status,
			"ip_address":           eip.IpAddress,
			"bandwidth":            eip.Bandwidth,
			"instance_id":          eip.InstanceId,
			"instance_type":        eip.InstanceType,
			"internet_charge_type": eip.InternetChargeType,
			"creation_time":        eip.AllocationTime,
		}
		ids = append(ids, eip.AllocationId)
		names = append(names, eip.Name)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("eips", s); err != nil {
		return WrapError(err)
	}
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
