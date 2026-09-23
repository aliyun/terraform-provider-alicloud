package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudSslCertificatesServiceCertificateApply() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSslCertificatesServiceCertificateApplyCreate,
		Read:   resourceAliCloudSslCertificatesServiceCertificateApplyRead,
		Delete: resourceAliCloudSslCertificatesServiceCertificateApplyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			// Deleting withdraws the application, and the withdrawal is processed asynchronously —
			// the instance can sit in pending for several minutes afterwards. It has to leave that
			// state before the instance itself can be refunded and removed, so the delete is given
			// room to wait it out.
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// The five fields below are the application configuration. ApplyCertificate accepts
			// none of them — they are written onto the certificate instance by UpdateInstance, and
			// ApplyCertificate submits whatever the instance holds at that moment. They are
			// declared here so that changing any of them replaces this resource, which resubmits
			// the application: a different configuration means a different certificate is being
			// requested. Reference the corresponding attribute of the instance rather than
			// hardcoding a value; Create rejects a mismatch.
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"csr": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"validation_method": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"key_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"generate_csr_method": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"domain_validation_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"root_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"validation_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"validation_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"validation_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cname_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"certificate_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pending_result": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cert_identifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// applicationConfigurationFields maps the application configuration declared on this resource to
// the field GetInstanceDetail reports it under.
var applicationConfigurationFields = []struct {
	schemaKey string
	apiKey    string
}{
	{"domain", "Domain"},
	{"csr", "Csr"},
	{"validation_method", "ValidationMethod"},
	{"key_algorithm", "KeyAlgorithm"},
	{"generate_csr_method", "GenerateCsrMethod"},
}

func resourceAliCloudSslCertificatesServiceCertificateApplyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sslCertificatesServiceServiceV2 := SslCertificatesServiceServiceV2{client}

	instanceId := d.Get("instance_id").(string)

	// ApplyCertificate submits the configuration held on the instance and accepts no application
	// parameters of its own, so a configured value that disagrees with the instance would be
	// silently ignored and then read back as the instance's value — leaving the resource in a diff
	// it can never converge out of. Refuse that up front instead.
	instanceObject, err := sslCertificatesServiceServiceV2.DescribeSslCertificatesServiceInstance(instanceId)
	if err != nil {
		return WrapError(err)
	}
	for _, field := range applicationConfigurationFields {
		configured, ok := d.GetOk(field.schemaKey)
		if !ok {
			continue
		}
		onInstance := ""
		if v, exist := instanceObject[field.apiKey]; exist && v != nil {
			onInstance = fmt.Sprint(v)
		}
		if onInstance != fmt.Sprint(configured) {
			return WrapError(Error("%s is set to %q, but certificate instance %s currently holds %q. "+
				"ApplyCertificate submits the configuration recorded on the instance, so the two must agree. "+
				"Reference the corresponding attribute of alicloud_ssl_certificates_service_instance instead of hardcoding a value.",
				field.schemaKey, fmt.Sprint(configured), instanceId, onInstance))
		}
	}

	action := "ApplyCertificate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	request = make(map[string]interface{})
	request["InstanceId"] = instanceId
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("cas", "2020-04-07", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ssl_certificates_service_certificate_apply", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(instanceId)

	// ApplyCertificate answers with nothing but a RequestId, and the instance takes a few seconds to
	// reflect the submission. Wait for the domain validation records to appear before returning —
	// they are the whole point of this resource, and returning early would hand callers an empty
	// list to drive their DNS records from.
	//
	// The wait keys on DomainValidationList rather than CertificateStatus: the latter only gets a
	// value once the certificate authority has actually issued a certificate, which cannot happen
	// until those very records have been published and validated.
	stateConf := BuildStateConf([]string{}, []string{"#CHECKSET"}, d.Timeout(schema.TimeoutCreate), 5*time.Second,
		sslCertificatesServiceServiceV2.SslCertificatesServiceCertificateApplyStateRefreshFunc(d.Id(), "#$.DomainValidationList", []string{}))
	// An application is not operable the moment its validation records appear: for several seconds
	// afterwards the instance still refuses to withdraw it, answering OperationDenied.StatusNotSupport.
	// Requiring the records to hold across consecutive reads keeps this create from handing back a
	// resource that cannot yet be destroyed, which is where that refusal would otherwise surface.
	stateConf.ContinuousTargetOccurence = 4
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudSslCertificatesServiceCertificateApplyRead(d, meta)
}

