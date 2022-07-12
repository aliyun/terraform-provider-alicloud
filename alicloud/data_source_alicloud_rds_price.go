package alicloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudRdsPrice() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRdsPriceRead,

		Schema: map[string]*schema.Schema{
			"commodity_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"bards", "rds", "rords", "rds_rordspre_public_cn", "bards_intl", "rds_intl", "rords_intl", "rds_rordspre_public_intl"}, false),
			},
			"engine": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"MySQL", "SQLServer", "PostgreSQL", "MariaDB"}, false),
			},
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"db_instance_class": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"db_instance_storage": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"pay_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Prepaid", "Postpaid"}, false),
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"used_time": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"time_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Year", "Month"}, false),
			},
			"quantity": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"instance_used_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 3}),
			},
			"order_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"BUY", "UPGRADE", "RENEW"}, false),
			},
			"db_instance_storage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"local_ssd", "cloud_ssd", "cloud_essd", "cloud_essd2", "cloud_essd3"}, false),
			},
			"db_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"price": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"price_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"original_price": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"discount_price": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"currency": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"trade_price": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"rule_ids": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"coupons": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"is_selected": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"coupon_no": {
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
											},
										},
									},
								},
							},
						},
						"rules": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"rule_id": {
										Type:     schema.TypeString,
										Computed: true,
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

func dataSourceAlicloudRdsPriceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DescribePrice"
	request := make(map[string]interface{})
	request["SourceIp"] = client.SourceIp
	request["Engine"] = d.Get("engine")
	request["EngineVersion"] = d.Get("engine_version")
	request["DBInstanceClass"] = d.Get("db_instance_class")
	request["Quantity"] = d.Get("quantity")
	if v, ok := d.GetOk("commodity_code"); ok {
		request["CommodityCode"] = v
	}
	if v, ok := d.GetOk("db_instance_storage"); ok {
		request["DBInstanceStorage"] = v
	}
	if v, ok := d.GetOk("pay_type"); ok {
		request["PayType"] = v
	}
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}
	if v, ok := d.GetOk("used_time"); ok {
		request["UsedTime"] = v
	}
	if v, ok := d.GetOk("time_type"); ok {
		request["TimeType"] = v
	}
	if v, ok := d.GetOk("instance_used_type"); ok {
		request["InstanceUsedType"] = v
	}
	if v, ok := d.GetOk("order_type"); ok {
		request["OrderType"] = v
	}
	if v, ok := d.GetOk("db_instance_storage_type"); ok {
		request["DBInstanceStorageType"] = v
	}
	if v, ok := d.GetOk("db_instance_id"); ok {
		request["DBInstanceId"] = v
	}
	var response map[string]interface{}
	conn, err := client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_rds_price", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.PriceInfo", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.PriceInfo", response)
	}
	priceInfo, _ := resp.(map[string]interface{})
	ruleIds := priceInfo["RuleIds"].(map[string]interface{})["RuleId"].([]interface{})
	couponsResp := priceInfo["Coupons"].(map[string]interface{})["Coupon"].([]interface{})
	coupons := make([]map[string]interface{}, 0)

	for _, v := range couponsResp {
		v := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"is_selected": v["IsSelected"],
			"coupon_no":   v["CouponNo"],
			"name":        v["Name"],
			"description": v["Description"],
		}
		coupons = append(coupons, mapping)
	}
	priceInfos := make([]map[string]interface{}, 0)
	priceInfoMapping := map[string]interface{}{
		"original_price": priceInfo["OriginalPrice"],
		"discount_price": priceInfo["DiscountPrice"],
		"currency":       priceInfo["Currency"],
		"trade_price":    priceInfo["TradePrice"],
		"rule_ids":       ruleIds,
		"coupons":        coupons,
	}
	priceInfos = append(priceInfos, priceInfoMapping)
	rulesResp, err := jsonpath.Get("$.Rules.Rule", response)
	rules := rulesResp.([]interface{})

	rulesMapping := make([]map[string]interface{}, 0)
	for _, v := range rules {
		v := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"name":        v["Name"],
			"description": v["Description"],
			"rule_id":     v["RuleId"],
		}
		rulesMapping = append(rulesMapping, mapping)
	}
	s := make([]map[string]interface{}, 0)
	finalMapping := make(map[string]interface{})
	finalMapping["price_info"] = priceInfos
	finalMapping["rules"] = rulesMapping
	s = append(s, finalMapping)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))
	if err := d.Set("price", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
