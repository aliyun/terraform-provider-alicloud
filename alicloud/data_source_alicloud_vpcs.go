package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudVpcs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVpcsRead,

		Schema: map[string]*schema.Schema{
			"cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
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
			"is_default": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"vpcs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
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
						"vpc_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vrouter_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_default": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func dataSourceAlicloudVpcsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := vpc.CreateDescribeVpcsRequest()
	args.RegionId = string(client.Region)
	args.PageSize = requests.NewInteger(PageSizeLarge)
	args.PageNumber = requests.NewInteger(1)

	var allVpcs []vpc.Vpc
	invoker := NewInvoker()
	for {
		var raw interface{}
		if err := invoker.Run(func() error {
			rsp, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeVpcs(args)
			})
			raw = rsp
			return err
		}); err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "vpcs", args.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		resp, _ := raw.(*vpc.DescribeVpcsResponse)
		if resp == nil || len(resp.Vpcs.Vpc) < 1 {
			break
		}

		allVpcs = append(allVpcs, resp.Vpcs.Vpc...)

		if len(resp.Vpcs.Vpc) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return WrapError(err)
		} else {
			args.PageNumber = page
		}
	}

	var filteredVpcs []vpc.Vpc
	var route_tables []string
	var r *regexp.Regexp
	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}

	for _, v := range allVpcs {
		if r != nil && !r.MatchString(v.VpcName) {
			continue
		}
		if cidrBlock, ok := d.GetOk("cidr_block"); ok && v.CidrBlock != cidrBlock.(string) {
			continue
		}

		if status, ok := d.GetOk("status"); ok && string(v.Status) != status.(string) {
			continue
		}

		if isDefault, ok := d.GetOk("is_default"); ok && v.IsDefault != isDefault.(bool) {
			continue
		}

		if vswitchId, ok := d.GetOk("vswitch_id"); ok && !vpcVswitchIdListContains(v.VSwitchIds.VSwitchId, vswitchId.(string)) {
			continue
		}

		request := vpc.CreateDescribeVRoutersRequest()
		request.VRouterId = v.VRouterId
		request.RegionId = string(client.Region)

		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVRouters(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "vpcs", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		vrs, _ := raw.(*vpc.DescribeVRoutersResponse)
		if vrs != nil && len(vrs.VRouters.VRouter) > 0 {
			route_tables = append(route_tables, vrs.VRouters.VRouter[0].RouteTableIds.RouteTableId[0])
		} else {
			route_tables = append(route_tables, "")
		}

		filteredVpcs = append(filteredVpcs, v)
	}

	return vpcsDecriptionAttributes(d, filteredVpcs, route_tables, meta)
}
func vpcVswitchIdListContains(vswitchIdList []string, vswitchId string) bool {
	for _, idListItem := range vswitchIdList {
		if idListItem == vswitchId {
			return true
		}
	}
	return false
}
func vpcsDecriptionAttributes(d *schema.ResourceData, vpcSetTypes []vpc.Vpc, route_tables []string, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for index, vpc := range vpcSetTypes {
		mapping := map[string]interface{}{
			"id":             vpc.VpcId,
			"region_id":      vpc.RegionId,
			"status":         vpc.Status,
			"vpc_name":       vpc.VpcName,
			"vswitch_ids":    vpc.VSwitchIds.VSwitchId,
			"cidr_block":     vpc.CidrBlock,
			"vrouter_id":     vpc.VRouterId,
			"route_table_id": route_tables[index],
			"description":    vpc.Description,
			"is_default":     vpc.IsDefault,
			"creation_time":  vpc.CreationTime,
		}
		ids = append(ids, vpc.VpcId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("vpcs", s); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
