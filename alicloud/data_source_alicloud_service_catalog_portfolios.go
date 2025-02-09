package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudServiceCatalogPortfolios() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudServiceCatalogPortfoliosRead,
		Schema: map[string]*schema.Schema{
			"product_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"scope": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Local", "Import", "All"}, false),
			},
			"sort_by": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"CreateTime"}, false),
			},
			"sort_order": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Asc", "Desc"}, false),
			},
			"ids": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name_regex": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.ValidateRegexp,
			},
			"names": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"page_number": {
				Optional: true,
				Type:     schema.TypeInt,
			},
			"page_size": {
				Optional: true,
				Type:     schema.TypeInt,
				Default:  10,
			},
			"portfolios": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"create_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"description": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"portfolio_arn": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"portfolio_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"portfolio_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"provider_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudServiceCatalogPortfoliosRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	if v, ok := d.GetOk("product_id"); ok {
		request["ProductId"] = v
	}
	if v, ok := d.GetOk("scope"); ok {
		request["Scope"] = v
	}
	if v, ok := d.GetOk("sort_by"); ok {
		request["SortBy"] = v
	}
	if v, ok := d.GetOk("sort_order"); ok {
		request["SortOrder"] = v
	}
	setPagingRequest(d, request, PageSizeLarge)

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var portfolioNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		portfolioNameRegex = r
	}

	var err error
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "ListPortfolios"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := client.RpcPost("servicecatalog", "2021-09-01", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			response = resp
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_service_catalog_portfolios", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.PortfolioDetails", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.PortfolioDetails", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["PortfolioId"])]; !ok {
					continue
				}
			}

			if portfolioNameRegex != nil && !portfolioNameRegex.MatchString(fmt.Sprint(item["PortfolioName"])) {
				continue
			}
			objects = append(objects, item)
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})

		mapping := map[string]interface{}{
			"id": fmt.Sprint(object["PortfolioId"]),
		}
		mapping["create_time"] = object["CreateTime"]
		mapping["description"] = object["Description"]
		mapping["portfolio_arn"] = object["PortfolioArn"]
		mapping["portfolio_id"] = object["PortfolioId"]
		mapping["portfolio_name"] = object["PortfolioName"]
		mapping["provider_name"] = object["ProviderName"]

		ids = append(ids, fmt.Sprint(object["PortfolioId"]))
		names = append(names, object["PortfolioName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("portfolios", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
