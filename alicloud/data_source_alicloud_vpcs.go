package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"regexp"
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
	conn := meta.(*AliyunClient).ecsconn

	args := &ecs.DescribeVpcsArgs{
		RegionId: getRegion(d, meta),
	}

	var allVpcs []ecs.VpcSetType

	for {
		vpcs, paginationResult, err := conn.DescribeVpcs(args)
		if err != nil {
			return err
		}

		allVpcs = append(allVpcs, vpcs...)

		pagination := paginationResult.NextPage()
		if pagination == nil {
			break
		}

		args.Pagination = *pagination
	}

	var filteredVpcsTemp []ecs.VpcSetType
	var route_tables []string

	for _, vpc := range allVpcs {
		if cidrBlock, ok := d.GetOk("cidr_block"); ok && vpc.CidrBlock != cidrBlock.(string) {
			continue
		}

		if status, ok := d.GetOk("status"); ok && string(vpc.Status) != status.(string) {
			continue
		}

		if isDefault, ok := d.GetOk("is_default"); ok && vpc.IsDefault != isDefault.(bool) {
			continue
		}

		if vswitchId, ok := d.GetOk("vswitch_id"); ok && !vpcVswitchIdListContains(vpc.VSwitchIds.VSwitchId, vswitchId.(string)) {
			continue
		}

		vrouters, _, err := meta.(*AliyunClient).vpcconn.DescribeVRouters(&ecs.DescribeVRoutersArgs{
			VRouterId: vpc.VRouterId,
			RegionId:  getRegion(d, meta),
		})
		if err != nil {
			return fmt.Errorf("Error DescribVRouters by vrouter_id %s: %#v", vpc.VRouterId, err)
		}
		if len(vrouters) > 0 && len(vrouters[0].RouteTableIds.RouteTableId) > 0 {
			route_tables = append(route_tables, vrouters[0].RouteTableIds.RouteTableId[0])
		} else {
			route_tables = append(route_tables, "")
		}

		filteredVpcsTemp = append(filteredVpcsTemp, vpc)
	}

	var filteredVpcs []ecs.VpcSetType

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
func vpcsDecriptionAttributes(d *schema.ResourceData, vpcSetTypes []ecs.VpcSetType, route_tables []string, meta interface{}) error {
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
			"creation_time":  vpc.CreationTime.String(),
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
