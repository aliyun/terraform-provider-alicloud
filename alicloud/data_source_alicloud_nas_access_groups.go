package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudNasAccessGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudNasAccessGroupsRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"type": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Field 'type' has been deprecated from provider version 1.95.0. New field 'access_group_type' replaces it.",
			},
			"access_group_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"access_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"file_system_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"standard", "extreme"}, false),
				Default:      "standard",
			},
			"useutc_date_time": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
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
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_group_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mount_target_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"rule_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudNasAccessGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := nas.CreateDescribeAccessGroupsRequest()
	if v, ok := d.GetOk("access_group_name"); ok {
		request.AccessGroupName = v.(string)
	}
	if v, ok := d.GetOk("file_system_type"); ok {
		request.FileSystemType = v.(string)
	}
	if v, ok := d.GetOkExists("useutc_date_time"); ok {
		request.UseUTCDateTime = requests.NewBoolean(v.(bool))
	}
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var objects []nas.AccessGroup
	var accessGroupNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		accessGroupNameRegex = r
	}
	var response *nas.DescribeAccessGroupsResponse
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
				return nasClient.DescribeAccessGroups(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"ServiceUnavailable", "Throttling"}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw)
			response, _ = raw.(*nas.DescribeAccessGroupsResponse)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_nas_access_groups", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		for _, item := range response.AccessGroups.AccessGroup {
			if accessGroupNameRegex != nil {
				if !accessGroupNameRegex.MatchString(item.AccessGroupName) {
					continue
				}
			}
			if v, ok := d.GetOk("type"); ok && v.(string) != "" && item.AccessGroupType != v.(string) {
				continue
			}
			if v, ok := d.GetOk("access_group_type"); ok && v.(string) != "" && item.AccessGroupType != v.(string) {
				continue
			}
			if v, ok := d.GetOk("description"); ok && v.(string) != "" && item.Description != v.(string) {
				continue
			}
			objects = append(objects, item)
		}
		if len(response.AccessGroups.AccessGroup) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}
	ids := make([]string, 0)
	names := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                 fmt.Sprintf("%v:%v", object.AccessGroupName, request.FileSystemType),
			"access_group_name":  object.AccessGroupName,
			"access_group_type":  object.AccessGroupType,
			"type":               object.AccessGroupType,
			"description":        object.Description,
			"mount_target_count": object.MountTargetCount,
			"rule_count":         object.RuleCount,
		}
		ids = append(ids, fmt.Sprintf("%v:%v", object.AccessGroupName, request.FileSystemType))
		names = append(names, object.AccessGroupName)
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
