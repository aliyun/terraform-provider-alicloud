package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type SecurityGroup struct {
	Attributes   ecs.DescribeSecurityGroupAttributeResponse
	CreationTime string
	Tags         ecs.TagsInDescribeSecurityGroups
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
			"tags": tagsSchema(),

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
						"tags": tagsSchema(),
					},
				},
			},
		},
	}
}

func dataSourceAlicloudSecurityGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

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
	if v, ok := d.GetOk("tags"); ok {
		var tags []ecs.DescribeSecurityGroupsTag

		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, ecs.DescribeSecurityGroupsTag{
				Key:   key,
				Value: value.(string),
			})
		}
		args.Tag = &tags
	}
	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeSecurityGroups(args)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "security_groups", args.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		resp, _ := raw.(*ecs.DescribeSecurityGroupsResponse)
		if resp == nil || len(resp.SecurityGroups.SecurityGroup) < 1 {
			break
		}

		for _, item := range resp.SecurityGroups.SecurityGroup {
			if nameRegex != nil {
				if !nameRegex.MatchString(item.SecurityGroupName) {
					continue
				}
			}

			attr, err := ecsService.DescribeSecurityGroupAttribute(item.SecurityGroupId)
			if err != nil {
				return WrapError(err)
			}

			sg = append(sg,
				SecurityGroup{
					Attributes:   attr,
					CreationTime: item.CreationTime,
					Tags:         item.Tags,
				},
			)
		}

		if len(resp.SecurityGroups.SecurityGroup) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return WrapError(err)
		} else {
			args.PageNumber = page
		}
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
			"tags":          tagsToMap(item.Tags.Tag),
		}

		ids = append(ids, string(item.Attributes.SecurityGroupId))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("groups", s); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
