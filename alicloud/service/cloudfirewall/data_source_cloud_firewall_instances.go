package cloudfirewall

import (
	"strconv"
	"time"

	jsonpath "github.com/PaesslerAG/jsonpath"
	connectivity "github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	errs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	helper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	validators "github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	retry "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceCloudFirewallInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceInstancesRead,
		Schema: map[string]*schema.Schema{
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"Subscription", "PayAsYouGo"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"renewal_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"renewal_duration_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "QueryAvailableInstances"
	request := make(map[string]interface{})
	var objects []map[string]interface{}
	var response map[string]interface{}
	var err error
	var productTypes []string
	paymentType := d.Get("payment_type").(string)
	productMapping := map[string]string{
		"vipcloudfw":                 "vipcloudfw",
		"cfw_elasticity_public_cn":   "cfw",
		"cfw_pre_intl":               "cfw",
		"cfw_elasticity_public_intl": "cfw",
	}
	paymentTypeMapping := map[string]string{
		"vipcloudfw":                 "Subscription",
		"cfw_elasticity_public_cn":   "PayAsYouGo",
		"cfw_pre_intl":               "Subscription",
		"cfw_elasticity_public_intl": "PayAsYouGo",
	}
	if paymentType == "Subscription" {
		if client.IsInternationalAccount() {
			productTypes = []string{"cfw_pre_intl"}
		} else {
			productTypes = []string{"vipcloudfw"}
		}
	} else if paymentType == "PayAsYouGo" {
		if client.IsInternationalAccount() {
			productTypes = []string{"cfw_elasticity_public_intl"}
		} else {
			productTypes = []string{"cfw_elasticity_public_cn"}
		}
	} else {
		if client.IsInternationalAccount() {
			productTypes = []string{"cfw_pre_intl", "cfw_elasticity_public_intl"}
		} else {
			productTypes = []string{"vipcloudfw", "cfw_elasticity_public_cn"}
		}
	}
	for _, productType := range productTypes {
		request["PageSize"] = helper.PageSizeLarge
		request["PageNum"] = 1
		request["ProductCode"] = productMapping[productType]
		request["ProductType"] = productType
		request["SubscriptionType"] = paymentTypeMapping[productType]
		for {
			wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
			err = retry.Retry(5*time.Minute, func() *retry.RetryError {
				response, err = client.RpcPost("BssOpenApi", "2017-12-14", action, nil, request, true)
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
				return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cloud_firewall_instances", action, errs.AlibabaCloudSdkGoERROR)
			}
			resp, err := jsonpath.Get("$.Data.InstanceList", response)
			if err != nil {
				return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Data.InstanceList", response)
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item := v.(map[string]interface{})
				objects = append(objects, item)
			}
			if len(result) < helper.PageSizeLarge {
				break
			}
			request["PageNum"] = request["PageNum"].(int) + 1
		}
	}
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		if v, ok := object["Status"]; ok && v.(string) != "Normal" {
			continue
		}
		mapping := map[string]interface{}{
			"status":                object["Status"],
			"end_time":              object["EndTime"],
			"create_time":           object["CreateTime"],
			"payment_type":          object["SubscriptionType"],
			"renewal_duration_unit": convertCloudFirewallInstanceRenewalDurationUnitResponse(object["RenewalDurationUnit"]),
			"id":                    object["InstanceID"],
			"renewal_status":        object["RenewStatus"],
			"instance_id":           object["InstanceID"],
		}
		s = append(s, mapping)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))

	if err := d.Set("instances", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
