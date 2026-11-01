package privatezone

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceEndpoints() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEndpointsRead,
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
				ValidateFunc: validation.StringInSlice([]string{"CHANGE_FAILED", "EXCEPTION", "FAILED", "INIT", "SUCCESS"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_configs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cidr_block": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vswitch_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"zone_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceEndpointsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeResolverEndpoints"
	request := make(map[string]interface{})
	request["Lang"] = "en"
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var endpointNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
		}
		endpointNameRegex = r
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
	var response map[string]interface{}
	var err error
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
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
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_pvtz_endpoints", action, errs.AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Endpoints", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Endpoints", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if endpointNameRegex != nil && !endpointNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["Id"])]; !ok {
					continue
				}
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
			"create_time":       object["CreateTime"],
			"endpoint_name":     object["Name"],
			"id":                object["Id"],
			"security_group_id": object["SecurityGroupId"],
			"status":            object["Status"],
			"vpc_id":            object["VpcId"],
			"vpc_name":          object["VpcName"],
			"vpc_region_id":     object["VpcRegionId"],
		}
		ipConfigsSli := make([]map[string]interface{}, 0)
		if ipConfigs, ok := object["IpConfigs"].([]interface{}); ok {
			for _, ipConfigArgs := range ipConfigs {
				ipConfigArg := ipConfigArgs.(map[string]interface{})
				ipConfigsMap := make(map[string]interface{})
				ipConfigsMap["zone_id"] = ipConfigArg["AzId"]
				ipConfigsMap["cidr_block"] = ipConfigArg["CidrBlock"]
				ipConfigsMap["ip"] = ipConfigArg["Ip"]
				ipConfigsMap["vswitch_id"] = ipConfigArg["VSwitchId"]
				ipConfigsSli = append(ipConfigsSli, ipConfigsMap)
			}
		}
		mapping["ip_configs"] = ipConfigsSli

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, fmt.Sprint(mapping["endpoint_name"]))
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("endpoints", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
