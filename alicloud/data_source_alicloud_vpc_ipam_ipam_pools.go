// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudVpcIpamIpamPools() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudVpcIpamIpamPoolRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"ipam_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ipam_pool_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ipam_scope_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"pool_region_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_ipam_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"pools": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allocation_default_cidr_mask": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"allocation_max_cidr_mask": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"allocation_min_cidr_mask": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"auto_import": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"has_sub_pool": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ip_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipam_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipam_pool_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipam_pool_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipam_pool_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipam_scope_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pool_depth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"pool_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_ipam_pool_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func dataSourceAliCloudVpcIpamIpamPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
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
	action := "ListIpamPools"
	conn, err := client.NewVpcipamClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("ipam_pool_name"); ok {
		request["IpamPoolName"] = v
	}
	request["IpamScopeId"] = d.Get("ipam_scope_id")
	if v, ok := d.GetOk("pool_region_id"); ok {
		request["PoolRegionId"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("source_ipam_pool_id"); ok {
		request["SourceIpamPoolId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	request["MaxResults"] = PageSizeLarge
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-02-28"), StringPointer("AK"), query, request, &runtime)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.IpamPools[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["IpamPoolName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["IpamPoolId"])]; !ok {
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

		mapping["id"] = objectRaw["IpamPoolId"]

		mapping["allocation_default_cidr_mask"] = objectRaw["AllocationDefaultCidrMask"]
		mapping["allocation_max_cidr_mask"] = objectRaw["AllocationMaxCidrMask"]
		mapping["allocation_min_cidr_mask"] = objectRaw["AllocationMinCidrMask"]
		mapping["auto_import"] = objectRaw["AutoImport"]
		mapping["create_time"] = objectRaw["CreateTime"]
		mapping["has_sub_pool"] = objectRaw["HasSubPool"]
		mapping["ip_version"] = objectRaw["IpVersion"]
		mapping["ipam_id"] = objectRaw["IpamId"]
		mapping["ipam_pool_description"] = objectRaw["IpamPoolDescription"]
		mapping["ipam_pool_name"] = objectRaw["IpamPoolName"]
		mapping["ipam_scope_id"] = objectRaw["IpamScopeId"]
		mapping["pool_depth"] = objectRaw["PoolDepth"]
		mapping["pool_region_id"] = objectRaw["PoolRegionId"]
		mapping["resource_group_id"] = objectRaw["ResourceGroupId"]
		mapping["source_ipam_pool_id"] = objectRaw["SourceIpamPoolId"]
		mapping["status"] = objectRaw["Status"]
		mapping["ipam_pool_id"] = objectRaw["IpamPoolId"]
		mapping["region_id"] = objectRaw["IpamRegionId"]

		tagsMaps := objectRaw["Tags"]
		mapping["tags"] = tagsToMap(tagsMaps)

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["IpamPoolName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("pools", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
