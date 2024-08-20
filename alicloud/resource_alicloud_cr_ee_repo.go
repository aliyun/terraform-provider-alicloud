package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCrEERepo() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCrEERepoCreate,
		Read:   resourceAliCloudCrEERepoRead,
		Update: resourceAliCloudCrEERepoUpdate,
		Delete: resourceAliCloudCrEERepoDelete,
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repo_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{RepoTypePublic, RepoTypePrivate}, false),
			},
			"summary": {
				Type:     schema.TypeString,
				Required: true,
			},
			"detail": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"repo_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCrEERepoCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := &CrService{client}
	response := &cr_ee.CreateRepositoryResponse{}
	request := cr_ee.CreateCreateRepositoryRequest()
	request.RegionId = crService.client.RegionId

	request.InstanceId = d.Get("instance_id").(string)
	request.RepoNamespaceName = d.Get("namespace").(string)
	request.RepoName = d.Get("name").(string)
	request.RepoType = d.Get("repo_type").(string)
	request.Summary = d.Get("summary").(string)

	if v, ok := d.GetOk("detail"); ok {
		request.Detail = v.(string)
	}

	var raw interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		raw, err = crService.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
			return creeClient.CreateRepository(request)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_ee_repo", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ = raw.(*cr_ee.CreateRepositoryResponse)
	if !response.CreateRepositoryIsSuccess {
		return WrapErrorf(fmt.Errorf("%v", response), DefaultErrorMsg, "alicloud_cr_ee_repo", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s:%s:%s", request.InstanceId, request.RepoNamespaceName, request.RepoName))

	return resourceAliCloudCrEERepoRead(d, meta)
}

func resourceAliCloudCrEERepoRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := &CrService{client}

	resp, err := crService.DescribeCrEERepo(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
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

func resourceAliCloudCrEERepoUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := &CrService{client}
	response := &cr_ee.UpdateRepositoryResponse{}

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}

	update := false
	request := cr_ee.CreateUpdateRepositoryRequest()
	request.RegionId = crService.client.RegionId
	request.InstanceId = parts[0]
	request.RepoId = d.Get("repo_id").(string)

	if d.HasChange("repo_type") {
		update = true
	}
	request.RepoType = d.Get("repo_type").(string)

	if d.HasChange("summary") {
		update = true
	}
	request.Summary = d.Get("summary").(string)

	if d.HasChange("detail") {
		update = true
	}
	if v, ok := d.GetOk("detail"); ok {
		request.Detail = v.(string)
	}

	if update {
		var raw interface{}
		var err error
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			raw, err = crService.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
				return creeClient.UpdateRepository(request)
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

		response, _ = raw.(*cr_ee.UpdateRepositoryResponse)
		if !response.UpdateRepositoryIsSuccess {
			return WrapErrorf(fmt.Errorf("%v", response), DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudCrEERepoRead(d, meta)
}

func resourceAliCloudCrEERepoDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := &CrService{client}

	repoId := d.Get("repo_id").(string)
	_, err := crService.DeleteCrEERepo(d.Id(), repoId)
	if err != nil {
		if NotFoundError(err) {
			return nil
		} else {
			return WrapError(err)
		}
	}

	return WrapError(crService.WaitForCrEERepo(d.Id(), Deleted, DefaultTimeout))
}
