// This file is maintained by hand. The certificate validation resource is a wait with no backing
// API of its own — something the generator cannot express — so unlike its sibling files it is not
// regenerated from a specification.
package alicloud

import (
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudSslCertificatesServiceCertificateValidation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSslCertificatesServiceCertificateValidationCreate,
		Read:   resourceAliCloudSslCertificatesServiceCertificateValidationRead,
		Delete: resourceAliCloudSslCertificatesServiceCertificateValidationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			// Domain validation and issuance are performed by the certificate authority. A DV
			// certificate is usually issued within minutes of the validation records going live,
			// while OV and EV certificates additionally go through a manual review that can take
			// hours. Raise this timeout accordingly for those.
			Create: schema.DefaultTimeout(75 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// The IDs of whatever published the domain ownership validation records — typically
			// alicloud_alidns_record resources created from the domain_validation_list reported by
			// alicloud_ssl_certificates_service_certificate_apply. The values map to no API field
			// and are never read; the attribute exists so that referencing the record IDs here
			// makes Terraform publish the records before this wait starts, the same way an explicit
			// depends_on would.
			"validation_record_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"certificate_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cert_identifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudSslCertificatesServiceCertificateValidationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sslCertificatesServiceServiceV2 := SslCertificatesServiceServiceV2{client}

	// This resource creates nothing. The application was already submitted by
	// alicloud_ssl_certificates_service_certificate_apply, and the identity of the wait is the
	// instance whose application is being waited on.
	d.SetId(d.Get("instance_id").(string))

	// Pending states are left empty on purpose: the intermediate states a certificate passes
	// through while under review are not documented, and no API reports the progress of the manual
	// review that OV and EV certificates go through. An empty pending list means "keep waiting
	// unless the status is a target state or a failure state".
	stateConf := BuildStateConf([]string{}, []string{"issued"}, d.Timeout(schema.TimeoutCreate), 30*time.Second,
		sslCertificatesServiceServiceV2.SslCertificatesServiceCertificateValidationStateRefreshFunc(d.Id(), "CertificateStatus", []string{"revoked", "expired"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudSslCertificatesServiceCertificateValidationRead(d, meta)
}

func resourceAliCloudSslCertificatesServiceCertificateValidationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sslCertificatesServiceServiceV2 := SslCertificatesServiceServiceV2{client}

	objectRaw, err := sslCertificatesServiceServiceV2.DescribeSslCertificatesServiceCertificateValidation(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ssl_certificates_service_certificate_validation DescribeSslCertificatesServiceCertificateValidation Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", d.Id())
	// CertificateId is typed inconsistently across this product: an integer here and on
	// GetCertificateDetail, a string on ListCertificates. formatInt normalises whichever form
	// comes back. It is only set once issuance has produced a certificate.
	if v, ok := objectRaw["CertificateId"]; ok && v != nil {
		d.Set("certificate_id", formatInt(v))
	}
	d.Set("cert_identifier", objectRaw["CertIdentifier"])
	d.Set("certificate_status", objectRaw["CertificateStatus"])

	return nil
}

func resourceAliCloudSslCertificatesServiceCertificateValidationDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Certificate Validation. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
