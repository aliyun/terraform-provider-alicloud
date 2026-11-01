package cr

import (
	"encoding/json"
	"regexp"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	crsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlicloudCRRepos() *schema.Resource {
	return &schema.Resource{
		Read:               dataSourceAlicloudCRReposRead,
		DeprecationMessage: "Data source `alicloud_cr_repos` has been deprecated and will be removed in the future.",
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			// Computed values
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"repos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"summary": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repo_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
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
						"tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"image_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"digest": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"image_size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"image_update": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"image_create": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
func dataSourceAlicloudCRReposRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	invoker := helper.NewInvoker()

	getRepoListRequest := crsdk.CreateGetRepoListRequest()
	getRepoListRequest.RegionId = string(client.Region)
	getRepoListRequest.PageSize = requests.NewInteger(helper.PageSizeMedium)
	getRepoListRequest.Page = requests.NewInteger(1)

	var repos []crRepo
	for {
		var getRepoListResponse *crsdk.GetRepoListResponse

		if err := invoker.Run(func() error {
			raw, err := client.WithCrClient(func(crClient *crsdk.Client) (interface{}, error) {
				return crClient.GetRepoList(getRepoListRequest)
			})
			getRepoListResponse, _ = raw.(*crsdk.GetRepoListResponse)
			return err
		}); err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cr_repos", "GetRepoList", errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(getRepoListRequest.GetActionName(), getRepoListResponse)
		var crResp crDescribeReposResponse

		err := json.Unmarshal(getRepoListResponse.GetHttpContentBytes(), &crResp)
		if err != nil {
			return errs.WrapError(err)
		}

		repos = append(repos, crResp.Data.Repos...)

		if len(crResp.Data.Repos) < helper.PageSizeMedium {
			break
		}

		if page, err := helper.GetNextpageNumber(getRepoListRequest.Page); err != nil {
			return errs.WrapError(err)
		} else {
			getRepoListRequest.Page = page
		}
	}

	var names []string
	var s []map[string]interface{}

	for _, repo := range repos {

		if namespace, ok := d.GetOk("namespace"); ok {
			if repo.RepoNamespace != namespace {
				continue
			}
		}

		if nameRegex, ok := d.GetOk("name_regex"); ok {
			r, err := regexp.Compile(nameRegex.(string))
			if err != nil {
				return errs.WrapError(err)
			}
			if !r.MatchString(repo.RepoName) {
				continue
			}
		}

		mapping := map[string]interface{}{
			"namespace": repo.RepoNamespace,
			"name":      repo.RepoName,
			"summary":   repo.Summary,
			"repo_type": repo.RepoType,
		}

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			names = append(names, repo.RepoName)
			s = append(s, mapping)
			continue
		}

		domainList := make(map[string]string)
		domainList["public"] = repo.RepoDomainList.Public
		domainList["internal"] = repo.RepoDomainList.Internal
		domainList["vpc"] = repo.RepoDomainList.Vpc

		mapping["domain_list"] = []interface{}{domainList}

		var tags []crTag

		getRepoTagsRequest := crsdk.CreateGetRepoTagsRequest()
		getRepoTagsRequest.RegionId = string(client.Region)
		getRepoTagsRequest.PageSize = requests.NewInteger(helper.PageSizeMedium)
		getRepoTagsRequest.Page = requests.NewInteger(1)
		getRepoTagsRequest.RepoNamespace = repo.RepoNamespace
		getRepoTagsRequest.RepoName = repo.RepoName

		for {
			var getRepoTagsResponse *crsdk.GetRepoTagsResponse

			if err := invoker.Run(func() error {
				raw, err := client.WithCrClient(func(crClient *crsdk.Client) (interface{}, error) {
					return crClient.GetRepoTags(getRepoTagsRequest)
				})
				getRepoTagsResponse, _ = raw.(*crsdk.GetRepoTagsResponse)
				return err
			}); err != nil {
				return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cr_repos", "GetRepoTags", errs.AlibabaCloudSdkGoERROR)
			}
			logs.AddDebug(getRepoTagsRequest.GetActionName(), getRepoTagsResponse)
			var crResp crDescribeRepoTagsResponse

			err := json.Unmarshal(getRepoTagsResponse.GetHttpContentBytes(), &crResp)
			if err != nil {
				return errs.WrapError(err)
			}

			tags = append(tags, crResp.Data.Tags...)

			if len(crResp.Data.Tags) < helper.PageSizeMedium {
				break
			}

			if page, err := helper.GetNextpageNumber(getRepoTagsRequest.Page); err != nil {
				return errs.WrapError(err)
			} else {
				getRepoTagsRequest.Page = page
			}
		}

		var tagList []map[string]interface{}
		for _, tag := range tags {
			tagList = append(tagList, map[string]interface{}{
				"tag":          tag.Tag,
				"image_id":     tag.ImageId,
				"digest":       tag.Digest,
				"status":       tag.Status,
				"image_size":   tag.ImageSize,
				"image_update": tag.ImageUpdate,
				"image_create": tag.ImageCreate,
			})
		}
		mapping["tags"] = tagList

		names = append(names, repo.RepoName)
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(names))
	if err := d.Set("repos", s); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("ids", names); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
