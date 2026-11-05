package cs

import (
	"fmt"
	"time"

	roaCS "github.com/alibabacloud-go/cs-20151215/v8/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceCSClusterCredential() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCSClusterCredentialRead,

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
				Type:     schema.TypeList,
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

func dataSourceCSClusterCredentialRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	roaClient, err := client.NewRoaCsClient()
	if err != nil {
		return errs.WrapError(err)
	}

	csClient := CsClient{roaClient}

	clusterId := d.Get("cluster_id").(string)
	cluster, err := csClient.client.DescribeClusterDetail(tea.String(clusterId))
	if err != nil {
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cs_cluster_credential", "DescribeClusterDetail", errs.AlibabaCloudSdkGoERROR)
	}

	return CsClusterAuthDescriptionAttributes(d, meta, cluster.Body)
}

func CsClusterAuthDescriptionAttributes(d *schema.ResourceData, meta interface{}, cluster *roaCS.DescribeClusterDetailResponseBody) error {
	client := meta.(*connectivity.AliyunClient)
	roaClient, err := client.NewRoaCsClient()
	if err != nil {
		return errs.WrapError(err)
	}

	csClient := CsClient{roaClient}

	var expiration int64 = 0
	if v, ok := d.GetOk("temporary_duration_minutes"); ok {
		expiration = int64(v.(int))
	}

	clusterId := tea.StringValue(cluster.ClusterId)
	var credential *roaCS.DescribeClusterUserKubeconfigResponseBody
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		credential, err = csClient.DescribeClusterKubeConfigWithExpiration(clusterId, expiration)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cs_cluster_credential", "DescribeClusterKubeConfigWithExpiration", errs.AlibabaCloudSdkGoERROR)
	}

	d.Set("cluster_id", clusterId)
	d.Set("cluster_name", tea.StringValue(cluster.Name))
	d.Set("kube_config", tea.StringValue(credential.Config))
	d.Set("expiration", tea.StringValue(credential.Expiration))
	certificateAuthority, err := FlattenAlicloudCSCertificate(credential)
	if err != nil {
		return errs.WrapError(fmt.Errorf("failed to parse kubeconfig for cluster %s: %s", clusterId, err))
	}
	if err := d.Set("certificate_authority", []interface{}{certificateAuthority}); err != nil {
		return errs.WrapError(fmt.Errorf("error setting certificate_authority: %s", err))
	}
	d.SetId(helper.DataResourceIdHash([]string{clusterId}))

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), tea.StringValue(credential.Config))
	}

	return nil
}
