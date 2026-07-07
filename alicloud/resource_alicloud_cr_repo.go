package alicloud

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
		DeprecationMessage: "This resource has been deprecated since v1.276.0 and will be removed in the future. Please use 'alicloud_cr_ee_repo' instead.",
		SchemaVersion:      1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type:    resourceAlicloudCRRepoV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceAlicloudCRRepoStateUpgradeV0,
			},
		},
		Schema: map[string]*schema.Schema{
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
				ValidateFunc: validation.StringLenBetween(2, 30),
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
			// computed
			"domain_list": {
				Type:     schema.TypeList,
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

	request := cr.CreateCreateRepoRequest()
	request.RegionId = client.RegionId
	request.SetContent(serialized)

	raw, err := client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
		return crClient.CreateRepo(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_repo", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)
	d.SetId(fmt.Sprintf("%s%s%s", repoNamespace, SLASH_SEPARATED, repoName))

	return resourceAlicloudCRRepoRead(d, meta)
}

func resourceAlicloudCRRepoUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if d.HasChange("summary") || d.HasChange("detail") || d.HasChange("repo_type") {
		payload := &crUpdateRepoRequestPayload{}
		payload.Repo.Summary = d.Get("summary").(string)
		payload.Repo.Detail = d.Get("detail").(string)
		payload.Repo.RepoType = d.Get("repo_type").(string)

		serialized, err := json.Marshal(payload)
		if err != nil {
			return WrapError(err)
		}
		request := cr.CreateUpdateRepoRequest()
		request.RegionId = client.RegionId
		request.SetContent(serialized)
		request.RepoName = d.Get("name").(string)
		request.RepoNamespace = d.Get("namespace").(string)

		raw, err := client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
			return crClient.UpdateRepo(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
	}
	return resourceAlicloudCRRepoRead(d, meta)
}

func resourceAlicloudCRRepoRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := CrService{client}

	object, err := crService.DescribeCrRepo(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	var response crDescribeRepoResponse
	err = json.Unmarshal(object.GetHttpContentBytes(), &response)
	if err != nil {
		return WrapError(err)
	}

	d.Set("namespace", response.Data.Repo.RepoNamespace)
	d.Set("name", response.Data.Repo.RepoName)
	d.Set("detail", response.Data.Repo.Detail)
	d.Set("summary", response.Data.Repo.Summary)
	d.Set("repo_type", response.Data.Repo.RepoType)

	domainList := make(map[string]string)
	domainList["public"] = response.Data.Repo.RepoDomainList.Public
	domainList["internal"] = response.Data.Repo.RepoDomainList.Internal
	domainList["vpc"] = response.Data.Repo.RepoDomainList.Vpc

	d.Set("domain_list", []interface{}{domainList})

	return nil
}

func resourceAlicloudCRRepoDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := CrService{client}

	sli := strings.Split(d.Id(), SLASH_SEPARATED)
	repoNamespace := sli[0]
	repoName := sli[1]

	request := cr.CreateDeleteRepoRequest()
	request.RegionId = client.RegionId
	request.RepoNamespace = repoNamespace
	request.RepoName = repoName

	raw, err := client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
		return crClient.DeleteRepo(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"REPO_NOT_EXIST"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)
	return WrapError(crService.WaitForCrRepo(d.Id(), Deleted, DefaultTimeout))
}

func resourceAlicloudCRRepoV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
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
			"summary": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"detail": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// lintignore: S022
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

func resourceAlicloudCRRepoStateUpgradeV0(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	if v, ok := rawState["domain_list"]; ok && v != nil {
		switch val := v.(type) {
		case map[string]interface{}:
			if len(val) > 0 {
				rawState["domain_list"] = []interface{}{val}
			} else {
				rawState["domain_list"] = []interface{}{}
			}
		}
	}
	return rawState, nil
}
