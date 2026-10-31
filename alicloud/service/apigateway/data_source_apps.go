package apigateway

import (
	"regexp"
	"strconv"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceApps() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudApigatewayAppsRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"tags": tags.TagsSchema(),
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
						"app_code": {
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
	cloudApiService := CloudApiService{client}

	request := cloudapi.CreateDescribeAppAttributesRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(helper.PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	var apps []cloudapi.AppAttribute

	for {
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DescribeAppAttributes(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_api_gateway_apps", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*cloudapi.DescribeAppAttributesResponse)

		apps = append(apps, response.Apps.AppAttribute...)

		if len(apps) < helper.PageSizeLarge {
			break
		}

		page, err := helper.GetNextpageNumber(request.PageNumber)
		if err != nil {
			return errs.WrapError(err)
		}
		request.PageNumber = page
	}

	var filteredAppsTemp []cloudapi.AppAttribute
	nameRegex, ok := d.GetOk("name_regex")
	var gatewayAppNameRegex *regexp.Regexp
	if ok && nameRegex.(string) != "" {
		r, err := regexp.Compile(nameRegex.(string))
		if err != nil {
			return errs.WrapError(err)
		}
		gatewayAppNameRegex = r
	}

	// ids
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	for _, app := range apps {
		if gatewayAppNameRegex != nil && !gatewayAppNameRegex.MatchString(app.AppName) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[strconv.FormatInt(app.AppId, 10)]; !ok {
				continue
			}
		}
		if value, ok := d.GetOk("tags"); ok {
			tags, err := cloudApiService.DescribeTags(strconv.FormatInt(app.AppId, 10), value.(map[string]interface{}), tags.TagResourceApp)
			if err != nil {
				return errs.WrapError(err)
			}
			if len(tags) < 1 {
				continue
			}
		}

		filteredAppsTemp = append(filteredAppsTemp, app)
	}

	return apigatewayAppsDecriptionAttributes(d, filteredAppsTemp, meta)
}

func apigatewayAppsDecriptionAttributes(d *schema.ResourceData, apps []cloudapi.AppAttribute, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var ids []string
	var s []map[string]interface{}
	var names []string
	for _, app := range apps {
		request := cloudapi.CreateDescribeAppSecurityRequest()
		request.RegionId = client.RegionId
		request.AppId = requests.NewInteger64(app.AppId)
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DescribeAppSecurity(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_api_gateway_apps", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*cloudapi.DescribeAppSecurityResponse)

		mapping := map[string]interface{}{
			"id":            app.AppId,
			"name":          app.AppName,
			"description":   app.Description,
			"created_time":  app.CreatedTime,
			"modified_time": app.ModifiedTime,
			"app_code":      response.AppCode,
		}
		ids = append(ids, strconv.FormatInt(app.AppId, 10))
		s = append(s, mapping)
		names = append(names, app.AppName)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("apps", s); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}
	return nil
}
