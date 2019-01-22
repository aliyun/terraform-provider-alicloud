package alicloud

import (
	"fmt"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
	client := meta.(*connectivity.AliyunClient)

	args := vpc.CreateDescribeVSwitchesRequest()
	args.RegionId = string(client.Region)
	// API DescribeVSwitches has some limitations
	// If there is no vpc_id, setting PageSizeSmall can avoid ServiceUnavailable Error
	args.PageSize = requests.NewInteger(PageSizeSmall)
	args.PageNumber = requests.NewInteger(1)
	if v, ok := d.GetOk("zone_id"); ok {
		args.ZoneId = Trim(v.(string))
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		args.VpcId = Trim(v.(string))
		args.PageSize = requests.NewInteger(PageSizeLarge)
	}

	var allVSwitches []vpc.VSwitch
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(Trim(v.(string))); err == nil {
			nameRegex = r
		}
	}
	invoker := NewInvoker()
	for {
		var raw interface{}
		if err := invoker.Run(func() error {
			rsp, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeVSwitches(args)
			})
			raw = rsp
			return err
		}); err != nil {
			return err
		}
		resp, _ := raw.(*vpc.DescribeVSwitchesResponse)
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

		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return err
		} else {
			args.PageNumber = page
		}
	}

	return VSwitchesDecriptionAttributes(d, allVSwitches, meta)
}

func VSwitchesDecriptionAttributes(d *schema.ResourceData, vsws []vpc.VSwitch, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var ids []string
	var s []map[string]interface{}
	instReq := ecs.CreateDescribeInstancesRequest()

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
		instReq.VpcId = vsw.VpcId
		instReq.VSwitchId = vsw.VSwitchId
		instReq.ZoneId = vsw.ZoneId
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstances(instReq)
		})
		if err != nil {
			return fmt.Errorf("DescribeInstances got an error: %#v.", err)
		}
		resp, _ := raw.(*ecs.DescribeInstancesResponse)
		if resp != nil && len(resp.Instances.Instance) > 0 {
			instance_ids := make([]string, len(resp.Instances.Instance))
			if len(instance_ids) > 0 {
				for _, inst := range resp.Instances.Instance {
					instance_ids = append(instance_ids, inst.InstanceId)
				}
			}
			mapping["instance_ids"] = instance_ids
		}

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
