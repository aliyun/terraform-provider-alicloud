package mse

import (
	"fmt"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceMseGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudMseGatewaysRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"gateway_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"0", "1", "2", "3", "4", "6", "8", "9", "10", "11", "12", "13"}, false),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gateways": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_unique_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"replica": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slb_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"associate_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"slb_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"slb_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"slb_port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"gmt_create": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"gateway_slb_mode": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"gateway_slb_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAlicloudMseGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListGateway"
	request := make(map[string]interface{})
	filterParams := make(map[string]interface{})

	if v, ok := d.GetOk("gateway_name"); ok {
		filterParams["Name"] = v
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		filterParams["Vpc"] = v
	}

	if len(filterParams) > 0 {
		if v, err := convert.ConvertMaptoJsonString(filterParams); err != nil {
			return errs.WrapError(err)
		} else {
			request["FilterParams"] = v
		}
	}

	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var gatewayNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
		}
		gatewayNameRegex = r
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
	status, statusOk := d.GetOk("status")
	var response map[string]interface{}
	var err error
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("mse", "2019-05-31", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_mse_gateways", action, errs.AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Data.Result", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Data.Result", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if gatewayNameRegex != nil && !gatewayNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["GatewayUniqueId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != fmt.Sprint(item["Status"]) {
				continue
			}
			objects = append(objects, item)
		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"gateway_name":      object["Name"],
			"id":                fmt.Sprint(object["GatewayUniqueId"]),
			"gateway_unique_id": fmt.Sprint(object["GatewayUniqueId"]),
			"payment_type":      object["ChargeType"],
			"replica":           fmt.Sprint(object["Replica"]),
			"spec":              object["Spec"],
			"status":            fmt.Sprint(object["Status"]),
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["Name"])
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["GatewayUniqueId"])
		mseService := NewService(client)
		getResp, err := mseService.DescribeMseGateway(id)
		if err != nil {
			return errs.WrapError(err)
		}
		mapping["backup_vswitch_id"] = getResp["Vswitch2"]
		mapping["vswitch_id"] = getResp["Vswitch"]
		mapping["vpc_id"] = getResp["Vpc"]

		slbListObject, err := mseService.ListGatewaySlb(id)
		if err != nil {
			return errs.WrapError(err)
		}
		slbList := slbListObject["SlbList"]

		slbMapList := make([]map[string]interface{}, 0)
		for _, v := range slbList.([]interface{}) {
			slbMap := v.(map[string]interface{})
			slb := map[string]interface{}{
				"associate_id":       slbMap["Id"],
				"slb_id":             slbMap["SlbId"],
				"slb_ip":             slbMap["SlbIp"],
				"slb_port":           slbMap["SlbPort"],
				"type":               slbMap["Type"],
				"gmt_create":         slbMap["GmtCreate"],
				"gateway_slb_mode":   slbMap["GatewaySlbMode"],
				"gateway_slb_status": slbMap["GatewaySlbStatus"],
			}
			slbMapList = append(slbMapList, slb)
		}
		mapping["slb_list"] = slbMapList
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("gateways", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
