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
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
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

	request := cloudapi.CreateDescribeAppAttributesRequest()
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	var apps []cloudapi.AppAttribute

	for {
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DescribeAppAttributes(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_api_gateway_apps", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*cloudapi.DescribeAppAttributesResponse)

		apps = append(apps, response.Apps.AppAttribute...)

		if len(apps) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	var filteredAppsTemp []cloudapi.AppAttribute
	nameRegex, ok := d.GetOk("name_regex")
	var r *regexp.Regexp
	if ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}

	// ids
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[vv.(string)] = vv.(string)
		}
	}

	for _, app := range apps {
		if r != nil && !r.MatchString(app.AppName) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[strconv.Itoa(app.AppId)]; !ok {
				continue
			}
		}

		filteredAppsTemp = append(filteredAppsTemp, app)
	}

	return apigatewayAppsDecriptionAttributes(d, filteredAppsTemp, meta)
}

func apigatewayAppsDecriptionAttributes(d *schema.ResourceData, apps []cloudapi.AppAttribute, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	var names []string
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
		names = append(names, app.AppName)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("apps", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
