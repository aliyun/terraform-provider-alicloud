package alicloud

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudVSwitches() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVSwitchesRead,

		Schema: map[string]*schema.Schema{
			"cidr_block": {
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
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"vswitches": {
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
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"cidr_block": {
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
func dataSourceAlicloudVSwitchesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).vpcconn

	args := vpc.CreateDescribeVSwitchesRequest()
	args.RegionId = string(getRegion(d, meta))
	args.PageSize = requests.NewInteger(PageSizeLarge)
	if v, ok := d.GetOk("zone_id"); ok {
		args.ZoneId = Trim(v.(string))
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		args.VpcId = Trim(v.(string))
	}

	var allVSwitches []vpc.VSwitch
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(Trim(v.(string))); err == nil {
			nameRegex = r
		}
	}
	for {
		resp, err := conn.DescribeVSwitches(args)
		if err != nil {
			return err
		}

		if resp == nil || len(resp.VSwitches.VSwitch) < 1 {
			break
		}

		for _, vsw := range resp.VSwitches.VSwitch {
			if v, ok := d.GetOk("cidr_block"); ok && vsw.CidrBlock != Trim(v.(string)) {
				continue
			}

			if v, ok := d.GetOk("is_default"); ok && vsw.IsDefault != v.(bool) {
				continue
			}

			if nameRegex != nil {
				if !nameRegex.MatchString(vsw.VSwitchName) {
					continue
				}
			}
			allVSwitches = append(allVSwitches, vsw)
		}

		if len(resp.VSwitches.VSwitch) < PageSizeLarge {
			break
		}

		args.PageNumber = args.PageNumber + requests.NewInteger(1)
	}

	if len(allVSwitches) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	log.Printf("[DEBUG] alicloud_vswitches - VSwitches found: %#v", allVSwitches)

	return VSwitchesDecriptionAttributes(d, allVSwitches, meta)
}

func VSwitchesDecriptionAttributes(d *schema.ResourceData, vsws []vpc.VSwitch, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, vsw := range vsws {
		mapping := map[string]interface{}{
			"id":            vsw.VSwitchId,
			"vpc_id":        vsw.VpcId,
			"zone_id":       vsw.ZoneId,
			"name":          vsw.VSwitchName,
			"cidr_block":    vsw.CidrBlock,
			"description":   vsw.Description,
			"is_default":    vsw.IsDefault,
			"creation_time": vsw.CreationTime,
		}
		instances, _, err := meta.(*AliyunClient).ecsconn.DescribeInstances(&ecs.DescribeInstancesArgs{
			RegionId:  getRegion(d, meta),
			VpcId:     vsw.VpcId,
			VSwitchId: vsw.VSwitchId,
			ZoneId:    vsw.ZoneId,
		})
		if err != nil {
			return fmt.Errorf("DescribeInstances got an error: %#v.", err)
		}
		instance_ids := make([]string, len(instances))
		if len(instance_ids) > 0 {
			for _, inst := range instances {
				instance_ids = append(instance_ids, inst.InstanceId)
			}
		}
		mapping["instance_ids"] = instance_ids

		log.Printf("[DEBUG] alicloud_vswitches - adding vswitch: %v", mapping)
		ids = append(ids, vsw.VSwitchId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("vswitches", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
