// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCrRepoSyncRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCrRepoSyncRuleCreate,
		Read:   resourceAliCloudCrRepoSyncRuleRead,
		Update: resourceAliCloudCrRepoSyncRuleUpdate,
		Delete: resourceAliCloudCrRepoSyncRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"repo_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"repo_sync_rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sync_rule_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"sync_rule_name", "name"},
				Computed:     true,
				ForceNew:     true,
			},
			"sync_scope": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"REPO", "NAMESPACE"}, false),
			},
			"sync_trigger": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"INITIATIVE", "PASSIVE"}, false),
			},
			"tag_filter": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_namespace_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_repo_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"target_user_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sync_direction": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field `name` has been deprecated from provider version 1.240.0. New field `sync_rule_name` instead.",
				ForceNew:   true,
			},
			"rule_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Field `rule_id` has been deprecated from provider version 1.240.0. New field `repo_sync_rule_id` instead.",
			},
		},
	}
}

func resourceAliCloudCrRepoSyncRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateRepoSyncRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewAcrClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["InstanceId"] = d.Get("instance_id")
	request["RegionId"] = client.RegionId

	request["NamespaceName"] = d.Get("namespace_name")
	request["TargetRegionId"] = d.Get("target_region_id")
	request["TargetInstanceId"] = d.Get("target_instance_id")
	request["TargetNamespaceName"] = d.Get("target_namespace_name")
	request["TagFilter"] = d.Get("tag_filter")

	if v, ok := d.GetOk("sync_rule_name"); ok {
		request["SyncRuleName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["SyncRuleName"] = v
	} else {
		return WrapError(Error(`[ERROR] Field "sync_rule_name" or "name" must be set one!`))
	}

	var repoName, targetRepoName string

	if v, ok := d.GetOk("repo_name"); ok {
		repoName = v.(string)
		request["RepoName"] = repoName
	}

	if v, ok := d.GetOk("target_repo_name"); ok {
		targetRepoName = v.(string)
		request["TargetRepoName"] = targetRepoName
	}

	if (repoName != "" && targetRepoName == "") || (repoName == "" && targetRepoName != "") {
		return WrapError(Error(`[ERROR] Field "repo_name" and "target_repo_name" must be set at the same time!`))
	}

	if v, ok := d.GetOk("sync_scope"); ok {
		request["SyncScope"] = v
	} else {
		if repoName != "" && targetRepoName != "" {
			request["SyncScope"] = "REPO"
		} else {
			request["SyncScope"] = "NAMESPACE"
		}
	}

	if v, ok := d.GetOk("sync_trigger"); ok {
		request["SyncTrigger"] = v
	} else {
		request["SyncTrigger"] = "PASSIVE"

	}

	if v, ok := d.GetOk("target_user_id"); ok {
		request["TargetUserId"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-12-01"), StringPointer("AK"), query, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_ee_sync_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s:%s:%s", request["InstanceId"], request["NamespaceName"], response["SyncRuleId"]))

	return resourceAliCloudCrRepoSyncRuleRead(d, meta)
}

func resourceAliCloudCrRepoSyncRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crServiceV2 := CrServiceV2{client}

	objectRaw, err := crServiceV2.DescribeCrRepoSyncRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cr_ee_sync_rule DescribeCrRepoSyncRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}
	if objectRaw["LocalNamespaceName"] != nil {
		d.Set("namespace_name", objectRaw["LocalNamespaceName"])
	}
	if objectRaw["LocalRegionId"] != nil {
		d.Set("region_id", objectRaw["LocalRegionId"])
	}
	if objectRaw["LocalRepoName"] != nil {
		d.Set("repo_name", objectRaw["LocalRepoName"])
	}
	if objectRaw["SyncRuleName"] != nil {
		d.Set("sync_rule_name", objectRaw["SyncRuleName"])
		d.Set("name", objectRaw["SyncRuleName"])
	}
	if objectRaw["SyncScope"] != nil {
		d.Set("sync_scope", objectRaw["SyncScope"])
	}
	if objectRaw["SyncTrigger"] != nil {
		d.Set("sync_trigger", objectRaw["SyncTrigger"])
	}
	if objectRaw["TagFilter"] != nil {
		d.Set("tag_filter", objectRaw["TagFilter"])
	}
	if objectRaw["TargetInstanceId"] != nil {
		d.Set("target_instance_id", objectRaw["TargetInstanceId"])
	}
	if objectRaw["TargetNamespaceName"] != nil {
		d.Set("target_namespace_name", objectRaw["TargetNamespaceName"])
	}
	if objectRaw["TargetRegionId"] != nil {
		d.Set("target_region_id", objectRaw["TargetRegionId"])
	}
	if objectRaw["TargetRepoName"] != nil {
		d.Set("target_repo_name", objectRaw["TargetRepoName"])
	}
	if objectRaw["LocalInstanceId"] != nil {
		d.Set("instance_id", objectRaw["LocalInstanceId"])
	}
	if objectRaw["SyncRuleId"] != nil {
		d.Set("repo_sync_rule_id", objectRaw["SyncRuleId"])
		d.Set("rule_id", objectRaw["SyncRuleId"])
	}
	if objectRaw["SyncDirection"] != nil {
		d.Set("sync_direction", objectRaw["SyncDirection"])
	}

	return nil
}

func resourceAliCloudCrRepoSyncRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource alicloud_cr_ee_sync_rule.")
	return nil
}

func resourceAliCloudCrRepoSyncRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteRepoSyncRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewAcrClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}

	syncDirection := d.Get("sync_direction").(string)
	if syncDirection != "FROM" {
		return WrapError(Error(DefaultErrorMsg, d.Id(), "delete", "[Please delete sync rule in the source instance]"))
	}

	request["RegionId"] = client.RegionId
	request["InstanceId"] = parts[0]
	request["SyncRuleId"] = parts[2]

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-12-01"), StringPointer("AK"), query, request, &runtime)

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
