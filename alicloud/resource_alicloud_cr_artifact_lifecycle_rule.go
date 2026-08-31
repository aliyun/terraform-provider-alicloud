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

func resourceAliCloudCrArtifactLifecycleRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCrArtifactLifecycleRuleCreate,
		Read:   resourceAliCloudCrArtifactLifecycleRuleRead,
		Update: resourceAliCloudCrArtifactLifecycleRuleUpdate,
		Delete: resourceAliCloudCrArtifactLifecycleRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"artifact_lifecycle_rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"modified_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"namespace_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"repo_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"retention_tag_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"schedule_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag_regexp": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudCrArtifactLifecycleRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateArtifactLifecycleRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("namespace_name"); ok {
		request["NamespaceName"] = v
	}
	request["EnableDeleteUntaggedManifest"] = "false"
	if v, ok := d.GetOk("schedule_time"); ok {
		request["ScheduleTime"] = v
	}
	if v, ok := d.GetOk("repo_name"); ok {
		request["RepoName"] = v
	}
	if v, ok := d.GetOk("scope"); ok {
		request["Scope"] = v
	}
	if v, ok := d.GetOkExists("retention_tag_count"); ok {
		request["RetentionTagCount"] = v
	}
	request["EnableDeleteTag"] = "true"
	if v, ok := d.GetOk("tag_regexp"); ok {
		request["TagRegexp"] = v
	}
	request["Auto"] = d.Get("auto")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_artifact_lifecycle_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["InstanceId"], response["RuleId"]))

	return resourceAliCloudCrArtifactLifecycleRuleRead(d, meta)
}

func resourceAliCloudCrArtifactLifecycleRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crServiceV2 := CrServiceV2{client}

	objectRaw, err := crServiceV2.DescribeCrArtifactLifecycleRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cr_artifact_lifecycle_rule DescribeCrArtifactLifecycleRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("auto", objectRaw["Auto"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("modified_time", objectRaw["ModifiedTime"])
	d.Set("namespace_name", objectRaw["NamespaceName"])
	d.Set("repo_name", objectRaw["RepoName"])
	d.Set("retention_tag_count", objectRaw["RetentionTagCount"])
	d.Set("schedule_time", objectRaw["ScheduleTime"])
	d.Set("scope", objectRaw["Scope"])
	d.Set("tag_regexp", objectRaw["TagRegexp"])
	d.Set("artifact_lifecycle_rule_id", objectRaw["RuleId"])
	d.Set("instance_id", objectRaw["InstanceId"])

	return nil
}

func resourceAliCloudCrArtifactLifecycleRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateArtifactLifecycleRule"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RuleId"] = parts[1]
	request["InstanceId"] = parts[0]
	request["RegionId"] = client.RegionId
	// The UpdateArtifactLifecycleRule API uses full-replacement semantics: any
	// attribute absent from the request is cleared on the server side, so every
	// attribute must be sent even when unchanged.
	if d.HasChange("namespace_name") {
		update = true
	}
	request["NamespaceName"] = d.Get("namespace_name")

	request["EnableDeleteUntaggedManifest"] = "false"
	if d.HasChange("schedule_time") {
		update = true
	}
	request["ScheduleTime"] = d.Get("schedule_time")

	if d.HasChange("repo_name") {
		update = true
	}
	request["RepoName"] = d.Get("repo_name")

	if d.HasChange("scope") {
		update = true
	}
	request["Scope"] = d.Get("scope")
	if d.HasChange("retention_tag_count") {
		update = true
	}
	if v, ok := d.GetOkExists("retention_tag_count"); ok {
		request["RetentionTagCount"] = v
	}

	request["EnableDeleteTag"] = "true"
	if d.HasChange("tag_regexp") {
		update = true
	}
	request["TagRegexp"] = d.Get("tag_regexp")

	if d.HasChange("auto") {
		update = true
	}
	request["Auto"] = d.Get("auto")
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

	return resourceAliCloudCrArtifactLifecycleRuleRead(d, meta)
}

func resourceAliCloudCrArtifactLifecycleRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteArtifactLifecycleRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RuleId"] = parts[1]
	request["InstanceId"] = parts[0]
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
