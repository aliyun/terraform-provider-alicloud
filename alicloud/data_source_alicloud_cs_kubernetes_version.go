package alicloud

import (
	"strconv"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	cs "github.com/alibabacloud-go/cs-20151215/v3/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCSKubernetesVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCSKubernetesVersionRead,
		Schema: map[string]*schema.Schema{
			"cluster_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Kubernetes", "ManagedKubernetes"}, false),
				Required:     true,
			},
			"kubernetes_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"profile": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Default", "Serverless", "Edge"}, false),
				Optional:     true,
			},
			"metadata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"runtime": {
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"version": {
										Type:     schema.TypeString,
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

func dataSourceAlicloudCSKubernetesVersionRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*connectivity.AliyunClient)

	resp, _err := describeKubernetesVersionMetadata(d, client)
	if _err != nil {
		return WrapErrorf(_err, DefaultErrorMsg, "DescribeKubernetesVersionMetadata", err)
	}

	var results []map[string]interface{}

	for _, metadata := range resp {
		result := make(map[string]interface{})
		result["version"] = *metadata.Version
		result["runtime"] = formatRuntime(metadata.Runtimes)
		results = append(results, result)
	}

	d.Set("cluster_type", d.Get("cluster_type").(string))
	if version, ok := d.GetOk("kubernetes_version"); ok {
		d.Set("kubernetes_version", version.(string))
	}
	if profile, ok := d.GetOk("profile"); ok {
		d.Set("profile", profile.(string))
	}
	d.Set("metadata", results)

	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))
	return nil
}

func describeKubernetesVersionMetadata(d *schema.ResourceData, client *connectivity.AliyunClient) ([]*cs.DescribeKubernetesVersionMetadataResponseBody, error) {
	csClient, err := client.NewRoaCsClient()
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "InitializeClient", err)
	}
	request := &cs.DescribeKubernetesVersionMetadataRequest{}
	request.Region = tea.String(client.RegionId)
	request.ClusterType = tea.String(d.Get("cluster_type").(string))
	request.Profile = tea.String(d.Get("profile").(string))
	request.KubernetesVersion = tea.String(d.Get("kubernetes_version").(string))
	resp, err := csClient.DescribeKubernetesVersionMetadata(request)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func formatRuntime(runtimes []*cs.Runtime) []map[string]interface{} {
	result := []map[string]interface{}{}
	for _, runtime := range runtimes {
		mapping := map[string]interface{}{
			"name":    *runtime.Name,
			"version": *runtime.Version,
		}
		result = append(result, mapping)
	}
	return result
}