func resourceAliCloudSslCertificatesServiceCertificateApplyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sslCertificatesServiceServiceV2 := SslCertificatesServiceServiceV2{client}

	objectRaw, err := sslCertificatesServiceServiceV2.DescribeSslCertificatesServiceCertificateApply(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ssl_certificates_service_certificate_apply DescribeSslCertificatesServiceCertificateApply Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", d.Id())
	d.Set("domain", objectRaw["Domain"])
	d.Set("csr", objectRaw["Csr"])
	d.Set("validation_method", objectRaw["ValidationMethod"])
	d.Set("key_algorithm", objectRaw["KeyAlgorithm"])
	d.Set("generate_csr_method", objectRaw["GenerateCsrMethod"])
	d.Set("certificate_status", objectRaw["CertificateStatus"])
	d.Set("pending_result", objectRaw["PendingResult"])
	d.Set("cert_identifier", objectRaw["CertIdentifier"])
	// Only populated once the certificate authority has issued the certificate.
	if v, ok := objectRaw["CertificateId"]; ok && v != nil {
		d.Set("certificate_id", formatInt(v))
	}

	domainValidationList := make([]map[string]interface{}, 0)
	if v, ok := objectRaw["DomainValidationList"]; ok {
		if items, ok := v.([]interface{}); ok {
			for _, item := range items {
				itemMap, ok := item.(map[string]interface{})
				if !ok {
					continue
				}
				domainValidationList = append(domainValidationList, map[string]interface{}{
					"domain":           itemMap["Domain"],
					"root_domain":      itemMap["RootDomain"],
					"validation_type":  itemMap["ValidationType"],
					"validation_key":   itemMap["ValidationKey"],
					"validation_value": itemMap["ValidationValue"],
					"cname":            itemMap["Cname"],
					"cname_key":        itemMap["CnameKey"],
				})
			}
		}
	}
	// The validation records are a one-time proof of domain control: once the certificate has been
	// issued the application stops reporting them. Keeping the last known set rather than blanking
	// the attribute preserves what was published — callers drive their DNS records off it, and
	// clearing it would tear those records down the moment the certificate arrives.
	if len(domainValidationList) > 0 {
		if err := d.Set("domain_validation_list", domainValidationList); err != nil {
			return WrapError(err)
		}
	}

	return nil
}

func resourceAliCloudSslCertificatesServiceCertificateApplyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sslCertificatesServiceServiceV2 := SslCertificatesServiceServiceV2{client}
	action := "CancelPendingCertificate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("cas", "2020-04-07", action, query, request, true)
		if err == nil {
			return nil
		}
		if NeedRetry(err) {
			wait()
			return resource.RetryableError(err)
		}
		return resource.NonRetryableError(err)
	})
	addDebug(action, response, request)

	if err != nil {
		// The instance is gone, so there is nothing left to withdraw from.
		if NotFoundError(err) {
			log.Printf("[WARN] %s CancelPendingCertificate tolerated failure: %v", d.Id(), err)
			return nil
		}
		// StatusNotSupport says the instance will not withdraw anything right now, which is the
		// state this delete is asking for only when there is nothing left to withdraw — the
		// certificate has been issued, or the application was withdrawn elsewhere. The same answer
		// comes back when an application really is outstanding and refuses to go, and reporting
		// that as success would leave the instance stuck in pending for whatever deletes it next.
		// The error code cannot tell the two apart, so the instance is asked instead.
		//
		// Revoking an issued certificate is deliberately not done here: the certificate is a
		// separate object from the application that produced it, other services may be serving it,
		// and revocation is irreversible. An instance left holding one cannot be refunded, which
		// its own delete reports.
		if IsExpectedErrors(err, []string{"OperationDenied.StatusNotSupport"}) &&
			!sslCertificatesServiceServiceV2.certificateApplicationIsOutstanding(d.Id()) {
			log.Printf("[WARN] %s CancelPendingCertificate tolerated failure: %v", d.Id(), err)
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	// Withdrawing is asynchronous: the call returns while the instance is still pending, and it
	// leaves that state only once the withdrawal has actually been processed. Wait for it here —
	// an instance still in pending cannot be refunded, so deleting it right after would fail with
	// CouldNotRefund.NotSupport.
	//
	// Any state other than pending means the withdrawal has landed, so all of them are accepted
	// rather than guessing at one: a withdrawn instance has been observed settling on normal, but
	// inactive, an expiring or expired subscription, a refunded or a disabled instance are all
	// equally "no longer applying for a certificate".
	stateConf := BuildStateConf([]string{"pending"},
		[]string{"normal", "inactive", "willExpire", "expired", "refund", "closed"},
		d.Timeout(schema.TimeoutDelete), 10*time.Second,
		sslCertificatesServiceServiceV2.SslCertificatesServiceCertificateApplyStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
