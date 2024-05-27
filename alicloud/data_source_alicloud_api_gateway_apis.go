package alicloud

import (
	"fmt"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudApiGatewayApis() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudApiGatewayApisRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"api_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"apis": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"api_id": {
							Type:     schema.TypeString,
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
						"group_name": {
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
		},
	}
}

func dataSourceAliCloudApiGatewayApisRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := cloudapi.CreateDescribeApisRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	if v, ok := d.GetOk("group_id"); ok {
		request.GroupId = v.(string)
	}

	if v, ok := d.GetOk("api_id"); ok {
		request.ApiId = v.(string)
	}

	var objects []cloudapi.ApiSummary
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var apiNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		apiNameRegex = r
	}

	for {
		response, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DescribeApis(request)
		})
		addDebug(request.GetActionName(), response, request.RpcRequest, request)

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_api_gateway_apis", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		result := response.(*cloudapi.DescribeApisResponse).ApiSummarys.ApiSummary
		for _, item := range result {
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprintf("%v:%v", item.GroupId, item.ApiId)]; !ok {
					continue
				}
			}

			if apiNameRegex != nil && !apiNameRegex.MatchString(fmt.Sprint(item.ApiName)) {
				continue
			}

			objects = append(objects, item)
		}

		if len(result) < PageSizeLarge {
			break
		}

		pageNum, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}

		request.PageNumber = pageNum
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":          fmt.Sprintf("%v:%v", object.GroupId, object.ApiId),
			"group_id":    fmt.Sprint(object.GroupId),
			"api_id":      fmt.Sprint(object.ApiId),
			"name":        object.ApiName,
			"description": object.Description,
			"group_name":  object.GroupName,
			"region_id":   object.RegionId,
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object.ApiName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("apis", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
