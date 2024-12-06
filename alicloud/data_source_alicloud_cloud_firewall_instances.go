package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudCloudFirewallInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCloudFirewallInstancesRead,
		Schema: map[string]*schema.Schema{
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Subscription", "PayAsYouGo"}, false),
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

func dataSourceAlicloudCloudFirewallInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "QueryAvailableInstances"
	request := make(map[string]interface{})
	var objects []map[string]interface{}
	var response map[string]interface{}
	var err error
	var endpoint string
	var productTypes []string
	paymentType := d.Get("payment_type").(string)
	productMapping := map[string]string{
		"vipcloudfw":               "vipcloudfw",
		"cfw_elasticity_public_cn": "cfw",
	}
	paymentTypeMapping := map[string]string{
		"vipcloudfw":               "Subscription",
		"cfw_elasticity_public_cn": "PayAsYouGo",
	}
	if paymentType == "Subscription" {
		delete(productMapping, "cfw_elasticity_public_cn")
		productTypes = []string{"vipcloudfw"}
	} else if paymentType == "PayAsYouGo" {
		delete(productMapping, "vipcloudfw")
		productTypes = []string{"cfw_elasticity_public_cn"}
	} else {
		productTypes = []string{"vipcloudfw", "cfw_elasticity_public_cn"}
	}
	for _, productType := range productTypes {
		request["PageSize"] = PageSizeLarge
		request["PageNum"] = 1
		request["ProductCode"] = productMapping[productType]
		request["ProductType"] = productType
		request["SubscriptionType"] = paymentTypeMapping[productType]
		attempts := 0
		for {
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					if IsExpectedErrors(err, []string{"NotApplicable", NotFoundArticle}) {
						request["ProductCode"] = "cfw"
						request["ProductType"] = "cfw_pre_intl"
						if productType == "cfw_elasticity_public_cn" {
							request["ProductType"] = "cfw_elasticity_public_intl"
						}
						endpoint = connectivity.BssOpenAPIEndpointInternational
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				} else if attempts < 1 {
					resp, err := jsonpath.Get("$.Data.InstanceList", response)
					if err != nil {
						log.Println(WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.InstanceList", response))
					} else if len(resp.([]interface{})) < 1 {
						attempts++
						request["ProductCode"] = "cfw"
						request["ProductType"] = "cfw_pre_intl"
						if productType == "cfw_elasticity_public_cn" {
							request["ProductType"] = "cfw_elasticity_public_intl"
						}
						endpoint = connectivity.BssOpenAPIEndpointInternational
						return resource.RetryableError(fmt.Errorf("try to fetching international instances"))
					}
				}
				return nil
			})
			addDebug(action, response, request)
			if err != nil {
				return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cloud_firewall_instances", action, AlibabaCloudSdkGoERROR)
			}
			resp, err := jsonpath.Get("$.Data.InstanceList", response)
			if err != nil {
				return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.InstanceList", response)
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item := v.(map[string]interface{})
				objects = append(objects, item)
			}
			if len(result) < PageSizeLarge {
				break
			}
			request["PageNum"] = request["PageNum"].(int) + 1
		}
	}
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
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
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
