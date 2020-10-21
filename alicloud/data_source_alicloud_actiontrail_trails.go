package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/actiontrail"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudActiontrailTrails() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudActiontrailTrailsRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"include_shadow_trails": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Disable", "Enable", "Fresh"}, false),
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"trails": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_rw": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"oss_bucket_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"oss_key_prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sls_project_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sls_write_role_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trail_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trail_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"actiontrails": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_rw": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"oss_bucket_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"oss_key_prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sls_project_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sls_write_role_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trail_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trail_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Deprecated: "Field 'actiontrails' has been deprecated from version 1.95.0. Use 'trails' instead.",
			},
		},
	}
}

func dataSourceAlicloudActiontrailTrailsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := actiontrail.CreateDescribeTrailsRequest()
	if v, ok := d.GetOkExists("include_shadow_trails"); ok {
		request.IncludeShadowTrails = requests.NewBoolean(v.(bool))
	}
	var objects []actiontrail.TrailListItem
	var trailNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		trailNameRegex = r
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
	status, statusOk := d.GetOk("status")
	var response *actiontrail.DescribeTrailsResponse
	raw, err := client.WithActiontrailClient(func(actiontrailClient *actiontrail.Client) (interface{}, error) {
		return actiontrailClient.DescribeTrails(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_actiontrail_trails", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ = raw.(*actiontrail.DescribeTrailsResponse)

	for _, item := range response.TrailList {
		if trailNameRegex != nil {
			if !trailNameRegex.MatchString(item.Name) {
				continue
			}
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[item.Name]; !ok {
				continue
			}
		}
		if statusOk && status != "" && status != item.Status {
			continue
		}
		objects = append(objects, item)
	}
	ids := make([]string, len(objects))
	names := make([]string, len(objects))
	s := make([]map[string]interface{}, len(objects))
	for i, object := range objects {
		mapping := map[string]interface{}{
			"event_rw":           object.EventRW,
			"oss_bucket_name":    object.OssBucketName,
			"oss_key_prefix":     object.OssKeyPrefix,
			"role_name":          object.RoleName,
			"sls_project_arn":    object.SlsProjectArn,
			"sls_write_role_arn": object.SlsWriteRoleArn,
			"status":             object.Status,
			"id":                 object.Name,
			"trail_name":         object.Name,
			"trail_region":       object.TrailRegion,
		}
		ids[i] = object.Name
		names[i] = object.Name
		s[i] = mapping
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("trails", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("actiontrails", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
