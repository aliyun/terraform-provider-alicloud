package alicloud

import (
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudApigatewayGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudApigatewayGroupsRead,

		Schema: map[string]*schema.Schema{
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
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sub_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modified_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"traffic_limit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"billing_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"illegal_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func dataSourceAlicloudApigatewayGroupsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).cloudapiconn

	args := cloudapi.CreateDescribeApiGroupsRequest()
	args.RegionId = string(getRegion(d, meta))
	args.PageSize = requests.NewInteger(PageSizeLarge)
	args.PageNumber = requests.NewInteger(1)

	var allGroups []cloudapi.ApiGroupAttribute

	for {
		resp, err := conn.DescribeApiGroups(args)
		if err != nil {
			return err
		}

		if resp == nil {
			break
		}

		allGroups = resp.ApiGroupAttributes.ApiGroupAttribute

		if len(allGroups) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return err
		} else {
			args.PageNumber = page
		}
	}

	if len(allGroups) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	log.Printf("[DEBUG] alicloud_apigateway - Groups found: %#v", allGroups)

	return apigatewayGroupsDecriptionAttributes(d, allGroups, meta)
}

func apigatewayGroupsDecriptionAttributes(d *schema.ResourceData, groupsSetTypes []cloudapi.ApiGroupAttribute, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, group := range groupsSetTypes {
		mapping := map[string]interface{}{
			"id":             group.GroupId,
			"name":           group.GroupName,
			"region_id":      group.RegionId,
			"sub_domain":     group.SubDomain,
			"description":    group.Description,
			"created_time":   group.CreatedTime,
			"modified_time":  group.ModifiedTime,
			"traffic_limit":  group.TrafficLimit,
			"billing_status": group.BillingStatus,
			"illegal_status": group.IllegalStatus,
		}
		log.Printf("[DEBUG] alicloud_apigateway - adding group: %v", mapping)
		ids = append(ids, group.GroupId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("groups", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
