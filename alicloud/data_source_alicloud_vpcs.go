package alicloud

import (
	"regexp"

	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
				ValidateFunc: validation.ValidateRegexp,
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
						"tags": {
							Type:     schema.TypeMap,
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
	vpcService := VpcService{client}

	request := vpc.CreateDescribeVpcsRequest()
	request.RegionId = string(client.Region)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	request.ResourceGroupId = d.Get("resource_group_id").(string)
	var allVpcs []vpc.Vpc
	invoker := NewInvoker()
	for {
		var raw interface{}
		var err error
		err = invoker.Run(func() error {
			raw, err = client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeVpcs(request)
			})
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return err
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_vpcs", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		response, _ := raw.(*vpc.DescribeVpcsResponse)
		if len(response.Vpcs.Vpc) < 1 {
			break
		}

		allVpcs = append(allVpcs, response.Vpcs.Vpc...)

		if len(response.Vpcs.Vpc) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}

	var filteredVpcs []vpc.Vpc
	var route_tables []string
	var r *regexp.Regexp
	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(nameRegex.(string))
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

	for _, v := range allVpcs {
		if r != nil && !r.MatchString(v.VpcName) {
			continue
		}

		if len(idsMap) > 0 {
			if _, ok := idsMap[v.VpcId]; !ok {
				continue
			}
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

		if value, ok := d.GetOk("tags"); ok && len(value.(map[string]interface{})) > 0 {
			tags, err := vpcService.DescribeTags(v.VpcId, value.(map[string]interface{}), TagResourceVpc)
			if err != nil {
				return WrapError(err)
			}
			if len(tags) < 1 {
				continue
			}

		}

		request := vpc.CreateDescribeVRoutersRequest()
		request.RegionId = client.RegionId
		request.VRouterId = v.VRouterId
		request.RegionId = string(client.Region)

		var response *vpc.DescribeVRoutersResponse
		wait := incrementalWait(1*time.Second, 1*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeVRouters(request)
			})
			if err != nil {
				if IsThrottling(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			response, _ = raw.(*vpc.DescribeVRoutersResponse)
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_vpcs", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		if len(response.VRouters.VRouter) > 0 {
			route_tables = append(route_tables, response.VRouters.VRouter[0].RouteTableIds.RouteTableId[0])
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
	var names []string
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
			"tags":           vpcTagsToMap(vpc.Tags.Tag),
		}
		ids = append(ids, vpc.VpcId)
		names = append(names, vpc.VpcName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("vpcs", s); err != nil {
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
