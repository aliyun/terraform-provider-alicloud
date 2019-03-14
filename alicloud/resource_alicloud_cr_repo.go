package alicloud

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudCRRepo() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCRRepoCreate,
		Read:   resourceAlicloudCRRepoRead,
		Update: resourceAlicloudCRRepoUpdate,
		Delete: resourceAlicloudCRRepoDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateContainerRegistryNamespaceName,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateContainerRegistryRepoName,
			},
			"summary": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 100),
			},
			"repo_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{string(RepoTypePublic), string(RepoTypePrivate)}),
			},
			"detail": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(0, 2000),
			},
			// computed
			"domain_list": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internal": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudCRRepoCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	invoker := NewInvoker()

	repoNamespace := d.Get("namespace").(string)
	repoName := d.Get("name").(string)

	payload := &crCreateRepoRequestPayload{}
	payload.Repo.RepoNamespace = repoNamespace
	payload.Repo.RepoName = repoName
	payload.Repo.Summary = d.Get("summary").(string)
	payload.Repo.Detail = d.Get("detail").(string)
	payload.Repo.RepoType = d.Get("repo_type").(string)
	serialized, err := json.Marshal(payload)
	if err != nil {
		return WrapError(err)
	}

	req := cr.CreateCreateRepoRequest()
	req.SetContent(serialized)

	if err := invoker.Run(func() error {
		_, err := client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
			return crClient.CreateRepo(req)
		})
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_repo", req.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", repoNamespace, SLASH_SEPARATED, repoName))

	return resourceAlicloudCRRepoRead(d, meta)
}

func resourceAlicloudCRRepoUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	invoker := NewInvoker()

	if d.HasChange("summary") || d.HasChange("detail") || d.HasChange("repo_type") {
		payload := &crUpdateRepoRequestPayload{}
		payload.Repo.Summary = d.Get("summary").(string)
		payload.Repo.Detail = d.Get("detail").(string)
		payload.Repo.RepoType = d.Get("repo_type").(string)

		serialized, err := json.Marshal(payload)
		if err != nil {
			return WrapError(err)
		}
		req := cr.CreateUpdateRepoRequest()
		req.SetContent(serialized)
		req.RepoName = d.Get("name").(string)
		req.RepoNamespace = d.Get("namespace").(string)

		if err := invoker.Run(func() error {
			_, err := client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
				return crClient.UpdateRepo(req)
			})
			return err
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), req.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudCRRepoRead(d, meta)
}

func resourceAlicloudCRRepoRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := CrService{client}

	raw, err := crService.DescribeRepo(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	var resp crDescribeRepoResponse
	err = json.Unmarshal(raw.GetHttpContentBytes(), &resp)
	if err != nil {
		return WrapError(err)
	}

	d.Set("namespace", resp.Data.Repo.RepoNamespace)
	d.Set("name", resp.Data.Repo.RepoName)
	d.Set("detail", resp.Data.Repo.Detail)
	d.Set("summary", resp.Data.Repo.Summary)
	d.Set("repo_type", resp.Data.Repo.RepoType)

	domainList := make(map[string]string)
	domainList["public"] = resp.Data.Repo.RepoDomainList.Public
	domainList["internal"] = resp.Data.Repo.RepoDomainList.Internal
	domainList["vpc"] = resp.Data.Repo.RepoDomainList.Vpc

	d.Set("domain_list", domainList)

	return nil
}

func resourceAlicloudCRRepoDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	invoker := NewInvoker()
	crService := CrService{client}

	sli := strings.Split(d.Id(), SLASH_SEPARATED)
	repoNamespace := sli[0]
	repoName := sli[1]

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		req := cr.CreateDeleteRepoRequest()
		req.RepoNamespace = repoNamespace
		req.RepoName = repoName

		if err := invoker.Run(func() error {
			_, err := client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
				return crClient.DeleteRepo(req)
			})
			return err
		}); err != nil {
			if IsExceptedError(err, ErrorRepoNotExist) {
				return nil
			}
			return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), req.GetActionName(), AlibabaCloudSdkGoERROR))
		}

		if err := invoker.Run(func() error {
			_, err := crService.DescribeRepo(d.Id())
			return err
		}); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), req.GetActionName(), AlibabaCloudSdkGoERROR))
		}

		time.Sleep(5 * time.Second)
		return resource.RetryableError(WrapError(Error("DeleteRepo timeout")))
	})
}
