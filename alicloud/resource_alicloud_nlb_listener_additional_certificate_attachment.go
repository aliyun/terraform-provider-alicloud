// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudNlbListenerAdditionalCertificateAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudNlbListenerAdditionalCertificateAttachmentCreate,
		Read:   resourceAliCloudNlbListenerAdditionalCertificateAttachmentRead,
		Update: resourceAliCloudNlbListenerAdditionalCertificateAttachmentUpdate,
		Delete: resourceAliCloudNlbListenerAdditionalCertificateAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"certificate_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudNlbListenerAdditionalCertificateAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AssociateAdditionalCertificatesWithListener"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["ListenerId"] = d.Get("listener_id")
	request["AdditionalCertificateIds.1"] = d.Get("certificate_id")
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), query, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nlb_listener_additional_certificate_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["ListenerId"], request["AdditionalCertificateIds.1"]))

	nlbServiceV2 := NlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Associated"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nlbServiceV2.NlbListenerAdditionalCertificateAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudNlbListenerAdditionalCertificateAttachmentRead(d, meta)
}

func resourceAliCloudNlbListenerAdditionalCertificateAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nlbServiceV2 := NlbServiceV2{client}

	objectRaw, err := nlbServiceV2.DescribeNlbListenerAdditionalCertificateAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nlb_listener_additional_certificate_attachment DescribeNlbListenerAdditionalCertificateAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}
	if objectRaw["CertificateId"] != nil {
		d.Set("certificate_id", objectRaw["CertificateId"])
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("listener_id", parts[0])

	return nil
}

func resourceAliCloudNlbListenerAdditionalCertificateAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Listener Additional Certificate Attachment.")
	return nil
}

func resourceAliCloudNlbListenerAdditionalCertificateAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DisassociateAdditionalCertificatesWithListener"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["ListenerId"] = parts[0]
	request["AdditionalCertificateIds.1"] = parts[1]
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

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
		if IsExpectedErrors(err, []string{"ResourceNotFound.listener", "ResourceNotFound.loadbalancer"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	nlbServiceV2 := NlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Succeeded"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, nlbServiceV2.DescribeAsyncNlbListenerAdditionalCertificateAttachmentStateRefreshFunc(d, response, "$.Status", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return nil
}
