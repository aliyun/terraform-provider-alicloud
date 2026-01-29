package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudSslCertificatesServicePcaCert() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSslCertificatesServicePcaCertCreate,
		Read:   resourceAliCloudSslCertificatesServicePcaCertRead,
		Update: resourceAliCloudSslCertificatesServicePcaCertUpdate,
		Delete: resourceAliCloudSslCertificatesServicePcaCertDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"after_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"alias_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"before_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"common_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"country_code": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"custom_identifier": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"days": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enable_crl": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"immediately": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"locality": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"months": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"organization": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"organization_unit": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"parent_identifier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"san_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"san_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"ISSUE", "PENDING", "REVOKE"}, false),
			},
			"tags": tagsSchema(),
			"upload_flag": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"years": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudSslCertificatesServicePcaCertCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateClientCertificate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("state"); ok {
		request["State"] = v
	}
	if v, ok := d.GetOk("country_code"); ok {
		request["Country"] = v
	}
	if v, ok := d.GetOk("organization_unit"); ok {
		request["OrganizationUnit"] = v
	}
	if v, ok := d.GetOk("locality"); ok {
		request["Locality"] = v
	}
	if v, ok := d.GetOkExists("after_time"); ok {
		request["AfterTime"] = v
	}
	request["ParentIdentifier"] = d.Get("parent_identifier")
	if v, ok := d.GetOkExists("enable_crl"); ok {
		request["EnableCrl"] = v
	}
	if v, ok := d.GetOk("organization"); ok {
		request["Organization"] = v
	}
	if v, ok := d.GetOkExists("days"); ok {
		request["Days"] = v
	}
	if v, ok := d.GetOkExists("years"); ok {
		request["Years"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMapWithTags(request, tagsMap)
	}

	if v, ok := d.GetOkExists("before_time"); ok {
		request["BeforeTime"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOkExists("immediately"); ok {
		request["Immediately"] = v
	}
	if v, ok := d.GetOkExists("months"); ok {
		request["Months"] = v
	}
	if v, ok := d.GetOk("custom_identifier"); ok {
		request["CustomIdentifier"] = v
	}
	if v, ok := d.GetOk("san_type"); ok {
		request["SanType"] = v
	}
	if v, ok := d.GetOk("san_value"); ok {
		request["SanValue"] = v
	}
	if v, ok := d.GetOk("alias_name"); ok {
		request["AliasName"] = v
	}
	if v, ok := d.GetOk("algorithm"); ok {
		request["Algorithm"] = v
	}
	if v, ok := d.GetOk("common_name"); ok {
		request["CommonName"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("cas", "2020-06-30", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ssl_certificates_service_pca_cert", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Identifier"]))

	return resourceAliCloudSslCertificatesServicePcaCertUpdate(d, meta)
}

func resourceAliCloudSslCertificatesServicePcaCertRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sslCertificatesServiceServiceV2 := SslCertificatesServiceServiceV2{client}

	objectRaw, err := sslCertificatesServiceServiceV2.DescribeSslCertificatesServicePcaCert(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ssl_certificates_service_pca_cert DescribeSslCertificatesServicePcaCert Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("algorithm", objectRaw["FullAlgorithm"])
	d.Set("alias_name", objectRaw["AliasName"])
	d.Set("common_name", objectRaw["CommonName"])
	d.Set("country_code", objectRaw["CountryCode"])
	d.Set("custom_identifier", objectRaw["CustomIdentifier"])
	d.Set("days", objectRaw["Days"])
	d.Set("locality", objectRaw["Locality"])
	d.Set("organization", objectRaw["Organization"])
	d.Set("organization_unit", objectRaw["OrganizationUnit"])
	d.Set("parent_identifier", objectRaw["ParentIdentifier"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("state", objectRaw["State"])
	d.Set("status", objectRaw["Status"])
	d.Set("upload_flag", objectRaw["UploadFlag"])

	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudSslCertificatesServicePcaCertUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	sslCertificatesServiceServiceV2 := SslCertificatesServiceServiceV2{client}
	objectRaw, _ := sslCertificatesServiceServiceV2.DescribeSslCertificatesServicePcaCert(d.Id())

	if d.HasChange("upload_flag") {
		var err error
		target := d.Get("upload_flag").(int)

		currentStatus, err := jsonpath.Get("UploadFlag", objectRaw)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "UploadFlag", objectRaw)
		}
		if formatInt(currentStatus) != target {
			if target == 1 {
				action := "UploadPcaCertToCas"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["Ids"] = d.Id()

				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("cas", "2020-06-30", action, query, request, true)
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
		}
	}

	if d.HasChange("status") {
		var err error
		target := d.Get("status").(string)

		currentStatus, err := jsonpath.Get("Status", objectRaw)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "Status", objectRaw)
		}
		if fmt.Sprint(currentStatus) != target {
			if target == "REVOKE" {
				action := "CreateRevokeClientCertificate"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["Identifier"] = d.Id()

				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("cas", "2020-06-30", action, query, request, true)
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
		}
	}

	var err error
	action := "MoveResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["ResourceGroupId"] = d.Get("resource_group_id")
	request["ResourceType"] = "PcaCertificate"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("cas", "2020-06-30", action, query, request, true)
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
	update = false
	action = "UpdatePcaCertificate"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Identifier"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("alias_name") {
		update = true
		request["AliasName"] = d.Get("alias_name")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("cas", "2020-06-30", action, query, request, true)
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

	if !d.IsNewResource() && d.HasChange("tags") {
		sslCertificatesServiceServiceV2 := SslCertificatesServiceServiceV2{client}
		if err := sslCertificatesServiceServiceV2.SetResourceTags(d, "PcaCertificate"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudSslCertificatesServicePcaCertRead(d, meta)
}

func resourceAliCloudSslCertificatesServicePcaCertDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteClientCertificate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["Identifier"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("cas", "2020-06-30", action, query, request, true)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
