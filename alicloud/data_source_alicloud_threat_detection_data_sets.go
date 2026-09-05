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

func dataSourceAlicloudThreatDetectionDataSets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudThreatDetectionDataSetsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"data_set_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_set_status": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"data_set_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"data_sets": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"data_set_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"data_set_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"data_set_field_key_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"data_set_file_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"data_set_description": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"data_set_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"data_set_status": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"role_for": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"lang": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"region_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"create_time": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"ip_whitelist_recognizers": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auto_recognize_status": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"recognize_scope": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"ip_whitelist_recognizer_type": {
										Computed: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudThreatDetectionDataSetsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListDataSets"
	request := make(map[string]interface{})

	if v, ok := d.GetOk("data_set_name"); ok {
		request["DataSetName"] = v
	}
	if v, ok := d.GetOk("data_set_status"); ok {
		request["DataSetStatus"] = v
	}
	if v, ok := d.GetOk("data_set_type"); ok {
		request["DataSetType"] = v
	}

	var objects []map[string]interface{}
	var dataSetNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		dataSetNameRegex = r
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

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("cloud-siem", "2024-12-12", action, nil, request, true)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_threat_detection_data_sets", action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.DataSets", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.DataSets", response)
	}

	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if dataSetNameRegex != nil && !dataSetNameRegex.MatchString(fmt.Sprint(item["DataSetName"])) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["DataSetId"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                       fmt.Sprint(object["DataSetId"]),
			"data_set_id":              fmt.Sprint(object["DataSetId"]),
			"data_set_name":            object["DataSetName"],
			"data_set_field_key_name":  object["DataSetFieldKeyName"],
			"data_set_file_name":       object["DataSetFileName"],
			"data_set_description":     object["DataSetDescription"],
			"data_set_type":            object["DataSetType"],
			"data_set_status":          formatInt(object["DataSetStatus"]),
			"role_for":                 formatInt(object["RoleFor"]),
			"lang":                     object["Lang"],
			"region_id":                object["RegionId"],
			"create_time":              formatInt(object["CreateTime"]),
			"ip_whitelist_recognizers": flattenThreatDetectionDataSetIpWhitelistRecognizersResponse(object["IpWhitelistRecognizers"]),
		}
		ids = append(ids, fmt.Sprint(object["DataSetId"]))
		names = append(names, object["DataSetName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("data_sets", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
