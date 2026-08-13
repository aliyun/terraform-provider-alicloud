// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAliCloudCrArtifactSubscriptionRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCrArtifactSubscriptionRuleRead,
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
			"namespace_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"repo_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accelerate": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"artifact_subscription_rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modified_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"override": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"platform": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repo_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_namespace_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_provider": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_repo_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tag_count": {
							Type:     schema.TypeInt,
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
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAliCloudCrArtifactSubscriptionRuleRead(d *schema.ResourceData, meta interface{}) error {
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
	action := "ListArtifactSubscriptionRule"
	var err error
	query = make(map[string]interface{})
	query["RegionId"] = client.RegionId
	if v, ok := d.GetOk("instance_id"); ok {
		query["InstanceId"] = v.(string)
	}

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

		mapping["accelerate"] = objectRaw["Accelerate"]
		mapping["create_time"] = fmt.Sprintf("%v", objectRaw["CreateTime"])
		mapping["modified_time"] = fmt.Sprintf("%v", objectRaw["ModifiedTime"])
		mapping["namespace_name"] = objectRaw["NamespaceName"]
		mapping["override"] = objectRaw["Override"]
		mapping["platform"] = objectRaw["Platform"]
		mapping["repo_name"] = objectRaw["RepoName"]
		mapping["source_domain"] = objectRaw["SourceDomain"]
		mapping["source_namespace_name"] = objectRaw["SourceNamespaceName"]
		mapping["source_provider"] = objectRaw["SourceProvider"]
		mapping["source_repo_name"] = objectRaw["SourceRepoName"]
		mapping["tag_count"] = objectRaw["TagCount"]
		mapping["tag_regexp"] = objectRaw["TagRegexp"]
		mapping["artifact_subscription_rule_id"] = objectRaw["RuleId"]
		mapping["instance_id"] = objectRaw["InstanceId"]

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(objectRaw["InstanceId"], ":", objectRaw["RuleId"])
		mapping, err = dataSourceAliCloudCrArtifactSubscriptionRuleReadDescription(d, id, mapping, meta)
		if err != nil {
			return WrapError(err)
		}

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

func dataSourceAliCloudCrArtifactSubscriptionRuleReadDescription(d *schema.ResourceData, id string, object map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	crServiceV2 := CrServiceV2{client}
	getResp, err := crServiceV2.DescribeCrArtifactSubscriptionRule(id)
	if err != nil {
		return nil, WrapError(err)
	}

	// Merge additional fields from Get API response to mapping
	// Reuse the response mapping template from Resource's read function
	mapping := object
	objectRaw := getResp

	mapping["accelerate"] = objectRaw["Accelerate"]
	mapping["create_time"] = fmt.Sprintf("%v", objectRaw["CreateTime"])
	mapping["modified_time"] = fmt.Sprintf("%v", objectRaw["ModifiedTime"])
	mapping["namespace_name"] = objectRaw["NamespaceName"]
	mapping["override"] = objectRaw["Override"]
	mapping["platform"] = objectRaw["Platform"]
	mapping["region_id"] = objectRaw["RegionId"]
	mapping["repo_name"] = objectRaw["RepoName"]
	mapping["source_domain"] = objectRaw["SourceDomain"]
	mapping["source_namespace_name"] = objectRaw["SourceNamespaceName"]
	mapping["source_provider"] = objectRaw["SourceProvider"]
	mapping["source_repo_name"] = objectRaw["SourceRepoName"]
	mapping["tag_count"] = objectRaw["TagCount"]
	mapping["tag_regexp"] = objectRaw["TagRegexp"]
	mapping["artifact_subscription_rule_id"] = objectRaw["RuleId"]
	mapping["instance_id"] = objectRaw["InstanceId"]

	return mapping, nil
}
