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

func dataSourceAlicloudServiceCatalogProvisionedProducts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudServiceCatalogProvisionedProductsRead,
		Schema: map[string]*schema.Schema{
			"access_level_filter": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"sort_by": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"sort_order": {
				Optional: true,
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
			"enable_details": {
				Optional: true,
				Type:     schema.TypeBool,
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
			"products": {
				Deprecated: "Field 'products' has been deprecated from provider version 1.197.0.",
				Computed:   true,
				Type:       schema.TypeList,
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
						"last_provisioning_task_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"last_successful_provisioning_task_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"last_task_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"outputs": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"output_key": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"output_value": {
										Computed: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"owner_principal_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"owner_principal_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"parameters": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"parameter_key": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"parameter_value": {
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
						"product_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"product_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"product_version_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"product_version_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"provisioned_product_arn": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"provisioned_product_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"provisioned_product_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"provisioned_product_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"stack_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"stack_region_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"status_message": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"tags": tagsSchema(),
					},
				},
			},
			"provisioned_products": {
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
						"last_provisioning_task_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"last_successful_provisioning_task_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"last_task_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"outputs": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"output_key": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"output_value": {
										Computed: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"owner_principal_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"owner_principal_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"parameters": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"parameter_key": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"parameter_value": {
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
						"product_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"product_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"product_version_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"product_version_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"provisioned_product_arn": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"provisioned_product_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"provisioned_product_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"provisioned_product_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"stack_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"stack_region_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"status_message": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"tags": tagsSchema(),
					},
				},
			},
		},
	}
}

func dataSourceAlicloudServiceCatalogProvisionedProductsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	if v, ok := d.GetOk("access_level_filter"); ok {
		request["AccessLevelFilter"] = v
	}
	if v, ok := d.GetOk("sort_by"); ok {
		request["SortBy"] = v
	}
	if v, ok := d.GetOk("sort_order"); ok {
		request["SortOrder"] = v
	}
	if v, ok := d.GetOk("page_number"); ok && v.(int) > 0 {
		request["PageNumber"] = v.(int)
	} else {
		request["PageNumber"] = 1
	}
	if v, ok := d.GetOk("page_size"); ok && v.(int) > 0 {
		request["PageSize"] = v.(int)
	} else {
		request["PageSize"] = PageSizeLarge
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

	var provisionedProductNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		provisionedProductNameRegex = r
	}

	var err error
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "ListProvisionedProducts"
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_service_catalog_provisioned_products", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.ProvisionedProductDetails", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ProvisionedProductDetails", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ProvisionedProductId"])]; !ok {
					continue
				}
			}

			if provisionedProductNameRegex != nil && !provisionedProductNameRegex.MatchString(fmt.Sprint(item["ProvisionedProductName"])) {
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
	servicecatalogService := ServicecatalogService{client}
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                                   fmt.Sprint(object["ProvisionedProductId"]),
			"create_time":                          object["CreateTime"],
			"last_provisioning_task_id":            object["LastProvisioningTaskId"],
			"last_successful_provisioning_task_id": object["LastSuccessfulProvisioningTaskId"],
			"last_task_id":                         object["LastTaskId"],
			"owner_principal_id":                   object["OwnerPrincipalId"],
			"owner_principal_type":                 object["OwnerPrincipalType"],
			"portfolio_id":                         object["PortfolioId"],
			"product_id":                           object["ProductId"],
			"product_name":                         object["ProductName"],
			"product_version_id":                   object["ProductVersionId"],
			"product_version_name":                 object["ProductVersionName"],
			"provisioned_product_arn":              object["ProvisionedProductArn"],
			"provisioned_product_id":               object["ProvisionedProductId"],
			"provisioned_product_name":             object["ProvisionedProductName"],
			"provisioned_product_type":             object["ProvisionedProductType"],
			"stack_id":                             object["StackId"],
			"stack_region_id":                      object["StackRegionId"],
			"status":                               object["Status"],
			"status_message":                       object["StatusMessage"],
		}

		ids = append(ids, fmt.Sprint(object["ProvisionedProductId"]))
		names = append(names, object["ProvisionedProductName"])

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["ProvisionedProductId"])
		object, err = servicecatalogService.DescribeServiceCatalogProvisionedProduct(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["create_time"] = object["CreateTime"]
		mapping["last_provisioning_task_id"] = object["LastProvisioningTaskId"]
		mapping["last_successful_provisioning_task_id"] = object["LastSuccessfulProvisioningTaskId"]
		mapping["last_task_id"] = object["LastTaskId"]
		outputs67Maps := make([]map[string]interface{}, 0)
		outputs67Raw := object["Outputs"]
		for _, value0 := range outputs67Raw.([]interface{}) {
			outputs67 := value0.(map[string]interface{})
			outputs67Map := make(map[string]interface{})
			outputs67Map["description"] = outputs67["Description"]
			outputs67Map["output_key"] = outputs67["OutputKey"]
			outputs67Map["output_value"] = outputs67["OutputValue"]
			outputs67Maps = append(outputs67Maps, outputs67Map)
		}
		mapping["outputs"] = outputs67Maps
		mapping["owner_principal_id"] = object["OwnerPrincipalId"]
		mapping["owner_principal_type"] = object["OwnerPrincipalType"]
		parameters31Maps := make([]map[string]interface{}, 0)
		parameters31Raw := object["Parameters"]
		for _, value0 := range parameters31Raw.([]interface{}) {
			parameters31 := value0.(map[string]interface{})
			parameters31Map := make(map[string]interface{})
			parameters31Map["parameter_key"] = parameters31["ParameterKey"]
			parameters31Map["parameter_value"] = parameters31["ParameterValue"]
			parameters31Maps = append(parameters31Maps, parameters31Map)
		}
		mapping["parameters"] = parameters31Maps
		mapping["portfolio_id"] = object["PortfolioId"]
		mapping["product_id"] = object["ProductId"]
		mapping["product_name"] = object["ProductName"]
		mapping["product_version_id"] = object["ProductVersionId"]
		mapping["product_version_name"] = object["ProductVersionName"]
		mapping["provisioned_product_arn"] = object["ProvisionedProductArn"]
		mapping["provisioned_product_id"] = object["ProvisionedProductId"]
		mapping["provisioned_product_name"] = object["ProvisionedProductName"]
		mapping["provisioned_product_type"] = object["ProvisionedProductType"]
		mapping["stack_id"] = object["StackId"]
		mapping["stack_region_id"] = object["StackRegionId"]
		mapping["status"] = object["Status"]
		mapping["status_message"] = object["StatusMessage"]
		tagsRaw, _ := jsonpath.Get("$.TaskTags", object)
		if tagsRaw != nil {
			mapping["tags"] = tagsToMap(tagsRaw)
		}
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("products", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("provisioned_products", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
