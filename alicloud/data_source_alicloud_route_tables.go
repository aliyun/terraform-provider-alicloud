package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudRouteTables() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRouteTablesRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
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
			"tables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"router_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_table_type": {
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
func dataSourceAlicloudRouteTablesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	request := vpc.CreateDescribeRouteTableListRequest()
	request.RegionId = string(client.Region)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	request.ResourceGroupId = d.Get("resource_group_id").(string)
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	var allRouteTables []vpc.RouterTableListType
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(Trim(v.(string))); err == nil {
			nameRegex = r
		} else {
			return WrapError(err)
		}
	}
	invoker := NewInvoker()
	for {
		var raw interface{}
		var err error
		if err := invoker.Run(func() error {
			raw, err = client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeRouteTableList(request)
			})
			return err
		}); err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_route_tables", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeRouteTableListResponse)
		if len(response.RouterTableList.RouterTableListType) < 1 {
			break
		}

		for _, tables := range response.RouterTableList.RouterTableListType {
			if vpc_id, ok := d.GetOk("vpc_id"); ok && tables.VpcId != vpc_id.(string) {
				continue
			}
			if nameRegex != nil {
				if !nameRegex.MatchString(tables.RouteTableName) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[tables.RouteTableId]; !ok {
					continue
				}
			}
			if value, ok := d.GetOk("tags"); ok && len(value.(map[string]interface{})) > 0 {
				tags, err := vpcService.DescribeTags(tables.RouteTableId, value.(map[string]interface{}), TagResourceRouteTable)
				if err != nil {
					return WrapError(err)
				}
				if len(tags) < 1 {
					continue
				}

			}
			allRouteTables = append(allRouteTables, tables)
		}

		if len(response.RouterTableList.RouterTableListType) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	return RouteTablesDecriptionAttributes(d, allRouteTables, meta)
}

func RouteTablesDecriptionAttributes(d *schema.ResourceData, tables []vpc.RouterTableListType, meta interface{}) error {
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, table := range tables {
		mapping := map[string]interface{}{
			"id":               table.RouteTableId,
			"router_id":        table.RouterId,
			"route_table_type": table.RouteTableType,
			"name":             table.RouteTableName,
			"description":      table.Description,
			"creation_time":    table.CreationTime,
		}
		names = append(names, table.RouteTableName)
		ids = append(ids, table.RouteTableId)
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("tables", s); err != nil {
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
