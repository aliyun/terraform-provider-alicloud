package tsdb

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceTsdbInstances() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceTsdbInstancesRead,
		Schema: map[string]*schema.Schema{
			"app_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"engine_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"tsdb_influxdb", "tsdb_tsdb"}, false),
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"query_str": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ACTIVATION", "CLASS_CHANGING", "CREATING", "DELETED", "LOCKED"}, false),
			},
			"status_list": {
				Type:     schema.TypeString,
				Optional: true,
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
						"auto_renew": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expired_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_class": {
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
						"instance_storage": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
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
						"vpc_connection_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
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
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func DataSourceTsdbInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeHiTSDBInstanceList"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("app_key"); ok {
		request["AppKey"] = v
	}
	if v, ok := d.GetOk("engine_type"); ok {
		request["EngineType"] = v
	}
	if v, ok := d.GetOk("query_str"); ok {
		request["QueryStr"] = v
	}
	if v, ok := d.GetOk("status_list"); ok {
		request["StatusList"] = v
	}
	request["RegionId"] = client.RegionId
	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1
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
	status, statusOk := d.GetOk("status")
	var response map[string]interface{}
	var err error
	for {
		response, err = client.RpcPost("hitsdb", "2017-06-01", action, nil, request, true)
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_tsdb_instances", action, errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(action, response, request)

		resp, err := jsonpath.Get("$.InstanceList", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.InstanceList", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["InstanceId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
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
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"engine_type":      object["EngineType"],
			"expired_time":     object["ExpiredTime"],
			"instance_alias":   object["InstanceAlias"],
			"instance_class":   object["InstanceClass"],
			"id":               fmt.Sprint(object["InstanceId"]),
			"instance_id":      fmt.Sprint(object["InstanceId"]),
			"instance_storage": object["InstanceStorage"],
			"network_type":     object["NetworkType"],
			"payment_type":     convertTsdbInstancePaymentTypeResponse(object["PaymentType"].(string)),
			"status":           object["Status"],
			"vswitch_id":       object["VswitchId"],
			"vpc_id":           object["VpcId"],
			"zone_id":          object["ZoneId"],
		}
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(object["InstanceId"]))
			s = append(s, mapping)
			continue
		}

		tsdbService := NewTsdbService(client)
		id := fmt.Sprint(object["InstanceId"])
		getResp, err := tsdbService.DescribeTsdbInstance(id)
		if err != nil {
			return errs.WrapError(err)
		}
		mapping["auto_renew"] = getResp["AutoRenew"]
		mapping["cpu_number"] = getResp["CpuNumber"]
		mapping["disk_category"] = getResp["DiskCategory"]
		mapping["memory_size"] = getResp["MemSize"]
		mapping["vpc_connection_address"] = getResp["ConnectionString"]
		ids = append(ids, fmt.Sprint(object["InstanceId"]))
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("instances", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		_ = helper.WriteToFile(output.(string), s)
	}

	return nil
}
