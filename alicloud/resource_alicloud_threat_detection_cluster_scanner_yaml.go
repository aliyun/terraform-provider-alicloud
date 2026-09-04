package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudThreatDetectionClusterScannerYaml() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudThreatDetectionClusterScannerYamlCreate,
		Read:   resourceAlicloudThreatDetectionClusterScannerYamlRead,
		Update: resourceAlicloudThreatDetectionClusterScannerYamlRead,
		Delete: resourceAlicloudThreatDetectionClusterScannerYamlDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"webhook_open": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"ca_cert_base64": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tls_key_base64": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tls_cert_base64": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_env_info": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudThreatDetectionClusterScannerYamlCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	request := make(map[string]interface{})

	if v, ok := d.GetOk("cluster_id"); ok {
		request["ClusterId"] = v
	}
	if v, ok := d.GetOkExists("webhook_open"); ok {
		request["WebhookOpen"] = v
	}

	var response map[string]interface{}
	action := "GenerateClusterScannerWebhookYaml"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_cluster_scanner_yaml", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["ClusterId"]))

	return resourceAlicloudThreatDetectionClusterScannerYamlRead(d, meta)
}

func resourceAlicloudThreatDetectionClusterScannerYamlRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sasService := SasService{client}

	object, err := sasService.DescribeThreatDetectionClusterScannerYaml(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_cluster_scanner_yaml sasService.DescribeThreatDetectionClusterScannerYaml Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("cluster_id", object["ClusterId"])
	d.Set("webhook_open", formatInt(object["WebhookOpen"]))
	d.Set("ca_cert_base64", object["CaCertBase64"])
	d.Set("tls_key_base64", object["TlsKeyBase64"])
	d.Set("tls_cert_base64", object["TlsCertBase64"])
	d.Set("cluster_env_info", object["ClusterEnvInfo"])
	d.Set("image", object["Image"])
	d.Set("region_id", client.RegionId)
	return nil
}

// resourceAlicloudThreatDetectionClusterScannerYamlDelete is a no-op because the
// GenerateClusterScannerWebhookYaml API does not provide a corresponding delete
// operation. Removing the resource from Terraform state does not modify the
// cloud-side scanner configuration.
func resourceAlicloudThreatDetectionClusterScannerYamlDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
