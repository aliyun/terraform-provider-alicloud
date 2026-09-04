// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAliCloudEcdDesktopGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudEcdDesktopGroupRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"desktop_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"desktop_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"office_site_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"period_unit": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_auto_setup": {
							Type:      schema.TypeInt,
							Computed:  true,
							Sensitive: true,
						},
						"allow_buffer_count": {
							Type:      schema.TypeInt,
							Computed:  true,
							Sensitive: true,
						},
						"bundle_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comments": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_disk_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_disk_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"desktop_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"desktop_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"directory_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"directory_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_user_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"end_user_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"expired_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gpu_count": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"gpu_spec": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"keep_duration": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_desktops_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"min_desktops_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"office_site_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"office_site_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"office_site_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"own_bundle_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pay_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"res_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"system_disk_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"system_disk_size": {
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
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAliCloudEcdDesktopGroupRead(d *schema.ResourceData, meta interface{}) error {
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
	action := "DescribeDesktopGroups"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("desktop_group_id"); ok {
		request["DesktopGroupId"] = v
	}
	if v, ok := d.GetOk("desktop_group_name"); ok {
		request["DesktopGroupName"] = v
	}
	if v, ok := d.GetOk("office_site_id"); ok {
		request["OfficeSiteId"] = v
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}
	request["MaxResults"] = PageSizeLarge
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RpcPost("ecd", "2020-09-30", action, query, request, true)

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

		resp, _ := jsonpath.Get("$.DesktopGroups[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["DesktopGroupName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["DesktopGroupId"])]; !ok {
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

		mapping["id"] = objectRaw["DesktopGroupId"]

		mapping["bundle_id"] = objectRaw["OwnBundleId"]
		mapping["comments"] = objectRaw["Comments"]
		mapping["cpu"] = objectRaw["Cpu"]
		mapping["create_time"] = objectRaw["CreateTime"]
		mapping["creator"] = objectRaw["Creator"]
		mapping["data_disk_category"] = objectRaw["DataDiskCategory"]
		mapping["data_disk_size"] = objectRaw["DataDiskSize"]
		mapping["desktop_group_name"] = objectRaw["DesktopGroupName"]
		mapping["directory_id"] = objectRaw["DirectoryId"]
		mapping["directory_type"] = objectRaw["DirectoryType"]
		mapping["expired_time"] = objectRaw["ExpiredTime"]
		mapping["gpu_count"] = objectRaw["GpuCount"]
		mapping["gpu_spec"] = objectRaw["GpuSpec"]
		mapping["keep_duration"] = objectRaw["KeepDuration"]
		mapping["max_desktops_count"] = objectRaw["MaxDesktopsCount"]
		mapping["memory"] = objectRaw["Memory"]
		mapping["min_desktops_count"] = objectRaw["MinDesktopsCount"]
		mapping["office_site_id"] = objectRaw["OfficeSiteId"]
		mapping["office_site_name"] = objectRaw["OfficeSiteName"]
		mapping["office_site_type"] = objectRaw["OfficeSiteType"]
		mapping["own_bundle_name"] = objectRaw["OwnBundleName"]
		mapping["pay_type"] = objectRaw["PayType"]
		mapping["policy_group_id"] = objectRaw["PolicyGroupId"]
		mapping["policy_group_name"] = objectRaw["PolicyGroupName"]
		mapping["system_disk_category"] = objectRaw["SystemDiskCategory"]
		mapping["system_disk_size"] = objectRaw["SystemDiskSize"]
		mapping["desktop_group_id"] = objectRaw["DesktopGroupId"]
		mapping["end_user_count"] = objectRaw["EndUserCount"]

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			names = append(names, objectRaw["DesktopGroupName"])
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(objectRaw["DesktopGroupId"])
		mapping, err = dataSourceAliCloudEcdDesktopGroupReadDescription(d, id, mapping, meta)
		if err != nil {
			return WrapError(err)
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["DesktopGroupName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
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

func dataSourceAliCloudEcdDesktopGroupReadDescription(d *schema.ResourceData, id string, object map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	ecdServiceV2 := EcdServiceV2{client}
	getResp, err := ecdServiceV2.DescribeEcdDesktopGroup(id)
	if err != nil {
		return nil, WrapError(err)
	}

	// Merge additional fields from Get API response to mapping
	// Reuse the response mapping template from Resource's read function
	mapping := object
	objectRaw := getResp

	mapping["allow_auto_setup"] = objectRaw["AllowAutoSetup"]
	mapping["allow_buffer_count"] = objectRaw["AllowBufferCount"]
	mapping["bundle_id"] = objectRaw["OwnBundleId"]
	mapping["comments"] = objectRaw["Comments"]
	mapping["cpu"] = objectRaw["Cpu"]
	mapping["create_time"] = objectRaw["CreationTime"]
	mapping["creator"] = objectRaw["Creator"]
	mapping["data_disk_category"] = objectRaw["DataDiskCategory"]
	mapping["data_disk_size"] = objectRaw["DataDiskSize"]
	mapping["desktop_group_name"] = objectRaw["DesktopGroupName"]
	mapping["directory_id"] = objectRaw["DirectoryId"]
	mapping["directory_type"] = objectRaw["DirectoryType"]
	mapping["expired_time"] = objectRaw["ExpiredTime"]
	mapping["gpu_count"] = objectRaw["GpuCount"]
	mapping["gpu_spec"] = objectRaw["GpuSpec"]
	mapping["keep_duration"] = objectRaw["KeepDuration"]
	mapping["max_desktops_count"] = objectRaw["MaxDesktopsCount"]
	mapping["memory"] = objectRaw["Memory"]
	mapping["min_desktops_count"] = objectRaw["MinDesktopsCount"]
	mapping["office_site_id"] = objectRaw["OfficeSiteId"]
	mapping["office_site_name"] = objectRaw["OfficeSiteName"]
	mapping["office_site_type"] = objectRaw["OfficeSiteType"]
	mapping["own_bundle_name"] = objectRaw["OwnBundleName"]
	mapping["pay_type"] = objectRaw["PayType"]
	mapping["policy_group_id"] = objectRaw["PolicyGroupId"]
	mapping["policy_group_name"] = objectRaw["PolicyGroupName"]
	mapping["res_type"] = objectRaw["ResType"]
	mapping["system_disk_category"] = objectRaw["SystemDiskCategory"]
	mapping["system_disk_size"] = objectRaw["SystemDiskSize"]
	mapping["desktop_group_id"] = objectRaw["DesktopGroupId"]

	usersObject, err := ecdServiceV2.DescribeDesktopGroupDescribeUsersInGroup(id)
	if err != nil {
		if !NotFoundError(err) {
			return nil, WrapError(err)
		}
	} else {
		mapping["end_user_ids"] = usersObject["EndUserIds"]
	}

	return mapping, nil
}
