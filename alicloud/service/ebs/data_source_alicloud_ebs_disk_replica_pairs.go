package ebs

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDiskReplicaPairs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDiskReplicaPairsRead,
		Schema: map[string]*schema.Schema{
			"replica_group_id": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"site": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"ids": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"pairs": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"bandwidth": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"description": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"destination_disk_id": {
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
						"disk_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"pair_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"payment_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"rpo": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"replica_pair_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"resource_group_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"source_zone_id": {
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

func dataSourceDiskReplicaPairsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	if v, ok := d.GetOk("replica_group_id"); ok {
		request["ReplicaGroupId"] = v
	}
	if v, ok := d.GetOk("site"); ok {
		request["Site"] = v
	}
	request["MaxResults"] = helper.PageSizeLarge

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var err error
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "DescribeDiskReplicaPairs"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			resp, err := client.RpcPost("ebs", "2021-07-30", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			response = resp
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_ebs_disk_replica_pairs", action, errs.AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.ReplicaPairs", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.ReplicaPairs", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ReplicaPairId"])]; !ok {
					continue
				}
			}
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
			"id":                    fmt.Sprint(object["ReplicaPairId"]),
			"bandwidth":             object["Bandwidth"],
			"description":           object["Description"],
			"rpo":                   object["RPO"],
			"replica_pair_id":       object["ReplicaPairId"],
			"resource_group_id":     object["ResourceGroupId"],
			"destination_disk_id":   object["ResourceGroupId"],
			"destination_region_id": object["DestinationRegion"],
			"destination_zone_id":   object["DestinationZoneId"],
			"disk_id":               object["SourceDiskId"],
			"pair_name":             object["PairName"],
			"payment_type":          object["ChargeType"],
			"source_zone_id":        object["SourceZoneId"],
			"status":                object["Status"],
		}

		ids = append(ids, fmt.Sprint(object["ReplicaPairId"]))
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("pairs", s); err != nil {
		return errs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}
	return nil
}
