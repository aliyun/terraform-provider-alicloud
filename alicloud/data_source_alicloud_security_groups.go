package alicloud

import (
	"fmt"
	"log"
	"regexp"

	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/util"
	"github.com/hashicorp/terraform/helper/schema"
)

type SecurityGroup struct {
	Attribute    ecs.DescribeSecurityGroupAttributeResponse
	CreationTime util.ISO6801Time
}

func dataSourceAlicloudSecurityGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSecurityGroupsRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"groups": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"inner_access": {
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

func dataSourceAlicloudSecurityGroupsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	regionId := getRegion(d, meta)

	args := &ecs.DescribeSecurityGroupsArgs{
		RegionId: regionId,
		VpcId:    d.Get("vpc_id").(string),
	}

	var sg []SecurityGroup

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(v.(string)); err == nil {
			nameRegex = r
		}
	}

	for {
		items, paginationResult, err := conn.DescribeSecurityGroups(args)
		if err != nil {
			return fmt.Errorf("DescribeSecurityGroups: %#v", err)
		}

		for _, item := range items {
			if nameRegex != nil {
				if !nameRegex.MatchString(item.SecurityGroupName) {
					continue
				}
			}

			attr, err := conn.DescribeSecurityGroupAttribute(
				&ecs.DescribeSecurityGroupAttributeArgs{
					SecurityGroupId: item.SecurityGroupId,
					RegionId:        regionId,
				},
			)
			if err != nil {
				return fmt.Errorf("DescribeSecurityGroupAttribute: %#v", err)
			}

			sg = append(sg,
				SecurityGroup{
					Attribute:    *attr,
					CreationTime: item.CreationTime,
				},
			)
		}

		pagination := paginationResult.NextPage()
		if pagination == nil {
			break
		}

		args.Pagination = *pagination
	}

	return securityGroupsDescription(d, sg)
}

func securityGroupsDescription(d *schema.ResourceData, sg []SecurityGroup) error {
	var ids []string
	var s []map[string]interface{}

	for _, item := range sg {
		mapping := map[string]interface{}{
			"id":            item.Attribute.SecurityGroupId,
			"name":          item.Attribute.SecurityGroupName,
			"description":   item.Attribute.Description,
			"vpc_id":        item.Attribute.VpcId,
			"inner_access":  item.Attribute.InnerAccessPolicy == ecs.GroupInnerAccept,
			"creation_time": item.CreationTime.String(),
		}

		log.Printf("alicloud_security_groups - adding security group mapping: %v", mapping)
		ids = append(ids, string(item.Attribute.SecurityGroupId))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("groups", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
