package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudSslCertificatesServiceInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSslCertificatesServiceInstanceCreate,
		Read:   resourceAliCloudSslCertificatesServiceInstanceRead,
		Update: resourceAliCloudSslCertificatesServiceInstanceUpdate,
		Delete: resourceAliCloudSslCertificatesServiceInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			// The refund issued before deletion is processed asynchronously, and the instance
			// refuses deletion until it lands, so the delete needs room to keep retrying.
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_reissue": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"enable", "disable"}, false),
			},
			"average_waiting_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"brand": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Read-only: the address on an instance is taken from the company it references, and
			// from a server-side default when it references none. Values supplied here were being
			// silently replaced, so the attribute reports what the instance actually holds rather
			// than pretending to accept input. Set the address on the company instead.
			"city": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Computed as well as optional: leaving it unset makes the server attach a default
			// company, and that value is reported back.
			"company_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// A set rather than a list: the order in which the server returns the contacts does not
			// follow the order they were submitted in, so an ordered type would report a difference
			// on every refresh even though the membership is identical.
			"contact_id_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"country_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// Computed for the mirror-image reason: with generate_csr_method = online the server
			// generates the CSR itself, so a configuration that leaves this unset must accept the
			// generated value rather than treat it as drift.
			"csr": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// Computed because the server decides this value when the CSR is uploaded rather
			// than generated: with generate_csr_method = upload the domain is taken from the CSR's
			// common name, so a configuration that leaves it unset must accept whatever comes back
			// instead of trying to remove it. With generate_csr_method = online the value set here
			// is used as-is. Supplying a domain that disagrees with an uploaded CSR is rejected by
			// the API.
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"full_domain_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"generate_csr_method": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"online", "upload"}, false),
			},
			"instance_end_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_start_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"order_end_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"order_start_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"parameter": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"code": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"pricing_cycle": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"product_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				// cas_dv_public_cn is kept in the CloudSpec enum only because removing an
				// already-published value is a breaking change; the order API rejects it with
				// COMMODITY.INVALID_COMPONENT, so it is not offered here.
				ValidateFunc: StringInSlice([]string{"cas", "cas_intl"}, false),
			},
			// Read-only for the same reason as city.
			"province": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"spec": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"upgrade_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"validation_method": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"DNS", "HTTP"}, false),
			},
			"wildcard_domain_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudSslCertificatesServiceInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("key_algorithm"); ok {
		request["KeyAlgorithm"] = v
	}
	request["ProductCode"] = "cas"
	if v, ok := d.GetOk("country_code"); ok {
		request["CountryCode"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("product_type"); ok {
		request["ProductType"] = v
	}
	if v, ok := d.GetOk("instance_name"); ok {
		request["CertificateName"] = v
	}
	if v, ok := d.GetOk("generate_csr_method"); ok {
		request["GenerateCsrMethod"] = v
	}
	if v, ok := d.GetOk("csr"); ok {
		request["Csr"] = v
	}
	// parameter is a list of order modules; every entry has to be forwarded, not just the first.
	if v, ok := d.GetOk("parameter"); ok {
		parameterMaps := make([]interface{}, 0)
		for _, item := range v.([]interface{}) {
			itemMap, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			parameterMaps = append(parameterMaps, map[string]interface{}{
				"Code":  itemMap["code"],
				"Value": itemMap["value"],
			})
		}
		request["Parameter"] = parameterMaps
	}

	if v, ok := d.GetOkExists("pricing_cycle"); ok {
		request["PricingCycle"] = v
	}
	if v, ok := d.GetOk("domain"); ok {
		request["Domain"] = v
	}
	request["SubscriptionType"] = "Subscription"
	if v, ok := d.GetOk("auto_reissue"); ok {
		request["AutoReissue"] = v
	}
	if v, ok := d.GetOk("company_id"); ok {
		request["CompanyId"] = v
	}
	if v, ok := d.GetOkExists("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("validation_method"); ok {
		request["ValidationMethod"] = v
	}
	var endpoint string
	request["ProductCode"] = "cas"
	if v, ok := d.GetOk("product_type"); ok {
		request["ProductType"] = v
	} else if client.IsInternationalAccount() {
		request["ProductType"] = "cas_intl"
	} else {
		request["ProductType"] = "cas"
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductType"] = "cas_intl"
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ssl_certificates_service_instance", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Data.InstanceId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudSslCertificatesServiceInstanceUpdate(d, meta)
}

func resourceAliCloudSslCertificatesServiceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sslCertificatesServiceServiceV2 := SslCertificatesServiceServiceV2{client}

	objectRaw, err := sslCertificatesServiceServiceV2.DescribeSslCertificatesServiceInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ssl_certificates_service_instance DescribeSslCertificatesServiceInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if v, ok := objectRaw["AutoReissue"]; ok {
		d.Set("auto_reissue", v)
	}
	d.Set("average_waiting_time", objectRaw["AverageWaitingTime"])
	d.Set("brand", objectRaw["Brand"])
	d.Set("certificate_type", objectRaw["CertificateType"])
	if v, ok := objectRaw["City"]; ok {
		d.Set("city", v)
	}
	if v, ok := objectRaw["CompanyId"]; ok {
		d.Set("company_id", v)
	}
	if v, ok := objectRaw["CountryCode"]; ok {
		d.Set("country_code", v)
	}
	if v, ok := objectRaw["Csr"]; ok {
		d.Set("csr", v)
	}
	if v, ok := objectRaw["Domain"]; ok {
		d.Set("domain", v)
	}
	d.Set("full_domain_count", objectRaw["FullDomainCount"])
	if v, ok := objectRaw["GenerateCsrMethod"]; ok {
		d.Set("generate_csr_method", v)
	}
	d.Set("instance_end_time", objectRaw["InstanceEndTime"])
	if v, ok := objectRaw["CertificateName"]; ok {
		d.Set("instance_name", v)
	}
	d.Set("instance_start_time", objectRaw["InstanceStartTime"])
	d.Set("instance_type", objectRaw["InstanceType"])
	if v, ok := objectRaw["KeyAlgorithm"]; ok {
		d.Set("key_algorithm", v)
	}
	d.Set("order_end_time", objectRaw["OrderEndTime"])
	d.Set("order_start_time", objectRaw["OrderStartTime"])
	if v, ok := objectRaw["Province"]; ok {
		d.Set("province", v)
	}
	if v, ok := objectRaw["ResourceGroupId"]; ok {
		d.Set("resource_group_id", v)
	}
	d.Set("spec", objectRaw["Spec"])
	d.Set("status", objectRaw["Status"])
	d.Set("upgrade_status", objectRaw["UpgradeStatus"])
	if v, ok := objectRaw["ValidationMethod"]; ok {
		d.Set("validation_method", v)
	}
	d.Set("wildcard_domain_count", objectRaw["WildcardDomainCount"])

	// GetInstanceDetail does not echo back every attribute UpdateInstance accepts, and which ones
	// it returns varies with the state of the instance. Writing an absent key would replace the
	// configured value in state with an empty one and produce a permanent diff, so the configurable
	// attributes above are refreshed only when the API actually reports them.
	if raw, ok := objectRaw["ContactIdList"]; ok && raw != nil {
		d.Set("contact_id_list", convertToInterfaceArray(raw))
	}
	if tagsMaps, ok := objectRaw["Tags"]; ok && tagsMaps != nil {
		d.Set("tags", tagsToMap(tagsMaps))
	}

	return nil
}

func resourceAliCloudSslCertificatesServiceInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateInstance"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("key_algorithm") || d.IsNewResource() {
		if v, ok := d.GetOk("key_algorithm"); ok {
			update = true
			request["KeyAlgorithm"] = v
		}
	}

	if d.HasChange("country_code") || d.IsNewResource() {
		if v, ok := d.GetOk("country_code"); ok {
			update = true
			request["CountryCode"] = v
		}
	}

	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["ResourceGroupId"] = d.Get("resource_group_id")
	}

	if d.HasChange("generate_csr_method") || d.IsNewResource() {
		if v, ok := d.GetOk("generate_csr_method"); ok {
			update = true
			request["GenerateCsrMethod"] = v
		}
	}

	if d.HasChange("instance_name") || d.IsNewResource() {
		if v, ok := d.GetOk("instance_name"); ok {
			update = true
			request["CertificateName"] = v
		}
	}

	if d.HasChange("contact_id_list") {
		update = true
		if v, ok := d.GetOk("contact_id_list"); ok || d.HasChange("contact_id_list") {
			contactIdListMapsArray := convertToInterfaceArray(v)

			request["ContactIdList"] = contactIdListMapsArray
		}
	}

	// Tags is a list of {TagKey, TagValue} pairs, one entry per tag.
	if d.HasChange("tags") {
		update = true
		tagsMaps := make([]interface{}, 0)
		if v, ok := d.GetOk("tags"); ok {
			for key, value := range v.(map[string]interface{}) {
				tagsMaps = append(tagsMaps, map[string]interface{}{
					"TagKey":   key,
					"TagValue": value,
				})
			}
		}
		request["Tags"] = tagsMaps
	}

	if d.HasChange("csr") || d.IsNewResource() {
		if v, ok := d.GetOk("csr"); ok {
			update = true
			request["Csr"] = v
		}
	}

	if d.HasChange("domain") || d.IsNewResource() {
		if v, ok := d.GetOk("domain"); ok {
			update = true
			request["Domain"] = v
		}
	}

	if d.HasChange("auto_reissue") || d.IsNewResource() {
		if v, ok := d.GetOk("auto_reissue"); ok {
			update = true
			request["AutoReissue"] = v
		}
	}

	if d.HasChange("company_id") || d.IsNewResource() {
		if v, ok := d.GetOk("company_id"); ok {
			update = true
			request["CompanyId"] = v
		}
	}

	if d.HasChange("validation_method") || d.IsNewResource() {
		if v, ok := d.GetOk("validation_method"); ok {
			update = true
			request["ValidationMethod"] = v
		}
	}

	if update {
		wait := incrementalWait(5*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("cas", "2020-04-07", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"InvalidStatus.UpdateProtection"}) || NeedRetry(err) {
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

	return resourceAliCloudSslCertificatesServiceInstanceRead(d, meta)
}

func resourceAliCloudSslCertificatesServiceInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	sslCertificatesServiceServiceV2 := SslCertificatesServiceServiceV2{client}
	action := "DeleteInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	// An instance with a certificate application in flight sits in the pending state, and in that
	// state it can neither be refunded nor deleted. Wait for it to settle rather than acting on the
	// application: withdrawing one belongs to the application resource, which is destroyed ahead of
	// this one and whose withdrawal is asynchronous, so what is left here is to wait for it to land.
	//
	// An instance still pending with no application resource to withdraw it — one applied for
	// outside Terraform, say — times out here rather than being withdrawn silently. Cancelling an
	// order the configuration never mentioned is not this resource's call to make.
	if object, describeErr := sslCertificatesServiceServiceV2.DescribeSslCertificatesServiceInstance(d.Id()); describeErr == nil {
		if fmt.Sprint(object["Status"]) == "pending" {
			stateConf := BuildStateConf([]string{"pending"},
				[]string{"normal", "inactive", "willExpire", "expired", "refund", "closed"},
				d.Timeout(schema.TimeoutDelete), 10*time.Second,
				sslCertificatesServiceServiceV2.SslCertificatesServiceInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
			if _, waitErr := stateConf.WaitForState(); waitErr != nil {
				return WrapErrorf(waitErr, IdMsg, d.Id())
			}
		}
	}

	// A purchased instance still within its subscription period has to be refunded before it can be
	// deleted. The refund is attempted first and its failure tolerated: an instance whose order is
	// no longer refundable — already expired, or already refunded on an earlier attempt — reports
	// that it cannot be refunded, and such an instance is deletable as it stands.
	//
	// A freshly purchased instance answers with that same CouldNotRefund.NotSupport for the first
	// seconds of its life, while its purchase order is still settling — a refund fired three
	// seconds after payment has been seen rejected where one fired six seconds after succeeded.
	// The code is therefore retried briefly before being taken at its word, or an instance
	// destroyed right after creation is left unrefunded and its deletion stuck behind delete
	// protection. The retry is bounded well below the delete timeout: for an instance that
	// genuinely cannot be refunded, it only defers the tolerant path below.
	refundRequest := map[string]interface{}{
		"InstanceId":  d.Id(),
		"ClientToken": buildClientToken("RefundInstance"),
	}
	refundWait := incrementalWait(3*time.Second, 5*time.Second)
	refundErr := resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, err := client.RpcPost("cas", "2020-04-07", "RefundInstance", make(map[string]interface{}), refundRequest, true)
		if err == nil {
			return nil
		}
		if NeedRetry(err) || IsExpectedErrors(err, []string{"CouldNotRefund.NotSupport"}) {
			refundWait()
			return resource.RetryableError(err)
		}
		return resource.NonRetryableError(err)
	})
	if refundErr != nil {
		if !NotFoundError(refundErr) && !IsExpectedErrors(refundErr, []string{"CouldNotRefund.NotSupport", "OperationDenied.StatusNotSupport"}) {
			return WrapErrorf(refundErr, DefaultErrorMsg, d.Id(), "RefundInstance", AlibabaCloudSdkGoERROR)
		}
		// The error itself is logged, not just the fact that one occurred: a refund that was
		// declined leaves the deletion below to fail on delete protection until the instance
		// becomes deletable some other way, and the reason for the decline is the only clue.
		log.Printf("[WARN] %s RefundInstance tolerated failure: %v", d.Id(), refundErr)
		addDebug("RefundInstance", nil, refundRequest)
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("cas", "2020-04-07", action, query, request, true)
		if err != nil {
			// The refund issued above is processed asynchronously, and until it lands the instance
			// reports that its current state does not permit deletion — the server's own message
			// says to retry later, so that is what happens here.
			if NeedRetry(err) || IsExpectedErrors(err, []string{"InvalidStatus.DeleteProtection"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
