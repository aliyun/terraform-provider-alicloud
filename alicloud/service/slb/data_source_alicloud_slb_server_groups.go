package slb

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceServerGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerGroupsRead,

		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"slb_server_groups": {
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
						"servers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"weight": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceServerGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := slb.CreateDescribeVServerGroupsRequest()
	request.RegionId = client.RegionId
	request.LoadBalancerId = d.Get("load_balancer_id").(string)

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[helper.Trim(vv.(string))] = helper.Trim(vv.(string))
		}
	}

	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeVServerGroups(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_slb_server_groups", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.DescribeVServerGroupsResponse)
	var filteredServerGroupsTemp []slb.VServerGroup
	nameRegex, ok := d.GetOk("name_regex")
	if (ok && nameRegex.(string) != "") || (len(idsMap) > 0) {
		var r *regexp.Regexp
		if nameRegex != "" {
			r, err = regexp.Compile(nameRegex.(string))
			if err != nil {
				return errs.WrapError(err)
			}
		}
		for _, serverGroup := range response.VServerGroups.VServerGroup {
			if r != nil && !r.MatchString(serverGroup.VServerGroupName) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[serverGroup.VServerGroupId]; !ok {
					continue
				}
			}

			filteredServerGroupsTemp = append(filteredServerGroupsTemp, serverGroup)
		}
	} else {
		filteredServerGroupsTemp = response.VServerGroups.VServerGroup
	}

	return slbServerGroupsDescriptionAttributes(d, filteredServerGroupsTemp, meta)
}

func slbServerGroupsDescriptionAttributes(d *schema.ResourceData, serverGroups []slb.VServerGroup, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var ids []string
	var names []string
	var s []map[string]interface{}

	for _, serverGroup := range serverGroups {
		mapping := map[string]interface{}{
			"id":   serverGroup.VServerGroupId,
			"name": serverGroup.VServerGroupName,
		}

		request := slb.CreateDescribeVServerGroupAttributeRequest()
		request.VServerGroupId = serverGroup.VServerGroupId
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DescribeVServerGroupAttribute(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_slb_server_groups", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*slb.DescribeVServerGroupAttributeResponse)
		if response != nil && len(response.BackendServers.BackendServer) > 0 {
			var backendServerMappings []map[string]interface{}
			for _, backendServer := range response.BackendServers.BackendServer {
				backendServerMapping := map[string]interface{}{
					"instance_id": backendServer.ServerId,
					"weight":      backendServer.Weight,
				}
				backendServerMappings = append(backendServerMappings, backendServerMapping)
			}
			mapping["servers"] = backendServerMappings
		}

		ids = append(ids, serverGroup.VServerGroupId)
		names = append(names, serverGroup.VServerGroupName)
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("slb_server_groups", s); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}
	return nil
}
