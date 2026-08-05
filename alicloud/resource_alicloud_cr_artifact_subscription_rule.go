// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCrArtifactSubscriptionRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCrArtifactSubscriptionRuleCreate,
		Read:   resourceAliCloudCrArtifactSubscriptionRuleRead,
		Update: resourceAliCloudCrArtifactSubscriptionRuleUpdate,
		Delete: resourceAliCloudCrArtifactSubscriptionRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"accelerate": {
				Type:     schema.TypeBool,
				Optional: true,
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
				Required: true,
				ForceNew: true,
			},
			"modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespace_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"override": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"platform": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"repo_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_namespace_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_provider": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_repo_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tag_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"tag_regexp": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAliCloudCrArtifactSubscriptionRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateArtifactSubscriptionRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	request["RegionId"] = client.RegionId

	request["SourceProvider"] = d.Get("source_provider")
	if v, ok := d.GetOkExists("override"); ok {
		request["Override"] = v
	}
	request["TagRegexp"] = d.Get("tag_regexp")
	request["SourceRepoName"] = d.Get("source_repo_name")
	if v, ok := d.GetOkExists("accelerate"); ok {
		request["Accelerate"] = v
	}
	request["NamespaceName"] = d.Get("namespace_name")
	request["RepoName"] = d.Get("repo_name")
	request["TagCount"] = d.Get("tag_count")
	if v, ok := d.GetOk("source_namespace_name"); ok {
		request["SourceNamespaceName"] = v
	}
	if v, ok := d.GetOk("platform"); ok {
		request["Platform"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_artifact_subscription_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["InstanceId"], response["RuleId"]))

	return resourceAliCloudCrArtifactSubscriptionRuleUpdate(d, meta)
}

func resourceAliCloudCrArtifactSubscriptionRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crServiceV2 := CrServiceV2{client}

	objectRaw, err := crServiceV2.DescribeCrArtifactSubscriptionRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cr_artifact_subscription_rule DescribeCrArtifactSubscriptionRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if objectRaw["RuleId"] == nil || objectRaw["RuleId"] == "" {
		if !d.IsNewResource() {
			log.Printf("[DEBUG] Resource alicloud_cr_artifact_subscription_rule not found with id %s", d.Id())
			d.SetId("")
			return nil
		}
	}
	d.Set("accelerate", objectRaw["Accelerate"])
	d.Set("create_time", fmt.Sprintf("%v", objectRaw["CreateTime"]))
	d.Set("modified_time", fmt.Sprintf("%v", objectRaw["ModifiedTime"]))
	d.Set("namespace_name", objectRaw["NamespaceName"])
	d.Set("override", objectRaw["Override"])
	d.Set("platform", objectRaw["Platform"])
	d.Set("repo_name", objectRaw["RepoName"])
	d.Set("source_domain", objectRaw["SourceDomain"])
	d.Set("source_namespace_name", objectRaw["SourceNamespaceName"])
	d.Set("source_provider", objectRaw["SourceProvider"])
	d.Set("source_repo_name", objectRaw["SourceRepoName"])
	d.Set("tag_count", objectRaw["TagCount"])
	d.Set("tag_regexp", objectRaw["TagRegexp"])
	d.Set("artifact_subscription_rule_id", objectRaw["RuleId"])
	d.Set("instance_id", objectRaw["InstanceId"])

	return nil
}

func resourceAliCloudCrArtifactSubscriptionRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateArtifactSubscriptionRule"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RuleId"] = parts[1]
	request["InstanceId"] = parts[0]
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("source_provider") {
		update = true
	}
	request["SourceProvider"] = d.Get("source_provider")
	if !d.IsNewResource() && d.HasChange("override") {
		update = true
		request["Override"] = d.Get("override")
	}

	if !d.IsNewResource() && d.HasChange("tag_regexp") {
		update = true
	}
	request["TagRegexp"] = d.Get("tag_regexp")
	if !d.IsNewResource() && d.HasChange("source_repo_name") {
		update = true
	}
	request["SourceRepoName"] = d.Get("source_repo_name")
	if !d.IsNewResource() && d.HasChange("accelerate") {
		update = true
		request["Accelerate"] = d.Get("accelerate")
	}

	if !d.IsNewResource() && d.HasChange("namespace_name") {
		update = true
	}
	request["NamespaceName"] = d.Get("namespace_name")
	if !d.IsNewResource() && d.HasChange("repo_name") {
		update = true
	}
	request["RepoName"] = d.Get("repo_name")
	if !d.IsNewResource() && d.HasChange("tag_count") {
		update = true
	}
	request["TagCount"] = d.Get("tag_count")
	if !d.IsNewResource() && d.HasChange("source_namespace_name") {
		update = true
		request["SourceNamespaceName"] = d.Get("source_namespace_name")
	}
	request["Platform"] = d.Get("platform")
	if !d.IsNewResource() && d.HasChange("platform") {
		update = true
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("cr", "2018-12-01", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudCrArtifactSubscriptionRuleRead(d, meta)
}

func resourceAliCloudCrArtifactSubscriptionRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteArtifactSubscriptionRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = parts[0]
	request["RuleId"] = parts[1]
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
