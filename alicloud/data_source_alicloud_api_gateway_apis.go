package alicloud

import (
	"fmt"
	"log"
	"regexp"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudApiGatewayApis() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudApigatewayApisRead,

		Schema: map[string]*schema.Schema{
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
			"apis": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func dataSourceAlicloudApigatewayApisRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	args := cloudapi.CreateDescribeApisRequest()
	args.RegionId = client.RegionId

	if groupId, ok := d.GetOk("group_id"); ok {
		args.GroupId = groupId.(string)
	}
	if apiId, ok := d.GetOk("id"); ok {
		args.ApiId = apiId.(string)
	}

	args.PageSize = requests.NewInteger(PageSizeLarge)
	args.PageNumber = requests.NewInteger(1)

	var allapis []cloudapi.ApiSummary

	for {
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DescribeApis(args)
		})
		if err != nil {
			return err
		}
		resp, _ := raw.(*cloudapi.DescribeApisResponse)
		if err != nil {
			return err
		}

		if resp == nil {
			break
		}

		allapis = append(allapis, resp.ApiSummarys.ApiSummary...)

		if len(allapis) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return err
		} else {
			args.PageNumber = page
		}
	}

	var filteredApisTemp []cloudapi.ApiSummary
	nameRegex, ok := d.GetOk("name_regex")
	if ok && nameRegex.(string) != "" {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, api := range allapis {
			if r != nil && !r.MatchString(api.ApiName) {
				continue
			}

			filteredApisTemp = append(filteredApisTemp, api)
		}
	} else {
		filteredApisTemp = allapis
	}

	if len(filteredApisTemp) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	log.Printf("[DEBUG] alicloud_apigateway - api found: %#v", filteredApisTemp)

	return apiGatewayApisDescribeSummarys(d, filteredApisTemp, meta)
}

func apiGatewayApisDescribeSummarys(d *schema.ResourceData, apis []cloudapi.ApiSummary, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, api := range apis {
		mapping := map[string]interface{}{
			"id":          api.ApiId,
			"name":        api.ApiName,
			"region_id":   api.RegionId,
			"group_id":    api.GroupId,
			"group_name":  api.GroupName,
			"description": api.Description,
		}
		log.Printf("[DEBUG] alicloud_apigateway - adding api: %v", mapping)
		ids = append(ids, api.ApiId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("apis", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
