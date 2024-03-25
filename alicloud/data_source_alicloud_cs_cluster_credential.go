package alicloud

import (
	"github.com/alibabacloud-go/cs-20151215/v5/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudCSClusterCredential() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCSClusterCredentialRead,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"temporary_duration_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed fields
			"cluster_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kube_config": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"certificate_authority": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_cert": {
							Type:      schema.TypeString,
							Sensitive: true,
							Computed:  true,
						},
						"client_cert": {
							Type:      schema.TypeString,
							Sensitive: true,
							Computed:  true,
						},
						"client_key": {
							Type:      schema.TypeString,
							Sensitive: true,
							Computed:  true,
						},
					},
				},
			},
			"expiration": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudCSClusterCredentialRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	roaClient, err := client.NewRoaCsClient()
	if err != nil {
		return WrapError(err)
	}

	csClient := CsClient{roaClient}

	clusterId := d.Get("cluster_id").(string)
	cluster, err := csClient.client.DescribeClusterDetail(tea.String(clusterId))
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cs_cluster_credential", "DescribeClusterDetail", AlibabaCloudSdkGoERROR)
	}

	return csClusterAuthDescriptionAttributes(d, meta, cluster.Body)
}

func csClusterAuthDescriptionAttributes(d *schema.ResourceData, meta interface{}, cluster *client.DescribeClusterDetailResponseBody) error {
	client := meta.(*connectivity.AliyunClient)
	roaClient, err := client.NewRoaCsClient()
	if err != nil {
		return WrapError(err)
	}

	csClient := CsClient{roaClient}

	var expiration int64 = 0
	if v, ok := d.GetOk("temporary_duration_minutes"); ok {
		expiration = int64(v.(int))
	}

	clusterId := tea.StringValue(cluster.ClusterId)
	credential, err := csClient.DescribeClusterKubeConfigWithExpiration(clusterId, expiration)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cs_cluster_credential", "DescribeClusterKubeConfigWithExpiration", AlibabaCloudSdkGoERROR)
	}

	d.Set("cluster_id", clusterId)
	d.Set("cluster_name", tea.StringValue(cluster.Name))
	d.Set("kube_config", tea.StringValue(credential.Config))
	d.Set("expiration", tea.StringValue(credential.Expiration))
	d.Set("certificate_authority", flattenAlicloudCSCertificate(credential))
	d.SetId(dataResourceIdHash([]string{clusterId}))

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), tea.StringValue(credential.Config))
	}

	return nil
}
