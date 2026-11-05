package lindorm

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

func DataSourceLindormInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLindormInstancesRead,
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
			"query_str": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ACTIVATION", "CLASS_CHANGING", "CREATING", "DELETED", "INSTANCE_LEVEL_MODIFY", "LOCKED", "MINOR_VERSION_TRANSING", "NET_MODIFYING", "RESIZING", "RESTARTING"}, false),
			},
			"support_engine": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 7),
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
						"cold_storage": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deletion_proection": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"disk_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_usage": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_warning_threshold": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"expired_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_engine_node_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"file_engine_specification": {
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
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_storage": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_white_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"lts_node_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"lts_node_specification": {
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
						"phoenix_node_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"phoenix_node_specification": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_owner_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"search_engine_node_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"search_engine_specification": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"table_engine_node_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"table_engine_specification": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"time_series_engine_node_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"time_serires_engine_specification": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
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

func dataSourceLindormInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "GetLindormInstanceList"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("query_str"); ok {
		request["QueryStr"] = v
	}
	if v, ok := d.GetOk("support_engine"); ok {
		request["SupportEngine"] = v
	}
	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var instanceNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
		}
		instanceNameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, val := range v.([]interface{}) {
			if val == nil {
				continue
			}
			idsMap[val.(string)] = val.(string)
		}
	}
	status, statusOk := d.GetOk("status")
	var response map[string]interface{}
	var err error
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("hitsdb", "2020-06-15", action, nil, request, true)
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
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_lindorm_instances", action, errs.AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.InstanceList", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.InstanceList", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if instanceNameRegex != nil && !instanceNameRegex.MatchString(fmt.Sprint(item["InstanceAlias"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["InstanceId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["InstanceStatus"].(string) {
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
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {

		engine_typeVal, err := convert.FormatInt(object["EngineType"])
		if err != nil {
			return err
		}
		mapping := map[string]interface{}{
			"create_time":      fmt.Sprint(object["CreateMilliseconds"]),
			"engine_type":      engine_typeVal,
			"expired_time":     fmt.Sprint(object["ExpiredMilliseconds"]),
			"id":               fmt.Sprint(object["InstanceId"]),
			"instance_id":      fmt.Sprint(object["InstanceId"]),
			"instance_name":    object["InstanceAlias"],
			"instance_storage": object["InstanceStorage"],
			"network_type":     object["NetworkType"],
			"payment_type":     convertLindormInstancePaymentTypeResponse(object["PayType"]),
			"service_type":     object["ServiceType"],
			"status":           object["InstanceStatus"],
			"vpc_id":           object["VpcId"],
			"zone_id":          object["ZoneId"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["InstanceAlias"])
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["InstanceId"])
		lindormService := NewService(client)
		ipWhiteList, err := lindormService.GetInstanceIpWhiteList(id)
		if err != nil {
			return errs.WrapError(err)
		}
		mapping["ip_white_list"] = convert.ConvertListToJsonString(ipWhiteList)
		getResp1, err := lindormService.DescribeLindormInstance(id)
		if err != nil {
			return errs.WrapError(err)
		}
		mapping["auto_renew"] = getResp1["AutoRenew"]
		if v, ok := getResp1["ColdStorage"]; ok {
			coldStorage, err := convert.FormatInt(v)
			if err != nil {
				return err
			}
			mapping["cold_storage"] = coldStorage
		}
		mapping["deletion_proection"] = getResp1["DeletionProtection"]
		mapping["disk_category"] = getResp1["DiskCategory"]
		mapping["disk_usage"] = getResp1["DiskUsage"]
		mapping["disk_warning_threshold"] = getResp1["DiskThreshold"]
		mapping["resource_owner_id"] = getResp1["AliUid"]
		mapping["vswitch_id"] = getResp1["VswitchId"]

		getResp2, err := lindormService.GetLindormInstanceEngineInfo(id)
		if err != nil {
			return errs.WrapError(err)
		}
		if v, ok := getResp2["FileEngineNodeCount"]; ok {
			fileEngineNodeCount, err := convert.FormatInt(v)
			if err != nil {
				return err
			}
			mapping["file_engine_node_count"] = fileEngineNodeCount
		}
		mapping["file_engine_specification"] = getResp2["FileEngineSpecification"]
		if v, ok := getResp2["LtsNodeCount"]; ok {
			ltsNodeCount, err := convert.FormatInt(v)
			if err != nil {
				return err
			}
			mapping["lts_node_count"] = ltsNodeCount
		}
		mapping["lts_node_specification"] = getResp2["LtsNodeSpecification"]
		if v, ok := getResp2["PhoenixNodeCount"]; ok {
			phoenixNodeCount, err := convert.FormatInt(v)
			if err != nil {
				return err
			}
			mapping["phoenix_node_count"] = phoenixNodeCount
		}
		mapping["phoenix_node_specification"] = getResp2["PhoenixNodeSpecification"]
		if v, ok := getResp2["SearchEngineNodeCount"]; ok {
			searchEngineNodeCount, err := convert.FormatInt(v)
			if err != nil {
				return err
			}
			mapping["search_engine_node_count"] = searchEngineNodeCount
		}
		mapping["search_engine_specification"] = getResp2["SearchEngineSpecification"]
		if v, ok := getResp2["TableEngineNodeCount"]; ok {
			tableEngineNodeCount, err := convert.FormatInt(v)
			if err != nil {
				return err
			}
			mapping["table_engine_node_count"] = tableEngineNodeCount
		}
		mapping["table_engine_specification"] = getResp2["TableEngineSpecification"]
		if v, ok := getResp2["TimeSeriesNodeCount"]; ok {
			timeSeriesEngineNodeCount, err := convert.FormatInt(v)
			if err != nil {
				return err
			}
			mapping["time_series_engine_node_count"] = timeSeriesEngineNodeCount
		}
		mapping["time_serires_engine_specification"] = getResp2["TimeSeriesSpecification"]

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
