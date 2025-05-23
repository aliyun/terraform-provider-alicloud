package alicloud

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	cs "github.com/alibabacloud-go/cs-20151215/v5/client"
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
				ValidateFunc: validation.StringInSlice([]string{"Default", "Serverless", "Edge", "Acs"}, false),
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
	client := meta.(*connectivity.AliyunClient)

	var resp []*cs.DescribeKubernetesVersionMetadataResponseBody
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = describeKubernetesVersionMetadata(d, client)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "DescribeKubernetesVersionMetadata", err)
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
	var resp *cs.DescribeKubernetesVersionMetadataResponse
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = csClient.DescribeKubernetesVersionMetadata(request)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

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
