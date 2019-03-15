package alicloud

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"regexp"
)

func dataSourceAlicloudCRRepos() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCRReposRead,

		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
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
	invoker := NewInvoker()

	req := cr.CreateGetRepoListRequest()
	req.RegionId = string(client.Region)
	req.PageSize = requests.NewInteger(PageSizeMedium)
	req.Page = requests.NewInteger(1)

	var repos []crRepo
	for {
		var resp *cr.GetRepoListResponse

		if err := invoker.Run(func() error {
			raw, err := client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
				return crClient.GetRepoList(req)
			})
			resp, _ = raw.(*cr.GetRepoListResponse)
			return err
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_repos", "GetRepoList", AlibabaCloudSdkGoERROR)
		}
		var crResp crDescribeReposResponse

		err := json.Unmarshal(resp.GetHttpContentBytes(), &crResp)
		if err != nil {
			return WrapError(err)
		}

		repos = append(repos, crResp.Data.Repos...)

		if len(crResp.Data.Repos) < PageSizeMedium {
			break
		}

		if page, err := getNextpageNumber(req.Page); err != nil {
			return WrapError(err)
		} else {
			req.Page = page
		}
	}

	var ids []string
	var s []map[string]interface{}

	for _, repo := range repos {

		if namespace, ok := d.GetOk("namespace"); ok {
			if repo.RepoNamespace != namespace {
				continue
			}
		}

		if nameRegex, ok := d.GetOk("name_regex"); ok {
			r := regexp.MustCompile(nameRegex.(string))
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

		if detailedEnabled, ok := d.GetOk("enable_details"); ok && !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprintf("%s%s%s", repo.RepoNamespace, SLASH_SEPARATED, repo.RepoName))
			s = append(s, mapping)
			continue
		}

		domainList := make(map[string]string)
		domainList["public"] = repo.RepoDomainList.Public
		domainList["internal"] = repo.RepoDomainList.Internal
		domainList["vpc"] = repo.RepoDomainList.Vpc

		mapping["domain_list"] = domainList

		var tags []crTag

		req := cr.CreateGetRepoTagsRequest()
		req.RegionId = string(client.Region)
		req.PageSize = requests.NewInteger(PageSizeMedium)
		req.Page = requests.NewInteger(1)
		req.RepoNamespace = repo.RepoNamespace
		req.RepoName = repo.RepoName

		for {
			var resp *cr.GetRepoTagsResponse

			if err := invoker.Run(func() error {
				raw, err := client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
					return crClient.GetRepoTags(req)
				})
				resp, _ = raw.(*cr.GetRepoTagsResponse)
				return err
			}); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_repos", "GetRepoTags", AlibabaCloudSdkGoERROR)
			}

			var crResp crDescribeRepoTagsResponse

			err := json.Unmarshal(resp.GetHttpContentBytes(), &crResp)
			if err != nil {
				return WrapError(err)
			}

			tags = append(tags, crResp.Data.Tags...)

			if len(crResp.Data.Tags) < PageSizeMedium {
				break
			}

			if page, err := getNextpageNumber(req.Page); err != nil {
				return WrapError(err)
			} else {
				req.Page = page
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

		ids = append(ids, fmt.Sprintf("%s%s%s", repo.RepoNamespace, SLASH_SEPARATED, repo.RepoName))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	d.Set("ids", ids)
	if err := d.Set("repos", s); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
