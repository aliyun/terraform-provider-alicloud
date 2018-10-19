package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudSlbServerCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSlbServerCertificateCreate,
		Read:   resourceAlicloudSlbServerCertificateRead,
		Update: resourceAlicloudSlbServerCertificateUpdate,
		Delete: resourceAlicloudSlbServerCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"server_certificate": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: slbServerCertificateDiffSuppressFunc,
			},
			"private_key": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: slbServerCertificateDiffSuppressFunc,
			},

			"server_certificate_file": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: slbServerCertificateFileDiffSuppressFunc,
			},
			"private_key_file": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: slbServerCertificateFileDiffSuppressFunc,
			},

			"alicloud_certifacte_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"alicloud_certifacte_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
		},
	}
}

func resourceAlicloudSlbServerCertificateCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient).slbconn

	request := slb.CreateUploadServerCertificateRequest()

	if val, ok := d.GetOk("name"); ok && val != "" {
		request.ServerCertificateName = val.(string)
	}

	if val, ok := d.GetOk("server_certificate"); ok && val != "" {
		request.ServerCertificate = val.(string)
	}

	if val, ok := d.GetOk("private_key"); ok && val != "" {
		request.PrivateKey = val.(string)
	}

	if val, ok := d.GetOk("alicloud_certificate_id"); ok && val != "" {
		request.AliCloudCertificateId = val.(string)
	}

	if val, ok := d.GetOk("alicloud_certificate_name"); ok && val != "" {
		request.AliCloudCertificateName = val.(string)
	}

	// check server_certificate and private_key
	if request.AliCloudCertificateId == "" {
		if val := strings.Trim(request.ServerCertificate, " "); val == "" {
			if file_name, ok := d.GetOk("server_certificate_file"); ok && file_name != "" {
				serverCertificateContent, err := readFileContent(file_name.(string))
				if err != nil {
					return fmt.Errorf("UploadServerCertificate got an error %#v, when read server certificate file %s content.", err, file_name)
				}
				request.ServerCertificate = serverCertificateContent
			} else {
				return fmt.Errorf("UploadServerCertificate got an error, as either server_certificate or server_certificate_file should be not null when alicloud_certificate_id is null.")
			}
		}

		if val := strings.Trim(request.PrivateKey, " "); val == "" {
			if file_name, ok := d.GetOk("private_key_file"); ok && file_name != "" {
				privateKeyContent, err := readFileContent(file_name.(string))
				if err != nil {
					return fmt.Errorf("UploadServerCertificate got an error %#v, when read private key file %s content.", err, file_name)
				}
				request.PrivateKey = privateKeyContent
			} else {
				return fmt.Errorf("UploadServerCertificate got an error, as either private_key or private_file  should be not null when alicloud_certificate_id is null.")
			}
		}
	}

	response, err := client.UploadServerCertificate(request)
	if err != nil {
		return fmt.Errorf("UploadServerCertificate got an error: %#v", err)
	}

	d.SetId(response.ServerCertificateId)

	return resourceAlicloudSlbServerCertificateUpdate(d, meta)
}

func resourceAlicloudSlbServerCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	serverCertificate, err := client.describeSlbServerCertificate(d.Id())
	if err != nil {
		if IsExceptedError(err, SlbServerCertificateIdNotFound) || NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	if err := d.Set("name", serverCertificate.ServerCertificateName); err != nil {
		return err
	}

	if serverCertificate.AliCloudCertificateId != "" {
		if err := d.Set("alicloud_certificate_id", serverCertificate.AliCloudCertificateId); err != nil {
			return err
		}
	}

	if serverCertificate.AliCloudCertificateName != "" {
		if err := d.Set("alicloud_certificate_name", serverCertificate.AliCloudCertificateName); err != nil {
			return err
		}
	}

	return nil
}

func resourceAlicloudSlbServerCertificateUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient).slbconn

	if !d.IsNewResource() && d.HasChange("name") {
		request := slb.CreateSetServerCertificateNameRequest()
		request.ServerCertificateId = d.Id()
		request.ServerCertificateName = d.Get("name").(string)
		if _, err := client.SetServerCertificateName(request); err != nil {
			return fmt.Errorf("SetServerCertificateName set %s  name %s got an error: %#v",
				d.Id(), request.ServerCertificateName, err)

		}
	}

	return resourceAlicloudSlbServerCertificateRead(d, meta)
}

func resourceAlicloudSlbServerCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := slb.CreateDeleteServerCertificateRequest()
		request.ServerCertificateId = d.Id()
		if _, err := client.slbconn.DeleteServerCertificate(request); err != nil {
			if IsExceptedError(err, SlbServerCertificateIdNotFound) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("DeleteServerCertificate %s got an error: %#v.", d.Id(), err))
		}

		if _, err := client.describeSlbServerCertificate(d.Id()); err != nil {
			if IsExceptedError(err, SlbServerCertificateIdNotFound) || NotFoundError(err) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("While DeleteServerCertificate，DescribeServerCertificates %s got an error: %#v.", d.Id(), err))
		}

		return resource.RetryableError(fmt.Errorf("DeleteServerCertificate %s timeout.", d.Id()))
	})
}
