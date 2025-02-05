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

func dataSourceAlicloudThreatDetectionAntiBruteForceRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudThreatDetectionAntiBruteForceRulesRead,
		Schema: map[string]*schema.Schema{
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
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.ValidateRegexp,
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
			"rules": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"anti_brute_force_rule_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"anti_brute_force_rule_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"default_rule": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"fail_count": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"forbidden_time": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"span": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"uuid_list": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudThreatDetectionAntiBruteForceRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var antiBruteForceRuleNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		antiBruteForceRuleNameRegex = r
	}

	var err error
	var objects []interface{}
	var response map[string]interface{}
	action := "DescribeAntiBruteForceRules"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_threat_detection_anti_brute_force_rules", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.Rules", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Rules", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["Id"])]; !ok {
				continue
			}
		}

		if antiBruteForceRuleNameRegex != nil && !antiBruteForceRuleNameRegex.MatchString(fmt.Sprint(item["Name"])) {
			continue
		}
		objects = append(objects, item)
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                         fmt.Sprint(object["Id"]),
			"anti_brute_force_rule_id":   object["Id"],
			"anti_brute_force_rule_name": object["Name"],
			"default_rule":               object["DefaultRule"],
			"fail_count":                 object["FailCount"],
			"forbidden_time":             object["ForbiddenTime"],
			"span":                       object["Span"],
			"uuid_list":                  object["UuidList"].([]interface{}),
		}

		ids = append(ids, fmt.Sprint(object["Id"]))
		names = append(names, object["Name"])

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("rules", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
