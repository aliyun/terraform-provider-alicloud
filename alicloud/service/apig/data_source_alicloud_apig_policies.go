// Package apig. This file is generated automatically. Please do not modify it manually, thank you!
package apig

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourcePolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePolicyRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"attach_resource_ids": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attach_resource_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"environment_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attach_resource_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"attach_resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"environment_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_attachment_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_class_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_class_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_config": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
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
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourcePolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
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

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	// ListPolicies
	action := fmt.Sprintf("/v1/policies")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)

	query["withAttachments"] = helper.StringPointer("true")

	if v, ok := d.GetOk("attach_resource_ids"); ok {
		query["attachResourceId"] = helper.StringPointer(v.(string))
	}

	if v, ok := d.GetOk("attach_resource_type"); ok {
		query["attachResourceType"] = helper.StringPointer(v.(string))
	}

	if v, ok := d.GetOk("environment_id"); ok {
		query["environmentId"] = helper.StringPointer(v.(string))
	}

	if v, ok := d.GetOk("gateway_id"); ok {
		query["gatewayId"] = helper.StringPointer(v.(string))
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, _ := jsonpath.Get("$.data.items[*]", response)

	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if nameRegex != nil {
			nameStr, _ := item["name"].(string)
			if nameStr == "" {
				apigServiceV2 := NewService(client)
				detail, dErr := apigServiceV2.DescribeApigPolicy(fmt.Sprint(item["policyId"]))
				if dErr == nil && detail != nil {
					if n, ok := detail["name"].(string); ok {
						item["name"] = n
						nameStr = n
					}
				}
			}
			if !nameRegex.MatchString(nameStr) {
				continue
			}
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["policyId"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["policyId"]

		mapping["policy_class_name"] = objectRaw["className"]
		mapping["policy_config"] = objectRaw["config"]
		mapping["policy_name"] = objectRaw["name"]
		mapping["policy_id"] = objectRaw["policyId"]

		attachmentsChildRawObj, _ := jsonpath.Get("$.attachments[*]", objectRaw)
		attachmentsChildRaw := make([]interface{}, 0)
		if attachmentsChildRawObj != nil {
			attachmentsChildRaw = convert.ConvertToInterfaceArray(attachmentsChildRawObj)
		}

		attachmentsRaw := make(map[string]interface{})
		if len(attachmentsChildRaw) > 0 {
			if m, ok := attachmentsChildRaw[0].(map[string]interface{}); ok {
				attachmentsRaw = m
			}
		}
		mapping["attach_resource_type"] = attachmentsRaw["attachResourceType"]
		mapping["environment_id"] = attachmentsRaw["environmentId"]
		mapping["gateway_id"] = attachmentsRaw["gatewayId"]
		mapping["policy_attachment_id"] = attachmentsRaw["policyAttachmentId"]

		attachResourceIds := make([]string, 0)
		for _, attachment := range attachmentsChildRaw {
			if m, ok := attachment.(map[string]interface{}); ok {
				for _, id := range convert.ConvertToInterfaceArray(m["attachResourceIds"]) {
					attachResourceIds = append(attachResourceIds, fmt.Sprint(id))
				}
				if len(attachResourceIds) == 0 {
					if id, ok := m["attachResourceId"].(string); ok && id != "" {
						attachResourceIds = append(attachResourceIds, id)
					}
				}
			}
		}
		mapping["attach_resource_ids"] = attachResourceIds

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			names = append(names, objectRaw["name"])
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(objectRaw["policyId"])
		mapping, err = dataSourcePolicyReadDescription(d, id, mapping, meta)
		if err != nil {
			return errs.WrapError(err)
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["name"])
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("policies", s); err != nil {
		return errs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}
	return nil
}

func dataSourcePolicyReadDescription(d *schema.ResourceData, id string, object map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	apigServiceV2 := NewService(client)
	getResp, err := apigServiceV2.DescribeApigPolicy(id)
	if err != nil {
		return nil, errs.WrapError(err)
	}

	// Merge additional fields from Get API response to mapping
	// Reuse the response mapping template from Resource's read function
	mapping := object
	objectRaw := getResp

	mapping["policy_class_id"] = objectRaw["classId"]
	mapping["policy_class_name"] = objectRaw["className"]
	mapping["policy_config"] = objectRaw["config"]
	mapping["policy_name"] = objectRaw["name"]
	mapping["policy_id"] = objectRaw["policyId"]

	return mapping, nil
}
