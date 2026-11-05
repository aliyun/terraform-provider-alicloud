// Package ens. This file is generated automatically. Please do not modify it manually, thank you!
package ens

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceEnsBucketLifecycles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEnsBucketLifecyclesRead,
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bucket_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expiration_days": {
							Type:     schema.TypeInt,
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

func dataSourceEnsBucketLifecyclesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "GetBucketLifecycle"
	request := make(map[string]interface{})
	query := make(map[string]interface{})
	bucketName := d.Get("bucket_name").(string)
	query["BucketName"] = bucketName
	if v, ok := d.GetOk("rule_id"); ok {
		query["RuleId"] = v
	}

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.NotFoundError(err) || errs.IsExpectedErrors(err, []string{"NoSuchBucket", "NoSuchLifecycle"}) {
			d.SetId(helper.DataResourceIdHash([]string{bucketName}))
			d.Set("ids", []string{})
			d.Set("rules", []interface{}{})
			return nil
		}
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_ens_bucket_lifecycles", action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Rule[*]", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Rule[*]", response)
	}
	result, _ := resp.([]interface{})

	ids := make([]string, 0)
	rules := make([]map[string]interface{}, 0)
	for _, v := range result {
		item := v.(map[string]interface{})
		ruleId := fmt.Sprint(item["ID"])
		rule := map[string]interface{}{
			"id":          fmt.Sprintf("%s:%s", bucketName, ruleId),
			"bucket_name": bucketName,
			"rule_id":     ruleId,
			"prefix":      item["Prefix"],
			"status":      item["Status"],
		}
		if expiration, ok := item["Expiration"].(map[string]interface{}); ok {
			if days, ok := expiration["Days"]; ok && days != nil {
				if daysStr := fmt.Sprint(days); daysStr != "" && daysStr != "<nil>" {
					if n, e := strconv.Atoi(daysStr); e == nil {
						rule["expiration_days"] = n
					}
				}
			}
		}
		ids = append(ids, fmt.Sprint(rule["id"]))
		rules = append(rules, rule)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("rules", rules); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), rules)
	}

	return nil
}
