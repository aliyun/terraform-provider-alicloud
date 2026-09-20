package alicloud

import (
	"encoding/json"
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudRosTemplateEstimateCost() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRosTemplateEstimateCostRead,
		Schema: map[string]*schema.Schema{
			"template_body": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"template_body", "template_url", "template_id", "template_scratch_id"},
			},
			"template_url": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"template_body", "template_url", "template_id", "template_scratch_id"},
			},
			"template_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"template_body", "template_url", "template_id", "template_scratch_id"},
			},
			"template_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_scratch_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"template_body", "template_url", "template_id", "template_scratch_id"},
			},
			"template_scratch_region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"stack_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_token": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parameters": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 200,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter_key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"parameter_value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alias_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"success": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"result": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"inquiry_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"order": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"currency": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"original_amount": {
													Type:     schema.TypeFloat,
													Computed: true,
												},
												"discount_amount": {
													Type:     schema.TypeFloat,
													Computed: true,
												},
												"trade_amount": {
													Type:     schema.TypeFloat,
													Computed: true,
												},
												"tax_amount": {
													Type:     schema.TypeFloat,
													Computed: true,
												},
												"total_cost_amount": {
													Type:     schema.TypeFloat,
													Computed: true,
												},
												"stand_price": {
													Type:     schema.TypeFloat,
													Computed: true,
												},
												"rule_ids": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"order_details": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"module_code": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"module_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"currency": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"original_amount": {
													Type:     schema.TypeFloat,
													Computed: true,
												},
												"discount_amount": {
													Type:     schema.TypeFloat,
													Computed: true,
												},
												"trade_amount": {
													Type:     schema.TypeFloat,
													Computed: true,
												},
											},
										},
									},
									"order_supplement": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"charge_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"period": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"period_unit": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"price_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"price_unit": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"quantity": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"auto_renew": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudRosTemplateEstimateCostRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "GetTemplateEstimateCost"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("template_body"); ok {
		request["TemplateBody"] = v
	}
	if v, ok := d.GetOk("template_url"); ok {
		request["TemplateURL"] = v
	}
	if v, ok := d.GetOk("template_id"); ok {
		request["TemplateId"] = v
	}
	if v, ok := d.GetOk("template_version"); ok {
		request["TemplateVersion"] = v
	}
	if v, ok := d.GetOk("template_scratch_id"); ok {
		request["TemplateScratchId"] = v
	}
	if v, ok := d.GetOk("template_scratch_region_id"); ok {
		request["TemplateScratchRegionId"] = v
	}
	if v, ok := d.GetOk("stack_id"); ok {
		request["StackId"] = v
	}
	if v, ok := d.GetOk("client_token"); ok {
		request["ClientToken"] = v
	}
	if v, ok := d.GetOk("parameters"); ok {
		parameters := make([]map[string]interface{}, 0)
		for _, item := range v.([]interface{}) {
			parameter := item.(map[string]interface{})
			parameters = append(parameters, map[string]interface{}{
				"ParameterKey":   parameter["parameter_key"],
				"ParameterValue": parameter["parameter_value"],
			})
		}
		request["Parameters"] = parameters
	}

	response, err := client.RpcPost("ROS", "2019-09-10", action, nil, request, true)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ros_template_estimate_cost", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	resp, err := jsonpath.Get("$.Resources", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Resources", response)
	}

	s := make([]map[string]interface{}, 0)
	if resourcesMap, ok := resp.(map[string]interface{}); ok {
		for resourceName, item := range resourcesMap {
			resource := item.(map[string]interface{})
			mapping := map[string]interface{}{
				"resource_name": resourceName,
				"resource_type": resource["Type"],
				"alias_type":    resource["AliasType"],
				"success":       resource["Success"],
			}
			result := make([]map[string]interface{}, 0)
			if resultRaw, ok := resource["Result"].(map[string]interface{}); ok {
				resultMapping := map[string]interface{}{
					"inquiry_type": resultRaw["InquiryType"],
				}
				order := make([]map[string]interface{}, 0)
				if orderRaw, ok := resultRaw["Order"].(map[string]interface{}); ok {
					ruleIds := make([]string, 0)
					if ruleIdsRaw, ok := orderRaw["RuleIds"].([]interface{}); ok {
						for _, ruleId := range ruleIdsRaw {
							ruleIds = append(ruleIds, fmt.Sprint(ruleId))
						}
					}
					order = append(order, map[string]interface{}{
						"currency":          orderRaw["Currency"],
						"original_amount":   orderRaw["OriginalAmount"],
						"discount_amount":   orderRaw["DiscountAmount"],
						"trade_amount":      orderRaw["TradeAmount"],
						"tax_amount":        orderRaw["TaxAmount"],
						"total_cost_amount": orderRaw["TotalCostAmount"],
						"stand_price":       orderRaw["StandPrice"],
						"rule_ids":          ruleIds,
					})
				}
				resultMapping["order"] = order

				orderDetails := make([]map[string]interface{}, 0)
				if orderDetailsRaw, ok := resultRaw["OrderDetails"].([]interface{}); ok {
					for _, detail := range orderDetailsRaw {
						detailMap := detail.(map[string]interface{})
						orderDetails = append(orderDetails, map[string]interface{}{
							"module_code":     detailMap["ModuleCode"],
							"module_name":     detailMap["ModuleName"],
							"currency":        detailMap["Currency"],
							"original_amount": detailMap["OriginalAmount"],
							"discount_amount": detailMap["DiscountAmount"],
							"trade_amount":    detailMap["TradeAmount"],
						})
					}
				}
				resultMapping["order_details"] = orderDetails

				orderSupplement := make([]map[string]interface{}, 0)
				if supplementRaw, ok := resultRaw["OrderSupplement"].(map[string]interface{}); ok {
					orderSupplement = append(orderSupplement, map[string]interface{}{
						"charge_type": supplementRaw["ChargeType"],
						"period":      supplementRaw["Period"],
						"period_unit": supplementRaw["PeriodUnit"],
						"price_type":  supplementRaw["PriceType"],
						"price_unit":  supplementRaw["PriceUnit"],
						"quantity":    supplementRaw["Quantity"],
						"auto_renew":  supplementRaw["AutoRenew"],
					})
				}
				resultMapping["order_supplement"] = orderSupplement
				result = append(result, resultMapping)
			}
			mapping["result"] = result
			s = append(s, mapping)
		}
	}

	requestJson, err := json.Marshal(request)
	if err != nil {
		return WrapError(err)
	}
	d.SetId(fmt.Sprintf("%d", helper.Hashcode(string(requestJson))))
	if err := d.Set("resources", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
