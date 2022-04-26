package alicloud

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	cs "github.com/alibabacloud-go/cs-20151215/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCSKubernetesVersion() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceAlicloudCSKubernetesVersionRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Kubernetes", "ManagedKubernetes", "Ask", "ExternalKubernetes"}, false),
				Required:     true,
			},
			//Computed value
			"kubernetes_versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeMap},
			},
		},
	}
}

func describeKubernetesVersionMetadata(d *schema.ResourceData, client *cs.Client) ([]*cs.DescribeKubernetesVersionMetadataResponseBody, error) {

	request := &cs.DescribeKubernetesVersionMetadataRequest{}
	request.Region = client.RegionId
	if v, ok := d.GetOk("cluster_type"); ok && v.(string) != "" {
		clusterType := v.(string)
		request.ClusterType = &clusterType
	}
	resp, err := client.DescribeKubernetesVersionMetadata(request)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func dataSourceAlicloudCSKubernetesVersionRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "InitializeClient", err)
	}

	resp, _err := describeKubernetesVersionMetadata(d, client)
	if _err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "DescribeKubernetesVersionMetadata", err)
	}
	res := make([]map[string]interface{}, 0)

	var s []map[string]interface{}

	for _, object_index := range resp {
		var runtume_list []map[string]interface{}
		for _, runtume_index := range object_index.Runtimes {
			mapping := map[string]interface{}{
				"name":    runtume_index.Name,
				"version": runtume_index.Version,
			}
			runtume_list = append(runtume_list, mapping)
		}
		mapping := map[string]interface{}{
			"runtimes": runtume_list,
			"version":  object_index.Version,
		}
		s = append(s, mapping)
	}
	cluster_type := d.Get("cluster_type").(string)
	res = append(res, map[string]interface{}{
		"cluster_type":        cluster_type,
		"kubernetes_versions": s,
	})
	d.Set("cluster_type", cluster_type)
	d.Set("kubernetes_versions", s)
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), res)
	}
	return nil
}
