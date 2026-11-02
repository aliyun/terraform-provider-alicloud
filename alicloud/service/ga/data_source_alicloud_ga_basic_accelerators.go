package ga

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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceAliCloudGaBasicAccelerators() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudGaBasicAcceleratorsRead,
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
			"accelerator_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bandwidth_billing_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"BandwidthPackage", "CDT", "CDT95"}, false),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"init", "active", "configuring", "binding", "unbinding", "deleting", "finacialLocked"}, false),
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  helper.PageSizeLarge,
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
			"accelerators": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"basic_accelerator_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"basic_accelerator_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"basic_endpoint_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"basic_ip_set_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth_billing_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"expired_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"basic_bandwidth_package": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"bandwidth": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"bandwidth_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"cross_domain_bandwidth_package": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"bandwidth": {
										Type:     schema.TypeInt,
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

func dataSourceAliCloudGaBasicAcceleratorsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListBasicAccelerators"
	request := make(map[string]interface{})
	helper.SetPagingRequest(d, request, helper.PageSizeLarge)
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("accelerator_id"); ok {
		request["AcceleratorId"] = v
	}

	if v, ok := d.GetOk("status"); ok {
		request["State"] = v
	}

	bandwidthBillingType, bandwidthBillingTypeOk := d.GetOk("bandwidth_billing_type")

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

	var basicAcceleratorNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
		}
		basicAcceleratorNameRegex = r
	}

	var response map[string]interface{}
	var err error
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_ga_basic_accelerators", action, errs.AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Accelerators", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Accelerators", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["AcceleratorId"])]; !ok {
					continue
				}
			}

			if basicAcceleratorNameRegex != nil && !basicAcceleratorNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}

			if bandwidthBillingTypeOk && bandwidthBillingType.(string) != "" && bandwidthBillingType.(string) != fmt.Sprint(item["BandwidthBillingType"]) {
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
	for _, object := range objects {

		create_timeVal, err := convert.FormatInt(object["CreateTime"])
		if err != nil {
			return err
		}
		expired_timeVal, err := convert.FormatInt(object["ExpiredTime"])
		if err != nil {
			return err
		}
		mapping := map[string]interface{}{
			"id":                      fmt.Sprint(object["AcceleratorId"]),
			"basic_accelerator_id":    fmt.Sprint(object["AcceleratorId"]),
			"basic_accelerator_name":  object["Name"],
			"basic_endpoint_group_id": object["BasicEndpointGroupId"],
			"basic_ip_set_id":         object["BasicIpSetId"],
			"bandwidth_billing_type":  object["BandwidthBillingType"],
			"instance_charge_type":    object["InstanceChargeType"],
			"description":             object["Description"],
			"region_id":               object["RegionId"],
			"create_time":             create_timeVal,
			"expired_time":            expired_timeVal,
			"status":                  object["State"],
		}

		if v, ok := object["BasicBandwidthPackage"]; ok {
			basicBandwidthPackageMaps := make([]map[string]interface{}, 0)
			basicBandwidthPackage := v.(map[string]interface{})
			basicBandwidthPackageMap := make(map[string]interface{})

			if instanceId, ok := basicBandwidthPackage["InstanceId"]; ok {
				basicBandwidthPackageMap["instance_id"] = instanceId
			}

			if bandwidth, ok := basicBandwidthPackage["Bandwidth"]; ok {
				basicBandwidthPackageMap["bandwidth"] = bandwidth
			}

			if bandwidthType, ok := basicBandwidthPackage["BandwidthType"]; ok {
				basicBandwidthPackageMap["bandwidth_type"] = bandwidthType
			}

			basicBandwidthPackageMaps = append(basicBandwidthPackageMaps, basicBandwidthPackageMap)
			mapping["basic_bandwidth_package"] = basicBandwidthPackageMaps
		}

		if v, ok := object["CrossDomainBandwidthPackage"]; ok {
			crossDomainBandwidthPackageMaps := make([]map[string]interface{}, 0)
			crossDomainBandwidthPackage := v.(map[string]interface{})
			crossDomainBandwidthPackageMap := make(map[string]interface{})

			if instanceId, ok := crossDomainBandwidthPackage["InstanceId"]; ok {
				crossDomainBandwidthPackageMap["instance_id"] = instanceId
			}

			if bandwidth, ok := crossDomainBandwidthPackage["Bandwidth"]; ok {
				crossDomainBandwidthPackageMap["bandwidth"] = bandwidth
			}

			crossDomainBandwidthPackageMaps = append(crossDomainBandwidthPackageMaps, crossDomainBandwidthPackageMap)
			mapping["cross_domain_bandwidth_package"] = crossDomainBandwidthPackageMaps
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["Name"])
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("accelerators", s); err != nil {
		return errs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
