// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudSslCertificatesServiceInstanceCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSslCertificatesServiceInstanceCertificateCreate,
		Read:   resourceAliCloudSslCertificatesServiceInstanceCertificateRead,
		Delete: resourceAliCloudSslCertificatesServiceInstanceCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"certificate_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cert_identifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_source": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"common_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subject_alternative_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"not_before": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"not_after": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"issuer": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"serial": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"finger_print": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"exist_private_key": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"using_product_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAliCloudSslCertificatesServiceInstanceCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	// Certificates are issued by the certificate authority once an application passes review and
	// cannot be created through an API, so this resource has no create operation. Declaring it
	// brings an already-issued certificate under management so that its attributes are tracked in
	// state; use alicloud_ssl_certificates_service_certificate_apply to request a certificate.
	d.SetId(fmt.Sprint(d.Get("certificate_id")))

	return resourceAliCloudSslCertificatesServiceInstanceCertificateRead(d, meta)
}

func resourceAliCloudSslCertificatesServiceInstanceCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sslCertificatesServiceServiceV2 := SslCertificatesServiceServiceV2{client}

	objectRaw, err := sslCertificatesServiceServiceV2.DescribeSslCertificatesServiceInstanceCertificate(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ssl_certificates_service_instance_certificate DescribeSslCertificatesServiceInstanceCertificate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	// CertificateId is typed inconsistently across this product: an integer here and on
	// GetInstanceDetail, a string on ListCertificates. formatInt normalises whichever form arrives.
	if v, ok := objectRaw["CertificateId"]; ok && v != nil {
		d.Set("certificate_id", formatInt(v))
	}
	d.Set("instance_id", objectRaw["InstanceId"])
	d.Set("cert_identifier", objectRaw["CertIdentifier"])
	d.Set("certificate_name", objectRaw["CertificateName"])
	d.Set("certificate_status", objectRaw["CertificateStatus"])
	d.Set("certificate_source", objectRaw["CertificateSource"])
	d.Set("common_name", objectRaw["CommonName"])
	d.Set("domain", objectRaw["Domain"])
	d.Set("algorithm", objectRaw["Algorithm"])
	d.Set("key_size", objectRaw["KeySize"])
	d.Set("not_before", objectRaw["NotBefore"])
	d.Set("not_after", objectRaw["NotAfter"])
	d.Set("issuer", objectRaw["Issuer"])
	d.Set("serial", objectRaw["Serial"])
	d.Set("finger_print", objectRaw["FingerPrint"])
	d.Set("exist_private_key", objectRaw["ExistPrivateKey"])
	d.Set("subject_alternative_names", objectRaw["SubjectAlternativeNames"])
	d.Set("using_product_list", objectRaw["UsingProductList"])

	return nil
}

func resourceAliCloudSslCertificatesServiceInstanceCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Instance Certificate. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
