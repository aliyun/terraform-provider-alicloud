package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudEbsDiskReplicaGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEbsDiskReplicaGroupsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Optional: true,
				Computed: true,
				ForceNew: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"groups": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"description": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"destination_region_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"destination_zone_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"group_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"last_recover_point": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"primary_region": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"primary_zone": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"rpo": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"replica_group_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"site": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"source_region_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"source_zone_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"standby_region": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"standby_zone": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEbsDiskReplicaGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	if v, ok := d.GetOk("ids"); ok {
		request["GroupIds"] = convertListToCommaSeparate(v.([]interface{}))
	}

	request["MaxResults"] = PageSizeLarge

	var err error
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "DescribeDiskReplicaGroups"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := client.RpcPost("ebs", "2021-07-30", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			response = resp
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ebs_disk_replica_groups", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.ReplicaGroups", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ReplicaGroups", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                    fmt.Sprint(object["ReplicaGroupId"]),
			"description":           object["Description"],
			"destination_region_id": object["DestinationRegionId"],
			"destination_zone_id":   object["DestinationZoneId"],
			"group_name":            object["GroupName"],
			"primary_region":        object["PrimaryRegion"],
			"primary_zone":          object["PrimaryZone"],
			"rpo":                   object["RPO"],
			"replica_group_id":      object["ReplicaGroupId"],
			"site":                  object["Site"],
			"source_region_id":      object["SourceRegionId"],
			"source_zone_id":        object["SourceZoneId"],
			"standby_region":        object["StandbyRegion"],
			"standby_zone":          object["StandbyZone"],
			"status":                object["Status"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("groups", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
