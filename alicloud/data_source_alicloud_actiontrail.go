package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/actiontrail"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudActiontrails() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudActionTrailsRead,

		Schema: map[string]*schema.Schema{
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"actiontrails": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
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
					},
				},
			},
		},
	}
}

func dataSourceAlicloudActionTrailsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := actiontrail.CreateDescribeTrailsRequest()
	raw, err := client.WithActionTrailClient(func(actiontrailClient *actiontrail.Client) (interface{}, error) {
		return actiontrailClient.DescribeTrails(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_actiontrails", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw)

	response, _ := raw.(*actiontrail.DescribeTrailsResponse)
	var filteredTrailList []actiontrail.TrailListItem
	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, trail := range response.TrailList {
			if r != nil && !r.MatchString(trail.Name) {
				continue
			}
			filteredTrailList = append(filteredTrailList, trail)
		}
	} else {
		filteredTrailList = response.TrailList
	}

	var ids []string
	var filteredActionTrail []map[string]interface{}
	for _, trail := range filteredTrailList {
		mapping := map[string]interface{}{
			"name":               trail.Name,
			"event_rw":           trail.EventRW,
			"oss_bucket_name":    trail.OssBucketName,
			"oss_key_prefix":     trail.OssKeyPrefix,
			"role_name":          trail.RoleName,
			"sls_project_arn":    trail.SlsProjectArn,
			"sls_write_role_arn": trail.SlsWriteRoleArn,
		}
		ids = append(ids, trail.Name)
		filteredActionTrail = append(filteredActionTrail, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("actiontrails", filteredActionTrail); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", ids); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), filteredActionTrail)
	}

	return nil
}
