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

func dataSourceAlicloudThreatDetectionVendors() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudThreatDetectionVendorsRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"vendor_type": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"lang": {
				Optional: true,
				Default:  "en",
				Type:     schema.TypeString,
			},
			"role_for": {
				Optional: true,
				Type:     schema.TypeInt,
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
			"vendors": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"vendor_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"vendor_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"vendor_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"create_time": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"update_time": {
							Computed: true,
							Type:     schema.TypeInt,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudThreatDetectionVendorsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("role_for"); ok {
		request["RoleFor"] = v
	}
	if v, ok := d.GetOk("vendor_type"); ok {
		request["VendorType"] = v
	}

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
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

	var objects []interface{}
	var response map[string]interface{}
	nextToken := ""
	maxResults := 50

	for {
		request["MaxResults"] = maxResults
		if nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			delete(request, "NextToken")
		}

		action := "ListVendors"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := client.RpcPost("cloud-siem", "2024-12-12", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_threat_detection_vendors", action, AlibabaCloudSdkGoERROR)
		}

		vendorsRaw, err := jsonpath.Get("$.Vendors", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Vendors", response)
		}
		result, ok := vendorsRaw.([]interface{})
		if !ok {
			result = []interface{}{}
		}

		if isPagingRequest(d) {
			objects = result
			break
		}

		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["VendorId"])]; !ok {
					continue
				}
			}
			if v, ok := d.GetOk("vendor_type"); ok && fmt.Sprint(item["VendorType"]) != v.(string) {
				continue
			}
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["VendorName"])) {
				continue
			}
			objects = append(objects, item)
		}

		nt, _ := jsonpath.Get("$.NextToken", response)
		nextToken = fmt.Sprint(nt)
		if nextToken == "" || nextToken == "<nil>" {
			break
		}
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":          fmt.Sprint(object["VendorId"]),
			"vendor_id":   object["VendorId"],
			"vendor_name": object["VendorName"],
			"vendor_type": object["VendorType"],
			"create_time": object["CreateTime"],
			"update_time": object["UpdateTime"],
		}
		ids = append(ids, fmt.Sprint(object["VendorId"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("vendors", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}
	return nil
}
