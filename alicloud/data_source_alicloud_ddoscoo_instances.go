package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlicloudDdoscooInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDdoscooInstancesRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
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
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"base_bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"service_bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"port_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"domain_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"remark": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"debt_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"edition": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ip_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"expire_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDdoscooInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeInstances"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeSmall
	request["PageNumber"] = 1

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(v.(string)); err == nil {
			nameRegex = r
		}
	}

	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
		request["InstanceIds"] = v
	}

	var response map[string]interface{}
	var err error
	var objects []map[string]interface{}
	// describe ddoscoo instance filtered by name_regex
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("ddoscoo", "2020-01-01", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ddoscoo_instances", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Instances", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Instances", response)
		}
		result, _ := resp.([]interface{})

		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["Remark"])) {
				continue
			}

			objects = append(objects, item)
		}

		if len(result) < PageSizeSmall {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	// describe instance spec filtered by instanceids
	var instanceIds []string
	for _, item := range objects {
		instanceIds = append(instanceIds, fmt.Sprint(item["InstanceId"]))
	}

	if len(instanceIds) < 1 {
		return WrapError(extractDdoscooInstance(d, map[string]map[string]interface{}{}, objects))
	}

	specAction := "DescribeInstanceSpecs"
	specReq := make(map[string]interface{})
	specReq["InstanceIds"] = instanceIds
	var specResponse map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		specResponse, err = client.RpcPost("ddoscoo", "2020-01-01", specAction, nil, specReq, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(specAction, specResponse, specReq)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ddoscoo_instances", specAction, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.InstanceSpecs", specResponse)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, specAction, "$.InstanceSpecs", specResponse)
	}
	result, _ := resp.([]interface{})
	specObjects := make(map[string]map[string]interface{}, len(instanceIds))
	for _, v := range result {
		item := v.(map[string]interface{})
		specObjects[fmt.Sprint(item["InstanceId"])] = item
	}

	return WrapError(extractDdoscooInstance(d, specObjects, objects))
}

func extractDdoscooInstance(d *schema.ResourceData, instanceSpecs map[string]map[string]interface{}, instance []map[string]interface{}) error {
	var instanceIds []string
	var names []string
	var s []map[string]interface{}

	for _, item := range instance {

		debt_statusVal, err234 := formatInt(item["DebtStatus"])
		if err234 != nil {
			return err234
		}
		editionVal, err235 := formatInt(item["Edition"])
		if err235 != nil {
			return err235
		}

		statusVal, err237 := formatInt(item["Status"])
		if err237 != nil {
			return err237
		}
		enabledVal, err238 := formatInt(item["Enabled"])
		if err238 != nil {
			return err238
		}
		expire_timeVal, err239 := formatInt(item["ExpireTime"])
		if err239 != nil {
			return err239
		}
		create_timeVal, err240 := formatInt(item["CreateTime"])
		if err240 != nil {
			return err240
		}
		mapping := map[string]interface{}{
			"id":          fmt.Sprint(item["InstanceId"]),
			"name":        item["Remark"],
			"remark":      item["Remark"],
			"ip_mode":     item["IpMode"],
			"debt_status": debt_statusVal,
			"edition":     editionVal,
			"ip_version":  item["IpVersion"],
			"status":      statusVal,
			"enabled":     enabledVal,
			"expire_time": expire_timeVal,
			"create_time": create_timeVal,
		}

		if v, ok := instanceSpecs[fmt.Sprint(item["InstanceId"])]; ok {
			v244, err244 := formatInt(v["ElasticBandwidth"])
			if err244 != nil {
				return err244
			}
			mapping["bandwidth"] = v244
			v245, err245 := formatInt(v["BaseBandwidth"])
			if err245 != nil {
				return err245
			}
			mapping["base_bandwidth"] = v245
			v246, err246 := formatInt(v["BandwidthMbps"])
			if err246 != nil {
				return err246
			}
			mapping["service_bandwidth"] = v246
			v247, err247 := formatInt(v["PortLimit"])
			if err247 != nil {
				return err247
			}
			mapping["port_count"] = v247
			v248, err248 := formatInt(v["DomainLimit"])
			if err248 != nil {
				return err248
			}
			mapping["domain_count"] = v248
		}

		instanceIds = append(instanceIds, fmt.Sprint(mapping["id"]))
		names = append(names, fmt.Sprint(mapping["name"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(instanceIds))
	if err := d.Set("ids", instanceIds); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
