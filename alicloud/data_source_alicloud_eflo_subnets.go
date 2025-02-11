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

func dataSourceAlicloudEfloSubnets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEfloSubnetsRead,
		Schema: map[string]*schema.Schema{
			"resource_group_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"status": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Available", "Not Available", "Executing", "Deleting"}, false),
			},
			"subnet_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"subnet_name": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"type": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"OOB", "LB"}, false),
			},
			"vpd_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"zone_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
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
				Default:  20,
			},
			"subnets": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"cidr": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"create_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"gmt_modified": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"message": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"resource_group_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"subnet_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"subnet_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"vpd_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"zone_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEfloSubnetsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		request["SubnetId"] = v
	}
	if v, ok := d.GetOk("subnet_name"); ok {
		request["SubnetName"] = v
	}
	if v, ok := d.GetOk("type"); ok {
		request["Type"] = v
	}
	if v, ok := d.GetOk("vpd_id"); ok {
		request["VpdId"] = v
	}
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
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

	var subnetNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		subnetNameRegex = r
	}

	var err error
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "ListSubnets"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := client.RpcPost("eflo", "2022-05-30", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_eflo_subnets", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Content.Data", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Content.Data", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["VpdId"], ":", item["SubnetId"])]; !ok {
					continue
				}
			}

			if subnetNameRegex != nil && !subnetNameRegex.MatchString(fmt.Sprint(item["SubnetName"])) {
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
	for _, v := range objects {
		object := v.(map[string]interface{})

		mapping := map[string]interface{}{
			"id": fmt.Sprint(object["VpdId"], ":", object["SubnetId"]),
		}
		mapping["cidr"] = object["Cidr"]
		mapping["create_time"] = object["CreateTime"]
		mapping["gmt_modified"] = object["GmtModified"]
		mapping["message"] = object["Message"]
		mapping["resource_group_id"] = object["ResourceGroupId"]
		mapping["status"] = object["Status"]
		mapping["subnet_id"] = object["SubnetId"]
		mapping["subnet_name"] = object["SubnetName"]
		mapping["type"] = object["Type"]
		mapping["vpd_id"] = object["VpdId"]
		mapping["zone_id"] = object["ZoneId"]

		ids = append(ids, fmt.Sprint(object["VpdId"], ":", object["SubnetId"]))
		names = append(names, object["SubnetName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("subnets", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
