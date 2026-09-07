// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudCmsAlertWebhooks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCmsAlertWebhooksRead,
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
			"alert_webhook_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"webhooks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alert_webhook_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alert_webhook_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"content_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"headers": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"lang": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"method": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workspace": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudCmsAlertWebhooksRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "/webhooks"
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

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
	}

	query := make(map[string]*string)
	query["pageNumber"] = StringPointer("1")
	query["pageSize"] = StringPointer(fmt.Sprintf("%d", PageSizeLarge))
	if v, ok := d.GetOk("alert_webhook_name"); ok {
		query["name"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("workspace"); ok {
		query["workspace"] = StringPointer(v.(string))
	}

	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, query)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cms_alert_webhooks", action, AlibabaCloudSdkGoERROR)
		}

		resp, gerr := jsonpath.Get("$.webhooks[*]", response)
		if gerr != nil {
			if response["webhooks"] == nil {
				break
			}
			return WrapErrorf(gerr, FailedGetAttributeMsg, action, "$.webhooks[*]", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			webhookName := fmt.Sprint(item["name"])
			if nameRegex != nil && !nameRegex.MatchString(webhookName) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["webhookId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		total, _ := jsonpath.Get("$.total", response)
		totalCount := 0
		if total != nil {
			totalCount, _ = strconv.Atoi(fmt.Sprint(total))
		}
		currentPage, _ := strconv.Atoi(fmt.Sprint(query["pageNumber"]))
		pageSize := PageSizeLarge
		if len(result) < pageSize || currentPage*pageSize >= totalCount {
			break
		}
		query["pageNumber"] = StringPointer(fmt.Sprintf("%d", currentPage+1))
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		webhookId := fmt.Sprint(object["webhookId"])
		mapping := map[string]interface{}{
			"id":                 webhookId,
			"alert_webhook_id":   webhookId,
			"alert_webhook_name": object["name"],
			"content_type":       object["contentType"],
			"headers":            object["headers"],
			"lang":               object["lang"],
			"method":             object["method"],
			"url":                object["url"],
			"workspace":          object["workspace"],
		}

		ids = append(ids, webhookId)
		names = append(names, object["name"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("webhooks", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
