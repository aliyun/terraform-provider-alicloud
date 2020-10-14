package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
				ValidateFunc: validation.ValidateRegexp,
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
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
	vpcService := VpcService{client}

	request := vpc.CreateDescribeVSwitchesRequest()
	request.RegionId = string(client.Region)
	request.ResourceGroupId = d.Get("resource_group_id").(string)
	// API DescribeVSwitches has some limitations
	// If there is no vpc_id, setting PageSizeSmall can avoid ServiceUnavailable Error
	request.PageSize = requests.NewInteger(PageSizeSmall)
	request.PageNumber = requests.NewInteger(1)
	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = Trim(v.(string))
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = Trim(v.(string))
	}

	var allVSwitches []vpc.VSwitch
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(Trim(v.(string))); err == nil {
			nameRegex = r
		}
	}

	// ids
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	invoker := NewInvoker()
	for {
		var raw interface{}
		var err error
		if err = invoker.Run(func() error {
			raw, err = client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeVSwitches(request)
			})
			return err
		}); err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_vswitches", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeVSwitchesResponse)
		if len(response.VSwitches.VSwitch) < 1 {
			break
		}

		for _, vsw := range response.VSwitches.VSwitch {
			if v, ok := d.GetOk("cidr_block"); ok && vsw.CidrBlock != Trim(v.(string)) {
				continue
			}

			if v, ok := d.GetOk("is_default"); ok && vsw.IsDefault != v.(bool) {
				continue
			}

			if len(idsMap) > 0 {
				if _, ok := idsMap[vsw.VSwitchId]; !ok {
					continue
				}
			}

			if nameRegex != nil {
				if !nameRegex.MatchString(vsw.VSwitchName) {
					continue
				}
			}
			if value, ok := d.GetOk("tags"); ok && len(value.(map[string]interface{})) > 0 {
				tags, err := vpcService.DescribeTags(vsw.VSwitchId, value.(map[string]interface{}), TagResourceVSwitch)
				if err != nil {
					return WrapError(err)
				}
				if len(tags) < 1 {
					continue
				}
			}
			allVSwitches = append(allVSwitches, vsw)
		}

		if len(response.VSwitches.VSwitch) < PageSizeSmall {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}

	return VSwitchesDecriptionAttributes(d, allVSwitches, meta)
}

func VSwitchesDecriptionAttributes(d *schema.ResourceData, vsws []vpc.VSwitch, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var ids []string
	var names []string
	var s []map[string]interface{}
	request := ecs.CreateDescribeInstancesRequest()

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
		request.VpcId = vsw.VpcId
		request.VSwitchId = vsw.VSwitchId
		request.ZoneId = vsw.ZoneId
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstances(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_vswitches", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		response, _ := raw.(*ecs.DescribeInstancesResponse)
		if len(response.Instances.Instance) > 0 {
			instanceIds := make([]string, 0, len(response.Instances.Instance))

			for _, inst := range response.Instances.Instance {
				instanceIds = append(instanceIds, inst.InstanceId)
			}
			mapping["instance_ids"] = instanceIds
		}

		ids = append(ids, vsw.VSwitchId)
		names = append(names, vsw.VSwitchName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("vswitches", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
