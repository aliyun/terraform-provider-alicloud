package alicloud

import (
	"regexp"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudApiGatewayApps() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudApigatewayAppsRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateNameRegex,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"apps": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
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
					},
				},
			},
		},
	}
}
func dataSourceAlicloudApigatewayAppsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := cloudapi.CreateDescribeAppAttributesRequest()
	args.PageSize = requests.NewInteger(PageSizeLarge)
	args.PageNumber = requests.NewInteger(1)

	var apps []cloudapi.AppAttribute

	for {
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DescribeAppAttributes(args)
		})
		if err != nil {
			return err
		}
		resp, _ := raw.(*cloudapi.DescribeAppAttributesResponse)

		if resp == nil {
			break
		}

		apps = append(apps, resp.Apps.AppAttribute...)

		if len(apps) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return err
		} else {
			args.PageNumber = page
		}
	}

	var filteredAppsTemp []cloudapi.AppAttribute
	nameRegex, ok := d.GetOk("name_regex")
	var r *regexp.Regexp
	if ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}
	for _, app := range apps {
		if r != nil && !r.MatchString(app.AppName) {
			continue
		}

		filteredAppsTemp = append(filteredAppsTemp, app)
	}

	return apigatewayAppsDecriptionAttributes(d, filteredAppsTemp, meta)
}

func apigatewayAppsDecriptionAttributes(d *schema.ResourceData, apps []cloudapi.AppAttribute, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, app := range apps {
		mapping := map[string]interface{}{
			"id":            app.AppId,
			"name":          app.AppName,
			"description":   app.Description,
			"created_time":  app.CreatedTime,
			"modified_time": app.ModifiedTime,
		}
		ids = append(ids, strconv.Itoa(app.AppId))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("apps", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
