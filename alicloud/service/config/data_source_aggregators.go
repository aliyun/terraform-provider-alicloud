package config

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceAggregators() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceAggregatorsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
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
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Creating", "Deleting", "Normal"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"aggregators": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"aggregator_accounts": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"account_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"account_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"account_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"aggregator_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"aggregator_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"aggregator_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
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

func DataSourceAggregatorsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListAggregators"
	request := make(map[string]interface{})
	request["MaxResults"] = helper.PageSizeLarge
	var objects []map[string]interface{}
	var aggregatorNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
		}
		aggregatorNameRegex = r
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
			response, err = client.RpcGet("Config", "2020-09-07", action, request, nil)
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
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_config_aggregators", action, errs.AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.AggregatorsResult.Aggregators", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.AggregatorsResult.Aggregators", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if aggregatorNameRegex != nil {
				if !aggregatorNameRegex.MatchString(fmt.Sprint(item["AggregatorName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["AggregatorId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" {
				aggregatorStatusVal, aggregatorStatusErr := convert.FormatInt(item["AggregatorStatus"])
				if aggregatorStatusErr != nil {
					return aggregatorStatusErr
				}
				if status.(string) != convertConfigAggregatorStatusResponse(aggregatorStatusVal) {
					continue
				}
			}
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {

		actualStatus, err := convert.FormatInt(object["AggregatorStatus"])
		if err != nil {
			return err
		}
		mapping := map[string]interface{}{
			"account_id":      fmt.Sprint(object["AccountId"]),
			"id":              fmt.Sprint(object["AggregatorId"]),
			"aggregator_id":   fmt.Sprint(object["AggregatorId"]),
			"aggregator_name": object["AggregatorName"],
			"aggregator_type": object["AggregatorType"],
			"description":     object["Description"],
			"status":          convertConfigAggregatorStatusResponse(actualStatus),
		}
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			names = append(names, object["AggregatorName"])
			s = append(s, mapping)
			continue
		}

		configService := NewConfigService(client)
		id := fmt.Sprint(object["AggregatorId"])
		getResp, err := configService.DescribeConfigAggregator(id)
		if err != nil {
			return errs.WrapError(err)
		}

		aggregatorAccounts := make([]map[string]interface{}, 0)
		if aggregatorAccountsList, ok := getResp["AggregatorAccounts"].([]interface{}); ok {
			for _, v := range aggregatorAccountsList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"account_id":   fmt.Sprint(m1["AccountId"]),
						"account_name": m1["AccountName"],
						"account_type": m1["AccountType"],
					}
					aggregatorAccounts = append(aggregatorAccounts, temp1)
				}
			}
		}
		mapping["aggregator_accounts"] = aggregatorAccounts
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["AggregatorName"])
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("aggregators", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
