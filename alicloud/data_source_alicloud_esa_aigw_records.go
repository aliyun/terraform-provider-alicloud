// Package alicloud
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

func dataSourceAliCloudEsaAigwRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudEsaAigwRecordsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"record_name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"fuzzy_search_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"record_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"site_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudEsaAigwRecordsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	instanceId := d.Get("instance_id").(string)
	action := "ListAIGWRecords"
	request := make(map[string]interface{})
	request["InstanceId"] = instanceId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	if v, ok := d.GetOk("fuzzy_search_key"); ok {
		request["FuzzySearchKey"] = v
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

	var objects []map[string]interface{}
	var response map[string]interface{}
	var err error

	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcGet("ESA", "2024-09-10", action, request, nil)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_esa_aigw_records", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Items", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Items", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			recordName := fmt.Sprint(item["RecordName"])
			compositeId := fmt.Sprintf("%s:%s", instanceId, recordName)
			if len(idsMap) > 0 {
				if _, ok := idsMap[compositeId]; !ok {
					if _, ok2 := idsMap[recordName]; !ok2 {
						continue
					}
				}
			}
			objects = append(objects, item)
		}

		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		recordName := fmt.Sprint(object["RecordName"])
		compositeId := fmt.Sprintf("%s:%s", instanceId, recordName)
		mapping := map[string]interface{}{
			"id":          compositeId,
			"instance_id": instanceId,
			"record_name": recordName,
			"site_id":     fmt.Sprint(object["SiteId"]),
			"create_time": fmt.Sprint(object["CreateTime"]),
		}
		ids = append(ids, compositeId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("records", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
