package vpcipam

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
)

func DataSourceVpcIpamIpamPoolAllocations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudVpcIpamIpamPoolAllocationRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipam_pool_allocation_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipam_pool_allocation_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipam_pool_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"allocations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipam_pool_allocation_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipam_pool_allocation_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipam_pool_allocation_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipam_pool_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_owner_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"resource_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"total_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAliCloudVpcIpamIpamPoolAllocationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
		}
		nameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "ListIpamPoolAllocations"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("cidr"); ok {
		request["Cidr"] = v
	}
	if v, ok := d.GetOk("ipam_pool_allocation_name"); ok {
		request["IpamPoolAllocationName"] = v
	}
	request["IpamPoolId"] = d.Get("ipam_pool_id")
	if v, ok := d.GetOk("ipam_pool_allocation_id"); ok {
		request["IpamPoolAllocationId"] = v
	}
	request["MaxResults"] = helper.PageSizeLarge
	for {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("VpcIpam", "2023-02-28", action, query, request, true)

			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.IpamPoolAllocations[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["IpamPoolAllocationName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["IpamPoolAllocationId"])]; !ok {
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
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["IpamPoolAllocationId"]

		mapping["cidr"] = objectRaw["Cidr"]
		mapping["create_time"] = objectRaw["CreationTime"]
		mapping["ipam_pool_allocation_description"] = objectRaw["IpamPoolAllocationDescription"]
		mapping["ipam_pool_allocation_name"] = objectRaw["IpamPoolAllocationName"]
		mapping["ipam_pool_id"] = objectRaw["IpamPoolId"]
		mapping["region_id"] = objectRaw["RegionId"]
		mapping["resource_id"] = objectRaw["ResourceId"]
		mapping["resource_owner_id"] = objectRaw["ResourceOwnerId"]
		mapping["resource_region_id"] = objectRaw["ResourceRegionId"]
		mapping["resource_type"] = objectRaw["ResourceType"]
		mapping["source_cidr"] = objectRaw["SourceCidr"]
		mapping["status"] = objectRaw["Status"]
		mapping["ipam_pool_allocation_id"] = objectRaw["IpamPoolAllocationId"]

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["IpamPoolAllocationName"])
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("allocations", s); err != nil {
		return errs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}
	return nil
}
