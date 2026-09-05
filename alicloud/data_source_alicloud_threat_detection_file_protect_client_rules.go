// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"
)

func dataSourceAliCloudThreatDetectionFileProtectClientRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudThreatDetectionFileProtectClientRuleRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"alert_level": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"platform": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule_action": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alert_level": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"exclude_users": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"file_ops": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"file_paths": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"file_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"platform": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"proc_paths": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"rule_action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"switch_id": {
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

func dataSourceAliCloudThreatDetectionFileProtectClientRuleRead(d *schema.ResourceData, meta interface{}) error {
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
	var query map[string]interface{}
	action := "ListFileProtectClientRule"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	if v, ok := d.GetOkExists("alert_level"); ok {
		request["AlertLevel"] = v
	}
	if v, ok := d.GetOk("platform"); ok {
		request["Platform"] = v
	}
	if v, ok := d.GetOk("rule_action"); ok {
		request["RuleAction"] = v
	}
	if v, ok := d.GetOk("rule_name"); ok {
		request["RuleName"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["CurrentPage"] = 1
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)

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

		resp, _ := jsonpath.Get("$.FileProtectList[*]", response)

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

		if len(result) < PageSizeLarge {
			break
		}
		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["Id"]

		mapping["alert_level"] = objectRaw["AlertLevel"]
		mapping["platform"] = objectRaw["Platform"]
		mapping["rule_action"] = objectRaw["RuleAction"]
		mapping["rule_name"] = objectRaw["RuleName"]
		mapping["status"] = objectRaw["Status"]
		mapping["switch_id"] = objectRaw["SwitchId"]

		excludeUsersRaw := make([]interface{}, 0)
		if objectRaw["ExcludeUsers"] != nil {
			excludeUsersRaw = convertToInterfaceArray(objectRaw["ExcludeUsers"])
		}

		mapping["exclude_users"] = excludeUsersRaw
		fileOpsRaw := make([]interface{}, 0)
		if objectRaw["FileOps"] != nil {
			fileOpsRaw = convertToInterfaceArray(objectRaw["FileOps"])
		}

		mapping["file_ops"] = fileOpsRaw
		filePathsRaw := make([]interface{}, 0)
		if objectRaw["FilePaths"] != nil {
			filePathsRaw = convertToInterfaceArray(objectRaw["FilePaths"])
		}

		mapping["file_paths"] = filePathsRaw
		fileTypesRaw := make([]interface{}, 0)
		if objectRaw["FileTypes"] != nil {
			fileTypesRaw = convertToInterfaceArray(objectRaw["FileTypes"])
		}

		mapping["file_types"] = fileTypesRaw
		procPathsRaw := make([]interface{}, 0)
		if objectRaw["ProcPaths"] != nil {
			procPathsRaw = convertToInterfaceArray(objectRaw["ProcPaths"])
		}

		mapping["proc_paths"] = procPathsRaw

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
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
