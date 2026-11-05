package amqpopen

import (
	"fmt"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceAmqpInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudAmqpInstancesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"DEPLOYING", "SERVING", "EXPIRED", "RELEASED"}, false),
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_end_point": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"support_eip": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"serverless_switch": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"renewal_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"renewal_duration": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"renewal_duration_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expire_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudAmqpInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListInstances"
	request := make(map[string]interface{})
	request["MaxResults"] = helper.PageSizeLarge

	status, statusOk := d.GetOk("status")

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

	var instanceNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
		}
		instanceNameRegex = r
	}

	var response map[string]interface{}
	var err error

	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcGet("amqp-open", "2019-12-12", action, request, nil)
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
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_amqp_instances", action, errs.AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Data.Instances", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Data.Instances", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["InstanceId"])]; !ok {
					continue
				}
			}

			if instanceNameRegex != nil && !instanceNameRegex.MatchString(fmt.Sprint(item["InstanceName"])) {
				continue
			}

			if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
				continue
			}
			objects = append(objects, item)
		}

		if nextToken, ok := response["Data"].(map[string]interface{})["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                fmt.Sprint(object["InstanceId"]),
			"instance_id":       fmt.Sprint(object["InstanceId"]),
			"instance_type":     convertAmqpInstanceInstanceTypeResponse(object["InstanceType"]),
			"instance_name":     object["InstanceName"],
			"public_endpoint":   object["PublicEndpoint"],
			"private_end_point": object["PrivateEndpoint"],
			"support_eip":       object["SupportEIP"],
			"serverless_switch": object["ServerlessSwitch"],
			"status":            object["Status"],
			"create_time":       fmt.Sprint(object["OrderCreateTime"]),
			"expire_time":       fmt.Sprint(object["ExpireTime"]),
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["InstanceName"])

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(object["InstanceId"])
		amqpServiceV2 := NewAmqpServiceV2(client)

		instancesAccountDetail, err := amqpServiceV2.QueryAvailableInstances(id, client.RegionId, "ons", "ons_onsproxy_pre", "ons", "ons_onsproxy_public_intl")
		if err != nil {
			return errs.WrapError(err)
		}

		mapping["payment_type"] = instancesAccountDetail["SubscriptionType"]
		mapping["renewal_status"] = instancesAccountDetail["RenewStatus"]
		if v, ok := instancesAccountDetail["RenewalDuration"]; ok && fmt.Sprint(v) != "0" {
			renewalDuration, err := convert.FormatInt(v)
			if err != nil {
				return err
			}
			mapping["renewal_duration"] = renewalDuration
		}
		mapping["renewal_duration_unit"] = convert.ConvertAmqpInstanceRenewalDurationUnitResponse(instancesAccountDetail["RenewalDurationUnit"])

		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("instances", s); err != nil {
		return errs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
