package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudSlbCACertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSlbCACertificateCreate,
		Read:   resourceAlicloudSlbCACertificateRead,
		Update: resourceAlicloudSlbCACertificateUpdate,
		Delete: resourceAlicloudSlbCACertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"ca_certificate": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudSlbCACertificateCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient).slbconn

	request := slb.CreateUploadCACertificateRequest()

	if val, ok := d.GetOk("name"); ok && val != "" {
		request.CACertificateName = val.(string)
	}

	if val, ok := d.GetOk("ca_certificate"); ok && val != "" {
		request.CACertificate = val.(string)
	}

	response, err := client.UploadCACertificate(request)
	if err != nil {
		return fmt.Errorf("UploadCACertificate got an error: %#v", err)
	}

	d.SetId(response.CACertificateId)

	return resourceAlicloudSlbCACertificateUpdate(d, meta)
}

func resourceAlicloudSlbCACertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	caCertificate, err := client.describeSlbCACertificate(d.Id())
	if err != nil {
		if IsExceptedError(err, SlbCACertificateIdNotFound) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", caCertificate.CACertificateName)

	return nil
}

func resourceAlicloudSlbCACertificateUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient).slbconn

	d.Partial(true)

	if !d.IsNewResource() && d.HasChange("name") {
		request := slb.CreateSetCACertificateNameRequest()
		request.CACertificateId = d.Id()
		request.CACertificateName = d.Get("name").(string)
		if _, err := client.SetCACertificateName(request); err != nil {

			return fmt.Errorf("SetCACertificateName set %s  name %s got an error: %#v",
				d.Id(), request.CACertificateName, err)

		}
		d.SetPartial("name")
	}

	d.Partial(false)

	return resourceAlicloudSlbCACertificateRead(d, meta)
}

func resourceAlicloudSlbCACertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := slb.CreateDeleteCACertificateRequest()
		request.CACertificateId = d.Id()
		if _, err := client.slbconn.DeleteCACertificate(request); err != nil {
			if IsExceptedError(err, SlbCACertificateIdNotFound) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("DeleteCACertificate %s got an error: %#v.", d.Id(), err))
		}

		if _, err := client.describeSlbCACertificate(d.Id()); err != nil {
			if IsExceptedError(err, SlbCACertificateIdNotFound) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("While DeleteCACertificate，DescribeCACertificates %s got an error: %#v.", d.Id(), err))
		}

		return resource.RetryableError(fmt.Errorf("DeleteCACertificate %s timeout.", d.Id()))
	})
}
