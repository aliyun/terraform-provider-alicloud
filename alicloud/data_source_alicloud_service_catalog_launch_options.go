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

func dataSourceAlicloudServiceCatalogLaunchOptions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudServiceCatalogLaunchOptionsRead,
		Schema: map[string]*schema.Schema{
			"product_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
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
			"options": {
				Deprecated: "Field 'options' has been deprecated from provider version 1.197.0.",
				Computed:   true,
				Type:       schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"constraint_summaries": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"constraint_type": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"description": {
										Computed: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"portfolio_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"portfolio_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"launch_options": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"constraint_summaries": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"constraint_type": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"description": {
										Computed: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"portfolio_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"portfolio_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudServiceCatalogLaunchOptionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	if v, ok := d.GetOk("product_id"); ok {
		request["ProductId"] = v
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

	var productNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		productNameRegex = r
	}

	var err error
	var objects []interface{}
	var response map[string]interface{}
	action := "ListLaunchOptions"
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_service_catalog_launch_options", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.LaunchOptionSummaries", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.LaunchOptionSummaries", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if productNameRegex != nil && !productNameRegex.MatchString(fmt.Sprint(item["PortfolioName"])) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(request["ProductId"], ":", item["PortfolioId"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":             fmt.Sprint(request["ProductId"], ":", object["PortfolioId"]),
			"portfolio_id":   object["PortfolioId"],
			"portfolio_name": object["PortfolioName"],
		}

		constraintSummaries91Maps := make([]map[string]interface{}, 0)
		constraintSummaries91Raw := object["ConstraintSummaries"]
		for _, value0 := range constraintSummaries91Raw.([]interface{}) {
			constraintSummaries91 := value0.(map[string]interface{})
			constraintSummaries91Map := make(map[string]interface{})
			constraintSummaries91Map["constraint_type"] = constraintSummaries91["ConstraintType"]
			constraintSummaries91Map["description"] = constraintSummaries91["Description"]
			constraintSummaries91Maps = append(constraintSummaries91Maps, constraintSummaries91Map)
		}
		mapping["constraint_summaries"] = constraintSummaries91Maps

		ids = append(ids, fmt.Sprint(request["ProductId"], ":", object["PortfolioId"]))

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("options", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("launch_options", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
