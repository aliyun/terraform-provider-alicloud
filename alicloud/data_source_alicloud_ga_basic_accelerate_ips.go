package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudGaBasicAccelerateIps() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudGaBasicAccelerateIpsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ip_set_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"accelerate_ip_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"accelerate_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"active", "binding", "bound", "unbinding", "deleting"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ips": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"accelerate_ip_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"accelerate_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"accelerator_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_set_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudGaBasicAccelerateIpsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListBasicAccelerateIps"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["MaxResults"] = PageSizeLarge
	request["IpSetId"] = d.Get("ip_set_id")

	if v, ok := d.GetOk("accelerate_ip_id"); ok {
		request["AccelerateIpId"] = v
	}

	if v, ok := d.GetOk("accelerate_ip_address"); ok {
		request["AccelerateIpAddress"] = v
	}

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

	var response map[string]interface{}
	var err error

	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ga_basic_accelerate_ips", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.AccelerateIps", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.AccelerateIps", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["AccelerateIpId"])]; !ok {
					continue
				}
			}

			if statusOk && status.(string) != "" && status.(string) != item["State"].(string) {
				continue
			}

			objects = append(objects, item)
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                    fmt.Sprint(object["AccelerateIpId"]),
			"accelerate_ip_id":      fmt.Sprint(object["AccelerateIpId"]),
			"accelerate_ip_address": object["AccelerateIpAddress"],
			"accelerator_id":        object["AcceleratorId"],
			"ip_set_id":             object["IpSetId"],
			"status":                object["State"],
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("ips", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
