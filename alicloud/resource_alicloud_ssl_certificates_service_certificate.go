package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudSslCertificatesServiceCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSslCertificatesServiceCertificateCreate,
		Read:   resourceAliCloudSslCertificatesServiceCertificateRead,
		Update: resourceAliCloudSslCertificatesServiceCertificateUpdate,
		Delete: resourceAliCloudSslCertificatesServiceCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cert": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"certificate_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				AtLeastOneOf: []string{"certificate_name", "name"},
			},
			"name": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Computed:   true,
				Deprecated: "attribute 'name' has been deprecated from provider version 1.129.0 and it will be remove in the future version. Please use the new attribute 'certificate_name' instead.",
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudSslCertificatesServiceCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	action := "CreateUserCertificate"
	request := make(map[string]interface{})
	request["Cert"] = d.Get("cert")
	request["Key"] = d.Get("key")

	if v, ok := d.GetOk("certificate_name"); ok {
		request["Name"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	} else {
		return WrapError(fmt.Errorf("Field 'certificate_name' is required"))
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("cas", "2018-07-13", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ssl_certificates_service_certificate", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["CertId"]))

	return resourceAliCloudSslCertificatesServiceCertificateRead(d, meta)
}

func resourceAliCloudSslCertificatesServiceCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	casService := CasService{client}

	object, err := casService.DescribeSslCertificatesServiceCertificate(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ssl_certificates_service_certificate casService.DescribeSslCertificatesServiceCertificate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cert", object["Cert"])
	d.Set("key", object["Key"])
	d.Set("certificate_name", object["Name"])
	d.Set("name", object["Name"])

	return nil
}

func resourceAliCloudSslCertificatesServiceCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAliCloudSslCertificatesServiceCertificateRead(d, meta)
}

func resourceAliCloudSslCertificatesServiceCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	casService := CasService{client}
	action := "DeleteUserCertificate"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"CertId": d.Id(),
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("cas", "2018-07-13", action, nil, request, true)
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

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, casService.SslCertificatesServiceCertificateStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
