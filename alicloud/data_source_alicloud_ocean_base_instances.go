package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudOceanBaseInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudOceanBaseInstancesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"instance_name": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"resource_group_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"search_key": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"status": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"PENDING_CREATE", "ONLINE", "TENANT_CREATING", "TENANT_SPEC_MODIFYING", "EXPANDING", "REDUCING", "SPEC_UPGRADING", "DISK_UPGRADING", "WHITE_LIST_MODIFYING", "PARAMETER_MODIFYING", "SSL_MODIFYING", "PREPAID_EXPIRE_CLOSED", "ARREARS_CLOSED", "PENDING_DELETE"}, false),
			},
			"ids": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name_regex": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.ValidateRegexp,
			},
			"names": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"page_number": {
				Optional: true,
				Type:     schema.TypeInt,
			},
			"page_size": {
				Optional: true,
				Type:     schema.TypeInt,
				Default:  50,
			},
			"enable_details": {
				Optional: true,
				Type:     schema.TypeBool,
			},
			"instances": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"commodity_code": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"cpu": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"create_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"disk_size": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"instance_class": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"instance_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"instance_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"node_num": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"payment_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"resource_group_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"series": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"zones": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudOceanBaseInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	if v, ok := d.GetOk("instance_name"); ok {
		request["InstanceName"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("search_key"); ok {
		request["SearchKey"] = v
	}
	setPagingRequest(d, request, PageSizeLarge)

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
			return WrapError(err)
		}
		instanceNameRegex = r
	}
	status, statusOk := d.GetOk("status")

	var err error
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "DescribeInstances"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := client.RpcPost("OceanBasePro", "2019-09-01", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			response = resp
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ocean_base_instances", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Instances", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Instances", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}
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

			if statusOk && status.(string) != "" && status.(string) != item["State"].(string) {
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
	oceanBaseProService := OceanBaseProService{client}
	for _, v := range objects {
		object := v.(map[string]interface{})

		mapping := map[string]interface{}{
			"id": fmt.Sprint(object["InstanceId"]),
		}

		mapping["commodity_code"] = object["CommodityCode"]
		mapping["cpu"] = object["Cpu"]
		mapping["create_time"] = object["CreateTime"]
		mapping["disk_size"] = object["DiskSize"]
		if v, ok := object["InstanceClass"]; ok && fmt.Sprint(v) != "" {
			mapping["instance_class"] = fmt.Sprint(v, "B")
		}
		mapping["instance_id"] = object["InstanceId"]
		mapping["instance_name"] = object["InstanceName"]
		mapping["payment_type"] = convertOceanBaseInstancePaymentTypeResponse(object["PayType"])
		mapping["resource_group_id"] = object["ResourceGroupId"]
		mapping["series"] = convertOceanBaseInstanceSeriesResponse(object["Series"])
		mapping["status"] = object["State"]
		mapping["zones"] = object["AvailableZones"].([]interface{})

		ids = append(ids, fmt.Sprint(object["InstanceId"]))
		names = append(names, object["InstanceName"])

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["InstanceId"])
		object, err = oceanBaseProService.DescribeOceanBaseInstance(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["node_num"] = object["NodeNum"]
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
