package alicloud

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/schema"
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
	conn := meta.(*AliyunClient).vpcconn

	args := vpc.CreateDescribeVpcsRequest()
	args.RegionId = string(getRegion(d, meta))
	args.PageSize = requests.NewInteger(PageSizeLarge)

	var allVpcs []vpc.Vpc

	for {
		resp, err := conn.DescribeVpcs(args)
		if err != nil {
			return err
		}

		if resp == nil || len(resp.Vpcs.Vpc) < 1 {
			break
		}

		allVpcs = append(allVpcs, resp.Vpcs.Vpc...)

		if len(resp.Vpcs.Vpc) < PageSizeLarge {
			break
		}

		args.PageNumber = args.PageNumber + requests.NewInteger(1)
	}

	var filteredVpcsTemp []vpc.Vpc
	var route_tables []string

	for _, v := range allVpcs {
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
		request.RegionId = string(getRegion(d, meta))

		vrs, err := conn.DescribeVRouters(request)
		if err != nil {
			return fmt.Errorf("Error DescribVRouters by vrouter_id %s: %#v", v.VRouterId, err)
		}
		if vrs != nil && len(vrs.VRouters.VRouter) > 0 {
			route_tables = append(route_tables, vrs.VRouters.VRouter[0].RouteTableIds.RouteTableId[0])
		} else {
			route_tables = append(route_tables, "")
		}

		filteredVpcsTemp = append(filteredVpcsTemp, v)
	}

	var filteredVpcs []vpc.Vpc

	if nameRegex, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(nameRegex.(string)); err == nil {
			for _, vpc := range filteredVpcsTemp {
				if r.MatchString(vpc.VpcName) {
					filteredVpcs = append(filteredVpcs, vpc)
				}
			}
		}
	} else {
		filteredVpcs = filteredVpcsTemp[:]
	}

	if len(filteredVpcs) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	log.Printf("[DEBUG] alicloud_vpc - VPCs found: %#v", allVpcs)

	return vpcsDecriptionAttributes(d, filteredVpcsTemp, route_tables, meta)
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
		log.Printf("[DEBUG] alicloud_vpc - adding vpc: %v", mapping)
		ids = append(ids, vpc.VpcId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("vpcs", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
