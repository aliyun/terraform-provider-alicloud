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

func dataSourceAlicloudThreatDetectionCustomCheckStandardPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudThreatDetectionCustomCheckStandardPoliciesRead,
		Schema: map[string]*schema.Schema{
			"dependent_policy_id": {
				Optional: true,
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
			"name_regex": {
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"names": {
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
			"policies": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"check_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"dependent_policy_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"policy_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"policy_show_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"policy_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"type": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"policy_type": {
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"STANDARD", "REQUIREMENT", "SECTION"}, false),
			},
			"type": {
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"AISPM", "IDENTITY_PERMISSION", "RISK", "COMPLIANCE"}, false),
			},
		},
	}
}

func dataSourceAlicloudThreatDetectionCustomCheckStandardPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"PageSize":    50,
		"CurrentPage": 1,
	}
	if v, ok := d.GetOk("policy_type"); ok {
		request["PolicyType"] = v
	}
	if v, ok := d.GetOk("dependent_policy_id"); ok && v.(string) != "" {
		request["DependentPolicyId"] = v
	}
	if v, ok := d.GetOk("type"); ok && v.(string) != "" {
		request["Type"] = v
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

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
	}

	var objects []map[string]interface{}
	var response map[string]interface{}
	action := "ListCheckPolicies"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_threat_detection_custom_check_standard_policies", action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Policies", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Policies", response)
	}
	policies, _ := resp.([]interface{})
	for _, policy := range policies {
		p := policy.(map[string]interface{})
		policyId := fmt.Sprint(p["PolicyId"])
		compositeId := fmt.Sprintf("%s:%s", policyId, p["PolicyType"])
		if len(idsMap) > 0 {
			if _, ok := idsMap[compositeId]; !ok {
				if _, ok2 := idsMap[policyId]; !ok2 {
					continue
				}
			}
		}
		if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(p["PolicyShowName"])) {
			continue
		}
		objects = append(objects, p)
	}

	ids := make([]string, 0, len(objects))
	names := make([]interface{}, 0, len(objects))
	s := make([]map[string]interface{}, 0, len(objects))
	for _, object := range objects {
		policyId := fmt.Sprint(object["PolicyId"])
		mapping := map[string]interface{}{
			"id":               fmt.Sprintf("%s:%s", policyId, object["PolicyType"]),
			"policy_id":        policyId,
			"policy_show_name": object["PolicyShowName"],
			"policy_type":      object["PolicyType"],
			"check_type":       object["CheckType"],
			"type":             object["Type"],
		}
		if object["DependentPolicyId"] != nil {
			mapping["dependent_policy_id"] = fmt.Sprint(object["DependentPolicyId"])
		}
		ids = append(ids, mapping["id"].(string))
		names = append(names, object["PolicyShowName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("policies", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
