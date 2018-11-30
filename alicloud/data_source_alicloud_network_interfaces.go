package alicloud

import (
	"fmt"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudNetworkInterfaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudNetworkInterfacesRead,
		Schema: map[string]*schema.Schema{
			"ids": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				MinItems: 1,
				MaxItems: 100,
			},
			"name_regex": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"vswitch_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"Primary", "Secondary"}),
			},
			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"interfaces": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"mac": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_groups": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
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

func dataSourceAlicloudNetworkInterfacesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	args := ecs.CreateDescribeNetworkInterfacesRequest()
	if networkInterfaceIds, ok := d.GetOk("ids"); ok {
		ids := expandStringList(networkInterfaceIds.(*schema.Set).List())
		args.NetworkInterfaceId = &ids
	}

	if vpcId, ok := d.GetOk("vpc_id"); ok {
		args.VpcId = vpcId.(string)
	}

	if vswitchId, ok := d.GetOk("vswitch_id"); ok {
		args.VSwitchId = vswitchId.(string)
	}

	if privateIp, ok := d.GetOk("private_ip"); ok {
		args.PrimaryIpAddress = privateIp.(string)
	}

	if securityGroupId, ok := d.GetOk("security_group_id"); ok {
		args.SecurityGroupId = securityGroupId.(string)
	}

	if typ, ok := d.GetOk("type"); ok {
		args.Type = typ.(string)
	}

	if instanceId, ok := d.GetOk("instance_id"); ok {
		args.InstanceId = instanceId.(string)
	}

	if v, ok := d.GetOk("tags"); ok {
		var tags []ecs.DescribeNetworkInterfacesTag
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, ecs.DescribeNetworkInterfacesTag{
				Key:   key,
				Value: value.(string),
			})
		}
		args.Tag = &tags
	}

	var allEnis []ecs.NetworkInterfaceSet
	args.PageSize = requests.NewInteger(PageSizeLarge)
	args.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeNetworkInterfaces(args)
		})
		if err != nil {
			return fmt.Errorf("Descirbe network interfaces failed, %#v", err)
		}

		resp := raw.(*ecs.DescribeNetworkInterfacesResponse)
		if resp == nil || len(resp.NetworkInterfaceSets.NetworkInterfaceSet) < 1 {
			break
		}

		allEnis = append(allEnis, resp.NetworkInterfaceSets.NetworkInterfaceSet...)

		if len(resp.NetworkInterfaceSets.NetworkInterfaceSet) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return err
		} else {
			args.PageNumber = page
		}
	}

	var filterEnis []ecs.NetworkInterfaceSet
	nameRegex, ok := d.GetOk("name_regex")
	if ok && nameRegex.(string) != "" {
		var r *regexp.Regexp
		r = regexp.MustCompile(nameRegex.(string))

		for i := range allEnis {
			if r.MatchString(allEnis[i].NetworkInterfaceName) {
				filterEnis = append(filterEnis, allEnis[i])
			}
		}
	} else {
		filterEnis = allEnis
	}

	return networkInterfaceDescriptionAttributes(d, filterEnis)
}

func networkInterfaceDescriptionAttributes(d *schema.ResourceData, enis []ecs.NetworkInterfaceSet) error {
	var ids []string
	var s []map[string]interface{}
	for _, eni := range enis {
		var ips []string
		for _, ip := range eni.PrivateIpSets.PrivateIpSet {
			if ip.Primary {
				continue
			}
			ips = append(ips, ip.PrivateIpAddress)
		}
		mapping := map[string]interface{}{
			"id":              eni.NetworkInterfaceId,
			"name":            eni.NetworkInterfaceName,
			"status":          eni.Status,
			"vpc_id":          eni.VpcId,
			"vswitch_id":      eni.VSwitchId,
			"zone_id":         eni.ZoneId,
			"public_ip":       eni.AssociatedPublicIp.PublicIpAddress,
			"private_ip":      eni.PrivateIpAddress,
			"private_ips":     ips,
			"mac":             eni.MacAddress,
			"security_groups": eni.SecurityGroupIds.SecurityGroupId,
			"description":     eni.Description,
			"instance_id":     eni.InstanceId,
			"creation_time":   eni.CreationTime,
			"tags":            tagsToMap(eni.Tags.Tag),
		}

		ids = append(ids, eni.NetworkInterfaceId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("interfaces", s); err != nil {
		return err
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
