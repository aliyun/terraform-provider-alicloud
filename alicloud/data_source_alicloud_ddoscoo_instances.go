package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudDdoscooInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDdoscooInstancesRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
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
	var objects []map[string]interface{}
	conn, err := client.NewDdoscooClient()
	if err != nil {
		return WrapError(err)
	}
	// describe ddoscoo instance filtered by name_regex
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &runtime)
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
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		specResponse, err = conn.DoRequest(StringPointer(specAction), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, specReq, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
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
		mapping := map[string]interface{}{
			"id":          fmt.Sprint(item["InstanceId"]),
			"name":        item["Remark"],
			"remark":      item["Remark"],
			"ip_mode":     item["IpMode"],
			"debt_status": formatInt(item["DebtStatus"]),
			"edition":     formatInt(item["Edition"]),
			"ip_version":  item["IpVersion"],
			"status":      formatInt(item["Status"]),
			"enabled":     formatInt(item["Enabled"]),
			"expire_time": formatInt(item["ExpireTime"]),
			"create_time": formatInt(item["CreateTime"]),
		}

		if v, ok := instanceSpecs[fmt.Sprint(item["InstanceId"])]; ok {
			mapping["bandwidth"] = formatInt(v["ElasticBandwidth"])
			mapping["base_bandwidth"] = formatInt(v["BaseBandwidth"])
			mapping["service_bandwidth"] = formatInt(v["BandwidthMbps"])
			mapping["port_count"] = formatInt(v["PortLimit"])
			mapping["domain_count"] = formatInt(v["DomainLimit"])
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
