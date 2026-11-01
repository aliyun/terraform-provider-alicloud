package cloudstoragegateway

import (
	"fmt"
	"strconv"
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

func DataSourceAlicloudCloudStorageGatewayStocks() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceAlicloudCloudStorageGatewayStocksRead,
		Schema: map[string]*schema.Schema{
			"gateway_class": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Basic", "Standard", "Enhanced", "Advanced"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"stocks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"available_gateway_classes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func DataSourceAlicloudCloudStorageGatewayStocksRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DescribeGatewayStock"
	request := make(map[string]interface{})
	request["GatewayRegionId"] = client.RegionId
	gatewayClass, gatewayClassOk := d.GetOk("gateway_class")
	var objects []map[string]interface{}
	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("sgw", "2018-05-11", action, nil, request, true)
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
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cloud_storage_gateway_stocks", action, errs.AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	resp, err := jsonpath.Get("$.Stocks.Stock", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Stocks.Stock", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		stockInfo, _ := convert.ConvertJsonStringToMap(item["StockInfo"].(string))
		if gatewayClassOk && gatewayClass.(string) != "" && !stockInfo[gatewayClass.(string)].(bool) {
			continue
		}

		stockInfoKeys := make([]string, 0)
		for k, v := range stockInfo {
			if v.(bool) {
				stockInfoKeys = append(stockInfoKeys, k)
			}
		}

		if len(stockInfoKeys) > 0 {
			object := map[string]interface{}{
				"zone_id":                   fmt.Sprint(item["ZoneId"]),
				"available_gateway_classes": stockInfoKeys,
			}

			objects = append(objects, object)
		}
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))
	if err := d.Set("stocks", objects); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), objects)
	}

	return nil
}
