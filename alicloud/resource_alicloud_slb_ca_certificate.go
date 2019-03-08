package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"ca_certificate": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudSlbCACertificateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	request := slb.CreateUploadCACertificateRequest()

	if val, ok := d.GetOk("name"); ok && val.(string) != "" {
		request.CACertificateName = val.(string)
	}

	if val, ok := d.GetOk("ca_certificate"); ok && val.(string) != "" {
		request.CACertificate = val.(string)
	} else {
		return fmt.Errorf("UploadCACertificate got an error, ca_certificate should be not null")
	}

	raw, err := slbService.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.UploadCACertificate(request)
	})
	if err != nil {
		return fmt.Errorf("UploadCACertificate got an error: %#v", err)
	}
	response := raw.(*slb.UploadCACertificateResponse)

	d.SetId(response.CACertificateId)

	return resourceAlicloudSlbCACertificateUpdate(d, meta)
}

func resourceAlicloudSlbCACertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	caCertificate, err := slbService.describeSlbCACertificate(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	if error := d.Set("name", caCertificate.CACertificateName); error != nil {
		return error
	}

	return nil
}

func resourceAlicloudSlbCACertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	if !d.IsNewResource() && d.HasChange("name") {
		request := slb.CreateSetCACertificateNameRequest()
		request.CACertificateId = d.Id()
		request.CACertificateName = d.Get("name").(string)
		_, err := slbService.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.SetCACertificateName(request)
		})
		if err != nil {
			return fmt.Errorf("SetCACertificateName set %s  name %s got an error: %#v",
				d.Id(), request.CACertificateName, err)
		}
	}

	return resourceAlicloudSlbCACertificateRead(d, meta)
}

func resourceAlicloudSlbCACertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := slb.CreateDeleteCACertificateRequest()
		request.CACertificateId = d.Id()
		_, err := slbService.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteCACertificate(request)
		})
		if err != nil {
			if IsExceptedError(err, SlbCACertificateIdNotFound) || NotFoundError(err) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("DeleteCACertificate %s got an error: %#v.", d.Id(), err))
		}

		if _, err := slbService.describeSlbCACertificate(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("While DeleteCACertificate，DescribeCACertificates %s got an error: %#v.", d.Id(), err))
		}

		return resource.RetryableError(fmt.Errorf("DeleteCACertificate %s timeout.", d.Id()))
	})
}
