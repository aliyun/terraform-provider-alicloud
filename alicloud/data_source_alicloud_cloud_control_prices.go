// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudCloudControlPrices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCloudControlPriceRead,
		Schema: map[string]*schema.Schema{
			"desire_attributes": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"product": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"prices": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"currency": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"discount_price": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"module_details": {
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
									"cost_after_discount": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"original_cost": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"price_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"invoice_discount": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
								},
							},
						},
						"original_price": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"promotion_details": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"promotion_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"promotion_desc": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"promotion_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"trade_price": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func dataSourceAliCloudCloudControlPriceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var request map[string]interface{}
	var response map[string]interface{}
	action := fmt.Sprintf("/api/v1/providers/%s/products/%s/price/%s", "aliyun", d.Get("product").(string), d.Get("resource_code").(string))
	conn, err := client.NewCloudcontrolClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query := make(map[string]*string)
	body := make(map[string]interface{})
	query["regionId"] = StringPointer(client.RegionId)
	if v, ok := d.GetOk("desire_attributes"); ok {
		query["resourceAttributes"] = StringPointer(convertObjectToJsonString(v))
	}

	body = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2022-08-30"), nil, StringPointer("GET"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	resp, _ := jsonpath.Get("$.body.price", response)
	objects = append(objects, resp.(map[string]interface{}))

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}
		mapping["currency"] = objectRaw["currency"]
		mapping["discount_price"] = objectRaw["discountPrice"]
		mapping["original_price"] = objectRaw["originalPrice"]
		mapping["trade_price"] = objectRaw["tradePrice"]

		moduleDetails1Raw := objectRaw["moduleDetails"]
		moduleDetailsMaps := make([]map[string]interface{}, 0)
		if moduleDetails1Raw != nil {
			for _, moduleDetailsChild1Raw := range moduleDetails1Raw.([]interface{}) {
				moduleDetailsMap := make(map[string]interface{})
				moduleDetailsChild1Raw := moduleDetailsChild1Raw.(map[string]interface{})
				moduleDetailsMap["cost_after_discount"] = moduleDetailsChild1Raw["costAfterDiscount"]
				moduleDetailsMap["invoice_discount"] = moduleDetailsChild1Raw["invoiceDiscount"]
				moduleDetailsMap["module_code"] = moduleDetailsChild1Raw["moduleCode"]
				moduleDetailsMap["module_name"] = moduleDetailsChild1Raw["moduleName"]
				moduleDetailsMap["original_cost"] = moduleDetailsChild1Raw["originalCost"]
				moduleDetailsMap["price_type"] = moduleDetailsChild1Raw["priceType"]

				moduleDetailsMaps = append(moduleDetailsMaps, moduleDetailsMap)
			}
		}
		mapping["module_details"] = moduleDetailsMaps
		promotionDetails1Raw := objectRaw["promotionDetails"]
		promotionDetailsMaps := make([]map[string]interface{}, 0)
		if promotionDetails1Raw != nil {
			for _, promotionDetailsChild1Raw := range promotionDetails1Raw.([]interface{}) {
				promotionDetailsMap := make(map[string]interface{})
				promotionDetailsChild1Raw := promotionDetailsChild1Raw.(map[string]interface{})
				promotionDetailsMap["promotion_desc"] = promotionDetailsChild1Raw["promotionDesc"]
				promotionDetailsMap["promotion_id"] = promotionDetailsChild1Raw["promotionId"]
				promotionDetailsMap["promotion_name"] = promotionDetailsChild1Raw["promotionName"]

				promotionDetailsMaps = append(promotionDetailsMaps, promotionDetailsMap)
			}
		}
		mapping["promotion_details"] = promotionDetailsMaps

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw[""])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("prices", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
