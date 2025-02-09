package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudServiceCatalogEndUserProducts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudServiceCatalogEndUserProductsRead,
		Schema: map[string]*schema.Schema{
			"sort_by": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
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
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
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
			"end_user_products": {
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
						"has_default_launch_option": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"product_arn": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"product_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"product_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"product_type": {
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

func dataSourceAlicloudServiceCatalogEndUserProductsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	if v, ok := d.GetOk("sort_by"); ok {
		request["SortBy"] = v
	}
	if v, ok := d.GetOk("sort_order"); ok {
		request["SortOrder"] = v
	}
	setPagingRequest(d, request, PageSizeLarge)

	var productNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		productNameRegex = r
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

	var err error
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "ListProductsAsEndUser"
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_service_catalog_end_user_products", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.ProductSummaries", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ProductSummaries", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}
		for _, v := range result {
			item := v.(map[string]interface{})
			if productNameRegex != nil && !productNameRegex.MatchString(fmt.Sprint(item["ProductName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ProductId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                        fmt.Sprint(object["ProductId"]),
			"create_time":               object["CreateTime"],
			"description":               object["Description"],
			"has_default_launch_option": object["HasDefaultLaunchOption"],
			"product_arn":               object["ProductArn"],
			"product_id":                object["ProductId"],
			"product_name":              object["ProductName"],
			"product_type":              object["ProductType"],
			"provider_name":             object["ProviderName"],
		}

		ids = append(ids, fmt.Sprint(object["ProductId"]))

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("end_user_products", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
