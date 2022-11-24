package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"
)

func dataSourceAlicloudQuotasTemplateQuotas() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudQuotasTemplateQuotasRead,
		Schema: map[string]*schema.Schema{
			"product_code": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"quota_action_code": {
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
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"quotas": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"applicable_range": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeFloat,
							},
						},
						"applicable_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"desire_value": {
							Computed: true,
							Type:     schema.TypeFloat,
						},
						"dimensions": {
							Computed: true,
							Type:     schema.TypeMap,
						},
						"env_language": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"notice_type": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"product_code": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"quota_action_code": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"quota_description": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"quota_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"template_quota_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudQuotasTemplateQuotasRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	if v, ok := d.GetOk("product_code"); ok {
		request["ProductCode"] = v
	}
	if v, ok := d.GetOk("quota_action_code"); ok {
		request["QuotaActionCode"] = v
	}
	request["MaxResults"] = PageSizeLarge

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	conn, err := client.NewQuotasClient()
	if err != nil {
		return WrapError(err)
	}
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "ListQuotaApplicationTemplates"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-05-10"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_quotas_template_quotas", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.QuotaApplicationTemplates", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.QuotaApplicationTemplates", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["Id"])]; !ok {
					continue
				}
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
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                fmt.Sprint(object["Id"]),
			"applicable_range":  object["ApplicableRange"].([]interface{}),
			"applicable_type":   object["ApplicableType"],
			"desire_value":      object["DesireValue"],
			"dimensions":        object["Dimensions"],
			"env_language":      object["EnvLanguage"],
			"notice_type":       object["NoticeType"],
			"product_code":      object["ProductCode"],
			"quota_action_code": object["QuotaActionCode"],
			"quota_description": object["QuotaDescription"],
			"quota_name":        object["QuotaName"],
			"template_quota_id": object["Id"],
		}

		ids = append(ids, fmt.Sprint(object["Id"]))

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("quotas", s); err != nil {
		return WrapError(err)
	}
	return nil
}
