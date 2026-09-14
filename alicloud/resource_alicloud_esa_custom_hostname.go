package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudEsaCustomHostname() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaCustomHostnameCreate,
		Read:   resourceAliCloudEsaCustomHostnameRead,
		Update: resourceAliCloudEsaCustomHostnameUpdate,
		Delete: resourceAliCloudEsaCustomHostnameDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cas_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"cas_region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cert_apply_code": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cert_apply_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cert_http_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cert_http_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cert_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cert_not_after": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cert_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cert_txt_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cert_txt_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cert_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"free", "upload", "cas"}, false),
			},
			"certificate": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"conflict_with": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"hostname_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name_match_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"offline_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"record_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"record_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"site_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"site_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssl_flag": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"on", "off"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"verify_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"verify_host": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudEsaCustomHostnameCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateCustomHostname"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["Hostname"] = d.Get("hostname")
	request["SiteId"] = d.Get("site_id")
	request["RecordId"] = d.Get("record_id")
	request["SslFlag"] = d.Get("ssl_flag")
	request["CertType"] = d.Get("cert_type")
	if v, ok := d.GetOk("certificate"); ok {
		request["Certificate"] = v
	}
	if v, ok := d.GetOk("private_key"); ok {
		request["PrivateKey"] = v
	}
	if v, ok := d.GetOk("cas_id"); ok {
		request["CasId"] = v
	}
	if v, ok := d.GetOk("cas_region"); ok {
		request["CasRegion"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Site.ServiceBusy", "Record.ServiceBusy", "TooManyRequests"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_custom_hostname", action, AlibabaCloudSdkGoERROR)
	}

	hostnameId, err := jsonpath.Get("$.Data.Content.HostnameId", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, "alicloud_esa_custom_hostname", "$.Data.Content.HostnameId", response)
	}

	d.SetId(fmt.Sprint(hostnameId))

	return resourceAliCloudEsaCustomHostnameRead(d, meta)
}

func resourceAliCloudEsaCustomHostnameRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	object, err := esaServiceV2.DescribeEsaCustomHostname(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_custom_hostname DescribeEsaCustomHostname Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("hostname_id", fmt.Sprint(object["HostnameId"]))
	d.Set("hostname", object["Hostname"])
	d.Set("site_id", object["SiteId"])
	d.Set("site_name", object["SiteName"])
	d.Set("record_id", object["RecordId"])
	d.Set("record_name", object["RecordName"])
	d.Set("status", object["Status"])
	d.Set("ssl_flag", object["SslFlag"])
	d.Set("cert_type", object["CertType"])
	d.Set("cas_id", object["CasId"])
	d.Set("cas_region", object["CasRegion"])
	d.Set("cert_id", object["CertId"])
	d.Set("certificate", object["Certificate"])
	d.Set("cert_status", object["CertStatus"])
	d.Set("cert_txt_key", object["CertTxtKey"])
	d.Set("cert_txt_value", object["CertTxtValue"])
	d.Set("cert_http_key", object["CertHttpKey"])
	d.Set("cert_http_value", object["CertHttpValue"])
	d.Set("cert_not_after", object["CertNotAfter"])
	d.Set("cert_apply_message", object["CertApplyMessage"])
	d.Set("cert_apply_code", object["CertApplyCode"])
	d.Set("verify_host", object["VerifyHost"])
	d.Set("verify_code", object["VerifyCode"])
	d.Set("offline_reason", object["OfflineReason"])
	d.Set("conflict_with", object["ConflictWith"])
	d.Set("create_time", object["GmtCreate"])
	d.Set("update_time", object["GmtModified"])
	d.Set("region_id", object["RegionId"])
	d.Set("name_match_type", object["NameMatchType"])

	return nil
}

func resourceAliCloudEsaCustomHostnameUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	update := false
	action := "UpdateCustomHostname"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	request = make(map[string]interface{})
	request["HostnameId"] = d.Id()

	if d.HasChange("record_id") {
		update = true
		request["RecordId"] = d.Get("record_id")
	}
	if d.HasChange("ssl_flag") {
		update = true
		request["SslFlag"] = d.Get("ssl_flag")
	}
	if d.HasChange("cert_type") {
		update = true
		request["CertType"] = d.Get("cert_type")
	}
	if d.HasChange("certificate") {
		update = true
		if v, ok := d.GetOk("certificate"); ok {
			request["Certificate"] = v
		}
	}
	if d.HasChange("private_key") {
		update = true
		if v, ok := d.GetOk("private_key"); ok {
			request["PrivateKey"] = v
		}
	}
	if d.HasChange("cas_id") {
		update = true
		if v, ok := d.GetOk("cas_id"); ok {
			request["CasId"] = v
		}
	}
	if d.HasChange("cas_region") {
		update = true
		if v, ok := d.GetOk("cas_region"); ok {
			request["CasRegion"] = v
		}
	}

	if update {
		var err error
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Site.ServiceBusy", "Record.ServiceBusy", "TooManyRequests"}) || NeedRetry(err) {
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

	return resourceAliCloudEsaCustomHostnameRead(d, meta)
}

func resourceAliCloudEsaCustomHostnameDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DeleteCustomHostname"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["HostnameId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Site.ServiceBusy", "Record.ServiceBusy", "TooManyRequests"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if NotFoundError(err) || IsExpectedErrors(err, []string{"CustomHostname.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
