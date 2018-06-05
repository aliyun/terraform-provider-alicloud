package alicloud

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/schema"
)

type SecurityGroup struct {
	Attributes   ecs.DescribeSecurityGroupAttributeResponse
	CreationTime string
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
	client := meta.(*AliyunClient)
	conn := client.ecsconn

	args := ecs.CreateDescribeSecurityGroupsRequest()
	args.VpcId = d.Get("vpc_id").(string)
	args.PageNumber = requests.NewInteger(1)
	args.PageSize = requests.NewInteger(PageSizeLarge)

	var sg []SecurityGroup

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(v.(string)); err == nil {
			nameRegex = r
		}
	}

	for {
		resp, err := conn.DescribeSecurityGroups(args)
		if err != nil {
			return fmt.Errorf("DescribeSecurityGroups: %#v", err)
		}
		if resp == nil || len(resp.SecurityGroups.SecurityGroup) < 1 {
			break
		}

		for _, item := range resp.SecurityGroups.SecurityGroup {
			if nameRegex != nil {
				if !nameRegex.MatchString(item.SecurityGroupName) {
					continue
				}
			}

			attr, err := client.DescribeSecurityGroupAttribute(item.SecurityGroupId)
			if err != nil {
				return fmt.Errorf("DescribeSecurityGroupAttribute: %#v", err)
			}

			sg = append(sg,
				SecurityGroup{
					Attributes:   attr,
					CreationTime: item.CreationTime,
				},
			)
		}

		if len(resp.SecurityGroups.SecurityGroup) < PageSizeLarge {
			break
		}

		args.PageNumber = args.PageNumber + requests.NewInteger(1)
	}

	return securityGroupsDescription(d, sg)
}

func securityGroupsDescription(d *schema.ResourceData, sg []SecurityGroup) error {
	var ids []string
	var s []map[string]interface{}

	for _, item := range sg {
		mapping := map[string]interface{}{
			"id":            item.Attributes.SecurityGroupId,
			"name":          item.Attributes.SecurityGroupName,
			"description":   item.Attributes.Description,
			"vpc_id":        item.Attributes.VpcId,
			"inner_access":  item.Attributes.InnerAccessPolicy == string(GroupInnerAccept),
			"creation_time": item.CreationTime,
		}

		log.Printf("alicloud_security_groups - adding security group mapping: %v", mapping)
		ids = append(ids, string(item.Attributes.SecurityGroupId))
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
