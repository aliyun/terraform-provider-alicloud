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

func dataSourceAlicloudEcdDesktopGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEcdDesktopGroupsRead,
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
			"desktop_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"excluded_end_user_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"office_site_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"own_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "Month",
				ValidateFunc: validation.StringInSlice([]string{"Month", "Year"}, false),
			},
			"status": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_auto_setup": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"allow_buffer_count": {
							Type:     schema.TypeInt,
							Computed: true,
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
						"id": {
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
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_desktops_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeString,
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
					},
				},
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAlicloudEcdDesktopGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDesktopGroups"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("desktop_group_name"); ok {
		request["DesktopGroupName"] = v
	}
	if v, ok := d.GetOk("excluded_end_user_ids"); ok {
		request["ExcludedEndUserIds"] = convertListToJsonString(v.(*schema.Set).List())
	}
	if v, ok := d.GetOk("office_site_id"); ok {
		request["OfficeSiteId"] = v
	}
	if v, ok := d.GetOk("own_type"); ok {
		request["OwnType"] = v
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var desktopGroupNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		desktopGroupNameRegex = r
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
	var response map[string]interface{}
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecd_desktop_groups", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.DesktopGroups", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.DesktopGroups", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if desktopGroupNameRegex != nil && !desktopGroupNameRegex.MatchString(fmt.Sprint(item["DesktopGroupName"])) {
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
	for _, object := range objects {
		mapping := map[string]interface{}{
			"bundle_id":            object["OwnBundleId"],
			"comments":             object["Comments"],
			"cpu":                  formatInt(object["Cpu"]),
			"create_time":          object["CreateTime"],
			"creator":              object["Creator"],
			"data_disk_category":   object["DataDiskCategory"],
			"data_disk_size":       object["DataDiskSize"],
			"id":                   fmt.Sprint(object["DesktopGroupId"]),
			"desktop_group_id":     fmt.Sprint(object["DesktopGroupId"]),
			"desktop_group_name":   object["DesktopGroupName"],
			"directory_id":         object["DirectoryId"],
			"directory_type":       object["DirectoryType"],
			"end_user_count":       formatInt(object["EndUserCount"]),
			"expired_time":         object["ExpiredTime"],
			"gpu_count":            object["GpuCount"],
			"gpu_spec":             object["GpuSpec"],
			"keep_duration":        fmt.Sprint(object["KeepDuration"]),
			"max_desktops_count":   formatInt(object["MaxDesktopsCount"]),
			"memory":               fmt.Sprint(object["Memory"]),
			"min_desktops_count":   formatInt(object["MinDesktopsCount"]),
			"office_site_id":       object["OfficeSiteId"],
			"office_site_name":     object["OfficeSiteName"],
			"office_site_type":     object["OfficeSiteType"],
			"own_bundle_name":      object["OwnBundleName"],
			"pay_type":             object["PayType"],
			"policy_group_id":      object["PolicyGroupId"],
			"policy_group_name":    object["PolicyGroupName"],
			"system_disk_category": object["SystemDiskCategory"],
			"system_disk_size":     formatInt(object["SystemDiskSize"]),
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["DesktopGroupName"])
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["DesktopGroupId"])
		ecdService := EcdService{client}
		getResp, err := ecdService.DescribeUsersInGroup(id)
		if err != nil {
			return WrapError(err)
		}
		endUserIdsItems := make([]string, 0)
		for _, endUser := range getResp {
			userObject := endUser.(map[string]interface{})
			if endUserId, ok := userObject["EndUserId"]; ok && endUserId != nil {
				endUserIdsItems = append(endUserIdsItems, fmt.Sprint(endUserId))
			}
		}
		mapping["end_user_ids"] = endUserIdsItems
		getResp1, err := ecdService.DescribeEcdDesktopGroup(id)
		if err != nil {
			return WrapError(err)
		}
		if v, ok := getResp1["AllowAutoSetup"]; ok && fmt.Sprint(v) != "0" {
			mapping["allow_auto_setup"] = formatInt(v)
		}
		if v, ok := getResp1["AllowBufferCount"]; ok && fmt.Sprint(v) != "0" {
			mapping["allow_buffer_count"] = formatInt(v)
		}
		if v, ok := getResp1["ResType"]; ok && fmt.Sprint(v) != "0" {
			mapping["res_type"] = formatInt(v)
		}

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
