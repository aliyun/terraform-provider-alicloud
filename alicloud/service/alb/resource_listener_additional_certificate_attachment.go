package alb

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceListenerAdditionalCertificateAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceListenerAdditionalCertificateAttachmentCreate,
		Read:   resourceListenerAdditionalCertificateAttachmentRead,
		Delete: resourceListenerAdditionalCertificateAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"certificate_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceListenerAdditionalCertificateAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AssociateAdditionalCertificatesWithListener"
	request := make(map[string]interface{})
	var err error

	request["Certificates.1.CertificateId"] = d.Get("certificate_id")
	request["ListenerId"] = d.Get("listener_id")
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken("AssociateAdditionalCertificatesWithListener")
	wait := helper.IncrementalWait(5*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"ResourceInConfiguring.Listener", "IncorrectStatus.Listener"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_alb_listener_additional_certificate_attachment", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(d.Get("listener_id"), ":", d.Get("certificate_id")))

	albService := AlbService{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"Associated"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, albService.AlbListenerAdditionalCertificateAttachmentStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return resourceListenerAdditionalCertificateAttachmentRead(d, meta)
}
func resourceListenerAdditionalCertificateAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}
	object, err := albService.DescribeAlbListenerAdditionalCertificateAttachment(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alb_listener_additional_certificate_attachment albService.DescribeAlbListenerAdditionalCertificateAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("listener_id", parts[0])
	d.Set("certificate_id", parts[1])
	d.Set("status", object["Status"])
	d.Set("certificate_type", object["CertificateType"])
	return nil
}
func resourceListenerAdditionalCertificateAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	action := "DissociateAdditionalCertificatesFromListener"
	var response map[string]interface{}
	request := map[string]interface{}{
		"ListenerId": parts[0],
	}
	request["Certificates.1.CertificateId"] = parts[1]
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken("DissociateAdditionalCertificatesFromListener")
	wait := helper.IncrementalWait(5*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"ResourceInConfiguring.Listener", "IncorrectStatus.Listener"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	albService := AlbService{client}
	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, albService.AlbListenerAdditionalCertificateAttachmentStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return nil
}
