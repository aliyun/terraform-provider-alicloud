package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudGaDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudGaDomainsRead,
		Schema: map[string]*schema.Schema{
			"accelerator_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"domain": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"status": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"illegal", "inactive", "active", "unknown"}, false),
			},
			"ids": {
				Optional: true,
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
			},
			"domains": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"accelerator_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"domain": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudGaDomainsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	if v, ok := d.GetOk("accelerator_id"); ok {
		request["AcceleratorId"] = v
	}
	if v, ok := d.GetOk("domain"); ok {
		request["Domain"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["State"] = v
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

	var err error
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "ListDomains"
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
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ga_domains", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Domains", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Domains", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}
		for _, v := range result {

			acceleratorsItem := v.(map[string]interface{})

			if _, ok := acceleratorsItem["Accelerators"]; !ok {
				continue
			}

			resp, err := jsonpath.Get("$.Accelerators", acceleratorsItem)
			if err != nil {
				return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Accelerators", response)
			}

			acceleratorsResult, _ := resp.([]interface{})

			for _, v := range acceleratorsResult {
				item := v.(map[string]interface{})
				item["Domain"] = acceleratorsItem["Domain"]
				item["State"] = acceleratorsItem["State"]
				if len(idsMap) > 0 {
					if _, ok := idsMap[fmt.Sprint(item["AcceleratorId"], ":", item["Domain"])]; !ok {
						continue
					}
				}
				objects = append(objects, item)
			}
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":             fmt.Sprint(object["AcceleratorId"], ":", object["Domain"]),
			"accelerator_id": object["AcceleratorId"],
			"domain":         object["Domain"],
			"status":         object["State"],
		}

		ids = append(ids, fmt.Sprint(object["AcceleratorId"], ":", object["Domain"]))

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("domains", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
