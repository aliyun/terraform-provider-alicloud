package alicloud

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudSlbServerGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSlbServerGroupsRead,

		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
				ForceNew:     true,
			},

			// Computed values
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

func dataSourceAlicloudSlbServerGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := slb.CreateDescribeVServerGroupsRequest()
	args.LoadBalancerId = d.Get("load_balancer_id").(string)

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeVServerGroups(args)
	})
	if err != nil {
		return fmt.Errorf("DescribeVServerGroups got an error: %#v", err)
	}
	resp, _ := raw.(*slb.DescribeVServerGroupsResponse)
	if resp == nil {
		return fmt.Errorf("there is no SLB with the ID %s. Please change your search criteria and try again", args.LoadBalancerId)
	}

	var filteredServerGroupsTemp []slb.VServerGroup
	nameRegex, ok := d.GetOk("name_regex")
	if (ok && nameRegex.(string) != "") || (len(idsMap) > 0) {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, serverGroup := range resp.VServerGroups.VServerGroup {
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
		filteredServerGroupsTemp = resp.VServerGroups.VServerGroup
	}

	if len(filteredServerGroupsTemp) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	log.Printf("[DEBUG] alicloud_slb_server_groups - Slb server groups found: %#v", filteredServerGroupsTemp)

	return slbServerGroupsDescriptionAttributes(d, filteredServerGroupsTemp, meta)
}

func slbServerGroupsDescriptionAttributes(d *schema.ResourceData, serverGroups []slb.VServerGroup, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var ids []string
	var s []map[string]interface{}

	for _, serverGroup := range serverGroups {
		mapping := map[string]interface{}{
			"id":   serverGroup.VServerGroupId,
			"name": serverGroup.VServerGroupName,
		}

		args := slb.CreateDescribeVServerGroupAttributeRequest()
		args.VServerGroupId = serverGroup.VServerGroupId
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DescribeVServerGroupAttribute(args)
		})
		if err == nil {
			resp, _ := raw.(*slb.DescribeVServerGroupAttributeResponse)
			if resp != nil && len(resp.BackendServers.BackendServer) > 0 {
				var backendServerMappings []map[string]interface{}
				for _, backendServer := range resp.BackendServers.BackendServer {
					backendServerMapping := map[string]interface{}{
						"instance_id": backendServer.ServerId,
						"weight":      backendServer.Weight,
					}
					backendServerMappings = append(backendServerMappings, backendServerMapping)
				}
				mapping["servers"] = backendServerMappings
			}
		}

		log.Printf("[DEBUG] alicloud_slb_server_groups - adding slb_server_group mapping: %v", mapping)
		ids = append(ids, serverGroup.VServerGroupId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("slb_server_groups", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
