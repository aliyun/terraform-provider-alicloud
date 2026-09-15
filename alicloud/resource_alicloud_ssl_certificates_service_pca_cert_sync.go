package alicloud

import (
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudSslCertificatesServicePcaCertSync() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSslCertificatesServicePcaCertSyncCreate,
		Read:   resourceAliCloudSslCertificatesServicePcaCertSyncRead,
		Delete: resourceAliCloudSslCertificatesServicePcaCertSyncDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAliCloudSslCertificatesServicePcaCertSyncCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	ids := expandStringList(d.Get("ids").([]interface{}))

	action := "UploadPcaCertToCas"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["Ids"] = strings.Join(ids, ",")

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("cas", "2020-06-30", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ssl_certificates_service_pca_cert_sync", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(strings.Join(ids, ","))

	return nil
}

func resourceAliCloudSslCertificatesServicePcaCertSyncRead(d *schema.ResourceData, meta interface{}) error {
	d.Set("ids", strings.Split(d.Id(), ","))

	return nil
}

func resourceAliCloudSslCertificatesServicePcaCertSyncDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Ssl Certificates Service Pca Cert Sync. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
