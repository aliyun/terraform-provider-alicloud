// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudCrArtifactLifecycleRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCrArtifactLifecycleRuleRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"artifact_lifecycle_rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"enable_delete_tag": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable_delete_untagged_manifest": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modified_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"namespace_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repo_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"retention_tag_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"schedule_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scope": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tag_regexp": {
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
		},
	}
}

func dataSourceAliCloudCrArtifactLifecycleRuleRead(d *schema.ResourceData, meta interface{}) error {
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

	var response map[string]interface{}
	var query map[string]interface{}
	action := "ListArtifactLifecycleRule"
	var err error
	query = make(map[string]interface{})
	query["RegionId"] = client.RegionId
	if v, ok := d.GetOk("instance_id"); ok {
		query["InstanceId"] = v.(string)
	}
	query["EnableDeleteTag"] = "true"
	query["EnableDeleteUntaggedManifest"] = "false"

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	query["PageSize"] = PageSizeLarge
	query["PageNo"] = 1
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RpcGet("cr", "2018-12-01", action, query, nil)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, query)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.Rules[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["InstanceId"], ":", item["RuleId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < PageSizeLarge {
			break
		}
		query["PageNo"] = query["PageNo"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = fmt.Sprint(objectRaw["InstanceId"], ":", objectRaw["RuleId"])

		mapping["auto"] = objectRaw["Auto"]
		mapping["create_time"] = objectRaw["CreateTime"]
		mapping["modified_time"] = objectRaw["ModifiedTime"]
		mapping["namespace_name"] = objectRaw["NamespaceName"]
		mapping["repo_name"] = objectRaw["RepoName"]
		mapping["retention_tag_count"] = objectRaw["RetentionTagCount"]
		mapping["schedule_time"] = objectRaw["ScheduleTime"]
		mapping["scope"] = objectRaw["Scope"]
		mapping["tag_regexp"] = objectRaw["TagRegexp"]
		mapping["artifact_lifecycle_rule_id"] = objectRaw["RuleId"]
		mapping["enable_delete_tag"] = objectRaw["EnableDeleteTag"]
		mapping["enable_delete_untagged_manifest"] = objectRaw["EnableDeleteUntaggedManifest"]
		mapping["instance_id"] = objectRaw["InstanceId"]

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
