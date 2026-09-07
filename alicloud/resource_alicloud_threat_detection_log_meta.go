// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudThreatDetectionLogMeta() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudThreatDetectionLogMetaCreate,
		Read:   resourceAliCloudThreatDetectionLogMetaRead,
		Update: resourceAliCloudThreatDetectionLogMetaUpdate,
		Delete: resourceAliCloudThreatDetectionLogMetaDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"log_meta_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"aegis-log-client", "aegis-log-crack", "aegis-log-dns-query", "aegis-log-login", "aegis-log-network", "aegis-log-process", "aegis-snapshot-host", "aegis-snapshot-port", "aegis-snapshot-process", "sas-net-block", "sas-log-dns", "sas-log-http", "sas-log-session", "sas-cspm-log", "sas-filedetect-log", "sas-hc-log", "sas-rasp-log", "sas-security-log", "sas-vul-log", "local-dns"}, false),
			},
			"status": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"disabled", "enabled"}, false),
			},
		},
	}
}

func resourceAliCloudThreatDetectionLogMetaCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "ModifyLogMetaStatus"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("log_meta_name"); ok {
		request["LogStore"] = v
	}

	request["Status"] = d.Get("status")
	request["From"] = "sas"
	request["Project"] = "sas"
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_log_meta", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(d.Get("log_meta_name").(string))

	return resourceAliCloudThreatDetectionLogMetaRead(d, meta)
}

func resourceAliCloudThreatDetectionLogMetaRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionServiceV2 := ThreatDetectionServiceV2{client}

	objectRaw, err := threatDetectionServiceV2.DescribeThreatDetectionLogMeta(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_log_meta DescribeThreatDetectionLogMeta Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("status", objectRaw["Status"])
	d.Set("log_meta_name", objectRaw["LogStore"])

	d.Set("log_meta_name", d.Id())

	return nil
}

func resourceAliCloudThreatDetectionLogMetaUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ModifyLogMetaStatus"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["LogStore"] = d.Id()

	if d.HasChange("status") {
		update = true
	}
	request["Status"] = d.Get("status")
	request["From"] = "sas"
	request["Project"] = "sas"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudThreatDetectionLogMetaRead(d, meta)
}

func resourceAliCloudThreatDetectionLogMetaDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Log Meta. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
