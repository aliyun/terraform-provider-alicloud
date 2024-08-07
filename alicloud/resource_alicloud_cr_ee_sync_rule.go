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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringLenBetween(2, 30),
			},
			"target_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_namespace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringLenBetween(2, 30),
			},
			"target_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
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
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringLenBetween(2, 64),
			},
			"target_repo_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringLenBetween(2, 64),
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

	instanceId := d.Get("instance_id").(string)
	namespaceName := d.Get("namespace_name").(string)
	targetInstanceId := d.Get("target_instance_id").(string)
	targetNamespaceName := d.Get("target_namespace_name").(string)
	targetRegionId := d.Get("target_region_id").(string)
	syncRuleName := d.Get("name").(string)
	tagFilter := d.Get("tag_filter").(string)

	var repoName, targetRepoName string
	if v, ok := d.GetOk("repo_name"); ok {
		repoName = v.(string)
	}

	if v, ok := d.GetOk("target_repo_name"); ok {
		targetRepoName = v.(string)
	}

	if (repoName != "" && targetRepoName == "") || (repoName == "" && targetRepoName != "") {
		return WrapError(Error(DefaultErrorMsg, syncRuleName, "create", "[Params repo_name or target_repo_name is empty]"))
	}

	request.RegionId = crService.client.RegionId
	request.InstanceId = instanceId
	request.NamespaceName = namespaceName
	request.TargetInstanceId = targetInstanceId
	request.TargetNamespaceName = targetNamespaceName
	request.TargetRegionId = targetRegionId
	request.SyncRuleName = syncRuleName
	request.TagFilter = tagFilter
	request.SyncTrigger = "PASSIVE"

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

	d.SetId(fmt.Sprintf("%v:%v:%v", instanceId, namespaceName, response.SyncRuleId))

	return resourceAliCloudCrEESyncRuleRead(d, meta)
}

func resourceAliCloudCrEESyncRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := &CrService{client}
	resp, err := crService.DescribeCrEESyncRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", resp.LocalInstanceId)
	d.Set("namespace_name", resp.LocalNamespaceName)
	d.Set("target_instance_id", resp.TargetInstanceId)
	d.Set("target_namespace_name", resp.TargetNamespaceName)
	d.Set("target_region_id", resp.TargetRegionId)
	d.Set("name", resp.SyncRuleName)
	d.Set("tag_filter", resp.TagFilter)
	d.Set("repo_name", resp.LocalRepoName)
	d.Set("target_repo_name", resp.TargetRepoName)
	d.Set("rule_id", resp.SyncRuleId)
	d.Set("sync_direction", resp.SyncDirection)
	d.Set("sync_scope", resp.SyncScope)

	return nil
}

func resourceAliCloudCrEESyncRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := &CrService{client}
	response := &cr_ee.DeleteRepoSyncRuleResponse{}
	request := cr_ee.CreateDeleteRepoSyncRuleRequest()

	syncDirection := d.Get("sync_direction").(string)
	if syncDirection != "FROM" {
		return WrapError(Error(DefaultErrorMsg, d.Id(), "delete", "[Please delete sync rule in the source instance]"))
	}

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}

	request.RegionId = crService.client.RegionId
	request.InstanceId = parts[0]
	request.SyncRuleId = parts[2]

	var raw interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
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
