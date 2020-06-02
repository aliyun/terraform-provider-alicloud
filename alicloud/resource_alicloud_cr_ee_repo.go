package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudCrEERepo() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCrEERepoCreate,
		Read:   resourceAlicloudCrEERepoRead,
		Update: resourceAlicloudCrEERepoUpdate,
		Delete: resourceAlicloudCrEERepoDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 30),
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 64),
			},
			"summary": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
			"repo_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{RepoTypePublic, RepoTypePrivate}, false),
			},
			"detail": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 2000),
			},

			//Computed
			"repo_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCrEERepoCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := &CrService{client}
	instanceId := d.Get("instance_id").(string)
	namespace := d.Get("namespace").(string)
	repoName := d.Get("name").(string)
	repoType := d.Get("repo_type").(string)
	summary := d.Get("summary").(string)
	detail := d.Get("detail").(string)

	response := &cr_ee.CreateRepositoryResponse{}
	request := cr_ee.CreateCreateRepositoryRequest()
	request.RegionId = crService.client.RegionId
	request.InstanceId = instanceId
	request.RepoNamespaceName = namespace
	request.RepoName = repoName
	request.RepoType = repoType
	request.Summary = summary
	request.Detail = detail
	resource := crService.GenResourceId(instanceId, namespace, repoName)
	action := request.GetActionName()

	raw, err := crService.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.CreateRepository(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, resource, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.CreateRepositoryResponse)
	if !response.CreateRepositoryIsSuccess {
		return crService.wrapCrServiceError(resource, action, response.Code)
	}

	d.SetId(crService.GenResourceId(instanceId, namespace, repoName))

	return resourceAlicloudCrEERepoRead(d, meta)
}

func resourceAlicloudCrEERepoRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := &CrService{client}
	resp, err := crService.DescribeCrEERepo(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", resp.InstanceId)
	d.Set("namespace", resp.RepoNamespaceName)
	d.Set("name", resp.RepoName)
	d.Set("repo_type", resp.RepoType)
	d.Set("summary", resp.Summary)
	d.Set("detail", resp.Detail)
	d.Set("repo_id", resp.RepoId)

	return nil
}

func resourceAlicloudCrEERepoUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := &CrService{client}
	instanceId := d.Get("instance_id").(string)
	namespace := d.Get("namespace").(string)
	repo := d.Get("name").(string)
	if d.HasChanges("repo_type", "summary", "detail") {
		repoId := d.Get("repo_id").(string)
		repoType := d.Get("repo_type").(string)
		summary := d.Get("summary").(string)
		detail := d.Get("detail").(string)

		response := &cr_ee.UpdateRepositoryResponse{}
		request := cr_ee.CreateUpdateRepositoryRequest()
		request.RegionId = crService.client.RegionId
		request.InstanceId = instanceId
		request.RepoId = repoId
		request.RepoType = repoType
		request.Summary = summary
		request.Detail = detail
		resource := crService.GenResourceId(instanceId, namespace, repo)
		action := request.GetActionName()

		raw, err := crService.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
			return creeClient.UpdateRepository(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, resource, action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, raw, request.RpcRequest, request)

		response, _ = raw.(*cr_ee.UpdateRepositoryResponse)
		if !response.UpdateRepositoryIsSuccess {
			return crService.wrapCrServiceError(resource, action, response.Code)
		}
	}

	return resourceAlicloudCrEERepoRead(d, meta)
}

func resourceAlicloudCrEERepoDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := &CrService{client}
	instanceId := d.Get("instance_id").(string)
	namespace := d.Get("namespace").(string)
	repo := d.Get("name").(string)
	repoId := d.Get("repo_id").(string)
	_, err := crService.DeleteCrEERepo(instanceId, namespace, repo, repoId)
	if err != nil {
		if NotFoundError(err) {
			return nil
		} else {
			return WrapError(err)
		}
	}

	return WrapError(crService.WaitForCrEERepo(instanceId, namespace, repo, Deleted, DefaultTimeout))
}
