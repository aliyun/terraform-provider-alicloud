// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/blues/jsonata-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudQuotasTemplateApplications() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudQuotasTemplateApplicationsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"batch_quota_application_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"product_code": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"quota_action_code": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"quota_category": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"applications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"aliyun_uids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"apply_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"audit_status_vos": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"batch_quota_application_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"desire_value": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"dimensions": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"effective_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expire_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota_action_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func dataSourceAliCloudQuotasTemplateApplicationsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var request map[string]interface{}
	var response map[string]interface{}
	action := "ListQuotaApplicationsForTemplate"
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("batch_quota_application_id"); ok {
		request["BatchQuotaApplicationId"] = v
	}
	if v, ok := d.GetOk("product_code"); ok {
		request["ProductCode"] = v
	}
	if v, ok := d.GetOk("quota_action_code"); ok {
		request["QuotaActionCode"] = v
	}
	if v, ok := d.GetOk("quota_category"); ok {
		request["QuotaCategory"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.Do("quotas", rpc("POST", "2020-05-10", action), nil, request, nil, nil, false)

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

	var objects []map[string]interface{}
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

	resp, _ := jsonpath.Get("$.QuotaBatchApplications[*]", response)

	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["Name"])) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["BatchQuotaApplicationId"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{
			"id":                objectRaw["BatchQuotaApplicationId"],
			"desire_value":      objectRaw["DesireValue"],
			"effective_time":    objectRaw["EffectiveTime"],
			"expire_time":       objectRaw["ExpireTime"],
			"product_code":      objectRaw["ProductCode"],
			"quota_action_code": objectRaw["QuotaActionCode"],
			"quota_category":    objectRaw["QuotaCategory"],
			"reason":            objectRaw["Reason"],
		}

		aliyunUids2Raw := make([]interface{}, 0)
		if objectRaw["AliyunUids"] != nil {
			aliyunUids2Raw = objectRaw["AliyunUids"].([]interface{})
		}
		mapping["aliyun_uids"] = aliyunUids2Raw
		e := jsonata.MustCompile("$each($.Dimensions, function($v, $k) {{\"value\":$v, \"key\": $k}})[]")
		evaluation, _ := e.Eval(objectRaw)
		mapping["dimensions"] = evaluation
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["AlertName"])
		s = append(s, mapping)
	}
	d.Set("applications", s)
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
