package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCrEESyncRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCrEESyncRuleCreate,
		Read:   resourceAliCloudCrEESyncRuleRead,
		Delete: resourceAliCloudCrEESyncRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
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
			"name": {
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
			"tag_filter": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repo_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"target_repo_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sync_direction": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sync_scope": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCrEESyncRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := &CrService{client}
	response := &cr_ee.CreateRepoSyncRuleResponse{}
	request := cr_ee.CreateCreateRepoSyncRuleRequest()
	request.RegionId = crService.client.RegionId
	request.SyncTrigger = "PASSIVE"
	request.InstanceId = d.Get("instance_id").(string)
	request.NamespaceName = d.Get("namespace_name").(string)
	request.SyncRuleName = d.Get("name").(string)
	request.TargetInstanceId = d.Get("target_instance_id").(string)
	request.TargetNamespaceName = d.Get("target_namespace_name").(string)
	request.TargetRegionId = d.Get("target_region_id").(string)
	request.TagFilter = d.Get("tag_filter").(string)

	var repoName, targetRepoName string

	if v, ok := d.GetOk("repo_name"); ok {
		repoName = v.(string)
	}

	if v, ok := d.GetOk("target_repo_name"); ok {
		targetRepoName = v.(string)
	}

	if (repoName != "" && targetRepoName == "") || (repoName == "" && targetRepoName != "") {
		return WrapError(Error("repo_name and target_repo_name must be set at the same time"))
	}

	if repoName != "" && targetRepoName != "" {
		request.SyncScope = "REPO"
		request.RepoName = repoName
		request.TargetRepoName = targetRepoName
	} else {
		request.SyncScope = "NAMESPACE"
	}

	var raw interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		raw, err = crService.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
			return creeClient.CreateRepoSyncRule(request)
		})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_ee_sync_rule", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ = raw.(*cr_ee.CreateRepoSyncRuleResponse)
	if !response.CreateRepoSyncRuleIsSuccess {
		return WrapErrorf(fmt.Errorf("%v", response), DefaultErrorMsg, "alicloud_cr_ee_sync_rule", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s:%s:%s", request.InstanceId, request.NamespaceName, response.SyncRuleId))

	return resourceAliCloudCrEESyncRuleRead(d, meta)
}

func resourceAliCloudCrEESyncRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := &CrService{client}

	object, err := crService.DescribeCrEESyncRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", object.LocalInstanceId)
	d.Set("namespace_name", object.LocalNamespaceName)
	d.Set("name", object.SyncRuleName)
	d.Set("target_instance_id", object.TargetInstanceId)
	d.Set("target_namespace_name", object.TargetNamespaceName)
	d.Set("target_region_id", object.TargetRegionId)
	d.Set("tag_filter", object.TagFilter)
	d.Set("repo_name", object.LocalRepoName)
	d.Set("target_repo_name", object.TargetRepoName)
	d.Set("rule_id", object.SyncRuleId)
	d.Set("sync_direction", object.SyncDirection)
	d.Set("sync_scope", object.SyncScope)

	return nil
}

func resourceAliCloudCrEESyncRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := &CrService{client}
	response := &cr_ee.DeleteRepoSyncRuleResponse{}

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}

	syncDirection := d.Get("sync_direction").(string)
	if syncDirection != "FROM" {
		return WrapError(Error(DefaultErrorMsg, d.Id(), "delete", "[Please delete sync rule in the source instance]"))
	}

	request := cr_ee.CreateDeleteRepoSyncRuleRequest()
	request.RegionId = crService.client.RegionId
	request.InstanceId = parts[0]
	request.SyncRuleId = parts[2]

	var raw interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		raw, err = crService.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
			return creeClient.DeleteRepoSyncRule(request)
		})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ = raw.(*cr_ee.DeleteRepoSyncRuleResponse)
	if !response.DeleteRepoSyncRuleIsSuccess {
		return WrapErrorf(fmt.Errorf("%v", response), DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return nil
}
