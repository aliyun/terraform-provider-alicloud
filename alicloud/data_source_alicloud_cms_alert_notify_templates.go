package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudCmsAlertNotifyTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCmsAlertNotifyTemplatesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"alert_notify_template_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"program_lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alert_notify_template_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alert_notify_template_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"templates": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"channel": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"title": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"content": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"program_lang": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAliCloudCmsAlertNotifyTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

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

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	action := "/alertNotifyTemplate"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	query["regionId"] = StringPointer(client.RegionId)
	if v, ok := d.GetOk("alert_notify_template_name"); ok {
		query["alertNotifyTemplateName"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("program_lang"); ok {
		query["programLang"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("type"); ok {
		query["type"] = StringPointer(v.(string))
	}

	pageNumber := 1
	pageSize := 50
	for {
		query["pageNumber"] = StringPointer(strconv.Itoa(pageNumber))
		query["pageSize"] = StringPointer(strconv.Itoa(pageSize))
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.alertNotifyTemplates", response)
		result, _ := resp.([]interface{})
		if len(result) == 0 {
			break
		}
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			if len(idsMap) > 0 {
				id := fmt.Sprint(item["alertNotifyTemplateId"])
				if _, ok := idsMap[id]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		total := 0
		if t, ok := response["total"]; ok && t != nil {
			total, _ = strconv.Atoi(fmt.Sprint(t))
		}
		if len(objects) >= total || len(result) < pageSize {
			break
		}
		pageNumber++
	}

	return cmsAlertNotifyTemplateDescription(d, objects)
}

func cmsAlertNotifyTemplateDescription(d *schema.ResourceData, objects []map[string]interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, object := range objects {
		mapping := map[string]interface{}{
			"alert_notify_template_id":   fmt.Sprint(object["alertNotifyTemplateId"]),
			"alert_notify_template_name": object["alertNotifyTemplateName"],
			"templates":                  flattenAlertNotifyTemplates(object["templates"]),
			"type":                       object["type"],
			"program_lang":               object["programLang"],
		}
		ids = append(ids, fmt.Sprint(object["alertNotifyTemplateId"]))
		s = append(s, mapping)
	}

	if err := d.Set("templates", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return WrapError(err)
		}
	}

	if len(ids) > 0 {
		d.SetId(ids[0])
	} else {
		d.SetId("alicloud_cms_alert_notify_templates")
	}

	return nil
}
